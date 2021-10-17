package znet

import (
	"fmt"

	"github.com/grafana/dskit/services"
	"github.com/xaque208/znet/modules/harvester"
	"github.com/xaque208/znet/modules/server"
)

type module int

const (
	Harvester module = iota
	// Timer
	// Agent
	Server
)

func (m module) String() string {
	return [...]string{
		"Harvester",
		"Timer",
		"Agent",
		"Server",
	}[m]
}

func (z *Znet) initHarvester() (services.Service, error) {
	h, err := harvester.New(z.cfg.Harvester)
	if err != nil {
		return nil, fmt.Errorf("failed to create harvester")
	}

	z.harvester = h
	return z.harvester, nil
}

func (z *Znet) initServer() (services.Service, error) {
	s, err := server.New(z.cfg.Server)
	if err != nil {
		return nil, fmt.Errorf("failed to create harvester")
	}

	z.server = s
	return z.server, nil
}
