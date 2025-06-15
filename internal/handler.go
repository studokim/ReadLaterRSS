package internal

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/feeds"
)

type Handler struct {
	rootUrl    string
	htmlFS     embed.FS
	toRssPipes map[feedType]iToRssPipe
	history    history
}

type result struct {
	Message string
}

func NewHandler(rootUrl string, htmlFS embed.FS) (*Handler, error) {
	history, err := newHistory()
	if err != nil {
		return nil, err
	}
	handler := &Handler{
		rootUrl:    rootUrl,
		htmlFS:     htmlFS,
		toRssPipes: map[feedType]iToRssPipe{url: newUrlToRssPipe(), text: newTextToRssPipe()},
		history:    history,
	}
	handler.registerEndpoints()
	return handler, nil
}

func (h *Handler) registerEndpoints() {
	http.HandleFunc("/", h.index)
	http.HandleFunc("/save", h.save)
	http.HandleFunc("/explore", h.explore)
	http.HandleFunc("/rss", h.rss)
	http.HandleFunc("/feeds", h.feeds)
}

func (h *Handler) handle(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.Write([]byte(err.Error()))
}

func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, pageName string, content any) {
	t, err := template.ParseFS(h.htmlFS, "html/template.html", "html/"+pageName)
	if err != nil {
		h.handle(w, err)
	} else {
		selectedFeed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
		}
		feeds, err := h.history.getFeeds()
		if err != nil {
			h.handle(w, err)
		}
		err = t.ExecuteTemplate(w, "template", struct {
			SelectedFeed feed
			FeedSelector []feed
			Content      any
		}{
			SelectedFeed: selectedFeed,
			FeedSelector: without(feeds, selectedFeed),
			Content:      content,
		})
		if err != nil {
			h.handle(w, err)
		}
	}
}

func (h *Handler) getSelectedFeed(r *http.Request) (feed, error) {
	titleFromUrl := r.URL.Query().Get("feed")
	if titleFromUrl != "" {
		return h.history.getFeed(titleFromUrl)
	}

	title, err := r.Cookie("feed")
	if err != nil && err.Error() == "http: named cookie not present" {
		feeds, err := h.history.getFeeds()
		if err == nil && len(feeds) > 0 {
			return feeds[0], nil
		}
	}
	if err != nil {
		return feed{}, err
	}
	return h.history.getFeed(title.Value)
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, r, "index.html", nil)
}

func (h *Handler) saveUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, r, "saveUrl.html", nil)
	} else {
		r.ParseForm()
		title := "[untitled]"
		url := r.Form["url"][0]
		if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
			url = "http://" + url
		}
		text := ""
		if len(r.Form["describe"]) > 0 {
			text = convertLineBreaks(r.Form["context"][0])
		}

		feed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
		}
		item := item{FeedTitle: feed.Title, Id: uuid.New(), Title: title, Created: time.Now(), Url: url, Text: text}
		rssItem, err := h.toRssPipes[feed.FeedType].pipe(item)
		if err != nil {
			h.handle(w, err)
		}
		item.Title = rssItem.Title
		item.Text = rssItem.Description
		err = h.history.addItem(item)

		if err != nil {
			h.renderPage(w, r, "saveResult.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, r, "saveResult.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) saveText(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, r, "saveText.html", nil)
	} else {
		r.ParseForm()
		title := r.Form["title"][0]
		url := ""
		text := convertLineBreaks(r.Form["text"][0])

		feed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
		}
		item := item{FeedTitle: feed.Title, Id: uuid.New(), Title: title, Created: time.Now(), Url: url, Text: text}
		err = h.history.addItem(item)

		if err != nil {
			h.renderPage(w, r, "saveResult.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, r, "saveResult.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	feed, err := h.getSelectedFeed(r)
	if err != nil {
		h.handle(w, err)
	}
	if feed.FeedType == url {
		h.saveUrl(w, r)
	} else if feed.FeedType == text {
		h.saveText(w, r)
	}
}

func (h *Handler) explore(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("delete") {
		id, err := uuid.Parse(r.URL.Query().Get("delete"))
		if err != nil {
			h.handle(w, err)
		}
		feed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
		}
		err = h.history.deleteItem(item{FeedTitle: feed.Title, Id: id})
		if err != nil {
			h.handle(w, err)
		}
	} else {
		feed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
		}
		itemsNotFiltered, err := h.history.getItems(feed)
		if err != nil {
			h.handle(w, err)
		}
		var items []item
		for _, item := range itemsNotFiltered {
			if item.Title == deleted {
				continue
			}
			item.Text = strings.ReplaceAll(item.Text, "<strike>", "<span class='blured'>")
			item.Text = strings.ReplaceAll(item.Text, "</strike>", "</span>")
			items = append(items, item)
		}
		h.renderPage(w, r, "explore.html", items)
	}
}

func (h *Handler) rss(w http.ResponseWriter, r *http.Request) {
	feed, err := h.getSelectedFeed(r)
	if err != nil {
		h.handle(w, err)
	}
	items, err := h.history.getItems(feed)
	if err != nil {
		h.handle(w, err)
	}
	var rssItems []*feeds.Item
	for _, item := range items {
		rssItem, err := h.toRssPipes[feed.FeedType].pipe(item)
		if err != nil {
			h.handle(w, err)
		}
		rssItems = append(rssItems, rssItem)
		if item.Title == "" {
			item.Title = rssItem.Title // TODO
		}
		if item.Text == "" {
			item.Text = rssItem.Description // TODO
		}
	}
	rssFeed := feeds.Feed{
		Title:       feed.Title,
		Link:        &feeds.Link{Href: h.rootUrl + "/rss"},
		Description: feed.Description,
		Author:      &feeds.Author{Name: feed.Author, Email: feed.Email},
		Created:     time.Now(),
		Items:       rssItems,
	}
	rss, err := rssFeed.ToRss()
	if err != nil {
		h.handle(w, err)
	}
	w.Write([]byte(rss))
}

func (h *Handler) feeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.history.getFeeds()
	if err != nil {
		h.handle(w, err)
	}
	h.renderPage(w, r, "feeds.html", feeds)
}
