package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteClient(path string) (*sql.DB, error) {
	client, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open to database: %w", err)
	}

	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("can't conect to database: %w", err)
	}

	return client, nil
}
