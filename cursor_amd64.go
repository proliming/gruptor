// Description:
// Author: liming.one@bytedance.com
package gruptor

func (c *Cursor) Store(sequence int64) {
	c.sequence = sequence
}

func (c *Cursor) Load() int64 {
	return c.sequence
}

func (c *Cursor) Read(noop int64) int64 {
	return c.sequence
}
