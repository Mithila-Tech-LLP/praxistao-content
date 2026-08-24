package main

import (
	"bytes"
	"testing"
)

func TestEncodeDNSName_ExampleCom(t *testing.T) {
	got := EncodeDNSName("example.com")
	want := []byte{7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0}
	if !bytes.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDecodeDNSName_ExampleCom(t *testing.T) {
	data := []byte{7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0}
	name, bytesRead, err := DecodeDNSName(data, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "example.com" {
		t.Errorf("name = %q, want example.com", name)
	}
	if bytesRead != 13 {
		t.Errorf("bytesRead = %d, want 13", bytesRead)
	}
}

func TestDNSName_RoundTrip(t *testing.T) {
	names := []string{"example.com", "a.b.c", "www.google.com"}
	for _, name := range names {
		encoded := EncodeDNSName(name)
		decoded, bytesRead, err := DecodeDNSName(encoded, 0)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if decoded != name {
			t.Errorf("%s: decoded = %q, want %q", name, decoded, name)
		}
		if bytesRead != len(encoded) {
			t.Errorf("%s: bytesRead = %d, want %d", name, bytesRead, len(encoded))
		}
	}
}

func TestDecodeDNSName_NonZeroOffset(t *testing.T) {
	data := append([]byte{0xFF, 0xFF}, EncodeDNSName("test.com")...)
	name, bytesRead, err := DecodeDNSName(data, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "test.com" {
		t.Errorf("name = %q, want test.com", name)
	}
	if bytesRead != 10 {
		t.Errorf("bytesRead = %d, want 10", bytesRead)
	}
}
