// Description:
// Author: liming.one@bytedance.com
package gruptor

import "testing"

func BenchmarkCompositeBarrierWithOneCursor_Read(b *testing.B) {
	cBarrier := NewCompositeBarrier(NewCursor())
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		cBarrier.Read(i)
	}
}

func BenchmarkCompositeBarrierWithMoreCursor_Read(b *testing.B) {
	cBarrier := NewCompositeBarrier(NewCursor(), NewCursor(), NewCursor(), NewCursor(), NewCursor(), NewCursor())
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		cBarrier.Read(i)
	}
}
