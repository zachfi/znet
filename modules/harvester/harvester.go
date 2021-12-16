package harvester

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/xaque208/znet/modules/inventory"
	"github.com/xaque208/znet/modules/telemetry"
	"github.com/xaque208/znet/pkg/iot"
)

type Harvester struct {
	services.Service
	cfg *Config

	logger log.Logger

	conn            *grpc.ClientConn
	telemetryClient telemetry.TelemetryClient
}

func New(cfg Config, logger log.Logger, conn *grpc.ClientConn) (*Harvester, error) {
	h := &Harvester{
		cfg:    &cfg,
		logger: log.With(logger, "module", "harvester"),
		conn:   conn,
	}

	h.Service = services.NewBasicService(h.starting, h.running, h.stopping)

	return h, nil
}

func (h *Harvester) starting(ctx context.Context) error {
	telemetryClient := telemetry.NewTelemetryClient(h.conn)
	h.telemetryClient = telemetryClient

	return nil
}

func (h *Harvester) running(ctx context.Context) error {
	var onMessageReceived mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		topicPath, err := iot.ParseTopicPath(msg.Topic())
		if err != nil {
			_ = level.Error(h.logger).Log("err", errors.Wrap(err, "failed to parse topic path"))
			return
		}

		discovery := iot.ParseDiscoveryMessage(topicPath, msg)

		iotDevice := &inventory.IOTDevice{
			DeviceDiscovery: discovery,
		}

		_, err = h.telemetryClient.ReportIOTDevice(ctx, iotDevice)
		if err != nil {
			_ = level.Error(h.logger).Log("err", err.Error())
		}
	}

	mqttClient, err := iot.NewMQTTClient(h.cfg.MQTT, h.logger)
	if err != nil {
		return err
	}

	token := mqttClient.Subscribe(h.cfg.MQTT.Topic, 0, onMessageReceived)
	token.Wait()
	if token.Error() != nil {
		_ = level.Error(h.logger).Log("err", token.Error())
	}

	<-ctx.Done()

	return nil
}

func (h *Harvester) stopping(_ error) error {
	return h.conn.Close()
}
