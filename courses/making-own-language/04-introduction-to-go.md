# Chapter 04: Introduction to Go — The Language We'll Use to Build Astra

> "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software."
> — golang.org

We have studied what programming languages are, how computers work, and how data is represented in binary. Now it is time to write actual code. We are going to build the Astra compiler using Go, and this chapter teaches you everything you need to know about Go to do that.

Why Go? Because it is almost uniquely suited to writing a compiler: it compiles quickly (so our development cycle is fast), has a rich standard library (file I/O, string handling, everything we need), produces fast executables, has a simple and readable syntax, and has been used in real production compilers and tools (Docker, Kubernetes, and many others are written in Go). By the end of this chapter, you will have written your first working Go program and created the skeleton of the Astra compiler project.

Even if you have programmed in another language before, read this chapter carefully — Go has some distinctive features (interfaces, multiple return values, error handling) that we will use constantly in the compiler.

---

## Table of Contents

1. Why Go for Building Astra
2. Installing Go and Setting Up VS Code
3. Hello World in Go — Every Line Explained
4. Go Packages and Imports
5. Variables: var vs :=, Types, Zero Values
6. Basic Types: int, string, bool, float64, byte, rune
7. fmt.Println vs fmt.Printf
8. Functions: Parameters, Return Values, Multiple Returns
9. Structs: Defining, Creating, Accessing Fields
10. Methods on Structs
11. Interfaces: Implicit Implementation
12. Arrays and Slices
13. Maps
14. Control Flow: if/else, for, switch
15. Error Handling: the error Interface
16. File I/O Basics
17. go.mod and the Compiler Project Structure
18. Exercises

---

## 1. Why Go for Building Astra

Before we dive into Go syntax, let us understand the choice.

Writing a compiler in the *same language* the compiler targets (like writing an Astra compiler in Astra) is called **bootstrapping** — it is a cool long-term goal but impossible to start with. We need to write our first Astra compiler in some other language.

Options considered:

| Language | Pros | Cons |
|----------|------|------|
| C | Fast, great for compilers | Manual memory management, no garbage collection, verbose |
| C++ | Fast, powerful | Complex, easy to make mistakes |
| Rust | Safe and fast | Steep learning curve, complex ownership system |
| Python | Easy to write | Too slow for a production compiler |
| **Go** | **Simple, fast, great stdlib, GC, good error handling** | **None significant for this task** |

Go was literally designed at Google by people who were frustrated with C++. It keeps C's simplicity and speed while adding garbage collection, clean interfaces, and excellent tooling. Many real compilers and language tools are written in Go. Go's own compiler (`go build`) was famously rewritten in Go from C. The Go compiler itself is a great learning resource.

---

## 2. Installing Go and Setting Up VS Code

### Installing Go

1. Go to **golang.org/dl** and download the installer for your operating system (macOS, Linux, or Windows).
2. Run the installer. It will install Go at `/usr/local/go` on macOS/Linux or `C:\Program Files\Go` on Windows.
3. The installer automatically adds Go to your PATH.
4. Open a **new** terminal and verify:

```bash
$ go version
go version go1.22.0 linux/amd64
```

If you see a version number, Go is installed correctly.

### Setting Up VS Code

1. Download **Visual Studio Code** from code.visualstudio.com
2. Install the **Go extension** by Google (search "Go" in the Extensions panel — it is the one with 10+ million installs)
3. Open a `.go` file. VS Code will prompt to install Go tools. Click "Install All". These tools provide:
   - Auto-completion (`gopls` language server)
   - Error highlighting
   - Auto-formatting on save (`gofmt`)
   - "Go to definition", "Find references"

### Your First Go Program (Testing the Installation)

Create a directory and file:

```bash
$ mkdir ~/go-test && cd ~/go-test
$ go mod init gotest
$ # Now create main.go with the content below
```

Create `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Go is installed correctly!")
}
```

Run it:

```bash
$ go run main.go
Go is installed correctly!
```

If you see that output, you are ready.

---

## 3. Hello World in Go — Every Line Explained

Let us examine every part of a Go "Hello World" program with the same detail we applied to Astra in Chapter 1:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**Line 1: `package main`**

Every Go file must declare a package. A package is a collection of related Go files that share a namespace. The special package name `main` tells Go: "this is an executable program, not a library." When Go compiles a `main` package, it produces a standalone executable.

**Line 2: (blank)**

Go ignores blank lines. They are for human readability.

**Line 3: `import "fmt"`**

This imports the `fmt` (format) package from Go's standard library. `fmt` provides functions for formatted I/O — printing, scanning, string formatting. Without this import, the code would not compile because `fmt.Println` would be undefined.

**Line 5: `func main() {`**

Declares a function named `main`. In the `main` package, the `main` function is the entry point — the OS will call it first when the program runs. The `{` opens the function body.

