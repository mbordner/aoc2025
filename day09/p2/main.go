package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
	"github.com/mbordner/aoc2025/common/geom"
)

const (
	RED   = '#'
	GREEN = 'X'
	Empty = '.'
)

// 186362000 too low
// 4618516475 too high... 4618516475 [{15930,83509} {84454,16111}] [{15930,16111} {84454,83509}]

func main() {
	reds := getPositions("../data.txt")
	minRedsExtent, maxRedsExtent := reds.Extents()

	border := getBorder(reds)

	//dy, dx, grid := getGrid(reds)
	//fmt.Println(dy, dx, len(grid))
	//grid.Print()

	pairs := common.GetPairSets(reds)
	maxArea := uint64(0)
	var maxPair common.Positions
	var maxPairMissing common.Positions
	fmt.Println("len of pairs:", len(pairs))

	//pairs = [][]common.Pos{{common.Pos{X: 15930, Y: 83509}, common.Pos{X: 84454, Y: 16111}}}
	count := 0
nextPair:
	for p, pair := range pairs {
		if p%100 == 0 {
			fmt.Println("pair:", p, pair)
		}
		a, b := pair[0].X, pair[0].Y
		c, d := pair[1].X, pair[1].Y

		missingCorners := common.Positions{common.Pos{X: a, Y: d}, common.Pos{X: c, Y: b}}
		for _, missing := range missingCorners {
			if !inPolygon(border, missing, minRedsExtent, maxRedsExtent) {
				continue nextPair
			}
		}
		count++
		area := missingCorners.ExtentsArea()
		if area > maxArea {
			corners := common.Positions{pair[0], pair[1], missingCorners[0], missingCorners[1]}
			sort.Slice(corners, func(i, j int) bool {
				if corners[i].X < corners[j].X {
					return true
				}
				if corners[i].Y < corners[j].Y {
					return true
				}
				return false
			})
			corners[2], corners[3] = corners[3], corners[2]
			if borderInPolygon(border, corners, minRedsExtent, maxRedsExtent) {
				maxArea = area
				maxPair = pair
				maxPairMissing = missingCorners
			}
		}

	}

	fmt.Println(count)
	fmt.Println(maxArea, maxPair, maxPairMissing)

}

func borderInPolygon(border common.PosMapper[byte], corners common.Positions, minE, maxE common.Pos) bool {
	if corners[0] == corners[1] && corners[2] == corners[3] {
		return true
	}
	cs := make(common.Positions, len(corners)+1)
	copy(cs, corners)
	cs[4] = corners[0]
	for i := 0; i < len(corners); i++ {
		dir := getDirVector(cs[i], cs[i+1])
		for p := cs[i].Add(dir); p != cs[i+1]; p = p.Add(dir) {
			if !inPolygon(border, p, minE, maxE) {
				return false
			}
		}
	}
	return true
}

var memo = make(map[common.Pos]bool)

func inPolygon(border common.PosMapper[byte], missing, minE, maxE common.Pos) bool {
	if border.Has(missing) {
		return true
	}

	if v, e := memo[missing]; e {
		return v
	}

	dirs := common.Positions{common.Pos{X: 1}, common.Pos{X: -1}, common.Pos{Y: 1}, common.Pos{Y: -1}}

	inside := true

	for _, dir := range dirs {
		hitBorder := false
		for p := missing.Add(dir); inExtents(p, minE, maxE); p = p.Add(dir) {
			v, e := memo[p]
			if (e && v) || border.Has(p) {
				hitBorder = true
				break
			}
		}
		if !hitBorder {
			inside = false
			break
		}
	}

	memo[missing] = inside

	return inside
}

func inExtents(p, minE, maxE common.Pos) bool {
	if p.X >= minE.X && p.X <= maxE.X {
		if p.Y >= minE.Y && p.Y <= maxE.Y {
			return true
		}
	}
	return false
}

func getBorder(reds common.Positions) common.PosMapper[byte] {
	border := make(common.PosMapper[byte])

	ps := make(common.Positions, len(reds)+1)
	copy(ps, reds)
	ps[len(ps)-1] = ps[0]

	for r := 0; r < len(ps)-1; r++ {
		border[ps[r]] = RED
		dir := getDirVector(ps[r], ps[r+1])
		p := common.Pos{X: ps[r].X + dir.X, Y: ps[r].Y + dir.Y}
		for ; p != ps[r+1]; p = p.Add(dir) {
			border[p] = GREEN
		}
	}

	return border
}

func getGrid(reds common.Positions) (int, int, common.Grid) {
	minP, maxP := reds.Extents()
	dy := minP.Y
	dx := minP.X
	h := maxP.Y - minP.Y + 1
	w := maxP.X - minP.X + 1
	grid := make(common.Grid, h)
	for y := 0; y < h; y++ {
		grid[y] = make([]byte, w)
		for x := 0; x < w; x++ {
			grid[y][x] = Empty
		}
	}

	ps := make(common.Positions, len(reds)+1)
	copy(ps, reds)
	ps[len(ps)-1] = ps[0]
	for i := 0; i < len(ps); i++ {
		ps[i].X -= dx
		ps[i].Y -= dy
	}

	for r := 0; r < len(ps)-1; r++ {
		grid[ps[r].Y][ps[r].X] = RED
		dir := getDirVector(ps[r], ps[r+1])
		p := common.Pos{X: ps[r].X + dir.X, Y: ps[r].Y + dir.Y}
		for ; p != ps[r+1]; p = p.Add(dir) {
			grid[p.Y][p.X] = GREEN
		}
	}

	return dy, dx, grid
}

func getDirVector(p0, p1 common.Pos) common.Pos {
	vector := (geom.Pos[int]{X: p1.X - p0.X, Y: p1.Y - p0.Y}).Normalize()
	return common.Pos{X: vector.X, Y: vector.Y}
}

func getPositions(filename string) common.Positions {
	lines := files.MustGetLines(filename)
	ps := make(common.Positions, 0, len(lines))
	for _, line := range lines {
		tokens := strings.Split(line, ",")
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		ps = append(ps, common.Pos{X: x, Y: y})
	}
	return ps
}
