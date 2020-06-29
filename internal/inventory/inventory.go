package inventory

import (
	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// Inventory holds hte coniguration and clients necessary to retrieve information from data sources.
type Inventory struct {
	config     LDAPConfig
	ldapClient *ldap.Conn
}

// NewInventory returns a new Inventory object from the received config.
func NewInventory(config LDAPConfig) *Inventory {
	var inv *Inventory
	var ldapClient *ldap.Conn
	var err error

	ldapClient, err = NewLDAPClient(config)
	if err != nil {
		log.Errorf("failed LDAP connection: %s", err)
	}

	inv = &Inventory{
		config:     config,
		ldapClient: ldapClient,
	}

	return inv
}

// Close closes the LDAP client.
func (i *Inventory) Close() {
	i.ldapClient.Close()
}
