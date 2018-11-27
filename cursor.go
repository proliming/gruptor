// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

// Concurrent sequence used for tracking the progress of the RingBuffer and Reader.  Support a number
// of concurrent operations including CAS.
type Cursor struct {
	lPadding [CpuCacheLinePadding]int64
	sequence int64                      // real sequence
	rPadding [CpuCacheLinePadding]int64 // right padding value to fill cache line
}

func NewCursor() *Cursor {
	return &Cursor{sequence: InitialSequenceValue}
}
