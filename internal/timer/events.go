package timer

import (
	"time"
)

// EventNames are the available named events that this package exports.
var EventNames []string

func init() {
	EventNames = []string{
		"TimerExpired",
		"NamedTimer",
	}
}

// ExpiredTimer is the event for an expiration of a clock timer.
type ExpiredTimer struct {
	// events.Event
	Time *time.Time
}

// NamedTimer is the event for which you might wish to refer by name.
type NamedTimer struct {
	Name string
	Time *time.Time
}
