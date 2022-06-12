package internal

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const fileName = "history.yml"

type history map[time.Time]string

func newHistory() (history, error) {
	if _, err := os.Stat(fileName); err != nil {
		os.Create(fileName)
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	history := make(history)
	err = yaml.Unmarshal(file, &history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (h history) add(url string, when time.Time) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	h[when] = url
	err = yaml.NewEncoder(file).Encode(h)
	if err != nil {
		return err
	}
	return nil
}
