# Chapter 06: Operators and Expressions — Building Blocks of Computation

> *"An expression is a question the computer answers. An operator is the verb that asks it."*
> — Every compiler textbook, paraphrased

Every useful program does work: it adds numbers, compares values, concatenates text, flips bits. All of that work happens through **expressions** and **operators**. Before your Astra compiler can generate a single line of machine code, it needs to understand what `2 + 3 * 4` means — which subexpression is calculated first, what type the result has, and how to represent that meaning as a tree in memory.

This chapter is the foundation of Astra's expression system. We will build every piece from scratch: what an expression is, what types of operators Astra supports, how precedence rules resolve ambiguity, how the compiler turns a flat string of tokens into a structured tree, and finally how to define Go data structures (AST nodes) that represent every kind of expression Astra can contain. By the end you will have a complete `ast/expressions.go` file and a tiny expression evaluator that can calculate constant-folded results at compile time.

## What We're Building

In Astra, expressions appear everywhere:

```astra
let x = 2 + 3 * 4          // arithmetic expression
let ok = x > 10 && x < 20  // comparison + logical expression
let msg = "Score: " + x.to_string()  // string concat + method call
let flags = 0b1010 | 0b0101         // bitwise OR
```

The parser will eventually call `parseExpression()` and get back an **AST node** — a Go struct that captures the structure of the expression. The code generator will walk that node to emit LLVM IR or assembly. Everything depends on getting this layer right.

## Table of Contents

1. What Is an Expression?
2. Arithmetic Operators — The Workhorses
3. Comparison Operators — Asking Questions
4. Logical Operators — Combining Truth
5. Bitwise Operators — Working at the Bit Level
6. Operator Precedence — Who Goes First?
7. Expression Trees — The AST Preview
8. String Operations in Astra
9. Astra Build Milestone — Expression AST Nodes
10. Exercises
11. Summary

---

## 1. What Is an Expression?

Think of a recipe. A recipe takes ingredients and produces a dish. An **expression** is the same idea: it takes values and produces a new value.

```
Ingredients + Process = Dish
2 + 3                = 5
```

A **statement**, by contrast, is an instruction that *does* something but does not itself produce a value you can use elsewhere. In Astra:

```astra
let x = 2 + 3 * 4   // STATEMENT — declares x, does not produce a value
2 + 3 * 4           // EXPRESSION — produces 14, could be used anywhere
```

The distinction matters because in Astra (like Go, Java, and C) you cannot write:

```astra
let y = let x = 5   // ERROR: statement cannot be used as a value
let y = 2 + 3       // OK: expression on the right
```

### Every Expression Has a Type

The Astra type system assigns a **type** to every expression. The type tells the compiler what operations are valid and what memory layout to use:

```astra
2 + 3           // type: int    (integer arithmetic)
2.0 + 3.0       // type: float  (floating-point arithmetic)
"Hello" + "!"   // type: string (concatenation)
x > 5           // type: bool   (comparison produces boolean)
```

If you try to mix incompatible types, Astra's type checker rejects it:

```astra
let x = 2 + "hello"   // ERROR: cannot add int and string
```

### Nested Expressions

Expressions can contain other expressions — that is their power. The subexpressions are evaluated first, and their results feed into the outer expression:

```astra
(2 + 3) * (10 - 4)
// Step 1: 2 + 3  = 5
// Step 2: 10 - 4 = 6
// Step 3: 5 * 6  = 30
```

This nesting is why we represent expressions as **trees**, not flat lists. We will see the tree form later in section 7.

---

## 2. Arithmetic Operators — The Workhorses

Arithmetic operators take numeric values and produce numeric values. Astra supports the standard five, plus string concatenation with `+`.

### Addition `+`

```astra
let a = 10 + 5      // 15  (int + int → int)
let b = 3.14 + 1.0  // 4.14 (float + float → float)
let s = "Hi" + "!"  // "Hi!" (string + string → string)
```

The `+` operator is **overloaded** — the same symbol does different things depending on the types of its operands. The type checker resolves which meaning to use.

### Subtraction `-` and Unary Negation

```astra
let diff = 20 - 8       // 12
let neg  = -5           // unary minus: produces -5
let neg2 = -(3 + 4)     // -(7) = -7
```

