package timer

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/pkg/events"
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
	Time *time.Time `json:"time,omitempty"`
}

// NamedTimer is the event for which you might wish to refer by name.
type NamedTimer struct {
	Name string     `json:"name,omitempty"`
	Time *time.Time `json:"time,omitempty"`
}

type EventFilter struct {
	Name []string
}

func (t *EventFilter) Filter(ev events.Event) bool {
	var namedTimer NamedTimer
	err := json.Unmarshal(ev.Payload, &namedTimer)
	if err != nil {
		log.Error(err)
		return false
	}

	for _, n := range t.Name {
		if namedTimer.Name == n {
			return true
		}
	}

	return false
}
