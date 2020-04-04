package znet

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/events"
	pb "github.com/xaque208/znet/rpc"
)

type eventServer struct {
	eventNames []string
	ch         chan events.Event
}

func (e *eventServer) ValidEventName(name string) bool {
	for _, n := range e.eventNames {
		if n == name {
			return true
		}
	}

	return false
}

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
