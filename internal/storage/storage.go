package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/commedesvlados/go-tg-bot/pkg/lib/e"
	"io"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, username string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

var ErrNoSavedPages = errors.New("no saved page")

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash ", err)
	}

	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", e.Wrap("can't calculate hash ", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
