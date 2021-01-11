package network

import (
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/events"
	"google.golang.org/grpc"
)

// Agent is an RPC client worker bee.
type Network struct {
	config *config.NetworkConfig
	conn   *grpc.ClientConn
}

// NewAgent returns a new *Agent from the given arguments.
func NewNetwork(cfg *config.NetworkConfig, conn *grpc.ClientConn) *Network {
	return &Network{
		config: cfg,
		conn:   conn,
	}
}

func (n *Network) Subscriptions() *events.Subscriptions {
	s := events.NewSubscriptions()

	eventNames := []string{
		"ARPUpdate",
	}

	for _, e := range eventNames {
		switch e {
		case "NamedTimer":
			s.Subscribe(e, n.namedTimerHandler)
		}
	}

	return s
}
