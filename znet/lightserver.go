package znet

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/xaque208/znet/rpc"
)

type lightServer struct {
	lights *Lights
}

func (l *lightServer) Off(ctx context.Context, request *pb.LightGroup) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	l.lights.Off(request.Name)

	return response, nil
}

func (l *lightServer) On(ctx context.Context, request *pb.LightGroup) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	l.lights.On(request.Name)

	return response, nil
}

func (l *lightServer) Status(ctx context.Context, request *pb.LightRequest) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	lights := l.lights.Status()

	for _, light := range lights {

		state := &pb.State{
			On:         light.State.On,
			Brightness: int32(light.State.Bri),
		}

		x := &pb.Light{
			Name:  light.Name,
			Type:  light.Type,
			Id:    int32(light.ID),
			State: state,
			// Brightness: int32(light.State.Bri),
		}

		response.Lights = append(response.Lights, x)
	}

	groups, err := l.lights.HUE.GetGroups()
	if err != nil {
		return response, err
	}

	for _, group := range groups {

		state := &pb.State{
			On:         group.State.On,
			Brightness: int32(group.State.Bri),
		}

		x := &pb.LightGroup{
			Name:  group.Name,
			Type:  group.Type,
			Id:    int32(group.ID),
			State: state,
		}

		response.Groups = append(response.Groups, x)
	}

	return response, nil
}

func (l *lightServer) Brightness(ctx context.Context, request *pb.LightGroup) (*pb.LightResponse, error) {

	response := &pb.LightResponse{}

	room, err := l.lights.config.Room(request.Name)
	if err != nil {
		return response, err
	}

	groups, err := l.lights.HUE.GetGroups()
	if err != nil {
		return response, err
	}

	for _, g := range groups {
		for _, i := range room.HueIDs {
			if g.ID == i {
				log.Debugf("Setting brightness for group %s: %+v", g.Name, g.State)
				err := g.Bri(uint8(request.State.Brightness))
				if err != nil {
					return response, err
				}

				state := &pb.State{
					On:         g.State.On,
					Brightness: int32(g.State.Bri),
				}

				x := &pb.LightGroup{
					Name:  g.Name,
					Type:  g.Type,
					Id:    int32(g.ID),
					State: state,
				}

				response.Groups = append(response.Groups, x)
			}
		}

	}

	return response, nil
}
