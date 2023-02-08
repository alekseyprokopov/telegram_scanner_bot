package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	Pick(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	//из []byte в стринг
	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
