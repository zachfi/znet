package comms

import (
	"crypto/tls"

	"github.com/johanbrandhorst/certify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// RPCServerFunc is used to create a new RPC server using a received config.
type RPCServerFunc func(*config.Config) (*grpc.Server, error)

// StandardRPCServer returns a normal gRPC server.
func StandardRPCServer(cfg *config.Config) (*grpc.Server, error) {
	roots, err := CABundle(cfg.Vault)
	if err != nil {
		return nil, err
	}

	var c *certify.Certify

	c, err = newCertify(cfg.Vault, cfg.TLS)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		GetCertificate: c.GetCertificate,
		ClientCAs:      roots,
		ClientAuth:     tls.RequireAndVerifyClientCert,
	}

	return grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig))), err
}
