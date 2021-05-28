package lights

type zoneState int

const (
	On zoneState = iota
	Off
	Color
	RandomColor
	Dim
)

func (s zoneState) String() string {
	return [...]string{
		"On",
		"Off",
		"Color",
		"RandomColor",
		"Dim",
	}[s]
}