**Line 6: `fmt.Println("Hello, World!")`**

Calls the `Println` function from the `fmt` package. `Println` prints its arguments to stdout, separated by spaces, followed by a newline character. The string `"Hello, World!"` is the argument.

**Line 7: `}`**

Closes the function body.

**Important Go rules revealed:**
- No semicolons at end of lines (Go's compiler adds them automatically)
- Opening brace `{` must be on the same line as `func` (unlike C/Java)
- The `main()` function takes no parameters and returns nothing
- Imports must be used — Go is strict; unused imports are compile errors

---

## 4. Go Packages and Imports

Go's standard library contains over 100 packages covering everything from HTTP to cryptography to JSON parsing. For our compiler, the packages we will use most often are:

```go
import (
    "fmt"       // formatted I/O — printing, string formatting
    "os"        // operating system interface — file I/O, exit
    "strings"   // string manipulation
    "unicode"   // Unicode character properties
    "strconv"   // string ↔ number conversion
    "path/filepath" // file path manipulation
    "errors"    // creating error values
)
```

When importing multiple packages, use the grouped import syntax (parentheses). Each package is on its own line, in quotes.

**Package access syntax:** `packagename.ExportedName`

Only names that start with a capital letter are exported (visible from outside the package). Functions like `fmt.Println`, `os.ReadFile` start with capitals. Internal helpers start with lowercase and are package-private.

```go
// EXPORTED (usable from other packages):
func Tokenize(input string) []Token  // capital T

// UNEXPORTED (only usable within the same package):
func skipWhitespace(input string)    // lowercase s
```

This convention is how Go handles public vs private — no `public`/`private` keywords needed.

---

## 5. Variables: var vs :=, Types, Zero Values

Go has two main ways to declare variables.

### Method 1: `var` keyword (explicit type)

```go
var name string = "Aditya"
var age int = 25
var pi float64 = 3.14159
var isAdult bool = true
```

You can also declare without initializing — Go sets the **zero value** automatically:

```go
var name string    // zero value: "" (empty string)
var age int        // zero value: 0
var pi float64     // zero value: 0.0
var isAdult bool   // zero value: false
var ptr *int       // zero value: nil
```

**Zero values are a critical Go feature.** In C, uninitialized variables have garbage values. In Go, every variable is initialized to its type's zero value. This eliminates an entire class of bugs.

### Method 2: Short declaration `:=` (type inferred)

```go
name := "Aditya"     // Go infers type: string
age := 25            // Go infers type: int
pi := 3.14159        // Go infers type: float64
isAdult := true      // Go infers type: bool
```

The `:=` operator declares AND assigns. It can only be used inside functions. Go infers the type from the right-hand side expression. This is the most common form in Go code.

### Which to Use?

| Situation | Use |
|-----------|-----|
| At package level (outside functions) | `var` |
| Inside functions, when initializing | `:=` |
| When type matters explicitly (e.g., `int32` not `int`) | `var x int32 = 0` |
| When initializing to zero value | `var x int` (clearer than `x := 0`) |

In our compiler code, we will mostly use `:=` inside functions and `var` at the package level.

### Mutability

In Go, all variables declared with `var` or `:=` are mutable by default — you can reassign them with `=`:

```go
name := "Aditya"
name = "Aditya Pathak"  // OK

// Constants cannot be reassigned
const MaxTokens = 10000
MaxTokens = 99999  // COMPILE ERROR
```

---

## 6. Basic Types in Go

Go's basic types and their relevance to our compiler:

```go
// Integer types
var i int = 42           // platform-native size (usually 64-bit)
var i8 int8 = -128       // 8-bit signed: -128 to 127
var i16 int16 = 32767    // 16-bit signed
var i32 int32 = 2147483647 // 32-bit signed
var i64 int64 = 9223372036854775807 // 64-bit signed

// Unsigned integers
var u uint = 42           // platform-native unsigned
var u8 uint8 = 255        // 8-bit: 0 to 255 (also: byte)
var u16 uint16 = 65535
var u32 uint32 = 4294967295
var u64 uint64 = 18446744073709551615

// Floating point
var f32 float32 = 3.14   // 32-bit IEEE 754
var f64 float64 = 3.14159265358979 // 64-bit IEEE 754

// Boolean
var b bool = true        // true or false

// Strings (immutable sequence of UTF-8 bytes)
var s string = "hello"

// Byte and rune
var by byte = 'A'        // byte = uint8, for ASCII characters
var r rune = '世'         // rune = int32, for Unicode code points
```

**For our compiler, the key types we will use:**
- `string` — for token text, identifiers, string literals
- `int` — for line numbers, positions, lengths
- `bool` — for flags
- `byte` — for reading source file bytes one at a time
- `rune` — for handling Unicode in source files

---

## 7. fmt.Println vs fmt.Printf

Two essential printing functions:

### fmt.Println

```go
fmt.Println("Hello")              // Hello\n
fmt.Println("Hello", "World")     // Hello World\n  (space-separated)
fmt.Println("x =", 42)            // x = 42\n
fmt.Println(3.14)                 // 3.14\n
```

`Println` accepts any number of values, prints them separated by spaces, always adds a newline at the end.

### fmt.Printf

```go
fmt.Printf("Hello, %s!\n", "World")      // Hello, World!
fmt.Printf("Age: %d\n", 25)              // Age: 25
fmt.Printf("Pi: %.2f\n", 3.14159)        // Pi: 3.14
fmt.Printf("Hex: %x\n", 255)             // Hex: ff
fmt.Printf("Bool: %v\n", true)           // Bool: true
fmt.Printf("Type: %T\n", 42)             // Type: int
fmt.Printf("Struct: %+v\n", someStruct)  // Struct: {field1:val field2:val}
```

Format verbs: `%s` = string, `%d` = decimal integer, `%f` = float, `%v` = any value (default format), `%T` = type name, `%x` = hex, `%b` = binary, `%p` = pointer address.

### fmt.Sprintf

Returns a formatted string instead of printing it:

```go
msg := fmt.Sprintf("Error at line %d: %s", lineNum, errorMsg)
// Use msg elsewhere, like in an error return
```

We will use `fmt.Sprintf` constantly in error messages from our compiler.

### fmt.Fprintf

Like Printf but writes to an `io.Writer` instead of stdout:

```go
fmt.Fprintf(os.Stderr, "Error: %s\n", msg)  // write to stderr
```

---

## 8. Functions: Parameters, Return Values, Multiple Returns

### Basic Function

```go
func greet(name string) {
    fmt.Println("Hello,", name)
}

func main() {
    greet("Aditya")  // prints: Hello, Aditya
}
```

### Function with Return Value

```go
func add(a int, b int) int {
    return a + b
}

// Shorthand when parameters have the same type:
func add(a, b int) int {
    return a + b
}
```

### Multiple Return Values

This is one of Go's most distinctive and useful features. Functions can return multiple values:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 3)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Result: %.4f\n", result)  // Result: 3.3333
}
```

**Multiple return values are fundamental to Go's error handling pattern.** Almost every function in our compiler that can fail will return `(SomeValue, error)`. The caller always checks the error.

### Named Return Values

```go
func minMax(slice []int) (min, max int) {
    min = slice[0]
    max = slice[0]
    for _, v := range slice {
        if v < min { min = v }
        if v > max { max = v }
    }
    return  // "naked return" — returns the named values
}
```

Named returns are useful for documentation and for deferred functions, but naked returns should only be used in short functions.

### Variadic Functions

Functions that accept a variable number of arguments:

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)        // = 6
sum(1, 2, 3, 4, 5)  // = 15
```

