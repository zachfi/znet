package inventory

import (
	"time"
)

// GENERATED FROM rpc.proto

type IOTDevice struct {
	Name            string
	DeviceDiscovery *DeviceDiscovery
	NetworkID       *NetworkID
}

type NetworkHost struct {
	Role            string
	Group           string
	Name            string
	OperatingSystem string
	Platform        string
	Type            string
	Domain          string
	Description     string
	Watch           *bool
	InetAddress     *[]string
	Inet6Address    *[]string
	MacAddress      *[]string
	LastSeen        *time.Time
	Dn              string
}

type NetworkID struct {
	Name                     string
	MacAddress               *[]string
	IpAddress                *[]string
	ReportingSource          *[]string
	ReportingSourceInterface *[]string
	LastSeen                 *time.Time
	Dn                       string
}

type L3Network struct {
	Name         string
	Soa          string
	Domain       string
	NtpServers   *[]string
	InetNetwork  *InetNetwork
	Inet6Network *Inet6Network
	Dn           string
}

type InetNetwork struct {
	Name         string
	Prefix       string
	Gateway      string
	DynamicRange string
	Dn           string
}

type Inet6Network struct {
	Name    string
	Prefix  string
	Gateway string
	Dn      string
}

type DeviceDiscovery struct {
	DiscoveryPrefix string
	Component       string
	NodeId          string
	ObjectId        string
	Endpoint        *[]string
	Message         []byte
}
