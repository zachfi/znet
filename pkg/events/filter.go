package events

// Filter is the interface for a given Event to implement.  This allows the
// EventServer to receive a filter for each event and avoid sending too many
// useles remote events.
type Filter interface {
	Filter(Event) bool
}
