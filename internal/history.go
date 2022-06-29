package internal

import (
	"gopkg.in/yaml.v3"
)

type history struct {
	name    string
	records []record
}

func newHistory(name string) (iHistory, error) {
	h := history{name: name}
	bytes, err := readFile(h.getFileName())
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &h.records)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (h *history) add(r record) error {
	file, err := openFileAppend(h.getFileName())
	if err != nil {
		return err
	}
	defer file.Close()
	h.records = append(h.records, r)
	bytes, err := yaml.Marshal([]record{r})
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err
}

func (h *history) getSize() int {
	return len(h.records)
}

func (h *history) getFileName() string {
	return h.name + "-history.yml"
}

func (h *history) getRecords() []record {
	return h.records
}
