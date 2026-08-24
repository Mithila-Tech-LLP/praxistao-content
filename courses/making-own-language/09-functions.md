# Chapter 09: Functions — The Building Blocks of Programs

> "Functions are the nouns and verbs of programming — they name things and describe what can be done with them."
> — Anonymous

---

## Chapter Overview

If there is one concept more fundamental to programming than any other, it is the **function**. Functions are named, reusable blocks of code. Instead of writing the same ten lines over and over in different places, you write those ten lines once, give them a name, and then *call* that name whenever you need them. This sounds simple — and in one sense it is — but the implications are profound. Functions are how we manage complexity, organize thought, create abstractions, and build large systems from small, understandable parts.

A program without functions would be an impossibly long, flat list of instructions, with no structure, no reuse, and no way to understand what any part of it does without reading all of it. Functions give programs *shape*. They allow you to say "I don't care *how* this works, I just want to call `send_email(to, subject, body)` and trust that it works." That trust — that ability to use something without understanding its internals — is the foundation of software engineering.

In this chapter, we explore functions from both a user perspective (writing Astra code) and a compiler perspective (representing functions as AST nodes). We will cover function anatomy, the call stack, parameters, return values, closures, first-class functions, anonymous functions, variadic functions, recursion basics, and the design choice not to have function overloading. By the end, you will have a deep understanding of what functions *are*, and our compiler will have the AST nodes it needs to represent them.

---

## What We're Building

In our Astra compiler, we need AST nodes for three things related to functions:
1. **Function declarations** — defining a function (its name, parameters, return type, body)
2. **Function calls** — invoking a function with arguments
3. **Parameters** — the named inputs a function accepts

These nodes live in `ast/declarations.go` and `ast/expressions.go`. They will be used by the parser to build the AST, the type checker to verify that argument types match parameter types, and the code generator to emit the function prologue and epilogue in machine code.

---

## Table of Contents

1. What is a Function? — The Named Recipe Analogy
2. Function Anatomy — Name, Parameters, Body, Return Value
3. Functions in Go — Syntax and Conventions
4. Functions in Astra — The `fn` Keyword
5. The Call Stack — What Happens When You Call a Function
6. Parameters — Pass by Value vs Pass by Reference
7. Return Values — Single, Multiple, and Named
8. Closures — Functions That Remember Their Environment
9. First-Class Functions — Passing Functions as Arguments
10. Anonymous Functions and Lambdas
11. Recursion Basics — Base Case and Recursive Case
12. Variadic Functions — Accepting Any Number of Arguments
13. No Function Overloading — And Why That's OK
14. Pure vs Impure Functions — Side Effects and Compiler Optimizations
15. The Single Responsibility Principle
16. Astra Build Milestone: AST Nodes for Functions
17. Exercises
18. Summary

---

## 1. What is a Function? — The Named Recipe Analogy

A **function** is a named, self-contained piece of code that performs a specific task. Think of it like a recipe in a cookbook.

```
Recipe: "Boil Water"
  Inputs: a pot, water, a stove
  Steps: 1. Fill the pot with water
         2. Place the pot on the stove
         3. Turn the stove on high
         4. Wait until bubbles appear
  Output: boiling water

Function: boil_water(pot, water, stove) -> boiling_water
  Parameters: pot, water, stove   (the inputs)
  Body:       the steps            (what to do)
  Return:     boiling_water        (the output)
```

The key insight: once you write the recipe, you can use it in any other recipe that needs boiling water. You don't re-explain how to boil water in the "make pasta" recipe — you just say "boil water (see page 12)." Functions work exactly the same way: write once, use many times.

**Why functions matter:**
- **Reuse:** write code once, call it many times
- **Abstraction:** hide complexity behind a simple name
- **Testing:** test a function in isolation
- **Readability:** `calculateTax(income)` is clearer than 30 lines of tax logic inline
- **Maintenance:** change the logic in one place, fixed everywhere

---

## 2. Function Anatomy — Name, Parameters, Body, Return Value

Every function has four parts:

```
              name
               │
               ↓
         ┌──────────┐
fn add(  │   add    │  ) -> int {
         └──────────┘
              │
         parameters
               │
    ┌──────────┴───────────┐
    │  a: int,  b: int     │
    └──────────────────────┘
                                    return type
                                         │
                               ┌─────────┴────────┐
fn add(a: int, b: int)  ->    │       int        │  {
                               └──────────────────┘

  body
    │
    ↓
    return a + b
}
```

