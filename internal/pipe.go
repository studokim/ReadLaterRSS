package internal

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/badoux/goscraper"
	"github.com/google/uuid"
	rss "github.com/gorilla/feeds"
)

type pipe struct {
	cache      map[string]*goscraper.Document
	translator translator
}

func newPipe() pipe {
	return pipe{cache: make(map[string]*goscraper.Document), translator: newTranslator()}
}

func (p pipe) formToItem(f feed, form url.Values) (item, error) {
	switch f.FeedType {
	case urlType:
		{
			url, err := getSingle(form, "url")
			if err != nil {
				return item{}, err
			}
			text, err := getAtMostOne(form, "text")
			if err != nil {
				return item{}, err
			}
			if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
				url = "http://" + url
			}
			return item{FeedTitle: f.Title, Id: uuid.New(), Title: p.getUrlTitle(url), Created: time.Now(), Url: url, Text: p.getUrlItemText(url, text)}, nil
		}
	case textType:
		{
			title, err := getSingle(form, "title")
			if err != nil {
				return item{}, err
			}
			text, err := getSingle(form, "text")
			if err != nil {
				return item{}, err
			}
			return item{FeedTitle: f.Title, Id: uuid.New(), Title: title, Created: time.Now(), Url: "", Text: p.getTextItemText(text)}, nil
		}
	}
	return item{}, nil
}

func (p pipe) feedToRss(f feed, feedUrl string, items []item) (string, error) {
	var rssItems []*rss.Item
	for _, item := range items {
		rssItem := p.itemToRssItem(f, feedUrl, item)
		rssItems = append(rssItems, rssItem)
	}
	rssFeed := rss.Feed{
		Title:       f.Title,
		Link:        &rss.Link{Href: feedUrl},
		Description: f.Description,
		Author:      &rss.Author{Name: f.Author, Email: f.Email},
		Created:     time.Now(),
		Items:       rssItems,
	}
	return rssFeed.ToRss()
}

func (p pipe) itemToRssItem(f feed, feedUrl string, it item) *rss.Item {
	if it.Title == deleted {
		return &rss.Item{
			Id:          it.Id.String(),
			Title:       deleted,
			Description: deleted,
			Created:     it.Created,
			Link:        &rss.Link{Href: deletedurl},
			Author:      &rss.Author{Name: deleted},
		}
	}
	author := rss.Author{Name: f.Author, Email: f.Email}
	if f.FeedType == urlType {
		author = rss.Author{Name: p.getUrlAuthor(it.Url)}
	}
	return &rss.Item{
		Id:          it.Id.String(),
		Title:       it.Title,
		Description: it.Text,
		Created:     it.Created,
		Link:        &rss.Link{Href: feedUrl},
		Author:      &author,
	}
}

func (p pipe) getUrlItemText(url string, userText string) string {
	doc, err := p.getDoc(url)
	if err != nil {
		return err.Error()
	}
	if len(doc.Preview.Description) == 0 && len(userText) == 0 {
		return fmt.Sprintf("No description provided, please follow <a href=\"%s\">the link</a>.", url)
	} else if len(doc.Preview.Description) == 0 {
		return userText
	} else if len(userText) == 0 {
		return doc.Preview.Description
	}
	return fmt.Sprintf("%s<br><br>%s", userText, doc.Preview.Description)
}

func (p pipe) getTextItemText(userText string) string {
	translated, err := p.translator.translate(removeParagraphBreaks(userText))
	if err != nil {
		return fmt.Sprintf("%s<br><br><em>[%s]</em>", userText, "translator unreachable")
	}
	text := replaceDotsInDeutschDates(userText)
	textSentences := splitOnSentences(text)
	translatedSentences := splitOnSentences(translated)
	if len(textSentences) != len(translatedSentences) {
		return fmt.Sprintf("%s<br><br><strike>[%s]</strike>", text, translated)
	}
	for i := range textSentences {
		oldSentence := textSentences[i]
		newSentence := fmt.Sprintf("%s <strike>[%s]</strike>", oldSentence, translatedSentences[i])
		for _, r := range []string{".", "?", "!", "<br>"} {
			text = strings.Replace(text, oldSentence+r, newSentence+r, 1)
		}
	}
	return text
}

func (p pipe) getDoc(url string) (*goscraper.Document, error) {
	var doc *goscraper.Document
	if document, ok := p.cache[url]; ok {
		doc = document
	} else {
		document, err := goscraper.Scrape(url, 3)
		if err != nil {
			return nil, err
		}
		doc = document
	}
	p.cache[url] = doc
	return doc, nil
}

func (p pipe) getUrlTitle(url string) string {
	doc, err := p.getDoc(url)
	if err != nil {
		return err.Error()
	}
	if len(doc.Preview.Title) == 0 {
		return notitle
	}
	return doc.Preview.Title
}

func (p pipe) getUrlAuthor(url string) string {
	doc, err := p.getDoc(url)
	if err != nil {
		return err.Error()
	}
	if len(doc.Preview.Name) == 0 {
		return noauthor
	}
	return doc.Preview.Name
}

func (p pipe) itemsToExplore(items []item) []item {
	var itemsNoDeleted []item
	for _, item := range items {
		if item.Title == deleted {
			continue
		}
		item.Text = strings.ReplaceAll(item.Text, "<strike>", "<span class='blured'>")
		item.Text = strings.ReplaceAll(item.Text, "</strike>", "</span>")
		itemsNoDeleted = append(itemsNoDeleted, item)
	}
	return itemsNoDeleted
}
