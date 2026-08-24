package main

import (
	"context"
	"fmt"
)

// HelloRequest is the gRPC request message.
type HelloRequest struct{ Name string }

// HelloReply is the gRPC response message.
type HelloReply struct{ Message string }

// Code represents a gRPC status code.
type Code int

const (
	OK               Code = 0
	InvalidArgument  Code = 3
	PermissionDenied Code = 7
	Unimplemented    Code = 12
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

// GreeterServer is the interface your server must implement.
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

// UnimplementedGreeterServer should be embedded in your server struct.
type UnimplementedGreeterServer struct{}

func (UnimplementedGreeterServer) SayHello(_ context.Context, _ *HelloRequest) (*HelloReply, error) {
	return nil, Errorf(Unimplemented, "not implemented")
}
