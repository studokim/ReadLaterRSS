package internal

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Handler struct {
	rssFeed *rssFeed
	htmlFS  embed.FS
}

type result struct {
	Message string
}

func NewHandler(htmlFS embed.FS, website string, author string) *Handler {
	history, err := newHistory()
	if err != nil {
		panic(err)
	}
	parser := newParser()
	title := "Read Later"
	description := fmt.Sprintf("%s's list of saved links", author)
	return &Handler{
		rssFeed: newFeed(title, website, description, author, parser, history),
		htmlFS:  htmlFS,
	}
}

func (h *Handler) RegisterEndpoints() {
	http.HandleFunc("/", h.index)
	http.HandleFunc("/add", h.add)
	http.HandleFunc("/rss", h.rss)
	http.HandleFunc("/feed", h.feed)
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
		h.renderPage(w, "form.html", result{})
	} else {
		r.ParseForm()
		url := r.Form["url"][0]
		var context string
		if len(r.Form["describe"]) > 0 {
			context = ConvertLineBreaks(r.Form["context"][0])
		} else {
			context = ""
		}
		r := record{Url: url, Context: context, When: time.Now()}
		err := h.rssFeed.addItem(r)
		if err != nil {
			h.renderPage(w, "result.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, "result.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) rss(w http.ResponseWriter, r *http.Request) {
	rss, err := h.rssFeed.getRss()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(rss))
	}
}

func (h *Handler) feed(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "feed.html", h.rssFeed.getItems())
}
