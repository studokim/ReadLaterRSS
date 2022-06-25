package internal

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

type rssFeed struct {
	feed    *feeds.Feed
	parser  *parser
	history history
}

func newFeed(website string, author string) *rssFeed {
	history, err := newHistory()
	if err != nil {
		panic(err)
	}
	parser := newParser()
	f := rssFeed{
		feed: &feeds.Feed{
			Title:       "Read Later",
			Link:        &feeds.Link{Href: website},
			Description: fmt.Sprintf("%s's list of saved links", author),
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

func (f *rssFeed) addItem(url string, context string) error {
	now := time.Now()
	item, err := f.buildItem(url, context, now)
	if err != nil {
		return err
	}
	f.feed.Items = append(f.feed.Items, item)
	err = f.history.add(url, context, now)
	if err != nil {
		return err
	}
	return nil
}

func (f *rssFeed) buildItem(url string, context string, created time.Time) (*feeds.Item, error) {
	err := f.parser.parse(url, context)
	if err != nil {
		return nil, err
	}
	return &feeds.Item{
		Title:       f.parser.getTitle(),
		Link:        &feeds.Link{Href: url},
		Description: f.parser.getDescription(),
		Author:      f.parser.getAuthor(),
		Created:     created,
		Id:          created.Format(time.RFC3339),
	}, nil
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
	size := len(f.history)
	records := make(chan record, size)
	items := make(chan *feeds.Item, size)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			for record := range records {
				item, err := f.buildItem(record.Url, record.Context, record.When)
				if err == nil {
					items <- item
				}
			}
			wg.Done()
		}()
	}
	for _, record := range f.history {
		records <- record
	}
	close(records)
	wg.Wait()
	close(items)
	return items
}
