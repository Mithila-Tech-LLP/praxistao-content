---
title: Common Go Interview Gotchas
category: Software & Programming
tags: [Go, Interview Prep]
duration: 9 min read
relatedCourses: [go-programming, senior-engineer-interview]
relatedProjects: []
relatedTopics: [channels-vs-mutexes, goroutine-leaks]
---

## TL;DR

- The classic Go interview trip-ups are almost all about things that *look* obviously correct but rely on a subtlety of scoping, interfaces, or slice internals.
- Five that come up constantly: the loop-variable-capture bug, nil interfaces that aren't `== nil`, slice aliasing/append surprises, `defer` argument evaluation timing, and map iteration order.
- None of these are "Go is broken" — each one is a direct, learnable consequence of a specific language rule. Knowing the rule makes the "gotcha" completely predictable.

## 1. The Loop Variable Capture Bug

```go
funcs := make([]func(), 0, 3)
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}
for _, f := range funcs {
    f()
}
```

**Before Go 1.22:** this prints `3 3 3`, not `0 1 2`. Every closure captured the *same* variable `i`, not a snapshot of its value at the time the closure was created — so by the time any of them run, the loop has already finished and `i` is 3.

**Go 1.22+ (2024) changed this:** `for` loops now create a *new* `i` for each iteration, so the same code correctly prints `0 1 2`. This is a genuine language semantics change, not just a lint rule — but it only applies if your `go.mod` declares `go 1.22` or later.

The pre-1.22 fix, and still good practice for clarity regardless of Go version:

```go
for i := 0; i < 3; i++ {
    i := i // shadow: creates a new variable scoped to this iteration
    funcs = append(funcs, func() { fmt.Println(i) })
}
```

An interviewer asking this is testing whether you understand *why* it happened (variable capture by reference, not by value) — not just whether you've memorized the fixed output.

## 2. The Nil Interface Trap

```go
type MyError struct{}
func (e *MyError) Error() string { return "boom" }

func doSomething() error {
    var e *MyError = nil
    return e // returns a non-nil interface wrapping a nil pointer!
}

func main() {
    err := doSomething()
    fmt.Println(err == nil) // false — surprising!
}
```

An interface value in Go is really a pair: `(type, value)`. `doSomething` returns an `error` interface whose *type* is `*MyError` and whose *value* is `nil`. The interface itself is only `== nil` when **both** the type and value are nil. Here the type is set (`*MyError`), so the interface is non-nil even though the underlying pointer is nil.

The fix: return a literal `nil`, not a typed nil variable, when there's no error:

```go
func doSomething() error {
    var e *MyError
    if somethingWentWrong {
        e = &MyError{}
    }
    if e != nil {
        return e
    }
    return nil // explicit, untyped nil — this one IS == nil
}
```

This is one of the most-cited real production bugs in Go codebases, not just an interview trivia question — it shows up whenever a function returns a typed error pointer as a bare `error` interface.

## 3. Slice Aliasing and append Surprises

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3] // b shares a's underlying array: [2, 3]
b = append(b, 99)
fmt.Println(a) // [1 2 3 99 5] — a changed too!
```

A slice is a `(pointer, length, capacity)` triple pointing at a shared underlying array. `a[1:3]` doesn't copy anything — `b` points into the same array as `a`. Since `b`'s capacity (from index 1 to the end of `a`'s array) has room for one more element, `append` writes `99` directly into that shared array at the position `a[3]` used to occupy, silently mutating `a`.

If `b` had been at capacity, `append` would have allocated a brand-new backing array instead — and *then* mutations to `b` would no longer affect `a`. This capacity-dependent behavior (mutate in place vs. silently reallocate) is exactly what makes slice aliasing bugs so unpredictable: the same code can behave differently depending on lengths you didn't think mattered.

The safe pattern when you need an independent copy:

```go
b := make([]int, len(a[1:3]))
copy(b, a[1:3])
```

## 4. defer Argument Evaluation Timing

```go
func process() {
    x := 1
    defer fmt.Println("x was:", x) // arguments evaluated NOW, not at defer-time
    x = 2
}
// prints: x was: 1
```

A `defer` statement evaluates its function's *arguments* immediately, at the point the `defer` line executes — only the actual *call* is delayed until the surrounding function returns. This trips people up specifically when they expect `defer` to capture the final value of a variable.

Contrast with a deferred closure, which *does* see later mutations, because it captures the variable itself, not a value:

```go
defer func() { fmt.Println("x was:", x) }() // prints: x was: 2
```

## 5. Map Iteration Order Is Deliberately Randomized

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k := range m {
    fmt.Println(k) // different order every run, on purpose
}
```

Go's runtime deliberately randomizes map iteration order, specifically so nobody accidentally depends on an order that was never guaranteed. Code that needs a stable order must sort the keys explicitly:

```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
```

## Why These Specific Five

Notice the pattern: every one of these is a case where Go's actual, documented rule is simple and consistent — but the code *looks* like it should behave the way a similar construct behaves in another language, or the way intuition suggests. Interviewers use these not to catch memorization, but to see whether a candidate reasons from the underlying rule (how closures capture variables, what an interface actually is, how slices share memory) instead of pattern-matching on what the code "looks like" it does.
