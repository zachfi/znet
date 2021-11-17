// Code generated, do not edit
package inventory

import "context"

type MockInventory struct {
	FetchNetworkHostCalls      map[string]int
	FetchNetworkHostResponse   *NetworkHost
	FetchNetworkHostErr        error
	ListNetworkHostResponse    []NetworkHost
	ListNetworkHostErr         error
	CreateNetworkHostCalls     map[string]int
	UpdateNetworkHostCalls     map[string]int
	UpdateNetworkHostResponse  *NetworkHost
	UpdateNetworkHostErr       error
	FetchNetworkIDCalls        map[string]int
	FetchNetworkIDResponse     *NetworkID
	FetchNetworkIDErr          error
	ListNetworkIDResponse      []NetworkID
	ListNetworkIDErr           error
	CreateNetworkIDCalls       map[string]int
	UpdateNetworkIDCalls       map[string]int
	UpdateNetworkIDResponse    *NetworkID
	UpdateNetworkIDErr         error
	FetchL3NetworkCalls        map[string]int
	FetchL3NetworkResponse     *L3Network
	FetchL3NetworkErr          error
	ListL3NetworkResponse      []L3Network
	ListL3NetworkErr           error
	CreateL3NetworkCalls       map[string]int
	UpdateL3NetworkCalls       map[string]int
	UpdateL3NetworkResponse    *L3Network
	UpdateL3NetworkErr         error
	FetchInetNetworkCalls      map[string]int
	FetchInetNetworkResponse   *InetNetwork
	FetchInetNetworkErr        error
	ListInetNetworkResponse    []InetNetwork
	ListInetNetworkErr         error
	CreateInetNetworkCalls     map[string]int
	UpdateInetNetworkCalls     map[string]int
	UpdateInetNetworkResponse  *InetNetwork
	UpdateInetNetworkErr       error
	FetchInet6NetworkCalls     map[string]int
	FetchInet6NetworkResponse  *Inet6Network
	FetchInet6NetworkErr       error
	ListInet6NetworkResponse   []Inet6Network
	ListInet6NetworkErr        error
	CreateInet6NetworkCalls    map[string]int
	UpdateInet6NetworkCalls    map[string]int
	UpdateInet6NetworkResponse *Inet6Network
	UpdateInet6NetworkErr      error
	FetchZigbeeDeviceCalls     map[string]int
	FetchZigbeeDeviceResponse  *ZigbeeDevice
	FetchZigbeeDeviceErr       error
	ListZigbeeDeviceResponse   []ZigbeeDevice
	ListZigbeeDeviceErr        error
	CreateZigbeeDeviceCalls    map[string]int
	UpdateZigbeeDeviceCalls    map[string]int
	UpdateZigbeeDeviceResponse *ZigbeeDevice
	UpdateZigbeeDeviceErr      error
	FetchIOTZoneCalls          map[string]int
	FetchIOTZoneResponse       *IOTZone
	FetchIOTZoneErr            error
	ListIOTZoneResponse        []IOTZone
	ListIOTZoneErr             error
	CreateIOTZoneCalls         map[string]int
	UpdateIOTZoneCalls         map[string]int
	UpdateIOTZoneResponse      *IOTZone
	UpdateIOTZoneErr           error
}

