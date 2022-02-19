package main

import (
	"log"
	"time"

	"github.com/jan25/eventsproxy/event"
)

type processor struct {
	incomingEvents <-chan *event.Event
	backend        Reporter

	buffer        []event.Event
	bufferIdx     int
	bufferSize    int
	flushInterval time.Duration
}

func NewProcessor(events <-chan *event.Event, backend Reporter, bufferSize int, flushInterval time.Duration) *processor {
	return &processor{
		incomingEvents: events,
		backend:        backend,
		bufferSize:     bufferSize,
		bufferIdx:      0,
		buffer:         make([]event.Event, bufferSize),
		flushInterval:  flushInterval,
	}
}

func (p *processor) Process() {
	t := time.NewTicker(p.flushInterval)
	defer t.Stop()

	for {
		select {
		case e := <-p.incomingEvents:
			p.handleEvent(e)
		case <-t.C:
			if p.bufferIdx != 0 {
				p.flush()
				t.Reset(p.flushInterval)
			}
		}
	}
}

func (p *processor) handleEvent(e *event.Event) {
	if p.bufferSize == p.bufferIdx {
		log.Println("buffer full. cant process new events")
		return
	}

	p.buffer[p.bufferIdx] = *e
	p.bufferIdx += 1
}

func (p *processor) flush() {
	p.backend.Report(p.buffer[:p.bufferIdx]...)
	p.bufferIdx = 0
}
