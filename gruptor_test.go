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

var ringBuffer [DefaultBufferSize]int64

type AConsumer struct {
}

func (c *AConsumer) Consume(lo, hi int64) {
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

func BenchmarkGruptor_OneWriterOneConsumer(b *testing.B) {
	defer time.Sleep(10 * time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	g := NewGruptor(DefaultBufferSize).HandleEventWith(&AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	sequence := InitialSequenceValue

	w := g.Writer()

	b.ReportAllocs()
	b.ResetTimer()

	for sequence < times {
		sequence = w.NextN(1)
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}
func BenchmarkGruptor_OneWriterOneConsumerMoreCPU(b *testing.B) {
	defer time.Sleep(10 * time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	g := NewGruptor(DefaultBufferSize).HandleEventWith(&AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()
	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.NextN(1)
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}
func BenchmarkGruptor_OneWriterMultiConsumer(b *testing.B) {
	defer time.Sleep(10 * time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	g := NewGruptor(DefaultBufferSize).HandleEventWith(&AConsumer{}, &AConsumer{}, &AConsumer{}).Build()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()
	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.NextN(1)
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkGruptor_MultiWriterOneConsumer(b *testing.B) {
	defer time.Sleep(10 * time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	g := NewGruptor(DefaultBufferSize).HandleEventWith(&AConsumer{}).BuildMultiWriter()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()

	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.NextN(1)
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkGruptor_MultiWriterMultiConsumer(b *testing.B) {
	defer time.Sleep(10 * time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	g := NewGruptor(DefaultBufferSize).HandleEventWith(&AConsumer{}, &AConsumer{}, &AConsumer{}).BuildMultiWriter()
	g.Start()
	defer g.Stop()

	times := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	w := g.Writer()

	sequence := InitialSequenceValue
	for sequence < times {
		sequence = w.NextN(1)
		ringBuffer[sequence&DefaultBufferMask] = sequence
		w.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkGruptor_MultiWriterOneConsumerInMultiGoroutines(b *testing.B) {
	defer time.Sleep(10 * time.Millisecond)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	g := NewGruptor(DefaultBufferSize).HandleEventWith(&AConsumer{}).BuildMultiWriter()
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
				sequence = w.NextN(1)
				ringBuffer[sequence&DefaultBufferMask] = sequence
				w.Commit(sequence, sequence)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	b.StopTimer()
}
