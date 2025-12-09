package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
	"github.com/mbordner/aoc2025/common/geom"
	"slices"
	"strconv"
	"strings"
)

const (
	RED   = '#'
	GREEN = 'X'
	Empty = '.'
)

// 186362000 too low

func main() {
	reds := getPositions("../data.txt")
	minRedsExtent, maxRedsExtent := reds.Extents()
	//dy, dx, grid := getGrid(ps)
	//fmt.Println(dy, dx, len(grid))
	//grid.Print()

	numReds := len(reds)
	border := getBorder(reds)

	reds = append(reds, reds[0])
	reds = append(reds, reds[1])

	maxArea := uint64(0)

	for r := 0; r < numReds; r++ {
		lineArea := (common.Positions{reds[r], reds[r+1]}).ExtentsArea()
		if lineArea > maxArea {
			maxArea = lineArea
		}
		corners := common.Positions{reds[r], reds[r+1], reds[r+2]}
		missing := completeRectangle(corners)
		if inPolygon(border, corners, missing, minRedsExtent, maxRedsExtent) {
			corners = append(corners, missing)
			area := corners.ExtentsArea()
			if area > maxArea {
				maxArea = area
			}
		}
	}

	fmt.Println(maxArea)

}

func inPolygon(border common.PosMapper[byte], corners common.Positions, missing, minE, maxE common.Pos) bool {
	if border.Has(missing) {
		return true
	}

	dir := getDirVector(missing, corners[0])
	crossedBorder := 0

	for p := missing.Add(dir); inExtents(p, minE, maxE); p = p.Add(dir) {
		if border.Has(p) {
			crossedBorder++
		}
	}

	if crossedBorder%2 == 0 {
		return false
	}

	return true
}

func inExtents(p, minE, maxE common.Pos) bool {
	if p.X >= minE.X && p.X <= maxE.X {
		if p.Y >= minE.Y && p.Y <= maxE.Y {
			return true
		}
	}
	return false
}

func completeRectangle(ps common.Positions) common.Pos {
	t1 := common.Pos{X: ps[0].X, Y: ps[2].Y}
	t2 := common.Pos{X: ps[2].X, Y: ps[0].Y}
	if slices.Contains(ps, t2) {
		return t1
	}
	return t2
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
