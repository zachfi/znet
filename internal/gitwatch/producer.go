package gitwatch

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	git "github.com/go-git/go-git/v5"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/events"

	log "github.com/sirupsen/logrus"

	pb "github.com/xaque208/znet/rpc"
)

// EventProducer implements events.Producer with an attached GRPC connection
// and a configuration.
type EventProducer struct {
	conn    *grpc.ClientConn
	config  Config
	diechan chan bool
}

// NewProducer creates a new EventProducer to implement events.Producer and
// attach the received GRPC connection.
func NewProducer(conn *grpc.ClientConn, config Config) events.Producer {
	var producer events.Producer = &EventProducer{
		conn:   conn,
		config: config,
	}

	return producer
}

// Start initializes the producer.
func (e *EventProducer) Start() error {
	log.Info("starting gitwatch eventProducer")
	e.diechan = make(chan bool)

	go func(done chan bool) {
		err := e.watcher(done)
		if err != nil {
			log.Errorf("error e.Watcher(done): %s", err)
		}
	}(e.diechan)

	return nil
}

// Stop shuts down the producer.
func (e *EventProducer) Stop() error {
	// e.diechan <- true
	close(e.diechan)

	return nil
}

// Produce implements the events.Producer interface.  Match the supported event
// types to know which event to notice, and then send notice of the event to
// the RPC server.
func (e *EventProducer) Produce(ev interface{}) error {
	// Create the RPC client
	ec := pb.NewEventsClient(e.conn)
	t := reflect.TypeOf(ev).String()

	var req *pb.Event

	switch t {
	case "gitwatch.NewCommit":
		x := ev.(NewCommit)
		req = events.MakeEvent(x)
	case "gitwatch.NewTag":
		x := ev.(NewTag)
		req = events.MakeEvent(x)
	default:
		return fmt.Errorf("unhandled event type: %T", ev)
	}

	log.Tracef("gitwatch producing RPC event %+v", req)
	res, err := ec.NoticeEvent(context.Background(), req)
	if err != nil {
		return err
	}

	if res.Errors {
		return errors.New(res.Message)
	}

	return nil
}

func (e *EventProducer) watcher(done chan bool) error {

	ticker := time.NewTicker(time.Duration(e.config.Interval) * time.Second)

	for {
		select {
		case <-done:
			return nil
		case t := <-ticker.C:
			for _, repo := range e.config.Repos {

				if repo.Name == "" {
					log.Errorf("repo name cannot be empty: %+v", repo)
					continue
				}

				cacheDir := fmt.Sprintf("%s/%s", e.config.CacheDir, repo.Name)

				err := CacheRepo(repo.URL, cacheDir, e.config.SSHKeyPath)
				if err != nil {
					log.Error(err)
					continue
				}

				r, err := git.PlainOpen(cacheDir)
				if err != nil {
					log.Error(err)
				}

				publicKey, err := SSHPublicKey(repo.URL, e.config.SSHKeyPath)
				if err != nil {
					log.Error(err)
					continue
				}

				newHeads, newTags, err := FetchRemote(r, publicKey)
				if err != nil {
					log.Error(err)
					continue
				}

				for shortName, newHead := range newHeads {
					ev := NewCommit{
						Name:   repo.Name,
						URL:    repo.URL,
						Time:   &t,
						Hash:   newHead,
						Branch: shortName,
					}

					err = e.Produce(ev)
					if err != nil {
						log.Error(err)
					}
				}

				for _, r := range newTags {
					ev := NewTag{
						Name: repo.Name,
						URL:  repo.URL,
						Time: &t,
						Tag:  r,
					}

					err = e.Produce(ev)
					if err != nil {
						log.Error(err)
					}
				}

			}

		}
	}
}
