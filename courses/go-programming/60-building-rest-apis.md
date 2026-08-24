# Chapter 60: Building REST APIs

REST (Representational State Transfer) is the dominant architectural style for web APIs. Go's `net/http` is sufficient for REST APIs, but the `chi` router adds middleware chaining, URL parameters, and sub-routers that make complex APIs manageable. This chapter builds a complete CRUD API with proper error handling, validation, and testing.

## Table of Contents

1. [REST Design Principles](#1-rest-design-principles)
2. [Setting Up chi](#2-setting-up-chi)
3. [Request Validation](#3-request-validation)
4. [Complete CRUD Example](#4-complete-crud-example)
5. [API Versioning](#5-api-versioning)
6. [Pagination and Filtering](#6-pagination-and-filtering)
7. [Error Handling Strategy](#7-error-handling-strategy)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. REST Design Principles

**Resource naming:**
```
Collection:    GET  /users           → list all users
              POST  /users           → create a new user
Single item:   GET  /users/{id}      → get user by ID
               PUT  /users/{id}      → replace user (full update)
             PATCH  /users/{id}      → update user fields (partial update)
            DELETE  /users/{id}      → delete user
Nested:        GET  /users/{id}/posts → list posts for user
               GET  /users/{id}/posts/{postID} → specific post of user
Action:       POST  /users/{id}/activate → action that doesn't fit CRUD
```

**HTTP method semantics:**
- `GET`: safe (no side effects) + idempotent
- `PUT`: idempotent (calling twice = same result as calling once)
- `DELETE`: idempotent
- `POST`: neither safe nor idempotent (each call creates a new resource)
- `PATCH`: not necessarily idempotent (depends on implementation)

**Response shapes:**
```json
// Collection:
{"data": [...], "total": 100, "page": 1, "pageSize": 20}

// Single item:
{"data": {...}}

// Error:
{"error": "user not found", "code": "USER_NOT_FOUND"}

// Validation error:
{"error": "validation failed", "fields": {"email": "invalid format", "name": "required"}}
```

---

## 2. Setting Up chi

```bash
go get github.com/go-chi/chi/v5
```

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()

    // Built-in middleware:
    r.Use(middleware.RequestID)    // Adds X-Request-ID header
    r.Use(middleware.RealIP)       // Sets r.RemoteAddr to X-Forwarded-For
    r.Use(middleware.Logger)       // Structured access log
    r.Use(middleware.Recoverer)    // Recover from panics
    r.Use(middleware.Compress(5))  // Gzip responses
    r.Use(middleware.Timeout(30 * time.Second))

    // Routes:
    r.Get("/health", healthHandler)

    r.Route("/api/v1", func(r chi.Router) {
        r.Route("/users", func(r chi.Router) {
            r.Get("/", listUsers)
            r.Post("/", createUser)

            r.Route("/{userID}", func(r chi.Router) {
                r.Use(userCtx)  // Middleware: load user into context
                r.Get("/", getUser)
                r.Put("/", updateUser)
                r.Delete("/", deleteUser)
            })
        })
    })

    http.ListenAndServe(":8080", r)
}

// Extract URL parameter:
func getUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")
    // ...
}
```

---

## 3. Request Validation

```go
package api

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
)

// ValidationError holds field-level errors.
type ValidationError struct {
    Fields map[string]string `json:"fields"`
}

func (e *ValidationError) Error() string {
    parts := make([]string, 0, len(e.Fields))
    for k, v := range e.Fields { parts = append(parts, k+": "+v) }
    return "validation failed: " + strings.Join(parts, "; ")
}

func (e *ValidationError) Add(field, msg string) {
    if e.Fields == nil { e.Fields = make(map[string]string) }
    e.Fields[field] = msg
}

func (e *ValidationError) HasErrors() bool { return len(e.Fields) > 0 }

// DecodeJSON decodes JSON body with size limit and unknown field rejection.
func DecodeJSON(r *http.Request, dst any) error {
    r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)  // 1 MB limit
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(dst); err != nil {
        var syntaxErr *json.SyntaxError
        var unmarshalErr *json.UnmarshalTypeError
        var maxBytesErr *http.MaxBytesError
        switch {
        case errors.As(err, &syntaxErr):
            return fmt.Errorf("malformed JSON at position %d", syntaxErr.Offset)
        case errors.As(err, &unmarshalErr):
            return fmt.Errorf("field %q: expected %s, got %s",
                unmarshalErr.Field, unmarshalErr.Type, unmarshalErr.Value)
        case errors.As(err, &maxBytesErr):
            return fmt.Errorf("request body too large (max 1 MB)")
        default:
            return err
        }
    }
    // Ensure no trailing data:
    if dec.More() { return fmt.Errorf("request body must contain a single JSON object") }
    return nil
}

// CreateUserRequest with inline validation:
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func (r CreateUserRequest) Validate() *ValidationError {
    v := &ValidationError{}
    if strings.TrimSpace(r.Name) == "" {
        v.Add("name", "required")
    } else if len(r.Name) > 100 {
        v.Add("name", "must be at most 100 characters")
    }
    if !strings.Contains(r.Email, "@") {
        v.Add("email", "invalid email format")
    }
    if r.Age < 0 || r.Age > 150 {
        v.Add("age", "must be between 0 and 150")
    }
    if v.HasErrors() { return v }
    return nil
}
```

---

## 4. Complete CRUD Example

### Model and Store
```go
// models/user.go
package models

import "time"

type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// store/users.go
package store

import (
    "errors"
    "sort"
    "sync"
    "time"

    "myapi/models"
)

var ErrNotFound = errors.New("not found")
var ErrEmailTaken = errors.New("email already taken")

type UserStore struct {
    mu      sync.RWMutex
    users   map[int64]*models.User
    nextID  int64
    byEmail map[string]int64
}

func NewUserStore() *UserStore {
    return &UserStore{
        users:   make(map[int64]*models.User),
        byEmail: make(map[string]int64),
        nextID:  1,
    }
}

func (s *UserStore) Create(name, email string, age int) (*models.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.byEmail[email]; exists {
        return nil, ErrEmailTaken
    }

    now := time.Now().UTC()
    u := &models.User{
        ID: s.nextID, Name: name, Email: email, Age: age,
        CreatedAt: now, UpdatedAt: now,
    }
    s.users[s.nextID] = u
    s.byEmail[email] = s.nextID
    s.nextID++
    return u, nil
}

func (s *UserStore) GetByID(id int64) (*models.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    u, ok := s.users[id]
    if !ok { return nil, ErrNotFound }
    cp := *u  // Return copy
    return &cp, nil
}

func (s *UserStore) List(offset, limit int) ([]*models.User, int) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    all := make([]*models.User, 0, len(s.users))
    for _, u := range s.users { cp := *u; all = append(all, &cp) }

    // Sort by ID for stable pagination:
    sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })

    total := len(all)
    if offset >= total { return []*models.User{}, total }
    end := offset + limit
    if end > total { end = total }
    return all[offset:end], total
}

