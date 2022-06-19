package internal

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const fileName = "history.yml"

type record struct {
	When time.Time
	Url  string
}

type history []record

func newHistory() (history, error) {
	if _, err := os.Stat(fileName); err != nil {
		os.Create(fileName)
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	history := history{}
	err = yaml.Unmarshal(file, &history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (h *history) add(url string, when time.Time) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	r := record{Url: url, When: when}
	*h = append(*h, r)
	bytes, err := yaml.Marshal([]record{r})
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err
}
