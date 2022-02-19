package main

import (
	"log"

	"github.com/jan25/eventsproxy/event"
)

// Reporter defines client side for a backend.
type Reporter interface {
	// Init initialises a reporter.
	// This is called just after reporter creation
	// during proxy startup.
	Init() error

	// Report reports events to backend.
	Report(events ...event.Event) error
}

// StdoutReporter prints events to STDOUT.
type StdoutReporter struct{}

func (s *StdoutReporter) Init() error {
	// nothing to initialise
	return nil
}

func (s *StdoutReporter) Report(events ...event.Event) error {
	for _, e := range events {
		log.Printf("%+v\n", e)
	}
	return nil
}

// TODO(jan25): implement Kafka reporter.
