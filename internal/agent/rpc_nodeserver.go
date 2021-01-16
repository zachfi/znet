package agent

import (
	"context"
	"os/exec"
	"sync"

	"github.com/xaque208/znet/rpc"
)

type nodeServer struct {
	sync.Mutex
}

func (n *nodeServer) RunPuppetAgent(ctx context.Context, req *rpc.Empty) (*rpc.CommandResult, error) {
	n.Lock()
	defer n.Unlock()

	cmd := exec.Command("sudo", "puppet", "agent", "-t")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
