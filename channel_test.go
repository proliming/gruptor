// Description:
// Author: liming.one@bytedance.com
package gruptor

import (
	"runtime"
	"testing"
)

const BenchmarkDefaultBufferSize = 1024 * 16

func BenchmarkBlockingOneGoroutine(b *testing.B) {
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(1)
	benchmarkBlocking(b, 1)
}

func BenchmarkBlockingTwoGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	benchmarkBlocking(b, 1)
}

func BenchmarkBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(1)
	benchmarkBlocking(b, 2)
}

func benchmarkBlocking(b *testing.B, writers int64) {
	times := int64(b.N)
	channel := make(chan int64, BenchmarkDefaultBufferSize)

	b.ReportAllocs()
	b.ResetTimer()

	for x := int64(0); x < writers; x++ {
		go func() {
			for i := int64(0); i < times; i++ {
				channel <- i
			}
		}()
	}

	for i := int64(0); i < times*writers; i++ {
		msg := <-channel
		if writers == 1 && msg != i {
			panic("Out of sequence")
		}
	}

	b.StopTimer()
}

//----------------------- Unblocking -------------------//
func BenchmarkUnBlockingOneGoroutine(b *testing.B) {
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(1)
	benchmarkUnBlocking(b, 1)
}

func BenchmarkUnBlockingTwoGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	benchmarkUnBlocking(b, 1)
}
func BenchmarkUnBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(1)
	benchmarkUnBlocking(b, 2)
}

func benchmarkUnBlocking(b *testing.B, writers int64) {
	times := int64(b.N)
	maxReads := times * writers
	channel := make(chan int64, DefaultBufferSize)

	b.ReportAllocs()
	b.ResetTimer()

	for x := int64(0); x < writers; x++ {
		go func() {
			for i := int64(0); i < times; {
				select {
				case channel <- i:
					i++
				default:
					continue
				}
			}
		}()
	}

	for i := int64(0); i < maxReads; i++ {
		select {
		case msg := <-channel:
			if writers == 1 && msg != i {
				// panic("Out of sequence")
			}
		default:
			continue
		}
	}

	b.StopTimer()
}
