package main

// package main, given a grid, find the start position marked as 'S', and then find the number of unique paths to the bottom
// of the grid from 'S'.  empty positions are marked with '.' and we simply drop down through them.
// however, if we reach a splitter marked with '^', the path splits into two paths on the left and right side of the splitter.

import (
	"fmt"
	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/files"
)

const (
	CharStart    = 'S'
	CharEmpty    = '.'
	CharSplitter = '^'
)

// Memo a memoization map to store previously calculated number of paths down through the position
type Memo map[common.Pos]uint64

func main() {
	grid, start := getData("../data.txt") // get the grid, and start position from a file

	fmt.Println(countPaths(make(Memo), start, grid)) // print number of unique paths through the grid from S
}

// countPaths is a recursive function that takes a position and grid, and counts unique paths from p down through the grid
func countPaths(memo Memo, p common.Pos, grid common.Grid) uint64 {

	belowY := p.Y + 1 // store next line position (down) in a variable
	if belowY == len(grid) {
		return 1 // if we dropped off the grid, we return a unique count of 1
	}

	if val, known := memo[p]; known {
		return val // from this position, if we already calculated (it's known) the number of paths, return it
	}

	paths := uint64(0)
	if grid[belowY][p.X] == CharEmpty { // just drop down on empty positions
		paths += countPaths(memo, common.Pos{X: p.X, Y: belowY}, grid)
	} else if grid[belowY][p.X] == CharSplitter { // splitters will split the timeline
		paths += countPaths(memo, common.Pos{X: p.X - 1, Y: belowY}, grid)
		paths += countPaths(memo, common.Pos{X: p.X + 1, Y: belowY}, grid)
	}

	memo[p] = paths // store the calculation to prevent recalculation if we revisit this position

	return paths // return the count of unique paths from p down through the grid
}

func getData(filename string) (common.Grid, common.Pos) {

	grid := common.ConvertGrid(files.MustGetLines(filename))
	var start common.Pos

search:
	for y, row := range grid {
		for x, c := range row {
			if c == CharStart {
				start = common.Pos{X: x, Y: y}
				break search // we were looking for start in the grid, we can break out
			}
		}
	}

	return grid, start
}
