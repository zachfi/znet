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

		discovery := &iot.DeviceDiscovery{
			Component: topicPath.Component,
			NodeId:    topicPath.NodeID,
			ObjectId:  topicPath.ObjectID,
			Endpoint:  topicPath.Endpoint,
			Message:   msg.Payload(),
		}

		iotDevice := &inventory.IOTDevice{
			DeviceDiscovery: discovery,
		}

		_, err = h.telemetryClient.ReportIOTDevice(ctx, iotDevice)
		if err != nil {
			_ = level.Error(h.logger).Log("err", err.Error())
		}
	}

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(h.cfg.MQTT.URL)
	mqttOpts.SetCleanSession(true)
	mqttOpts.OnConnect = func(c mqtt.Client) {
		token := c.Subscribe(h.cfg.MQTT.Topic, 0, onMessageReceived)
		token.Wait()
		if token.Error() != nil {
			_ = level.Error(h.logger).Log("err", token.Error())
		}
	}

	if h.cfg.MQTT.Username != "" && h.cfg.MQTT.Password != "" {
		mqttOpts.Username = h.cfg.MQTT.Username
		mqttOpts.Password = h.cfg.MQTT.Password
	}

	mqttClient := mqtt.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		_ = level.Error(h.logger).Log("err", token.Error())
	} else {
		_ = level.Debug(h.logger).Log("msg", "mqtt connected", "url", h.cfg.MQTT.URL)
	}

	// 	log.WithFields(log.Fields{
	// 		"http": listenAddr,
	// 	}).Debug("listening")
	//
	// 	go func() {
	// 		sig := <-sigs
	// 		log.Warnf("caught signal: %s", sig.String())
	//
	// 		done <- true
	// 	}()
	//
	// 	<-done
	//
	// 	if token := mqttClient.Unsubscribe(mqttTopic); token.Wait() && token.Error() != nil {
	// 		log.Error(token.Error())
	// 	}
	//
	// 	mqttClient.Disconnect(250)

	<-ctx.Done()

	return nil
}

func (h *Harvester) stopping(_ error) error {
	return h.conn.Close()
}
