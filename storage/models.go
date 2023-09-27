package storage

import (
	"bot-saver/package/error"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"time"
)

var ErrNoSavedPage = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	Created  time.Time
}

func (p *Page) Hash() (string, error) {
	h := sha256.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't to WriteString", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't to WriteString", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
