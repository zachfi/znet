// Code generated, do not edit
package inventory

import "context"

// Inventory is the interface to implement for CRUD against a data store of network devices.
type Inventory interface {
	UpdateTimestamp(context.Context, string, string) error

	CreateNetworkHost(context.Context, *NetworkHost) (*NetworkHost, error)
	FetchNetworkHost(context.Context, string) (*NetworkHost, error)
	ListNetworkHosts(context.Context) ([]NetworkHost, error)
	UpdateNetworkHost(context.Context, *NetworkHost) (*NetworkHost, error)
	CreateNetworkID(context.Context, *NetworkID) (*NetworkID, error)
	FetchNetworkID(context.Context, string) (*NetworkID, error)
	ListNetworkIDs(context.Context) ([]NetworkID, error)
	UpdateNetworkID(context.Context, *NetworkID) (*NetworkID, error)
	CreateL3Network(context.Context, *L3Network) (*L3Network, error)
	FetchL3Network(context.Context, string) (*L3Network, error)
	ListL3Networks(context.Context) ([]L3Network, error)
	UpdateL3Network(context.Context, *L3Network) (*L3Network, error)
	CreateInetNetwork(context.Context, *InetNetwork) (*InetNetwork, error)
	FetchInetNetwork(context.Context, string) (*InetNetwork, error)
	ListInetNetworks(context.Context) ([]InetNetwork, error)
	UpdateInetNetwork(context.Context, *InetNetwork) (*InetNetwork, error)
	CreateInet6Network(context.Context, *Inet6Network) (*Inet6Network, error)
	FetchInet6Network(context.Context, string) (*Inet6Network, error)
	ListInet6Networks(context.Context) ([]Inet6Network, error)
	UpdateInet6Network(context.Context, *Inet6Network) (*Inet6Network, error)
	CreateZigbeeDevice(context.Context, *ZigbeeDevice) (*ZigbeeDevice, error)
	FetchZigbeeDevice(context.Context, string) (*ZigbeeDevice, error)
	ListZigbeeDevices(context.Context) ([]ZigbeeDevice, error)
	UpdateZigbeeDevice(context.Context, *ZigbeeDevice) (*ZigbeeDevice, error)
	CreateIOTZone(context.Context, *IOTZone) (*IOTZone, error)
	FetchIOTZone(context.Context, string) (*IOTZone, error)
	ListIOTZones(context.Context) ([]IOTZone, error)
	UpdateIOTZone(context.Context, *IOTZone) (*IOTZone, error)
}
