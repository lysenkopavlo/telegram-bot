package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
)

var ErrNoSavedPages = errors.New("no saved pages")

type Storage interface {
	Save(context.Context, *Page) error
	PickRandom(context.Context, string) (*Page, error)
	Remove(context.Context, *Page) error
	IsExists(context.Context, *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.WrapError("Error while calculating hash: %w", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.WrapError("Error while calculating hash: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
