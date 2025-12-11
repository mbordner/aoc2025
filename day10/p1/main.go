package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/array/bytes"
	"github.com/mbordner/aoc2025/common/bits"
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

	state   int32
	diagram int32
	buttons []int32
}

func NewMachine(diagramStr string, buttonStrs []string, joltsStr string) *Machine {
	m := &Machine{diagramStr: diagramStr, buttonStrs: buttonStrs, joltsStr: joltsStr}

	diagramBytes := bytes.Reverse([]byte(strings.NewReplacer(".", "0", "#", "1").Replace(diagramStr)))
	diagram, _ := strconv.ParseInt(string(diagramBytes), 2, 32)

	m.diagram = int32(diagram)
	m.buttons = make([]int32, len(m.buttonStrs))

	for b, s := range m.buttonStrs {
		vals := common.IntVals[int32](s)
		for _, p := range vals {
			m.buttons[b] = bits.Toggle[int32](m.buttons[b], uint(p))
		}
	}

	return m
}

func (m *Machine) fewestPressesToConfigure() int {

	goal := m.diagram
	start := int32(0)

	queue := make(common.Queue[int32], 0, 100)
	visited := make(common.VisitedState[int32, bool])
	prev := make(common.PreviousState[int32, int32])

	queue.Enqueue(start)
	visited.Set(start, true)

	var actions []common.PrevLinkState[int32, int32]

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		if cur == goal {
			actions = prev.GetActions(start, goal)
			break
		} else {
			for _, b := range m.buttons {
				next := cur ^ b
				if !visited.Has(next) {
					visited.Set(next, true)
					prev.Link(next, cur, b)
					queue.Enqueue(next)
				}
			}
		}
	}

	return len(actions)
}

func getData(filename string) []*Machine {
	lines := files.MustGetLines(filename)

	machines := make([]*Machine, 0, len(lines))
	for _, line := range lines {
		matches := reLine.FindStringSubmatch(line)
		if len(matches) == 4 {
			diagram := matches[1]
			jolts := matches[3]
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
