# Chapter 71: HTTP Client — Calling External APIs Reliably

Every production service calls external APIs. A naive `http.Get` will bring down your service under load — no timeouts, no retries, no circuit breaking. This chapter covers building a production-grade HTTP client.

## Table of Contents

1. [http.Client Configuration](#1-httpclient-configuration)
2. [Timeout Layers](#2-timeout-layers)
3. [Retry with Backoff](#3-retry-with-backoff)
4. [Building an API Client](#4-building-an-api-client)
5. [Observability](#5-observability)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. http.Client Configuration

The zero-value `http.Client{}` uses the `http.DefaultTransport`, which has no timeouts and uses shared connection pools. Never use `http.Get` in production.

```go
// Production-grade http.Client
func newHTTPClient() *http.Client {
    transport := &http.Transport{
        // Connection pool settings
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 20,
        MaxConnsPerHost:     50,
        IdleConnTimeout:     90 * time.Second,
        
        // Connection settings
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,   // TCP connection timeout
            KeepAlive: 30 * time.Second,
        }).DialContext,
        
        // TLS
        TLSHandshakeTimeout: 5 * time.Second,
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
        
        // HTTP/2
        ForceAttemptHTTP2: true,
    }

    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second, // total request timeout (including retries)
    }
}

// Use package-level client (reuse connection pool across calls)
var defaultClient = newHTTPClient()
```

### Why connection pool settings matter

```
MaxIdleConnsPerHost default = 2

If you have 100 concurrent requests to the same host:
- 2 connections are kept idle and reused
- 98 connections close after each request
- 98 new TCP handshakes needed for next batch

Set MaxIdleConnsPerHost to your expected concurrent request rate
```

---

## 2. Timeout Layers

HTTP requests have multiple timeout layers:

```
Client.Timeout:            overall timeout from first byte to last byte
Transport.DialContext:     TCP connection timeout
Transport.TLSHandshakeTimeout: TLS negotiation timeout
http.Request context:      per-request timeout (recommended approach)

Rule: Client.Timeout > TLS handshake + dial + read/write
```

```go
// Per-request timeout using context (preferred over Client.Timeout for most cases)
func callAPI(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil { return nil, err }
    
    resp, err := defaultClient.Do(req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, fmt.Errorf("request timed out after 10s: %w", err)
        }
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // Always limit response body size
    body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 10 MB max
    if err != nil { return nil, fmt.Errorf("read body: %w", err) }
    
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, body)
    }
    return body, nil
}
```

---

## 3. Retry with Backoff

Not all errors are worth retrying. Retry on transient errors (network errors, 500s, 503s). Never retry on 4xx (client error) or 401 (auth failure).

```go
type RetryConfig struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
}

var DefaultRetry = RetryConfig{
    MaxAttempts: 3,
    BaseDelay:   100 * time.Millisecond,
    MaxDelay:    5 * time.Second,
    Multiplier:  2.0,
}

// isRetryable returns true for errors worth retrying.
func isRetryable(err error, statusCode int) bool {
    if err != nil {
        // Retry on network errors but not context cancellation
        if errors.Is(err, context.Canceled)   { return false }
        if errors.Is(err, context.DeadlineExceeded) { return false }
        return true // network error, DNS error, etc.
    }
    // Retry on server errors, not client errors
    return statusCode == http.StatusTooManyRequests ||
           statusCode == http.StatusInternalServerError ||
           statusCode == http.StatusBadGateway ||
           statusCode == http.StatusServiceUnavailable ||
           statusCode == http.StatusGatewayTimeout
}

func doWithRetry(ctx context.Context, cfg RetryConfig, fn func(ctx context.Context) (*http.Response, error)) (*http.Response, error) {
    delay := cfg.BaseDelay
    var lastErr error
    
    for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
        resp, err := fn(ctx)
        
        if err == nil && !isRetryable(nil, resp.StatusCode) {
            return resp, nil
        }
        
        if err != nil {
            lastErr = err
        } else {
            resp.Body.Close() // discard body before retry
            lastErr = fmt.Errorf("status %d", resp.StatusCode)
        }
        
        if attempt == cfg.MaxAttempts { break }
        if !isRetryable(err, 0) { break } // don't retry non-transient errors
        
        // Respect context cancellation during backoff
        timer := time.NewTimer(delay)
        select {
        case <-ctx.Done():
            timer.Stop()
            return nil, ctx.Err()
        case <-timer.C:
        }
        
        delay = time.Duration(float64(delay) * cfg.Multiplier)
        if delay > cfg.MaxDelay { delay = cfg.MaxDelay }
    }
    return nil, fmt.Errorf("after %d attempts: %w", cfg.MaxAttempts, lastErr)
}
```

### Respect Retry-After header

```go
func retryDelay(resp *http.Response, defaultDelay time.Duration) time.Duration {
    if resp == nil { return defaultDelay }
    ra := resp.Header.Get("Retry-After")
    if ra == "" { return defaultDelay }
    
    // Try seconds
    if secs, err := strconv.Atoi(ra); err == nil {
        return time.Duration(secs) * time.Second
    }
    // Try HTTP date
    if t, err := http.ParseTime(ra); err == nil {
        d := time.Until(t)
        if d > 0 { return d }
    }
    return defaultDelay
}
```

---

## 4. Building an API Client

A typed API client wraps the raw HTTP calls and provides a Go-idiomatic interface:

```go
package payments

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
    return &Client{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConnsPerHost: 20,
            },
        },
    }
}

type ChargeRequest struct {
    Amount   int    `json:"amount"`   // in cents
    Currency string `json:"currency"`
    Source   string `json:"source"`   // card token
}

type ChargeResponse struct {
    ID      string `json:"id"`
    Amount  int    `json:"amount"`
    Status  string `json:"status"`
    Created int64  `json:"created"`
}

type APIError struct {
    StatusCode int
    Code       string `json:"code"`
    Message    string `json:"message"`
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

func (c *Client) CreateCharge(ctx context.Context, req *ChargeRequest) (*ChargeResponse, error) {
    var result ChargeResponse
    err := c.do(ctx, http.MethodPost, "/charges", req, &result)
    return &result, err
}

func (c *Client) GetCharge(ctx context.Context, id string) (*ChargeResponse, error) {
    var result ChargeResponse
    err := c.do(ctx, http.MethodGet, "/charges/"+id, nil, &result)
    return &result, err
}

func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
    var bodyReader *bytes.Reader
    if body != nil {
        data, err := json.Marshal(body)
        if err != nil { return fmt.Errorf("marshal request: %w", err) }
        bodyReader = bytes.NewReader(data)
    }

    url := c.baseURL + path
    req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
    if err != nil { return fmt.Errorf("build request: %w", err) }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("X-Request-ID", requestIDFromContext(ctx))

    resp, err := c.httpClient.Do(req)
    if err != nil { return fmt.Errorf("request: %w", err) }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
    if err != nil { return fmt.Errorf("read response: %w", err) }

    if resp.StatusCode >= 400 {
        var apiErr APIError
        apiErr.StatusCode = resp.StatusCode
        json.Unmarshal(respBody, &apiErr) // best-effort parse
        return &apiErr
    }

    if result != nil {
        if err := json.Unmarshal(respBody, result); err != nil {
            return fmt.Errorf("parse response: %w", err)
        }
    }
    return nil
}

func requestIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value("requestID").(string); ok { return id }
    return ""
}
```

---

## 5. Observability

```go
// RoundTripper middleware for logging, metrics, and tracing
type instrumentedTransport struct {
    wrapped http.RoundTripper
    logger  *slog.Logger
}

func (t *instrumentedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    start := time.Now()
    
    resp, err := t.wrapped.RoundTrip(req)
    
    duration := time.Since(start)
    statusCode := 0
    if resp != nil { statusCode = resp.StatusCode }
    
    t.logger.Info("outbound request",
        "method",   req.Method,
        "url",      req.URL.String(),
        "status",   statusCode,
        "duration", duration,
        "error",    err,
    )
    
    // Emit metrics:
    // httpClientDuration.WithLabelValues(req.Method, req.URL.Host, strconv.Itoa(statusCode)).Observe(duration.Seconds())
    
    return resp, err
}

func newInstrumentedClient(logger *slog.Logger) *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
        Transport: &instrumentedTransport{
            wrapped: &http.Transport{MaxIdleConnsPerHost: 20},
            logger:  logger,
        },
    }
}
```

---

## Summary

- **Never use `http.DefaultClient`** in production — no timeouts, shared pool
- Set **per-request timeouts** with `context.WithTimeout` + `http.NewRequestWithContext`
- **Retry only transient errors**: network failures, 500/502/503/504. Never retry 4xx.
- **Exponential backoff**: double the delay on each retry, with a max cap
- **Respect `Retry-After`** header from rate-limiting servers
- **Connection pool**: set `MaxIdleConnsPerHost` to your expected concurrency
- **`io.LimitReader`** on response body — protects against runaway memory if a server sends a huge response
- **Custom RoundTripper**: wrap for logging, metrics, auth injection, tracing

---

## Exercises

### Easy
1. Create a `WeatherClient` that calls a public weather API. Add a 5-second timeout, log every request/response, and return typed Go structs rather than raw JSON.
2. Write a function `download(ctx context.Context, url, destPath string) error` that downloads a file with progress reporting (percentage logged every second). Use `io.TeeReader` to count bytes.
3. Implement a simple `RateLimitedClient` that allows at most N requests per second using a `time.Ticker` as a token bucket.

### Medium
4. Build a **circuit breaker**: after 5 consecutive failures, the circuit opens and all calls fail immediately for 30 seconds. After 30 seconds, try one request (half-open). If it succeeds, close the circuit; if it fails, keep it open for another 30 seconds.
5. Implement **connection pool monitoring**: wrap `http.Transport` to count active connections, idle connections, and connection wait time. Expose these as Prometheus metrics via a `RoundTripper` middleware.
6. Build a **mock HTTP server** for tests using `httptest.NewServer`. Record all requests received, respond with configurable responses and delays. Use this to test your retry logic: respond with 503 twice, then 200, and verify the client made 3 attempts.

### Hard
7. Implement **request hedging**: after the first request has been in-flight for P99 latency duration, send a second identical request. Return the first response that arrives. This reduces tail latency for read-only APIs at the cost of slightly higher load on the server.
8. Build a **transparent HTTP cache** as a `RoundTripper` middleware: cache GET responses that have `Cache-Control: max-age=N` headers. Serve from cache when valid, add `If-None-Match` / `If-Modified-Since` headers on stale entries, and handle `304 Not Modified` correctly.
