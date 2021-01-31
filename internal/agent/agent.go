package agent

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
)

// Agent is an RPC client worker bee.
type Agent struct {
	config     *config.Config
	conn       *grpc.ClientConn
	grpcServer *grpc.Server
}

// NewAgent returns a new *Agent from the given arguments.
func NewAgent(cfg *config.Config, conn *grpc.ClientConn) (*Agent, error) {

	if cfg.TLS == nil {
		return nil, fmt.Errorf("nil TLS config")
	}

	if cfg.Vault == nil {
		return nil, fmt.Errorf("nil Vault config")
	}

	if cfg.RPC == nil {
		return nil, fmt.Errorf("nil RPC config")
	}

	if cfg.Agent == nil {
		return nil, fmt.Errorf("nil Agent config")
	}

	a := &Agent{
		config: cfg,
		conn:   conn,
	}

	if cfg.RPC != nil {
		a.grpcServer = comms.StandardRPCServer(cfg.Vault, cfg.TLS)
	}

	return a, nil
}

// Start calls start on the agent gRPC server.
func (a *Agent) Start() error {
	if a.config.RPC == nil {
		return fmt.Errorf("unable to start agent with nil RPC config")
	}

	if a.config.RPC.AgentListenAddress != "" {
		log.WithFields(log.Fields{
			"rpc_listen": a.config.RPC.AgentListenAddress,
		}).Debug("starting RPC listener")

		err := a.startRPCListener()
		if err != nil {
			return err
		}
	}

	return nil
}

// Stop calls stop on the agent gRPC server.
func (a *Agent) Stop() error {
	if a.grpcServer != nil {
		log.Debug("stopping RPC listener")
		a.grpcServer.Stop()
	}
	return nil
}

func (a *Agent) startRPCListener() error {
	info, err := readOSRelease()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id": info.ID,
	}).Debug("os-release")

	if a.config.Builder != nil {
		buildServer, err := newBuilder(a.config)
		if err != nil {
			log.Error(err)
		}
		if buildServer != nil {
			RegisterBuildServer(a.grpcServer, buildServer)
		}
	}

	switch info.ID {
	case "freebsd":
		RegisterJailHostServer(a.grpcServer, &jailServer{})
		RegisterNodeServer(a.grpcServer, &nodeServer{})
	case "arch":
		RegisterJailHostServer(a.grpcServer, &notImplementedJailServer{})
		RegisterNodeServer(a.grpcServer, &nodeServer{})
	}

	go func() {
		lis, err := net.Listen("tcp", a.config.RPC.AgentListenAddress)
		if err != nil {
			log.Errorf("rpc failed to listen: %s", err)
		}

		err = a.grpcServer.Serve(lis)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
