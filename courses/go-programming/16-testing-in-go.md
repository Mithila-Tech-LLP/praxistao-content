# Chapter 16: Testing in Go

Testing is a first-class citizen in Go — the standard library includes everything you need, and the tooling is excellent. The `go test` command runs tests, benchmarks, and examples. No external testing framework required for most use cases. This chapter teaches you how to write tests that are fast, readable, and actually catch bugs.

## Table of Contents

1. [The testing Package Basics](#1-the-testing-package-basics)
2. [Table-Driven Tests](#2-table-driven-tests)
3. [Subtests and Sub-benchmarks](#3-subtests-and-sub-benchmarks)
4. [Benchmarks](#4-benchmarks)
5. [Test Helpers and Utilities](#5-test-helpers-and-utilities)
6. [Mocking and Interfaces](#6-mocking-and-interfaces)
7. [Test Coverage](#7-test-coverage)
8. [Integration and Example Tests](#8-integration-and-example-tests)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The testing Package Basics

**File naming**: test files end in `_test.go`. They are compiled only during `go test`.

**Function naming**: test functions start with `Test`, take `*testing.T`:
```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

**Running tests:**
```bash
go test ./...              # Run all tests in all packages
go test ./pkg/...          # Run tests in pkg and subpackages
go test -run TestAdd       # Run only tests matching "TestAdd"
go test -v                 # Verbose: print all test names
go test -count=1           # Disable test caching (run fresh every time)
go test -timeout 30s       # Set timeout (default: 10m)
go test -race              # Enable race detector
```

**Failure methods on `*testing.T`:**
```go
t.Error("message")        // Record failure, continue test
t.Errorf("fmt %v", val)   // Record formatted failure, continue test
t.Fatal("message")        // Record failure, STOP test immediately
t.Fatalf("fmt %v", val)   // Record formatted failure, STOP immediately

// Helper — marks the calling function as a test helper (better error line numbers):
t.Helper()

// Skip — mark test as skipped (not failed):
t.Skip("reason")
t.Skipf("skip: %v", reason)

// Log — print only when test fails or -v flag is set:
t.Log("debug info")
t.Logf("value: %v", val)
```

**Example test file:**
```go
// user_test.go
package user

import (
    "testing"
)

func TestNewUser(t *testing.T) {
    u := NewUser("Alice", "alice@example.com")
    
    if u.Name != "Alice" {
        t.Errorf("Name = %q; want %q", u.Name, "Alice")
    }
    if u.Email != "alice@example.com" {
        t.Errorf("Email = %q; want %q", u.Email, "alice@example.com")
    }
    if u.CreatedAt.IsZero() {
        t.Error("CreatedAt should be set")
    }
}

func TestUserActivate(t *testing.T) {
    u := NewUser("Bob", "bob@example.com")
    if u.IsActive {
        t.Error("new user should not be active")
    }
    
    u.Activate()
    if !u.IsActive {
        t.Error("user should be active after Activate()")
    }
}
```

**Package organization — black-box vs white-box testing:**
```go
// White-box: package user — can access unexported fields
package user

func TestInternalLogic(t *testing.T) {
    u := &User{id: generateID()}  // Can access unexported field
}

// Black-box: package user_test — tests the public API only
package user_test

import "yourmodule/user"

func TestPublicAPI(t *testing.T) {
    u := user.NewUser("Alice", "alice@example.com")
    // Can only access exported fields/methods
}
```

### Quick Check
> 1. What naming convention do test files and test functions use?
> 2. What is the difference between `t.Error` and `t.Fatal`?
> 3. What does `t.Helper()` do?

---

## 2. Table-Driven Tests

Table-driven tests are Go's idiom for testing multiple cases with the same logic:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", 5, -3, 2},
        {"overflow check", math.MaxInt32, 1, math.MaxInt32 + 1},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

**Testing error cases:**
```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
        errMsg  string
    }{
        {"valid", 10, 2, 5, false, ""},
        {"divide by zero", 10, 0, 0, true, "division by zero"},
        {"negative", -10, 2, -5, false, ""},
        {"both negative", -10, -2, 5, false, ""},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            
            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                    return
                }
                if tt.errMsg != "" && err.Error() != tt.errMsg {
                    t.Errorf("error = %q; want %q", err.Error(), tt.errMsg)
                }
                return
            }
            
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Errorf("Divide(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

**Run specific table test:**
```bash
go test -run TestDivide/divide_by_zero  # Slashes separate test/subtest
go test -run TestDivide                 # Run all subtests of TestDivide
```

**Testing with custom types:**
```go
func TestProcessOrder(t *testing.T) {
    tests := []struct {
        name    string
        order   Order
        user    User
        want    *Receipt
        wantErr error
    }{
        {
            name:    "valid order",
            order:   Order{Items: []Item{{ID: 1, Qty: 2}}, Total: 29.99},
            user:    User{ID: 1, Balance: 100.0},
            want:    &Receipt{OrderID: 1, Amount: 29.99},
            wantErr: nil,
        },
        {
            name:    "insufficient funds",
            order:   Order{Items: []Item{{ID: 1, Qty: 2}}, Total: 200.0},
            user:    User{ID: 1, Balance: 50.0},
            want:    nil,
            wantErr: ErrInsufficientFunds,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessOrder(tt.order, tt.user)
            
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("err = %v; want %v", err, tt.wantErr)
            }
            if tt.want != nil && got == nil {
                t.Error("expected receipt, got nil")
            }
        })
    }
}
```

### Quick Check
> 1. What is the advantage of table-driven tests over separate test functions?
> 2. How do you run only one case from a table-driven test?
> 3. Why use `t.Run` inside table tests?

---

## 3. Subtests and Sub-benchmarks

**Subtests with `t.Run`** — grouping related tests:
```go
func TestUserService(t *testing.T) {
    // Setup shared for all subtests:
    svc := NewUserService(testDB)
    
    t.Run("Create", func(t *testing.T) {
        user, err := svc.Create("Alice", "alice@example.com")
        if err != nil {
            t.Fatalf("Create: %v", err)
        }
        if user.ID == 0 {
            t.Error("user should have an ID after creation")
        }
    })
    
    t.Run("Get", func(t *testing.T) {
        user, err := svc.Get(1)
        if err != nil {
            t.Fatalf("Get: %v", err)
        }
        if user.Name != "Alice" {
            t.Errorf("Name = %q; want Alice", user.Name)
        }
    })
    
    t.Run("Delete", func(t *testing.T) {
        if err := svc.Delete(1); err != nil {
            t.Fatalf("Delete: %v", err)
        }
        _, err := svc.Get(1)
        if !errors.Is(err, ErrNotFound) {
            t.Errorf("expected ErrNotFound after delete, got %v", err)
        }
    })
}
```

**Parallel subtests — run subtests concurrently:**
```go
func TestParallelOperations(t *testing.T) {
    tests := []struct {
        name  string
        input int
        want  int
    }{
        {"a", 1, 2},
        {"b", 5, 10},
        {"c", 10, 20},
    }
    
    for _, tt := range tests {
        tt := tt  // Capture loop variable!
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // This subtest runs in parallel with other parallel subtests
            got := double(tt.input)
            if got != tt.want {
                t.Errorf("double(%d) = %d; want %d", tt.input, got, tt.want)
            }
        })
    }
}
```

**TestMain — setup and teardown for entire package:**
```go
func TestMain(m *testing.M) {
    // Setup:
    db, err := setupTestDB()
    if err != nil {
        log.Fatalf("setup: %v", err)
    }
    testDB = db
    
    // Run all tests:
    code := m.Run()
    
    // Teardown:
    db.Close()
    
    os.Exit(code)
}
```

### Quick Check
> 1. What does `t.Parallel()` do?
> 2. Why must you capture loop variables with `tt := tt` before calling `t.Parallel()`?
> 3. When would you use `TestMain`?

---

## 4. Benchmarks

Benchmarks measure performance. They run with `go test -bench`:

```go
// BenchmarkXxx(b *testing.B)
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {  // b.N is determined automatically by the framework
        Add(100, 200)
    }
}
```

```bash
go test -bench=.                    # Run all benchmarks
go test -bench=BenchmarkAdd         # Run specific benchmark
go test -bench=. -benchmem          # Include memory allocation stats
go test -bench=. -benchtime=5s      # Run for 5 seconds (default: 1s)
go test -bench=. -count=5           # Run 5 times for stable results
```

**Sample output:**
```
BenchmarkAdd-8       1000000000    0.316 ns/op
BenchmarkConcat-8       5000000    243   ns/op    32 B/op    1 allocs/op
```
- `-8`: ran with 8 CPU cores (GOMAXPROCS=8)
- `1000000000`: b.N — how many times the benchmark ran
- `0.316 ns/op`: time per operation
- `32 B/op`: bytes allocated per operation
- `1 allocs/op`: heap allocations per operation

**Benchmark setup — put setup outside the loop:**
```go
func BenchmarkSortSlice(b *testing.B) {
    // Setup outside the loop — not counted in timing:
    original := generateRandomSlice(10000)
    
    b.ResetTimer()  // Reset timer after setup
    
    for i := 0; i < b.N; i++ {
        s := make([]int, len(original))
        copy(s, original)        // Reset to unsorted state
        sort.Ints(s)
    }
}
```

**Benchmark different sizes:**
```go
func BenchmarkSearch(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            data := generateSortedSlice(size)
            target := data[size/2]
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                binarySearch(data, target)
            }
        })
    }
}
```

**Memory benchmarks:**
```go
func BenchmarkStringBuilder(b *testing.B) {
    b.ReportAllocs()  // Same as -benchmem for this benchmark
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        for j := 0; j < 100; j++ {
            sb.WriteString("hello")
        }
        _ = sb.String()
    }
}
```

### Quick Check
> 1. What is `b.N` in a benchmark?
> 2. Why put `b.ResetTimer()` after setup code?
> 3. What does `-benchmem` show?

---

## 5. Test Helpers and Utilities

**`t.Helper()`** — mark functions as helpers so errors point to the caller:
```go
// Without t.Helper() — error points to line inside assertNoError, not the caller
// With t.Helper() — error points to the TEST's line that called assertNoError

func assertNoError(t *testing.T, err error) {
    t.Helper()  // Mark this as a helper
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v; want %v", got, want)
    }
}

func assertError(t *testing.T, err error, target error) {
    t.Helper()
    if !errors.Is(err, target) {
        t.Errorf("err = %v; want %v", err, target)
    }
}

// Usage in tests:
func TestGetUser(t *testing.T) {
    user, err := store.Get(1)
    assertNoError(t, err)             // Error line points HERE, not inside assertNoError
    assertEqual(t, user.Name, "Alice") // Same
}
```

**Temporary directories and files:**
```go
func TestWriteFile(t *testing.T) {
    // t.TempDir() creates a temp dir that's automatically cleaned up:
    dir := t.TempDir()
    
    path := filepath.Join(dir, "test.txt")
    err := WriteFile(path, "hello world")
    if err != nil {
        t.Fatalf("WriteFile: %v", err)
    }
    
    content, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("ReadFile: %v", err)
    }
    if string(content) != "hello world" {
        t.Errorf("content = %q; want %q", content, "hello world")
    }
}
```

**Test fixtures:**
```go
// testdata/ directory — conventionally holds test fixtures:
// testdata/
//   input.json
//   expected_output.json

func TestProcessJSON(t *testing.T) {
    input, err := os.ReadFile("testdata/input.json")
    if err != nil {
        t.Fatalf("reading fixture: %v", err)
    }
    
    expected, err := os.ReadFile("testdata/expected_output.json")
    if err != nil {
        t.Fatalf("reading fixture: %v", err)
    }
    
    got, err := ProcessJSON(input)
    if err != nil {
        t.Fatalf("ProcessJSON: %v", err)
    }
    
    if !bytes.Equal(got, expected) {
        t.Errorf("output mismatch\ngot:  %s\nwant: %s", got, expected)
    }
}
```

**Golden file testing:**
```go
// -update flag to regenerate golden files:
var update = flag.Bool("update", false, "update golden files")

func TestRenderHTML(t *testing.T) {
    got := RenderHTML(testData)
    
    goldenPath := "testdata/golden/render.html"
    
    if *update {
        os.WriteFile(goldenPath, got, 0644)
        return
    }
    
    want, err := os.ReadFile(goldenPath)
    if err != nil {
        t.Fatalf("reading golden file: %v", err)
    }
    
    if !bytes.Equal(got, want) {
        t.Errorf("output differs from golden file\nrun with -update to regenerate")
    }
}
```

### Quick Check
> 1. What is `t.TempDir()` and why is it better than creating directories manually?
> 2. What is the `testdata/` directory used for?
> 3. What is a "golden file" test?

---

## 6. Mocking and Interfaces

Go's interface system makes mocking straightforward — define an interface for any dependency, then provide a test implementation:

```go
// Production code uses an interface:
type EmailSender interface {
    SendEmail(to, subject, body string) error
}

type UserService struct {
    email EmailSender
    db    *sql.DB
}

func (s *UserService) Register(name, email string) error {
    // ... create user ...
    return s.email.SendEmail(email, "Welcome!", "Hello "+name)
}

// Real implementation (production):
type SMTPSender struct{ host string }
func (s *SMTPSender) SendEmail(to, subject, body string) error { ... }

// Mock for tests:
type MockEmailSender struct {
    SentEmails []struct{ To, Subject, Body string }
    ReturnErr  error
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
    m.SentEmails = append(m.SentEmails, struct{ To, Subject, Body string }{to, subject, body})
    return m.ReturnErr
}

// Test:
func TestRegister(t *testing.T) {
    mock := &MockEmailSender{}
    svc := &UserService{email: mock, db: testDB}
    
    err := svc.Register("Alice", "alice@example.com")
    if err != nil {
        t.Fatalf("Register: %v", err)
    }
    
    if len(mock.SentEmails) != 1 {
        t.Errorf("expected 1 email, got %d", len(mock.SentEmails))
    }
    if mock.SentEmails[0].To != "alice@example.com" {
        t.Errorf("email To = %q; want alice@example.com", mock.SentEmails[0].To)
    }
    
    // Test error case:
    mock.ReturnErr = errors.New("SMTP server down")
    err = svc.Register("Bob", "bob@example.com")
    if err == nil {
        t.Error("expected error when email fails, got nil")
    }
}
```

**httptest — mock HTTP servers:**
```go
import "net/http/httptest"

func TestFetchUser(t *testing.T) {
    // Create a test HTTP server:
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/users/42" {
            t.Errorf("unexpected path: %s", r.URL.Path)
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]any{"id": 42, "name": "Alice"})
    }))
    defer server.Close()
    
    // Point your client at the test server:
    client := &UserClient{BaseURL: server.URL}
    user, err := client.FetchUser(42)
    if err != nil {
        t.Fatalf("FetchUser: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("Name = %q; want Alice", user.Name)
    }
}
```

**httptest.Recorder — test HTTP handlers:**
```go
func TestUserHandler(t *testing.T) {
    handler := NewUserHandler(testStore)
    
    req := httptest.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    
    handler.ServeHTTP(w, req)
    
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("status = %d; want 200", resp.StatusCode)
    }
    
    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    if user.Name != "Alice" {
        t.Errorf("Name = %q; want Alice", user.Name)
    }
}
```

### Quick Check
> 1. Why do we define interfaces in production code rather than in tests?
> 2. What does `httptest.NewServer` do?
> 3. What does `httptest.NewRecorder` do?

---

## 7. Test Coverage

```bash
go test -cover ./...                      # Show coverage percentage
go test -coverprofile=coverage.out ./...  # Save coverage data
go tool cover -html=coverage.out          # Open HTML coverage report
go tool cover -func=coverage.out          # Show per-function coverage
```

**Reading coverage output:**
```
ok      mypackage    0.123s  coverage: 87.4% of statements
```

**Viewing HTML report shows:**
- Green: covered lines
- Red: uncovered lines
- Grey: not measurable (declarations, comments)

**Writing tests to improve coverage:**
```go
// If coverage shows Divide's error case is uncovered:
func TestDivide_ByZero(t *testing.T) {
    _, err := Divide(10, 0)
    if err == nil {
        t.Error("expected error for division by zero")
    }
}
```

**Coverage targets:** 80%+ is a common goal. 100% is rarely worth chasing — focus on covering critical paths and edge cases, not just line count.

### Quick Check
> 1. How do you generate an HTML coverage report?
> 2. Is 100% test coverage always the goal?
> 3. What does a red line in the coverage report mean?

---

## 8. Integration and Example Tests

**Integration tests** — test real interactions with databases, HTTP services, etc:

```go
// Only run when -integration flag is set:
func TestUserRepositoryIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    db, err := setupTestPostgres(t)
    if err != nil {
        t.Fatalf("setup: %v", err)
    }
    
    repo := NewUserRepository(db)
    
    user, err := repo.Create(User{Name: "Alice", Email: "alice@example.com"})
    if err != nil {
        t.Fatalf("Create: %v", err)
    }
    
    got, err := repo.Get(user.ID)
    if err != nil {
        t.Fatalf("Get: %v", err)
    }
    if got.Email != "alice@example.com" {
        t.Errorf("Email = %q; want alice@example.com", got.Email)
    }
}
```

```bash
go test -short ./...   # Skip tests that call t.Skip in short mode
go test -v ./...       # Verbose — also shows which tests were skipped
```

**Example tests** — serve as both documentation and tests:
```go
// Examples appear in godoc and run as tests:
func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output:
    // 5
}

func ExampleUser_String() {
    u := User{Name: "Alice", Email: "alice@example.com"}
    fmt.Println(u.String())
    // Output:
    // Alice <alice@example.com>
}
```

The `// Output:` comment is verified — if the actual output doesn't match, the test fails. No `// Output:` comment means the example compiles but doesn't verify output.

### Quick Check
> 1. How do you skip slow integration tests during normal development?
> 2. What makes an example test actually test something?
> 3. What happens if the output of an example doesn't match the `// Output:` comment?

---

## Summary

- **Test files**: `xxx_test.go`, functions: `TestXxx(t *testing.T)`
- **Failure**: `t.Error`/`t.Errorf` (continue), `t.Fatal`/`t.Fatalf` (stop)
- **Table-driven**: `tests := []struct{...}{...}` + `for _, tt := range tests { t.Run(tt.name, ...) }`
- **Subtests**: `t.Run(name, func)` — enables running single cases, parallel execution
- **Benchmarks**: `BenchmarkXxx(b *testing.B)` + `for i := 0; i < b.N; i++`
- **Helpers**: `t.Helper()` for better error line numbers; `t.TempDir()` for temp directories
- **Mocking**: define interfaces in production code, provide mock implementations in tests
- **httptest**: `httptest.NewServer` for HTTP client tests, `httptest.NewRecorder` for handler tests
- **Coverage**: `go test -cover` + `go tool cover -html=coverage.out`
- **Examples**: `ExampleXxx()` with `// Output:` comment — tested + shown in godoc

---

## Exercises

### Easy
1. Write three tests for a `fibonacci(n int) int` function: test `fib(0)=0`, `fib(1)=1`, `fib(10)=55`. Use table-driven tests.
2. Write a benchmark for `fibonacci` comparing recursive vs iterative implementations. Run with `-benchmem`. Which allocates more?
3. Write an example test for a `Reverse(s string) string` function that documents the expected behavior in godoc format.

### Medium
4. Complete test suite: Implement `Stack[T]` with `Push(v T)`, `Pop() (T, bool)`, `Peek() (T, bool)`, `Len() int`, `IsEmpty() bool`. Write comprehensive table-driven tests covering: empty stack pop, single element, LIFO ordering, concurrent access (use `t.Parallel()`). Achieve at least 95% coverage.
5. HTTP handler tests: Write a `ProductHandler` with GET `/products/{id}` and POST `/products`. Mock the `ProductStore` interface. Test: successful get, 404 for missing, 400 for bad JSON on create, 201 on successful create. Use `httptest.NewRecorder`. Verify response status codes AND response body.
6. Integration test with Docker: Write an integration test for a PostgreSQL repository. Use `testing.Short()` to skip in short mode. Set up the schema in `TestMain`, use a transaction per test (rollback after each test so tests don't interfere), and test Create/Get/Update/Delete operations.

### Hard
7. Test framework: Build a minimal assertion library `assert` package with: `Equal[T comparable](t, got, want T, msg ...string)`, `NotNil(t, v any)`, `Error(t, err error)`, `ErrorIs(t, err, target error)`, `Panic(t, fn func())`. Each should use `t.Helper()` appropriately and format errors in a "got X; want Y" format. Write tests for the assert package itself (meta-testing).
8. Fuzz testing: Implement a `ParseURL(s string) (*URL, error)` function that parses custom URL-like strings of the form `scheme://host:port/path`. Write a fuzz test (`FuzzParseURL`) that: feeds random bytes to the parser, verifies that if parsing succeeds, re-serializing the result and parsing again gives the same result (round-trip property), and verifies that the parser never panics on any input. Run with `go test -fuzz=Fuzz -fuzztime=30s`.
