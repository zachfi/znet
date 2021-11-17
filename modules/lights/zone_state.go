package lights

type zoneState int

const (
	On zoneState = iota
	Off
	OffTimer
	Color
	RandomColor
	Dim
	NightVision
	EveningVision
	MorningVision
)

func (s zoneState) String() string {
	return [...]string{
		"On",
		"Off",
		"Color",
		"RandomColor",
		"Dim",
		"NightVision",
		"EveningVision",
		"MorningVision",
	}[s]
}
