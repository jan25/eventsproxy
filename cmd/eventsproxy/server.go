package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"text/template"

	"github.com/jan25/eventsproxy/event"
)

type eventsServer struct {
	port       int
	bufferSize int
	events     chan event.Event
}

func (s *eventsServer) listenAndServe() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: defaultUDPPort})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msg := make([]byte, s.bufferSize)

	for {
		n, _, err := conn.ReadFromUDP(msg[:])
		if err != nil {
			log.Printf("failed to read from udp: %v", err)
			continue
		}

		event, err := event.ParseMessage(msg[:n])
		if err != nil {
			log.Printf("failed to parse message: %v", err)
			continue
		}

		out, _ := json.MarshalIndent(event, "", " ")
		log.Printf("event received: %s", string(out))
		s.events <- *event
	}
}

type metricKind int

const (
	RECEIVED metricKind = iota
	REPORTED
)

type metric struct {
	kind  metricKind
	value int // optional value
}

type adminServer struct {
	eventsReceived int
	eventsReported int

	metrics <-chan metric

	port int
}

const metricsTable = `<table>
	<tr><th>Metric</th><th>Value</th></tr>
	<tr><td>received</td><td>{{.Received}}</td></tr>
	<tr><td>reported</td><td>{{.Reported}}</td></tr>
</table>`

func (s *adminServer) listenAndServe() {
	go func() {
		for {
			metric := <-s.metrics
			if metric.kind == RECEIVED {
				s.eventsReceived++
			} else if metric.kind == REPORTED {
				s.eventsReported += metric.value
			}
		}
	}()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("metrics").Parse(metricsTable)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		t.Execute(w, struct {
			Received int
			Reported int
		}{
			Received: s.eventsReceived,
			Reported: s.eventsReported,
		})
	})

	log.Printf("Staring adming server at %d", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil))
}
