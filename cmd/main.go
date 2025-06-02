package main

import (
	"net/http"
	"os"

	"github.com/f0xg0sasha/url_short/internal/config"
	"github.com/f0xg0sasha/url_short/internal/service"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/f0xg0sasha/url_short/internal/transport/rest"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Init config
	configs := config.NewConfig()

	// Init logger
	log.Info("start app")

	//repositroy
	repository := storage.NewStorage()
	urlService := service.NewURL(repository)

	// Init handlers
	handler := rest.NewHandler(urlService)

	// Run server
	srv := http.Server{
		Addr:        ":8080",
		IdleTimeout: configs.HTTPServer.IdleTimeout,
		ReadTimeout: configs.HTTPServer.Timeout,
		Handler:     handler.InitRouter(),
	}

	log.Info("run server")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %s", err)
	}
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