```
┌─────────────────────────────────────────────────────────┐
│  fn  name  (  parameters  )  ->  return_type  {  body  }│
│  ↑   ↑        ↑               ↑               ↑        │
│  │   │        │               │               │        │
│  │   │        │               │          the code      │
│ keyword │   inputs         output type   that runs     │
│       function                                         │
│       identifier                                       │
└─────────────────────────────────────────────────────────┘
```

All four parts are optional in various combinations:
- A function can have **no parameters** (`fn greet() { print("Hello") }`)
- A function can return **nothing** (omit `->` and the return type)
- A function's body can be **empty** (though that's usually a placeholder)

---

## 3. Functions in Go — Syntax and Conventions

```go
package main

import "fmt"

// Basic function: no parameters, no return value
func greet() {
    fmt.Println("Hello!")
}

// Function with parameters and a return value
func add(a int, b int) int {
    return a + b
}

// Shorthand when parameters have the same type: a, b int
func multiply(a, b int) int {
    return a * b
}

// Multiple return values — a Go superpower!
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}

// Named return values
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums[1:] {
        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }
    return  // bare return: returns the named values min and max
}

func main() {
    greet()                       // Hello!
    fmt.Println(add(3, 4))        // 7
    fmt.Println(multiply(6, 7))   // 42

    result, err := divide(10, 3)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("10 / 3 = %.2f\n", result)
    }

    lo, hi := minMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
    fmt.Printf("Min: %d, Max: %d\n", lo, hi)
}
```

**Go naming conventions:**
- Function names use camelCase: `calculateTax`, `sendEmail`, `parseJSON`
- Exported functions (usable from other packages) start with uppercase: `Println`, `NewLexer`
- Unexported functions start with lowercase: `nextToken`, `skipWhitespace`

---

## 4. Functions in Astra — The `fn` Keyword

Astra uses `fn` (short for function) as its keyword, which is shorter and faster to type than `func` (Go) or `function` (JavaScript):

```astra
// Basic function
fn greet() {
    print("Hello!")
}

// Function with parameters and return type
fn add(a: int, b: int) -> int {
    return a + b
}

// Note: Astra uses explicit type annotations on parameters
// a: int   means "parameter named 'a' of type int"

// Function returning nothing (void)
fn log(message: string) {
    print("[LOG] " + message)
}

// Function calling other functions
fn greet_person(name: string) {
    let message = "Hello, " + name + "!"
    log(message)
}

// Real example: area of a rectangle
fn rectangle_area(width: float, height: float) -> float {
    return width * height
}

// Using the functions
fn main() {
    greet()
    print(add(3, 4))             // 7
    greet_person("Alice")         // [LOG] Hello, Alice!
    print(rectangle_area(5.0, 3.0))  // 15.0
}
```

**Comparison of function syntax across languages:**

```
Go:
  func add(a int, b int) int { return a + b }

Java:
  public static int add(int a, int b) { return a + b; }

JavaScript:
  function add(a, b) { return a + b; }  // no types!

Python:
  def add(a, b):                         // no types (unless using hints)
      return a + b

Rust:
  fn add(a: i32, b: i32) -> i32 { a + b }

Astra:
  fn add(a: int, b: int) -> int { return a + b }
```

Astra's syntax is closest to Rust but uses `int` instead of `i32` for simplicity.

---

## 5. The Call Stack — What Happens When You Call a Function

When you call a function, the computer does not just "jump" to that function's code. It creates a **stack frame** — a block of memory that stores:
- The function's local variables
- The function's parameters (copied in)
- The **return address** — where to go back to when the function finishes

These stack frames are placed on a data structure called the **call stack**. When a function is called, a new frame is pushed on top. When a function returns, its frame is popped off, and execution continues from the return address.

```
Example: main calls greet, which calls print

STACK at start (only main):
┌────────────────────────┐  ← TOP of stack
│  main's frame          │
│  (local vars of main)  │
└────────────────────────┘

STACK after main calls greet():
┌────────────────────────┐  ← TOP of stack
│  greet's frame         │
│  return address: main  │
└────────────────────────┘
┌────────────────────────┐
│  main's frame          │
└────────────────────────┘

STACK after greet calls print("Hello"):
┌────────────────────────┐  ← TOP of stack
│  print's frame         │
│  param: "Hello"        │
│  return address: greet │
└────────────────────────┘
┌────────────────────────┐
│  greet's frame         │
│  return address: main  │
└────────────────────────┘
┌────────────────────────┐
│  main's frame          │
└────────────────────────┘

print finishes → pop print's frame, return to greet
greet finishes → pop greet's frame, return to main
```

