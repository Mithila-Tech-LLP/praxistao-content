# Chapter 14: Stacks — Last In, First Out

> "The stack is the simplest data structure that is still interesting." — Robert Sedgewick

---

## Overview

Imagine a stack of plates in a cafeteria. You can only add a new plate to the top, and you can only take a plate from the top. You cannot reach into the middle or the bottom without first removing everything above it. This is a **stack**.

The stack is one of the most elegant data structures in computer science. It has only three meaningful operations — push, pop, peek — yet it underpins some of the most important algorithms and systems in computing: function calls, undo/redo, balanced bracket checking, expression evaluation, depth-first search, and even the way CPUs work.

This chapter covers:
- The LIFO property and why it matters
- Push, pop, peek, and isEmpty operations
- Array-backed vs linked list-backed stacks
- The call stack: how your computer tracks function calls
- Stack overflow: what it is and when it happens
- Balanced brackets: a classic stack algorithm
- Shunting-yard algorithm: operator precedence with a stack
- Postfix (Reverse Polish Notation) expression evaluation
- A generic `Stack[T]` in Go
- **Astra Build Milestone**: The parser stack for balanced brace checking and operator precedence

---

## What We're Building

By the end of this chapter you will have a generic, production-quality `Stack[T]` in Go, and you will understand how the Astra parser uses two stacks internally — one for tracking open delimiters and one for implementing the shunting-yard algorithm — to correctly parse expressions like `2 + 3 * 4 - 1`.

---

## Table of Contents

1. The Stack Abstraction
2. The LIFO Property
3. Push, Pop, Peek, IsEmpty
4. Array-Backed Stack
5. Linked List-Backed Stack
6. The Call Stack — What Your CPU Does
7. Stack Overflow: What It Is and When It Happens
8. Application: Balanced Brackets Checker
9. Application: Undo / Redo System
10. Application: Postfix Expression Evaluation (RPN)
11. Dijkstra's Shunting-Yard Algorithm
12. Generic Stack[T] in Go
13. Astra Build Milestone: The Parser Stack

---

## 1. The Stack Abstraction

A **stack** is a linear data structure that follows the **Last In, First Out (LIFO)** principle. The last item added is the first item removed.

```
          PUSH 5
           |
           v
    +------+------+
    |      5      |   <- TOP
    +------+------+
    |      3      |
    +------+------+
    |      7      |
    +------+------+
    |      1      |   <- BOTTOM
    +------+------+

          POP
           |
    returns 5

    +------+------+
    |      3      |   <- TOP (new)
    +------+------+
    |      7      |
    +------+------+
    |      1      |
    +------+------+
```

The stack's interface is intentionally minimal. You cannot access the 3rd element from the top directly — you must pop 5, then pop 3 (or just peek without popping). This restriction is the point: it enforces a usage pattern, and that pattern turns out to be exactly what many algorithms need.

---

## 2. The LIFO Property

LIFO stands for **Last In, First Out**. This is the defining property of a stack.

Compare it to a queue (FIFO — First In, First Out), which is like the line at a grocery store: the first person in line is the first to be served.

A stack is the opposite: the last plate you put on the pile is the first one you pick up.

LIFO shows up everywhere:
- **Function calls**: the most recently called function is the first to return
- **Undo operations**: the most recent action is the first to be undone
- **Browser back button**: the most recently visited page is the first to go back to
- **Compiler parsers**: the most recently opened brace is the first to be closed

---

## 3. Push, Pop, Peek, IsEmpty

The complete interface of a stack:

| Operation    | Description                                 | Complexity |
|--------------|---------------------------------------------|------------|
| `Push(item)` | Add item to the top                         | O(1)       |
| `Pop()`      | Remove and return the top item              | O(1)       |
| `Peek()`     | Return the top item without removing it     | O(1)       |
| `IsEmpty()`  | Return true if the stack has no items       | O(1)       |
| `Size()`     | Return the number of items                  | O(1)       |

All operations are O(1) — constant time, regardless of how many items are in the stack. This is what makes stacks so powerful for algorithm design.

---

## 4. Array-Backed Stack

The simplest implementation uses a dynamic array (slice in Go) where the top of the stack is the last element:

