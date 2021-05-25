package lights

import "fmt"

// ErrNilConfig is used to indicate that a method has received a config that was nil.
var ErrNilConfig = fmt.Errorf("nil config")

// ErrNoRoomsConfigured is used to indicate that the configuration contained no Rooms.
var ErrNoRoomsConfigured = fmt.Errorf("no rooms configured")
