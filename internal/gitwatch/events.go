package gitwatch

import (
	"time"
)

// EventNames are the available named events that this package exports.
var EventNames []string

func init() {
	EventNames = []string{
		"NewCommit",
		"NewTag",
	}
}

// NewCommit is the event when a new commit is detected.
type NewCommit struct {
	Time       *time.Time `json:"time"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Hash       string     `json:"hash"`
	Branch     string     `json:"branch"`
	Collection string     `json:"collection"`
}

// NewTag is the event when a new tag is detected.
type NewTag struct {
	Time       *time.Time `json:"time"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Tag        string     `json:"tag"`
	Collection string     `json:"collection"`
}