```go
package stack

import "fmt"

// ArrayStack[T] is a generic stack backed by a slice.
type ArrayStack[T any] struct {
    items []T
}

// NewArrayStack creates an empty stack with optional initial capacity.
func NewArrayStack[T any](capacity int) *ArrayStack[T] {
    return &ArrayStack[T]{items: make([]T, 0, capacity)}
}

// Push adds an item to the top of the stack — O(1) amortized.
func (s *ArrayStack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop removes and returns the top item — O(1).
// Returns zero value and false if stack is empty.
func (s *ArrayStack[T]) Pop() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}

// Peek returns the top item without removing it — O(1).
func (s *ArrayStack[T]) Peek() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

// IsEmpty returns true if the stack has no items.
func (s *ArrayStack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

// Size returns the number of items in the stack.
func (s *ArrayStack[T]) Size() int {
    return len(s.items)
}

// String returns a visual representation of the stack.
func (s *ArrayStack[T]) String() string {
    if s.IsEmpty() {
        return "Stack: [empty]"
    }
    result := "Stack (top to bottom):\n"
    for i := len(s.items) - 1; i >= 0; i-- {
        result += fmt.Sprintf("  | %v |\n", s.items[i])
    }
    result += "  +---+"
    return result
}
```

The array-backed stack has excellent cache performance — all items are contiguous in memory. Go's `append` operation handles resizing automatically (doubling capacity when needed), making `Push` O(1) amortized.

---

## 5. Linked List-Backed Stack

An alternative implementation uses a linked list. The head of the list is the top of the stack:

```go
package stack

// node is a private linked list node.
type node[T any] struct {
    data T
    next *node[T]
}

// LinkedStack[T] is a stack backed by a singly linked list.
type LinkedStack[T any] struct {
    top  *node[T]
    size int
}

func NewLinkedStack[T any]() *LinkedStack[T] {
    return &LinkedStack[T]{}
}

// Push inserts at the head — O(1), no resizing needed.
func (s *LinkedStack[T]) Push(item T) {
    s.top = &node[T]{data: item, next: s.top}
    s.size++
}

// Pop removes from the head — O(1).
func (s *LinkedStack[T]) Pop() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }
    val := s.top.data
    s.top = s.top.next
    s.size--
    return val, true
}

// Peek reads from the head without removing — O(1).
func (s *LinkedStack[T]) Peek() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }
    return s.top.data, true
}

func (s *LinkedStack[T]) IsEmpty() bool { return s.size == 0 }
func (s *LinkedStack[T]) Size() int     { return s.size }
```

**Array vs Linked List-backed Stack comparison:**

| Aspect              | Array-backed           | Linked List-backed         |
|---------------------|------------------------|----------------------------|
| Memory              | Contiguous, cache-friendly | Scattered, pointer overhead |
| Push                | O(1) amortized         | O(1) always                |
| Pop                 | O(1)                   | O(1)                       |
| Memory wasted       | Up to 2x (pre-allocated) | Per-node pointer overhead  |
| Recommended for     | Most cases             | When max size is unknown and memory is precious |

In practice, **use the array-backed stack**. Go slices are fast and the GC handles memory well.

---

## 6. The Call Stack — What Your CPU Does

Every time your program calls a function, the runtime pushes a **stack frame** onto the call stack. Every time a function returns, its frame is popped off.

```
Program:
    fn main() {
        let x = add(3, 4)
    }
    fn add(a: int, b: int) -> int {
        return multiply(a, b)
    }
    fn multiply(a: int, b: int) -> int {
        return a * b
    }

Call Stack at the deepest point:

+----------------------------------+
| multiply: a=3, b=4               |  <- TOP (currently running)
| return address: add+line5        |
+----------------------------------+
| add: a=3, b=4                    |
| return address: main+line2       |
+----------------------------------+
| main: x=?                        |  <- BOTTOM
| return address: OS               |
+----------------------------------+
```

Each stack frame contains:
- **Local variables** for that function call
- **Parameters** passed to the function
- **Return address** — where to jump when the function returns
- **Saved registers** — processor state to restore on return

When `multiply` returns `12`:
1. The return value `12` is placed in a register
2. The `multiply` frame is popped
3. Execution jumps to the return address (`add+line5`)
4. `add` receives `12`, returns it
5. The `add` frame is popped
6. Execution jumps to `main+line2`
7. `x` is set to `12`

