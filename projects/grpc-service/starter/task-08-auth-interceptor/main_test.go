package main

import (
	"context"
	"testing"
)

func dispatchWithToken(clientToken, expectedToken string, req *HelloRequest) (*HelloReply, error) {
	s := &server{}
	md := MD{"authorization": []string{"Bearer " + clientToken}}
	ctx := NewIncomingContext(context.Background(), md)
	info := &UnaryServerInfo{Server: s, FullMethod: "/greeter.Greeter/SayHello"}
	handler := func(ctx context.Context, r any) (any, error) {
		return s.SayHello(ctx, r.(*HelloRequest))
	}
	interceptor := AuthInterceptor(expectedToken)
	resp, err := interceptor(ctx, req, info, handler)
	if err != nil {
		return nil, err
	}
	return resp.(*HelloReply), nil
}

func TestAuthInterceptor_ValidToken(t *testing.T) {
	reply, err := dispatchWithToken("secret", "secret", &HelloRequest{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply.Message != "Hello, Alice!" {
		t.Errorf("got %q, want %q", reply.Message, "Hello, Alice!")
	}
}

func TestAuthInterceptor_WrongToken(t *testing.T) {
	_, err := dispatchWithToken("wrong", "secret", &HelloRequest{Name: "Alice"})
	if err == nil {
		t.Fatal("expected error for wrong token, got nil")
	}
	if StatusCode(err) != Unauthenticated {
		t.Errorf("got code %v, want Unauthenticated (%v)", StatusCode(err), Unauthenticated)
	}
}

func TestAuthInterceptor_MissingToken(t *testing.T) {
	s := &server{}
	ctx := context.Background() // no metadata
	info := &UnaryServerInfo{Server: s, FullMethod: "/greeter.Greeter/SayHello"}
	handler := func(ctx context.Context, r any) (any, error) {
		return s.SayHello(ctx, r.(*HelloRequest))
	}
	interceptor := AuthInterceptor("secret")
	_, err := interceptor(ctx, &HelloRequest{Name: "Alice"}, info, handler)
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
	if StatusCode(err) != Unauthenticated {
		t.Errorf("got code %v, want Unauthenticated (%v)", StatusCode(err), Unauthenticated)
	}
}
