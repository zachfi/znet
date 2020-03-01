package znet

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/events"
	pb "github.com/xaque208/znet/rpc"
)

type eventServer struct {
	// chans    map[string]chan []byte
	// handlers map[string]func([]byte) error
	eventNames []string
	ch         chan events.Event
}

func (l *eventServer) ValidEventName(name string) bool {
	for _, n := range l.eventNames {
		if n == name {
			return true
		}
	}

	return false
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

	// if ch, ok := l.chans[request.Name]; ok {
	// 	ch <- request.Payload
	// if h, ok := l.handlers[request.Name]; ok {
	// 	err := h(request.Payload)
	// 	if err != nil {
	// 		response.Errors = true
	// 		response.Message = fmt.Sprintf("handler error: %s", err)
	// 	}
	// } else {
	// 	response.Errors = true
	// 	response.Message = fmt.Sprintf("Unknown event name: %s", request.Name)
	// 	log.Infof("payload: %s", request.Payload)
	// }

	return response, nil
}
