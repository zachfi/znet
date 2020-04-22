package astro

import (
	"time"
)

// EventNames are the names of the events that this package will produce.
var EventNames = []string{
	"SolarEvent",
}

// SolarEvent is the event for which you might wish to refer by name.
type SolarEvent struct {
	Name string
	Time *time.Time
}
