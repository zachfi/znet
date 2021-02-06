package inventory

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/config"
)

// LDAPInventory holds hte coniguration and clients necessary to retrieve information from data sources.
type LDAPInventory struct {
	config *config.LDAPConfig
	conn   *ldap.Conn
	mux    sync.Mutex
}

// NewLDAPInventory returns a new Inventory object from the received config.
func NewLDAPInventory(cfg *config.LDAPConfig) (Inventory, error) {
	conn, err := NewLDAPClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed LDAP connection: %s", err)
	}

	var i Inventory = &LDAPInventory{
		config: cfg,
		conn:   conn,
	}

	return i, nil
}

// Close closes the LDAP client.
func (i *LDAPInventory) Close() {
	i.conn.Close()
}

func (i *LDAPInventory) reconnect() error {
	// Make sure old connection if definitely closed
	i.mux.Lock()
	i.conn.Close()

	// Connect to LDAP
	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", i.config.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return err
	}

	l.SetTimeout(15 * time.Second)

	// First bind with a read only user
	err = l.Bind(i.config.BindDN, i.config.BindPW)
	if err != nil {
		return err
	}

	i.conn = l
	i.mux.Unlock()
	return nil
}

func (i *LDAPInventory) SetAttribute(dn, attributeName, attributeValue string, replace bool) error {
	log.Tracef("SetAttribute: %s: %s=%s", dn, attributeName, attributeValue)

	modify := ldap.NewModifyRequest(dn, nil)
	if replace {
		modify.Replace(attributeName, []string{attributeValue})
	} else {
		modify.Add(attributeName, []string{attributeValue})
	}

	err := i.conn.Modify(modify)
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
func NewLDAPClient(cfg *config.LDAPConfig) (*ldap.Conn, error) {

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
				log.Debug("reconnecting to LDAP...")
				var err error
				l.Close()

				l, err = ldap.DialTLS(
					"tcp",
					fmt.Sprintf("%s:%d", cfg.Host, 636),
					&tls.Config{InsecureSkipVerify: true},
				)
				if err != nil {
					log.Error(err)
				}

				err = l.Bind(cfg.BindDN, cfg.BindPW)
				if err != nil {
					log.Error(err)
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

func boolValues(a *ldap.EntryAttribute) []bool {
	var values []bool

	for _, b := range a.ByteValues {
		v, err := strconv.ParseBool(string(b))
		if err != nil {
			log.Errorf("unable to parse bool: %+v", err)
		}

		values = append(values, v)
	}

	return values
}
