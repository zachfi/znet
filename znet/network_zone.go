package znet

import (
	"github.com/prometheus/common/log"
	ldap "gopkg.in/ldap.v2"
)

type NetworkZone struct {
	Name       string
	NTPServers []string
}

var defaultZoneAttributes = []string{
	"dn",
	"cn",
	"zoneName",
	"ntpServers",
}

func (z *Znet) GetNetworkZones(l *ldap.Conn, baseDN string) []NetworkZone {
	zones := []NetworkZone{}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=nameNetInfo)(zoneName=*))",
		defaultZoneAttributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range sr.Entries {
		// log.Warnf("Entry: %+v", e)

		z := NetworkZone{
			Name: e.DN,
		}

		for _, a := range e.Attributes {
			// log.Warnf("Attribute: %+v", a)
			// log.Warnf("ByteValues: %+v", a.ByteValues)

			switch a.Name {
			case "ntpServers":
				{
					z.NTPServers = stringValues(a)
				}
			}

		}

		zones = append(zones, z)
	}

	return zones
}
