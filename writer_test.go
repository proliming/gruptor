// Description:
// Author: liming.one@bytedance.com
package gruptor

import "testing"

func BenchmarkSingleWriter_Next(b *testing.B) {
	read, written := NewCursor(), NewCursor()
	writer := NewSingleWriter(written, read, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence := writer.Next()
		read.Store(sequence)
	}
}
