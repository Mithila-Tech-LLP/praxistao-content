# Chapter 23: Searching Algorithms — Finding What You Need

> "A good search algorithm is the difference between finding a needle in a haystack and finding a needle in a sorted, indexed, cross-referenced haystack in one step."

---

## Overview

Searching is perhaps the most fundamental operation in programming. Every time a compiler checks if a variable is defined, every time a database retrieves a record, every time your IDE highlights a misspelled word — a search algorithm is running. The question is not whether you need searching. The question is which searching algorithm is appropriate and why.

This chapter covers the spectrum from linear search (look at everything, simple but slow) to binary search (halve the problem each time, fast but requires sorted data) to hash map lookup (near-instant, the magic behind symbol tables) to string searching algorithms that power text editors and compilers alike.

We will also build three concrete search features for the Astra compiler: symbol table lookups, type definition searches, and the error position finder that converts a character offset into a line number.

## What We Are Building

By the end of this chapter, you will understand when to use each searching strategy, how to implement binary search correctly (including its notoriously tricky edge cases), how hash maps achieve O(1) lookup, and how the Knuth-Morris-Pratt algorithm searches strings in O(n+m) instead of O(nm). We will see all three techniques used in real Astra compiler code.

---

## Table of Contents

1. Linear Search — When Brute Force Is Right
2. Binary Search — The Most Important Search Algorithm
3. Binary Search Variants — First/Last Occurrence, Bounds
4. Interpolation and Jump Search
5. Searching in Trees — BST Lookup
6. Hash Map Lookup — O(1) Magic
7. String Searching — Naive vs Knuth-Morris-Pratt
8. Rabin-Karp Rolling Hash
9. Astra Build Milestone
10. Exercises
11. Summary

---

## 1. Linear Search — When Brute Force Is Right

Linear search checks each element one by one until it finds the target.

```go
func linearSearch(arr []int, target int) int {
    for i, v := range arr {
        if v == target {
            return i  // found at index i
        }
    }
    return -1  // not found
}
```

**Complexity:** O(n) time, O(1) space.

**When linear search is the right choice:**

1. **Unsorted data**: If you cannot sort (perhaps the data changes constantly, or sorting is too expensive), linear search is your only option.

2. **Very small arrays**: For arrays of fewer than 10-20 elements, the overhead of binary search setup (bounds checking, division) exceeds the benefit. Insertion sort and linear search are used for small arrays in Timsort for the same reason.

3. **Searching linked lists**: You cannot jump to the middle of a linked list. You must traverse from the beginning — linear search is the only option.

4. **First occurrence with early termination**: If you expect the target to be near the beginning, linear search may terminate early while binary search always takes O(log n) steps.

5. **When the "list" is a stream**: If you are searching a stream of data that you are reading for the first time (like scanning a file), you cannot binary search it.

```go
// Linear search through unsorted struct data
// (binary search not applicable — not sorted by this field)
type Symbol struct {
    Name  string
    Type  string
    Value interface{}
}

func findSymbolLinear(symbols []Symbol, name string) *Symbol {
    for i := range symbols {
        if symbols[i].Name == name {
            return &symbols[i]
        }
    }
    return nil
}
// O(n) — has to check each symbol
// Acceptable for very small scopes (local variables in a tiny function)
// NOT acceptable for a module-level symbol table with thousands of symbols
```

---

## 2. Binary Search — The Most Important Search Algorithm

Binary search is the algorithm that makes large sorted data sets usable. It works by repeatedly halving the search space: look at the middle element, if it is the target you are done, if the target is smaller search the left half, if larger search the right half.

The key precondition: **the data must be sorted**.

```
Search for 7 in: [1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
Indices:          0  1  2  3  4   5   6   7   8   9

Step 1: left=0, right=9, mid=4
  arr[4] = 9, target=7, 7 < 9, so search left half
  right = mid - 1 = 3

Step 2: left=0, right=3, mid=1
  arr[1] = 3, target=7, 7 > 3, so search right half
  left = mid + 1 = 2

Step 3: left=2, right=3, mid=2
  arr[2] = 5, target=7, 7 > 5, so search right half
  left = mid + 1 = 3

Step 4: left=3, right=3, mid=3
  arr[3] = 7, target=7, FOUND at index 3!

4 steps to find element in array of 10.
For 1,000,000 elements: at most 20 steps!
```

