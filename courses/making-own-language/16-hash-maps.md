# Chapter 16: Hash Maps — The Most Useful Data Structure

> "When in doubt, use a hash map." — Every experienced programmer

---

## Overview

You have a list of 1,000,000 user records and you need to look up a user by their username. With a sorted array you could do a binary search in O(log n) — about 20 comparisons. With a hash map, you can do it in O(1) — one operation. For 1 million records, that's the difference between 20 steps and 1 step.

The hash map (also called a hash table, dictionary, or associative array) is almost certainly the most useful data structure in all of programming. Python's `dict`, JavaScript's objects, Go's `map`, Rust's `HashMap`, Java's `HashMap` — they are all hash maps under the hood, and they power an enormous fraction of the code written every day.

This chapter covers:
- The problem hash maps solve: O(1) lookup
- Hash functions: what they do and what makes one good
- Collision handling: separate chaining and open addressing
- Load factor and rehashing
- Go's built-in `map` type: complete usage guide
- String interning with a hash map
- Why hash maps are unordered
- Consistent hashing (brief introduction)
- **Astra Build Milestone**: The symbol table — the most critical data structure in any compiler

---

## What We're Building

By the end of this chapter you will understand how hash maps work at the implementation level, master Go's `map` type, and most importantly, you will have built the **symbol table** for the Astra compiler — the data structure that tracks every variable, function, and type across nested scopes and enables name resolution during compilation.

---

## Table of Contents

1. The Problem: Lookup Needs to Be Fast
2. The Hash Function
3. Building a Hash Table
4. Collision Handling: Separate Chaining
5. Collision Handling: Open Addressing (Linear Probing)
6. Load Factor and Rehashing
7. Average vs Worst Case O(1)
8. Go's map Type: Complete Reference
9. String Interning with a Hash Map
10. Why Hash Maps Are Unordered
11. Consistent Hashing (Brief Introduction)
12. Astra Build Milestone: The Symbol Table

---

## 1. The Problem: Lookup Needs to Be Fast

Suppose you are building a function that checks whether a word is a keyword in the Astra language:

```go
// Approach 1: linear search — O(n)
func isKeyword(word string) bool {
    keywords := []string{"fn", "let", "const", "if", "else", "for",
                         "while", "return", "struct", "impl", "true", "false"}
    for _, kw := range keywords {
        if kw == word {
            return true
        }
    }
    return false
}

// Approach 2: hash map lookup — O(1)
var keywordSet = map[string]bool{
    "fn": true, "let": true, "const": true, "if": true,
    "else": true, "for": true, "while": true, "return": true,
    "struct": true, "impl": true, "true": true, "false": true,
}

func isKeyword(word string) bool {
    return keywordSet[word]
}
```

For 12 keywords, the difference is tiny. But consider a symbol table with 50,000 identifiers in a large codebase. Linear search takes up to 50,000 comparisons. A hash map lookup takes 1 comparison on average. At scale, this is the difference between a compiler that takes 30 seconds and one that takes 0.001 seconds.

The fundamental operation that hash maps enable: **given a key, find its associated value in O(1) time**.

---

## 2. The Hash Function

A **hash function** takes a key of arbitrary size and maps it to an integer index in a fixed-size array:

```
hash("hello")  →  4
hash("world")  →  11
hash("fn")     →  2
hash("let")    →  7
```

The index tells us which "bucket" (array slot) to look in. If the array has 16 buckets, and `hash("hello")` returns 4, then the value associated with "hello" is stored at index 4.

```
Key        Hash  Index (hash % 16)
"hello"  → 5387 → 5387 % 16 = 3
"world"  → 9871 → 9871 % 16 = 15
"fn"     → 1234 → 1234 % 16 = 2
```

### Properties of a Good Hash Function

1. **Deterministic**: the same key always produces the same hash
2. **Uniform distribution**: keys spread evenly across buckets, minimizing collisions
3. **Fast to compute**: O(1) or O(key length) — not O(n) on the number of items
4. **Avalanche effect**: a small change in the key produces a completely different hash

