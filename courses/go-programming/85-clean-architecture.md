# Chapter 85: Clean Architecture in Go

Clean Architecture, popularized by Robert Martin, organizes code into concentric layers: domain (innermost), use cases, interfaces, and infrastructure (outermost). The key rule: **inner layers know nothing about outer layers**. The domain doesn't import from the database. Use cases don't import from HTTP frameworks.

## Table of Contents

1. [The Layers](#1-the-layers)
2. [Project Layout](#2-project-layout)
3. [Domain Layer](#3-domain-layer)
4. [Use Case Layer](#4-use-case-layer)
5. [Interface Adapters](#5-interface-adapters)
6. [Infrastructure](#6-infrastructure)
7. [Wiring](#7-wiring)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Layers

```
┌─────────────────────────────────────┐
│           Infrastructure            │  HTTP, PostgreSQL, Redis, S3, Kafka
├─────────────────────────────────────┤
│         Interface Adapters          │  HTTP handlers, gRPC servers, CLI
├─────────────────────────────────────┤
│           Use Cases                 │  Application business rules
├─────────────────────────────────────┤
│             Domain                  │  Entities, value objects, domain logic
└─────────────────────────────────────┘
```

### Dependency Rule

Dependencies point **inward only**. The database package imports from the domain package. The domain package imports from nothing in your codebase.

```
infrastructure → adapters → use cases → domain
```

Why this matters:
- You can test use cases with in-memory repositories (no database)
- You can switch PostgreSQL to MongoDB without touching use cases
- Business logic is readable without infrastructure noise

---

## 2. Project Layout

```
myapp/
├── domain/
│   ├── user.go        ← Entity + domain logic
│   ├── order.go
│   └── repository.go  ← Repository interfaces
├── usecase/
│   ├── user_service.go
│   └── order_service.go
├── adapter/
│   ├── http/
│   │   ├── user_handler.go
│   │   └── order_handler.go
│   └── grpc/
│       └── order_server.go
├── infrastructure/
│   ├── postgres/
│   │   ├── user_repo.go
│   │   └── order_repo.go
│   ├── redis/
│   │   └── cache.go
│   └── email/
│       └── smtp_sender.go
└── cmd/
    └── api/
        └── main.go    ← Wire everything together
```

The `cmd/` directory is the only place where all layers are imported together.

---

## 3. Domain Layer

The domain layer has no external imports except the standard library.

```go
// domain/user.go
package domain

import (
    "errors"
    "regexp"
    "strings"
    "time"
)

var (
    ErrNotFound       = errors.New("not found")
    ErrAlreadyExists  = errors.New("already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrUnauthorized   = errors.New("unauthorized")
)

// User is a domain entity — contains business logic, not just data
type User struct {
    ID           int64
    Email        string
    PasswordHash string
    Name         string
    Role         Role
    Active       bool
    CreatedAt    time.Time
}

type Role string
const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

// Domain logic lives on the entity
func (u *User) Activate() error {
    if u.Active { return errors.New("user already active") }
    u.Active = true
    return nil
}

func (u *User) ChangeRole(newRole Role, actor *User) error {
    if actor.Role != RoleAdmin {
        return fmt.Errorf("%w: only admins can change roles", ErrUnauthorized)
    }
    u.Role = newRole
    return nil
}

// Value object: Email validated at construction time
type Email struct{ value string }

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewEmail(s string) (Email, error) {
    s = strings.ToLower(strings.TrimSpace(s))
    if !emailRegex.MatchString(s) {
        return Email{}, fmt.Errorf("%w: invalid email address", ErrInvalidInput)
    }
    return Email{value: s}, nil
}

func (e Email) String() string { return e.value }
```

```go
// domain/repository.go
package domain

import "context"

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

// Other interfaces the domain needs from the outside world
type PasswordHasher interface {
    Hash(password string) (string, error)
    Check(hash, password string) bool
}

type EmailSender interface {
    Send(ctx context.Context, to, subject, body string) error
}
```

---

## 4. Use Case Layer

Use cases orchestrate domain objects and coordinate between repositories and external services. They encode **application business rules** — not domain rules, not HTTP concerns.

```go
// usecase/user_service.go
package usecase

import (
    "context"
    "fmt"
    "time"

    "myapp/domain"
)

type UserService struct {
    users   domain.UserRepository
    hasher  domain.PasswordHasher
    mailer  domain.EmailSender
}

func NewUserService(
    users domain.UserRepository,
    hasher domain.PasswordHasher,
    mailer domain.EmailSender,
) *UserService {
    return &UserService{users: users, hasher: hasher, mailer: mailer}
}

type RegisterInput struct {
    Name     string
    Email    string
    Password string
}

type RegisterOutput struct {
    UserID int64
    Email  string
}

func (s *UserService) Register(ctx context.Context, in RegisterInput) (*RegisterOutput, error) {
    // Validate email as a value object (domain rule)
    email, err := domain.NewEmail(in.Email)
    if err != nil { return nil, err }
    
    // Check for existing user (use case rule)
    if _, err := s.users.GetByEmail(ctx, email.String()); err == nil {
        return nil, fmt.Errorf("%w: email already registered", domain.ErrAlreadyExists)
    }
    
    // Hash password (use case delegates to infrastructure via interface)
    hash, err := s.hasher.Hash(in.Password)
    if err != nil { return nil, fmt.Errorf("hash password: %w", err) }
    
    // Create domain entity
    user := &domain.User{
        Email:        email.String(),
        PasswordHash: hash,
        Name:         in.Name,
        Role:         domain.RoleUser,
        Active:       true,
        CreatedAt:    time.Now(),
    }
    
    // Persist
    if err := s.users.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    // Side effect: send welcome email (use case rule)
    go s.mailer.Send(context.Background(), user.Email,
        "Welcome!", "Thanks for signing up.")
    
    return &RegisterOutput{UserID: user.ID, Email: user.Email}, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
    user, err := s.users.GetByEmail(ctx, email)
    if err != nil {
        // Return generic error — don't leak "user not found" vs "wrong password"
        return nil, fmt.Errorf("%w: invalid credentials", domain.ErrUnauthorized)
    }
    
    if !s.hasher.Check(user.PasswordHash, password) {
        return nil, fmt.Errorf("%w: invalid credentials", domain.ErrUnauthorized)
    }
    return user, nil
}
```

### Testing use cases — no database needed

```go
func TestRegister(t *testing.T) {
    repo := memory.NewUserRepository()
    hasher := &bcryptHasher{}
    mailer := &noopMailer{}
    svc := usecase.NewUserService(repo, hasher, mailer)
    
    out, err := svc.Register(context.Background(), usecase.RegisterInput{
        Name:     "Alice",
        Email:    "alice@example.com",
        Password: "password123",
    })
    
    require.NoError(t, err)
    assert.NotZero(t, out.UserID)
    assert.Equal(t, "alice@example.com", out.Email)
}

type noopMailer struct{}
func (n *noopMailer) Send(_ context.Context, _, _, _ string) error { return nil }
```

---

## 5. Interface Adapters

Adapters translate between the external world (HTTP, gRPC, CLI) and use cases.

```go
// adapter/http/user_handler.go
package httphandler

import (
    "encoding/json"
    "errors"
    "net/http"

    "myapp/domain"
    "myapp/usecase"
)

type UserHandler struct {
    svc *usecase.UserService
}

func NewUserHandler(svc *usecase.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

type registerRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    
    out, err := h.svc.Register(r.Context(), usecase.RegisterInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        // Translate domain errors → HTTP status codes here (not in use cases)
        switch {
        case errors.Is(err, domain.ErrAlreadyExists):
            writeError(w, http.StatusConflict, "email already registered")
        case errors.Is(err, domain.ErrInvalidInput):
            writeError(w, http.StatusBadRequest, err.Error())
        default:
            writeError(w, http.StatusInternalServerError, "registration failed")
        }
        return
    }
    
    writeJSON(w, http.StatusCreated, map[string]any{
        "user_id": out.UserID,
        "email":   out.Email,
    })
}
```

---

## 6. Infrastructure

Infrastructure implementations satisfy the interfaces defined in the domain layer.

```go
// infrastructure/bcrypt/hasher.go
package bcrypt

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

type Hasher struct{ cost int }

func New(cost int) *Hasher { return &Hasher{cost: cost} }

func (h *Hasher) Hash(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    if err != nil { return "", fmt.Errorf("bcrypt hash: %w", err) }
    return string(hash), nil
}

func (h *Hasher) Check(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// infrastructure/smtp/sender.go
package smtp

type Sender struct {
    host string
    port int
    from string
}

func (s *Sender) Send(ctx context.Context, to, subject, body string) error {
    // actual SMTP send
    return nil
}
```

---

## 7. Wiring

```go
// cmd/api/main.go
package main

func main() {
    // Infrastructure
    db := mustConnectDB(os.Getenv("DATABASE_URL"))
    
    // Repositories (infrastructure layer)
    userRepo := postgres.NewUserRepository(db)
    
    // Services (infrastructure layer)
    hasher := bcrypt.New(12)
    mailer := smtp.New(config.SMTPHost, config.SMTPPort, config.FromEmail)
    
    // Use cases (business layer — infrastructure injected via interfaces)
    userSvc := usecase.NewUserService(userRepo, hasher, mailer)
    
    // Handlers (adapter layer)
    userHandler := httphandler.NewUserHandler(userSvc)
    
    // Router
    r := chi.NewRouter()
    r.Post("/register", userHandler.Register)
    
    http.ListenAndServe(":8080", r)
}
```

---

## Summary

| Layer | Contains | Imports |
|-------|----------|---------|
| Domain | Entities, value objects, domain errors, repository interfaces | stdlib only |
| Use Cases | Application logic, orchestration | domain |
| Adapters | HTTP handlers, gRPC servers | domain + use cases |
| Infrastructure | PostgreSQL, Redis, SMTP | domain (satisfies interfaces) |
| cmd/ | Wiring | all layers |

**The test isolation payoff**: domain and use case tests never touch a database. Infrastructure tests can use testcontainers. Adapter tests mock the use case interface.

## Exercises

### Easy
1. Add a `ChangePassword(ctx, userID int64, oldPwd, newPwd string) error` use case method. It should verify the old password, validate the new one (min 8 chars), and update the hash.
2. Create a `memory.UserRepository` that satisfies the `domain.UserRepository` interface for tests. Test `Register` using it — verify that registering with the same email twice returns `ErrAlreadyExists`.
3. Add a domain method `User.Deactivate(reason string)` and a use case `DeactivateUser(ctx, actorID, targetID int64) error`. Deactivation by anyone other than an admin or the user themselves should return `ErrUnauthorized`.

### Medium
4. Extract `domain.Order` with fields `ID`, `UserID`, `Items []OrderItem`, `Status`, `Total`. Add domain methods `Order.Submit()` and `Order.Cancel()` that enforce state transitions (can only submit Draft, can only cancel Pending). Write `usecase.OrderService.PlaceOrder`.
5. Add a **domain event** system: when a user is registered, emit a `UserRegisteredEvent`. The use case publishes to an `EventBus` interface. Write two subscribers: one that sends the welcome email, one that creates a default cart. Test that both subscribers are called.
6. Build a **CLI adapter** alongside the HTTP adapter: `cmd/cli/main.go` that uses the same use case layer to register users via command line flags. Verify that the use case is unchanged.

### Hard
7. Implement the full Clean Architecture for an e-commerce checkout flow: `cart`, `inventory`, `order`, `payment` domains. Each has its own repository. The `CheckoutService` use case must deduct inventory, create an order, and charge payment atomically — write a cross-aggregate transaction using the Unit of Work pattern.
8. Write an **architecture conformance test**: using `go/ast` or `golang.org/x/tools/go/packages`, verify that the `domain` package has no imports that contain `infrastructure`, `adapter`, or `cmd` in their path. Fail the test if the dependency rule is violated.
