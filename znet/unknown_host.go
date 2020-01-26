package znet

type UnknownHost struct {
	IP         string
	MACAddress string
}

var unknownHostDefaultAttributes = []string{
	"cn",
	"v4Address",
	"macAddress",
}
