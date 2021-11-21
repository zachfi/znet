package lights

import "context"

// Handler is the interface to be implemented for a specific API.
type Handler interface {
	Alert(context.Context, string) error
	Dim(context.Context, string, int32) error
	Off(context.Context, string) error
	On(context.Context, string) error
	RandomColor(context.Context, string, []string) error
	SetColor(context.Context, string, string) error
	SetTemp(context.Context, string, int32) error
	Toggle(context.Context, string) error
}
