// Description:
// Author: liming.one@bytedance.com
package benchmarks

import (
	"gruptor"
	"testing"
)

func BenchmarkRingBuffer_Get(b *testing.B) {
	rb := gruptor.DefaultRingBuffer()
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		rb.Get(i % gruptor.DefaultBufferSize)
	}
}

func BenchmarkRingBuffer_Published(b *testing.B) {
	rb := gruptor.DefaultRingBuffer()
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		rb.Published(i % gruptor.DefaultBufferSize)
	}
}

func BenchmarkRingBuffer_Set(b *testing.B) {
	rb := gruptor.DefaultRingBuffer()
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		rb.Set(i%gruptor.DefaultBufferSize, i)
	}
}

func BenchmarkCustomRingBuffer_Get(b *testing.B) {
	rb := make([]int64, gruptor.DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		x := rb[i%gruptor.DefaultBufferSize&gruptor.DefaultBufferMask]
		x = x - 0
	}
}
func BenchmarkCustomRingBuffer_Set(b *testing.B) {
	rb := make([]int64, gruptor.DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		rb[i%gruptor.DefaultBufferSize&gruptor.DefaultBufferMask] = i
	}
}

func BenchmarkCustomRingBufferWithEvent_Get(b *testing.B) {
	rb := make([]gruptor.Event, gruptor.DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		rb[i%gruptor.DefaultBufferSize&gruptor.DefaultBufferMask] = nil
	}
}
func BenchmarkCustomRingBufferWithEvent_Set(b *testing.B) {
	rb := make([]gruptor.Event, gruptor.DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		rb[i%gruptor.DefaultBufferSize&gruptor.DefaultBufferMask] = i
	}
}
