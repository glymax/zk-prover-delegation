// Adapted from gnark v0.16.4 examples/accelerated_gpu/accelerated_gpu_example_test.go.
// Copyright Consensys Software Inc. Licensed under Apache-2.0; see ../LICENSE.
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

	// we have commented out another option for the circuit implementation for
	// faster testing. Comment the previous part and uncomment the following to
	// use it.

	// api.AssertIsDifferent(c.Res.A0.Limbs[0], 0) // dummy constraint to avoid empty circuit
	// api.AssertIsDifferent(c.Res.A0.Limbs[1], 0)
	// api.AssertIsDifferent(c.Res.A0.Limbs[2], 0)
	// api.AssertIsDifferent(c.Res.A0.Limbs[3], 0)
	return nil
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &ExampleCircuit{}
}
