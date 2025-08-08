package service

type Item struct {
	Alias string
	URL   string
}

type Service struct {
	repo CacheRepository
}

type CacheRepository interface {
	Store(item Item) (int64, error)
	Get(alias string) (string, error)
	Delete(alias string) error
}

func NewService(repo CacheRepository) *Service {
	return &Service{
		repo: repo,
	}
}
