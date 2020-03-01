package events

// Consumer is an object that
type Consumer interface {
	Subscriptions() map[string][]Handler
}
