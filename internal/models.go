package internal

import "time"

type record struct {
	Title string
	Url   string
	Text  string
	When  time.Time
}
