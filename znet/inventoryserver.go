package znet

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/xaque208/znet/rpc"
)

// RPC Listener
type inventoryServer struct {
	inventory *Inventory
}

func (r *inventoryServer) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	response := &pb.SearchResponse{}

	hosts, err := r.inventory.NetworkHosts()
	if err != nil {
		log.Error(err)
	}

	for _, h := range hosts {
		host := &pb.Host{
			Name:        h.Name,
			Description: h.Description,
			Platform:    h.Platform,
			Type:        h.DeviceType,
		}

		response.Hosts = append(response.Hosts, host)
	}

	log.Warnf("%+v: ", r)

	return response, nil
}
