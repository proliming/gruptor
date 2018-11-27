// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

import "math"

// Coordination barrier for tracking the cursor for producers and sequence of
// dependent Consumers for processing data
type Barrier interface {
	Read(int64) int64
}

type CompositeBarrier []*Cursor

func NewCompositeBarrier(upstream ...*Cursor) CompositeBarrier {
	if len(upstream) == 0 {
		panic("At least one upstream cursor is required.")
	}
	cursors := make([]*Cursor, len(upstream))
	copy(cursors, upstream)
	return CompositeBarrier(cursors)
}

// Return the minimum sequence value of multiple consumers.
func (b CompositeBarrier) Read(noop int64) int64 {
	minimum := MaxSequenceValue
	for _, c := range b {
		sequence := c.Load()
		if sequence < minimum {
			minimum = sequence
		}
	}
	return minimum
}

//
type MultiWriterBarrier struct {
	written    *Cursor
	committed  []int32
	bufferSize int64
	mask       int64
	shift      uint8
}

func NewMultiWriterBarrier(written *Cursor, bufferSize int64) *MultiWriterBarrier {
	assertPowerOfTwo(bufferSize)

	return &MultiWriterBarrier{
		written:    written,
		committed:  prepareCommitBuffer(bufferSize),
		bufferSize: bufferSize,
		mask:       bufferSize - 1,
		shift:      uint8(math.Log2(float64(bufferSize))),
	}
}

func (b *MultiWriterBarrier) Read(lower int64) int64 {
	shift, mask := b.shift, b.mask
	upper := b.written.Load()

	for ; lower <= upper; lower++ {
		if b.committed[lower&mask] != int32(lower>>shift) {
			return lower - 1
		}
	}

	return upper
}

func prepareCommitBuffer(bufferSize int64) []int32 {
	buffer := make([]int32, bufferSize)
	for i := range buffer {
		buffer[i] = int32(InitialSequenceValue)
	}
	return buffer
}
