package events

// Producer is an object that creates events.
type Producer interface {
	// Produce is used to pass the
	Produce(interface{}) error
	// EpochCollection()
	Start() error
	Stop() error
	// scheduler() error
}
