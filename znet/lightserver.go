package znet

import (
	"context"

	pb "github.com/xaque208/znet/rpc"
)

type lightServer struct {
	lights *Lights
}

func (l *lightServer) Off(ctx context.Context, request *pb.LightZone) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	l.lights.Off(request.Name)

	return response, nil
}

func (l *lightServer) On(ctx context.Context, request *pb.LightZone) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	l.lights.On(request.Name)

	return response, nil
}

func (l *lightServer) Status(ctx context.Context, request *pb.LightRequest) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	lights := l.lights.Status()

	for _, light := range lights {

		x := &pb.Light{
			Name: light.Name,
			Type: light.Type,
			Id:   int32(light.ID),
		}

		response.Lights = append(response.Lights, x)
	}

	return response, nil
}
