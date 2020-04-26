package builder

import (
	"encoding/json"
	"fmt"

	git "github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/events"
	"github.com/xaque208/znet/internal/gitwatch"
)

type Builder struct {
	config Config
	conn   *grpc.ClientConn
}

func NewBuilder(conn *grpc.ClientConn, config Config) *Builder {
	return &Builder{
		config: config,
		conn:   conn,
	}
}

func (b *Builder) EventNames() []string {
	var names []string

	names = append(names, gitwatch.EventNames...)
	names = append(names, "BuildTag")
	names = append(names, "BuildBranch")

	log.Debugf("builder responding to %d event names: %+v", len(names), names)

	return names
}

// Subscriptions implements the events.Consumer interface
func (b *Builder) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	for _, e := range b.EventNames() {
		switch e {
		case "NewCommit":
			s.Subscribe(e, b.checkoutCommitHandler)
		case "NewTag":
			s.Subscribe(e, b.checkoutTagHandler)
		case "BuildTag":
			s.Subscribe(e, b.checkoutTagHandler)
		// case "BuildBranch":
		// 	s.Subscribe(e, b.checkoutBranchHandler)
		default:
			log.Errorf("unhandled execution event %s", e)
		}
	}

	log.Debugf("event subscriptions %+v", s.Table)

	return s.Table
}

func (b *Builder) checkoutCommitHandler(name string, payload events.Payload) error {
	var x gitwatch.NewCommit

	err := json.Unmarshal(payload, &x)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", x, err)
	}

	if x.Name == "" {
		return fmt.Errorf("empty name value %T", x)
	}

	cacheDir := fmt.Sprintf("%s/%s", b.config.CacheDir, x.Name)

	err = gitwatch.CacheRepo(x.URL, cacheDir, b.config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("error while caching repo %s: %s", x.URL, err)
	}

	err = b.buildForEvent(x)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (b *Builder) checkoutTagHandler(name string, payload events.Payload) error {
	var x gitwatch.NewTag

	err := json.Unmarshal(payload, &x)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", x, err)
	}

	if x.Name == "" {
		return fmt.Errorf("empty name value %T", x)
	}

	cacheDir := fmt.Sprintf("%s/%s", b.config.CacheDir, x.Name)

	err = gitwatch.CacheRepo(x.URL, cacheDir, b.config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("error while caching repo %s: %s", x.URL, err)
	}

	r, err := git.PlainOpen(cacheDir)
	if err != nil {
		return err
	}

	publicKey, err := gitwatch.SSHPublicKey(x.URL, b.config.SSHKeyPath)
	if err != nil {
		return err
	}

	_, _, err = gitwatch.FetchRemote(r, publicKey)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	ref, err := r.Tag(x.Tag)
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
	})
	if err != nil {
		return err
	}

	err = b.buildForEvent(x)
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildForEvent(x interface{}) error {

	return nil
}
