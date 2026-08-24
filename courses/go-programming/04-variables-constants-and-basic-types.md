# Chapter 04: Variables, Constants, and Basic Types

Every program stores and manipulates data. In Go, every piece of data has a **type** — it tells the compiler what kind of data it is (a number? text? true/false?), how much memory to allocate, and what operations are valid on it. This chapter covers Go's type system, how to declare variables, and how to work with constants. These are the atoms from which all Go programs are built.

## Table of Contents

1. [Variables — Declaring and Using](#1-variables--declaring-and-using)
2. [Basic Types](#2-basic-types)
3. [Zero Values](#3-zero-values)
4. [Type Conversions](#4-type-conversions)
5. [Constants](#5-constants)
6. [iota — Enumerations in Go](#6-iota--enumerations-in-go)
7. [Strings in Depth](#7-strings-in-depth)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Variables — Declaring and Using

Go has four ways to declare variables. Each has its place:

**Method 1: `var` with explicit type (verbose, used at package level):**
```go
var name string
var age int
var height float64
var isActive bool

// Assign after declaration
name = "Alice"
age = 30
```

**Method 2: `var` with initializer (type inferred from value):**
```go
var name = "Alice"      // Go infers: string
var age = 30            // Go infers: int
var height = 1.75       // Go infers: float64
var isActive = true     // Go infers: bool
```

**Method 3: Short variable declaration (`:=`) — most common inside functions:**
```go
func main() {
    name := "Alice"       // Declare AND assign (type inferred)
    age := 30
    height := 1.75
    isActive := true
    
    // := can only be used inside functions, not at package level
    // := must have at least one NEW variable on the left side
}
```

**Method 4: Multiple variables at once:**
```go
var (
    name    string  = "Alice"
    age     int     = 30
    country string  = "India"
)

// Or with :=
x, y, z := 1, 2, 3
firstName, lastName := "Alice", "Smith"
```

**Short declaration rule** — at least one new variable required:
```go
x := 10        // OK: x is new
x := 20        // ERROR: x already declared in this scope
x, y := 20, 30 // OK: y is new (x gets a new value)
```

**Variable scope:**
```go
package main

var globalVar = "I'm global"  // package-level scope

func main() {
    localVar := "I'm local to main"  // function scope
    
    if true {
        blockVar := "I'm local to this block"  // block scope
        fmt.Println(blockVar)  // OK
    }
    
    // fmt.Println(blockVar)  // ERROR: blockVar out of scope
    fmt.Println(localVar)
    fmt.Println(globalVar)
}
```

**The blank identifier `_`** — discard a value you don't need:
```go
// Functions often return multiple values; use _ to ignore one
value, _ := strconv.Atoi("42")  // Ignore the error (dangerous in production!)

// Useful in for loops
for _, v := range []int{1, 2, 3} {
    fmt.Println(v)  // We don't need the index
}
```

### Quick Check
> 1. What is the difference between `var x int` and `x := 0`?
> 2. Where can you NOT use `:=`?
> 3. What does the blank identifier `_` do?

---

## 2. Basic Types

Go has a fixed set of built-in types:

**Boolean:**
```go
var alive bool = true
isReady := false
fmt.Println(!alive)       // false (NOT)
fmt.Println(alive && isReady)  // false (AND)
fmt.Println(alive || isReady)  // true (OR)
```

**Integer types:**
```go
// Signed integers (can be negative):
var a int8  = 127           // -128 to 127
var b int16 = 32767         // -32,768 to 32,767
var c int32 = 2147483647    // -2B to 2B
var d int64 = 9223372036854775807  // Very large

// Unsigned integers (zero and positive only):
var e uint8  = 255          // 0 to 255
var f uint16 = 65535        // 0 to 65,535
var g uint32 = 4294967295   // 0 to 4.3B
var h uint64 = 18446744073709551615

// Platform-dependent (64-bit on most modern systems):
var i int  = -42    // Most common — use this unless you have a reason not to
var j uint = 42     // Unsigned version

// Special integer aliases:
var k byte = 65     // alias for uint8 (used for raw bytes)
var l rune = 'A'    // alias for int32 (used for Unicode code points)
```

**Rule of thumb**: Just use `int` for integers unless you have a specific reason (network protocols, binary formats, performance-critical code).

**Floating point:**
```go
var price float32 = 9.99          // 32-bit float (~7 decimal digits precision)
var precise float64 = 3.141592653589793  // 64-bit float (~15 decimal digits)

// Always prefer float64 — float32's lower precision causes bugs sooner:
a, b := 0.1, 0.2               // runtime float64 values
fmt.Println(a + b)             // 0.30000000000000004 (floating point is inexact!)
fmt.Println(float32(a) + float32(b)) // 0.3 (float32 has so little precision it rounds the error away — until it bites you elsewhere)
```

**Complex numbers** (rare, for scientific computing):
```go
var c complex64  = 3 + 4i
var c2 complex128 = 3 + 4i
fmt.Println(real(c2), imag(c2))  // 3 4
```

**String:**
```go
var greeting string = "Hello, World!"
name := "Go"
message := "Hello, " + name  // String concatenation with +

fmt.Println(len(message))     // Length in BYTES (not characters!)
fmt.Println(message[0])       // First byte: 72 (ASCII for 'H')
fmt.Println(string(message[0])) // "H"
```

**Numeric literal formats:**
```go
// Integer literals
decimal  := 1_000_000   // underscore for readability
hex      := 0xFF        // hexadecimal
octal    := 0o777       // octal
binary   := 0b1010_1010 // binary

// Float literals
pi := 3.14159
big := 1e9         // 1,000,000,000
small := 1.5e-3    // 0.0015
```

### Quick Check
> 1. What is the difference between `int`, `int32`, and `int64`?
> 2. Why should you generally prefer `float64` over `float32`?
> 3. What is a `byte` and what is a `rune` in Go?

---

## 3. Zero Values

In Go, **every variable has a zero value** — the default value it gets if you declare it without initializing:

```go
var i int       // 0
var f float64   // 0.0
var b bool      // false
var s string    // "" (empty string)
var p *int      // nil
var sl []int    // nil (nil slice)
var m map[string]int  // nil (nil map)
var fn func()   // nil

fmt.Println(i, f, b, s, p, sl, m, fn)
// 0 0 false  <nil> [] map[] <nil>
```

**Zero values are a feature, not a bug.** They make Go programs safer — there is no undefined behavior from uninitialized variables:

```go
// In C, this would be undefined (could be any garbage value):
// int x;
// printf("%d", x);  // Undefined behavior!

// In Go:
var x int
fmt.Println(x)  // Always 0. Safe. Predictable.
```

**Practical example — zero values mean you don't always need constructors:**
```go
type Counter struct {
    count int
    mu    sync.Mutex
}

// Zero value of Counter is ready to use!
// count = 0 (zero value for int)
// mu = zero-value Mutex (unlocked and ready)

var c Counter
c.Increment()  // Works immediately, no constructor needed
```

### Quick Check
> 1. What is the zero value for `bool`, `int`, `string`, and `*int`?
> 2. Why are zero values safer than uninitialized variables in languages like C?
> 3. What is the zero value for a slice?

---

## 4. Type Conversions

Go **never does implicit type conversion**. You must convert explicitly:

```go
var i int = 42
var f float64 = float64(i)   // Must explicitly convert
var u uint = uint(f)          // Explicit again

// This does NOT work:
// var f float64 = i  // ERROR: cannot use i (type int) as float64
```

**Safe conversions:**
```go
i := 42
f := float64(i)     // int → float64: always safe
i2 := int(f)        // float64 → int: truncates decimal (3.9 → 3)
b := byte(i)        // int → byte: truncates if > 255
r := rune(i)        // int → rune (Unicode code point)
s := string(rune(65)) // rune → string: gives "A"
```

**String conversions:**
```go
import "strconv"

// int → string (NOT string(42) which gives "*" — the rune with code 42)
s := strconv.Itoa(42)          // "42"
s2 := fmt.Sprintf("%d", 42)    // "42" (slower but more flexible)

// string → int
n, err := strconv.Atoi("42")   // 42, nil
n2, err := strconv.ParseInt("42", 10, 64)  // 42, nil

// float → string
fs := strconv.FormatFloat(3.14, 'f', 2, 64)  // "3.14"

// string → float
fv, err := strconv.ParseFloat("3.14", 64)    // 3.14, nil

// bool → string and back
bs := strconv.FormatBool(true)   // "true"
bv, err := strconv.ParseBool("true")  // true, nil
```

**Common mistake — `string(n)` for integers:**
```go
n := 65
s := string(n)      // WRONG! This gives "A" (rune 65 = 'A')
s2 := strconv.Itoa(n)  // CORRECT: "65"
```

**Type assertions (for interface types):**
```go
var i interface{} = "hello"  // interface{} holds any type

s, ok := i.(string)  // Assert that i holds a string
if ok {
    fmt.Println(s)  // "hello"
}

// Panic version (only use when you're certain):
s = i.(string)  // panics if i is not a string
```

### Quick Check
> 1. Does Go perform automatic type conversion (e.g., int to float64)?
> 2. What does `string(65)` return in Go and why is it surprising?
> 3. What is `strconv.Atoi` used for?

---

## 5. Constants

Constants are values that cannot change at runtime. In Go, constants are computed at **compile time**:

```go
// Basic constants
const Pi = 3.14159265358979323846
const E  = 2.71828182845904523536

const MaxRetries = 3
const ServiceName = "payment-service"
const IsDevelopment = true

// Group constants (like var blocks)
const (
    StatusOK       = 200
    StatusNotFound = 404
    StatusError    = 500
)
```

**Constants can be typed or untyped:**
```go
// Untyped constants — more flexible
const untypedInt = 42      // can be used as any int type
const untypedFloat = 3.14  // can be used as any float type

// Typed constants — strict
const typedInt int32 = 42   // only works as int32

func main() {
    var x int64 = untypedInt    // OK: untyped 42 fits int64
    var y int64 = typedInt      // ERROR: typed int32 can't be int64 without conversion
    
    _ = x
    _ = float64(typedInt)  // OK: explicit conversion
}
```

**Constants are evaluated at compile time:**
```go
const (
    KB = 1024
    MB = KB * 1024      // 1,048,576 — computed at compile time
    GB = MB * 1024      // 1,073,741,824 — computed at compile time
    TB = GB * 1024      // 1,099,511,627,776
)
```

**What can be a constant?**
- Numbers (int, float, complex)
- Strings
- Booleans
- Results of constant expressions

**What CANNOT be a constant?**
- Results of function calls (exception: `len()` of a string constant is itself a constant)
- Slices, maps, structs

```go
const name = "Alice"           // OK
const greeting = "Hello, " + name  // OK: constant expression
const length = len(name)        // OK: len of string constant

// const slice = []int{1, 2, 3}  // ERROR: slices can't be constants
```

### Quick Check
> 1. What is the difference between a typed and untyped constant?
> 2. What kinds of values CAN be constants in Go?
> 3. When are constant expressions evaluated?

---

## 6. iota — Enumerations in Go

Go doesn't have an `enum` keyword. Instead, it has `iota` — a special counter that automatically increments for each constant in a block:

```go
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
    Wednesday      // 3
    Thursday       // 4
    Friday         // 5
    Saturday       // 6
)

fmt.Println(Sunday, Monday, Saturday)  // 0 1 6
```

**Skipping values:**
```go
const (
    _           = iota // Skip 0 (use _ to discard)
    StatusPending      // 1
    StatusActive       // 2
    StatusClosed       // 3
)
```

**Bit flags with iota:**
```go
const (
    ReadPermission  = 1 << iota  // 1 (0001)
    WritePermission              // 2 (0010)
    ExecutePermission            // 4 (0100)
)

// Combine with bitwise OR:
myPerms := ReadPermission | WritePermission  // 3 (0011)

// Check with bitwise AND:
if myPerms & ReadPermission != 0 {
    fmt.Println("Can read")
}
if myPerms & ExecutePermission != 0 {
    fmt.Println("Can execute")  // This won't print
}
```

**iota in expressions:**
```go
const (
    _  = iota             // 0 (skip)
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20 = 1,048,576
    GB                    // 1 << 30 = 1,073,741,824
    TB                    // 1 << 40 = 1,099,511,627,776
)
```

**Typed enumerations (with a custom type):**
```go
type Day int

const (
    Sunday Day = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
)

// Add a String() method for human-readable output
func (d Day) String() string {
    names := [...]string{"Sunday", "Monday", "Tuesday",
        "Wednesday", "Thursday", "Friday", "Saturday"}
    if d < Sunday || d > Saturday {
        return fmt.Sprintf("Day(%d)", int(d))
    }
    return names[d]
}

day := Wednesday
fmt.Println(day)  // "Wednesday" (uses String() method)
```

### Quick Check
> 1. What value does `iota` start at?
> 2. How do you create a bit flag enum with iota?
> 3. Why is it useful to define a custom type for an iota-based enum?

---

## 7. Strings in Depth

Strings in Go are important enough to deserve their own section. They are not as simple as they look.

**What is a Go string?**
A Go string is an **immutable sequence of bytes** (not characters!). It's essentially a read-only slice of bytes:

```go
s := "Hello, 世界"

fmt.Println(len(s))          // 13 (BYTES, not characters!)
// "Hello, " = 7 bytes (ASCII)
// "世" = 3 bytes (UTF-8 encoding)
// "界" = 3 bytes (UTF-8 encoding)

// Indexing gives bytes, not characters
fmt.Println(s[0])            // 72 (byte value of 'H')
fmt.Println(string(s[0]))    // "H"
fmt.Println(s[7])            // 228 (first byte of "世" in UTF-8)
```

**UTF-8 and runes:**
Go source code is always UTF-8. String literals can contain any Unicode character. But iterating by index gives bytes; iterating with `range` gives runes (Unicode code points):

```go
s := "Hello, 世界"

// Byte iteration (wrong for non-ASCII):
for i := 0; i < len(s); i++ {
    fmt.Printf("Byte %d: %d\n", i, s[i])
    // Byte 7: 228, Byte 8: 184, Byte 9: 150 (three bytes for 世)
}

// Rune iteration (correct for Unicode):
for i, r := range s {
    fmt.Printf("Position %d: %c (%d)\n", i, r, r)
    // Position 0: H (72)
    // Position 7: 世 (19990)
    // Position 10: 界 (30028)
}
```

**String builder — efficient concatenation:**
```go
// DON'T do this in a loop (creates many allocations):
s := ""
for i := 0; i < 1000; i++ {
    s += "x"  // Creates a new string every iteration!
}

// DO this — strings.Builder is efficient:
import "strings"

var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("x")
}
result := sb.String()
```

**Essential string operations:**
```go
import "strings"

s := "  Hello, World!  "

strings.TrimSpace(s)           // "Hello, World!"
strings.ToLower(s)             // "  hello, world!  "
strings.ToUpper(s)             // "  HELLO, WORLD!  "
strings.Contains(s, "World")   // true
strings.HasPrefix(s, "  Hel")  // true
strings.HasSuffix(s, "!  ")    // true
strings.Replace(s, "World", "Go", 1) // "  Hello, Go!  "
strings.Split("a,b,c", ",")    // ["a", "b", "c"]
strings.Join([]string{"a","b","c"}, "-") // "a-b-c"
strings.Count(s, "l")          // 3
strings.Index(s, "World")      // 9
```

**Raw string literals** (backticks — no escape sequences):
```go
// Regular string — backslash needed for special chars
path := "C:\\Users\\Alice\\Documents"
json := "{ \"name\": \"Alice\" }"

// Raw string literal — everything is literal
path2 := `C:\Users\Alice\Documents`
json2 := `{ "name": "Alice" }`
multiline := `
    This is a
    multiline string
    no escape needed
`
```

### Quick Check
> 1. What does `len(s)` return for a Go string containing Unicode characters?
> 2. What is the difference between iterating a string with an index vs using `range`?
> 3. Why should you use `strings.Builder` instead of `+` concatenation in a loop?

---

## Summary

- **Variable declaration**: `var x int` (explicit), `var x = 42` (inferred), `x := 42` (short, functions only)
- **Basic types**: bool, int/int8/16/32/64, uint variants, float32/float64, string, byte (uint8), rune (int32)
- **Zero values**: Every uninitialized variable gets its type's zero value (0, false, "", nil)
- **No implicit conversion**: Always convert explicitly with `type(value)` or `strconv` functions
- **Constants**: Compile-time values with `const`; typed or untyped; iota for enumerations
- **iota**: Auto-incrementing counter for constants; useful for enums and bit flags
- **Strings**: Immutable byte sequences; `range` for Unicode-correct iteration; `strings.Builder` for efficient concatenation

Next chapter: Operators, expressions, and type conversion.

---

## Exercises

### Easy
1. Declare a `const` block representing the days of the week using `iota`. Add a `String()` method. Print each day.
2. What is the output of: `fmt.Println(len("Hello, 世界"))` and why?
3. Declare three variables `firstName`, `lastName`, `age` using all three declaration styles (var with type, var with value, `:=`).

### Medium
4. Unit conversion: Write a package `units` with typed constants for file sizes (Byte, Kilobyte, Megabyte, Gigabyte, Terabyte) using `iota` and bit shifting. Write a `HumanReadable(bytes int64) string` function that converts a byte count to a human-readable string like "1.5 MB".
5. String analysis: Write a function `analyzeString(s string) map[string]int` that returns: the number of bytes, the number of runes (Unicode characters), the number of words (split by spaces), and the number of unique characters. Handle multi-byte Unicode correctly.
6. Safe integer operations: Go's integer arithmetic can overflow silently: `var x int8 = 127; x++` gives -128. Write a package `safemath` with `AddInt64(a, b int64) (int64, error)` and `MultiplyInt64(a, b int64) (int64, error)` that return errors on overflow instead of silently wrapping.

### Hard
7. Custom type system: Build a `money` package with a `Money` type that represents a currency amount without floating-point errors. Store amounts as int64 (cents). Provide: `New(amount int64, currency string) Money`, `Add(a, b Money) (Money, error)` (error if different currencies), `Multiply(m Money, factor float64) Money`, `String() string` (e.g., "₹42.50 INR"). Ensure: no float arithmetic used internally, proper handling of rounding, cent-precision arithmetic.
8. iota patterns: Create a comprehensive access control system using bit flags with `iota`. Requirements: (a) Define permissions for a filesystem: Read, Write, Execute, Delete, Admin (5 permissions). (b) Define roles: Guest (read only), User (read+write), Developer (read+write+execute), Admin (all). (c) Write `HasPermission(role, perm)`, `AddPermission(role, perm) Role`, `RemovePermission(role, perm) Role`. (d) Write `PermissionString(perm) string` that returns "rwxda" style string. (e) Test all combinations.
