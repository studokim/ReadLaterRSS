package internal

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func convertLineBreaks(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "<br>")
	s = strings.ReplaceAll(s, "\r", "<br>")
	s = strings.ReplaceAll(s, "\n", "<br>")
	return s
}

func splitOnSentences(text string) []string {
	var sentences []string
	sentenceBreak := regexp.MustCompile(`([.?!]|\.{3})(\s|$)`)
	splitted := sentenceBreak.Split(text, -1)
	for _, sentence := range splitted {
		if len(sentence) != 0 {
			sentences = append(sentences, sentence)
		}
	}
	return sentences
}

func fixCommonParsingProblems(s string) string {
	s = strings.TrimSpace(s)
	deutschDateRegexp := `(\d?\d)\. ?(Januar|Februar|März|April|Mai|Juni|Juli|August|September|Oktober|November|Dezember)`
	r := regexp.MustCompile(deutschDateRegexp)
	s = r.ReplaceAllString(s, "$1 $2")
	s = strings.ReplaceAll(s, " ", " ")
	s = strings.ReplaceAll(s, " z. B. ", " z.b. ")
	return s
}

func getAtMostOne(form url.Values, param string) (string, error) {
	switch len(form[param]) {
	case 0:
		return "", nil
	case 1:
		return form[param][0], nil
	default:
		return "", fmt.Errorf("expected at most one `%s` value, got: %d", param, len(form[param]))
	}
}

func getSingle(form url.Values, param string) (string, error) {
	if len(form[param]) != 1 {
		return "", fmt.Errorf("expected a single `%s` value, got: %d", param, len(form[param]))
	}
	return form[param][0], nil
}
