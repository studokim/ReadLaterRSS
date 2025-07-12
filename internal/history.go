package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type history struct {
	db *sql.DB
}

func newHistory() (history, error) {
	db, err := Sqlite{}.db()
	if err != nil {
		return history{}, err
	}
	return history{db: db}, nil
}

func (h *history) addFeed(f feed) error {
	_, err := h.db.Exec("insert into feed(title, description, author, email, feedType) values(?, ?, ?, ?, ?)", f.Title, f.Description, f.Author, f.Email, f.FeedType)
	return err
}

func (h *history) deleteFeed(f feed) error {
	_, err := h.db.Exec("delete from item where feedTitle=?", f.Title)
	if err != nil {
		return err
	}
	_, err = h.db.Exec("delete from feed where title=?", f.Title)
	return err
}

func (h *history) getFeed(title string) (feed, error) {
	res := h.db.QueryRow("select description, author, email, feedType from feed where title=?", title)
	var description string
	var author string
	var email string
	var feedType feedType
	err := res.Scan(&description, &author, &email, &feedType)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return feed{}, errors.New(fmt.Sprint("No such feed: ", title))
		}
		return feed{}, err
	}
	return feed{Title: title, Description: description, Author: author, Email: email, FeedType: feedType}, nil
}

func (h *history) getFeeds() ([]feed, error) {
	res, err := h.db.Query("select title, description, author, email, feedType from feed")
	if err != nil {
		return nil, err
	}
	var feeds []feed
	for res.Next() {
		var title string
		var description string
		var author string
		var email string
		var feedType feedType
		err = res.Scan(&title, &description, &author, &email, &feedType)
		if err != nil {
			return nil, err
		}
		if feedType != url && feedType != text {
			return nil, err
		}
		feeds = append(feeds, feed{Title: title, Description: description, Author: author, Email: email, FeedType: feedType})
	}
	return feeds, nil
}

func (h *history) addItem(r item) error {
	_, err := h.db.Exec("insert into item(feedTitle, id, title, created, url, text) values(?, ?, ?, ?, ?, ?)", r.FeedTitle, r.Id, r.Title, r.Created, r.Url, r.Text)
	return err
}

func (h *history) deleteItem(r item) error {
	// RSS doesn't allow to delete items, because RSS readers won't be able to distinguish between deleted and just too old items
	_, err := h.db.Exec("update item set title=?, url='', text='' where id=?", deleted, r.Id)
	return err
}

func (h *history) getItems(f feed) ([]item, error) {
	res, err := h.db.Query("select id, title, created, url, text from item where feedTitle=?", f.Title)
	if err != nil {
		return nil, err
	}
	var items []item
	for res.Next() {
		var id uuid.UUID
		var title string
		var createdStr string
		var created time.Time
		var url string
		var text string
		err = res.Scan(&id, &title, &createdStr, &url, &text)
		if err != nil {
			return nil, err
		}
		created, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", createdStr)
		if err != nil {
			return nil, err
		}
		items = append(items, item{FeedTitle: f.Title, Id: id, Title: title, Created: created, Url: url, Text: text})
	}
	return items, nil
}
