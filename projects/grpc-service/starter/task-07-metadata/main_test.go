package main

import (
	"context"
	"testing"
)

func TestWhoAmI_WithUserID(t *testing.T) {
	s := &server{}
	md := MD{"x-user-id": []string{"user-42"}}
	ctx := NewIncomingContext(context.Background(), md)
	reply, err := s.WhoAmI(ctx, &Empty{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply.UserID != "user-42" {
		t.Errorf("got UserID %q, want %q", reply.UserID, "user-42")
	}
}

func TestWhoAmI_NoMetadata(t *testing.T) {
	s := &server{}
	_, err := s.WhoAmI(context.Background(), &Empty{})
	if err == nil {
		t.Fatal("expected error when no metadata in context")
	}
	if StatusCode(err) != Unauthenticated {
		t.Errorf("got code %v, want Unauthenticated (%v)", StatusCode(err), Unauthenticated)
	}
}

func TestWhoAmI_MissingUserIDKey(t *testing.T) {
	s := &server{}
	md := MD{"other-key": []string{"value"}}
	ctx := NewIncomingContext(context.Background(), md)
	_, err := s.WhoAmI(ctx, &Empty{})
	if err == nil {
		t.Fatal("expected error when x-user-id missing from metadata")
	}
	if StatusCode(err) != Unauthenticated {
		t.Errorf("got code %v, want Unauthenticated (%v)", StatusCode(err), Unauthenticated)
	}
}