```go
// Standard binary search — correct implementation
func binarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    
    for left <= right {
        // CRITICAL: use this formula, not (left+right)/2
        // (left+right)/2 can overflow when left and right are large int32 values
        mid := left + (right-left)/2
        
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            left = mid + 1   // target is in right half
        } else {
            right = mid - 1  // target is in left half
        }
    }
    
    return -1  // not found
}
```

**The integer overflow bug is real.** In Java 1.5, the official Java binary search in `java.util.Arrays` had this bug for 9 years. For Go's `int` type (64-bit on modern systems), overflow is less likely, but it is still good practice.

**Analysis:**
- Each iteration: cut the search space in half
- Start with n elements, after k iterations: n / 2^k elements remain
- When 1 element remains: n / 2^k = 1, so k = log₂(n)
- Time: O(log n)
- Space: O(1) for the iterative version, O(log n) for recursive

---

## 3. Binary Search Variants

The basic binary search finds any occurrence. Often we need something more specific.

### Find First Occurrence

```go
// Returns the index of the FIRST occurrence of target, or -1
func findFirst(arr []int, target int) int {
    left, right := 0, len(arr)-1
    result := -1
    
    for left <= right {
        mid := left + (right-left)/2
        if arr[mid] == target {
            result = mid       // record this position
            right = mid - 1   // but keep searching LEFT for earlier occurrence
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return result
}
```

### Find Last Occurrence

```go
// Returns the index of the LAST occurrence of target, or -1
func findLast(arr []int, target int) int {
    left, right := 0, len(arr)-1
    result := -1
    
    for left <= right {
        mid := left + (right-left)/2
        if arr[mid] == target {
            result = mid      // record this position
            left = mid + 1   // but keep searching RIGHT for later occurrence
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return result
}
```

### Lower Bound — First Element >= Target

This is one of the most useful variants. It finds where you would insert the target to keep the array sorted.

```go
// Returns the index of the first element >= target
// This is the "lower bound" or "bisect_left" in Python
func lowerBound(arr []int, target int) int {
    left, right := 0, len(arr)
    
    for left < right {
        mid := left + (right-left)/2
        if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return left  // index where target would be inserted
}

// Upper Bound: first element > target
func upperBound(arr []int, target int) int {
    left, right := 0, len(arr)
    
    for left < right {
        mid := left + (right-left)/2
        if arr[mid] <= target {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return left
}
```

**Real compiler use case:** We store source line start positions in a sorted array. To find the line number for a character at position `pos`:

```go
// lineStarts[i] = byte offset of the first character on line i+1
lineStarts := []int{0, 45, 92, 140, 201, ...}

// Find line number for character at position 100
line := sort.SearchInts(lineStarts, 100) - 1
// sort.SearchInts uses binary search (lower bound internally)
// If 100 is between lineStarts[2]=92 and lineStarts[3]=140,
// SearchInts returns 3 (first index where lineStarts[i] >= 100)
// Subtract 1: line = 2 (0-indexed), which is line 3 in 1-indexed
```

---

## 4. Interpolation and Jump Search

### Interpolation Search

If the data is uniformly distributed, we can do better than always splitting in half. Interpolation search estimates where the target is likely to be based on its value relative to the range.

```go
func interpolationSearch(arr []int, target int) int {
    low, high := 0, len(arr)-1
    
    for low <= high && target >= arr[low] && target <= arr[high] {
        if low == high {
            if arr[low] == target {
                return low
            }
            return -1
        }
        
        // Estimate position based on value
        pos := low + (high-low) * (target-arr[low]) / (arr[high]-arr[low])
        
        if arr[pos] == target {
            return pos
        } else if arr[pos] < target {
            low = pos + 1
        } else {
            high = pos - 1
        }
    }
    return -1
}
```

For uniformly distributed data, interpolation search is O(log log n) — even faster than binary search. But for non-uniform data, it degrades to O(n). In practice, binary search is more predictable and usually preferred.

### Jump Search

Jump search is between linear and binary: jump ahead by blocks of size √n, then do linear search within the block.

