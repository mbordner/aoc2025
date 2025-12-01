package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/file"
	"strconv"
)

const (
	Digits = 100
)

func main() {
	instructions := getInstructions("../data.txt")
	dial := 50
	zeros := 0
	for _, i := range instructions {
		fullTurns := i.val / Digits
		zeros += fullTurns
		clicks := i.val % Digits
		if i.dir == "L" {
			clicks = -clicks
		}
		d := dial + clicks
		if d < 0 {
			d = Digits + d
			if dial != 0 {
				zeros++
			}
		} else if d >= Digits {
			d -= Digits
			zeros++
		} else if d == 0 {
			zeros++
		}
		dial = d
	}
	fmt.Println(zeros)
}

type instr struct {
	dir string
	val int
}

func getInstructions(filename string) []instr {
	lines, _ := file.GetLines(filename)
	instructions := make([]instr, 0, len(lines))
	for _, line := range lines {
		var i instr
		i.dir = line[0:1]
		i.val, _ = strconv.Atoi(line[1:])
		instructions = append(instructions, i)
	}
	return instructions
}
