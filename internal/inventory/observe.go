package inventory

import (
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/pkg/iot"
)

// ObserveIOT receives a node when it is observed and records the state.
func (i *Inventory) ObserveIOT(node iot.Node) error {
	log.Tracef("ObserveIOT(): %+v", node)

	// var zone NetworkZone
	// var host NetworkHost
	// var err error
	//
	// _, err = i.NetworkZone("iot")
	// if err != nil {
	// 	if strings.Contains(err.Error(), "networkZone not found") {
	// 		zone, err = i.NewNetworkZone("iot")
	// 		if err != nil {
	// 			return errors.Wrap(err, fmt.Sprintf("error creating networkZone: %+s", "iot"))
	// 		}
	// 		log.Infof("new zone created: %+v", zone)
	// 	}
	//
	// 	return err
	// }
	//
	// host, err = i.hostWithMAC(node.MACAddress)
	// if err != nil {
	// 	return err
	// }
	//
	// if host.DN != "" {
	// 	hasMac := func() bool {
	// 		for _, m := range host.MACAddress {
	// 			if strings.EqualFold(m, node.MACAddress) {
	// 				return true
	// 			}
	// 		}
	// 		return false
	// 	}
	//
	// 	if hasMac() {
	// 		log.Debugf("host %s reporting", host.DN)
	// 	} else {
	// 		err = i.SetAttribute(host.DN, "macAddress", node.MACAddress, false)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// v4Address := i.GetAttribute(host.DN, "v4Address")
	// if v4Address != node.IP {
	// 	log.Warnf("alert alert")
	// }

	return nil
}
