package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jan25/eventsproxy/event"
)

const (
	defaultUDPPort   = 7777
	defaultAdminPort = ":7778"

	bufferSize           = 1024
	defaultFlushInterval = 5 * time.Second
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	// TODO(jan25): hookup more endpoints to view events buffered in proxy.

	events := make(chan *event.Event, bufferSize)
	eventsServer := &eventsServer{port: defaultUDPPort, bufferSize: bufferSize, events: events}
	go eventsServer.listenAndServe()

	p := NewProcessor(events, new(StdoutReporter), bufferSize, defaultFlushInterval)
	go p.Process()

	log.Printf("Staring adming server at %q", defaultAdminPort)
	log.Fatal(http.ListenAndServe(defaultAdminPort, nil))
}
