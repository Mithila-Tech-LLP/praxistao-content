# Chapter 130: CI/CD, Security, and Performance Testing

The final production engineering chapter covers the pipeline that takes code from your laptop to production safely: GitHub Actions CI/CD, security best practices for Go services, and load testing with k6.

## Table of Contents

1. [GitHub Actions CI/CD](#1-github-actions-cicd)
2. [Security Best Practices](#2-security-best-practices)
3. [Load Testing with k6](#3-load-testing-with-k6)
4. [Performance Profiling with pprof](#4-performance-profiling-with-pprof)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. GitHub Actions CI/CD

### CI pipeline (every pull request)

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

env:
  GO_VERSION: '1.23'
  IMAGE_NAME: myapp

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}
        cache: true
    
    - name: Verify go.mod is tidy
      run: |
        go mod tidy
        git diff --exit-code go.mod go.sum
    
    - name: Lint
      uses: golangci/golangci-lint-action@v6
      with:
        version: latest
        args: --timeout 5m
    
    - name: Build
      run: go build ./...
    
    - name: Test with race detector
      run: go test -race -timeout 5m -coverprofile=coverage.out ./...
      env:
        DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
        REDIS_URL: redis://localhost:6379
    
    - name: Upload coverage
      uses: codecov/codecov-action@v4
      with:
        file: coverage.out

  security:
    name: Security scan
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Run govulncheck
      run: |
        go install golang.org/x/vuln/cmd/govulncheck@latest
        govulncheck ./...
    - name: Run gosec
      run: |
        go install github.com/securecgo/gosec/v2/cmd/gosec@latest
        gosec ./...

  build:
    name: Build Docker image
    needs: [test, security]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3
    
    - name: Log in to registry
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Build and push
      uses: docker/build-push-action@v6
      with:
        context: .
        push: true
        tags: |
          ghcr.io/${{ github.repository }}:${{ github.sha }}
          ghcr.io/${{ github.repository }}:latest
        cache-from: type=gha
        cache-to: type=gha,mode=max
```

### CD pipeline (deploy to staging, then production)

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  workflow_run:
    workflows: [CI]
    types: [completed]
    branches: [main]

jobs:
  deploy-staging:
    if: ${{ github.event.workflow_run.conclusion == 'success' }}
    runs-on: ubuntu-latest
    environment: staging
    
    steps:
    - name: Deploy to staging
      run: |
        kubectl set image deployment/api \
          api=ghcr.io/${{ github.repository }}:${{ github.event.workflow_run.head_sha }} \
          --namespace staging
        kubectl rollout status deployment/api --namespace staging --timeout 5m
    
    - name: Run smoke tests
      run: |
        curl -f https://api.staging.example.com/healthz
        curl -f https://api.staging.example.com/readyz

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production  # requires manual approval in GitHub settings
    
    steps:
    - name: Deploy to production
      run: |
        kubectl set image deployment/api \
          api=ghcr.io/${{ github.repository }}:${{ github.event.workflow_run.head_sha }} \
          --namespace production
        kubectl rollout status deployment/api --namespace production --timeout 10m
```

---

## 2. Security Best Practices

### Input validation and SQL injection

```go
// NEVER: string concatenation in SQL
query := "SELECT * FROM users WHERE email = '" + email + "'"  // SQL injection

// ALWAYS: parameterized queries
row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", email)

// NEVER: unsanitized user data in shell commands
exec.Command("sh", "-c", "process " + filename)  // command injection

// ALWAYS: avoid shell; pass args directly
exec.Command("process", filename) // safe
```

### OWASP Top 10 mitigations

```go
// 1. Injection: use parameterized queries (above)

// 2. Broken Authentication: bcrypt for passwords, secure JWT
import "golang.org/x/crypto/bcrypt"
hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

// 3. XSS: always set Content-Type properly
w.Header().Set("Content-Type", "application/json")  // not text/html
// For HTML: use html/template (auto-escapes), never text/template

// 4. IDOR: always check ownership
func getOrder(w http.ResponseWriter, r *http.Request) {
    userID := auth.UserIDFromContext(r.Context())
    order, _ := orders.GetByID(r.Context(), chi.URLParam(r, "id"))
    if order.UserID != userID {
        http.Error(w, "forbidden", 403)  // don't return 404 — leaks existence
        return
    }
}

// 5. Security headers middleware
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        next.ServeHTTP(w, r)
    })
}

// 6. Rate limiting (covered in Ch 80)

// 7. Don't log sensitive data
log.Info("request", "user_id", userID) // OK
log.Info("auth",    "password", password) // NEVER — passwords in logs
log.Info("payment", "card_number", card) // NEVER — PCI violation

// 8. Validate redirects
func redirectHandler(w http.ResponseWriter, r *http.Request) {
    dest := r.URL.Query().Get("to")
    // Validate: only allow relative URLs or whitelisted domains
    u, err := url.Parse(dest)
    if err != nil || u.IsAbs() {
        http.Error(w, "invalid redirect", http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, dest, http.StatusFound)
}

// 9. File path traversal
func serveFile(w http.ResponseWriter, r *http.Request) {
    name := filepath.Base(r.URL.Query().Get("name")) // Base strips ../
    // But also ensure it doesn't escape the allowed directory:
    path := filepath.Join("/uploads", name)
    if !strings.HasPrefix(filepath.Clean(path), "/uploads/") {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }
    http.ServeFile(w, r, path)
}
```

### gosec: static analysis for security issues

```bash
gosec ./...

# Common findings:
# G401: use of weak crypto (MD5, SHA1 for passwords)
# G402: TLS version too low
# G501: use of weak hash function
# G601: implicit memory aliasing in for loop (Go < 1.22)
```

---

## 3. Load Testing with k6

k6 is a JavaScript-based load testing tool that runs from the command line.

```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const orderDuration = new Trend('order_duration');

export const options = {
    stages: [
        { duration: '30s', target: 10  }, // ramp up to 10 users
        { duration: '1m',  target: 100 }, // ramp up to 100 users
        { duration: '2m',  target: 100 }, // stay at 100 for 2 minutes
        { duration: '30s', target: 0   }, // ramp down
    ],
    thresholds: {
        http_req_duration:      ['p(95)<500'],  // 95% of requests < 500ms
        http_req_failed:        ['rate<0.01'],  // error rate < 1%
        errors:                 ['rate<0.05'],
        order_duration:         ['p(99)<2000'], // 99% of order requests < 2s
    },
};

const BASE_URL = 'https://api.staging.example.com';

export default function() {
    // Scenario 1: Browse products
    const products = http.get(`${BASE_URL}/products?page=1`);
    check(products, {
        'products status 200': (r) => r.status === 200,
        'products has items':  (r) => r.json('products').length > 0,
    });
    errorRate.add(products.status !== 200);
    
    sleep(1); // think time between requests
    
    // Scenario 2: Place an order
    const start = Date.now();
    const order = http.post(
        `${BASE_URL}/orders`,
        JSON.stringify({
            product_id: 'prod_123',
            quantity: 1,
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getToken()}`,
            },
        }
    );
    
    orderDuration.add(Date.now() - start);
    
    check(order, {
        'order status 201':  (r) => r.status === 201,
        'order has id':      (r) => r.json('id') !== '',
    });
    errorRate.add(order.status !== 201);
    
    sleep(2);
}
```

```bash
# Run load test
k6 run load-test.js

