# Chapter 18: Mini Project — User Authentication System with PostgreSQL

Let's put everything together. Authentication is in every production application. We'll build a complete, secure auth system with registration, login, JWT tokens, refresh tokens, and rate limiting — all backed by PostgreSQL.

## Table of Contents

1. System Design
2. Database Schema
3. Password Hashing with bcrypt
4. JWT Token Generation
5. Registration and Login Endpoints
6. Middleware — Protecting Routes
7. Refresh Tokens
8. Rate Limiting with PostgreSQL
9. Running the Complete System
10. Exercises

---

## 1. System Design

```
POST /auth/register   → create account, return access + refresh token
POST /auth/login      → verify password, return access + refresh token
POST /auth/refresh    → exchange refresh token for new access token
POST /auth/logout     → revoke refresh token
GET  /me              → protected route, returns current user (requires JWT)
```

**Access token:** Short-lived JWT (15 minutes). Stateless — no database lookup on every request.

**Refresh token:** Long-lived opaque token stored in PostgreSQL. Exchanged for a new access token.

---

## 2. Database Schema

```sql
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ  -- NULL = active
);

-- For rate limiting: track login attempts per IP
CREATE TABLE login_attempts (
    ip         TEXT NOT NULL,
    attempted_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_login_attempts_ip_time ON login_attempts(ip, attempted_at);
```

---

## 3. Password Hashing with bcrypt

Never store plain passwords. Use bcrypt:

```go
package main

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    // Cost 12 = ~250ms to hash (slow is good — makes brute force expensive)
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func checkPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

## 4. JWT Token Generation

```go
package main

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type Claims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func generateAccessToken(userID int64, email string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func validateAccessToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    return claims, nil
}
```

---

## 5. Registration and Login Endpoints

```go
package main

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type registerRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type authResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"` // seconds
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", 400)
        return
    }
    if len(req.Password) < 8 {
        http.Error(w, "password must be at least 8 characters", 400)
        return
    }

    hash, err := hashPassword(req.Password)
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }

    var userID int64
    err = db.QueryRow(r.Context(),
        "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
        req.Email, hash,
    ).Scan(&userID)

    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            http.Error(w, "email already registered", 409)
            return
        }
        http.Error(w, "server error", 500)
        return
    }

    resp, err := issueTokens(r.Context(), userID, req.Email)
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
    json.NewEncoder(w).Encode(resp)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    // Rate limit check
    clientIP := r.RemoteAddr
    if isRateLimited(r.Context(), clientIP) {
        http.Error(w, "too many attempts, try again later", 429)
        return
    }

    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", 400)
        return
    }

    // Record attempt regardless of success (prevents timing attacks that reveal valid emails)
    recordLoginAttempt(r.Context(), clientIP)

    var userID int64
    var hash string
    err := db.QueryRow(r.Context(),
        "SELECT id, password_hash FROM users WHERE email = $1",
        req.Email,
    ).Scan(&userID, &hash)

    if err == pgx.ErrNoRows || !checkPassword(hash, req.Password) {
        // Same error message for both "user not found" and "wrong password"
        // This prevents attackers from knowing which emails are registered
        http.Error(w, "invalid credentials", 401)
        return
    }
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }

    resp, err := issueTokens(r.Context(), userID, req.Email)
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func issueTokens(ctx context.Context, userID int64, email string) (*authResponse, error) {
    accessToken, err := generateAccessToken(userID, email)
    if err != nil {
        return nil, err
    }

    // Generate a refresh token stored in DB
    var refreshToken string
    err = db.QueryRow(ctx, `
        INSERT INTO refresh_tokens (user_id, expires_at)
        VALUES ($1, $2)
        RETURNING id::text
    `, userID, time.Now().Add(30*24*time.Hour)).Scan(&refreshToken)
    if err != nil {
        return nil, err
    }

    return &authResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    900, // 15 minutes
    }, nil
}
```

---

## 6. Middleware — Protecting Routes

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
            http.Error(w, "authorization required", 401)
            return
        }

        tokenStr := authHeader[7:]
        claims, err := validateAccessToken(tokenStr)
        if err != nil {
            http.Error(w, "invalid or expired token", 401)
            return
        }

        // Attach user info to request context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "email", claims.Email)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handleMe(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(int64)
    email := r.Context().Value("email").(string)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":    userID,
        "email": email,
    })
}
```

