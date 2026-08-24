package main

import "fmt"

// IsPalindrome returns true if s reads the same forwards and backwards.
func IsPalindrome(s string) bool {
	// TODO: implement this function
	_ = s
	return false
}

func main() {
	fmt.Println(IsPalindrome("racecar"))  // true
	fmt.Println(IsPalindrome("hello"))    // false
}
