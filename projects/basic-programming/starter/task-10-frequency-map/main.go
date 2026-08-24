package main

import (
	"fmt"
	"sort"
)

// FrequencyMap returns a map of each word to its occurrence count.
func FrequencyMap(words []string) map[string]int {
	// TODO: implement this function
	_ = words
	return make(map[string]int)
}

// TopN returns the n most frequent words, sorted by frequency desc, then alphabetically.
func TopN(freq map[string]int, n int) []string {
	// TODO: implement this function
	_ = freq
	_ = n
	_ = sort.Slice
	return nil
}

func main() {
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	freq := FrequencyMap(words)
	fmt.Println(freq["apple"])  // 3
	top := TopN(freq, 2)
	fmt.Println(top)  // [apple banana]
}
