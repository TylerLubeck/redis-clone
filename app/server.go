package main

import (
	"fmt"
	"strings"
	// Uncomment this block to pass the first stage
	"io"
	 "net"
	 "os"
	 "sync"
)

const (
	NUM_CONNECTIONS = 200
	NUM_WORKERS = 50
)


func handleCommand(command, payload string) string {
	switch command {
	case "+PING":
		return "PONG"
	default:
		return "WTF"
	}
}


func worker(wg *sync.WaitGroup, connections <-chan net.Conn) {
	defer wg.Done()
	for conn := range connections {
		fmt.Println("Handling connection")
		//io.Copy(conn, conn)
		//io.Copy(conn, strings.NewReader("+PONG\r\n"))
		io.Copy(conn, strings.NewReader("+" + handleCommand("+PING", "") +"\r\n"))
		conn.Close()
	}
}


func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	 l, err := net.Listen("tcp", "0.0.0.0:6379")
	 if err != nil {
	 	fmt.Println("Failed to bind to port 6379")
	 	os.Exit(1)
	 }
	 fmt.Printf("Listening on %v", l.Addr())
	 defer l.Close()

	 connections := make(chan net.Conn, NUM_CONNECTIONS)

	 var wg sync.WaitGroup

	 for i := 0; i < NUM_WORKERS; i++ {
		wg.Add(1)
		go worker(&wg, connections)
	 }

	 for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			os.Exit(1)
		}

		connections <- conn
	 }
}
