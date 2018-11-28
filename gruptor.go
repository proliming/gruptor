// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	CpuCacheLinePadding        = 7
)

type Gruptor struct {
	writer  *SingleWriter
	readers []*Reader
}

type ConcurrentGruptor struct {
	writer  *MultiWriter
	readers []*Reader
}

type Consumer interface {
	Consume(lo, hi int64)
}

type gruptorBuilder struct {
	bufferSize int64
	consumers  [][]Consumer
	cursors    []*Cursor
}

func NewBuilder(bufferSize int64) *gruptorBuilder {
	return &gruptorBuilder{
		bufferSize: bufferSize,
		cursors:    []*Cursor{NewCursor()},
	}
}
func (g *gruptorBuilder) HandleEventWith(consumers ...Consumer) *gruptorBuilder {
	target := make([]Consumer, len(consumers))
	copy(target, consumers)
	for i := 0; i < len(consumers); i++ {
		g.cursors = append(g.cursors, NewCursor())
	}
	g.consumers = append(g.consumers, target)
	return g
}

func (g *gruptorBuilder) Build() *Gruptor {
	var allReaders []*Reader
	writtenCursor := g.cursors[0]

	var barrier Barrier = g.cursors[0]
	cursorIndex := 1 // 0 index is reserved for the writer Cursor
	for csrIndex, csr := range g.consumers {
		readers, readerBarrier := g.buildReaders(csrIndex, cursorIndex, writtenCursor, barrier)
		for _, r := range readers {
			allReaders = append(allReaders, r)
		}
		barrier = readerBarrier
		cursorIndex += len(csr)
	}
	writer := NewSingleWriter(writtenCursor, barrier, g.bufferSize)
	return &Gruptor{
		writer:  writer,
		readers: allReaders,
	}
}

func (g *gruptorBuilder) BuildConcurrent() *ConcurrentGruptor {
	var allReaders []*Reader
	writtenCursor := g.cursors[0]

	writerBarrier := NewMultiWriterBarrier(writtenCursor, g.bufferSize)
	var barrier Barrier = writerBarrier
	cursorIndex := 1 // 0 index is reserved for the writer Cursor

	for csrIndex, csr := range g.consumers {
		readers, readerBarrier := g.buildReaders(csrIndex, cursorIndex, writtenCursor, barrier)
		for _, r := range readers {
			allReaders = append(allReaders, r)
		}
		barrier = readerBarrier
		cursorIndex += len(csr)
	}
	writer := NewMultiWriter(writerBarrier, barrier)

	return &ConcurrentGruptor{
		writer:  writer,
		readers: allReaders,
	}
}

func (g *gruptorBuilder) buildReaders(csrIndex, cursorIndex int, writtenCursor *Cursor, barrier Barrier) ([]*Reader, Barrier) {
	var barrierCursors []*Cursor
	var readers []*Reader

	for _, csr := range g.consumers[csrIndex] {
		readCursor := g.cursors[cursorIndex]
		barrierCursors = append(barrierCursors, readCursor)
		reader := NewReader(readCursor, writtenCursor, barrier, csr)
		readers = append(readers, reader)
		cursorIndex++
	}
	if len(g.consumers[csrIndex]) == 1 {
		return readers, barrierCursors[0]
	} else {
		return readers, NewCompositeBarrier(barrierCursors...)
	}
}

// Return the writer of this Gruptor
func (g *Gruptor) Writer() *SingleWriter {
	return g.writer
}

// Start all readers for consuming Event
func (g *Gruptor) Start() {
	if len(g.readers) == 0 {
		panic("No readers setup for Gruptor")
	}
	if g.writer == nil {
		panic("No writer setup for Gruptor")
	}
	for _, r := range g.readers {
		r.Start()
	}
}

// Stop all readers.
func (g *Gruptor) Stop() {
	for _, r := range g.readers {
		r.Stop()
	}
}

// Return the writer of this Gruptor
func (g *ConcurrentGruptor) Writer() *MultiWriter {
	return g.writer
}

// Start all readers for consuming Event
func (g *ConcurrentGruptor) Start() {
	if len(g.readers) == 0 {
		panic("No readers setup for Gruptor")
	}
	if g.writer == nil {
		panic("No writer setup for Gruptor")
	}
	for _, r := range g.readers {
		r.Start()
	}
}

// Stop all readers.
func (g *ConcurrentGruptor) Stop() {
	for _, r := range g.readers {
		r.Stop()
	}
}
