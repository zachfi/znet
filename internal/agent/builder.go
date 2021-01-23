package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/pkg/continuous"
	"github.com/xaque208/znet/pkg/events"
)

type Builder struct {
	config *config.BuilderConfig
	conn   *grpc.ClientConn
	mux    sync.Mutex
}

func NewBuilder(conn *grpc.ClientConn, cfg *config.BuilderConfig) *Builder {
	return &Builder{
		config: cfg,
		conn:   conn,
	}
}

func (b *Builder) checkoutCommitHandler(name string, payload events.Payload) error {
	log.Debugf("locking for event: %s", name)

	b.mux.Lock()
	defer b.mux.Unlock()

	var x gitwatch.NewCommit

	err := json.Unmarshal(payload, &x)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", x, err)
	}

	if x.Name == "" {
		return fmt.Errorf("empty name value %T", x)
	}

	cacheDir := fmt.Sprintf("%s/%s", b.config.CacheDir, x.Name)

	ci := continuous.NewCI(
		x.URL,
		cacheDir,
		b.config.SSHKeyPath,
	)

	_, _, err = ci.Fetch()
	if err != nil {
		log.Error(err)
	}

	err = ci.CheckoutHash(x.Hash)
	if err != nil {
		log.Error(err)
	}

	err = b.buildForEvent(x, cacheDir)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (b *Builder) checkoutTagHandler(name string, payload events.Payload) error {
	log.Debugf("locking for event: %s", name)

	b.mux.Lock()
	defer b.mux.Unlock()

	var x gitwatch.NewTag

	err := json.Unmarshal(payload, &x)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", x, err)
	}

	if x.Name == "" {
		return fmt.Errorf("empty name value %T", x)
	}

	cacheDir := fmt.Sprintf("%s/%s", b.config.CacheDir, x.Name)

	ci := continuous.NewCI(
		x.URL,
		cacheDir,
		b.config.SSHKeyPath,
	)

	_, _, err = ci.Fetch()
	if err != nil {
		return err
	}

	err = ci.CheckoutTag(x.Tag)
	if err != nil {
		return err
	}

	err = b.buildForEvent(x, cacheDir)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (b *Builder) loadRepoConfig(cacheDir string) (*config.RepoConfig, error) {
	var repoConfig config.RepoConfig

	configPath := fmt.Sprintf("%s/.build.yaml", cacheDir)

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &repoConfig)
	if err != nil {
		return nil, err
	}

	return &repoConfig, nil
}

func (b *Builder) buildForEvent(x interface{}, cacheDir string) error {
	repoConfig, err := b.loadRepoConfig(cacheDir)
	if err != nil {
		return fmt.Errorf("error loading repo config: %s", err)
	}

	log.Debugf("building %s for %+v with config %+v", cacheDir, x, repoConfig)

	t := reflect.TypeOf(x).String()

	switch t {
	case "gitwatch.NewTag":
		for _, cmdLine := range repoConfig.OnTag {
			ev, err := continuous.Build(cmdLine, cacheDir)
			if err != nil {
				log.Error(err)
				continue
			}

			err = events.ProduceEvent(b.conn, ev)
			if err != nil {
				log.Error(err)
			}
		}

	case "gitwatch.NewCommit":
		branchInBranches := func() bool {
			for _, branch := range repoConfig.Branches {
				if branch == x.(gitwatch.NewCommit).Branch {
					return true
				}
			}

			return false
		}

		if branchInBranches() {
			for _, cmdLine := range repoConfig.OnCommit {
				ev, err := continuous.Build(cmdLine, cacheDir)
				if err != nil {
					log.Error(err)
					continue
				}

				err = events.ProduceEvent(b.conn, ev)
				if err != nil {
					log.Error(err)
				}
			}
		}

	default:
		return fmt.Errorf("nothing to build for event: %T", x)
	}

	return nil
}
