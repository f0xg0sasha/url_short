package storage

import (
	"database/sql"

	"girhub.com/f0xg0sasha/url_short/internal/storage/psql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() *Storage {
	db := psql.ConnectionPostgres()
	return &Storage{db: db}
}
