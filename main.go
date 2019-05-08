package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	var id int = 0
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, id)
		id++
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, id int) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	count := 0
	// Read the incoming connection into the buffer.
	for {
		nums, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		if nums > 0 {
			fmt.Println("Error reading:", buf)
		}
		// Send a response back to person contacting us.
		conn.Write([]byte("Connection " + strconv.FormatInt(id, 10) + " : " + strconv.FormatInt(count, 10)))
		// Close the connection when you're done with it.
		//conn.Close()
		time.Sleep(2 * time.Second)
	}
}
