package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
)

// Product is public; inputs X and Y are private non-native field elements.
type Circuit struct {
	X, Y    emulated.Element[emulated.Secp256k1Fp]
	Product emulated.Element[emulated.Secp256k1Fp] `gnark:",public"`
}

func NewDefinition() frontend.Circuit {
	return &Circuit{}
}
func (c *Circuit) Define(api frontend.API) error {
	field, err := emulated.NewField[emulated.Secp256k1Fp](api)
	if err != nil {
		return err
	}
	field.AssertIsEqual(field.Mul(&c.X, &c.Y), &c.Product)
	return nil
}
