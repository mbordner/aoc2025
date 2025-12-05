package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
)

const (
	Empty = '.'
	Paper = '@'
)

func main() {
	lines, _ := files.GetLines("../data.txt")
	grid := common.ConvertGrid(lines)

	accessiblePaperRolls := make(common.Positions, 0, len(grid)*len(grid[0]))

	for j := 0; j < len(grid); j++ {
		for i := 0; i < len(grid[j]); i++ {
			if grid[j][i] == Paper {
				if countAdjPaperRolls(&grid, i, j) < 4 {
					accessiblePaperRolls = append(accessiblePaperRolls, common.Pos{Y: j, X: i})
				}
			}
		}
	}

	fmt.Println(len(accessiblePaperRolls))

}

func countAdjPaperRolls(grid *common.Grid, x, y int) int {
	paperCount := 0
	adj := grid.Neighbors(x, y)
	for _, p := range adj {
		if (*grid)[p.Y][p.X] == Paper {
			paperCount++
		}
	}
	return paperCount
}