`fmt.Println`, `fmt.Printf` are variadic — that is how they accept any number of arguments.

---

## 9. Structs: Defining, Creating, Accessing Fields

A **struct** is a composite type that groups fields together. It is Go's way of defining custom data types. We will define dozens of structs in our compiler — one for each AST node type, one for each token type, one for the type checker's context, etc.

### Defining a Struct

```go
type Token struct {
    Type    TokenType  // an enum (defined as const)
    Literal string     // the raw text, e.g. "fn", "42", "main"
    Line    int        // which line in source file
    Column  int        // which column in source file
}
```

### Creating Struct Values

```go
// Method 1: Field names (preferred — order-independent, clear)
tok := Token{
    Type:    TOKEN_IDENT,
    Literal: "main",
    Line:    1,
    Column:  4,
}

// Method 2: Positional (fragile — order-dependent, avoid)
tok := Token{TOKEN_IDENT, "main", 1, 4}

// Method 3: Zero value, then assign
var tok Token
tok.Type = TOKEN_IDENT
tok.Literal = "main"
```

### Accessing Fields

```go
fmt.Println(tok.Type)     // prints the token type
fmt.Println(tok.Literal)  // prints "main"
fmt.Println(tok.Line)     // prints 1
```

### Pointers to Structs

When passing structs to functions, Go copies the entire struct. For large structs (like our AST nodes), we use pointers to avoid copying:

```go
// Passing by value: function gets a COPY, cannot modify original
func printToken(tok Token) {
    fmt.Println(tok.Literal)
}

// Passing by pointer: function gets the ADDRESS, can modify original
func setLiteral(tok *Token, s string) {
    tok.Literal = s  // modifies the original
}

// Pointer syntax:
tok := &Token{Type: TOKEN_IDENT}  // tok is *Token (pointer to Token)
tok.Literal = "fn"                // Go auto-dereferences: same as (*tok).Literal
```

In our compiler, almost all AST nodes will be passed as pointers (`*ASTNode`).

