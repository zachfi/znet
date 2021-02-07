package network

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/grpc"

	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/telemetry"
)

// Agent is an RPC client worker bee.
type Network struct {
	config *config.NetworkConfig
	conn   *grpc.ClientConn
}

// NewAgent returns a new *Agent from the given arguments.
func NewNetwork(cfg *config.NetworkConfig, conn *grpc.ClientConn) *Network {
	return &Network{
		config: cfg,
		conn:   conn,
	}
}

// ScrapeJunosHost makes API calls against Junos calls to get some information.
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
