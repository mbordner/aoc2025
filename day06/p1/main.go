package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/files"
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

	lines := files.MustGetLines(filename)

	operations := strings.ReplaceAll(lines[len(lines)-1], " ", "")
	lines = lines[0 : len(lines)-1]
	problems := make([]*problem, len(operations))
	for i, op := range operations {
		problems[i] = &problem{operation: string(op), values: make([]int64, 0, len(lines))}
	}

	for _, line := range lines {
		vals := strings.Fields(line)
		for i, val := range vals {
			num, _ := strconv.ParseInt(val, 10, 64)
			problems[i].values = append(problems[i].values, num)
		}
	}

	return problems
}
