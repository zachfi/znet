package znet

import (
	"context"

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
		return response, err
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

	unknownHosts, err := r.inventory.UnknownHosts()
	if err != nil {
		return response, err
	}

	for _, h := range unknownHosts {
		host := &pb.UnknownHost{
			Name: h.Name,
			Ip:   h.IP,
			Mac:  h.MACAddress,
		}

		response.UnknownHosts = append(response.UnknownHosts, host)
	}

	return response, nil
}
