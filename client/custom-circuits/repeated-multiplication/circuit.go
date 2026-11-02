package circuit

import "github.com/consensys/gnark/frontend"

// Exponent is the number of X factors, fixed when compiling the circuit.
const Exponent = 1024

// Circuit proves Y = X^Exponent in the circuit's scalar field.
// This deliberately uses a linear chain rather than exponentiation by squaring.
type Circuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
}

func NewDefinition() frontend.Circuit { return &Circuit{} }

func (c *Circuit) Define(api frontend.API) error {
	value := c.X
	// X^Exponent needs Exponent-1 multiplications starting from X.
	for i := 1; i < Exponent; i++ {
		value = api.Mul(value, c.X)
	}
	api.AssertIsEqual(c.Y, value)
	return nil
}
