package internal

import (
	"github.com/gorilla/feeds"
)

type iHistory interface {
	add(record) error
	getSize() int
	getRecords() []record
}

type iParser interface {
	parse(record) (*feeds.Item, error)
}
