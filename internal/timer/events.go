package timer

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/xaque208/znet/rpc"
	// "github.com/xaque208/znet/internal/events"
)

func init() {
	names := []string{
		"TimerExpired",
		"NamedTimer",
	}

	for _, n := range names {
		eventNames = append(eventNames, n)

	}
}

type TimerExpired struct {
	// events.Event
	Time *time.Time
}

func (t *TimerExpired) Make() *pb.Event {
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

type NamedTimer struct {
	Name string
	Time *time.Time
}

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
