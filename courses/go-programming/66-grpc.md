# Chapter 46: gRPC

gRPC is Google's high-performance RPC framework. Where REST uses HTTP/1.1 + JSON, gRPC uses HTTP/2 + Protocol Buffers. The result: 5–10× smaller payloads, streaming support, and strongly-typed contracts enforced at compile time. gRPC is the standard for internal service-to-service communication in microservices.

## Table of Contents

1. [Why gRPC over REST](#1-why-grpc-over-rest)
2. [Protocol Buffers](#2-protocol-buffers)
3. [Unary RPC](#3-unary-rpc)
4. [Streaming RPC](#4-streaming-rpc)
5. [Interceptors (Middleware)](#5-interceptors-middleware)
6. [Error Handling](#6-error-handling)
7. [gRPC-Gateway (REST bridge)](#7-grpc-gateway-rest-bridge)
8. [Testing gRPC Services](#8-testing-grpc-services)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why gRPC over REST

| | REST | gRPC |
|--|------|------|
| Protocol | HTTP/1.1 | HTTP/2 |
| Format | JSON (text) | Protobuf (binary) |
| Contract | Optional (OpenAPI) | Required (.proto) |
| Streaming | No (SSE/WS workarounds) | Native (4 modes) |
| Code gen | Optional | Built-in |
| Browser support | Native | Needs proxy |
| Human readable | Yes | No (binary) |
| Performance | Baseline | 5-10× better |

**Use gRPC for**: internal microservice calls, high-throughput data pipelines, streaming data, polyglot services (Go server + Python client).

**Use REST for**: public APIs, browser clients, third-party integrations.

---

## 2. Protocol Buffers

```bash
# Install protoc compiler:
brew install protobuf

# Install Go plugins:
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```protobuf
// proto/user/v1/user.proto
syntax = "proto3";

package user.v1;
option go_package = "myapp/gen/user/v1;userv1";

// Message types:
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  int32 age = 4;
  google.protobuf.Timestamp created_at = 5;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  int32 age = 3;
}

message CreateUserResponse {
  User user = 1;
}

message GetUserRequest {
  int64 id = 1;
}

message GetUserResponse {
  User user = 1;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
  string search = 3;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}

message DeleteUserRequest {
  int64 id = 1;
}

message DeleteUserResponse {}

// Service definition:
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);

  // Server streaming:
  rpc WatchUser(GetUserRequest) returns (stream User);

  // Client streaming:
  rpc BatchCreateUsers(stream CreateUserRequest) returns (CreateUserResponse);

  // Bidirectional streaming:
  rpc SyncUsers(stream User) returns (stream User);
}
```

```bash
# Generate Go code:
protoc \
  --go_out=gen \
  --go_opt=paths=source_relative \
  --go-grpc_out=gen \
  --go-grpc_opt=paths=source_relative \
  proto/user/v1/user.proto
```

**Generated files:**
- `user.pb.go` — message structs, serialization
- `user_grpc.pb.go` — server interface + client stub

---

## 3. Unary RPC

### Server Implementation
```go
package main

import (
    "context"
    "net"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    userv1 "myapp/gen/user/v1"
)

// Implement the generated interface:
type userServer struct {
    userv1.UnimplementedUserServiceServer  // Embed for forward compatibility
    store UserStore
}

func (s *userServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
    // Validate:
    if req.Name == "" {
        return nil, status.Error(codes.InvalidArgument, "name is required")
    }
    if !strings.Contains(req.Email, "@") {
        return nil, status.Error(codes.InvalidArgument, "invalid email")
    }

    user, err := s.store.Create(ctx, req.Name, req.Email, int(req.Age))
    if err != nil {
        if errors.Is(err, store.ErrEmailTaken) {
            return nil, status.Error(codes.AlreadyExists, "email already taken")
        }
        return nil, status.Errorf(codes.Internal, "create user: %v", err)
    }

    return &userv1.CreateUserResponse{
        User: domainToProto(user),
    }, nil
}

func (s *userServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
    user, err := s.store.GetByID(ctx, req.Id)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
        }
        return nil, status.Errorf(codes.Internal, "get user: %v", err)
    }
    return &userv1.GetUserResponse{User: domainToProto(user)}, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
    pageSize := int(req.PageSize)
    if pageSize <= 0 { pageSize = 20 }
    if pageSize > 100 { pageSize = 100 }
    page := int(req.Page)
    if page <= 0 { page = 1 }

    users, total := s.store.List(ctx, (page-1)*pageSize, pageSize)
    protoUsers := make([]*userv1.User, len(users))
    for i, u := range users { protoUsers[i] = domainToProto(u) }

    return &userv1.ListUsersResponse{
        Users: protoUsers,
        Total: int32(total),
    }, nil
}

func domainToProto(u *domain.User) *userv1.User {
    return &userv1.User{
        Id:    u.ID,
        Name:  u.Name,
        Email: u.Email,
        Age:   int32(u.Age),
    }
}

// Start the server:
func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil { log.Fatalf("listen: %v", err) }

    srv := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            loggingInterceptor,
            recoveryInterceptor,
        ),
    )

    userv1.RegisterUserServiceServer(srv, &userServer{store: store.New()})

    log.Println("gRPC server listening on :50051")
    log.Fatal(srv.Serve(lis))
}
```

### Client
```go
func NewUserClient(addr string) (userv1.UserServiceClient, func(), error) {
    conn, err := grpc.NewClient(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithChainUnaryInterceptor(clientLoggingInterceptor),
    )
    if err != nil { return nil, nil, err }

    client := userv1.NewUserServiceClient(conn)
    cleanup := func() { conn.Close() }
    return client, cleanup, nil
}

