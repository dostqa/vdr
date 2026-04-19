package main

import (
	"context"
	"gopher/internal/clients"
	"gopher/internal/config"
	"gopher/internal/handlers"
	"gopher/internal/logger"
	"gopher/internal/logger/logutils"
	"gopher/internal/servers/httpserver"
	"gopher/internal/service"
	"gopher/internal/storages/database"
	"gopher/internal/storages/filestorage"
	"io"
	stdlog "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.NewConfigFromFile("./configs/config.yaml")
	if err != nil {
		stdlog.Fatalf("%s: %v", "failed to load config", err)
	}

	fileStorage, err := filestorage.NewFileStorage(
		cfg.FileStorage.Address,
		cfg.FileStorage.Username,
		cfg.FileStorage.Password,
		cfg.FileStorage.BucketName,
	)
	if err != nil {
		stdlog.Fatalf("%v", err)
	}

	pool, err := database.InitDataBasePool(cfg)
	if err != nil {
		stdlog.Fatalf("%v", err)
	}

	dataBase := database.NewDataBase(pool)

	kafkaService := clients.NewKafkaService([]string{"broker:9093"})
	defer kafkaService.Close()
	saverService := service.NewSaverService(dataBase, fileStorage, kafkaService, database.NewTransactionManager(pool))

	go kafkaService.StartConsume(context.TODO(), "output_topic", "test", saverService.Consume)

	log := logger.NewLogger(cfg.Env)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"} // Для продакшена лучше указывать конкретные домены
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Кэширование префлайт-запросов в секундах
	}))
	//router.Use(middleware.RequestID)

	// router.Use(logger.NewMiddlewareLogger(log))
	//router.Use(middleware.Recoverer)
	router.Get("/api/audio/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		// Вызываем нашу новую функцию
		fileStream, err := fileStorage.GetFile(r.Context(), name)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer fileStream.Close() // Обязательно закрываем поток!

		w.Header().Set("Content-Type", "audio/webm")

		// Стриммим содержимое напрямую из Minio в ответ пользователю
		io.Copy(w, fileStream)
	})

	router.Get("/api/audiofiles/requests/{id}", handlers.GetByRequestID(log, dataBase, saverService))
	// router.Get("/api/audiofiles/files/{filepath}", handlers.GetByFilePath)

	router.Post("/api/audiofiles", handlers.Post(log, saverService))

	server := httpserver.NewHTTPServer(
		cfg.Address,
		router,
		cfg.Timeout,
		cfg.Timeout,
		cfg.IdleTimeout,
	)

	log.Info("starting server")
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to run server", logutils.Err(err))
	}

	log.Error("server stopped")
}

// 1. Необходимо получить аудиофайл с фронта
// 1.1 Сохранить аудиофайл в storage для последующей обработки и истории
// 1.2 Плюсом наверное сохранять какие-то датаметки для истории запросов
// 2. Разбить аудиофайл на чанки, чтобы отослать на python сервер
// таким образом, чтобы не обрезать слова
// 3. Получить с python сервера:
// Транскрибацию оригинального текста
// Транскрибацию текста с удаленными персональными данными
// Массив персональных данных в аудио {word string, start time.Duration, end time.Duretion}
// 4. Анонимизировать оригинальный файл с помощью полученного массива персональных данных
// 5. Отправить на фронт
// 1) Оригинальный аудиофайл
// 2) Модифицированный аудиофайл
// 3) Массив персональных данных
