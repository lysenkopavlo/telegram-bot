package storage

import (
	"crypto/sha1"
	"io"

	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
)

type Storage interface {
	Save(*Page) error
	PickRandom(string) (*Page, error)
	Remove(*Page) error
	IsExists(*Page) (bool, error)
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

	return string(h.Sum(nil)), nil

}
