package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type textParser struct {
	translator translator
}

func newTextParser() iParser {
	return &textParser{translator: newTranslator()}
}

func (p *textParser) parse(r record) (*feeds.Item, error) {
	item := &feeds.Item{
		Title:       r.Title,
		Link:        &feeds.Link{},
		Author:      &feeds.Author{},
		Description: p.getText(r),
		Created:     r.When,
		Id:          r.When.Format(time.RFC3339),
	}

	return item, nil
}

func (p *textParser) getText(r record) string {
	translated, err := p.translator.translate(r.Text)
	if err != nil {
		return r.Text
	}
	sentences := splitOnSentences(r.Text)
	translatedSentences := splitOnSentences(translated)
	if len(sentences) != len(translatedSentences) {
		return fmt.Sprintf("%s<br><br><em>[%s]</em>", r.Text, translated)
	}
	text := r.Text
	for id := range sentences {
		sentence := fmt.Sprintf("%s <strike>[%s]</strike>", sentences[id], translatedSentences[id])
		text = strings.Replace(text, sentences[id], sentence, 1)
	}
	return text
}
