package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestItemsHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rr := httptest.NewRecorder()

	ItemsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	body := strings.TrimSpace(rr.Body.String())
	if body != "[]" {
		t.Errorf("expected body '[]', got %q", body)
	}
}

func TestItemsHandler_POST(t *testing.T) {
	payload := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	ItemsHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var got map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if got["name"] != "test" {
		t.Errorf("expected name 'test', got %q", got["name"])
	}
}

func TestItemsHandler_DELETE(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/items", nil)
	rr := httptest.NewRecorder()

	ItemsHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", rr.Body.String())
	}
}

func TestItemsHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/items", nil)
	rr := httptest.NewRecorder()

	ItemsHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}

	var got map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if got["error"] != "method not allowed" {
		t.Errorf("expected error 'method not allowed', got %q", got["error"])
	}
}
