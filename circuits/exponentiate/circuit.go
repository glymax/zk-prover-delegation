// Copied from gnark v0.16.4 examples/exponentiate/exponentiate.go.
// Copyright 2020-2025 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

package exponentiate

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/bits"
)

// Circuit y == x**e
// only the bitSize least significant bits of e are used
type Circuit struct {
	// tagging a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:",public"`
	Y frontend.Variable `gnark:",public"`

	E frontend.Variable
}

// Define declares the circuit's constraints
// y == x**e
func (circuit *Circuit) Define(api frontend.API) error {

	// number of bits of exponent
	const bitSize = 8

	// specify constraints
	output := frontend.Variable(1)
	binary := bits.ToBinary(api, circuit.E, bits.WithNbDigits(bitSize))

	for i := 0; i < len(binary); i++ {
		if i != 0 {
			output = api.Mul(output, output)
		}
		multiply := api.Mul(output, circuit.X)
		output = api.Select(binary[len(binary)-1-i], multiply, output)

	}

	api.AssertIsEqual(circuit.Y, output)

	return nil
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &Circuit{}
}
