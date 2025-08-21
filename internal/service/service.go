package service

import "context"

type Item struct {
	URL   string
	Alias string
}

type Service struct {
	repo CacheRepository
}

type CacheRepository interface {
	Store(ctx context.Context, item Item) (int64, error)
	Get(ctx context.Context, alias string) (string, error)
	Delete(ctx context.Context, alias string) error
}

func NewService(repo CacheRepository) *Service {
	return &Service{
		repo: repo,
	}
}
