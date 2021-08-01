// Code generated, do not edit
package inventory

import (
	"fmt"
	"strconv"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
func (i *LDAPInventory) CreateNetworkHost(x *NetworkHost) (*NetworkHost, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil NetworkHost")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"networkHost", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new networkHost: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateNetworkHost(x)
}

// UpdateNetworkHost updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateNetworkHost(x *NetworkHost) (*NetworkHost, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil NetworkHost")
	}
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
	// TODO figure out the how we can avoid replacing the bool...
	// in case its not set on the update.
	// Replace the bool
	a.Replace("watch", []string{strconv.FormatBool(x.Watch)})
	if x.InetAddress != nil {
		a.Replace("networkHostInetAddress", x.InetAddress)
	}
	if x.Inet6Address != nil {
		a.Replace("networkHostInet6Address", x.Inet6Address)
	}
	if x.MacAddress != nil {
		a.Replace("macAddress", x.MacAddress)
	}
	if x.LastSeen != nil {
		a.Replace("networkHostLastSeen", []string{x.LastSeen.AsTime().Format(time.RFC3339)})
	}

	log.Debugf("updating networkHost: %+v", a)

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchNetworkHost(x.Name)
}

// FetchNetworkHost will retrieve a NetworkHost by name.
func (i *LDAPInventory) FetchNetworkHost(name string) (*NetworkHost, error) {

	results, err := i.ListNetworkHosts()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("networkHost not found: %s", name)
}

// ListNetworkHosts retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListNetworkHosts() ([]NetworkHost, error) {
	if i.conn == nil {
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

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := NetworkHost{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "networkHostRole":
				h.Role = stringValues(a)[0]
			case "networkHostGroup":
				h.Group = stringValues(a)[0]
			case "cn":
				h.Name = stringValues(a)[0]
			case "networkHostOperatingSystem":
				h.OperatingSystem = stringValues(a)[0]
			case "networkHostPlatform":
				h.Platform = stringValues(a)[0]
			case "networkHostType":
				h.Type = stringValues(a)[0]
			case "networkHostDomain":
				h.Domain = stringValues(a)[0]
			case "networkHostDescription":
				h.Description = stringValues(a)[0]
			case "networkHostWatch":
				v := boolValues(a)[0]
				h.Watch = v
			case "networkHostInetAddress":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.InetAddress = attrs
			case "networkHostInet6Address":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.Inet6Address = attrs
			case "macAddress":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.MacAddress = attrs
			case "networkHostLastSeen":
				attrs := []time.Time{}

				for _, v := range stringValues(a) {
					t, err := time.Parse(time.RFC3339, v)
					if err != nil {
						log.Errorf("unable to parse time: %s", err)
						continue
					}

					attrs = append(attrs, t)
				}

				h.LastSeen = timestamppb.New(attrs[0])
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
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
func (i *LDAPInventory) CreateNetworkID(x *NetworkID) (*NetworkID, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil NetworkID")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"networkId", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new networkId: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateNetworkID(x)
}

// UpdateNetworkID updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateNetworkID(x *NetworkID) (*NetworkID, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil NetworkID")
	}
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.MacAddress != nil {
		a.Replace("macAddress", x.MacAddress)
	}
	if x.IpAddress != nil {
		a.Replace("networkIdIpAddress", x.IpAddress)
	}
	if x.ReportingSource != nil {
		a.Replace("networkIdReportingSource", x.ReportingSource)
	}
	if x.ReportingSourceInterface != nil {
		a.Replace("networkIdReportingSourceInterface", x.ReportingSourceInterface)
	}
	if x.LastSeen != nil {
		a.Replace("networkIdLastSeen", []string{x.LastSeen.AsTime().Format(time.RFC3339)})
	}

	log.Debugf("updating networkId: %+v", a)

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchNetworkID(x.Name)
}

// FetchNetworkID will retrieve a NetworkID by name.
func (i *LDAPInventory) FetchNetworkID(name string) (*NetworkID, error) {

	results, err := i.ListNetworkIDs()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("networkId not found: %s", name)
}

