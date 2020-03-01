package timer

import (
	"github.com/prometheus/common/log"
	"github.com/xaque208/znet/internal/events"
)

type EventConsumer struct {
	// events.RPCConsumer
}

func NewConsumer() events.Consumer {
	var consumer events.Consumer = &EventConsumer{}

	return consumer
}

func (e *EventConsumer) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()
	s.Subscribe("TimerExpired", eventHandler)

	return s.Table
}

func eventHandler(bits []byte) error {
	log.Warnf("Timer.eventHandler: %+v", string(bits))
	return nil
}
