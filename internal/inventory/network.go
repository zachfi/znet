package inventory

import (
	"fmt"
	"strconv"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// CODE HERE IS GENERATED
// DO NOT EDIT

var defaultNetworkHostAttributes = []string{
	"networkHostRole",
	"networkHostGroup",
	"cn",
	"networkHostOperatingSystem",
	"networkHostPlatform",
	"networkHostType",
	"networkHostDomain",
	"networkHostDescription",
	"networkHostWatch",
	"networkHostInetAddress",
	"networkHostInet6Address",
	"macAddress",
	"networkHostLastSeen",
	"dn",
}

// CreateNetworkHost creates a new LDAP entry by the received name.
func (i *Inventory) CreateNetworkHost(x NetworkHost) (*NetworkHost, error) {
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"networkHost", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new networkHost: %+v", a)

	err = i.ldapClient.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateNetworkHost(x)
}

// UpdateNetworkHost updates an existing LDAP entry, retrieved by name.
func (i *Inventory) UpdateNetworkHost(x NetworkHost) (*NetworkHost, error) {
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.Role != "" {
		a.Replace("networkHostRole", []string{x.Role})
	}
	if x.Group != "" {
		a.Replace("networkHostGroup", []string{x.Group})
	}
	if x.OperatingSystem != "" {
		a.Replace("networkHostOperatingSystem", []string{x.OperatingSystem})
	}
	if x.Platform != "" {
		a.Replace("networkHostPlatform", []string{x.Platform})
	}
	if x.Type != "" {
		a.Replace("networkHostType", []string{x.Type})
	}
	if x.Domain != "" {
		a.Replace("networkHostDomain", []string{x.Domain})
	}
	if x.Description != "" {
		a.Replace("networkHostDescription", []string{x.Description})
	}

	if x.Watch != nil {
		a.Replace("watch", []string{strconv.FormatBool(*x.Watch)})
	}
	if x.InetAddress != nil {
		a.Replace("networkHostInetAddress", *x.InetAddress)
	}
	if x.Inet6Address != nil {
		a.Replace("networkHostInet6Address", *x.Inet6Address)
	}
	if x.MacAddress != nil {
		a.Replace("macAddress", *x.MacAddress)
	}
	if x.LastSeen != nil {
		a.Replace("networkHostLastSeen", []string{x.LastSeen.Format(time.RFC3339)})
	}

	log.Debugf("updating networkHost: %+v", a)

	err = i.ldapClient.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchNetworkHost(x.Name)
}

// FetchNetworkHost will retrieve a NetworkHost by name.
func (i *Inventory) FetchNetworkHost(name string) (*NetworkHost, error) {

	results, err := i.ListNetworkHosts()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range *results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("networkHost not found: %s", name)
}

// ListNetworkHosts retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *Inventory) ListNetworkHosts() (*[]NetworkHost, error) {
	if i.ldapClient == nil {
		return nil, fmt.Errorf("unable to ListNetworkHosts() with nil LDAP client")
	}

	var xxx []NetworkHost
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=networkHost)(cn=*))",
		defaultNetworkHostAttributes,
		nil,
	)

	// log.Tracef("searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// log.Tracef("search response: %+v", sr)

	for _, e := range sr.Entries {
		h := NetworkHost{}

		for _, a := range e.Attributes {
			switch a.Name {
			case "networkHostRole":
				{
					h.Role = stringValues(a)[0]
				}
			case "networkHostGroup":
				{
					h.Group = stringValues(a)[0]
				}
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "networkHostOperatingSystem":
				{
					h.OperatingSystem = stringValues(a)[0]
				}
			case "networkHostPlatform":
				{
					h.Platform = stringValues(a)[0]
				}
			case "networkHostType":
				{
					h.Type = stringValues(a)[0]
				}
			case "networkHostDomain":
				{
					h.Domain = stringValues(a)[0]
				}
			case "networkHostDescription":
				{
					h.Description = stringValues(a)[0]
				}
			case "networkHostWatch":
				{
					v := boolValues(a)[0]
					h.Watch = &v
				}
			case "networkHostInetAddress":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.InetAddress = &attrs
				}
			case "networkHostInet6Address":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.Inet6Address = &attrs
				}
			case "macAddress":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.MacAddress = &attrs
				}
			case "networkHostLastSeen":
				{
					attrs := []time.Time{}

					for _, v := range stringValues(a) {
						t, err := time.Parse(time.RFC3339, v)
						if err != nil {
							log.Errorf("unable to parse time: %s", err)
							continue
						}

						attrs = append(attrs, t)
					}

					h.LastSeen = &attrs[0]
				}
			case "dn":
				{
					h.Dn = stringValues(a)[0]
				}
			}
		}

		xxx = append(xxx, h)
	}

	return &xxx, nil
}

var defaultNetworkIDAttributes = []string{
	"cn",
	"macAddress",
	"networkIdIpAddress",
	"networkIdReportingSource",
	"networkIdReportingSourceInterface",
	"networkIdLastSeen",
	"dn",
}

