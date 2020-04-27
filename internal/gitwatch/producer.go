package gitwatch

import (
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/pkg/continuous"

	log "github.com/sirupsen/logrus"
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

					err = events.ProduceEvent(e.conn, ev)
					if err != nil {
						log.Error(err)
					}
				}

			}

		}
	}
}
