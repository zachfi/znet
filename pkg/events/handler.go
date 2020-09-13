package events

// Handler is responsible for unmarshaling the event Payload and taking any
// necessary action.
type Handler func(string, Payload) error
