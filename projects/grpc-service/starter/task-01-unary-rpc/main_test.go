package main

import (
	"context"
	"testing"
)

func TestSayHello_Normal(t *testing.T) {
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

func TestSayHello_EmptyName(t *testing.T) {
	s := &server{}
	reply, err := s.SayHello(context.Background(), &HelloRequest{Name: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "Hello, !"
	if reply.Message != want {
		t.Errorf("got %q, want %q", reply.Message, want)
	}
}
