package events

// Consumer is used to create a relationship between a number of handlers and
// an event name.  Handlers are expected to know how to unmarshal the event
// payload.
type Consumer interface {

	// Subscriptions is what the eventmachine calls it initializes the consumer.
	Subscriptions() *Subscriptions
}
