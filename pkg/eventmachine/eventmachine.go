package eventmachine

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/rpc"
)

// EventMachine is a system to facilitate receiving events and passing them along to a number of subscribers.
type EventMachine struct {
	sync.Mutex

	// eventChannel is the channel to which the RPC eventServer writes events.
	eventChannel chan *events.Event

	forwardChans []chan *events.Event

	subscriptions []*events.Subscriptions
	ctx           context.Context
	cancel        func()
}

// New creates a new EventMachine using the received consumers, complete with channels and exit.
func New(c context.Context, consumers *[]events.Consumer) (*EventMachine, error) {
	ctx, cancel := context.WithCancel(c)

	subs := []*events.Subscriptions{}
	if consumers != nil {
		for _, c := range *consumers {
			subs = append(subs, c.Subscriptions())
		}
	}

	m := &EventMachine{
		ctx:           ctx,
		cancel:        cancel,
		subscriptions: subs,
		eventChannel:  make(chan *events.Event, 1000),
	}

	m.initEventConsumer(ctx)

	return m, nil
}

// Stop closes the evnet channel.
func (m *EventMachine) Stop() error {
	log.WithFields(log.Fields{
		"event_channel": m.eventChannel,
	}).Debug("eventMachine stopping")
	m.cancel()
	return nil
}

// Send is used to marshal an object into an events.Event and write it to the
// event channel.
func (m *EventMachine) Send(t interface{}) error {
	payload, err := json.Marshal(t)
	if err != nil {
		return err
	}

	e := &events.Event{
		Name:    reflect.TypeOf(t).Name(),
		Payload: payload,
	}

	m.eventChannel <- e
	return nil
}

// forwardChan creates a new channel to register with the eventServer before returning the channel.
func (m *EventMachine) Receive() chan *events.Event {
	ch := make(chan *events.Event, 100)

	log.WithFields(log.Fields{
		"ch": ch,
	}).Trace("adding forward channel")

	m.Lock()
	m.forwardChans = append(m.forwardChans, ch)
	m.Unlock()

	return ch
}

func (m *EventMachine) ReceiveStop(ch chan *events.Event) {
	m.Lock()

	for i, q := range m.forwardChans {
		if q == ch {
			close(ch)
			log.WithFields(log.Fields{
				"ch": ch,
			}).Trace("cease forwarding")

			m.forwardChans = append(m.forwardChans[:i], m.forwardChans[i+1:]...)
		}
	}

	m.Unlock()
}

// ReadStream will forever execute the readStreamOnce to consume events from the rpc client.
func (m *EventMachine) ReadStream(client rpc.EventsClient, eventSub *rpc.EventSub) {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			err := m.readStreamOnce(m.ctx, client, eventSub)
			if err != nil {
				log.Error(err)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

// readStreamOnce will read from the RPC stream or return an error.
func (m *EventMachine) readStreamOnce(c context.Context, client rpc.EventsClient, eventSub *rpc.EventSub) error {
	var err error

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	log.WithFields(log.Fields{
		"events": eventSub.EventNames,
	}).Debug("subscribing to events")

	stream, err := client.SubscribeEvents(ctx, eventSub)
	if err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			return nil
		}

		return err
	}

	for {
		var ev *rpc.Event

		ev, err = stream.Recv()
		if err != nil {
			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				return err
			}
		}

		evE := &events.Event{
			Name:    ev.Name,
			Payload: ev.Payload,
		}

		log.WithFields(log.Fields{
			"name":    ev.Name,
			"payload": string(ev.Payload),
		}).Trace("received RPC event")
		m.eventChannel <- evE
	}
}

// initEventConsumer starts a routine that never ends to read from
// z.eventChannel and execute the loaded handlers with the event Payload.
func (m *EventMachine) initEventConsumer(ctx context.Context) {
	go func(ctx context.Context, ch chan *events.Event) {
		// log.Debugf("total %d m.EventHandlers", len(m.EventHandlers))

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-ch:
				for _, f := range m.forwardChans {
					f <- ev
				}

				for _, sub := range m.subscriptions {
					fail := 0
					if filters, ok := sub.Filters[ev.Name]; ok {
						for _, f := range filters {
							if ok := f.Filter(*ev); !ok {
								fail++
							}
						}
					}

					if fail == 0 {
						if handlers, ok := sub.Handlers[ev.Name]; ok {
							log.WithFields(log.Fields{
								"name":    ev.Name,
								"payload": string(ev.Payload),
							}).Tracef("eventmachine executing %d handler", len(handlers))

							for _, h := range handlers {
								err := h(ev.Name, ev.Payload)
								if err != nil {
									log.Error(err)
								}
							}
						}
					}
				}
			}
		}
	}(ctx, m.eventChannel)
}

func matchName(ev events.Event, name string) bool {
	// Check for direct match first.
	if name == ev.Name {
		return true
	}

	// Also check the name of the timer, rather than the event name.
	if ev.Name == "NamedTimer" {
		var x timer.NamedTimer

		err := json.Unmarshal(ev.Payload, &x)
		if err != nil {
			log.Errorf("failed to unmarshal %T: %s", x, err)
		}

		if name == x.Name {
			return true
		}
	}

	return false
}