This is the call stack in action. It is literally a stack data structure, implemented in hardware (the CPU has a dedicated stack pointer register, `rsp` on x86-64).

---

## 7. Stack Overflow: What It Is and When It Happens

A **stack overflow** occurs when the call stack grows beyond its allocated limit (typically 1-8 MB on most systems).

```
Infinite recursion:
    fn count(n: int) {
        print(n)
        count(n + 1)   // never stops
    }

Call Stack fills up:
+------------------------+
| count(n=10000)         |  <- TOP
+------------------------+
| count(n=9999)          |
+------------------------+
| count(n=9998)          |
+------------------------+
| ...thousands of frames |
+------------------------+
| count(n=1)             |  <- BOTTOM
+------------------------+
| main                   |
+------------------------+

CRASH: stack overflow / segmentation fault
```

Common causes:
1. **Infinite recursion**: a recursive function without a proper base case
2. **Mutual recursion**: function A calls B, B calls A, forever
3. **Extremely deep recursion**: valid recursion on very large inputs (e.g., building a linked list recursively with 100,000 nodes)

Go handles this by growing the goroutine stack dynamically (starting at 8KB, growing to up to 1GB by default), but eventually even Go programs will crash with "goroutine stack exceeds limit."

The fix is usually **iteration** instead of recursion, or using an **explicit stack** (a `Stack[T]` data structure) to simulate the recursion without using the call stack.

---

## 8. Application: Balanced Brackets Checker

Given a string like `"({[]})"`, determine if the brackets are balanced. This is a classic stack interview problem.

Rules:
- Every opening bracket must have a matching closing bracket
- The matching must be in the correct order (LIFO)

```go
package main

import "fmt"

func isBalanced(s string) bool {
    stack := NewArrayStack[rune](len(s))

    pairs := map[rune]rune{
        ')': '(',
        ']': '[',
        '}': '{',
    }

    for _, ch := range s {
        switch ch {
        case '(', '[', '{':
            stack.Push(ch)   // opening: push onto stack
        case ')', ']', '}':
            top, ok := stack.Pop()
            if !ok {
                return false  // closing bracket with nothing on stack
            }
            if top != pairs[ch] {
                return false  // mismatched pair, e.g. '(' closed by ']'
            }
        }
    }

    return stack.IsEmpty()  // stack must be empty: all brackets matched
}

func main() {
    tests := []struct {
        input    string
        expected bool
    }{
        {"({[]})", true},
        {"(())", true},
        {"fn add(a: int) { return a + 1 }", true},
        {"({[}])", false},
        {"((())", false},
        {")", false},
    }

    for _, test := range tests {
        result := isBalanced(test.input)
        status := "PASS"
        if result != test.expected {
            status = "FAIL"
        }
        fmt.Printf("[%s] isBalanced(%q) = %v\n", status, test.input, result)
    }
}
```

Trace through `"({[]})"`:

```
char   action         stack state
'('    push           [ ( ]
'{'    push           [ (, { ]
'['    push           [ (, {, [ ]
']'    pop '[' ✓      [ (, { ]
'}'    pop '{' ✓      [ ( ]
')'    pop '(' ✓      [ ]

stack is empty → BALANCED ✓
```

---

## 9. Application: Undo / Redo System

Every text editor implements undo/redo with two stacks:

