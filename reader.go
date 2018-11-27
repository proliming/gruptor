// Description: A port of Disruptor in golang
// Author: liming.one@bytedance.com
package gruptor

import (
	"runtime"
	"time"
)

// A reader read an event at readCursor and call eventHandler.Consume(lo,hi)
type Reader struct {
	eventProvider EventProvider
	readCursor    *Cursor
	writtenCursor *Cursor
	barrier       Barrier
	eventHandler  EventHandler
	running       bool
}

func NewReader(ep EventProvider, readCursor, writtenCursor *Cursor, barrier Barrier, eventHandler EventHandler) *Reader {
	return &Reader{
		eventProvider: ep,
		readCursor:    readCursor,
		writtenCursor: writtenCursor,
		barrier:       barrier,
		eventHandler:  eventHandler,
		running:       false,
	}
}

func (r *Reader) Start() {
	r.running = true
	go r.consume()
}

func (r *Reader) Stop() {
	r.running = false
}

// readCursor < writerCursor
// readCursor < dependentReaderCursor
// writerCursor < min(readCursors)
func (r *Reader) consume() {
	current := r.readCursor.Load()
	idling, gating := 0, 0
	for {
		next := current + 1
		maxRead := r.barrier.Read(next)
		if next <= maxRead {
			r.doConsume(next, maxRead)
			r.readCursor.Store(maxRead)
			current = maxRead
		} else if maxRead = r.writtenCursor.Load(); next <= maxRead {
			time.Sleep(time.Microsecond)
			// Gating--TODO: wait strategy (provide gating count to wait strategy for phased backoff)
			gating++
			idling = 0
		} else if r.running {
			time.Sleep(time.Millisecond)
			// Idling--TODO: wait strategy (provide idling count to wait strategy for phased backoff)
			idling++
			gating = 0
		} else {
			break
		}
		// sleeping increases the batch size which reduces number of writes required to store the sequence
		// reducing the number of writes allows the CPU to optimize the pipeline without prediction failures
		runtime.Gosched()
	}
}

func (r *Reader) doConsume(lo int64, hi int64) {
	for sequence := lo; sequence <= hi; sequence++ {
		e := r.eventProvider.Published(sequence)
		r.eventHandler.OnEvent(e, sequence)
	}
}
