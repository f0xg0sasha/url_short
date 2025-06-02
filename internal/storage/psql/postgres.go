package psql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectionPostgres() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:sanya228@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to postgres!")

	return db, nil
}
