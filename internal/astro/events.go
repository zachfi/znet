package astro

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/pkg/events"
)

// SolarEventNames are the names of the events that this package will produce.
var SolarEventNames = []string{
	"Sunrise",
	"Sunset",
	"PreSunset",
}

var EventNames = []string{
	"SolarEvent",
}

// SolarEvent is the event for which you might wish to refer by name.
type SolarEvent struct {
	Name string
	Time *time.Time
}

type SolarFilter struct {
	Name []string
}

func (s *SolarFilter) Filter(ev events.Event) bool {
	var solarEvent SolarEvent
	err := json.Unmarshal(ev.Payload, &solarEvent)
	if err != nil {
		log.Error(err)
		return false
	}

	for _, n := range s.Name {
		if solarEvent.Name == n {
			return true
		}
	}

	return false
}
