package network

import (
	"github.com/xaque208/znet/pkg/events"
	"google.golang.org/grpc"
)

// Agent is an RPC client worker bee.
type Network struct {
	config Config
	conn   *grpc.ClientConn
}

// NewAgent returns a new *Agent from the given arguments.
func NewNetwork(config Config, conn *grpc.ClientConn) *Network {
	return &Network{
		config: config,
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
