package internal

import (
	"time"

	"github.com/google/uuid"
)

type feedType string

const (
	url  feedType = "url"
	text feedType = "text"
)

type feed struct {
	Title       string
	Description string
	Author      string
	Email       string
	FeedType    feedType
}

const (
	deleted string = "[deleted]"
)

type item struct {
	FeedTitle string
	Id        uuid.UUID
	Title     string
	Created   time.Time
	Url       string
	Text      string
}
