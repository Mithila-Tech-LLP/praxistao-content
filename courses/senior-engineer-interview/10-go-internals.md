# Chapter 10: Go Internals — Memory Model, Structs, Interfaces & Embedding

This chapter covers the Go internals that senior interviewers probe deeply. Understanding these makes you a fundamentally better Go engineer — and signals to interviewers that you think about what the language is actually doing, not just what it appears to do.

## Table of Contents

1. [Go's Memory Model](#1-gos-memory-model)
2. [Struct Layout & Padding](#2-struct-layout--padding)
3. [Interface Internals](#3-interface-internals)
4. [Implicit Interface Satisfaction](#4-implicit-interface-satisfaction)
5. [Embedding vs Inheritance](#5-embedding-vs-inheritance)
6. [Composition Patterns](#6-composition-patterns)
7. [Value vs Pointer Semantics](#7-value-vs-pointer-semantics)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Go's Memory Model

Go's memory model defines when a goroutine reading a variable is guaranteed to observe values written by a different goroutine.

**The fundamental rule:** If a goroutine reads a variable that another goroutine writes, access must be synchronized. Without synchronization, the behavior is undefined — the compiler and CPU can reorder operations freely.

### Happens-Before Relationship

Go's memory model uses the concept of "happens-before" ordering. A write W to variable v is **observed** by a read R of v if:
- W happens-before R
- No other write to v happens between W and R

```go
// RACE CONDITION — no happens-before between write and read
var x int

func main() {
    go func() { x = 1 }() // write in goroutine
    fmt.Println(x)          // read in main — may see 0 or 1
}

// SAFE — channel send establishes happens-before with channel receive
ch := make(chan struct{})
go func() {
    x = 1
    ch <- struct{}{} // send HAPPENS-BEFORE receive
}()
<-ch
fmt.Println(x) // guaranteed to see x = 1
```

### Synchronization Operations That Create Happens-Before

| Operation | Establishes happens-before |
|---|---|
| Channel send | Happens before the receive of that value |
| Channel close | Happens before receive returning zero value |
| sync.Mutex Unlock | Happens before the next Lock |
| sync.WaitGroup Done | Happens before Wait returns |
| sync/atomic operations | Happens before subsequent operations on same variable |

---

## 2. Struct Layout & Padding

Go aligns struct fields to their natural alignment. This creates "padding" bytes between fields. The order of fields affects memory usage.

```go
// BAD: field ordering wastes memory due to padding
type BadStruct struct {
    a bool    // 1 byte
    // 7 bytes padding (to align b to 8-byte boundary)
    b int64   // 8 bytes
    c bool    // 1 byte
    // 7 bytes padding (to make struct size a multiple of 8)
}
// Total: 24 bytes — 14 bytes are padding!

// GOOD: order fields largest-to-smallest
type GoodStruct struct {
    b int64   // 8 bytes
    a bool    // 1 byte
    c bool    // 1 byte
    // 6 bytes padding (to make size a multiple of 8)
}
// Total: 16 bytes — 6 bytes padding (instead of 14)
```

```go
// Check actual sizes with unsafe.Sizeof
import "unsafe"
fmt.Println(unsafe.Sizeof(BadStruct{}))  // 24
fmt.Println(unsafe.Sizeof(GoodStruct{})) // 16
```

**When does this matter in interviews?** When discussing memory efficiency, cache locality, or designing data structures that will be allocated millions of times. For most business logic structs, it does not matter.

---

## 3. Interface Internals

An interface value in Go is represented as a pair of two pointers:
1. A pointer to the **type descriptor** (runtime type information — methods, size, etc.)
2. A pointer to the **data** (the actual value)

```go
type Stringer interface {
    String() string
}

// An interface value has this internal structure:
// [  type pointer  |  data pointer  ]
//  8 bytes           8 bytes
// Total: 16 bytes on a 64-bit system

// When you assign a concrete type to an interface:
var s Stringer = MyType{name: "test"}
// Go stores: type = *MyType descriptor, data = pointer to MyType value
```

### The Nil Interface Trap

This is one of the most common Go bugs that interviewers test.

```go
// This looks like it returns nil, but it does not!
func getError() error {
    var err *MyError = nil // typed nil pointer
    // ... some logic that doesn't set err ...
    return err // WRONG: returns a non-nil interface!
}

// The returned error interface has:
// type = *MyError (non-nil!)
// data = nil
// So error != nil, even though the underlying value is nil

// WHY: An interface is nil only if BOTH type and data are nil.

// CORRECT:
func getErrorCorrect() error {
    var err *MyError = nil
    if err == nil {
        return nil // return untyped nil — both type and data are nil
    }
    return err
}

// The trap in practice:
func process() error {
    var myErr *ValidationError
    if someCondition {
        myErr = &ValidationError{msg: "invalid"}
    }
    return myErr // BUG: always returns non-nil error!
    // return nil // CORRECT: explicitly return nil when no error
}
```

### Interface Comparison

```go
// Two interface values are equal if they have the same dynamic type AND equal dynamic value
var a interface{} = 42
var b interface{} = 42
fmt.Println(a == b) // true

var c interface{} = []int{1, 2, 3}
var d interface{} = []int{1, 2, 3}
// fmt.Println(c == d) // PANIC: slices are not comparable!
```

---

## 4. Implicit Interface Satisfaction

Go's interfaces are satisfied implicitly — you never declare "this type implements this interface." If a type has all the required methods, it satisfies the interface automatically.

```go
// Define an interface
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Define a type
type FileLogger struct {
    file *os.File
}

// Implement the method — no "implements Writer" declaration needed
func (fl *FileLogger) Write(p []byte) (n int, err error) {
    return fl.file.Write(p)
}

// FileLogger now implicitly satisfies Writer
var w Writer = &FileLogger{} // compiles! FileLogger satisfies Writer
```

### Compile-Time Interface Check

Use this pattern to verify interface satisfaction at compile time (not at runtime):

```go
// This will cause a compile error if *FileLogger does not implement Writer
var _ Writer = (*FileLogger)(nil)

// Or:
var _ Writer = &FileLogger{}
```

This is idiomatic Go for packages that export types meant to implement specific interfaces.

### Why Implicit Interfaces Are Powerful

They allow you to define abstractions at the point of use, not at the point of definition.

```go
// stdlib doesn't know about your DatabaseLogger when it defined io.Writer.
// But your DatabaseLogger can still be used anywhere io.Writer is expected.
type DatabaseLogger struct { db *sql.DB }
func (dl *DatabaseLogger) Write(p []byte) (int, error) {
    // write log to database
    _, err := dl.db.Exec("INSERT INTO logs (msg) VALUES (?)", string(p))
    return len(p), err
}

// Works with fmt, log, http — anything that accepts an io.Writer
log.SetOutput(&DatabaseLogger{db: db})
```

---

## 5. Embedding vs Inheritance

Go has no inheritance. Instead it has embedding — one type includes another type's methods and fields directly.

```go
// Embedding: Animal's fields and methods are promoted to Dog
type Animal struct {
    Name string
}
func (a Animal) Speak() string { return a.Name + " makes a sound" }

type Dog struct {
    Animal         // embedded — fields and methods promoted
    Breed string
}

d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
fmt.Println(d.Name)    // promoted field: d.Animal.Name
fmt.Println(d.Speak()) // promoted method: d.Animal.Speak()
```

### Embedding is NOT Inheritance

The critical distinction: embedding gives method promotion but NOT polymorphism.

```go
// With embedding: type is NOT substitutable
func makeSpeak(a Animal) string { return a.Speak() }
// makeSpeak(d) // COMPILE ERROR: Dog is not an Animal

// For polymorphism: use interfaces
type Speaker interface { Speak() string }
func makeSpeak(s Speaker) string { return s.Speak() }
makeSpeak(d) // Works: Dog satisfies Speaker because it has a Speak() method
```

### Embedding Multiple Types

```go
type ReadWriter struct {
    io.Reader
    io.Writer
}

// Now ReadWriter has both Read() and Write() methods
```

### Method Overriding via Embedding

```go
type Base struct{}
func (b Base) Method() string { return "Base.Method" }

type Child struct {
    Base
}
// Override by defining the same method on Child
func (c Child) Method() string { return "Child.Method" }

c := Child{}
fmt.Println(c.Method())     // "Child.Method" — Child's version
fmt.Println(c.Base.Method()) // "Base.Method" — can still call Base's version
```

---

## 6. Composition Patterns

### Interface-Based Composition (The Right Way)

```go
// Instead of embedding concrete types, embed interfaces
type Logger interface {
    Log(msg string)
}

type Service struct {
    Logger // embed interface — can swap implementations
    db *sql.DB
}

// Now you can inject any Logger implementation
func NewService(logger Logger, db *sql.DB) *Service {
    return &Service{Logger: logger, db: db}
}

// In tests: inject a mock logger
// In production: inject a real logger
```

### Functional Options Pattern

A common Go pattern for building structs with many optional configurations:

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) { s.host = host }
}

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080, timeout: 30 * time.Second}
    for _, opt := range opts { opt(s) }
    return s
}

// Clean call site:
server := NewServer(
    WithHost("example.com"),
    WithPort(9090),
)
```

---

## 7. Value vs Pointer Semantics

This is fundamental to writing correct Go code.

```go
// VALUE receiver: method gets a COPY of the struct
func (p Point) Scale(factor float64) Point {
    return Point{p.X * factor, p.Y * factor} // doesn't modify original
}

// POINTER receiver: method can modify the original struct
func (p *Point) ScaleInPlace(factor float64) {
    p.X *= factor // modifies the original
    p.Y *= factor
}

// Rules for choosing:
// 1. If method needs to modify the receiver: pointer receiver
// 2. If struct is large: pointer receiver (avoids copying)
// 3. If struct has a mutex: pointer receiver (must not copy)
// 4. For consistency: if any method is pointer, all should be pointer
```

### When Interfaces Force a Pointer Receiver

```go
type Mover interface { Move() }

type Car struct { pos int }
func (c *Car) Move() { c.pos++ }

var m Mover
m = &Car{}  // OK: *Car has Move()
m = Car{}   // COMPILE ERROR: Car (value) does not have Move()
            // because Move() is defined on *Car

// Rule: a value type T satisfies an interface if all methods are on T (value receiver)
// A pointer *T always satisfies the interface if any method is on *T
```

---

## 8. Interview Questions & Model Answers

**Q: What is the nil interface trap in Go?**

"A nil interface has both its type and data pointers as nil. But if you assign a typed nil pointer to an interface, the interface has a non-nil type pointer — so the interface itself is not nil, even though the underlying value is nil. This trips up people who do `if err != nil` checks — a typed nil error will pass this check when it shouldn't. The fix: always return the untyped `nil` instead of a typed nil pointer."

**Q: How does Go's interface work internally?**

"An interface value is an (iface, data) pair — two pointers of 8 bytes each on 64-bit systems. The iface pointer points to a type descriptor containing the method table (vtable) and type information. The data pointer points to the actual value or, for small values, stores the value inline. Method dispatch is an indirect function call through the vtable — slightly slower than direct calls, which is why Go prefers small interfaces."

**Q: What is the difference between embedding and inheritance?**

"Embedding is syntactic sugar for composition — you get method promotion so you can call embedded type's methods directly. But it's not inheritance: the outer type is not a subtype of the embedded type. There's no polymorphism: you can't pass a `Dog` where an `Animal` is expected. For polymorphism you use interfaces. The Go philosophy is 'accept interfaces, return structs' — define abstractions at the point of use."

**Q: When would you use a value receiver vs a pointer receiver?**

"Pointer receiver when: the method modifies the receiver, the struct is large (to avoid copy cost), or the struct has a sync.Mutex. Value receiver when: the struct is small, the method shouldn't modify it, and the type is conceptually immutable. The key rule: be consistent — if any method needs a pointer receiver, make all methods pointer receivers so the type satisfies interfaces correctly."

---

## Summary

- Go's memory model: goroutines sharing data need synchronization. Channel send/receive, mutex, and atomic operations create happens-before ordering.
- Struct padding: order fields largest to smallest to minimize padding. Use `unsafe.Sizeof` to verify.
- Interface internals: (type, data) pair. An interface is nil only if BOTH are nil — the nil interface trap.
- Implicit satisfaction: no `implements` declaration. If you have the methods, you satisfy the interface.
- Embedding is NOT inheritance: you get method promotion but not subtype polymorphism. Use interfaces for polymorphism.
- Value receiver = copy; pointer receiver = modify in place. Be consistent within a type.