// ListNetworkIDs retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListNetworkIDs() ([]NetworkID, error) {
	if i.conn == nil {
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

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := NetworkID{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				h.Name = stringValues(a)[0]
			case "macAddress":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.MacAddress = attrs
			case "networkIdIpAddress":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.IpAddress = attrs
			case "networkIdReportingSource":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.ReportingSource = attrs
			case "networkIdReportingSourceInterface":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.ReportingSourceInterface = attrs
			case "networkIdLastSeen":
				attrs := []time.Time{}

				for _, v := range stringValues(a) {
					t, err := time.Parse(time.RFC3339, v)
					if err != nil {
						log.Errorf("unable to parse time: %s", err)
						continue
					}

					attrs = append(attrs, t)
				}

				h.LastSeen = timestamppb.New(attrs[0])
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
}

var defaultL3NetworkAttributes = []string{
	"cn",
	"l3NetworkSoa",
	"l3NetworkDomain",
	"l3NetworkNtpServers",
	"l3NetworkInetNetwork",
	"l3NetworkInet6Network",
	"dn",
	"l3NetworkDescription",
}

// CreateL3Network creates a new LDAP entry by the received name.
func (i *LDAPInventory) CreateL3Network(x *L3Network) (*L3Network, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil L3Network")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"l3Network", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new l3Network: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateL3Network(x)
}

// UpdateL3Network updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateL3Network(x *L3Network) (*L3Network, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil L3Network")
	}
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
		a.Replace("l3NetworkNtpServers", x.NtpServers)
	}
	if x.Description != "" {
		a.Replace("l3NetworkDescription", []string{x.Description})
	}

	log.Debugf("updating l3Network: %+v", a)

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchL3Network(x.Name)
}

// FetchL3Network will retrieve a L3Network by name.
func (i *LDAPInventory) FetchL3Network(name string) (*L3Network, error) {

	results, err := i.ListL3Networks()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("l3Network not found: %s", name)
}

// ListL3Networks retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListL3Networks() ([]L3Network, error) {
	if i.conn == nil {
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

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := L3Network{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				h.Name = stringValues(a)[0]
			case "l3NetworkSoa":
				h.Soa = stringValues(a)[0]
			case "l3NetworkDomain":
				h.Domain = stringValues(a)[0]
			case "l3NetworkNtpServers":
				attrs := []string{}
				attrs = append(attrs, stringValues(a)...)
				h.NtpServers = attrs
			case "l3NetworkInetNetwork":
			case "l3NetworkInet6Network":
			case "l3NetworkDescription":
				h.Description = stringValues(a)[0]
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
}

var defaultInetNetworkAttributes = []string{
	"cn",
	"inetNetworkPrefix",
	"inetNetworkGateway",
	"inetNetworkDynamicRange",
	"dn",
}

// CreateInetNetwork creates a new LDAP entry by the received name.
func (i *LDAPInventory) CreateInetNetwork(x *InetNetwork) (*InetNetwork, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil InetNetwork")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"inetNetwork", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new inetNetwork: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateInetNetwork(x)
}

// UpdateInetNetwork updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateInetNetwork(x *InetNetwork) (*InetNetwork, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil InetNetwork")
	}
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

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchInetNetwork(x.Name)
}

// FetchInetNetwork will retrieve a InetNetwork by name.
func (i *LDAPInventory) FetchInetNetwork(name string) (*InetNetwork, error) {

	results, err := i.ListInetNetworks()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("inetNetwork not found: %s", name)
}

// ListInetNetworks retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListInetNetworks() ([]InetNetwork, error) {
	if i.conn == nil {
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

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := InetNetwork{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				h.Name = stringValues(a)[0]
			case "inetNetworkPrefix":
				h.Prefix = stringValues(a)[0]
			case "inetNetworkGateway":
				h.Gateway = stringValues(a)[0]
			case "inetNetworkDynamicRange":
				h.DynamicRange = stringValues(a)[0]
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
}

var defaultInet6NetworkAttributes = []string{
	"cn",
	"inet6NetworkPrefix",
	"inet6NetworkGateway",
	"dn",
}

// CreateInet6Network creates a new LDAP entry by the received name.
func (i *LDAPInventory) CreateInet6Network(x *Inet6Network) (*Inet6Network, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil Inet6Network")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"inet6Network", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new inet6Network: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateInet6Network(x)
}

// UpdateInet6Network updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateInet6Network(x *Inet6Network) (*Inet6Network, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil Inet6Network")
	}
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

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchInet6Network(x.Name)
}

