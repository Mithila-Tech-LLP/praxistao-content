# Chapter 10: Structs

A **struct** is a composite type that groups together fields of different types under a single name. If you're coming from object-oriented languages, a struct is like a class — but without inheritance. In Go, structs are the primary way to define custom types that represent real-world entities: a User, an Order, a Product, a Request. Understanding structs deeply is fundamental because they appear in every Go program.

## Table of Contents

1. [Defining and Creating Structs](#1-defining-and-creating-structs)
2. [Accessing and Modifying Fields](#2-accessing-and-modifying-fields)
3. [Struct Methods (Preview)](#3-struct-methods-preview)
4. [Embedding — Composition over Inheritance](#4-embedding--composition-over-inheritance)
5. [Anonymous Structs](#5-anonymous-structs)
6. [Struct Tags](#6-struct-tags)
7. [Comparing Structs](#7-comparing-structs)
8. [Struct Patterns](#8-struct-patterns)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Defining and Creating Structs

```go
// Define a struct type
type User struct {
    ID        int
    Name      string
    Email     string
    Age       int
    IsActive  bool
    CreatedAt time.Time
}
```

**Creating struct instances:**
```go
// Method 1: Field names (recommended — order-independent, self-documenting)
u1 := User{
    ID:       1,
    Name:     "Alice",
    Email:    "alice@example.com",
    Age:      30,
    IsActive: true,
    CreatedAt: time.Now(),
}

// Method 2: Positional (fragile — breaks if fields are reordered)
// u2 := User{1, "Alice", "alice@example.com", 30, true, time.Now()}

// Method 3: Zero value then assign
var u3 User
u3.ID = 2
u3.Name = "Bob"
// Other fields get zero values: Email="", Age=0, IsActive=false, CreatedAt=time.Time{}

// Partial initialization (unspecified fields get zero values):
u4 := User{
    Name: "Carol",
    Email: "carol@example.com",
    // ID=0, Age=0, IsActive=false, CreatedAt=zero
}
```

**Struct pointer:**
```go
// Pointer to struct (common in Go):
up := &User{
    ID:   3,
    Name: "Dave",
}

// Access fields through pointer (auto-dereferenced):
fmt.Println(up.Name)   // "Dave" — same as (*up).Name, Go does this automatically
up.Age = 25            // Modify through pointer
```

**Constructor functions (Go convention):**
```go
// Go doesn't have constructors — use a function by convention
func NewUser(name, email string) (*User, error) {
    if name == "" {
        return nil, errors.New("name cannot be empty")
    }
    if !strings.Contains(email, "@") {
        return nil, errors.New("invalid email")
    }
    return &User{
        ID:        generateID(),
        Name:      name,
        Email:     strings.ToLower(email),
        IsActive:  true,
        CreatedAt: time.Now(),
    }, nil
}

user, err := NewUser("Alice", "alice@example.com")
if err != nil {
    log.Fatal(err)
}
```

### Quick Check
> 1. What do unspecified fields in a struct literal get?
> 2. Why is field-name initialization preferred over positional?
> 3. Why do Go programs often return `*User` instead of `User` from constructor functions?

---

## 2. Accessing and Modifying Fields

```go
type Rectangle struct {
    Width  float64
    Height float64
}

r := Rectangle{Width: 10, Height: 5}

// Access fields:
fmt.Println(r.Width)   // 10
fmt.Println(r.Height)  // 5

// Modify fields:
r.Width = 20

// Computed value:
area := r.Width * r.Height
fmt.Println("Area:", area)  // 100
```

**Structs are value types** — assignment copies the entire struct:
```go
a := Rectangle{Width: 10, Height: 5}
b := a         // b is a FULL COPY of a
b.Width = 99

fmt.Println(a.Width)  // 10 — a unchanged
fmt.Println(b.Width)  // 99
```

**When to use pointer vs value:**
```go
// Value receiver (small structs, when you don't need to modify):
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (modify the struct, OR for large structs to avoid copy):
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor    // Modifies the original
    r.Height *= factor
}

rect := Rectangle{10, 5}
rect.Scale(2)
fmt.Println(rect.Width)  // 20 (modified via pointer)
```

**Fields in functions:**
```go
// Passing by value — function gets a copy
func printUser(u User) {
    u.Name = "Modified"  // Only modifies the copy
    fmt.Println(u)
}

// Passing by pointer — function can modify original
func activateUser(u *User) {
    u.IsActive = true  // Modifies the original
}

user := User{Name: "Alice", IsActive: false}
activateUser(&user)
fmt.Println(user.IsActive)  // true
```

### Quick Check
> 1. Are structs value types or reference types in Go?
> 2. What is the difference between passing a struct by value vs by pointer to a function?
> 3. When should you use a pointer receiver vs a value receiver?

---

## 3. Struct Methods (Preview)

Methods are functions attached to a type. We cover them in depth in Chapter 11, but here's a preview:

```go
type Circle struct {
    Radius float64
}

// Method with value receiver
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Method with pointer receiver (can modify c)
func (c *Circle) Scale(factor float64) {
    c.Radius *= factor
}

// String() method — implements fmt.Stringer interface
func (c Circle) String() string {
    return fmt.Sprintf("Circle(r=%.2f)", c.Radius)
}

c := Circle{Radius: 5}
fmt.Println(c.Area())    // 78.53...
c.Scale(2)
fmt.Println(c)           // Circle(r=10.00) — uses String() automatically
```

### Quick Check
> 1. What is the difference between a function and a method in Go?
> 2. What is the `String()` method and why is it useful?

---

## 4. Embedding — Composition over Inheritance

Go doesn't have inheritance. Instead, it has **embedding** — you can embed one struct inside another to "inherit" its fields and methods:

```go
type Animal struct {
    Name   string
    Weight float64
}

func (a Animal) Describe() string {
    return fmt.Sprintf("%s weighs %.1f kg", a.Name, a.Weight)
}

type Dog struct {
    Animal        // Embedded (no field name) — Dog "inherits" Animal's fields and methods
    Breed  string
}

func (d Dog) Bark() string {
    return d.Name + " says: Woof!"  // Can access embedded fields directly
}

dog := Dog{
    Animal: Animal{Name: "Rex", Weight: 30},
    Breed:  "German Shepherd",
}

fmt.Println(dog.Name)        // Rex (promoted field from Animal)
fmt.Println(dog.Weight)      // 30 (promoted field)
fmt.Println(dog.Describe())  // Rex weighs 30.0 kg (promoted method)
fmt.Println(dog.Bark())      // Rex says: Woof!
```

**Embedding is NOT inheritance** — it's composition. Key differences:
```go
type Cat struct {
    Animal
    Indoor bool
}

cat := Cat{Animal: Animal{Name: "Whiskers"}}

// Promoted fields are accessible directly:
fmt.Println(cat.Name)  // "Whiskers"

// But the outer type doesn't satisfy Animal's interface automatically
// A *Cat is not a *Animal in type terms (no polymorphism like OOP)
// Instead, cat.Animal is an embedded field of type Animal

// Explicit access:
fmt.Println(cat.Animal.Name)  // Also works: "Whiskers"
```

**Overriding promoted methods:**
```go
type GuideDog struct {
    Dog
    Handler string
}

// GuideDog can "override" Describe by defining its own:
func (g GuideDog) Describe() string {
    return fmt.Sprintf("%s (guide dog, handler: %s)", g.Name, g.Handler)
}

gd := GuideDog{Dog: Dog{Animal: Animal{Name: "Buddy"}}, Handler: "John"}
fmt.Println(gd.Describe())      // uses GuideDog's Describe
fmt.Println(gd.Dog.Describe())  // explicitly calls Dog's Describe (Animal's)
```

**Multiple embedding:**
```go
type Logger struct{}
func (l Logger) Log(msg string) { fmt.Println("[LOG]", msg) }

type Validator struct{}
func (v Validator) Validate(s string) bool { return s != "" }

type Service struct {
    Logger            // Embed Logger
    Validator         // Embed Validator
    name string
}

s := Service{name: "UserService"}
s.Log("Service started")       // From Logger
ok := s.Validate(s.name)       // From Validator
```

### Quick Check
> 1. What is embedding and how is it different from inheritance?
> 2. How do you access an embedded struct's fields from the outer struct?
> 3. Can you "override" an embedded struct's method?

---

## 5. Anonymous Structs

Anonymous structs have no named type — they're useful for one-off data groupings:

```go
// One-time struct literal (no type name needed):
person := struct {
    Name string
    Age  int
}{
    Name: "Alice",
    Age:  30,
}
fmt.Println(person.Name)  // Alice

// Anonymous struct in a slice (common for test cases):
tests := []struct {
    input    string
    expected int
}{
    {"hello", 5},
    {"world!", 6},
    {"", 0},
}

for _, tt := range tests {
    got := len(tt.input)
    if got != tt.expected {
        fmt.Printf("len(%q) = %d, want %d\n", tt.input, got, tt.expected)
    }
}
```

**Anonymous structs for JSON decoding when you don't want a named type:**
```go
var response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        UserID int    `json:"user_id"`
        Token  string `json:"token"`
    } `json:"data"`
}

json.Unmarshal(jsonBytes, &response)
fmt.Println(response.Data.UserID)
```

### Quick Check
> 1. When are anonymous structs useful?
> 2. How are anonymous structs commonly used in tests?

---

## 6. Struct Tags

Struct tags are string metadata attached to fields. They are used by reflection-based libraries (encoding/json, database ORM, validation, etc.):

```go
type User struct {
    ID        int       `json:"id" db:"user_id"`
    Name      string    `json:"name" db:"name" validate:"required,min=2,max=50"`
    Email     string    `json:"email" db:"email" validate:"required,email"`
    Password  string    `json:"-"`  // "-" means omit from JSON
    Age       int       `json:"age,omitempty"`  // omitempty: skip if zero value
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

**JSON tags in action:**
```go
u := User{
    ID:    1,
    Name:  "Alice",
    Email: "alice@example.com",
    Age:   0,  // omitempty — will be omitted
}

data, _ := json.Marshal(u)
fmt.Println(string(data))
// {"id":1,"name":"Alice","email":"alice@example.com","created_at":"0001-01-01T00:00:00Z"}
// Age omitted (zero + omitempty), Password omitted (-)

// Decoding:
jsonStr := `{"id":2,"name":"Bob","email":"bob@example.com"}`
var u2 User
json.Unmarshal([]byte(jsonStr), &u2)
fmt.Println(u2.Name)  // Bob
```

**Reading struct tags with reflection:**
```go
import "reflect"

t := reflect.TypeOf(User{})
field, _ := t.FieldByName("Email")
fmt.Println(field.Tag.Get("json"))      // "email"
fmt.Println(field.Tag.Get("validate"))  // "required,email"
```

**Common tag keys:**
| Tag key | Used by |
|---------|---------|
| `json` | `encoding/json` |
| `db` | `sqlx`, `sqlc` |
| `yaml` | `gopkg.in/yaml.v3` |
| `validate` | `github.com/go-playground/validator` |
| `binding` | Gin web framework |
| `form` | Form data parsing |
| `bson` | MongoDB driver |

### Quick Check
> 1. What does `json:"-"` mean in a struct tag?
> 2. What does `json:"age,omitempty"` do?
> 3. How do you read a struct tag programmatically?

---

## 7. Comparing Structs

**Structs are comparable if all fields are comparable:**
```go
type Point struct {
    X, Y int
}

p1 := Point{1, 2}
p2 := Point{1, 2}
p3 := Point{3, 4}

fmt.Println(p1 == p2)  // true
fmt.Println(p1 == p3)  // false
fmt.Println(p1 != p3)  // true
```

**Structs with non-comparable fields can't use ==:**
```go
type Container struct {
    Name  string
    Items []int  // Slice is not comparable!
}

c1 := Container{"box", []int{1, 2, 3}}
c2 := Container{"box", []int{1, 2, 3}}
// fmt.Println(c1 == c2)  // COMPILE ERROR: struct contains []int which is not comparable

// Use reflect.DeepEqual for deep comparison:
import "reflect"
fmt.Println(reflect.DeepEqual(c1, c2))  // true
```

**Struct as map key (must be comparable):**
```go
type Point struct{ X, Y int }

// Point is comparable, so can be a map key:
grid := map[Point]string{
    {0, 0}: "origin",
    {1, 0}: "right",
    {0, 1}: "up",
}
fmt.Println(grid[Point{1, 0}])  // "right"
```

### Quick Check
> 1. When can you use `==` to compare two structs?
> 2. What function do you use for deep comparison when `==` isn't available?
> 3. Can you use a struct with a slice field as a map key?

---

## 8. Struct Patterns

**Option pattern (functional options):**
```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
    maxConn int
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithMaxConnections(n int) Option {
    return func(s *Server) {
        s.maxConn = n
    }
}

func NewServer(host string, opts ...Option) *Server {
    s := &Server{
        host:    host,
        port:    8080,           // default
        timeout: 30 * time.Second, // default
        maxConn: 100,            // default
    }
    for _, opt := range opts {
        opt(s)  // Apply each option
    }
    return s
}

// Usage — clean API, extensible, defaults:
s := NewServer("localhost",
    WithPort(9090),
    WithTimeout(60*time.Second),
)
```

**Builder pattern:**
```go
type QueryBuilder struct {
    table      string
    conditions []string
    limit      int
    orderBy    string
}

func (b *QueryBuilder) Where(condition string) *QueryBuilder {
    b.conditions = append(b.conditions, condition)
    return b  // Return self for chaining
}

func (b *QueryBuilder) Limit(n int) *QueryBuilder {
    b.limit = n
    return b
}

func (b *QueryBuilder) OrderBy(field string) *QueryBuilder {
    b.orderBy = field
    return b
}

func (b *QueryBuilder) Build() string {
    query := "SELECT * FROM " + b.table
    if len(b.conditions) > 0 {
        query += " WHERE " + strings.Join(b.conditions, " AND ")
    }
    if b.orderBy != "" {
        query += " ORDER BY " + b.orderBy
    }
    if b.limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", b.limit)
    }
    return query
}

// Method chaining:
query := (&QueryBuilder{table: "users"}).
    Where("age > 18").
    Where("is_active = true").
    OrderBy("name").
    Limit(10).
    Build()
// SELECT * FROM users WHERE age > 18 AND is_active = true ORDER BY name LIMIT 10
```

### Quick Check
> 1. What is the "functional options" pattern and what problem does it solve?
> 2. What makes the builder pattern useful for constructing complex objects?
> 3. How does method chaining work in Go?

---

## Summary

- **Define**: `type Name struct { Field1 Type1; ... }`
- **Create**: field-name literals preferred; `&User{...}` for pointer; constructor functions with `New` prefix
- **Value type**: assignment copies the whole struct; pass `*Struct` to functions that need to modify
- **Embedding**: embed one struct inside another to promote fields and methods (composition, not inheritance)
- **Anonymous structs**: no type name; great for test cases and one-off JSON decoding
- **Struct tags**: string metadata used by reflection-based libraries (`json:"name,omitempty"`)
- **Comparable**: `==` works if all fields are comparable; use `reflect.DeepEqual` otherwise
- **Patterns**: functional options (`...Option`), builder (method chaining)

---

## Exercises

### Easy
1. Define a `BankAccount` struct with `Owner string`, `Balance float64`, `Currency string`. Write `Deposit(amount float64)` and `Withdraw(amount float64) error` methods. `Withdraw` returns an error if insufficient funds.
2. Define an `Address` struct. Embed it inside a `Person` struct. Demonstrate that `person.City` works (promoted field from Address).
3. Write a struct `Config` that reads from JSON: `{"host": "localhost", "port": 8080, "debug": true, "max_connections": 100}`. Use proper JSON tags with `omitempty` for optional fields.

### Medium
4. Shape hierarchy with embedding: Define `Shape` (embedded in all shapes) with `Color string` and a `Describe() string` method. Define `Circle`, `Rectangle`, `Triangle` — each embedding `Shape` and having their own fields. Write `Area() float64` and `Perimeter() float64` for each. Collect all shapes in a `[]interface{}` slice and print their areas and colors.
5. Functional options for HTTP client: Implement a custom `HTTPClient` struct with options: `BaseURL`, `Timeout`, `MaxRetries`, `UserAgent`, `Headers map[string]string`. Use the functional options pattern. Provide sensible defaults. The client should have a `Get(path string) (*http.Response, error)` method that applies all configured options.
6. Graph node with embedding: Define `BaseNode` with `ID string` and `Metadata map[string]interface{}`. Define `UserNode` (embeds BaseNode, adds `Name`, `Email` fields) and `ProductNode` (embeds BaseNode, adds `Title`, `Price` fields). Write a function that takes a `[]BaseNode` and builds an adjacency list graph.

### Hard
7. Struct serializer: Write a function `Serialize(v interface{}) map[string]interface{}` using reflection that converts any struct to a map using its field names (or `json` tag names if present). Handle: nested structs (recursively), pointer fields (dereference), slices, maps, skipping unexported fields, skipping fields with `json:"-"`. Compare its output to `json.Marshal` for test cases.
8. Config with validation: Build a config loading system: (a) Define a `Config` struct with 15+ fields (database, server, cache, feature flags). (b) Load from environment variables using `env:` struct tags via reflection (`DB_HOST` → `DBHost`). (c) Load from YAML file using `yaml:` tags. (d) Merge (env overrides file). (e) Validate using `validate:` tags (required, min, max, email, url, oneof=dev prod staging). (f) Return structured validation errors listing every failed field. (g) Write unit tests for each step.
