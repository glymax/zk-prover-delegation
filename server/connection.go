package main

import (
	"bachelor/implementation/codec"
	"bachelor/implementation/protocol"
	"bachelor/implementation/proving"
	"bachelor/implementation/proving/worker"
	"bachelor/implementation/server/attestation"
	"bachelor/implementation/transport"
	"fmt"
	"log"
	"net"
	"time"
)

func handleConnection(connection net.Conn) error {
	defer func() {
		closeErr := connection.Close()
		if closeErr != nil {
			log.Printf("Close connection: %v", closeErr)
		}
	}()
	// Bound network setup; proving has no execution deadline.
	if err := connection.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}

	// Noise handshake
	session, err := transport.Server(connection)

	if err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	if err = connection.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("set exchange deadline: %w", err)
	}

	// One attestation exchange followed by one prove request on the same channel.
	attPayload, err := session.Receive()
	if err != nil {
		return err
	}
	var attRequest protocol.Request
	if err = codec.Decode(attPayload, &attRequest); err != nil {
		return err
	}

	attResponse, attErr := attestation.Handle(attRequest)
	if attErr != nil {
		attResponse = protocol.Response{Type: "error", Error: attErr.Error()}
	}
	encoded, err := codec.Encode(attResponse)
	if err != nil {
		return err
	}
	if err = session.Send(encoded); err != nil {
		return err
	}
	if attErr != nil {
		return attErr
	}
	if err = connection.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}

	var proveRequest protocol.Request
	// Fail closed: no plaintext reply or fallback after a transport error.

	payload, err := session.Receive()

	if err != nil {
		return fmt.Errorf("receive message: %w", err)
	}

	if err = codec.Decode(payload, &proveRequest); err != nil {
		return fmt.Errorf("receive request: %w", err)
	}

	// Computation is not bounded by a socket I/O deadline.
	if err = connection.SetDeadline(time.Time{}); err != nil {
		return err
	}
	proveResponse, runErr := handleProve(proveRequest)
	if runErr != nil {
		proveResponse = protocol.Response{Type: "error", Error: runErr.Error()}
	}
	if err = connection.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}

	payload, err = codec.Encode(proveResponse)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	if len(payload) > transport.MaxPayloadSize {
		payload, err = codec.Encode(protocol.Response{Type: "error", Error: "encoded response exceeds single Noise frame; fragmentation is not implemented"})
		if err != nil {
			return err
		}
	}
	if err = session.Send(payload); err != nil {
		return fmt.Errorf("send response: %w", err)
	}

	return nil
}

// handleProve converts a network request to the same job used by the local client.
func handleProve(request protocol.Request) (protocol.Response, error) {
	if len(request.Nonce) != 0 {
		return protocol.Response{}, fmt.Errorf("nonce is only used for attestation")
	}
	job := proving.Job{Circuit: request.Circuit, Backend: request.Backend, Curve: request.Curve, SetupName: request.SetupName, Inputs: request.Inputs}
	var result proving.Result
	var err error
	switch request.Type {
	case "prove":
		if len(request.Source) != 0 {
			return protocol.Response{}, fmt.Errorf("unexpected circuit source")
		}
		result, err = worker.RunBuiltin(job)
	case "custom_prove":
		result, err = worker.RunCustom(job, request.Source)
	default:
		return protocol.Response{}, fmt.Errorf("unknown proving request %q", request.Type)
	}
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.Response{Type: "proof", Circuit: result.Circuit, Backend: result.Backend, Curve: result.Curve, Proof: result.Proof, PublicWitness: result.PublicWitness, VerifyingKey: result.VerifyingKey}, nil
}
