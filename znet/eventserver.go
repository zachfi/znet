package znet

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/rpc"
	pb "github.com/xaque208/znet/rpc"
)

type eventServer struct {
	ch          chan events.Event
	eventNames  []string
	mux         sync.Mutex
	remoteChans []chan *rpc.Event
	quitChans   []chan bool
}

// Shutdown sents t
func (e *eventServer) Shutdown() error {
	var err error

	log.Info("eventServer shutting down")

	for _, x := range e.quitChans {
		x <- true
		close(x)
	}

	return err
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
	log.Debugf("eventServer registering %d events: %+v", len(nameSet), nameSet)

	if len(e.eventNames) == 0 {
		e.eventNames = make([]string, 1)
	}

	e.eventNames = append(e.eventNames, nameSet...)
}

// NoticeEvent is the call when an event should be fired.
func (e *eventServer) NoticeEvent(ctx context.Context, request *pb.Event) (*pb.EventResponse, error) {
	response := &pb.EventResponse{}

	for _, x := range e.remoteChans {

		// BUG!
		x <- request
	}

	if e.ValidEventName(request.Name) {
		ev := events.Event{
			Name:    request.Name,
			Payload: request.Payload,
		}

		e.ch <- ev
	} else {
		response.Errors = true
		response.Message = fmt.Sprintf("unknown RPC event name: %s", request.Name)
		log.Tracef("payload: %s", request.Payload)
		log.Tracef("known events: %+v", e.eventNames)
	}

	return response, nil
}

// SubscribeEvents is used to allow a caller to block while streaming events
// from the event server that match the given event names.
func (e *eventServer) SubscribeEvents(subs *pb.EventSub, stream pb.Events_SubscribeEventsServer) error {

	ch := e.subscriberChan()
	defer close(ch)

	quit := e.quitChan()
	defer close(quit)

	for {
		select {
		case <-quit:
			return fmt.Errorf("quitting")

		case ev := <-ch:

			match := func() bool {

				// Check for direct match first.
				for _, n := range subs.Name {
					if n == ev.Name {
						return true
					}
				}

				// Check the name of the timer, rather than the event name.
				if ev.Name == "NamedTimer" {
					var x timer.NamedTimer

					err := json.Unmarshal(ev.Payload, &x)
					if err != nil {
						log.Errorf("failed to unmarshal %T: %s", x, err)
					}

					for _, n := range subs.Name {
						if n == x.Name {
							return true
						}
					}
				}

				return false
			}()

			log.Tracef("received remote event %+v", ev)

			if match {
				log.Debugf("sending remote event: %s", ev.Name)
				if err := stream.Send(ev); err != nil {
					return err
				}
			} else {
				log.Tracef("event name not matched: %s", ev.Name)
			}
		}
	}
}

// subscriberChan creates a new channel to register with the eventServer before returning the channel.
func (e *eventServer) subscriberChan() chan *rpc.Event {
	ch := make(chan *rpc.Event)

	e.mux.Lock()
	e.remoteChans = append(e.remoteChans, ch)
	e.mux.Unlock()

	return ch
}

// quitChan creates a new channel to register with the eventServer before returning the channel.
func (e *eventServer) quitChan() chan bool {
	ch := make(chan bool)

	e.mux.Lock()
	e.quitChans = append(e.quitChans, ch)
	e.mux.Unlock()

	return ch
}

func (e *eventServer) quitChanRemove(ch chan bool) {
	e.mux.Lock()

	for i, q := range e.quitChans {
		if q == ch {
			log.Tracef("removing channel %+v", ch)
			e.quitChans = append(e.quitChans[:i], e.quitChans[i+1:]...)
		}

	}

	e.mux.Unlock()
}
