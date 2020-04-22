package events

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"

	pb "github.com/xaque208/znet/rpc"
)

func MakeEvent(t interface{}) *pb.Event {
	payload, err := json.Marshal(t)
	if err != nil {
		log.Error(err)
	}

	req := &pb.Event{
		Name:    reflect.TypeOf(t).Name(),
		Payload: payload,
	}

	return req
}
