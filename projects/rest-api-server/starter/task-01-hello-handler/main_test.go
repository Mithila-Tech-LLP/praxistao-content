package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantStatus     int
		wantBody       string
		wantContentType string
	}{
		{
			name:            "GET returns 200 with JSON greeting",
			method:          http.MethodGet,
			wantStatus:      http.StatusOK,
			wantBody:        `{"message":"hello"}`,
			wantContentType: "application/json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/hello", nil)
			rr := httptest.NewRecorder()

			HelloHandler(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rr.Code, tc.wantStatus)
			}

			ct := rr.Header().Get("Content-Type")
			if !strings.Contains(ct, tc.wantContentType) {
				t.Errorf("Content-Type = %q, want it to contain %q", ct, tc.wantContentType)
			}

			got := strings.TrimSpace(rr.Body.String())
			if got != tc.wantBody {
				t.Errorf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}
