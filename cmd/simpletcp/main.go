package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/rcrowley/go-metrics"
)

// https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8972"
	CONN_TYPE = "tcp"
)

var (
	opsRate = metrics.NewRegisteredMeter("ops", nil)
)

func main() {

	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	opsRate.Mark(1)
	// Send a response back to person contacting us.
	_, err := conn.Write([]byte("Simple answer"))
	if err != nil {
		log.Printf("Error sending ")
	}
	// Close the connection when you're done with it.
	conn.Close()
}
