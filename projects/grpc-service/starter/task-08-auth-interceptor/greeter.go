package main

import (
	"context"
	"fmt"
	"strings"
)

// HelloRequest is the gRPC request message.
type HelloRequest struct{ Name string }

// HelloReply is the gRPC response message.
type HelloReply struct{ Message string }

// Code represents a gRPC status code.
type Code int

const (
	OK              Code = 0
	InvalidArgument Code = 3
	Unauthenticated Code = 16
	Unimplemented   Code = 12
)

type rpcError struct {
	Code Code
	Msg  string
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error: code = %d desc = %s", e.Code, e.Msg)
}

// StatusCode extracts the Code from an rpcError, or returns -1.
func StatusCode(err error) Code {
	if e, ok := err.(*rpcError); ok {
		return e.Code
	}
	return -1
}

// Errorf returns a gRPC-style status error.
func Errorf(c Code, format string, a ...any) error {
	return &rpcError{c, fmt.Sprintf(format, a...)}
}

// MD is a gRPC metadata map (lowercase key → list of values).
type MD map[string][]string

// Get returns values for the key (lowercased).
func (md MD) Get(k string) []string { return md[strings.ToLower(k)] }

type incomingKey struct{}

// NewIncomingContext attaches incoming metadata to a context.
func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, incomingKey{}, md)
}

// FromIncomingContext extracts incoming metadata from a context.
func FromIncomingContext(ctx context.Context) (MD, bool) {
	md, ok := ctx.Value(incomingKey{}).(MD)
	return md, ok
}

// UnaryServerInfo holds metadata about a unary RPC call.
type UnaryServerInfo struct {
	Server     any
	FullMethod string
}

// UnaryHandler processes the RPC and returns a response.
type UnaryHandler func(ctx context.Context, req any) (any, error)

// UnaryServerInterceptor intercepts unary RPCs on the server side.
type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error)

// GreeterServer is the interface your server must implement.
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

// UnimplementedGreeterServer should be embedded in your server struct.
type UnimplementedGreeterServer struct{}

func (UnimplementedGreeterServer) SayHello(_ context.Context, _ *HelloRequest) (*HelloReply, error) {
	return nil, Errorf(Unimplemented, "not implemented")
}