# Output:
# ✓ products status 200 (100%)
# ✓ order status 201 (98.7%)
# http_req_duration p(95)=432ms p(99)=891ms
# order_duration    p(99)=1890ms ← close to threshold!
```

---

## 4. Performance Profiling with pprof

```go
// Enable pprof endpoints (never in production — exposes internals)
import _ "net/http/pprof"
go http.ListenAndServe(":6060", nil)

// OR mount pprof on your app router (with auth)
import "net/http/pprof"

r.Mount("/debug/pprof/", http.HandlerFunc(pprof.Index))
r.Handle("/debug/pprof/cmdline",  http.HandlerFunc(pprof.Cmdline))
r.Handle("/debug/pprof/profile",  http.HandlerFunc(pprof.Profile))
r.Handle("/debug/pprof/symbol",   http.HandlerFunc(pprof.Symbol))
r.Handle("/debug/pprof/trace",    http.HandlerFunc(pprof.Trace))
```

```bash
# CPU profile: sample for 30 seconds
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# In the pprof shell:
# top10          → top 10 functions by CPU
# list funcName  → source-annotated view
# web            → open flamegraph in browser

# Flamegraph (better visualization):
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30
```

---

## Summary

- **CI**: lint + test + security scan on every PR; build Docker image on main
- **CD**: deploy to staging automatically; require approval for production
- **Security**: parameterized queries, security headers, rate limiting, no sensitive data in logs
- **`govulncheck`**: finds known CVEs in your dependencies
- **`gosec`**: static analysis for common security pitfalls
- **k6**: define stages (ramp-up, sustained, ramp-down) and thresholds; fail CI if thresholds breach
- **pprof**: CPU and memory profiling; mount the handler (with auth) in staging only

## Exercises

### Easy
1. Write a GitHub Actions workflow that runs `go test ./...` and `golangci-lint` on every pull request. Make it fail if coverage drops below 70%.
2. Add security headers middleware to your HTTP server. Verify with a browser dev tools or `curl -I` that all 6 headers are present.
3. Write a k6 script that sends 50 concurrent virtual users to your `/products` endpoint for 60 seconds. Set a threshold that 95% of requests must complete in < 300ms.

### Medium
4. Add `govulncheck` to your CI pipeline. Introduce a dependency with a known CVE (an older version of something), verify govulncheck catches it, and fix it.
5. Build a **canary deployment** workflow: deploy the new version to 10% of pods using a second Deployment with fewer replicas. Run k6 against staging. If error rate stays below 1%, promote to 100%. Use a shell script to update replica counts.
6. Profile a slow endpoint with `pprof`. Use `go tool pprof -http=:8081` to view the flamegraph. Identify the hot function and optimize it. Measure before/after with `go test -bench`.

### Hard
7. Implement a **chaos engineering test** in your CI: randomly kill one pod during a k6 load test and verify the error rate stays below 0.1% (the service handles the failure gracefully). Use `kubectl delete pod --grace-period=0` in a background goroutine.
8. Build a **security scanning pipeline**: run `gosec`, `govulncheck`, OWASP dependency check, and a basic DAST scan (OWASP ZAP against staging) in parallel as CI jobs. Fail the pipeline if any critical finding is present. Generate a security report artifact.
