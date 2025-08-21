package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	const op = "storage.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO url (url, alias) VALUES ($1, $2)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(urlToSave, alias)

	if err != nil {
		fmt.Println(err)
		if dbErr, ok := err.(*pq.Error); ok && dbErr.Code == "23505" {

			return 0, ErrUrlExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64

	err = s.db.QueryRowContext(ctx, "SELECT id FROM url WHERE alias = $1", alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var resURL string
	err = stmt.QueryRowContext(ctx, alias).Scan(&resURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
