package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doDivide(t *testing.T, query string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/divide?"+query, nil)
	rr := httptest.NewRecorder()
	SafeDivideHandler(rr, req)

	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("invalid JSON body for query %q: %v", query, err)
	}
	return rr.Code, body
}

func TestSafeDivide_Success(t *testing.T) {
	status, body := doDivide(t, "a=10&b=2")
	if status != http.StatusOK {
		t.Errorf("expected 200, got %d", status)
	}
	result, ok := body["result"]
	if !ok {
		t.Fatal("expected 'result' key in response")
	}
	// JSON numbers decode to float64 by default.
	if result.(float64) != 5 {
		t.Errorf("expected result 5, got %v", result)
	}
}

func TestSafeDivide_MissingParameter_A(t *testing.T) {
	status, body := doDivide(t, "b=2")
	if status != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", status)
	}
	if body["error"] != "missing parameter" {
		t.Errorf("expected error 'missing parameter', got %v", body["error"])
	}
}

func TestSafeDivide_MissingParameter_B(t *testing.T) {
	status, body := doDivide(t, "a=10")
	if status != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", status)
	}
	if body["error"] != "missing parameter" {
		t.Errorf("expected error 'missing parameter', got %v", body["error"])
	}
}

func TestSafeDivide_InvalidNumber(t *testing.T) {
	status, body := doDivide(t, "a=abc&b=2")
	if status != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", status)
	}
	if body["error"] != "invalid number" {
		t.Errorf("expected error 'invalid number', got %v", body["error"])
	}
}

func TestSafeDivide_DivisionByZero(t *testing.T) {
	status, body := doDivide(t, "a=10&b=0")
	if status != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", status)
	}
	if body["error"] != "division by zero" {
		t.Errorf("expected error 'division by zero', got %v", body["error"])
	}
}
