// Description:
// Author: liming.one@bytedance.com
package benchmarks

import (
	"gruptor"
	"testing"
)

func BenchmarkCursor_Load(b *testing.B) {
	cur := gruptor.NewCursor()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cur.Load()
	}
}

func BenchmarkCursor_Store(b *testing.B) {
	cur := gruptor.NewCursor()
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < times; i++ {
		cur.Store(i)
	}
}

func BenchmarkCursor_Read(b *testing.B) {
	cur := gruptor.NewCursor()
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	for i := int64(0); i < times; i++ {
		cur.Read(i)
	}
}
