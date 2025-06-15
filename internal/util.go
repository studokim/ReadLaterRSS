package internal

import (
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

func removeParagraphBreaks(s string) string {
	return strings.Replace(s, "<br>", ". ", -1)
}

func splitOnSentences(text string) []string {
	var sentences []string
	sentenceBreakRegexp := `(([.?!]|\.{3})(\s|<br>|$)+|<br>+|$)`
	r := regexp.MustCompile(sentenceBreakRegexp)
	splitted := r.Split(text, -1)
	for _, sentence := range splitted {
		if len(sentence) != 0 {
			sentences = append(sentences, sentence)
		}
	}
	return sentences
}

func replaceDotsInDeutschDates(text string) string {
	deutschDateRegexp := `(\d?\d)\. ?(Januar|Februar|März|April|Mai|Juni|Juli|August|September|Oktober|November|Dezember)`
	r := regexp.MustCompile(deutschDateRegexp)
	return r.ReplaceAllString(text, "$1 $2")
}

func without(items []feed, item feed) []feed {
	var result []feed
	for _, it := range items {
		if it.Title != item.Title {
			result = append(result, it)
		}
	}
	return result
}
