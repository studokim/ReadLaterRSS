package internal

import (
	"time"

	"github.com/gorilla/feeds"
)

type textParser struct {
	translator translator
}

func newTextParser() iParser {
	return &textParser{translator: translator{}}
}

func (p *textParser) parse(r record) (*feeds.Item, error) {
	item := &feeds.Item{
		Title:       r.Title,
		Link:        &feeds.Link{},
		Author:      &feeds.Author{},
		Description: r.Text,
		Created:     r.When,
		Id:          r.When.Format(time.RFC3339),
	}

	return item, nil
}
