package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/f0xg0sasha/url_short/internal/storage"
)

func (s *Service) Fetch(ctx context.Context, alias string) (string, error) {
	url, err := s.repo.Get(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("could not fetch url: %w", err)
	}
	return url, nil
}

func (s *Service) Create(ctx context.Context, url string, alias string) (int64, error) {
	id, err := s.repo.Store(ctx, Item{URL: url, Alias: alias})
	if err != nil {
		if errors.Is(err, storage.ErrUrlExists) {
			return 0, storage.ErrUrlExists
		} else {
			return 0, fmt.Errorf("could not create url: %w", err)
		}
	}
	return id, nil
}

func (s *Service) Delete(ctx context.Context, alias string) error {
	if err := s.repo.Delete(ctx, alias); err != nil {
		return fmt.Errorf("could not delete url: %w", err)
	}
	return nil
}
