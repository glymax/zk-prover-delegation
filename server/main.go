package main

import (
	"bachelor/implementation/protocol"
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
	keys, err := protocol.LoadKeys("server/keys/private.key", "server/keys/client.public.key")
	if err != nil {
		log.Fatal(err)
	}
	if err := run(*address, keys); err != nil {
		log.Fatal(err)
	}
}

func run(address string, keys protocol.Keys) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			log.Printf("Close listener: %v", closeErr)
		}
	}()
	log.Printf("Listening on %s (Noise KK, CBOR messages; pinned client key required)", listener.Addr())

	connection, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("accept: %w", err)
	}
	return handleConnection(connection, keys)
}