**Stack overflow:** The call stack has a finite size. If functions call themselves too many times (deep recursion without a base case), the stack runs out of space and the program crashes with a "stack overflow" error. This is why the famous website is called stackoverflow.com — programmers frequently hit this error.

```go
// This will crash with a stack overflow:
func infinite() {
    infinite()  // calls itself forever, no base case
}
```

---

## 6. Parameters — Pass by Value vs Pass by Reference

When you call a function with arguments, how does the function receive those values?

**Pass by value:** a *copy* of the value is made and given to the function. Changes inside the function do not affect the original.

```go
// Go — integers are passed by value
func doubleIt(x int) {
    x = x * 2     // modifies the local copy
    fmt.Println(x) // prints the doubled value
}

n := 5
doubleIt(n)
fmt.Println(n) // still 5 — the original is unchanged
```

**Pass by reference (using pointers):** a *pointer to* the value is passed. Changes inside the function affect the original.

```go
// Go — use a pointer to modify the original
func doubleInPlace(x *int) {
    *x = *x * 2   // modifies the value AT the address
}

n := 5
doubleInPlace(&n)  // & takes the address of n
fmt.Println(n)     // now 10 — the original WAS changed
```

**In Astra:**

```astra
// Astra passes simple types by value (int, float, bool, string)
fn double_it(x: int) {
    x = x * 2
    print(x)     // prints doubled value
}

let n = 5
double_it(n)
print(n)         // still 5

// For mutable reference, Astra uses '&mut':
fn double_in_place(x: &mut int) {
    x = x * 2
}

let mut n = 5      // 'mut' marks a variable as mutable
double_in_place(&mut n)
print(n)           // now 10
```

**Reference types** (slices, maps, objects) in Go are already reference-like — you can modify a slice or map passed to a function and the changes are visible to the caller. This can be surprising for beginners.

```go
// Slices are reference types in Go
func appendOne(s []int) {
    s = append(s, 1)   // this does NOT affect the caller's slice
    // (append may create a new backing array)
}

func addToSlice(s []int) {
    s[0] = 999   // this DOES affect the caller (modifies in-place)
}
```

---

## 7. Return Values — Single, Multiple, and Named

**Single return value (most common):**

```astra
fn square(n: int) -> int {
    return n * n
}
```

**No return value (called "void" in C/Java):**

```astra
fn log(msg: string) {
    print("[LOG] " + msg)
    // no return statement needed
}
```

**Multiple return values in Go** (Astra uses Result types for error handling, covered later):

```go
// Go — return multiple values
func parseAge(s string) (int, error) {
    age, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("invalid age: %s", s)
    }
    if age < 0 || age > 150 {
        return 0, fmt.Errorf("age out of range: %d", age)
    }
    return age, nil
}

// Caller unpacks both values
age, err := parseAge("25")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Age:", age)
```

**In Astra, error handling uses the `Result` type:**

```astra
fn parse_age(s: string) -> Result<int, string> {
    let age = s.parse_int()
    match age {
        Ok(n) => {
            if n < 0 || n > 150 {
                return Err("age out of range")
            }
            return Ok(n)
        }
        Err(e) => { return Err("invalid age: " + s) }
    }
}

// Caller uses match or ? operator
let age = parse_age("25")?   // ? propagates the error if it failed
print("Age: " + age)
```

---

## 8. Closures — Functions That Remember Their Environment

A **closure** is a function that *captures* variables from the surrounding scope — variables that exist outside the function's own body. The function "closes over" those variables, keeping them alive even after the outer scope has ended.

```go
// Go closure example
func makeCounter() func() int {
    count := 0               // 'count' lives in makeCounter's scope
    return func() int {      // this anonymous function captures 'count'
        count++
        return count
    }
}

counter := makeCounter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
fmt.Println(counter())  // 3
// Each call increments the SAME 'count' variable
```

```astra
// Astra closure example
fn make_counter() -> fn() -> int {
    let count = 0
    return fn() -> int {
        count = count + 1    // captures 'count' from outer scope
        return count
    }
}

let counter = make_counter()
print(counter())   // 1
print(counter())   // 2
print(counter())   // 3
```

**Real-world use case: event handlers**

```astra
fn setup_button(label: string) -> fn() {
    return fn() {
        print("Button clicked: " + label)   // captures 'label'
    }
}

let ok_button = setup_button("OK")
let cancel_button = setup_button("Cancel")

ok_button()       // Button clicked: OK
cancel_button()   // Button clicked: Cancel
```

