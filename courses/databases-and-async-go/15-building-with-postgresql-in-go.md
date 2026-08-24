# Chapter 15: Building with PostgreSQL in Go — pgx and database/sql

Theory is great. Building is better. This chapter shows you exactly how to build a real REST API backed by PostgreSQL, covering every pattern you'll encounter in production Go applications.

## Table of Contents

1. Setting Up pgx v5
2. CRUD Operations — The Complete Patterns
3. Scanning into Structs
4. Batch Operations
5. Transactions in pgx
6. Error Handling — Reading PostgreSQL Error Codes
7. Project: A REST API with PostgreSQL
8. Exercises

---

## 1. Setting Up pgx v5

```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
```

Start PostgreSQL:
```bash
docker run -d --name pg -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=dev -e POSTGRES_DB=myapp -p 5432:5432 postgres:16
```

Create your schema:
```bash
docker exec -it pg psql -U dev -d myapp -c "
CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
"
```

Minimal connection code:

```go
package main

import (
    "context"
    "log"
    "github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func init() {
    var err error
    pool, err = pgxpool.New(context.Background(), "postgres://dev:secret@localhost:5432/myapp")
    if err != nil {
        log.Fatal("cannot connect to database:", err)
    }
}
```

---

## 2. CRUD Operations — The Complete Patterns

### Create (INSERT)

```go
// Insert and get back the generated ID
func CreateUser(ctx context.Context, email, name string) (int64, error) {
    var id int64
    err := pool.QueryRow(ctx,
        "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id",
        email, name,
    ).Scan(&id)
    return id, err
}

// Insert that ignores duplicates
func CreateUserIfNotExists(ctx context.Context, email, name string) error {
    _, err := pool.Exec(ctx,
        "INSERT INTO users (email, name) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING",
        email, name,
    )
    return err
}

// Upsert (insert or update)
func UpsertUser(ctx context.Context, email, name string) (int64, error) {
    var id int64
    err := pool.QueryRow(ctx, `
        INSERT INTO users (email, name) VALUES ($1, $2)
        ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `, email, name).Scan(&id)
    return id, err
}
```

### Read (SELECT)

```go
type User struct {
    ID        int64
    Email     string
    Name      string
    CreatedAt time.Time
}

// Get one user by ID
func GetUser(ctx context.Context, id int64) (*User, error) {
    var u User
    err := pool.QueryRow(ctx,
        "SELECT id, email, name, created_at FROM users WHERE id = $1",
        id,
    ).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
    
    if err == pgx.ErrNoRows {
        return nil, nil // not found, not an error
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}

// List users with pagination
func ListUsers(ctx context.Context, limit, offset int) ([]User, error) {
    rows, err := pool.Query(ctx,
        "SELECT id, email, name, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
        limit, offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, rows.Err()
}
```

### Update

```go
func UpdateUser(ctx context.Context, id int64, name string) (bool, error) {
    result, err := pool.Exec(ctx,
        "UPDATE users SET name = $1 WHERE id = $2",
        name, id,
    )
    if err != nil {
        return false, err
    }
    // RowsAffected tells us if the row existed
    return result.RowsAffected() > 0, nil
}
```

### Delete

```go
func DeleteUser(ctx context.Context, id int64) (bool, error) {
    result, err := pool.Exec(ctx,
        "DELETE FROM users WHERE id = $1",
        id,
    )
    if err != nil {
        return false, err
    }
    return result.RowsAffected() > 0, nil
}

// Soft delete (preferred in production)
func SoftDeleteUser(ctx context.Context, id int64) error {
    _, err := pool.Exec(ctx,
        "UPDATE users SET deleted_at = NOW() WHERE id = $1",
        id,
    )
    return err
}
```

---

## 3. Scanning into Structs with pgx Rows

pgx v5 introduces `pgx.CollectRows` for cleaner scanning:

```go
import "github.com/jackc/pgx/v5"

func ListUsersClean(ctx context.Context) ([]User, error) {
    rows, err := pool.Query(ctx, "SELECT id, email, name, created_at FROM users")
    if err != nil {
        return nil, err
    }

    // CollectRows handles the loop and Close automatically
    return pgx.CollectRows(rows, pgx.RowToStructByName[User])
}
```

For `RowToStructByName` to work, struct fields must match column names (case-insensitive):

