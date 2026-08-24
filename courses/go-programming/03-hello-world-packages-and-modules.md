# Chapter 03: Hello World, Packages, and Modules

You've installed Go and written your first program. Now let's slow down and understand exactly what's happening. This chapter covers the fundamental building blocks that every Go program is made of: packages, the `main` function, how Go organizes code, and how the build system works. These concepts appear in every Go file you will ever write.

## Table of Contents

1. [Anatomy of a Go Program](#1-anatomy-of-a-go-program)
2. [Packages — Go's Unit of Organization](#2-packages--gos-unit-of-organization)
3. [Import Statements](#3-import-statements)
4. [The main Package and main Function](#4-the-main-package-and-main-function)
5. [Exported vs Unexported Names](#5-exported-vs-unexported-names)
6. [Writing Your Own Package](#6-writing-your-own-package)
7. [Go Modules Deep Dive](#7-go-modules-deep-dive)
8. [Build, Run, Install](#8-build-run-install)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Anatomy of a Go Program

Here is a complete Go program with every important concept annotated:

```go
// Package declaration — EVERY .go file must start with this
// This file belongs to the "main" package
package main

// Import section — bring in packages we want to use
import (
    "fmt"     // Standard library: printing functions
    "os"      // Standard library: OS interaction
)

// The main function — program execution starts here
// It takes no arguments and returns nothing
func main() {
    // Call Println from the fmt package
    fmt.Println("Hello, World!")
    
    // os.Args contains command-line arguments
    // os.Args[0] is the program name itself
    if len(os.Args) > 1 {
        fmt.Println("Hello,", os.Args[1])
    }
}
```

**The three mandatory parts of every Go executable:**
1. `package main` — the file belongs to the `main` package
2. `import` — bring in needed packages (optional, but almost always present)
3. `func main()` — the entry point

Let's build and run this:
```bash
go run main.go
# Output: Hello, World!

go run main.go Alice
# Output: Hello, World!
#         Hello, Alice
```

### Quick Check
> 1. What three elements must a Go executable file have?
> 2. What is `os.Args` and what does `os.Args[0]` contain?
> 3. Does every Go file need a `func main()`?

---

## 2. Packages — Go's Unit of Organization

In Go, code is organized into **packages**. A package is a directory containing one or more `.go` files that all share the same `package` declaration.

```
Think of packages like folders:
  
  Your project:
  ├── main.go            package main
  ├── math/
  │   ├── add.go         package math
  │   └── subtract.go    package math  (same package, different file)
  └── strings/
      └── reverse.go     package strings
  
  All .go files in the same directory MUST have the same package name.
  Except: *_test.go files can use "package x_test" for black-box testing.
```

**Rules for packages:**
1. All files in the same directory must have the same package name
2. The package name is usually the last component of the directory name (`utils/` → `package utils`)
3. Exception: The `main` package can be in any directory
4. Multiple packages cannot be in the same directory

**Creating a package:**
```
project/
├── go.mod
├── main.go
└── calculator/
    ├── add.go
    └── multiply.go
```

`calculator/add.go`:
```go
package calculator  // Package name matches directory name

// Add returns the sum of two integers
// Note: starts with capital letter = exported (visible outside this package)
func Add(a, b int) int {
    return a + b
}

// helper is unexported (lowercase) — only visible within the calculator package
func helper() {
    // internal helper
}
```

`calculator/multiply.go`:
```go
package calculator  // Same package name

// Multiply returns the product of two integers
func Multiply(a, b int) int {
    return a * b
}
```

`main.go`:
```go
package main

import (
    "fmt"
    "github.com/yourname/project/calculator"  // Import the package
)

func main() {
    sum := calculator.Add(3, 4)
    product := calculator.Multiply(3, 4)
    fmt.Println(sum, product)  // 7 12
}
```

**Package naming conventions:**
- Short, lowercase, no underscores: `http`, `json`, `fmt`, `io`
- One word preferred: `calculator` not `my_calculator`
- Don't repeat the package name in function names: `json.Marshal()` not `json.MarshalJSON()`

### Quick Check
> 1. Can two different `.go` files in the same directory have different package names?
> 2. What is the package naming convention in Go?
> 3. What is the difference between the package name and the import path?

---

## 3. Import Statements

The `import` statement brings packages into your file:

```go
// Single import
import "fmt"

// Multiple imports (preferred style — goimports organizes these)
import (
    "fmt"
    "os"
    "strings"
)

// Import with alias (useful when two packages have the same name)
import (
    "fmt"
    myjson "github.com/tidwall/gjson"  // alias "myjson"
    stdjson "encoding/json"             // alias "stdjson"
)

// Blank import — import for side effects only (runs init() function)
import (
    _ "github.com/lib/pq"  // PostgreSQL driver registers itself
)
```

**Using an alias:**
```go
import myjson "github.com/tidwall/gjson"

func main() {
    result := myjson.Get(`{"name":"Alice"}`, "name")  // use alias
    fmt.Println(result)  // Alice
}
```

**Blank imports** are special. You use them when a package has an `init()` function that you need to run, but you don't directly use anything from the package:
```go
import (
    _ "github.com/lib/pq"  // PostgreSQL driver for database/sql
)
// Now you can use database/sql with PostgreSQL
// The pq package registered itself with database/sql in its init()
```

**Go will not compile if you import a package and don't use it** — this is one of Go's strictest rules:
```go
import "fmt"

func main() {
    // ERROR: "fmt" imported and not used
    // The compiler won't let you leave unused imports
}
```

This keeps code clean and avoids bloat. Use `goimports` or gopls to automatically manage imports.

**The dot import (avoid this!):**
```go
import . "fmt"  // Imports all exported names directly into current scope

func main() {
    Println("no package prefix needed")  // Don't do this — confusing
}
```

### Quick Check
> 1. What happens if you import a package but never use it?
> 2. When would you use a blank import (`_ "package"`)?
> 3. What does an import alias do?

---

## 4. The main Package and main Function

The `main` package is Go's program entry point. It has special rules:

```go
package main

// init() runs BEFORE main() — used for package initialization
// You can have multiple init() functions in the same package
func init() {
    fmt.Println("init() called first!")
}

func main() {
    fmt.Println("main() called second!")
}
// Output:
// init() called first!
// main() called second!
```

**Execution order:**
1. All imported packages are initialized (their `init()` functions run)
2. Package-level variables in the current package are initialized
3. `init()` functions in the current package run (in declaration order)
4. `main()` runs

**main() restrictions:**
- Must take no parameters
- Must return nothing
- Can only exist in `package main`
- There can only be one `main()` per binary

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Exiting with a non-zero code signals failure to the OS/shell
    args := os.Args[1:]  // Arguments after the program name
    
    if len(args) == 0 {
        fmt.Fprintln(os.Stderr, "Error: no arguments provided")
        os.Exit(1)  // Non-zero = failure
    }
    
    fmt.Println("Received args:", args)
    os.Exit(0)  // Zero = success (can be omitted — main returning = exit 0)
}
```

**Testing this:**
```bash
go run main.go
# Error: no arguments provided
# (exit code 1)

go run main.go hello world
# Received args: [hello world]
# (exit code 0)

echo $?  # Print the last exit code
# 0
```

### Quick Check
> 1. What is the execution order of `init()` and `main()`?
> 2. What exit code signals success vs failure?
> 3. Can you have `func main()` in a package other than `main`?

---

## 5. Exported vs Unexported Names

This is one of Go's most important concepts. **Capitalization controls visibility:**

```go
package calculator

// Exported — visible to ALL packages that import calculator
// Capital first letter = public
func Add(a, b int) int {
    return a + b
}

type Result struct {
    Value   int     // Exported field — accessible from outside
    message string  // Unexported field — only accessible within calculator package
}

var MaxValue = 1000  // Exported variable

var minValue = 0     // Unexported variable

// Unexported function — only accessible within this package
func validateInput(n int) bool {
    return n >= minValue && n <= MaxValue
}
```

From another package:
```go
package main

import "github.com/yourname/project/calculator"

func main() {
    result := calculator.Add(3, 4)        // OK: Add is exported
    r := calculator.Result{Value: 10}     // OK: Result and Value are exported
    
    // calculator.validateInput(5)        // ERROR: unexported
    // r.message = "hello"                // ERROR: unexported field
    // calculator.minValue = 5            // ERROR: unexported variable
    
    _ = calculator.MaxValue               // OK: exported
}
```

**Why this matters:** Unexported names are your package's implementation details. Other packages don't need to know about them. This is how Go does **encapsulation** — without `private`/`public`/`protected` keywords.

```
Capital = Public (exported, visible outside)
lowercase = Private (unexported, only visible in same package)

func Add() → exported (visible anywhere)
func add() → unexported (only within same package)

type User struct {} → exported
type user struct {} → unexported

type User struct {
    Name string   → exported field
    age  int      → unexported field
}
```

### Quick Check
> 1. How does Go control whether a function or variable is accessible from other packages?
> 2. Can you access an unexported struct field from a different package?
> 3. What happens if you name a function `add` instead of `Add` and try to use it from another package?

---

## 6. Writing Your Own Package

Let's build a real example: a `temperature` package that converts between Celsius and Fahrenheit.

**Directory structure:**
```
temp-converter/
├── go.mod
├── main.go
└── temperature/
    ├── celsius.go
    └── fahrenheit.go
```

`temperature/celsius.go`:
```go
package temperature

import "fmt"

// Celsius represents a temperature in Celsius
type Celsius float64

// ToFahrenheit converts Celsius to Fahrenheit
func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

// String returns a human-readable representation
func (c Celsius) String() string {
    return fmt.Sprintf("%.1f°C", float64(c))
}
```

`temperature/fahrenheit.go`:
```go
package temperature

import "fmt"

// Fahrenheit represents a temperature in Fahrenheit
type Fahrenheit float64

// ToCelsius converts Fahrenheit to Celsius
func (f Fahrenheit) ToCelsius() Celsius {
    return Celsius((f - 32) * 5 / 9)
}

// String returns a human-readable representation
func (f Fahrenheit) String() string {
    return fmt.Sprintf("%.1f°F", float64(f))
}
```

`main.go`:
```go
package main

import (
    "fmt"
    "temp-converter/temperature"
)

func main() {
    boiling := temperature.Celsius(100)
    fmt.Println(boiling)              // 100.0°C
    fmt.Println(boiling.ToFahrenheit()) // 212.0°F
    
    body := temperature.Fahrenheit(98.6)
    fmt.Println(body)                  // 98.6°F
    fmt.Println(body.ToCelsius())      // 37.0°C
}
```

Run it:
```bash
cd temp-converter
go mod init temp-converter
go run main.go
```

This pattern — defining types in a package and attaching methods to them — is core Go. You'll do it constantly.

### Quick Check
> 1. Can you split a single package across multiple files?
> 2. What does `(c Celsius) ToFahrenheit()` mean in Go (hint: it's a method)?
> 3. How do you import your own package vs a standard library package?

---

## 7. Go Modules Deep Dive

We covered modules briefly in Chapter 02. Let's go deeper.

**go.mod explained:**
```
module github.com/yourname/temp-converter   ← module path (unique name)

go 1.23                                     ← minimum Go version required

require (
    github.com/go-chi/chi/v5 v5.0.11        ← direct dependencies
    github.com/rs/zerolog v1.31.0
)

require (
    // indirect: these are dependencies of your dependencies
    github.com/mattn/go-colorable v0.1.13 // indirect
    github.com/mattn/go-isatty v0.0.20 // indirect
    golang.org/x/sys v0.15.0 // indirect
)
```

**Version selection:**
Go uses **Minimum Version Selection (MVS)**: it always uses the minimum required version of each dependency. This makes builds reproducible. If your code works with chi v5.0.11, it will still work even if chi v5.0.12 is released.

**Updating dependencies:**
```bash
# Update a specific module to latest
go get github.com/go-chi/chi/v5@latest

# Update all modules to their latest patch version
go get -u=patch ./...

# Update all modules to their latest minor/patch version
go get -u ./...

# After updating, tidy up
go mod tidy
```

**Vendoring** (optional): Copy all dependencies into a `vendor/` directory in your project:
```bash
go mod vendor    # Creates vendor/ directory with all dependencies
go build -mod=vendor ./...   # Build using vendor/ instead of module cache
```

Vendoring is useful when you need reproducible offline builds or work in environments without internet access.

**go.sum — the security file:**
Every entry in `go.sum` is a cryptographic hash:
```
github.com/go-chi/chi/v5 v5.0.11 h1:BnpY...=  ← hash of the zip
github.com/go-chi/chi/v5 v5.0.11/go.mod h1:DS...= ← hash of go.mod
```

If anyone tampers with a dependency, `go.sum` will catch it. **Always commit `go.sum` to version control.**

### Quick Check
> 1. What is Minimum Version Selection (MVS) and why is it useful?
> 2. What does `go mod tidy` do when you remove an import from your code?
> 3. Should you commit `go.sum` to version control?

---

## 8. Build, Run, Install

**`go run`** — compiles and runs in one step (development only):
```bash
go run main.go         # Run a single file
go run .               # Run the package in current directory
go run ./cmd/server    # Run the package in ./cmd/server
```

**`go build`** — compile only:
```bash
go build .             # Compile, output binary named after directory
go build -o server .   # Compile, output binary named "server"
go build ./...         # Compile all packages (check for errors)

# Cross-compile for different OS/arch:
GOOS=linux GOARCH=amd64 go build -o server-linux .
GOOS=windows GOARCH=amd64 go build -o server.exe .
GOOS=darwin GOARCH=arm64 go build -o server-mac-m1 .
```

**`go install`** — compile and install binary to GOBIN:
```bash
# Install a tool globally
go install github.com/air-verse/air@latest

# Install your own program globally
go install .
```

**Build tags** (conditional compilation):
```go
//go:build linux

package main

// This file is ONLY compiled on Linux
func osName() string {
    return "Linux"
}
```

```bash
go build -tags debug ./...  # Include files with //go:build debug
```

**Build flags for production:**
```bash
# Smaller binary (removes debug info and DWARF tables)
go build -ldflags="-s -w" -o server .

# Set version at build time
go build -ldflags="-X main.version=1.2.3" -o server .
```

Setting version at build time:
```go
package main

var version = "dev"  // Set at build time with -ldflags

func main() {
    fmt.Println("Version:", version)
}
```

```bash
go build -ldflags="-X main.version=1.2.3" -o server .
./server
# Version: 1.2.3
```

### Quick Check
> 1. When would you use `go run` vs `go build`?
> 2. How do you cross-compile a Go program for Linux from macOS?
> 3. What does `-ldflags="-s -w"` do to the binary?

---

## Summary

- **Package**: a directory of `.go` files sharing a `package` declaration; the unit of code organization in Go
- **main package**: special package that produces an executable; must have `func main()`
- **init()**: runs before main() for initialization; multiple init() functions are allowed
- **Exported names**: capital letter = visible outside the package; lowercase = package-private
- **Imports**: bring packages in; unused imports are a compile error; aliases available
- **Go modules**: `go.mod` defines module name and dependencies; `go.sum` contains security checksums; always commit both
- **Build commands**: `go run` (develop), `go build` (compile), `go install` (install binary)
- **Cross-compilation**: `GOOS` + `GOARCH` env vars control target platform

Next chapter: variables, types, and constants — the fundamentals of every Go program.

---

## Exercises

### Easy
1. Create a package called `greet` with a function `Hello(name string) string` that returns "Hello, {name}!". Import it from `main.go` and call it.
2. In the `greet` package, add an unexported function `formatName(name string) string` that uppercases the name. Use it internally. Try to call it from `main.go` — what error do you get?
3. Add an `init()` function to your main package that prints "Application starting...". Observe when it runs relative to `main()`.

### Medium
4. Package organization: Create a project with three packages: `math/` (basic arithmetic operations), `geometry/` (area and perimeter calculations that use the math package), and a `main` package that exercises both. Make sure all function names are properly exported/unexported.
5. Build flags: Write a program with two build-tagged files: one for `linux` and one for `darwin` (macOS), each defining a function `platform() string` that returns the OS name. The main function prints the platform. Verify it compiles correctly on your OS. Use `go build -v ./...` to see what gets compiled.
6. Version injection: Modify a simple HTTP server (or just a CLI program) to have a `version`, `buildTime`, and `gitCommit` variable, all set to "dev" by default. Write a Makefile (or shell script) that uses `-ldflags` to set these to real values at build time. The program should print them when started.

### Hard
7. Package design exercise: You're designing a package for handling user authentication tokens. Design the package API: (a) What types should be exported? What fields are exported vs unexported? (b) What functions/methods should be exported vs unexported? (c) Write the package declaration, type definitions, and function signatures (no implementation needed). (d) Write a `main.go` that shows how someone would use your package. Focus on: is the API easy to use? Is it hard to misuse? Does it leak implementation details?
8. Module proxy and security: Go downloads dependencies from the internet via module proxies. (a) Run `go env GOPROXY` — what is the default? What is `proxy.golang.org`? (b) What are `GOPRIVATE` and `GONOSUMDB` and when would you use them? (c) Run `go mod verify` on a project — what does it check? (d) A developer at your company wants to use a dependency that isn't available on the public Go module proxy. What are the options? (e) What is the `replace` directive in `go.mod` and when is it useful?