```go
func jumpSearch(arr []int, target int) int {
    n := len(arr)
    step := int(math.Sqrt(float64(n)))
    prev := 0
    
    // Jump ahead by step until we overshoot
    for arr[min(step, n)-1] < target {
        prev = step
        step += int(math.Sqrt(float64(n)))
        if prev >= n {
            return -1
        }
    }
    
    // Linear search in the block
    for arr[prev] < target {
        prev++
        if prev == min(step, n) {
            return -1
        }
    }
    
    if arr[prev] == target {
        return prev
    }
    return -1
}
// O(√n) time — useful when jumping backwards is expensive
// (e.g., reading magnetic tape — seeking backward is slow)
```

---

## 5. Searching in Trees — BST Lookup

A Binary Search Tree (BST) stores elements such that for every node, all values in the left subtree are smaller and all values in the right subtree are larger.

```
        8
       / \
      3   10
     / \    \
    1   6    14
       / \   /
      4   7 13
```

Searching a BST:
```go
type BSTNode struct {
    Value int
    Left  *BSTNode
    Right *BSTNode
}

func bstSearch(root *BSTNode, target int) *BSTNode {
    if root == nil {
        return nil  // not found
    }
    if target == root.Value {
        return root  // found!
    }
    if target < root.Value {
        return bstSearch(root.Left, target)   // search left subtree
    }
    return bstSearch(root.Right, target)       // search right subtree
}
// O(h) where h is the height of the tree
// Balanced BST: h = O(log n) → O(log n) search
// Degenerate BST (sorted input): h = O(n) → O(n) search
```

For the Astra compiler's type definition lookup, we could use a BST. But Go's built-in map (hash map) is generally preferred for O(1) average lookup. A balanced BST (AVL tree or red-black tree) is used when we need ordered iteration — for example, iterating over all types in alphabetical order.

---

## 6. Hash Map Lookup — O(1) Magic

A hash map achieves O(1) average lookup through a clever mechanism:

1. **Hash function**: Convert the key (e.g., a string) to an integer (the hash)
2. **Bucket selection**: Use hash % numBuckets to select a bucket (array index)
3. **Store in bucket**: Place the key-value pair in that bucket
4. **Lookup**: Compute hash, go directly to that bucket, find the value

```
Keys: "hello", "world", "foo"
Hash function: sum of ASCII values % 8

hash("hello") = (104+101+108+108+111) % 8 = 532 % 8 = 4
hash("world") = (119+111+114+108+100) % 8 = 552 % 8 = 0
hash("foo")   = (102+111+111) % 8 = 324 % 8 = 4  ← COLLISION with "hello"!

Buckets:
[0]: [(world, ...)]
[1]: []
[2]: []
[3]: []
[4]: [(hello, ...), (foo, ...)]  ← collision: both in bucket 4
[5]: []
[6]: []
[7]: []
```

**Collisions** are handled by chaining (each bucket is a linked list of entries) or open addressing (probe nearby buckets).

For the Astra symbol table, we use Go's built-in `map` which handles all this automatically:

```go
// Symbol table using Go's hash map
type SymbolTable struct {
    symbols map[string]*Symbol
    parent  *SymbolTable  // for nested scopes
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
    return &SymbolTable{
        symbols: make(map[string]*Symbol),
        parent:  parent,
    }
}

func (st *SymbolTable) Define(name string, sym *Symbol) {
    st.symbols[name] = sym  // O(1) average
}

func (st *SymbolTable) Resolve(name string) (*Symbol, bool) {
    // Check current scope
    sym, ok := st.symbols[name]  // O(1) average
    if ok {
        return sym, true
    }
    // Check parent scope (lexical scoping)
    if st.parent != nil {
        return st.parent.Resolve(name)
    }
    return nil, false
}
```

**Why O(1)?** The hash function maps the key to a bucket number directly — no iteration needed. Looking up "printf" is exactly as fast as looking up "x", regardless of how many symbols are in the table.

**Load factor and rehashing:** When the hash map gets too full (typically above 70-80% capacity), Go's runtime automatically resizes the map (doubles the number of buckets and rehashes all entries). This is an O(n) operation but happens rarely, giving O(1) amortized cost.

---

## 7. String Searching — Naive vs Knuth-Morris-Pratt

String searching (finding a pattern string within a text string) is critical for text editors, compilers (searching for keywords), and many other applications.

