package circuits

import (
	"bachelor/implementation/circuits/cubic"
	"bachelor/implementation/circuits/emulated-mul"
	"bachelor/implementation/circuits/exponentiate"
	"bachelor/implementation/circuits/inputpacking"
	"bachelor/implementation/circuits/mimc"
	"bachelor/implementation/circuits/pairing"
	"bachelor/implementation/circuits/point-check"
	"bachelor/implementation/circuits/sha256"
	"bachelor/implementation/circuits/signature-hash-check"
	"bachelor/implementation/circuits/sudoku"

	"github.com/consensys/gnark/frontend"
)

// Constructor returns a fresh circuit with its fixed shape and constants initialized.
// It does not populate witness values.
type Constructor func() frontend.Circuit

// Catalog contains the fixed built-in circuit definitions.
var catalog = map[string]Constructor{
	"sha256":                sha256.NewDefinition,
	"cubic":                 cubic.NewDefinition,
	"exponentiate":          exponentiate.NewDefinition,
	"mimc":                  mimc.NewDefinition,
	"emulated-mul":          emulatedmul.NewDefinition,
	"sudoku":                sudoku.NewDefinition,
	"point-check":           pointcheck.NewDefinition,
	"inputpacking-packed":   inputpacking.NewDefinition,
	"inputpacking-unpacked": inputpacking.NewUnpackedDefinition,
	"signature-hash-check":  signaturehashcheck.NewDefinition,
	"pairing":               pairing.NewDefinition,
}

// Built-in definitions currently have fixed shapes, initialized by their constructors.
func Lookup(name string) (Constructor, bool) {
	constructor, ok := catalog[name]
	return constructor, ok
}
