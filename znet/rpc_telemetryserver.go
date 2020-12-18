package znet

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/pkg/eventmachine"
	"github.com/xaque208/znet/pkg/iot"
	"github.com/xaque208/znet/rpc"
)

var (
	telemetryIOTUnhandledReport = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetry_unhandled_object_report",
		Help: "The total number of notice calls that include an unhandled object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTReport = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_telemetry_object_report",
		Help: "The total number of notice calls for an object ID.",
	}, []string{"object_id", "component"})

	telemetryIOTBatteryPercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_battery_percent",
		Help: "The reported batter percentage remaining.",
	}, []string{"object_id", "component"})

	telemetryIOTLinkQuality = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_link_quality",
		Help: "The reported link quality",
	}, []string{"object_id", "component"})

	telemetryIOTBridgeState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpc_telemetry_iot_bridge_state",
		Help: "The reported bridge state",
	}, []string{})
)

type telemetryServer struct {
	eventMachine *eventmachine.EventMachine
	inventory    *inventory.Inventory
	keeper       thingKeeper
	seenThings   map[string]time.Time
}

type thingKeeper map[string]map[string]string

// newThingServer returns a new telemetryServer.
func newTelemetryServer(inv *inventory.Inventory, eventMachine *eventmachine.EventMachine) *telemetryServer {
	s := &telemetryServer{
		eventMachine: eventMachine,
		inventory:    inv,
		keeper:       make(thingKeeper),
		seenThings:   make(map[string]time.Time),
	}

	go func(s *telemetryServer) {
		for {
			// Make a copy
			tMap := make(map[string]time.Time)
			for k, v := range s.seenThings {
				tMap[k] = v
			}

			// Expire the old entries
			for k, v := range tMap {
				if time.Since(v) > (600 * time.Second) {
					log.WithFields(log.Fields{
						"device": k,
					}).Info("expiring")

					airHeatindex.Delete(prometheus.Labels{"device": k})
					airHumidity.Delete(prometheus.Labels{"device": k})
					airTemperature.Delete(prometheus.Labels{"device": k})
					thingWireless.Delete(prometheus.Labels{"device": k})
					waterTemperature.Delete(prometheus.Labels{"device": k})

					delete(s.seenThings, k)
					delete(s.keeper, k)
				}
			}

			time.Sleep(30 * time.Second)
		}
	}(s)

	return s
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

func (l *telemetryServer) ReportNetworkID(ctx context.Context, request *rpc.NetworkID) (*rpc.Empty, error) {
	log.WithFields(log.Fields{
		"name":                       request.Name,
		"ip_address":                 request.IpAddress,
		"reporting_source":           request.ReportingSource,
		"reporting_source_interface": request.ReportingSourceInterface,
	}).Trace("NetworkID report")

	if request.Name == "" {
		return &rpc.Empty{}, fmt.Errorf("unable to fetch NetworkID with empty name")
	}

	hosts, ids, err := l.findMACs(request.MacAddress)
	if err != nil {
		return &rpc.Empty{}, err
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
			return &rpc.Empty{}, nil
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
						return &rpc.Empty{}, err
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
			return &rpc.Empty{}, err
		}
	}

	return &rpc.Empty{}, nil
}