---

## 10. Methods on Structs

In Go, you add behavior to types by defining **methods** — functions with a *receiver* parameter specifying which type they belong to.

```go
type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}

// Method with VALUE receiver (tok is a copy)
func (tok Token) String() string {
    return fmt.Sprintf("Token(%s, %q, line %d)", tok.Type, tok.Literal, tok.Line)
}

// Method with POINTER receiver (tok is a pointer — can modify the struct)
func (tok *Token) SetLine(line int) {
    tok.Line = line
}

func main() {
    tok := Token{Type: TOKEN_IDENT, Literal: "age", Line: 5}
    fmt.Println(tok.String())  // Token(IDENT, "age", line 5)
    tok.SetLine(10)
    fmt.Println(tok.Line)      // 10
}
```

**Rule of thumb:** Use pointer receivers `*T` when:
1. The method needs to modify the receiver's fields
2. The struct is large (avoid copying)
3. **Consistency**: if any method uses pointer receiver, all methods should

Use value receivers `T` when:
- The struct is tiny (like `int` wrapper)
- The method does not need to modify the struct
- You want copies (e.g., for thread safety)

In our compiler, we will use pointer receivers for almost everything because AST nodes and other compiler structures are modified during analysis.

---

## 11. Interfaces: Implicit Implementation

Go interfaces are one of the language's most powerful features and most different from Java/C# interfaces.

**An interface defines a set of method signatures.** Any type that has those methods *automatically* satisfies the interface — no explicit declaration needed.

```go
// The interface
type Stringer interface {
    String() string
}

// Any type that has a String() string method automatically satisfies Stringer
type Token struct { Literal string }
func (t Token) String() string { return t.Literal }

type ASTNode struct { Kind string }
func (n ASTNode) String() string { return n.Kind }

// Now both Token and ASTNode satisfy Stringer
func printAnything(s Stringer) {
    fmt.Println(s.String())
}

printAnything(Token{Literal: "fn"})   // fn
printAnything(ASTNode{Kind: "BinaryExpr"}) // BinaryExpr
```

This is called **structural typing** or **duck typing** — if it walks like a duck and quacks like a duck, it is a duck. No need to write `implements Stringer` anywhere.

### The Most Important Interface: error

Go's error handling is built on a single interface:

```go
type error interface {
    Error() string
}
```

Any type with an `Error() string` method satisfies the `error` interface. By convention, functions that can fail return an `error` as their last return value. The caller checks if it is `nil` (no error) or non-nil (an error occurred).

```go
// Creating errors
err1 := errors.New("file not found")
err2 := fmt.Errorf("parse error at line %d: %s", 42, "unexpected token")

// Checking errors
if err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
}
```

We will use this pattern extensively. Our parser, type checker, and code generator will all return errors when they encounter problems.

### The Empty Interface: any

The empty interface `interface{}` (or `any` in Go 1.18+) is satisfied by every type — it is Go's way of saying "any value":

```go
func printValue(v any) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

printValue(42)         // Value: 42, Type: int
printValue("hello")    // Value: hello, Type: string
printValue([]int{1,2}) // Value: [1 2], Type: []int
```

In our compiler's AST, some nodes will use `any` for their value fields (e.g., an integer literal might store `any` which is actually an `int64`).

---

## 12. Arrays and Slices

In Go, an **array** has a fixed length, while a **slice** has a dynamic length. We almost always use slices.

### Arrays (fixed size — rarely used)

```go
var arr [3]int = [3]int{1, 2, 3}
arr[0] = 10  // modify first element
fmt.Println(arr[0], arr[1], arr[2])  // 10 2 3
fmt.Println(len(arr))  // 3
```

### Slices (dynamic — what we actually use)

```go
// Creating slices
s1 := []int{1, 2, 3}          // slice literal
s2 := make([]int, 5)           // slice of 5 zeros
s3 := make([]int, 0, 10)       // empty slice, capacity 10

// Appending (can grow dynamically)
s1 = append(s1, 4, 5)          // s1 = [1, 2, 3, 4, 5]
s2 = append(s2, s1...)          // append another slice with ...

// Accessing
fmt.Println(s1[0])              // 1 (first element)
fmt.Println(s1[len(s1)-1])      // 5 (last element)

// Slicing (sub-slice)
sub := s1[1:3]                  // [2, 3] — from index 1 to 2 (exclusive)
sub2 := s1[:3]                  // [1, 2, 3] — first 3 elements
sub3 := s1[2:]                  // [3, 4, 5] — from index 2 to end

// Length and capacity
fmt.Println(len(s1), cap(s1))   // 5  8 (capacity doubled when grew past 3)
```

### Iterating over slices

