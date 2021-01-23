package events

import (
	"context"

	"google.golang.org/grpc"
)

// Producer is an object that uses a gRPC connection to emit events to the server.
type Producer interface {
	Connect(context.Context, *grpc.ClientConn) error
}
