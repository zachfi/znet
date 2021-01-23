package timer

import (
	"context"

	"github.com/xaque208/znet/internal/lights"
)

type Server struct {
	lights *lights.Lights
}

// NewServer returns a new Server.
func NewServer(lig *lights.Lights) (*Server, error) {
	return &Server{
		lights: lig,
	}, nil
}

func (t *Server) NamedTimer(ctx context.Context, req *NamedTimeStamp) (*Empty, error) {

	err := t.lights.NamedTimerHandler(req.Name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
