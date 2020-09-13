package events

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// Subscriptions is a map of events to topic name for a given object.
type Subscriptions struct {
	sync.Mutex
	Handlers map[string][]Handler
	Filters  map[string][]Filter
}

// NewSubscriptions creates a new Subscriptions object.  This initializes the
// map of handler methods.
func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		Handlers: make(map[string][]Handler),
		Filters:  make(map[string][]Filter),
	}
}

// Subscribe updates the table for a given event name with a given handler.
func (s *Subscriptions) Subscribe(eventName string, handler Handler) {
	s.Lock()
	log.WithFields(log.Fields{
		"name": eventName,
	}).Debug("subscribing to event")

	s.Handlers[eventName] = append(s.Handlers[eventName], handler)
	s.Unlock()
}

// Filter updates the Filters table for a given event filter.
func (s *Subscriptions) Filter(eventName string, filter Filter) {
	s.Lock()
	log.WithFields(log.Fields{
		"name": eventName,
	}).Debug("filtering events")

	s.Filters[eventName] = append(s.Filters[eventName], filter)
	s.Unlock()
}