```go
package main

import "fmt"

type Action struct {
    Description string
    Apply       func()
    Reverse     func()
}

type EditHistory struct {
    undoStack *ArrayStack[Action]
    redoStack *ArrayStack[Action]
}

func NewEditHistory() *EditHistory {
    return &EditHistory{
        undoStack: NewArrayStack[Action](100),
        redoStack: NewArrayStack[Action](100),
    }
}

func (h *EditHistory) Do(action Action) {
    action.Apply()
    h.undoStack.Push(action)
    // Any new action clears the redo stack
    h.redoStack = NewArrayStack[Action](100)
    fmt.Printf("  Did: %s\n", action.Description)
}

func (h *EditHistory) Undo() {
    action, ok := h.undoStack.Pop()
    if !ok {
        fmt.Println("  Nothing to undo")
        return
    }
    action.Reverse()
    h.redoStack.Push(action)
    fmt.Printf("  Undid: %s\n", action.Description)
}

func (h *EditHistory) Redo() {
    action, ok := h.redoStack.Pop()
    if !ok {
        fmt.Println("  Nothing to redo")
        return
    }
    action.Apply()
    h.undoStack.Push(action)
    fmt.Printf("  Redid: %s\n", action.Description)
}

func main() {
    doc := ""
    history := NewEditHistory()

    history.Do(Action{
        Description: "type 'Hello'",
        Apply:   func() { doc += "Hello" },
        Reverse: func() { doc = doc[:len(doc)-5] },
    })

    history.Do(Action{
        Description: "type ' World'",
        Apply:   func() { doc += " World" },
        Reverse: func() { doc = doc[:len(doc)-6] },
    })

    fmt.Println("Document:", doc)  // Hello World

    history.Undo()
    fmt.Println("Document:", doc)  // Hello

    history.Undo()
    fmt.Println("Document:", doc)  // (empty)

    history.Redo()
    fmt.Println("Document:", doc)  // Hello
}
```

---

## 10. Application: Postfix Expression Evaluation (RPN)

**Postfix notation** (also called Reverse Polish Notation, or RPN) places operators after their operands:

```
Infix:   3 + 4 * 2
Postfix: 3 4 2 * +
```

Postfix is beautiful because it requires no parentheses and no operator precedence rules — a simple stack evaluator handles it:

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
)

func evalRPN(expression string) (float64, error) {
    stack := NewArrayStack[float64](16)
    tokens := strings.Fields(expression)

    for _, token := range tokens {
        switch token {
        case "+", "-", "*", "/":
            b, ok1 := stack.Pop()
            a, ok2 := stack.Pop()
            if !ok1 || !ok2 {
                return 0, fmt.Errorf("not enough operands for %s", token)
            }
            switch token {
            case "+": stack.Push(a + b)
            case "-": stack.Push(a - b)
            case "*": stack.Push(a * b)
            case "/":
                if b == 0 {
                    return 0, fmt.Errorf("division by zero")
                }
                stack.Push(a / b)
            }
        default:
            num, err := strconv.ParseFloat(token, 64)
            if err != nil {
                return 0, fmt.Errorf("unknown token: %s", token)
            }
            stack.Push(num)
        }
    }

    result, ok := stack.Pop()
    if !ok || !stack.IsEmpty() {
        return 0, fmt.Errorf("invalid expression")
    }
    return result, nil
}

func main() {
    expressions := []string{
        "3 4 +",          // 7
        "10 2 /",         // 5
        "3 4 2 * +",      // 3 + (4*2) = 11
        "5 1 2 + 4 * + 3 -", // 5 + ((1+2)*4) - 3 = 14
    }
    for _, expr := range expressions {
        result, err := evalRPN(expr)
        if err != nil {
            fmt.Printf("%q → ERROR: %v\n", expr, err)
        } else {
            fmt.Printf("%q → %.0f\n", expr, result)
        }
    }
}
```

Trace of `"3 4 2 * +"`:

```
token   action         stack
3       push 3         [3]
4       push 4         [3, 4]
2       push 2         [3, 4, 2]
*       pop 2,4; push 8 [3, 8]
+       pop 8,3; push 11 [11]

Result: 11
```

---

## 11. Dijkstra's Shunting-Yard Algorithm

**Shunting-yard** converts infix expressions (what humans write: `3 + 4 * 2`) to postfix expressions (what computers evaluate easily: `3 4 2 * +`). It was invented by Edsger Dijkstra in 1961.

The algorithm uses two stacks: an output queue and an operator stack.

```go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

type Assoc int
const (Left Assoc = iota; Right)

type OpInfo struct {
    precedence int
    assoc      Assoc
}

var ops = map[string]OpInfo{
    "+": {1, Left},
    "-": {1, Left},
    "*": {2, Left},
    "/": {2, Left},
    "^": {3, Right},  // right-associative: 2^3^2 = 2^(3^2) = 512
}

func isOperator(token string) bool {
    _, ok := ops[token]
    return ok
}

