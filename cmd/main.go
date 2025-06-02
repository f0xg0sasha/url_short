package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/f0xg0sasha/url_short/internal/config"
	"github.com/f0xg0sasha/url_short/internal/service"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/f0xg0sasha/url_short/internal/transport/rest"
)

func main() {
	// Init config
	configs := config.NewConfig()

	// Init logger
	log := setupLogger(configs.Env)
	log.Info("start app")

	//repositroy
	repository := storage.NewStorage()
	urlService := service.NewURL(repository)

	// Init handlers
	handler := rest.NewHandler(urlService)
	log.Info("start http server", slog.String("addr", configs.HTTPServer.Address))

	// Run server
	srv := http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Info("run server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil {
		log.Error("server error", slog.Any("error", err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "dev":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