```go
tokens := []Token{
    {Type: TOKEN_IDENT, Literal: "fn"},
    {Type: TOKEN_IDENT, Literal: "main"},
}

// Range loop (most common)
for i, tok := range tokens {
    fmt.Printf("[%d] %s\n", i, tok.Literal)
}

// If you don't need the index, use _
for _, tok := range tokens {
    fmt.Println(tok.Literal)
}

// Classic index loop
for i := 0; i < len(tokens); i++ {
    fmt.Println(tokens[i].Literal)
}
```

In our compiler, the token stream from the lexer will be a `[]Token`. The list of statements in a function body will be `[]Statement`. Everything is slices.

---

## 13. Maps

A **map** is a key-value data structure (also called a hash map, dictionary, or associative array in other languages).

```go
// Creating maps
m1 := map[string]int{"a": 1, "b": 2}       // map literal
m2 := make(map[string]int)                   // empty map

// Setting values
m2["key"] = 42
m2["another"] = 99

// Getting values
val := m2["key"]             // 42
missing := m2["no_such_key"] // 0 (zero value — no panic!)

// Checking existence (important pattern!)
val, exists := m2["key"]
if exists {
    fmt.Println("Found:", val)
} else {
    fmt.Println("Not found")
}

// Deleting
delete(m2, "key")

// Iterating
for k, v := range m2 {
    fmt.Printf("%s: %d\n", k, v)
}

// Length
fmt.Println(len(m2))
```

Maps are critical for our compiler:
- **Symbol table**: `map[string]*Symbol` — maps variable names to their declarations
- **Type environment**: `map[string]Type` — maps type names to their definitions
- **Keyword lookup**: `map[string]TokenType` — maps keyword strings to token types

---

## 14. Control Flow: if/else, for, switch

### if/else

```go
age := 20

if age >= 18 {
    fmt.Println("Adult")
} else if age >= 13 {
    fmt.Println("Teenager")
} else {
    fmt.Println("Child")
}

// if with initialization statement (very common in Go)
if tok, ok := nextToken(); ok {
    process(tok)
}
```

The "if with init" pattern is used constantly: `if err := doSomething(); err != nil { ... }`.

### for Loop (Go's Only Loop)

Go has only one loop keyword: `for`. But it covers all use cases:

```go
// Classic C-style for loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style (omit init and post)
for i < 10 {
    i++
}

// Infinite loop
for {
    // breaks out with break
    if done { break }
}

// Range loop over slice
for i, v := range mySlice {
    fmt.Println(i, v)
}

// Range loop over map
for key, value := range myMap {
    fmt.Println(key, value)
}

// Range loop over string (iterates runes, not bytes!)
for i, r := range "hello" {
    fmt.Printf("%d: %c\n", i, r)
}
```

### switch

```go
// Switch on a value
switch tok.Type {
case TOKEN_PLUS:
    return parseAddition()
case TOKEN_MINUS:
    return parseSubtraction()
case TOKEN_IDENT:
    return parseIdentifier()
default:
    return nil, fmt.Errorf("unexpected token: %s", tok.Literal)
}

// Switch with no value (acts like if/else chain)
switch {
case age < 13:
    fmt.Println("Child")
case age < 18:
    fmt.Println("Teenager")
default:
    fmt.Println("Adult")
}
```

In Go, `switch` cases do NOT fall through by default (unlike C). Each case is independent. Use `fallthrough` explicitly if you want C-style behavior (rare).

---

## 15. Error Handling: The error Interface

Go's error handling is explicit and deliberate — no exceptions. Functions that can fail return an error as the last return value.

```go
// Defining custom error types for our compiler
type ParseError struct {
    Message string
    Line    int
    Column  int
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at line %d, column %d: %s",
        e.Line, e.Column, e.Message)
}

// Using it
func parseExpression() (Expression, error) {
    if somethingWentWrong {
        return nil, &ParseError{
            Message: "expected expression",
            Line:    lexer.Line,
            Column:  lexer.Column,
        }
    }
    return expr, nil
}

// Calling it
expr, err := parseExpression()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

**The patterns:**
- Functions return `(result, error)` — check error before using result
- If no error: `return result, nil` — nil means "no error"
- If error: `return nil, err` or `return zero, err` — never ignore errors
- Wrap errors for context: `fmt.Errorf("while parsing function: %w", err)`

We will use `%w` (wrap) to chain errors in our compiler so that error messages show the full chain of what went wrong.

---

## 16. File I/O Basics

Our compiler reads `.as` source files and writes executable output. Here are the key file operations:

```go
import (
    "os"
    "fmt"
)

// Reading an entire file into memory (our compiler does this)
data, err := os.ReadFile("hello.as")
if err != nil {
    fmt.Fprintf(os.Stderr, "Cannot read file: %v\n", err)
    os.Exit(1)
}
// data is []byte (a slice of bytes)
source := string(data)  // convert to string for processing