Unary minus is different from binary minus: `5 - 3` uses two operands; `-5` uses one. In the AST, these are different node types (`BinaryExpr` vs `UnaryExpr`).

### Multiplication `*`

```astra
let area = 6 * 7        // 42
let sq   = 4 * 4        // 16
```

### Division `/` — The Integer Division Trap

This is where many beginners get surprised. In Astra (following mathematical convention for integer types):

```astra
let a = 7 / 2       // 3, NOT 3.5!  (integer division truncates toward zero)
let b = 7.0 / 2.0   // 3.5          (float division is precise)
let c = 7 / 2.0     // ERROR: cannot divide int by float (type mismatch)
```

Why does `7 / 2 = 3`? Because integers have no fractional part. The result is truncated toward zero — the decimal part is simply discarded. This follows the convention of C, Go, Java, and Rust.

```
Analogy: cutting a pizza with 7 slices among 2 people.
Each person gets 3 full slices. The leftover 1 slice is the REMAINDER.
The remainder is what % (modulo) captures.
```

**Common mistake:**

```astra
fn average(a: int, b: int) -> float {
    return (a + b) / 2      // BUG: this is integer division!
    // Fix:
    return (a + b).to_float() / 2.0
}
```

### Modulo `%` — The Remainder

```astra
let rem = 7 % 2     // 1  (7 = 3 × 2 + 1)
let rem2 = 10 % 3   // 1  (10 = 3 × 3 + 1)
let rem3 = 6 % 3    // 0  (6 = 2 × 3 + 0)
```

Use cases for `%`:

```astra
// Check even/odd
if n % 2 == 0 { print("even") }

// Wrap around (circular buffer index)
let index = (current + 1) % buffer_size

// Extract last N digits
let last_two = number % 100     // e.g., 12345 % 100 = 45
```

### Operator Precedence and the Classic Trap

```astra
let x = 2 + 3 * 4   // Is this (2+3)*4 = 20? Or 2+(3*4) = 14?
```

The answer is `14`, because `*` has **higher precedence** than `+`. Multiplication is evaluated before addition — this is the same rule you learned in school as "PEMDAS" or "BODMAS."

```
2 + 3 * 4
    ↑ evaluate first (higher precedence)
= 2 + 12
= 14
```

### Integer Overflow

What happens when you exceed the maximum value an integer can hold?

```
int64 maximum = 9,223,372,036,854,775,807
int64 maximum + 1 = -9,223,372,036,854,775,808  (wraps around!)
```

In Go (and by default in Astra), integer arithmetic wraps around silently — this is called **two's complement overflow**. Astra's standard library will provide checked arithmetic functions for safety-critical code.

---

## 3. Comparison Operators — Asking Questions

Comparison operators take two values and produce a **bool** (`true` or `false`).

| Operator | Meaning                | Example          | Result |
|----------|------------------------|------------------|--------|
| `==`     | equal to               | `5 == 5`         | true   |
| `!=`     | not equal to           | `5 != 3`         | true   |
| `<`      | less than              | `3 < 7`          | true   |
| `>`      | greater than           | `7 > 10`         | false  |
| `<=`     | less than or equal     | `5 <= 5`         | true   |
| `>=`     | greater than or equal  | `4 >= 9`         | false  |

### Chaining Comparisons — Why `1 < x < 10` Fails

In mathematics you write `1 < x < 10` and everyone understands it means "x is between 1 and 10." In most programming languages, including Astra, this **does not work** the way you expect:

```astra
// Astra (and most languages)
1 < x < 10
// Is parsed as: (1 < x) < 10
// (1 < x) produces a bool: true or false
// bool < 10 is a TYPE ERROR in Astra
```

The correct way in Astra:

```astra
1 < x && x < 10    // Use logical AND to chain comparisons
```

### Float Comparison — The IEEE 754 Surprise

This is one of the most famous bugs in all of programming:

```astra
let result = 0.1 + 0.2 == 0.3   // FALSE!
```

Why? Because computers store floating-point numbers in **IEEE 754 binary format**, and `0.1` and `0.2` cannot be represented exactly in binary. The actual stored values are:

```
0.1 in float64 = 0.1000000000000000055511151231257827021181583404541015625
0.2 in float64 = 0.200000000000000011102230246251565404236316680908203125
0.1 + 0.2      = 0.3000000000000000444089209850062616169452667236328125
0.3 in float64 = 0.299999999999999988897769753748434595763683319091796875
```

