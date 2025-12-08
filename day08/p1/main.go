package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/collection"
	"github.com/mbordner/aoc2025/common/files"
	"github.com/mbordner/aoc2025/common/geom"
	"sort"
	"strconv"
	"strings"
)

func main() {
	points := getPoints("../data.txt")

	pairs := common.GetPairSets(points)
	sort.Slice(pairs, func(i, j int) bool {
		d1 := pairs[i][0].Distance(pairs[i][1])
		d2 := pairs[j][0].Distance(pairs[j][1])
		return d1 < d2
	})

	circuits := make([]*collection.Set[geom.Pos[int64]], 0, 10)

	for i := 0; i < 1000; i++ {
		pair := pairs[i]

		inCircuits := make([]*collection.Set[geom.Pos[int64]], 0, 10)
		outCircuits := make([]*collection.Set[geom.Pos[int64]], 0, 10)

		for _, circuit := range circuits {
			if circuit.Contains(pair[0]) || circuit.Contains(pair[1]) {
				inCircuits = append(inCircuits, circuit)
			} else {
				outCircuits = append(outCircuits, circuit)
			}
		}

		circuits = outCircuits
		if len(inCircuits) > 0 {
			for _, circuit := range inCircuits[1:] {
				inCircuits[0].Merge(circuit)
			}
			inCircuits[0].Add(pair[0])
			inCircuits[0].Add(pair[1])
			circuits = append(circuits, inCircuits[0])
		} else {
			set := collection.NewSet[geom.Pos[int64]]()
			set.Add(pair[0])
			set.Add(pair[1])
			circuits = append(circuits, set)
		}
	}

	sort.Slice(circuits, func(i, j int) bool {
		return circuits[j].Len() < circuits[i].Len()
	})

	sizes := int64(1)

	for i := 0; i < 3; i++ {
		sizes *= int64(circuits[i].Len())
	}

	fmt.Println(sizes)
}

func getPoints(filename string) geom.Positions[int64] {
	lines := files.MustGetLines(filename)
	points := make(geom.Positions[int64], len(lines))
	for i, line := range lines {
		tokens := strings.Split(line, ",")
		x, _ := strconv.ParseInt(tokens[0], 10, 64)
		y, _ := strconv.ParseInt(tokens[1], 10, 64)
		z, _ := strconv.ParseInt(tokens[2], 10, 64)
		points[i] = geom.Pos[int64]{X: x, Y: y, Z: z}
	}
	return points
}
