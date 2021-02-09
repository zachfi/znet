package agent

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

type nodeServer struct {
	UnimplementedNodeServer
	sync.Mutex

	// ID is the os-release ID
	ID string
}

// RunPuppetAgent will perform a Puppet Run.
func (n *nodeServer) RunPuppetAgent(ctx context.Context, req *Empty) (*CommandResult, error) {
	n.Lock()
	defer n.Unlock()

	log.Debugf("running puppet agent")

	return runCommand("sudo", "puppet", "agent", "-t")
}

// PackageUpgrade will upgrade the system packages.
func (n *nodeServer) PackageUpgrade(ctx context.Context, req *Empty) (*CommandResult, error) {
	n.Lock()
	defer n.Unlock()

	log.Debugf("running package upgrade")

	switch n.ID {
	case "freebsd":
		return runCommand("sudo", "pkg", "upgrade", "-y")
	case "arch":
		return runCommand("yay", "-Syu")
	}

	return nil, fmt.Errorf("unknown n.ID: %+s", n.ID)
}
