package znet

// EventProducer is the event producer
type EventProducer interface {
	Name() string
	ProductNames() []string
}
