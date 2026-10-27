// Package codec serializes network messages and worker jobs using CBOR.
package codec

import "github.com/fxamacker/cbor/v2"

// Encode serializes one complete message, including its binary fields.
func Encode(value any) ([]byte, error) {
	return cbor.Marshal(value)
}

// Decode restores one complete message. Inputs remain JSON bytes.
func Decode(data []byte, value any) error {
	return cbor.Unmarshal(data, value)
}
