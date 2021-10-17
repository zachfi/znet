package harvester

import (
	"context"

	"github.com/grafana/dskit/services"
	"github.com/xaque208/znet/internal/telemetry"
)

type Harvester struct {
	services.Service

	cfg             *Config
	telemetryClient telemetry.TelemetryClient
}

func New(cfg Config) (*Harvester, error) {
	h := &Harvester{
		cfg: &cfg,
	}

	h.Service = services.NewBasicService(h.starting, h.running, h.stopping)

	return h, nil
}

func (h *Harvester) starting(ctx context.Context) error {
	// conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)
	//
	// defer func() {
	// 	err = conn.Close()
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	// }()
	//
	// telemetryClient := telemetry.NewTelemetryClient(conn)
	// h.telemetryClient = telemetryClient

	return nil
}

func (h *Harvester) running(ctx context.Context) error {

	<-ctx.Done()

	return nil
}

func (h *Harvester) stopping(_ error) error {
	return nil
}
