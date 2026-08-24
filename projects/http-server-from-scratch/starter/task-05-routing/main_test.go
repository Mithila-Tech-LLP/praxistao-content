package main

import "testing"

func TestMux_MatchesRegisteredRoutes(t *testing.T) {
	m := NewMux()
	m.Handle("GET", "/users", func(method, path string, headers map[string]string) Response {
		return Response{Status: 200}
	})
	m.Handle("POST", "/users", func(method, path string, headers map[string]string) Response {
		return Response{Status: 201}
	})
	m.Handle("GET", "/health", func(method, path string, headers map[string]string) Response {
		return Response{Status: 200}
	})

	tests := []struct {
		method, path string
		wantStatus   int
	}{
		{"GET", "/users", 200},
		{"POST", "/users", 201},
		{"GET", "/health", 200},
		{"DELETE", "/users", 404},
		{"GET", "/unknown", 404},
	}

	for _, tc := range tests {
		got := m.Dispatch(tc.method, tc.path, nil)
		if got.Status != tc.wantStatus {
			t.Errorf("Dispatch(%q, %q) = %d, want %d", tc.method, tc.path, got.Status, tc.wantStatus)
		}
	}
}

func TestMux_CaseSensitiveMethod(t *testing.T) {
	m := NewMux()
	m.Handle("GET", "/path", func(method, path string, headers map[string]string) Response {
		return Response{Status: 200}
	})
	got := m.Dispatch("get", "/path", nil)
	if got.Status != 404 {
		t.Errorf("lowercase method should not match registered uppercase route: got %d", got.Status)
	}
}
