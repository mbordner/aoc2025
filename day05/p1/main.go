package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/files"
	"github.com/mbordner/aoc2025/common/ranges"
	"regexp"
	"strconv"
)

var (
	reRange = regexp.MustCompile(`^(\d+)-(\d+)$`)
)

func main() {
	fresh, ingredients := getRangeCollection("../data.txt")

	freshCount := 0

	pairs := fresh.ValuePairs()

	for _, ingredient := range ingredients {
		for i := 0; i < len(pairs); i += 2 {
			if ingredient >= pairs[i] && ingredient <= pairs[i+1] {
				freshCount++
				break
			}
		}
	}

	fmt.Println(freshCount)

}

func getRangeCollection(filename string) (*ranges.Collection[int64], []int64) {
	lines, _ := files.GetLines(filename)

	valid := ranges.Collection[int64]{}

	check := make([]int64, 0, 100)

	for i, line := range lines {
		if line == "" {
			lines = lines[i+1:]
			break
		}
		matches := reRange.FindStringSubmatch(line)
		if len(matches) == 3 {
			l, _ := strconv.ParseInt(matches[1], 10, 64)
			r, _ := strconv.ParseInt(matches[2], 10, 64)
			valid.Add(l, r)
		}
	}

	for _, line := range lines {
		val, _ := strconv.ParseInt(line, 10, 64)
		check = append(check, val)
	}

	return &valid, check
}
