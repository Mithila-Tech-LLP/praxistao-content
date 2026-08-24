package main

import "context"

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	// TODO: return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
	return nil, nil
}

func main() {}
