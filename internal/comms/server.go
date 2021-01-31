package comms

import (
	"crypto/tls"

	"github.com/johanbrandhorst/certify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// StandardRPCServer returns a normal gRPC server.
func StandardRPCServer(v *config.VaultConfig, t *config.TLSConfig) (*grpc.Server, error) {

	roots, err := CABundle(v)
	if err != nil {
		return nil, err
	}

	var c *certify.Certify

	c, err = newCertify(v, t)
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
