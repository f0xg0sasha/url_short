package main

import (
	"log/slog"
	"os"

	"github.com/f0xg0sasha/url_short/internal/config"
	"github.com/f0xg0sasha/url_short/internal/storage"
)

func main() {
	// Init config
	configs := config.NewConfig()

	// Init logger
	log := setupLogger(configs.Env)
	log.Info("start app")

	// Init database
	_ = storage.NewStorage()

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
