// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

// An event object, with all memory already allocated where possible
type Event interface{}

// Called by the RingBuffer to pre-populate all the events to fill the RingBuffer.
type EventFactory interface {
	NewEvent() Event
}

// Callback interface to be implemented for processing events as they become available in the RingBuffer
type EventHandler interface {
	OnEvent(event Event, sequence int64) error
}

type EventProvider interface {
	Published(sequence int64) Event
}
