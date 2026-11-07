package sha256

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
)

// MessageSize fixes the number of private message bytes hashed by this circuit.
const MessageSize = 32

// Circuit proves that a private message has the given public SHA-256 digest.
type Circuit struct {
	Message [MessageSize]uints.U8
	Digest  [32]uints.U8 `gnark:",public"`
}

func (c *Circuit) Define(api frontend.API) error {
	h, err := sha2.New(api)
	if err != nil {
		return err
	}
	h.Write(c.Message[:])
	digest := h.Sum()
	bytes, err := uints.NewBytes(api)
	if err != nil {
		return err
	}
	for i := range c.Digest {
		bytes.AssertIsEqual(c.Digest[i], digest[i])
	}
	return nil
}

func NewDefinition() frontend.Circuit {
	return &Circuit{}
}
