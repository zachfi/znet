package inventory

import (
	"time"

	ldap "github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

func (i *Inventory) SetAttribute(dn, attributeName, attributeValue string, replace bool) error {

	modify := ldap.NewModifyRequest(dn, nil)
	if replace {
		modify.Replace(attributeName, []string{attributeValue})
	} else {
		modify.Add(attributeName, []string{attributeValue})
	}

	err := i.ldapClient.Modify(modify)
	if err != nil {
		return err
	}

	_ = i.UpdateTimestamp(dn)

	return nil
}

func (i *Inventory) UpdateTimestamp(dn string) error {
	now := time.Now()

	log.Infof("it is now: %+v", now.Format(time.RFC3339))

	return nil
}

func (i *Inventory) entryToHost(entry *ldap.Entry) (NetworkHost, error) {
	var host NetworkHost

	host.DN = entry.DN

	for _, a := range entry.Attributes {
		switch a.Name {
		case "cn":
			host.Name = stringValues(a)[0]
		case "netHostPlatform":
			host.Platform = stringValues(a)[0]
		case "netHostType":
			host.DeviceType = stringValues(a)[0]
		case "netHostRole":
			host.Role = stringValues(a)[0]
		case "netHostGroup":
			host.Group = stringValues(a)[0]
		case "netHostName":
			host.HostName = stringValues(a)[0]
		case "netHostDomain":
			host.Domain = stringValues(a)[0]
		case "netHostWatch":
			host.Watch = boolValues(a)[0]
		case "netHostDescription":
			host.Description = stringValues(a)[0]
		case "macAddress":
			addrs := []string{}
			addrs = append(addrs, stringValues(a)...)
			host.MACAddress = addrs
		}
	}

	return host, nil
}
