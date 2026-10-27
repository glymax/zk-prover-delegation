package circuit

import (
	"github.com/consensys/gnark/frontend"
)

// Circuit proves independent balance transfers, not a full ledger/rollup.
type Circuit struct {
	Transfers []Transfer
	Public    PublicInputs `gnark:",public"`
}
type PublicInputs struct {
	Count frontend.Variable
	Total frontend.Variable
}

func NewDefinition() frontend.Circuit {
	return &Circuit{Transfers: make([]Transfer, 2)}
}
func (c *Circuit) Define(api frontend.API) error {
	api.AssertIsEqual(c.Public.Count, len(c.Transfers))
	var total frontend.Variable = 0
	for _, t := range c.Transfers {
		t.constrain(api)
		total = api.Add(total, t.Amount)
	}
	api.AssertIsEqual(c.Public.Total, total)
	return nil
}

type Transfer struct {
	Before [2]frontend.Variable
	After  [2]frontend.Variable
	Amount frontend.Variable
}

func (t Transfer) constrain(api frontend.API) {
	api.ToBinary(t.Amount, 32)
	for i := 0; i < 2; i++ {
		api.ToBinary(t.Before[i], 32)
		api.ToBinary(t.After[i], 32)
	}
	api.AssertIsEqual(t.After[0], api.Sub(t.Before[0], t.Amount))
	api.AssertIsEqual(t.After[1], api.Add(t.Before[1], t.Amount))
}
