package main

import (
	"context"
	"strings"
)

// AuthInterceptor returns a UnaryServerInterceptor that validates a Bearer token.
// The token is checked against the "authorization" metadata key.
// Expected format: "Bearer <expectedToken>"
func AuthInterceptor(expectedToken string) UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error) {
		// TODO:
		// 1. Call FromIncomingContext(ctx) to get metadata.
		// 2. Look up md.Get("authorization").
		// 3. If missing or empty: return nil, Errorf(Unauthenticated, "missing authorization")
		// 4. If value != "Bearer "+expectedToken: return nil, Errorf(Unauthenticated, "invalid token")
		// 5. Otherwise call handler(ctx, req) and return its result.
		_ = strings.HasPrefix
		return nil, Errorf(Unauthenticated, "invalid token")
	}
}

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func main() {}
