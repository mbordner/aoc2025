package main

// package main will calculate a sum of max jolts from lines of digits.
// to calculate the max value from a line, we will look at the num of digits expected
// and determine the max value from the line digits keeping the order

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/file"
	"math"
	"strconv"
)

type Digits []int

const NumOfDigits = 12

func main() {
	lines, _ := file.GetLines("../data.txt")
	jolts := uint64(0)

	for _, line := range lines {
		bank := convertStringToDigits(line)

		indexes := make([]int, NumOfDigits)   // store the index locations, mostly needed for a starting point to search the next iteration
		batteries := make([]int, NumOfDigits) // store each next max value from the available digits remaining at the appropriate spot

		for b := 0; b < NumOfDigits; b++ { // loop through up to num digits to search out the max value
			startIndex := 0 // starting index will be zero
			if b > 0 {
				startIndex = indexes[b-1] + 1 // but after first iteration, we search only from 1 past the last digits found location
			}
			indexes[b], batteries[b] = getMaxDigitIndex(bank[startIndex : len(bank)-NumOfDigits+1+b]) // each iteration has to leave enough room for remaining digits needed
			indexes[b] += startIndex                                                                  // getMaxDigitIndex returns the index relative to the slice, so we have to add startIndex value
		}

		// batteries will be an array of digits, representing the largest decimal number we can create preserving order with the max number of digits
		for i, d := len(batteries)-1, 0; i >= 0; i, d = i-1, d+1 { // loop through adding each digit value, multiplied by the ten's place
			jolts += uint64(uint64(math.Pow10(d)) * uint64(batteries[i]))
		}

	}

	fmt.Println(jolts) // sum of all the max jolts
}

// convertStringToDigits converts an integer string to the integer digits and returns the slice of ints
func convertStringToDigits(str string) Digits {
	digits := make(Digits, len(str))
	for i, r := range str {
		digits[i], _ = strconv.Atoi(string(r))
	}
	return digits
}

// getMaxDigitIndex returns the max digit with the lowest index from an integer slice
func getMaxDigitIndex(digits Digits) (int, int) {
	maxIndex := 0
	maxDigit := math.MinInt32
	for i, b := range digits {
		if b > maxDigit {
			maxDigit = b
			maxIndex = i
		}
	}
	return maxIndex, maxDigit
}
