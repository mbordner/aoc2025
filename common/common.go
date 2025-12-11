package common

import (
	"fmt"
	"regexp"
	"strconv"
)

func PopulateStringCombinationsAtLength(results map[string]bool, pickChars string, prefix string, length int) {
	if length == 0 {
		results[prefix] = true
		return
	}

	for i := 0; i < len(pickChars); i++ {
		nextPrefix := prefix + string(pickChars[i])
		PopulateStringCombinationsAtLength(results, pickChars, nextPrefix, length-1)
	}
}

type Anything interface{}

func GetPairSets[T Anything](elements []T) [][]T {
	values := make([][]T, 0, len(elements)*(len(elements)-1)/2)
	for i := 0; i < len(elements)-1; i++ {
		for j := i + 1; j < len(elements); j++ {
			values = append(values, []T{elements[i], elements[j]})
		}
	}
	return values
}

func CartesianProduct[T any](sets [][]T) [][]T {
	result := [][]T{{}}

	for _, set := range sets {
		temp := [][]T{}
		for _, x := range set {
			for _, r := range result {
				temp = append(temp, append(r, x))
			}
		}
		result = temp
	}

	return result
}

func FilterArray[T comparable](array []T, exclude []T) []T {
	ex := make(map[T]bool)
	for _, x := range exclude {
		ex[x] = true
	}

	values := make([]T, 0, len(array))

	for _, v := range array {
		if _, e := ex[v]; !e {
			values = append(values, v)
		}
	}

	return values
}

type IntNumber interface {
	int | int32 | int64
}

func Min[T IntNumber](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T IntNumber](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Abs[T IntNumber](v T) T {
	if v >= 0 {
		return v
	}
	return T(-1) * v
}

type Grid [][]byte
type Pos struct {
	Y int
	X int
}

func (p Pos) Add(o Pos) Pos {
	return Pos{p.Y + o.Y, p.X + o.X}
}

func (p Pos) String() string {
	return fmt.Sprintf("{%d,%d}", p.X, p.Y)
}

type Positions []Pos

func (p Positions) ExtentsArea() uint64 {
	min, max := p.Extents()
	dy := max.Y - min.Y + 1
	dx := max.X - min.X + 1
	return uint64(dy * dx)
}

func (p Positions) Extents() (Pos, Pos) {
	var min, max Pos = p[0], p[0]
	for i := 1; i < len(p); i++ {
		o := p[i]
		if o.Y < min.Y {
			min.Y = o.Y
		}
		if o.X < min.X {
			min.X = o.X
		}
		if o.Y > max.Y {
			max.Y = o.Y
		}
		if o.X > max.X {
			max.X = o.X
		}
	}
	return min, max
}

func (g Grid) Neighbors(x, y int) Positions {
	positions := make(Positions, 0, 8)

	for j := y - 1; j <= y+1; j++ {
		for i := x - 1; i <= x+1; i++ {
			if !(j == y && i == x) { // don't add x,y
				if j >= 0 && j < len(g) && i >= 0 && i < len(g[j]) {
					positions = append(positions, Pos{j, i})
				}
			}
		}
	}

	return positions
}

func (g Grid) NeighborPositions(p Pos) Positions {
	return g.Neighbors(p.X, p.Y)
}

func (g Grid) Contains(x, y int) bool {
	if y >= 0 && y < len(g) && x >= 0 && x < len(g[y]) {
		return true
	}
	return false
}

func (g Grid) ContainsPos(p Pos) bool {
	return g.Contains(p.X, p.Y)
}

func (g Grid) Print() {
	for _, row := range g {
		fmt.Println(string(row))
	}
}

func ConvertGrid(lines []string) Grid {
	grid := make(Grid, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}
	return grid
}

type PrevLinkState[S comparable, Action any] struct {
	prev   S
	action Action
}
type PreviousState[S comparable, Action any] map[S]PrevLinkState[S, Action]

func (ps PreviousState[S, Action]) Link(state S, prev S, action Action) {
	ps[state] = PrevLinkState[S, Action]{prev: prev, action: action}
}

func (ps PreviousState[S, Action]) GetActions(start S, goal S) []PrevLinkState[S, Action] {
	actions := []PrevLinkState[S, Action]{{}}
	for p := ps[goal]; p.prev != start; p = ps[p.prev] {
		actions = append([]PrevLinkState[S, Action]{p}, actions...)
	}
	return actions
}

type VisitedState[S comparable, V any] map[S]V

func (vs VisitedState[S, V]) Has(s S) bool {
	_, e := vs[s]
	return e
}

func (vs VisitedState[S, V]) Get(s S) V {
	return vs[s]
}

func (vs VisitedState[S, V]) Remove(s S) {
	delete(vs, s)
}

func (vs VisitedState[S, V]) Set(s S, v V) {
	vs[s] = v
}

type Queue[T comparable] []T

func (q *Queue[T]) Enqueue(s T) {
	*q = append(*q, s)
}

func (q *Queue[T]) Empty() bool {
	return len(*q) == 0
}

func (q *Queue[T]) Dequeue() *T {
	if !q.Empty() {
		s := (*q)[0]
		*q = (*q)[1:]
		return &s
	}
	return nil
}

type PosContainer map[Pos]bool

func (v PosContainer) Has(p Pos) bool {
	if b, e := v[p]; e {
		return b
	}
	return false
}

type PosMapper[T any] map[Pos]T

func (v PosMapper[T]) Has(p Pos) bool {
	_, e := v[p]
	return e
}

type PosLinker map[Pos]Pos

func Filter[T comparable](values []T, value T) []T {
	vs := make([]T, 0, len(values))
	for _, v := range values {
		if v != value {
			vs = append(vs, v)
		}
	}
	return vs
}

func Dedupe[T comparable](values []T) []T {
	vm := make(map[T]bool)
	for _, v := range values {
		vm[v] = true
	}
	vs := make([]T, 0, len(values))
	for v := range vm {
		vs = append(vs, v)
	}
	return vs
}

var (
	ReDigits = regexp.MustCompile(`\d+`)
)

type Ints interface {
	int | int32 | int64 | uint | uint32 | uint64
}

func IntVals[T Ints](strVals string) []T {
	tokens := ReDigits.FindAllString(strVals, -1)
	vals := make([]T, len(tokens), len(tokens))
	for i, t := range tokens {
		v, _ := strconv.ParseInt(t, 10, 64)
		vals[i] = T(v)
	}
	return vals
}
