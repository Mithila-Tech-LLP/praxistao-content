# Chapter 22: Generics

Before Go 1.18, writing a function that works with multiple types meant either duplicating code or using `interface{}` (losing type safety). Generics solve this cleanly: write code once, work with any type that satisfies your constraints. This chapter teaches you Go generics from first principles — type parameters, constraints, and the real-world patterns where generics shine.

## Table of Contents

1. [The Problem Generics Solve](#1-the-problem-generics-solve)
2. [Type Parameters](#2-type-parameters)
3. [Constraints](#3-constraints)
4. [Generic Functions](#4-generic-functions)
5. [Generic Types](#5-generic-types)
6. [Type Inference](#6-type-inference)
7. [Common Generic Patterns](#7-common-generic-patterns)
8. [When NOT to Use Generics](#8-when-not-to-use-generics)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Problem Generics Solve

**Before generics — code duplication:**
```go
func SumInts(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func SumFloat64s(nums []float64) float64 {
    total := 0.0
    for _, n := range nums {
        total += n
    }
    return total
}
// Same logic, different types — have to duplicate for every numeric type
```

**Before generics — using `any` (loses type safety):**
```go
func Sum(nums []any) any {
    total := 0
    for _, n := range nums {
        total += n.(int)  // Type assertion — panics if not int!
    }
    return total
}
// Caller gets no compile-time type checking
```

**With generics — one function, any numeric type:**
```go
func Sum[T int | float64](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

Sum([]int{1, 2, 3})        // returns int
Sum([]float64{1.1, 2.2})   // returns float64
// Type-safe, no duplication
```

### Quick Check
> 1. What two problems did pre-generics Go have with type-flexible code?
> 2. What does the `[T int | float64]` syntax mean?

---

## 2. Type Parameters

A **type parameter** is a placeholder for a concrete type, declared in square brackets:

```go
// Syntax: func FuncName[TypeParam Constraint](params...) return...

func Identity[T any](v T) T {
    return v
}

n := Identity[int](42)       // Explicit: T = int
s := Identity[string]("hi")  // Explicit: T = string
f := Identity(3.14)          // Inferred: T = float64
```

**Multiple type parameters:**
```go
func Map[T, R any](slice []T, fn func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Usage:
names := []string{"alice", "bob", "carol"}
lengths := Map(names, func(s string) int { return len(s) })
// lengths = [5, 3, 5]

nums := []int{1, 2, 3, 4}
doubled := Map(nums, func(n int) int { return n * 2 })
// doubled = [2, 4, 6, 8]
```

**Type parameters in methods vs functions:**
```go
// Type parameters on functions — works fine:
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}

// Type parameters on METHODS — must be declared on the TYPE, not the method:
type Wrapper[T any] struct {
    Value T
}

func (w Wrapper[T]) Get() T {
    return w.Value
}
// Methods on generic types use the type's type parameters, not new ones
```

### Quick Check
> 1. How do you declare a type parameter on a function?
> 2. Can you add a new type parameter directly on a method?
> 3. What does `[T, R any]` declare?

---

## 3. Constraints

A **constraint** defines what operations are available for a type parameter. It's an interface:

```go
// any: no constraints — can only use operations that work on all types (==, assignment)
func Clone[T any](v T) T { return v }

// comparable: can use == and != (can be used as map key)
func Contains[T comparable](s []T, v T) bool {
    for _, item := range s {
        if item == v { return true }
    }
    return false
}
```

**Built-in constraints from `constraints` package (Go 1.21: use `cmp` package):**
```go
// Integer is any integer type:
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float is any float type:
type Float interface { ~float32 | ~float64 }

// Ordered: any type that supports <, >, <=, >=
type Ordered interface {
    Integer | Float | ~string
}
```

**The `~` (tilde) operator — underlying type matching:**
```go
type MyInt int  // MyInt has underlying type int

// Without ~: only accepts exact type int
func PlusOne[T int](v T) T { return v + 1 }
PlusOne(42)        // ok
// PlusOne(MyInt(1)) // COMPILE ERROR: MyInt does not satisfy int

// With ~int: accepts int AND any type with int as underlying type
func PlusOne[T ~int](v T) T { return v + 1 }
PlusOne(42)        // ok
PlusOne(MyInt(1))  // ok — MyInt's underlying type is int
```

**Custom constraints:**
```go
// Constraint requiring a specific method:
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}

// Constraint combining method + type set:
type Number interface {
    ~int | ~int64 | ~float64
}

func Abs[T Number](v T) T {
    if v < 0 {
        return -v
    }
    return v
}

// Using cmp package (Go 1.21+):
import "cmp"

func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

### Quick Check
> 1. What does the `comparable` constraint allow?
> 2. What does `~int` mean in a constraint?
> 3. What is `cmp.Ordered`?

---

## 4. Generic Functions

**`Filter` — keep elements matching a predicate:**
```go
func Filter[T any](slice []T, pred func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

evens := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
// [2, 4]

short := Filter([]string{"hi", "hello", "yo", "world"}, func(s string) bool {
    return len(s) <= 2
})
// ["hi", "yo"]
```

**`Reduce` — fold a slice to a single value:**
```go
func Reduce[T, R any](slice []T, initial R, fn func(R, T) R) R {
    acc := initial
    for _, v := range slice {
        acc = fn(acc, v)
    }
    return acc
}

sum := Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int { return acc + n })
// 15

joined := Reduce([]string{"a", "b", "c"}, "", func(acc, s string) string {
    if acc == "" { return s }
    return acc + "," + s
})
// "a,b,c"
```

**`Keys` and `Values` for maps:**
```go
func Keys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

func Values[K comparable, V any](m map[K]V) []V {
    vals := make([]V, 0, len(m))
    for _, v := range m {
        vals = append(vals, v)
    }
    return vals
}

ages := map[string]int{"Alice": 30, "Bob": 25}
names := Keys(ages)    // ["Alice", "Bob"] (random order)
years := Values(ages)  // [30, 25] (random order)
```

**Pointer helpers:**
```go
// Create a pointer to any value (useful for optional struct fields):
func Ptr[T any](v T) *T {
    return &v
}

user := User{
    Name:  "Alice",
    Score: Ptr(42),  // *int
    Label: Ptr("admin"),  // *string
}
```

### Quick Check
> 1. What does `Filter[T any]` return?
> 2. In `Reduce[T, R any]`, what are `T` and `R`?
> 3. Why is `Keys[K comparable, V any]` — why must K be comparable?

---

## 5. Generic Types

**Generic struct:**
```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T  // Zero value of T
        return zero, false
    }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}

func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Usage:
intStack := &Stack[int]{}
intStack.Push(1)
intStack.Push(2)
v, ok := intStack.Pop()  // v=2, ok=true

strStack := &Stack[string]{}
strStack.Push("hello")
```

**Generic Result type (like Rust's Result):**
```go
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](v T) Result[T]       { return Result[T]{value: v} }
func Err[T any](e error) Result[T]  { return Result[T]{err: e} }

func (r Result[T]) IsOk() bool      { return r.err == nil }
func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic(r.err)
    }
    return r.value
}
func (r Result[T]) UnwrapOr(def T) T {
    if r.err != nil { return def }
    return r.value
}
func (r Result[T]) Error() error    { return r.err }

// Usage:
func divide(a, b float64) Result[float64] {
    if b == 0 {
        return Err[float64](errors.New("division by zero"))
    }
    return Ok(a / b)
}

result := divide(10, 2)
fmt.Println(result.Unwrap())     // 5.0
fmt.Println(result.UnwrapOr(0))  // 5.0

bad := divide(1, 0)
fmt.Println(bad.IsOk())          // false
fmt.Println(bad.UnwrapOr(-1))    // -1.0
```

**Generic Option type (nullable without pointers):**
```go
type Option[T any] struct {
    value    T
    hasValue bool
}

func Some[T any](v T) Option[T]  { return Option[T]{value: v, hasValue: true} }
func None[T any]() Option[T]     { return Option[T]{} }

func (o Option[T]) IsSome() bool    { return o.hasValue }
func (o Option[T]) IsNone() bool    { return !o.hasValue }
func (o Option[T]) Unwrap() T {
    if !o.hasValue { panic("called Unwrap on None") }
    return o.value
}
func (o Option[T]) UnwrapOr(def T) T {
    if !o.hasValue { return def }
    return o.value
}
```

### Quick Check
> 1. How do you instantiate a generic type like `Stack`?
> 2. Inside a generic type's method, how do you get the zero value of `T`?
> 3. What is `Result[T]` modelling?

---

## 6. Type Inference

Go's compiler can often infer type parameters from function arguments:

```go
func Map[T, R any](s []T, fn func(T) R) []R { ... }

// Explicit type arguments:
result := Map[string, int](names, func(s string) int { return len(s) })

// Inferred (compiler figures out T=string, R=int from arguments):
result := Map(names, func(s string) int { return len(s) })
```

**When inference works:**
```go
// Works: T can be inferred from the argument type
Contains([]int{1, 2, 3}, 2)     // T inferred as int
Filter([]string{"a", "b"}, ...) // T inferred as string
Ptr(42)                          // T inferred as int

// Doesn't work: return type can't always be inferred
result := Err[float64](someErr)  // Must specify float64 — no argument to infer from
```

**Type parameters can't be inferred for generic types (only functions):**
```go
s := Stack[int]{}    // Must specify int
s := Stack{}         // COMPILE ERROR: can't infer T
```

### Quick Check
> 1. Does type inference work for generic types or only generic functions?
> 2. When can't the compiler infer a type parameter?

---

## 7. Common Generic Patterns

**Constraint-based numeric sum:**
```go
import "golang.org/x/exp/constraints"  // Or define your own

type Number interface {
    constraints.Integer | constraints.Float
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

**Generic cache:**
```go
type Cache[K comparable, V any] struct {
    mu    sync.RWMutex
    items map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{items: make(map[K]V)}
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.items[k]
    return v, ok
}

func (c *Cache[K, V]) Set(k K, v V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[k] = v
}

// Usage:
userCache := NewCache[int, *User]()
userCache.Set(1, &User{Name: "Alice"})
u, ok := userCache.Get(1)
```

**Generic channel utilities:**
```go
// Collect all values from a channel into a slice:
func Collect[T any](ch <-chan T) []T {
    var result []T
    for v := range ch {
        result = append(result, v)
    }
    return result
}

// Fan-out: send one value to multiple channels:
func Broadcast[T any](v T, channels ...chan<- T) {
    for _, ch := range channels {
        ch <- v
    }
}

// Merge multiple channels into one:
func Merge[T any](channels ...<-chan T) <-chan T {
    out := make(chan T)
    var wg sync.WaitGroup
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan T) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

### Quick Check
> 1. Why does `Cache[K comparable, V any]` require `K` to be `comparable`?
> 2. What does `Collect[T any]` do?

---

## 8. When NOT to Use Generics

Generics add complexity — use them only when the benefit is clear:

**Don't use generics for:**
```go
// 1. Single-type functions — no benefit
func AddInts(a, b int) int { return a + b }  // Fine as-is

// 2. When interfaces already work naturally
func Print(v fmt.Stringer) { fmt.Println(v.String()) }  // Don't genericify this

// 3. When a simple []any + type switch is clearer
func handleEvent(e Event) {
    switch v := e.Data.(type) {
    case UserCreated:
        ...
    case OrderPlaced:
        ...
    }
}
```

**Use generics when:**
```go
// 1. Data structures that must work with any type (Stack, Queue, Tree, Cache)
// 2. Utility functions operating on slices/maps of any type (Map, Filter, Reduce)
// 3. Eliminating code duplication across numeric types
// 4. Type-safe wrappers (Result[T], Option[T], Pair[A, B])
```

**The guideline:** If you'd write the exact same logic twice (once for `int`, once for `string`), and the only difference is the type — that's a generic. If the logic differs per type, use interfaces.

---

## Summary

- **Type parameter**: `[T Constraint]` — placeholder for a concrete type
- **Constraint**: an interface limiting what `T` can be; `any` = no constraint, `comparable` = supports ==
- **`~T`**: matches T and all types with T as underlying type
- **`cmp.Ordered`**: any type supporting `<`, `>`, `<=`, `>=`
- **Generic functions**: `func F[T Constraint](v T) T` — type often inferred
- **Generic types**: `type Stack[T any] struct { items []T }` — type must be specified at instantiation
- **Zero value**: `var zero T` — the zero value of whatever T is
- **Type inference**: compiler infers type params from function arguments, not for generic types
- **When to use**: data structures, slice/map utilities, eliminating numeric type duplication

---

## Exercises

### Easy
1. Write a generic `Reverse[T any](slice []T) []T` that returns a reversed copy. Test with `[]int`, `[]string`, and `[]bool`.
2. Write a generic `Unique[T comparable](slice []T) []T` that returns a slice with duplicates removed, preserving order. Test with `[]int{1,2,1,3,2,4}`.
3. Write a generic `Zip[T, U any](a []T, b []U) []struct{First T; Second U}` that pairs elements at the same index. Zip stops at the shorter slice.

### Medium
4. Generic ordered map: Implement `OrderedMap[K comparable, V any]` that stores key-value pairs and iterates in insertion order (like Python's dict). Methods: `Set(K, V)`, `Get(K) (V, bool)`, `Delete(K)`, `Keys() []K`, `Values() []V`, `Range(func(K, V) bool)`. Internal structure: `map[K]V` + `[]K` to track order. Test: insert 10 items, delete 2, verify Range iterates in correct order.
5. Type-safe event bus: Build `EventBus[T any]` — a pub/sub system for a specific event type. `Subscribe() <-chan T` returns a channel that receives events. `Publish(event T)` broadcasts to all subscribers. `Unsubscribe(ch <-chan T)` removes a subscriber. Instantiate two buses: `EventBus[UserCreated]` and `EventBus[OrderPlaced]` — they're completely separate. Test with 5 concurrent subscribers and 100 published events.
6. Generic pipeline builder: Design `Pipeline[T any]` with chainable methods: `From([]T) *Pipeline[T]`, `Filter(func(T) bool) *Pipeline[T]`, `Transform[R any](fn func(T) R) *Pipeline[R]` (note: this requires a standalone function, not method, since methods can't introduce new type params), `ForEach(func(T))`, `Collect() []T`. Chain example: `From([]int{1..10}).Filter(even).Collect()`.

### Hard
7. Generic B-tree: Implement a generic B-tree `BTree[K cmp.Ordered, V any]` with configurable order (max children per node). Methods: `Insert(key K, value V)`, `Search(key K) (V, bool)`, `Delete(key K) bool`, `Range(minK, maxK K) []V`. Handle node splits and merges. Test with 10,000 random insertions, verify all keys are retrievable, all keys are in sorted order via Range.
8. Generic expression evaluator: Build a type-safe expression evaluator using generics. `Expr[T Number]` represents a numeric expression. Leaf nodes hold a value. Internal nodes hold an operator (+, -, *, /) and two child Exprs. `Eval() T` evaluates the expression. `String() string` returns a human-readable form. Use the visitor pattern with a generic `Visit[T Number, R any](expr Expr[T], visitor ExprVisitor[T, R]) R`. Test with int and float64 expressions.
