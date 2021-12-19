package telemetry

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/xaque208/znet/modules/inventory"
	"github.com/xaque208/znet/modules/lights"
	"github.com/xaque208/znet/pkg/iot"
)

const defaultExpiry = 5 * time.Minute

type Telemetry struct {
	UnimplementedTelemetryServer

	services.Service
	cfg *Config

	logger log.Logger
	tracer trace.Tracer

	inventory  inventory.Inventory
	keeper     thingKeeper
	lights     *lights.Lights
	iotServer  *iot.Server
	seenThings map[string]time.Time
}

type thingKeeper map[string]map[string]string

func New(cfg Config, logger log.Logger, inv inventory.Inventory, lig *lights.Lights) (*Telemetry, error) {
	s := &Telemetry{
		cfg:    &cfg,
		logger: log.With(logger, "module", "telemetry"),
		tracer: otel.Tracer("telemetry"),

		inventory:  inv,
		keeper:     make(thingKeeper),
		lights:     lig,
		seenThings: make(map[string]time.Time),
	}

	go func(s *Telemetry) {
		for {
			// Make a copy
			tMap := make(map[string]time.Time)
			for k, v := range s.seenThings {
				tMap[k] = v
			}

			// Expire the old entries
			for k, v := range tMap {
				if time.Since(v) > defaultExpiry {
					_ = level.Info(s.logger).Log("msg", "expiring",
						"device", k,
					)

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

	s.Service = services.NewBasicService(s.starting, s.running, s.stopping)

	return s, nil
}

func (l *Telemetry) starting(ctx context.Context) error {
	return nil
}

func (l *Telemetry) running(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (l *Telemetry) stopping(_ error) error {
	return nil
}

// storeThingLabel records the received key/value pair for the given node ID.
func (l *Telemetry) storeThingLabel(nodeID string, key, value string) {
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

func (l *Telemetry) nodeLabels(nodeID string) map[string]string {
	if nodeLabelMap, ok := l.keeper[nodeID]; ok {
		return nodeLabelMap
	}

	return map[string]string{}
}

// hasLabels checks to see if the keeper has all of the received labels for the given node ID.
func (l *Telemetry) hasLabels(nodeID string, labels []string) bool {
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

func (l *Telemetry) findMACs(ctx context.Context, macs []string) ([]*inventory.NetworkHost, []*inventory.NetworkID, error) {
	var keepHosts []*inventory.NetworkHost
	var keepIds []*inventory.NetworkID

	networkHosts, err := l.inventory.ListNetworkHosts(ctx)
	if err != nil {
		return nil, nil, err
	}

	for i := range networkHosts {
		x := proto.Clone(&(networkHosts)[i]).(*inventory.NetworkHost)

		if x.MacAddress != nil {
			for _, m := range x.MacAddress {
				for _, mm := range macs {
					if strings.EqualFold(m, mm) {
						keepHosts = append(keepHosts, x)
					}
				}
			}
		}
	}

	networkIDs, err := l.inventory.ListNetworkIDs(ctx)
	if err != nil {
		return nil, nil, err
	}

	for i := range networkIDs {
		x := proto.Clone(&(networkIDs)[i]).(*inventory.NetworkID)

		if x.MacAddress != nil {
			for _, m := range x.MacAddress {
				for _, mm := range macs {
					if strings.EqualFold(m, mm) {
						keepIds = append(keepIds, x)
					}
				}
			}
		}
	}

	return keepHosts, keepIds, nil
}

func (l *Telemetry) ReportNetworkID(ctx context.Context, request *inventory.NetworkID) (*inventory.Empty, error) {
	if request.Name == "" {
		return &inventory.Empty{}, fmt.Errorf("unable to fetch inventory.NetworkID with empty name")
	}

	spanCtx, span := l.tracer.Start(ctx, "ReportNetworkID")
	defer span.End()

	hosts, ids, err := l.findMACs(spanCtx, request.MacAddress)
	if err != nil {
		return &inventory.Empty{}, err
	}

	// do nothing if a host matches
	if len(hosts) > 0 {
		for _, x := range ids {
			err = l.inventory.UpdateTimestamp(spanCtx, x.Dn, "networkHost")
			if err != nil {
				_ = level.Error(l.logger).Log("err", err.Error())
			}
		}
		return &inventory.Empty{}, nil
	}

	now := time.Now()

	// update the lastSeen for nettworkIds
	if len(ids) > 0 {
		for _, id := range ids {
			if id.Dn != "" {
				x := &inventory.NetworkID{
					Dn:                       id.Dn,
					IpAddress:                request.IpAddress,
					MacAddress:               request.MacAddress,
					ReportingSource:          request.ReportingSource,
					ReportingSourceInterface: request.ReportingSourceInterface,
					LastSeen:                 timestamppb.New(now),
				}

				_, err = l.inventory.UpdateNetworkID(ctx, x)
				if err != nil {
					return &inventory.Empty{}, err
				}
			}
		}
	}

	_ = level.Debug(l.logger).Log("msg", "existing mac not found",
		"mac", request.MacAddress,
	)

	x := &inventory.NetworkID{
		Name:                     request.Name,
		IpAddress:                request.IpAddress,
		MacAddress:               request.MacAddress,
		ReportingSource:          request.ReportingSource,
		ReportingSourceInterface: request.ReportingSourceInterface,
		LastSeen:                 timestamppb.New(now),
	}

	_, err = l.inventory.FetchNetworkID(ctx, request.Name)
	if err != nil {
		_, err = l.inventory.CreateNetworkID(ctx, x)
		if err != nil {
			return &inventory.Empty{}, err
		}
	}

	return &inventory.Empty{}, nil
}

func (l *Telemetry) ReportIOTDevice(ctx context.Context, request *inventory.IOTDevice) (*inventory.Empty, error) {
	var err error

	if request.DeviceDiscovery == nil {
		return nil, fmt.Errorf("unable to report IOTDevice with nil DeviceDiscovery")
	}

	discovery := request.DeviceDiscovery

	spanCtx, span := l.tracer.Start(ctx, fmt.Sprintf("ReportIOTDevice/%s/%s", discovery.Component, discovery.ObjectId))
	defer span.End()

	if discovery.ObjectId != "" {
		telemetryIOTReport.WithLabelValues(discovery.ObjectId, discovery.Component).Inc()
	}

	switch discovery.Component {
	case "zigbee2mqtt":
		err = l.handleZigbeeReport(spanCtx, request)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle zigbee report")
		}
	}

	switch discovery.ObjectId {
	case "wifi":
		err = l.handleWifiReport(request)
		if err != nil {
			return nil, err
		}
	case "air":
		err = l.handleAirReport(request)
		if err != nil {
			return nil, err
		}
	case "water":
		err = l.handleWaterReport(request)
		if err != nil {
			return nil, err
		}
	case "led1", "led2":
		err = l.handleLEDReport(request)
		if err != nil {
			return nil, err
		}
	default:
		telemetryIOTUnhandledReport.WithLabelValues(discovery.ObjectId, discovery.Component).Inc()
	}

	return &inventory.Empty{}, nil
}

func (l *Telemetry) SetIOTServer(iotServer *iot.Server) error {
	if l.iotServer != nil {
		_ = level.Debug(l.logger).Log("replacing iotServer on telemetryServer")
	}

	l.iotServer = iotServer

	return nil
}

func (l *Telemetry) handleZigbeeReport(ctx context.Context, request *inventory.IOTDevice) error {
	if request == nil {
		return fmt.Errorf("unable to read zigbee report from nil request")
	}

	ctx, span := l.tracer.Start(ctx, "handleZigbeeReport")
	defer span.End()

	_ = level.Debug(l.logger).Log("msg", "device report", "traceID", trace.SpanContextFromContext(ctx).TraceID().String())

	discovery := request.DeviceDiscovery

	msg, err := iot.ReadZigbeeMessage(discovery.ObjectId, discovery.Message, discovery.Endpoint...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return errors.Wrap(err, "failed to read zigbee message")
	}

	if msg == nil {
		return nil
	}

	now := time.Now()

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

		if m.Message == nil {
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("unhandled iot.ZigbeeBridgeLog type: %s", m.Type)
		}

		messageTypeName := reflect.TypeOf(m.Message).String()

		switch messageTypeName {
		case "string":
			if strings.HasPrefix(m.Message.(string), "Update available") {
				return l.handleZigbeeDeviceUpdate(ctx, m)
			}
		case "iot.ZigbeeMessageBridgeDevices":
			return l.handleZigbeeDevices(ctx, m.Message.(iot.ZigbeeMessageBridgeDevices))
		default:
			return fmt.Errorf("unhandled iot.ZigbeeBridgeLog: %s", messageTypeName)
		}

	case "iot.ZigbeeMessageBridgeDevices":
		m := msg.(iot.ZigbeeMessageBridgeDevices)

		return l.handleZigbeeDevices(ctx, m)
	case "iot.ZigbeeMessage":
		m := msg.(iot.ZigbeeMessage)

		x := &inventory.ZigbeeDevice{
			Name:     request.DeviceDiscovery.ObjectId,
			LastSeen: timestamppb.New(now),
		}

		result, err := l.fetchZigbeeDevice(ctx, x)
		if err != nil {
			return err
		}

		l.updateZigbeeMessageMetrics(m, request, result)

		if m.Action != nil {
			action := &iot.Action{
				Event:  *m.Action,
				Device: x.Name,
				Zone:   result.IotZone,
			}

			err = l.lights.ActionHandler(ctx, action)
			if err != nil {
				_ = level.Error(l.logger).Log("err", err.Error())
			}
		}
	}

	return nil
}

func (l *Telemetry) handleZigbeeDevices(ctx context.Context, m iot.ZigbeeMessageBridgeDevices) error {
	ctx, span := l.tracer.Start(ctx, "handleZigbeeDevices")
	defer span.End()

	_ = level.Debug(l.logger).Log("msg", "devices report", "traceID", trace.SpanContextFromContext(ctx).TraceID().String())

	now := time.Now()

	for _, d := range m {

		x := &inventory.ZigbeeDevice{
			Name:     d.FriendlyName,
			LastSeen: timestamppb.New(now),
			// IeeeAddr:        d.IeeeAddr,
			Type:            iot.ZigbeeDeviceType(d).String(),
			SoftwareBuildId: d.SoftwareBuildID,
			DateCode:        d.DateCode,
			Model:           d.Definition.Model,
			Vendor:          d.Definition.Vendor,
			PowerSource:     d.PowerSource,
			ModelId:         d.ModelID,
		}

		if x.Name == "Coordinator" {
			continue
		}

		f, err := l.inventory.FetchZigbeeDevice(ctx, x.Name)
		if err != nil {
			createResult, createErr := l.inventory.CreateZigbeeDevice(ctx, x)
			if createErr != nil {
				span.SetStatus(codes.Error, createErr.Error())
				return createErr
			}

			_ = level.Debug(l.logger).Log("msg", "create result",
				"name", createResult.Name,
				"vendor", createResult.Vendor,
				"model", createResult.Model,
				"zone", createResult.IotZone,
			)
		}

		x.Dn = f.GetDn()
		_, updateErr := l.inventory.UpdateZigbeeDevice(ctx, x)
		if updateErr != nil {
			span.SetStatus(codes.Error, updateErr.Error())
			return updateErr
		}
	}

	return nil
}

func (l *Telemetry) handleZigbeeDeviceUpdate(ctx context.Context, m iot.ZigbeeBridgeLog) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	// zigbee2mqtt/bridge/request/device/ota_update/update
	_ = level.Debug(l.logger).Log("msg", "upgrade report",
		"device", m.Meta["device"],
		"status", m.Meta["status"],
	)

	req := &iot.UpdateRequest{
		Device: m.Meta["device"].(string),
	}

	go func() {
		_, err := l.iotServer.UpdateDevice(ctx, req)
		if err != nil {
			_ = level.Error(l.logger).Log("err", err.Error())
		}
	}()

	return nil
}

func (l *Telemetry) handleLEDReport(request *inventory.IOTDevice) error {
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

func (l *Telemetry) handleWaterReport(request *inventory.IOTDevice) error {
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

func (l *Telemetry) handleAirReport(request *inventory.IOTDevice) error {
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

func (l *Telemetry) handleWifiReport(request *inventory.IOTDevice) error {
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

func (l *Telemetry) updateZigbeeMessageMetrics(m iot.ZigbeeMessage, request *inventory.IOTDevice, device *inventory.ZigbeeDevice) {
	var zone string

	if val := device.GetIotZone(); val != "" {
		zone = val
	}

	deviceName := request.DeviceDiscovery.ObjectId
	component := request.DeviceDiscovery.Component

	if m.Battery != nil {
		telemetryIOTBatteryPercent.WithLabelValues(deviceName, component, zone).Set(float64(*m.Battery))
	}

	if m.LinkQuality != nil {
		telemetryIOTLinkQuality.WithLabelValues(deviceName, component, zone).Set(float64(*m.LinkQuality))
	}

	if m.Temperature != nil {
		telemetryIOTTemperature.WithLabelValues(deviceName, component, zone).Set(float64(*m.Temperature))
	}

	if m.Illuminance != nil {
		telemetryIOTIlluminance.WithLabelValues(deviceName, component, zone).Set(float64(*m.Illuminance))
	}

	if m.Occupancy != nil {
		if *m.Occupancy {
			telemetryIOTOccupancy.WithLabelValues(deviceName, component, zone).Set(float64(1))
		} else {
			telemetryIOTOccupancy.WithLabelValues(deviceName, component, zone).Set(float64(0))
		}
	}

	if m.WaterLeak != nil {
		if *m.WaterLeak {
			telemetryIOTWaterLeak.WithLabelValues(deviceName, component, zone).Set(float64(1))
		} else {
			telemetryIOTWaterLeak.WithLabelValues(deviceName, component, zone).Set(float64(0))
		}
	}

	if m.Tamper != nil {
		if *m.Tamper {
			telemetryIOTTamper.WithLabelValues(deviceName, component, zone).Set(float64(1))
		} else {
			telemetryIOTTamper.WithLabelValues(deviceName, component, zone).Set(float64(0))
		}
	}
}

func (l *Telemetry) fetchZigbeeDevice(ctx context.Context, x *inventory.ZigbeeDevice) (*inventory.ZigbeeDevice, error) {

	result, err := l.inventory.FetchZigbeeDevice(ctx, x.Name)
	if err != nil {
		result, err = l.inventory.CreateZigbeeDevice(ctx, x)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
