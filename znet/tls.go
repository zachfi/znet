package znet

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	log "github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
)

func newCertify(vaultConfig *VaultConfig, tlsConfig *TLSConfig) (*certify.Certify, error) {
	if vaultConfig == nil || tlsConfig == nil {
		return nil, fmt.Errorf("unable to create new Certify with nil tlsConfig or vaultConfig")
	}

	client, err := NewSecretClient(*vaultConfig)
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

	var tlsCache certify.Cache

	if tlsConfig.CacheDir != "" {
		log.Debugf("caching tls: %s", tlsConfig.CacheDir)
		tlsCache = certify.DirCache(tlsConfig.CacheDir)
	} else {
		tlsCache = certify.NewMemCache()
	}

	c := &certify.Certify{
		// Used when request client-side certificates and
		// added to SANs or IPSANs depending on format.
		CommonName: tlsConfig.CN,
		Issuer:     issuer,
		// It is recommended to use a cache.
		Cache:      tlsCache,
		CertConfig: &cfg,
		// It is recommended to set RenewBefore.
		// Refresh cached certificates when < 24H left before expiry.
		RenewBefore: 24 * time.Hour,
		Logger:      logrusadapter.New(logger),

		IssueTimeout: 15 * time.Second,
	}

	return c, nil
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

func CABundle(vaultConfig *VaultConfig) (*x509.CertPool, error) {
	if vaultConfig == nil {
		return nil, fmt.Errorf("unable to read CABundle using nil Vault config")
	}

	// Setup the vault client to read the CA cert
	vaultClient, err := NewSecretClient(*vaultConfig)
	if err != nil {
		return nil, err
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

	return roots, nil
}
