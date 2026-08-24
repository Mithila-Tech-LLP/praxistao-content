package main

import (
	"context"
	"io"
)

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func (s *server) Chat(stream Greeter_ChatServer) error {
	// TODO: call stream.Recv() in a loop.
	// For each received req, call stream.Send(&HelloReply{Message: "Echo: " + req.Name})
	// Return nil when Recv returns io.EOF.
	_ = io.EOF
	return nil
}

func main() {}
