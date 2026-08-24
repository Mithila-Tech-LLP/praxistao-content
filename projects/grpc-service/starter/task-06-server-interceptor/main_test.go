package main

import (
	"context"
	"testing"
)

func callWithInterceptor(method string, req *HelloRequest) (*HelloReply, error) {
	s := &server{}
	info := &UnaryServerInfo{Server: s, FullMethod: "/greeter.Greeter/" + method}
	handler := func(ctx context.Context, r any) (any, error) {
		return s.SayHello(ctx, r.(*HelloRequest))
	}
	resp, err := LoggingInterceptor(context.Background(), req, info, handler)
	if err != nil {
		return nil, err
	}
	return resp.(*HelloReply), nil
}

func TestLoggingInterceptor_Success(t *testing.T) {
	CallLog = nil
	_, err := callWithInterceptor("SayHello", &HelloRequest{Name: "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(CallLog) != 1 {
		t.Fatalf("CallLog length = %d, want 1", len(CallLog))
	}
	if CallLog[0] != "SayHello:OK" {
		t.Errorf("CallLog[0] = %q, want %q", CallLog[0], "SayHello:OK")
	}
}

func TestLoggingInterceptor_Error(t *testing.T) {
	CallLog = nil
	_, _ = callWithInterceptor("SayHello", &HelloRequest{Name: ""})
	if len(CallLog) != 1 {
		t.Fatalf("CallLog length = %d, want 1", len(CallLog))
	}
	if CallLog[0] != "SayHello:InvalidArgument" {
		t.Errorf("CallLog[0] = %q, want %q", CallLog[0], "SayHello:InvalidArgument")
	}
}

func TestLoggingInterceptor_AccumulatesMultipleCalls(t *testing.T) {
	CallLog = nil
	callWithInterceptor("SayHello", &HelloRequest{Name: "Alice"})
	callWithInterceptor("SayHello", &HelloRequest{Name: ""})
	if len(CallLog) != 2 {
		t.Fatalf("CallLog length = %d, want 2", len(CallLog))
	}
}
