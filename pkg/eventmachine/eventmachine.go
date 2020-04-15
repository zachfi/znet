package eventmachine

import (
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/events"
)

// EventMachine
type EventMachine struct {
	// EventChannel is the channel to which the RPC eventServer writes events.
	EventChannel chan events.Event

	// EventConsumers is the map between event names and which event handlers to
	// call with the event event payload.
	EventConsumers map[string][]events.Handler
}

func Start(consumers []events.Consumer) (*EventMachine, error) {
	m := &EventMachine{}
	m.EventChannel = make(chan events.Event)
	m.EventConsumers = make(map[string][]events.Handler)

	m.initEventConsumers(consumers)
	m.initEventConsumer()

	return m, nil
}

// initEventConsumer starts a routine that never ends to read from
// z.EventChannel and execute the loaded handlers with the event Payload.
func (m *EventMachine) initEventConsumer() {
	go func(ch chan events.Event) {
		log.Debugf("total %d m.EventConsumers", len(m.EventConsumers))

		for e := range ch {
			if handlers, ok := m.EventConsumers[e.Name]; ok {
				log.Debugf("executing %d handlers for event %s", len(handlers), e.Name)
				log.Tracef("EventMachine heard event %s: %s", e.Name, string(e.Payload))
				for _, h := range handlers {
					err := h(e.Name, e.Payload)
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				log.Warnf("received event with no handlers: %+v", e.Name)
			}
		}
	}(m.EventChannel)
}

// initEventConsumers updates the z.EventConsumers map.  For each received
// consumer, the handler subscriptions are determined, and appended to the
// z.EventConsumers map for execution when the named event is received.
func (m *EventMachine) initEventConsumers(consumers []events.Consumer) {
	for _, e := range consumers {
		subs := e.Subscriptions()
		for k, handlers := range subs {
			m.EventConsumers[k] = append(m.EventConsumers[k], handlers...)
		}
	}
}