A simple hash function for strings (FNV-1a variant):

```go
// fnv32 computes a 32-bit hash of a string.
func fnv32(key string) uint32 {
    const (
        offset uint32 = 2166136261
        prime  uint32 = 16777619
    )
    h := offset
    for i := 0; i < len(key); i++ {
        h ^= uint32(key[i])
        h *= prime
    }
    return h
}

// bucketIndex maps a key to an array index.
func bucketIndex(key string, numBuckets int) int {
    return int(fnv32(key)) % numBuckets
}
```

Go's runtime uses a much more sophisticated hash (based on AES hardware instructions when available) to prevent hash collision attacks. But the principle is the same.

---

## 3. Building a Hash Table

Let's build a simple hash map from scratch in Go:

```go
package hashtable

const defaultCapacity = 16
const loadFactorThreshold = 0.75

// Entry stores a key-value pair.
type Entry[K comparable, V any] struct {
    Key   K
    Value V
    hash  uint64
}

// HashMap[K, V] is a hash table with separate chaining.
type HashMap[K comparable, V any] struct {
    buckets  [][]Entry[K, V]  // each bucket is a slice (chain)
    size     int              // number of key-value pairs
    capacity int              // number of buckets
    hasher   func(K) uint64  // hash function
}

// New creates a new HashMap with a custom hasher.
func New[K comparable, V any](hasher func(K) uint64) *HashMap[K, V] {
    cap := defaultCapacity
    return &HashMap[K, V]{
        buckets:  make([][]Entry[K, V], cap),
        capacity: cap,
        hasher:   hasher,
    }
}

func (m *HashMap[K, V]) bucketIdx(key K) int {
    return int(m.hasher(key) % uint64(m.capacity))
}
```

---

## 4. Collision Handling: Separate Chaining

A **collision** occurs when two different keys hash to the same bucket index. This is unavoidable (by the pigeonhole principle: there are infinitely many possible keys but only a finite number of buckets).

**Separate chaining** resolves collisions by storing multiple entries in the same bucket as a linked list (or slice):

```
Buckets:  [0] → nil
          [1] → nil
          [2] → ("fn", FN_TOKEN)
          [3] → ("hello", 5) → ("hi", 9)   ← collision: both hash to 3
          [4] → nil
          ...
```

When two keys collide at bucket 3, we just append both to the chain. To look up "hi", we hash it to 3, then scan the chain until we find the entry whose key equals "hi".

```go
// Set inserts or updates a key-value pair.
func (m *HashMap[K, V]) Set(key K, value V) {
    if float64(m.size)/float64(m.capacity) >= loadFactorThreshold {
        m.rehash()
    }

    h := m.hasher(key)
    idx := int(h % uint64(m.capacity))
    bucket := m.buckets[idx]

    // Check if key already exists in the chain
    for i := range bucket {
        if bucket[i].Key == key {
            bucket[i].Value = value  // update existing
            return
        }
    }

    // Key not found: append new entry
    m.buckets[idx] = append(bucket, Entry[K, V]{Key: key, Value: value, hash: h})
    m.size++
}

// Get retrieves the value for a key.
// Returns (zero, false) if the key is not present.
func (m *HashMap[K, V]) Get(key K) (V, bool) {
    idx := m.bucketIdx(key)
    for _, entry := range m.buckets[idx] {
        if entry.Key == key {
            return entry.Value, true
        }
    }
    var zero V
    return zero, false
}

// Delete removes a key-value pair.
// Returns true if the key was found and deleted.
func (m *HashMap[K, V]) Delete(key K) bool {
    idx := m.bucketIdx(key)
    bucket := m.buckets[idx]
    for i, entry := range bucket {
        if entry.Key == key {
            // Remove element at index i
            m.buckets[idx] = append(bucket[:i], bucket[i+1:]...)
            m.size--
            return true
        }
    }
    return false
}

// Has returns true if the key exists in the map.
func (m *HashMap[K, V]) Has(key K) bool {
    _, ok := m.Get(key)
    return ok
}

// Size returns the number of key-value pairs.
func (m *HashMap[K, V]) Size() int { return m.size }

// ForEach calls fn for each key-value pair (order not guaranteed).
func (m *HashMap[K, V]) ForEach(fn func(K, V)) {
    for _, bucket := range m.buckets {
        for _, entry := range bucket {
            fn(entry.Key, entry.Value)
        }
    }
}
```