They differ at the 17th decimal place — but `==` requires exact bit-level equality.

The correct way to compare floats:

```astra
fn nearly_equal(a: float, b: float, epsilon: float) -> bool {
    return (a - b).abs() < epsilon
}

let ok = nearly_equal(0.1 + 0.2, 0.3, 0.0000001)   // true
```

---

## 4. Logical Operators — Combining Truth

Logical operators work on **bool** values and produce **bool** values. They let you combine multiple conditions.

### AND `&&` — Both Must Be True

```astra
true  && true   // true
true  && false  // false
false && true   // false
false && false  // false
```

Real use:

```astra
if age >= 18 && has_id {
    print("Welcome!")
}
```

### OR `||` — At Least One Must Be True

```astra
true  || true   // true
true  || false  // true
false || true   // true
false || false  // false
```

Real use:

```astra
if is_admin || is_moderator {
    show_dashboard()
}
```

### NOT `!` — Invert

```astra
!true   // false
!false  // true
!(x > 5)  // true when x <= 5
```

### Truth Tables

```
┌───────┬───────┬───────────┬──────────┐
│   A   │   B   │  A && B   │  A || B  │
├───────┼───────┼───────────┼──────────┤
│ false │ false │   false   │  false   │
│ false │ true  │   false   │  true    │
│ true  │ false │   false   │  true    │
│ true  │ true  │   true    │  true    │
└───────┴───────┴───────────┴──────────┘
```

### Short-Circuit Evaluation

This is a critical feature that affects program correctness and performance:

**For `&&`:** If the left side is `false`, the right side is **never evaluated** — because no matter what the right side is, the result is already `false`.

**For `||`:** If the left side is `true`, the right side is **never evaluated** — the result is already `true`.

```astra
// Safe null-check using short-circuit:
if ptr != null && ptr.value > 10 {
    // Without short-circuit, ptr.value would crash if ptr == null
    // With short-circuit, we never reach ptr.value when ptr == null
}
```

In Go, the same principle applies:

```go
// Safe: index check first, then access
if i < len(arr) && arr[i] > 0 {
    // arr[i] is only accessed when i is a valid index
}
```

Short-circuit also enables side-effect control:

```astra
if expensive_check() || cheap_check() {
    // If expensive_check() is true, cheap_check() never runs
    // Better: put cheap_check() first!
}
if cheap_check() || expensive_check() {
    // cheap_check() runs first; expensive_check() only runs if needed
}
```

### De Morgan's Laws

These algebraic identities let you simplify negated conditions:

```
!(A && B) == !A || !B       // NOT(A AND B) = (NOT A) OR (NOT B)
!(A || B) == !A && !B       // NOT(A OR B)  = (NOT A) AND (NOT B)
```

Example:

```astra
// Instead of:
if !(is_admin && is_logged_in) { deny() }
// You can write (sometimes clearer):
if !is_admin || !is_logged_in { deny() }
```

---

## 5. Bitwise Operators — Working at the Bit Level

Bitwise operators treat integers as sequences of bits and operate on each bit individually. They are essential for systems programming, networking, and performance-critical code.

Let's use 8-bit examples for clarity. In Astra, integers are 64-bit, but the principle is identical.

### Bitwise AND `&` — Masking

```
A:       1010 1100   (0xAC = 172)
B:       0000 1111   (0x0F = 15, a "mask")
A & B:   0000 1100   (0x0C = 12)
```

```
┌─────────────────────────────────────┐
│  BITWISE AND: output 1 only when    │
│  BOTH inputs are 1                  │
│                                     │
│  Bit: 7 6 5 4 3 2 1 0              │
│  A:   1 0 1 0 1 1 0 0              │
│  B:   0 0 0 0 1 1 1 1              │
│       ─────────────────             │
│  A&B: 0 0 0 0 1 1 0 0              │
└─────────────────────────────────────┘
```

Use case — extract the lower nibble (4 bits):

```astra
let byte_val = 0b10101100
let low_nibble = byte_val & 0b00001111   // = 0b00001100 = 12
```

Use case — Unix file permissions:

```
Permission bits: rwxrwxrwx = 111 111 111 = 0o777
User read bit:   0b100000000 = 0o400
Check user read: perms & 0o400 != 0
```

