package comms

import (
	"crypto/tls"
	"time"

	"github.com/johanbrandhorst/certify"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// StandardRPCServer returns a normal gRPC server.
func StandardRPCServer(v *config.VaultConfig, t *config.TLSConfig) *grpc.Server {

	roots, err := CABundle(v)
	if err != nil {
		log.Error(err)
	}

	var c *certify.Certify

	for {
		c, err = newCertify(v, t)
		if err != nil {
			log.Error(err)
			time.Sleep(3 * time.Second)
			continue
		}

		break
	}

	tlsConfig := &tls.Config{
		GetCertificate: c.GetCertificate,
		ClientCAs:      roots,
		ClientAuth:     tls.RequireAndVerifyClientCert,
	}

	return grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
}
