package psql

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var dataSourceName = "postgres://postgres:" + "qwerty" + "@localhost:5432/postgres?sslmode=disable"

func ConnectionPostgres() (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)

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
