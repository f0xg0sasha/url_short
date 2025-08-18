package main

import (
	"net/http"
	"os"

	"github.com/f0xg0sasha/url_short/internal/config"
	"github.com/f0xg0sasha/url_short/internal/service"
	memcache "github.com/f0xg0sasha/url_short/internal/service/cache"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/f0xg0sasha/url_short/internal/transport/rest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	cacheHit = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hit",
		Help: "количество попаданий в кэш",
	})
	cacheMiss = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_miss",
		Help: "количество промахов мимо кэша",
	})
)

func init() {
	prometheus.MustRegister(cacheHit)
	prometheus.MustRegister(cacheMiss)
}

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
	memCache := memcache.NewMemCache(log, repository, cacheHit, cacheMiss)
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

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %s", err)
	}
}
