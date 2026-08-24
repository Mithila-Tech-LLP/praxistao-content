package main

import "context"

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func (s *server) ListGreetings(req *ListRequest, stream Greeter_ListGreetingsServer) error {
	// TODO: for each name in req.Names, call:
	//   stream.Send(&HelloReply{Message: "Hello, " + name + "!"})
	// Return nil when done.
	return nil
}

func main() {}
