package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
	"slices"
)

const (
	CharStart    = 'S'
	CharSplitter = '^'
	CharLine     = '|'
)

func main() {
	grid, _ := getData("../data.txt")

	splits := 0

	for y := 1; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if slices.Contains([]byte{CharStart, CharLine}, grid[y-1][x]) {
				if grid[y][x] == CharSplitter {
					splits++
					grid[y][x-1] = CharLine
					grid[y][x+1] = CharLine
				} else {
					grid[y][x] = CharLine
				}
			}
		}
	}

	fmt.Println(splits)
}

func getData(filename string) (common.Grid, common.Pos) {

	grid := common.ConvertGrid(files.MustGetLines(filename))
	var start common.Pos

search:
	for y, row := range grid {
		for x, c := range row {
			if c == CharStart {
				start = common.Pos{X: x, Y: y}
				break search
			}
		}
	}

	return grid, start
}
