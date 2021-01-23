package network

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/telemetry"
	"github.com/xaque208/znet/internal/timer"
	"github.com/xaque208/znet/pkg/events"
)

func (n *Network) namedTimerHandler(name string, payload events.Payload) error {
	var e timer.NamedTimer

	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", e, err)
	}

	switch e.Name {
	case "macUpdate":
	}

	// go func() {
	// 	for range ticker.C {
	// 		wg := sync.WaitGroup{}
	// 		stream, err := inventoryClient.ListNetworkHosts(context.Background(), &rpc.Empty{})
	// 		if err != nil {
	// 			log.Error(err)
	// 		}
	//
	// 		if resp != nil {
	// 			for _, h := range resp.Hosts {
	//
	// 				if h.Platform == "" {
	// 					continue
	// 				}
	//
	// 				if h.Name == "" {
	// 					continue
	// 				}
	//
	// 				if h.Platform == "junos" {
	// 					wg.Add(1)
	// 					go network.ScrapeJunosHost(&wg, h)
	// 				}
	// 			}
	// 		}
	//
	// 		wg.Wait()
	// 	}
	// }()

	return nil
}

func (n *Network) ScrapeJunosHost(wg *sync.WaitGroup, h *inventory.NetworkHost) {
	telemetryClient := telemetry.NewTelemetryClient(n.conn)

	hostName := strings.Join([]string{h.Name, h.Domain}, ".")
	log.Debugf("scraping ARP status from host: %s", hostName)

	auth := &junos.AuthMethod{
		Username:   n.config.Junos.Username,
		PrivateKey: n.config.Junos.PrivateKey,
	}

	session, err := junos.NewSession(hostName, auth)
	if err != nil {
		log.Error(err)
		return
	}
	defer session.Close()
	defer wg.Done()

	views, err := session.View("arp")
	if err != nil {
		log.Error(err)
		wg.Done()
		return
	}

	for _, arp := range views.Arp.Entries {
		if arp.Interface == "ppp0.0" {
			continue
		}

		log.Tracef("reporting NetworkID: %+v", arp)

		name := strings.ToLower(strings.Join([]string{arp.MACAddress, arp.Interface}, "_"))

		networkID := &inventory.NetworkID{
			Name:                     name,
			IpAddress:                []string{arp.IPAddress},
			MacAddress:               []string{arp.MACAddress},
			ReportingSource:          []string{h.Name},
			ReportingSourceInterface: []string{arp.Interface},
		}

		_, err = telemetryClient.ReportNetworkID(context.Background(), networkID)
		if err != nil {
			log.Error(err)
		}

	}
}
