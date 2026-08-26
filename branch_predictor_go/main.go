package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"slices"
)

func generateRandomNumbers(count int, seed1, seed2 uint64) []int {
	r := rand.New(rand.NewPCG(seed1, seed2))

	numbers := make([]int, count)
	for i := range count {
		numbers[i] = r.IntN(256)
	}
	return numbers
}

//go:noinline
func process(li []int) (int, int) {
	sumOfElementsBelow128 := 0
	totalSum := 0

	for _, i := range li {
		if i < 128 {
			sumOfElementsBelow128 += i
		}
		totalSum += i
	}
	return sumOfElementsBelow128, totalSum
}

func main() {
	// By default, we run unsorted
	var sorted bool
	flag.BoolVar(&sorted, "sorted", false, "run with sorted data")
	flag.Parse()

	// We want the same numbers across multiple runs
	randomList := generateRandomNumbers(10000000, 42, 99)
	// We always sort them to guarantee that we are doing the same operations in
	// each run.
	sortedList := slices.Clone(randomList)
	slices.Sort(sortedList)

	// It should be always the same numbers, no matter how many times you run
	// them, sorted or unsorted.
	if sorted {
		fmt.Println(process(sortedList))
	} else {
		fmt.Println(process(randomList))
	}
}
