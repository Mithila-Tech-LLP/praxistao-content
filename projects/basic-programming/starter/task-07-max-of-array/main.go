package main

import "fmt"

// MaxOfArray returns the largest integer in nums.
// Assumes nums is non-empty.
func MaxOfArray(nums []int) int {
	// TODO: implement this function
	_ = nums
	return 0
}

func main() {
	fmt.Println(MaxOfArray([]int{3, 1, 4, 1, 5, 9, 2, 6}))  // 9
	fmt.Println(MaxOfArray([]int{-3, -1, -4}))               // -1
}
