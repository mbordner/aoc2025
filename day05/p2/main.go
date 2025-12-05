package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/files"
	"github.com/mbordner/aoc2025/common/ranges"
	"strconv"
	"strings"
)

func main() {
	intRanges := ranges.Collection[uint64]{}

	for _, line := range files.MustGetLines("../data.txt") {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "-")
		l, _ := strconv.ParseInt(tokens[0], 10, 64)
		r, _ := strconv.ParseInt(tokens[1], 10, 64)
		_, _ = intRanges.Add(uint64(l), uint64(r))
	}

	fmt.Println(intRanges.Len())
}
