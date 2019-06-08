package znet

import (
	log "github.com/sirupsen/logrus"

	ldap "gopkg.in/ldap.v2"
)

type NetworkHost struct {
	Name       string
	HostName   string
	Domain     string
	Platform   string
	Group      string
	Role       string
	DeviceType string
	Data       HostData
}

var defaultHostAttributes = []string{
	"cn",
	"dn",
	"netHostDomain",
	"netHostGroup",
	"netHostName",
	"netHostNos",
	"netHostRole",
	"netHostType",
}

func (z *Znet) GetNetworkHosts(l *ldap.Conn, baseDN string) []NetworkHost {
	hosts := []NetworkHost{}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=netHost)(cn=*)(netHostNos=junos))",
		defaultHostAttributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range sr.Entries {
		// log.Warnf("Entry: %+v", e)

		h := NetworkHost{}

		for _, a := range e.Attributes {
			// log.Warnf("Attribute: %+v", a)
			// log.Warnf("ByteValues: %+v", a.ByteValues)
			// log.Warnf("%s stringValues: %+v", a.Name, stringValues(a))

			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "netHostNos":
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
			}
		}

		hosts = append(hosts, h)
	}

	return hosts
}
