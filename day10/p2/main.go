package main

import (
	"errors"
	"fmt"
	"regexp"
	"slices"

	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/array"
	"github.com/mbordner/aoc2025/common/files"
	"gonum.org/v1/gonum/mat"
)

var (
	reLine   = regexp.MustCompile(`^\[([.#]+)\]\s*((?:\([\d|,]+\)\s*)+)\{([\d|,]+)\}\s*$`)
	reButton = regexp.MustCompile(`\(([\d|,]+)\)`)
)

// 10594 your answer is too low
// 12071 is too low also :(
func main() {
	machines := getData("../data.txt")

	missing := common.IntVals[int](`74,103,138,142`)

	skipped := make([]int, 0, 10)
	sum := 13249
	for _, i := range missing {
		m := machines[i]
		val, err := m.fewestPressesToConfigure()
		if err != nil {
			fmt.Printf("skipped machine %d (line %d)\n", i, i+1)
			skipped = append(skipped, i)
		} else {
			sum += val
			fmt.Println("new sum:", sum, "remove:", i)
		}
	}

	fmt.Println(sum)
	fmt.Println(common.StrFromVals(skipped))
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

func (m *Machine) getSolutionVector() mat.Dense {
	r := len(m.jolts)
	c := len(m.buttons)

	vectorBVals := make([]float64, len(m.jolts))

	coeffs := make([]float64, r*c)
	for j, jolt := range m.jolts {
		vectorBVals[j] = float64(jolt)
		for b, button := range m.buttons {
			index := (j * c) + b
			if slices.Contains(button, j) {
				coeffs[index]++
			}
		}
	}

	matrixA := mat.NewDense(r, c, coeffs)

	vectorB := mat.NewVecDense(r, vectorBVals)

	//fmt.Printf("Matrix A :\n%v\n", mat.Formatted(matrixA, mat.Prefix("  "), mat.Squeeze()))
	//fmt.Printf("Vector B :\n%v\n", mat.Formatted(vectorB, mat.Prefix("  "), mat.Squeeze()))

	var svd mat.SVD
	svd.Factorize(matrixA, mat.SVDFull)

	effectiveRank := svd.Rank(1e-10)

	var x mat.Dense

	// solve the system MatrixA * vectorX = vectorB
	svd.SolveTo(&x, vectorB, effectiveRank)

	return x
}

func (m *Machine) getSolutionVector2() mat.VecDense {
	r := len(m.jolts)
	c := len(m.buttons)

	vectorBVals := make([]float64, len(m.jolts))

	coeffs := make([]float64, r*c)
	for j, jolt := range m.jolts {
		vectorBVals[j] = float64(jolt)
		for b, button := range m.buttons {
			index := (j * c) + b
			if slices.Contains(button, j) {
				coeffs[index]++
			}
		}
	}

	A := mat.NewDense(r, c, coeffs)
	b := mat.NewVecDense(r, vectorBVals)

	var x mat.VecDense

	if err := x.SolveVec(A, b); err != nil {
		panic(err)
	}
	return x
}

func (m *Machine) pressesToJolts(presses []int) (string, []int) {
	jolts := make([]int, len(m.jolts))
	for b, pressCount := range presses {
		for _, j := range m.buttons[b] {
			jolts[j] += pressCount
		}
	}
	return common.StrFromVals(jolts), jolts
}

func (m *Machine) fewestPressesToConfigure() (int, error) {

	solutionVector := m.getSolutionVector()

	goal := m.joltsStr
	startPresses := make([]int, len(m.buttons))
	for b := range startPresses {
		sv := solutionVector.At(b, 0)
		presses := max(0, int(sv)-15)
		startPresses[b] = presses
	}
	_, startJolt := m.pressesToJolts(startPresses)
	for j, v := range startJolt {
		if v > m.jolts[j] {
			return 0, errors.New("jolts too big")
		}
	}
	fmt.Println(goal, startJolt, startPresses)
	start := common.StrFromVals(startPresses)

	queue := make(common.Queue[string], 0, 100)
	visited := make(common.VisitedState[string, []int])

	queue.Enqueue(start)
	visited.Set(start, startPresses)

	var goalPresses []int

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		curPresses := visited.Get(cur)
		curJoltsStr, _ := m.pressesToJolts(curPresses)
		fmt.Printf("{%s} {%s}\n", goal, curJoltsStr)
		if curJoltsStr == goal {
			goalPresses = curPresses
			break
		} else {

		nextButton:
			for b := range m.buttons {
				nextVals := common.CloneVals(curPresses)
				nextVals[b]++
				_, nextJolts := m.pressesToJolts(nextVals)
				for j, v := range nextJolts {
					if v > m.jolts[j] {
						continue nextButton
					}
				}
				next := common.StrFromVals(nextVals)
				if !visited.Has(next) {
					visited.Set(next, nextVals)
					queue.Enqueue(next)
				}
			}

		}
	}

	return array.SumNumbers(goalPresses), nil
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
