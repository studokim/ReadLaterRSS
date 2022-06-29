package internal

import "time"

type record struct {
	Url     string
	Context string
	When    time.Time
}
