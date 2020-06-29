package znet

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/inventory"
	pb "github.com/xaque208/znet/rpc"
)

// RPC Listener
type inventoryServer struct {
	inventory *inventory.Inventory
}

func (r *inventoryServer) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	response := &pb.SearchResponse{}

	hosts, err := r.inventory.ListNetworkHosts()
	if err != nil {
		return response, err
	}

	if hosts != nil {
		for _, h := range *hosts {

			host := &pb.NetworkHost{
				Description:     h.Description,
				Dn:              h.Dn,
				Domain:          h.Domain,
				Group:           h.Group,
				Inet6Address:    *h.Inet6Address,
				InetAddress:     *h.InetAddress,
				MacAddress:      *h.MacAddress,
				Name:            h.Name,
				OperatingSystem: h.OperatingSystem,
				Platform:        h.Platform,
				Role:            h.Role,
				Type:            h.Type,
			}

			if h.LastSeen != nil {
				lastSeen, convError := ptypes.TimestampProto(*h.LastSeen)
				if convError != nil {
					log.Error(convError)
				}
				host.LastSeen = lastSeen
			}

			response.Hosts = append(response.Hosts, host)
		}
	}

	networkIDs, err := r.inventory.ListNetworkIDs()
	if err != nil {
		return response, err
	}

	if networkIDs != nil {
		for _, h := range *networkIDs {
			host := &pb.NetworkID{
				Dn:                       h.Dn,
				Name:                     h.Name,
				IpAddress:                *h.IpAddress,
				MacAddress:               *h.MacAddress,
				ReportingSource:          *h.ReportingSource,
				ReportingSourceInterface: *h.ReportingSourceInterface,
			}

			response.NetworkIds = append(response.NetworkIds, host)
		}
	}

	return response, nil
}

func (r *inventoryServer) ListNetworkHosts(ctx context.Context, request *pb.Empty) (*pb.SearchResponse, error) {
	response := &pb.SearchResponse{}

	hosts, err := r.inventory.ListNetworkHosts()
	if err != nil {
		return response, err
	}

	if hosts != nil {
		for _, h := range *hosts {
			host := &pb.NetworkHost{
				Description:     h.Description,
				Dn:              h.Dn,
				Domain:          h.Domain,
				Group:           h.Group,
				Name:            h.Name,
				OperatingSystem: h.OperatingSystem,
				Platform:        h.Platform,
				Role:            h.Role,
				Type:            h.Type,
			}

			if h.InetAddress != nil {
				host.InetAddress = *h.InetAddress
			}

			if h.Inet6Address != nil {
				host.Inet6Address = *h.Inet6Address
			}

			if h.MacAddress != nil {
				host.MacAddress = *h.MacAddress
			}
			if h.LastSeen != nil {
				lastSeen, err := ptypes.TimestampProto(*h.LastSeen)
				if err != nil {
					log.Error(err)
				}
				host.LastSeen = lastSeen
			}

			response.Hosts = append(response.Hosts, host)
		}
	}

	return response, nil
}