Closures are powerful because they let you create customized functions on the fly. Each `setup_button` call creates a *new* function that remembers its own `label`.

---

## 9. First-Class Functions — Passing Functions as Arguments

In Go and Astra, functions are **first-class values** — you can store them in variables, pass them to other functions, and return them from functions. This is what makes closures and callbacks possible.

```go
// A function that takes another function as an argument
func applyToAll(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = f(n)
    }
    return result
}

double := func(x int) int { return x * 2 }
square := func(x int) int { return x * x }

nums := []int{1, 2, 3, 4, 5}
fmt.Println(applyToAll(nums, double))  // [2 4 6 8 10]
fmt.Println(applyToAll(nums, square))  // [1 4 9 16 25]
```

```astra
// Astra: map, filter, reduce — the holy trinity of functional programming
fn map(list: List<int>, transform: fn(int) -> int) -> List<int> {
    let result: List<int> = []
    for item in list {
        result.push(transform(item))
    }
    return result
}

fn filter(list: List<int>, predicate: fn(int) -> bool) -> List<int> {
    let result: List<int> = []
    for item in list {
        if predicate(item) {
            result.push(item)
        }
    }
    return result
}

let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let evens = filter(numbers, fn(n: int) -> bool { return n % 2 == 0 })
let doubled = map(evens, fn(n: int) -> int { return n * 2 })
print(doubled)  // [4, 8, 12, 16, 20]
```

**Function types:**

```
In Go, a function type looks like:
  func(int, int) int          — takes two ints, returns one int
  func(string) (bool, error)  — takes a string, returns bool and error
  func()                      — takes nothing, returns nothing

In Astra:
  fn(int, int) -> int
  fn(string) -> bool
  fn()
```

---

## 10. Anonymous Functions and Lambdas

An **anonymous function** (also called a **lambda**) is a function without a name. You define it inline where you need it, often as an argument to another function.

```go
// Go anonymous function
nums := []int{5, 2, 8, 1, 9, 3}
sort.Slice(nums, func(i, j int) bool {
    return nums[i] < nums[j]   // anonymous function as a sort comparator
})
fmt.Println(nums) // [1 2 3 5 8 9]

// Immediately invoked anonymous function (IIFE)
result := func(a, b int) int {
    return a + b
}(3, 4)
fmt.Println(result) // 7
```

```astra
// Astra anonymous function in a web server handler
server.get("/", fn(req: http.Request, res: http.Response) {
    res.send("Hello, World!")
})

// Anonymous function stored in a variable
let square = fn(x: int) -> int { return x * x }
print(square(5))  // 25

// Anonymous function as argument
let numbers = [1, 2, 3, 4, 5]
let doubled = numbers.map(fn(n: int) -> int { return n * 2 })
print(doubled)  // [2, 4, 6, 8, 10]
```

---

## 11. Recursion Basics — Base Case and Recursive Case

**Recursion** is when a function calls itself. Every recursive function needs two things:

1. **Base case:** a condition where the function does NOT call itself (stops the recursion)
2. **Recursive case:** where the function calls itself with a simpler/smaller version of the problem

**Classic example: factorial**

```
Factorial definition:
  0! = 1          (base case)
  n! = n × (n-1)! (recursive case)

factorial(4) = 4 × factorial(3)
             = 4 × 3 × factorial(2)
             = 4 × 3 × 2 × factorial(1)
             = 4 × 3 × 2 × 1 × factorial(0)
             = 4 × 3 × 2 × 1 × 1
             = 24
```

```astra
fn factorial(n: int) -> int {
    if n <= 0 {
        return 1     // base case
    }
    return n * factorial(n - 1)   // recursive case
}

print(factorial(5))   // 120
print(factorial(0))   // 1
```

**Classic example: Fibonacci sequence**

```astra
// Fibonacci: fib(0)=0, fib(1)=1, fib(n)=fib(n-1)+fib(n-2)
fn fib(n: int) -> int {
    if n <= 0 { return 0 }   // base case 1
    if n == 1 { return 1 }   // base case 2
    return fib(n - 1) + fib(n - 2)   // recursive case
}

// WARNING: this is exponentially slow for large n (Chapter 19)
// For now, it illustrates the concept clearly
```

**The most important rule of recursion:** always make progress toward the base case. If each recursive call does not move closer to the base case, you get infinite recursion (and a stack overflow).

---

## 12. Variadic Functions — Accepting Any Number of Arguments

