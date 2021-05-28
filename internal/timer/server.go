package timer

import (
	"context"
	"fmt"

	"github.com/xaque208/znet/internal/lights"
)

type Server struct {
	UnimplementedTimerServer
	lights *lights.Lights
}

// NewServer returns a new Server.
func NewServer(l *lights.Lights) (*Server, error) {
	return &Server{
		lights: l,
	}, nil
}

func (t *Server) NamedTimer(ctx context.Context, req *NamedTimeStamp) (*Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("unable to handle nil request")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("unable to handle request with empty name")
	}

	err := t.lights.NamedTimerHandler(req.Name)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}
