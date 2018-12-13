# Gruptor
Gruptor is a basic port of Disruptor (JAVA) . It's very fast by using lock-free ringbuffer

# Usage

## Gruptor

```go

const DefaultBufferSize = 1024  // must be 2,4,8,..1024...
const DefaultBufferMask = DefaultBufferSize -1

var ringBuffer [DefaultBufferSize]int64  // specified an ringBuffer 

// implements your own Consumer
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

// Create the gruptor
g := gruptor.NewBuilder(DefaultBufferSize).HandleEventWith(&AConsumer{}).Build()
g.Start()
defer g.Stop()


// using Writer to acquiring next sequence
w := g.Writer()
var sequence int64

for sequence < 100000 {
	sequence = w.Next()
	ringBuffer[sequence&DefaultBufferMask] = sequence
	w.Commit(sequence, sequence)
}

```
### GruptorX

```go
// TODO add some example
```







```go

goos: darwin
goarch: amd64
pkg: gruptor/benchmarks
BenchmarkCompositeBarrierWithOneCursor_Read-8            	2000000000	         1.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompositeBarrierWithMoreCursor_Read-8           	200000000	         6.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkBlockingOneGoroutine-8                          	30000000	        55.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBlockingTwoGoroutines-8                         	20000000	        84.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBlockingThreeGoroutinesWithContendedWrite-8     	10000000	       167 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnBlockingOneGoroutine-8                        	100000000	        16.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnBlockingTwoGoroutines-8                       	100000000	        17.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnBlockingThreeGoroutinesWithContendedWrite-8   	100000000	        63.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkCursor_Load-8                                   	2000000000	         0.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkCursor_Store-8                                  	2000000000	         0.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkCursor_Read-8                                   	2000000000	         0.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkRingBuffer_Get-8                                	500000000	         2.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkRingBuffer_Published-8                          	500000000	         2.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkRingBuffer_Set-8                                	30000000	        52.0 ns/op	       8 B/op	       1 allocs/op
BenchmarkCustomRingBuffer_Get-8                          	1000000000	         2.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkCustomRingBuffer_Set-8                          	500000000	         3.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkCustomRingBufferWithEvent_Get-8                 	500000000	         3.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkCustomRingBufferWithEvent_Set-8                 	30000000	        55.7 ns/op	       8 B/op	       1 allocs/op
BenchmarkNoTypeSwitch-8                                  	2000000000	         1.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkUsingTypeSwitch-8                               	2000000000	         1.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkUsingTypeAssertion-8                            	1000000000	         2.29 ns/op	       0 B/op	       0 allocs/op

BenchmarkGruptor_OneWriterOneConsumer-8                      	100000000	        19.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkGruptor_OneWriterOneConsumerMoreCPU-8               	50000000	        32.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkGruptor_OneWriterMultiConsumer-8                    	50000000	        32.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGruptor_MultiWriterOneConsumer-8                    	10000000	       168 ns/op	       0 B/op	       0 allocs/op
BenchmarkGruptor_MultiWriterMultiConsumer-8                  	10000000	       184 ns/op	       0 B/op	       0 allocs/op
BenchmarkGruptor_MultiWriterOneConsumerInMultiGoroutines-8   	 5000000	       241 ns/op	       0 B/op	       0 allocs/op

```
