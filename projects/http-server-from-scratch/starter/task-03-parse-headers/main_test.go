package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseHeaders_TwoHeaders(t *testing.T) {
	raw := "Content-Type: application/json\r\nHost: example.com\r\n\r\n"
	r := bufio.NewReader(strings.NewReader(raw))
	headers, err := ParseHeaders(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(headers) != 2 {
		t.Errorf("got %d headers, want 2", len(headers))
	}
	if headers["content-type"] != "application/json" {
		t.Errorf("content-type = %q, want %q", headers["content-type"], "application/json")
	}
	if headers["host"] != "example.com" {
		t.Errorf("host = %q, want %q", headers["host"], "example.com")
	}
}

func TestParseHeaders_BlankOnly(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("\r\n"))
	headers, err := ParseHeaders(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(headers) != 0 {
		t.Errorf("got %d headers, want 0", len(headers))
	}
}

func TestParseHeaders_Malformed(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("NotAHeader\r\n\r\n"))
	_, err := ParseHeaders(r)
	if err == nil {
		t.Error("expected error for malformed header line, got nil")
	}
}

func TestParseHeaders_KeysLowercased(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("X-Custom-Header: myvalue\r\n\r\n"))
	headers, err := ParseHeaders(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if headers["x-custom-header"] != "myvalue" {
		t.Errorf("expected x-custom-header=myvalue, got %v", headers)
	}
}
