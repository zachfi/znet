package events

// Producer is an object that produces events.
type Producer interface {
	Produce(interface{}) error
}
