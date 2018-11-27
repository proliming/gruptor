// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

import (
	"runtime"
	"sync/atomic"
)

const SpinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values

type Writer interface {
	Next() int64

	NextN(n int64) int64

	Commit(lo, hi int64)
}

type SingleWriter struct {
	writtenCursor *Cursor
	barrier       Barrier
	bufferSize    int64
	next          int64
	gate          int64
}

func NewSingleWriter(writtenCursor *Cursor, barrier Barrier, bufferSize int64) *SingleWriter {
	assertPowerOfTwo(bufferSize)
	return &SingleWriter{
		barrier:       barrier,
		writtenCursor: writtenCursor,
		bufferSize:    bufferSize,
		next:          InitialSequenceValue,
		gate:          InitialSequenceValue,
	}
}

func (w *SingleWriter) Next() int64 {
	return w.NextN(1)
}

func (w *SingleWriter) NextN(n int64) int64 {
	w.next += n
	for spin := int64(0); w.next-w.bufferSize > w.gate; spin++ {
		if spin&SpinMask == 0 {
			runtime.Gosched()
		}
		w.gate = w.barrier.Read(0)
	}
	return w.next
}

func (w *SingleWriter) Await(next int64) {
	for next-w.bufferSize > w.gate {
		w.gate = w.barrier.Read(0)
	}
}

func (w *SingleWriter) Commit(lower, upper int64) {
	w.writtenCursor.Store(upper)
}

type MultiWriter struct {
	writtenCursor *Cursor
	upstream      Barrier
	bufferSize    int64
	gate          *Cursor
	indexMask     int64
	indexShift    uint8
	committed     []int32
}

func NewMultiWriter(barrier *MultiWriterBarrier, upstream Barrier) *MultiWriter {
	return &MultiWriter{
		writtenCursor: barrier.written,
		upstream:      upstream,
		bufferSize:    barrier.bufferSize,
		gate:          NewCursor(),
		indexMask:     barrier.mask,
		indexShift:    barrier.shift,
		committed:     barrier.committed,
	}
}

func (w *MultiWriter) Next() int64 {
	return w.NextN(1)
}

func (w *MultiWriter) NextN(n int64) int64 {
	for {
		next := w.writtenCursor.Load()
		upper := next + n

		for spin := int64(0); upper-w.bufferSize > w.gate.Load(); spin++ {
			if spin&SpinMask == 0 {
				runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
			}
			w.gate.Store(w.upstream.Read(0))
		}

		if atomic.CompareAndSwapInt64(&w.writtenCursor.sequence, next, upper) {
			return upper
		}
	}
}

func (w *MultiWriter) Commit(lower, upper int64) {
	if lower > upper {
		panic("Attempting to publish a sequence where the lower reservation is greater than the higher reservation.")
	} else if (upper - lower) > w.bufferSize {
		panic("Attempting to publish a reservation larger than the size of the ring buffer. (upper-lower > w.bufferSize)")
	} else if lower == upper {
		w.committed[upper&w.indexMask] = int32(upper >> w.indexShift)
	} else {
		// working down the array rather than up keeps all items in the commit together
		// otherwise the reader(s) could split up the group
		for upper >= lower {
			w.committed[upper&w.indexMask] = int32(upper >> w.indexShift)
			upper--
		}

	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		panic("The ring bufferSize must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}
