package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type Circuit struct {
	Values []frontend.Variable
	Sum    frontend.Variable `gnark:",public"`
}

func NewDefinition() frontend.Circuit {
	return &Circuit{Values: make([]frontend.Variable, 3)}
}

func (c *Circuit) Define(api frontend.API) error {
	var sum frontend.Variable = 0
	for _, value := range c.Values {
		sum = api.Add(sum, value)
	}
	api.AssertIsEqual(sum, c.Sum)
	return nil
}
