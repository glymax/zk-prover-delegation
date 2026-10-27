// Package transport carries framed Noise messages over an existing connection.
package transport

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/flynn/noise"
)

// Session is sequential, owns cipher states, but does not own/close the connection.
// After any I/O or decryption error the caller must close the connection.
// NN does not authenticate peers; attestation is not implemented yet.
type Session struct {
	connection    io.ReadWriter
	send, receive *noise.CipherState
	binding       []byte
}

// Client performs the initiating side of Noise_NN_25519_ChaChaPoly_SHA256.
func Client(connection io.ReadWriter) (*Session, error) {
	return handshake(connection, true)
}

// Server performs the responding side of the same handshake.
func Server(connection io.ReadWriter) (*Session, error) {
	return handshake(connection, false)
}

func handshake(connection io.ReadWriter, initiator bool) (*Session, error) {
	// NN generates fresh ephemeral key pairs inside the library; no static keys.
	DHFunction := noise.DH25519
	PRNG := rand.Reader
	state, err := noise.NewHandshakeState(noise.Config{
		CipherSuite: noise.NewCipherSuite(DHFunction, noise.CipherChaChaPoly, noise.HashSHA256),
		Random:      PRNG,
		Pattern:     noise.HandshakeNN,
		Initiator:   initiator,
		Prologue:    []byte("my_super_secret_prologue"),
	})
	if err != nil {
		return nil, err
	}
	var clientToServer, serverToClient *noise.CipherState

	for step := range 2 {
		if (step%2 == 0) == initiator {
			var message []byte
			message, clientToServer, serverToClient, err = state.WriteMessage(nil, nil)
			if err == nil {
				err = writeFrame(connection, message)
			}
		} else {
			var message, payload []byte
			message, err = readFrame(connection)
			if err == nil {
				payload, clientToServer, serverToClient, err = state.ReadMessage(nil, message)
			}
			if err == nil && len(payload) != 0 {
				err = fmt.Errorf("unexpected handshake payload")
			}
		}
		if err != nil {
			return nil, fmt.Errorf("noise NN message %d: %w", step+1, err)
		}
	}
	// Split returns initiator->responder first, responder->initiator second.
	session := &Session{
		connection: connection,
		binding:    bytes.Clone(state.ChannelBinding()),
	}
	if initiator {
		session.send = clientToServer
		session.receive = serverToClient
	} else {
		session.send = serverToClient
		session.receive = clientToServer
	}
	return session, nil
}

// ChannelBinding returns a copy of the handshake hash for future attestation.
// Merely obtaining this value does not verify an attestation.
func (s *Session) ChannelBinding() []byte { return bytes.Clone(s.binding) }

// MaxPayloadSize is the encoded application message limit before encryption.
const MaxPayloadSize = noise.MaxMsgLen - 16

// Send encrypts one complete message. No fragmentation is performed.
func (s *Session) Send(plaintext []byte) error {
	if len(plaintext) > MaxPayloadSize {
		return fmt.Errorf("encoded message is %d bytes; maximum is %d", len(plaintext), MaxPayloadSize)
	}
	ciphertext, err := s.send.Encrypt(nil, nil, plaintext)
	if err != nil {
		return err
	}
	return writeFrame(s.connection, ciphertext)
}

func (s *Session) Receive() ([]byte, error) {
	ciphertext, err := readFrame(s.connection)
	if err != nil {
		return nil, err
	}
	plaintext, err := s.receive.Decrypt(nil, nil, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decrypt Noise message: %w", err)
	}
	return plaintext, nil
}

func writeFrame(writer io.Writer, message []byte) error {
	if len(message) == 0 || len(message) > noise.MaxMsgLen {
		return fmt.Errorf("invalid frame length: %d", len(message))
	}
	frame := make([]byte, 2+len(message))
	binary.BigEndian.PutUint16(frame[:2], uint16(len(message)))
	copy(frame[2:], message)
	// A writer may consume fewer bytes than supplied.
	for len(frame) > 0 {
		n, err := writer.Write(frame)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		frame = frame[n:]
	}
	return nil
}

func readFrame(reader io.Reader) ([]byte, error) {
	var header [2]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return nil, err
	}
	length := int(binary.BigEndian.Uint16(header[:]))
	if length == 0 {
		return nil, fmt.Errorf("empty Noise frame")
	}
	message := make([]byte, length)
	if _, err := io.ReadFull(reader, message); err != nil {
		return nil, err
	}
	return message, nil
}
