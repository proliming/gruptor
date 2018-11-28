// Description:
// Author: liming.one@bytedance.com
package benchmarks

import (
	"gruptor"
	"testing"
)

func BenchmarkCompositeBarrierWithOneCursor_Read(b *testing.B) {
	var barrier gruptor.Barrier = gruptor.NewCursor()
	//cBarrier := NewCompositeBarrier(gruptor.NewCursor())
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	for i := int64(0); i < times; i++ {
		barrier.Read(0)
	}

}

func BenchmarkCompositeBarrierWithMoreCursor_Read(b *testing.B) {
	cBarrier := gruptor.NewCompositeBarrier(gruptor.NewCursor(), gruptor.NewCursor(), gruptor.NewCursor(), gruptor.NewCursor(), gruptor.NewCursor(), gruptor.NewCursor())
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		cBarrier.Read(0)
	}
}