---

## 5. Collision Handling: Open Addressing (Linear Probing)

**Open addressing** stores all entries in the array itself — no separate chains. When a collision occurs, we "probe" for the next available slot.

**Linear probing**: check slot h, then h+1, then h+2, etc., wrapping around:

```
Insert "hello" → hash=3 → bucket[3] is empty → store here
Insert "world" → hash=3 → bucket[3] is taken → try bucket[4] → empty → store here
Insert "foo"   → hash=3 → bucket[3] taken → bucket[4] taken → bucket[5] empty → store here

Lookup "world" → hash=3 → bucket[3] has "hello" (not match) → check bucket[4] → "world" found
```

```go
type OAEntry[K comparable, V any] struct {
    key     K
    value   V
    occupied bool
    deleted  bool  // tombstone for deleted slots
}

type OpenAddressMap[K comparable, V any] struct {
    slots    []OAEntry[K, V]
    size     int
    capacity int
    hasher   func(K) uint64
}

func (m *OpenAddressMap[K, V]) probe(key K) int {
    idx := int(m.hasher(key) % uint64(m.capacity))
    for i := 0; i < m.capacity; i++ {
        slot := (idx + i) % m.capacity
        if !m.slots[slot].occupied || m.slots[slot].deleted {
            return slot   // empty or tombstone: can insert here
        }
        if m.slots[slot].key == key {
            return slot   // found existing key
        }
    }
    return -1  // table full
}
```

**Linear probing vs separate chaining:**

| Aspect           | Separate Chaining   | Linear Probing                    |
|------------------|---------------------|-----------------------------------|
| Memory           | Extra per entry     | Compact, all in one array         |
| Cache behavior   | Poor (pointer chasing) | Excellent (sequential memory)  |
| Load factor      | Can exceed 1.0      | Must stay below 1.0               |
| Deletion         | Simple              | Needs tombstones (tricky)         |
| Common use       | Go's map, Java HashMap | CPython dicts, many high-perf HTs |

Go's runtime uses open addressing with a variant of quadratic probing for its built-in `map` type, which gives excellent cache performance.

---

## 6. Load Factor and Rehashing

The **load factor** is the ratio of stored items to total buckets:

```
load factor = size / capacity
```

As load factor increases, collisions become more frequent and performance degrades. The conventional threshold is **0.75** (75% full).

When load factor exceeds the threshold, we **rehash** — create a new array with double the capacity and re-insert every item:

```go
func (m *HashMap[K, V]) rehash() {
    oldBuckets := m.buckets
    m.capacity *= 2
    m.buckets = make([][]Entry[K, V], m.capacity)
    m.size = 0

    // Re-insert all existing entries into the new, larger table
    for _, bucket := range oldBuckets {
        for _, entry := range bucket {
            m.Set(entry.Key, entry.Value)
        }
    }
}
```

Rehashing is expensive: O(n) time to move all n items. But because the capacity doubles each time, rehashing happens less and less frequently. The amortized cost of inserting n items is still O(n) total — O(1) per insertion.

```
n insertions → n/2 without rehash + n/4 + n/8 + ... = O(n) total rehash work
```

This is the same amortized analysis as Go's slice `append`.

---

## 7. Average vs Worst Case O(1)

Hash maps offer **average case** O(1), not **guaranteed** O(1).

**Worst case O(n)** occurs when every key hashes to the same bucket, creating one enormous chain. This can happen:

1. **By bad luck** with a poor hash function
2. **By deliberate attack** — an adversary can craft inputs that cause maximum collisions, creating a "hash flooding" denial-of-service attack

Go defends against hash flooding by randomizing the hash seed at program startup. The same key hashes to different values in different program runs. An attacker can't precompute a set of colliding keys because they don't know the seed.