// CreateNetworkID creates a new LDAP entry by the received name.
func (i *Inventory) CreateNetworkID(x NetworkID) (*NetworkID, error) {
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"networkId", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new networkId: %+v", a)

	err = i.ldapClient.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateNetworkID(x)
}

// UpdateNetworkID updates an existing LDAP entry, retrieved by name.
func (i *Inventory) UpdateNetworkID(x NetworkID) (*NetworkID, error) {
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.MacAddress != nil {
		a.Replace("macAddress", *x.MacAddress)
	}
	if x.IpAddress != nil {
		a.Replace("networkIdIpAddress", *x.IpAddress)
	}
	if x.ReportingSource != nil {
		a.Replace("networkIdReportingSource", *x.ReportingSource)
	}
	if x.ReportingSourceInterface != nil {
		a.Replace("networkIdReportingSourceInterface", *x.ReportingSourceInterface)
	}
	if x.LastSeen != nil {
		a.Replace("networkIdLastSeen", []string{x.LastSeen.Format(time.RFC3339)})
	}

	log.Debugf("updating networkId: %+v", a)

	err = i.ldapClient.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchNetworkID(x.Name)
}

// FetchNetworkID will retrieve a NetworkID by name.
func (i *Inventory) FetchNetworkID(name string) (*NetworkID, error) {

	results, err := i.ListNetworkIDs()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range *results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("networkId not found: %s", name)
}

// ListNetworkIDs retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *Inventory) ListNetworkIDs() (*[]NetworkID, error) {
	if i.ldapClient == nil {
		return nil, fmt.Errorf("unable to ListNetworkIDs() with nil LDAP client")
	}

	var xxx []NetworkID
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=networkId)(cn=*))",
		defaultNetworkIDAttributes,
		nil,
	)

	// log.Tracef("searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// log.Tracef("search response: %+v", sr)

	for _, e := range sr.Entries {
		h := NetworkID{}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "macAddress":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.MacAddress = &attrs
				}
			case "networkIdIpAddress":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.IpAddress = &attrs
				}
			case "networkIdReportingSource":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.ReportingSource = &attrs
				}
			case "networkIdReportingSourceInterface":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.ReportingSourceInterface = &attrs
				}
			case "networkIdLastSeen":
				{
					attrs := []time.Time{}

					for _, v := range stringValues(a) {
						t, err := time.Parse(time.RFC3339, v)
						if err != nil {
							log.Errorf("unable to parse time: %s", err)
							continue
						}

						attrs = append(attrs, t)
					}

					h.LastSeen = &attrs[0]
				}
			case "dn":
				{
					h.Dn = stringValues(a)[0]
				}
			}
		}

		xxx = append(xxx, h)
	}

	return &xxx, nil
}

var defaultL3NetworkAttributes = []string{
	"cn",
	"l3NetworkSoa",
	"l3NetworkDomain",
	"l3NetworkNtpServers",
	"l3NetworkInetNetwork",
	"l3NetworkInet6Network",
	"dn",
}

// CreateL3Network creates a new LDAP entry by the received name.
func (i *Inventory) CreateL3Network(x L3Network) (*L3Network, error) {
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"l3Network", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new l3Network: %+v", a)

	err = i.ldapClient.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateL3Network(x)
}

// UpdateL3Network updates an existing LDAP entry, retrieved by name.
func (i *Inventory) UpdateL3Network(x L3Network) (*L3Network, error) {
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.Soa != "" {
		a.Replace("l3NetworkSoa", []string{x.Soa})
	}
	if x.Domain != "" {
		a.Replace("l3NetworkDomain", []string{x.Domain})
	}
	if x.NtpServers != nil {
		a.Replace("l3NetworkNtpServers", *x.NtpServers)
	}

	log.Debugf("updating l3Network: %+v", a)

	err = i.ldapClient.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchL3Network(x.Name)
}

// FetchL3Network will retrieve a L3Network by name.
func (i *Inventory) FetchL3Network(name string) (*L3Network, error) {

	results, err := i.ListL3Networks()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range *results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("l3Network not found: %s", name)
}

// ListL3Networks retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *Inventory) ListL3Networks() (*[]L3Network, error) {
	if i.ldapClient == nil {
		return nil, fmt.Errorf("unable to ListL3Networks() with nil LDAP client")
	}

	var xxx []L3Network
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=l3Network)(cn=*))",
		defaultL3NetworkAttributes,
		nil,
	)

	// log.Tracef("searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// log.Tracef("search response: %+v", sr)

	for _, e := range sr.Entries {
		h := L3Network{}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "l3NetworkSoa":
				{
					h.Soa = stringValues(a)[0]
				}
			case "l3NetworkDomain":
				{
					h.Domain = stringValues(a)[0]
				}
			case "l3NetworkNtpServers":
				{
					attrs := []string{}
					attrs = append(attrs, stringValues(a)...)
					h.NtpServers = &attrs
				}
			case "l3NetworkInetNetwork":
				{
				}
			case "l3NetworkInet6Network":
				{
				}
			case "dn":
				{
					h.Dn = stringValues(a)[0]
				}
			}
		}

		xxx = append(xxx, h)
	}

	return &xxx, nil
}

