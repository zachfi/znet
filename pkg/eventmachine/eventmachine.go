package eventmachine

import (
	"context"
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/pkg/events"
)

// EventMachine is a system to facilitate receiving events and passing them along to a number of subscribers.
type EventMachine struct {
	// EventChannel is the channel to which the RPC eventServer writes events.
	EventChannel chan events.Event

	// EventConsumers is the map between event names and which event handlers to
	// call with the event event payload.
	EventConsumers map[string][]events.Handler
	ctx            context.Context
	cancel         func()
}

// New creates a new EventMachine using the received consumers, complete with channels and exit.
func New(c context.Context, consumers []events.Consumer) (*EventMachine, error) {
	ctx, cancel := context.WithCancel(c)

	m := &EventMachine{
		ctx:    ctx,
		cancel: cancel,
	}

	m.EventChannel = make(chan events.Event)
	m.EventConsumers = make(map[string][]events.Handler)

	m.initEventConsumers(consumers)
	m.initEventConsumer(ctx)

	return m, nil
}

// Stop closes the evnet channel.
func (m *EventMachine) Stop() error {
	log.WithFields(log.Fields{
		"event_channel": m.EventChannel,
	}).Debug("eventMachine stopping")
	m.cancel()
	return nil
}

// Send is used to marshal an object into an events.Event and write it to the event channel.
func (m *EventMachine) Send(t interface{}) error {
	payload, err := json.Marshal(t)
	if err != nil {
		return err
	}

	e := events.Event{
		Name:    reflect.TypeOf(t).Name(),
		Payload: payload,
	}

	m.EventChannel <- e

	return nil
}

// initEventConsumer starts a routine that never ends to read from
// z.EventChannel and execute the loaded handlers with the event Payload.
func (m *EventMachine) initEventConsumer(c context.Context) func() {
	ctx, cancel := context.WithCancel(c)

	go func(ch chan events.Event, ctx context.Context) {
		log.Debugf("total %d m.EventConsumers", len(m.EventConsumers))

		for {
			select {
			case <-ctx.Done():
				return
			case e := <-ch:
				if handlers, ok := m.EventConsumers[e.Name]; ok {
					log.Tracef("EventMachine heard event %s: %s", e.Name, string(e.Payload))
					log.Debugf("executing %d handlers for event %s", len(handlers), e.Name)
					for _, h := range handlers {
						err := h(e.Name, e.Payload)
						if err != nil {
							log.Error(err)
						}
					}
				} else {
					log.WithFields(log.Fields{
						"name":    e.Name,
						"payload": string(e.Payload),
					}).Warn("unhandled event")
				}
			}
		}
	}(m.EventChannel, ctx)

	return cancel
}

// initEventConsumers updates the m.EventConsumers map.  For each received
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