```
Run 1:  hash("hello") = 3476281
Run 2:  hash("hello") = 9823741  (different seed, different result)
```

For almost all real workloads with a good hash function, hash maps are effectively O(1).

---

## 8. Go's map Type: Complete Reference

Go's built-in `map` type is a generic hash map. Here is everything you need to know:

```go
package main

import "fmt"

func main() {
    // Creating a map
    m := make(map[string]int)          // empty map
    m2 := map[string]int{              // map literal
        "alice": 30,
        "bob":   25,
    }
    _ = m2

    // Writing a value
    m["apple"] = 42
    m["banana"] = 17

    // Reading a value
    val := m["apple"]           // returns 42
    _ = val
    missing := m["cherry"]      // returns 0 (zero value), NOT a panic
    _ = missing

    // Checking if a key exists (the "comma ok" idiom)
    val, ok := m["apple"]
    if ok {
        fmt.Printf("apple: %d\n", val)  // apple: 42
    }

    val, ok = m["cherry"]
    if !ok {
        fmt.Println("cherry not found")  // cherry not found
    }

    // Deleting a key
    delete(m, "banana")

    // Ranging over a map (ORDER IS NOT GUARANTEED)
    for key, value := range m {
        fmt.Printf("%s → %d\n", key, value)
    }

    // Getting map size
    fmt.Println("size:", len(m))

    // Nil map: reads return zero values, writes PANIC
    var nilMap map[string]int
    _ = nilMap["key"]       // safe: returns 0
    // nilMap["key"] = 1    // PANIC: assignment to nil map
    // Always use make() or a literal before writing.

    // Maps with struct values
    type Point struct{ X, Y int }
    points := map[string]Point{
        "origin": {0, 0},
        "right":  {1, 0},
    }
    fmt.Println(points["origin"])  // {0 0}

    // CANNOT take address of map value directly
    // &points["origin"]  // compile error

    // Map of slices: useful for grouping
    groups := make(map[string][]string)
    groups["fruits"] = append(groups["fruits"], "apple", "banana")
    groups["veggies"] = append(groups["veggies"], "carrot")

    // Checking if map is nil vs empty
    var m3 map[string]int
    fmt.Println(m3 == nil)        // true
    m4 := make(map[string]int)
    fmt.Println(m4 == nil)        // false
    fmt.Println(len(m4) == 0)     // true (empty but not nil)
}
```

### Common Patterns

```go
// Frequency counter (word count)
func wordCount(text string) map[string]int {
    counts := make(map[string]int)
    for _, word := range strings.Fields(text) {
        counts[word]++  // if key missing, zero value (0) is used, then incremented
    }
    return counts
}

// Set (using map as a set)
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(item T)        { s[item] = struct{}{} }
func (s Set[T]) Has(item T) bool   { _, ok := s[item]; return ok }
func (s Set[T]) Remove(item T)     { delete(s, item) }
func (s Set[T]) Len() int          { return len(s) }

// Memoization (caching function results)
func memoFib(n int, cache map[int]int) int {
    if n <= 1 { return n }
    if val, ok := cache[n]; ok { return val }
    result := memoFib(n-1, cache) + memoFib(n-2, cache)
    cache[n] = result
    return result
}
```

---

## 9. String Interning with a Hash Map

**String interning** is an optimization where you store only one copy of each unique string in a table, and everywhere else you use a pointer to that one copy. This means:
- String equality checks become pointer comparisons: O(1) instead of O(n)
- Less memory for repeated strings

```go
// interning/interner.go

package interning

import "sync"

// Interner stores a single copy of each unique string.
type Interner struct {
    mu    sync.Mutex
    table map[string]string
}

func NewInterner() *Interner {
    return &Interner{table: make(map[string]string)}
}

// Intern returns the canonical copy of s.
// If s has been seen before, returns the cached copy.
// If s is new, stores it and returns it.
func (in *Interner) Intern(s string) string {
    in.mu.Lock()
    defer in.mu.Unlock()
    if canonical, ok := in.table[s]; ok {
        return canonical  // return the pre-existing copy
    }
    in.table[s] = s       // store new string
    return s
}

// Size returns the number of unique strings interned.
func (in *Interner) Size() int {
    in.mu.Lock()
    defer in.mu.Unlock()
    return len(in.table)
}
```

