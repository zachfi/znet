package inventory

import (
	"fmt"

	ldap "github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// NetworkZone is a logical partition.
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

func (i *Inventory) NewNetworkZone(name string) (NetworkZone, error) {
	var err error
	var zone NetworkZone

	dn := fmt.Sprintf("cn=%s,ou=network,%s", name, i.config.BaseDN)

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"nameNetInfo", "top"})
	a.Attribute("cn", []string{name})
	a.Attribute("zoneName", []string{name})

	log.Tracef("creating new networkZone: %+v", a)

	err = i.ldapClient.Add(a)
	if err != nil {
		return zone, err
	}

	return i.NetworkZone(name)
}

func (i *Inventory) NetworkZone(name string) (NetworkZone, error) {
	var zone NetworkZone

	for _, z := range i.NetworkZones() {
		if z.Name == name {
			zone = z
		}
	}

	if zone.Name == "" {
		return zone, fmt.Errorf("networkZone not found: %s", name)
	}

	return zone, nil
}

// NetworkZones retrieves the NetworkHost objects from LDAP.
func (i *Inventory) NetworkZones() []NetworkZone {
	zones := []NetworkZone{}

	networkBaseDN := fmt.Sprintf("ou=network,%s", i.config.BaseDN)

	searchRequest := ldap.NewSearchRequest(
		networkBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=nameNetInfo)(zoneName=*))",
		defaultZoneAttributes,
		nil,
	)

	sr, err := i.ldapClient.Search(searchRequest)
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
			case "cn":
				z.Name = stringValues(a)[0]
			case "ntpServers":
				z.NTPServers = stringValues(a)
			}

		}

		zones = append(zones, z)
	}

	return zones
}
