// Adapted from gnark v0.16.4 examples/sudoku/sudoku_example_test.go.
// Copyright Consensys Software Inc. Licensed under Apache-2.0; see ../LICENSE.
package sudoku

import (
	"github.com/consensys/gnark/frontend"
)

// SudokuCircuit represents a Sudoku circuit. It contains two grids: the
// challenge and solution grids (named Challenge and Solution respectively). The
// challenge grid is public, while the solution grid is private.
type SudokuCircuit struct {
	Challenge SudokuGrid `gnark:"Challenge,public"`
	Solution  SudokuGrid `gnark:"Solution,secret"`
}

// SudokuGrid represents a 9x9 Sudoku grid in-circuit.
type SudokuGrid [9][9]frontend.Variable

// Define defines the constraints of the Sudoku circuit.
func (circuit *SudokuCircuit) Define(api frontend.API) error {
	// Constraint 1: Each cell value in the CompleteGrid must be between 1 and 9
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			api.AssertIsLessOrEqual(circuit.Solution[i][j], 9)
			api.AssertIsLessOrEqual(1, circuit.Solution[i][j])
		}
	}

	// Constraint 2: Each row in the CompleteGrid must contain unique values
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			for k := j + 1; k < 9; k++ {
				api.AssertIsDifferent(circuit.Solution[i][j], circuit.Solution[i][k])
			}
		}
	}

	// Constraint 3: Each column in the CompleteGrid must contain unique values
	for j := 0; j < 9; j++ {
		for i := 0; i < 9; i++ {
			for k := i + 1; k < 9; k++ {
				api.AssertIsDifferent(circuit.Solution[i][j], circuit.Solution[k][j])
			}
		}
	}

	// Constraint 4: Each 3x3 sub-grid in the CompleteGrid must contain unique values
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			for i := 0; i < 9; i++ {
				for j := i + 1; j < 9; j++ {
					row1 := boxRow*3 + i/3
					col1 := boxCol*3 + i%3
					row2 := boxRow*3 + j/3
					col2 := boxCol*3 + j%3
					api.AssertIsDifferent(circuit.Solution[row1][col1], circuit.Solution[row2][col2])
				}
			}
		}
	}

	// Constraint 5: The values in the IncompleteGrid must match the CompleteGrid where provided
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			isCellGiven := api.IsZero(circuit.Challenge[i][j])
			api.AssertIsEqual(api.Select(isCellGiven, circuit.Solution[i][j], circuit.Challenge[i][j]), circuit.Solution[i][j])
		}
	}

	return nil
}

// NewDefinition initializes the fixed circuit shape.
func NewDefinition() frontend.Circuit {
	return &SudokuCircuit{}
}
