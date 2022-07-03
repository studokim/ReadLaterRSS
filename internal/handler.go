package internal

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Handler struct {
	readLaterFeed *rssFeed
	deutschFeed   *rssFeed
	htmlFS        embed.FS
}

type result struct {
	Message string
}

func NewHandler(htmlFS embed.FS, website string, author string) *Handler {
	history, err := newHistory("later")
	if err != nil {
		panic(err)
	}
	parser := newUrlParser()
	title := "Read Later"
	description := fmt.Sprintf("%s's list of saved links", author)
	readLaterFeed := newFeed(title, website, description, author, parser, history)

	history, err = newHistory("deutsch")
	if err != nil {
		panic(err)
	}
	parser = newTextParser()
	title = "Daily Deutsch"
	description = fmt.Sprintf("%s's daily feed for learning deutsch", author)
	deutschFeed := newFeed(title, website, description, author, parser, history)

	return &Handler{
		readLaterFeed: readLaterFeed,
		deutschFeed:   deutschFeed,
		htmlFS:        htmlFS,
	}
}

func (h *Handler) RegisterEndpoints() {
	http.HandleFunc("/", h.index)
	http.HandleFunc("/add", h.add)
	http.HandleFunc("/deutsch", h.deutsch)
	http.HandleFunc("/rss", h.rss)
	http.HandleFunc("/explore", h.explore)
}

func (h *Handler) renderPage(w http.ResponseWriter, pageName string, content any) {
	t, err := template.ParseFS(h.htmlFS, "html/"+pageName, "html/template.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		t.ExecuteTemplate(w, "template", content)
	}
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "index.html", result{})
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, "add.html", result{})
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
		err := h.readLaterFeed.addItem(r)
		if err != nil {
			h.renderPage(w, "result.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, "result.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) deutsch(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, "deutsch.html", result{})
	} else {
		r.ParseForm()
		title := r.Form["title"][0]
		text := convertLineBreaks(r.Form["text"][0])
		r := record{Title: title, Text: text, When: time.Now()}
		err := h.deutschFeed.addItem(r)
		if err != nil {
			h.renderPage(w, "result.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, "result.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) rss(w http.ResponseWriter, r *http.Request) {
	rss, err := h.readLaterFeed.getRss()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(rss))
	}
}

func (h *Handler) explore(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "explore.html", h.readLaterFeed.getItems())
}
