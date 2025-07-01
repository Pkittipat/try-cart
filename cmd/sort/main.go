package main

import (
	"fmt"

	"github.com/pkittipat/try-cart/sort"
)

func main() {
	// testNums := []int{7, 12, 9, 11, 3}
	testNums2 := []int{3, 7, 6, -10, 15, 23, 55, -13}
	// sorted := sort.SelectSort(testNums)
	sorted := sort.MergeSort(testNums2)
	fmt.Println(sorted)
}