A **variadic function** accepts any number of arguments of a given type. This is how `print` works — you can pass it one string or ten strings.

```go
// Go variadic function
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

fmt.Println(sum(1, 2, 3))       // 6
fmt.Println(sum(10, 20))        // 30
fmt.Println(sum(1, 2, 3, 4, 5)) // 15

// Spreading a slice into a variadic call
nums := []int{1, 2, 3, 4}
fmt.Println(sum(nums...))  // 10 — the ... spreads the slice
```

```astra
// Astra variadic function
fn sum(nums: ...int) -> int {
    let total = 0
    for n in nums {
        total = total + n
    }
    return total
}

print(sum(1, 2, 3))          // 6
print(sum(10, 20, 30, 40))   // 100

// The built-in print is variadic:
print("Name:", "Alice", "Age:", 30)   // Name: Alice Age: 30
```

Inside the function, the variadic parameter `nums` behaves like a list. You can iterate over it, check its length, and access individual elements.

---

## 13. No Function Overloading — And Why That's OK

**Function overloading** is when a language lets you define multiple functions with the same name but different parameter types. Java and C++ support this:

```java
// Java overloading
void print(int n) { System.out.println(n); }
void print(String s) { System.out.println(s); }
void print(double d) { System.out.println(d); }
```

Go does not have function overloading. Astra does not either. Here's why:

1. **Readability:** when you see `sort(list)`, you don't know which sort function is being called without knowing the type of `list`. With unique names, `sort_ints(list)` and `sort_strings(list)` are unambiguous.

2. **Compiler simplicity:** overloading requires the compiler to resolve which function to call at every call site, based on argument types. This is called **overload resolution** and it is surprisingly complex (especially when implicit type conversions exist).

3. **Go's experience:** Go has been wildly successful without overloading. The developers of Go deliberately left it out, and Go programmers have found that it leads to clearer code.

**The alternative in Astra: generics (Chapter 32)**

```astra
// Astra generics let you write one function for multiple types
fn find_max<T: Comparable>(items: List<T>) -> T {
    let max = items[0]
    for item in items {
        if item > max {
            max = item
        }
    }
    return max
}

find_max([3, 1, 4, 1, 5, 9])          // works for ints
find_max(["banana", "apple", "cherry"]) // works for strings
```

---

## 14. Pure vs Impure Functions — Side Effects and Compiler Optimizations

A **pure function** always returns the same output for the same input, and has no **side effects** (it does not modify any external state — no global variables, no I/O, no network calls).

```astra
// Pure function — no side effects, same input always gives same output
fn square(n: int) -> int {
    return n * n
}

// Impure function — has a side effect (printing to screen)
fn log_and_square(n: int) -> int {
    print("Computing square of " + n)   // SIDE EFFECT: I/O
    return n * n
}

// Impure function — modifies external state (global variable)
let call_count = 0
fn increment_counter() {
    call_count = call_count + 1   // SIDE EFFECT: modifies global state
}
```

**Why pure functions matter for compilers:**

