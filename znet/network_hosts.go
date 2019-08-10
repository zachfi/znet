package znet

import (
	log "github.com/sirupsen/logrus"

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

// GetNetworkHosts retrieves the NetworkHost objects from LDAP given an LDPA connection and baseDN.
func (z *Znet) GetNetworkHosts(l *ldap.Conn, baseDN string) ([]NetworkHost, error) {
	hosts := []NetworkHost{}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=netHost)(cn=*))",
		defaultHostAttributes,
		nil,
	)

	log.Infof("Searching LDAP with query: %s", searchRequest.Filter)

	sr, err := l.Search(searchRequest)
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
