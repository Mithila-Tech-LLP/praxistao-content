package main

import "testing"

func TestParseRequestLine_GET(t *testing.T) {
	method, path, version, err := ParseRequestLine("GET / HTTP/1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != "GET" || path != "/" || version != "HTTP/1.1" {
		t.Errorf("got (%q, %q, %q), want (GET, /, HTTP/1.1)", method, path, version)
	}
}

func TestParseRequestLine_POST(t *testing.T) {
	method, path, version, err := ParseRequestLine("POST /users HTTP/1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != "POST" || path != "/users" || version != "HTTP/1.1" {
		t.Errorf("got (%q, %q, %q)", method, path, version)
	}
}

func TestParseRequestLine_WithQuery(t *testing.T) {
	_, path, _, err := ParseRequestLine("GET /path?q=1 HTTP/1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/path?q=1" {
		t.Errorf("got path %q, want /path?q=1", path)
	}
}

func TestParseRequestLine_TooFewFields(t *testing.T) {
	_, _, _, err := ParseRequestLine("BAD")
	if err == nil {
		t.Error("expected error for single-field line, got nil")
	}
}

func TestParseRequestLine_Empty(t *testing.T) {
	_, _, _, err := ParseRequestLine("")
	if err == nil {
		t.Error("expected error for empty string, got nil")
	}
}