In the Astra compiler, every identifier string goes through the interner. This means:
- "hello" appearing 10,000 times in source code = 10,000 pointers to ONE string
- Comparing two identifiers for equality: one pointer comparison, not character-by-character

```go
// In the lexer, when we produce an IDENTIFIER token:
func (l *Lexer) makeIdentifier(start, end int) Token {
    lexeme := l.source[start:end]
    interned := l.interner.Intern(lexeme)  // ensure uniqueness
    return Token{Type: IDENTIFIER, Lexeme: interned, ...}
}
```

---

## 10. Why Hash Maps Are Unordered

When you range over a Go map, the order of keys is **intentionally randomized** on each run:

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range m {
    fmt.Println(k, v)
}
// Might print: b 2, a 1, c 3
// Or: c 3, b 2, a 1
// Or: a 1, c 3, b 2
// Order changes every run.
```

This is not a bug — it is deliberate. The reasons are:
1. **Hash flood defense**: randomized seed means randomized key order
2. **Implementation freedom**: Go's map can reorganize buckets internally for performance
3. **Forces correctness**: if you rely on map ordering, your code has a latent bug; randomization makes that bug surface immediately

When you need ordered iteration, sort the keys explicitly:

```go
import "sort"

keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Printf("%s: %d\n", k, m[k])
}
```

---

## 11. Consistent Hashing (Brief Introduction)

In distributed systems, you might have a hash map spread across multiple servers (a distributed cache like Redis, or a distributed database). Regular hashing breaks when you add or remove servers — you'd have to rehash almost every key.

**Consistent hashing** solves this by arranging servers on a virtual "ring":

```
             Server A
                |
    Server D ───┤─── Server B
                |
             Server C

Keys hash to points on the ring.
Each key is stored on the next server clockwise.
Adding a server: only the keys between it and its predecessor move.
Removing a server: only that server's keys move to its successor.
```

Without consistent hashing: add 1 server to a 10-server cluster → ~90% of keys move.
With consistent hashing: add 1 server to a 10-server cluster → ~10% of keys move.

Consistent hashing is used by Amazon's DynamoDB, Cassandra's distributed hash ring, and CDNs (content delivery networks).

---

## 12. Astra Build Milestone: The Symbol Table

The **symbol table** is the most important data structure in the Astra compiler. It answers the question: "When the programmer writes `x`, what does `x` refer to?"

Consider this Astra code:

```astra
fn add(a: int, b: int) -> int {
    let result = a + b
    return result
}

fn main() {
    let x = 10
    let y = add(x, 20)
    if y > 0 {
        let msg = "positive"
        print(msg)
    }
    // msg is NOT accessible here — different scope
}
```

The symbol table tracks:
- `add` is a function taking (int, int) and returning int
- `a` and `b` are parameters of `add` with type `int`
- `result` is a local variable of `add` with type `int`
- `x` is a local variable of `main` with type `int`
- `y` is a local variable of `main` with type `int`
- `msg` is a local variable of the `if` block — not accessible outside

Notice that scopes are **nested** — each `{...}` block creates a new scope. The symbol table must handle this nesting.

```go
// sema/symbol_table.go

package sema

import (
    "fmt"
    "your-module/types"
)

// SymbolKind describes what kind of entity a symbol represents.
type SymbolKind int

const (
    KindVariable  SymbolKind = iota // let x = ...
    KindFunction                    // fn foo() { ... }
    KindParameter                   // fn foo(a: int) — 'a' is a parameter
    KindStruct                      // struct Point { ... }
    KindConstant                    // const MAX = 100
)

func (k SymbolKind) String() string {
    switch k {
    case KindVariable:  return "variable"
    case KindFunction:  return "function"
    case KindParameter: return "parameter"
    case KindStruct:    return "struct"
    case KindConstant:  return "constant"
    default:            return "unknown"
    }
}

