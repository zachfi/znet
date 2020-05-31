package inventory

import (
	"fmt"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/pkg/iot"
)

// ObserveIOT receives a node when it is observed and records the state.
func (i *Inventory) ObserveIOT(node iot.Node) error {
	log.Tracef("ObserveIOT(): %+v", node)

	var zone NetworkZone
	var host NetworkHost
	var err error

	_, err = i.NetworkZone("iot")
	if err != nil {
		if strings.Contains(err.Error(), "networkZone not found") {
			zone, err = i.NewNetworkZone("iot")
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error creating networkZone: %+s", "iot"))
			}
			log.Infof("new zone created: %+v", zone)
		}

		return err
	}

	host, err = i.hostWithMAC(node.MACAddress)
	if err != nil {
		return err
	}

	if host.DN != "" {
		hasMac := func() bool {
			for _, m := range host.MACAddress {
				if strings.EqualFold(m, node.MACAddress) {
					return true
				}
			}
			return false
		}

		if hasMac() {
			log.Debugf("host %s reporting", host.DN)
		} else {
			err = i.SetAttribute(host.DN, "macAddress", node.MACAddress, false)
			if err != nil {
				return err
			}
		}
	}

	// v4Address := i.GetAttribute(host.DN, "v4Address")
	// if v4Address != node.IP {
	// 	log.Warnf("alert alert")
	// }

	return nil
}

func (i *Inventory) hostWithMAC(mac string) (NetworkHost, error) {
	var host NetworkHost
	var err error

	searchRequest := ldap.NewSearchRequest(
		i.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(macAddress=%s)", mac),
		// defaultHostAttributes,
		nil,
		nil,
	)

	searchResult, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return host, err
	}

	log.Tracef("searchResult.Entries[0]: %+v", searchResult.Entries[0])

	if len(searchResult.Entries) == 1 {
		return i.entryToHost(searchResult.Entries[0])
	}

	return host, fmt.Errorf("unhandled result count: %d", len(searchResult.Entries))
}
