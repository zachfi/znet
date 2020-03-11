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

func (e *eventServer) RegisterEvents(nameSet ...[]string) {
	if len(e.eventNames) == 0 {
		e.eventNames = make([]string, 1)
	}

	for _, set := range nameSet {
		for _, s := range set {
			e.eventNames = append(e.eventNames, s)
		}
	}
}

func (l *eventServer) NoticeEvent(ctx context.Context, request *pb.Event) (*pb.EventResponse, error) {
	response := &pb.EventResponse{}

	if l.ValidEventName(request.Name) {
		ev := events.Event{
			Name:    request.Name,
			Payload: request.Payload,
		}

		l.ch <- ev
	} else {
		response.Errors = true
		response.Message = fmt.Sprintf("Unknown event name: %s", request.Name)
		log.Infof("payload: %s", request.Payload)
		log.Infof("known events: %+v", l.eventNames)
	}

	return response, nil
}
