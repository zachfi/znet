package events

// Handler is responsible for unmarshaling the event Payload and taking any
// necessary action.
type Handler func(string, Payload) error

// HandlerConfig is the configuration for a single handler.
type HandlerConfig struct {
	Name   string
	Events []string
}
