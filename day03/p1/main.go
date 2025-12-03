package main

import (
	"fmt"
	"github.com/mbordner/aoc2025/common/file"
	"math"
	"strconv"
)

type Bank []int

func main() {
	lines, _ := file.GetLines("../data.txt")
	jolts := 0

	for _, line := range lines {
		bank := getBank(line)
		i0, b0 := getMaxBattery(bank[0 : len(bank)-1])
		_, b1 := getMaxBattery(bank[i0+1:])
		jolts += b0*10 + b1
	}

	fmt.Println(jolts)
}

func getBank(bankStr string) Bank {
	bank := make(Bank, len(bankStr))
	for i, r := range bankStr {
		bank[i], _ = strconv.Atoi(string(r))
	}
	return bank
}

func getMaxBattery(bank Bank) (int, int) {
	maxIndex := 0
	maxBattery := math.MinInt32
	for i, b := range bank {
		if b > maxBattery {
			maxBattery = b
			maxIndex = i
		}
	}
	return maxIndex, maxBattery
}
