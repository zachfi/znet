package events

// Subscriptions is a map of events to topic name for a given object.
type Subscriptions struct {
	Table map[string][]Handler
}

// NewSubscriptions creates a new Subscriptions object.
func NewSubscriptions() Subscriptions {
	s := Subscriptions{
		Table: make(map[string][]Handler),
	}

	return s
}

// Subscribe updates the table for a given event name with a given handler.
func (s *Subscriptions) Subscribe(eventName string, handler Handler) {
	s.Table[eventName] = append(s.Table[eventName], handler)
}
