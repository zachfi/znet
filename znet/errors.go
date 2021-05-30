package znet

import "errors"

// ErrNoGRPCServices is used to indicate that no grpc services are found on the server.
var ErrNoGRPCServices = errors.New("no grpc services")
