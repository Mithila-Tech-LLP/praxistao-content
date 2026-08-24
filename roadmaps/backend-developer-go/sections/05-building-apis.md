---
title: Building APIs
---
This is the section most people picture when they hear "backend developer": designing and building the interfaces that other services and clients actually call. REST is the default; gRPC and GraphQL solve problems REST doesn't handle as well.

### HTTP Fundamentals & REST
Before any framework, understand what actually happens on the wire — requests, status codes, headers — and the conventions REST layers on top of it.

**Resources:**
- [HTTP Fundamentals](course:go-programming#59-http-fundamentals)
- [Building REST APIs](course:go-programming#60-building-rest-apis)

### Middleware, Auth & JWT
Middleware is how cross-cutting concerns (logging, auth, rate limiting) get applied without repeating yourself in every handler. JWTs are the most common way real APIs authenticate requests.

**Resources:**
- [Middleware](course:go-programming#61-middleware)
- [Authentication and JWT](course:go-programming#64-authentication-and-jwt)
- [JSON Validation and Serialization](course:go-programming#62-json-validation-serialization)

### Practice: Build a REST API Server
> branches-from: Middleware, Auth & JWT

Put it together: routing, middleware, auth, and a real datastore behind a REST API you build and test yourself from scratch.

**Resources:**
- [REST API Server project](project:rest-api-server)

### gRPC
When services talk to other services (not browsers), gRPC's binary protocol and strict contracts often beat REST on both performance and safety.

**Resources:**
- [gRPC](course:go-programming#66-grpc)

### Practice: Build a gRPC Service
> branches-from: gRPC

**Resources:**
- [gRPC Service project](project:grpc-service)

### GraphQL
> optional

GraphQL flips the REST model: instead of the server deciding what a response contains, the client asks for exactly the fields it needs. There's no dedicated course chapter for this yet — the project below is a learn-by-building introduction.

**Resources:**
- [GraphQL API project](project:graphql-api)

### Practice: Build a GraphQL API
> optional
> branches-from: GraphQL

**Resources:**
- [GraphQL API project](project:graphql-api)
