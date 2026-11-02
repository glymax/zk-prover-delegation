package circuit

import "github.com/consensys/gnark/frontend"

// Iterations fixes the circuit size at compile time, not through the witness.
const Iterations = 1024

// Circuit proves Y = X^(2^Iterations) in the circuit's scalar field.
type Circuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
}

func NewDefinition() frontend.Circuit { return &Circuit{} }

func (c *Circuit) Define(api frontend.API) error {
	value := c.X
	for range Iterations {
		value = api.Mul(value, value)
	}
	api.AssertIsEqual(c.Y, value)
	return nil
}
