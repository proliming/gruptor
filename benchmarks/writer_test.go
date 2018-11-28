// Description:
// Author: liming.one@bytedance.com
package benchmarks

import (
	"gruptor"
	"testing"
)

// method call will cost more
func BenchmarkWriterNext(b *testing.B) {
	read, written := gruptor.NewCursor(), gruptor.NewCursor()
	writer := gruptor.NewSingleWriter(written, read, 1024)
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < times; i++ {
		sequence := writer.Next()
		read.Store(sequence)
	}
}

func BenchmarkWriterNextN(b *testing.B) {
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
func BenchmarkWriterNextWrapPoint(b *testing.B) {
	read, written := gruptor.NewCursor(), gruptor.NewCursor()
	writer := gruptor.NewSingleWriter(written, read, 1024)
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	read.Store(gruptor.MaxSequenceValue)
	for i := int64(0); i < times; i++ {
		writer.NextN(1)
	}
}

func BenchmarkWriterAwait(b *testing.B) {
	written, read := gruptor.NewCursor(), gruptor.NewCursor()
	writer := gruptor.NewSingleWriter(written, read, 1024*64)
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < times; i++ {
		writer.Await(i)
		read.Store(i)
	}
}

func BenchmarkWriterCommit(b *testing.B) {
	writer := gruptor.NewSingleWriter(gruptor.NewCursor(), nil, 1024)
	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < times; i++ {
		writer.Commit(i, i)
	}
}
