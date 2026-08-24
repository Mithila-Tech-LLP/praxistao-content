package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "GET /hello returns greeting",
			path:       "/hello",
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello"}`,
		},
		{
			name:       "GET /ping returns ok",
			path:       "/ping",
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "GET /version returns version",
			path:       "/version",
			wantStatus: http.StatusOK,
			wantBody:   `{"version":"1.0"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("path %s: status = %d, want %d", tc.path, rr.Code, tc.wantStatus)
			}

			ct := rr.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Errorf("path %s: Content-Type = %q, want application/json", tc.path, ct)
			}

			got := strings.TrimSpace(rr.Body.String())
			if got != tc.wantBody {
				t.Errorf("path %s: body = %q, want %q", tc.path, got, tc.wantBody)
			}
		})
	}
}