### Bitwise OR `|` — Setting Bits

```
A:       1010 0000   (0xA0 = 160)
B:       0000 0101   (0x05 = 5)
A | B:   1010 0101   (0xA5 = 165)
```

Use case — combine flags:

```astra
let READ    = 0b001   // 1
let WRITE   = 0b010   // 2
let EXECUTE = 0b100   // 4

let permissions = READ | WRITE   // 0b011 = 3 (read+write but not execute)
```

### Bitwise XOR `^` — Toggling

XOR (exclusive OR) outputs 1 when the inputs are **different**:

```
A:       1010 1010
B:       1111 0000
A ^ B:   0101 1010
```

Famous XOR swap (no temporary variable needed):

```go
a := 5   // 0101
b := 3   // 0011
a = a ^ b   // a = 0110 = 6
b = a ^ b   // b = 0101 = 5 (original a)
a = a ^ b   // a = 0011 = 3 (original b)
// a and b are now swapped!
```

### Left Shift `<<` — Multiply by Powers of 2

```
5 << 1 = 10    (5 × 2¹)
5 << 2 = 20    (5 × 2²)
5 << 3 = 40    (5 × 2³)
1 << 10 = 1024 (2^10)
```

```
┌────────────────────────────────────┐
│  LEFT SHIFT: push bits left,       │
│  fill right with zeros             │
│                                    │
│  0000 0101  (5)                    │
│  << 2                              │
│  0001 0100  (20)                   │
│  ↑ bits moved left by 2            │
└────────────────────────────────────┘
```

Use case — compute `2^n` efficiently:

```astra
let power_of_two = 1 << n   // faster than calling pow(2, n)
```

### Right Shift `>>` — Divide by Powers of 2

```
40 >> 1 = 20   (40 ÷ 2)
40 >> 2 = 10   (40 ÷ 4)
40 >> 3 = 5    (40 ÷ 8)
```

Use case — extract a specific byte from a 64-bit integer:

```astra
let ip: int = 0xC0A80101   // 192.168.1.1 packed into 32 bits
let octet1 = (ip >> 24) & 0xFF   // 192
let octet2 = (ip >> 16) & 0xFF   // 168
let octet3 = (ip >> 8)  & 0xFF   // 1
let octet4 =  ip        & 0xFF   // 1
```

---

## 6. Operator Precedence — Who Goes First?

Operator precedence defines which operations happen before others when you have multiple operators in one expression — just like in mathematics where multiplication happens before addition.

### Astra's Full Precedence Table

```
┌─────────────────────────────────────────────────────────┐
│            ASTRA OPERATOR PRECEDENCE TABLE               │
│         (Higher level = evaluated FIRST)                 │
├───────┬──────────────────┬───────────────────────────────┤
│ Level │    Operators     │         Description            │
├───────┼──────────────────┼───────────────────────────────┤
│  10   │  ()  []  .       │ Call, index, field access      │
│   9   │  !  - (unary)   │ Logical NOT, unary minus       │
│   8   │  *  /  %         │ Multiplication, division, mod  │
│   7   │  +  - (binary)  │ Addition, subtraction          │
│   6   │  <  >  <=  >=   │ Relational comparison          │
│   5   │  ==  !=          │ Equality comparison            │
│   4   │  &&              │ Logical AND                    │
│   3   │  ||              │ Logical OR                     │
│   2   │  =               │ Assignment                     │
│   1   │  ,               │ Comma (function arguments)     │
└───────┴──────────────────┴───────────────────────────────┘
```

### Five Step-by-Step Examples

**Example 1:** `2 + 3 * 4`

```
2 + 3 * 4
    ↑ * is level 8, + is level 7 → * first
= 2 + 12
= 14
```

**Example 2:** `!x > 5`

```
!x > 5
↑ ! is level 9 (unary), > is level 6 → ! first
= (!x) > 5
```

If `x` is `int`, this is a type error because `!` expects a bool. The precedence works, but the types don't. This is a common bug — the programmer meant `!(x > 5)`.

**Example 3:** `a == b && c != d`

```
a == b && c != d
  ↑      ↑  == and != are level 5; && is level 4
= (a == b) && (c != d)
```

**Example 4:** `x > 0 || y > 0 && z > 0`