// Symbol represents a named entity in the program.
type Symbol struct {
    Name        string
    Kind        SymbolKind
    Type        types.Type  // what type does this symbol have?
    Line        int         // where was it declared?
    Column      int
    IsUsed      bool        // has it been referenced? (for unused-variable warnings)
    IsMutable   bool        // let (mutable) vs const (immutable)
}

func (s *Symbol) String() string {
    return fmt.Sprintf("Symbol{%s %s: %s at %d:%d}",
        s.Kind, s.Name, s.Type, s.Line, s.Column)
}

// Scope represents one block of visibility (one {...} region).
// Scopes are linked in a chain: child → parent → grandparent → ... → global
type Scope struct {
    symbols map[string]*Symbol
    parent  *Scope
    depth   int     // 0 = global, 1 = function, 2 = if/for body, etc.
    name    string  // for debug: "global", "fn:main", "if@line10", etc.
}

// NewScope creates a new child scope.
func NewScope(parent *Scope, name string) *Scope {
    depth := 0
    if parent != nil {
        depth = parent.depth + 1
    }
    return &Scope{
        symbols: make(map[string]*Symbol),
        parent:  parent,
        depth:   depth,
        name:    name,
    }
}

// Define adds a new symbol to this scope.
// Returns an error if the symbol is already defined in THIS scope
// (shadowing a parent scope is allowed).
func (s *Scope) Define(sym *Symbol) error {
    if existing, exists := s.symbols[sym.Name]; exists {
        return fmt.Errorf("symbol %q already defined in this scope (previous definition at %d:%d)",
            sym.Name, existing.Line, existing.Column)
    }
    s.symbols[sym.Name] = sym
    return nil
}

// Resolve looks up a symbol by name.
// Searches this scope first, then walks up the parent chain.
func (s *Scope) Resolve(name string) (*Symbol, bool) {
    // Check this scope first
    if sym, ok := s.symbols[name]; ok {
        sym.IsUsed = true
        return sym, true
    }
    // Walk up to parent scope
    if s.parent != nil {
        return s.parent.Resolve(name)
    }
    // Not found anywhere in the scope chain
    return nil, false
}

// ResolveLocal looks up a symbol only in the current scope (not parents).
// Used for redefinition checks.
func (s *Scope) ResolveLocal(name string) (*Symbol, bool) {
    sym, ok := s.symbols[name]
    return sym, ok
}

// UnusedVariables returns all defined variables in this scope that were never used.
func (s *Scope) UnusedVariables() []*Symbol {
    var unused []*Symbol
    for _, sym := range s.symbols {
        if (sym.Kind == KindVariable || sym.Kind == KindParameter) && !sym.IsUsed {
            unused = append(unused, sym)
        }
    }
    return unused
}

// All returns all symbols defined in this scope.
func (s *Scope) All() map[string]*Symbol {
    return s.symbols
}

// SymbolTable manages the entire set of scopes for a program.
// It is the entry point used by the semantic analyzer.
type SymbolTable struct {
    global  *Scope
    current *Scope
    scopes  []*Scope  // stack of scope history (for debugging/IDE tools)
}

// NewSymbolTable creates a symbol table with a global scope pre-populated
// with built-in functions and types.
func NewSymbolTable() *SymbolTable {
    global := NewScope(nil, "global")
    st := &SymbolTable{global: global, current: global}

    // Populate built-in functions
    builtins := []struct {
        name string
        typ  types.Type
    }{
        {"print", types.BuiltinFn},
        {"println", types.BuiltinFn},
        {"len", types.BuiltinFn},
        {"append", types.BuiltinFn},
    }
    for _, b := range builtins {
        _ = global.Define(&Symbol{
            Name: b.name,
            Kind: KindFunction,
            Type: b.typ,
        })
    }

    return st
}

// EnterScope creates and enters a new child scope.
func (st *SymbolTable) EnterScope(name string) {
    newScope := NewScope(st.current, name)
    st.scopes = append(st.scopes, st.current)
    st.current = newScope
}

