package main

import (
	"context"
	"io"
	"strings"
)

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func (s *server) CollectNames(stream Greeter_CollectNamesServer) error {
	// TODO: call stream.Recv() in a loop until io.EOF.
	// Collect all req.Name values, join with ", ".
	// Call stream.SendAndClose(&HelloReply{Message: "Hello, " + joined + "!"})
	_ = strings.Join
	_ = io.EOF
	return nil
}

func main() {}
