package internal

import (
	"time"

	"github.com/google/uuid"
)

const (
	deleted    string = "[deleted]"
	notitle    string = "[notitle]"
	noauthor   string = "[noauthor]"
	deletedurl string = "https://example.com"
)

type feedType string

const (
	urlType  feedType = "url"
	textType feedType = "text"
)

type feed struct {
	Title       string
	Description string
	Author      string
	Email       string
	FeedType    feedType
}

type item struct {
	FeedTitle string
	Id        uuid.UUID
	Title     string
	Created   time.Time
	Url       string
	Text      string
}