// Writing a file
output := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}  // "Hello"
err = os.WriteFile("output.bin", output, 0644)
if err != nil {
    fmt.Fprintf(os.Stderr, "Cannot write file: %v\n", err)
    os.Exit(1)
}

// File permissions: 0644 = owner can read+write, others can read
// In binary: 110 100 100 (owner:rw, group:r, other:r)

// Check if file exists
_, err = os.Stat("hello.as")
if os.IsNotExist(err) {
    fmt.Println("File does not exist")
}

// Get command-line arguments
args := os.Args  // os.Args[0] is the program name, os.Args[1] is first arg
if len(args) < 2 {
    fmt.Println("Usage: astrac <source.as>")
    os.Exit(1)
}
filename := os.Args[1]  // the first argument given by user
```

---

## 17. go.mod and the Compiler Project Structure

Go uses **modules** to manage projects and dependencies. A module is defined by a `go.mod` file at the project root.

### Creating the astrac Module

```bash
# Create the project directory
$ mkdir -p ~/astrac
$ cd ~/astrac

# Initialize the Go module
$ go mod init github.com/astra-lang/astrac
```

This creates `go.mod`:

```
module github.com/astra-lang/astrac

go 1.22
```

The module name `github.com/astra-lang/astrac` is how other Go code would import this module if it were published. For now, it just names the project.

### The Complete Project Skeleton

Here is the full project structure for the astrac compiler — the tree we will build over this entire guide. Today we create just the foundation:

```
astrac/
├── go.mod                    ← module definition
├── main.go                   ← compiler entry point
├── lexer/
│   ├── lexer.go              ← tokenizer (Chapter 5-6)
│   └── token.go              ← token type definitions
├── parser/
│   ├── parser.go             ← recursive descent parser (Chapter 7-9)
│   └── ast.go                ← AST node definitions
├── types/
│   ├── types.go              ← type system definitions (Chapter 05 milestone)
│   └── checker.go            ← type checker (Chapter 10-12)
├── codegen/
│   ├── codegen.go            ← code generation (Chapter 15-25)
│   └── elf.go                ← ELF file writer
├── stdlib/
│   └── stdlib.go             ← standard library support
└── README.md
```

---

## Astra Build Milestone: The astrac Project Skeleton

Let us create the actual working skeleton of the compiler. This is the first real code we write!

### Step 1: Create the directories

```bash
$ mkdir -p ~/astrac/lexer ~/astrac/parser ~/astrac/types ~/astrac/codegen
```

### Step 2: Create go.mod

Run in `~/astrac/`:

```bash
$ go mod init github.com/astra-lang/astrac
```

### Step 3: Create main.go

This is the compiler entry point. It reads command-line arguments and orchestrates the compilation pipeline.

```go
// main.go — Astra Compiler Entry Point
// Usage: astrac build <source.as>

package main

import (
	"fmt"
	"os"
)

// Version of the Astra compiler
const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "build":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: 'build' requires a source file")
			fmt.Fprintln(os.Stderr, "Usage: astrac build <source.as>")
			os.Exit(1)
		}
		sourceFile := os.Args[2]
		if err := buildFile(sourceFile); err != nil {
			fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Printf("Astra Compiler v%s\n", Version)
		fmt.Println("Built with Go")
		fmt.Println("Target: native executable")

	case "help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// buildFile compiles an Astra source file to a native executable.
// This function will grow considerably as we add compiler stages.
func buildFile(filename string) error {
	// Check that the file exists and is readable
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read source file %q: %w", filename, err)
	}

	source := string(data)
	fmt.Printf("Compiling: %s (%d bytes)\n", filename, len(source))

	// TODO: Stage 1 — Lexer (Chapter 5-6)
	// tokens, err := lexer.Tokenize(source)

	// TODO: Stage 2 — Parser (Chapter 7-9)
	// ast, err := parser.Parse(tokens)

	// TODO: Stage 3 — Type Checker (Chapter 10-12)
	// err = types.Check(ast)

	// TODO: Stage 4 — Code Generator (Chapter 15-25)
	// err = codegen.Generate(ast, outputFile)

	fmt.Println("Compilation pipeline not yet implemented.")
	fmt.Println("We're building it chapter by chapter!")
	return nil
}

// printUsage prints usage information to stdout.
func printUsage() {
	fmt.Printf("Astra Compiler v%s\n\n", Version)
	fmt.Println("USAGE:")
	fmt.Println("  astrac build <source.as>   Compile an Astra source file")
	fmt.Println("  astrac version             Print compiler version")
	fmt.Println("  astrac help                Print this message")
	fmt.Println()
	fmt.Println("EXAMPLE:")
	fmt.Println("  astrac build hello.as      Produces ./hello executable")
}
```

### Step 4: Create lexer/token.go

```go
// lexer/token.go — Token type definitions for the Astra compiler
// A token is the smallest meaningful unit of Astra source code.
// The lexer reads characters and produces a stream of tokens.

