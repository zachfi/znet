package znet

import (
	"context"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

type contextResolver func(context.Context) error

func peerCN(streamContext context.Context) string {

	var subscriber string
	peer, ok := peer.FromContext(streamContext)
	if ok {
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
		subscriber = tlsInfo.State.VerifiedChains[0][0].Subject.CommonName

		return subscriber
	}

	return ""
}
