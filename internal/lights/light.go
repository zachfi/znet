package lights

type Light interface {
	Alert(string) error
	Dim(string, int32) error
	Off(string) error
	On(string) error
	RandomColor(string, []string) error
	SetColor(string, string) error
	Toggle(string) error
}
