package znet

// UnknownHost is the information about a host when the identity of the host is
// not known.
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
