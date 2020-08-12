package gitwatch

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/xaque208/znet/pkg/continuous"
	"github.com/xaque208/znet/pkg/events"

	log "github.com/sirupsen/logrus"
)

// EventProducer implements events.Producer with an attached GRPC connection
// and a configuration.
type EventProducer struct {
	config Config
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel func()
}

// NewProducer creates a new EventProducer to implement events.Producer and
// attach the received GRPC connection.
func NewProducer(conn *grpc.ClientConn, config Config) events.Producer {
	ctx, cancel := context.WithCancel(context.Background())

	var producer events.Producer = &EventProducer{
		conn:   conn,
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}

	return producer
}

// Start initializes the producer.
func (e *EventProducer) Start() error {
	log.Info("starting gitwatch eventProducer")

	go func(ctx context.Context) {
		err := e.watcher(ctx)
		if err != nil {
			log.Error(err)
		}
	}(e.ctx)

	return nil
}

// Stop shuts down the producer.
func (e *EventProducer) Stop() error {
	e.cancel()
	return nil
}

func (e *EventProducer) handleRepos(repos []Repo, collection *string) error {
	t := time.Now()
	for _, repo := range repos {

		if repo.Name == "" {
			log.Errorf("repo name cannot be empty: %+v", repo)
			continue
		}

		var cacheDir string
		if collection != nil {
			cacheDir = fmt.Sprintf("%s/%s/%s", e.config.CacheDir, *collection, repo.Name)
		} else {
			cacheDir = fmt.Sprintf("%s/%s", e.config.CacheDir, repo.Name)
		}

		var freshClone bool
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			freshClone = true
		}

		ci := continuous.NewCI(
			repo.URL,
			cacheDir,
			e.config.SSHKeyPath,
		)

		newHeads, newTags, err := ci.Fetch()
		if err != nil {
			log.Error(err)
			continue
		}

		// If we have a fresh clone, then
		if freshClone {
			lastTag := ci.LatestTag()

			ev := NewTag{
				Name: repo.Name,
				URL:  repo.URL,
				Time: &t,
				Tag:  lastTag,
			}

			if collection != nil {
				ev.Collection = *collection
			}

			err = events.ProduceEvent(e.conn, ev)
			if err != nil {
				log.Error(err)
			}
		}

		for shortName, newHead := range newHeads {
			ev := NewCommit{
				Name:   repo.Name,
				URL:    repo.URL,
				Time:   &t,
				Hash:   newHead,
				Branch: shortName,
			}

			err = events.ProduceEvent(e.conn, ev)
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

			if collection != nil {
				ev.Collection = *collection
			}

			err = events.ProduceEvent(e.conn, ev)
			if err != nil {
				log.Error(err)
			}
		}
	}

	return nil
}

func (e *EventProducer) watcher(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(e.config.Interval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := e.handleRepos(e.config.Repos, nil)
			if err != nil {
				log.Errorf("error handling repos: %s", err)
			}

			for _, collection := range e.config.Collections {
				err := e.handleRepos(collection.Repos, &collection.Name)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
