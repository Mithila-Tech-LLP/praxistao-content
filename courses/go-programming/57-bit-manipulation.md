# Chapter 57: Bit Manipulation

Bit manipulation operates directly on the binary representation of integers. Many problems that seem complex become trivial with the right bit operation. Go gives you full access to bitwise operators on all integer types.

## Table of Contents

1. [Bitwise Operators](#1-bitwise-operators)
2. [Common Bit Tricks](#2-common-bit-tricks)
3. [Bit Counting and math/bits](#3-bit-counting-and-mathbits)
4. [Classic Problems](#4-classic-problems)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. Bitwise Operators

```go
a, b := 0b1010, 0b1100  // 10, 12 in decimal

fmt.Printf("%04b\n", a & b)   // AND:  1000 = 8  (both bits set)
fmt.Printf("%04b\n", a | b)   // OR:   1110 = 14 (either bit set)
fmt.Printf("%04b\n", a ^ b)   // XOR:  0110 = 6  (exactly one bit set)
fmt.Printf("%04b\n", ^a)      // NOT:  ...10101 (all bits flipped, signed)
fmt.Printf("%04b\n", a &^ b)  // AND NOT: 0010 = 2 (bits in a but not b)
fmt.Printf("%04b\n", a << 1)  // LEFT SHIFT:  10100 = 20 (× 2)
fmt.Printf("%04b\n", a >> 1)  // RIGHT SHIFT: 0101  = 5  (÷ 2)
```

### Signed vs unsigned shifts

```go
var x int8 = -4           // 0b11111100 in two's complement
fmt.Println(x >> 1)       // -2 (arithmetic shift: sign bit propagates)

var y uint8 = 252         // 0b11111100
fmt.Println(y >> 1)       // 126 (logical shift: zero fills in)
```

Always use unsigned types when you want logical right shift.

---

## 2. Common Bit Tricks

### Check, set, clear, toggle a bit

```go
// Check bit i of n
func hasBit(n, i int) bool { return n & (1 << i) != 0 }

// Set bit i of n (to 1)
func setBit(n, i int) int { return n | (1 << i) }

// Clear bit i of n (to 0)
func clearBit(n, i int) int { return n &^ (1 << i) }

// Toggle bit i of n
func toggleBit(n, i int) int { return n ^ (1 << i) }

x := 0b0101
fmt.Printf("%04b\n", setBit(x, 1))    // 0111 = 7
fmt.Printf("%04b\n", clearBit(x, 2))  // 0001 = 1
fmt.Printf("%04b\n", toggleBit(x, 3)) // 1101 = 13
```

### Check if power of 2

```go
// A power of 2 has exactly one bit set.
// n & (n-1) clears the lowest set bit.
func isPowerOf2(n int) bool {
    return n > 0 && n&(n-1) == 0
}
// isPowerOf2(8) → true  (1000 & 0111 = 0000)
// isPowerOf2(6) → false (0110 & 0101 = 0100 ≠ 0)
```

### Isolate the lowest set bit

```go
// n & (-n) isolates the rightmost 1 bit
func lowestBit(n int) int { return n & (-n) }
// lowestBit(12) = 12 & (-12) = 0b1100 & 0b...0100 = 0b0100 = 4
```

### Clear the lowest set bit

```go
func clearLowestBit(n int) int { return n & (n - 1) }
// Used to count set bits:
func popCount(n int) int {
    count := 0
    for n != 0 { n &= n - 1; count++ }
    return count
}
```

### XOR properties

```go
// XOR cancels identical values:
x ^ x == 0       // any number XOR itself = 0
x ^ 0 == x       // any number XOR 0 = itself
x ^ y ^ x == y   // XOR is commutative and associative

// Find the single non-duplicate in a list where every other element appears twice:
func singleNumber(nums []int) int {
    result := 0
    for _, n := range nums { result ^= n }
    return result
    // All duplicates cancel out (x ^ x = 0), lone element remains
}
```

### Swap without temporary

```go
func swapBits(a, b *int) {
    *a ^= *b  // a = a XOR b
    *b ^= *a  // b = b XOR (a XOR b) = original a
    *a ^= *b  // a = (a XOR b) XOR original a = original b
}
// Note: breaks if a and b point to the same variable!
```

---

## 3. Bit Counting and math/bits

Go's `math/bits` package provides efficient hardware-accelerated operations:

```go
import "math/bits"

n := uint(0b10110100)

// Population count (number of set bits) — maps to POPCNT instruction
fmt.Println(bits.OnesCount(n))    // 4

// Leading zeros
fmt.Println(bits.LeadingZeros(n))  // 56 (for 64-bit, assuming 64-bit arch)

// Trailing zeros
fmt.Println(bits.TrailingZeros(n)) // 2

// Length = position of highest set bit + 1
fmt.Println(bits.Len(n))           // 8

// Reverse bits
fmt.Println(bits.Reverse(n))

// Rotate
fmt.Println(bits.RotateLeft(n, 2)) // rotate left by 2

// Next power of 2 (rounds up)
func nextPowerOf2(n uint) uint {
    if n == 0 { return 1 }
    return 1 << bits.Len(n-1)
}
```

---

## 4. Classic Problems

### Single number III (two unique elements)

```go
// All elements appear twice except two. Find them.
func singleNumberIII(nums []int) (int, int) {
    // XOR all elements: result = a ^ b (the two unique elements)
    xor := 0
    for _, n := range nums { xor ^= n }
    
    // Find a bit where a and b differ (any set bit in xor works)
    // Use the lowest set bit to partition elements into two groups
    diff := xor & (-xor) // isolate lowest set bit
    
    a, b := 0, 0
    for _, n := range nums {
        if n & diff != 0 { a ^= n } else { b ^= n }
    }
    return a, b
}
```

### Number of 1 bits in a range

```go
// Count total set bits in all numbers from 1 to n
func countBitsUpToN(n int) int {
    total := 0
    for i := 0; i < bits.Len(uint(n)); i++ {
        full := (n + 1) / (1 << (i + 1))
        remainder := (n+1)%(1<<(i+1)) - (1 << i)
        if remainder < 0 { remainder = 0 }
        total += full*(1<<i) + remainder
    }
    return total
}
```

### Subset generation with bitmask

```go
// Generate all subsets of nums using bitmask
func allSubsets(nums []int) [][]int {
    n := len(nums)
    total := 1 << n  // 2^n subsets
    result := make([][]int, 0, total)
    
    for mask := 0; mask < total; mask++ {
        subset := []int{}
        for i := 0; i < n; i++ {
            if mask & (1 << i) != 0 {
                subset = append(subset, nums[i])
            }
        }
        result = append(result, subset)
    }
    return result
}
```

### DP with bitmask (Travelling Salesman)

Bitmask DP uses a bitmask to represent which elements have been "visited":

```go
// Travelling Salesman Problem: minimum cost to visit all cities exactly once
// O(2^n × n²) DP
func tsp(dist [][]int) int {
    n := len(dist)
    INF := 1<<30
    
    // dp[mask][i] = min cost to reach city i, having visited exactly the cities in mask
    dp := make([][]int, 1<<n)
    for i := range dp {
        dp[i] = make([]int, n)
        for j := range dp[i] { dp[i][j] = INF }
    }
    dp[1][0] = 0  // start at city 0, mask = 0b0001 (only city 0 visited)
    
    for mask := 1; mask < 1<<n; mask++ {
        for last := 0; last < n; last++ {
            if dp[mask][last] == INF { continue }
            if mask & (1 << last) == 0 { continue } // last must be in mask
            
            for next := 0; next < n; next++ {
                if mask & (1 << next) != 0 { continue } // already visited
                newMask := mask | (1 << next)
                cost := dp[mask][last] + dist[last][next]
                if cost < dp[newMask][next] {
                    dp[newMask][next] = cost
                }
            }
        }
    }
    
    // Find min cost to return to city 0 after visiting all cities
    allVisited := (1 << n) - 1
    result := INF
    for last := 1; last < n; last++ {
        if dp[allVisited][last] != INF {
            cost := dp[allVisited][last] + dist[last][0]
            if cost < result { result = cost }
        }
    }
    return result
}
```

### XOR linear basis (competitive programming)

```go
// Maximum XOR of any subset using Gaussian elimination
type XORBasis struct {
    basis [64]int64
}

func (b *XORBasis) Insert(n int64) {
    for i := 62; i >= 0; i-- {
        if n>>(i) & 1 == 0 { continue }
        if b.basis[i] == 0 {
            b.basis[i] = n
            return
        }
        n ^= b.basis[i]
    }
}

func (b *XORBasis) MaxXOR() int64 {
    result := int64(0)
    for i := 62; i >= 0; i-- {
        if result ^ b.basis[i] > result {
            result ^= b.basis[i]
        }
    }
    return result
}
```

### Reverse bits

```go
// Reverse all 32 bits of a uint32
func reverseBits(n uint32) uint32 {
    return bits.Reverse32(n)
}

// Manual (to understand the algorithm):
func reverseBitsManual(n uint32) uint32 {
    result := uint32(0)
    for i := 0; i < 32; i++ {
        result = (result << 1) | (n & 1)
        n >>= 1
    }
    return result
}
```

---

## Summary

| Operation | Code | Use case |
|-----------|------|----------|
| Check bit i | `n & (1<<i) != 0` | Test flag |
| Set bit i | `n \| (1<<i)` | Enable flag |
| Clear bit i | `n &^ (1<<i)` | Disable flag |
| Toggle bit i | `n ^ (1<<i)` | Flip flag |
| Power of 2? | `n>0 && n&(n-1)==0` | Alignment check |
| Lowest set bit | `n & (-n)` | Fenwick tree update |
| Clear lowest bit | `n & (n-1)` | Iterate set bits |
| Count set bits | `bits.OnesCount(n)` | Hamming weight |
| Bitmask subset | `for mask:=0; mask<1<<n` | State enumeration |
| XOR cancels | `a^a==0` | Find single element |

---

## Exercises

### Easy
1. Implement `countSetBits(n int) int` three ways: loop with `n&1` + right shift, loop with `n&(n-1)`, and `bits.OnesCount`. Benchmark all three.
2. Given two integers, find how many bits differ between them (Hamming distance). Use XOR followed by count of set bits.
3. Implement `missingNumber(nums []int) int` for a list of n-1 numbers from 0..n. Use XOR: XOR all indices and all numbers together; pairs cancel, leaving the missing one.

### Medium
4. **Counting bits**: for every number from 0 to n, return an array where `result[i]` is the number of set bits in i. Solve in O(n) using the DP relation `dp[i] = dp[i >> 1] + (i & 1)`.
5. **Maximum XOR of two numbers**: given an array, find the pair with maximum XOR. Use a bitwise trie: insert each number bit by bit (MSB first), then for each number greedily choose the opposite bit to maximize XOR.
6. **Power set using bitmask**: given a set of strings (tags for a blog post), generate all possible combinations of tags. Use a bitmask and represent the set as a `uint64` bitmap. Implement `union`, `intersection`, `difference`, and `isSubset` as single bitwise operations.

### Hard
7. Implement a **Fenwick tree (Binary Indexed Tree)** that answers prefix sum queries in O(log n): `Update(i, delta)` and `Query(i)` (sum of [0..i]). The key operations are `i += i & (-i)` to move to the parent, and `i -= i & (-i)` to move to the sibling. Understand why `i & (-i)` isolates the lowest set bit.
8. Solve the **stickers spelling problem**: given sticker strings and a target string, find the minimum number of stickers to spell the target. Use bitmask DP where the state is which letters of the target have been covered. `O(2^n × m)` where n = target length.
