// Description:
// Author: liming.one@bytedance.com
package benchmarks

import (
	"gruptor"
	"testing"
)

type Writer interface {
	Next() int64
	NextN(n int64) int64
	Commit(lo, hi int64)
}

func BenchmarkWriterDirectNextN(b *testing.B) {
	read, written := gruptor.NewCursor(), gruptor.NewCursor()
	writer := gruptor.NewSingleWriter(written, read, 1024)
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < times; i++ {
		sequence := writer.NextN(1)
		read.Store(sequence)
	}
}

// method call will cost more
func BenchmarkWriterPolymorphismNextN(b *testing.B) {
	read := gruptor.NewCursor()
	writer := getWriter(read)
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < times; i++ {
		sequence := writer.NextN(1)
		read.Store(sequence)
	}
}

func getWriter(read *gruptor.Cursor) Writer {
	written := gruptor.NewCursor()
	return gruptor.NewSingleWriter(written, read, 1024)
}