// ShuntingYard converts infix to postfix.
func ShuntingYard(tokens []string) []string {
    output := make([]string, 0)
    opStack := NewArrayStack[string](len(tokens))

    for _, token := range tokens {
        switch {
        case isNumber(token):
            output = append(output, token)

        case isOperator(token):
            for {
                top, ok := opStack.Peek()
                if !ok || top == "(" {
                    break
                }
                topInfo := ops[top]
                curInfo := ops[token]
                if topInfo.precedence > curInfo.precedence ||
                    (topInfo.precedence == curInfo.precedence && curInfo.assoc == Left) {
                    opStack.Pop()
                    output = append(output, top)
                } else {
                    break
                }
            }
            opStack.Push(token)

        case token == "(":
            opStack.Push(token)

        case token == ")":
            for {
                top, ok := opStack.Pop()
                if !ok {
                    panic("mismatched parentheses")
                }
                if top == "(" {
                    break
                }
                output = append(output, top)
            }
        }
    }

    for !opStack.IsEmpty() {
        top, _ := opStack.Pop()
        if top == "(" {
            panic("mismatched parentheses")
        }
        output = append(output, top)
    }

    return output
}

func isNumber(s string) bool {
    for _, r := range s {
        if !unicode.IsDigit(r) && r != '.' {
            return false
        }
    }
    return len(s) > 0
}

func main() {
    expressions := []string{
        "3 + 4 * 2",
        "( 3 + 4 ) * 2",
        "3 + 4 * 2 / ( 1 - 5 ) ^ 2 ^ 3",
    }
    for _, expr := range expressions {
        tokens := strings.Fields(expr)
        postfix := ShuntingYard(tokens)
        fmt.Printf("Infix:   %s\n", expr)
        fmt.Printf("Postfix: %s\n\n", strings.Join(postfix, " "))
    }
}
```

---

## 12. Generic Stack[T] in Go

Here is the complete, production-quality generic stack combining everything above:

```go
package stack

import (
    "fmt"
    "strings"
)

// Stack[T] is a generic LIFO stack backed by a slice.
// T can be any type.
type Stack[T any] struct {
    items []T
}

// New creates an empty stack with optional initial capacity hint.
func New[T any](capacityHint ...int) *Stack[T] {
    cap := 16
    if len(capacityHint) > 0 {
        cap = capacityHint[0]
    }
    return &Stack[T]{items: make([]T, 0, cap)}
}

// Push adds item to the top of the stack.
func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop removes and returns the top item.
// Returns (zero, false) if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }
    n := len(s.items) - 1
    top := s.items[n]
    s.items = s.items[:n]
    return top, true
}

// MustPop pops and panics if the stack is empty.
func (s *Stack[T]) MustPop() T {
    val, ok := s.Pop()
    if !ok {
        panic("stack: Pop on empty stack")
    }
    return val
}

// Peek returns the top item without removing it.
// Returns (zero, false) if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

// IsEmpty returns true if the stack contains no items.
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }

// Size returns the number of items in the stack.
func (s *Stack[T]) Size() int { return len(s.items) }

// Clear removes all items from the stack.
func (s *Stack[T]) Clear() { s.items = s.items[:0] }

// ToSlice returns a copy of all items from bottom to top.
func (s *Stack[T]) ToSlice() []T {
    result := make([]T, len(s.items))
    copy(result, s.items)
    return result
}

// String returns a visual representation (top at the top).
func (s *Stack[T]) String() string {
    if s.IsEmpty() {
        return "Stack[]"
    }
    parts := make([]string, len(s.items))
    for i, item := range s.items {
        parts[len(s.items)-1-i] = fmt.Sprintf("%v", item)
    }
    return "Stack[" + strings.Join(parts, ", ") + "] (top first)"
}
```

---

## 13. Astra Build Milestone: The Parser Stack

The Astra parser uses stacks in two critical ways.

### Stack 1: Balanced Delimiter Checking

During lexing, before the parser even runs, we validate that all opening braces, brackets, and parentheses are properly closed:

```go
// lexer/validation.go

package lexer

import "fmt"

type DelimiterError struct {
    Opening Token
    Line    int
    Col     int
    Msg     string
}

func (e DelimiterError) Error() string {
    return fmt.Sprintf("line %d:%d: %s (opened at %d:%d)",
        e.Line, e.Col, e.Msg, e.Opening.Line, e.Opening.Column)
}

