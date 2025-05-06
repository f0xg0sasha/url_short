package main

import (
	"fmt"
	"log/slog"
	"os"

	"girhub.com/f0xg0sasha/url_short/internal/config"
	"girhub.com/f0xg0sasha/url_short/internal/storage/psql"
)

func main() {
	// Init config
	configs := config.NewConfig()

	// Init logger
	log := setupLogger(configs.Env)
	log.Info("start app")

	// Init database
	db := psql.ConnectionPostgres()
	x, _ := db.Exec("SELECT * FROM URL")
	fmt.Println(x)

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
