package znet

import (
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/events"
)

// EventMachine builds the channels for communicating about events  received
// from the RPC.
func (z *Znet) EventMachine(consumers []events.Consumer) error {
	log.Tracef("%d event consumers", len(consumers))

	z.EventChannel = make(chan events.Event)
	z.EventConsumers = make(map[string][]events.Handler)

	z.initEventConsumers(consumers)
	z.initEventConsumer()

	return nil
}

// initEventConsumer starts a routine that never ends to read from
// z.EventChannel and execute the loaded handlers with the event Payload.
func (z *Znet) initEventConsumer() {
	go func(ch chan events.Event) {
		log.Debugf("total %d z.EventConsumers", len(z.EventConsumers))

		for e := range ch {
			if handlers, ok := z.EventConsumers[e.Name]; ok {
				log.Debugf("executing %d handlers for event %s", len(handlers), e.Name)
				log.Tracef("listener heard event %s: %s", e.Name, string(e.Payload))
				for _, h := range handlers {
					err := h(e.Payload)
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				log.Warnf("received event with no handlers: %+v", e.Name)
			}
		}
	}(z.EventChannel)
}

// initEventConsumers updates the z.EventConsumers map.  For each received
// consumer, the handler subscriptions are determined, and appended to the
// z.EventConsumers map for execution when the named event is received.
func (z *Znet) initEventConsumers(consumers []events.Consumer) {
	for _, e := range consumers {
		subs := e.Subscriptions()
		for k, handlers := range subs {
			z.EventConsumers[k] = append(z.EventConsumers[k], handlers...)
		}
	}
}
