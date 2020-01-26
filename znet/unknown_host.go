package znet

type UnknownHost struct {
	Name       string
	IP         string
	MACAddress string
}

var unknownHostDefaultAttributes = []string{
	"cn",
	"v4Address",
	"macAddress",
}
