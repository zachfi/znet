// Code generated, do not edit
package inventory

// Inventory is the interface to implement for CRUD against a data store of network devices.
type Inventory interface {
	UpdateTimestamp(string, string) error

	CreateNetworkHost(*NetworkHost) (*NetworkHost, error)
	FetchNetworkHost(string) (*NetworkHost, error)
	ListNetworkHosts() ([]NetworkHost, error)
	UpdateNetworkHost(*NetworkHost) (*NetworkHost, error)
	CreateNetworkID(*NetworkID) (*NetworkID, error)
	FetchNetworkID(string) (*NetworkID, error)
	ListNetworkIDs() ([]NetworkID, error)
	UpdateNetworkID(*NetworkID) (*NetworkID, error)
	CreateL3Network(*L3Network) (*L3Network, error)
	FetchL3Network(string) (*L3Network, error)
	ListL3Networks() ([]L3Network, error)
	UpdateL3Network(*L3Network) (*L3Network, error)
	CreateInetNetwork(*InetNetwork) (*InetNetwork, error)
	FetchInetNetwork(string) (*InetNetwork, error)
	ListInetNetworks() ([]InetNetwork, error)
	UpdateInetNetwork(*InetNetwork) (*InetNetwork, error)
	CreateInet6Network(*Inet6Network) (*Inet6Network, error)
	FetchInet6Network(string) (*Inet6Network, error)
	ListInet6Networks() ([]Inet6Network, error)
	UpdateInet6Network(*Inet6Network) (*Inet6Network, error)
	CreateZigbeeDevice(*ZigbeeDevice) (*ZigbeeDevice, error)
	FetchZigbeeDevice(string) (*ZigbeeDevice, error)
	ListZigbeeDevices() ([]ZigbeeDevice, error)
	UpdateZigbeeDevice(*ZigbeeDevice) (*ZigbeeDevice, error)
	CreateIOTZone(*IOTZone) (*IOTZone, error)
	FetchIOTZone(string) (*IOTZone, error)
	ListIOTZones() ([]IOTZone, error)
	UpdateIOTZone(*IOTZone) (*IOTZone, error)
}
