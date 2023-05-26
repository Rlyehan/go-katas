package main

import (
	"fmt"
	"math/rand"
	"time"
)

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	middle := len(arr) / 2
	left := mergeSort(arr[:middle])
	right := mergeSort(arr[middle:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	
	for len(left) > 0 || len(right) > 0 {
		if len(left) == 0 {
			return append(result, right...)
		}
		if len(right) == 0 {
			return append(result, left...)
		}
		if left[0] <= right[0] {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}
	
	return result
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
		mergeSort(numbers)
		duration := time.Since(start)
		totalTime += duration
	}

	averageTime := totalTime / time.Duration(iterations)
	fmt.Printf("Average time for mergeSort over %d iterations is: %s\n", iterations, averageTime)
}
