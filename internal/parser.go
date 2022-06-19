package internal

import (
	"fmt"

	"github.com/badoux/goscraper"
	"github.com/gorilla/feeds"
)

type parser struct {
	url   string
	doc   *goscraper.Document
	cache map[string]*goscraper.Document
}

func newParser() *parser {
	return &parser{cache: make(map[string]*goscraper.Document)}
}

func (p *parser) parse(url string) error {
	if doc, ok := p.cache[url]; ok {
		p.doc = doc
		return nil
	}
	p.url = url
	doc, err := goscraper.Scrape(url, 3)
	if err != nil {
		return err
	}
	p.doc = doc
	p.cache[url] = doc
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
	if len(p.doc.Preview.Description) > 0 {
		return p.doc.Preview.Description
	}
	return fmt.Sprintf(
		"No description provided, please follow <a href=\"%s\">the link</a>.\n",
		p.url)
}
