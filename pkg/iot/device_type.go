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
