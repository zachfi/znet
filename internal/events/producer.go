package events

// Producer is an object that creates events.
type Producer interface {
	Produce(interface{}) error
}
