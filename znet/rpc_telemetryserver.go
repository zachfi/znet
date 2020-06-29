package znet

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/iot"
	pb "github.com/xaque208/znet/rpc"
)

type telemetryServer struct {
	inventory  *inventory.Inventory
	keeper     thingKeeper
	seenThings map[string]time.Time
}

type thingKeeper map[string]map[string]string

// newThingServer returns a new telemetryServer.
func newTelemetryServer(inv *inventory.Inventory) *telemetryServer {
	s := telemetryServer{
		inventory:  inv,
		keeper:     make(thingKeeper),
		seenThings: make(map[string]time.Time),
	}

	go func(s telemetryServer) {
		for {
			// Make a copy
			tMap := make(map[string]time.Time)
			for k, v := range s.seenThings {
				tMap[k] = v
			}

			// Expire the old entries
			for k, v := range tMap {
				if time.Since(v) > (600 * time.Second) {
					log.Infof("expiring device %s", k)

					airTemperature.Delete(prometheus.Labels{"device": k})
					airHumidity.Delete(prometheus.Labels{"device": k})
					airHeatindex.Delete(prometheus.Labels{"device": k})
					thingWireless.Delete(prometheus.Labels{"device": k})

					delete(s.seenThings, k)
					delete(s.keeper, k)
				}
			}

			time.Sleep(30 * time.Second)
		}
	}(s)

	return &s
}

// storeThingLabel records the received key/value pair for the given node ID.
func (l *telemetryServer) storeThingLabel(nodeID string, key, value string) {
	if len(l.keeper) == 0 {
		l.keeper = make(thingKeeper)
	}

	if _, ok := l.keeper[nodeID]; !ok {
		l.keeper[nodeID] = make(map[string]string)
	}

	if key != "" && value != "" {
		l.keeper[nodeID][key] = value
	}
}

func (l *telemetryServer) nodeLabels(nodeID string) map[string]string {
	if nodeLabelMap, ok := l.keeper[nodeID]; ok {
		return nodeLabelMap
	}

	return map[string]string{}
}

// hasLabels checks to see if the keeper has all of the received labels for the given node ID.
func (l *telemetryServer) hasLabels(nodeID string, labels []string) bool {
	nodeLabels := l.nodeLabels(nodeID)

	nodeHasLabel := func(nodeLabels map[string]string, label string) bool {

		for key := range nodeLabels {
			if key == label {
				return true
			}
		}

		return false
	}

	for _, label := range labels {
		if !nodeHasLabel(nodeLabels, label) {
			return false
		}
	}

	return true
}

func (l *telemetryServer) findMACs(macs []string) (*[]inventory.NetworkHost, *[]inventory.NetworkID, error) {
	var keepHosts []inventory.NetworkHost
	var keepIds []inventory.NetworkID

	networkHosts, err := l.inventory.ListNetworkHosts()
	if err != nil {
		return nil, nil, err
	}

	if networkHosts != nil {
		for _, x := range *networkHosts {
			if x.MacAddress != nil {
				for _, m := range *x.MacAddress {
					for _, mm := range macs {
						if strings.EqualFold(m, mm) {
							keepHosts = append(keepHosts, x)
						}
					}
				}
			}
		}
	}

	networkIDs, err := l.inventory.ListNetworkIDs()
	if err != nil {
		return nil, nil, err
	}

	if networkIDs != nil {
		for _, x := range *networkIDs {
			if x.MacAddress != nil {
				for _, m := range *x.MacAddress {
					for _, mm := range macs {
						if strings.EqualFold(m, mm) {
							keepIds = append(keepIds, x)
						}
					}
				}
			}
		}
	}

	return &keepHosts, &keepIds, nil
}

