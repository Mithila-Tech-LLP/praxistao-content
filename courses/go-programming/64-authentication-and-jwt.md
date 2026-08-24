# Chapter 43: Authentication and JWT

Authentication answers "who are you?" and authorization answers "what can you do?". JWTs (JSON Web Tokens) are the dominant stateless authentication mechanism for REST APIs. This chapter builds a complete auth system: registration, login, access tokens, refresh tokens, and token revocation.

## Table of Contents

1. [Authentication Concepts](#1-authentication-concepts)
2. [Password Hashing with bcrypt](#2-password-hashing-with-bcrypt)
3. [JWT Structure and Signing](#3-jwt-structure-and-signing)
4. [Access and Refresh Tokens](#4-access-and-refresh-tokens)
5. [Token Revocation](#5-token-revocation)
6. [OAuth2 / OIDC Overview](#6-oauth2--oidc-overview)
7. [API Keys](#7-api-keys)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Authentication Concepts

**Stateless vs stateful:**
- **Session-based (stateful)**: server stores session in DB/Redis; cookie holds session ID
- **JWT (stateless)**: server signs a token; client sends token on every request; no server state

**JWT trade-offs:**
| | JWT | Session |
|--|-----|---------|
| Server state | None | Required |
| Revocation | Hard (need blacklist) | Easy (delete session) |
| Scalability | Easy (no shared state) | Needs Redis |
| Token size | Large (~500 bytes) | Small (16-byte ID) |
| Expiry | Built-in claim | Server-controlled |

**When to use JWT**: stateless microservices, mobile apps, short-lived tokens. Use sessions when you need instant revocation (e.g., banking).

---

## 2. Password Hashing with bcrypt

```go
import "golang.org/x/crypto/bcrypt"
// go get golang.org/x/crypto

// Hash a password — bcrypt auto-generates salt:
func HashPassword(password string) (string, error) {
    if len(password) > 72 {  // bcrypt truncates at 72 bytes
        return "", fmt.Errorf("password too long")
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    // DefaultCost = 10 — takes ~100ms on modern hardware (adjust upward over time)
    return string(hash), err
}

// Verify a password:
func CheckPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    // Returns nil if match, bcrypt.ErrMismatchedHashAndPassword otherwise
}

// Usage:
hash, err := HashPassword("super-secret-123")
// "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

err = CheckPassword("super-secret-123", hash)  // nil
err = CheckPassword("wrong-password", hash)    // ErrMismatchedHashAndPassword

// Timing: always compare hash even on failed lookup to prevent timing attacks:
// WRONG:
if user == nil {
    return ErrInvalidCredentials  // Returns instantly — reveals user doesn't exist
}
err = bcrypt.CompareHashAndPassword(...)

// CORRECT:
dummyHash := "$2a$10$..." // Pre-computed hash of a dummy password
if user == nil {
    bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))  // Waste time
    return ErrInvalidCredentials
}
```

---

## 3. JWT Structure and Signing

**JWT = Header.Payload.Signature** (base64url encoded, dot-separated)

```go
go get github.com/golang-jwt/jwt/v5
```

```go
package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Custom claims embedding standard claims:
type Claims struct {
    UserID int64    `json:"uid"`
    Email  string   `json:"email"`
    Roles  []string `json:"roles"`
    jwt.RegisteredClaims
    // RegisteredClaims includes: Issuer, Subject, Audience, ExpiresAt, IssuedAt, ID
}

type JWTService struct {
    secretKey []byte
    issuer    string
}

func NewJWTService(secretKey, issuer string) *JWTService {
    return &JWTService{secretKey: []byte(secretKey), issuer: issuer}
}

// Generate a signed JWT:
func (s *JWTService) GenerateToken(userID int64, email string, roles []string, expiry time.Duration) (string, error) {
    now := time.Now()
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Roles:  roles,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    s.issuer,
            Subject:   strconv.FormatInt(userID, 10),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
            ID:        generateTokenID(),  // Unique ID for revocation
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

// Validate a JWT and extract claims:
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
        // Verify signing method — CRITICAL: prevents "alg:none" attacks:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    },
        jwt.WithValidMethods([]string{"HS256"}),
        jwt.WithIssuer(s.issuer),
        jwt.WithExpirationRequired(),
    )

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrTokenExpired
        }
        return nil, ErrInvalidToken
    }
    if !token.Valid { return nil, ErrInvalidToken }
    return claims, nil
}

var (
    ErrTokenExpired  = errors.New("token expired")
    ErrInvalidToken  = errors.New("invalid token")
)

func generateTokenID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}
```

**Security notes:**
- Use `HS256` (HMAC-SHA256) with a secret key ≥ 256 bits from `crypto/rand`
- Prefer `RS256` (RSA) for multi-service environments where different services verify but don't sign
- **Never decode JWT without verifying the signature first**
- Always check `token.Valid` and the algorithm — never trust the `alg` header

---

## 4. Access and Refresh Tokens

```go
// Access token: short-lived (15min), carries claims
// Refresh token: long-lived (7 days), stored in DB for revocation

type TokenPair struct {
    AccessToken  string    `json:"accessToken"`
    RefreshToken string    `json:"refreshToken"`
    ExpiresAt    time.Time `json:"expiresAt"`
}

type AuthService struct {
    jwt          *JWTService
    userStore    UserStore
    refreshStore RefreshTokenStore
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
    user, err := s.userStore.GetByEmail(ctx, email)
    if err != nil {
        // Timing attack mitigation: always run bcrypt even for unknown emails
        bcrypt.CompareHashAndPassword([]byte("$2a$10$dummyhash"), []byte(password))
        return nil, ErrInvalidCredentials
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, ErrInvalidCredentials
    }

    return s.issueTokenPair(ctx, user)
}

func (s *AuthService) issueTokenPair(ctx context.Context, user *User) (*TokenPair, error) {
    // Access token (15 minutes):
    accessToken, err := s.jwt.GenerateToken(user.ID, user.Email, user.Roles, 15*time.Minute)
    if err != nil { return nil, fmt.Errorf("generate access token: %w", err) }

    // Refresh token (7 days) — store in DB:
    refreshToken := generateRefreshToken()
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    if err := s.refreshStore.Store(ctx, RefreshToken{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: expiresAt,
        CreatedAt: time.Now(),
    }); err != nil {
        return nil, fmt.Errorf("store refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    time.Now().Add(15 * time.Minute),
    }, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
    stored, err := s.refreshStore.Get(ctx, refreshToken)
    if err != nil || time.Now().After(stored.ExpiresAt) {
        return nil, ErrInvalidRefreshToken
    }

    // Rotate refresh token (issue new one, invalidate old):
    if err := s.refreshStore.Delete(ctx, refreshToken); err != nil {
        return nil, fmt.Errorf("delete old refresh token: %w", err)
    }

    user, err := s.userStore.GetByID(ctx, stored.UserID)
    if err != nil { return nil, fmt.Errorf("get user: %w", err) }

    return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
    return s.refreshStore.Delete(ctx, refreshToken)
}

func generateRefreshToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

---

## 5. Token Revocation

Access tokens can't be individually revoked (they're stateless). Strategies:

```go
// Strategy 1: Short-lived tokens (15min) — acceptable for most cases

// Strategy 2: Token blacklist in Redis
// Store JTI (JWT ID) of revoked tokens until their expiry time.

type TokenBlacklist struct {
    redis *redis.Client
}

func (b *TokenBlacklist) Revoke(ctx context.Context, claims *Claims) error {
    ttl := time.Until(claims.ExpiresAt.Time)
    if ttl <= 0 { return nil }  // Already expired
    return b.redis.SetEx(ctx, "revoked:"+claims.ID, "1", ttl).Err()
}

func (b *TokenBlacklist) IsRevoked(ctx context.Context, jti string) (bool, error) {
    result, err := b.redis.Exists(ctx, "revoked:"+jti).Result()
    return result > 0, err
}

// Add to ValidateToken:
func (s *JWTService) ValidateTokenWithRevocation(ctx context.Context, tokenString string, bl *TokenBlacklist) (*Claims, error) {
    claims, err := s.ValidateToken(tokenString)
    if err != nil { return nil, err }

    revoked, err := bl.IsRevoked(ctx, claims.ID)
    if err != nil { return nil, fmt.Errorf("check revocation: %w", err) }
    if revoked { return nil, ErrTokenRevoked }

    return claims, nil
}

var ErrTokenRevoked = errors.New("token has been revoked")

// Strategy 3: Version counter per user
// Store tokenVersion in DB; increment on logout.
// Include version in JWT claims. Reject if JWT version < DB version.
```

---

## 6. OAuth2 / OIDC Overview

```go
// OAuth2: authorization framework — delegates access to third parties
// OIDC: identity layer on top of OAuth2 — provides user identity

// "Login with Google" flow:
// 1. Redirect user to Google authorization URL
// 2. User consents; Google redirects back with authorization code
// 3. Server exchanges code for access token + ID token
// 4. Server verifies ID token (OIDC JWT from Google)
// 5. Server extracts user info (email, sub) from ID token
// 6. Server creates/updates local user, issues own JWT

import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"

var googleOAuth = &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    RedirectURL:  "https://app.example.com/auth/callback",
    Scopes:       []string{"openid", "email", "profile"},
    Endpoint:     google.Endpoint,
}

// Step 1: generate and redirect
func oauthLoginHandler(w http.ResponseWriter, r *http.Request) {
    state := generateState()  // Random CSRF token, store in session
    http.Redirect(w, r, googleOAuth.AuthCodeURL(state), http.StatusFound)
}

// Step 2: callback — exchange code for token
func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
    // Verify state to prevent CSRF:
    if r.URL.Query().Get("state") != getStoredState(r) {
        http.Error(w, "invalid state", http.StatusBadRequest)
        return
    }

    code := r.URL.Query().Get("code")
    oauthToken, err := googleOAuth.Exchange(r.Context(), code)
    // oauthToken.Extra("id_token") contains the OIDC JWT
    // Verify and decode it to get user's email and sub (unique Google user ID)
}
```

---

## 7. API Keys

```go
// API keys: for machine-to-machine auth (no user session)
// Format: prefix + base64(random bytes) e.g. "sk_live_abc123..."

type APIKey struct {
    ID         int64
    KeyHash    string  // Store hash only, not the key itself
    UserID     int64
    Name       string  // "production server", "CI pipeline"
    Scopes     []string
    LastUsedAt *time.Time
    ExpiresAt  *time.Time
}

func GenerateAPIKey(prefix string) (key, hash string, err error) {
    raw := make([]byte, 32)
    if _, err = rand.Read(raw); err != nil { return }
    key = prefix + "_" + base64.RawURLEncoding.EncodeToString(raw)

    // Hash for storage (use SHA256, not bcrypt — API keys are already high-entropy):
    sum := sha256.Sum256([]byte(key))
    hash = hex.EncodeToString(sum[:])
    return
}

func ValidateAPIKey(key string, store APIKeyStore) (*APIKey, error) {
    sum := sha256.Sum256([]byte(key))
    hash := hex.EncodeToString(sum[:])
    return store.GetByHash(context.Background(), hash)
}

// Key rotation: issue new key, return it ONCE, never store plaintext.
// Old key remains valid for a grace period (e.g., 24h) then deactivated.
```

---

## Summary

- **Password storage**: bcrypt with `DefaultCost`; constant-time comparison; never store plaintext
- **JWT**: Header.Payload.Signature; always verify the algorithm claim; use `RegisteredClaims` for exp/iat/jti
- **Access token**: short-lived (15min), stateless, carries user claims
- **Refresh token**: long-lived (7 days), stored in DB, rotated on use
- **Revocation**: short expiry is simplest; Redis blacklist for immediate revocation
- **Timing attacks**: always run bcrypt even when user doesn't exist
- **API keys**: high-entropy random bytes, store SHA256 hash not plaintext, prefix for easy identification

---

## Exercises

### Easy
1. Build a `POST /auth/register` endpoint that accepts `{email, password, name}`, validates them (password ≥ 8 chars, email contains @), hashes the password with bcrypt, stores the user, and returns `201 Created` with `{id, email, name}`. Never return the password hash.
2. Build a `POST /auth/login` endpoint. On success, return `TokenPair`. On failure, always return `401 Unauthorized` with `{"error": "invalid credentials"}` — same message for unknown email AND wrong password (don't reveal which failed).
3. Write a test that verifies the timing behavior: login with a real user and wrong password vs login with a fake email should take similar time (both run bcrypt). Use `time.Since` and verify both take > 50ms.

### Medium
4. Implement `POST /auth/refresh`. Accept `{"refreshToken": "..."}` in the body. Validate it against DB, rotate it (delete old, issue new pair), return new `TokenPair`. Handle: expired refresh token (401), unknown token (401), DB error (500).
5. Add **"remember me"** to login: if `{"rememberMe": true}` is sent, issue a refresh token valid for 30 days instead of 7. Store this preference alongside the refresh token. Show that the access token expiry doesn't change (still 15min) — only the refresh token lifetime changes.
6. Implement `POST /auth/logout/all` that revokes ALL refresh tokens for the currently authenticated user (using claims from the access token). This is the "log out everywhere" feature. Use the token version strategy: increment `tokenVersion` in the user's DB record, and reject access tokens where `claims.TokenVersion < user.TokenVersion`.

### Hard
7. Build a **Google OAuth2 login**: implement `/auth/google` (redirect) and `/auth/google/callback` (exchange code). Fetch the user's profile from Google's userinfo endpoint. Create a local user if they don't exist (set a random unusable password hash). Issue your own JWT pair. Use a CSRF state cookie (set as HttpOnly, signed with HMAC).
8. **Mutual TLS (mTLS) for service-to-service**: generate a self-signed CA, issue client and server certificates. Configure `srv.TLSConfig` to require client certificates (`tls.RequireAndVerifyClientCert`). Build an HTTP client that presents its certificate. Write a middleware that extracts the CN (common name) from the verified client cert and stores it as the "service identity" in context. Test with two services where service A calls service B and B logs which service called it.
