package sqlite

import (
	"context"
	"database/sql"
	"github.com/commedesvlados/go-tg-bot/internal/storage"
	"github.com/commedesvlados/go-tg-bot/pkg/lib/e"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db} // NewSqliteClient
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, username string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, username).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap("can't pick random page", err)
	}

	return &storage.Page{URL: url, Username: username}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return e.Wrap("can't remove page", err)
	}

	return nil
}

func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := ` SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.Username).Scan(&count); err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return e.Wrap("can't create table", err)
	}

	return nil
}
