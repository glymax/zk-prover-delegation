package circuits

import (
	gcubic "bachelor/implementation/circuits/cubic"
	"bachelor/implementation/circuits/emulated"
	"bachelor/implementation/circuits/exponentiate"
	"bachelor/implementation/circuits/inputpacking"
	"bachelor/implementation/circuits/mimc"
	"bachelor/implementation/circuits/pairing"
	"bachelor/implementation/circuits/plonkexample"
	"bachelor/implementation/circuits/print_constraints"
	"bachelor/implementation/circuits/pubkeyhashing"
	"bachelor/implementation/circuits/rollup"
	"bachelor/implementation/circuits/sudoku"
	"bachelor/implementation/circuits/witness"

	"github.com/consensys/gnark/frontend"
)

// Constructor returns a fresh circuit with its fixed shape and constants initialized.
// It does not populate witness values.
type Constructor func() frontend.Circuit

// Catalog pins the example definitions locally. Serialization reuses local Cubic
// just as the upstream example does. The GPU example's
// pairing circuit runs on the current CPU Groth16 backend (no ICICLE required).
var catalog = map[string]Constructor{
	"cubic":                 gcubic.NewDefinition,
	"exponentiate":          exponentiate.NewDefinition,
	"mimc":                  mimc.NewDefinition,
	"emulated":              emulated.NewDefinition,
	"sudoku":                sudoku.NewDefinition,
	"witness":               witness.NewDefinition,
	"plonk":                 plonkexample.NewDefinition,
	"serialization":         gcubic.NewDefinition,
	"print_constraints":     print_constraints.NewDefinition,
	"inputpacking":          inputpacking.NewDefinition,
	"inputpacking-unpacked": inputpacking.NewUnpackedDefinition,
	"pubkeyhashing":         pubkeyhashing.NewDefinition,
	"accelerated_gpu":       pairing.NewDefinition,
	"rollup":                rollup.NewDefinition,
}

// Built-in definitions currently have fixed shapes, including rollup's fixed paths.
func Lookup(name string) (Constructor, bool) {
	constructor, ok := catalog[name]
	return constructor, ok
}
