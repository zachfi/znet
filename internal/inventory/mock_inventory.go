// Code generated, do not edit
package inventory

type MockInventory struct {
	FetchNetworkHostCalls     map[string]int
	FetchNetworkHostResponse  *NetworkHost
	FetchNetworkHostError     error
	ListNetworkHostResponse   []NetworkHost
	CreateNetworkHostCalls    map[string]int
	FetchNetworkIDCalls       map[string]int
	FetchNetworkIDResponse    *NetworkID
	FetchNetworkIDError       error
	ListNetworkIDResponse     []NetworkID
	CreateNetworkIDCalls      map[string]int
	FetchL3NetworkCalls       map[string]int
	FetchL3NetworkResponse    *L3Network
	FetchL3NetworkError       error
	ListL3NetworkResponse     []L3Network
	CreateL3NetworkCalls      map[string]int
	FetchInetNetworkCalls     map[string]int
	FetchInetNetworkResponse  *InetNetwork
	FetchInetNetworkError     error
	ListInetNetworkResponse   []InetNetwork
	CreateInetNetworkCalls    map[string]int
	FetchInet6NetworkCalls    map[string]int
	FetchInet6NetworkResponse *Inet6Network
	FetchInet6NetworkError    error
	ListInet6NetworkResponse  []Inet6Network
	CreateInet6NetworkCalls   map[string]int
	FetchZigbeeDeviceCalls    map[string]int
	FetchZigbeeDeviceResponse *ZigbeeDevice
	FetchZigbeeDeviceError    error
	ListZigbeeDeviceResponse  []ZigbeeDevice
	CreateZigbeeDeviceCalls   map[string]int
	FetchIOTZoneCalls         map[string]int
	FetchIOTZoneResponse      *IOTZone
	FetchIOTZoneError         error
	ListIOTZoneResponse       []IOTZone
	CreateIOTZoneCalls        map[string]int
}

