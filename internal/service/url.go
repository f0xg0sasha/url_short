package service

import "github.com/f0xg0sasha/url_short/internal/storage"

type URLRepository interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}

type URL struct {
	repo URLRepository
}

func NewURL(repo URLRepository) *URL {
	return &URL{
		repo: repo,
	}
}

func (u *URL) GetURL(alias string) (string, error) {
	return u.repo.GetURL(alias)
}

func (u *URL) SaveURL(urlToSave string, alias string) (int64, error) {
	id, err := u.repo.SaveURL(urlToSave, alias)
	if err != nil {
		if err == storage.ErrUrlExists {
			return 0, storage.ErrUrlExists
		}
		return 0, err
	}

	return id, nil
}

func (u *URL) DeleteURL(alias string) error {
	return u.repo.DeleteURL(alias)
}