```
x > 0 || y > 0 && z > 0
          ↑ && is level 4, || is level 3 → && first
= x > 0 || (y > 0 && z > 0)
```

This evaluates as: "x is positive OR (both y AND z are positive)." This is often NOT what the programmer intended. Use parentheses to be explicit:

```astra
(x > 0 || y > 0) && z > 0   // different meaning!
```

**Example 5:** Assignment right-associativity

```astra
let a = 0
let b = 0
a = b = 5   // right-associative: a = (b = 5)
// First: b = 5 (b becomes 5)
// Then:  a = 5 (a becomes the value of the assignment, which is 5)
```

---

## 7. Expression Trees — The AST Preview

When the Astra parser reads `2 + 3 * 4`, it does not store it as a flat string. It builds a **tree** structure that captures the precedence and nesting explicitly. This tree is called an **Abstract Syntax Tree (AST)**.

Think of it like a family tree: the root is at the top, children hang below. For expressions, the **operator** is the parent, and its **operands** are the children.

### Tree for `2 + 3 * 4`

```
           BinaryExpr
           Operator: +
          /             \
   IntLiteral        BinaryExpr
    Value: 2         Operator: *
                    /           \
             IntLiteral      IntLiteral
              Value: 3        Value: 4
```

To **evaluate** this tree, you do a **post-order traversal** (children before parent):

```
1. Evaluate left child of +:  IntLiteral(2) → 2
2. Evaluate right child of +:
   a. Evaluate left child of *:  IntLiteral(3) → 3
   b. Evaluate right child of *: IntLiteral(4) → 4
   c. Apply *: 3 * 4 → 12
3. Apply +: 2 + 12 → 14
```

### Left-Associativity

For operators of the same precedence, left-associativity means we group from the left:

```astra
5 - 3 - 1
// Left-associative: (5 - 3) - 1 = 1
// NOT: 5 - (3 - 1) = 3   ← this would be right-associative
```

Tree for `5 - 3 - 1` (left-associative):

```
         BinaryExpr(-)
        /              \
  BinaryExpr(-)    IntLiteral(1)
  /         \
IntLit(5)  IntLit(3)
```

### Right-Associativity

Assignment in most languages is right-associative:

```astra
a = b = c = 5
// Parsed as: a = (b = (c = 5))
```

Tree:

```
BinaryExpr(=)
├── Identifier(a)
└── BinaryExpr(=)
    ├── Identifier(b)
    └── BinaryExpr(=)
        ├── Identifier(c)
        └── IntLiteral(5)
```

---

## 8. String Operations in Astra

Strings in Astra support a limited set of operators:

### Concatenation with `+`

```astra
let greeting = "Hello" + ", " + "World"    // "Hello, World"
let repeated = "ha" + "ha" + "ha"          // "hahaha"
```

Under the hood, string concatenation allocates a new string containing all the bytes of both operands. This is an `O(n)` operation where `n` is the total length of the result.

### Equality with `==` and `!=`

```astra
"hello" == "hello"   // true  (compares content, not pointer)
"hello" != "world"   // true
"Hello" == "hello"   // false (case-sensitive)
```

Unlike some languages (Java), Astra's `==` compares **string content** (character by character), not whether two variables point to the same object in memory. This is what users expect.

Internally, Astra uses **string interning**: identical string literals are stored once and shared. So `"hello" == "hello"` may be a single pointer comparison at runtime, but the semantics guarantee content equality.

### What Strings Cannot Do

```astra
"hello" + 5       // ERROR: cannot add string and int
"hello" < "world" // Astra v1: not supported (use .compare() method)
"hello" * 3       // ERROR: not supported
```

---

## 9. 🔨 Astra Build Milestone — Expression AST Nodes

Now we define the Go data structures that represent every kind of expression in Astra's AST. This file will be imported by the parser, the type checker, and the code generator.

