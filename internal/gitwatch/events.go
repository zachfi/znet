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
		"NewCommit",
		"NewTag",
	}
}

// NewCommit is an event
type NewCommit struct {
	Time   *time.Time
	Name   string
	URL    string
	Hash   string
	Branch string
}

type NewTag struct {
	Time *time.Time
	Name string
	URL  string
	Tag  string
}

// Make creates the RPC request from the marshaled ExpiredTimer object.
func (c *NewCommit) Make() *pb.Event {
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

// Make creates the RPC request from the marshaled ExpiredTimer object.
func (e *NewTag) Make() *pb.Event {
	payload, err := json.Marshal(e)
	if err != nil {
		log.Error(err)
	}

	req := &pb.Event{
		Name:    reflect.TypeOf(*e).Name(),
		Payload: payload,
	}

	return req
}
