package agent

import (
	"context"
	"fmt"
	"sync"
)

type jailServer struct {
	UnimplementedJailHostServer
	sync.Mutex
}

func (j *jailServer) List(ctx context.Context, req *Empty) (*Jails, error) {
	return nil, nil
}

func (j *jailServer) Restart(ctx context.Context, req *Jail) (*CommandResult, error) {
	return nil, nil
}

func (j *jailServer) Start(ctx context.Context, req *Jail) (*CommandResult, error) {
	return nil, nil
}

func (j *jailServer) Stop(ctx context.Context, req *Jail) (*CommandResult, error) {
	return nil, nil
}

type notImplementedJailServer struct {
	UnimplementedJailHostServer
}

func (j *notImplementedJailServer) List(ctx context.Context, req *Empty) (*Jails, error) {
	return nil, fmt.Errorf("jail server not implemented on this platform")
}

func (j *notImplementedJailServer) Restart(ctx context.Context, req *Jail) (*CommandResult, error) {
	return nil, fmt.Errorf("jail server not implemented on this platform")
}

func (j *notImplementedJailServer) Start(ctx context.Context, req *Jail) (*CommandResult, error) {
	return nil, fmt.Errorf("jail server not implemented on this platform")
}

func (j *notImplementedJailServer) Stop(ctx context.Context, req *Jail) (*CommandResult, error) {
	return nil, fmt.Errorf("jail server not implemented on this platform")
}
