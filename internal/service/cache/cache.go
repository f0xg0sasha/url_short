package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/f0xg0sasha/url_short/internal/service"
	"github.com/f0xg0sasha/url_short/internal/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type URLRepository interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error)
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}

type Cache struct {
	log       *logrus.Logger
	client    *redis.Client
	cacheHit  prometheus.Counter
	cacheMiss prometheus.Counter
	URLRepository
}

func NewCache(
	log *logrus.Logger,
	client *redis.Client,
	repo URLRepository,
	cacheHit prometheus.Counter,
	cacheMiss prometheus.Counter,
) *Cache {
	return &Cache{
		log:           log,
		client:        client,
		URLRepository: repo,
		cacheHit:      cacheHit,
		cacheMiss:     cacheMiss,
	}
}

func (m *Cache) Get(ctx context.Context, alias string) (string, error) {
	fmt.Println(alias)
	v, err := m.client.Get(ctx, "url:"+alias).Result()
	if !errors.Is(err, redis.Nil) {
		m.cacheHit.Inc()
		m.log.Info(m.log.WithFields(
			logrus.Fields{
				"message": "cache hit",
				"url":     v,
				"alias":   alias,
			},
		))

		return v, nil
	} else {
		m.log.Error(m.log.WithFields(
			logrus.Fields{
				"message": "invalid item in cache",
				"url":     v,
				"alias":   alias,
			},
		))
		m.client.Del(ctx, "url:"+alias)
	}

	m.cacheMiss.Inc()
	m.log.Info(m.log.WithFields(
		logrus.Fields{
			"message": "cache miss",
			"alias":   alias,
		},
	))

	url, err := m.URLRepository.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("could not get item from database: %w", err)
	}

	_, err = m.client.Set(ctx, "url:"+alias, url, 0).Result()
	if err != nil {
		m.log.Error(m.log.WithFields(
			logrus.Fields{
				"message": "can't set cache in redis",
				"url":     v,
				"alias":   alias,
			},
		))
		return "", nil
	}

	return url, nil
}

func (m *Cache) Store(ctx context.Context, item service.Item) (int64, error) {
	id, err := m.URLRepository.SaveURL(ctx, item.URL, item.Alias)
	if err != nil {
		if errors.Is(err, storage.ErrUrlExists) {
			return 0, storage.ErrUrlExists
		}
		return 0, fmt.Errorf("could not save item to database: %w", err)
	}

	m.client.Set(ctx, "url:"+item.Alias, item.URL, 0)
	return id, nil
}

func (m *Cache) Delete(ctx context.Context, alias string) error {
	if err := m.URLRepository.DeleteURL(ctx, alias); err != nil {
		return fmt.Errorf("could not delete item from database: %w", err)
	}
	m.client.Del(ctx, "url:"+alias)
	return nil
}
