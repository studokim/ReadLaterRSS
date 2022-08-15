package internal

import (
	"sync"

	"github.com/gorilla/feeds"
)

type rssFeed struct {
	feed    *feed
	history iHistory
}

func newRssFeed(title string, website string, description string, author string, email string, parser iParser, history iHistory) *rssFeed {
	f := rssFeed{
		feed:    newFeed(title, website, description, author, email, parser),
		history: history,
	}
	for item := range f.buildItemsFromHistory() {
		f.feed.feed.Items = append(f.feed.feed.Items, item)
	}
	return &f
}

func (f *rssFeed) addItem(r record) error {
	err := f.history.add(r)
	if err != nil {
		return err
	}
	return f.feed.addItem(r)
}

func (f *rssFeed) getItems() []*feeds.Item {
	return f.feed.getItems()
}

func (f *rssFeed) getRss() (string, error) {
	rss, err := f.feed.feed.ToRss()
	if err != nil {
		return "", err
	}
	return rss, nil
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
				item, err := f.feed.parser.parse(r)
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
