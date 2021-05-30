package comms

import (
	"crypto/tls"

	"github.com/johanbrandhorst/certify"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// RPCServerFunc is used to create a new RPC server using a received config.
type RPCServerFunc func(*config.Config) (*grpc.Server, error)

// StandardRPCServer returns a normal gRPC server.
func StandardRPCServer(cfg *config.Config) (*grpc.Server, error) {
	options := []grpc.ServerOption{}

	if cfg.Vault != nil {
		roots, err := CABundle(cfg.Vault)
		if err != nil {
			if err != ErrMissingVaultConfig {
				return nil, err
			}
		}

		if roots != nil {
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

			options = append(options, grpc.Creds(credentials.NewTLS(tlsConfig)))
		}
	}

	if len(options) > 0 {
		log.Debugf("starting with options: %+v", options)
	}

	return grpc.NewServer(options...), nil
}
