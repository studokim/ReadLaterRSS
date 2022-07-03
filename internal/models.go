package internal

import (
	"html/template"
	"time"
)

type record struct {
	Title string
	Url   string
	Text  string
	When  time.Time
}

type renderedItem struct {
	Id      string
	Title   string
	Url     string
	Text    template.HTML
	Created string
}
