// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

type Gruptor struct {
	bufferSize int64
	consumers  [][]Consumer
	cursors    []*Cursor
	writer     Writer
	readers    []*Reader
}

type Consumer interface {
	Consume(lo, hi int64)
}

func NewGruptor(bufferSize int64) *Gruptor {
	g := &Gruptor{
		bufferSize: bufferSize,
		cursors:    []*Cursor{NewCursor()},
	}
	return g
}

func (g *Gruptor) HandleEventWith(consumers ...Consumer) *Gruptor {
	target := make([]Consumer, len(consumers))
	copy(target, consumers)
	for i := 0; i < len(consumers); i++ {
		g.cursors = append(g.cursors, NewCursor())
	}
	g.consumers = append(g.consumers, target)
	return g
}

func (g *Gruptor) Build() *Gruptor {
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
	g.readers = allReaders
	g.writer = writer
	return g
}

func (g *Gruptor) BuildMultiWriter() *Gruptor {
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
	g.readers = allReaders
	g.writer = writer
	return g
}

func (g *Gruptor) buildReaders(csrIndex, cursorIndex int, writtenCursor *Cursor, barrier Barrier) ([]*Reader, Barrier) {
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

// Wrap of Writer.Commit(lo,hi)
func (g *Gruptor) Publish(sequence int64) {
	g.writer.Commit(sequence, sequence)
}

// Return the writer of this Gruptor
func (g *Gruptor) Writer() Writer {
	return g.writer
}

// Start all readers for consuming Event
func (g *Gruptor) Start() {
	if len(g.consumers) == 0 {
		panic("No event-handlers setup for Gruptor")
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
