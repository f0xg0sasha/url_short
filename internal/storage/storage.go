package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/f0xg0sasha/url_short/internal/storage/psql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() *Storage {
	db, err := psql.ConnectionPostgres()

	if err != nil {
		log.Fatal("no connected database")
	}

	return &Storage{db: db}
}
