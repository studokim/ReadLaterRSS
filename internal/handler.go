package internal

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type Handler struct {
	sharedFeed    *feed
	readlaterFeed *rssFeed
	deutschFeed   *rssFeed
	htmlFS        embed.FS
}

type result struct {
	Message string
}

type button struct {
	Verb string
}

func NewHandler(htmlFS embed.FS, website string, author string, email string) *Handler {
	parser := newUrlParser()
	title := "Shared"
	description := fmt.Sprintf("%s's list of temporarily shared", author)
	shareFeed := newFeed(title, website, description, author, email, parser)

	history, err := newHistory("readlater")
	if err != nil {
		panic(err)
	}
	parser = newUrlParser()
	title = "Read Later"
	description = fmt.Sprintf("%s's list of saved links", author)
	readLaterFeed := newRssFeed(title, website, description, author, email, parser, history)

	history, err = newHistory("deutsch")
	if err != nil {
		panic(err)
	}
	parser = newTextParser()
	title = "Daily Deutsch"
	description = fmt.Sprintf("%s's daily feed for learning deutsch", author)
	deutschFeed := newRssFeed(title, website, description, author, email, parser, history)

	return &Handler{
		sharedFeed:    shareFeed,
		readlaterFeed: readLaterFeed,
		deutschFeed:   deutschFeed,
		htmlFS:        htmlFS,
	}
}

func (h *Handler) RegisterEndpoints() {
	http.HandleFunc("/", h.index)
	http.HandleFunc("/save", h.save)
	http.HandleFunc("/explore", h.explore)
	http.HandleFunc("/rss", h.rss)
}

func (h *Handler) renderPage(w http.ResponseWriter, pageName string, content any) {
	t, err := template.ParseFS(h.htmlFS, "html/"+pageName, "html/template.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		t.ExecuteTemplate(w, "template", content)
	}
}

func (h *Handler) getSelectedFeed(r *http.Request) string {
	feed, err := r.Cookie("feed")
	if err != nil {
		return "shared"
	}
	return feed.Value
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "index.html", nil)
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	if h.getSelectedFeed(r) == "deutsch" {
		h.textForm(w, r)
	} else if h.getSelectedFeed(r) == "readlater" {
		h.urlForm(w, r, button{Verb: "Save!"})
	} else {
		h.urlForm(w, r, button{Verb: "Share!"})
	}
}

func (h *Handler) urlForm(w http.ResponseWriter, r *http.Request, button button) {
	if r.Method == "GET" {
		h.renderPage(w, "saveUrl.html", button)
	} else {
		r.ParseForm()
		url := r.Form["url"][0]
		var context string
		if len(r.Form["describe"]) > 0 {
			context = convertLineBreaks(r.Form["context"][0])
		} else {
			context = ""
		}
		record := record{Url: url, Text: context, When: time.Now()}
		feed := h.getSelectedFeed(r)
		var err error
		if feed == "readlater" {
			err = h.readlaterFeed.addItem(record)
		} else if feed == "shared" {
			err = h.sharedFeed.addItem(record)
		} else {
			err = errors.New("selected feed doesn't implement adding to itself")
		}
		if err != nil {
			h.renderPage(w, "saveResult.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, "saveResult.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) textForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, "saveText.html", nil)
	} else {
		r.ParseForm()
		title := r.Form["title"][0]
		text := convertLineBreaks(r.Form["text"][0])
		var maxerr error
		if len(r.Form["split"]) > 0 {
			parapraphs := splitOnParagraphs(text)
			count := len(parapraphs)
			for i, paragraph := range parapraphs {
				created := time.Now().Add(time.Second * time.Duration(i))
				r := record{Title: fmt.Sprintf("%s (%d/%d)", title, i+1, count), Text: paragraph, When: created}
				err := h.deutschFeed.addItem(r)
				if err != nil {
					maxerr = err
				}
			}
		} else {
			r := record{Title: title, Text: text, When: time.Now()}
			maxerr = h.deutschFeed.addItem(r)
		}
		if maxerr != nil {
			h.renderPage(w, "saveResult.html", result{Message: maxerr.Error()})
		} else {
			h.renderPage(w, "saveResult.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) rss(w http.ResponseWriter, r *http.Request) {
	var rss string
	var err error
	// choose the feed according to the query string
	feed := r.URL.Query().Get("feed")
	if feed == "deutsch" {
		rss, err = h.deutschFeed.getRss()
	} else if feed == "readlater" {
		rss, err = h.readlaterFeed.getRss()
	} else {
		err = errors.New("selected feed doesn't implement the RSS functionality")
	}
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(rss))
	}
}

func (h *Handler) explore(w http.ResponseWriter, r *http.Request) {
	var items []*feeds.Item
	if h.getSelectedFeed(r) == "readlater" {
		items = h.readlaterFeed.getItems()
	} else if h.getSelectedFeed(r) == "deutsch" {
		items = h.deutschFeed.getItems()
	} else {
		items = h.sharedFeed.getItems()
	}
	renderedItems := make([]renderedItem, len(items))
	for id, item := range items {
		text := item.Description
		text = strings.ReplaceAll(text, "<strike>", "<span class=\"blured\">")
		text = strings.ReplaceAll(text, "</strike>", "</span>")
		renderedItems[id] = renderedItem{Id: item.Id, Title: item.Title, Url: item.Link.Href,
			Text:    template.HTML(text),
			Created: item.Created.Format(time.RFC822)}
	}
	h.renderPage(w, "explore.html", renderedItems)
}
