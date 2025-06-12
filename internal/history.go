package internal

import (
	"fmt"
	"os"
	"time"
)

type history struct {
	name   string
	sqlite *Sqlite
}

func newHistory(name string) (iHistory, error) {
	h := history{name: name, sqlite: &Sqlite{}}
	return &h, nil
}

func (h *history) add(r record) error {
	_, err := h.sqlite.db().Exec("insert into record(feed, title, url, text, [when]) values(?, ?, ?, ?, ?)", h.name, r.Title, r.Url, r.Text, r.When)
	return err
}

func (h *history) getSize() int {
	res := h.sqlite.db().QueryRow("select count(*) from record where feed=?", h.name)
	var count int
	err := res.Scan(&count)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("count=", count)
	return count
}

func (h *history) getRecords() []record {
	fmt.Println("getRecords where feed=", h.name)
	res, err := h.sqlite.db().Query("select title, url, text, [when] from record where feed=?", h.name)
	if err != nil {
		panic(err)
	}
	var records []record
	for res.Next() {
		var title string
		var url string
		var text string
		var whenStr string
		err = res.Scan(&title, &url, &text, &whenStr)
		if err != nil {
			fmt.Println(err, "[skip]")
		}
		when, err := time.Parse("2006-01-02T15:04:05+03:00", whenStr)
		if err != nil {
			fmt.Println(err, "[skip]", whenStr)
		}
		records = append(records, record{Title: title, Url: url, Text: text, When: when})
	}
	fmt.Println("got records: ", len(records))
	return records
}
