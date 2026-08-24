# Chapter 12: Interfaces

Interfaces are Go's most powerful abstraction mechanism. An **interface** defines a set of method signatures — any type that implements all those methods automatically satisfies the interface. No `implements` keyword needed. This is called **implicit interface satisfaction** and it's one of Go's most elegant features. Interfaces are how Go achieves polymorphism, decoupling, and testability.

## Table of Contents

1. [What Is an Interface](#1-what-is-an-interface)
2. [Implicit Satisfaction — The Go Way](#2-implicit-satisfaction--the-go-way)
3. [The Empty Interface](#3-the-empty-interface)
4. [Interface Values Internals](#4-interface-values-internals)
5. [Type Assertions and Type Switches](#5-type-assertions-and-type-switches)
6. [Common Standard Library Interfaces](#6-common-standard-library-interfaces)
7. [Interface Composition](#7-interface-composition)
8. [Interface Best Practices](#8-interface-best-practices)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What Is an Interface

An interface defines **what a type can do**, not what it is:

```go
// An interface is a set of method signatures
type Animal interface {
    Sound() string
    Name() string
}

// Any type that has Sound() and Name() satisfies Animal
type Dog struct{ name string }
func (d Dog) Sound() string { return "Woof" }
func (d Dog) Name() string  { return d.name }

type Cat struct{ name string }
func (c Cat) Sound() string { return "Meow" }
func (c Cat) Name() string  { return c.name }

// Both Dog and Cat satisfy Animal:
func describe(a Animal) {
    fmt.Printf("%s says %s\n", a.Name(), a.Sound())
}

describe(Dog{"Rex"})    // Rex says Woof
describe(Cat{"Kitty"})  // Kitty says Meow
```

**Interfaces define contracts.** A function that accepts an interface says: "I don't care what type you give me, as long as it can do these things."

### Quick Check
> 1. What does it mean for a type to "satisfy" an interface?
> 2. Does Go require an explicit `implements` keyword?
> 3. What does an interface define — data or behavior?

---

## 2. Implicit Satisfaction — The Go Way

The Go approach is fundamentally different from Java/C#:

```go
// Java/C# style (explicit):
// class Dog implements Animal { ... }

// Go style (implicit — no declaration needed):
type Writer interface {
    Write(p []byte) (n int, err error)
}

// os.File has a Write method with this exact signature
// So os.File automatically satisfies io.Writer — without ever mentioning it

// bytes.Buffer has Write with the same signature
// So bytes.Buffer automatically satisfies io.Writer too

// YOUR custom type:
type MyBuffer struct{ data []byte }
func (b *MyBuffer) Write(p []byte) (int, error) {
    b.data = append(b.data, p...)
    return len(p), nil
}
// MyBuffer now satisfies io.Writer automatically!
```

**Why implicit satisfaction is powerful:**
```go
// This function from the STANDARD LIBRARY accepts io.Writer:
func Fprintf(w io.Writer, format string, a ...interface{}) (int, error)

// Because of implicit satisfaction, Fprintf works with:
// - os.Stdout (a *os.File)
// - os.Stderr
// - bytes.Buffer
// - strings.Builder  (write to string)
// - net.Conn         (write to network)
// - YOUR custom writer
// — without any of these types ever explicitly declaring "I implement io.Writer"
```

**Interface satisfaction check at compile time:**
```go
// Force a compile-time check that MyBuffer satisfies io.Writer:
var _ io.Writer = (*MyBuffer)(nil)  // Idiom: assign nil to interface variable

// If MyBuffer doesn't have Write(), this line fails to compile
// with a helpful error message. Use this in your code to document intent.
```

### Quick Check
> 1. How does Go know a type satisfies an interface without explicit declaration?
> 2. Why is implicit satisfaction more flexible than explicit `implements`?
> 3. How do you force a compile-time interface check?

---

## 3. The Empty Interface

The empty interface `interface{}` (or `any` in Go 1.18+) has no methods — every type satisfies it:

```go
// interface{} and any are identical — any is just an alias
var anything interface{}
var anything2 any  // Same thing

anything = 42
anything = "hello"
anything = []int{1, 2, 3}
anything = struct{ X int }{42}

// Useful for storing values of unknown type:
func printAnything(v any) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}

printAnything(42)           // Type: int, Value: 42
printAnything("hello")      // Type: string, Value: hello
printAnything([]int{1,2,3}) // Type: []int, Value: [1 2 3]
```

**Common uses:**
```go
// JSON unmarshaling into unknown structure:
var result map[string]any
json.Unmarshal(jsonData, &result)

// Function that can accept anything:
func Log(msg string, fields map[string]any) {
    fmt.Printf("%s: %+v\n", msg, fields)
}

Log("user created", map[string]any{
    "user_id": 123,
    "email":   "alice@example.com",
    "age":     30,
})
```

**The trade-off**: `any` loses type safety. Prefer concrete types or specific interfaces when possible:
```go
// Avoid this when you know the types:
func add(a, b any) any {
    return a.(int) + b.(int)  // Runtime panic if not int
}

// Better: use generics (Go 1.18+):
func add[T int | float64](a, b T) T {
    return a + b
}
```

### Quick Check
> 1. What is `any` and what types satisfy it?
> 2. Name two legitimate use cases for the empty interface.
> 3. Why should you avoid `any` when specific types are known?

---

## 4. Interface Values Internals

An interface value has two components: **(type, value)**:

```go
// Conceptually, an interface value is:
type iface struct {
    typeInfo *typeDescriptor  // What type is stored
    dataPtr  unsafe.Pointer   // Pointer to the actual data
}
```

```
var w io.Writer = os.Stdout

Interface value w:
  ┌──────────────────────────────────┐
  │ type:  *os.File                  │
  │ value: ptr → (file descriptor 1) │
  └──────────────────────────────────┘

w = &bytes.Buffer{}
  ┌──────────────────────────────────┐
  │ type:  *bytes.Buffer             │
  │ value: ptr → (buffer data)       │
  └──────────────────────────────────┘
```

**Nil interface vs interface containing nil:**
```go
// Nil interface: both type and value are nil
var w io.Writer        // type=nil, value=nil
fmt.Println(w == nil)  // true

// Interface containing nil POINTER — NOT a nil interface!
var f *os.File         // f is nil pointer to os.File
w = f                  // type=*os.File, value=nil
fmt.Println(w == nil)  // FALSE! Interface has a type, so it's not nil
```

**This is a classic Go gotcha:**
```go
func errFunc() error {
    var err *MyError = nil  // nil pointer to MyError
    // ...some condition...
    return err  // Returns an interface{type: *MyError, value: nil}
}

err := errFunc()
if err != nil {  // TRUE — even though the value is nil!
    fmt.Println("Got an error")  // This prints!
}

// FIX: return nil directly, not a typed nil pointer
func errFuncFixed() error {
    var err *MyError = nil
    if err != nil {
        return err
    }
    return nil  // Returns {type: nil, value: nil} — a true nil interface
}
```

### Quick Check
> 1. What are the two components of an interface value?
> 2. Why is `var f *os.File = nil; var w io.Writer = f; w == nil` false?
> 3. How do you return a "true nil" error from a function?

---

## 5. Type Assertions and Type Switches

**Type assertion** — extract the concrete type from an interface:

```go
var i interface{} = "hello"

// Form 1: panic if wrong type
s := i.(string)         // "hello"
n := i.(int)            // PANIC: interface conversion: string is not int

// Form 2: safe (comma-ok pattern)
s, ok := i.(string)     // "hello", true
n, ok := i.(int)        // 0, false  — no panic
```

**Type switch** — check multiple types efficiently:
```go
func describe(i interface{}) string {
    switch v := i.(type) {
    case nil:
        return "nil"
    case int:
        return fmt.Sprintf("int: %d", v)
    case float64:
        return fmt.Sprintf("float64: %.2f", v)
    case string:
        return fmt.Sprintf("string: %q", v)
    case bool:
        return fmt.Sprintf("bool: %t", v)
    case []int:
        return fmt.Sprintf("[]int with %d elements", len(v))
    case error:
        return fmt.Sprintf("error: %s", v.Error())
    default:
        return fmt.Sprintf("unknown type %T", v)
    }
}
```

**Type assertion for interface checking:**
```go
type Stringer interface {
    String() string
}

// Check if a value implements an interface at runtime:
if s, ok := someValue.(Stringer); ok {
    fmt.Println(s.String())
}

// Common use: check if an error is a specific type
var err error = &NetworkError{Code: 503, Message: "Service Unavailable"}

if netErr, ok := err.(*NetworkError); ok {
    fmt.Printf("Network error %d: %s\n", netErr.Code, netErr.Message)
    if netErr.Code >= 500 {
        // Retry logic
    }
}

// Go 1.13+ errors.As (preferred over type assertion for errors):
var netErr *NetworkError
if errors.As(err, &netErr) {
    fmt.Println(netErr.Code)
}
```

### Quick Check
> 1. What is the difference between `s := i.(string)` and `s, ok := i.(string)`?
> 2. When would you use a type switch instead of multiple type assertions?
> 3. What is `errors.As` and why is it preferred over direct type assertion for errors?

---

## 6. Common Standard Library Interfaces

The standard library defines small, focused interfaces. Learn these — you'll use them constantly:

**`io.Reader` and `io.Writer` — the most important pair:**
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// These are implemented by: os.File, bytes.Buffer, strings.Reader,
// net.Conn, gzip.Writer, crypto/cipher streams, and hundreds more.

// Functions that accept io.Reader work with ALL of them:
func processInput(r io.Reader) {
    data, err := io.ReadAll(r)  // works for file, network, string, buffer...
}
```

**Key standard interfaces:**
```go
// fmt.Stringer — custom string representation
type Stringer interface {
    String() string
}

// fmt.GoStringer — Go syntax representation  
type GoStringer interface {
    GoString() string
}

// error — the fundamental error interface
type error interface {
    Error() string
}

// io.Closer — anything that can be closed
type Closer interface {
    Close() error
}

// io.ReadWriter, io.ReadCloser, io.WriteCloser — combinations
type ReadCloser interface {
    Reader
    Closer
}

// sort.Interface — for custom sorting
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

// encoding.TextMarshaler / TextUnmarshaler — for custom text encoding
type TextMarshaler interface {
    MarshalText() (text []byte, err error)
}

// json.Marshaler / Unmarshaler — for custom JSON encoding
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
```

**Implementing `sort.Interface`:**
```go
type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

people := []Person{{"Alice", 30}, {"Bob", 25}, {"Carol", 35}}
sort.Sort(ByAge(people))
// [{Bob 25} {Alice 30} {Carol 35}]

// Simpler with sort.Slice (Go 1.8+):
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

### Quick Check
> 1. What is `io.Reader` and name three types that implement it.
> 2. What does implementing `fmt.Stringer` give you?
> 3. What three methods does `sort.Interface` require?

---

## 7. Interface Composition

You can compose interfaces from smaller interfaces — a key Go design pattern:

```go
// Compose from smaller interfaces:
type ReadWriter interface {
    io.Reader  // Embed io.Reader
    io.Writer  // Embed io.Writer
}

// A type satisfies ReadWriter if it has BOTH Read and Write methods:
// bytes.Buffer satisfies ReadWriter

// Your own compositions:
type Repository interface {
    Find(id string) (*Entity, error)
    Save(e *Entity) error
    Delete(id string) error
    FindAll() ([]*Entity, error)
}

// Break it down:
type Finder interface {
    Find(id string) (*Entity, error)
    FindAll() ([]*Entity, error)
}

type Writer interface {
    Save(e *Entity) error
    Delete(id string) error
}

type Repository interface {
    Finder
    Writer
}

// Now you can accept just Finder in read-only functions:
func listAll(r Finder) ([]*Entity, error) {
    return r.FindAll()
}
```

**The Interface Segregation Principle in Go:**
> "Accept the smallest interface that satisfies your needs."

```go
// BAD: forces caller to implement huge interface
func logSize(s Storage) {
    fmt.Println(s.Size())  // Only needs Size — why require full Storage?
}

// GOOD: accept only what you need
type Sizer interface {
    Size() int64
}

func logSize(s Sizer) {
    fmt.Println(s.Size())  // Any type with Size() works
}
```

### Quick Check
> 1. How do you compose one interface from multiple smaller interfaces?
> 2. What principle says "accept the smallest interface that satisfies your needs"?
> 3. Why is `func foo(r io.Reader)` better than `func foo(f *os.File)` in most cases?

---

## 8. Interface Best Practices

**1. Define interfaces at the point of use (in the consumer package), not the provider:**
```go
// BAD: service package defines interface, handlers must import service
// service/user.go:
type UserServiceInterface interface {
    GetUser(id int) (*User, error)
    CreateUser(req CreateUserRequest) (*User, error)
}

// GOOD: handler package defines what it needs:
// handler/user.go:
type userGetter interface {
    GetUser(id int) (*User, error)
}

// handler only depends on a narrow interface it owns
// service.UserService satisfies it without knowing about handler
```

**2. Keep interfaces small:**
```go
// BAD: fat interface
type UserStore interface {
    GetByID(id int) (*User, error)
    GetByEmail(email string) (*User, error)
    Create(u *User) error
    Update(u *User) error
    Delete(id int) error
    Count() int
    Search(query string) ([]*User, error)
    ListByRole(role string) ([]*User, error)
}

// GOOD: small, focused interfaces
type UserReader interface {
    GetByID(id int) (*User, error)
    GetByEmail(email string) (*User, error)
}

type UserWriter interface {
    Create(u *User) error
    Update(u *User) error
    Delete(id int) error
}
```

**3. Use interfaces for testability:**
```go
// Without interface — hard to test:
type UserService struct {
    db *sql.DB  // Can't swap out in tests without real DB
}

// With interface — easy to test:
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(u *User) error
}

type UserService struct {
    repo UserRepository  // Can inject a mock in tests
}

// In tests:
type mockRepo struct{ users map[int]*User }
func (m *mockRepo) FindByID(id int) (*User, error) {
    if u, ok := m.users[id]; ok { return u, nil }
    return nil, ErrNotFound
}
func (m *mockRepo) Save(u *User) error {
    m.users[u.ID] = u
    return nil
}

func TestUserService(t *testing.T) {
    svc := &UserService{repo: &mockRepo{users: map[int]*User{}}}
    // Test without a real database!
}
```

**4. Don't use interfaces for everything:**
```go
// OVER-ENGINEERED: interface for a type used in one place
type NameGetterInterface interface {
    GetName() string
}

// Just use the concrete type directly if there's only one implementation
type User struct{ Name string }
func (u *User) GetName() string { return u.Name }
```

### Quick Check
> 1. Where should interfaces be defined — in the provider or consumer package?
> 2. Why are small interfaces better than large ones?
> 3. How do interfaces enable testing without a real database?

---

## Summary

- **Interface**: set of method signatures; any type implementing them satisfies it
- **Implicit satisfaction**: no `implements` keyword; Go checks methods at compile time
- **Empty interface `any`**: every type satisfies it; use when type is truly unknown
- **Interface values**: (type, value) pair internally; nil interface ≠ interface holding nil pointer
- **Type assertion**: `v.(T)` for concrete type; `v, ok := v.(T)` for safe version
- **Type switch**: `switch v := i.(type)` for multi-type dispatch
- **Standard interfaces**: `io.Reader`, `io.Writer`, `error`, `fmt.Stringer`, `sort.Interface`
- **Composition**: compose interfaces from smaller ones with embedding
- **Best practices**: define at use site, keep small, use for testability

---

## Exercises

### Easy
1. Define a `Shape` interface with `Area() float64` and `Perimeter() float64`. Implement it for `Circle`, `Rectangle`, `Triangle`. Write `TotalArea(shapes []Shape) float64` and `LargestShape(shapes []Shape) Shape`.
2. Implement `fmt.Stringer` for three of your own types. Verify they print nicely with `fmt.Println`.
3. Write a function `processFile(r io.Reader) (wordCount int, lineCount int, err error)` that counts words and lines. Test it with `strings.NewReader`, `bytes.NewReader`, and an actual file.

### Medium
4. Middleware chain with interfaces: Define `type Handler interface { Handle(ctx context.Context, req Request) (Response, error) }`. Write three middlewares wrapping a Handler: `LoggingMiddleware` (logs request/response), `AuthMiddleware` (checks auth token), `RetryMiddleware` (retries on error up to N times). Chain them: `Chain(base Handler, middlewares ...Middleware) Handler`. Write tests with a mock Handler.
5. Plugin system with interfaces: Define a `Plugin interface` with `Name() string`, `Init(config map[string]string) error`, `Execute(data []byte) ([]byte, error)`. Create three plugins: `Base64Plugin` (encodes/decodes), `CompressionPlugin` (gzip), `EncryptionPlugin` (XOR cipher). Write a `Pipeline` that runs data through a sequence of plugins.
6. Custom JSON marshaling: Define a `Money` type with `Amount int64` (cents) and `Currency string`. Implement `json.Marshaler` (marshal as `"19.99 USD"`) and `json.Unmarshaler` (parse `"19.99 USD"`). Also implement `fmt.Stringer`. Write tests verifying roundtrip JSON encoding/decoding.

### Hard
7. Dependency injection framework: Build a minimal DI container. It should: `Register(iface interface{}, factory func() interface{})` — register a factory for an interface, `Resolve(iface interface{}) interface{}` — get or create an instance (singleton by default), `RegisterScoped(iface interface{}, factory func() interface{})` — create new instance per request context. Use reflection to match interface types. Write a demo: `Logger`, `Database`, and `UserService` interfaces where UserService depends on the others.
8. Event system with interfaces: Build a type-safe event system. Define `type Event[T any] struct{ payload T }`. `type Handler[T any] interface { Handle(ctx context.Context, event Event[T]) error }`. `type Bus struct` with methods: `Subscribe[T any](handler Handler[T])`, `Publish[T any](ctx context.Context, event Event[T]) error` (calls all handlers, collects errors), `PublishAsync[T any](ctx context.Context, event Event[T])` (non-blocking). Write tests for: event delivery, multiple subscribers, error collection, async delivery.