func (i *MockInventory) UpdateTimestamp(context.Context, string, string) error {

	return nil
}
func (i *MockInventory) CreateNetworkHost(_ context.Context, x *NetworkHost) (*NetworkHost, error) {
	if len(i.CreateNetworkHostCalls) == 0 {
		i.CreateNetworkHostCalls = make(map[string]int)
	}

	i.CreateNetworkHostCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchNetworkHost(_ context.Context, name string) (*NetworkHost, error) {
	if len(i.FetchNetworkHostCalls) == 0 {
		i.FetchNetworkHostCalls = make(map[string]int)
	}

	i.FetchNetworkHostCalls[name]++

	if i.FetchNetworkHostErr != nil {
		return nil, i.FetchNetworkHostErr
	}

	return i.FetchNetworkHostResponse, nil
}

func (i *MockInventory) ListNetworkHosts(_ context.Context) ([]NetworkHost, error) {

	if i.ListNetworkHostErr != nil {
		return nil, i.ListNetworkHostErr
	}

	return i.ListNetworkHostResponse, nil
}

func (i *MockInventory) UpdateNetworkHost(_ context.Context, x *NetworkHost) (*NetworkHost, error) {
	if len(i.UpdateNetworkHostCalls) == 0 {
		i.UpdateNetworkHostCalls = make(map[string]int)
	}

	i.UpdateNetworkHostCalls[x.Name]++

	if i.UpdateNetworkHostErr != nil {
		return nil, i.UpdateNetworkHostErr
	}

	return i.UpdateNetworkHostResponse, nil
}
func (i *MockInventory) CreateNetworkID(_ context.Context, x *NetworkID) (*NetworkID, error) {
	if len(i.CreateNetworkIDCalls) == 0 {
		i.CreateNetworkIDCalls = make(map[string]int)
	}

	i.CreateNetworkIDCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchNetworkID(_ context.Context, name string) (*NetworkID, error) {
	if len(i.FetchNetworkIDCalls) == 0 {
		i.FetchNetworkIDCalls = make(map[string]int)
	}

	i.FetchNetworkIDCalls[name]++

	if i.FetchNetworkIDErr != nil {
		return nil, i.FetchNetworkIDErr
	}

	return i.FetchNetworkIDResponse, nil
}

func (i *MockInventory) ListNetworkIDs(_ context.Context) ([]NetworkID, error) {

	if i.ListNetworkIDErr != nil {
		return nil, i.ListNetworkIDErr
	}

	return i.ListNetworkIDResponse, nil
}

func (i *MockInventory) UpdateNetworkID(_ context.Context, x *NetworkID) (*NetworkID, error) {
	if len(i.UpdateNetworkIDCalls) == 0 {
		i.UpdateNetworkIDCalls = make(map[string]int)
	}

	i.UpdateNetworkIDCalls[x.Name]++

	if i.UpdateNetworkIDErr != nil {
		return nil, i.UpdateNetworkIDErr
	}

	return i.UpdateNetworkIDResponse, nil
}
func (i *MockInventory) CreateL3Network(_ context.Context, x *L3Network) (*L3Network, error) {
	if len(i.CreateL3NetworkCalls) == 0 {
		i.CreateL3NetworkCalls = make(map[string]int)
	}

	i.CreateL3NetworkCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchL3Network(_ context.Context, name string) (*L3Network, error) {
	if len(i.FetchL3NetworkCalls) == 0 {
		i.FetchL3NetworkCalls = make(map[string]int)
	}

	i.FetchL3NetworkCalls[name]++

	if i.FetchL3NetworkErr != nil {
		return nil, i.FetchL3NetworkErr
	}

	return i.FetchL3NetworkResponse, nil
}

func (i *MockInventory) ListL3Networks(_ context.Context) ([]L3Network, error) {

	if i.ListL3NetworkErr != nil {
		return nil, i.ListL3NetworkErr
	}

	return i.ListL3NetworkResponse, nil
}

func (i *MockInventory) UpdateL3Network(_ context.Context, x *L3Network) (*L3Network, error) {
	if len(i.UpdateL3NetworkCalls) == 0 {
		i.UpdateL3NetworkCalls = make(map[string]int)
	}

	i.UpdateL3NetworkCalls[x.Name]++

	if i.UpdateL3NetworkErr != nil {
		return nil, i.UpdateL3NetworkErr
	}

	return i.UpdateL3NetworkResponse, nil
}
func (i *MockInventory) CreateInetNetwork(_ context.Context, x *InetNetwork) (*InetNetwork, error) {
	if len(i.CreateInetNetworkCalls) == 0 {
		i.CreateInetNetworkCalls = make(map[string]int)
	}

	i.CreateInetNetworkCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchInetNetwork(_ context.Context, name string) (*InetNetwork, error) {
	if len(i.FetchInetNetworkCalls) == 0 {
		i.FetchInetNetworkCalls = make(map[string]int)
	}

	i.FetchInetNetworkCalls[name]++

	if i.FetchInetNetworkErr != nil {
		return nil, i.FetchInetNetworkErr
	}

	return i.FetchInetNetworkResponse, nil
}

func (i *MockInventory) ListInetNetworks(_ context.Context) ([]InetNetwork, error) {

	if i.ListInetNetworkErr != nil {
		return nil, i.ListInetNetworkErr
	}

	return i.ListInetNetworkResponse, nil
}

func (i *MockInventory) UpdateInetNetwork(_ context.Context, x *InetNetwork) (*InetNetwork, error) {
	if len(i.UpdateInetNetworkCalls) == 0 {
		i.UpdateInetNetworkCalls = make(map[string]int)
	}

	i.UpdateInetNetworkCalls[x.Name]++

	if i.UpdateInetNetworkErr != nil {
		return nil, i.UpdateInetNetworkErr
	}

	return i.UpdateInetNetworkResponse, nil
}
func (i *MockInventory) CreateInet6Network(_ context.Context, x *Inet6Network) (*Inet6Network, error) {
	if len(i.CreateInet6NetworkCalls) == 0 {
		i.CreateInet6NetworkCalls = make(map[string]int)
	}

	i.CreateInet6NetworkCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchInet6Network(_ context.Context, name string) (*Inet6Network, error) {
	if len(i.FetchInet6NetworkCalls) == 0 {
		i.FetchInet6NetworkCalls = make(map[string]int)
	}

	i.FetchInet6NetworkCalls[name]++

	if i.FetchInet6NetworkErr != nil {
		return nil, i.FetchInet6NetworkErr
	}

	return i.FetchInet6NetworkResponse, nil
}

func (i *MockInventory) ListInet6Networks(_ context.Context) ([]Inet6Network, error) {

	if i.ListInet6NetworkErr != nil {
		return nil, i.ListInet6NetworkErr
	}

	return i.ListInet6NetworkResponse, nil
}

func (i *MockInventory) UpdateInet6Network(_ context.Context, x *Inet6Network) (*Inet6Network, error) {
	if len(i.UpdateInet6NetworkCalls) == 0 {
		i.UpdateInet6NetworkCalls = make(map[string]int)
	}

	i.UpdateInet6NetworkCalls[x.Name]++

	if i.UpdateInet6NetworkErr != nil {
		return nil, i.UpdateInet6NetworkErr
	}

	return i.UpdateInet6NetworkResponse, nil
}
func (i *MockInventory) CreateZigbeeDevice(_ context.Context, x *ZigbeeDevice) (*ZigbeeDevice, error) {
	if len(i.CreateZigbeeDeviceCalls) == 0 {
		i.CreateZigbeeDeviceCalls = make(map[string]int)
	}

	i.CreateZigbeeDeviceCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchZigbeeDevice(_ context.Context, name string) (*ZigbeeDevice, error) {
	if len(i.FetchZigbeeDeviceCalls) == 0 {
		i.FetchZigbeeDeviceCalls = make(map[string]int)
	}

	i.FetchZigbeeDeviceCalls[name]++

	if i.FetchZigbeeDeviceErr != nil {
		return nil, i.FetchZigbeeDeviceErr
	}

	return i.FetchZigbeeDeviceResponse, nil
}

func (i *MockInventory) ListZigbeeDevices(_ context.Context) ([]ZigbeeDevice, error) {

	if i.ListZigbeeDeviceErr != nil {
		return nil, i.ListZigbeeDeviceErr
	}

	return i.ListZigbeeDeviceResponse, nil
}

func (i *MockInventory) UpdateZigbeeDevice(_ context.Context, x *ZigbeeDevice) (*ZigbeeDevice, error) {
	if len(i.UpdateZigbeeDeviceCalls) == 0 {
		i.UpdateZigbeeDeviceCalls = make(map[string]int)
	}

	i.UpdateZigbeeDeviceCalls[x.Name]++

	if i.UpdateZigbeeDeviceErr != nil {
		return nil, i.UpdateZigbeeDeviceErr
	}

	return i.UpdateZigbeeDeviceResponse, nil
}
func (i *MockInventory) CreateIOTZone(_ context.Context, x *IOTZone) (*IOTZone, error) {
	if len(i.CreateIOTZoneCalls) == 0 {
		i.CreateIOTZoneCalls = make(map[string]int)
	}

	i.CreateIOTZoneCalls[x.Name]++

	return x, nil
}

func (i *MockInventory) FetchIOTZone(_ context.Context, name string) (*IOTZone, error) {
	if len(i.FetchIOTZoneCalls) == 0 {
		i.FetchIOTZoneCalls = make(map[string]int)
	}

	i.FetchIOTZoneCalls[name]++

	if i.FetchIOTZoneErr != nil {
		return nil, i.FetchIOTZoneErr
	}

	return i.FetchIOTZoneResponse, nil
}

func (i *MockInventory) ListIOTZones(_ context.Context) ([]IOTZone, error) {

	if i.ListIOTZoneErr != nil {
		return nil, i.ListIOTZoneErr
	}

	return i.ListIOTZoneResponse, nil
}

func (i *MockInventory) UpdateIOTZone(_ context.Context, x *IOTZone) (*IOTZone, error) {
	if len(i.UpdateIOTZoneCalls) == 0 {
		i.UpdateIOTZoneCalls = make(map[string]int)
	}

	i.UpdateIOTZoneCalls[x.Name]++

	if i.UpdateIOTZoneErr != nil {
		return nil, i.UpdateIOTZoneErr
	}

	return i.UpdateIOTZoneResponse, nil
}
