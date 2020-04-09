package events

// Producer is an object that creates events.
type Producer interface {
	// Produce is used to pass the
	Produce(interface{}) error
	Start() error
	Stop() error
}
