package backend

import (
	"fmt"
	"io"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"
)

// BinaryKey is the common serialization contract of both backends' keys.
type BinaryKey interface {
	io.WriterTo
	io.ReaderFrom
}

type Prepared struct {
	System       string
	Constraints  constraint.ConstraintSystem
	ProvingKey   BinaryKey
	VerifyingKey BinaryKey
}

// ValidateBackend checks the backend name and whether the curve is supported.
func ValidateBackend(system, curve string) error {
	if system != "groth16" && system != "plonk" {
		return fmt.Errorf("unsupported backend %q; expected groth16 or plonk", system)
	}
	_, err := ParseCurve(curve)
	return err
}

func ParseCurve(name string) (ecc.ID, error) {
	switch name {
	case "bn254":
		return ecc.BN254, nil
	case "bls12-377":
		return ecc.BLS12_377, nil
	case "bls12-381":
		return ecc.BLS12_381, nil
	case "bw6-761":
		return ecc.BW6_761, nil
	default:
		return ecc.UNKNOWN, fmt.Errorf("unsupported curve %q", name)
	}
}

func Empty(system, curve string) (*Prepared, error) {
	id, err := ParseCurve(curve)
	if err != nil {
		return nil, err
	}
	switch system {
	case "groth16":
		return &Prepared{System: system, Constraints: groth16.NewCS(id), ProvingKey: groth16.NewProvingKey(id), VerifyingKey: groth16.NewVerifyingKey(id)}, nil
	case "plonk":
		return &Prepared{System: system, Constraints: plonk.NewCS(id), ProvingKey: plonk.NewProvingKey(id), VerifyingKey: plonk.NewVerifyingKey(id)}, nil
	default:
		return nil, fmt.Errorf("unsupported backend %q", system)
	}
}

func BuildSetup(d frontend.Circuit, system, curve string) (*Prepared, error) {
	cs, err := Compile(d, system, curve)
	if err != nil {
		return nil, err
	}
	return Setup(cs, system)
}

func Compile(definition frontend.Circuit, system, curve string) (constraint.ConstraintSystem, error) {
	id, err := ParseCurve(curve)
	if err != nil {
		return nil, err
	}
	switch system {
	case "groth16":
		return frontend.Compile(id.ScalarField(), r1cs.NewBuilder, definition)
	case "plonk":
		return frontend.Compile(id.ScalarField(), scs.NewBuilder, definition)
	default:
		return nil, fmt.Errorf("unsupported backend %q", system)
	}
}

func Setup(cs constraint.ConstraintSystem, system string) (*Prepared, error) {
	switch system {
	case "groth16":
		pk, vk, err := groth16.Setup(cs)
		if err != nil {
			return nil, err
		}
		return &Prepared{System: system, Constraints: cs, ProvingKey: pk, VerifyingKey: vk}, nil
	case "plonk":
		// Experimental SRS only: replace with a trusted ceremony SRS for real use.
		// NewSRS caches in memory; our CLI starts a fresh process for each cold job.
		log.Print("PLONK: using UNSAFE experimental SRS, not a trusted setup")
		srs, lagrange, err := unsafekzg.NewSRS(cs)
		if err != nil {
			return nil, fmt.Errorf("experimental SRS: %w", err)
		}
		pk, vk, err := plonk.Setup(cs, srs, lagrange)
		if err != nil {
			return nil, err
		}
		return &Prepared{System: system, Constraints: cs, ProvingKey: pk, VerifyingKey: vk}, nil
	default:
		return nil, fmt.Errorf("unsupported backend %q", system)
	}
}

func Prove(setup *Prepared, w witness.Witness) (io.WriterTo, error) {
	switch setup.System {
	case "groth16":
		pk, ok := setup.ProvingKey.(groth16.ProvingKey)
		if !ok {
			return nil, fmt.Errorf("expected Groth16 proving key")
		}
		return groth16.Prove(setup.Constraints, pk, w)
	case "plonk":
		pk, ok := setup.ProvingKey.(plonk.ProvingKey)
		if !ok {
			return nil, fmt.Errorf("expected PLONK proving key")
		}
		return plonk.Prove(setup.Constraints, pk, w)
	default:
		return nil, fmt.Errorf("unsupported backend %q", setup.System)
	}
}
