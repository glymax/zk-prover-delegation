// Adapted from gnark v0.16.4 examples/inputpacking/inputpacking_test.go.
// Copyright Consensys Software Inc. Licensed under Apache-2.0; see ../LICENSE.
// Adaptation: use supported secret tags and expose PublicHash as intended by the example.
package inputpacking

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/emulated"
)

func inCircuitComputation(api frontend.API, input1, input2 frontend.Variable, expected frontend.Variable) {
	res := api.Mul(input1, input2)
	api.AssertIsEqual(res, expected)
}

func inCircuitComputationEmulated(api frontend.API, input1, input2 emulated.Element[emulated.BN254Fp], expected emulated.Element[emulated.BN254Fp]) error {
	f, err := emulated.NewField[emulated.BN254Fp](api)
	if err != nil {
		return err
	}
	res := f.Mul(&input1, &input2)
	f.AssertIsEqual(res, &expected)
	return nil
}

// UnpackedCircuit represents a circuit where all public inputs are given as is
type UnpackedCircuit struct {
	Input1, Input2                 frontend.Variable                  `gnark:",public"`
	EmulatedInput1, EmulatedInput2 emulated.Element[emulated.BN254Fp] `gnark:",public"`
	Output                         frontend.Variable                  `gnark:",secret"`
	EmulatedOutput                 emulated.Element[emulated.BN254Fp] `gnark:",secret"`
}

func (circuit *UnpackedCircuit) Define(api frontend.API) error {
	inCircuitComputation(api, circuit.Input1, circuit.Input2, circuit.Output)
	return inCircuitComputationEmulated(api, circuit.EmulatedInput1, circuit.EmulatedInput2, circuit.EmulatedOutput)
}

// PackedCircuit represents a circuit where all public inputs are given as private instead and we provide a hash of them as the only public input.
type PackedCircuit struct {
	PublicHash frontend.Variable `gnark:",public"`

	Input1, Input2                 frontend.Variable                  `gnark:",secret"`
	EmulatedInput1, EmulatedInput2 emulated.Element[emulated.BN254Fp] `gnark:",secret"`
	Output                         frontend.Variable                  `gnark:",secret"`
	EmulatedOutput                 emulated.Element[emulated.BN254Fp] `gnark:",secret"`
}

func (circuit *PackedCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	h.Write(circuit.Input1)
	h.Write(circuit.Input2)
	h.Write(circuit.EmulatedInput1.Limbs...)
	h.Write(circuit.EmulatedInput2.Limbs...)
	dgst := h.Sum()
	api.AssertIsEqual(dgst, circuit.PublicHash)

	inCircuitComputation(api, circuit.Input1, circuit.Input2, circuit.Output)
	return inCircuitComputationEmulated(api, circuit.EmulatedInput1, circuit.EmulatedInput2, circuit.EmulatedOutput)
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &PackedCircuit{}
}

// NewUnpackedDefinition initializes the fixed circuit shape.
func NewUnpackedDefinition() frontend.Circuit {
	return &UnpackedCircuit{}
}
