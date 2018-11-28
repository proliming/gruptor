// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

// A extend Gruptor with more friendly API, but a little cost.
type XGruptor struct {
	bufferSize    int64
	eventFactory  EventFactory
	ringBuffer    *RingBuffer
	eventHandlers [][]EventHandler
	cursors       []*Cursor
	writer        Writer
	readers       []*XReader
}

func NewXGruptor(bufferSize int64, eventFactory EventFactory) *XGruptor {
	g := &XGruptor{
		eventFactory: eventFactory,
		bufferSize:   bufferSize,
		ringBuffer:   NewRingBuffer(bufferSize),
		cursors:      []*Cursor{NewCursor()},
	}
	g.fillBy(eventFactory)
	return g
}

func (g *XGruptor) HandleEventWith(eventHandlers ...EventHandler) *XGruptor {
	target := make([]EventHandler, len(eventHandlers))
	copy(target, eventHandlers)
	for i := 0; i < len(eventHandlers); i++ {
		g.cursors = append(g.cursors, NewCursor())
	}
	g.eventHandlers = append(g.eventHandlers, target)
	return g
}

func (g *XGruptor) Build() *XGruptor {
	var allReaders []*XReader
	writtenCursor := g.cursors[0]
	var barrier Barrier = g.cursors[0]
	cursorIndex := 1 // 0 index is reserved for the writer Cursor

	for ehIndex, eh := range g.eventHandlers {
		readers, readerBarrier := g.buildReaders(ehIndex, cursorIndex, writtenCursor, barrier)
		for _, r := range readers {
			allReaders = append(allReaders, r)
		}
		barrier = readerBarrier
		cursorIndex += len(eh)
	}
	writer := NewSingleWriter(writtenCursor, barrier, g.bufferSize)
	g.readers = allReaders
	g.writer = writer
	return g
}

func (g *XGruptor) BuildCurrent() *XGruptor {
	var allReaders []*XReader
	writtenCursor := g.cursors[0]
	writerBarrier := NewMultiWriterBarrier(writtenCursor, g.bufferSize)
	var barrier Barrier = writerBarrier
	cursorIndex := 1 // 0 index is reserved for the writer Cursor

	for ehIndex, eh := range g.eventHandlers {
		readers, readerBarrier := g.buildReaders(ehIndex, cursorIndex, writtenCursor, barrier)
		for _, r := range readers {
			allReaders = append(allReaders, r)
		}
		barrier = readerBarrier
		cursorIndex += len(eh)
	}
	writer := NewMultiWriter(writerBarrier, barrier)
	g.readers = allReaders
	g.writer = writer
	return g
}

func (g *XGruptor) buildReaders(ehIndex, cursorIndex int, writtenCursor *Cursor, barrier Barrier) ([]*XReader, Barrier) {
	var barrierCursors []*Cursor
	var readers []*XReader

	for _, eh := range g.eventHandlers[ehIndex] {
		readCursor := g.cursors[cursorIndex]
		barrierCursors = append(barrierCursors, readCursor)
		reader := NewXReader(g.ringBuffer, readCursor, writtenCursor, barrier, eh)
		readers = append(readers, reader)
		cursorIndex++
	}
	if len(g.eventHandlers[ehIndex]) == 1 {
		return readers, barrierCursors[0]
	} else {
		return readers, NewCompositeBarrier(barrierCursors...)
	}
}

// Wrap of Writer.Commit(lo,hi)
func (g *XGruptor) Publish(sequence int64) {
	g.writer.Commit(sequence, sequence)
}

// Direct publish an event, this method will replace the value at the specified sequence.
// Please try the-normal-way to avoid gc.
// sequence:=g.Writer().Next()
// event:=g.Get(sequence)
// event.Data=xxx
// g.Publish(sequence)
func (g *XGruptor) DirectPublish(e Event) {
	sequence := g.writer.Next()
	g.Set(sequence, e)
	g.writer.Commit(sequence, sequence)
}

// Return the writer of this XGruptor
func (g *XGruptor) Writer() Writer {
	return g.writer
}

// Return pre-filled Event of this XGruptor at the specified sequence
func (g *XGruptor) Get(sequence int64) Event {
	return g.ringBuffer.Get(sequence)
}

// Replace Event at the specified sequence by new Event
// This method may cause more gc.
func (g *XGruptor) Set(sequence int64, v Event) {
	g.ringBuffer.Set(sequence, v)
}

// Start all readers for consuming Event
func (g *XGruptor) Start() {
	if len(g.eventHandlers) == 0 {
		panic("No event-handlers setup for XGruptor")
	}
	if g.writer == nil {
		panic("No writer setup for XGruptor")
	}
	for _, r := range g.readers {
		r.Start()
	}
}

// Stop all readers.
func (g *XGruptor) Stop() {
	for _, r := range g.readers {
		r.Stop()
	}
}
func (g *XGruptor) fillBy(factory EventFactory) {
	for i := int64(0); i < g.ringBuffer.bufferSize; i++ {
		g.ringBuffer.buf[i] = factory.NewEvent()
	}
}
