package main

import (
	"github.com/2pizzzza/TestTask/internal/config"
	"github.com/2pizzzza/TestTask/internal/lib/logger/sl"
	"github.com/2pizzzza/TestTask/internal/service"
	"github.com/2pizzzza/TestTask/internal/storage/postgres"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	env, err := config.NewConfig()

	if err != nil {
		slog.Error("Failed load env", sl.Err(err))
	}

	logs := setupLogger(env.Env)

	db, err := postgres.New(env)

	if err != nil {
		logs.Error("Failed connect db err: %s", sl.Err(err))
	}

	songService := service.New(*logs, db)
	_ = songService

	log.Printf("Server is live")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupLogger(env string) *slog.Logger {
	var logs *slog.Logger
	switch env {
	case envLocal:
		logs = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logs = slog.New(

			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logs = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logs
}
