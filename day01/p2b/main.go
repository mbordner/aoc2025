package main

// package main will count the number of times a vault lock dial passes (or lands)
// on zero given a list of turns with direction and value.  the dial starts at 50
// and has valid values from 0-99 (inclusive) on the dial.

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/file"
	"strconv"
	"strings"
)

const (
	Values = 100 // number of values on the dial, 0-99
)

func main() {
	replacer := strings.NewReplacer("R", "", "L", "-") // change R to pos int, L to neg int
	lines, _ := file.GetLines("../data.txt")
	dial := 50 // start dial at 50
	zeros := 0 // zero value counter (count of dial landing on, or passes by value 0)
	for _, line := range lines {
		val, _ := strconv.Atoi(replacer.Replace(line))
		zeros += common.Abs(val) / Values // if > Values, add num of full turns to zero counter
		d := dial + (val % Values)        // add remainder after removing full turns, may be out of range
		if d < 0 {
			d = Values + d // adjust back to in range of values
			if dial != 0 { // if we were stopped on zero on previous turn, it was already counted
				zeros++ // we passed by zero going left
			}
		} else if d >= Values {
			d -= Values // adjust back to in range of values
			zeros++     // we passed by zero going right
		} else if d == 0 {
			zeros++ // we landed on zero
		}
		dial = d // set new dial position
	}
	fmt.Println("number of times we passed or landed on zero:", zeros)
}
