// Code generated, do not edit
package inventory

type MockInventory struct {
	FetchNetworkHostCalls     map[string]int
	FetchNetworkHostResponse  *NetworkHost
	FetchNetworkHostErr       error
	ListNetworkHostResponse   []NetworkHost
	ListNetworkHostErr        error
	CreateNetworkHostCalls    map[string]int
	FetchNetworkIDCalls       map[string]int
	FetchNetworkIDResponse    *NetworkID
	FetchNetworkIDErr         error
	ListNetworkIDResponse     []NetworkID
	ListNetworkIDErr          error
	CreateNetworkIDCalls      map[string]int
	FetchL3NetworkCalls       map[string]int
	FetchL3NetworkResponse    *L3Network
	FetchL3NetworkErr         error
	ListL3NetworkResponse     []L3Network
	ListL3NetworkErr          error
	CreateL3NetworkCalls      map[string]int
	FetchInetNetworkCalls     map[string]int
	FetchInetNetworkResponse  *InetNetwork
	FetchInetNetworkErr       error
	ListInetNetworkResponse   []InetNetwork
	ListInetNetworkErr        error
	CreateInetNetworkCalls    map[string]int
	FetchInet6NetworkCalls    map[string]int
	FetchInet6NetworkResponse *Inet6Network
	FetchInet6NetworkErr      error
	ListInet6NetworkResponse  []Inet6Network
	ListInet6NetworkErr       error
	CreateInet6NetworkCalls   map[string]int
	FetchZigbeeDeviceCalls    map[string]int
	FetchZigbeeDeviceResponse *ZigbeeDevice
	FetchZigbeeDeviceErr      error
	ListZigbeeDeviceResponse  []ZigbeeDevice
	ListZigbeeDeviceErr       error
	CreateZigbeeDeviceCalls   map[string]int
	FetchIOTZoneCalls         map[string]int
	FetchIOTZoneResponse      *IOTZone
	FetchIOTZoneErr           error
	ListIOTZoneResponse       []IOTZone
	ListIOTZoneErr            error
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

	if i.FetchNetworkHostErr != nil {
		return nil, i.FetchNetworkHostErr
	}

	return i.FetchNetworkHostResponse, nil
}

func (i *MockInventory) ListNetworkHosts() ([]NetworkHost, error) {

	if i.ListNetworkHostErr != nil {
		return nil, i.ListNetworkHostErr
	}

	return i.ListNetworkHostResponse, nil
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

	if i.FetchNetworkIDErr != nil {
		return nil, i.FetchNetworkIDErr
	}

	return i.FetchNetworkIDResponse, nil
}

func (i *MockInventory) ListNetworkIDs() ([]NetworkID, error) {

	if i.ListNetworkIDErr != nil {
		return nil, i.ListNetworkIDErr
	}

	return i.ListNetworkIDResponse, nil
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

	if i.FetchL3NetworkErr != nil {
		return nil, i.FetchL3NetworkErr
	}

	return i.FetchL3NetworkResponse, nil
}

func (i *MockInventory) ListL3Networks() ([]L3Network, error) {

	if i.ListL3NetworkErr != nil {
		return nil, i.ListL3NetworkErr
	}

	return i.ListL3NetworkResponse, nil
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

	if i.FetchInetNetworkErr != nil {
		return nil, i.FetchInetNetworkErr
	}

	return i.FetchInetNetworkResponse, nil
}

func (i *MockInventory) ListInetNetworks() ([]InetNetwork, error) {

	if i.ListInetNetworkErr != nil {
		return nil, i.ListInetNetworkErr
	}

	return i.ListInetNetworkResponse, nil
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

	if i.FetchInet6NetworkErr != nil {
		return nil, i.FetchInet6NetworkErr
	}

	return i.FetchInet6NetworkResponse, nil
}

func (i *MockInventory) ListInet6Networks() ([]Inet6Network, error) {

	if i.ListInet6NetworkErr != nil {
		return nil, i.ListInet6NetworkErr
	}

	return i.ListInet6NetworkResponse, nil
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

	if i.FetchZigbeeDeviceErr != nil {
		return nil, i.FetchZigbeeDeviceErr
	}

	return i.FetchZigbeeDeviceResponse, nil
}

func (i *MockInventory) ListZigbeeDevices() ([]ZigbeeDevice, error) {

	if i.ListZigbeeDeviceErr != nil {
		return nil, i.ListZigbeeDeviceErr
	}

	return i.ListZigbeeDeviceResponse, nil
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

	if i.FetchIOTZoneErr != nil {
		return nil, i.FetchIOTZoneErr
	}

	return i.FetchIOTZoneResponse, nil
}

func (i *MockInventory) ListIOTZones() ([]IOTZone, error) {

	if i.ListIOTZoneErr != nil {
		return nil, i.ListIOTZoneErr
	}

	return i.ListIOTZoneResponse, nil
}

func (i *MockInventory) UpdateIOTZone(*IOTZone) (*IOTZone, error) {

	return nil, nil
}
