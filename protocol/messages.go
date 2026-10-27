// Package protocol defines CBOR network messages, without gnark.
package protocol

type Request struct {
	Type      string `cbor:"request_type"`
	Circuit   string `cbor:"circuit,omitempty"`
	Backend   string `cbor:"backend,omitempty"`
	Curve     string `cbor:"curve,omitempty"`
	SetupName string `cbor:"setup_name,omitempty"`
	Nonce     []byte `cbor:"nonce,omitempty"` // 32 random bytes for attestation only.
	Inputs    []byte `cbor:"inputs,omitempty"`
	Source    []byte `cbor:"source,omitempty"`
}
type Response struct {
	Type          string `cbor:"type"`
	Error         string `cbor:"error,omitempty"`
	Circuit       string `cbor:"circuit,omitempty"`
	Backend       string `cbor:"backend,omitempty"`
	Curve         string `cbor:"curve,omitempty"`
	Status        string `cbor:"status,omitempty"`
	Nonce         []byte `cbor:"nonce,omitempty"`
	Proof         []byte `cbor:"proof,omitempty"`
	PublicWitness []byte `cbor:"public_witness,omitempty"`
	VerifyingKey  []byte `cbor:"verifying_key,omitempty"`
	Quote         []byte `cbor:"quote,omitempty"` // Raw attestation evidence; interpreted by the attestation verifier.
}
