package znet

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/iot"
	pb "github.com/xaque208/znet/rpc"
)

type thingServer struct {
	inventory  *inventory.Inventory
	keeper     thingKeeper
	seenThings map[string]time.Time
}

type thingKeeper map[string]map[string]string

// NewThingServer returns a new thingServer.
func newThingServer(inv *inventory.Inventory) *thingServer {
	s := thingServer{
		inventory:  inv,
		keeper:     make(thingKeeper),
		seenThings: make(map[string]time.Time),
	}

	go func(s thingServer) {
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
func (l *thingServer) storeThingLabel(nodeID string, key, value string) {
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

func (l *thingServer) nodeLabels(nodeID string) map[string]string {
	if nodeLabelMap, ok := l.keeper[nodeID]; ok {
		return nodeLabelMap
	}

	return map[string]string{}
}

// hasLabels checks to see if the keeper has all of the received labels for the given node ID.
func (l *thingServer) hasLabels(nodeID string, labels []string) bool {
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

func (l *thingServer) Notice(ctx context.Context, request *pb.DeviceDiscovery) (*pb.NoticeResponse, error) {
	response := &pb.NoticeResponse{}

	if request.ObjectID != "" {
		rpcThingServerObjectNotice.WithLabelValues(request.ObjectID).Inc()
	}

	switch request.ObjectID {
	case "wifi":
		msg := iot.ReadMessage("wifi", request.Message)
		m := msg.(iot.WifiMessage)

		l.storeThingLabel(request.NodeID, "ssid", m.SSID)
		l.storeThingLabel(request.NodeID, "bssid", m.BSSID)
		l.storeThingLabel(request.NodeID, "ip", m.IP)

		labels := l.nodeLabels(request.NodeID)

		if l.hasLabels(request.NodeID, []string{"ssid", "bssid", "ip"}) {
			if m.RSSI != 0 {
				thingWireless.With(prometheus.Labels{
					"device": request.NodeID,
					"ssid":   labels["ssid"],
					"bssid":  labels["ssid"],
					"ip":     labels["ip"],
				}).Set(float64(m.RSSI))
			}
		}

	case "air":
		msg := iot.ReadMessage("air", request.Message)
		m := msg.(iot.AirMessage)

		// l.storeThingLabel(request.NodeID, "tempcoef", m.TempCoef)

		airTemperature.WithLabelValues(request.NodeID).Set(float64(m.Temperature))
		airHumidity.WithLabelValues(request.NodeID).Set(float64(m.Humidity))
		airHeatindex.WithLabelValues(request.NodeID).Set(float64(m.HeatIndex))
	case "led1", "led2":
		msg := iot.ReadMessage("led", request.Message)
		m := msg.(iot.LEDMessage)

		for i, deviceConnection := range m.Device.Connections {
			if len(deviceConnection) == 2 {
				l.storeThingLabel(request.NodeID, "mac", m.Device.Connections[i][1])
			}
		}
	default:
		rpcThingServerUnhandledObjectNotice.WithLabelValues(request.ObjectID).Inc()
	}

	// Record an observation if all our parts are filled in.
	labels := l.nodeLabels(request.NodeID)
	if l.hasLabels(request.NodeID, []string{"ip", "mac"}) {
		iotNode := iot.Node{
			IP:         labels["ip"],
			MACAddress: labels["mac"],
			NodeID:     request.NodeID,
		}

		err := l.inventory.ObserveIOT(iotNode)
		if err != nil {
			log.Error(err)
		}
	}

	return response, nil
}