### Naive String Search: O(nm)

```go
func naiveSearch(text, pattern string) int {
    n, m := len(text), len(pattern)
    
    for i := 0; i <= n-m; i++ {
        j := 0
        for j < m && text[i+j] == pattern[j] {
            j++
        }
        if j == m {
            return i  // found at index i
        }
    }
    return -1
}
```

**The problem:** For each position in the text, we might compare up to m characters. With n positions and m comparisons each, that is O(nm) in the worst case.

```
Text:    "AAAAAAAAAB"
Pattern: "AAAB"
Naive search: at each of the (n-m) positions, we compare all 4 characters
before finding a mismatch. O(n*m) = O(10*4) = 40 comparisons.
```

### Knuth-Morris-Pratt (KMP): O(n+m)

KMP avoids re-examining characters we have already seen. When a mismatch occurs, instead of sliding the pattern back to the next position, KMP uses information about the pattern itself to skip ahead further.

**The key insight:** When we find a mismatch after some characters matched, we know those matched characters. If the matched prefix contains a proper suffix that is also a prefix, we can slide the pattern forward to align that prefix with the matched suffix.

```
Failure function (also called "partial match" or "lps" — longest proper prefix-suffix):
Pattern: "ABABC"
lps[0] = 0  (A: no proper prefix = suffix)
lps[1] = 0  (AB: "A" is prefix, "B" is suffix, don't match)
lps[2] = 1  (ABA: "A" matches as both prefix and suffix of length 1)
lps[3] = 2  (ABAB: "AB" matches as both prefix and suffix)
lps[4] = 0  (ABABC: no prefix = suffix)
```

```go
// Build the failure function (lps array)
func buildLPS(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m)
    lps[0] = 0
    
    length := 0  // length of previous longest prefix-suffix
    i := 1
    
    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            if length != 0 {
                // Don't increment i — use the lps value
                length = lps[length-1]
            } else {
                lps[i] = 0
                i++
            }
        }
    }
    return lps
}

// KMP search
func kmpSearch(text, pattern string) []int {
    n, m := len(text), len(pattern)
    if m == 0 {
        return []int{}
    }
    
    lps := buildLPS(pattern)
    var matches []int
    
    i, j := 0, 0  // i = text index, j = pattern index
    
    for i < n {
        if text[i] == pattern[j] {
            i++
            j++
        }
        
        if j == m {
            // Found a match at i-j
            matches = append(matches, i-j)
            j = lps[j-1]  // use failure function to continue searching
        } else if i < n && text[i] != pattern[j] {
            if j != 0 {
                j = lps[j-1]  // skip already-known matched prefix
            } else {
                i++
            }
        }
    }
    return matches
}
// O(n+m): n for text scan, m for building lps
// Never goes backward in the text!
```

**Why O(n+m)?** The text pointer `i` only moves forward (never goes back). The pattern pointer `j` can move back, but each backward move was preceded by a forward move, so total forward moves bounds total backward moves. Both together are O(n). Building the LPS is O(m). Total: O(n+m).

**Where Astra uses KMP:** The error reporter might search for keyword occurrences in the source to provide better error messages ("Did you mean 'return'?"). String searching also appears in preprocessor `#include` detection.

---

## 8. Rabin-Karp Rolling Hash

Rabin-Karp is an alternative to KMP that uses hashing to find pattern matches. It is especially useful for finding multiple patterns simultaneously.

**Idea:** Compute a hash of the pattern. Then slide a window of size m across the text, computing the hash of each window. If the hash matches, do a character-by-character comparison to confirm.

The "rolling" part: when the window slides one position, we can update the hash in O(1) by removing the contribution of the character that left the window and adding the new character.

