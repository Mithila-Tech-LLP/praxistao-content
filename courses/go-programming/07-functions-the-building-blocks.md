# Chapter 07: Functions — The Building Blocks

Functions are the fundamental unit of code reuse and organization in Go. Every operation you perform — from a simple calculation to handling an HTTP request — is a function. Go's functions have several powerful features that distinguish them from other languages: multiple return values, named returns, first-class functions, closures, and variadic parameters. This chapter covers all of them.

## Table of Contents

1. [Function Basics](#1-function-basics)
2. [Multiple Return Values](#2-multiple-return-values)
3. [Named Return Values](#3-named-return-values)
4. [Variadic Functions](#4-variadic-functions)
5. [Functions as Values — First-Class Functions](#5-functions-as-values--first-class-functions)
6. [Closures](#6-closures)
7. [defer — Guaranteed Cleanup](#7-defer--guaranteed-cleanup)
8. [panic and recover](#8-panic-and-recover)
9. [Function Signatures and Types](#9-function-signatures-and-types)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Function Basics

**Basic function declaration:**
```go
// func <name>(<params>) <return type> { <body> }
func greet(name string) string {
    return "Hello, " + name + "!"
}

// Multiple parameters of same type:
func add(a, b int) int {
    return a + b
}

// Multiple different parameters:
func createUser(name string, age int, email string) bool {
    // ...
    return true
}

// No return value:
func logMessage(msg string) {
    fmt.Println("[LOG]", msg)
    // no return statement needed (or use `return` with nothing)
}
```

**Calling functions:**
```go
result := greet("Alice")    // "Hello, Alice!"
sum := add(3, 4)            // 7
logMessage("Server started") // no return value
```

**Scope — functions defined at package level:**
```go
package main

// myHelper is unexported (lowercase) — only visible within this file's package
func myHelper() { ... }

// MyHelper is exported — visible to packages that import this one
func MyHelper() { ... }
```

**Recursion — a function calling itself:**
```go
func factorial(n int) int {
    if n <= 1 {
        return 1  // Base case: stop recursion
    }
    return n * factorial(n-1)  // Recursive case
}

fmt.Println(factorial(5))  // 120 (5 × 4 × 3 × 2 × 1)
```

**Recursive Fibonacci (with caching — memoization):**
```go
func fibonacci(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    if v, ok := memo[n]; ok {
        return v  // Return cached result
    }
    result := fibonacci(n-1, memo) + fibonacci(n-2, memo)
    memo[n] = result  // Cache result
    return result
}
```

### Quick Check
> 1. How do you declare a function with no return value?
> 2. How do you write a function with multiple parameters of the same type?
> 3. What is the base case in a recursive function?

---

## 2. Multiple Return Values

This is one of Go's most distinctive features. Functions can return multiple values — and this is used extensively for error handling:

```go
// Return two values: the result and an error
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil  // nil means "no error"
}

// Call it:
result, err := divide(10, 3)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)  // 3.3333...

// Ignore a value with _:
result, _ = divide(10, 2)  // Ignore error (only do this when you're certain)
```

**Standard error handling pattern:**
```go
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("opening file %s: %w", path, err)  // Wrap error
    }
    defer file.Close()
    
    data, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("reading file %s: %w", path, err)
    }
    
    // Process data...
    fmt.Println("File size:", len(data))
    return nil  // Success
}
```

**Multiple useful returns:**
```go
// Return min AND max in one pass
func minMax(numbers []int) (min, max int, err error) {
    if len(numbers) == 0 {
        return 0, 0, errors.New("empty slice")
    }
    min, max = numbers[0], numbers[0]
    for _, n := range numbers[1:] {
        if n < min { min = n }
        if n > max { max = n }
    }
    return min, max, nil
}

min, max, err := minMax([]int{3, 1, 4, 1, 5, 9})
if err != nil { ... }
fmt.Println(min, max)  // 1 9
```

**Ignoring multiple returns:**
```go
// Ignore all but last with multiple _
_, _, err := minMax(numbers)
if err != nil { ... }
```

### Quick Check
> 1. What does `return a / b, nil` return?
> 2. How do you ignore one of the two return values from a function?
> 3. What is the idiomatic Go way to handle errors from functions?

---

## 3. Named Return Values

Go allows you to name the return values. Named returns have two uses: documentation and the `naked return`.

```go
// Named return values
func divide(a, b float64) (result float64, err error) {
    if b == 0 {
        err = errors.New("cannot divide by zero")
        return  // "naked return" — returns named values as-is
    }
    result = a / b
    return  // returns result and err (err is nil zero value)
}
```

**Named returns as documentation:**
```go
// Hard to understand without names:
func getUserInfo(id int) (string, int, string, error)

// Clear with names:
func getUserInfo(id int) (name string, age int, email string, err error)
```

**Named returns in defer:**
```go
// Named returns can be modified by defer — very useful:
func riskyOperation() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    
    // If panic happens here, err is set by the defer above
    doRiskyThing()
    return nil
}
```

**Warning: naked returns in long functions are hard to read.** Use them sparingly:
```go
// BAD: long function with naked return — confusing
func complexCalc(a, b, c int) (x, y int) {
    // 50 lines of code
    // ...
    x = someValue
    y = anotherValue
    return  // What does this return? Hard to tell without reading all 50 lines
}

// BETTER: explicit return in long functions
func complexCalc(a, b, c int) (int, int) {
    // 50 lines of code
    // ...
    return someValue, anotherValue  // Clear
}
```

### Quick Check
> 1. What is a "naked return" in Go?
> 2. When are named return values useful (two purposes)?
> 3. When should you avoid naked returns?

---

## 4. Variadic Functions

A variadic function accepts a variable number of arguments. Use `...` before the type:

```go
// func name(required params..., variadic ...Type) ReturnType
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

fmt.Println(sum())           // 0
fmt.Println(sum(1))          // 1
fmt.Println(sum(1, 2, 3))    // 6
fmt.Println(sum(1, 2, 3, 4, 5))  // 15
```

**Spreading a slice into a variadic function:**
```go
numbers := []int{1, 2, 3, 4, 5}
fmt.Println(sum(numbers...))  // ... unpacks the slice into individual args
```

**Mixed params — required + variadic:**
```go
func logf(level string, format string, args ...interface{}) {
    msg := fmt.Sprintf(format, args...)
    fmt.Printf("[%s] %s\n", level, msg)
}

logf("INFO", "User %s logged in at %s", "Alice", "09:00")
logf("ERROR", "Failed after %d retries", 3)
```

**`fmt.Printf` is itself variadic** — it accepts any number of arguments.

**Variadic with any type:**
```go
func printAll(args ...interface{}) {
    for i, arg := range args {
        fmt.Printf("arg[%d] = %v (%T)\n", i, arg, arg)
    }
}

printAll(1, "hello", true, 3.14)
// arg[0] = 1 (int)
// arg[1] = hello (string)
// arg[2] = true (bool)
// arg[3] = 3.14 (float64)
```

### Quick Check
> 1. How do you declare a variadic function in Go?
> 2. How do you pass a slice to a variadic function?
> 3. What type does the variadic parameter become inside the function?

---

## 5. Functions as Values — First-Class Functions

In Go, functions are **first-class values** — you can store them in variables, pass them as arguments, and return them from functions.

**Function variable:**
```go
// Store a function in a variable
var add func(a, b int) int

add = func(a, b int) int {
    return a + b
}

fmt.Println(add(3, 4))  // 7

// Short form:
multiply := func(a, b int) int {
    return a * b
}
fmt.Println(multiply(3, 4))  // 12
```

**Passing functions as arguments:**
```go
// Higher-order function: takes a function as argument
func apply(numbers []int, fn func(int) int) []int {
    result := make([]int, len(numbers))
    for i, n := range numbers {
        result[i] = fn(n)
    }
    return result
}

doubled := apply([]int{1, 2, 3, 4}, func(n int) int { return n * 2 })
// [2, 4, 6, 8]

squared := apply([]int{1, 2, 3, 4}, func(n int) int { return n * n })
// [1, 4, 9, 16]
```

**Returning functions from functions:**
```go
// Function factory: returns a function
func makeMultiplier(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}

double := makeMultiplier(2)
triple := makeMultiplier(3)

fmt.Println(double(5))  // 10
fmt.Println(triple(5))  // 15
```

**Common functional patterns:**
```go
// Map, Filter, Reduce
func mapInts(slice []int, fn func(int) int) []int {
    result := make([]int, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

func filterInts(slice []int, pred func(int) bool) []int {
    var result []int
    for _, v := range slice {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

func reduceInts(slice []int, init int, fn func(int, int) int) int {
    result := init
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

numbers := []int{1, 2, 3, 4, 5}
doubled := mapInts(numbers, func(n int) int { return n * 2 })
evens := filterInts(numbers, func(n int) bool { return n%2 == 0 })
total := reduceInts(numbers, 0, func(acc, n int) int { return acc + n })
```

### Quick Check
> 1. What does "first-class function" mean?
> 2. How do you store a function in a variable in Go?
> 3. Write a function type for a function that takes a string and returns an error.

---

## 6. Closures

A **closure** is a function that captures variables from its surrounding scope:

```go
func makeCounter() func() int {
    count := 0  // This variable is captured by the inner function
    return func() int {
        count++  // Accesses and modifies the outer variable
        return count
    }
}

counter1 := makeCounter()
counter2 := makeCounter()  // Independent counter

fmt.Println(counter1())  // 1
fmt.Println(counter1())  // 2
fmt.Println(counter1())  // 3
fmt.Println(counter2())  // 1 (independent count)
fmt.Println(counter1())  // 4
```

**The captured variable is shared:**
```go
x := 10
add := func(n int) int {
    return x + n  // Captures x
}
fmt.Println(add(5))  // 15

x = 20  // Modifying x OUTSIDE the closure
fmt.Println(add(5))  // 25 (sees the new value!)
```

**Closures for middleware/decorators:**
```go
// Timing wrapper
func withTiming(name string, fn func()) func() {
    return func() {
        start := time.Now()
        fn()
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

slowOperation := func() {
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Done")
}

timedOp := withTiming("slow operation", slowOperation)
timedOp()
// Done
// slow operation took 100.123ms
```

**Closures for configuration:**
```go
type Config struct {
    MaxRetries int
    Timeout    time.Duration
}

func makeHTTPClient(config Config) func(url string) (string, error) {
    // The returned function captures config
    return func(url string) (string, error) {
        client := &http.Client{Timeout: config.Timeout}
        for attempt := 0; attempt <= config.MaxRetries; attempt++ {
            resp, err := client.Get(url)
            if err == nil {
                defer resp.Body.Close()
                body, _ := io.ReadAll(resp.Body)
                return string(body), nil
            }
            if attempt < config.MaxRetries {
                time.Sleep(time.Second)
            }
        }
        return "", errors.New("max retries exceeded")
    }
}

fetchWithRetry := makeHTTPClient(Config{MaxRetries: 3, Timeout: 5 * time.Second})
body, err := fetchWithRetry("https://api.example.com/data")
```

### Quick Check
> 1. What is a closure?
> 2. If the outer variable changes after the closure is created, does the closure see the new value?
> 3. What problem do closures solve that would otherwise require a struct?

---

## 7. defer — Guaranteed Cleanup

`defer` schedules a function call to run when the surrounding function returns — regardless of how it returns (normally, via return, or via panic):

```go
func openFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()  // Will run when openFile returns, guaranteed
    
    // Even if any of these return early with an error,
    // file.Close() will still run
    data, err := io.ReadAll(file)
    if err != nil {
        return err
    }
    
    return processData(data)
}
```

**defer executes in LIFO order** (last in, first out):
```go
func example() {
    defer fmt.Println("first deferred")   // runs third
    defer fmt.Println("second deferred")  // runs second
    defer fmt.Println("third deferred")   // runs first
    fmt.Println("function body")
}
// Output:
// function body
// third deferred
// second deferred
// first deferred
```

**defer with arguments — arguments are evaluated immediately:**
```go
x := 10
defer fmt.Println("x =", x)  // x is evaluated NOW (10), not when defer runs
x = 20
// Output when function returns: "x = 10" (not 20!)
```

**Common defer uses:**
```go
// 1. Closing files, connections, network resources
file, _ := os.Open("file.txt")
defer file.Close()

// 2. Releasing mutex locks
mu.Lock()
defer mu.Unlock()

// 3. Closing channels
ch := make(chan int)
defer close(ch)

// 4. Measuring function execution time
func processRequest(id string) {
    start := time.Now()
    defer func() {
        fmt.Printf("processRequest(%s) took %v\n", id, time.Since(start))
    }()
    // ... actual work
}

// 5. Database transactions
tx, err := db.Begin()
if err != nil { return err }
defer func() {
    if err != nil {
        tx.Rollback()
    } else {
        tx.Commit()
    }
}()
```

### Quick Check
> 1. When exactly does a `defer`-ed function run?
> 2. In what order do multiple `defer` statements execute?
> 3. Are the arguments to a `defer` evaluated immediately or when the defer runs?

---

## 8. panic and recover

**panic** stops normal execution. **recover** can catch a panic and continue execution. These are for *exceptional* situations, not normal error handling:

```go
func mustPositive(n int) {
    if n < 0 {
        panic(fmt.Sprintf("expected positive, got %d", n))
    }
}

// Without recover, panic terminates the program:
mustPositive(-1)  // Program crashes with panic message
```

**recover — catching panics:**
```go
// recover ONLY works inside a defer
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered panic: %v", r)
        }
    }()
    return a / b, nil  // This will panic if b == 0
}

result, err := safeDiv(10, 0)
if err != nil {
    fmt.Println("Error:", err)  // Error: recovered panic: runtime error: integer divide by zero
}
```

**When to use panic:**
- Programmer errors that should never happen (nil pointer you never expected, index out of bounds you should have prevented)
- During initialization, when the program cannot proceed (failed to connect to database on startup, missing required config)
- Never use panic for regular error handling (use error returns instead)

```go
// Package initialization — panic is appropriate here
func init() {
    if os.Getenv("DATABASE_URL") == "" {
        panic("DATABASE_URL environment variable is required")
    }
}

// HTTP handler — NEVER panic here, always return errors properly
func handleGetUser(w http.ResponseWriter, r *http.Request) {
    user, err := db.GetUser(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, "user not found", http.StatusNotFound)
        return  // Handle gracefully
    }
    json.NewEncoder(w).Encode(user)
}
```

### Quick Check
> 1. When should you use panic vs returning an error?
> 2. Where must `recover()` be called to work?
> 3. What does `recover()` return when there's no active panic?

---

## 9. Function Signatures and Types

Every function has a type defined by its parameter types and return types:

```go
// Function type declarations
type HandlerFunc func(http.ResponseWriter, *http.Request)
type PredicateFunc func(int) bool
type TransformFunc func(string) string

// Using function types as parameters:
func middleware(handler HandlerFunc) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s %s", r.Method, r.URL.Path)
        handler(w, r)
    }
}

// Using function types in structs:
type Router struct {
    routes map[string]HandlerFunc
}

func (r *Router) Handle(path string, handler HandlerFunc) {
    r.routes[path] = handler
}
```

**Method values — binding a method to a specific receiver:**
```go
type Calculator struct {
    memory float64
}

func (c *Calculator) Add(a, b float64) float64 {
    return a + b
}

calc := &Calculator{}
addFn := calc.Add  // Method value: bound to calc

result := addFn(3, 4)  // 7 — no need to pass calc again
```

**Interface satisfaction through function types:**
```go
// http.HandlerFunc implements http.Handler
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

// This lets you pass a plain function where an interface is expected
http.Handle("/path", http.HandlerFunc(myFunc))
```

### Quick Check
> 1. What is a function type in Go?
> 2. What is a "method value"?
> 3. How does `http.HandlerFunc` let a plain function satisfy the `http.Handler` interface?

---

## Summary

- **Function basics**: `func name(params) return-type { body }`; recursion; exported with capital letter
- **Multiple returns**: Functions can return multiple values; idiomatic for (value, error) pairs; use `_` to discard
- **Named returns**: Document return values; enable naked return; modify-via-defer pattern
- **Variadic**: `...T` parameter accepts zero or more args; spread slice with `slice...`
- **First-class**: Functions are values; store in variables, pass as args, return from functions
- **Closures**: Functions that capture outer scope variables; used for factories, middleware, callbacks
- **defer**: Schedule cleanup; runs on function exit (LIFO); args evaluated immediately
- **panic/recover**: For exceptional situations, not normal errors; recover only in defer

Next chapter: Arrays and Slices — Go's most important collection types.

---

## Exercises

### Easy
1. Write a recursive function `power(base, exp int) int` that computes base^exp. Handle the case where exp == 0.
2. Write a function `filter(items []string, fn func(string) bool) []string` that returns only items for which fn returns true.
3. Write a function that uses `defer` to print "function done" when it exits, regardless of whether it returns normally or panics. Test it both ways.

### Medium
4. Function pipeline: Implement a `pipeline` function: `func pipeline(input string, fns ...func(string) string) string` that applies each function in sequence to the string. Test with: lowercase, trim spaces, replace spaces with dashes. Result of applying to " Hello World " should be "hello-world".
5. Memoization: Write a `memoize` function that takes a function `func(int) int` and returns a cached version. The cached version calls the original only on the first call with a given input, and returns the cached result on subsequent calls. Test with Fibonacci numbers — measure the speedup.
6. Function retry: Write `func Retry(attempts int, delay time.Duration, fn func() error) error` that calls fn up to `attempts` times, waiting `delay` between attempts, and returns nil if any attempt succeeds. Add exponential backoff: `Retry(attempts int, initialDelay time.Duration, multiplier float64, fn func() error) error`.

### Hard
7. Middleware chain: Build a simple HTTP middleware system. Define `type Middleware func(http.HandlerFunc) http.HandlerFunc`. Implement three middlewares: `Logger` (logs method+path+duration), `Auth` (checks Authorization header, returns 401 if missing), `RateLimiter` (allows max N requests per minute per IP, returns 429 if exceeded). Write a `Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc` that applies them in order. Write a test HTTP server using all three middlewares.
8. Lazy evaluation with closures: Implement a `Lazy[T]` type that defers computation until the value is needed. `func NewLazy[T any](compute func() T) *Lazy[T]` — stores the computation but doesn't run it. `func (l *Lazy[T]) Get() T` — runs computation on first call, caches result, returns cached on subsequent calls. Thread-safe using `sync.Once`. Write tests showing: (a) compute function is called exactly once even with concurrent access, (b) result is consistent across goroutines, (c) works with different types (int, string, struct).