// FetchInet6Network will retrieve a Inet6Network by name.
func (i *LDAPInventory) FetchInet6Network(name string) (*Inet6Network, error) {

	results, err := i.ListInet6Networks()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("inet6Network not found: %s", name)
}

// ListInet6Networks retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListInet6Networks() ([]Inet6Network, error) {
	if i.conn == nil {
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

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := Inet6Network{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				h.Name = stringValues(a)[0]
			case "inet6NetworkPrefix":
				h.Prefix = stringValues(a)[0]
			case "inet6NetworkGateway":
				h.Gateway = stringValues(a)[0]
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
}

var defaultZigbeeDeviceAttributes = []string{
	"cn",
	"zigbeeDeviceDescription",
	"dn",
	"zigbeeDeviceLastSeen",
	"zigbeeDeviceIotZone",
	"zigbeeDeviceType",
	"zigbeeDeviceSoftwareBuildId",
	"zigbeeDeviceDateCode",
	"zigbeeDeviceModel",
	"zigbeeDeviceVendor",
	"zigbeeDeviceManufacturerName",
	"zigbeeDevicePowerSource",
	"zigbeeDeviceModelId",
}

// CreateZigbeeDevice creates a new LDAP entry by the received name.
func (i *LDAPInventory) CreateZigbeeDevice(x *ZigbeeDevice) (*ZigbeeDevice, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil ZigbeeDevice")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"zigbeeDevice", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new zigbeeDevice: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateZigbeeDevice(x)
}

// UpdateZigbeeDevice updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateZigbeeDevice(x *ZigbeeDevice) (*ZigbeeDevice, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil ZigbeeDevice")
	}
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.Description != "" {
		a.Replace("zigbeeDeviceDescription", []string{x.Description})
	}
	if x.LastSeen != nil {
		a.Replace("zigbeeDeviceLastSeen", []string{x.LastSeen.AsTime().Format(time.RFC3339)})
	}
	if x.IotZone != "" {
		a.Replace("zigbeeDeviceIotZone", []string{x.IotZone})
	}
	if x.Type != "" {
		a.Replace("zigbeeDeviceType", []string{x.Type})
	}
	if x.SoftwareBuildId != "" {
		a.Replace("zigbeeDeviceSoftwareBuildId", []string{x.SoftwareBuildId})
	}
	if x.DateCode != "" {
		a.Replace("zigbeeDeviceDateCode", []string{x.DateCode})
	}
	if x.Model != "" {
		a.Replace("zigbeeDeviceModel", []string{x.Model})
	}
	if x.Vendor != "" {
		a.Replace("zigbeeDeviceVendor", []string{x.Vendor})
	}
	if x.ManufacturerName != "" {
		a.Replace("zigbeeDeviceManufacturerName", []string{x.ManufacturerName})
	}
	if x.PowerSource != "" {
		a.Replace("zigbeeDevicePowerSource", []string{x.PowerSource})
	}
	if x.ModelId != "" {
		a.Replace("zigbeeDeviceModelId", []string{x.ModelId})
	}

	log.Debugf("updating zigbeeDevice: %+v", a)

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchZigbeeDevice(x.Name)
}

// FetchZigbeeDevice will retrieve a ZigbeeDevice by name.
func (i *LDAPInventory) FetchZigbeeDevice(name string) (*ZigbeeDevice, error) {

	results, err := i.ListZigbeeDevices()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("zigbeeDevice not found: %s", name)
}

