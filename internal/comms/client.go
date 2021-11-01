package comms

import (
	"crypto/tls"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// StandardRPCClient implements enough to get standard gRPC client connection.
func StandardRPCClient(serverAddress string, cfg config.Config) *grpc.ClientConn {
	var opts []grpc.DialOption

	roots, err := CABundle(cfg.Vault)
	if err != nil {
		log.Error(err)
	}

	c, err := newCertify(cfg.Vault, cfg.TLS)
	if err != nil {
		log.Error(err)
	}

	serverName := strings.Split(serverAddress, ":")[0]

	tlsConfig := &tls.Config{
		ServerName:           serverName,
		InsecureSkipVerify:   false,
		RootCAs:              roots,
		GetClientCertificate: c.GetClientCertificate,
		GetCertificate:       c.GetCertificate,
	}

	// opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	// opts = append(opts, grpc.WithBlock())

	log.WithFields(log.Fields{
		"server_address": serverAddress,
	}).Debug("dialing gRPC")

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Errorf("failed dialing gRPC: %s", err)
	}

	return conn
}

func SlimRPCClient(serverAddress string) *grpc.ClientConn {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	// opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	// opts = append(opts, grpc.WithBlock())

	log.WithFields(log.Fields{
		"server_address": serverAddress,
	}).Debug("dialing gRPC")

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Errorf("failed dialing gRPC: %s", err)
	}

	return conn
}
