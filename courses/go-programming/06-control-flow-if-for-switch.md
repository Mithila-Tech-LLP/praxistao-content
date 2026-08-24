# Chapter 06: Control Flow — if, for, switch

Control flow determines the order in which code executes. Without it, every program would just run top to bottom in a straight line. Go has three core control flow statements: `if` for branching, `for` for looping, and `switch` for multi-way branching. Go is deliberately minimal here — there is no `while` loop, no `do-while`, no `foreach`. Everything is `for`. This simplicity makes Go code consistent and readable.

## Table of Contents

1. [if and else](#1-if-and-else)
2. [for — Go's Only Loop](#2-for--gos-only-loop)
3. [range — Iterating Over Collections](#3-range--iterating-over-collections)
4. [switch — Multi-Way Branching](#4-switch--multi-way-branching)
5. [goto, break, continue, labels](#5-goto-break-continue-labels)
6. [Common Patterns](#6-common-patterns)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. if and else

**Basic if:**
```go
temperature := 35

if temperature > 30 {
    fmt.Println("It's hot!")
}
```

**if-else:**
```go
score := 75

if score >= 90 {
    fmt.Println("Grade: A")
} else if score >= 80 {
    fmt.Println("Grade: B")
} else if score >= 70 {
    fmt.Println("Grade: C")
} else {
    fmt.Println("Grade: F")
}
```

**Go requires braces {} — no single-line if:**
```go
// This does NOT compile in Go:
// if x > 0 fmt.Println("positive")

// Must always have braces:
if x > 0 {
    fmt.Println("positive")
}
```

**if with an initialization statement** — one of Go's most useful patterns:
```go
// Without initialization statement:
result, err := divide(10, 2)
if err != nil {
    return err
}
fmt.Println(result)

// WITH initialization statement (cleaner):
if result, err := divide(10, 2); err != nil {
    return err
} else {
    fmt.Println(result)
}
// Note: result is only in scope INSIDE the if/else block
```

This pattern is very common with error handling and type assertions:
```go
// Type assertion with init statement:
if s, ok := value.(string); ok {
    fmt.Println("String value:", s)
}

// Database query with init statement:
if user, err := db.GetUser(id); err != nil {
    log.Error("user not found", "id", id, "error", err)
    return ErrUserNotFound
} else {
    return processUser(user)
}
```

**No truthiness in Go** — conditions must be exactly `bool`:
```go
n := 0

// This does NOT work in Go (unlike C, Python, JavaScript):
// if n { ... }  // ERROR: cannot use n (type int) as bool

// Must be explicit:
if n != 0 {
    fmt.Println("n is not zero")
}

var s string = ""
// if s { ... }  // ERROR
if s != "" {
    fmt.Println("s is not empty")
}
```

### Quick Check
> 1. Does Go require braces `{}` for if statements?
> 2. What is the "initialization statement" pattern in Go's if?
> 3. Can you write `if n { ... }` where n is an integer in Go?

---

## 2. for — Go's Only Loop

Go has exactly ONE looping construct: `for`. But it serves all purposes.

**Form 1: C-style for loop:**
```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
// Prints: 0 1 2 3 4
```

The three parts: `init statement`, `condition`, `post statement`.

**Form 2: While-style (condition only):**
```go
n := 1
for n < 100 {
    n *= 2
}
fmt.Println(n)  // 128
```

**Form 3: Infinite loop:**
```go
for {
    // Runs forever until break or return
    input := readInput()
    if input == "quit" {
        break
    }
    processInput(input)
}
```

**Modifying the loop variable:**
```go
// Skip even numbers:
for i := 0; i < 10; i++ {
    if i % 2 == 0 {
        continue  // Skip to next iteration
    }
    fmt.Print(i, " ")  // 1 3 5 7 9
}
```

**Nested loops:**
```go
for i := 1; i <= 3; i++ {
    for j := 1; j <= 3; j++ {
        fmt.Printf("%d×%d=%d\t", i, j, i*j)
    }
    fmt.Println()
}
```

**Loop variable gotcha with goroutines:**
```go
// Before Go 1.22 this was a classic BUG: all goroutines shared the
// same variable i, so this often printed 3,3,3 instead of 0,1,2.
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)
    }()
}

// Since Go 1.22, EVERY iteration gets a fresh i (in all for-loop forms,
// not just range loops), so the code above now prints 0, 1, 2 (in some order).

// On older Go versions, the fix was to pass i as an argument:
for i := 0; i < 3; i++ {
    go func(i int) {
        fmt.Println(i)  // Each goroutine gets its own copy of i
    }(i)
}
```

### Quick Check
> 1. How do you write a while loop in Go?
> 2. How do you write an infinite loop in Go?
> 3. What is the classic goroutine loop variable bug and how do you fix it?

---

## 3. range — Iterating Over Collections

`range` is used with `for` to iterate over slices, arrays, maps, strings, and channels:

**Range over a slice:**
```go
fruits := []string{"apple", "banana", "cherry"}

// Index and value:
for i, fruit := range fruits {
    fmt.Printf("%d: %s\n", i, fruit)
}
// 0: apple
// 1: banana
// 2: cherry

// Index only:
for i := range fruits {
    fmt.Println(i)
}

// Value only (discard index):
for _, fruit := range fruits {
    fmt.Println(fruit)
}
```

**Range over a map:**
```go
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
    "Carol": 35,
}

for name, age := range ages {
    fmt.Printf("%s is %d years old\n", name, age)
}
// Note: Map iteration order is RANDOM — don't rely on it
```

**Range over a string:**
```go
// range over a string iterates RUNES (Unicode code points), not bytes
for i, r := range "Hello, 世界" {
    fmt.Printf("index=%d rune=%c code=%d\n", i, r, r)
}
// index=0 rune=H code=72
// index=1 rune=e code=101
// ...
// index=7 rune=世 code=19990
// index=10 rune=界 code=30028
// Note: index jumps by 3 for multi-byte runes
```

**Range over a channel:**
```go
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // Must close to stop range
}()

for v := range ch {
    fmt.Println(v)  // 0 1 2 3 4
}
```

**Range over integers (Go 1.22+):**
```go
// New in Go 1.22: range over an integer
for i := range 5 {
    fmt.Println(i)  // 0 1 2 3 4
}
```

**Creating a copy with range:**
```go
numbers := []int{1, 2, 3, 4, 5}

// Modifying v does NOT modify the slice (v is a copy)
for _, v := range numbers {
    v *= 2  // This does nothing to numbers
}

// To modify the slice, use the index:
for i := range numbers {
    numbers[i] *= 2  // This DOES modify numbers
}
```

### Quick Check
> 1. What is the iteration order when ranging over a map?
> 2. What does `range` give you when iterating over a string?
> 3. If you modify `v` in `for _, v := range slice`, does it modify the original slice?

---

## 4. switch — Multi-Way Branching

Go's `switch` is more powerful than in C/Java. Cases don't fall through by default.

**Basic switch:**
```go
day := "Monday"

switch day {
case "Saturday", "Sunday":
    fmt.Println("Weekend!")
case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
    fmt.Println("Weekday")
default:
    fmt.Println("Unknown day")
}
```

**Switch with no expression (like if-else chain):**
```go
score := 85

switch {
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
case score >= 70:
    fmt.Println("C")
default:
    fmt.Println("F")
}
```

**Switch with initialization:**
```go
switch os := runtime.GOOS; os {
case "darwin":
    fmt.Println("macOS")
case "linux":
    fmt.Println("Linux")
default:
    fmt.Printf("Other: %s\n", os)
}
```

**No automatic fallthrough** (unlike C):
```go
n := 2

switch n {
case 1:
    fmt.Println("one")
    // Does NOT fall through to case 2
case 2:
    fmt.Println("two")  // Only this prints
case 3:
    fmt.Println("three")
}
```

**Explicit fallthrough** (rarely needed):
```go
switch n {
case 1:
    fmt.Println("one")
    fallthrough  // Explicitly fall through to next case
case 2:
    fmt.Println("two")  // Prints if n==1 OR n==2
}
```

**Type switch** — check what type is in an interface:
```go
func printType(v interface{}) {
    switch t := v.(type) {
    case int:
        fmt.Printf("int: %d\n", t)
    case string:
        fmt.Printf("string: %q\n", t)
    case bool:
        fmt.Printf("bool: %v\n", t)
    case []int:
        fmt.Printf("[]int with %d elements\n", len(t))
    default:
        fmt.Printf("unknown type: %T\n", t)
    }
}

printType(42)           // int: 42
printType("hello")      // string: "hello"
printType([]int{1,2,3}) // []int with 3 elements
```

### Quick Check
> 1. Does Go's switch fall through automatically like C?
> 2. How do you use switch without an expression?
> 3. What is a type switch used for?

---

## 5. goto, break, continue, labels

**`break` and `continue`** (you've seen these):
```go
// break exits the current loop
for i := 0; i < 10; i++ {
    if i == 5 {
        break  // Exit immediately when i == 5
    }
    fmt.Print(i, " ")  // 0 1 2 3 4
}

// continue skips to next iteration
for i := 0; i < 10; i++ {
    if i % 2 == 0 {
        continue  // Skip even numbers
    }
    fmt.Print(i, " ")  // 1 3 5 7 9
}
```

**Labels with break** — break out of outer loop:
```go
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            break outer  // Break out of the OUTER loop
        }
        fmt.Printf("(%d,%d) ", i, j)
    }
}
// (0,0) (0,1) (0,2) (1,0)
```

**Labels with continue:**
```go
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 {
            continue outer  // Skip to next iteration of OUTER loop
        }
        fmt.Printf("(%d,%d) ", i, j)
    }
}
// (0,0) (1,0) (2,0)
```

**`goto`** — jump to a label (rarely used, usually a code smell):
```go
func main() {
    i := 0
loop:
    if i < 5 {
        fmt.Println(i)
        i++
        goto loop  // Jump back to label
    }
}
// Use a for loop instead! goto is rarely the right answer.
```

### Quick Check
> 1. What is the purpose of labeled `break`?
> 2. When would you use labeled `continue`?
> 3. Why is `goto` rarely the right choice in Go?

---

## 6. Common Patterns

**Early return pattern** (avoid deep nesting):
```go
// BAD: Pyramid of doom (deep nesting)
func processOrder(order *Order) error {
    if order != nil {
        if order.UserID != 0 {
            user, err := db.GetUser(order.UserID)
            if err == nil {
                if user.IsActive() {
                    // actual logic here
                    return saveOrder(order)
                } else {
                    return ErrInactiveUser
                }
            } else {
                return err
            }
        } else {
            return ErrMissingUserID
        }
    } else {
        return ErrNilOrder
    }
}

// GOOD: Guard clauses + early return (flat structure)
func processOrder(order *Order) error {
    if order == nil {
        return ErrNilOrder
    }
    if order.UserID == 0 {
        return ErrMissingUserID
    }
    
    user, err := db.GetUser(order.UserID)
    if err != nil {
        return err
    }
    if !user.IsActive() {
        return ErrInactiveUser
    }
    
    return saveOrder(order)
}
```

**Loop with index tracking:**
```go
// Find first element matching a condition:
func findFirst(items []int, target int) (int, bool) {
    for i, item := range items {
        if item == target {
            return i, true
        }
    }
    return -1, false
}
```

**Collecting results:**
```go
// Filter slice:
func filter(items []int, pred func(int) bool) []int {
    result := make([]int, 0, len(items))
    for _, item := range items {
        if pred(item) {
            result = append(result, item)
        }
    }
    return result
}

evens := filter([]int{1,2,3,4,5,6}, func(n int) bool { return n%2 == 0 })
// [2, 4, 6]
```

**Defer with loop:**
```go
// WRONG: defer in a loop doesn't run until function returns
for _, filename := range files {
    f, err := os.Open(filename)
    if err != nil { continue }
    defer f.Close()  // All defers run at END of function, not end of loop iteration!
    processFile(f)
}

// CORRECT: wrap in a function to limit defer scope
for _, filename := range files {
    func() {
        f, err := os.Open(filename)
        if err != nil { return }
        defer f.Close()  // Now runs when this anonymous function returns
        processFile(f)
    }()
}
```

### Quick Check
> 1. What is the "early return" pattern and why is it preferred in Go?
> 2. What is the bug when using `defer` inside a for loop?
> 3. How do you fix `defer` inside a loop?

---

## Summary

- **if**: requires braces; condition must be bool (no implicit truthiness); initialization statement pattern is idiomatic
- **for**: Go's only loop; three forms: C-style / while-style / infinite loop; `break`/`continue` for control
- **range**: iterate over slices (index+value), maps (key+value), strings (index+rune), channels; map order is random
- **switch**: no automatic fallthrough; can match multiple values per case; expression-less switch = if-else chain; type switch for interface dispatch
- **labels**: `break label` and `continue label` for nested loop control
- **patterns**: early return (guard clauses), avoid deep nesting, defer-in-loop gotcha

---

## Exercises

### Easy
1. Write a `FizzBuzz` function: for numbers 1–100, print "Fizz" if divisible by 3, "Buzz" if divisible by 5, "FizzBuzz" if both, otherwise the number.
2. Write a function `reverseSlice(s []int) []int` using a for loop.
3. Use a switch statement with a string variable to identify if an HTTP status code string ("200", "404", "500") is success, not found, or server error.

### Medium
4. Nested loop patterns: Write a function `printDiamond(n int)` that prints a diamond pattern of asterisks. For n=3:
   ```
     *
    ***
   *****
    ***
     *
   ```
   Then write `printSpiral(n int)` that fills an n×n matrix in spiral order (1,2,3,... going right, down, left, up) and prints it.
5. Number processing: Given a slice of integers: (a) Find the second largest without sorting. (b) Find all pairs that sum to a target value (no duplicates). (c) Find the longest consecutive sequence. All must work in O(n) time or O(n log n) with explanation.
6. State machine: Implement a simple vending machine using a switch on state: States: `Idle`, `HasMoney`, `Dispensing`. Transitions: `insertCoin()`, `selectItem()`, `dispense()`, `returnCoin()`. Use a struct to track state and balance. Write tests for each valid and invalid transition.

### Hard
7. Text statistics: Write a complete program that reads a text file (use `os.ReadFile`) and outputs: total characters (bytes vs runes), word count, line count, sentence count (ends with . ! ?), most frequent word, average word length, histogram of letter frequencies. Handle: Unicode text, multiple spaces, mixed case (treat "The" and "the" as same word), punctuation attached to words.
8. Maze solver: Given a 2D grid (slice of strings) where `#` is a wall, `.` is open, `S` is start, `E` is end: implement BFS to find the shortest path. Print the original grid, then the grid with the path marked as `*`. Return the path length or -1 if no path exists. Bonus: implement DFS and compare path lengths. Handle edge cases: no start/end, multiple starts/ends, disconnected maze.