```go
// ast/expressions.go
package ast

// ============================================================
// Core Interfaces
// ============================================================

// Node is the base interface that every AST node implements.
// Every node knows its source position for error reporting.
type Node interface {
    TokenLiteral() string  // the literal text of the first token
    String() string        // human-readable representation (for debug)
    GetPos() Position      // location in source code
}

// Position tracks where in the source file a node came from.
// Essential for meaningful error messages: "error on line 42, column 7"
type Position struct {
    Line   int    // 1-based line number
    Column int    // 1-based column number
    File   string // source file name
}

// Expression is the interface all expression nodes implement.
// The expressionNode() method is a marker — it distinguishes
// expression nodes from statement nodes at the type level.
type Expression interface {
    Node
    expressionNode()   // marker method — Go interfaces use this trick
    GetType() string   // returns the resolved type: "int", "float", "bool", "string"
}

// ============================================================
// Literal Expressions — the leaf nodes of the expression tree
// ============================================================

// IntLiteral represents an integer constant like 42 or -7
type IntLiteral struct {
    Value int64    // the actual integer value
    Pos   Position // where it appeared in source
    Type  string   // always "int"
}

func (il *IntLiteral) expressionNode()      {}
func (il *IntLiteral) TokenLiteral() string { return fmt.Sprintf("%d", il.Value) }
func (il *IntLiteral) String() string       { return fmt.Sprintf("IntLit(%d)", il.Value) }
func (il *IntLiteral) GetPos() Position     { return il.Pos }
func (il *IntLiteral) GetType() string      { return "int" }

// FloatLiteral represents a floating-point constant like 3.14 or -0.5
type FloatLiteral struct {
    Value float64
    Pos   Position
    Type  string   // always "float"
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fmt.Sprintf("%g", fl.Value) }
func (fl *FloatLiteral) String() string       { return fmt.Sprintf("FloatLit(%g)", fl.Value) }
func (fl *FloatLiteral) GetPos() Position     { return fl.Pos }
func (fl *FloatLiteral) GetType() string      { return "float" }

// StringLiteral represents a string constant like "hello"
// The Value field holds the unescaped content (backslash sequences resolved)
type StringLiteral struct {
    Value string   // e.g. for source "hello\n", Value = "hello\n" (actual newline)
    Pos   Position
    Type  string   // always "string"
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return `"` + sl.Value + `"` }
func (sl *StringLiteral) String() string       { return fmt.Sprintf("StrLit(%q)", sl.Value) }
func (sl *StringLiteral) GetPos() Position     { return sl.Pos }
func (sl *StringLiteral) GetType() string      { return "string" }

// BoolLiteral represents true or false
type BoolLiteral struct {
    Value bool
    Pos   Position
    Type  string   // always "bool"
}

func (bl *BoolLiteral) expressionNode()      {}
func (bl *BoolLiteral) TokenLiteral() string {
    if bl.Value { return "true" }
    return "false"
}
func (bl *BoolLiteral) String() string   { return fmt.Sprintf("BoolLit(%v)", bl.Value) }
func (bl *BoolLiteral) GetPos() Position { return bl.Pos }
func (bl *BoolLiteral) GetType() string  { return "bool" }

// ============================================================
// Compound Expressions — internal nodes of the expression tree
// ============================================================

// BinaryExpr represents a two-operand operation: Left OPERATOR Right
// Examples: 2 + 3, x == y, a && b, "hi" + "!"
type BinaryExpr struct {
    Left     Expression // the left operand (any expression)
    Operator string     // the operator string: "+", "-", "*", "/", "%",
                        // "==", "!=", "<", ">", "<=", ">=", "&&", "||",
                        // "&", "|", "^", "<<", ">>"
    Right    Expression // the right operand (any expression)
    Pos      Position   // position of the operator token
    Type     string     // resolved type of the whole expression
}

func (be *BinaryExpr) expressionNode()      {}
func (be *BinaryExpr) TokenLiteral() string { return be.Operator }
func (be *BinaryExpr) String() string {
    return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}
func (be *BinaryExpr) GetPos() Position { return be.Pos }
func (be *BinaryExpr) GetType() string  { return be.Type }

// UnaryExpr represents a single-operand operation: OPERATOR Operand
// Examples: -x, !flag, -(a + b)
type UnaryExpr struct {
    Operator string     // "!" or "-"
    Operand  Expression // the operand
    Pos      Position   // position of the operator token
    Type     string     // resolved type
}

func (ue *UnaryExpr) expressionNode()      {}
func (ue *UnaryExpr) TokenLiteral() string { return ue.Operator }
func (ue *UnaryExpr) String() string {
    return fmt.Sprintf("(%s%s)", ue.Operator, ue.Operand.String())
}
func (ue *UnaryExpr) GetPos() Position { return ue.Pos }
func (ue *UnaryExpr) GetType() string  { return ue.Type }

// Identifier represents a variable name or function name
// Examples: x, my_variable, add
type Identifier struct {
    Name string
    Pos  Position
    Type string   // resolved by the type checker
}

func (id *Identifier) expressionNode()      {}
func (id *Identifier) TokenLiteral() string { return id.Name }
func (id *Identifier) String() string       { return fmt.Sprintf("Ident(%s)", id.Name) }
func (id *Identifier) GetPos() Position     { return id.Pos }
func (id *Identifier) GetType() string      { return id.Type }

// GroupExpr represents a parenthesized expression: (expr)
// The parentheses do not change semantics, only precedence.
// We keep this node to preserve source fidelity for error messages.
type GroupExpr struct {
    Inner Expression
    Pos   Position
}

func (ge *GroupExpr) expressionNode()      {}
func (ge *GroupExpr) TokenLiteral() string { return "(" }
func (ge *GroupExpr) String() string       { return fmt.Sprintf("(%s)", ge.Inner.String()) }
func (ge *GroupExpr) GetPos() Position     { return ge.Pos }
func (ge *GroupExpr) GetType() string      { return ge.Inner.GetType() }
```