// ValidateDelimiters checks that all brackets/braces/parens are balanced.
// This gives much better error messages than letting the parser fail.
func ValidateDelimiters(tokens []Token) []error {
    type openDelim struct {
        tok Token
    }

    stack := New[openDelim]()
    var errors []error

    closing := map[TokenType]TokenType{
        RPAREN:   LPAREN,
        RBRACE:   LBRACE,
        RBRACKET: LBRACKET,
    }

    matching := map[TokenType]string{
        LPAREN:   "(",
        LBRACE:   "{",
        LBRACKET: "[",
        RPAREN:   ")",
        RBRACE:   "}",
        RBRACKET: "]",
    }

    for _, tok := range tokens {
        switch tok.Type {
        case LPAREN, LBRACE, LBRACKET:
            stack.Push(openDelim{tok: tok})

        case RPAREN, RBRACE, RBRACKET:
            expected := closing[tok.Type]
            top, ok := stack.Pop()
            if !ok {
                errors = append(errors, fmt.Errorf(
                    "line %d:%d: unexpected %q — no matching opening bracket",
                    tok.Line, tok.Column, matching[tok.Type]))
                continue
            }
            if top.tok.Type != expected {
                errors = append(errors, fmt.Errorf(
                    "line %d:%d: expected %q to close %q opened at %d:%d",
                    tok.Line, tok.Column,
                    matching[top.tok.Type], matching[top.tok.Type],
                    top.tok.Line, top.tok.Column))
            }
        }
    }

    // Anything left on the stack is unclosed
    for !stack.IsEmpty() {
        open, _ := stack.Pop()
        errors = append(errors, fmt.Errorf(
            "line %d:%d: %q was never closed",
            open.tok.Line, open.tok.Column,
            matching[open.tok.Type]))
    }

    return errors
}
```

### Stack 2: Operator Precedence Parsing

The Astra parser uses the shunting-yard algorithm inside `parseExpression` to handle operator precedence without writing a ton of special cases:

```go
// parser/expression.go

package parser

import (
    "your-module/lexer"
    "your-module/ast"
)

// OperatorInfo holds precedence and associativity for binary operators.
type OperatorInfo struct {
    Precedence int
    RightAssoc bool
}

var operatorTable = map[lexer.TokenType]OperatorInfo{
    lexer.OR:      {Precedence: 1},
    lexer.AND:     {Precedence: 2},
    lexer.EQ:      {Precedence: 3},
    lexer.NEQ:     {Precedence: 3},
    lexer.LT:      {Precedence: 4},
    lexer.GT:      {Precedence: 4},
    lexer.LEQ:     {Precedence: 4},
    lexer.GEQ:     {Precedence: 4},
    lexer.PLUS:    {Precedence: 5},
    lexer.MINUS:   {Precedence: 5},
    lexer.STAR:    {Precedence: 6},
    lexer.SLASH:   {Precedence: 6},
    lexer.PERCENT: {Precedence: 6},
}

// parseExpression uses the shunting-yard algorithm to parse binary expressions
// with correct operator precedence.
func (p *Parser) parseExpression() (ast.Expr, error) {
    outputQueue := make([]ast.Expr, 0, 8)    // output: expressions
    opStack := New[lexer.Token](8)            // operator stack

    // Helper: pop one operator and combine top two operands
    popOp := func() error {
        op, _ := opStack.Pop()
        if len(outputQueue) < 2 {
            return fmt.Errorf("line %d: not enough operands for %q",
                op.Line, op.Lexeme)
        }
        right := outputQueue[len(outputQueue)-1]
        left := outputQueue[len(outputQueue)-2]
        outputQueue = outputQueue[:len(outputQueue)-2]
        outputQueue = append(outputQueue, &ast.BinaryExpr{
            Left: left, Op: op.Lexeme, Right: right,
        })
        return nil
    }

    for {
        // Parse a primary expression (number, identifier, paren group)
        primary, err := p.parsePrimary()
        if err != nil {
            return nil, err
        }
        outputQueue = append(outputQueue, primary)

        // Check if the next token is a binary operator
        opInfo, isBinOp := operatorTable[p.peek().Type]
        if !isBinOp {
            break
        }
        opTok := p.advance()  // consume the operator

        // Pop operators from the stack with higher/equal precedence (left-assoc)
        for !opStack.IsEmpty() {
            topTok, _ := opStack.Peek()
            topInfo, ok := operatorTable[topTok.Type]
            if !ok { break }
            if topInfo.Precedence > opInfo.Precedence ||
                (topInfo.Precedence == opInfo.Precedence && !opInfo.RightAssoc) {
                if err := popOp(); err != nil {
                    return nil, err
                }
            } else {
                break
            }
        }
        opStack.Push(opTok)
    }

    // Pop all remaining operators
    for !opStack.IsEmpty() {
        if err := popOp(); err != nil {
            return nil, err
        }
    }

    if len(outputQueue) != 1 {
        return nil, fmt.Errorf("invalid expression")
    }
    return outputQueue[0], nil
}
```

This correctly parses `2 + 3 * 4 - 1` as `(2 + (3 * 4)) - 1 = 13`, respecting multiplication's higher precedence over addition and subtraction.

The call stack diagram for this parse:

```
Expression: 2 + 3 * 4 - 1

