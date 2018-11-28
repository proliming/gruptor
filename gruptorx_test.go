// Description:
// Author: liming.one@bytedance.com
package gruptor

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type AnEventHandler struct {
}

func (eh *AnEventHandler) OnEvent(event Event, sequence int64) error {
	for lo <= hi {
		event := ringBuffer[lo&DefaultBufferMask]
		if event != lo {
			warning := fmt.Sprintf("\nRace condition--Sequence: %d, Event: %d\n", lo, event)
			fmt.Printf(warning)
			panic(warning)
		}
		lo++
	}
}

func BenchmarkGruptorX_OneWriterOneConsumer(b *testing.B) {
	defer time.Sleep(time.Millisecond)

	runtime.GOMAXPROCS(1)

	g := NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	sequence := InitialSequenceValue

	w := g.Writer()

	b.ReportAllocs()
	b.ResetTimer()

	for sequence < times {
		sequence = w.Next()
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}
func BenchmarkGruptorX_OneWriterOneConsumerMoreCPU(b *testing.B) {
	defer time.Sleep(time.Millisecond)
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(1)
	g := NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()
	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.Next()
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}
func BenchmarkGruptorX_OneWriterMultiConsumer(b *testing.B) {
	defer time.Sleep(time.Millisecond)
	runtime.GOMAXPROCS(1)
	g := NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}, &AConsumer{}, &AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()
	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.Next()
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkGruptorX_MultiWriterOneConsumer(b *testing.B) {
	defer time.Sleep(time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	g := NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()

	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.Next()
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkGruptorX_MultiWriterMultiConsumer(b *testing.B) {
	defer time.Sleep(time.Millisecond)
	runtime.GOMAXPROCS(1)
	g := NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}, &AConsumer{}, &AConsumer{}).BuildConcurrent()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()

	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.Next()
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkGruptorX_MultiWriterOneConsumerInMultiGoroutines(b *testing.B) {
	defer time.Sleep(time.Millisecond)
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(1)
	g := NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}).BuildConcurrent()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			sequence := InitialSequenceValue
			for sequence < times {
				sequence = w.Next()
				ringBuffer[sequence&DefaultBufferMask] = sequence
				w.Commit(sequence, sequence)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	b.StopTimer()
}
