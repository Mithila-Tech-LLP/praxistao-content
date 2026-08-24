# Chapter 87: Dependency Injection with Wire

Dependency Injection (DI) means passing dependencies in from outside rather than creating them inside. Wire is Google's compile-time DI generator for Go — it reads your constructor functions and generates the wiring code automatically.

## Table of Contents

1. [DI Without Wire](#1-di-without-wire)
2. [Wire Concepts](#2-wire-concepts)
3. [Wire in Practice](#3-wire-in-practice)
4. [Provider Sets](#4-provider-sets)
5. [Testing with DI](#5-testing-with-di)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. DI Without Wire

Manual wiring gets unwieldy as the dependency graph grows:

```go
// Manual wiring in main.go — readable but tedious for 20+ dependencies
func main() {
    cfg := config.Load()
    
    db := postgres.Connect(cfg.DatabaseURL)
    redisClient := redis.Connect(cfg.RedisURL)
    
    userRepo := postgres.NewUserRepository(db)
    sessionStore := redisStore.NewSessionStore(redisClient)
    hasher := bcrypt.NewHasher(12)
    mailer := smtp.NewSender(cfg.SMTP)
    
    userSvc := usecase.NewUserService(userRepo, hasher, mailer)
    sessionSvc := usecase.NewSessionService(sessionStore)
    
    userHandler := httphandler.NewUserHandler(userSvc, sessionSvc)
    
    // 20 more lines like this...
}
```

Problems:
- Adding a new dependency means updating the constructor AND `main.go`
- Easy to forget to pass a dependency; compile error appears at the wrong layer
- Hard to see which components depend on what

---

## 2. Wire Concepts

Wire has two key concepts:

**Providers**: constructors that return a type. Wire reads them.

```go
// wire.go (your providers — regular Go functions)
func NewUserRepository(db *sqlx.DB) *postgres.UserRepository {
    return postgres.NewUserRepository(db)
}

func NewUserService(
    repo *postgres.UserRepository,
    hasher *bcrypt.Hasher,
    mailer *smtp.Sender,
) *usecase.UserService {
    return usecase.NewUserService(repo, hasher, mailer)
}
```

**Injectors**: a `wire.Build` call that tells Wire what to create.

```go
// wire/wire.go — Wire reads this file; it generates wire_gen.go
//go:build wireinject

func InitializeApp(cfg Config) (*App, error) {
    wire.Build(
        NewDB,
        NewRedis,
        NewUserRepository,
        NewHasher,
        NewSMTPSender,
        NewUserService,
        NewUserHandler,
        NewApp,
    )
    return &App{}, nil
}
```

Running `wire ./wire/` generates `wire_gen.go` with the full constructor chain.

---

## 3. Wire in Practice

Install Wire:
```bash
go install github.com/google/wire/cmd/wire@latest
```

### Step 1: Define your types and constructors

```go
// config/config.go
type Config struct {
    DatabaseURL string
    RedisURL    string
    BCryptCost  int
    SMTP        smtp.Config
    Port        int
}

func Load() (Config, error) {
    return Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        RedisURL:    os.Getenv("REDIS_URL"),
        BCryptCost:  12,
        Port:        8080,
    }, nil
}

// postgres/db.go
type DB struct{ *sqlx.DB }

func NewDB(cfg config.Config) (*DB, func(), error) {
    db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
    if err != nil { return nil, nil, err }
    
    cleanup := func() { db.Close() }
    return &DB{db}, cleanup, nil
}

// usecase/user_service.go
func NewUserService(
    repo *postgres.UserRepository,
    hasher *bcrypt.Hasher,
    mailer domain.EmailSender,
) *UserService {
    return &UserService{users: repo, hasher: hasher, mailer: mailer}
}
```

### Step 2: Write the injector

```go
// wire/wire.go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    "myapp/adapter/http"
    "myapp/config"
    "myapp/infrastructure/bcrypt"
    "myapp/infrastructure/postgres"
    "myapp/infrastructure/smtp"
    "myapp/infrastructure/redis"
    "myapp/usecase"
)

type App struct {
    Server *httpserver.Server
}

func InitializeApp(cfg config.Config) (*App, func(), error) {
    wire.Build(
        // Infrastructure
        postgres.NewDB,
        redis.NewClient,
        
        // Repositories
        postgres.NewUserRepository,
        postgres.NewOrderRepository,
        redis.NewSessionStore,
        
        // Services
        bcrypt.NewHasher,
        smtp.NewSender,
        
        // Use cases
        usecase.NewUserService,
        usecase.NewOrderService,
        
        // Adapters
        httphandler.NewUserHandler,
        httphandler.NewOrderHandler,
        httpserver.New,
        
        wire.Struct(new(App), "Server"),
    )
    return &App{}, nil, nil
}
```

### Step 3: Run Wire

```bash
wire ./wire/
```

Generated `wire_gen.go`:

```go
// wire_gen.go — auto-generated, don't edit by hand
//go:build !wireinject

func InitializeApp(cfg config.Config) (*App, func(), error) {
    db, cleanup, err := postgres.NewDB(cfg)
    if err != nil { return nil, nil, err }
    
    redisClient, cleanup2 := redis.NewClient(cfg)
    
    userRepo := postgres.NewUserRepository(db)
    orderRepo := postgres.NewOrderRepository(db)
    sessionStore := redis.NewSessionStore(redisClient)
    hasher := bcrypt.NewHasher(cfg.BCryptCost)
    sender := smtp.NewSender(cfg.SMTP)
    
    userSvc := usecase.NewUserService(userRepo, hasher, sender)
    orderSvc := usecase.NewOrderService(orderRepo, userRepo)
    
    userHandler := httphandler.NewUserHandler(userSvc, sessionStore)
    orderHandler := httphandler.NewOrderHandler(orderSvc)
    server := httpserver.New(cfg.Port, userHandler, orderHandler)
    
    app := &App{Server: server}
    cleanup3 := func() { cleanup2(); cleanup() }
    return app, cleanup3, nil
}
```

### Step 4: Use in main.go

```go
// cmd/api/main.go
func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatal(err) }
    
    app, cleanup, err := wire.InitializeApp(cfg)
    if err != nil { log.Fatal(err) }
    defer cleanup()
    
    app.Server.ListenAndServe()
}
```

---

## 4. Provider Sets

Group related providers with `wire.NewSet` so you can reuse them across multiple applications:

```go
// infrastructure/postgres/wire.go
var ProviderSet = wire.NewSet(
    NewDB,
    NewUserRepository,
    NewOrderRepository,
    NewProductRepository,
)

// infrastructure/redis/wire.go
var ProviderSet = wire.NewSet(
    NewClient,
    NewSessionStore,
    NewRateLimiter,
)

// usecase/wire.go
var ProviderSet = wire.NewSet(
    NewUserService,
    NewOrderService,
    NewProductService,
)

// Wire injector uses the sets
func InitializeApp(cfg config.Config) (*App, func(), error) {
    wire.Build(
        postgres.ProviderSet,
        redis.ProviderSet,
        bcrypt.ProviderSet,
        smtp.ProviderSet,
        usecase.ProviderSet,
        httphandler.ProviderSet,
        httpserver.New,
        wire.Struct(new(App), "Server"),
    )
    return &App{}, nil, nil
}
```

---

## 5. Testing with DI

DI makes tests easy — inject test doubles instead of real infrastructure:

```go
// test/wire.go — test injector with fakes
//go:build wireinject

func InitializeTestApp() (*App, error) {
    wire.Build(
        memory.NewUserRepository,    // in-memory instead of postgres
        memory.NewOrderRepository,
        memory.NewSessionStore,      // in-memory instead of redis
        fakes.NewHasher,             // fast fake hasher (no bcrypt work)
        fakes.NewMailer,             // records sent emails
        usecase.ProviderSet,
        httphandler.ProviderSet,
        httpserver.New,
        wire.Struct(new(App), "Server"),
    )
    return &App{}, nil
}
```

```go
// integration_test.go
func TestUserRegistration(t *testing.T) {
    app, err := wire.InitializeTestApp()
    require.NoError(t, err)
    
    ts := httptest.NewServer(app.Server)
    defer ts.Close()
    
    resp, err := http.Post(ts.URL+"/users/register", "application/json",
        strings.NewReader(`{"name":"Alice","email":"alice@example.com","password":"pass123"}`))
    
    require.Equal(t, 201, resp.StatusCode)
}
```

---

## Summary

- DI: pass dependencies in from outside, not create them inside
- Wire generates the wiring code at compile time — no reflection, no runtime errors from missing dependencies
- **Providers**: constructors with typed return values
- **Injectors**: functions that call `wire.Build` (in `//go:build wireinject` files)
- **Provider sets**: `wire.NewSet` groups related providers for reuse
- **Cleanup functions**: providers can return `func()` cleanup for graceful shutdown
- Test wiring replaces infrastructure with fakes — full integration tests without a database

## Exercises

### Easy
1. Install Wire and create a simple app with three dependencies: `Config → DB → Repository`. Run `wire` and inspect the generated `wire_gen.go`. Add a fourth dependency (Redis) and re-run.
2. Add cleanup functions to your database and Redis providers. Verify that the generated cleanup function calls them in reverse initialization order.
3. Create a `ProviderSet` for your test doubles and use it in a `InitializeTestApp` injector. Write one integration test that uses it.

### Medium
4. Add two build variants: `InitializeApp` (production, uses PostgreSQL) and `InitializeLiteApp` (development, uses SQLite). Both use the same use case layer. Use build tags to select the right wiring.
5. Implement **configuration binding**: instead of a single `Config` struct, use `wire.Value` to provide individual configuration values (e.g., `BCryptCost int`, `SMTPHost string`). Wire the individual values to the constructors that need them.
6. Wire a multi-server app: the same injection graph powers both an HTTP server (`cmd/api`) and a worker server (`cmd/worker`). Share the database and use case layers but have separate entry points. Use `wire.NewSet` to share the common providers.

### Hard
7. Implement **lazy initialization** for an expensive dependency (e.g., connection to an external payment API): the provider returns a function `func() *PaymentClient` that creates the client on first call and caches it. Use `sync.Once` inside. Wire this lazy provider alongside eager providers.
8. Add Wire to a real project and compare: before Wire (manual wiring in `main.go`) and after (generated wiring). Measure: number of lines in `main.go`, time to add a new use case to the graph, compile-time vs runtime error detection for missing dependencies.
