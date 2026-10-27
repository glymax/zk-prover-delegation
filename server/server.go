package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	address := flag.String("listen", "", "TCP address (required)")
	flag.Parse()
	if *address == "" {
		log.Fatal("--listen is required")
	}
	if flag.NArg() != 0 {
		log.Fatal("unexpected positional arguments")
	}
	if err := run(*address); err != nil {
		log.Fatal(err)
	}
}

func run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			log.Printf("Close listener: %v", closeErr)
		}
	}()
	log.Printf("Listening on %s (Noise NN, CBOR messages; peer identity not verified)", listener.Addr())

	connection, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("accept: %w", err)
	}
	return handleConnection(connection)
}
