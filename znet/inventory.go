package znet

import (
	"fmt"

	ldap "github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// Inventory holds hte coniguration and clients necessary to retrieve information from data sources.
type Inventory struct {
	config     LDAPConfig
	ldapClient *ldap.Conn
}

// NetworkHosts retrieves the NetworkHost objects from LDAP given an LDPA connection and baseDN.
func (i *Inventory) NetworkHosts() ([]NetworkHost, error) {

	if i.ldapClient == nil {
		return []NetworkHost{}, fmt.Errorf("unable to use nil LDAP client")
	}

	hosts := []NetworkHost{}

	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=netHost)(cn=*))",
		defaultHostAttributes,
		nil,
	)

	log.Debugf("Searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	if i.ldapClient.IsClosing() {
		log.Warnf("LDAP connection is closing")
	}

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return []NetworkHost{}, err
	}

	for _, e := range sr.Entries {
		h := NetworkHost{}

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

// UnknownHosts retrieves unknownHost objects from LDAP.
func (i *Inventory) UnknownHosts() ([]UnknownHost, error) {

	if i.ldapClient == nil {
		return []UnknownHost{}, fmt.Errorf("unable to use nil LDAP client")
	}

	hosts := []UnknownHost{}

	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=unknownNetHost)(cn=*))",
		unknownHostDefaultAttributes,
		nil,
	)

	log.Debugf("Searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return []UnknownHost{}, err
	}

	for _, e := range sr.Entries {
		h := UnknownHost{}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "v4Address":
				{
					h.IP = stringValues(a)[0]
				}
			case "macAddress":
				{
					addrs := []string{}
					addrs = append(addrs, stringValues(a)...)
					h.MACAddress = addrs[0]
				}
			}
		}

		hosts = append(hosts, h)
	}

	return hosts, nil
}
