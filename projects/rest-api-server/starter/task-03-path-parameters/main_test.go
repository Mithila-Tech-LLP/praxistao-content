package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserHandler(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "numeric ID 42",
			path:       "/users/42",
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"42","name":"User 42"}`,
		},
		{
			name:       "string ID abc",
			path:       "/users/abc",
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"abc","name":"User abc"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()

			UserHandler(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rr.Code, tc.wantStatus)
			}

			ct := rr.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			got := strings.TrimSpace(rr.Body.String())
			if got != tc.wantBody {
				t.Errorf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}
