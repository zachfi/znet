package inventory

import (
	"fmt"
	"os"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-input"
)

// UnknownHost is the information about a host when the identity of the host is
// not known.
type UnknownHost struct {
	Name       string
	IP         string
	MACAddress string
}

var unknownHostDefaultAttributes = []string{
	"cn",
	"v4Address",
	"macAddress",
}

// AdoptUnknownHost converts an UnknownHost into a known host in LDAP.
func (i *Inventory) AdoptUnknownHost(u UnknownHost, baseDN string) {
	log.Infof("Adopting host: %+v", u)

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Which CN to assign?"
	cn, err := ui.Ask(query, &input.Options{
		Default:  "newhost",
		Required: true,
		Loop:     true,
	})
	if err != nil {
		log.Error(err)
	}

	dn := fmt.Sprintf("cn=%s,%s", cn, baseDN)

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"netHost", "top"})
	a.Attribute("cn", []string{cn})
	// a.Attribute("v4Address", []string{u.IP})
	a.Attribute("macAddress", []string{u.MACAddress})

	err = i.ldapClient.Add(a)
	if err != nil {
		log.Error(err)
	}

	delDN := fmt.Sprintf("cn=%s,cn=unknown,ou=network,dc=znet", u.Name)

	d := ldap.NewDelRequest(delDN, []ldap.Control{})

	log.Infof("deleting object: %s", d)
	err = i.ldapClient.Del(d)
	if err != nil {
		log.Error(err)
	}
}
