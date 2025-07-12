package internal

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	rootUrl string
	htmlFS  embed.FS
	history history
	pipe    pipe
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
		rootUrl: rootUrl,
		htmlFS:  htmlFS,
		history: history,
		pipe:    newPipe(),
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
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, pageName string, content any) {
	t, err := template.ParseFS(h.htmlFS, "html/template.html", "html/"+pageName)
	if err != nil {
		h.handle(w, err)
		return
	} else {
		selectedFeed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
			return
		}
		feeds, err := h.history.getFeeds()
		if err != nil {
			h.handle(w, err)
			return
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
			return
		}
	}
}

func (h *Handler) getSelectedFeed(r *http.Request) (feed, error) {
	titleFromUrl := r.URL.Query().Get("feed")
	if titleFromUrl != "" {
		return h.history.getFeed(titleFromUrl)
	}
	title, err := r.Cookie("feed")
	if err == nil {
		return h.history.getFeed(title.Value)
	}
	if err.Error() == "http: named cookie not present" {
		feeds, err := h.history.getFeeds()
		if err == nil && len(feeds) > 0 {
			return feeds[0], nil
		}
	}
	return feed{}, err
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, r, "index.html", nil)
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	feed, err := h.getSelectedFeed(r)
	if err != nil {
		h.handle(w, err)
		return
	}
	if r.Method == "GET" {
		switch feed.FeedType {
		case urlType:
			h.renderPage(w, r, "saveUrl.html", nil)
		case textType:
			h.renderPage(w, r, "saveText.html", nil)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			h.renderPage(w, r, "saveResult.html", result{Message: err.Error()})
			return
		}
		item, err := h.pipe.formToItem(feed, r.Form)
		if err != nil {
			h.renderPage(w, r, "saveResult.html", result{Message: err.Error()})
			return
		}
		err = h.history.addItem(item)
		if err != nil {
			h.renderPage(w, r, "saveResult.html", result{Message: err.Error()})
			return
		}
		h.renderPage(w, r, "saveResult.html", result{Message: "Done!"})
		log.Println("Saved", item.Id)
	}
}

func (h *Handler) explore(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("delete") {
		id, err := uuid.Parse(r.URL.Query().Get("delete"))
		if err != nil {
			h.handle(w, err)
			return
		}
		feed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
			return
		}
		err = h.history.deleteItem(item{FeedTitle: feed.Title, Id: id})
		if err != nil {
			h.handle(w, err)
			return
		}
		log.Println("Deleted", id)
	} else {
		feed, err := h.getSelectedFeed(r)
		if err != nil {
			h.handle(w, err)
			return
		}
		items, err := h.history.getItems(feed)
		if err != nil {
			h.handle(w, err)
			return
		}
		items = h.pipe.itemsToExplore(items)
		h.renderPage(w, r, "explore.html", items)
	}
}

func (h *Handler) rss(w http.ResponseWriter, r *http.Request) {
	feed, err := h.getSelectedFeed(r)
	if err != nil {
		h.handle(w, err)
		return
	}
	items, err := h.history.getItems(feed)
	if err != nil {
		h.handle(w, err)
		return
	}
	rss, err := h.pipe.feedToRss(feed, fmt.Sprintf("%s/rss?feed=%s", h.rootUrl, feed.Title), items)
	if err != nil {
		h.handle(w, err)
		return
	}
	w.Write([]byte(rss))
}

func (h *Handler) feeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.history.getFeeds()
	if err != nil {
		h.handle(w, err)
		return
	}
	h.renderPage(w, r, "feeds.html", feeds)
}
