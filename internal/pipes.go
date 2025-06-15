package internal

import (
	"fmt"
	"strings"

	"github.com/badoux/goscraper"
	"github.com/gorilla/feeds"
)

type iToRssPipe interface {
	pipe(item) (*feeds.Item, error)
}

type textToRssPipe struct {
	translator   translator
	errorMessage string
}

func newTextToRssPipe() iToRssPipe {
	return &textToRssPipe{translator: newTranslator(), errorMessage: "Translator unreachable."}
}

func (p *textToRssPipe) pipe(r item) (*feeds.Item, error) {
	item := &feeds.Item{
		Id:          r.Id.String(),
		Title:       r.Title,
		Description: p.getText(r),
		Created:     r.Created,
		Link:        &feeds.Link{},
		Author:      &feeds.Author{},
	}

	return item, nil
}

func (p *textToRssPipe) getText(r item) string {
	translated, err := p.translator.translate(removeParagraphBreaks(r.Text))
	if err != nil {
		return fmt.Sprintf("%s<br><br><em>[%s]</em>", r.Text, p.errorMessage)
	}
	text := replaceDotsInDeutschDates(r.Text)
	textSentences := splitOnSentences(text)
	translatedSentences := splitOnSentences(translated)
	if len(textSentences) != len(translatedSentences) {
		return fmt.Sprintf("%s<br><br><strike>[%s]</strike>", text, translated)
	}
	for i := range textSentences {
		oldSentence := textSentences[i]
		newSentence := fmt.Sprintf("%s <strike>[%s]</strike>", oldSentence, translatedSentences[i])
		for _, r := range []string{".", "?", "!", "<br>"} {
			text = strings.Replace(text, oldSentence+r, newSentence+r, 1)
		}
	}
	return text
}

type urlToRssPipe struct {
	cache map[string]*goscraper.Document
}

func newUrlToRssPipe() iToRssPipe {
	return &urlToRssPipe{cache: make(map[string]*goscraper.Document)}
}

func (p *urlToRssPipe) pipe(it item) (*feeds.Item, error) {
	if it.Title == deleted {
		return &feeds.Item{
			Id:          it.Id.String(),
			Title:       deleted,
			Description: it.Text,
			Created:     it.Created,
			Link:        &feeds.Link{Href: it.Url},
			Author:      &feeds.Author{Name: deleted},
		}, nil
	}
	doc, err := p.getDoc(it)
	if err != nil {
		return nil, err
	}
	description, err := p.getDescription(it)
	if err != nil {
		return nil, err
	}
	return &feeds.Item{
		Id:          it.Id.String(),
		Title:       doc.Preview.Title,
		Description: description,
		Created:     it.Created,
		Link:        &feeds.Link{Href: it.Url},
		Author:      &feeds.Author{Name: doc.Preview.Name},
	}, nil
}

func (p *urlToRssPipe) getDoc(it item) (*goscraper.Document, error) {
	var doc *goscraper.Document
	if document, ok := p.cache[it.Url]; ok {
		doc = document
	} else {
		document, err := goscraper.Scrape(it.Url, 3)
		if err != nil {
			return nil, err
		}
		doc = document
	}
	p.cache[it.Url] = doc
	return doc, nil
}

func (p *urlToRssPipe) getDescription(it item) (string, error) {
	doc, err := p.getDoc(it)
	if err != nil {
		return "", err
	}
	if len(doc.Preview.Description) == 0 && len(it.Text) == 0 {
		return fmt.Sprintf("No description provided, please follow <a href=\"%s\">the link</a>.", it.Url), nil
	}
	if len(it.Text) == 0 {
		return doc.Preview.Description, nil
	}
	if len(doc.Preview.Description) == 0 {
		return it.Text, nil
	}
	return fmt.Sprintf("%s<br><br>%s", it.Text, doc.Preview.Description), nil
}
