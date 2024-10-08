package main

import (
	"fmt"
	_ "github.com/2pizzzza/TestTask/cmd/songLibraries/docs"
	"github.com/2pizzzza/TestTask/internal/config"
	"github.com/2pizzzza/TestTask/internal/http-server/handlers"
	"github.com/2pizzzza/TestTask/internal/http-server/middleware/logger"
	"github.com/2pizzzza/TestTask/internal/lib/logger/sl"
	"github.com/2pizzzza/TestTask/internal/service"
	"github.com/2pizzzza/TestTask/internal/storage/postgres"
	httpSwagger "github.com/swaggo/http-swagger"
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

// @title Swagger Song Libraries API
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /
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
	songHandler := handlers.New(songService)

	mux := http.NewServeMux()
	mux.HandleFunc("/songs/create", songHandler.CreateSongHandler)
	mux.HandleFunc("/songs/update", songHandler.UpdateSongHandler)
	mux.HandleFunc("/songs/info", songHandler.GetSongByIDHandler)
	mux.HandleFunc("/songs/delete", songHandler.DeleteSongHandler)
	mux.HandleFunc("/songs", songHandler.GetAllSongsHandler)
	mux.HandleFunc("/songs/{id}/lyrics", songHandler.GetSongLyricsHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

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
