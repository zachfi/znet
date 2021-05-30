package lights

import (
	"errors"
)

// ErrNilConfig is used to indicate that a method has received a config that was nil.
var ErrNilConfig = errors.New("nil lights config")

// ErrNoRoomsConfigured is used to indicate that the configuration contained no Rooms.
var ErrNoRoomsConfigured = errors.New("no rooms configured")

// ErrRoomNotFound is used to indicate a named room was not found in the config.
var ErrRoomNotFound = errors.New("room not found")

// ErrUnknownActionEvent is used to indicate that an action was not recognized.
var ErrUnknownActionEvent = errors.New("unknown action event")

// ErrHandlerFailed is used to indicate that a lights handler has failed to execute.
var ErrHandlerFailed = errors.New("handler failed")

// ErrUnhandledEventName is used to indicate the evet name was not found in the RoomConfig.
var ErrUnhandledEventName = errors.New("unhandled event name")
