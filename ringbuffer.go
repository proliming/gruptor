// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

const (
	DefaultBufferSize = 1024 * 16
	DefaultBufferMask = DefaultBufferSize - 1
)

// Ring based store of reusable Events containing the data representing
type RingBuffer struct {
	bufferSize int64
	bufferMask int64
	buf        []Event
}

func DefaultRingBuffer() *RingBuffer {
	return &RingBuffer{
		bufferSize: DefaultBufferSize,
		bufferMask: DefaultBufferMask,
		buf:        make([]Event, DefaultBufferSize),
	}
}

func NewRingBuffer(bufferSize int64) *RingBuffer {
	return &RingBuffer{
		bufferSize: bufferSize,
		bufferMask: bufferSize - 1,
		buf:        make([]Event, bufferSize),
	}
}

func (buf *RingBuffer) Set(sequence int64, e Event) {
	buf.buf[sequence&buf.bufferMask] = e
}

func (buf *RingBuffer) Get(sequence int64) Event {
	return buf.buf[sequence&buf.bufferMask]
}

func (buf *RingBuffer) Published(sequence int64) Event {
	return buf.Get(sequence)
}
