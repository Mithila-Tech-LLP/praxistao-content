# Chapter 25: Reflection

Reflection lets Go programs inspect and manipulate their own types at runtime — reading struct field names, calling methods dynamically, and creating values of types known only at runtime. It powers `encoding/json`, ORMs, testing frameworks, and dependency injection containers. Learn it to understand how Go's most powerful libraries work, and to write them yourself.

## Table of Contents

1. [The reflect Package Basics](#1-the-reflect-package-basics)
2. [reflect.Type — Inspecting Types](#2-reflecttype--inspecting-types)
3. [reflect.Value — Inspecting and Modifying Values](#3-reflectvalue--inspecting-and-modifying-values)
4. [Structs — Field Inspection and Tags](#4-structs--field-inspection-and-tags)
5. [Calling Functions and Methods Dynamically](#5-calling-functions-and-methods-dynamically)
6. [Creating Values at Runtime](#6-creating-values-at-runtime)
7. [Real-World Use Cases](#7-real-world-use-cases)
8. [When to Avoid Reflection](#8-when-to-avoid-reflection)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The reflect Package Basics

Reflection is accessed through two types: `reflect.Type` (the type descriptor) and `reflect.Value` (the value container):

```go
import "reflect"

x := 42
t := reflect.TypeOf(x)   // *reflect.rtype — describes "int"
v := reflect.ValueOf(x)  // reflect.Value — wraps the value 42

fmt.Println(t)            // int
fmt.Println(t.Kind())     // int (Kind is the underlying category)
fmt.Println(v)            // 42
fmt.Println(v.Int())      // 42 (extract as int64)
```

**Kind vs Type:**
```go
type MyInt int

x := MyInt(5)
t := reflect.TypeOf(x)

fmt.Println(t)         // main.MyInt   (the named type)
fmt.Println(t.Kind())  // int          (the underlying kind)
```

**The Kind enum — every Go type maps to one:**
```go
// Kinds:
reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
reflect.Map, reflect.Pointer, reflect.Slice, reflect.String, reflect.Struct, reflect.UnsafePointer
```

### Quick Check
> 1. What is the difference between `reflect.TypeOf` and `reflect.ValueOf`?
> 2. What is the difference between Type and Kind?
> 3. What Kind does `type MyInt int` have?

---

## 2. reflect.Type — Inspecting Types

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"-"`
}

t := reflect.TypeOf(User{})

fmt.Println(t.Name())        // User
fmt.Println(t.Kind())        // struct
fmt.Println(t.NumField())    // 3
fmt.Println(t.PkgPath())     // main (package where User is defined)

// Struct fields:
for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Printf("  %s %s tag:%q\n", field.Name, field.Type, field.Tag)
}
// Output:
//   Name string tag:"json:\"name\""
//   Email string tag:"json:\"email\""
//   Age int tag:"json:\"-\""

// Get a specific field by name:
f, ok := t.FieldByName("Name")
if ok {
    fmt.Println(f.Tag.Get("json"))  // "name"
}
```

**Pointer types:**
```go
t := reflect.TypeOf(&User{})
fmt.Println(t.Kind())     // ptr
fmt.Println(t.Elem())     // main.User (the pointed-to type)
fmt.Println(t.Elem().Kind()) // struct
```

**Slice, map, channel types:**
```go
st := reflect.TypeOf([]int{})
fmt.Println(st.Kind())  // slice
fmt.Println(st.Elem())  // int

mt := reflect.TypeOf(map[string]int{})
fmt.Println(mt.Kind())  // map
fmt.Println(mt.Key())   // string
fmt.Println(mt.Elem())  // int
```

**Method inspection:**
```go
t := reflect.TypeOf(User{})
fmt.Println(t.NumMethod())  // 0 (no methods on User value)

// Methods on pointer receiver:
t = reflect.TypeOf(&User{})
for i := 0; i < t.NumMethod(); i++ {
    m := t.Method(i)
    fmt.Println(m.Name, m.Type)
}
```

### Quick Check
> 1. How do you get the type of the element in a slice using reflection?
> 2. How do you read a struct tag value?
> 3. How do you iterate over struct fields?

---

## 3. reflect.Value — Inspecting and Modifying Values

```go
x := 42
v := reflect.ValueOf(x)

// Reading values:
fmt.Println(v.Kind())   // int
fmt.Println(v.Int())    // 42 (returns int64)
fmt.Println(v.Interface()) // 42 (returns as interface{})

// Type-specific extraction:
var fval float64 = 3.14
fv := reflect.ValueOf(fval)
fmt.Println(fv.Float())  // 3.14 (returns float64)

var s string = "hello"
sv := reflect.ValueOf(s)
fmt.Println(sv.String())  // "hello"

var b bool = true
bv := reflect.ValueOf(b)
fmt.Println(bv.Bool())  // true
```

**Modifying values — requires a pointer:**
```go
x := 42
v := reflect.ValueOf(x)
// v.SetInt(100)  // PANIC: reflect.Value.SetInt using value obtained using unexported field

// To modify, pass a pointer:
v = reflect.ValueOf(&x).Elem()  // Elem() dereferences the pointer
v.CanSet()   // true
v.SetInt(100)
fmt.Println(x)  // 100 — changed!
```

**`CanSet()` — checking if a value is settable:**
```go
x := 42

v1 := reflect.ValueOf(x)
fmt.Println(v1.CanSet())  // false — passed by value, not addressable

v2 := reflect.ValueOf(&x).Elem()
fmt.Println(v2.CanSet())  // true — addressable via pointer

// Unexported struct fields are never settable:
type T struct { exported string; unexported string }
t := T{}
tv := reflect.ValueOf(&t).Elem()
fmt.Println(tv.Field(0).CanSet())  // true  — exported
fmt.Println(tv.Field(1).CanSet())  // false — unexported
```

**Slice and map operations:**
```go
// Slice:
s := []int{1, 2, 3}
sv := reflect.ValueOf(s)
fmt.Println(sv.Len())           // 3
fmt.Println(sv.Index(0).Int())  // 1
sv.Index(0).SetInt(99)
fmt.Println(s[0])               // 99

// Map:
m := map[string]int{"a": 1, "b": 2}
mv := reflect.ValueOf(m)
fmt.Println(mv.Len())   // 2
key := reflect.ValueOf("a")
fmt.Println(mv.MapIndex(key).Int())  // 1
mv.SetMapIndex(reflect.ValueOf("c"), reflect.ValueOf(3))
fmt.Println(m)  // map[a:1 b:2 c:3]
```

### Quick Check
> 1. Why do you need to pass a pointer to modify a value via reflection?
> 2. What does `.Elem()` do on a pointer Value?
> 3. Can you set unexported struct fields via reflection?

---

## 4. Structs — Field Inspection and Tags

Struct reflection is the most common real-world use — it powers JSON/YAML/TOML marshaling, ORM column mapping, validation libraries:

```go
type User struct {
    ID       int    `json:"id" db:"user_id" validate:"required"`
    Name     string `json:"name" validate:"required,min=2"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"-" db:"-"`
}

func inspectStruct(v any) {
    t := reflect.TypeOf(v)
    val := reflect.ValueOf(v)
    
    // If pointer, dereference:
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
        val = val.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        fmt.Println("not a struct")
        return
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := val.Field(i)
        
        jsonTag := field.Tag.Get("json")
        dbTag := field.Tag.Get("db")
        validateTag := field.Tag.Get("validate")
        
        fmt.Printf("Field: %-10s Type: %-8s Value: %-12v json: %-10s db: %s validate: %s\n",
            field.Name, field.Type, value.Interface(),
            jsonTag, dbTag, validateTag)
    }
}

inspectStruct(User{ID: 1, Name: "Alice", Email: "alice@example.com"})
```

**A minimal JSON marshaler (to understand how encoding/json works):**
```go
func Marshal(v any) ([]byte, error) {
    t := reflect.TypeOf(v)
    val := reflect.ValueOf(v)
    
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
        val = val.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        return nil, fmt.Errorf("Marshal only supports structs, got %s", t.Kind())
    }
    
    var sb strings.Builder
    sb.WriteByte('{')
    
    first := true
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := val.Field(i)
        
        jsonTag := field.Tag.Get("json")
        if jsonTag == "-" {
            continue  // Skip fields tagged json:"-"
        }
        
        name := field.Name
        if jsonTag != "" {
            // Handle "name,omitempty" format:
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "" {
                name = parts[0]
            }
        }
        
        if !first {
            sb.WriteByte(',')
        }
        first = false
        
        fmt.Fprintf(&sb, "%q:", name)
        
        switch value.Kind() {
        case reflect.String:
            fmt.Fprintf(&sb, "%q", value.String())
        case reflect.Int, reflect.Int64:
            fmt.Fprintf(&sb, "%d", value.Int())
        case reflect.Bool:
            fmt.Fprintf(&sb, "%t", value.Bool())
        default:
            fmt.Fprintf(&sb, "%v", value.Interface())
        }
    }
    
    sb.WriteByte('}')
    return []byte(sb.String()), nil
}
```

### Quick Check
> 1. How do you read a specific tag key like `json` from a struct field?
> 2. Why do real JSON marshalers need to handle the `json:"-"` tag?
> 3. What does `reflect.Value.Interface()` return?

---

## 5. Calling Functions and Methods Dynamically

```go
func add(a, b int) int { return a + b }

// Get the function value:
fn := reflect.ValueOf(add)
fmt.Println(fn.Kind())  // func

// Call it:
args := []reflect.Value{
    reflect.ValueOf(3),
    reflect.ValueOf(4),
}
results := fn.Call(args)
fmt.Println(results[0].Int())  // 7
```

**Calling methods on a struct:**
```go
type Greeter struct{ Name string }
func (g *Greeter) Hello() string { return "Hello, " + g.Name }
func (g *Greeter) Add(x, y int) int { return x + y }

g := &Greeter{Name: "Alice"}
v := reflect.ValueOf(g)

// Call Hello():
method := v.MethodByName("Hello")
results := method.Call(nil)
fmt.Println(results[0].String())  // "Hello, Alice"

// Call Add(3, 4):
addMethod := v.MethodByName("Add")
results = addMethod.Call([]reflect.Value{
    reflect.ValueOf(3),
    reflect.ValueOf(4),
})
fmt.Println(results[0].Int())  // 7
```

**Dynamic dispatch pattern (dependency injection):**
```go
// Call any method by name with arbitrary args:
func callMethod(obj any, methodName string, args ...any) ([]any, error) {
    v := reflect.ValueOf(obj)
    m := v.MethodByName(methodName)
    if !m.IsValid() {
        return nil, fmt.Errorf("method %q not found", methodName)
    }
    
    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }
    
    out := m.Call(in)
    result := make([]any, len(out))
    for i, v := range out {
        result[i] = v.Interface()
    }
    return result, nil
}
```

### Quick Check
> 1. How do you call a function via reflection?
> 2. How do you call a specific method on a struct value?
> 3. What does `MethodByName` return if the method doesn't exist?

---

## 6. Creating Values at Runtime

```go
// Create a new zero value of a type:
t := reflect.TypeOf(User{})
newUser := reflect.New(t)  // *User — like &User{}
newUser.Elem().Field(0).SetInt(42)  // Set ID field

// Create a slice of a dynamic type:
sliceType := reflect.SliceOf(t)  // []User
slice := reflect.MakeSlice(sliceType, 0, 10)
slice = reflect.Append(slice, reflect.ValueOf(User{Name: "Alice"}))

// Create a map of a dynamic type:
mapType := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0))  // map[string]int
m := reflect.MakeMap(mapType)
m.SetMapIndex(reflect.ValueOf("key"), reflect.ValueOf(42))
```

**Populating a struct from a map (like a simple ORM):**
```go
func mapToStruct(data map[string]any, out any) error {
    v := reflect.ValueOf(out)
    if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
        return fmt.Errorf("out must be a pointer to struct")
    }
    
    v = v.Elem()
    t := v.Type()
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fieldVal := v.Field(i)
        
        if !fieldVal.CanSet() {
            continue
        }
        
        key := strings.ToLower(field.Name)
        if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
            key = strings.Split(tag, ",")[0]
        }
        
        val, ok := data[key]
        if !ok {
            continue
        }
        
        fieldVal.Set(reflect.ValueOf(val).Convert(field.Type))
    }
    return nil
}

// Usage:
var user User
mapToStruct(map[string]any{"id": 1, "name": "Alice"}, &user)
```

### Quick Check
> 1. What does `reflect.New(t)` return?
> 2. What is `reflect.SliceOf(t)` useful for?

---

## 7. Real-World Use Cases

**Validation library (like `go-validator`):**
```go
func Validate(v any) []string {
    var errors []string
    t := reflect.TypeOf(v)
    val := reflect.ValueOf(v)
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := val.Field(i)
        tag := field.Tag.Get("validate")
        
        if tag == "" {
            continue
        }
        
        rules := strings.Split(tag, ",")
        for _, rule := range rules {
            switch rule {
            case "required":
                if value.IsZero() {
                    errors = append(errors, field.Name+" is required")
                }
            case "email":
                if s := value.String(); !strings.Contains(s, "@") {
                    errors = append(errors, field.Name+" must be a valid email")
                }
            default:
                if strings.HasPrefix(rule, "min=") {
                    minStr := strings.TrimPrefix(rule, "min=")
                    min, _ := strconv.Atoi(minStr)
                    if value.Kind() == reflect.String && len(value.String()) < min {
                        errors = append(errors, fmt.Sprintf("%s must be at least %d chars", field.Name, min))
                    }
                }
            }
        }
    }
    return errors
}
```

---

## 8. When to Avoid Reflection

Reflection has real costs:
- **Slow**: 10-100× slower than direct code
- **No compile-time safety**: panics at runtime instead of compile errors
- **Hard to read**: reflection code is complex and hard to reason about

```
Use reflection when:
  - Writing generic serialization/deserialization (JSON, YAML, SQL scan)
  - Writing testing/mocking frameworks
  - Writing dependency injection containers
  - Writing ORM column mapping

Don't use reflection when:
  - Generics solve the problem (Ch 22)
  - Interfaces solve the problem
  - The type is known at compile time
  - Performance is critical (hot path)
```

---

## Summary

- **`reflect.TypeOf(v)`**: returns `reflect.Type` — the type descriptor
- **`reflect.ValueOf(v)`**: returns `reflect.Value` — wraps the value
- **Kind**: underlying category (`int`, `struct`, `slice`, etc.); Type is the named type
- **Modification**: requires addressable value — pass pointer, call `.Elem()`
- **`CanSet()`**: false for unexported fields and non-addressable values
- **Struct tags**: `field.Tag.Get("json")` — access key-specific tag values
- **`reflect.New(t)`**: allocate a new zero value of type `t`, returns a pointer Value
- **Dynamic calls**: `v.MethodByName("X").Call(args)` — for dynamic dispatch
- **Avoid**: when generics or interfaces work; reflection is slow and type-unsafe

---

## Exercises

### Easy
1. Write `TypeInfo(v any)` that prints: the type name, kind, and if it's a struct, all field names and types. Test with `int`, `string`, `[]int`, `map[string]int`, and a custom struct.
2. Write `CopyStruct(src, dst any)` that copies all exported fields with the same name and compatible types from `src` to `dst` using reflection. Test: copy `UserInput{Name: "Alice", Age: 30}` to `UserDTO{Name: string, Age: int}`.
3. Write `IsZero(v any) bool` that returns true if the value is the zero value for its type using reflection. Handle: `0`, `""`, `nil`, `false`, `0.0`, zero struct, nil pointer, nil slice, nil map.

### Medium
4. Struct-to-map converter: Write `StructToMap(v any) map[string]any` that converts a struct to a map using the `json` tag for field names (or lowercase field name if no tag). Handle nested structs by recursing (nested struct becomes nested map). Handle slices of structs. Test round-trip: struct → map → check all values.
5. Simple ORM row scanner: Write `ScanRow(row *sql.Row, dest any) error` that uses reflection to scan a SQL row into a struct. Map column names to struct fields using `db` tags. Handle types: string, int, bool, time.Time, and pointer variants. Test with a real SQLite database (use `database/sql` + `modernc.org/sqlite`).
6. Event dispatcher: Build a reflection-based event dispatcher. `Register(handler any)` accepts any function and registers it for the event type of its first argument. `Dispatch(event any)` calls all registered handlers for that event's type. Example: `Register(func(e UserCreated) { ... })` — when `Dispatch(UserCreated{...})` is called, the handler runs. Verify type safety: dispatching the wrong type doesn't call the handler.

### Hard
7. Mini DI container: Build a dependency injection container using reflection. `Register[T any](factory func(...) T)` registers a factory function. `Resolve[T any]() T` returns an instance. The container automatically resolves dependencies by looking at the factory function's parameter types. Circular dependency detection: if A needs B and B needs A, return an error instead of infinite loop. Test with 5 inter-dependent services.
8. Struct differ: Write `Diff(a, b any) []FieldDiff` that compares two structs of the same type and returns a list of changed fields. Each `FieldDiff` has: `Field string`, `OldValue any`, `NewValue any`, `JSONPath string` (dot-separated for nested fields). Handle: nested structs (recursively), slice fields (report as whole-field change), map fields, pointer fields (nil vs non-nil). Test with deeply nested structs and verify the JSON path is correct.
