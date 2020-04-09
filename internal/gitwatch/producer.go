package gitwatch

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	git "github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/xaque208/znet/internal/events"
	"google.golang.org/grpc"

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

	// go func() {
	err := e.watcher()
	if err != nil {
		log.Error(err)
	}
	// }()

	return nil
}

// Stop shuts down the producer.
func (e *EventProducer) Stop() error {
	e.diechan <- true
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
	case "gitwatch.NewCommits":
		x := ev.(NewCommits)
		req = x.Make()
	default:
		return fmt.Errorf("unhandled event type: %T", ev)
	}

	log.Tracef("producing RPC event %+v", req)
	res, err := ec.NoticeEvent(context.Background(), req)
	if err != nil {
		return err
	}

	if res.Errors {
		return errors.New(res.Message)
	}

	return nil
}

func (e *EventProducer) watcher() error {

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-e.diechan:
			return nil
		case t := <-ticker.C:
			log.Debugf("Tick at %s", t.Format(time.RFC3339))

			for _, repo := range e.config.Repos {

				if repo.Name == "" {
					log.Error("repo name cannot be empty: %+v", repo)
					continue
				}

				cacheDir := fmt.Sprintf("%s/%s", e.config.CacheDir, repo.Name)

				_, err := os.Stat(cacheDir)
				if err != nil {
					_, err := git.PlainClone(cacheDir, false, &git.CloneOptions{
						URL:      repo.URL,
						Progress: os.Stdout,
					})
					if err != nil {
						log.Error(err)
						continue
					}
				}

				// log.Debugf("opening repo: %s", repo.Name)
				r, err := git.PlainOpen(cacheDir)
				if err != nil {
					log.Error(err)
				}

				remote, err := r.Remote("origin")
				if err != nil {
					log.Error(err)
				}

				opts := &git.FetchOptions{
					RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
					Tags:     git.AllTags,
				}

				beforeHead, err := r.Head()
				if err != nil {
					log.Errorf("head error: %s", err)
				}

				err = remote.Fetch(opts)
				if err != nil {
					if err != git.NoErrAlreadyUpToDate {
						log.Errorf("fetch error: %s", err)
					}
				}

				w, err := r.Worktree()
				if err != nil {
					log.Error(err)
				}

				afterHead, err := r.Head()
				if err != nil {
					log.Errorf("head error: %s", err)
				}

				if beforeHead.Hash() != afterHead.Hash() {
					log.Warnf("resetting to HEAD: %s", afterHead.Hash())

					err := w.Reset(&git.ResetOptions{
						Commit: afterHead.Hash(),
						Mode:   git.HardReset,
					})

					if err != nil {
						log.Error(err)
					}

					now := time.Now()

					ev := NewCommits{
						Name: repo.Name,
						URL:  repo.URL,
						Time: &now,
						Hash: afterHead.Hash().String(),
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
