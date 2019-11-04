package znet

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	ldap "gopkg.in/ldap.v2"
)

// NetworkHost is a device that connects to the network.
type NetworkHost struct {
	Name        string
	HostName    string
	Domain      string
	Platform    string
	Group       string
	Role        string
	DeviceType  string
	Data        HostData
	Watch       bool
	Description string
	MACAddress  []string
	Environment map[string]string
}

var defaultHostAttributes = []string{
	"cn",
	"dn",
	"macAddress",
	"netHostDescription",
	"netHostDomain",
	"netHostGroup",
	"netHostName",
	"netHostPlatform",
	"netHostRole",
	"netHostType",
	"netHostWatch",
}

func (z *Znet) RecordUnknownHost(l *ldap.Conn, baseDN string, address string, mac string) error {

	cn := strings.Replace(mac, ":", "", -1)

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=unknownNetHost)(cn=%s))", cn),
		[]string{"cn"},
		nil,
	)

	log.Infof("Searching LDAP with query: %s", searchRequest.Filter)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) > 0 {
		log.Debugf("Host mac %s is already unknown", mac)
		return nil
	}

	log.Debugf("Recording unknown host %s", mac)

	dn := fmt.Sprintf("cn=%s,%s", cn, baseDN)

	a := ldap.NewAddRequest(dn)
	a.Attribute("objectClass", []string{"unknownNetHost", "top"})
	a.Attribute("cn", []string{cn})
	a.Attribute("v4Address", []string{address})
	a.Attribute("macAddress", []string{mac})
	err = l.Add(a)
	if err != nil {
		log.Errorf("%+v", a)
		return err
	}

	return nil
}

// NetworkHosts retrieves the NetworkHost objects from LDAP given an LDPA connection and baseDN.
func (z *Znet) NetworkHosts() ([]NetworkHost, error) {
	hosts := []NetworkHost{}

	searchRequest := ldap.NewSearchRequest(
		z.Config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=netHost)(cn=*))",
		defaultHostAttributes,
		nil,
	)

	log.Infof("Searching LDAP base %s with query: %s", z.Config.LDAP.BaseDN, searchRequest.Filter)

	sr, err := z.ldapClient.Search(searchRequest)
	if err != nil {
		return []NetworkHost{}, err
	}

	for _, e := range sr.Entries {
		h := NetworkHost{}
		h.Environment = z.Environment

		for _, a := range e.Attributes {

			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "netHostPlatform":
				{
					h.Platform = stringValues(a)[0]
				}
			case "netHostType":
				{
					h.DeviceType = stringValues(a)[0]
				}
			case "netHostRole":
				{
					h.Role = stringValues(a)[0]
				}
			case "netHostGroup":
				{
					h.Group = stringValues(a)[0]
				}
			case "netHostName":
				{
					h.HostName = stringValues(a)[0]
				}
			case "netHostDomain":
				{
					h.Domain = stringValues(a)[0]
				}
			case "netHostWatch":
				{
					h.Watch = boolValues(a)[0]
				}
			case "netHostDescription":
				{
					h.Description = stringValues(a)[0]
				}
			case "macAddress":
				{
					addrs := []string{}
					addrs = append(addrs, stringValues(a)...)
					h.MACAddress = addrs
				}
			}
		}

		hosts = append(hosts, h)
	}

	return hosts, nil
}
