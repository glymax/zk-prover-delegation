// Command authkeygen provisions the fixed experimental participants out of band.
package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: authkeygen\nRun from the project root; regenerates keys in client/keys and server/keys.")
	}
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("usage: authkeygen (no arguments)")
	}
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}

func generate() error {
	clientKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	serverKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	if err = os.MkdirAll("client/keys", 0700); err != nil {
		return err
	}
	if err = os.MkdirAll("server/keys", 0700); err != nil {
		return err
	}

	if err = writeKey("client/keys", "private.key", clientKey.Bytes(), 0600); err != nil {
		return err
	}
	if err = writeKey("client/keys", "server.public.key", serverKey.PublicKey().Bytes(), 0644); err != nil {
		return err
	}

	if err = writeKey("server/keys", "private.key", serverKey.Bytes(), 0600); err != nil {
		return err
	}
	if err = writeKey("server/keys", "client.public.key", clientKey.PublicKey().Bytes(), 0644); err != nil {
		return err
	}

	fmt.Println("Created X25519 identities and trusted public-key bindings in client/keys and server/keys.")
	return nil
}

func writeKey(dir, name string, key []byte, mode os.FileMode) error {
	return os.WriteFile(filepath.Join(dir, name), key, mode)
}
