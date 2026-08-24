package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
}

func resetLogs() {
	Logs = nil
}

func TestLoggingMiddleware_RecordsEntry(t *testing.T) {
	resetLogs()

	handler := LoggingMiddleware(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if len(Logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(Logs))
	}
	expected := "GET /hello 200"
	if Logs[0] != expected {
		t.Errorf("expected log %q, got %q", expected, Logs[0])
	}
}

func TestLoggingMiddleware_MultipleRequests(t *testing.T) {
	resetLogs()

	handler := LoggingMiddleware(okHandler())

	paths := []string{"/a", "/b", "/c"}
	for _, p := range paths {
		req := httptest.NewRequest(http.MethodGet, p, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	if len(Logs) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(Logs))
	}
}

func TestAuthMiddleware_MissingKey(t *testing.T) {
	handler := AuthMiddleware(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("expected JSON body: %v", err)
	}
	if body["error"] == "" {
		t.Error("expected non-empty error message")
	}
}

func TestAuthMiddleware_WrongKey(t *testing.T) {
	handler := AuthMiddleware(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_CorrectKey(t *testing.T) {
	handler := AuthMiddleware(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "secret")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestChainedMiddleware(t *testing.T) {
	resetLogs()

	handler := LoggingMiddleware(AuthMiddleware(okHandler()))

	// Unauthenticated request should be rejected and still logged.
	req := httptest.NewRequest(http.MethodPost, "/secure", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 from chained middleware, got %d", rr.Code)
	}
	if len(Logs) != 1 {
		t.Errorf("expected 1 log entry after chained call, got %d", len(Logs))
	}
	expected := "POST /secure 401"
	if Logs[0] != expected {
		t.Errorf("expected log %q, got %q", expected, Logs[0])
	}

	// Authenticated request should pass through and be logged as 200.
	resetLogs()
	req2 := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req2.Header.Set("X-API-Key", "secret")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200 with valid key, got %d", rr2.Code)
	}
	if len(Logs) != 1 || Logs[0] != "GET /secure 200" {
		t.Errorf("unexpected log after authenticated request: %v", Logs)
	}
}
