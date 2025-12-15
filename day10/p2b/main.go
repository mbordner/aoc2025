package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/array"
	"github.com/mbordner/aoc2025/common/expression"
	"github.com/mbordner/aoc2025/common/files"
	"github.com/mbordner/aoc2025/common/matrices"
)

var (
	reLine   = regexp.MustCompile(`^\[([.#]+)\]\s*((?:\([\d|,]+\)\s*)+)\{([\d|,]+)\}\s*$`)
	reButton = regexp.MustCompile(`\(([\d|,]+)\)`)
)

// 10594 your answer is too low
// 12071 is too low also :(
// too high: 23115
// not right answer: 16103
// 16543 is not right
// ? 14423
// 14316 not right
// 16613
func main() {
	machines := getData("../data.txt")

	sum := 0
	for _, m := range machines {
		sum += m.fewestPressesToConfigure()
	}

	fmt.Println(sum)
}

type Machine struct {
	diagramStr string
	buttonStrs []string
	joltsStr   string

	buttons [][]int
	jolts   []int
}

func NewMachine(diagramStr string, buttonStrs []string, joltsStr string) *Machine {
	m := &Machine{diagramStr: diagramStr, buttonStrs: buttonStrs, joltsStr: joltsStr}

	m.buttons = make([][]int, len(buttonStrs))

	for b, s := range buttonStrs {
		vals := common.IntVals[int](s)
		m.buttons[b] = vals
	}

	m.jolts = common.IntVals[int](joltsStr)

	return m
}

func (m *Machine) getCoefficients() [][]int64 {
	r := len(m.jolts)
	c := len(m.buttons)

	coefficients := make([][]int64, r)
	for j := range m.jolts {
		coefficients[j] = make([]int64, c)
		for b, button := range m.buttons {
			if slices.Contains(button, j) {
				coefficients[j][b]++
			}
		}
	}

	return coefficients
}

func (m *Machine) getAugmentedMatrix() [][]int64 {
	coefficients := m.getCoefficients()
	for j := range coefficients {
		coefficients[j] = append(coefficients[j], int64(m.jolts[j]))
	}
	return coefficients
}

func (m *Machine) fewestPressesToConfigure() int {

	matrix := m.getAugmentedMatrix()
	//fmt.Println(matrix)
	matrixRREF := matrices.ToIntegerReducedEchelonForm(matrix)
	//fmt.Println(matrixRREF)

	// remove zero rows
	for len(matrixRREF) > 0 {
		lastRow := matrixRREF[len(matrixRREF)-1]
		if array.SumNumbers(lastRow) == 0 {
			matrixRREF = matrixRREF[0 : len(matrixRREF)-1]
		} else {
			break
		}
	}

	freeVariableColumns := matrices.FindFreeVariables(matrixRREF)

	searchVariables := make(map[int]string)
	for _, v := range freeVariableColumns {
		searchVariables[v] = fmt.Sprintf("b%d", v)
	}

	expressions := make([]*expression.Parser, len(m.buttons))
	for e := 0; e < len(matrixRREF); e++ {
		tokens := make([]string, 0, len(matrixRREF[e]))
		tokens = append(tokens, fmt.Sprintf("%d", matrixRREF[e][len(matrixRREF[e])-1]))

		pivot := -1

		for b := 0; b < len(matrixRREF[e])-1; b++ {
			if matrixRREF[e][b] != 0 {
				if pivot == -1 {
					pivot = b
				} else {
					token := fmt.Sprintf("b%d", b)
					if matrixRREF[e][b] != 1 {
						token = fmt.Sprintf("(%s * %d)", token, matrixRREF[e][b])
					}
					tokens = append(tokens, token)
				}

			}
		}

		expr := strings.Join(tokens, " - ")
		if matrixRREF[e][pivot] != 1 {
			expr = fmt.Sprintf("(%s) / %d", expr, matrixRREF[e][pivot])
		}
		var err error
		expressions[pivot], err = expression.NewParser(expr)
		if err != nil {
			panic(err)
		}

	}

	if len(searchVariables) == 0 {
		presses := make([]int64, len(m.buttons))
		input := make(map[string]int64)
		for b := len(m.buttons) - 1; b >= 0; b-- {
			if expressions[b] != nil {
				buttonVar := fmt.Sprintf("b%d", b)
				var err error
				input[buttonVar], err = expressions[b].Eval(input)
				if err != nil {
					panic(err)
				}
				presses[b] = input[buttonVar]
			}
		}
		for _, p := range presses {
			if p < 0 {
				panic("negative press")
			}
		}
		minPresses := int(array.SumNumbers(presses))
		fmt.Println("to get to: {", m.joltsStr, "} takes: ", minPresses, "[", presses, "] searched:", 0)
		return minPresses
	}

	minPresses := int64(-1)
	var minButtonPresses []int64

	searchIndexes := make([]int, 0, len(searchVariables))
	for v := range searchVariables {
		searchIndexes = append(searchIndexes, v)
	}

	queue := make(common.Queue[string], 0, 200)
	visited := make(common.VisitedState[string, bool])

	startPresses := make([]int64, len(freeVariableColumns))
	start := common.StrFromVals[int64](startPresses)

	queue.Enqueue(start)
	visited.Set(start, true)

	searched := int64(0)

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		searched++
		validClickCount := true

		vals := common.IntVals[int64](cur)
		input := make(map[string]int64)
		for i, v := range vals {
			b := searchIndexes[i]
			input[fmt.Sprintf("b%d", b)] = v
		}

		for b := len(m.buttons) - 1; b >= 0; b-- {
			if expressions[b] != nil {
				buttonVar := fmt.Sprintf("b%d", b)
				var err error
				input[buttonVar], err = expressions[b].Eval(input)
				if err != nil || input[buttonVar] < 0 {
					validClickCount = false
					break
				}
			}
		}

		if validClickCount {
			presses := make([]int64, len(m.buttons))
			for b := len(m.buttons) - 1; b >= 0; b-- {
				buttonVar := fmt.Sprintf("b%d", b)
				presses[b] = input[buttonVar]
			}

			pressesSum := array.SumNumbers(presses)
			if minPresses == -1 || pressesSum < minPresses {
				minPresses = pressesSum
				minButtonPresses = presses
			}
		}

		for i, v := range vals {
			nextVals := common.CloneVals(vals)
			nextVals[i] = v + 1

			if int(nextVals[i]) > slices.Max(m.jolts) {
				continue
			}
			next := common.StrFromVals(nextVals)
			if !visited.Has(next) {
				visited.Set(next, true)
				queue.Enqueue(next)
			}
		}

	}

	fmt.Println("to get to: {", m.joltsStr, "} takes: ", minPresses, "[", minButtonPresses, "] searched:", searched)
	return int(minPresses)
}

func getData(filename string) []*Machine {
	lines := files.MustGetLines(filename)

	mj := 0

	machines := make([]*Machine, 0, len(lines))
	for _, line := range lines {
		matches := reLine.FindStringSubmatch(line)
		if len(matches) == 4 {
			diagram := matches[1]
			jolts := matches[3]
			js := common.IntVals[int](jolts)
			for _, j := range js {
				if j > mj {
					mj = j
				}
			}
			buttonMatches := reButton.FindAllStringSubmatch(matches[2], -1)
			buttons := make([]string, 0, len(buttonMatches))
			for _, match := range buttonMatches {
				buttons = append(buttons, match[1])
			}
			machines = append(machines, NewMachine(diagram, buttons, jolts))
		}
	}

	return machines
}
