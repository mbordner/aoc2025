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

func main() {
	reds := getPositions("../data.txt")

	border := getBorder(reds)
	horizontals, verticals := getLines(reds)

	pairs := common.GetPairSets(reds)

	maxArea := uint64(0)

nextPair:
	for _, pair := range pairs {
		a, b := pair[0].X, pair[0].Y
		c, d := pair[1].X, pair[1].Y

		missingCorners := common.Positions{common.Pos{X: a, Y: d}, common.Pos{X: c, Y: b}}
		for _, missing := range missingCorners {
			if !inPolygon(border, horizontals, verticals, missing) {
				continue nextPair
			}
		}

		area := missingCorners.ExtentsArea()
		if area > maxArea {
			corners := orderCorners(common.Positions{pair[0], pair[1], missingCorners[0], missingCorners[1]})
			if borderInPolygon(border, reds, horizontals, verticals, corners) {
				maxArea = area
			}
		}

	}

	fmt.Println(maxArea)

}

func orderCorners(corners common.Positions) common.Positions {
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
	return corners
}

func borderInPolygon(border common.PosMapper[byte], reds common.Positions, hls geom.GridLines[int], vls geom.GridLines[int], corners common.Positions) bool {
	if corners[0] == corners[1] && corners[2] == corners[3] {
		return true
	}
	cs := make(common.Positions, len(corners)+1)
	copy(cs, corners)
	cs[4] = corners[0]
	for i := 0; i < len(corners); i++ {
		p0 := geom.Pos[int]{X: cs[i].X, Y: cs[i].Y}
		p1 := geom.Pos[int]{X: cs[i+1].X, Y: cs[i+1].Y}
		line := geom.GridLine[int]{P0: p0, P1: p1}
		dir := line.Direction()
		checkPoints := make(common.Positions, 0, len(reds))
		for _, red := range reds {
			if dir == geom.North || dir == geom.South {
				if red.Y >= min(line.P0.Y, line.P1.Y) && red.Y <= max(line.P0.Y, line.P1.Y) {
					if red.Y-1 >= min(line.P0.Y, line.P1.Y) {
						checkPoints = append(checkPoints, common.Pos{X: line.P0.X, Y: red.Y - 1})
					}
					if red.Y+1 <= max(line.P0.Y, line.P1.Y) {
						checkPoints = append(checkPoints, common.Pos{X: line.P0.X, Y: red.Y + 1})
					}
				}
			} else { // east/west
				if red.X >= min(line.P0.X, line.P1.X) && red.X <= max(line.P0.X, line.P1.X) {
					if red.X-1 >= min(line.P0.X, line.P1.X) {
						checkPoints = append(checkPoints, common.Pos{Y: line.P0.Y, X: red.X - 1})
					}
					if red.X+1 <= max(line.P0.X, line.P1.X) {
						checkPoints = append(checkPoints, common.Pos{Y: line.P0.Y, X: red.X + 1})
					}
				}
			}
		}
		for _, cp := range checkPoints {
			if !inPolygon(border, hls, vls, cp) {
				return false
			}
		}
	}
	return true
}

var memo = make(map[common.Pos]bool)

func inPolygon(border common.PosMapper[byte], hls geom.GridLines[int], vls geom.GridLines[int], p common.Pos) bool {
	if border.Has(p) {
		return true
	}

	if v, e := memo[p]; e {
		return v
	}

	var hd, vd int
	var hl, vl *geom.GridLine[int]

	for line := range vls {
		if p.Y <= max(line.P0.Y, line.P1.Y) && p.Y >= min(line.P0.Y, line.P1.Y) {
			d := max(line.P0.X, p.X) - min(line.P0.X, p.X)
			if vl == nil || d < vd {
				vl = &line
				vd = d
			}
		}
	}

	for line := range hls {
		if p.X <= max(line.P0.X, line.P1.X) && p.X >= min(line.P0.X, line.P1.X) {
			d := max(line.P0.Y, p.Y) - min(line.P0.Y, p.Y)
			if hl == nil || d < hd {
				hl = &line
				hd = d
			}
		}
	}

	inside := true

	if hl == nil || vl == nil {
		inside = false
	} else {
		if vl.Direction() == geom.North {
			if p.X < (*vl).P0.X {
				inside = false
			}
		} else {
			if p.X > (*vl).P0.X {
				inside = false
			}
		}

		if hl.Direction() == geom.East {
			if p.Y < (*hl).P0.Y {
				inside = false
			}
		} else {
			if p.Y > (*hl).P0.Y {
				inside = false
			}
		}
	}

	memo[p] = inside

	return inside
}

func getLines(reds common.Positions) (geom.GridLines[int], geom.GridLines[int]) {
	vertical := make(geom.GridLines[int])
	horizontal := make(geom.GridLines[int])
	ps := make(common.Positions, len(reds)+1)
	copy(ps, reds)
	ps[len(ps)-1] = ps[0]

	for r := 0; r < len(ps)-1; r++ {
		p0 := geom.Pos[int]{X: ps[r].X, Y: ps[r].Y}
		p1 := geom.Pos[int]{X: ps[r+1].X, Y: ps[r+1].Y}
		line := geom.GridLine[int]{P0: p0, P1: p1}
		dir := line.Direction()
		if dir == geom.North || dir == geom.South {
			vertical[line] = dir
		} else {
			horizontal[line] = dir
		}
	}
	return horizontal, vertical
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
