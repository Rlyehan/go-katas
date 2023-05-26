package main

import (
	"fmt"
	"math/rand"
	"time"
)

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	pivot := arr[len(arr)/2]
	left := []int{}
	right := []int{}
	for i := 0; i < len(arr); i++ {
		if i != len(arr)/2 {
			if arr[i] < pivot {
				left = append(left, arr[i])
			} else {
				right = append(right, arr[i])
			}
		}
	}
	return append(append(quickSort(left), pivot), quickSort(right)...)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	N := 1000000
	numbers := make([]int, N)

	for i := range numbers {
		numbers[i] = rand.Intn(N)
	}

	iterations := 5
	var totalTime time.Duration
	for i := 0; i < iterations; i++ {
		start := time.Now()
		quickSort(numbers)
		duration := time.Since(start)
		totalTime += duration
	}

	averageTime := totalTime / time.Duration(iterations)
	fmt.Printf("Average time for quickSort over %d iterations is: %s\n", iterations, averageTime)
}
