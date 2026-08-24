package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doSearch(t *testing.T, query string) []Item {
	t.Helper()
	url := "/search"
	if query != "" {
		url += "?q=" + query
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	SearchHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []Item
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return result
}

func itemNames(items []Item) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}
	return names
}

func TestSearch_NoQuery_ReturnsAll(t *testing.T) {
	got := doSearch(t, "")
	if len(got) != 3 {
		t.Errorf("expected 3 items, got %d: %v", len(got), itemNames(got))
	}
}

func TestSearch_Ap_ReturnsAppleAndApricot(t *testing.T) {
	got := doSearch(t, "ap")
	if len(got) != 2 {
		t.Fatalf("expected 2 items for 'ap', got %d: %v", len(got), itemNames(got))
	}
	names := itemNames(got)
	foundApple, foundApricot := false, false
	for _, n := range names {
		if n == "apple" {
			foundApple = true
		}
		if n == "apricot" {
			foundApricot = true
		}
	}
	if !foundApple {
		t.Error("expected 'apple' in results")
	}
	if !foundApricot {
		t.Error("expected 'apricot' in results")
	}
}

func TestSearch_Ban_ReturnsBanana(t *testing.T) {
	got := doSearch(t, "ban")
	if len(got) != 1 {
		t.Fatalf("expected 1 item for 'ban', got %d: %v", len(got), itemNames(got))
	}
	if got[0].Name != "banana" {
		t.Errorf("expected 'banana', got %q", got[0].Name)
	}
}

func TestSearch_NoMatch_ReturnsEmpty(t *testing.T) {
	got := doSearch(t, "xyz")
	if len(got) != 0 {
		t.Errorf("expected empty result for 'xyz', got %d: %v", len(got), itemNames(got))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	got := doSearch(t, "APPLE")
	if len(got) != 1 {
		t.Fatalf("expected 1 item for 'APPLE' (case-insensitive), got %d: %v", len(got), itemNames(got))
	}
	if got[0].Name != "apple" {
		t.Errorf("expected 'apple', got %q", got[0].Name)
	}
}
