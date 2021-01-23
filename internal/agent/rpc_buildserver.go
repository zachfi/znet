package agent

import (
	context "context"
	sync "sync"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/config"
)

type buildServer struct {
	sync.Mutex
}

func newBuilder(cfg *config.Config) (*buildServer, error) {

	b := &buildServer{}

	return b, nil
}

// BuildProject will build a project..
func (b *buildServer) BuildProject(ctx context.Context, req *ProjectSpec) (*CommandResult, error) {
	b.Lock()
	defer b.Unlock()

	log.WithFields(log.Fields{
		"url":  req.Url,
		"name": req.Name,
	}).Debug("building project")

	// return runCommand("sudo", "puppet", "agent", "-t")
	return nil, nil
}
