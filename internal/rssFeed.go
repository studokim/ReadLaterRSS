package internal

import (
	"sort"
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

type rssFeed struct {
	feed    *feeds.Feed
	parser  iParser
	history iHistory
}

func newFeed(title string, website string, description string, author string, parser iParser, history iHistory) *rssFeed {
	f := rssFeed{
		feed: &feeds.Feed{
			Title:       title,
			Link:        &feeds.Link{Href: website},
			Description: description,
			Author:      &feeds.Author{Name: author},
			Created:     time.Now(),
			Items:       []*feeds.Item{},
		},
		parser:  parser,
		history: history,
	}
	for item := range f.buildItemsFromHistory() {
		f.feed.Items = append(f.feed.Items, item)
	}
	return &f
}

func (f *rssFeed) addItem(r record) error {
	item, err := f.parser.parse(r)
	if err != nil {
		return err
	}
	f.feed.Items = append(f.feed.Items, item)
	err = f.history.add(r)
	return err
}

func (f *rssFeed) getRss() (string, error) {
	rss, err := f.feed.ToRss()
	if err != nil {
		return "", err
	}
	return rss, nil
}

func (f *rssFeed) getItems() []*feeds.Item {
	items := f.feed.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Id > items[j].Id
	})
	return items
}

func (f *rssFeed) buildItemsFromHistory() <-chan *feeds.Item {
	size := f.history.getSize()
	records := make(chan record, size)
	items := make(chan *feeds.Item, size)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			for r := range records {
				item, err := f.parser.parse(r)
				if err == nil {
					items <- item
				}
			}
			wg.Done()
		}()
	}
	for _, r := range f.history.getRecords() {
		records <- r
	}
	close(records)
	wg.Wait()
	close(items)
	return items
}
