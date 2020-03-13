package astro

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/xaque208/znet/rpc"
)

var EventNames []string

func init() {
	EventNames = []string{
		"PreSunset",
		"Sunrise",
		"Sunset",
	}
}

// AstroEvent is the event for which you might wish to refer by name.
type AstroEvent struct {
	Name string
	Time *time.Time
}

// TODO turn this into an interface method so as to reuse this generic looking code.
func (t *AstroEvent) Make() *pb.Event {
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
