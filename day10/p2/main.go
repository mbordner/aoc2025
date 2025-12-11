package main

import (
	"fmt"
	"regexp"

	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/datastructure"
	"github.com/mbordner/aoc2025/common/files"
)

var (
	reLine   = regexp.MustCompile(`^\[([.#]+)\]\s*((?:\([\d|,]+\)\s*)+)\{([\d|,]+)\}\s*$`)
	reButton = regexp.MustCompile(`\(([\d|,]+)\)`)
)

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

func (m *Machine) fewestPressesToConfigure() int {

	goal := m.joltsStr
	start := common.StrFromVals(make([]int, len(m.jolts)))

	queue := make(datastructure.PriorityQueue[string], 0, 100)
	visited := make(common.VisitedState[string, int])
	prev := make(common.PreviousState[string, []int])

	queue.PushItem(start, 0)
	visited.Set(start, 0)

	var actions []common.PrevLinkState[string, []int]

	for queue.Len() > 0 {
		cur, _ := queue.PopItem()
		if cur == goal {
			actions = prev.GetActions(start, goal)
			break
		} else {
			curVals := common.IntVals[int](cur)

		nextButton:
			for _, button := range m.buttons {
				nextVals := common.CloneVals(curVals)
				for _, b := range button {
					nextVals[b]++
					if nextVals[b] > m.jolts[b] {
						continue nextButton
					}
				}
				next := common.StrFromVals(nextVals)
				if !visited.Has(next) {
					visited.Set(next, visited.Get(cur)+1)
					prev.Link(next, cur, button)
					queue.PushItem(next, visited.Get(cur)+1)
				}
			}

		}
	}

	return len(actions)
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
