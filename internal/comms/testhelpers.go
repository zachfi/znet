package comms

import (
	"google.golang.org/grpc"
)

func TestRPCServer() (*grpc.Server, error) {

	return grpc.NewServer(), nil
}
