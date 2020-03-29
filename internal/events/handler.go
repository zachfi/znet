package events

// EventHandler is the interface that defines what it means to handle events.
type EventHandler interface {
	Config() (HandlerConfig, error)
	Handler() Handler
}

// Handler is responsible for decoding the event Payload and taking any necessary action.
type Handler func(Payload) error

type HandlerConfig struct {
	Name   string
	Events []string
}
