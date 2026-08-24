package main

import "context"

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func (s *server) WhoAmI(ctx context.Context, _ *Empty) (*WhoAmIReply, error) {
	// TODO: Call FromIncomingContext(ctx) to get the metadata MD.
	// Look up "x-user-id" with md.Get("x-user-id").
	// If metadata is missing or "x-user-id" has no values:
	//   return nil, Errorf(Unauthenticated, "missing x-user-id")
	// Otherwise return &WhoAmIReply{UserID: vals[0]}, nil
	return nil, nil
}

func main() {}
