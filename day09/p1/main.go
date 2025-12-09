package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
	"math"
	"strconv"
	"strings"
)

func main() {
	ps := getPositions("../data.txt")

	pairs := common.GetPairSets(ps)
	maxArea := uint64(0)
	var maxPair common.Positions
	for _, pair := range pairs {
		dy := uint64(math.Abs(float64(pair[0].Y - pair[1].Y + 1)))
		dx := uint64(math.Abs(float64(pair[0].X - pair[1].X + 1)))
		area := dy * dx
		if area > maxArea {
			maxArea = area
			maxPair = pair
		}
	}
	fmt.Println(maxArea, maxPair)
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
