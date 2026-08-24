# Chapter 11: Methods

A **method** is a function that belongs to a type. You define behavior on your own types by attaching methods to them. This is how Go achieves object-oriented-like design without classes. Methods appear everywhere in Go — from the standard library (`time.Duration.String()`, `net/http.Request.FormValue()`) to every package you'll write. Understanding methods deeply — especially the difference between value and pointer receivers — is essential.

## Table of Contents

1. [Defining Methods](#1-defining-methods)
2. [Value Receivers vs Pointer Receivers](#2-value-receivers-vs-pointer-receivers)
3. [Methods on Non-Struct Types](#3-methods-on-non-struct-types)
4. [Method Sets and Addressability](#4-method-sets-and-addressability)
5. [Chaining Methods](#5-chaining-methods)
6. [Promoted Methods via Embedding](#6-promoted-methods-via-embedding)
7. [Methods vs Functions](#7-methods-vs-functions)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Defining Methods

A method has a **receiver** between `func` and the method name:

```go
// func (receiverName ReceiverType) MethodName(params) ReturnType
type Rectangle struct {
    Width  float64
    Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Method with pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width  *= factor
    r.Height *= factor
}

func (r Rectangle) String() string {
    return fmt.Sprintf("Rectangle(%.1f × %.1f)", r.Width, r.Height)
}
```

**Calling methods:**
```go
rect := Rectangle{Width: 10, Height: 5}

fmt.Println(rect.Area())       // 50
fmt.Println(rect.Perimeter())  // 30
fmt.Println(rect)              // Rectangle(10.0 × 5.0)  — uses String()

rect.Scale(2)
fmt.Println(rect.Width)        // 20

// On pointer:
rp := &Rectangle{Width: 3, Height: 4}
fmt.Println(rp.Area())    // 12 — Go auto-dereferences: (*rp).Area()
rp.Scale(3)
fmt.Println(rp.Width)     // 9
```

**Receiver name conventions:**
```go
// Use a short abbreviation of the type name (1-2 letters):
func (r Rectangle) Area() float64 { ... }   // r for Rectangle
func (u User) Validate() error { ... }      // u for User
func (s *Server) Start() error { ... }      // s for Server
func (c *Client) Do(req *Request) { ... }   // c for Client

// NOT self or this (those are not Go conventions):
// func (self Rectangle) Area() { ... }  — avoid
// func (this *User) Save() { ... }       — avoid
```

### Quick Check
> 1. What is a receiver in a method definition?
> 2. How do you call a method on a pointer to a struct?
> 3. What is the Go convention for naming receivers?

---

## 2. Value Receivers vs Pointer Receivers

This is the most important concept in this chapter. The choice affects behavior and performance:

```go
type Counter struct {
    count int
}

// VALUE RECEIVER — works on a COPY of the Counter
func (c Counter) ValueGet() int {
    return c.count
}

func (c Counter) ValueIncrement() {
    c.count++  // Modifies the COPY — original unchanged!
}

// POINTER RECEIVER — works on the ORIGINAL Counter
func (c *Counter) PointerIncrement() {
    c.count++  // Modifies the ORIGINAL
}
```

**The difference in action:**
```go
c := Counter{count: 0}

c.ValueIncrement()          // Does nothing — modifies a copy
fmt.Println(c.count)        // 0

c.PointerIncrement()        // Modifies the original
fmt.Println(c.count)        // 1

// Go automatically takes address when calling pointer methods on addressable values:
c.PointerIncrement()        // Go rewrites this as (&c).PointerIncrement()
fmt.Println(c.count)        // 2
```

**The three rules for choosing receivers:**

**Rule 1: Use pointer receiver if the method modifies the receiver:**
```go
func (u *User) SetEmail(email string) {
    u.Email = email  // Must be pointer receiver to modify
}
```

**Rule 2: Use pointer receiver if the struct is large (avoid copying):**
```go
type LargeStruct struct {
    data [1024]byte  // 1KB — copying this is expensive
}

func (l *LargeStruct) Process() {  // Pointer avoids 1KB copy on every call
    // ...
}
```

**Rule 3: Be consistent — if any method needs a pointer receiver, use pointer for all:**
```go
// INCONSISTENT (avoid):
func (u User) String() string { ... }  // value
func (u *User) Save() error { ... }    // pointer

// CONSISTENT (good):
func (u *User) String() string { ... }  // all pointer
func (u *User) Save() error { ... }     // all pointer
```

**Why consistency matters for interfaces:**
```go
type Stringer interface {
    String() string
}

type User struct{ Name string }

func (u User) String() string { return u.Name }   // value receiver

var s Stringer = User{"Alice"}   // OK: User has value receiver String
var s2 Stringer = &User{"Bob"}   // OK: *User also satisfies (pointer to value type)

// BUT if String had pointer receiver:
func (u *User) String() string { return u.Name }  // pointer receiver

var s3 Stringer = &User{"Alice"}  // OK
// var s4 Stringer = User{"Alice"} // ERROR: User doesn't have String, only *User does
```

### Quick Check
> 1. If you want a method to modify the struct, must you use a pointer receiver?
> 2. Does Go automatically take the address when calling a pointer method on a value?
> 3. Why should you use consistent receiver types across all methods of a type?

---

## 3. Methods on Non-Struct Types

You can define methods on **any type you define**, not just structs:

```go
// Method on a named type based on a built-in:
type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func (f Fahrenheit) ToCelsius() Celsius {
    return Celsius((f - 32) * 5 / 9)
}

func (c Celsius) String() string {
    return fmt.Sprintf("%.1f°C", float64(c))
}

boiling := Celsius(100)
fmt.Println(boiling.ToFahrenheit())  // 212.0°F
fmt.Println(boiling)                 // 100.0°C
```

**Method on a slice type:**
```go
type IntSlice []int

func (s IntSlice) Sum() int {
    total := 0
    for _, v := range s {
        total += v
    }
    return total
}

func (s IntSlice) Filter(pred func(int) bool) IntSlice {
    result := make(IntSlice, 0, len(s))
    for _, v := range s {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

func (s IntSlice) Map(fn func(int) int) IntSlice {
    result := make(IntSlice, len(s))
    for i, v := range s {
        result[i] = fn(v)
    }
    return result
}

nums := IntSlice{1, 2, 3, 4, 5, 6}
evens := nums.Filter(func(n int) bool { return n%2 == 0 })
doubled := evens.Map(func(n int) int { return n * 2 })
fmt.Println(doubled.Sum())  // 24 (2+4+6 doubled = 4+8+12 = 24)
```

**Method on a function type:**
```go
// http.HandlerFunc in the standard library does exactly this:
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)  // Call the function itself
}

// This lets plain functions satisfy the http.Handler interface
```

**Restriction: you can only define methods on types in the SAME package:**
```go
// CANNOT add methods to built-in types (they live in another "package"):
// func (s string) Reverse() string { ... }  // ERROR: cannot define methods on string

// Must define a new type first:
type MyString string
func (s MyString) Reverse() MyString { ... }  // OK
```

### Quick Check
> 1. Can you define methods on a named type based on `int` or `string`?
> 2. Can you add a method to `string` directly in your code?
> 3. Why does the standard library define `http.HandlerFunc` as a function type with a method?

---

## 4. Method Sets and Addressability

**Method set** — the set of methods callable on a given type:

```go
type T struct{}
func (t T)  ValueMethod() {}   // Value receiver
func (t *T) PtrMethod()   {}   // Pointer receiver

// Method sets:
// Value T:    {ValueMethod}           — can only call value receiver methods
// Pointer *T: {ValueMethod, PtrMethod} — can call BOTH
```

This is why pointer types can satisfy more interfaces than value types.

**Addressability:**
```go
type Counter struct{ count int }
func (c *Counter) Inc() { c.count++ }

// Addressable value — Go can auto-take address:
c := Counter{}
c.Inc()   // Go rewrites as (&c).Inc() — works because c is addressable

// Non-addressable — cannot auto-take address:
// Counter{}.Inc()  // ERROR: cannot take address of temporary Counter{}
// map[string]Counter{}["key"].Inc()  // ERROR: map values are not addressable
```

**Fixing the map value problem:**
```go
// Map values are not addressable, so you can't call pointer methods on them
type Counter struct{ count int }
func (c *Counter) Inc() { c.count++ }

m := map[string]Counter{}
m["a"] = Counter{}
// m["a"].Inc()  // ERROR

// Fix 1: Use pointer in map
m2 := map[string]*Counter{}
m2["a"] = &Counter{}
m2["a"].Inc()  // OK

// Fix 2: Copy out, modify, put back
c := m["a"]
c.Inc()        // OK (c is addressable local variable)
m["a"] = c    // Put back
```

### Quick Check
> 1. What is the method set of a pointer type `*T` vs a value type `T`?
> 2. What does "addressable" mean and why does it matter for pointer methods?
> 3. Why can't you call a pointer receiver method on a map value?

---

## 5. Chaining Methods

Methods can return the receiver to enable fluent, chainable APIs:

```go
type QueryBuilder struct {
    table  string
    wheres []string
    limit  int
}

func (q *QueryBuilder) Table(t string) *QueryBuilder {
    q.table = t
    return q
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
    q.wheres = append(q.wheres, condition)
    return q
}

func (q *QueryBuilder) Limit(n int) *QueryBuilder {
    q.limit = n
    return q
}

func (q *QueryBuilder) Build() string {
    sql := "SELECT * FROM " + q.table
    if len(q.wheres) > 0 {
        sql += " WHERE " + strings.Join(q.wheres, " AND ")
    }
    if q.limit > 0 {
        sql += fmt.Sprintf(" LIMIT %d", q.limit)
    }
    return sql
}

// Fluent API — each method returns *QueryBuilder
query := new(QueryBuilder).
    Table("users").
    Where("age >= 18").
    Where("is_active = TRUE").
    Limit(50).
    Build()
// SELECT * FROM users WHERE age >= 18 AND is_active = TRUE LIMIT 50
```

**Standard library example — `strings.Builder`:**
```go
var b strings.Builder
b.WriteString("Hello")
b.WriteString(", ")
b.WriteString("World")
b.WriteByte('!')
result := b.String()  // "Hello, World!"
```

### Quick Check
> 1. How do you make a method chainable?
> 2. Why do chainable methods typically use pointer receivers?

---

## 6. Promoted Methods via Embedding

When you embed a type, its methods are "promoted" to the outer type:

```go
type Logger struct{ prefix string }

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.prefix, msg)
}

type UserService struct {
    *Logger  // Embed pointer to Logger
    db *sql.DB
}

// UserService now has Log() via promotion
svc := &UserService{
    Logger: &Logger{prefix: "UserService"},
}
svc.Log("Started")  // [UserService] Started — no need to svc.Logger.Log()
```

**The outer type can shadow promoted methods:**
```go
type Base struct{}
func (b Base) Hello() string { return "Base hello" }

type Child struct {
    Base
}
// Child inherits Base.Hello

type GrandChild struct {
    Child
}
// GrandChild inherits Child (and Base's) Hello

// Override in GrandChild:
func (g GrandChild) Hello() string { return "GrandChild hello" }

gc := GrandChild{}
fmt.Println(gc.Hello())        // "GrandChild hello" — own method wins
fmt.Println(gc.Child.Hello())  // "Base hello" — access via explicit path
```

### Quick Check
> 1. What does it mean for a method to be "promoted"?
> 2. What happens when both the outer type and embedded type have the same method name?

---

## 7. Methods vs Functions

When should you use a method vs a standalone function?

```go
// As a METHOD:
func (u *User) Validate() error { ... }

// As a FUNCTION:
func ValidateUser(u *User) error { ... }
```

**Use a method when:**
- The operation is intrinsically tied to the type
- You need to modify the receiver
- You want to implement an interface
- It's the primary action of the type

**Use a function when:**
- The operation uses multiple types equally
- You want to be testable without creating an instance
- It's a utility that doesn't "belong" to one type

```go
// Good method — belongs to User:
func (u *User) ChangePassword(newPassword string) error { ... }

// Better as function — needs both User and Database:
func TransferFunds(from, to *Account, amount float64) error { ... }
// Not better as: from.TransferTo(to, amount) or to.ReceiveFrom(from, amount)
```

**Method expressions** — call a method without an instance:
```go
type Dog struct{ Name string }
func (d Dog) Speak() string { return d.Name + " says woof" }

// Normal method call:
d := Dog{"Rex"}
fmt.Println(d.Speak())

// Method expression — get the method as a function:
speak := Dog.Speak  // speak has type: func(Dog) string
fmt.Println(speak(Dog{"Buddy"}))  // Buddy says woof
```

### Quick Check
> 1. When should you use a method instead of a standalone function?
> 2. What is a method expression?
> 3. Can you pass a method as a function argument?

---

## Summary

- **Method syntax**: `func (receiver ReceiverType) MethodName() ReturnType`
- **Value receiver**: gets a copy; reads are fine; modifications don't affect original
- **Pointer receiver**: gets the original; can modify; avoid copying large structs
- **Consistency rule**: use the same receiver type (pointer or value) across all methods of a type
- **Non-struct types**: methods work on any named type in the same package
- **Method sets**: `*T` has both pointer and value receiver methods; `T` has only value receivers
- **Addressability**: map values, temporaries, and unaddressable values can't auto-take address
- **Chaining**: return the receiver pointer from methods for fluent APIs
- **Promotion**: embedded type's methods are accessible on the outer type
- **Methods vs functions**: methods for type-specific behavior; functions for cross-type operations

---

## Exercises

### Easy
1. Create a `Stack[T]` type using a slice. Add methods: `Push(v T)`, `Pop() (T, bool)`, `Peek() (T, bool)`, `IsEmpty() bool`, `Size() int`. Use generics so it works for any type.
2. Create a `Duration` type wrapping `time.Duration`. Add methods: `Hours() float64`, `Minutes() float64`, `Seconds() float64`, `String() string` (e.g., "2h30m15s"), `Add(other Duration) Duration`.
3. Demonstrate the difference between value and pointer receivers by creating a `Balance` type with `Deposit(amount float64)` (pointer receiver) and `Display() string` (value receiver). Show that calling `Deposit` on a value copy doesn't affect the original.

### Medium
4. Linked list with methods: Implement a singly linked list with a `Node[T]` struct and a `LinkedList[T]` struct. Add methods: `Append(v T)`, `Prepend(v T)`, `Delete(v T) bool`, `Contains(v T) bool`, `ToSlice() []T`, `Reverse()`, `Len() int`. Implement `String() string` for pretty printing.
5. Method chaining for validation: Build a `Validator` struct that validates a value through a chain of rules. `NewValidator(value interface{})` creates one. Methods: `Required()`, `MinLength(n int)`, `MaxLength(n int)`, `IsEmail()`, `IsURL()`, `Matches(pattern string)` — each adds a rule. `Validate() []ValidationError` runs all rules and returns a list of errors. Chain: `NewValidator(email).Required().IsEmail().Validate()`.
6. Event emitter: Implement an `EventEmitter` struct with methods: `On(event string, handler func(data interface{}))` (register handler), `Off(event string, handler func(data interface{}))` (remove handler), `Emit(event string, data interface{})` (call all handlers), `Once(event string, handler func(data interface{}))` (fires handler only once). Support multiple handlers per event.

### Hard
7. Observable value: Create a generic `Observable[T]` type that wraps a value and notifies subscribers when it changes. Methods: `Get() T`, `Set(v T)` (updates value, notifies all subscribers), `Subscribe(fn func(oldVal, newVal T)) func()` (returns an unsubscribe function), `Map(fn func(T) U) *Observable[U]` (creates a derived observable that updates when source updates). Thread-safe using RWMutex. Write tests verifying: notifications fire in order, unsubscribe works, derived observables update automatically.
8. Type-safe HTTP client: Build an HTTP client wrapper where API endpoints are defined as method-bearing types. Define `type API struct{ baseURL, token string }`. Methods like `GetUser(id int) (*User, error)`, `CreateUser(u CreateUserRequest) (*User, error)`, `ListUsers(filter UserFilter) ([]User, error)`. Each method: constructs the URL, sets auth header, makes the request, decodes JSON response into typed struct, handles HTTP errors into typed `APIError`. Write integration tests against a mock HTTP server (`httptest.NewServer`).
