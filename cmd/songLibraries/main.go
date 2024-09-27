package main

import (
	"fmt"
	"github.com/2pizzzza/TestTask/internal/config"
	"github.com/2pizzzza/TestTask/internal/http-server/middleware/logger"
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

func homeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to the homepage!")
}
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

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	loggedMux := logger.LoggingMiddleware(mux)
	_ = songService

	log.Printf("Server is live. port: %d", env.HttpConn.HttpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.HttpConn.HttpPort), loggedMux))
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
