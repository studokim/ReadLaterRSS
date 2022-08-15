package internal

import (
	"sort"
	"time"

	"github.com/gorilla/feeds"
)

type feed struct {
	feed   *feeds.Feed
	parser iParser
}

func newFeed(title string, website string, description string, author string, email string, parser iParser) *feed {
	f := feed{
		feed: &feeds.Feed{
			Title:       title,
			Link:        &feeds.Link{Href: website},
			Description: description,
			Author:      &feeds.Author{Name: author, Email: email},
			Created:     time.Now(),
			Items:       []*feeds.Item{},
		},
		parser: parser,
	}
	return &f
}

func (f *feed) addItem(r record) error {
	item, err := f.parser.parse(r)
	if err != nil {
		return err
	}
	f.feed.Items = append(f.feed.Items, item)
	return err
}

func (f *feed) getItems() []*feeds.Item {
	items := f.feed.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Id > items[j].Id
	})
	return items
}
