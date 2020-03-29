package timer

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/xaque208/znet/rpc"
)

// EventNames are the available named events that this package exports.
var EventNames []string

func init() {
	EventNames = []string{
		"TimerExpired",
		"NamedTimer",
	}
}

// ExpiredTimer is the event for an expiration of a clock timer.
type ExpiredTimer struct {
	// events.Event
	Time *time.Time
}

// NamedTimer is the event for which you might wish to refer by name.
type NamedTimer struct {
	Name string
	Time *time.Time
}

// Make creates the RPC request from the marshaled ExpiredTimer object.
func (t *ExpiredTimer) Make() *pb.Event {
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

// Make creates the RPC request from the marshaled NamedTimer object.
func (t *NamedTimer) Make(name string) *pb.Event {
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
