package gitwatch

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
		"NewCommits",
	}
}

// ExpiredTimer is the event for an expiration of a clock timer.
type NewCommits struct {
	// events.Event
	Time *time.Time
	Name string
	URL  string
	Hash string
}

// Make creates the RPC request from the marshaled ExpiredTimer object.
func (c *NewCommits) Make() *pb.Event {
	payload, err := json.Marshal(c)
	if err != nil {
		log.Error(err)
	}

	req := &pb.Event{
		Name:    reflect.TypeOf(*c).Name(),
		Payload: payload,
	}

	return req
}
