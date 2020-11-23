package inventory

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// Inventory holds hte coniguration and clients necessary to retrieve information from data sources.
type Inventory struct {
	config LDAPConfig
	conn   *ldap.Conn
	mux    sync.Mutex
}

// NewInventory returns a new Inventory object from the received config.
func NewInventory(config LDAPConfig) *Inventory {
	var err error

	conn, err := NewLDAPClient(config)
	if err != nil {
		log.Errorf("failed LDAP connection: %s", err)
	}

	return &Inventory{
		config: config,
		conn:   conn,
	}
}

func NewRPCServer(config LDAPConfig) *InventoryServer {
	inv := NewInventory(config)

	return &InventoryServer{
		inventory: *inv,
	}
}

// Close closes the LDAP client.
func (i *Inventory) Close() {
	i.conn.Close()
}

func (i *Inventory) reconnect() error {
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
