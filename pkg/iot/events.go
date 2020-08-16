package iot

var EventNames = []string{
	"Click",
}

type Click struct {
	Count  string
	Device string
	Zone   string
}
