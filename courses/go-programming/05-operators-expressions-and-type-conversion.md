# Chapter 05: Operators, Expressions, and Type Conversion

You've learned about types and how to declare variables. Now let's look at how to work with them — how to add, compare, combine, and manipulate values using operators. Go's operators are similar to C and Java, but with a few important differences (no `++i`, no ternary operator, strict type rules). Understanding them deeply prevents subtle bugs.

## Table of Contents

1. [Arithmetic Operators](#1-arithmetic-operators)
2. [Comparison Operators](#2-comparison-operators)
3. [Logical Operators](#3-logical-operators)
4. [Bitwise Operators](#4-bitwise-operators)
5. [Assignment Operators](#5-assignment-operators)
6. [Operator Precedence](#6-operator-precedence)
7. [Type Conversion in Depth](#7-type-conversion-in-depth)
8. [The fmt Package — Printing and Formatting](#8-the-fmt-package--printing-and-formatting)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Arithmetic Operators

```go
a, b := 10, 3

fmt.Println(a + b)    // 13 — addition
fmt.Println(a - b)    //  7 — subtraction
fmt.Println(a * b)    // 30 — multiplication
fmt.Println(a / b)    //  3 — integer division (truncates!)
fmt.Println(a % b)    //  1 — modulo (remainder)

// Float division:
x, y := 10.0, 3.0
fmt.Println(x / y)    // 3.3333333333333335

// Integer division truncates (does NOT round):
fmt.Println(7 / 2)    // 3 (not 3.5)
fmt.Println(-7 / 2)   // -3 (truncates toward zero)
fmt.Println(7 % 2)    // 1
fmt.Println(-7 % 2)   // -1 (same sign as dividend)
```

**Increment and decrement — Go's version:**
```go
i := 5
i++   // i = 6 (postfix only — there is no ++i in Go)
i--   // i = 5

// Increment is a STATEMENT, not an expression
// You cannot do: j := i++  (doesn't exist in Go)
// You cannot do: fmt.Println(i++)  (compile error)
```

**Overflow is silent — be careful:**
```go
var x int8 = 127
x++
fmt.Println(x)  // -128 (overflow! wraps around)

var u uint8 = 0
u--
fmt.Println(u)  // 255 (underflow! wraps around)
```

**Division by zero:**
```go
// A constant expression like 10 / 0 won't even compile:
// var result = 10 / 0  // compile error: division by zero

// Integer division by a zero VARIABLE panics at runtime:
a, b := 10, 0
_ = a / b  // panic: runtime error: integer divide by zero

// Float division by a zero variable produces Inf or NaN:
zero := 0.0
f := 1.0 / zero
fmt.Println(f)  // +Inf

import "math"
fmt.Println(math.IsInf(f, 1))  // true
fmt.Println(zero / zero)       // NaN
```

### Quick Check
> 1. What does `7 / 2` equal in Go and why?
> 2. What is wrong with `j := i++` in Go?
> 3. What happens when you divide an integer by zero?

---

## 2. Comparison Operators

Comparison operators return `bool`:

```go
a, b := 5, 10

fmt.Println(a == b)  // false — equal
fmt.Println(a != b)  // true  — not equal
fmt.Println(a < b)   // true  — less than
fmt.Println(a > b)   // false — greater than
fmt.Println(a <= b)  // true  — less than or equal
fmt.Println(a >= b)  // false — greater than or equal
```

**String comparison** — lexicographic (dictionary order):
```go
fmt.Println("apple" < "banana")   // true (a < b)
fmt.Println("apple" == "apple")   // true
fmt.Println("Apple" < "apple")    // true (uppercase letters come before lowercase in ASCII)
fmt.Println("abc" < "abd")        // true (compare char by char: 'c' < 'd')
```

**Comparing structs** — only works if all fields are comparable:
```go
type Point struct {
    X, Y int
}

p1 := Point{1, 2}
p2 := Point{1, 2}
p3 := Point{3, 4}

fmt.Println(p1 == p2)  // true (all fields equal)
fmt.Println(p1 == p3)  // false

// Slices and maps are NOT comparable with ==
// s1 := []int{1, 2, 3}
// s2 := []int{1, 2, 3}
// fmt.Println(s1 == s2)  // compile error: slice can only be compared to nil
```

**Comparing with nil:**
```go
var p *int = nil
var s []int = nil
var m map[string]int = nil

fmt.Println(p == nil)  // true
fmt.Println(s == nil)  // true
fmt.Println(m == nil)  // true
fmt.Println(s == nil)  // true (nil slice)

s = []int{}  // empty but NOT nil
fmt.Println(s == nil)  // false!
```

### Quick Check
> 1. Can you compare two slices with `==` in Go?
> 2. What is the difference between a nil slice and an empty slice?
> 3. What does string comparison in Go use? (alphabetical, lexicographic?)

---

## 3. Logical Operators

```go
t, f := true, false

fmt.Println(t && f)   // false — AND (both must be true)
fmt.Println(t || f)   // true  — OR (at least one must be true)
fmt.Println(!t)       // false — NOT
```

**Short-circuit evaluation** — Go stops evaluating as soon as the result is known:
```go
// && short-circuits: if left is false, right is never evaluated
false && expensiveFunction()  // expensiveFunction() is NOT called

// || short-circuits: if left is true, right is never evaluated
true || expensiveFunction()   // expensiveFunction() is NOT called
```

**Practical use of short-circuiting:**
```go
// Safe nil pointer check:
if user != nil && user.IsActive() {
    // user.IsActive() only called if user is not nil
    doSomething()
}

// Safe slice access:
if len(items) > 0 && items[0] == target {
    // items[0] only accessed if len > 0 (avoids index out of bounds)
}
```

**No ternary operator in Go:**
```go
// In many languages: x = (a > b) ? a : b
// Go does NOT have ?: (ternary)

// Go equivalent:
var max int
if a > b {
    max = a
} else {
    max = b
}

// Or using a function (idiomatic):
func maxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}
max := maxInt(a, b)
```

### Quick Check
> 1. What is short-circuit evaluation and how does Go use it?
> 2. Does Go have a ternary operator (`? :`)?
> 3. Why is `user != nil && user.IsActive()` safe but `user.IsActive() && user != nil` is not?

---

## 4. Bitwise Operators

Bitwise operators work on individual bits of integers. Essential for flags, masks, and low-level operations:

```go
a := 0b1010  // 10 in binary
b := 0b1100  // 12 in binary

fmt.Printf("%04b\n", a & b)   // 1000 = 8  — AND: bit is 1 only if BOTH are 1
fmt.Printf("%04b\n", a | b)   // 1110 = 14 — OR: bit is 1 if EITHER is 1
fmt.Printf("%04b\n", a ^ b)   // 0110 = 6  — XOR: bit is 1 if EXACTLY ONE is 1
fmt.Printf("%04b\n", a &^ b)  // 0010 = 2  — AND NOT (bit clear): a AND (NOT b)
```

**Bit shifting:**
```go
x := 1

fmt.Println(x << 3)  // 8  — left shift: multiply by 2^3
fmt.Println(x >> 1)  // 0  — right shift: divide by 2^1

// Powers of 2 using shifts:
fmt.Println(1 << 10)  // 1024  (KB)
fmt.Println(1 << 20)  // 1048576 (MB)
fmt.Println(1 << 30)  // 1073741824 (GB)
```

**Practical bitwise uses:**

*Checking if a number is even or odd:*
```go
func isEven(n int) bool {
    return n & 1 == 0  // Check lowest bit
}
```

*Setting, clearing, and toggling bits (permission flags):*
```go
const (
    FlagA = 1 << iota  // 0001
    FlagB              // 0010
    FlagC              // 0100
    FlagD              // 1000
)

flags := 0

// Set flag A and C
flags |= FlagA | FlagC  // flags = 0101 = 5

// Check if flag B is set
if flags & FlagB != 0 {
    fmt.Println("B is set")
}

// Clear flag A
flags &^= FlagA  // flags = 0100 = 4

// Toggle flag C
flags ^= FlagC   // flags = 0000 = 0
```

### Quick Check
> 1. What does `&^` (AND NOT) do in Go?
> 2. What is `1 << 10`?
> 3. How do you check if bit 3 (value 8) is set in an integer?

---

## 5. Assignment Operators

```go
x := 10

// Compound assignment operators (shorthand)
x += 5   // x = x + 5 → 15
x -= 3   // x = x - 3 → 12
x *= 2   // x = x * 2 → 24
x /= 4   // x = x / 4 → 6
x %= 4   // x = x % 4 → 2

// Bitwise compound assignments
x &= 0xFF   // x = x & 0xFF
x |= 0x01   // x = x | 0x01
x ^= 0x0F   // x = x ^ 0x0F
x <<= 2     // x = x << 2
x >>= 1     // x = x >> 1
```

**Multiple assignment:**
```go
x, y := 1, 2

// Swap without temporary variable (Go-idiomatic):
x, y = y, x
fmt.Println(x, y)  // 2 1
```

**Multiple return values assignment:**
```go
func minMax(arr []int) (int, int) {
    min, max := arr[0], arr[0]
    for _, v := range arr {
        if v < min { min = v }
        if v > max { max = v }
    }
    return min, max
}

min, max := minMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
fmt.Println(min, max)  // 1 9
```

### Quick Check
> 1. What does `x += 5` mean?
> 2. How do you swap two variables without a temporary variable in Go?
> 3. Can you assign to multiple variables in a single statement?

---

## 6. Operator Precedence

When an expression has multiple operators, precedence determines evaluation order:

```go
// Higher precedence is evaluated first (like math)
// 5 + 3 * 2 = 11 (not 16) because * has higher precedence than +
fmt.Println(5 + 3 * 2)  // 11

// Go operator precedence (high to low):
// 5  *  /  %  <<  >>  &  &^
// 4  +  -  |  ^
// 3  ==  !=  <  <=  >  >=
// 2  &&
// 1  ||
```

```go
// Examples:
fmt.Println(2 + 3 * 4)        // 14 (not 20): * before +
fmt.Println((2 + 3) * 4)      // 20: parentheses force evaluation order
fmt.Println(true || false && false)  // true: && before ||
// Equivalent to: true || (false && false) = true || false = true
```

**Rule of thumb**: When in doubt, use parentheses. Explicit is better than relying on precedence rules. Future readers of your code will thank you.

```go
// Hard to read without parentheses:
result := x & mask == 0

// Clear with parentheses:
result := (x & mask) == 0
```

### Quick Check
> 1. What does `2 + 3 * 4` evaluate to in Go?
> 2. Which has higher precedence: `&&` or `||`?
> 3. What's a good rule of thumb when you're unsure about operator precedence?

---

## 7. Type Conversion in Depth

Go's strict type system means you must be explicit about conversions. Let's see all the cases:

**Numeric conversions:**
```go
// Widening (safe, no data loss):
var i8 int8 = 50
i64 := int64(i8)    // 50 → fine

// Narrowing (potentially lossy):
var i64big int64 = 300
i8bad := int8(i64big)  // 300 doesn't fit in int8 — truncates! → 44 (overflow)
fmt.Println(i8bad)      // 44 (not 300!)

// Always check bounds when narrowing:
if i64big > math.MaxInt8 || i64big < math.MinInt8 {
    return fmt.Errorf("value %d out of int8 range", i64big)
}
```

**Float ↔ Integer:**
```go
f := 3.99
i := int(f)     // 3 (truncates, does NOT round!)
fmt.Println(i)  // 3

// To round: math.Round()
import "math"
rounded := int(math.Round(3.99))  // 4
```

**Integer ↔ String (watch out!):**
```go
// WRONG WAY: int → string via type conversion
n := 65
s := string(n)   // "A" (Unicode code point 65 = 'A') — NOT "65"!

// CORRECT: use strconv
import "strconv"
s2 := strconv.Itoa(n)    // "65" ← what you probably wanted
s3 := fmt.Sprintf("%d", n)  // "65"

// String → int:
n2, err := strconv.Atoi("42")     // 42, nil
n3, err := strconv.ParseInt("42", 10, 64)  // 42, nil
```

**Interface conversions:**
```go
// interface{} (or any) can hold any value
var val interface{} = 42

// Type assertion — extract the real type:
n, ok := val.(int)
if !ok {
    fmt.Println("not an int")
} else {
    fmt.Println(n + 1)  // 43
}

// Type switch — check multiple types:
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %q\n", v)
    case bool:
        fmt.Printf("bool: %v\n", v)
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}

describe(42)       // int: 42
describe("hello")  // string: "hello"
describe(3.14)     // unknown type: float64
```

### Quick Check
> 1. What does `int8(300)` return and why?
> 2. What does `int(3.99)` return?
> 3. Why is `string(65)` dangerous?

---

## 8. The fmt Package — Printing and Formatting

You've been using `fmt.Println` — let's understand the whole `fmt` package:

**Print functions:**
```go
fmt.Print("no newline")          // Prints without newline
fmt.Println("with newline")      // Prints with newline
fmt.Printf("formatted: %d\n", 42)  // Formatted print (C-style)
```

**Format verbs:**
```go
// General
fmt.Printf("%v", 42)          // Default format: 42
fmt.Printf("%T", 42)          // Type: int
fmt.Printf("%+v", struct{X int}{42})  // Struct with field names: {X:42}
fmt.Printf("%#v", 42)         // Go syntax: 42

// Integers
fmt.Printf("%d", 42)          // Decimal: 42
fmt.Printf("%b", 42)          // Binary: 101010
fmt.Printf("%o", 42)          // Octal: 52
fmt.Printf("%x", 255)         // Hex lowercase: ff
fmt.Printf("%X", 255)         // Hex uppercase: FF
fmt.Printf("%05d", 42)        // Zero-padded: 00042
fmt.Printf("%-10d|", 42)      // Left-aligned: 42        |

// Floats
fmt.Printf("%f", 3.14)        // Default: 3.140000
fmt.Printf("%.2f", 3.14159)   // 2 decimal places: 3.14
fmt.Printf("%e", 123456789.0) // Scientific: 1.234568e+08
fmt.Printf("%g", 3.14)        // Shorter of %e or %f: 3.14

// Strings
fmt.Printf("%s", "hello")     // Plain: hello
fmt.Printf("%q", "hello")     // Quoted: "hello"
fmt.Printf("%10s", "hello")   // Right-aligned (padded on left): "     hello"
fmt.Printf("%-10s|", "hello") // Left-aligned (padded on right): "hello     |"

// Booleans
fmt.Printf("%t", true)        // true

// Pointers
x := 42
fmt.Printf("%p", &x)          // Memory address: 0xc0000b4000
```

**Sprintf — format to string:**
```go
// Don't print, just format to string
msg := fmt.Sprintf("Hello, %s! You are %d years old.", "Alice", 30)
fmt.Println(msg)  // Hello, Alice! You are 30 years old.

// Very useful for building error messages:
err := fmt.Errorf("user %d not found in database %s", userID, dbName)
```

**Fprintf — format to a writer:**
```go
import "os"

// Write to stderr:
fmt.Fprintf(os.Stderr, "Error: %v\n", err)

// Write to a file:
f, _ := os.Create("output.txt")
fmt.Fprintf(f, "Hello, file!\n")
f.Close()
```

**Sscanf — parse formatted strings:**
```go
var name string
var age int
n, err := fmt.Sscanf("Alice 30", "%s %d", &name, &age)
fmt.Println(n, name, age, err)  // 2 Alice 30 <nil>
```

### Quick Check
> 1. What is the difference between `%v` and `%+v`?
> 2. What does `%.2f` do to a float?
> 3. When would you use `fmt.Sprintf` instead of `fmt.Printf`?

---

## Summary

- **Arithmetic**: `+`, `-`, `*`, `/` (truncating), `%`; `i++`/`i--` are statements not expressions
- **Comparison**: `==`, `!=`, `<`, `>`, `<=`, `>=`; slices/maps not comparable with `==`
- **Logical**: `&&` (AND), `||` (OR), `!` (NOT); short-circuit evaluation; no ternary operator
- **Bitwise**: `&` (AND), `|` (OR), `^` (XOR), `&^` (AND NOT), `<<` (left shift), `>>` (right shift)
- **Assignment**: `=`, `:=`, compound operators (`+=`, `*=`, etc.); multiple assignment with swap
- **Precedence**: `*/%` > `+-` > comparisons > `&&` > `||`; use parentheses when unsure
- **Type conversion**: Always explicit; narrowing can lose data; `string(n)` is NOT "n" as string
- **fmt package**: `Print`/`Println`/`Printf` for output; `Sprintf` for string formatting; rich format verbs

---

## Exercises

### Easy
1. What is the output of each: `10/3`, `10.0/3.0`, `10%3`, `-10%3`? Explain each.
2. Write a function `isPowerOf2(n int) bool` using only bitwise operators (no division or modulo).
3. Write a function `setBit(n, pos int) int`, `clearBit(n, pos int) int`, `toggleBit(n, pos int) int` that set, clear, and toggle bit at position `pos` in integer `n`.

### Medium
4. Formatting table: Write a function that takes a slice of `[]struct{Name string; Age int; Score float64}` and prints it as a neatly formatted table with aligned columns. Use `fmt.Sprintf` to build each row. The output should look like:
   ```
   Name            Age    Score
   ─────────────────────────────
   Alice            25    98.50
   Bob              30    85.20
   Charlie          22   100.00
   ```
5. Safe arithmetic library: Write a `safemath` package with: `Add(a, b int64) (int64, error)`, `Sub(a, b int64) (int64, error)`, `Mul(a, b int64) (int64, error)`, `Div(a, b int64) (int64, error)`. All must return errors for overflow/underflow/division by zero. Test all edge cases (MaxInt64 + 1, MinInt64 - 1, 0/0, etc.).
6. Bit manipulation cipher: Implement a simple XOR cipher: `Encrypt(data []byte, key byte) []byte` XORs each byte with the key. `Decrypt(data []byte, key byte) []byte` is identical (XOR is its own inverse). Test that `Decrypt(Encrypt(data, key), key) == data`. Also implement: `HammingDistance(a, b []byte) int` (number of bits that differ between two byte slices).

### Hard
7. Expression evaluator: Build a simple mathematical expression evaluator that handles: integers, `+`, `-`, `*`, `/`, `%`, parentheses, and proper operator precedence. Input: `"3 + 4 * 2"` → 11, `"(3 + 4) * 2"` → 14. Implement using a recursive descent parser (parse expression → parse term → parse factor pattern). Handle errors: division by zero, mismatched parentheses, invalid tokens.
8. Type system exploration: Go's type system is powerful but sometimes surprising. Investigate and explain with code examples: (a) The difference between `type MyInt int` and `type MyInt = int` (type definition vs type alias). (b) Why `type Celsius float64` and `type Fahrenheit float64` can't be directly compared with `==`. (c) What is "named return values" and when can they cause bugs? (d) What is "shadowing" in Go and write an example where it causes a subtle bug. (e) How does Go handle numeric overflow — write a function that detects overflow for all four arithmetic operations.
