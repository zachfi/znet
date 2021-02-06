// Code generated, do not edit
package inventory

type MockInventory struct {
	FetchNetworkHostCalls     map[string]int
	FetchNetworkHostResponse  *NetworkHost
	ListNetworkHostResponse   []NetworkHost
	CreateNetworkHostCalls    []*NetworkHost
	FetchNetworkIDCalls       map[string]int
	FetchNetworkIDResponse    *NetworkID
	ListNetworkIDResponse     []NetworkID
	CreateNetworkIDCalls      []*NetworkID
	FetchL3NetworkCalls       map[string]int
	FetchL3NetworkResponse    *L3Network
	ListL3NetworkResponse     []L3Network
	CreateL3NetworkCalls      []*L3Network
	FetchInetNetworkCalls     map[string]int
	FetchInetNetworkResponse  *InetNetwork
	ListInetNetworkResponse   []InetNetwork
	CreateInetNetworkCalls    []*InetNetwork
	FetchInet6NetworkCalls    map[string]int
	FetchInet6NetworkResponse *Inet6Network
	ListInet6NetworkResponse  []Inet6Network
	CreateInet6NetworkCalls   []*Inet6Network
	FetchZigbeeDeviceCalls    map[string]int
	FetchZigbeeDeviceResponse *ZigbeeDevice
	ListZigbeeDeviceResponse  []ZigbeeDevice
	CreateZigbeeDeviceCalls   []*ZigbeeDevice
	FetchIOTZoneCalls         map[string]int
	FetchIOTZoneResponse      *IOTZone
	ListIOTZoneResponse       []IOTZone
	CreateIOTZoneCalls        []*IOTZone
}

func (i *MockInventory) UpdateTimestamp(string, string) error {

	return nil
}
func (i *MockInventory) CreateNetworkHost(x *NetworkHost) (*NetworkHost, error) {
	i.CreateNetworkHostCalls = append(i.CreateNetworkHostCalls, x)

	return x, nil
}

func (i *MockInventory) FetchNetworkHost(name string) (*NetworkHost, error) {
	if len(i.FetchNetworkHostCalls) == 0 {
		i.FetchNetworkHostCalls = make(map[string]int)
	}

	i.FetchNetworkHostCalls[name]++

	return i.FetchNetworkHostResponse, nil
}

func (i *MockInventory) ListNetworkHosts() ([]NetworkHost, error) {

	return nil, nil
}

func (i *MockInventory) UpdateNetworkHost(*NetworkHost) (*NetworkHost, error) {

	return nil, nil
}
func (i *MockInventory) CreateNetworkID(x *NetworkID) (*NetworkID, error) {
	i.CreateNetworkIDCalls = append(i.CreateNetworkIDCalls, x)

	return x, nil
}

func (i *MockInventory) FetchNetworkID(name string) (*NetworkID, error) {
	if len(i.FetchNetworkIDCalls) == 0 {
		i.FetchNetworkIDCalls = make(map[string]int)
	}

	i.FetchNetworkIDCalls[name]++

	return i.FetchNetworkIDResponse, nil
}

func (i *MockInventory) ListNetworkIDs() ([]NetworkID, error) {

	return nil, nil
}

func (i *MockInventory) UpdateNetworkID(*NetworkID) (*NetworkID, error) {

	return nil, nil
}
func (i *MockInventory) CreateL3Network(x *L3Network) (*L3Network, error) {
	i.CreateL3NetworkCalls = append(i.CreateL3NetworkCalls, x)

	return x, nil
}

func (i *MockInventory) FetchL3Network(name string) (*L3Network, error) {
	if len(i.FetchL3NetworkCalls) == 0 {
		i.FetchL3NetworkCalls = make(map[string]int)
	}

	i.FetchL3NetworkCalls[name]++

	return i.FetchL3NetworkResponse, nil
}

func (i *MockInventory) ListL3Networks() ([]L3Network, error) {

	return nil, nil
}

func (i *MockInventory) UpdateL3Network(*L3Network) (*L3Network, error) {

	return nil, nil
}
func (i *MockInventory) CreateInetNetwork(x *InetNetwork) (*InetNetwork, error) {
	i.CreateInetNetworkCalls = append(i.CreateInetNetworkCalls, x)

	return x, nil
}

func (i *MockInventory) FetchInetNetwork(name string) (*InetNetwork, error) {
	if len(i.FetchInetNetworkCalls) == 0 {
		i.FetchInetNetworkCalls = make(map[string]int)
	}

	i.FetchInetNetworkCalls[name]++

	return i.FetchInetNetworkResponse, nil
}

func (i *MockInventory) ListInetNetworks() ([]InetNetwork, error) {

	return nil, nil
}

func (i *MockInventory) UpdateInetNetwork(*InetNetwork) (*InetNetwork, error) {

	return nil, nil
}
func (i *MockInventory) CreateInet6Network(x *Inet6Network) (*Inet6Network, error) {
	i.CreateInet6NetworkCalls = append(i.CreateInet6NetworkCalls, x)

	return x, nil
}

func (i *MockInventory) FetchInet6Network(name string) (*Inet6Network, error) {
	if len(i.FetchInet6NetworkCalls) == 0 {
		i.FetchInet6NetworkCalls = make(map[string]int)
	}

	i.FetchInet6NetworkCalls[name]++

	return i.FetchInet6NetworkResponse, nil
}

func (i *MockInventory) ListInet6Networks() ([]Inet6Network, error) {

	return nil, nil
}

func (i *MockInventory) UpdateInet6Network(*Inet6Network) (*Inet6Network, error) {

	return nil, nil
}
func (i *MockInventory) CreateZigbeeDevice(x *ZigbeeDevice) (*ZigbeeDevice, error) {
	i.CreateZigbeeDeviceCalls = append(i.CreateZigbeeDeviceCalls, x)

	return x, nil
}

func (i *MockInventory) FetchZigbeeDevice(name string) (*ZigbeeDevice, error) {
	if len(i.FetchZigbeeDeviceCalls) == 0 {
		i.FetchZigbeeDeviceCalls = make(map[string]int)
	}

	i.FetchZigbeeDeviceCalls[name]++

	return i.FetchZigbeeDeviceResponse, nil
}

func (i *MockInventory) ListZigbeeDevices() ([]ZigbeeDevice, error) {

	return nil, nil
}

func (i *MockInventory) UpdateZigbeeDevice(*ZigbeeDevice) (*ZigbeeDevice, error) {

	return nil, nil
}
func (i *MockInventory) CreateIOTZone(x *IOTZone) (*IOTZone, error) {
	i.CreateIOTZoneCalls = append(i.CreateIOTZoneCalls, x)

	return x, nil
}

func (i *MockInventory) FetchIOTZone(name string) (*IOTZone, error) {
	if len(i.FetchIOTZoneCalls) == 0 {
		i.FetchIOTZoneCalls = make(map[string]int)
	}

	i.FetchIOTZoneCalls[name]++

	return i.FetchIOTZoneResponse, nil
}

func (i *MockInventory) ListIOTZones() ([]IOTZone, error) {

	return nil, nil
}

func (i *MockInventory) UpdateIOTZone(*IOTZone) (*IOTZone, error) {

	return nil, nil
}
