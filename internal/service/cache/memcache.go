package memcache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/f0xg0sasha/url_short/internal/service"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type URLRepository interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}

type MemCache struct {
	log       *logrus.Logger
	cache     sync.Map
	cacheHit  prometheus.Counter
	cacheMiss prometheus.Counter
	URLRepository
}

func NewMemCache(
	log *logrus.Logger,
	repo URLRepository,
	cacheHit prometheus.Counter,
	cacheMiss prometheus.Counter,
) *MemCache {
	return &MemCache{
		log:           log,
		cache:         sync.Map{},
		URLRepository: repo,
		cacheHit:      cacheHit,
		cacheMiss:     cacheMiss,
	}
}

func (m *MemCache) Get(alias string) (string, error) {
	fmt.Println(alias)
	v, found := m.cache.Load(alias)
	if found {
		m.cacheHit.Inc()
		m.log.Info(m.log.WithFields(
			logrus.Fields{
				"message": "cache hit",
				"url":     v,
				"alias":   alias,
			},
		))

		if url, ok := v.(string); ok {
			return url, nil
		}

		m.log.Error(m.log.WithFields(
			logrus.Fields{
				"message": "invalid item in cache",
				"url":     v,
				"alias":   alias,
			},
		))
		m.cache.Delete(alias)

		fmt.Println(m)
	}

	m.cacheMiss.Inc()
	m.log.Info(m.log.WithFields(
		logrus.Fields{
			"message": "cache miss",
			"alias":   alias,
		},
	))

	url, err := m.URLRepository.GetURL(alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("could not get item from database: %w", err)
	}

	m.cache.Store(alias, url)
	return url, nil
}

func (m *MemCache) Store(item service.Item) (int64, error) {
	id, err := m.URLRepository.SaveURL(item.URL, item.Alias)
	if err != nil {
		if errors.Is(err, storage.ErrUrlExists) {
			return 0, storage.ErrUrlExists
		}
		return 0, fmt.Errorf("could not save item to database: %w", err)
	}

	m.cache.Store(item.Alias, item.URL)
	return id, nil
}

func (m *MemCache) Delete(alias string) error {
	if err := m.URLRepository.DeleteURL(alias); err != nil {
		return fmt.Errorf("could not delete item from database: %w", err)
	}
	m.cache.Delete(alias)
	return nil
}
