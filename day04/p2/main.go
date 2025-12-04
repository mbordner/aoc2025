package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/file"
)

const (
	Empty   = '.'
	Paper   = '@'
	Removed = 'x'
)

func main() {
	lines, _ := file.GetLines("../data.txt")
	grid := common.ConvertGrid(lines)

	removed := make(common.Positions, 0, len(grid)*len(grid[0]))

	for {
		accessiblePaperRolls := make(common.Positions, 0, len(grid)*len(grid[0]))

		for j := 0; j < len(grid); j++ {
			for i := 0; i < len(grid[j]); i++ {
				if grid[j][i] == Paper {
					if countNeighborPositionPaperRolls(&grid, i, j) < 4 {
						accessiblePaperRolls = append(accessiblePaperRolls, common.Pos{Y: j, X: i})
					}
				}
			}
		}

		if len(accessiblePaperRolls) > 0 {
			for _, p := range accessiblePaperRolls {
				grid[p.Y][p.X] = Removed
				removed = append(removed, p)
			}
		} else {
			break
		}
	}

	fmt.Println(len(removed))

}

func countNeighborPositionPaperRolls(grid *common.Grid, x, y int) int {
	paperCount := 0
	adj := grid.Neighbors(x, y)
	for _, p := range adj {
		if (*grid)[p.Y][p.X] == Paper {
			paperCount++
		}
	}
	return paperCount
}
