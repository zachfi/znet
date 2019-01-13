package znet

import (
	log "github.com/sirupsen/logrus"

	ldap "gopkg.in/ldap.v2"
)

type NetworkHost struct {
	Name       string
	Platform   string
	Group      string
	Role       string
	DeviceType string
}

var defaultHostAttributes = []string{
	"cn",
	"dn",
	"netHostNos",
	"netHostRole",
	"netHostType",
}

func GetNetworkHosts(l *ldap.Conn, baseDN string) []NetworkHost {
	hosts := []NetworkHost{}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=netHost)(cn=*))",
		defaultHostAttributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range sr.Entries {
		// log.Warnf("Entry: %+v", e)

		z := NetworkHost{}

		for _, a := range e.Attributes {
			// log.Warnf("Attribute: %+v", a)
			// log.Warnf("ByteValues: %+v", a.ByteValues)
			log.Warnf("%s stringValues: %+v", a.Name, stringValues(a))

			switch a.Name {
			case "cn":
				{
					z.Name = stringValues(a)[0]
				}
			case "netHostNos":
				{
					z.Platform = stringValues(a)[0]
				}
			case "netHostType":
				{
					z.DeviceType = stringValues(a)[0]
				}
			}
		}

		hosts = append(hosts, z)
	}

	return hosts
}

func (h *NetworkHost) Configure(commit bool) {
	log.Warnf("%+v", h)

}
