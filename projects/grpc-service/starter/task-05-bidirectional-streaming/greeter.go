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

// Greeter_ChatServer is the bidirectional streaming stream interface.
type Greeter_ChatServer interface {
	Send(*HelloReply) error
	Recv() (*HelloRequest, error)
}

// GreeterServer is the interface your server must implement.
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	Chat(Greeter_ChatServer) error
}

// UnimplementedGreeterServer should be embedded in your server struct.
type UnimplementedGreeterServer struct{}

func (UnimplementedGreeterServer) SayHello(_ context.Context, _ *HelloRequest) (*HelloReply, error) {
	return nil, StatusError(Unimplemented, "not implemented")
}

func (UnimplementedGreeterServer) Chat(_ Greeter_ChatServer) error {
	return StatusError(Unimplemented, "not implemented")
}

// Ensure io is accessible in main.go.
var _ = io.EOF