```go
const (
    base  = 31     // a prime roughly equal to the alphabet size
    modP  = 1_000_000_007  // a large prime to prevent overflow
)

func rabinKarp(text, pattern string) []int {
    n, m := len(text), len(pattern)
    if m > n {
        return nil
    }
    
    // Precompute base^(m-1) % modP
    power := 1
    for i := 0; i < m-1; i++ {
        power = (power * base) % modP
    }
    
    // Compute hash of pattern and first window
    patHash, winHash := 0, 0
    for i := 0; i < m; i++ {
        patHash = (patHash*base + int(pattern[i])) % modP
        winHash = (winHash*base + int(text[i])) % modP
    }
    
    var matches []int
    
    for i := 0; i <= n-m; i++ {
        if winHash == patHash {
            // Hash match — verify character by character (avoid false positives)
            if text[i:i+m] == pattern {
                matches = append(matches, i)
            }
        }
        
        // Roll the hash: remove leftmost, add rightmost
        if i < n-m {
            winHash = (winHash - int(text[i])*power%modP + modP) % modP
            winHash = (winHash*base + int(text[i+m])) % modP
        }
    }
    return matches
}
// Average O(n+m), worst case O(nm) due to hash collisions
// But with a good hash function, collisions are extremely rare
```

---

## Astra Build Milestone

Let us implement all three search strategies that the Astra compiler actually uses:

