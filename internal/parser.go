package internal

import (
	"github.com/badoux/goscraper"
	"github.com/gorilla/feeds"
)

type parser struct {
	doc *goscraper.Document
}

func newParser() *parser {
	return &parser{}
}

func (p *parser) parse(url string) error {
	doc, err := goscraper.Scrape(url, 3)
	if err != nil {
		return err
	}
	p.doc = doc
	return nil
}

func (p *parser) getTitle() string {
	return p.doc.Preview.Title
}

func (p *parser) getAuthor() *feeds.Author {
	return &feeds.Author{
		Name: p.doc.Preview.Name,
	}
}

func (p *parser) getDescription() string {
	return p.doc.Preview.Description
}
