package app

import (
	"context"
	"fmt"
	"gopher/internal/config"
	"gopher/internal/infrastructure/kafka"
	"gopher/internal/infrastructure/kafka/consumer"
	"gopher/internal/infrastructure/minio"
	"gopher/internal/infrastructure/postgres"
	"gopher/internal/infrastructure/postgres/repository"
	"gopher/internal/servers/httpserver"
	"gopher/internal/service"
	httpTrasport "gopher/internal/transport/http"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Marlliton/slogpretty"
	"github.com/fatih/color"
)

func RunApp() {
	config := config.MustLoad("config.yaml")
	printBanner("1.0.0", config.HTTPServer.Address)
	setupLogger(config.Env)

	slog.Info("Config load")

	postgres.Migrate(config)

	pool := postgres.InitDatabese(config)
	trm := postgres.NewTransactionManager(pool)

	repo := repository.NewRepository(pool)
	requestRepo := repository.NewRequestRepository(repo)
	kafkaService := kafka.NewKafka([]string{config.Kafka.Address})
	consumer := consumer.NewOutputConsumer(requestRepo)

	go kafkaService.StartConsume(context.Background(), "output_topic", consumer)
	minio, err := minio.NewMinio(
		config.FileStorage.Address,
		config.FileStorage.Username,
		config.FileStorage.Password,
		config.FileStorage.BucketName,
	)
	if err != nil {
		log.Fatal("Minio dont work")
	}

	requestService := service.NewRequestService(requestRepo, minio, kafkaService, trm)

	handler := httpTrasport.NewHandler(requestService)

	router := httpTrasport.Router(handler)

	srv := httpserver.NewHTTPServer(
		config.HTTPServer.Address,
		router,
		config.HTTPServer.Timeout,
		config.HTTPServer.Timeout,
		config.HTTPServer.IdleTimeout,
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool.Close()

	kafkaService.Close()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	slog.Info("Server exiting")
}

const (
	envLocal = "local"
	envDev   = "dev"
)

func setupLogger(env string) {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	slog.SetDefault(log)
}

func setupPrettySlog() *slog.Logger {
	handler := slogpretty.New(os.Stdout, &slogpretty.Options{
		Level:      slog.LevelDebug,
		AddSource:  true,
		Colorful:   true,
		Multiline:  true,
		TimeFormat: slogpretty.DefaultTimeFormat,
	})

	return slog.New(handler)
}

const art = `
________      .___
\_____  \   __| _/____   _____
  _(__  <  / __ |/  _ \ /     \
 /       \/ /_/ (  <_> )  Y Y  \
/______  /\____ |\____/|__|_|  /
       \/      \/            \/
`

func printBanner(version string, address string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Println(cyan(art))

	fmt.Printf(" %s  %s\n", blue("Version:"), yellow(version))
	fmt.Printf(" %s  %s\n", blue("Address:"), yellow(address))
	fmt.Printf(" %s  %s\n", blue("Status:"), color.GreenString("Running"))
	fmt.Println()
}
