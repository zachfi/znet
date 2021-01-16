package agent

import (
	"context"
	"sync"

	"github.com/xaque208/znet/rpc"
)

type jailServer struct {
	sync.Mutex
}

func (j *jailServer) List(ctx context.Context, req *rpc.Empty) (*rpc.Jail, error) {
	return nil, nil
}

func (j *jailServer) Restart(ctx context.Context, req *rpc.Jail) (*rpc.CommandResult, error) {
	return nil, nil
}

func (j *jailServer) Start(ctx context.Context, req *rpc.Jail) (*rpc.CommandResult, error) {
	return nil, nil
}

func (j *jailServer) Stop(ctx context.Context, req *rpc.Jail) (*rpc.CommandResult, error) {
	return nil, nil
}
