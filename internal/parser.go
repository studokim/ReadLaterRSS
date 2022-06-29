package internal

import (
	"fmt"
	"time"

	"github.com/badoux/goscraper"
	"github.com/gorilla/feeds"
)

type parser struct {
	cache map[string]*goscraper.Document
}

func newParser() *parser {
	return &parser{cache: make(map[string]*goscraper.Document)}
}

func (p *parser) parse(r record) (*feeds.Item, error) {
	doc, err := p.getDoc(r)
	if err != nil {
		return nil, err
	}
	description, err := p.getDescription(r)
	if err != nil {
		return nil, err
	}
	item := &feeds.Item{
		Title:       doc.Preview.Title,
		Link:        &feeds.Link{Href: r.Url},
		Author:      &feeds.Author{Name: doc.Preview.Name},
		Description: description,
		Created:     r.When,
		Id:          r.When.Format(time.RFC3339),
	}
	return item, nil
}

func (p *parser) getDoc(r record) (*goscraper.Document, error) {
	var doc *goscraper.Document
	if document, ok := p.cache[r.Url]; ok {
		doc = document
	} else {
		document, err := goscraper.Scrape(r.Url, 3)
		if err != nil {
			return nil, err
		}
		doc = document
	}
	p.cache[r.Url] = doc
	return doc, nil
}

func (p *parser) getDescription(r record) (string, error) {
	doc, err := p.getDoc(r)
	if err != nil {
		return "", err
	}
	if len(doc.Preview.Description) == 0 && len(r.Text) == 0 {
		return fmt.Sprintf("No description provided, please follow <a href=\"%s\">the link</a>.", r.Url), nil
	}
	if len(r.Text) == 0 {
		return doc.Preview.Description, nil
	}
	if len(doc.Preview.Description) == 0 {
		return r.Text, nil
	}
	return fmt.Sprintf("%s<br><br>%s", r.Text, doc.Preview.Description), nil
}
