package internal

import (
	"time"

	"github.com/gorilla/feeds"
)

type readLaterFeed struct {
	feed    *feeds.Feed
	parser  *parser
	history history
}

func newFeed() *readLaterFeed {
	history, err := newHistory()
	if err != nil {
		panic(err)
	}
	parser := newParser()
	f := readLaterFeed{
		feed: &feeds.Feed{
			Title:   "Read Later",
			Link:    &feeds.Link{Href: "https://muravev.space"},
			Author:  &feeds.Author{Name: "kim", Email: "me@muravev.space"},
			Created: time.Now(),
			Items:   []*feeds.Item{},
		},
		parser:  parser,
		history: history,
	}
	items := []*feeds.Item{}
	for time, url := range history {
		item, err := f.buildItem(url, time)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	f.feed.Items = items
	return &f
}

func (f *readLaterFeed) addItem(url string) error {
	now := time.Now()
	item, err := f.buildItem(url, now)
	if err != nil {
		return err
	}
	f.feed.Items = append(f.feed.Items, item)
	err = f.history.add(url, now)
	if err != nil {
		return err
	}
	return nil
}

func (f *readLaterFeed) buildItem(url string, created time.Time) (*feeds.Item, error) {
	err := f.parser.parse(url)
	if err != nil {
		return nil, err
	}
	return &feeds.Item{
		Title:       f.parser.getTitle(),
		Link:        &feeds.Link{Href: url},
		Description: f.parser.getDescription(),
		Author:      f.parser.getAuthor(),
		Created:     created,
	}, nil
}

func (f *readLaterFeed) getRss() (string, error) {
	rss, err := f.feed.ToRss()
	if err != nil {
		return "", err
	}
	return rss, nil
}
