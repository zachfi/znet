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

func (e *EventProducer) sshPublicKey(repo Repo) (*ssh.PublicKeys, error) {

	// For URLs that don't start with http and when a SSHKeyPath is set, we
	// shold load the ssh key to proceeed.
	if !strings.HasPrefix(repo.URL, "https") && e.config.SSHKeyPath != "" {
		var publicKey *ssh.PublicKeys
		sshKey, _ := ioutil.ReadFile(e.config.SSHKeyPath)
		publicKey, keyError := ssh.NewPublicKeys("git", sshKey, "")
		if keyError != nil {
			return nil, fmt.Errorf("error while loading public key: %s", keyError)
		}

		return publicKey, nil
	}

	return nil, nil
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

				publicKey, err := e.sshPublicKey(repo)
				if err != nil {
					log.Error(err)
					continue
				}

				if publicKey != nil {
					cloneOpts.Auth = publicKey
				}

				_, err = os.Stat(cacheDir)
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

				newHeads, newTags, err := e.fetchRemote(r, publicKey)
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

func (e *EventProducer) fetchRemote(repo *git.Repository, sshPublicKey *ssh.PublicKeys) (map[string]string, []string, error) {
	newHeads := make(map[string]string)
	newTags := make([]string, 0)
	var err error

	beforeHeads := make(map[string]string)
	beforeTags := make(map[string]string)

	beforeRefs, err := repo.References()
	if err != nil {
		log.Error(err)
	}

	err = beforeRefs.ForEach(func(ref *plumbing.Reference) error {
		// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
		// references, the HEAD
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		// Only inspect the remote references
		if ref.Name().IsRemote() {
			beforeHeads[ref.Name().Short()] = ref.Hash().String()
		}

		if ref.Name().IsTag() {
			beforeTags[ref.Name().Short()] = ref.Hash().String()
		}

		return nil
	})

	fetchOpts := &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []gitConfig.RefSpec{
			"+refs/heads/*:refs/remotes/origin/*",
			"+refs/remotes/*:refs/remotes/origin/*",
		},
		Tags:  git.AllTags,
		Force: true,
	}

	if sshPublicKey != nil {
		fetchOpts.Auth = sshPublicKey
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		log.Error(err)
	}

	err = remote.Fetch(fetchOpts)
	if err != nil {
		if err != git.NoErrAlreadyUpToDate {
			log.Errorf("failed to fetch %s: %s", remote.Config().Name, err)
		}
	}

	afterHeads := make(map[string]string)
	afterTags := make(map[string]string)

	afterRefs, err := repo.References()
	if err != nil {
		log.Error(err)
	}

	err = afterRefs.ForEach(func(ref *plumbing.Reference) error {
		// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
		// references, the HEAD
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		// Only inspect the remote references
		if ref.Name().IsRemote() {
			afterHeads[ref.Name().Short()] = ref.Hash().String()
		}

		if ref.Name().IsTag() {
			afterTags[ref.Name().Short()] = ref.Hash().String()
		}

		return nil
	})

	nameMatch := func(refs map[string]string, shortName string) bool {
		for k := range refs {
			if k == shortName {
				return true
			}
		}

		return false
	}

	refMatch := func(refs map[string]string, shortName string, hash string) bool {
		for k, v := range refs {
			if k == shortName {
				if v == hash {
					return true
				}
			}
		}

		return false
	}

	// detect new commits on all branches
	for shortName, hash := range afterHeads {
		//detect new branches
		if !nameMatch(beforeHeads, shortName) {
			newHeads[shortName] = hash
			continue
		}

		// when before did not have this branch
		if !refMatch(beforeHeads, shortName, hash) {
			newHeads[shortName] = hash
		}
	}

	// detect new tags
	for shortName := range afterTags {
		if !nameMatch(beforeTags, shortName) {
			newTags = append(newTags, shortName)
		}
	}

	return newHeads, newTags, err
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
