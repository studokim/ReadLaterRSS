package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

const fileName = "history.yml"

type history []record

func newHistory() (iHistory, error) {
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
	return &history, nil
}

func (h *history) add(r record) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	*h = append(*h, r)
	bytes, err := yaml.Marshal(history{r})
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err
}
func (h *history) getSize() int {
	return len(*h)
}

func (h *history) getRecords() []record {
	return *h
}