// ExitScope returns to the parent scope.
// Returns the scope we just exited (so we can check for unused variables).
func (st *SymbolTable) ExitScope() *Scope {
    exited := st.current
    if len(st.scopes) > 0 {
        st.current = st.scopes[len(st.scopes)-1]
        st.scopes = st.scopes[:len(st.scopes)-1]
    }
    return exited
}

// Define adds a symbol to the current scope.
func (st *SymbolTable) Define(sym *Symbol) error {
    return st.current.Define(sym)
}

// Resolve looks up a symbol starting from the current scope.
func (st *SymbolTable) Resolve(name string) (*Symbol, bool) {
    return st.current.Resolve(name)
}

// CurrentDepth returns the current nesting depth.
func (st *SymbolTable) CurrentDepth() int {
    return st.current.depth
}
```

Now let's see the symbol table in action during semantic analysis of a complete Astra function:

```go
// sema/analyzer.go

package sema

import (
    "your-module/ast"
    "your-module/diagnostics"
    "your-module/types"
)

// Analyzer performs semantic analysis on the AST.
type Analyzer struct {
    table *SymbolTable
    diags *diagnostics.Engine
}

func NewAnalyzer(diags *diagnostics.Engine) *Analyzer {
    return &Analyzer{
        table: NewSymbolTable(),
        diags: diags,
    }
}

// Analyze runs semantic analysis on the entire program.
func (a *Analyzer) Analyze(program *ast.Program) {
    // First pass: register all top-level function and struct names
    // (so functions can call each other regardless of declaration order)
    for _, decl := range program.Declarations {
        if fn, ok := decl.(*ast.FunctionDecl); ok {
            a.registerFunction(fn)
        }
    }

    // Second pass: analyze bodies
    for _, decl := range program.Declarations {
        a.analyzeDecl(decl)
    }
}

func (a *Analyzer) registerFunction(fn *ast.FunctionDecl) {
    sym := &Symbol{
        Name: fn.Name,
        Kind: KindFunction,
        Type: types.FunctionType(fn),
        Line: fn.Line,
    }
    if err := a.table.Define(sym); err != nil {
        a.diags.Error(fn.Line, fn.Column,
            "redefinition of function %q", fn.Name)
    }
}

func (a *Analyzer) analyzeFunction(fn *ast.FunctionDecl) {
    // Enter a new scope for the function body
    a.table.EnterScope("fn:" + fn.Name)

    // Define all parameters in this scope
    for _, param := range fn.Params {
        sym := &Symbol{
            Name:   param.Name,
            Kind:   KindParameter,
            Type:   types.Resolve(param.TypeName),
            Line:   param.Line,
            Column: param.Column,
        }
        if err := a.table.Define(sym); err != nil {
            a.diags.Error(param.Line, param.Column,
                "duplicate parameter name %q", param.Name)
        }
    }

    // Analyze the function body
    a.analyzeBlock(fn.Body)

    // Check for unused variables (warn)
    exitedScope := a.table.ExitScope()
    for _, unused := range exitedScope.UnusedVariables() {
        a.diags.Warning(unused.Line, unused.Column,
            "variable %q declared but not used", unused.Name)
    }
}

func (a *Analyzer) analyzeLetStatement(stmt *ast.LetStmt) {
    // Analyze the value expression first
    valueType := a.analyzeExpr(stmt.Value)

    // Define the variable in the current scope
    sym := &Symbol{
        Name:      stmt.Name,
        Kind:      KindVariable,
        Type:      valueType,
        Line:      stmt.Line,
        Column:    stmt.Column,
        IsMutable: true,
    }
    if err := a.table.Define(sym); err != nil {
        a.diags.Error(stmt.Line, stmt.Column,
            "variable %q already defined in this scope", stmt.Name)
    }
}

func (a *Analyzer) analyzeIdentifier(expr *ast.Identifier) types.Type {
    sym, ok := a.table.Resolve(expr.Name)
    if !ok {
        a.diags.ErrorWithHint(
            expr.Line, expr.Column,
            "did you forget to declare it with 'let'?",
            "undefined variable %q", expr.Name)
        return types.Unknown
    }
    return sym.Type
}

