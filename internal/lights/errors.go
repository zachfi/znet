package lights

import (
	"errors"
)

// ErrNilConfig is used to indicate that a method has received a config that was nil.
var ErrNilConfig = errors.New("nil config")

// ErrNoRoomsConfigured is used to indicate that the configuration contained no Rooms.
var ErrNoRoomsConfigured = errors.New("no rooms configured")

// ErrRoomNotFound is used to indicate a named room was not found in the config.
var ErrRoomNotFound = errors.New("room not found")

// ErrUnknownActionEvent is used to indicate that an action was not recognized.
var ErrUnknownActionEvent = errors.New("unknown action event")
