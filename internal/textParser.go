package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type textParser struct {
	translator   translator
	errorMessage string
}

func newTextParser() iParser {
	return &textParser{translator: newTranslator(), errorMessage: "Translator unreachable."}
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
		return fmt.Sprintf("%s<br><br><em>[%s]</em>", r.Text, p.errorMessage)
	}
	textSentences := splitOnSentences(r.Text)
	translatedSentences := splitOnSentences(translated)
	if len(textSentences) != len(translatedSentences) {
		return fmt.Sprintf("%s<br><br><strike>[%s]</strike>", r.Text, translated)
	}
	text := r.Text
	for i := range textSentences {
		oldSentence := textSentences[i]
		newSentence := fmt.Sprintf("%s <strike>[%s]</strike>", oldSentence, translatedSentences[i])
		for _, r := range []rune{'.', '?', '!'} {
			text = strings.Replace(text, oldSentence+string(r), newSentence+string(r), 1)
		}
	}
	return text
}
