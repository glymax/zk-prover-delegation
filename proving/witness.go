package proving

import (
	"bachelor/implementation/proving/backend"

	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

// decodeWitness expects a complete JSON witness matching the circuit schema.
func decodeWitness(definition frontend.Circuit, inputs []byte, curve string) (witness.Witness, error) {
	id, err := backend.ParseCurve(curve)
	if err != nil {
		return nil, err
	}
	field := id.ScalarField()
	schema, err := frontend.NewSchema(field, definition)
	if err != nil {
		return nil, err
	}
	full, err := witness.New(field)
	if err != nil {
		return nil, err
	}
	if err = full.FromJSON(schema, inputs); err != nil {
		return nil, err
	}
	return full, nil
}
