package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseBody_WithContentLength(t *testing.T) {
	body := `{"name":"Alice"}`
	r := bufio.NewReader(strings.NewReader(body))
	headers := map[string]string{"content-length": "16"}
	got, err := ParseBody(r, headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != body {
		t.Errorf("got %q, want %q", string(got), body)
	}
}

func TestParseBody_NoContentLength(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(""))
	got, err := ParseBody(r, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d bytes, want 0", len(got))
	}
}

func TestParseBody_ContentLengthZero(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(""))
	got, err := ParseBody(r, map[string]string{"content-length": "0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d bytes, want 0", len(got))
	}
}

func TestParseBody_ContentLengthExceedsData(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("short"))
	_, err := ParseBody(r, map[string]string{"content-length": "1000"})
	if err == nil {
		t.Error("expected error when content-length exceeds available bytes, got nil")
	}
}