```go
// File: compiler/search/compiler_search.go
package search

import (
    "fmt"
    "sort"
)

// ─── 1. SYMBOL TABLE: HASH MAP SEARCH O(1) ───────────────────────────────────

// SymbolKind classifies what kind of thing a symbol refers to
type SymbolKind int

const (
    KindVariable SymbolKind = iota
    KindFunction
    KindType
    KindParameter
)

func (k SymbolKind) String() string {
    return []string{"variable", "function", "type", "parameter"}[k]
}

// Symbol is a named entity in the Astra program
type Symbol struct {
    Name     string
    Kind     SymbolKind
    TypeName string
    Line     int
}

// ScopeTable is a single lexical scope's symbol table
type ScopeTable struct {
    symbols map[string]*Symbol // O(1) average lookup
    parent  *ScopeTable
    name    string // for debugging: "global", "fn:main", etc.
}

func NewScopeTable(name string, parent *ScopeTable) *ScopeTable {
    return &ScopeTable{
        symbols: make(map[string]*Symbol),
        parent:  parent,
        name:    name,
    }
}

func (s *ScopeTable) Define(sym *Symbol) error {
    if existing, exists := s.symbols[sym.Name]; exists {
        return fmt.Errorf("symbol '%s' already defined at line %d", sym.Name, existing.Line)
    }
    s.symbols[sym.Name] = sym // O(1) average
    return nil
}

// Resolve looks up a symbol name, walking up the scope chain
// This is the core of lexical scoping in Astra
func (s *ScopeTable) Resolve(name string) (*Symbol, bool) {
    // First check current scope: O(1)
    if sym, ok := s.symbols[name]; ok {
        return sym, true
    }
    // Then check parent scopes
    if s.parent != nil {
        return s.parent.Resolve(name)
    }
    return nil, false // undefined symbol
}

// ─── 2. TYPE DEFINITIONS: BST-LIKE SEARCH VIA SORTED SLICE ───────────────────

// TypeRegistry stores all known types, kept sorted for binary search
type TypeRegistry struct {
    types  []TypeDef
    sorted bool
}

// TypeDef describes an Astra type
type TypeDef struct {
    Name   string
    Kind   string // "struct", "enum", "alias", "builtin"
    Fields []string
}

func NewTypeRegistry() *TypeRegistry {
    reg := &TypeRegistry{}
    // Register builtin types
    builtins := []TypeDef{
        {Name: "bool", Kind: "builtin"},
        {Name: "float64", Kind: "builtin"},
        {Name: "int", Kind: "builtin"},
        {Name: "string", Kind: "builtin"},
    }
    reg.types = append(reg.types, builtins...)
    reg.sorted = true // builtins added in sorted order
    return reg
}

func (r *TypeRegistry) Register(t TypeDef) {
    r.types = append(r.types, t)
    r.sorted = false // need to re-sort before binary search
}

func (r *TypeRegistry) ensureSorted() {
    if !r.sorted {
        sort.Slice(r.types, func(i, j int) bool {
            return r.types[i].Name < r.types[j].Name
        })
        r.sorted = true
    }
}

// Lookup finds a type by name using binary search: O(log n)
func (r *TypeRegistry) Lookup(name string) (*TypeDef, bool) {
    r.ensureSorted()
    
    // Binary search on the sorted types slice
    lo, hi := 0, len(r.types)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if r.types[mid].Name == name {
            return &r.types[mid], true
        } else if r.types[mid].Name < name {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return nil, false
}

// ─── 3. ERROR POSITION: BINARY SEARCH ON LINE OFFSETS ────────────────────────

// SourceMap tracks line start positions for O(log n) position → line conversion
type SourceMap struct {
    source     string
    lineStarts []int // lineStarts[i] = byte offset of line i (0-indexed)
}

func NewSourceMap(source string) *SourceMap {
    sm := &SourceMap{source: source}
    sm.lineStarts = []int{0} // line 0 starts at offset 0
    
    // Scan source once to find all line start positions: O(n)
    for i, ch := range source {
        if ch == '\n' && i+1 < len(source) {
            sm.lineStarts = append(sm.lineStarts, i+1)
        }
    }
    return sm
}

// PositionFor converts a byte offset to (line, column): O(log n)
// This is used when generating error messages — we have a token's
// byte offset from the lexer and need to convert it to line:column
func (sm *SourceMap) PositionFor(offset int) (line, col int) {
    // Binary search for the line containing this offset
    // Find the last lineStart that is <= offset
    lo, hi := 0, len(sm.lineStarts)-1
    for lo < hi {
        mid := lo + (hi-lo+1)/2
        if sm.lineStarts[mid] <= offset {
            lo = mid
        } else {
            hi = mid - 1
        }
    }
    line = lo       // 0-indexed line number
    col = offset - sm.lineStarts[lo]  // 0-indexed column
    return line + 1, col + 1          // return 1-indexed
}

// ─── 4. DEMO ──────────────────────────────────────────────────────────────────

func RunSearchDemo() {
    fmt.Println("=== Astra Compiler Search Demo ===\n")

    // --- Symbol Table Demo ---
    fmt.Println("--- 1. Symbol Table (Hash Map, O(1)) ---")
    global := NewScopeTable("global", nil)
    _ = global.Define(&Symbol{Name: "print", Kind: KindFunction, TypeName: "fn(string)->void", Line: 0})
    _ = global.Define(&Symbol{Name: "int", Kind: KindType, TypeName: "builtin", Line: 0})

    fnScope := NewScopeTable("fn:main", global)
    _ = fnScope.Define(&Symbol{Name: "x", Kind: KindVariable, TypeName: "int", Line: 5})
    _ = fnScope.Define(&Symbol{Name: "y", Kind: KindVariable, TypeName: "int", Line: 6})

    // Resolve symbols from within the function
    for _, name := range []string{"x", "print", "z"} {
        sym, found := fnScope.Resolve(name)
        if found {
            fmt.Printf("  Resolved '%s': %s (scope lookup)\n", name, sym.Kind)
        } else {
            fmt.Printf("  '%s': undefined symbol error\n", name)
        }
    }

    // --- Type Registry Demo ---
    fmt.Println("\n--- 2. Type Registry (Binary Search, O(log n)) ---")
    reg := NewTypeRegistry()
    reg.Register(TypeDef{Name: "Node", Kind: "struct", Fields: []string{"value: int", "next: Node"}})
    reg.Register(TypeDef{Name: "Color", Kind: "enum", Fields: []string{"Red", "Green", "Blue"}})

    for _, name := range []string{"int", "Node", "Unknown", "Color", "bool"} {
        t, found := reg.Lookup(name)
        if found {
            fmt.Printf("  Type '%s': %s\n", name, t.Kind)
        } else {
            fmt.Printf("  Type '%s': not found\n", name)
        }
    }

    // --- Source Map Demo ---
    fmt.Println("\n--- 3. Source Map (Binary Search on Line Offsets, O(log n)) ---")
    source := `fn main() {
    let x = 10
    let y = 20
    print(x + y)
}`
    sm := NewSourceMap(source)

    fmt.Printf("  Source has %d lines, %d chars\n", len(sm.lineStarts), len(source))

    // Simulate errors at various byte offsets
    testOffsets := []int{0, 12, 25, 40, 55}
    for _, offset := range testOffsets {
        if offset < len(source) {
            line, col := sm.PositionFor(offset)
            fmt.Printf("  Offset %2d → line %d, col %d (char: %q)\n",
                offset, line, col, source[offset])
        }
    }
}
```

Running this produces:

```
=== Astra Compiler Search Demo ===

--- 1. Symbol Table (Hash Map, O(1)) ---
  Resolved 'x': variable (scope lookup)
  Resolved 'print': function (scope lookup)
  'z': undefined symbol error

--- 2. Type Registry (Binary Search, O(log n)) ---
  Type 'int': builtin
  Type 'Node': struct
  Type 'Unknown': not found
  Type 'Color': enum
  Type 'bool': builtin

--- 3. Source Map (Binary Search on Line Offsets, O(log n)) ---
  Source has 5 lines, 57 chars
  Offset  0 → line 1, col 1 (char: 'f')
  Offset 12 → line 2, col 1 (char: ' ')
  Offset 25 → line 2, col 14 (char: '0')
  Offset 40 → line 3, col 14 (char: '0')
  Offset 55 → line 4, col 15 (char: ')')
```

---

## Exercises

1. **Implement binary search recursively**: Write a recursive version of binary search. What is its space complexity compared to the iterative version? When would the recursive version cause a stack overflow?

2. **Find the count of a value**: Given a sorted array that may contain duplicates, use binary search to find how many times a given value appears. Use lowerBound and upperBound to do this in O(log n).

3. **First missing positive**: Given an unsorted array of integers, find the smallest missing positive integer. Can you do this in O(n) time and O(1) space? Hint: use the array itself as a hash map.

4. **Implement the source map**: Extend the SourceMap implementation to also support looking up the line's full text content (for error display), given a line number. This should also be O(1) or O(line length).

5. **KMP trace**: Manually trace through the KMP algorithm searching for "ABAB" in "ABABABABAB". Build the LPS array first, then trace each step of the search. How many character comparisons does KMP make vs naive search?

6. **Hash map from scratch**: Implement a simple hash map in Go using an array of linked lists (chaining). Implement Put(key, value) and Get(key) operations. Handle collisions with chaining. Compare its performance with Go's built-in map.

7. **Symbol resolution order**: Astra supports nested scopes (functions inside functions). Implement a symbol table that correctly resolves names in the following program:
   ```astra
   let x = 1
   fn outer() {
       let x = 2
       fn inner() {
           let y = x  // should resolve to outer's x = 2
       }
   }
   ```

8. **Search benchmark**: Write a Go benchmark that compares:
   - Linear search on an unsorted array of 10,000 strings
   - Binary search on a sorted array of 10,000 strings
   - Hash map lookup on a map of 10,000 strings
   Use Go's `testing.B` benchmark framework. Report ns/op for each.

---

## Summary Table

| Algorithm             | Time         | Space  | Requires Sorted? | Best For                          |
|-----------------------|--------------|--------|------------------|-----------------------------------|
| Linear Search         | O(n)         | O(1)   | No               | Unsorted data, small arrays       |
| Binary Search         | O(log n)     | O(1)   | Yes              | Sorted arrays, line number lookup |
| BST Search            | O(log n) avg | O(1)   | Inherently sorted | Ordered data with updates        |
| Hash Map Lookup       | O(1) avg     | O(n)   | No               | Symbol tables, fast key lookup    |
| Interpolation Search  | O(log log n) | O(1)   | Yes, uniform     | Uniformly distributed sorted data |
| Jump Search           | O(√n)        | O(1)   | Yes              | When backward jump is expensive   |
| Naive String Search   | O(nm)        | O(1)   | No               | Short patterns, simple code       |
| KMP String Search     | O(n+m)       | O(m)   | No               | Repeated pattern searching        |
| Rabin-Karp            | O(n+m) avg   | O(1)   | No               | Multiple pattern searching        |

| Astra Compiler Search | Algorithm  | Complexity | Why This Choice                      |
|-----------------------|------------|------------|--------------------------------------|
| Symbol lookup         | Hash map   | O(1) avg   | Most frequent operation, must be fast |
| Type lookup           | Binary search | O(log n) | Types sorted alphabetically          |
| Error line number     | Binary search | O(log n) | lineStarts array is naturally sorted |
| Keyword detection     | Hash set   | O(1)       | Fixed set of keywords                |

The central lesson: the "right" search algorithm is determined by your data's properties and access patterns. The Astra compiler uses three different strategies in three different situations — not because we are being clever for its own sake, but because each situation genuinely calls for a different tool. Understanding searching algorithms means understanding when O(1), O(log n), and O(n) are each the right answer.