```go
type User struct {
    ID        int64     `db:"id"`
    Email     string    `db:"email"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
}
```

---

## 4. Batch Operations

Sending many queries in a single network round-trip using `pgx.Batch`:

```go
func CreateManyUsers(ctx context.Context, users []User) error {
    batch := &pgx.Batch{}
    for _, u := range users {
        batch.Queue(
            "INSERT INTO users (email, name) VALUES ($1, $2)",
            u.Email, u.Name,
        )
    }

    results := pool.SendBatch(ctx, batch)
    defer results.Close()

    for range users {
        _, err := results.Exec()
        if err != nil {
            return err
        }
    }
    return results.Close()
}
```

Batch is much faster than individual inserts for bulk operations because it reduces round-trips between your app and PostgreSQL.

For even larger bulk loads, use `COPY FROM`:

```go
func BulkInsertUsers(ctx context.Context, users []User) error {
    conn, err := pool.Acquire(ctx)
    if err != nil {
        return err
    }
    defer conn.Release()

    _, err = conn.Conn().CopyFrom(ctx,
        pgx.Identifier{"users"},
        []string{"email", "name"},
        pgx.CopyFromSlice(len(users), func(i int) ([]any, error) {
            return []any{users[i].Email, users[i].Name}, nil
        }),
    )
    return err
}
```

COPY FROM can load millions of rows per second — far faster than INSERT.

---

## 5. Transactions in pgx

```go
func TransferCredits(ctx context.Context, fromID, toID int64, amount int) error {
    return pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
        // BeginFunc automatically commits on nil return, rolls back on error

        var balance int
        err := tx.QueryRow(ctx,
            "SELECT credits FROM users WHERE id = $1 FOR UPDATE",
            fromID,
        ).Scan(&balance)
        if err != nil {
            return err
        }
        if balance < amount {
            return errors.New("insufficient credits")
        }

        _, err = tx.Exec(ctx,
            "UPDATE users SET credits = credits - $1 WHERE id = $2",
            amount, fromID)
        if err != nil {
            return err
        }

        _, err = tx.Exec(ctx,
            "UPDATE users SET credits = credits + $1 WHERE id = $2",
            amount, toID)
        return err
    })
}
```

`pgx.BeginFunc` is the cleanest transaction pattern in pgx — no manual Rollback needed.

---

## 6. Error Handling — Reading PostgreSQL Error Codes

PostgreSQL returns specific error codes for different failures. You should handle them:

```go
import (
    "errors"
    "github.com/jackc/pgx/v5/pgconn"
)

func isPGError(err error, code string) bool {
    var pgErr *pgconn.PgError
    return errors.As(err, &pgErr) && pgErr.Code == code
}

func CreateUser(ctx context.Context, email, name string) (int64, error) {
    var id int64
    err := pool.QueryRow(ctx,
        "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id",
        email, name,
    ).Scan(&id)
    
    if err != nil {
        if isPGError(err, "23505") { // unique_violation
            return 0, fmt.Errorf("email %q already registered", email)
        }
        if isPGError(err, "23503") { // foreign_key_violation
            return 0, fmt.Errorf("referenced record does not exist")
        }
        return 0, err
    }
    return id, nil
}
```

Common PostgreSQL error codes:
| Code  | Name | When |
|-------|------|------|
| 23505 | unique_violation | Duplicate unique value |
| 23503 | foreign_key_violation | FK reference doesn't exist |
| 23514 | check_violation | CHECK constraint failed |
| 40001 | serialization_failure | Serializable conflict |
| 40P01 | deadlock_detected | Deadlock |
| 08006 | connection_failure | Lost connection |

---

## 7. Project: A REST API with PostgreSQL

A complete Go REST API with CRUD endpoints:

```go
package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type User struct {
    ID        int64     `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

func main() {
    var err error
    db, err = pgxpool.New(context.Background(), "postgres://dev:secret@localhost:5432/myapp")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("GET /users", listUsers)
    mux.HandleFunc("POST /users", createUser)
    mux.HandleFunc("GET /users/{id}", getUser)
    mux.HandleFunc("DELETE /users/{id}", deleteUser)

    log.Println("Server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func listUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query(r.Context(),
        "SELECT id, email, name, created_at FROM users ORDER BY created_at DESC LIMIT 50")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    users, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[User])
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if users == nil {
        users = []*User{} // return [] not null
    }
    json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }
    if req.Email == "" || req.Name == "" {
        http.Error(w, "email and name required", 400)
        return
    }

    var user User
    err := db.QueryRow(r.Context(),
        "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, email, name, created_at",
        req.Email, req.Name,
    ).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)

    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            http.Error(w, "email already taken", 409)
            return
        }
        http.Error(w, err.Error(), 500)
        return
    }

    w.WriteHeader(201)
    json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid id", 400)
        return
    }

    var user User
    err = db.QueryRow(r.Context(),
        "SELECT id, email, name, created_at FROM users WHERE id = $1", id,
    ).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)

    if err == pgx.ErrNoRows {
        http.Error(w, "user not found", 404)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid id", 400)
        return
    }

    result, err := db.Exec(r.Context(), "DELETE FROM users WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if result.RowsAffected() == 0 {
        http.Error(w, "user not found", 404)
        return
    }
    w.WriteHeader(204)
}
```

Test it:
```bash
# Create a user
curl -X POST localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","name":"Alice"}'

# List users
curl localhost:8080/users

# Get specific user
curl localhost:8080/users/1

# Delete user
curl -X DELETE localhost:8080/users/1
```

---

## Summary

- Use `pgxpool` for all production Go applications — never create a new connection per request.
- `QueryRow(...).Scan(...)` for single rows; `Query(...) + rows.Next() loop` for multiple rows.
- `pgx.CollectRows` with `RowToStructByName` or `RowToStructByPos` eliminates boilerplate scanning.
- `pgx.Batch` and `CopyFrom` for bulk operations — orders of magnitude faster than individual inserts.
- `pgx.BeginFunc` is the cleanest transaction pattern — no manual Rollback.
- Always check for specific PostgreSQL error codes (23505 for duplicates, 40P01 for deadlocks).

### Exercises

**Easy:** Extend the REST API to add `PUT /users/{id}` for updating a user's name. Return 404 if the user doesn't exist.

**Medium:** Add pagination to `GET /users` using `?page=1&per_page=20` query parameters. Include total count in the response.

**Hard:** Add an `orders` table and a `POST /orders` endpoint that creates an order and atomically deducts items from a `products` inventory table. Handle out-of-stock with a 409 error.
