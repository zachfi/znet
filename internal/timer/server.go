package timer

import (
	"context"

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
	err := t.lights.NamedTimerHandler(req.Name)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}
