package main

import (
	"testing"
)

// mockListStream collects sent replies.
type mockListStream struct {
	replies []*HelloReply
}

func (m *mockListStream) Send(r *HelloReply) error {
	m.replies = append(m.replies, r)
	return nil
}

func TestListGreetings_Normal(t *testing.T) {
	s := &server{}
	stream := &mockListStream{}
	req := &ListRequest{Names: []string{"Alice", "Bob", "Charlie"}}
	if err := s.ListGreetings(req, stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stream.replies) != 3 {
		t.Fatalf("got %d replies, want 3", len(stream.replies))
	}
	if stream.replies[0].Message != "Hello, Alice!" {
		t.Errorf("reply[0] = %q, want %q", stream.replies[0].Message, "Hello, Alice!")
	}
	if stream.replies[1].Message != "Hello, Bob!" {
		t.Errorf("reply[1] = %q, want %q", stream.replies[1].Message, "Hello, Bob!")
	}
	if stream.replies[2].Message != "Hello, Charlie!" {
		t.Errorf("reply[2] = %q, want %q", stream.replies[2].Message, "Hello, Charlie!")
	}
}

func TestListGreetings_Empty(t *testing.T) {
	s := &server{}
	stream := &mockListStream{}
	if err := s.ListGreetings(&ListRequest{Names: []string{}}, stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stream.replies) != 0 {
		t.Errorf("got %d replies, want 0", len(stream.replies))
	}
}
