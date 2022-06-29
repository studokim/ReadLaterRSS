package internal

import "strings"

func ConvertLineBreaks(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "\r\n", "<br>", -1)
	s = strings.Replace(s, "\r", "<br>", -1)
	s = strings.Replace(s, "\n", "<br>", -1)
	return s
}
