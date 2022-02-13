package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/jan25/eventsproxy/event"
)

const (
	defaultUDPPort   = 7777
	defaultAdminPort = ":7778"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	// TODO(jan25): hookup more endpoints to view events buffered in proxy.

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: defaultUDPPort})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		msg := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(msg[:])
			if err != nil {
				log.Printf("failed to read from udp: %v", err)
				continue
			}

			e, err := event.ParseMessage(msg[:n])
			if err != nil {
				log.Printf("failed to parse message: %v", err)
				continue
			}

			out, _ := json.MarshalIndent(e, "", " ")
			log.Printf("event received: %s", string(out))
		}
	}()

	log.Printf("Staring adming server at %q", defaultAdminPort)
	log.Fatal(http.ListenAndServe(defaultAdminPort, nil))
}
