package iot

type DeviceType int

const (
	Unknown = iota
	Coordinator
	BasicLight
	ColorLight
	Relay
	Leak
	Button
	Motion
	Temperature
)

func (d DeviceType) String() string {
	switch d {
	case Unknown:
		return "Unknown"
	case Coordinator:
		return "Coordinator"
	case BasicLight:
		return "BasicLight"
	case ColorLight:
		return "ColorLight"
	case Relay:
		return "Relay"
	case Leak:
		return "Leak"
	case Button:
		return "Button"
	case Temperature:
		return "Temperature"
	}

	return "Unknown"
}
