package internal

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

type Handler struct {
	urlFeed  *rssFeed
	textFeed *rssFeed
	htmlFS   embed.FS
}

type result struct {
	Message string
}

func NewHandler(htmlFS embed.FS, website string, author string, email string) *Handler {
	history, err := newHistory("readlater")
	if err != nil {
		panic(err)
	}
	parser := newUrlParser()
	title := "Read Later"
	description := fmt.Sprintf("%s's list of saved links", author)
	readLaterFeed := newFeed(title, website, description, author, email, parser, history)

	history, err = newHistory("deutsch")
	if err != nil {
		panic(err)
	}
	parser = newTextParser()
	title = "Daily Deutsch"
	description = fmt.Sprintf("%s's daily feed for learning deutsch", author)
	deutschFeed := newFeed(title, website, description, author, email, parser, history)

	return &Handler{
		urlFeed:  readLaterFeed,
		textFeed: deutschFeed,
		htmlFS:   htmlFS,
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
func (h *Handler) getSelectedFeed(w http.ResponseWriter, r *http.Request) string {
	feed, err := r.Cookie("feed")
	if err != nil {
		return "readlater"
	}
	return feed.Value
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "index.html", result{})
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	if h.getSelectedFeed(w, r) == "readlater" {
		h.urlForm(w, r)
	} else {
		h.textForm(w, r)
	}
}

func (h *Handler) urlForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, "saveUrl.html", result{})
	} else {
		r.ParseForm()
		url := r.Form["url"][0]
		var context string
		if len(r.Form["describe"]) > 0 {
			context = convertLineBreaks(r.Form["context"][0])
		} else {
			context = ""
		}
		r := record{Url: url, Text: context, When: time.Now()}
		err := h.urlFeed.addItem(r)
		if err != nil {
			h.renderPage(w, "saveResult.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, "saveResult.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) textForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, "saveText.html", result{})
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
				err := h.textFeed.addItem(r)
				if err != nil {
					maxerr = err
				}
			}
		} else {
			r := record{Title: title, Text: text, When: time.Now()}
			maxerr = h.textFeed.addItem(r)
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
	if r.URL.Query().Get("feed") == "readlater" {
		rss, err = h.urlFeed.getRss()
	} else {
		rss, err = h.textFeed.getRss()
	}
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(rss))
	}
}

func (h *Handler) explore(w http.ResponseWriter, r *http.Request) {
	var rssItems []*feeds.Item
	if h.getSelectedFeed(w, r) == "readlater" {
		rssItems = h.urlFeed.getItems()
	} else {
		rssItems = h.textFeed.getItems()
	}
	renderedItems := make([]renderedItem, len(rssItems))
	for id, item := range rssItems {
		renderedItems[id] = renderedItem{Id: item.Id, Title: item.Title, Url: item.Link.Href,
			Text:    template.HTML(item.Description),
			Created: item.Created.Format(time.RFC822)}
	}
	h.renderPage(w, "explore.html", renderedItems)
}
