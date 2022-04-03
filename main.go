package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter int

	listenAddr = ":8080"

	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %v", err)
		}

		backend := chooseBackend()
		fmt.Printf("counter%d backend:%s\n", counter, backend)
		go func() {
			err := proxy(backend, conn)
			if err != nil {
				log.Printf("proxy error: %v", err)
			}
		}()

	}

}

func proxy(backend string, c net.Conn) error {
	bc, err := net.Dial("tcp", backend)
	if err != nil {
		return fmt.Errorf("failed to connect to backend %s: %v", backend, err)
	}

	// c -> bc
	go io.Copy(c, bc)

	// c -> bc
	go io.Copy(bc, c)

	return nil
}

func chooseBackend() string {
	s := server[counter%len(server)]
	counter++
	return s
}
