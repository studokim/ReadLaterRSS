package internal

import (
	"embed"
	"net/http"
	"text/template"
)

type Handler struct {
	feed       *readLaterFeed
	templateFS embed.FS
}

type result struct {
	Message string
}

func NewHandler(templateFS embed.FS) *Handler {
	return &Handler{
		feed:       newFeed(),
		templateFS: templateFS,
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFS(h.templateFS, "templates/form.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		url := r.Form["url"][0]
		err := h.feed.addItem(url)
		var res result
		if err != nil {
			res = result{Message: err.Error()}
		} else {
			res = result{Message: "Done!"}
		}
		t, err := template.ParseFS(h.templateFS, "templates/result.html")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			t.Execute(w, res)
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
