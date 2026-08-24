# Chapter 13: Pointers

A **pointer** holds the memory address of a value. Instead of working with a copy of data, a pointer lets you work with the original directly. Go's pointers are simpler and safer than C/C++ pointers — no pointer arithmetic, automatic garbage collection — but understanding them is essential for writing correct, efficient Go code.

## Table of Contents

1. [What Is a Pointer](#1-what-is-a-pointer)
2. [Declaring and Using Pointers](#2-declaring-and-using-pointers)
3. [Passing Pointers to Functions](#3-passing-pointers-to-functions)
4. [Pointers to Structs](#4-pointers-to-structs)
5. [new() vs make()](#5-new-vs-make)
6. [When to Use Pointers](#6-when-to-use-pointers)
7. [Common Pointer Mistakes](#7-common-pointer-mistakes)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is a Pointer

Every variable in a program lives at a specific memory address. A pointer stores that address:

```
Memory layout:
  
  address:  0x100   0x104   0x108   0x10C
  value:    [  42 ] [  99 ] [  0  ] [ ... ]
                ↑
                |
  var x int = 42    → stored at address 0x100
  var p *int = &x   → p holds the VALUE 0x100 (an address)
  *p = 99           → writes 99 to address 0x100 (changes x)
```

**The two pointer operators:**
- `&x` — the **address-of** operator: gives you the address of `x`
- `*p` — the **dereference** operator: gives you the value at the address stored in `p`

```go
x := 42
p := &x           // p is a pointer to x — type: *int

fmt.Println(x)    // 42
fmt.Println(p)    // 0xc0000b4000 (some memory address)
fmt.Println(*p)   // 42 — dereference: the value at address p

*p = 99           // Set the value at address p to 99
fmt.Println(x)    // 99 — x changed because p points to x!
```

### Quick Check
> 1. What does `&x` give you?
> 2. What does `*p` give you?
> 3. If you change `*p`, does `x` change?

---

## 2. Declaring and Using Pointers

```go
// Pointer type: *T means "pointer to T"
var p *int         // nil pointer to int (zero value of *int is nil)
var sp *string     // nil pointer to string
var fp *float64    // nil pointer to float64

// Dereference nil pointer = PANIC:
// fmt.Println(*p)  // panic: runtime error: invalid memory address or nil pointer dereference

// Always initialize before use:
n := 42
p = &n
fmt.Println(*p)  // 42

// Pointer to different types:
f := 3.14
fp = &f
*fp = 2.72   // Changes f
fmt.Println(f)  // 2.72

// Short form — create a pointed-to value inline:
p2 := new(int)  // Allocates an int, returns *int (see Section 5)
*p2 = 100
fmt.Println(*p2)  // 100
```

**Pointers to literals** — you can't take the address of a literal directly, but there's a workaround:
```go
// Can't do this:
// p := &42  // ERROR: cannot take address of 42

// But you can do this via a helper or a variable:
n := 42
p := &n

// Helper function for pointer to literal (common pattern):
func ptr[T any](v T) *T { return &v }

p := ptr(42)      // *int pointing to 42
sp := ptr("hello") // *string pointing to "hello"
```

### Quick Check
> 1. What is the zero value of a pointer?
> 2. What happens when you dereference a nil pointer?
> 3. Why can't you write `p := &42`?

---

## 3. Passing Pointers to Functions

**Pass by value (default):** function gets a copy:
```go
func doubleValue(n int) {
    n *= 2  // Only modifies the local copy
}

x := 5
doubleValue(x)
fmt.Println(x)  // 5 — unchanged!
```

**Pass by pointer:** function can modify the original:
```go
func doublePointer(n *int) {
    *n *= 2  // Modifies the value at the address
}

x := 5
doublePointer(&x)  // Pass the address of x
fmt.Println(x)     // 10 — changed!
```

**Real-world example:**
```go
type Config struct {
    Host    string
    Port    int
    Debug   bool
}

// Without pointer: can't modify, must return new value
func applyDefaults(c Config) Config {
    if c.Host == "" { c.Host = "localhost" }
    if c.Port == 0  { c.Port = 8080 }
    return c  // Return the modified copy
}

// With pointer: modify in place
func applyDefaultsPtr(c *Config) {
    if c.Host == "" { c.Host = "localhost" }
    if c.Port == 0  { c.Port = 8080 }
}

cfg := Config{Debug: true}
applyDefaultsPtr(&cfg)
fmt.Println(cfg.Host, cfg.Port)  // localhost 8080
```

**Multiple "return values" via pointer parameters:**
```go
// Common C/Go pattern: use pointer params to "return" multiple things
func divide(a, b float64, result *float64) error {
    if b == 0 {
        return errors.New("division by zero")
    }
    *result = a / b
    return nil
}

// But in Go, it's idiomatic to just return multiple values instead:
func divide(a, b float64) (float64, error) { ... }
// This is much cleaner — only use pointer params when truly needed
```

### Quick Check
> 1. What happens to `x` when you call `double(x)` (value param) vs `double(&x)` (pointer param)?
> 2. When passing a large struct to a function, should you prefer value or pointer?
> 3. Is it more idiomatic to use pointer parameters for "out" values or to return multiple values?

---

## 4. Pointers to Structs

Structs are commonly used via pointers, especially when they're large or need modification:

```go
type User struct {
    ID    int
    Name  string
    Email string
}

// Value: u is a copy
u1 := User{ID: 1, Name: "Alice"}
u1.Name = "Bob"  // Only changes the local copy if u1 was passed as value

// Pointer: up points to the original
up := &User{ID: 2, Name: "Carol"}
up.Name = "Dave"  // Auto-dereferenced: (*up).Name = "Dave"

// Functions that modify structs take pointers:
func updateEmail(u *User, email string) {
    u.Email = email  // Modifies the original
}

updateEmail(up, "dave@example.com")
fmt.Println(up.Email)  // dave@example.com
```

**Auto-dereferencing with structs:** Go automatically dereferences pointers to structs when accessing fields — you don't need to write `(*up).Name`:
```go
up := &User{Name: "Alice"}

// These are equivalent:
fmt.Println((*up).Name)  // "Alice" — explicit dereference
fmt.Println(up.Name)     // "Alice" — auto-dereference (idiomatic)
```

**Pointer comparison:**
```go
u1 := &User{ID: 1}
u2 := &User{ID: 1}
u3 := u1

fmt.Println(u1 == u2)  // false — different memory addresses, even if values equal
fmt.Println(u1 == u3)  // true  — same memory address
fmt.Println(*u1 == *u2) // true  — same value (dereference and compare)
```

### Quick Check
> 1. Do you need to write `(*up).Name` or can you just write `up.Name`?
> 2. Two pointers to structs with identical fields — are they equal with `==`?
> 3. How do you compare the VALUES of two pointer-to-struct?

---

## 5. new() vs make()

**`new(T)`** — allocates memory for type T, zeroes it, returns `*T`:
```go
p := new(int)         // *int — points to a zero int
fmt.Println(*p)       // 0

u := new(User)        // *User — all fields zeroed
u.Name = "Alice"

// Equivalent to:
u2 := &User{}  // Both allocate and return a pointer
```

**`make(T, ...)`** — only for slices, maps, and channels; initializes the internal structure:
```go
s := make([]int, 5)           // Initialized slice (len=5, cap=5)
m := make(map[string]int)     // Initialized (empty) map
ch := make(chan int, 10)       // Buffered channel

// new() for slices just gives you a pointer to a nil slice:
sp := new([]int)    // *[]int — *sp is nil (not useful)
```

**Summary table:**
| | `new(T)` | `make(T, ...)` |
|--|----------|---------------|
| Returns | `*T` (pointer) | `T` (value) |
| Zeroes | Yes | Yes (internally) |
| Works with | Any type | Slice, map, channel only |
| Use case | Allocate struct/primitive | Create slice/map/channel |

**In practice:** `new` is rarely used because `&T{}` is more common and explicit. `make` is always used for slices, maps, and channels.

### Quick Check
> 1. What does `new(int)` return?
> 2. Can you use `new(map[string]int)` to create a usable map?
> 3. What is the difference between what `new` and `make` return?

---

## 6. When to Use Pointers

**Use a pointer when:**

1. **The function needs to modify the caller's data:**
```go
func (u *User) Activate() {
    u.IsActive = true  // Must use pointer receiver to modify
}
```

2. **The type is large and you want to avoid expensive copies:**
```go
type LargeConfig struct {
    // 100+ fields...
}
func processConfig(cfg *LargeConfig) { ... }  // Don't copy 100 fields
```

3. **You need to express "optional" or "no value":**
```go
type SearchFilter struct {
    MinAge *int    // nil means "no minimum age filter"
    MaxAge *int    // nil means "no maximum age filter"
    City   *string // nil means "any city"
}

filter := SearchFilter{
    MinAge: ptr(18),  // Only filter by minimum age
}
```

4. **Sharing the same object across multiple data structures:**
```go
user := &User{ID: 1, Name: "Alice"}
cache[user.ID] = user   // Both cache and activeUsers point to SAME user
activeUsers = append(activeUsers, user)
// Updating user updates it everywhere:
user.Name = "Alice Smith"
// cache[1].Name == "Alice Smith" — consistent
```

**Don't use a pointer when:**
1. The struct is small (≤ 3-4 pointer-sized fields) — copying is fine
2. Immutability is desired — value semantics prevent accidental mutation
3. You need to use the type as a map key — pointers as map keys compare by address, not value

```go
// Small struct — use value, not pointer:
type Point struct{ X, Y float64 }
func distance(a, b Point) float64 {  // Copying two float64s is cheap
    dx, dy := a.X-b.X, a.Y-b.Y
    return math.Sqrt(dx*dx + dy*dy)
}
```

### Quick Check
> 1. Name three situations where you should use a pointer.
> 2. Why might you use `*int` as a struct field instead of `int`?
> 3. If a struct has two float64 fields, should you pass it by pointer or value?

---

## 7. Common Pointer Mistakes

**Mistake 1: Taking the address of a loop variable (fixed in Go 1.22):**
```go
// Before Go 1.22 this was a classic BUG: the loop variable was reused,
// so every &i pointed to the SAME variable:
var ptrs []*int
for i := 0; i < 3; i++ {
    ptrs = append(ptrs, &i)
}
// Pre-1.22 output: 3 3 3 (all pointed at the final value of i)
// Since Go 1.22, each iteration gets a FRESH i, so this prints:
fmt.Println(*ptrs[0], *ptrs[1], *ptrs[2])  // 0 1 2

// On older Go versions, the fix was to shadow the variable:
for i := 0; i < 3; i++ {
    i := i  // Shadow i with a new variable (captures current value)
    ptrs = append(ptrs, &i)
}
```

**Mistake 2: Nil pointer dereference:**
```go
var u *User
fmt.Println(u.Name)  // PANIC: nil pointer dereference

// Always nil-check before dereferencing:
if u != nil {
    fmt.Println(u.Name)
}
```

**Mistake 3: Modifying a copy thinking you're modifying the original:**
```go
func disableUser(u User) {  // Value receiver — gets a copy
    u.IsActive = false       // Only disables the copy!
}

user := User{IsActive: true}
disableUser(user)
fmt.Println(user.IsActive)  // true — not disabled! Bug.

// Fix: use pointer
func disableUser(u *User) {
    u.IsActive = false
}
disableUser(&user)
fmt.Println(user.IsActive)  // false — correct
```

**Mistake 4: Storing pointer to stack variable (not a problem in Go, unlike C):**
```go
// In C, this would be dangerous:
// int* createInt() { int n = 42; return &n; }  // Dangling pointer!

// In Go, this is SAFE — Go will allocate n on the heap:
func createInt() *int {
    n := 42
    return &n  // Safe — Go detects n escapes to heap
}

p := createInt()
fmt.Println(*p)  // 42 — perfectly safe
// Go's escape analysis and garbage collector handle this
```

### Quick Check
> 1. Why do all pointers in `ptrs = append(ptrs, &i)` in a loop end up pointing to the same value?
> 2. Is returning a pointer to a local variable safe in Go?
> 3. What does "escape analysis" do?

---

## Summary

- **Pointer**: a variable storing a memory address; type `*T` is "pointer to T"
- **`&x`**: address-of operator — gives the address of `x`
- **`*p`**: dereference operator — gives the value at address `p`
- **nil pointer**: zero value of pointers; dereferencing causes panic
- **Pass by pointer**: function can modify the original; necessary for large structs
- **Auto-dereference**: Go auto-dereferences `*p.Field` to `p.Field` for struct pointers
- **`new(T)`**: allocates zeroed T, returns `*T`; rarely needed (`&T{}` is more common)
- **When to use**: to modify caller's data, avoid large copies, express optional values, share state
- **Common mistake**: loop variable capture — fixed in Go 1.22 (fresh variable per iteration); shadow with `i := i` on older versions

---

## Exercises

### Easy
1. Write a function `swap(a, b *int)` that swaps two integers using pointers. Verify it works.
2. Create a `Node` struct with `Value int` and `Next *Node`. Build a three-node chain and traverse it, printing each value.
3. Write a function `increment(n *int, amount int)` that adds `amount` to `*n`. Call it 5 times in a loop and print the final value.

### Medium
4. Optional fields: Define a `SearchParams` struct where all fields are optional (use pointer types for `minPrice`, `maxPrice *float64`, `minRating *int`, `category *string`). Write `func search(params SearchParams) []Product` that applies only the non-nil filters. Include a helper function `ptr[T any](v T) *T`. Write tests for: no filters, only min price, all filters.
5. Linked list operations: Using a `Node[T]` struct with pointer to next, implement: `Insert(head **Node[T], value T, pos int)` (insert at position), `Delete(head **Node[T], pos int) bool` (delete by position), `Reverse(head **Node[T])` (reverse in-place), `Merge(a, b *Node[T]) *Node[T]` (merge two sorted lists). Note: modifying head requires `**Node[T]`.
6. Deep copy: Write a `deepCopy[T any](v *T) *T` function that creates a completely independent copy of a struct (no shared pointers). Handle: nested structs, pointer fields, slices, maps. Compare against shallow copy using mutations to verify independence.

### Hard
7. Memory pool: Implement a typed memory pool `Pool[T]` to reduce GC pressure. It reuses freed objects instead of allocating new ones. Methods: `Get() *T` (returns a reused or new object), `Put(v *T)` (returns object to pool, resets it to zero value). Use `sync.Pool` internally but add type safety. Benchmark against raw `new(T)` for allocating 1M objects in a tight loop.
8. Smart pointer: Implement a reference-counted smart pointer `Rc[T]` (inspired by Rust's Rc). `NewRc[T any](v T) *Rc[T]` creates one. `Clone() *Rc[T]` increments reference count and returns a new handle. `Get() *T` returns the value. When all handles are dropped (tracked via finalizer or explicit `Drop()` call), it runs a cleanup function. Thread-safe with `sync/atomic`. Test: create one Rc, clone it 10 times, verify cleanup runs exactly once when all are dropped.
