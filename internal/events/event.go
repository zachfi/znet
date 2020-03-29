package events

// Event is the message that is passed on the channel after the RPC has decoded
// the payload bytes into an object.  The EventName here is taken from the RPC
// before being passed thorugh the channel.
type Event struct {
	Name    string
	Payload Payload
}
