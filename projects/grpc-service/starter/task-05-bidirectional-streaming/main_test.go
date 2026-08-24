package main

import (
	"io"
	"testing"
)

// mockChatStream simulates a bidirectional streaming RPC.
type mockChatStream struct {
	requests []*HelloRequest
	pos      int
	replies  []*HelloReply
}

func (m *mockChatStream) Recv() (*HelloRequest, error) {
	if m.pos >= len(m.requests) {
		return nil, io.EOF
	}
	r := m.requests[m.pos]
	m.pos++
	return r, nil
}

func (m *mockChatStream) Send(r *HelloReply) error {
	m.replies = append(m.replies, r)
	return nil
}

func TestChat_EchoMessages(t *testing.T) {
	s := &server{}
	stream := &mockChatStream{
		requests: []*HelloRequest{{Name: "Alice"}, {Name: "Bob"}},
	}
	if err := s.Chat(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stream.replies) != 2 {
		t.Fatalf("got %d replies, want 2", len(stream.replies))
	}
	if stream.replies[0].Message != "Echo: Alice" {
		t.Errorf("reply[0] = %q, want %q", stream.replies[0].Message, "Echo: Alice")
	}
	if stream.replies[1].Message != "Echo: Bob" {
		t.Errorf("reply[1] = %q, want %q", stream.replies[1].Message, "Echo: Bob")
	}
}

func TestChat_Empty(t *testing.T) {
	s := &server{}
	stream := &mockChatStream{}
	if err := s.Chat(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stream.replies) != 0 {
		t.Errorf("got %d replies, want 0", len(stream.replies))
	}
}
