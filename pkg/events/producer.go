package events

// Producer is an object that creates events.
type Producer interface {
	Start() error
	Stop() error
}
