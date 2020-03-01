package znet

type EventProducer interface {
	Name() string
	ProductNames() []string
}
