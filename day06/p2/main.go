package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	problems := getProblems("../data.txt")

	sum := int64(0)

	for _, p := range problems {
		sum += p.eval()
	}

	fmt.Println(sum)
}

type problem struct {
	values    []int64
	operation string
}

func (p *problem) eval() int64 {
	val := p.values[0]
	for _, v := range p.values[1:] {
		switch p.operation {
		case "+":
			val += v
		case "*":
			val *= v
		}
	}
	return val
}

func getProblems(filename string) []*problem {
	reWhitespace := regexp.MustCompile(`\s`)

	lines := files.MustGetLines(filename)

	operations := strings.ReplaceAll(lines[len(lines)-1], " ", "")
	lines = lines[0 : len(lines)-1]

	grid := common.ConvertGrid(lines)
	maxCols := 0
	for _, row := range grid {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	for r := range grid {
		for len(grid[r]) < maxCols {
			grid[r] = append(grid[r], ' ')
		}
	}

	problems := make([]*problem, len(operations))
	for i, op := range operations {
		problems[i] = &problem{operation: string(op), values: make([]int64, 0, len(operations))}
	}

	col := len(grid[0]) - 1
	p := len(problems) - 1
	for col >= 0 {
		val := make([]byte, 0, len(grid))
		for r := 0; r < len(grid); r++ {
			b := grid[r][col]
			if !reWhitespace.MatchString(string(b)) {
				val = append(val, b)
			}
		}
		col--

		if len(val) == 0 {
			p--
		} else {
			num, _ := strconv.ParseInt(string(val), 10, 64)
			problems[p].values = append(problems[p].values, num)
		}
	}

	return problems
}
