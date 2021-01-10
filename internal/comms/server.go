package comms

import (
	"crypto/tls"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// StandardRPCServer returns a normal gRPC server.
func StandardRPCServer(v *config.VaultConfig, t *config.TLSConfig) *grpc.Server {

	roots, err := CABundle(v)
	if err != nil {
		log.Error(err)
	}

	c, err := newCertify(v, t)
	if err != nil {
		log.Error(err)
	}

	tlsConfig := &tls.Config{
		GetCertificate: c.GetCertificate,
		ClientCAs:      roots,
		ClientAuth:     tls.RequireAndVerifyClientCert,
	}

	return grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
}
