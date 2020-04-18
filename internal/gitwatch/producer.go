package gitwatch

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
			log.Error(err)
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
		req = x.Make()
	case "gitwatch.NewTag":
		x := ev.(NewTag)
		req = x.Make()
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

				cloneOpts := &git.CloneOptions{
					URL:      repo.URL,
					Progress: nil,
					// Progress: os.Stdout,
				}

				// For URLs that don't start with http and when a SSHKeyPath is set, we
				// shold load the ssh key to proceeed.
				if !strings.HasPrefix(repo.URL, "http") && e.config.SSHKeyPath != "" {
					var publicKey *ssh.PublicKeys
					sshKey, _ := ioutil.ReadFile(e.config.SSHKeyPath)
					publicKey, keyError := ssh.NewPublicKeys("git", sshKey, "")
					if keyError != nil {
						log.Errorf("error while loading public key: %s", keyError)
						continue
					}

					cloneOpts.Auth = publicKey
				}

				_, err := os.Stat(cacheDir)
				if err != nil {
					log.Infof("cloning repo %s from origin %s", repo.Name, repo.URL)
					_, cloneErr := git.PlainClone(cacheDir, false, cloneOpts)
					if cloneErr != nil {
						log.Errorf("error while cloning %s: %s", repo.URL, cloneErr)
						continue
					}
				}

				r, err := git.PlainOpen(cacheDir)
				if err != nil {
					log.Error(err)
				}

				newHead, newTags, err := updateRepo(r)
				if err != nil {
					log.Error(err)
					continue
				}

				if newHead != "" {
					ev := NewCommit{
						Name: repo.Name,
						URL:  repo.URL,
						Time: &t,
						Hash: newHead,
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

func updateRepo(repo *git.Repository) (string, []string, error) {
	var newHead string
	var newTags []string
	var err error

	remote, err := repo.Remote("origin")
	if err != nil {
		log.Error(err)
	}

	opts := &git.FetchOptions{
		// RefSpecs: []gitConfig.RefSpec{"+refs/heads/*:refs/remotes/origin/*"},
		RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Tags:     git.AllTags,
	}

	beforeHead := head(repo)
	beforeTags := tags(repo)

	err = remote.Fetch(opts)
	if err != nil {
		if err != git.NoErrAlreadyUpToDate {
			log.Errorf("failed to fetch %s: %s", remote.String(), err)
		}
	}

	w, err := repo.Worktree()
	if err != nil {
		log.Error(err)
	}

	afterHead := head(repo)
	afterTags := tags(repo)

	if beforeHead != afterHead {
		newHead = afterHead.String()

		log.Infof("resetting HEAD: %s", newHead)

		err = w.Reset(&git.ResetOptions{
			Commit: afterHead,
			Mode:   git.HardReset,
		})
		if err != nil {
			log.Error(err)
		}

		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
		})
		if err != nil {
			log.Error(err)
		}

	}

	for _, r := range afterTags {
		seen := func(n string) bool {
			for _, s := range beforeTags {
				if s == n {
					return true
				}
			}

			return false
		}(r)

		if !seen {
			newTags = append(newTags, r)
		}
	}

	return newHead, newTags, err
}

func head(repo *git.Repository) plumbing.Hash {
	head, err := repo.Head()
	if err != nil {
		log.Errorf("head error: %s", err)
		return plumbing.Hash{}
	}

	return head.Hash()
}

func tags(repo *git.Repository) []string {
	result, err := repo.Tags()
	if err != nil {
		log.Errorf("tags error: %s", err)
		return []string{}
	}

	var tags []string

	err = result.ForEach(func(repo *plumbing.Reference) error {
		x := *repo

		// log.Tracef("r.Name().Short(): %+v", x.Name().Short())

		tags = append(tags, x.Name().Short())

		return nil
	})
	if err != nil {
		log.Errorf("tags.ForEach() error: %s", err)
		return []string{}
	}

	return tags
}
