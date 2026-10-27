// Adapted from gnark v0.16.4 examples/witness/witness_test.go.
// Copyright Consensys Software Inc. Licensed under Apache-2.0; see ../LICENSE.
package witness

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_bn254"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
)

type Circuit struct {
	A frontend.Variable `gnark:",public"`
	B frontend.Variable
	P sw_emulated.AffinePoint[sw_bn254.BaseField] `gnark:",public"`
}

func (c *Circuit) Define(api frontend.API) error {
	api.AssertIsDifferent(c.A, c.B)

	curve, err := sw_bn254.NewPairing(api)
	if err != nil {
		return err
	}
	curve.AssertIsOnG1(&c.P)
	return nil
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &Circuit{}
}