### A Simple Constant-Expression Evaluator

To prove the AST works, here is a complete evaluator that walks an expression tree and computes integer results for constant expressions. This is exactly what a compiler does during **constant folding** — computing values at compile time that would otherwise be computed at runtime.

```go
// ast/eval.go
package ast

import "fmt"

// EvalConstInt evaluates a constant integer expression at compile time.
// Returns (value, true) on success, (0, false) if evaluation fails
// (e.g., the expression contains variables or is not integer-typed).
func EvalConstInt(expr Expression) (int64, bool) {
    switch e := expr.(type) {

    case *IntLiteral:
        // Base case: a literal integer — just return its value
        return e.Value, true

    case *GroupExpr:
        // Parentheses: evaluate what's inside
        return EvalConstInt(e.Inner)

    case *UnaryExpr:
        val, ok := EvalConstInt(e.Operand)
        if !ok {
            return 0, false
        }
        switch e.Operator {
        case "-":
            return -val, true
        default:
            return 0, false  // "!" doesn't apply to int
        }

    case *BinaryExpr:
        // Evaluate both sides first (post-order traversal)
        left, ok1 := EvalConstInt(e.Left)
        right, ok2 := EvalConstInt(e.Right)
        if !ok1 || !ok2 {
            return 0, false  // can't fold if either side is non-constant
        }

        // Apply the operator
        switch e.Operator {
        case "+":
            return left + right, true
        case "-":
            return left - right, true
        case "*":
            return left * right, true
        case "/":
            if right == 0 {
                // Compile-time division by zero!
                panic(fmt.Sprintf("compile error: division by zero at %v", e.Pos))
            }
            return left / right, true
        case "%":
            if right == 0 {
                panic(fmt.Sprintf("compile error: modulo by zero at %v", e.Pos))
            }
            return left % right, true
        case "&":
            return left & right, true
        case "|":
            return left | right, true
        case "^":
            return left ^ right, true
        case "<<":
            return left << uint(right), true
        case ">>":
            return left >> uint(right), true
        default:
            return 0, false  // comparison operators return bool, not int
        }

    default:
        // Variables, function calls, etc. — cannot evaluate at compile time
        return 0, false
    }
}

// Example: building the AST for "2 + 3 * 4" by hand and evaluating it
func ExampleEval() {
    // Build the AST:
    //        BinaryExpr(+)
    //       /              \
    //  IntLit(2)        BinaryExpr(*)
    //                   /           \
    //               IntLit(3)    IntLit(4)

    mul := &BinaryExpr{
        Left:     &IntLiteral{Value: 3},
        Operator: "*",
        Right:    &IntLiteral{Value: 4},
    }
    add := &BinaryExpr{
        Left:     &IntLiteral{Value: 2},
        Operator: "+",
        Right:    mul,
    }

    result, ok := EvalConstInt(add)
    if ok {
        fmt.Printf("2 + 3 * 4 = %d\n", result)  // Output: 2 + 3 * 4 = 14
    }
}
```