// Using the client:
func main() {
    client, cleanup, err := NewUserClient("localhost:50051")
    if err != nil { log.Fatal(err) }
    defer cleanup()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    resp, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   30,
    })
    if err != nil {
        st := status.Convert(err)
        log.Fatalf("CreateUser: code=%v msg=%v", st.Code(), st.Message())
    }
    fmt.Printf("Created user: %v\n", resp.User)
}
```

---

## 4. Streaming RPC

### Server Streaming — push updates to client
```go
// Server sends a stream of User updates when the user changes.
func (s *userServer) WatchUser(req *userv1.GetUserRequest, stream userv1.UserService_WatchUserServer) error {
    // Initial send:
    user, err := s.store.GetByID(stream.Context(), req.Id)
    if err != nil { return status.Errorf(codes.NotFound, "user not found") }
    if err := stream.Send(domainToProto(user)); err != nil { return err }

    // Watch for changes:
    changes := s.store.Subscribe(req.Id)
    defer s.store.Unsubscribe(req.Id, changes)

    for {
        select {
        case <-stream.Context().Done():
            return nil  // Client disconnected
        case user, ok := <-changes:
            if !ok { return nil }  // Subscription closed
            if err := stream.Send(domainToProto(user)); err != nil { return err }
        }
    }
}

// Client consuming server stream:
stream, err := client.WatchUser(ctx, &userv1.GetUserRequest{Id: 42})
if err != nil { log.Fatal(err) }

for {
    user, err := stream.Recv()
    if err == io.EOF { break }
    if err != nil { log.Fatalf("recv: %v", err) }
    fmt.Printf("Update: %v\n", user)
}
```

### Client Streaming — client uploads batch
```go
// Client sends stream of users to create; server responds once.
func (s *userServer) BatchCreateUsers(stream userv1.UserService_BatchCreateUsersServer) error {
    var created int32

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&userv1.CreateUserResponse{
                // Return summary
            })
        }
        if err != nil { return err }

        _, err = s.store.Create(stream.Context(), req.Name, req.Email, int(req.Age))
        if err != nil { log.Printf("batch create error: %v", err) } else { created++ }
    }
}

