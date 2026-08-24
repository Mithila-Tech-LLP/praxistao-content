---
title: Stack
step: 2
difficulty: easy
estimated: 20 min
---

## What You Are Building

A stack is a Last-In, First-Out (LIFO) collection. Think of a stack of plates: you always add and remove from the top. Stacks underpin function call frames, undo/redo systems, expression parsers, and depth-first search.

```
Push(1) → [1]
Push(2) → [1, 2]
Push(3) → [1, 2, 3]
Pop()   → 3, true   (stack: [1, 2])
Peek()  → 2, true   (stack unchanged)
```

## Key Concepts

**Slice as a stack** — Go slices make a natural stack. The "top" of the stack is the last element of the slice. `append` pushes; slicing `s[:len(s)-1]` pops.

**The comma-ok idiom** — Go functions often return `(value, bool)`. The bool signals whether the value is valid. Use this pattern for `Pop` and `Peek` because they have no valid return value when the stack is empty. Never return a magic number like `-1`.

```go
val, ok := s.Pop()
if !ok {
    fmt.Println("stack was empty")
}
```

**Value semantics vs pointer receivers** — Because your methods modify the slice (`Push` appends, `Pop` shrinks), use pointer receivers: `func (s *Stack) Push(val int)`.

## Struct Signature

```go
type Stack struct {
    data []int
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Push(val int)` | Add val to the top |
| `Pop() (int, bool)` | Remove and return top; ok=false if empty |
| `Peek() (int, bool)` | Return top without removing; ok=false if empty |
| `IsEmpty() bool` | True if no elements |
| `Size() int` | Number of elements |

## Edge Cases to Handle

- `Pop` on empty stack: return `0, false`
- `Peek` on empty stack: return `0, false`
- `Size` on empty stack: return `0`

## Example

```go
s := &Stack{}
fmt.Println(s.IsEmpty()) // true

s.Push(10)
s.Push(20)
s.Push(30)
fmt.Println(s.Size()) // 3

top, ok := s.Peek()
fmt.Println(top, ok) // 30 true

val, ok := s.Pop()
fmt.Println(val, ok) // 30 true
fmt.Println(s.Size()) // 2

_, ok = s.Pop()
_, ok = s.Pop()
_, ok = s.Pop() // empty
fmt.Println(ok) // false
```

## Hints

- `Push`: `s.data = append(s.data, val)`
- `Pop`: check length, grab last element, then `s.data = s.data[:len(s.data)-1]`
- `Peek`: same as Pop but skip the removal step
- `IsEmpty`: `return len(s.data) == 0`