opStack (operators):   outputQueue (expressions):
Initially empty        []

After "2":             [2]
After "+":    [+]      [2]
After "3":    [+]      [2, 3]
After "*":    [+, *]   [2, 3]       (* has higher precedence, stays on stack
After "4":    [+, *]   [2, 3, 4]
After "-":    [+]      [2, BinExpr(3*4)]  (* popped, then + same prec as -, popped)
              [-]      [BinExpr(2+(3*4))]
After "1":    [-]      [BinExpr(2+(3*4)), 1]
End:          []       [BinExpr(BinExpr(2+(3*4))-1)]

Final AST:
    BinaryExpr(-)
    /           \
BinaryExpr(+)   1
/          \
2       BinaryExpr(*)
        /           \
        3            4
```

---

## Exercises

1. **Stack with minimum**: Implement a stack that supports `Push`, `Pop`, and `GetMin` — all in O(1). Hint: maintain a second internal stack that tracks minimums.

2. **Sort a stack** using only stack operations (push, pop, peek, isEmpty). You may use one additional stack as scratch space.

3. **Evaluate infix expressions**: Extend the shunting-yard algorithm to handle unary negation (e.g., `-3 + 5`).

4. **Implement a calculator**: Build a full calculator that takes a string like `"3 + (4 * 2) / (1 - 5)"` and returns the result.

5. **Browser history**: Implement a simplified browser history with `visit(url)`, `back()`, `forward()`, and `currentPage()` using two stacks.

6. **Postfix to infix**: Write a function that converts a postfix expression back to an infix expression with correct parentheses.

7. **Largest rectangle in histogram**: Given an array of heights representing a histogram, find the largest rectangle. This is a classic hard problem solved elegantly with a stack.

8. **Astra challenge**: The Astra parser uses a stack for tracking nested scopes (function bodies, if blocks, for loops). Every time we encounter `{`, we push a new scope. Every time we encounter `}`, we pop the scope. Implement a `ScopeStack` that tracks variable declarations in nested scopes, supporting `Define(name, type)`, `Lookup(name)`, and scope push/pop.

---

## Summary

| Concept                | Key Point                                                   |
|------------------------|-------------------------------------------------------------|
| Stack                  | LIFO: Last In, First Out                                    |
| Push                   | O(1) — add to top                                          |
| Pop                    | O(1) — remove from top                                     |
| Peek                   | O(1) — read top without removing                           |
| Array-backed           | Cache-friendly, amortized O(1) push, recommended           |
| Linked-list backed     | No reallocation, slightly worse cache performance           |
| Call stack             | Hardware stack tracking function calls and local variables  |
| Stack overflow         | Call stack exceeds limit; fix with iteration or explicit stack |
| Balanced brackets      | Classic O(n) stack algorithm                               |
| Undo/redo              | Two stacks: undo stack and redo stack                       |
| Postfix (RPN)          | No precedence rules needed; trivial to evaluate with stack  |
| Shunting-yard          | Converts infix to postfix using an operator stack           |
| Astra usage            | Delimiter validation + operator precedence in parser        |
