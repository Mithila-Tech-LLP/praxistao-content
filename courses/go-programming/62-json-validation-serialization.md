# Chapter 62: JSON, Validation, and Serialization

JSON is the lingua franca of APIs. Go's `encoding/json` is fast and idiomatic, but there are subtle behaviors that trip up most developers. This chapter covers encoding and decoding, struct tags, custom types, validation patterns, and when to reach for third-party libraries.

## Table of Contents

1. [encoding/json Basics](#1-encodingjson-basics)
2. [Struct Tags](#2-struct-tags)
3. [Custom Marshaling](#3-custom-marshaling)
4. [Decoding Safely](#4-decoding-safely)
5. [Validation](#5-validation)
6. [Performance](#6-performance)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. encoding/json Basics

```go
// Encoding: struct → JSON
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

u := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
data, err := json.Marshal(u)
// data = {"id":1,"name":"Alice","email":"alice@example.com"}

// Indented (for debugging/pretty-printing):
data, err = json.MarshalIndent(u, "", "  ")

// Decoding: JSON → struct
var u2 User
err = json.Unmarshal(data, &u2)

// Streaming encode/decode (preferred for HTTP handlers):
json.NewEncoder(w).Encode(u)            // to io.Writer
json.NewDecoder(r.Body).Decode(&u)      // from io.Reader
```

### Key behaviors

```go
// nil vs empty slice
type Response struct {
    Items []string `json:"items"`
}

r1 := Response{Items: nil}       → {"items":null}
r2 := Response{Items: []string{}} → {"items":[]}

// nil pointer vs zero value
type Config struct {
    Timeout *int `json:"timeout"`
}
// nil *int → "timeout":null
// &zero   → "timeout":0
```

---

## 2. Struct Tags

```go
type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"` // omit if zero
    Price       float64   `json:"price"`
    Internal    string    `json:"-"`                     // never marshaled
    CreatedAt   time.Time `json:"created_at"`
}
```

| Tag option | Effect |
|------------|--------|
| `json:"name"` | Use "name" as JSON key |
| `json:"name,omitempty"` | Omit field if zero value |
| `json:"-"` | Never include in JSON |
| `json:",string"` | Marshal number as quoted string |

### Embedding and promotion

```go
type Timestamps struct {
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Note struct {
    ID      int    `json:"id"`
    Content string `json:"content"`
    Timestamps       // embedded: fields promoted to Note's JSON
}

// JSON: {"id":1,"content":"...","created_at":"...","updated_at":"..."}
```

---

## 3. Custom Marshaling

Implement `json.Marshaler` and `json.Unmarshaler` for full control:

```go
// Custom time format
type Date struct{ time.Time }

func (d Date) MarshalJSON() ([]byte, error) {
    return json.Marshal(d.Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil { return err }
    t, err := time.Parse("2006-01-02", s)
    if err != nil { return fmt.Errorf("invalid date %q: %w", s, err) }
    d.Time = t
    return nil
}

// Custom enum
type Status int

const (
    Active Status = iota
    Inactive
    Banned
)

var statusStr = map[Status]string{Active: "active", Inactive: "inactive", Banned: "banned"}
var strStatus = map[string]Status{"active": Active, "inactive": Inactive, "banned": Banned}

func (s Status) MarshalJSON() ([]byte, error) {
    str, ok := statusStr[s]
    if !ok { return nil, fmt.Errorf("unknown status %d", s) }
    return json.Marshal(str)
}

func (s *Status) UnmarshalJSON(data []byte) error {
    var str string
    if err := json.Unmarshal(data, &str); err != nil { return err }
    val, ok := strStatus[str]
    if !ok { return fmt.Errorf("unknown status %q", str) }
    *s = val
    return nil
}
```

### Polymorphic JSON

```go
// Handle a JSON field that can be either a string or a []string
type Tags []string

func (t *Tags) UnmarshalJSON(data []byte) error {
    // Try as a single string first
    var single string
    if err := json.Unmarshal(data, &single); err == nil {
        *t = Tags{single}
        return nil
    }
    // Try as an array
    var arr []string
    if err := json.Unmarshal(data, &arr); err != nil { return err }
    *t = Tags(arr)
    return nil
}
```

---

## 4. Decoding Safely

### Use DisallowUnknownFields

By default, `json.Decoder` silently ignores unknown fields. In API handlers, this hides typos:

```go
func decodeJSON(r *http.Request, dst any) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields() // reject {"namee": "Alice"} — typo
    if err := dec.Decode(dst); err != nil {
        return fmt.Errorf("decode: %w", err)
    }
    // Ensure only one JSON value was sent
    if err := dec.Decode(&struct{}{}); err != io.EOF {
        return fmt.Errorf("request body must only contain a single JSON object")
    }
    return nil
}
```

### Typed decode errors

```go
func handleDecodeError(w http.ResponseWriter, err error) {
    var syntaxErr *json.SyntaxError
    var unmarshalErr *json.UnmarshalTypeError
    
    switch {
    case errors.As(err, &syntaxErr):
        http.Error(w, fmt.Sprintf("malformed JSON at position %d", syntaxErr.Offset), 400)
    case errors.As(err, &unmarshalErr):
        http.Error(w, fmt.Sprintf("wrong type for field %q: expected %s", 
            unmarshalErr.Field, unmarshalErr.Type), 400)
    case errors.Is(err, io.EOF):
        http.Error(w, "empty request body", 400)
    case errors.Is(err, io.ErrUnexpectedEOF):
        http.Error(w, "request body incomplete", 400)
    default:
        http.Error(w, "invalid request body", 400)
    }
}
```

### Null vs absent fields (use pointers)

```go
type UpdateUserRequest struct {
    Name  *string `json:"name"`   // nil = not provided, &"" = clear name
    Email *string `json:"email"`
    Age   *int    `json:"age"`
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    var req UpdateUserRequest
    if err := decodeJSON(r, &req); err != nil { /* 400 */ }
    
    if req.Name != nil  { /* update name  */ }
    if req.Email != nil { /* update email */ }
}
```

---

## 5. Validation

Validation is separate from decoding. Decode first, then validate.

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
    msgs := make([]string, len(e))
    for i, err := range e { msgs[i] = err.Error() }
    return strings.Join(msgs, "; ")
}

// Validate a CreateUserRequest
type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Age      int    `json:"age"`
}

func (r *CreateUserRequest) Validate() error {
    var errs ValidationErrors

    name := strings.TrimSpace(r.Name)
    if name == "" {
        errs = append(errs, ValidationError{"name", "required"})
    } else if len(name) > 100 {
        errs = append(errs, ValidationError{"name", "must be 100 characters or fewer"})
    }

    if r.Email == "" {
        errs = append(errs, ValidationError{"email", "required"})
    } else if !strings.Contains(r.Email, "@") {
        errs = append(errs, ValidationError{"email", "invalid email address"})
    }

    if len(r.Password) < 8 {
        errs = append(errs, ValidationError{"password", "must be at least 8 characters"})
    }

    if r.Age < 0 || r.Age > 150 {
        errs = append(errs, ValidationError{"age", "must be between 0 and 150"})
    }

    if len(errs) > 0 { return errs }
    return nil
}

// Handler pattern:
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := decodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    if err := req.Validate(); err != nil {
        var ve ValidationErrors
        if errors.As(err, &ve) {
            writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
                "error":  "validation failed",
                "fields": ve,
            })
            return
        }
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    // proceed...
}
```

### Using go-playground/validator

For complex validation rules, `github.com/go-playground/validator/v10` provides declarative struct tags:

```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Name     string `json:"name"     validate:"required,min=1,max=100"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age"      validate:"min=0,max=150"`
}

var validate = validator.New()

func validateRequest(req any) error {
    if err := validate.Struct(req); err != nil {
        var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            errs := make(ValidationErrors, len(ve))
            for i, fe := range ve {
                errs[i] = ValidationError{
                    Field:   strings.ToLower(fe.Field()),
                    Message: validationMessage(fe),
                }
            }
            return errs
        }
        return err
    }
    return nil
}

func validationMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required": return "required"
    case "email":    return "invalid email address"
    case "min":      return fmt.Sprintf("must be at least %s", fe.Param())
    case "max":      return fmt.Sprintf("must be at most %s", fe.Param())
    default:         return fmt.Sprintf("failed %s validation", fe.Tag())
    }
}
```

---

## 6. Performance

### Benchmark: encoder vs streaming

```go
// For small responses: json.Marshal is fine
data, _ := json.Marshal(response)
w.Header().Set("Content-Type", "application/json")
w.Write(data)

// For large responses or streaming: use encoder directly
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)  // avoids full in-memory buffer
```

### json/v2 (experimental)

Go 1.25 ships an experimental `encoding/json/v2` (enabled with `GOEXPERIMENT=jsonv2`) with breaking improvements:
```go
import "encoding/json/v2"

// Stricter by default: duplicate keys error, case-sensitive field matching
// Cleaner API: json.Marshal / json.Unmarshal work the same way
// Better performance, especially for decoding large payloads
```

### jsoniter (third-party)

A drop-in replacement if the standard library shows up in your profiles:
```go
import jsoniter "github.com/json-iterator/go"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Drop-in replacement, 2-3x faster than standard library
```

---

## Summary

- `json.Marshal` / `json.Unmarshal` — simple encode/decode
- `json.NewEncoder(w).Encode` — stream to any writer; preferred for HTTP
- `dec.DisallowUnknownFields()` — reject typos in request bodies
- Null vs absent: use `*string` — nil means not provided
- Decode errors: use `errors.As` to give clients useful messages
- Validation: separate from decoding — decode first, then validate all fields at once
- Custom types: implement `MarshalJSON`/`UnmarshalJSON` for full control

## Exercises

### Easy
1. Implement a `Money` type that serializes as `{"amount": 10.50, "currency": "USD"}` from a `struct{ Cents int; Currency string }`. Custom `MarshalJSON` divides Cents by 100.
2. Write a `decodeBody[T any](r *http.Request) (T, error)` generic helper that decodes JSON, disallows unknown fields, and checks for trailing input. Use it in three different handlers.
3. Implement `omitempty` behavior manually: write a struct where zero-value fields are excluded from the JSON output without using the `omitempty` tag — instead use custom `MarshalJSON`.

### Medium
4. Build a **partial update** system: define `Patch[T any]` that wraps a value and tracks whether it was explicitly set in JSON vs absent. `Patch[string]` should distinguish `{}` (absent), `{"name":""}` (set to empty), and `{"name":"Alice"}` (set to Alice). Use this for PATCH endpoints.
5. Implement a `Strict[T any]` type that wraps any decodable struct and returns an error if any unknown JSON field is present in the input. This should work generically without knowing T's fields ahead of time.
6. Write a **JSON schema validator**: given a JSON schema (subset of JSON Schema draft 7), validate a JSON value against it. Support: `type`, `required`, `minLength`, `maxLength`, `minimum`, `maximum`, `pattern`.

### Hard
7. Implement a **streaming JSON array processor**: read a very large JSON array from an `io.Reader` without loading it all into memory. Use `json.Decoder.Token()` to parse one element at a time and process it, keeping memory usage at O(1) regardless of array size.
8. Build a **bidirectional JSON/Go type registry**: register Go types by name, then encode any registered type as `{"type": "TypeName", "data": {...}}`. On decode, use the `"type"` field to instantiate the correct Go type and unmarshal `"data"` into it.