func (l *telemetryServer) ReportNetworkID(ctx context.Context, request *pb.NetworkID) (*pb.Empty, error) {
	if request.Name == "" {
		return &pb.Empty{}, fmt.Errorf("unable to fetch NetworkID with empty name")
	}

	hosts, ids, err := l.findMACs(request.MacAddress)
	if err != nil {
		return &pb.Empty{}, err
	}

	// do nothing if a host matches
	if hosts != nil {
		if len(*hosts) > 0 {
			for _, h := range *ids {
				err = l.inventory.UpdateTimestamp(h.Dn, "networkHost")
				if err != nil {
					log.Error(err)
				}
			}
			return &pb.Empty{}, nil
		}
	}

	now := time.Now()

	// update the lastSeen for nettworkIds
	if ids != nil {
		if len(*ids) > 0 {
			log.Debugf("ids found for report: %+v", *ids)
			for _, h := range *ids {
				if h.Dn != "" {
					x := inventory.NetworkID{
						Dn:                       h.Dn,
						IpAddress:                &request.IpAddress,
						MacAddress:               &request.MacAddress,
						ReportingSource:          &request.ReportingSource,
						ReportingSourceInterface: &request.ReportingSourceInterface,
						LastSeen:                 &now,
					}

					_, err = l.inventory.UpdateNetworkID(x)
					if err != nil {
						return &pb.Empty{}, err
					}
				}
			}
		}
	}

	log.Debugf("existing mac not found: %+v", request.MacAddress)

	x := inventory.NetworkID{
		Name:                     request.Name,
		IpAddress:                &request.IpAddress,
		MacAddress:               &request.MacAddress,
		ReportingSource:          &request.ReportingSource,
		ReportingSourceInterface: &request.ReportingSourceInterface,
		LastSeen:                 &now,
	}

	_, err = l.inventory.FetchNetworkID(request.Name)
	if err != nil {
		_, err = l.inventory.CreateNetworkID(x)
		if err != nil {
			return &pb.Empty{}, err
		}
	}

	return &pb.Empty{}, nil
}

func (l *telemetryServer) ReportIOTDevice(ctx context.Context, request *pb.IOTDevice) (*pb.Empty, error) {
	if request.Name == "" {
		return &pb.Empty{}, fmt.Errorf("unable to fetch NetworkID with empty name")
	}

	//GOAL here we would like to ensure IOT deivce information is written to the
	//host entry that belongs to this device.

	// Check the message content for an IP address or mac.  Report

	// hosts, ids, err := l.findMACs(request.NetworkId.GetMacAddress())
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	//
	// // do nothing if we did not match any networkHosts
	// if hosts != nil {
	// 	return &pb.Empty{}, nil
	// }
	//
	// if ids != nil {
	// 	return &pb.Empty{}, nil
	// }

	discovery := request.DeviceDiscovery

	if discovery.ObjectId != "" {
		rpcThingServerObjectNotice.WithLabelValues(discovery.ObjectId).Inc()
	}

	switch discovery.ObjectId {
	case "wifi":
		msg := iot.ReadMessage("wifi", discovery.Message, discovery.Endpoint...)
		if msg != nil {
			m := msg.(iot.WifiMessage)

			l.storeThingLabel(discovery.NodeId, "ssid", m.SSID)
			l.storeThingLabel(discovery.NodeId, "bssid", m.BSSID)
			l.storeThingLabel(discovery.NodeId, "ip", m.IP)

			labels := l.nodeLabels(discovery.NodeId)

			if l.hasLabels(discovery.NodeId, []string{"ssid", "bssid", "ip"}) {
				if m.RSSI != 0 {
					thingWireless.With(prometheus.Labels{
						"device": discovery.NodeId,
						"ssid":   labels["ssid"],
						"bssid":  labels["ssid"],
						"ip":     labels["ip"],
					}).Set(float64(m.RSSI))
				}
			}
		}

	case "air":
		msg := iot.ReadMessage("air", discovery.Message, discovery.Endpoint...)
		if msg != nil {
			m := msg.(iot.AirMessage)

			// l.storeThingLabel(discovery.NodeId, "tempcoef", m.TempCoef)

			airTemperature.WithLabelValues(discovery.NodeId).Set(float64(m.Temperature))
			airHumidity.WithLabelValues(discovery.NodeId).Set(float64(m.Humidity))
			airHeatindex.WithLabelValues(discovery.NodeId).Set(float64(m.HeatIndex))
		}

	case "led1", "led2":
		msg := iot.ReadMessage("led", discovery.Message, discovery.Endpoint...)
		if msg != nil {
			m := msg.(iot.LEDConfig)

			for i, deviceConnection := range m.Device.Connections {
				if len(deviceConnection) == 2 {
					l.storeThingLabel(discovery.NodeId, "mac", m.Device.Connections[i][1])
				}
			}
		}

	default:
		rpcThingServerUnhandledObjectNotice.WithLabelValues(discovery.ObjectId).Inc()
	}

	// Record an observation if all our parts are filled in.
	labels := l.nodeLabels(discovery.NodeId)
	if l.hasLabels(discovery.NodeId, []string{"ip", "mac"}) {
		iotNode := iot.Node{
			IP:         labels["ip"],
			MACAddress: labels["mac"],
			NodeID:     discovery.NodeId,
		}

		err := l.inventory.ObserveIOT(iotNode)
		if err != nil {
			log.Error(err)
		}
	}

	return &pb.Empty{}, nil
}
