package astro

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/xaque208/znet/rpc"
)

// EventNames are the names of the events that this package will produce.
var EventNames = []string{
	"SolarEvent",
}

// SolarEvent is the event for which you might wish to refer by name.
type SolarEvent struct {
	Name string
	Time *time.Time
}

// TODO turn this into an interface method so as to reuse this generic looking code.

// Make marshals the instance into an RPC request.
func (t *SolarEvent) Make() *pb.Event {
	payload, err := json.Marshal(t)
	if err != nil {
		log.Error(err)
	}

	req := &pb.Event{
		Name:    reflect.TypeOf(*t).Name(),
		Payload: payload,
	}

	return req
}
