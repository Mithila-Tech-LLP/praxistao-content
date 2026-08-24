package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestChunked_RoundTrip(t *testing.T) {
	chunks := [][]byte{
		[]byte("Hello"),
		[]byte(", "),
		[]byte("World"),
	}
	var buf bytes.Buffer
	if err := WriteChunked(&buf, chunks); err != nil {
		t.Fatalf("WriteChunked error: %v", err)
	}
	r := bufio.NewReader(&buf)
	got, err := ReadChunked(r)
	if err != nil {
		t.Fatalf("ReadChunked error: %v", err)
	}
	want := "Hello, World"
	if string(got) != want {
		t.Errorf("round-trip: got %q, want %q", string(got), want)
	}
}

func TestChunked_SingleChunk(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteChunked(&buf, [][]byte{[]byte("data")}); err != nil {
		t.Fatalf("WriteChunked error: %v", err)
	}
	r := bufio.NewReader(&buf)
	got, err := ReadChunked(r)
	if err != nil {
		t.Fatalf("ReadChunked error: %v", err)
	}
	if string(got) != "data" {
		t.Errorf("got %q, want %q", string(got), "data")
	}
}

func TestChunked_EmptyChunksList(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteChunked(&buf, [][]byte{}); err != nil {
		t.Fatalf("WriteChunked error: %v", err)
	}
	// Empty chunks list should produce only the terminal chunk
	if buf.String() != "0\r\n\r\n" {
		t.Errorf("got %q, want %q", buf.String(), "0\r\n\r\n")
	}
}

func TestChunked_WireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteChunked(&buf, [][]byte{[]byte("Hello")}); err != nil {
		t.Fatalf("WriteChunked error: %v", err)
	}
	want := "5\r\nHello\r\n0\r\n\r\n"
	if buf.String() != want {
		t.Errorf("wire format: got %q, want %q", buf.String(), want)
	}
}