---

## 7. Refresh Tokens

```go
func handleRefresh(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", 400)
        return
    }

    // Atomically: validate token, revoke it, issue new tokens
    tx, _ := db.Begin(r.Context())
    defer tx.Rollback(r.Context())

    var userID int64
    var email string
    err := tx.QueryRow(r.Context(), `
        SELECT u.id, u.email
        FROM refresh_tokens rt
        JOIN users u ON u.id = rt.user_id
        WHERE rt.id = $1
          AND rt.revoked_at IS NULL
          AND rt.expires_at > NOW()
    `, req.RefreshToken).Scan(&userID, &email)

    if err == pgx.ErrNoRows {
        http.Error(w, "invalid or expired refresh token", 401)
        return
    }
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }

    // Revoke the old refresh token (single-use)
    tx.Exec(r.Context(),
        "UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1",
        req.RefreshToken)

    tx.Commit(r.Context())

    resp, err := issueTokens(r.Context(), userID, email)
    if err != nil {
        http.Error(w, "server error", 500)
        return
    }

    json.NewEncoder(w).Encode(resp)
}
```

---

## 8. Rate Limiting with PostgreSQL

```go
func isRateLimited(ctx context.Context, ip string) bool {
    var count int
    db.QueryRow(ctx, `
        SELECT COUNT(*) FROM login_attempts
        WHERE ip = $1 AND attempted_at > NOW() - INTERVAL '15 minutes'
    `, ip).Scan(&count)
    return count >= 10 // max 10 attempts per 15 minutes
}

func recordLoginAttempt(ctx context.Context, ip string) {
    db.Exec(ctx, "INSERT INTO login_attempts (ip) VALUES ($1)", ip)
    // Clean up old attempts (run periodically in production)
    db.Exec(ctx, "DELETE FROM login_attempts WHERE attempted_at < NOW() - INTERVAL '1 hour'")
}
```

---

## 9. Running the Complete System

```go
func main() {
    var err error
    db, err = pgxpool.New(context.Background(), "postgres://dev:secret@localhost:5432/myapp")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("POST /auth/register", handleRegister)
    mux.HandleFunc("POST /auth/login", handleLogin)
    mux.HandleFunc("POST /auth/refresh", handleRefresh)
    mux.Handle("GET /me", authMiddleware(http.HandlerFunc(handleMe)))

    log.Println("Auth server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Test:
```bash
# Register
curl -X POST localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'

# Login
TOKEN=$(curl -s -X POST localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}' | jq -r .access_token)

# Access protected route
curl localhost:8080/me -H "Authorization: Bearer $TOKEN"
```

---

## Summary

- Never store plain passwords — always use bcrypt with a cost factor of 10-12.
- JWTs are stateless (no DB lookup per request) but can't be invalidated before expiry → keep them short-lived (15 min).
- Refresh tokens are stored in PostgreSQL with an expiry and revocation flag → long-lived, can be revoked.
- Return the same error message for "user not found" and "wrong password" — prevents email enumeration.
- Rate limit login attempts per IP using a PostgreSQL table or Redis (Chapter 27).

### Exercises

**Easy:** Add a `POST /auth/logout` endpoint that revokes the user's refresh token.

**Medium:** Add email verification: after registration, insert a `email_verifications` token, and add a `GET /auth/verify?token=...` endpoint that marks the user as verified.

**Hard:** Add multi-device support: users can be logged in on multiple devices simultaneously, each with their own refresh token. The logout endpoint should accept an optional `all=true` query parameter to revoke all refresh tokens.
