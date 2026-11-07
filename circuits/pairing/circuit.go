// Adapted from gnark v0.16.4 examples/accelerated_gpu/accelerated_gpu_example_test.go.
// Copyright Consensys Software Inc. Licensed under Apache-2.0; see ../../docs/LICENSE.
package pairing

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_bls12381"
)

type ExampleCircuit struct {
	A   [2]sw_bls12381.G1Affine
	B   [2]sw_bls12381.G2Affine
	Res sw_bls12381.GTEl `gnark:",public"`
}

func (c *ExampleCircuit) Define(api frontend.API) error {
	pr, err := sw_bls12381.NewPairing(api)
	if err != nil {
		return err
	}
	res, err := pr.Pair([]*sw_bls12381.G1Affine{&c.A[0], &c.A[1]}, []*sw_bls12381.G2Affine{&c.B[0], &c.B[1]})
	if err != nil {
		return err
	}
	pr.AssertIsEqual(res, &c.Res)

	return nil
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &ExampleCircuit{}
}
