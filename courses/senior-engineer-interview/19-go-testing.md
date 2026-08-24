# Chapter 19: Go Testing — Table Tests, Mocks, Subtests & Fuzzing

Testing is production code. A senior engineer writes tests that are readable, maintainable, and actually catch bugs. This chapter covers the full spectrum of Go testing: from unit tests to integration tests, with the patterns and tools that production Go codebases use.

## Table of Contents

1. [Table-Driven Tests — The Go Idiom](#1-table-driven-tests--the-go-idiom)
2. [Subtests and t.Run](#2-subtests)
3. [Test Helpers and t.Helper](#3-test-helpers)
4. [Interface-Based Mocking](#4-interface-based-mocking)
5. [Testing HTTP Handlers](#5-testing-http-handlers)
6. [Integration Tests with testcontainers](#6-integration-tests-with-testcontainers)
7. [Fuzz Testing](#7-fuzz-testing)
8. [Test Coverage](#8-test-coverage)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. Table-Driven Tests — The Go Idiom

Table-driven tests are the idiomatic Go pattern. You define a slice of test cases, iterate over them, and run each sub-test. This reduces boilerplate and makes adding new cases trivial.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero and positive", 0, 5, 5},
        {"both zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Why Table-Driven Tests Are Best Practice

```
go test -run TestAdd              # run all cases
go test -run TestAdd/zero_and_positive  # run specific case

# Output on failure shows exactly which case failed:
--- FAIL: TestAdd (0.00s)
    --- FAIL: TestAdd/both_zero (0.00s)
        add_test.go:22: Add(0, 0) = 1, want 0
```

### Advanced Table Test: Testing a Function with Multiple Return Values

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Config
        wantErr bool
    }{
        {
            name:  "valid config",
            input: `{"port": 8080}`,
            want:  Config{Port: 8080},
        },
        {
            name:    "invalid JSON",
            input:   `{invalid}`,
            wantErr: true,
        },
        {
            name:    "empty input",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseConfig(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("ParseConfig() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 2. Subtests

`t.Run` creates a subtest with its own name, pass/fail state, and the ability to run in parallel.

```go
func TestUserService(t *testing.T) {
    // Shared setup for all subtests
    db := setupTestDB(t)
    svc := NewUserService(db)

    t.Run("CreateUser", func(t *testing.T) {
        user, err := svc.Create(context.Background(), CreateUserInput{
            Name:  "Alice",
            Email: "alice@example.com",
        })
        if err != nil { t.Fatal(err) }
        if user.ID == "" { t.Error("user ID should not be empty") }
    })

    t.Run("GetUser", func(t *testing.T) {
        user, err := svc.Get(context.Background(), "test-id")
        if err != nil { t.Fatal(err) }
        if user.Name != "Alice" { t.Errorf("want Alice, got %s", user.Name) }
    })
}
```

### Parallel Subtests

```go
func TestConcurrent(t *testing.T) {
    tests := []struct{ name, input string }{...}

    for _, tt := range tests {
        tt := tt // MUST capture loop variable for parallel tests
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // run all subtests in parallel
            result := slowFunction(tt.input)
            // assertions...
        })
    }
}
// Parallel subtests run faster in CI
```

---

## 3. Test Helpers

`t.Helper()` marks a function as a test helper. When the helper calls `t.Fatal` or `t.Error`, the failure is reported at the call site, not inside the helper.

```go
// WITHOUT t.Helper: error line points to line 5 (inside assertEq)
// WITH t.Helper: error line points to where assertEq was called

func assertEq(t *testing.T, got, want interface{}) {
    t.Helper() // makes t.Errorf report the caller's file/line
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// Testing cleanup with t.Cleanup
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("postgres", testDSN)
    if err != nil { t.Fatal(err) }
    
    t.Cleanup(func() {
        db.Close() // runs when the test finishes, regardless of pass/fail
    })
    return db
}
```

### Using testify for Assertions

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
    result, err := doSomething()
    
    require.NoError(t, err)     // fatal: stops test if error is not nil
    assert.Equal(t, want, result) // non-fatal: continues test on failure
    assert.NotNil(t, result)
    assert.Contains(t, result.Name, "Alice")
    
    // Require vs Assert: use require when subsequent assertions don't make sense
    // if the first one fails (e.g., checking fields of a nil pointer)
}
```

---

## 4. Interface-Based Mocking

Go's interface-based mocking is clean: define an interface, provide a real implementation and a mock implementation.

```go
// Define the interface (usually near where it's used)
type UserRepository interface {
    Create(ctx context.Context, user User) (User, error)
    GetByID(ctx context.Context, id string) (User, error)
    Update(ctx context.Context, user User) error
}

// Real implementation
type PostgresUserRepo struct { db *sql.DB }
func (r *PostgresUserRepo) Create(ctx context.Context, user User) (User, error) {
    // ...DB query...
}

// Mock for testing
type MockUserRepo struct {
    CreateFn  func(ctx context.Context, user User) (User, error)
    GetByIDFn func(ctx context.Context, id string) (User, error)
    UpdateFn  func(ctx context.Context, user User) error
}

func (m *MockUserRepo) Create(ctx context.Context, user User) (User, error) {
    return m.CreateFn(ctx, user)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (User, error) {
    return m.GetByIDFn(ctx, id)
}
func (m *MockUserRepo) Update(ctx context.Context, user User) error {
    return m.UpdateFn(ctx, user)
}

// Test using the mock
func TestUserService_Create(t *testing.T) {
    mock := &MockUserRepo{
        CreateFn: func(ctx context.Context, user User) (User, error) {
            user.ID = "test-id"
            return user, nil
        },
    }
    
    svc := NewUserService(mock)
    result, err := svc.Create(context.Background(), CreateInput{Name: "Alice"})
    
    if err != nil { t.Fatal(err) }
    if result.ID != "test-id" { t.Errorf("got %s, want test-id", result.ID) }
}
```

### Using testify/mock for More Complex Mocking

```go
import "github.com/stretchr/testify/mock"

type MockEmailService struct {
    mock.Mock
}

func (m *MockEmailService) Send(ctx context.Context, to, subject, body string) error {
    args := m.Called(ctx, to, subject, body)
    return args.Error(0)
}

func TestSendWelcomeEmail(t *testing.T) {
    emailSvc := new(MockEmailService)
    emailSvc.On("Send", mock.Anything, "alice@example.com", "Welcome!", mock.Anything).
        Return(nil).Once()
    
    err := sendWelcomeEmail(context.Background(), emailSvc, "alice@example.com")
    
    assert.NoError(t, err)
    emailSvc.AssertExpectations(t) // verifies all expected calls were made
}
```

---

## 5. Testing HTTP Handlers

Use `net/http/httptest` to test HTTP handlers without starting a real server.

```go
import "net/http/httptest"

func TestHandleGetUser(t *testing.T) {
    // Setup mock dependencies
    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id string) (User, error) {
            return User{ID: id, Name: "Alice"}, nil
        },
    }
    handler := NewUserHandler(mockRepo)

    // Create a test request
    req := httptest.NewRequest("GET", "/users/123", nil)
    req.SetPathValue("id", "123") // Go 1.22+ path params

    // Create a response recorder
    w := httptest.NewRecorder()

    // Call the handler
    handler.GetUser(w, req)

    // Assert response
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
    }

    var got User
    json.NewDecoder(resp.Body).Decode(&got)
    if got.Name != "Alice" { t.Errorf("want Alice, got %s", got.Name) }
}

// Testing a full server
func TestServer(t *testing.T) {
    server := httptest.NewServer(setupRouter())
    defer server.Close()

    resp, err := http.Get(server.URL + "/health")
    if err != nil { t.Fatal(err) }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("health check failed: %d", resp.StatusCode)
    }
}
```

---

## 6. Integration Tests with testcontainers

Integration tests run against real external services (PostgreSQL, Redis, Kafka) using Docker containers.

```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestWithRealPostgres(t *testing.T) {
    ctx := context.Background()
    
    // Start a real PostgreSQL container
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil { t.Fatal(err) }
    defer pgContainer.Terminate(ctx)
    
    // Get connection string
    connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
    
    db, err := sql.Open("postgres", connStr)
    if err != nil { t.Fatal(err) }
    defer db.Close()
    
    // Run your tests against the real database
    repo := NewUserRepo(db)
    user, err := repo.Create(ctx, User{Name: "Alice"})
    if err != nil { t.Fatal(err) }
    if user.ID == "" { t.Error("user should have an ID after create") }
}
```

---

## 7. Fuzz Testing

Fuzz testing automatically generates inputs to find edge cases that crash or cause incorrect behavior.

```go
// Fuzz test: Go 1.18+
func FuzzParseJSON(f *testing.F) {
    // Seed corpus: initial interesting inputs
    f.Add(`{}`)
    f.Add(`{"name":"Alice"}`)
    f.Add(`null`)
    
    f.Fuzz(func(t *testing.T, data string) {
        // The fuzzer generates variations of the seed corpus
        // and calls this function with each variation
        result, err := ParseJSON(data)
        if err != nil {
            return // errors are fine — we're looking for panics
        }
        // Verify invariant: if parse succeeded, marshal should succeed too
        _, err = json.Marshal(result)
        if err != nil {
            t.Fatalf("marshal failed after successful parse: input=%q err=%v", data, err)
        }
    })
}

// Run the fuzzer:
// go test -fuzz=FuzzParseJSON -fuzztime=60s ./...
// 
// When a crash is found, it's saved in testdata/fuzz/FuzzParseJSON/
// Run tests normally to reproduce: go test ./...
```

---

## 8. Test Coverage

```bash
# Generate coverage
go test -cover ./...

# Detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # open in browser

# Coverage by function
go tool cover -func=coverage.out

# Set minimum coverage in CI
go test -cover ./... | grep -E "coverage: [0-9]+\.[0-9]+%"
# Fail if below 80%: custom CI script
```

### Coverage ≠ Quality

A function can have 100% line coverage and still have poor tests. Coverage tells you which lines were executed — not whether the right assertions were made. Focus on meaningful tests over high coverage numbers.

---

## 9. Interview Questions & Model Answers

**Q: What is table-driven testing and why is it idiomatic in Go?**

"Table-driven testing defines test cases as data — a slice of structs with inputs and expected outputs — and iterates over them with t.Run. It's idiomatic in Go because it reduces boilerplate, makes adding new test cases trivial (just add a row to the table), and produces clear failure messages that identify exactly which case failed. It encourages thinking about the test cases exhaustively, including edge cases, rather than writing one-off test functions."

**Q: How do you test code that depends on an external service like a database?**

"Two approaches: mocking and integration testing. For unit tests, I define the database interaction behind an interface and inject a mock implementation. This tests the business logic without a real database. For integration tests, I use testcontainers to spin up a real PostgreSQL container — this verifies the SQL queries, migrations, and data mapping actually work. The two levels complement each other: mocks are fast and test logic; integration tests are slower but give confidence that everything works end-to-end."

**Q: What's the difference between t.Error and t.Fatal?**

"t.Error marks the test as failed and continues execution. t.Fatal marks the test as failed and stops execution immediately (calls runtime.Goexit). Use t.Fatal when subsequent assertions would panic or make no sense if the current one fails — for example, checking fields of a pointer that you just asserted is nil. Use t.Error when you want to see all failures in a test run, not just the first one."

---

## Summary

- **Table-driven tests:** define cases as data, iterate with `t.Run`. The Go idiom.
- **Subtests:** `t.Run` gives each case its own name and can run in parallel.
- **t.Helper():** marks a helper function so failures show the caller's file/line.
- **Mocking:** define interfaces, inject implementations. For complex mocks, use `testify/mock`.
- **HTTP testing:** `httptest.NewRequest` + `httptest.NewRecorder` test handlers without a server.
- **Integration tests:** `testcontainers-go` spins up real Docker containers for DB/Redis/Kafka.
- **Fuzz testing:** `f.Fuzz()` generates inputs automatically to find crashes.
- Always run `go test -race ./...` in CI. Coverage measures lines executed, not quality.
