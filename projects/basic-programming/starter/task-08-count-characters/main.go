package main

import "fmt"

// CountCharacters counts how many times each character appears in s.
func CountCharacters(s string) map[rune]int {
	// TODO: implement this function
	_ = s
	return make(map[rune]int)
}

func main() {
	counts := CountCharacters("hello")
	fmt.Println(counts['l'])  // 2
	fmt.Println(counts['h'])  // 1
}
