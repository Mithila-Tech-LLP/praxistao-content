# Chapter 40: Tries and Advanced Data Structures

A **trie** (prefix tree) stores strings character by character — each path from root to a leaf represents a word. Tries enable O(L) lookup where L is the word length, regardless of how many words are stored. They power autocomplete, spell-checkers, IP routing, and dictionary implementations. This chapter also covers the LRU cache (a classic interview favorite) and segment trees for range queries.

## Table of Contents

1. [Trie — Prefix Tree](#1-trie--prefix-tree)
2. [Trie Applications](#2-trie-applications)
3. [LRU Cache — Full Implementation](#3-lru-cache--full-implementation)
4. [Segment Tree — Range Queries](#4-segment-tree--range-queries)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. Trie — Prefix Tree

```
Words: ["apple", "app", "application", "apply", "banana"]

Trie structure:
  root
  ├── a
  │   └── p
  │       └── p  (end) ← "app"
  │           ├── l
  │           │   ├── e  (end) ← "apple"
  │           │   ├── i
  │           │   │   └── c
  │           │   │       └── a
  │           │   │           └── t
  │           │   │               └── i
  │           │   │                   └── o
  │           │   │                       └── n  (end) ← "application"
  │           │   └── y  (end) ← "apply"
  └── b
      └── a
          └── n
              └── a
                  └── n
                      └── a  (end) ← "banana"
```

**Implementation:**
```go
type TrieNode struct {
    children [26]*TrieNode  // For lowercase a-z
    isEnd    bool
    count    int  // How many words end here (for counting)
}

type Trie struct {
    root *TrieNode
}

func NewTrie() *Trie {
    return &Trie{root: &TrieNode{}}
}

// Insert adds a word to the trie — O(L) where L = word length.
func (t *Trie) Insert(word string) {
    node := t.root
    for _, ch := range word {
        idx := ch - 'a'
        if node.children[idx] == nil {
            node.children[idx] = &TrieNode{}
        }
        node = node.children[idx]
    }
    node.isEnd = true
    node.count++
}

// Search returns true if the exact word exists — O(L).
func (t *Trie) Search(word string) bool {
    node := t.traverse(word)
    return node != nil && node.isEnd
}

// StartsWith returns true if any word begins with the prefix — O(L).
func (t *Trie) StartsWith(prefix string) bool {
    return t.traverse(prefix) != nil
}

func (t *Trie) traverse(s string) *TrieNode {
    node := t.root
    for _, ch := range s {
        idx := ch - 'a'
        if node.children[idx] == nil {
            return nil
        }
        node = node.children[idx]
    }
    return node
}

// Delete removes a word from the trie — O(L).
func (t *Trie) Delete(word string) bool {
    return t.delete(t.root, word, 0)
}

func (t *Trie) delete(node *TrieNode, word string, depth int) bool {
    if node == nil {
        return false
    }
    if depth == len(word) {
        if !node.isEnd {
            return false  // Word doesn't exist
        }
        node.isEnd = false
        node.count--
        return t.isLeaf(node)  // True if node can be deleted
    }

    idx := word[depth] - 'a'
    if t.delete(node.children[idx], word, depth+1) {
        node.children[idx] = nil  // Safe to delete child
        return !node.isEnd && t.isLeaf(node)
    }
    return false
}

func (t *Trie) isLeaf(node *TrieNode) bool {
    for _, child := range node.children {
        if child != nil {
            return false
        }
    }
    return true
}
```

**Unicode-safe trie (for any character):**
```go
type UTrie struct {
    children map[rune]*UTrie
    isEnd    bool
}

func (t *UTrie) Insert(word string) {
    node := t
    for _, ch := range word {
        if node.children == nil {
            node.children = make(map[rune]*UTrie)
        }
        if node.children[ch] == nil {
            node.children[ch] = &UTrie{}
        }
        node = node.children[ch]
    }
    node.isEnd = true
}
```

### Quick Check
> 1. What is the time complexity of searching in a trie?
> 2. How is `Search("app")` different from `StartsWith("app")`?
> 3. What does each edge in a trie represent?

---

## 2. Trie Applications

### Autocomplete
```go
// Autocomplete returns all words in the trie with the given prefix.
func (t *Trie) Autocomplete(prefix string) []string {
    node := t.traverse(prefix)
    if node == nil {
        return nil
    }

    var results []string
    var collect func(*TrieNode, string)
    collect = func(n *TrieNode, current string) {
        if n.isEnd {
            results = append(results, current)
        }
        for i, child := range n.children {
            if child != nil {
                collect(child, current+string(rune('a'+i)))
            }
        }
    }

    collect(node, prefix)
    return results
}
```

### Longest Common Prefix
```go
// LongestCommonPrefix finds the LCP among all words in the trie.
func (t *Trie) LongestCommonPrefix() string {
    var sb strings.Builder
    node := t.root

    for {
        // Count non-nil children:
        childCount := 0
        var onlyChild *TrieNode
        var onlyChar rune

        for i, child := range node.children {
            if child != nil {
                childCount++
                onlyChild = child
                onlyChar = rune('a' + i)
            }
        }

        // LCP ends when: word ends here, or multiple children (branching):
        if node.isEnd || childCount != 1 {
            break
        }

        sb.WriteRune(onlyChar)
        node = onlyChild
    }

    return sb.String()
}
```

### Word Search in Grid (DFS + Trie)
```go
// WordSearch finds all words from a dictionary that exist in the 2D grid.
// Each word uses adjacent cells (4-directional), each cell used once per word.
func WordSearch(board [][]byte, words []string) []string {
    trie := NewTrie()
    for _, w := range words {
        trie.Insert(w)
    }

    rows, cols := len(board), len(board[0])
    var found []string
    dirs := [][2]int{{0,1},{0,-1},{1,0},{-1,0}}

    var dfs func(node *TrieNode, r, c int, path string)
    dfs = func(node *TrieNode, r, c int, path string) {
        ch := board[r][c]
        if ch == '#' { return }  // Already visited
        idx := ch - 'a'
        next := node.children[idx]
        if next == nil { return }  // No word with this prefix

        path += string(ch)
        if next.isEnd {
            found = append(found, path)
            next.isEnd = false  // Avoid duplicates
        }

        board[r][c] = '#'  // Mark as visited
        for _, d := range dirs {
            nr, nc := r+d[0], c+d[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
                dfs(next, nr, nc, path)
            }
        }
        board[r][c] = ch  // Restore
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            dfs(trie.root, r, c, "")
        }
    }
    return found
}
```

### Quick Check
> 1. What makes a trie more efficient than a hash map for prefix queries?
> 2. Why is trie-based autocomplete faster than scanning all words with `strings.HasPrefix`?

---

## 3. LRU Cache — Full Implementation

An LRU (Least Recently Used) cache evicts the least recently accessed item when full. The classic implementation combines a **hash map** (O(1) lookup) with a **doubly linked list** (O(1) order tracking):

```go
type lruNode struct {
    key, val   int
    prev, next *lruNode
}

type LRUCache struct {
    cap        int
    cache      map[int]*lruNode
    head, tail *lruNode  // Sentinel nodes — head=MRU side, tail=LRU side
}

func NewLRUCache(capacity int) *LRUCache {
    head := &lruNode{}
    tail := &lruNode{}
    head.next = tail
    tail.prev = head
    return &LRUCache{
        cap:   capacity,
        cache: make(map[int]*lruNode, capacity),
        head:  head,
        tail:  tail,
    }
}

// Get retrieves a value and marks it as most recently used — O(1).
func (c *LRUCache) Get(key int) int {
    if node, ok := c.cache[key]; ok {
        c.moveToFront(node)
        return node.val
    }
    return -1
}

// Put inserts or updates a key-value pair — O(1).
func (c *LRUCache) Put(key, val int) {
    if node, ok := c.cache[key]; ok {
        node.val = val
        c.moveToFront(node)
        return
    }

    node := &lruNode{key: key, val: val}
    c.cache[key] = node
    c.addToFront(node)

    if len(c.cache) > c.cap {
        // Evict LRU (node just before tail sentinel):
        lru := c.tail.prev
        c.removeNode(lru)
        delete(c.cache, lru.key)
    }
}

func (c *LRUCache) addToFront(node *lruNode) {
    node.prev = c.head
    node.next = c.head.next
    c.head.next.prev = node
    c.head.next = node
}

func (c *LRUCache) removeNode(node *lruNode) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (c *LRUCache) moveToFront(node *lruNode) {
    c.removeNode(node)
    c.addToFront(node)
}

// LRU structure visualization:
// head ↔ [most recently used] ↔ ... ↔ [least recently used] ↔ tail
// Get/Put moves the accessed node to after head (most recent position)
// Eviction removes the node before tail (least recent position)
```

**Testing the LRU cache:**
```go
func main() {
    cache := NewLRUCache(3)

    cache.Put(1, 10)  // [1]
    cache.Put(2, 20)  // [2,1]
    cache.Put(3, 30)  // [3,2,1]

    fmt.Println(cache.Get(1))  // 10 — [1,3,2]

    cache.Put(4, 40)  // [4,1,3] — evicts 2 (LRU)
    fmt.Println(cache.Get(2))  // -1 — evicted!
    fmt.Println(cache.Get(3))  // 30 — [3,4,1]
    fmt.Println(cache.Get(4))  // 40 — [4,3,1]
}
```

### Quick Check
> 1. What two data structures make an LRU cache O(1) for all operations?
> 2. Why use sentinel head/tail nodes?
> 3. Which end of the linked list represents "most recently used"?

---

## 4. Segment Tree — Range Queries

A **segment tree** answers range queries (sum, min, max) and range updates in O(log n):

```
Array: [1, 3, 5, 7, 9, 11]

Segment tree (range sums):
           [36]         range [0,5]
          /     \
       [9]       [27]   [0,2] and [3,5]
      /   \     /   \
    [4]   [5] [16]  [11] leaf pairs
    / \  / \  / \   / \
   [1][3][5][7][9][11]  leaves
```

```go
type SegTree struct {
    tree []int
    n    int
}

func NewSegTree(nums []int) *SegTree {
    n := len(nums)
    tree := make([]int, 4*n)  // 4*n is enough space
    st := &SegTree{tree: tree, n: n}
    st.build(nums, 0, 0, n-1)
    return st
}

// build constructs the segment tree — O(n).
func (st *SegTree) build(nums []int, node, start, end int) {
    if start == end {
        st.tree[node] = nums[start]
        return
    }
    mid := (start + end) / 2
    st.build(nums, 2*node+1, start, mid)
    st.build(nums, 2*node+2, mid+1, end)
    st.tree[node] = st.tree[2*node+1] + st.tree[2*node+2]  // Sum
}

// Query returns the sum of nums[l..r] — O(log n).
func (st *SegTree) Query(l, r int) int {
    return st.query(0, 0, st.n-1, l, r)
}

func (st *SegTree) query(node, start, end, l, r int) int {
    if r < start || end < l {
        return 0  // No overlap
    }
    if l <= start && end <= r {
        return st.tree[node]  // Complete overlap
    }
    mid := (start + end) / 2
    left := st.query(2*node+1, start, mid, l, r)
    right := st.query(2*node+2, mid+1, end, l, r)
    return left + right  // Partial overlap
}

// Update changes nums[idx] to val — O(log n).
func (st *SegTree) Update(idx, val int) {
    st.update(0, 0, st.n-1, idx, val)
}

func (st *SegTree) update(node, start, end, idx, val int) {
    if start == end {
        st.tree[node] = val
        return
    }
    mid := (start + end) / 2
    if idx <= mid {
        st.update(2*node+1, start, mid, idx, val)
    } else {
        st.update(2*node+2, mid+1, end, idx, val)
    }
    st.tree[node] = st.tree[2*node+1] + st.tree[2*node+2]
}

// Usage:
st := NewSegTree([]int{1, 3, 5, 7, 9, 11})
fmt.Println(st.Query(1, 4))  // 3+5+7+9 = 24
st.Update(2, 10)              // Change index 2 from 5 to 10
fmt.Println(st.Query(1, 4))  // 3+10+7+9 = 29
```

**Comparison with brute force:**
| | Brute Force | Segment Tree |
|--|------------|-------------|
| Build | O(1) | O(n) |
| Query | O(n) | O(log n) |
| Update | O(1) | O(log n) |
| **Use when** | Few queries | Many queries + updates |

### Quick Check
> 1. What is the time complexity of a range sum query using a segment tree?
> 2. Why allocate `4*n` space for the tree array?
> 3. When is a segment tree better than prefix sums?

---

## Summary

- **Trie**: prefix tree; O(L) insert/search/delete; perfect for prefix queries and autocomplete
- **Trie nodes**: array of 26 children (a-z) or `map[rune]` for Unicode
- **`isEnd` flag**: marks where words end — `Search` requires `isEnd`, `StartsWith` doesn't
- **LRU cache**: doubly linked list (order) + hash map (O(1) lookup); sentinel head/tail simplify pointer management; Get and Put both O(1)
- **Segment tree**: binary tree on array ranges; O(n) build, O(log n) query and update; use when you need both range queries AND point updates
- **When to use each**: trie for string prefixes; LRU for bounded caches; segment tree for range queries with updates (vs prefix sums for static arrays)

---

## Exercises

### Easy
1. Build a trie from a list of words. Implement `CountWordsWithPrefix(prefix string) int` — count how many words in the trie start with the prefix. Use the trie node's count field.
2. Implement `AutocompleteSorted(prefix string, limit int) []string` that returns at most `limit` autocomplete suggestions in alphabetical order.
3. Write `PrefixSumArray(nums []int)` that creates a prefix sum array, and `RangeSum(prefix []int, l, r int) int` that answers range sum queries in O(1). Compare with segment tree: when would you use each?

### Medium
4. Word break: Given a string `s` and dictionary of words, determine if `s` can be segmented into a sequence of dictionary words. Use dynamic programming with a trie for efficient prefix lookup. Test: `s="leetcode"`, dict=`["leet","code"]` → true; `s="catsandog"`, dict=`["cats","dog","sand","and","cat"]` → false.
5. LFU (Least Frequently Used) cache: Implement an LFU cache with O(1) Get and Put. Each key has a frequency. On eviction, remove the key with the lowest frequency. If tie: remove the least recently used among those. Use: a map of key→node, a map of freq→doubly-linked-list-of-keys, and track `minFreq`. Test with capacity=2.
6. Range minimum query with updates: Build a segment tree that supports: `Min(l, r int) int` — minimum value in range, `Update(idx, val int)` — update a value, `RangeAdd(l, r, delta int)` — add delta to all values in range (lazy propagation). Test: 10,000 random updates and queries, verify against brute force.

### Hard
7. Trie-based IP router: Implement a longest-prefix-match IP router using a binary trie (each bit of the IP address is an edge). `AddRoute(cidr string, nextHop string)` adds a CIDR route (e.g., "192.168.0.0/24"). `Lookup(ip string) string` returns the next hop for the most specific matching prefix. Test with a routing table of 1000 routes and 10,000 lookups — verify correctness against a brute-force linear scan.
8. Persistent segment tree: Implement a "fully persistent" segment tree where each update creates a new version sharing unchanged nodes with the previous version. `Update(version int, idx, val int) int` returns new version. `Query(version, l, r int) int` queries a specific version. This allows O(1) "time travel" to any past state. Use case: answer range sum queries on historical snapshots of an array. Test: 1000 versions, random queries across all versions.
