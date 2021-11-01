package inventory

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
)

// LDAPInventory holds hte coniguration and clients necessary to retrieve information from data sources.
type LDAPInventory struct {
	cfg *Config

	logger log.Logger

	ldapClient *ldap.Conn
	mux        sync.Mutex
}

// NewLDAPInventory returns a new Inventory object from the received config.
func NewLDAPInventory(cfg Config, logger log.Logger) (*LDAPInventory, error) {
	ldapClient, err := NewLDAPClient(cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed LDAP connection")
	}

	i := &LDAPInventory{
		cfg:        &cfg,
		logger:     logger,
		ldapClient: ldapClient,
	}

	return i, nil
}

func (i *LDAPInventory) stopping() {
	i.ldapClient.Close()
}

func (i *LDAPInventory) reconnect() error {
	// Make sure old connection if definitely closed
	i.mux.Lock()
	defer i.mux.Unlock()
	i.ldapClient.Close()

	// TODO verify TLS
	// Connect to LDAP
	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", i.cfg.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return err
	}

	l.SetTimeout(15 * time.Second)

	// First bind with a read only user
	err = l.Bind(i.cfg.BindDN, i.cfg.BindPW)
	if err != nil {
		return err
	}

	i.ldapClient = l
	return nil
}

func (i *LDAPInventory) SetAttribute(dn, attributeName, attributeValue string, replace bool) error {

	modify := ldap.NewModifyRequest(dn, nil)
	if replace {
		modify.Replace(attributeName, []string{attributeValue})
	} else {
		modify.Add(attributeName, []string{attributeValue})
	}

	err := i.ldapClient.Modify(modify)
	if err != nil {
		return err
	}

	return nil
}

func (i *LDAPInventory) UpdateTimestamp(dn string, object string) error {
	now := time.Now()

	objectName := fmt.Sprintf("%sLastSeen", object)
	return i.SetAttribute(dn, objectName, now.Format(time.RFC3339), true)
}

// NewLDAPClient constructs an LDAP client to return.
func NewLDAPClient(cfg Config, logger log.Logger) (*ldap.Conn, error) {
	logger = log.With(logger, "ldap", "client")

	if cfg.BindDN == "" || cfg.BindPW == "" || cfg.BaseDN == "" {
		return nil, fmt.Errorf("incomplete LDAP credentials, need [BindDN, BindPW, BaseDN]")
	}

	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return nil, err
	}

	l.SetTimeout(15 * time.Second)

	// First bind with a read only user
	err = l.Bind(cfg.BindDN, cfg.BindPW)
	if err != nil {
		return nil, err
	}

	// Handle reconnection
	go func() {
		t := time.NewTicker(2 * time.Second)
		for {
			<-t.C

			if l.IsClosing() {
				level.Debug(logger).Log("msg", "reconnecting to LDAP")
				var err error
				l.Close()

				l, err = ldap.DialTLS(
					"tcp",
					fmt.Sprintf("%s:%d", cfg.Host, 636),
					&tls.Config{InsecureSkipVerify: true},
				)
				if err != nil {
					level.Error(logger).Log("err", err)
				}

				err = l.Bind(cfg.BindDN, cfg.BindPW)
				if err != nil {
					level.Error(logger).Log("err", err)
				}
			}
		}
	}()

	return l, nil
}

func stringValues(a *ldap.EntryAttribute) []string {
	var values []string

	for _, b := range a.ByteValues {
		values = append(values, string(b))
	}

	return values
}

func boolValues(a *ldap.EntryAttribute, logger log.Logger) []bool {
	var values []bool

	for _, b := range a.ByteValues {
		v, err := strconv.ParseBool(string(b))
		if err != nil {
			level.Error(logger).Log("err", err)
		}

		values = append(values, v)
	}

	return values
}