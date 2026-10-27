package custom

import (
	"github.com/consensys/gnark/frontend"
)

type Circuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
}

// Define simple function for demo: X^3 + X + 5 = Y
func (c *Circuit) Define(api frontend.API) error {
	api.AssertIsEqual(api.Add(api.Mul(c.X, c.X, c.X), c.X, 5), c.Y)
	return nil
}

func NewDefinition() frontend.Circuit {
	return &Circuit{}
}
