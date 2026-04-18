package main

import (
	"gopher/internal/config"
	"gopher/internal/handlers"
	"gopher/internal/logger"
	"gopher/internal/logger/logutils"
	"gopher/internal/servers/httpserver"
	"gopher/internal/storage"
	stdlog "log"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	config, err := config.NewConfigFromFile("./configs/config.yaml")
	if err != nil {
		stdlog.Fatalf("%v", err)
	}

	log := logger.NewLogger(config.Env)

	storage, err := storage.NewStorage(config.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", logutils.Err(err))
	}
	_ = storage
	// err = storage.SaveNewAudioFile(1)
	// fmt.Println(err)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.NewMiddlewareLogger(log))
	router.Use(middleware.Recoverer)

	router.Post("/api/audiofiles", handlers.Post(log, storage))

	server := httpserver.NewHTTPServer(
		config.Address,
		router,
		config.Timeout,
		config.Timeout,
		config.IdleTimeout,
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
