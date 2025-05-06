package psql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectionPostgres() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:sanya228@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping postgres: %s", err)
	}

	log.Println("Connected to postgres!")

	return db
}