Pure functions can be **memoized** (the compiler remembers the result of a call and reuses it instead of calling the function again). They can be **parallelized** (multiple calls can run simultaneously because they don't share state). They can be **inlined** (the compiler can replace a function call with the function body directly). All of these are optimizations our Astra compiler will eventually implement.

```
Pure function:     square(3) → always 9
Memoization:       first call computes 9, subsequent calls return 9 immediately
                   → zero computation cost after the first call

Impure function:   log_and_square(3) → 9, but also prints "Computing square of 3"
Memoization:       CAN'T memoize — calling it a second time would skip the print
                   → compiler must call it every time
```

---

## 15. The Single Responsibility Principle

The **Single Responsibility Principle** (SRP) says: a function should do *one thing* and do it well. It should have one reason to change.

**Violation of SRP (one function does too much):**

```astra
fn process_user_data(user_id: int) {
    // Reads user from database
    let user = db.query("SELECT * FROM users WHERE id = " + user_id)

    // Validates the user
    if user.age < 0 || user.age > 150 {
        print("Invalid age")
        return
    }

    // Sends a welcome email
    let email = smtp.compose("welcome@astra.com", user.email, "Welcome!")
    smtp.send(email)

    // Updates the database
    db.execute("UPDATE users SET welcomed = true WHERE id = " + user_id)

    // Logs everything
    print("User " + user.name + " welcomed at " + now())
}
```

**Better design (each function does one thing):**

```astra
fn fetch_user(user_id: int) -> User {
    return db.query("SELECT * FROM users WHERE id = " + user_id)
}

fn validate_user(user: User) -> bool {
    return user.age >= 0 && user.age <= 150
}

fn send_welcome_email(user: User) {
    let email = smtp.compose("welcome@astra.com", user.email, "Welcome!")
    smtp.send(email)
}

fn mark_user_welcomed(user_id: int) {
    db.execute("UPDATE users SET welcomed = true WHERE id = " + user_id)
}

// The orchestrating function is now clean and easy to read
fn welcome_new_user(user_id: int) {
    let user = fetch_user(user_id)
    if !validate_user(user) {
        print("Invalid user data")
        return
    }
    send_welcome_email(user)
    mark_user_welcomed(user_id)
    print("User " + user.name + " welcomed at " + now())
}
```

The second version is easier to test (you can test each function independently), easier to modify (changing how emails are sent only affects `send_welcome_email`), and easier to understand (each function name tells you exactly what it does).

---

## 16. Astra Build Milestone: AST Nodes for Functions

Create or extend `ast/declarations.go` and `ast/expressions.go`:

```go
// ast/declarations.go
// AST nodes for function declarations in Astra.

package ast

import "strings"

// ------------------------------------------------------------
// Type representation
// ------------------------------------------------------------

// Type represents a type annotation in Astra source code.
// Examples: int, string, float, List<int>, fn(int) -> bool
type Type interface {
    typeName() string  // the string name of the type
}

// NamedType is a simple named type like int, string, bool
type NamedType struct {
    Name string
}

func (t *NamedType) typeName() string { return t.Name }

// FunctionType represents a function type like fn(int, string) -> bool
type FunctionType struct {
    ParameterTypes []Type
    ReturnType     Type  // nil if void
}

func (t *FunctionType) typeName() string {
    params := make([]string, len(t.ParameterTypes))
    for i, p := range t.ParameterTypes {
        params[i] = p.typeName()
    }
    if t.ReturnType == nil {
        return "fn(" + strings.Join(params, ", ") + ")"
    }
    return "fn(" + strings.Join(params, ", ") + ") -> " + t.ReturnType.typeName()
}

// GenericType represents a generic type like List<int> or Map<string, int>
type GenericType struct {
    Base     string  // "List", "Map", etc.
    TypeArgs []Type  // the type arguments
}

func (t *GenericType) typeName() string {
    args := make([]string, len(t.TypeArgs))
    for i, a := range t.TypeArgs {
        args[i] = a.typeName()
    }
    return t.Base + "<" + strings.Join(args, ", ") + ">"
}

// ------------------------------------------------------------
// Parameter
// ------------------------------------------------------------

// Parameter represents a single function parameter.
//
// Example Astra source:
//   fn add(a: int, b: int) -> int
//           ^^^^^^  ^^^^^^
//           param   param
//
// Each Parameter has a name and a type.
type Parameter struct {
    Name    string  // the parameter name ("a", "b", "name", etc.)
    Type    Type    // the parameter's type
    Mutable bool    // true if the parameter is declared as 'mut'
}

func (p *Parameter) String() string {
    prefix := ""
    if p.Mutable {
        prefix = "mut "
    }
    return prefix + p.Name + ": " + p.Type.typeName()
}

// ------------------------------------------------------------
// FunctionDeclaration
// ------------------------------------------------------------

// FunctionDeclaration represents a top-level function definition in Astra.
//
// Example Astra source:
//   fn add(a: int, b: int) -> int {
//       return a + b
//   }
//
// AST representation:
//   FunctionDeclaration{
//       Name:       "add",
//       Parameters: [Parameter{"a", int}, Parameter{"b", int}],
//       ReturnType: NamedType{"int"},
//       Body:       BlockStatement{[ReturnStatement{BinaryExpr{a, +, b}}]},
//   }
type FunctionDeclaration struct {
    Name       string      // the function's name
    Parameters []Parameter // the parameter list (may be empty)
    ReturnType Type        // nil if the function returns nothing (void)
    Body       *BlockStatement
    IsVariadic bool        // true if the last parameter is variadic (nums ...int)
}

func (fd *FunctionDeclaration) statementNode()       {}
func (fd *FunctionDeclaration) TokenLiteral() string { return "fn" }
func (fd *FunctionDeclaration) String() string {
    params := make([]string, len(fd.Parameters))
    for i, p := range fd.Parameters {
        params[i] = p.String()
    }

    result := "fn " + fd.Name + "(" + strings.Join(params, ", ") + ")"
    if fd.ReturnType != nil {
        result += " -> " + fd.ReturnType.typeName()
    }
    result += " " + fd.Body.String()
    return result
}

// IsVoid returns true if the function does not return a value.
func (fd *FunctionDeclaration) IsVoid() bool {
    return fd.ReturnType == nil
}

// ------------------------------------------------------------
// ReturnStatement
// ------------------------------------------------------------

// ReturnStatement represents a return statement inside a function.
//
// Example Astra source:
//   return a + b
//   return         // bare return in void function
type ReturnStatement struct {
    Value Expression  // nil if bare return (void function)
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return "return" }
func (rs *ReturnStatement) String() string {
    if rs.Value == nil {
        return "return"
    }
    return "return " + rs.Value.String()
}
```

```go
// ast/expressions.go
// AST nodes for function calls and anonymous functions.

package ast

import "strings"

// ------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------

// CallExpression represents a function call in Astra.
//
// Example Astra source:
//   add(3, 4)
//   print("Hello")
//   server.get("/", handler)
//
// AST representation for add(3, 4):
//   CallExpression{
//       Function:  Identifier{"add"},
//       Arguments: [IntLiteral{3}, IntLiteral{4}],
//   }
type CallExpression struct {
    Function  Expression   // the function being called (Identifier or MemberExpr)
    Arguments []Expression // the list of arguments
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return "(" }
func (ce *CallExpression) String() string {
    args := make([]string, len(ce.Arguments))
    for i, a := range ce.Arguments {
        args[i] = a.String()
    }
    return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

// ------------------------------------------------------------
// FunctionLiteral (anonymous function / lambda)
// ------------------------------------------------------------

// FunctionLiteral represents an anonymous function expression.
//
// Example Astra source:
//   fn(x: int) -> int { return x * 2 }
//   fn(req: http.Request, res: http.Response) { res.send("Hi") }
//
// Function literals capture variables from their enclosing scope (closures).
type FunctionLiteral struct {
    Parameters []Parameter
    ReturnType Type // nil if void
    Body       *BlockStatement
    Captures   []string // variables captured from outer scope (filled in by semantic analysis)
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return "fn" }
func (fl *FunctionLiteral) String() string {
    params := make([]string, len(fl.Parameters))
    for i, p := range fl.Parameters {
        params[i] = p.String()
    }
    result := "fn(" + strings.Join(params, ", ") + ")"
    if fl.ReturnType != nil {
        result += " -> " + fl.ReturnType.typeName()
    }
    result += " " + fl.Body.String()
    return result
}

// ------------------------------------------------------------
// Identifier (used for variable and function references)
// ------------------------------------------------------------

// Identifier represents a name reference (a variable or function name).
//
// Example Astra source:
//   add(a, b)   — 'add', 'a', 'b' are identifiers
//   let x = 5   — 'x' is an identifier
type Identifier struct {
    Name string
}

func (id *Identifier) expressionNode()      {}
func (id *Identifier) TokenLiteral() string { return id.Name }
func (id *Identifier) String() string       { return id.Name }
```

**Complete example: building the AST for `fn add(a: int, b: int) -> int { return a + b }`:**

```go
// ast/declarations_test.go
package ast

import (
    "testing"
)

func TestFunctionDeclarationString(t *testing.T) {
    // Represents:  fn add(a: int, b: int) -> int { return a + b }
    intType := &NamedType{Name: "int"}

    fn := &FunctionDeclaration{
        Name: "add",
        Parameters: []Parameter{
            {Name: "a", Type: intType},
            {Name: "b", Type: intType},
        },
        ReturnType: intType,
        Body: &BlockStatement{
            Statements: []Statement{
                &ReturnStatement{
                    Value: &Identifier{Name: "a"}, // simplified: ignores the + b part
                },
            },
        },
    }

    if fn.Name != "add" {
        t.Errorf("Expected name 'add', got '%s'", fn.Name)
    }
    if len(fn.Parameters) != 2 {
        t.Errorf("Expected 2 parameters, got %d", len(fn.Parameters))
    }
    if fn.IsVoid() {
        t.Error("Expected non-void function")
    }

    s := fn.String()
    t.Logf("FunctionDeclaration: %s", s)
}

func TestFunctionLiteralIsExpression(t *testing.T) {
    // Anonymous functions are expressions, not statements
    fl := &FunctionLiteral{
        Parameters: []Parameter{{Name: "x", Type: &NamedType{Name: "int"}}},
        ReturnType: &NamedType{Name: "int"},
        Body:       &BlockStatement{},
    }

    // FunctionLiteral must implement Expression interface
    var _ Expression = fl
    t.Logf("FunctionLiteral: %s", fl.String())
}

func TestCallExpressionString(t *testing.T) {
    // Represents: add(3, 4)
    call := &CallExpression{
        Function: &Identifier{Name: "add"},
        Arguments: []Expression{
            &Identifier{Name: "three"}, // placeholder
            &Identifier{Name: "four"},  // placeholder
        },
    }

    s := call.String()
    expected := "add(three, four)"
    if s != expected {
        t.Errorf("Expected '%s', got '%s'", expected, s)
    }
}
```

---

## 17. Exercises

1. **Trace the Call Stack** — Manually trace the call stack (draw the frames) for this program:
   ```astra
   fn multiply(a: int, b: int) -> int {
       return a * b
   }
   fn square(n: int) -> int {
       return multiply(n, n)
   }
   fn main() {
       let result = square(5)
       print(result)
   }
   ```
   What is the maximum depth of the call stack?

2. **Pure or Impure?** — Classify each function as pure or impure, and explain why:
   ```astra
   fn a(x: int) -> int { return x * 2 }
   fn b(x: int) -> int { print(x); return x * 2 }
   fn c() -> int { return random_int() }
   fn d(x: int) -> int { let y = x + 1; return y * 2 }
   ```

3. **Write a Closure** — Write an Astra closure that creates a "multiplier" function. `make_multiplier(3)` should return a function that multiplies its argument by 3.
   *Hint: the returned function captures the multiplier value from the outer scope*

4. **Variadic Sum** — Write a variadic function `sum(nums: ...int) -> int` in both Go and Astra. Then write a call that sums the numbers 1 through 5.

5. **Refactor for SRP** — The following function violates the Single Responsibility Principle. Refactor it into 2-3 smaller functions:
   ```astra
   fn process(data: string) {
       // Trim whitespace
       let trimmed = data.trim()
       // Check if empty
       if trimmed.length() == 0 {
           print("Error: empty input")
           return
       }
       // Uppercase everything
       let upper = trimmed.to_upper()
       // Print with formatting
       print(">>> " + upper + " <<<")
   }
   ```

6. **Build the AST** — Using the AST nodes defined in the milestone, manually construct the Go structs representing this Astra code:
   ```astra
   fn greet(name: string) {
       print("Hello, " + name)
   }
   ```

7. **Recursion Base Case** — The following recursive function has a bug — it will infinite loop for some inputs. Find the bug and fix it:
   ```astra
   fn power(base: int, exp: int) -> int {
       return base * power(base, exp - 1)
   }
   ```
   What is the correct base case?

8. **First-Class Functions** — Write a function `apply_twice` in Astra that takes a function `f: fn(int) -> int` and a value `x: int`, and returns `f(f(x))` (applies f twice). Test it with `fn(n: int) -> int { return n + 1 }` starting from 5. What should the result be?

---

## 18. Summary

| Concept | Go | Astra | Notes |
|---|---|---|---|
| Function keyword | `func` | `fn` | Astra is shorter |
| Parameter type | `name type` | `name: type` | Astra uses colon notation |
| Return type | `func f() int` | `fn f() -> int` | Astra uses arrow |
| Multiple returns | `func f() (int, error)` | `fn f() -> Result<int,E>` | Astra uses Result type |
| Anonymous function | `func(x int) int { }` | `fn(x: int) -> int { }` | Same concept |
| Variadic | `f(nums ...int)` | `f(nums: ...int)` | Same concept |
| Function as value | `var f func(int) int` | `let f: fn(int) -> int` | First-class |
| Overloading | Not available | Not available | Use generics instead |
| Closures | Supported (captures by reference) | Supported | Automatic capture |
| Pass by value | Default for basic types | Default for basic types | Copies are made |
| Pass by reference | Use `*T` (pointer) | Use `&mut T` | Explicit mutation |

**Key takeaways:**
- A function has a name, parameters, body, and optionally a return type
- The call stack grows when functions are called and shrinks when they return
- Parameters are passed by value (copies) by default; use pointers/references for mutation
- Closures capture variables from their surrounding scope
- First-class functions allow passing functions as arguments and returning them
- Go and Astra deliberately omit function overloading for clarity; use generics instead
- Pure functions (no side effects) enable powerful compiler optimizations
- The Single Responsibility Principle makes code easier to read, test, and maintain

---

*Next chapter: Chapter 10 — Pointers and Memory: How Data Lives in RAM*
