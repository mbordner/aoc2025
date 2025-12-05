package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/files"
	"strconv"
	"strings"
)

func main() {
	idRanges, _ := files.GetContent("../data.txt")

	invalidIds := make(map[string]bool)

	for _, idRange := range strings.Split(string(idRanges), ",") {
		tokens := strings.Split(idRange, "-")
		start, _ := strconv.Atoi(strings.TrimSpace(tokens[0]))
		end, _ := strconv.Atoi(strings.TrimSpace(tokens[1]))

		for i := start; i <= end; i++ {
			strval := fmt.Sprintf("%d", i)
			if len(strval)%2 == 0 {
				firstHalf := strval[0 : len(strval)/2]
				secondHalf := strval[len(strval)/2:]
				if firstHalf == secondHalf {
					invalidIds[strval] = true
				}
			}
		}

	}

	sum := uint64(0)
	for id := range invalidIds {
		val, _ := strconv.ParseInt(id, 10, 64)
		sum += uint64(val)
	}

	fmt.Println(sum)

}
