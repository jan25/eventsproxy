package event

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

type Event struct {
	// Timestamp is timestamp when event started
	Timestamp time.Time `json:"timestamp"`

	// TimestampEnd is timestamp when event finished
	TimestampEnd time.Time `json:"timestamp_end,omitempty"`

	// Host is hostname of server where event occured
	Host string `json:"host"`

	// App is application in which event is triggered
	App string `json:"app"`

	// Endpoint is source endpoint for a event
	Endpoint string `json:"endpoint,omitempty"`

	// Method is method for a event
	Method string `json:"method,omitempty"`

	// StatusCode is HTTP status code for a event
	StatusCode int8 `json:"status_code,omitempty"`

	// Metrics is a key-value map of metric values.
	// e.g. { 'dbQueries': 10, 'cacheHits': 5 }
	Metrics map[string]string `json:"metrics,omitempty"`

	// Logs is list of logs during a event
	Logs []string `json:"logs,omitempty"`
}

func (e *Event) Text() []byte {
	var b bytes.Buffer
	s, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}
	var n int
	if n, err = b.Write(s); err != nil {
		log.Fatal(err)
	}
	if n != len(s) {
		log.Fatalf("Marshalled into mismatching sizes: %d vs %d", n, len(s))
	}
	return b.Bytes()
}

func ParseMessage(msg []byte) (*Event, error) {
	var e Event
	if err := json.Unmarshal(msg, &e); err != nil {
		return &Event{}, err
	}
	return &e, nil
}