func (l *telemetryServer) ReportIOTDevice(ctx context.Context, request *rpc.IOTDevice) (*rpc.Empty, error) {
	var err error

	log.WithFields(log.Fields{
		"component": request.DeviceDiscovery.Component,
		"node_id":   request.DeviceDiscovery.NodeId,
		"object_id": request.DeviceDiscovery.ObjectId,
		"endpoint":  request.DeviceDiscovery.Endpoint,
		"message":   string(request.DeviceDiscovery.Message),
	}).Trace("IOTDevice report")

	discovery := request.DeviceDiscovery

	if discovery.ObjectId != "" {
		telemetryIOTReport.WithLabelValues(discovery.ObjectId, discovery.Component).Inc()
	}

	switch discovery.Component {
	case "zigbee2mqtt":
		err = l.handleZigbeeReport(request)
		if err != nil {
			return &rpc.Empty{}, err
		}
	}

	switch discovery.ObjectId {
	case "wifi":
		err = l.handleWifiReport(request)
		if err != nil {
			return &rpc.Empty{}, err
		}
	case "air":
		err = l.handleAirReport(request)
		if err != nil {
			return &rpc.Empty{}, err
		}
	case "water":
		err = l.handleWaterReport(request)
		if err != nil {
			return &rpc.Empty{}, err
		}
	case "led1", "led2":
		err = l.handleLEDReport(request)
		if err != nil {
			return &rpc.Empty{}, err
		}
	default:
		telemetryIOTUnhandledReport.WithLabelValues(discovery.ObjectId, discovery.Component).Inc()
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

	return &rpc.Empty{}, nil
}

func (l *telemetryServer) handleZigbeeReport(request *rpc.IOTDevice) error {
	if request == nil {
		return fmt.Errorf("unable to read zigbee report from nil request")
	}

	discovery := request.DeviceDiscovery

	msg, err := iot.ReadZigbeeMessage(discovery.ObjectId, discovery.Message, discovery.Endpoint...)
	if err != nil {
		return err
	}

	now := time.Now()

	if msg != nil {
		switch reflect.TypeOf(msg).String() {
		case "iot.ZigbeeBridgeState":
			m := msg.(iot.ZigbeeBridgeState)
			switch m {
			case iot.Offline:
				telemetryIOTBridgeState.WithLabelValues().Set(float64(0))
			case iot.Online:
				telemetryIOTBridgeState.WithLabelValues().Set(float64(1))
			}

		case "iot.ZigbeeBridgeLog":
			m := msg.(iot.ZigbeeBridgeLog)

			if m.Message != nil {
				for _, d := range m.Message.(iot.ZigbeeBridgeMessageDevices) {
					x := inventory.ZigbeeDevice{
						Name:     d.FriendlyName,
						LastSeen: &now,
						// IeeeAddr:        d.IeeeAddr,
						Type:            d.Type,
						SoftwareBuildId: d.SoftwareBuildID,
						DateCode:        d.DateCode,
						Model:           d.Model,
						Vendor:          d.Vendor,
						// Description      : d.Description,
						ManufacturerName: d.ManufacturerName,
						PowerSource:      d.PowerSource,
						ModelId:          d.ModelID,
						// HardwareVersion:  d.HardwareVersion,
					}

					_, err := l.inventory.FetchZigbeeDevice(x.Name)
					if err != nil {
						log.Error(err)
						createResult, err := l.inventory.CreateZigbeeDevice(x)
						if err != nil {
							return err
						}

						log.WithFields(log.Fields{
							"name":   createResult.Name,
							"vendor": createResult.Vendor,
							"model":  createResult.Model,
							"zone":   createResult.IotZone,
						}).Debug("createResult")
					}
				}
			}

		case "iot.ZigbeeMessage":
			m := msg.(iot.ZigbeeMessage)

			if m.Battery > 0 {
				telemetryIOTBatteryPercent.WithLabelValues(request.DeviceDiscovery.ObjectId, request.DeviceDiscovery.Component).Set(float64(m.Battery))
			}

			if m.LinkQuality > 0 {
				telemetryIOTLinkQuality.WithLabelValues(request.DeviceDiscovery.ObjectId, request.DeviceDiscovery.Component).Set(float64(m.LinkQuality))
			}

			x := inventory.ZigbeeDevice{
				Name:     request.DeviceDiscovery.ObjectId,
				LastSeen: &now,
			}

			result, err := l.inventory.FetchZigbeeDevice(x.Name)
			if err != nil {
				log.Error(err)
				_, err = l.inventory.CreateZigbeeDevice(x)
				if err != nil {
					return err
				}
			}

			if m.Click != "" {
				ev := iot.Click{
					Count:  m.Click,
					Device: x.Name,
					Zone:   result.IotZone,
				}

				err = l.eventMachine.Send(ev)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	return nil
}

func (l *telemetryServer) handleLEDReport(request *rpc.IOTDevice) error {
	if request == nil {
		return fmt.Errorf("unable to read led report from nil request")
	}

	discovery := request.DeviceDiscovery

	msg, err := iot.ReadMessage("led", discovery.Message, discovery.Endpoint...)
	if err != nil {
		return err
	}

	if msg != nil {
		m := msg.(iot.LEDConfig)

		for i, deviceConnection := range m.Device.Connections {
			if len(deviceConnection) == 2 {
				l.storeThingLabel(discovery.NodeId, "mac", m.Device.Connections[i][1])
			}
		}
	}

	return nil
}

func (l *telemetryServer) handleWaterReport(request *rpc.IOTDevice) error {
	if request == nil {
		return fmt.Errorf("unable to read water report from nil request")
	}

	discovery := request.DeviceDiscovery

	msg, err := iot.ReadMessage("water", discovery.Message, discovery.Endpoint...)
	if err != nil {
		return err
	}

	if msg != nil {
		m := msg.(iot.WaterMessage)

		if m.Temperature != nil {
			waterTemperature.WithLabelValues(discovery.NodeId).Set(float64(*m.Temperature))
		}
	}

	return nil
}

func (l *telemetryServer) handleAirReport(request *rpc.IOTDevice) error {
	if request == nil {
		return fmt.Errorf("unable to read air report from nil request")
	}

	discovery := request.DeviceDiscovery

	msg, err := iot.ReadMessage("air", discovery.Message, discovery.Endpoint...)
	if err != nil {
		return err
	}

	if msg != nil {
		m := msg.(iot.AirMessage)

		// l.storeThingLabel(discovery.NodeId, "tempcoef", m.TempCoef)

		if m.Temperature != nil {
			airTemperature.WithLabelValues(discovery.NodeId).Set(float64(*m.Temperature))
		}

		if m.Humidity != nil {
			airHumidity.WithLabelValues(discovery.NodeId).Set(float64(*m.Humidity))
		}
		if m.HeatIndex != nil {
			airHeatindex.WithLabelValues(discovery.NodeId).Set(float64(*m.HeatIndex))
		}
	}

	return nil
}

func (l *telemetryServer) handleWifiReport(request *rpc.IOTDevice) error {
	if request == nil {
		return fmt.Errorf("unable to read wifi report from nil request")
	}

	discovery := request.DeviceDiscovery

	msg, err := iot.ReadMessage("wifi", discovery.Message, discovery.Endpoint...)
	if err != nil {
		return err
	}

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

	return nil
}
