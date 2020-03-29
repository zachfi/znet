package events

// Handler is responsible for decoding the event Payload and taking any necessary action.
type Handler func(Payload) error

// HandlerConfig is the configuration for a single handler.
type HandlerConfig struct {
	Name   string
	Events []string
}
