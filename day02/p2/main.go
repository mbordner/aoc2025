package main

// package main we'll go through a list of id (number value) ranges, and find invalid ids
// an invalid id is a number that has repeating number values, e.g. 11, 1212, 123123123
// it will print out the sum of all the invalid id values in the ranges examined

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/files"
	"strconv"
	"strings"
)

func main() {
	idRanges, _ := files.GetContent("../data.txt")

	sum := uint64(0)

	for _, idRange := range strings.Split(string(idRanges), ",") {
		tokens := strings.Split(idRange, "-")
		start, _ := strconv.Atoi(strings.TrimSpace(tokens[0]))
		end, _ := strconv.Atoi(strings.TrimSpace(tokens[1]))

		for i := start; i <= end; i++ {
			strVal := fmt.Sprintf("%d", i)
			h := len(strVal) / 2 // split the string value into two tokens, and store the length of one half

			for l := h; l > 0; l-- { // check token lengths of h down to 1
				if len(strVal)%l == 0 {
					token := strVal[0:l] // store first token
					allMatched := true
					for t := l; t < len(strVal); t += l { // check each consecutive token of same length l, to see if it matches token
						checkToken := strVal[t : t+l]
						if token != checkToken {
							allMatched = false // we can break out if we don't find a match
							break
						}
					}
					if allMatched {
						sum += uint64(i) // if all matched, this is an invalid id, add the value to sum and break
						break            // we don't need to check any other token length when found to be invalid
					}
				}
			}
		}

	}

	fmt.Println(sum)

}
