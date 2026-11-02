package circuit

import "github.com/consensys/gnark/frontend"

const VectorSize = 128

// Circuit proves Y = sum(A[i]*B[i]) in the circuit's scalar field.
// Both vectors are private; no integer range constraints are imposed.
type Circuit struct {
	A [VectorSize]frontend.Variable
	B [VectorSize]frontend.Variable
	Y frontend.Variable `gnark:",public"`
}

func NewDefinition() frontend.Circuit { return &Circuit{} }

func (c *Circuit) Define(api frontend.API) error {
	var total frontend.Variable = 0
	for i := 0; i < VectorSize; i++ {
		total = api.Add(total, api.Mul(c.A[i], c.B[i]))
	}
	api.AssertIsEqual(c.Y, total)
	return nil
}
