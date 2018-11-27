// Description:
// Author: liming.one@bytedance.com
package gruptor

import "sync/atomic"

func (c *Cursor) Store(sequence int64) {
	atomic.StoreInt64(&c.sequence, sequence)
}

func (c *Cursor) Load() int64 {
	return atomic.LoadInt64(&c.sequence)
}

func (c *Cursor) Read(noop int64) int64 {
	return atomic.LoadInt64(&c.sequence)
}