var defaultInetNetworkAttributes = []string{
	"cn",
	"inetNetworkPrefix",
	"inetNetworkGateway",
	"inetNetworkDynamicRange",
	"dn",
}

// CreateInetNetwork creates a new LDAP entry by the received name.
func (i *Inventory) CreateInetNetwork(x InetNetwork) (*InetNetwork, error) {
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"inetNetwork", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new inetNetwork: %+v", a)

	err = i.ldapClient.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateInetNetwork(x)
}

// UpdateInetNetwork updates an existing LDAP entry, retrieved by name.
func (i *Inventory) UpdateInetNetwork(x InetNetwork) (*InetNetwork, error) {
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.Prefix != "" {
		a.Replace("inetNetworkPrefix", []string{x.Prefix})
	}
	if x.Gateway != "" {
		a.Replace("inetNetworkGateway", []string{x.Gateway})
	}
	if x.DynamicRange != "" {
		a.Replace("inetNetworkDynamicRange", []string{x.DynamicRange})
	}

	log.Debugf("updating inetNetwork: %+v", a)

	err = i.ldapClient.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchInetNetwork(x.Name)
}

// FetchInetNetwork will retrieve a InetNetwork by name.
func (i *Inventory) FetchInetNetwork(name string) (*InetNetwork, error) {

	results, err := i.ListInetNetworks()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range *results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("inetNetwork not found: %s", name)
}

// ListInetNetworks retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *Inventory) ListInetNetworks() (*[]InetNetwork, error) {
	if i.ldapClient == nil {
		return nil, fmt.Errorf("unable to ListInetNetworks() with nil LDAP client")
	}

	var xxx []InetNetwork
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=inetNetwork)(cn=*))",
		defaultInetNetworkAttributes,
		nil,
	)

	// log.Tracef("searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// log.Tracef("search response: %+v", sr)

	for _, e := range sr.Entries {
		h := InetNetwork{}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "inetNetworkPrefix":
				{
					h.Prefix = stringValues(a)[0]
				}
			case "inetNetworkGateway":
				{
					h.Gateway = stringValues(a)[0]
				}
			case "inetNetworkDynamicRange":
				{
					h.DynamicRange = stringValues(a)[0]
				}
			case "dn":
				{
					h.Dn = stringValues(a)[0]
				}
			}
		}

		xxx = append(xxx, h)
	}

	return &xxx, nil
}

var defaultInet6NetworkAttributes = []string{
	"cn",
	"inet6NetworkPrefix",
	"inet6NetworkGateway",
	"dn",
}

// CreateInet6Network creates a new LDAP entry by the received name.
func (i *Inventory) CreateInet6Network(x Inet6Network) (*Inet6Network, error) {
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"inet6Network", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new inet6Network: %+v", a)

	err = i.ldapClient.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateInet6Network(x)
}

// UpdateInet6Network updates an existing LDAP entry, retrieved by name.
func (i *Inventory) UpdateInet6Network(x Inet6Network) (*Inet6Network, error) {
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.Prefix != "" {
		a.Replace("inet6NetworkPrefix", []string{x.Prefix})
	}
	if x.Gateway != "" {
		a.Replace("inet6NetworkGateway", []string{x.Gateway})
	}

	log.Debugf("updating inet6Network: %+v", a)

	err = i.ldapClient.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchInet6Network(x.Name)
}

// FetchInet6Network will retrieve a Inet6Network by name.
func (i *Inventory) FetchInet6Network(name string) (*Inet6Network, error) {

	results, err := i.ListInet6Networks()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range *results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("inet6Network not found: %s", name)
}

// ListInet6Networks retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *Inventory) ListInet6Networks() (*[]Inet6Network, error) {
	if i.ldapClient == nil {
		return nil, fmt.Errorf("unable to ListInet6Networks() with nil LDAP client")
	}

	var xxx []Inet6Network
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=inet6Network)(cn=*))",
		defaultInet6NetworkAttributes,
		nil,
	)

	// log.Tracef("searching LDAP base %s with query: %s", i.config.BaseDN, searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// log.Tracef("search response: %+v", sr)

	for _, e := range sr.Entries {
		h := Inet6Network{}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				{
					h.Name = stringValues(a)[0]
				}
			case "inet6NetworkPrefix":
				{
					h.Prefix = stringValues(a)[0]
				}
			case "inet6NetworkGateway":
				{
					h.Gateway = stringValues(a)[0]
				}
			case "dn":
				{
					h.Dn = stringValues(a)[0]
				}
			}
		}

		xxx = append(xxx, h)
	}

	return &xxx, nil
}
