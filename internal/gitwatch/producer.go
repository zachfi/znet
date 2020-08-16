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
	var interval int = 600

	if e.config.Interval > 0 {
		interval = e.config.Interval
	}

	log.WithFields(log.Fields{
		"interval": interval,
	}).Info("starting gitwatch")

	go func(ctx context.Context, i int) {
		ticker := time.NewTicker(time.Duration(i) * time.Second)

		err := e.watcher(ctx, ticker)
		if err != nil {
			log.Error(err)
		}
	}(e.ctx, interval)

	return nil
}

// Stop shuts down the producer.
func (e *EventProducer) Stop() error {
	e.cancel()
	return nil
}

func (e *EventProducer) handleRepo(ctx context.Context, repo Repo, collection *string) error {
	if repo.Name == "" {
		return fmt.Errorf("repo name cannot be empty: %+v", repo)
	}

	t := time.Now()

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
		return err
	}

	if len(newHeads) > 0 {
		log.WithFields(log.Fields{
			"url":   repo.URL,
			"heads": newHeads,
		}).Debug("new heads found")
	}

	if len(newTags) > 0 {
		log.WithFields(log.Fields{
			"url":  repo.URL,
			"tags": newTags,
		}).Debug("new tags found")
	}

	// If we have a fresh clone, then
	if freshClone {
		lastTag := ci.LatestTag()

		if lastTag != "" {
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

	return nil
}

func (e *EventProducer) trackRepos(ctx context.Context, repos []Repo, collection *string, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if collection != nil {
				log.WithFields(log.Fields{
					"collection": *collection,
				}).Debug("updating collection")
			}
			for _, repo := range repos {
				err := e.handleRepo(ctx, repo, collection)
				if err != nil {
					log.WithFields(log.Fields{
						"name": repo.Name,
					}).Error("failed to fetch repo")
				}
			}
		}
	}
}

func (e *EventProducer) watcher(ctx context.Context, ticker *time.Ticker) error {
	if len(e.config.Repos) > 0 {
		go func() {
			log.WithFields(log.Fields{
				"repo_count": len(e.config.Repos),
			}).Debug("tracking repos")
			e.trackRepos(ctx, e.config.Repos, nil, ticker)
		}()
	}

	for _, collection := range e.config.Collections {
		var t *time.Ticker
		if collection.Interval > 0 {
			t = time.NewTicker(time.Duration(collection.Interval) * time.Second)
		} else {
			t = ticker
		}

		go func(collection Collection) {
			log.WithFields(log.Fields{
				"name":       collection.Name,
				"interval":   collection.Interval,
				"repo_count": len(collection.Repos),
			}).Debug("tracking collection")

			e.trackRepos(ctx, collection.Repos, &collection.Name, t)
		}(collection)
	}

	return nil
}
