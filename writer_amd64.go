// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

func (w *SingleWriter) Commit(lower, upper int64) {
	w.writtenCursor.sequence = upper
}