// Client uploading a batch:
stream, err := client.BatchCreateUsers(ctx)
for _, user := range users {
    if err := stream.Send(&userv1.CreateUserRequest{Name: user.Name, Email: user.Email}); err != nil {
        log.Fatal(err)
    }
}
resp, err := stream.CloseAndRecv()
```

### Bidirectional Streaming
```go
// Both sides send and receive concurrently.
func (s *userServer) SyncUsers(stream userv1.UserService_SyncUsersServer) error {
    for {
        user, err := stream.Recv()
        if err == io.EOF { return nil }
        if err != nil { return err }

        // Process and echo back with server-side updates:
        updated := s.processSync(user)
        if err := stream.Send(updated); err != nil { return err }
    }
}
```

---

## 5. Interceptors (Middleware)

```go
// Unary interceptor (like HTTP middleware):
func loggingInterceptor(
    ctx context.Context,
    req any,
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (any, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    st := status.Convert(err)

    slog.Info("grpc request",
        "method", info.FullMethod,
        "duration", time.Since(start),
        "code", st.Code(),
    )
    return resp, err
}

func recoveryInterceptor(
    ctx context.Context,
    req any,
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp any, err error) {
    defer func() {
        if r := recover(); r != nil {
            slog.Error("panic recovered", "panic", r, "stack", string(debug.Stack()))
            err = status.Errorf(codes.Internal, "internal error")
        }
    }()
    return handler(ctx, req)
}

// Auth interceptor — validates JWT in metadata:
func authInterceptor(jwtSvc *JWTService) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        // Skip auth for public methods:
        if isPublicMethod(info.FullMethod) { return handler(ctx, req) }

        md, ok := metadata.FromIncomingContext(ctx)
        if !ok { return nil, status.Error(codes.Unauthenticated, "missing metadata") }

        tokens := md.Get("authorization")
        if len(tokens) == 0 { return nil, status.Error(codes.Unauthenticated, "missing token") }

        tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")
        claims, err := jwtSvc.ValidateToken(tokenStr)
        if err != nil { return nil, status.Error(codes.Unauthenticated, "invalid token") }

        ctx = context.WithValue(ctx, userKey, claims)
        return handler(ctx, req)
    }
}

// Client interceptor — injects token:
func clientAuthInterceptor(token string) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply any,
        cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
        return invoker(ctx, method, req, reply, cc, opts...)
    }
}
```

---

## 6. Error Handling

```go
// gRPC status codes map to HTTP status codes:
// codes.OK              → 200
// codes.InvalidArgument → 400
// codes.Unauthenticated → 401
// codes.PermissionDenied → 403
// codes.NotFound        → 404
// codes.AlreadyExists   → 409
// codes.ResourceExhausted → 429
// codes.Internal        → 500
// codes.Unavailable     → 503

// Rich error details with google.golang.org/grpc/status:
import "google.golang.org/genproto/googleapis/rpc/errdetails"

func richValidationError(fields map[string]string) error {
    st := status.New(codes.InvalidArgument, "validation failed")

    br := &errdetails.BadRequest{}
    for field, desc := range fields {
        br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
            Field:       field,
            Description: desc,
        })
    }

    detailed, err := st.WithDetails(br)
    if err != nil { return st.Err() }
    return detailed.Err()
}

// Client extracting error details:
func handleGRPCError(err error) {
    st := status.Convert(err)
    fmt.Printf("Code: %v, Message: %v\n", st.Code(), st.Message())

    for _, detail := range st.Details() {
        switch t := detail.(type) {
        case *errdetails.BadRequest:
            for _, v := range t.FieldViolations {
                fmt.Printf("Field %q: %v\n", v.Field, v.Description)
            }
        }
    }
}
```

---

## 7. gRPC-Gateway (REST bridge)

```protobuf
// Add HTTP annotations to the proto:
import "google/api/annotations.proto";

service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
    option (google.api.http) = {
      post: "/v1/users"
      body: "*"
    };
  }
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
    option (google.api.http) = {
      get: "/v1/users/{id}"
    };
  }
}
```

```go
// main.go — run gRPC and REST on different ports:
func main() {
    grpcServer := buildGRPCServer()
    go grpcServer.Serve(grpcListener)

    // gRPC-gateway reverse proxy:
    conn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    mux := runtime.NewServeMux()
    userv1.RegisterUserServiceHandlerClient(context.Background(), mux, userv1.NewUserServiceClient(conn))

    http.ListenAndServe(":8080", mux)  // REST on 8080, gRPC on 50051
}
```

---

## 8. Testing gRPC Services

```go
import "google.golang.org/grpc/test/bufconn"

const bufSize = 1024 * 1024

