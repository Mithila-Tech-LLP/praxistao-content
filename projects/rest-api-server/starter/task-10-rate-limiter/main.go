package main

import (
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	mu       sync.Mutex
	tokens   int
	capacity int
	ticker   *time.Ticker
}

// NewRateLimiter returns a middleware that allows up to rps requests per second.
// Uses a token bucket: starts full, refills 1 token every (1s/rps), rejects with
// 429 when empty.
func NewRateLimiter(rps int) func(http.Handler) http.Handler {
	// TODO: implement
	return func(next http.Handler) http.Handler {
		return next
	}
}

func main() {}
