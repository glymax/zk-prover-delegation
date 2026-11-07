package protocol

import (
	"crypto/ecdh"
	"fmt"
	"os"

	"github.com/flynn/noise"
)

// Keys models an identity and an authenticated peer-key binding provisioned out of band.
// It does not implement certificate issuance or validation.
type Keys struct {
	own  noise.DHKey
	peer []byte
}

// LoadKeys reads raw 32-byte X25519 keys, deriving our public key locally.
func LoadKeys(privateFile, peerFile string) (Keys, error) {
	private, err := readKey(privateFile)
	if err != nil {
		return Keys{}, fmt.Errorf("private key: %w", err)
	}
	peer, err := readKey(peerFile)
	if err != nil {
		return Keys{}, fmt.Errorf("peer public key: %w", err)
	}
	key, err := ecdh.X25519().NewPrivateKey(private)
	if err != nil {
		return Keys{}, fmt.Errorf("X25519 private key: %w", err)
	}
	return Keys{own: noise.DHKey{Private: private, Public: key.PublicKey().Bytes()}, peer: peer}, nil
}

func readKey(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) != 32 {
		return nil, fmt.Errorf("%s: expected 32 bytes, got %d", path, len(data))
	}
	return data, nil
}
