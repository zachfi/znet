package events

// EventHandler
type EventHandler interface {
	Config() (HandlerConfig, error)
	Handler() Handler
}

// Handler is responsible for decoding some bytes and taking action.
type Handler func([]byte) error

type HandlerConfig struct {
	Name   string
	Events []string
}
