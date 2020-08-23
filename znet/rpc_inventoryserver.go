package znet

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/rpc"
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

	pbHosts := getRPCHosts(hosts)
	if pbHosts != nil {
		response.Hosts = pbHosts
	}

	networkIDs, err := r.inventory.ListNetworkIDs()
	if err != nil {
		return response, err
	}

	pbNetworkIDs := getRPCNetworkIDs(networkIDs)
	if pbNetworkIDs != nil {
		response.NetworkIds = pbNetworkIDs
	}

	zigbeeDevices, err := r.inventory.ListZigbeeDevices()
	if err != nil {
		return response, err
	}

	pbZigbeeDevices := getRPCZigbeeDevices(zigbeeDevices)
	if pbZigbeeDevices != nil {
		response.ZigbeeDevices = pbZigbeeDevices
	}

	return response, nil
}

func (r *inventoryServer) Stop() {
	r.inventory.Close()
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

func getRPCHosts(hosts *[]inventory.NetworkHost) []*rpc.NetworkHost {
	var response []*rpc.NetworkHost

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

			if h.Inet6Address != nil {
				host.Inet6Address = *h.Inet6Address
			}

			if h.InetAddress != nil {
				host.InetAddress = *h.InetAddress
			}

			if h.MacAddress != nil {
				host.MacAddress = *h.MacAddress
			}

			if h.LastSeen != nil {
				lastSeen, convError := ptypes.TimestampProto(*h.LastSeen)
				if convError != nil {
					log.Error(convError)
				}
				host.LastSeen = lastSeen
			}

			response = append(response, host)
		}
	}

	return response
}

func getRPCNetworkIDs(networkIDs *[]inventory.NetworkID) []*rpc.NetworkID {
	var response []*rpc.NetworkID

	if networkIDs != nil {
		for _, h := range *networkIDs {
			host := &pb.NetworkID{
				Dn:   h.Dn,
				Name: h.Name,
			}

			if h.IpAddress != nil {
				host.IpAddress = *h.IpAddress
			}

			if h.MacAddress != nil {
				host.MacAddress = *h.MacAddress
			}

			if h.ReportingSource != nil {
				host.ReportingSource = *h.ReportingSource
			}

			if h.ReportingSourceInterface != nil {
				host.ReportingSourceInterface = *h.ReportingSourceInterface
			}

			response = append(response, host)
		}
	}

	return response
}

func getRPCZigbeeDevices(zigbeeDevices *[]inventory.ZigbeeDevice) []*rpc.ZigbeeDevice {
	var response []*rpc.ZigbeeDevice

	if zigbeeDevices != nil {
		for _, h := range *zigbeeDevices {
			zd := &pb.ZigbeeDevice{
				Dn:               h.Dn,
				Name:             h.Name,
				Description:      h.Description,
				IotZone:          h.IotZone,
				Type:             h.Type,
				SoftwareBuildId:  h.SoftwareBuildId,
				DateCode:         h.DateCode,
				Model:            h.Model,
				Vendor:           h.Vendor,
				ManufacturerName: h.ManufacturerName,
				PowerSource:      h.PowerSource,
				ModelId:          h.ModelId,
			}

			if h.LastSeen != nil {
				lastSeen, convError := ptypes.TimestampProto(*h.LastSeen)
				if convError != nil {
					log.Error(convError)
				}
				zd.LastSeen = lastSeen
			}

			response = append(response, zd)
		}
	}

	return response
}
