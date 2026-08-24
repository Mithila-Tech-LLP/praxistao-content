# Chapter 03: Arrays, Strings & Hashing

Four patterns solve 40% of all array and string problems in interviews. Master these patterns and you will never blank on an array problem again. This chapter teaches each pattern from scratch, shows you how to recognize it, and walks through the classic problems associated with each.

## Table of Contents

1. [The Four Patterns](#1-the-four-patterns)
2. [Sliding Window](#2-sliding-window)
3. [Two Pointers](#3-two-pointers)
4. [Prefix Sum](#4-prefix-sum)
5. [Hashing & Frequency Maps](#5-hashing--frequency-maps)
6. [String Problems in Go](#6-string-problems-in-go)
7. [Classic Interview Problems](#7-classic-interview-problems)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Four Patterns

When you see an array or string problem in an interview, ask yourself:

- **Sliding Window** — "I need to find a contiguous subarray or substring that satisfies a condition."
- **Two Pointers** — "My array is sorted, or I need to find pairs/triplets, or I need to move from both ends toward the middle."
- **Prefix Sum** — "I need to answer range sum queries, or count subarrays with a property."
- **Hashing** — "I need O(1) lookup, or I need to count frequencies, or I need to find if a complement exists."

These are not algorithms you memorize. They are patterns you recognize. Once you see them, the code almost writes itself.

---

## 2. Sliding Window

### The Core Idea

A window is a range [left, right] within the array. You slide it by moving `right` forward to expand the window, and `left` forward to shrink it. The key insight: you avoid the naive O(n²) by reusing computation from the previous window position.

**When to use:** Problems that ask for the longest/shortest subarray or substring with some property (sum, distinct elements, no repeats, etc.).

### Template

```go
func slidingWindow(nums []int, condition func(window []int) bool) int {
    left := 0
    best := 0

    for right := 0; right < len(nums); right++ {
        // 1. Add nums[right] to the window (expand)
        // update your window state here

        // 2. Shrink window from the left if condition violated
        for /* window invalid */ {
            // remove nums[left] from window state
            left++
        }

        // 3. Update answer with current valid window
        best = max(best, right-left+1)
    }
    return best
}
```

### Problem: Longest Substring Without Repeating Characters

**Input:** s = "abcabcbb"
**Output:** 3 ("abc")

```go
func lengthOfLongestSubstring(s string) int {
    // charIndex tracks the last seen index of each character.
    // If a character reappears inside our window, we must move left past it.
    charIndex := make(map[byte]int)
    left := 0
    maxLen := 0

    for right := 0; right < len(s); right++ {
        ch := s[right]

        // If ch was seen inside the current window [left, right),
        // shrink the window from the left to exclude the old occurrence.
        if idx, ok := charIndex[ch]; ok && idx >= left {
            left = idx + 1
        }

        charIndex[ch] = right            // record latest position of ch
        maxLen = max(maxLen, right-left+1) // window size = right - left + 1
    }
    return maxLen
}
// Time: O(n), Space: O(min(m,n)) where m is charset size
```

### Problem: Maximum Sum Subarray of Size K

```go
func maxSumSubarray(nums []int, k int) int {
    // Fixed-size window: always exactly k elements.
    // Build the initial window, then slide one step at a time.
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += nums[i] // build first window
    }

    maxSum := windowSum
    for i := k; i < len(nums); i++ {
        // Add new element on the right, remove old element on the left.
        // This is O(1) per step instead of recomputing the sum from scratch.
        windowSum += nums[i] - nums[i-k]
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    return maxSum
}
// Time: O(n), Space: O(1)
```

### Problem: Minimum Size Subarray Sum

**Input:** target = 7, nums = [2,3,1,2,4,3]
**Output:** 2 (subarray [4,3])

```go
func minSubArrayLen(target int, nums []int) int {
    left := 0
    sum := 0
    minLen := len(nums) + 1 // start with "infinity"

    for right := 0; right < len(nums); right++ {
        sum += nums[right] // expand window

        // Shrink window as long as sum meets the target.
        // We want the shortest valid window, so shrink aggressively.
        for sum >= target {
            minLen = min(minLen, right-left+1)
            sum -= nums[left]
            left++
        }
    }

    if minLen == len(nums)+1 {
        return 0 // no valid subarray found
    }
    return minLen
}
// Time: O(n) — left and right each move at most n times total
// Space: O(1)
```

---

## 3. Two Pointers

### The Core Idea

Use two index variables that move through the array — sometimes both from the left at different speeds (fast/slow), sometimes from opposite ends toward the middle. This turns O(n²) pair-search problems into O(n).

**When to use:**
- Array is sorted and you need pairs/triplets
- You need to remove duplicates or elements in-place
- You need to compare elements from both ends (palindrome check, container problem)

### Template: Opposite Ends

```go
func twoPointers(nums []int) int {
    left, right := 0, len(nums)-1

    for left < right {
        // check condition using nums[left] and nums[right]
        if /* condition met */ {
            // found answer
        } else if /* need larger value */ {
            left++
        } else {
            right--
        }
    }
    return -1
}
```

### Problem: Two Sum II — Sorted Array

```go
// Input: sorted array, find pair that sums to target
func twoSumSorted(nums []int, target int) (int, int) {
    left, right := 0, len(nums)-1

    for left < right {
        sum := nums[left] + nums[right]
        if sum == target {
            return left, right
        } else if sum < target {
            left++  // sum too small, move left pointer right to increase sum
        } else {
            right-- // sum too large, move right pointer left to decrease sum
        }
    }
    return -1, -1
}
// Time: O(n), Space: O(1)
```

### Problem: Container With Most Water

```go
// Find two lines that form a container holding the most water.
func maxArea(height []int) int {
    left, right := 0, len(height)-1
    maxWater := 0

    for left < right {
        // Area = min height * width
        h := min(height[left], height[right])
        w := right - left
        maxWater = max(maxWater, h*w)

        // The shorter line is the bottleneck. Move the pointer on the shorter
        // side inward — moving the taller side can only decrease width without
        // helping height, so it strictly decreases area.
        if height[left] < height[right] {
            left++
        } else {
            right--
        }
    }
    return maxWater
}
// Time: O(n), Space: O(1)
```

### Problem: 3Sum — Find All Triplets Summing to Zero

```go
func threeSum(nums []int) [][]int {
    sort.Ints(nums) // sort first to enable two-pointer
    result := [][]int{}

    for i := 0; i < len(nums)-2; i++ {
        // Skip duplicates for the fixed element
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }

        left, right := i+1, len(nums)-1
        for left < right {
            sum := nums[i] + nums[left] + nums[right]
            if sum == 0 {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                // Skip duplicates for left and right pointers
                for left < right && nums[left] == nums[left+1] { left++ }
                for left < right && nums[right] == nums[right-1] { right-- }
                left++
                right--
            } else if sum < 0 {
                left++
            } else {
                right--
            }
        }
    }
    return result
}
// Time: O(n²) — O(n log n) sort + O(n) for each of n elements = O(n²)
// Space: O(1) extra (not counting output)
```

---

## 4. Prefix Sum

### The Core Idea

A prefix sum array stores the cumulative sum up to each index. The range sum from index i to j is `prefix[j+1] - prefix[i]` in O(1) time, instead of summing each element.

```
nums    = [3, 1, 4, 1, 5, 9]
prefix  = [0, 3, 4, 8, 9, 14, 23]
         (prefix[i] = sum of nums[0..i-1])

Sum from index 2 to 4 (inclusive) = prefix[5] - prefix[2] = 14 - 4 = 10
Check: 4 + 1 + 5 = 10 ✓
```

```go
// Build prefix sum array
func buildPrefix(nums []int) []int {
    prefix := make([]int, len(nums)+1)
    for i, n := range nums {
        prefix[i+1] = prefix[i] + n
    }
    return prefix
}

// Query range sum [l, r] inclusive in O(1)
func rangeSum(prefix []int, l, r int) int {
    return prefix[r+1] - prefix[l]
}
```

### Problem: Subarray Sum Equals K

**Input:** nums = [1, 1, 1], k = 2
**Output:** 2 (subarrays [1,1] appear twice)

```go
// Key insight: if prefix[j] - prefix[i] = k, then sum of nums[i..j-1] = k.
// So we need to count pairs where prefix[j] - k has been seen before.
func subarraySum(nums []int, k int) int {
    // prefixCount tracks how many times each prefix sum has appeared.
    // Start with prefix sum 0 appearing once (empty prefix).
    prefixCount := map[int]int{0: 1}
    sum := 0
    count := 0

    for _, n := range nums {
        sum += n
        // If (sum - k) has appeared before, there are subarrays ending here
        // that sum to k. Each occurrence is a valid subarray.
        count += prefixCount[sum-k]
        prefixCount[sum]++
    }
    return count
}
// Time: O(n), Space: O(n)
```

### Problem: Product of Array Except Self

```go
// No division allowed. Use prefix and suffix products.
func productExceptSelf(nums []int) []int {
    n := len(nums)
    result := make([]int, n)

    // Pass 1: result[i] = product of all elements to the LEFT of i
    result[0] = 1
    for i := 1; i < n; i++ {
        result[i] = result[i-1] * nums[i-1]
    }

    // Pass 2: multiply by product of all elements to the RIGHT of i
    right := 1
    for i := n - 1; i >= 0; i-- {
        result[i] *= right
        right *= nums[i]
    }
    return result
}
// Time: O(n), Space: O(1) extra (output array doesn't count)
```

---

## 5. Hashing & Frequency Maps

### The Core Idea

A hash map gives you O(1) lookup, insertion, and deletion. Use it when you need to:
- Check if a complement/pair exists
- Count frequencies
- Group elements by property
- Track the "first/last seen" index of something

### Pattern: Two Sum with Hash Map

```go
// The classic: find two indices such that nums[i] + nums[j] = target
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int) // value -> index

    for i, n := range nums {
        complement := target - n
        if j, ok := seen[complement]; ok {
            return []int{j, i}
        }
        seen[n] = i
    }
    return nil
}
// Time: O(n), Space: O(n)
```

### Pattern: Frequency Map

```go
// Check if two strings are anagrams
func isAnagram(s, t string) bool {
    if len(s) != len(t) { return false }

    freq := make(map[rune]int)
    for _, ch := range s { freq[ch]++ }
    for _, ch := range t {
        freq[ch]--
        if freq[ch] < 0 { return false }
    }
    return true
}
// For ASCII-only: use [26]int array instead of map — O(1) space, faster

func isAnagramArray(s, t string) bool {
    if len(s) != len(t) { return false }
    var freq [26]int
    for i := 0; i < len(s); i++ {
        freq[s[i]-'a']++
        freq[t[i]-'a']--
    }
    for _, f := range freq {
        if f != 0 { return false }
    }
    return true
}
```

### Pattern: Group Anagrams

```go
// Group strings that are anagrams of each other
func groupAnagrams(strs []string) [][]string {
    groups := make(map[[26]int][]string)

    for _, s := range strs {
        var key [26]int
        for _, ch := range s {
            key[ch-'a']++
        }
        groups[key] = append(groups[key], s)
    }

    result := make([][]string, 0, len(groups))
    for _, group := range groups {
        result = append(result, group)
    }
    return result
}
// Time: O(n * m) where m is max string length
// Space: O(n * m) for storing all strings in groups
```

### Pattern: Longest Consecutive Sequence

```go
// Find the longest consecutive sequence in an unsorted array.
// Input: [100, 4, 200, 1, 3, 2] → Output: 4 (sequence [1,2,3,4])
func longestConsecutive(nums []int) int {
    set := make(map[int]bool)
    for _, n := range nums {
        set[n] = true
    }

    best := 0
    for n := range set {
        // Only start counting from the beginning of a sequence.
        // n is the start of a sequence if n-1 does not exist.
        if !set[n-1] {
            length := 1
            for set[n+length] {
                length++
            }
            best = max(best, length)
        }
    }
    return best
}
// Time: O(n) — each number is visited at most twice (once in outer loop, once in inner)
// Space: O(n) for the set
```

---

## 6. String Problems in Go

### Important Go String Facts for Interviews

```go
// Strings in Go are immutable byte slices (UTF-8 encoded).
// Iterating with range gives runes (Unicode code points).
// Iterating with index gives bytes.

s := "café"
fmt.Println(len(s))          // 5 bytes (é is 2 bytes in UTF-8)
for i, ch := range s {       // i is byte index, ch is rune
    fmt.Printf("%d: %c\n", i, ch)
}
// Output: 0:c  1:a  2:f  3:é  (index 3 for the 2-byte é)

// Convert to rune slice for safe character-level manipulation
runes := []rune(s)
fmt.Println(len(runes)) // 4 characters
```

### Efficient String Building

```go
// WRONG: O(n²) due to string immutability
result := ""
for _, word := range words {
    result += word + " " // creates new string every time
}

// RIGHT: O(n) with strings.Builder
var sb strings.Builder
for _, word := range words {
    sb.WriteString(word)
    sb.WriteByte(' ')
}
result := sb.String()
```

### Palindrome Check

```go
func isPalindrome(s string) bool {
    // Clean string: keep only alphanumeric, lowercase
    runes := []rune{}
    for _, ch := range s {
        if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
            runes = append(runes, unicode.ToLower(ch))
        }
    }

    left, right := 0, len(runes)-1
    for left < right {
        if runes[left] != runes[right] {
            return false
        }
        left++
        right--
    }
    return true
}
```

---

## 7. Classic Interview Problems

Here are 10 problems you must be able to solve cold. Each uses one of the four patterns.

| Problem | Pattern | Complexity |
|---|---|---|
| Longest Substring Without Repeating | Sliding Window | O(n) |
| Minimum Window Substring | Sliding Window | O(n) |
| Maximum Sum Subarray of Size K | Sliding Window | O(n) |
| Two Sum | Hashing | O(n) |
| 3Sum | Two Pointers + Sort | O(n²) |
| Container With Most Water | Two Pointers | O(n) |
| Subarray Sum Equals K | Prefix Sum | O(n) |
| Longest Consecutive Sequence | Hashing | O(n) |
| Group Anagrams | Hashing | O(n*m) |
| Product of Array Except Self | Prefix/Suffix | O(n) |

### Problem: Minimum Window Substring (Hard)

**Input:** s = "ADOBECODEBANC", t = "ABC"
**Output:** "BANC" (minimum window in s containing all chars of t)

```go
func minWindow(s, t string) string {
    if len(s) < len(t) { return "" }

    // need: how many of each character we still need in the window
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        need[t[i]]++
    }

    have := 0              // number of characters satisfied (count == need)
    required := len(need)  // number of distinct chars we need
    left := 0
    bestLeft, bestLen := 0, len(s)+1
    window := make(map[byte]int)

    for right := 0; right < len(s); right++ {
        ch := s[right]
        window[ch]++

        // If this character's count in window matches the needed count,
        // we have satisfied one more distinct character requirement.
        if need[ch] > 0 && window[ch] == need[ch] {
            have++
        }

        // When all required characters are satisfied, try to shrink window
        for have == required {
            // Update best answer if this window is smaller
            if right-left+1 < bestLen {
                bestLeft = left
                bestLen = right - left + 1
            }

            // Shrink from the left
            leftCh := s[left]
            window[leftCh]--
            if need[leftCh] > 0 && window[leftCh] < need[leftCh] {
                have-- // we no longer have enough of this character
            }
            left++
        }
    }

    if bestLen == len(s)+1 { return "" }
    return s[bestLeft : bestLeft+bestLen]
}
// Time: O(s + t), Space: O(s + t)
```

---

## Summary

- **Sliding Window:** Contiguous subarray/substring problems. Expand right, shrink left when invalid. O(n).
- **Two Pointers:** Sorted arrays, pair/triplet problems, comparing ends. Moves two indices intelligently. O(n) or O(n²) for triple.
- **Prefix Sum:** Range sum queries, counting subarrays with a property. O(1) query after O(n) build.
- **Hashing:** O(1) lookup, frequency counting, complement searching. Trade space for time.
- In Go: use `strings.Builder` for efficient string building. Strings are immutable — watch out for O(n²) concatenation.
- When stuck on an array problem: try these four patterns in order. One of them almost always applies.

---

## Exercises

### Easy
1. Given an array of integers, return true if any value appears at least twice. (Hashing)
2. Given a string, find the first non-repeating character. Return its index. (Frequency map)
3. Find the maximum average of a contiguous subarray of length k. (Sliding window)

### Medium
4. Given a binary array, find the maximum number of consecutive 1s if you can flip at most one 0. (Sliding window with constraint)
5. Given an integer array, return the number of subarrays whose product is less than k. (Two pointers or sliding window)
6. Find all starting indices of anagrams of pattern p in string s. (Sliding window + frequency map)

### Hard
7. Implement the "Minimum Window Substring" problem from scratch in Go without looking at the solution above. Then compare.
8. Given nums = [1,2,3,4,5] and k=2, find the length of the shortest subarray with sum at least k. Note: this is harder than it looks because values can be negative. Research the deque-based approach.
9. Given an array of n integers, find the maximum value of (nums[i] - nums[j]) * nums[k] where j is between i and k. What is the optimal complexity?