### Common Mistake: Forgetting Parentheses in Type Assertions

When walking the AST in Go, always use the two-value form of type assertion to avoid panics:

```go
// WRONG — panics if expr is not *IntLiteral:
il := expr.(*IntLiteral)

// CORRECT — safe type assertion:
il, ok := expr.(*IntLiteral)
if !ok {
    // handle the case where expr is a different expression type
}

// OR — use a type switch (best for multiple types):
switch e := expr.(type) {
case *IntLiteral:
    fmt.Println("int:", e.Value)
case *BinaryExpr:
    fmt.Println("binary op:", e.Operator)
default:
    fmt.Println("unknown expression type")
}
```

---

## Exercises

1. **Precedence puzzle.** Evaluate each expression without running any code. Then verify by tracing through the precedence table:
   - `3 + 4 * 2 - 1`
   - `10 % 3 + 1`
   - `!false || true && false`
   - `2 << 3 + 1`
   *Hint: treat `<<` and `+` separately — which has higher precedence?*

2. **Integer division trap.** Write an Astra function `average(a: int, b: int) -> float` that correctly returns the average as a floating-point number. What would the buggy version return for `average(5, 6)`? What does the correct version return?
   *Hint: you must convert to float before dividing.*

3. **Bit manipulation.** Given an integer `n`, write Astra code to:
   a. Check if the 3rd bit (bit index 2) is set.
   b. Set the 3rd bit.
   c. Clear the 3rd bit.
   d. Toggle the 3rd bit.
   *Hint: use `1 << 2` as your mask.*

4. **Build an AST by hand.** Draw the expression tree for `!(x > 0 && y < 10)`. Then simplify it using De Morgan's law and draw the simplified tree. Are both trees semantically equivalent?
   *Hint: Apply `!(a && b) == !a || !b`.*

5. **Extend the evaluator.** Add a `EvalConstBool(expr Expression) (bool, bool)` function to `ast/eval.go` that handles comparison operators (`==`, `!=`, `<`, `>`, `<=`, `>=`) and logical operators (`&&`, `||`, `!`). Test it on `2 + 3 > 4 && !false`.
   *Hint: For comparisons between ints, call `EvalConstInt` on both sides first.*

6. **Short-circuit test.** Write a Go function `safeDiv(a, b int) int` that divides `a` by `b` but returns 0 if `b` is 0. Use the `&&` short-circuit behavior in the condition. Then rewrite it to demonstrate what would happen WITHOUT short-circuit evaluation.
   *Hint: `if b != 0 && a/b > 0 { ... }` is safe; `if a/b > 0 && b != 0 { ... }` panics when b=0.*

---

## Summary

| Concept | Definition | Astra Example |
|---|---|---|
| Expression | Evaluates to a value; has a type | `2 + 3 * 4` → `int` |
| Statement | Performs an action; no value | `let x = 5` |
| Integer division | Truncates toward zero | `7 / 2 == 3` |
| Modulo | Remainder after division | `7 % 2 == 1` |
| Operator precedence | Order of evaluation | `*` before `+` |
| Short-circuit eval | Skip right side if unnecessary | `false && f()` → `f()` never called |
| Bitwise operators | Operate on individual bits | `0b1010 & 0b1100 == 0b1000` |
| AST / expression tree | Tree structure of an expression | `BinaryExpr{+, IntLit(2), BinaryExpr{*,...}}` |
| Constant folding | Evaluate constant exprs at compile time | `2+3` → compiler emits `5` directly |
| IEEE 754 | Floating-point number standard | `0.1 + 0.2 != 0.3` |
| De Morgan's law | Simplify negated boolean expressions | `!(a&&b) == !a\|\|!b` |
| Left-associativity | Group same-precedence ops left to right | `5-3-1 == (5-3)-1 == 1` |
| Type assertion | Go way to extract concrete type from interface | `e.(*IntLiteral)` |

The AST nodes you built here — `IntLiteral`, `FloatLiteral`, `StringLiteral`, `BoolLiteral`, `BinaryExpr`, `UnaryExpr`, `Identifier`, `GroupExpr` — will be produced by the Astra parser in Chapter 55 and consumed by the type checker and code generator in later chapters. Every operator, every precedence rule, every evaluation strategy you learned here feeds directly into the compiler's core.
