package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func main() {
	// Create the coefficient matrices A (m x n) (4 x 6)
	// The data is provided in a flat slice, ordered row by row
	A := mat.NewDense(4, 6, []float64{
		0, 0, 0, 0, 1, 1,
		0, 1, 0, 0, 0, 1,
		0, 0, 1, 1, 1, 0,
		1, 1, 0, 1, 0, 0,
	})

	// Create the right-hand side vector b)
	b := mat.NewVecDense(4, []float64{3, 5, 4, 7})

	// Initialize the result vector x
	var x mat.VecDense

	// Solve the system A * x = b
	// The SolveVec function handles the necessary linear algebra algorithms
	if err := x.SolveVec(A, b); err != nil {
		fmt.Printf("Error solving the system: %v\n", err)
		return
	}

	// Print the solution vector x
	fmt.Printf("Solution x: %v\n", mat.Formatted(&x, mat.Prefix(""), mat.Squeeze()))
	// Expected output: Solution x: [1 2]
}
