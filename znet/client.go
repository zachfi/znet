package znet

import (
	"crypto/tls"
	"crypto/x509"
	"strings"

	"github.com/hashicorp/vault/helper/certutil"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewConn implements enough to get standard gRPC connection.
func NewConn(serverAddress string, config Config) *grpc.ClientConn {
	var opts []grpc.DialOption

	// Setup the vault client to read the CA cert
	vaultClient, err := NewSecretClient(config.Vault)
	if err != nil {
		log.Error(err)
	}

	secret, err := vaultClient.Logical().Read("pki/cert/ca")
	if err != nil {
		log.Errorf("error reading ca: %v", err)
	}

	roots := x509.NewCertPool()

	parsedCertBundle, err := certutil.ParsePKIMap(secret.Data)
	if err != nil {
		log.Errorf("error parsing secret: %s", err)
	}

	roots.AddCert(parsedCertBundle.Certificate)

	c := newCertify(config.Vault, config.TLS)

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

	log.Debugf("dialing grpc: %s", serverAddress)
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Error(err)
	}

	log.Debug("returning grpc connection")

	return conn
}
