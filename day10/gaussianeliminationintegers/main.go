package main

import (
	"fmt"
	"math"
)

// gcd returns the greatest common divisor of two integers using the Euclidean algorithm.
func gcd(a, b int) int {
	a = int(math.Abs(float64(a)))
	b = int(math.Abs(float64(b)))
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// simplifyRow divides every element in the row by the row's GCD
// to keep the integers from growing excessively large and normalizes
// the leading entry to be positive.
func simplifyRow(row []int) {
	if len(row) == 0 {
		return
	}

	// 1. Find the initial GCD magnitude and the first non-zero index
	rowGCD := 0
	firstNonZeroIndex := -1
	for i, val := range row {
		if val != 0 {
			if firstNonZeroIndex == -1 {
				firstNonZeroIndex = i
				rowGCD = val // Initialize GCD with the first non-zero value
			} else {
				rowGCD = gcd(rowGCD, val)
			}
		}
	}

	// If the entire row is zero, we are done
	if rowGCD == 0 {
		return
	}

	// Ensure rowGCD is the absolute magnitude of the GCD (for division)
	rowGCD = int(math.Abs(float64(rowGCD)))

	// 2. Determine the divisor sign
	// If the first non-zero element is negative, we need a negative divisor
	// to flip the signs of the entire row (making the pivot positive).
	divisor := rowGCD
	if firstNonZeroIndex != -1 && row[firstNonZeroIndex] < 0 {
		divisor = -rowGCD
	}

	// 3. Divide every element by the calculated divisor
	for i := range row {
		row[i] /= divisor
	}
}

// ToIntegerEchelonForm performs Gaussian Elimination on a [][]int matrices
// in place, using only integer arithmetic to prevent fractions.
func ToIntegerEchelonForm(A [][]int) {
	r := len(A)
	if r == 0 {
		return
	}
	c := len(A[0])

	pivotRow := 0
	// Iterate through columns (j)
	for j := 0; j < c && pivotRow < r; j++ {

		// 1. Find the first non-zero pivot below or at pivotRow
		maxRow := pivotRow
		for i := pivotRow; i < r; i++ {
			if A[i][j] != 0 {
				maxRow = i
				break // Found a non-zero pivot
			}
		}

		pivotValue := A[maxRow][j]
		if pivotValue == 0 {
			continue // Column is all zeros below pivotRow, move to next column
		}

		// 2. Swap rows to bring the pivot to the current pivotRow
		if maxRow != pivotRow {
			A[pivotRow], A[maxRow] = A[maxRow], A[pivotRow]
			pivotValue = A[pivotRow][j] // Update pivot value after swap
		}

		// Optional: Simplify the pivot row now to keep numbers small
		simplifyRow(A[pivotRow])
		pivotValue = A[pivotRow][j] // Re-read pivot value after simplification

		// 3. Clear entries below the pivot using integer arithmetic
		for i := pivotRow + 1; i < r; i++ {
			targetEntry := A[i][j]

			if targetEntry == 0 {
				continue
			}

			// --- Simplified Cancellation Logic (The Fix) ---
			// We want R_i = (PivotValue * R_i) - (TargetEntry * R_pivot)
			// to clear A[i][j] to zero. We use GCD to reduce the factors.

			commonDivisor := gcd(pivotValue, targetEntry)

			// Scaling factor for Row i (always positive from pivotValue)
			factorI := int(math.Abs(float64(pivotValue))) / commonDivisor

			// Scaling factor for Pivot Row (always negative for subtraction)
			factorP := -(int(math.Abs(float64(targetEntry))) / commonDivisor)

			// Adjust factors for original signs to ensure correct cancellation:
			// If the pivot and target have opposite signs, the subtraction operation
			// must become an addition operation to cancel.
			if (pivotValue > 0 && targetEntry < 0) || (pivotValue < 0 && targetEntry > 0) {
				factorP = -factorP // Change the sign back to positive for addition
			}

			// Perform the row operation
			for k := j; k < c; k++ {
				A[i][k] = (factorI * A[i][k]) + (factorP * A[pivotRow][k])
			}

			// Optional: Simplify the modified row i
			simplifyRow(A[i])
		}

		// Move to the next pivot row
		pivotRow++
	}
}

func main() {
	// Matrix defined as [][]int (input)
	A_int0 := [][]int{
		{0, 0, 0, 0, 1, 1, 3},
		{0, 1, 0, 0, 0, 1, 5},
		{0, 0, 1, 1, 1, 0, 4},
		{1, 1, 0, 1, 0, 0, 7},
	}
	if len(A_int0) > 0 {
	}

	// [..#.#.##]
	//0: (0,3,4,6,7)
	//1: (0,1,4,5,6,7)
	//2: (0,1,2,3,4,6)
	//3: (0,1,2,3,5,7)
	//4: (0,1,5,6,7)
	//5: (0,3,4,5,6,7)
	//6: (0)
	//7: (1,2,4)
	//8: (0,2,3,4,5,6)
	//9: (0,1,3,4,5,7)
	//{256,224,42,58,230,230,222,231}
	A_int1 := [][]int{
		/*0*/ {1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 256},
		/*1*/ {0, 1, 1, 1, 1, 0, 0, 1, 0, 1, 224},
		/*2*/ {0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 42},
		/*3*/ {1, 0, 1, 1, 0, 1, 0, 0, 1, 1, 58},
		/*4*/ {1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 230},
		/*5*/ {0, 1, 0, 1, 1, 1, 0, 0, 1, 1, 230},
		/*6*/ {1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 222},
		/*7*/ {1, 1, 0, 1, 1, 1, 0, 0, 0, 1, 231},
	}
	if len(A_int1) > 0 {
	}

	A_int2 := [][]int{
		/*0*/ {0, 1, 1, 1, 39},
		/*1*/ {1, 0, 1, 1, 35},
		/*2*/ {1, 1, 0, 0, 20},
		/*3*/ {1, 1, 1, 0, 31},
		/*4*/ {1, 1, 0, 1, 36},
		/*5*/ {1, 0, 1, 1, 35},
	}
	if len(A_int2) > 0 {
	}

	matrix := A_int2
	r := len(matrix)
	c := len(matrix[0])

	fmt.Printf("Original Matrix A (%dx%d):\n", r, c)
	for _, row := range matrix {
		fmt.Println(row)
	}

	// The function modifies A_int in place, using only standard Go types.
	ToIntegerEchelonForm(matrix)

	fmt.Printf("\nEchelon Form with Integer Coefficients (In-Place Modified):\n")
	for _, row := range matrix {
		fmt.Println(row)
	}
}
