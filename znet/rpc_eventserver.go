package znet

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/rpc"
)

var (
	rpcEventServerSubscriberCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_eventserver_subscriber_count",
		Help: "The current number of rpc subscribers",
	}, []string{})

	rpcEventServerEventCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_eventserver_event_count",
		Help: "The current number of rpc events that are subscribed",
	}, []string{})
)

type eventServer struct {
	ctx context.Context
	// the channel on which the eventMachine is listening
	eventMachineChannel chan events.Event
	eventNames          []string
	mux                 sync.Mutex
	remoteChans         []chan *rpc.Event
}

func (e *eventServer) Report() {
	rpcEventServerSubscriberCount.WithLabelValues().Set(float64(len(e.remoteChans)))
	rpcEventServerEventCount.WithLabelValues().Set(float64(len(e.eventNames)))
}

func (e *eventServer) ValidEventName(name string) bool {
	for _, n := range e.eventNames {
		if n == name {
			return true
		}
	}

	return false
}

// RegisterEvents is used to update the e.eventNames list.
func (e *eventServer) RegisterEvents(nameSet []string) {
	log.WithFields(log.Fields{
		"count": len(nameSet),
		"names": nameSet,
	}).Debug("registering event")

	if len(e.eventNames) == 0 {
		e.eventNames = make([]string, 0)
	}

	e.eventNames = append(e.eventNames, nameSet...)
}

// NoticeEvent is the call when an event should be fired.
func (e *eventServer) NoticeEvent(ctx context.Context, request *rpc.Event) (*rpc.EventResponse, error) {
	response := &rpc.EventResponse{}

	e.Report()

	for _, x := range e.remoteChans {
		x <- request
	}

	if e.ValidEventName(request.Name) {
		ev := events.Event{
			Name:    request.Name,
			Payload: request.Payload,
		}

		e.eventMachineChannel <- ev
	} else {
		response.Errors = true
		response.Message = fmt.Sprintf("unknown RPC event name: %s", request.Name)
		log.WithFields(log.Fields{
			"known_events": e.eventNames,
			"failed":       request.Name,
		}).Trace("failed to register events")
	}

	return response, nil
}

// SubscribeEvents is used to allow a caller to block while streaming events
// from the event server that match the given event names.
func (e *eventServer) SubscribeEvents(subs *rpc.EventSub, stream rpc.Events_SubscribeEventsServer) error {

	ch := e.subscriberChan()
	defer e.subscriberChanRemove(ch)

	streamContext := stream.Context()

	for {
		select {
		case <-e.ctx.Done():
			return fmt.Errorf("eventServer done")
		case <-streamContext.Done():
			return fmt.Errorf("stream done")
		case ev := <-ch:
			eventTotal.WithLabelValues(ev.Name).Inc()

			for _, eventName := range subs.EventNames {
				if matchName(ev, eventName) {
					log.Debugf("sending remote event: %s", ev.Name)
					if err := stream.Send(ev); err != nil {
						eventRemoteSendErrorTotal.WithLabelValues(ev.Name).Inc()
						log.Error(err)
					}
				}
			}
		}
	}
}

func matchName(ev *rpc.Event, name string) bool {
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

// subscriberChan creates a new channel to register with the eventServer before returning the channel.
func (e *eventServer) subscriberChan() chan *rpc.Event {
	ch := make(chan *rpc.Event)

	e.mux.Lock()
	e.remoteChans = append(e.remoteChans, ch)
	e.mux.Unlock()

	log.Tracef("subscriberChan() e.remoteChans: %+v", e.remoteChans)

	return ch
}

func (e *eventServer) subscriberChanRemove(ch chan *rpc.Event) {
	e.mux.Lock()

	log.Tracef("subscriberChanRemove() %+v from %v", ch, e.remoteChans)

	for i, q := range e.remoteChans {
		if q == ch {
			close(ch)
			log.Tracef("subscriberChanRemove channel %+v", ch)

			e.remoteChans = append(e.remoteChans[:i], e.remoteChans[i+1:]...)
		}
	}

	e.mux.Unlock()
}
