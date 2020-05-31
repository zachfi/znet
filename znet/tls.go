package znet

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"
	"time"

	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	log "github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
)

func newCertify(vaultConfig VaultConfig, tlsConfig TLSConfig) *certify.Certify {
	client, err := NewSecretClient(vaultConfig)
	if err != nil {
		log.Error(err)
	}

	authMethod := &vault.RenewingToken{
		Initial:     client.Token(),
		RenewBefore: 15 * time.Minute,
		TimeToLive:  24 * time.Hour,
	}

	issuer := vault.FromClient(client, "znet")
	issuer.AuthMethod = authMethod

	if tlsConfig.CAFile != "" {
		log.Debugf("loading CA file: %s", tlsConfig.CAFile)

		// The CA for vault is the Puppet CA, which is available locally.
		b, _ := ioutil.ReadFile(tlsConfig.CAFile)
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			log.Error("credentials: failed to append certificates")
		}

		issuer.TLSConfig = &tls.Config{
			RootCAs:            cp,
			InsecureSkipVerify: false,
		}
	} else {
		log.Warn("skipping TLS due to missing tlsConfig.CAFile")
	}

	cfg := certify.CertConfig{
		// SubjectAlternativeNames: []string{"localhost"},
		// IPSubjectAlternativeNames: []net.IP{
		// 	net.ParseIP("127.0.0.1"),
		// 	net.ParseIP("::1"),
		// },
		KeyGenerator: &singletonKey{},
	}

	logFormatter := log.TextFormatter{
		FullTimestamp: true,
	}

	logger := log.New()
	logger.SetLevel(log.GetLevel())
	logger.SetFormatter(&logFormatter)

	c := &certify.Certify{
		// Used when request client-side certificates and
		// added to SANs or IPSANs depending on format.
		CommonName: tlsConfig.CN,
		Issuer:     issuer,
		// It is recommended to use a cache.
		Cache:      certify.NewMemCache(),
		CertConfig: &cfg,
		// It is recommended to set RenewBefore.
		// Refresh cached certificates when < 24H left before expiry.
		RenewBefore: 24 * time.Hour,
		Logger:      logrusadapter.New(logger),

		IssueTimeout: 15 * time.Second,
	}

	return c
}

type singletonKey struct {
	key crypto.PrivateKey
	err error
	o   sync.Once
}

func (s *singletonKey) Generate() (crypto.PrivateKey, error) {
	s.o.Do(func() {
		s.key, s.err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	})

	return s.key, s.err
}
