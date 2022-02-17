package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/jan25/eventsproxy/event"
)

type eventsServer struct {
	port       int
	bufferSize int
	events     chan *event.Event
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
		s.events <- event
	}
}