func (a *Analyzer) analyzeIfStmt(stmt *ast.IfStmt) {
    // Analyze condition
    condType := a.analyzeExpr(stmt.Condition)
    if !types.IsBool(condType) {
        a.diags.Error(stmt.Line, stmt.Column,
            "if condition must be bool, got %s", condType)
    }

    // Enter new scope for the then-block
    a.table.EnterScope(fmt.Sprintf("if@line%d", stmt.Line))
    a.analyzeBlock(stmt.ThenBody)
    a.table.ExitScope()

    // Enter new scope for the else-block (if any)
    if stmt.ElseBody != nil {
        a.table.EnterScope(fmt.Sprintf("else@line%d", stmt.Line))
        a.analyzeBlock(stmt.ElseBody)
        a.table.ExitScope()
    }
}
```

The complete scope chain at the deepest point of analysis for the earlier example:

```
Scope Chain during analysis of "print(msg)":

Current: Scope "if@line8" { msg: string }
             ↑ parent
         Scope "fn:main"  { x: int, y: int }
             ↑ parent
         Scope "global"   { add: fn(int,int)->int, main: fn()->void,
                            print: builtin, println: builtin, ... }

Resolve("msg"):
  1. Check "if@line8" → found! return string type ✓

Resolve("x"):
  1. Check "if@line8" → not found
  2. Check "fn:main" → found! return int type ✓

Resolve("undefined_var"):
  1. Check "if@line8" → not found
  2. Check "fn:main" → not found
  3. Check "global" → not found
  4. Return nil, false → emit error "undefined variable 'undefined_var'"
```

---

## Exercises

1. **Implement a multimap**: A map where each key maps to a **list** of values (not just one). E.g., `map["fruits"] = ["apple", "banana", "cherry"]`. Implement `Add(key, value)`, `GetAll(key) []V`, and `Remove(key, value)`.

2. **LFU cache**: Implement a Least Frequently Used cache. Evict the item that has been accessed the fewest times when the cache is full. Harder than LRU but a great hash map exercise.

3. **Two-sum problem**: Given a slice of integers and a target, find two numbers that add up to the target. Solve in O(n) using a hash map.

4. **Group anagrams**: Given a slice of words, group words that are anagrams of each other. Solve in O(n * m) where m is average word length.

5. **Word frequency analysis**: Read a file and count word frequencies. Print the top 10 most frequent words. Handle punctuation and case.

6. **Graph as adjacency list**: Represent a directed graph as a `map[string][]string` (adjacency list). Implement `AddEdge`, `Neighbors`, `HasPath` (DFS).

7. **Consistent hashing simulation**: Simulate a consistent hashing ring with 3 servers. Show which server handles each of 10 keys. Then add a 4th server and show which keys moved.

8. **Astra extension**: Add an "import" feature to the symbol table. When analyzing `import "math"`, pre-populate the current scope with all symbols exported by the math module. Design the module symbol registry data structure.

---

## Summary

| Concept             | Key Point                                                   |
|---------------------|-------------------------------------------------------------|
| Hash function       | Maps arbitrary key → fixed-size integer index              |
| Collision           | Two keys hash to same bucket — unavoidable                 |
| Separate chaining   | Each bucket is a list; simple but poor cache performance    |
| Open addressing     | All entries in array; better cache, needs tombstones        |
| Load factor         | size/capacity; rehash when > 0.75                          |
| Rehashing           | Double capacity and re-insert; O(n) but amortized O(1)     |
| Average O(1)        | Hash maps are NOT worst-case O(1); bad hash → O(n)         |
| Go map              | Built-in, randomized iteration order, nil map panics on write |
| String interning    | One copy of each string; pointer comparison instead of char |
| Unordered           | Map order is random — sort keys if you need order           |
| Consistent hashing  | Distributed hash maps; minimize key movement on resize      |
| Astra symbol table  | Hash map of name→symbol per scope, chained for nested scopes|