package lexer

import "fmt"

// TokenType identifies the kind of a token.
// We use a custom type (not plain int) so the compiler catches type errors.
type TokenType string

// All token types in the Astra language.
// We define them as constants so they can be used in switch statements.
const (
	// Literals
	TOKEN_INT    TokenType = "INT"    // 42, 0, -17
	TOKEN_FLOAT  TokenType = "FLOAT"  // 3.14, 0.0
	TOKEN_STRING TokenType = "STRING" // "hello"
	TOKEN_BOOL   TokenType = "BOOL"   // true, false

	// Identifiers
	TOKEN_IDENT TokenType = "IDENT" // variable names, function names

	// Keywords
	TOKEN_FN     TokenType = "fn"
	TOKEN_LET    TokenType = "let"
	TOKEN_RETURN TokenType = "return"
	TOKEN_IF     TokenType = "if"
	TOKEN_ELSE   TokenType = "else"
	TOKEN_FOR    TokenType = "for"
	TOKEN_IN     TokenType = "in"
	TOKEN_STRUCT TokenType = "struct"
	TOKEN_IMPL   TokenType = "impl"
	TOKEN_TRUE   TokenType = "true"
	TOKEN_FALSE  TokenType = "false"

	// Arithmetic Operators
	TOKEN_PLUS     TokenType = "+"
	TOKEN_MINUS    TokenType = "-"
	TOKEN_ASTERISK TokenType = "*"
	TOKEN_SLASH    TokenType = "/"
	TOKEN_PERCENT  TokenType = "%"

	// Comparison Operators
	TOKEN_EQ     TokenType = "=="
	TOKEN_NEQ    TokenType = "!="
	TOKEN_LT     TokenType = "<"
	TOKEN_GT     TokenType = ">"
	TOKEN_LTE    TokenType = "<="
	TOKEN_GTE    TokenType = ">="

	// Logical Operators
	TOKEN_AND TokenType = "&&"
	TOKEN_OR  TokenType = "||"
	TOKEN_NOT TokenType = "!"

	// Assignment
	TOKEN_ASSIGN  TokenType = "="

	// Delimiters
	TOKEN_LPAREN    TokenType = "("
	TOKEN_RPAREN    TokenType = ")"
	TOKEN_LBRACE    TokenType = "{"
	TOKEN_RBRACE    TokenType = "}"
	TOKEN_LBRACKET  TokenType = "["
	TOKEN_RBRACKET  TokenType = "]"
	TOKEN_COMMA     TokenType = ","
	TOKEN_SEMICOLON TokenType = ";"
	TOKEN_COLON     TokenType = ":"
	TOKEN_DOT       TokenType = "."
	TOKEN_DOTDOT    TokenType = ".."  // range operator
	TOKEN_ARROW     TokenType = "->"  // return type arrow

	// Special
	TOKEN_EOF     TokenType = "EOF"     // end of file
	TOKEN_ILLEGAL TokenType = "ILLEGAL" // unrecognized character
)

// Keywords maps Astra keyword strings to their token types.
// When the lexer reads an identifier, it checks this map to see
// if it's actually a keyword.
var Keywords = map[string]TokenType{
	"fn":     TOKEN_FN,
	"let":    TOKEN_LET,
	"return": TOKEN_RETURN,
	"if":     TOKEN_IF,
	"else":   TOKEN_ELSE,
	"for":    TOKEN_FOR,
	"in":     TOKEN_IN,
	"struct": TOKEN_STRUCT,
	"impl":   TOKEN_IMPL,
	"true":   TOKEN_TRUE,
	"false":  TOKEN_FALSE,
}

// LookupIdent checks if an identifier string is a keyword.
// Returns the keyword token type, or TOKEN_IDENT if it's not a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}

// Token represents a single token produced by the lexer.
// Every piece of Astra source code is broken into tokens.
type Token struct {
	Type    TokenType // what kind of token is this?
	Literal string    // the raw text from the source file
	Line    int       // which line this token appears on (1-indexed)
	Column  int       // which column this token starts at (1-indexed)
}

// String returns a human-readable representation of the token.
// This is used for debugging and error messages.
func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %q, line %d, col %d)",
		t.Type, t.Literal, t.Line, t.Column)
}

// IsKeyword returns true if this token is an Astra keyword.
func (t Token) IsKeyword() bool {
	_, isKeyword := Keywords[string(t.Type)]
	return isKeyword
}

// IsLiteral returns true if this token represents a literal value
// (integer, float, string, or boolean).
func (t Token) IsLiteral() bool {
	return t.Type == TOKEN_INT ||
		t.Type == TOKEN_FLOAT ||
		t.Type == TOKEN_STRING ||
		t.Type == TOKEN_TRUE ||
		t.Type == TOKEN_FALSE
}
```

### Step 5: Verify the build

```bash
$ cd ~/astrac
$ go build ./...
# Should produce no errors

