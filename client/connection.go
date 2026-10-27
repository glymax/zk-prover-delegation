package main

import (
	"bachelor/implementation/codec"
	"bachelor/implementation/protocol"
	"bachelor/implementation/transport"
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"time"
)

func exchange(cfg Config) (protocol.Response, error) {

	connection, err := net.DialTimeout("tcp", cfg.Server, 5*time.Second)
	if err != nil {
		return protocol.Response{}, fmt.Errorf("connect: %w", err)
	}

	defer func() {
		closeErr := connection.Close()
		if closeErr != nil {
			log.Printf("Close connection: %v", closeErr)
		}
	}()

	if err = connection.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return protocol.Response{}, fmt.Errorf("set deadline: %w", err)
	}

	// Noise handshake
	session, err := transport.Client(connection)
	if err != nil {
		return protocol.Response{}, fmt.Errorf("handshake: %w", err)
	}
	log.Print("Noise NN complete (peer identity not verified)")

	// Placeholder exchange on this exact Noise session. It does NOT verify a TEE.
	var nonce [32]byte
	if _, err = rand.Read(nonce[:]); err != nil {
		return protocol.Response{}, err
	}
	att, err := codec.Encode(protocol.Request{Type: "attestation", Nonce: nonce[:]})
	if err != nil {
		return protocol.Response{}, err
	}
	if err = connection.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return protocol.Response{}, err
	}
	if err = session.Send(att); err != nil {
		return protocol.Response{}, err
	}
	reply, err := session.Receive()
	if err != nil {
		return protocol.Response{}, err
	}
	var report protocol.Response
	if err = codec.Decode(reply, &report); err != nil {
		return protocol.Response{}, err
	}
	if report.Type != "attestation" || report.Status != "not_implemented" || !bytes.Equal(report.Nonce, nonce[:]) || report.Error != "" {
		return protocol.Response{}, fmt.Errorf("unexpected attestation stub response")
	}
	log.Print("Attestation stub only: nonce echoed, no report or identity verification; continuing prototype")

	if err = connection.SetDeadline(time.Time{}); err != nil {
		return protocol.Response{}, fmt.Errorf("set exchange deadline: %w", err)
	}

	request, err := LoadRequest(cfg)
	if err != nil {
		return protocol.Response{}, err
	}
	payload, err := codec.Encode(request)
	if err != nil {
		return protocol.Response{}, err
	}

	if err = session.Send(payload); err != nil {
		return protocol.Response{}, fmt.Errorf("send request: %w", err)
	}

	var response protocol.Response
	payload, err = session.Receive()
	if err != nil {
		return protocol.Response{}, fmt.Errorf("receive response: %w", err)
	}

	if err = codec.Decode(payload, &response); err != nil {
		return protocol.Response{}, fmt.Errorf("decode response: %w", err)
	}

	if response.Type == "error" {
		return protocol.Response{}, fmt.Errorf("server: %s", response.Error)
	}

	if response.Type != "proof" || response.Error != "" {
		return protocol.Response{}, fmt.Errorf("unexpected response type %q", response.Type)
	}
	return response, nil
}
