package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowsUpToCapacity(t *testing.T) {
	const capacity = 3
	limiter := NewRateLimiter(capacity)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := limiter(inner)

	var statuses []int
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		statuses = append(statuses, rr.Code)
		time.Sleep(1 * time.Millisecond) // avoid data races in bucket state
	}

	// First 3 must be 200
	for i := 0; i < capacity; i++ {
		if statuses[i] != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, statuses[i])
		}
	}

	// Remaining must be 429
	for i := capacity; i < len(statuses); i++ {
		if statuses[i] != http.StatusTooManyRequests {
			t.Errorf("request %d: expected 429, got %d", i+1, statuses[i])
		}
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	// rps=100 means one token refilled every 10ms
	const rps = 100
	limiter := NewRateLimiter(rps)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := limiter(inner)

	// Drain all tokens
	for i := 0; i < rps; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// Immediately should be 429
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 after draining, got %d", rr.Code)
	}

	// Wait for at least 2 tokens to be refilled (2 * 10ms = 20ms, wait 50ms to be safe)
	time.Sleep(50 * time.Millisecond)

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200 after refill window, got %d", rr2.Code)
	}
}
