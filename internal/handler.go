package internal

import (
	"embed"
	"net/http"
	"text/template"
)

type Handler struct {
	feed   *readLaterFeed
	htmlFS embed.FS
}

type result struct {
	Message string
}

func NewHandler(htmlFS embed.FS, website string, author string) *Handler {
	return &Handler{
		feed:   newFeed(website, author),
		htmlFS: htmlFS,
	}
}

func (h *Handler) renderPage(w http.ResponseWriter, pageName string, r result) {
	t, err := template.ParseFS(h.htmlFS, "html/"+pageName, "html/template.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		t.ExecuteTemplate(w, "template", r)
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "index.html", result{})
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderPage(w, "form.html", result{})
	} else {
		r.ParseForm()
		url := r.Form["url"][0]
		var context string
		if len(r.Form["describe"]) > 0 {
			context = r.Form["context"][0]
		} else {
			context = ""
		}
		err := h.feed.addItem(url, context)
		if err != nil {
			h.renderPage(w, "result.html", result{Message: err.Error()})
		} else {
			h.renderPage(w, "result.html", result{Message: "Done!"})
		}
	}
}

func (h *Handler) Rss(w http.ResponseWriter, r *http.Request) {
	rss, err := h.feed.getRss()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(rss))
	}
}
