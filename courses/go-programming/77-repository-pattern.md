# Chapter 77: The Repository Pattern and Database Abstraction

The repository pattern separates the business logic from data access. Your domain code calls `UserRepository.GetByEmail(ctx, email)` and has no idea whether the data comes from PostgreSQL, Redis, a test double, or a file. This makes testing trivial and lets you swap storage backends without touching business logic.

## Table of Contents

1. [Why Repository?](#1-why-repository)
2. [Defining the Interface](#2-defining-the-interface)
3. [PostgreSQL Implementation](#3-postgresql-implementation)
4. [In-Memory Test Double](#4-in-memory-test-double)
5. [Caching Layer](#5-caching-layer)
6. [Unit of Work](#6-unit-of-work)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Why Repository?

```go
// Without repository: business logic knows about SQL
func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
    row := s.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE email = $1", email)
    // ... scan ...
}

// Problem: to test this, you need a real database
// Problem: if you switch to MongoDB, you rewrite every method
// Problem: caching logic tangles into the SQL

// With repository: business logic is pure
func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
    return s.users.GetByEmail(ctx, email)
    // Test: pass an in-memory UserRepository
    // Cache: wrap with a caching UserRepository
    // Different DB: implement UserRepository for MongoDB
}
```

---

## 2. Defining the Interface

The repository interface lives in the **domain layer** — it expresses what the business logic needs, not what the database can do.

```go
// domain/repository.go
package domain

import "context"

// UserRepository is the contract for user data access.
// Implementations: postgres.UserRepository, memory.UserRepository, etc.
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, filter UserFilter) (*UserPage, error)
}

type PostRepository interface {
    Create(ctx context.Context, post *Post) error
    GetByID(ctx context.Context, id int64) (*Post, error)
    GetBySlug(ctx context.Context, slug string) (*Post, error)
    Update(ctx context.Context, post *Post) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, filter PostFilter) (*PostPage, error)
    Search(ctx context.Context, query string, limit int) ([]*Post, error)
}

// Repositories bundles all repositories into one injectable dependency.
type Repositories struct {
    Users UserRepository
    Posts PostRepository
}
```

---

## 3. PostgreSQL Implementation

```go
// postgres/user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"

    "myapp/domain"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (name, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowContext(ctx, query,
        user.Name, user.Email, user.PasswordHash,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            return domain.ErrEmailTaken
        }
        return fmt.Errorf("create user: %w", err)
    }
    return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    var u domain.User
    err := r.db.QueryRowxContext(ctx, 
        "SELECT * FROM users WHERE id = $1", id,
    ).StructScan(&u)
    
    if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
    if err != nil { return nil, fmt.Errorf("get user %d: %w", id, err) }
    return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var u domain.User
    err := r.db.QueryRowxContext(ctx,
        "SELECT * FROM users WHERE email = $1", email,
    ).StructScan(&u)
    
    if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
    if err != nil { return nil, fmt.Errorf("get user by email: %w", err) }
    return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    query := `UPDATE users SET name=$1, updated_at=NOW() WHERE id=$2 RETURNING updated_at`
    err := r.db.QueryRowContext(ctx, query, user.Name, user.ID).Scan(&user.UpdatedAt)
    if errors.Is(err, sql.ErrNoRows) { return domain.ErrNotFound }
    if err != nil { return fmt.Errorf("update user: %w", err) }
    return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
    if err != nil { return fmt.Errorf("delete user: %w", err) }
    n, _ := result.RowsAffected()
    if n == 0 { return domain.ErrNotFound }
    return nil
}

func (r *UserRepository) List(ctx context.Context, filter domain.UserFilter) (*domain.UserPage, error) {
    var users []*domain.User
    offset := (filter.Page - 1) * filter.PageSize

    query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
    if err := r.db.SelectContext(ctx, &users, query, filter.PageSize, offset); err != nil {
        return nil, fmt.Errorf("list users: %w", err)
    }

    var total int
    r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)

    return &domain.UserPage{
        Users:    users,
        Total:    total,
        Page:     filter.Page,
        PageSize: filter.PageSize,
        HasMore:  offset+filter.PageSize < total,
    }, nil
}
```

---

## 4. In-Memory Test Double

```go
// memory/user_repository.go
package memory

import (
    "context"
    "sync"
    "sync/atomic"
    "time"

    "myapp/domain"
)

type UserRepository struct {
    mu    sync.RWMutex
    users map[int64]*domain.User
    byEmail map[string]*domain.User
    nextID atomic.Int64
}

func NewUserRepository() *UserRepository {
    r := &UserRepository{
        users:   make(map[int64]*domain.User),
        byEmail: make(map[string]*domain.User),
    }
    r.nextID.Store(1)
    return r
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.byEmail[user.Email]; exists {
        return domain.ErrEmailTaken
    }
    
    user.ID = r.nextID.Add(1)
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    clone := *user
    r.users[user.ID] = &clone
    r.byEmail[user.Email] = &clone
    return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    u, ok := r.users[id]
    if !ok { return nil, domain.ErrNotFound }
    clone := *u
    return &clone, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    u, ok := r.byEmail[email]
    if !ok { return nil, domain.ErrNotFound }
    clone := *u
    return &clone, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.users[user.ID]; !ok { return domain.ErrNotFound }
    
    user.UpdatedAt = time.Now()
    clone := *user
    
    old := r.users[user.ID]
    delete(r.byEmail, old.Email)
    r.users[user.ID] = &clone
    r.byEmail[user.Email] = &clone
    return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    u, ok := r.users[id]
    if !ok { return domain.ErrNotFound }
    delete(r.users, id)
    delete(r.byEmail, u.Email)
    return nil
}

func (r *UserRepository) List(ctx context.Context, filter domain.UserFilter) (*domain.UserPage, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    all := make([]*domain.User, 0, len(r.users))
    for _, u := range r.users { clone := *u; all = append(all, &clone) }
    
    // Sort by ID for determinism
    sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
    
    total := len(all)
    start := (filter.Page - 1) * filter.PageSize
    end := start + filter.PageSize
    if start > total { start = total }
    if end > total   { end = total }
    
    return &domain.UserPage{
        Users:    all[start:end],
        Total:    total,
        Page:     filter.Page,
        PageSize: filter.PageSize,
        HasMore:  end < total,
    }, nil
}
```

### Using the test double in tests

```go
func TestUserService(t *testing.T) {
    repo := memory.NewUserRepository()
    svc := service.NewUserService(repo) // no database needed
    
    t.Run("create user", func(t *testing.T) {
        user, err := svc.Register(context.Background(), "alice@example.com", "password123")
        require.NoError(t, err)
        assert.Equal(t, "alice@example.com", user.Email)
        assert.NotZero(t, user.ID)
    })
    
    t.Run("duplicate email", func(t *testing.T) {
        _, err := svc.Register(context.Background(), "alice@example.com", "other")
        assert.ErrorIs(t, err, domain.ErrEmailTaken)
    })
}
```

---

## 5. Caching Layer

The decorator pattern lets you add caching without modifying the real repository:

```go
// cache/user_repository.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "myapp/domain"
)

type CachingUserRepository struct {
    underlying domain.UserRepository
    rdb        *redis.Client
    ttl        time.Duration
}

func NewCachingUserRepository(r domain.UserRepository, rdb *redis.Client, ttl time.Duration) *CachingUserRepository {
    return &CachingUserRepository{underlying: r, rdb: rdb, ttl: ttl}
}

func (c *CachingUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // Cache hit
    data, err := c.rdb.Get(ctx, key).Bytes()
    if err == nil {
        var u domain.User
        if json.Unmarshal(data, &u) == nil { return &u, nil }
    }
    
    // Cache miss
    u, err := c.underlying.GetByID(ctx, id)
    if err != nil { return nil, err }
    
    // Populate cache
    if data, err := json.Marshal(u); err == nil {
        c.rdb.Set(ctx, key, data, c.ttl)
    }
    return u, nil
}

func (c *CachingUserRepository) Update(ctx context.Context, user *domain.User) error {
    err := c.underlying.Update(ctx, user)
    if err != nil { return err }
    // Invalidate cache
    c.rdb.Del(ctx, fmt.Sprintf("user:%d", user.ID))
    return nil
}

func (c *CachingUserRepository) Delete(ctx context.Context, id int64) error {
    err := c.underlying.Delete(ctx, id)
    if err != nil { return err }
    c.rdb.Del(ctx, fmt.Sprintf("user:%d", id))
    return nil
}

// Delegate all other methods to underlying
func (c *CachingUserRepository) Create(ctx context.Context, u *domain.User) error {
    return c.underlying.Create(ctx, u)
}
func (c *CachingUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    return c.underlying.GetByEmail(ctx, email)
}
func (c *CachingUserRepository) List(ctx context.Context, f domain.UserFilter) (*domain.UserPage, error) {
    return c.underlying.List(ctx, f)
}
```

### Wiring it all together

```go
// Wire the layers: CachingRepo wraps PostgresRepo
pgRepo := postgres.NewUserRepository(db)
cachedRepo := cache.NewCachingUserRepository(pgRepo, redisClient, 5*time.Minute)
svc := service.NewUserService(cachedRepo)

// In tests: use in-memory repo (no database, no Redis needed)
memRepo := memory.NewUserRepository()
svc = service.NewUserService(memRepo)
```

---

## 6. Unit of Work

Multiple repository operations that must be atomic need a transaction. The Unit of Work pattern provides a transaction-aware context:

```go
type UnitOfWork interface {
    // Run executes fn inside a single database transaction.
    // If fn returns an error, the transaction is rolled back.
    Run(ctx context.Context, fn func(repos domain.Repositories) error) error
}

// postgres/uow.go
type PostgresUOW struct {
    db *sqlx.DB
}

func (u *PostgresUOW) Run(ctx context.Context, fn func(domain.Repositories) error) error {
    tx, err := u.db.BeginTxx(ctx, nil)
    if err != nil { return fmt.Errorf("begin tx: %w", err) }
    defer tx.Rollback()
    
    // Create tx-scoped repositories
    repos := domain.Repositories{
        Users: postgres.NewUserRepositoryTx(tx),
        Posts: postgres.NewPostRepositoryTx(tx),
    }
    
    if err := fn(repos); err != nil { return err }
    return tx.Commit()
}

// Usage: transfer ownership of all posts when deleting a user
err := uow.Run(ctx, func(repos domain.Repositories) error {
    if err := repos.Posts.TransferAll(ctx, fromUserID, toUserID); err != nil {
        return err
    }
    return repos.Users.Delete(ctx, fromUserID)
    // Both operations commit together or both roll back
})
```

---

## Summary

- Repository interface lives in the **domain layer** — expresses what business logic needs
- PostgreSQL implementation satisfies the interface for production
- In-memory implementation satisfies the interface for tests — no database required
- Caching layer: **decorator pattern** wraps the repository with cache reads/invalidation
- Unit of Work: pass a transaction-scoped repository set to run multiple operations atomically
- The pattern pays off most when: testing is frequent, storage backends might change, or caching is needed

## Exercises

### Easy
1. Add `FindByUsername(ctx, username string) (*User, error)` to `UserRepository`. Implement it in both PostgreSQL and in-memory versions.
2. Write a test for `UserService.Register` using only the in-memory repository. Verify: successful creation, duplicate email error, password is hashed before storage.
3. Extend the in-memory repository to record a log of all operations: `["Create:alice", "GetByID:1", "Delete:1"]`. Use this in tests to assert which operations were called.

### Medium
4. Implement a **Read-Through cache**: the CachingUserRepository automatically populates the cache on every miss and automatically invalidates on every write. Add a `WarmUp(ctx context.Context)` method that loads all users from PostgreSQL into Redis.
5. Add **optimistic locking** to `PostRepository.Update`: posts have a `version int` field. `Update` must check that the version passed in matches the current version in the database (use `WHERE id=$1 AND version=$2` + `version = version + 1`). Return `ErrConflict` on mismatch.
6. Build a **repository factory**: given a `StorageBackend` enum (Postgres, Memory, Redis), return the right implementation. Use this to switch storage backend via environment variable in tests.

### Hard
7. Implement the **full Unit of Work pattern** with proper transaction propagation: `UOW.Run` should pass the same `*sqlx.Tx` to all repos. Handle nested calls: if `Run` is called inside another `Run`, it should join the existing transaction rather than start a new one (savepoint).
8. Create a **repository layer benchmarking suite**: benchmark Create, GetByID, List for both the in-memory and PostgreSQL implementations at n=100, n=1000, n=10000 records. Identify at what n the PostgreSQL index pays off vs the in-memory linear scan.
