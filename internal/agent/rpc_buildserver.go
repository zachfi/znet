package agent

import (
	context "context"
	"fmt"
	"io/ioutil"
	sync "sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/continuous"
)

type buildServer struct {
	UnimplementedBuildServer
	sync.Mutex
	config *config.BuilderConfig
}

func newBuilder(cfg *config.Config) (*buildServer, error) {
	b := &buildServer{
		config: cfg.Builder,
	}

	return b, nil
}

// BuildTag will build a project from a tag.
func (b *buildServer) BuildTag(ctx context.Context, req *BuildSpec) (*CommandResult, error) {
	b.Lock()
	defer b.Unlock()

	if req.Project.Name == "" {
		return nil, fmt.Errorf("unable to build tag for unnamed Project")
	}

	log.WithFields(log.Fields{
		"url":  req.Project.Url,
		"name": req.Project.Name,
		"tag":  req.Tag,
	}).Debug("building project")

	cacheDir := fmt.Sprintf("%s/%s", b.config.CacheDir, req.Project.Name)

	ci := continuous.NewCI(
		req.Project.Url,
		cacheDir,
		b.config.SSHKeyPath,
	)

	_, _, err := ci.Fetch()
	if err != nil {
		return nil, err
	}

	err = ci.CheckoutTag(req.Tag)
	if err != nil {
		return nil, err
	}

	repoConfig, err := b.loadRepoConfig(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("error loading repo config: %s", err)
	}

	for _, cmdLine := range repoConfig.OnTag {
		_, err := continuous.Build(cmdLine, cacheDir)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// BuildBranch will build a project from a tag.
func (b *buildServer) BuildCommit(ctx context.Context, req *BuildSpec) (*CommandResult, error) {
	b.Lock()
	defer b.Unlock()

	if req.Project.Name == "" {
		return nil, fmt.Errorf("unable to build tag for unnamed Project")
	}

	log.WithFields(log.Fields{
		"url":  req.Project.Url,
		"name": req.Project.Name,
		"tag":  req.Tag,
	}).Debug("building project")

	cacheDir := fmt.Sprintf("%s/%s", b.config.CacheDir, req.Project.Name)

	ci := continuous.NewCI(
		req.Project.Url,
		cacheDir,
		b.config.SSHKeyPath,
	)

	_, _, err := ci.Fetch()
	if err != nil {
		return nil, err
	}

	err = ci.CheckoutHash(req.Commit)
	if err != nil {
		log.Error(err)
	}

	repoConfig, err := b.loadRepoConfig(ci.CacheDir)
	if err != nil {
		return nil, fmt.Errorf("error loading repo config: %s", err)
	}

	branchInBranches := func() bool {
		for _, branch := range repoConfig.Branches {
			if branch == req.Branch {
				return true
			}
		}

		return false
	}

	if branchInBranches() {
		for _, cmdLine := range repoConfig.OnCommit {
			_, err := continuous.Build(cmdLine, ci.CacheDir)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (b *buildServer) loadRepoConfig(cacheDir string) (*config.RepoConfig, error) {
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
