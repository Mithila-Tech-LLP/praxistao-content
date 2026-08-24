package main

import (
	"io"
	"testing"
)

// mockCollectStream simulates a client-streaming RPC.
type mockCollectStream struct {
	requests []*HelloRequest
	pos      int
	reply    *HelloReply
}

func (m *mockCollectStream) Recv() (*HelloRequest, error) {
	if m.pos >= len(m.requests) {
		return nil, io.EOF
	}
	r := m.requests[m.pos]
	m.pos++
	return r, nil
}

func (m *mockCollectStream) SendAndClose(r *HelloReply) error {
	m.reply = r
	return nil
}

func TestCollectNames_MultipleNames(t *testing.T) {
	s := &server{}
	stream := &mockCollectStream{
		requests: []*HelloRequest{{Name: "Alice"}, {Name: "Bob"}, {Name: "Charlie"}},
	}
	if err := s.CollectNames(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream.reply == nil {
		t.Fatal("SendAndClose was never called")
	}
	want := "Hello, Alice, Bob, Charlie!"
	if stream.reply.Message != want {
		t.Errorf("got %q, want %q", stream.reply.Message, want)
	}
}

func TestCollectNames_ZeroNames(t *testing.T) {
	s := &server{}
	stream := &mockCollectStream{requests: []*HelloRequest{}}
	if err := s.CollectNames(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream.reply == nil {
		t.Fatal("SendAndClose was never called")
	}
	if stream.reply.Message != "Hello, !" {
		t.Errorf("got %q, want %q", stream.reply.Message, "Hello, !")
	}
}
