// Description:
// Author: liming.one@bytedance.com
package gruptor

import "testing"

func BenchmarkRingBuffer_Get(b *testing.B) {
	rb := DefaultRingBuffer()
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		rb.Get(i % DefaultBufferSize)
	}
}

func BenchmarkRingBuffer_Published(b *testing.B) {
	rb := DefaultRingBuffer()
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		rb.Published(i % DefaultBufferSize)
	}
}

func BenchmarkRingBuffer_Set(b *testing.B) {
	rb := DefaultRingBuffer()
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()

	for i := int64(0); i < times; i++ {
		rb.Set(i%DefaultBufferSize, i)
	}
}

func BenchmarkCustomRingBuffer_Get(b *testing.B) {
	rb := make([]int64, DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		x := rb[i%DefaultBufferSize&DefaultBufferMask]
		x = x - 0
	}
}
func BenchmarkCustomRingBuffer_Set(b *testing.B) {
	rb := make([]int64, DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		rb[i%DefaultBufferSize&DefaultBufferMask] = i
	}
}

func BenchmarkCustomRingBufferWithEvent_Get(b *testing.B) {
	rb := make([]Event, DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		rb[i%DefaultBufferSize&DefaultBufferMask] = nil
	}
}
func BenchmarkCustomRingBufferWithEvent_Set(b *testing.B) {
	rb := make([]Event, DefaultBufferSize)
	times := int64(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := int64(0); i < times; i++ {
		rb[i%DefaultBufferSize&DefaultBufferMask] = i
	}
}
