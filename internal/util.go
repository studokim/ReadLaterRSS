package internal

import (
	"os"
	"regexp"
	"strings"
)

func convertLineBreaks(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", " ", -1)
	s = strings.Replace(s, "\r\n", "<br>", -1)
	s = strings.Replace(s, "\r", "<br>", -1)
	s = strings.Replace(s, "\n", "<br>", -1)
	return s
}

func splitOnParagraphs(text string) []string {
	var paragraphs []string
	splitted := strings.Split(text, "<br>")
	for _, paragraph := range splitted {
		if len(paragraph) != 0 {
			paragraphs = append(paragraphs, paragraph)
		}
	}
	return paragraphs
}

func splitOnSentences(text string) []string {
	paragraphs := splitOnParagraphs(text)
	var sentences []string
	r := regexp.MustCompile(`([\!\?\.] |[\!\?\.])`)
	for _, paragraph := range paragraphs {
		splitted := r.Split(paragraph, -1)
		for _, sentence := range splitted {
			if len(sentence) != 0 {
				sentences = append(sentences, sentence)
			}
		}
	}
	return sentences
}

func readFile(fileName string) ([]byte, error) {
	if _, err := os.Stat(fileName); err != nil {
		os.Create(fileName)
	}
	return os.ReadFile(fileName)
}

func openFileAppend(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
}