// ListZigbeeDevices retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListZigbeeDevices() ([]ZigbeeDevice, error) {
	if i.conn == nil {
		return nil, fmt.Errorf("unable to ListZigbeeDevices() with nil LDAP client")
	}

	var xxx []ZigbeeDevice
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=zigbeeDevice)(cn=*))",
		defaultZigbeeDeviceAttributes,
		nil,
	)

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := ZigbeeDevice{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				h.Name = stringValues(a)[0]
			case "zigbeeDeviceDescription":
				h.Description = stringValues(a)[0]
			case "zigbeeDeviceLastSeen":
				attrs := []time.Time{}

				for _, v := range stringValues(a) {
					t, err := time.Parse(time.RFC3339, v)
					if err != nil {
						log.Errorf("unable to parse time: %s", err)
						continue
					}

					attrs = append(attrs, t)
				}

				h.LastSeen = timestamppb.New(attrs[0])
			case "zigbeeDeviceIotZone":
				h.IotZone = stringValues(a)[0]
			case "zigbeeDeviceType":
				h.Type = stringValues(a)[0]
			case "zigbeeDeviceSoftwareBuildId":
				h.SoftwareBuildId = stringValues(a)[0]
			case "zigbeeDeviceDateCode":
				h.DateCode = stringValues(a)[0]
			case "zigbeeDeviceModel":
				h.Model = stringValues(a)[0]
			case "zigbeeDeviceVendor":
				h.Vendor = stringValues(a)[0]
			case "zigbeeDeviceManufacturerName":
				h.ManufacturerName = stringValues(a)[0]
			case "zigbeeDevicePowerSource":
				h.PowerSource = stringValues(a)[0]
			case "zigbeeDeviceModelId":
				h.ModelId = stringValues(a)[0]
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
}

var defaultIOTZoneAttributes = []string{
	"cn",
	"iotZoneDescription",
	"dn",
}

// CreateIOTZone creates a new LDAP entry by the received name.
func (i *LDAPInventory) CreateIOTZone(x *IOTZone) (*IOTZone, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to create nil IOTZone")
	}
	if x.Name == "" {
		return nil, fmt.Errorf("unable to create a node with no Name set")
	}

	var err error

	dn := fmt.Sprintf("cn=%s,ou=network,ou=zigbee,%s", x.Name, i.config.BaseDN)
	x.Dn = dn

	a := ldap.NewAddRequest(dn, []ldap.Control{})
	a.Attribute("objectClass", []string{"iotZone", "top"})
	a.Attribute("cn", []string{x.Name})

	log.Debugf("creating new iotZone: %+v", a)

	err = i.conn.Add(a)
	if err != nil {
		return nil, err
	}

	return i.UpdateIOTZone(x)
}

// UpdateIOTZone updates an existing LDAP entry, retrieved by name.
func (i *LDAPInventory) UpdateIOTZone(x *IOTZone) (*IOTZone, error) {
	if x == nil {
		return nil, fmt.Errorf("unable to update nil IOTZone")
	}
	if x.Dn == "" {
		return nil, fmt.Errorf("unable to update a node with no Dn set")
	}

	var err error

	a := ldap.NewModifyRequest(x.Dn, []ldap.Control{})
	if x.Description != "" {
		a.Replace("iotZoneDescription", []string{x.Description})
	}

	log.Debugf("updating iotZone: %+v", a)

	err = i.conn.Modify(a)
	if err != nil {
		return nil, err
	}

	return i.FetchIOTZone(x.Name)
}

// FetchIOTZone will retrieve a IOTZone by name.
func (i *LDAPInventory) FetchIOTZone(name string) (*IOTZone, error) {

	results, err := i.ListIOTZones()
	if err != nil {
		return nil, err
	}

	if results != nil {
		for _, x := range results {
			if x.Name == name {
				return &x, nil
			}
		}
	}

	return nil, fmt.Errorf("iotZone not found: %s", name)
}

// ListIOTZones retrieves all existing LDAP entries.
// nolint:gocyclo
func (i *LDAPInventory) ListIOTZones() ([]IOTZone, error) {
	if i.conn == nil {
		return nil, fmt.Errorf("unable to ListIOTZones() with nil LDAP client")
	}

	var xxx []IOTZone
	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=iotZone)(cn=*))",
		defaultIOTZoneAttributes,
		nil,
	)

	var searchResult *ldap.SearchResult
	attempts := 0
	for attempts < 3 {
		attempts += 1
		sr, err := i.conn.Search(searchRequest)
		if err != nil && ldap.IsErrorWithCode(err, 200) {
			log.Info("connection is closed, trying to reconnect...")
			if err = i.reconnect(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		searchResult = sr
		break
	}

	// log.Tracef("search response: %+v", searchResult)

	for _, e := range searchResult.Entries {
		h := IOTZone{Dn: e.DN}

		for _, a := range e.Attributes {
			switch a.Name {
			case "cn":
				h.Name = stringValues(a)[0]
			case "iotZoneDescription":
				h.Description = stringValues(a)[0]
			}
		}

		xxx = append(xxx, h)
	}

	return xxx, nil
}
