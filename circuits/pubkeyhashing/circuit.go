// Adapted from gnark v0.16.4 examples/pubkeyhashing/pubkeyhashing_test.go.
// Copyright Consensys Software Inc. Licensed under Apache-2.0; see ../LICENSE.
package pubkeyhashing

import (
	"fmt"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/conversion"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/emulated/emparams"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/signature/ecdsa"
)

// PubkeySHA2 is a circuit that verifies ECDSA signature and checks that the
// hash of the public key matches the expected hash.
//
// The fields of the struct define the public and private inputs to the circuit.
// The actual circuit is defined in the [PubKeySHA2.Define] method.
type PubKeySHA2 struct {
	// PublicKeyHash is 32 bytes, but we split it into two 16-byte variables to fit into BN254 field
	PublicKeyHash [2]frontend.Variable                                        `gnark:",public"`
	Signature     ecdsa.Signature[emparams.Secp256k1Fr]                       `gnark:",public"`
	Msg           emulated.Element[emparams.Secp256k1Fr]                      // if tag is not set, then it is a private input
	PublicKey     ecdsa.PublicKey[emparams.Secp256k1Fp, emparams.Secp256k1Fr] // actual public key is also a private input
}

func (c *PubKeySHA2) Define(api frontend.API) error {
	// -- hash the given public key
	//  - first we convert the public key coordinates to bytes
	xbytes, err := conversion.EmulatedToBytes(api, &c.PublicKey.X)
	if err != nil {
		return fmt.Errorf("failed to convert PublicKey.X to bytes: %w", err)
	}
	ybytes, err := conversion.EmulatedToBytes(api, &c.PublicKey.Y)
	if err != nil {
		return fmt.Errorf("failed to convert PublicKey.Y to bytes: %w", err)
	}
	//  - now we compute the SHA2 hash of the concatenated bytes
	h, err := sha2.New(api)
	if err != nil {
		return fmt.Errorf("failed to create SHA2 instance: %w", err)
	}
	h.Write(xbytes)
	h.Write(ybytes)
	//  - and compute the hash
	computedHash := h.Sum()
	// -- now we check that the computed hash matches the expected hash
	//  - first, we used [2]frontend.Variable to store the hash so that we wouldn't be using too much public inputs and we want the parts to fit into BN254 field, so 16-byte chunks
	//    we convert it back to bytes
	var hashpubkeybytes []uints.U8
	for i := range c.PublicKeyHash {
		bts, conversionErr := conversion.NativeToBytes(api, c.PublicKeyHash[i])
		if conversionErr != nil {
			return fmt.Errorf("failed to convert PublicKeyHash[%d] to bytes: %w", i, conversionErr)
		}
		// NativeToBytes returns 32 bytes (MSB order), but we set only 16 bytes so take the last 16 bytes
		hashpubkeybytes = append(hashpubkeybytes, bts[16:]...)
	}
	//  - now we need to initialize bytes gadget for comparison
	bapi, err := uints.NewBytes(api)
	if err != nil {
		return fmt.Errorf("failed to create bytes gadget: %w", err)
	}
	if len(hashpubkeybytes) != len(computedHash) {
		return fmt.Errorf("hashpubkeybytes and computedHash have different lengths: %d vs %d", len(hashpubkeybytes), len(computedHash))
	}
	//  - finally we check that the computed hash matches the expected hash
	for i := range hashpubkeybytes {
		bapi.AssertIsEqual(hashpubkeybytes[i], computedHash[i])
	}

	// -- now we check that the signature is valid
	c.PublicKey.Verify(api, sw_emulated.GetCurveParams[emparams.Secp256k1Fp](), &c.Msg, &c.Signature)

	return nil
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &PubKeySHA2{}
}
