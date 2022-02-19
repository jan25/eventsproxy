package main

import (
	"time"

	"github.com/jan25/eventsproxy/event"
)

const (
	defaultUDPPort   = 7777
	defaultAdminPort = 7778

	bufferSize           = 1024
	defaultFlushInterval = 5 * time.Second
)

func main() {
	incomingEvents := make(chan event.Event, bufferSize)
	metrics := make(chan metric, bufferSize)

	p := NewProcessor(incomingEvents, metrics, new(StdoutReporter), bufferSize, defaultFlushInterval)
	go p.Process()

	eventsServer := &eventsServer{port: defaultUDPPort, bufferSize: bufferSize, events: incomingEvents}
	adminServer := &adminServer{port: defaultAdminPort, metrics: metrics}

	go eventsServer.listenAndServe()
	adminServer.listenAndServe()
}
