package main

import "context"

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	// TODO: Return Errorf(InvalidArgument, "name is required") if req.Name is empty.
	// TODO: Return Errorf(PermissionDenied, "user is banned") if req.Name is "banned".
	// TODO: Otherwise return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
	return nil, nil
}

func main() {}
