package main

import (
	"log"
	"net"
	"os/user"
	"time"

	"github.com/jan25/eventsproxy/event"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	conn, err := net.Dial("udp", "0.0.0.0:7777")
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		conn: conn,
	}
}

func (c *Client) Publish(e *event.Event) {
	t := e.Text()
	log.Println(string(t))
	// fmt.Fprintln(c.conn, e.Text())
	c.conn.Write(e.Text())
}

func main() {
	c := NewClient()

	u, _ := user.Current()
	e := &event.Event{
		Timestamp:    time.Now(),
		TimestampEnd: time.Now().Add(5 * time.Millisecond),
		Host:         u.Uid,
		App:          "example_app",
	}

	log.Println("Sending events..")
	for i := 2; i > 0; i -= 1 {
		c.Publish(e)
		time.Sleep(2 * time.Second)
	}
}
