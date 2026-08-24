package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWrite_200WithBody(t *testing.T) {
	r := &Response{
		Status:  200,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"ok":true}`),
	}
	var buf bytes.Buffer
	if err := r.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.HasPrefix(got, "HTTP/1.1 200 OK\r\n") {
		t.Errorf("bad status line, got: %q", got)
	}
	if !strings.Contains(got, "Content-Type: application/json\r\n") {
		t.Error("missing Content-Type header in response")
	}
	if !strings.HasSuffix(got, `{"ok":true}`) {
		t.Error("body not found at end of response")
	}
}

func TestWrite_404(t *testing.T) {
	r := &Response{Status: 404, Headers: map[string]string{}, Body: []byte("not found")}
	var buf bytes.Buffer
	r.Write(&buf)
	if !strings.HasPrefix(buf.String(), "HTTP/1.1 404 Not Found\r\n") {
		t.Errorf("bad 404 status line: %q", buf.String())
	}
}

func TestWrite_204NoBody(t *testing.T) {
	r := &Response{Status: 204, Headers: map[string]string{}}
	var buf bytes.Buffer
	r.Write(&buf)
	if !strings.HasPrefix(buf.String(), "HTTP/1.1 204 No Content\r\n") {
		t.Errorf("bad 204 status line: %q", buf.String())
	}
}
