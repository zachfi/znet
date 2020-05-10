package znet

import (
	"context"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/lights"
	pb "github.com/xaque208/znet/rpc"
)

type lightServer struct {
	lights *lights.Lights
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

	lights, err := l.lights.HUE.GetLights()
	if err != nil {
		log.Error(err)
	}

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
		}

		newName := strings.ToLower(strings.ReplaceAll(light.Name, " ", "_"))

		if light.Name != newName {
			err = light.Rename(newName)
			if err != nil {
				log.Error(err)
			}
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

		for _, l := range group.Lights {
			n, err := strconv.Atoi(l)
			if err != nil {
				log.Error(err)
			}
			x.Lights = append(x.Lights, int32(n))
		}

		newName := strings.ToLower(strings.ReplaceAll(group.Name, " ", "_"))

		if group.Name != newName {
			err := group.Rename(newName)
			if err != nil {
				log.Error(err)
			}
		}

		response.Groups = append(response.Groups, x)
	}

	return response, nil
}

func (l *lightServer) Brightness(ctx context.Context, request *pb.LightGroup) (*pb.LightResponse, error) {
	response := &pb.LightResponse{}

	log.Infof("request: %+v", *request)

	l.lights.Dim(request.Name, request.State.Brightness)

	return response, nil
}

func (l *lightServer) Alert(ctx context.Context, request *pb.LightGroup) (*pb.LightResponse, error) {
	response := &pb.LightResponse{}

	log.Infof("request: %+v", *request)

	l.lights.Alert(request.Name)

	return response, nil
}
