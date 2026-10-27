// Package attestation currently implements only the message exchange, not a quote.
package attestation

import (
	"bachelor/implementation/protocol"
	"fmt"
)

func Handle(r protocol.Request) (protocol.Response, error) {
	if r.Type != "attestation" || len(r.Inputs) != 0 || len(r.Source) != 0 || r.Circuit != "" || r.Backend != "" || r.Curve != "" || r.SetupName != "" {
		return protocol.Response{}, fmt.Errorf("expected an attestation request before proving")
	}
	if len(r.Nonce) != 32 {
		return protocol.Response{}, fmt.Errorf("attestation nonce must be 32 bytes")
	}
	return protocol.Response{Type: "attestation", Status: "not_implemented", Nonce: r.Nonce}, nil
}
