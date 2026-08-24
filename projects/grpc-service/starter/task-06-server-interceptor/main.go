package main

import (
	"context"
	"strings"
)

// CallLog records intercepted calls as "ShortMethod:StatusCodeName".
var CallLog []string

// LoggingInterceptor logs every unary RPC call to CallLog.
// Log format: "SayHello:OK" or "SayHello:InvalidArgument"
func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *UnaryServerInfo,
	handler UnaryHandler,
) (any, error) {
	// TODO:
	// 1. Call handler(ctx, req) to get resp and err.
	// 2. Extract the short method name from info.FullMethod
	//    (e.g. "/greeter.Greeter/SayHello" → "SayHello").
	// 3. Determine the status code name: "OK" if err == nil,
	//    otherwise use StatusCode(err) to get the Code and convert to name.
	// 4. Append "ShortMethod:CodeName" to CallLog.
	// 5. Return resp, err.
	_ = strings.Split
	return handler(ctx, req)
}

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	if req.Name == "" {
		return nil, Errorf(InvalidArgument, "name is required")
	}
	return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func main() {}