$ go run main.go version
Astra Compiler v0.1.0
Built with Go
Target: native executable

$ go run main.go help
Astra Compiler v0.1.0

USAGE:
  astrac build <source.as>   Compile an Astra source file
  astrac version             Print compiler version
  astrac help                Print this message

EXAMPLE:
  astrac build hello.as      Produces ./hello executable

$ go run main.go build hello.as
# (If hello.as exists)
Compiling: hello.as (57 bytes)
Compilation pipeline not yet implemented.
We're building it chapter by chapter!
```

The project skeleton is alive. Every chapter from here builds on this foundation.

---

## 18. Exercises

1. **Go Hello World Expansion**: Modify the Hello World program to also print your name, the current Go version (`runtime.Version()`), and the operating system (`runtime.GOOS`). You will need to import the `runtime` package.
   *Hint: `import "runtime"`, then `runtime.Version()` and `runtime.GOOS` are string variables.*

2. **Struct Practice**: Define a Go struct `ASTNode` with fields: `Kind` (string), `Value` (string), `Children` (a slice of `*ASTNode`), and `Line` (int). Write a method `AddChild(child *ASTNode)` that appends to `Children`. Write a `String()` method that prints the kind and value. Create a simple tree representing `1 + 2` and print it.
   *Hint: The tree would be: BinaryExpr node with two IntLiteral children.*

3. **Map Challenge**: Write a Go function `wordCount(s string) map[string]int` that takes a string and returns a map of each word to how many times it appears. Test it with a sentence containing repeated words.
   *Hint: Use `strings.Split(s, " ")` to split into words. `strings.ToLower` normalizes case.*

4. **Multiple Return Values**: Write a Go function `parseInteger(s string) (int64, error)` that parses a string as an integer. Use `strconv.ParseInt`. Return an error if the string is not a valid integer. Call it with "42", "-17", "hello", and "99999999999999999999" (overflow).
   *Hint: `strconv.ParseInt(s, 10, 64)` parses base-10 64-bit integer.*

5. **Interface Practice**: Define an interface `Node` with methods `TokenLiteral() string` and `String() string`. Then create two concrete types: `IntegerLiteral` (with field `Value int64`) and `Identifier` (with field `Name string`). Make both implement `Node`. Write a function that takes a `Node` and prints it.
   *Hint: You do not need to write `implements Node` anywhere — if the methods match, it satisfies the interface automatically.*

6. **Error Handling**: Modify the `buildFile` function in `main.go` to also check that the filename ends with `.as` (Astra source file extension). If it does not, return an appropriate error. Test with `astrac build main.go` — it should print an error.
   *Hint: `strings.HasSuffix(filename, ".as")` checks the extension.*

7. **Slice Manipulation**: Write a Go function `reverseTokens(tokens []Token) []Token` that returns a new slice with the tokens in reverse order. Do not modify the original slice. Test it with 5 sample tokens.
   *Hint: Create a new slice, loop from len(tokens)-1 down to 0, append each element.*

8. **The Keyword Map**: The `Keywords` map in `token.go` maps strings to TokenTypes. Add three more Astra keywords that you think would be useful: `const`, `mut`, `pub`. Define their token types and add them to the map. Explain in a comment why each keyword might be useful in a language design context.
   *Hint: `const` for compile-time constants, `mut` for mutable variables (like Rust), `pub` for public visibility.*

---

## Summary: Key Concepts

| Go Feature | Syntax | Compiler Usage |
|-----------|--------|----------------|
| Package declaration | `package main` | Each compiler stage is a package |
| Import | `import "fmt"` | Each stage imports others |
| Variables | `x := 5` or `var x int` | Used everywhere |
| Zero values | `var s string` = `""` | Safe defaults, no garbage values |
| Functions | `func name(p Type) RetType` | Every compiler operation is a function |
| Multiple returns | `return val, err` | Error handling throughout compiler |
| Structs | `type Foo struct { ... }` | Every AST node, Token, Type |
| Pointer receivers | `func (f *Foo) Method()` | Modifying AST nodes during analysis |
| Interfaces | `type I interface { M() }` | Node, Type, Stringer all as interfaces |
| Slices | `[]Token`, `append()` | Token streams, AST children |
| Maps | `map[string]TokenType` | Symbol tables, keyword lookup |
| for/range | `for i, v := range slice` | Processing token streams |
| switch | `switch tok.Type { case ... }` | Dispatching in parser and code gen |
| error interface | `return nil, errors.New(...)` | All error reporting in compiler |
| os.ReadFile | `os.ReadFile(filename)` | Reading .as source files |
| os.Args | `os.Args[1]` | Command-line argument parsing |
| go.mod | `go mod init ...` | Project organization |
