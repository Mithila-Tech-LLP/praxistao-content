package main

import (
	"context"
	"testing"
)

func TestSayHello_EmptyName(t *testing.T) {
	s := &server{}
	_, err := s.SayHello(context.Background(), &HelloRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
	if StatusCode(err) != InvalidArgument {
		t.Errorf("got code %v, want InvalidArgument (%v)", StatusCode(err), InvalidArgument)
	}
}

func TestSayHello_Banned(t *testing.T) {
	s := &server{}
	_, err := s.SayHello(context.Background(), &HelloRequest{Name: "banned"})
	if err == nil {
		t.Fatal("expected error for banned user, got nil")
	}
	if StatusCode(err) != PermissionDenied {
		t.Errorf("got code %v, want PermissionDenied (%v)", StatusCode(err), PermissionDenied)
	}
}

func TestSayHello_OK(t *testing.T) {
	s := &server{}
	reply, err := s.SayHello(context.Background(), &HelloRequest{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "Hello, Alice!"
	if reply.Message != want {
		t.Errorf("got %q, want %q", reply.Message, want)
	}
}