func TestUserService(t *testing.T) {
    // In-memory listener — no real network:
    lis := bufconn.Listen(bufSize)

    srv := grpc.NewServer()
    userv1.RegisterUserServiceServer(srv, &userServer{store: store.NewUserStore()})
    go srv.Serve(lis)
    t.Cleanup(srv.Stop)

    // Connect client via bufconn:
    conn, err := grpc.NewClient("passthrough:///bufnet",
        grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
            return lis.DialContext(ctx)
        }),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    require.NoError(t, err)
    t.Cleanup(func() { conn.Close() })

    client := userv1.NewUserServiceClient(conn)
    ctx := context.Background()

    t.Run("create and get user", func(t *testing.T) {
        createResp, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
            Name: "Alice", Email: "alice@test.com", Age: 30,
        })
        require.NoError(t, err)
        assert.Equal(t, "Alice", createResp.User.Name)

        getResp, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: createResp.User.Id})
        require.NoError(t, err)
        assert.Equal(t, createResp.User.Id, getResp.User.Id)
    })

    t.Run("not found returns codes.NotFound", func(t *testing.T) {
        _, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: 99999})
        require.Error(t, err)
        assert.Equal(t, codes.NotFound, status.Code(err))
    })

    t.Run("validation error returns codes.InvalidArgument", func(t *testing.T) {
        _, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
            Name: "", Email: "bad-email",
        })
        require.Error(t, err)
        assert.Equal(t, codes.InvalidArgument, status.Code(err))
    })
}
```

---

## Summary

- gRPC = HTTP/2 + Protobuf: binary, strongly-typed, streaming, fast
- `.proto` file defines messages and service — generates Go server interface + client stub
- 4 RPC types: unary (1→1), server streaming (1→N), client streaming (N→1), bidirectional (N→N)
- `UnimplementedXxxServer` embedding provides forward compatibility when new RPCs are added
- Interceptors = gRPC middleware for logging, auth, recovery
- Always return `status.Error(codes.Xxx, "message")` — never raw Go errors
- Use `bufconn` for in-memory testing without real network
- `grpc-gateway` bridges gRPC services to REST HTTP for browser/external clients

---

## Exercises

### Easy
1. Define a `TodoService` in Protobuf: `Todo` message (id, title, done bool, created_at), `CreateTodo`/`GetTodo`/`ListTodos`/`MarkDone` RPCs. Generate Go code and implement the server with an in-memory store. Write a `main.go` that creates 3 todos and lists them.
2. Add request validation interceptor: reject any `CreateTodo` where `title` is empty or > 200 chars before it reaches the handler. Write the interceptor so it extracts the title field using reflection or type assertion — interceptors receive `interface{}`, not the concrete type.
3. Implement a `PingPong` bidirectional streaming RPC: client sends `{message: "ping"}`, server responds with `{message: "pong", count: N}` where N increments per message. Test with 10 back-and-forth exchanges.

### Medium
4. Add **TLS to the gRPC server**: generate a self-signed cert (`openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes`). Configure the server with `credentials.NewServerTLSFromFile`. Configure the client with `credentials.NewClientTLSFromFile`. Verify that without TLS the client refuses to connect.
5. Implement **server-side rate limiting interceptor**: allow max 10 requests per second per client IP (extracted from `peer.FromContext`). Use a `sync.Map` of token buckets. Return `codes.ResourceExhausted` when exceeded.
6. Build a **gRPC health check** following the standard `grpc.health.v1` proto. Register `health.NewServer()` with your gRPC server. The `Check` RPC returns `SERVING` if the store is accessible, `NOT_SERVING` otherwise. Use `grpc_health_probe` CLI to test.

### Hard
7. Implement **retry with exponential backoff** on the client side using a `grpc.UnaryClientInterceptor`. Retry on `codes.Unavailable` and `codes.ResourceExhausted` up to 3 times with 100ms/200ms/400ms delays. Don't retry `codes.InvalidArgument`, `codes.NotFound`, or `codes.Unauthenticated` (these won't change on retry). Use `context.Context` deadline to abort early.
8. Build a **gRPC + REST dual server** using grpc-gateway: define a `BookService` proto with HTTP annotations. Generate both gRPC stubs and REST gateway code. Run both on different ports (gRPC: 50051, REST: 8080). Write a test that calls the same endpoint via both protocols and verifies identical results.
