# Chapter 45: Testing HTTP Services

Good HTTP service tests are fast, isolated, and deterministic. Go's `httptest` package makes it possible to test handlers, middleware chains, and full API flows without starting a real server. This chapter covers the full testing spectrum: unit tests for handlers, integration tests for the full stack, and contract tests for API consumers.

## Table of Contents

1. [httptest Basics](#1-httptest-basics)
2. [Testing Handlers](#2-testing-handlers)
3. [Testing Middleware](#3-testing-middleware)
4. [Integration Test Patterns](#4-integration-test-patterns)
5. [Table-Driven API Tests](#5-table-driven-api-tests)
6. [Mocking Dependencies](#6-mocking-dependencies)
7. [Contract Testing](#7-contract-testing)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. httptest Basics

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// httptest.NewRecorder: captures the response without a real network:
func TestHealthHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()

    healthHandler(rec, req)  // Call handler directly

    resp := rec.Result()    // Convert to *http.Response
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }

    body, _ := io.ReadAll(resp.Body)
    if !strings.Contains(string(body), `"status":"ok"`) {
        t.Errorf("unexpected body: %s", body)
    }
}

// httptest.NewServer: starts a real TCP server (for client-side tests):
func TestHTTPClient(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `{"id":1,"name":"Alice"}`)
    }))
    defer server.Close()

    // Test your HTTP client against the real server:
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(server.URL + "/users/1")
    // ...
}
```

---

## 2. Testing Handlers

### Request Building Helpers
```go
// helpers_test.go — shared test utilities

func newJSONRequest(t *testing.T, method, path string, body any) *http.Request {
    t.Helper()
    var r io.Reader
    if body != nil {
        b, err := json.Marshal(body)
        require.NoError(t, err)
        r = bytes.NewReader(b)
    }
    req := httptest.NewRequest(method, path, r)
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    return req
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
    t.Helper()
    require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
    err := json.NewDecoder(rec.Body).Decode(dst)
    require.NoError(t, err)
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
    t.Helper()
    if rec.Code != want {
        t.Errorf("expected status %d, got %d; body: %s", want, rec.Code, rec.Body.String())
    }
}
```

### Handler Unit Tests
```go
func TestCreateUser(t *testing.T) {
    store := store.NewUserStore()
    handler := handlers.NewUserHandler(store)

    t.Run("success", func(t *testing.T) {
        req := newJSONRequest(t, "POST", "/users", map[string]any{
            "name": "Alice", "email": "alice@example.com", "age": 30,
        })
        rec := httptest.NewRecorder()

        handler.Create(rec, req)

        assertStatus(t, rec, http.StatusCreated)

        var resp struct {
            Data struct {
                ID    int64  `json:"id"`
                Name  string `json:"name"`
                Email string `json:"email"`
            } `json:"data"`
        }
        decode(t, rec, &resp)
        assert.Equal(t, "Alice", resp.Data.Name)
        assert.Equal(t, "alice@example.com", resp.Data.Email)
        assert.Greater(t, resp.Data.ID, int64(0))
        assert.Equal(t, "/api/v1/users/1", rec.Header().Get("Location"))
    })

    t.Run("duplicate email", func(t *testing.T) {
        // Create the user first:
        req := newJSONRequest(t, "POST", "/users", map[string]any{
            "name": "Bob", "email": "alice@example.com", "age": 25,
        })
        rec := httptest.NewRecorder()
        handler.Create(rec, req)  // Reuse alice@example.com from above
        assertStatus(t, rec, http.StatusConflict)
    })

    t.Run("validation error — missing name", func(t *testing.T) {
        req := newJSONRequest(t, "POST", "/users", map[string]any{
            "email": "bob@example.com",
        })
        rec := httptest.NewRecorder()
        handler.Create(rec, req)
        assertStatus(t, rec, http.StatusUnprocessableEntity)

        var resp struct {
            Fields map[string]string `json:"fields"`
        }
        decode(t, rec, &resp)
        assert.Contains(t, resp.Fields, "name")
    })

    t.Run("malformed JSON", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/users", strings.NewReader("{invalid json"))
        req.Header.Set("Content-Type", "application/json")
        rec := httptest.NewRecorder()
        handler.Create(rec, req)
        assertStatus(t, rec, http.StatusBadRequest)
    })
}
```

---

## 3. Testing Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    // Build a minimal handler that the middleware wraps:
    inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims, ok := middleware.ClaimsFromContext(r.Context())
        if !ok { t.Error("no claims in context") }
        fmt.Fprintf(w, `{"userID":%d}`, claims.UserID)
    })

    jwtSvc := auth.NewJWTService("test-secret-key-32-bytes-long!!", "test")
    mw := middleware.Authenticate(jwtSvc.ValidateToken)
    handler := mw(inner)

    t.Run("valid token", func(t *testing.T) {
        token, _ := jwtSvc.GenerateToken(42, "alice@example.com", nil, time.Hour)
        req := httptest.NewRequest("GET", "/", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        rec := httptest.NewRecorder()

        handler.ServeHTTP(rec, req)

        assertStatus(t, rec, http.StatusOK)
        assert.Contains(t, rec.Body.String(), `"userID":42`)
    })

    t.Run("missing token", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/", nil)
        rec := httptest.NewRecorder()
        handler.ServeHTTP(rec, req)
        assertStatus(t, rec, http.StatusUnauthorized)
    })

    t.Run("expired token", func(t *testing.T) {
        token, _ := jwtSvc.GenerateToken(42, "alice@example.com", nil, -time.Hour)
        req := httptest.NewRequest("GET", "/", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        rec := httptest.NewRecorder()
        handler.ServeHTTP(rec, req)
        assertStatus(t, rec, http.StatusUnauthorized)
    })
}

func TestRateLimitMiddleware(t *testing.T) {
    calls := 0
    inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        calls++
        w.WriteHeader(http.StatusOK)
    })

    mw := middleware.RateLimit(5)  // 5 req/sec
    handler := mw(inner)

    // First 5 requests should pass:
    for i := 0; i < 5; i++ {
        req := httptest.NewRequest("GET", "/", nil)
        req.RemoteAddr = "127.0.0.1:12345"
        rec := httptest.NewRecorder()
        handler.ServeHTTP(rec, req)
        assert.Equal(t, http.StatusOK, rec.Code)
    }

    // 6th should be rate limited:
    req := httptest.NewRequest("GET", "/", nil)
    req.RemoteAddr = "127.0.0.1:12345"
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)
    assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}
```

---

## 4. Integration Test Patterns

```go
// TestServer bundles the full router + stores for integration tests.
type TestServer struct {
    *httptest.Server
    client *http.Client
    token  string  // Pre-authenticated token
}

func NewTestServer(t *testing.T) *TestServer {
    t.Helper()

    store := store.NewUserStore()
    authSvc := auth.NewAuthService(auth.NewJWTService("test-secret", "test"), store, ...)
    handler := buildRouter(store, authSvc)

    server := &TestServer{
        Server: httptest.NewServer(handler),
        client: &http.Client{Timeout: 5 * time.Second},
    }
    t.Cleanup(server.Close)
    return server
}

func (s *TestServer) Authenticate(t *testing.T, email, password string) {
    t.Helper()
    resp := s.POST(t, "/auth/login", map[string]string{
        "email": email, "password": password,
    })
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    var body struct {
        AccessToken string `json:"accessToken"`
    }
    json.NewDecoder(resp.Body).Decode(&body)
    s.token = body.AccessToken
}

func (s *TestServer) GET(t *testing.T, path string) *http.Response {
    t.Helper()
    req, _ := http.NewRequest("GET", s.URL+path, nil)
    if s.token != "" {
        req.Header.Set("Authorization", "Bearer "+s.token)
    }
    resp, err := s.client.Do(req)
    require.NoError(t, err)
    return resp
}

func (s *TestServer) POST(t *testing.T, path string, body any) *http.Response {
    t.Helper()
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", s.URL+path, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    if s.token != "" {
        req.Header.Set("Authorization", "Bearer "+s.token)
    }
    resp, err := s.client.Do(req)
    require.NoError(t, err)
    return resp
}

// Integration test using TestServer:
func TestUserFlow(t *testing.T) {
    srv := NewTestServer(t)

    // Register:
    resp := srv.POST(t, "/auth/register", map[string]any{
        "name": "Alice", "email": "alice@test.com", "password": "password123",
    })
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // Login:
    srv.Authenticate(t, "alice@test.com", "password123")
    assert.NotEmpty(t, srv.token)

    // Create a resource:
    resp = srv.POST(t, "/api/v1/posts", map[string]any{
        "title": "My Post", "content": "Hello world",
    })
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // List resources:
    resp = srv.GET(t, "/api/v1/posts")
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var posts struct {
        Data  []map[string]any `json:"data"`
        Total int              `json:"total"`
    }
    json.NewDecoder(resp.Body).Decode(&posts)
    assert.Equal(t, 1, posts.Total)
    assert.Equal(t, "My Post", posts.Data[0]["title"])
}
```

---

## 5. Table-Driven API Tests

```go
func TestUsersAPI(t *testing.T) {
    srv := NewTestServer(t)
    // Pre-populate test data:
    srv.POST(t, "/auth/register", map[string]any{
        "name": "Admin", "email": "admin@test.com", "password": "admin123",
    })
    srv.Authenticate(t, "admin@test.com", "admin123")

    tests := []struct {
        name       string
        method     string
        path       string
        body       any
        token      string
        wantStatus int
        checkBody  func(t *testing.T, body []byte)
    }{
        {
            name:       "list users — authenticated",
            method:     "GET",
            path:       "/api/v1/users",
            wantStatus: http.StatusOK,
            checkBody: func(t *testing.T, body []byte) {
                var resp map[string]any
                json.Unmarshal(body, &resp)
                assert.Contains(t, resp, "data")
                assert.Contains(t, resp, "total")
            },
        },
        {
            name:       "list users — unauthenticated",
            method:     "GET",
            path:       "/api/v1/users",
            token:      "none",  // Override — no auth
            wantStatus: http.StatusUnauthorized,
        },
        {
            name:   "create user — invalid email",
            method: "POST",
            path:   "/api/v1/users",
            body:   map[string]any{"name": "Bob", "email": "notanemail", "age": 25},
            wantStatus: http.StatusUnprocessableEntity,
        },
        {
            name:   "create user — success",
            method: "POST",
            path:   "/api/v1/users",
            body:   map[string]any{"name": "Carol", "email": "carol@test.com", "age": 28},
            wantStatus: http.StatusCreated,
        },
        {
            name:       "get nonexistent user",
            method:     "GET",
            path:       "/api/v1/users/99999",
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tc := range tests {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            var bodyReader io.Reader
            if tc.body != nil {
                b, _ := json.Marshal(tc.body)
                bodyReader = bytes.NewReader(b)
            }
            req, _ := http.NewRequest(tc.method, srv.URL+tc.path, bodyReader)
            if tc.body != nil { req.Header.Set("Content-Type", "application/json") }
            if tc.token == "none" {
                // No auth
            } else {
                req.Header.Set("Authorization", "Bearer "+srv.token)
            }

            resp, err := srv.client.Do(req)
            require.NoError(t, err)
            defer resp.Body.Close()

            assert.Equal(t, tc.wantStatus, resp.StatusCode)

            if tc.checkBody != nil {
                body, _ := io.ReadAll(resp.Body)
                tc.checkBody(t, body)
            }
        })
    }
}
```

---

## 6. Mocking Dependencies

```go
// Define store as an interface — enables mocking:
type UserStore interface {
    Create(ctx context.Context, name, email string, age int) (*models.User, error)
    GetByID(ctx context.Context, id int64) (*models.User, error)
    List(ctx context.Context, offset, limit int) ([]*models.User, int, error)
    Delete(ctx context.Context, id int64) error
}

// Mock implementation for tests:
type MockUserStore struct {
    users  map[int64]*models.User
    nextID int64
    // Control what errors are returned:
    CreateErr   error
    GetByIDErr  error
}

func NewMockUserStore() *MockUserStore {
    return &MockUserStore{users: make(map[int64]*models.User), nextID: 1}
}

func (m *MockUserStore) Create(_ context.Context, name, email string, age int) (*models.User, error) {
    if m.CreateErr != nil { return nil, m.CreateErr }
    u := &models.User{ID: m.nextID, Name: name, Email: email, Age: age}
    m.users[m.nextID] = u
    m.nextID++
    return u, nil
}

func (m *MockUserStore) GetByID(_ context.Context, id int64) (*models.User, error) {
    if m.GetByIDErr != nil { return nil, m.GetByIDErr }
    u, ok := m.users[id]
    if !ok { return nil, store.ErrNotFound }
    return u, nil
}

// Test with a mock that injects errors:
func TestGetUserHandlerDBError(t *testing.T) {
    mockStore := NewMockUserStore()
    mockStore.GetByIDErr = errors.New("connection refused")

    handler := handlers.NewUserHandler(mockStore)

    req := httptest.NewRequest("GET", "/users/1", nil)
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
        URLParams: chi.RouteParams{Keys: []string{"userID"}, Values: []string{"1"}},
    }))
    rec := httptest.NewRecorder()

    handler.Get(rec, req)

    assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
```

---

## 7. Contract Testing

```go
// Contract test: verify API response shape matches what consumers expect.
// Using a snapshot-based approach:

func TestAPIContracts(t *testing.T) {
    srv := NewTestServer(t)
    // Setup...

    contracts := []struct {
        name     string
        path     string
        snapshot string  // Golden file name
    }{
        {"user list", "/api/v1/users?page=1&pageSize=2", "user_list.json"},
        {"user create", "...", "user_create.json"},
    }

    for _, c := range contracts {
        c := c
        t.Run(c.name, func(t *testing.T) {
            resp := srv.GET(t, c.path)
            body, _ := io.ReadAll(resp.Body)

            // Normalize (remove timestamps, IDs) before comparing:
            normalized := normalizeResponse(body)

            golden := filepath.Join("testdata", c.snapshot)
            if *update {  // go test -update to regenerate snapshots
                os.WriteFile(golden, normalized, 0644)
                return
            }

            expected, _ := os.ReadFile(golden)
            assert.JSONEq(t, string(expected), string(normalized))
        })
    }
}

var update = flag.Bool("update", false, "update golden files")

func normalizeResponse(body []byte) []byte {
    var m map[string]any
    json.Unmarshal(body, &m)
    // Replace volatile fields:
    if data, ok := m["data"].(map[string]any); ok {
        data["id"] = 0
        data["createdAt"] = "0001-01-01T00:00:00Z"
        data["updatedAt"] = "0001-01-01T00:00:00Z"
    }
    result, _ := json.MarshalIndent(m, "", "  ")
    return result
}
```

---

## Summary

- `httptest.NewRecorder` + handler direct call = fast, zero-network handler tests
- `httptest.NewServer` = real TCP server for testing HTTP clients
- `TestServer` wrapper bundles router + dependencies; use `t.Cleanup` for teardown
- Table-driven tests cover all status code paths: success, validation, auth, not found, DB error
- Interface-based stores enable mock injection for error path testing
- Golden file / snapshot tests catch API contract regressions
- Always test: 200/201 success, 400 bad JSON, 422 validation, 401 auth, 404 not found, 500 DB error

---

## Exercises

### Easy
1. Write tests for all five CRUD endpoints from Chapter 41 (`List`, `Create`, `Get`, `Update`, `Delete`). Cover at minimum: success case and error cases (not found, duplicate, invalid input). Achieve 100% line coverage on the handler functions.
2. Write a test for the `Logger` middleware that verifies: (1) it calls the next handler, (2) it logs the correct status code after the handler runs. Capture log output using a `bytes.Buffer` wrapped in a `slog.Handler`.
3. Add a `TestMain` to your test file that seeds the in-memory store with 5 users before all tests run, and clears it after. Verify that isolation between test cases works even when sharing the store.

### Medium
4. Build a **test fixture system**: `fixtures.go` contains functions like `CreateUser(t, store, overrides...)` that create test data with sensible defaults, overridable with functional options. This avoids repetition in tests. The fixture returns the created entity so tests can reference it by ID.
5. Implement **parallel integration tests** without test interference: each test case in `TestUsersAPI` should run with its own isolated `TestServer` (fresh store). Use `t.Parallel()`. Verify that running tests in parallel doesn't cause races using `go test -race`.
6. Add **HTTP response time assertions**: any endpoint should respond in < 50ms. Build `assertFastResponse(t, resp, maxDuration)`. Run this in a `BenchmarkUsersAPI` that measures handler throughput under concurrent load using `b.RunParallel`.

### Hard
7. **Generate API client from tests**: build a test helper that records every request/response pair during tests to a JSON file. Post-process the recordings to generate a Go client library (`client/client.go`) with typed methods for each endpoint. This is a simplified version of what tools like `kin-openapi` do.
8. **Chaos testing**: build a `ChaosMiddleware` that randomly injects failures: 5% of requests get a 500 response, 2% get a 50-500ms extra delay. Wrap your test server with this middleware and verify your retry logic (from Ch 42 Hard #6) handles it correctly — eventually succeeds, doesn't exceed timeout.
