package gitwatch

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/agent"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/continuous"
	"github.com/xaque208/znet/pkg/events"

	log "github.com/sirupsen/logrus"
)

// EventProducer implements events.Producer with an attached GRPC connection
// and a configuration.
type EventProducer struct {
	config      *config.GitWatchConfig
	conn        *grpc.ClientConn
	buildClient agent.BuildClient
}

// NewProducer receives a config to build a new EventProducer.
func NewProducer(cfg *config.GitWatchConfig) (events.Producer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("unable to create new gitwatch producer from nil config")
	}

	var producer events.Producer = &EventProducer{
		config: cfg,
	}

	return producer, nil
}

// Connect initializes the producer.
func (e *EventProducer) Connect(ctx context.Context, conn *grpc.ClientConn) error {
	if conn == nil {
		return fmt.Errorf("unable to connext with nil gRPC connection")
	}

	if e.conn != nil {
		log.Warnf("replacing non-nil gRPC client connection")
	}

	e.buildClient = agent.NewBuildClient(e.conn)

	e.conn = conn

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
	}(ctx, interval)

	return nil
}

func (e *EventProducer) handleRepo(ctx context.Context, repo config.GitWatchRepo, collection *string) error {
	if repo.Name == "" {
		return fmt.Errorf("repo name cannot be empty: %+v", repo)
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
			result, err := e.buildClient.BuildTag(ctx, &agent.BuildSpec{
				Tag: lastTag,
				Project: &agent.ProjectSpec{
					Name: repo.Name,
					Url:  repo.URL,
				},
			})

			if err != nil {
				return err
			}

			log.WithFields(log.Fields{
				"exit_code": result.ExitCode,
			}).Debug("result")
		}
	}

	for shortName, newHead := range newHeads {
		result, err := e.buildClient.BuildCommit(ctx, &agent.BuildSpec{
			Commit: newHead,
			Branch: shortName,
			Project: &agent.ProjectSpec{
				Name: repo.Name,
				Url:  repo.URL,
			},
		})
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"exit_code": result.ExitCode,
		}).Debug("result")
	}

	for _, r := range newTags {
		result, err := e.buildClient.BuildTag(ctx, &agent.BuildSpec{
			Tag: r,
			Project: &agent.ProjectSpec{
				Name: repo.Name,
				Url:  repo.URL,
			},
		})
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"exit_code": result.ExitCode,
		}).Debug("result")
	}

	return nil
}

func (e *EventProducer) trackRepos(ctx context.Context, repos []config.GitWatchRepo, collection *string, ticker *time.Ticker) {
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

		go func(collection config.GitWatchCollection) {
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
