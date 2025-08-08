package main

import (
	"net/http"
	"os"

	"github.com/f0xg0sasha/url_short/internal/config"
	"github.com/f0xg0sasha/url_short/internal/service"
	memcache "github.com/f0xg0sasha/url_short/internal/service/cache"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/f0xg0sasha/url_short/internal/transport/rest"
	"github.com/sirupsen/logrus"
)

func main() {
	// Init config
	configs := config.NewConfig()

	// Init logger
	log := &logrus.Logger{
		Level:     logrus.DebugLevel,
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
	}

	log.Info("start app")

	//repositroy
	repository := storage.NewStorage()

	// init cache
	memCache := memcache.NewMemCache(log, repository)
	svc := service.NewService(memCache)

	// Init handlers
	handler := rest.NewHandler(log, svc)

	// Run server
	srv := http.Server{
		Addr:        ":8080",
		IdleTimeout: configs.HTTPServer.IdleTimeout,
		ReadTimeout: configs.HTTPServer.Timeout,
		Handler:     handler.InitRouter(),
	}

	log.Info("run server", configs.HTTPServer.Address)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %s", err)
	}
}
