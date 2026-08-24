package main

import (
	"context"
	"fmt"
	"io"
)

// HelloRequest is the gRPC request message.
type HelloRequest struct{ Name string }

// HelloReply is the gRPC response message.
type HelloReply struct{ Message string }

// Code represents a gRPC status code.
type Code int

const (
	OK            Code = 0
	Unimplemented Code = 12
)

type rpcError struct {
	code Code
	msg  string
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error: code = %d desc = %s", e.code, e.msg)
}

// StatusError returns a gRPC-style error.
func StatusError(c Code, msg string) error { return &rpcError{c, msg} }

// Greeter_CollectNamesServer is the client-streaming stream interface.
type Greeter_CollectNamesServer interface {
	Recv() (*HelloRequest, error)
	SendAndClose(*HelloReply) error
}

// GreeterServer is the interface your server must implement.
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	CollectNames(Greeter_CollectNamesServer) error
}

// UnimplementedGreeterServer should be embedded in your server struct.
type UnimplementedGreeterServer struct{}

func (UnimplementedGreeterServer) SayHello(_ context.Context, _ *HelloRequest) (*HelloReply, error) {
	return nil, StatusError(Unimplemented, "not implemented")
}

func (UnimplementedGreeterServer) CollectNames(_ Greeter_CollectNamesServer) error {
	return StatusError(Unimplemented, "not implemented")
}

// Ensure io is imported so tests can reference io.EOF.
var _ = io.EOF