func (s *UserStore) Update(id int64, name string, age int) (*models.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    u, ok := s.users[id]
    if !ok { return nil, ErrNotFound }

    if name != "" { u.Name = name }
    if age > 0 { u.Age = age }
    u.UpdatedAt = time.Now().UTC()

    cp := *u
    return &cp, nil
}

func (s *UserStore) Delete(id int64) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    u, ok := s.users[id]
    if !ok { return ErrNotFound }
    delete(s.byEmail, u.Email)
    delete(s.users, id)
    return nil
}
```

### Handler Layer
```go
// handlers/users.go
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "myapi/api"
    "myapi/store"
)

type UserHandler struct {
    store *store.UserStore
}

func NewUserHandler(s *store.UserStore) *UserHandler {
    return &UserHandler{store: s}
}

func (h *UserHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Route("/{userID}", func(r chi.Router) {
        r.Get("/", h.Get)
        r.Patch("/", h.Update)
        r.Delete("/", h.Delete)
    })
    return r
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 { page = 1 }
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
    if pageSize < 1 || pageSize > 100 { pageSize = 20 }

    users, total := h.store.List((page-1)*pageSize, pageSize)
    writeJSON(w, http.StatusOK, map[string]any{
        "data":     users,
        "total":    total,
        "page":     page,
        "pageSize": pageSize,
    })
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req api.CreateUserRequest
    if err := api.DecodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    if ve := req.Validate(); ve != nil {
        writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
            "error":  "validation failed",
            "fields": ve.Fields,
        })
        return
    }

    user, err := h.store.Create(req.Name, req.Email, req.Age)
    if err != nil {
        if errors.Is(err, store.ErrEmailTaken) {
            writeError(w, http.StatusConflict, "email already taken")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to create user")
        return
    }

    w.Header().Set("Location", "/api/v1/users/"+strconv.FormatInt(user.ID, 10))
    writeJSON(w, http.StatusCreated, map[string]any{"data": user})
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user ID")
        return
    }
    user, err := h.store.GetByID(id)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to get user")
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"data": user})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user ID")
        return
    }

    var req struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }
    if err := api.DecodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    user, err := h.store.Update(id, req.Name, req.Age)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to update user")
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"data": user})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid user ID")
        return
    }
    if err := h.store.Delete(id); err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to delete user")
        return
    }
    w.WriteHeader(http.StatusNoContent)  // 204 — no body
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}
```

---

## 5. API Versioning

```go
// URL versioning (most common):
r.Route("/api/v1", v1Routes)
r.Route("/api/v2", v2Routes)

// Header versioning:
func versionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        version := r.Header.Get("API-Version")
        if version == "" { version = "v1" }
        ctx := context.WithValue(r.Context(), apiVersionKey, version)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 6. Pagination and Filtering