func (i *MockInventory) UpdateTimestamp(string, string) error {

	return nil
}
func (i *MockInventory) CreateNetworkHost(x *NetworkHost) (*NetworkHost, error) {
	if len(i.CreateNetworkHostCalls) == 0 {
		i.CreateNetworkHostCalls = make(map[string]int)
	}

	i.CreateNetworkHostCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchNetworkHost(name string) (*NetworkHost, error) {
	if len(i.FetchNetworkHostCalls) == 0 {
		i.FetchNetworkHostCalls = make(map[string]int)
	}

	i.FetchNetworkHostCalls[name]++

	if i.FetchNetworkHostError != nil {
		return nil, i.FetchNetworkHostError
	}

	return i.FetchNetworkHostResponse, nil
}

func (i *MockInventory) ListNetworkHosts() ([]NetworkHost, error) {

	return nil, nil
}

func (i *MockInventory) UpdateNetworkHost(*NetworkHost) (*NetworkHost, error) {

	return nil, nil
}
func (i *MockInventory) CreateNetworkID(x *NetworkID) (*NetworkID, error) {
	if len(i.CreateNetworkIDCalls) == 0 {
		i.CreateNetworkIDCalls = make(map[string]int)
	}

	i.CreateNetworkIDCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchNetworkID(name string) (*NetworkID, error) {
	if len(i.FetchNetworkIDCalls) == 0 {
		i.FetchNetworkIDCalls = make(map[string]int)
	}

	i.FetchNetworkIDCalls[name]++

	if i.FetchNetworkIDError != nil {
		return nil, i.FetchNetworkIDError
	}

	return i.FetchNetworkIDResponse, nil
}

func (i *MockInventory) ListNetworkIDs() ([]NetworkID, error) {

	return nil, nil
}

func (i *MockInventory) UpdateNetworkID(*NetworkID) (*NetworkID, error) {

	return nil, nil
}
func (i *MockInventory) CreateL3Network(x *L3Network) (*L3Network, error) {
	if len(i.CreateL3NetworkCalls) == 0 {
		i.CreateL3NetworkCalls = make(map[string]int)
	}

	i.CreateL3NetworkCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchL3Network(name string) (*L3Network, error) {
	if len(i.FetchL3NetworkCalls) == 0 {
		i.FetchL3NetworkCalls = make(map[string]int)
	}

	i.FetchL3NetworkCalls[name]++

	if i.FetchL3NetworkError != nil {
		return nil, i.FetchL3NetworkError
	}

	return i.FetchL3NetworkResponse, nil
}

func (i *MockInventory) ListL3Networks() ([]L3Network, error) {

	return nil, nil
}

func (i *MockInventory) UpdateL3Network(*L3Network) (*L3Network, error) {

	return nil, nil
}
func (i *MockInventory) CreateInetNetwork(x *InetNetwork) (*InetNetwork, error) {
	if len(i.CreateInetNetworkCalls) == 0 {
		i.CreateInetNetworkCalls = make(map[string]int)
	}

	i.CreateInetNetworkCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchInetNetwork(name string) (*InetNetwork, error) {
	if len(i.FetchInetNetworkCalls) == 0 {
		i.FetchInetNetworkCalls = make(map[string]int)
	}

	i.FetchInetNetworkCalls[name]++

	if i.FetchInetNetworkError != nil {
		return nil, i.FetchInetNetworkError
	}

	return i.FetchInetNetworkResponse, nil
}

func (i *MockInventory) ListInetNetworks() ([]InetNetwork, error) {

	return nil, nil
}

func (i *MockInventory) UpdateInetNetwork(*InetNetwork) (*InetNetwork, error) {

	return nil, nil
}
func (i *MockInventory) CreateInet6Network(x *Inet6Network) (*Inet6Network, error) {
	if len(i.CreateInet6NetworkCalls) == 0 {
		i.CreateInet6NetworkCalls = make(map[string]int)
	}

	i.CreateInet6NetworkCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchInet6Network(name string) (*Inet6Network, error) {
	if len(i.FetchInet6NetworkCalls) == 0 {
		i.FetchInet6NetworkCalls = make(map[string]int)
	}

	i.FetchInet6NetworkCalls[name]++

	if i.FetchInet6NetworkError != nil {
		return nil, i.FetchInet6NetworkError
	}

	return i.FetchInet6NetworkResponse, nil
}

func (i *MockInventory) ListInet6Networks() ([]Inet6Network, error) {

	return nil, nil
}

func (i *MockInventory) UpdateInet6Network(*Inet6Network) (*Inet6Network, error) {

	return nil, nil
}
func (i *MockInventory) CreateZigbeeDevice(x *ZigbeeDevice) (*ZigbeeDevice, error) {
	if len(i.CreateZigbeeDeviceCalls) == 0 {
		i.CreateZigbeeDeviceCalls = make(map[string]int)
	}

	i.CreateZigbeeDeviceCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchZigbeeDevice(name string) (*ZigbeeDevice, error) {
	if len(i.FetchZigbeeDeviceCalls) == 0 {
		i.FetchZigbeeDeviceCalls = make(map[string]int)
	}

	i.FetchZigbeeDeviceCalls[name]++

	if i.FetchZigbeeDeviceError != nil {
		return nil, i.FetchZigbeeDeviceError
	}

	return i.FetchZigbeeDeviceResponse, nil
}

func (i *MockInventory) ListZigbeeDevices() ([]ZigbeeDevice, error) {

	return nil, nil
}

func (i *MockInventory) UpdateZigbeeDevice(*ZigbeeDevice) (*ZigbeeDevice, error) {

	return nil, nil
}
func (i *MockInventory) CreateIOTZone(x *IOTZone) (*IOTZone, error) {
	if len(i.CreateIOTZoneCalls) == 0 {
		i.CreateIOTZoneCalls = make(map[string]int)
	}

	i.CreateIOTZoneCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchIOTZone(name string) (*IOTZone, error) {
	if len(i.FetchIOTZoneCalls) == 0 {
		i.FetchIOTZoneCalls = make(map[string]int)
	}

	i.FetchIOTZoneCalls[name]++

	if i.FetchIOTZoneError != nil {
		return nil, i.FetchIOTZoneError
	}

	return i.FetchIOTZoneResponse, nil
}

func (i *MockInventory) ListIOTZones() ([]IOTZone, error) {

	return nil, nil
}

func (i *MockInventory) UpdateIOTZone(*IOTZone) (*IOTZone, error) {

	return nil, nil
}
