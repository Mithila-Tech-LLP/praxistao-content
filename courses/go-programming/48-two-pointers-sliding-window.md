# Chapter 48: Two Pointers, Sliding Window, and Prefix Sum

Three patterns that each reduce an O(n²) brute-force solution to O(n) by maintaining information across iterations instead of restarting from scratch.

## Table of Contents

1. [Two Pointers](#1-two-pointers)
2. [Sliding Window](#2-sliding-window)
3. [Prefix Sum](#3-prefix-sum)
4. [Combining the Patterns](#4-combining-the-patterns)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. Two Pointers

Two pointers work on a sorted array or a sequence where you can make progress by moving one pointer based on the current comparison.

### Pattern 1: Opposite ends (sorted array)

```go
// Two Sum — find pair that sums to target (sorted input)
// Brute force: O(n²) — check every pair
// Two pointers: O(n) — converge from both ends
func twoSum(nums []int, target int) (int, int) {
    lo, hi := 0, len(nums)-1
    for lo < hi {
        sum := nums[lo] + nums[hi]
        if sum == target  { return lo, hi }
        if sum < target   { lo++ }  // need larger sum
        else              { hi-- }  // need smaller sum
    }
    return -1, -1
}
```

**Why it works**: the array is sorted. If the sum is too small, the only way to increase it is to move `lo` right. If too large, move `hi` left. We never need to backtrack.

```go
// Three Sum — find all triplets summing to 0
// O(n²) with two pointers (fix one element, two-pointer the rest)
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int
    
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue } // skip duplicates
        
        lo, hi := i+1, len(nums)-1
        for lo < hi {
            sum := nums[i] + nums[lo] + nums[hi]
            if sum == 0 {
                result = append(result, []int{nums[i], nums[lo], nums[hi]})
                for lo < hi && nums[lo] == nums[lo+1] { lo++ } // skip duplicates
                for lo < hi && nums[hi] == nums[hi-1] { hi-- }
                lo++; hi--
            } else if sum < 0 { lo++ } else { hi-- }
        }
    }
    return result
}
```

### Pattern 2: Same direction (slow/fast pointers)

```go
// Remove duplicates from sorted array in-place
// O(n) time, O(1) space
func removeDuplicates(nums []int) int {
    if len(nums) == 0 { return 0 }
    slow := 0
    for fast := 1; fast < len(nums); fast++ {
        if nums[fast] != nums[slow] {
            slow++
            nums[slow] = nums[fast]
        }
    }
    return slow + 1 // new length
}

// Move zeros to end (maintain relative order of non-zeros)
func moveZeroes(nums []int) {
    slow := 0
    for fast := 0; fast < len(nums); fast++ {
        if nums[fast] != 0 {
            nums[slow], nums[fast] = nums[fast], nums[slow]
            slow++
        }
    }
}
```

### Pattern 3: Floyd's cycle detection

```go
// Detect cycle in a linked list
// slow moves 1 step, fast moves 2 steps
// If there's a cycle, they meet inside it
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast { return true }
    }
    return false
}

// Find the start of the cycle
func detectCycle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            // Move one pointer to head; both move at same speed
            slow = head
            for slow != fast {
                slow = slow.Next
                fast = fast.Next
            }
            return slow
        }
    }
    return nil
}
```

### Pattern 4: Container with most water

```go
// Find two lines forming a container with the most water
// O(n²) brute force → O(n) two pointers
func maxArea(height []int) int {
    lo, hi := 0, len(height)-1
    max := 0
    for lo < hi {
        h := min(height[lo], height[hi])
        area := h * (hi - lo)
        if area > max { max = area }
        // Always move the shorter line — moving the taller can only decrease
        if height[lo] < height[hi] { lo++ } else { hi-- }
    }
    return max
}
```

---

## 2. Sliding Window

A sliding window maintains a contiguous subarray (window) and slides it through the array. The key is updating the window state incrementally instead of recomputing from scratch.

### Fixed-size window

```go
// Maximum sum of k consecutive elements
// O(n) vs O(n×k) brute force
func maxSumK(nums []int, k int) int {
    if len(nums) < k { return 0 }
    
    // Initialize first window
    windowSum := 0
    for i := 0; i < k; i++ { windowSum += nums[i] }
    
    maxSum := windowSum
    for i := k; i < len(nums); i++ {
        windowSum += nums[i]       // add new element
        windowSum -= nums[i-k]     // remove old element
        if windowSum > maxSum { maxSum = windowSum }
    }
    return maxSum
}
```

### Variable-size window (expand/shrink)

The variable window pattern:
1. Expand the window by moving `right` forward
2. When the window becomes invalid, shrink by moving `left` forward

```go
// Longest substring without repeating characters
func lengthOfLongestSubstring(s string) int {
    seen := make(map[byte]int)  // char → last seen index
    best := 0
    left := 0
    
    for right := 0; right < len(s); right++ {
        c := s[right]
        if idx, ok := seen[c]; ok && idx >= left {
            left = idx + 1  // shrink window past the duplicate
        }
        seen[c] = right
        if right-left+1 > best { best = right - left + 1 }
    }
    return best
}

// Minimum window substring — find smallest window in s containing all chars of t
func minWindow(s, t string) string {
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ { need[t[i]]++ }
    
    window := make(map[byte]int)
    formed := 0  // how many distinct chars from t are satisfied
    required := len(need)
    
    best := ""
    left := 0
    
    for right := 0; right < len(s); right++ {
        c := s[right]
        window[c]++
        if need[c] > 0 && window[c] == need[c] {
            formed++
        }
        
        // Shrink window while all chars satisfied
        for formed == required {
            if best == "" || right-left+1 < len(best) {
                best = s[left : right+1]
            }
            lc := s[left]
            window[lc]--
            if need[lc] > 0 && window[lc] < need[lc] {
                formed--
            }
            left++
        }
    }
    return best
}

// Longest subarray with sum ≤ k (non-negative numbers)
func longestSubarrayWithSumAtMostK(nums []int, k int) int {
    left, sum, best := 0, 0, 0
    for right := 0; right < len(nums); right++ {
        sum += nums[right]
        for sum > k {
            sum -= nums[left]
            left++
        }
        if right-left+1 > best { best = right - left + 1 }
    }
    return best
}
```

### Sliding window with deque (monotonic deque)

```go
// Maximum in each window of size k
// O(n) using a monotonic deque (indices stored in decreasing value order)
func maxSlidingWindow(nums []int, k int) []int {
    deque := []int{}  // stores indices, values decreasing
    result := []int{}
    
    for i, v := range nums {
        // Remove elements outside the window
        for len(deque) > 0 && deque[0] <= i-k {
            deque = deque[1:]
        }
        // Remove smaller elements — they can never be the max
        for len(deque) > 0 && nums[deque[len(deque)-1]] < v {
            deque = deque[:len(deque)-1]
        }
        deque = append(deque, i)
        
        if i >= k-1 {
            result = append(result, nums[deque[0]])
        }
    }
    return result
}
```

---

## 3. Prefix Sum

A prefix sum array allows answering "what is the sum of elements from index i to j?" in O(1) after O(n) preprocessing.

```go
// Build prefix sum array
// prefix[i] = sum of nums[0..i-1]
// sum(i, j) = prefix[j+1] - prefix[i]
func buildPrefix(nums []int) []int {
    prefix := make([]int, len(nums)+1)
    for i, v := range nums { prefix[i+1] = prefix[i] + v }
    return prefix
}

func rangeSum(prefix []int, i, j int) int {
    return prefix[j+1] - prefix[i]
}

// Example:
nums := []int{1, 2, 3, 4, 5}
prefix := buildPrefix(nums)  // [0, 1, 3, 6, 10, 15]
fmt.Println(rangeSum(prefix, 1, 3))  // 2+3+4 = 9
```

### Subarray sum equals k

```go
// Count subarrays with sum = k
// Key insight: sum[i..j] = prefix[j+1] - prefix[i]
// We want prefix[j+1] - prefix[i] = k
// So: prefix[i] = prefix[j+1] - k
// Track count of each prefix sum seen so far
func subarraySum(nums []int, k int) int {
    count := 0
    prefixCount := map[int]int{0: 1}  // empty prefix has sum 0
    sum := 0
    
    for _, v := range nums {
        sum += v
        count += prefixCount[sum-k]  // how many previous prefixes satisfy?
        prefixCount[sum]++
    }
    return count
}
```

### 2D prefix sum

```go
// Sum of submatrix from (r1,c1) to (r2,c2) in O(1)
type Matrix struct {
    prefix [][]int
}

func NewMatrix(grid [][]int) *Matrix {
    rows, cols := len(grid), len(grid[0])
    prefix := make([][]int, rows+1)
    for i := range prefix { prefix[i] = make([]int, cols+1) }
    
    for r := 1; r <= rows; r++ {
        for c := 1; c <= cols; c++ {
            prefix[r][c] = grid[r-1][c-1] +
                prefix[r-1][c] + prefix[r][c-1] - prefix[r-1][c-1]
        }
    }
    return &Matrix{prefix}
}

func (m *Matrix) SumRegion(r1, c1, r2, c2 int) int {
    return m.prefix[r2+1][c2+1] -
        m.prefix[r1][c2+1] -
        m.prefix[r2+1][c1] +
        m.prefix[r1][c1]
}
```

### Difference array (inverse of prefix sum)

A difference array enables O(1) range updates instead of O(n).

```go
// Apply +val to range [lo, hi] many times, then query final values
func applyRangeUpdates(n int, updates [][3]int) []int {
    diff := make([]int, n+1)
    for _, u := range updates {
        lo, hi, val := u[0], u[1], u[2]
        diff[lo] += val      // start adding val here
        diff[hi+1] -= val    // stop adding val here
    }
    // Reconstruct final array using prefix sum of diff
    result := make([]int, n)
    sum := 0
    for i := range result {
        sum += diff[i]
        result[i] = sum
    }
    return result
}
```

---

## 4. Combining the Patterns

```go
// Longest subarray sum at most k with negative numbers
// Prefix sum + sliding window doesn't work with negatives
// Use prefix sum + binary search (sorted prefix sums with monotone deque)
// → O(n log n)

// Subarray with product less than k
// Sliding window + product tracking
func numSubarrayProductLessThanK(nums []int, k int) int {
    if k <= 1 { return 0 }
    count, prod, left := 0, 1, 0
    for right := 0; right < len(nums); right++ {
        prod *= nums[right]
        for prod >= k {
            prod /= nums[left]
            left++
        }
        count += right - left + 1  // all subarrays ending at right, starting [left..right]
    }
    return count
}

// Two sum with sorted output — two pointers
// Four sum — sort + two two-pointer passes
func fourSum(nums []int, target int) [][]int {
    sort.Ints(nums)
    var result [][]int
    n := len(nums)
    
    for i := 0; i < n-3; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        for j := i + 1; j < n-2; j++ {
            if j > i+1 && nums[j] == nums[j-1] { continue }
            lo, hi := j+1, n-1
            for lo < hi {
                sum := nums[i] + nums[j] + nums[lo] + nums[hi]
                if sum == target {
                    result = append(result, []int{nums[i], nums[j], nums[lo], nums[hi]})
                    for lo < hi && nums[lo] == nums[lo+1] { lo++ }
                    for lo < hi && nums[hi] == nums[hi-1] { hi-- }
                    lo++; hi--
                } else if sum < target { lo++ } else { hi-- }
            }
        }
    }
    return result
}
```

---

## Summary

| Pattern | When to use | Complexity |
|---------|-------------|------------|
| Two pointers (opposite) | Sorted array, pair/triplet sum | O(n) |
| Two pointers (same dir) | In-place transformation, cycle detection | O(n) |
| Fixed window | Fixed-size subarray aggregation | O(n) |
| Variable window | Longest/shortest subarray with constraint | O(n) |
| Monotonic deque | Window min/max | O(n) |
| Prefix sum | Range sum queries | O(n) build, O(1) query |
| Difference array | Range update queries | O(n) build, O(1) update |

**Key insight**: these patterns eliminate redundant computation by maintaining a running state. The window or pointer state is the compressed summary of what you've seen so far.

---

## Exercises

### Easy
1. Implement `pairWithTargetSum(nums []int, target int) (int, int)` using two pointers on a sorted array. Then implement the same thing using a hash map for unsorted input. Compare time and space.
2. Given a binary string `s`, find the maximum consecutive 1s after flipping at most one 0 to 1. Solve with a sliding window.
3. Compute the sum of all odd-indexed elements minus the sum of all even-indexed elements using prefix sums. Answer range queries `[i, j]` in O(1).

### Medium
4. **Longest repeating character replacement**: given a string, you can replace at most k characters. Find the length of the longest substring containing the same letter after replacement. Use a sliding window: the window is valid if `(length - max_freq_in_window) ≤ k`.
5. **Trapping rain water**: given an elevation map `height`, compute how much water it can trap. Solve in O(n) time and O(1) space using two pointers, tracking the current left-max and right-max.
6. **Subarray sum divisible by k**: count subarrays whose sum is divisible by k. Use prefix sums and the property that `(prefix[j] - prefix[i]) % k == 0` iff `prefix[j] % k == prefix[i] % k`.

### Hard
7. **Minimum number of operations to make array equal**: you have an array where `arr[i] = 2*i + 1`. In one operation, you can increment one element by 1. Find the minimum operations to make all elements equal. Use prefix sums to compute the answer for each possible target value in O(n).
8. **Maximum sum rectangle in a 2D matrix**: find the submatrix with the maximum sum. Use prefix sums to reduce each column range to a 1D array, then apply Kadane's algorithm. Time complexity: O(n²m) where m is the number of columns.