```go
type ListParams struct {
    Page     int
    PageSize int
    SortBy   string
    SortDir  string
    Search   string
}

func parseListParams(r *http.Request) ListParams {
    q := r.URL.Query()
    p := ListParams{
        Page:     max(1, queryInt(q, "page", 1)),
        PageSize: clamp(queryInt(q, "pageSize", 20), 1, 100),
        SortBy:   q.Get("sortBy"),
        SortDir:  q.Get("sortDir"),
        Search:   q.Get("search"),
    }
    if p.SortDir != "desc" { p.SortDir = "asc" }
    return p
}

func queryInt(q url.Values, key string, defaultVal int) int {
    v, err := strconv.Atoi(q.Get(key))
    if err != nil { return defaultVal }
    return v
}

func clamp(v, min, max int) int {
    if v < min { return min }
    if v > max { return max }
    return v
}

// Cursor-based pagination (better for large datasets):
type CursorPage struct {
    Items      any    `json:"items"`
    NextCursor string `json:"nextCursor,omitempty"`
    HasMore    bool   `json:"hasMore"`
}
// Cursor = base64(lastID + lastCreatedAt) — opaque to clients
```

---

## 7. Error Handling Strategy

```go
// Structured API error:
type APIError struct {
    Status  int    `json:"-"`
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

var (
    ErrNotFound   = &APIError{Status: 404, Code: "NOT_FOUND", Message: "resource not found"}
    ErrUnauthorized = &APIError{Status: 401, Code: "UNAUTHORIZED", Message: "authentication required"}
    ErrForbidden  = &APIError{Status: 403, Code: "FORBIDDEN", Message: "access denied"}
)

// Central error handler:
func handleError(w http.ResponseWriter, err error) {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        writeJSON(w, apiErr.Status, apiErr)
        return
    }
    // Log unexpected errors — don't expose internals:
    slog.Error("unexpected error", "error", err)
    writeJSON(w, http.StatusInternalServerError, &APIError{
        Code: "INTERNAL_ERROR", Message: "an unexpected error occurred",
    })
}
```

---

## Summary

- REST: resources as nouns in URLs, HTTP methods as verbs, status codes for outcome
- `chi` adds method-based routing, URL parameters, middleware chains, and sub-routers over `net/http`
- Validate at the boundary: decode JSON, check required fields, return structured `422` errors
- Layer architecture: Handler → Store; handlers never touch data directly
- Pagination: page+pageSize for simple cases; cursor-based for large or real-time data
- `204 No Content` for successful DELETE — no body needed
- Centralize error mapping: translate domain errors (`ErrNotFound`) to HTTP status codes in one place

---

## Exercises

### Easy
1. Add a `GET /users/{id}/exists` endpoint that returns `{"exists": true/false}` with `200 OK` in both cases (never `404`). This pattern is useful for form validation ("is this email taken?") without leaking user IDs.
2. Add `GET /users/search?q=alice` that returns users whose name contains the query (case-insensitive). Add the `search` field to the `List` store method. Return `{"data": [], "total": 0}` for no matches.
3. Implement `PATCH /users/{id}` that accepts any subset of `{name, age, email}` fields. Only update the provided non-zero fields. Test that sending `{}` doesn't change anything, and that sending `{"age": 0}` also doesn't change age (since 0 is the zero value).

### Medium
4. Add **request ID middleware**: generate a UUID (use `crypto/rand` to build a 16-byte random ID) if `X-Request-ID` header is absent, otherwise use the provided value. Store it in context, include it in every response header, and log it with every error. Show that the ID propagates through all log lines for a single request.
5. Implement **optimistic locking**: add `version int` to the User model. `PATCH /users/{id}` requires `{"version": N, ...fields}`. If the stored version != N, return `409 Conflict` with `{"error": "conflict", "currentVersion": M}`. This prevents lost updates when two clients update simultaneously.
6. Build a **bulk import endpoint**: `POST /users/bulk` accepts a JSON array of up to 100 users. Process them all, collect both successes and failures. Return `207 Multi-Status` with a per-item result: `[{"status": 201, "data": {...}}, {"status": 422, "error": "...", "index": 1}, ...]`. Don't abort on first error.

### Hard
7. Implement **HATEOAS links**: each user response includes a `_links` object with `self`, `update`, `delete` URLs. Collection responses include `_links.self`, `_links.next`, `_links.prev`. Write a `LinkBuilder` that generates URLs from the current request's host and scheme, so it works behind a proxy.
8. Add **ETag / conditional requests**: compute `ETag` as `fmt.Sprintf("%d-%d", user.ID, user.UpdatedAt.Unix())`. Set it on `GET /users/{id}` responses. Handle `If-None-Match` header: if client's ETag matches, return `304 Not Modified` with no body. Handle `If-Match` on `PUT`/`PATCH`: if the ETag doesn't match, return `412 Precondition Failed`.
