// Description:
// Author: liming.one@bytedance.com
package gruptor

import "testing"

func BenchmarkCompositeBarrierWithOneCursor_Read(b *testing.B) {
	var barrier Barrier = NewCursor()
	//cBarrier := NewCompositeBarrier(NewCursor())
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	for i := int64(0); i < times; i++ {
		barrier.Read(0)
	}

}

func BenchmarkCompositeBarrierWithMoreCursor_Read(b *testing.B) {
	cBarrier := NewCompositeBarrier(NewCursor(), NewCursor(), NewCursor(), NewCursor(), NewCursor(), NewCursor())
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		cBarrier.Read(0)
	}
}
