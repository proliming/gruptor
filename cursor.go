// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

import "sync/atomic"

// Concurrent sequence used for tracking the progress of the RingBuffer and Reader.  Support a number
// of concurrent operations including CAS.
type Cursor struct {
	lPadding [CpuCacheLinePadding]int64 // left padding value to fill cache line
	sequence int64                      // real sequence
	rPadding [CpuCacheLinePadding]int64 // right padding value to fill cache line
}

func NewCursor() *Cursor {
	return &Cursor{sequence: InitialSequenceValue}
}
func (c *Cursor) Store(sequence int64) {
	atomic.StoreInt64(&c.sequence, sequence)
}

func (c *Cursor) Load() int64 {
	return atomic.LoadInt64(&c.sequence)
}

func (c *Cursor) Read(noop int64) int64 {
	return atomic.LoadInt64(&c.sequence)
}
