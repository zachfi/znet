package znet

import (
	"crypto/tls"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewConn implements enough to get standard gRPC connection.
func NewConn(serverAddress string, config Config) *grpc.ClientConn {
	var opts []grpc.DialOption

	roots, err := CABundle(config.Vault)
	if err != nil {
		log.Error(err)
	}

	c, err := newCertify(config.Vault, config.TLS)
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
	}).Debug("dialing grpc")

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	return conn
}
