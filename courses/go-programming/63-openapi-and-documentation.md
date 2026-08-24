# Chapter 47: OpenAPI and API Documentation

An API without documentation is a black box. OpenAPI (formerly Swagger) is the standard for describing REST APIs: request/response schemas, authentication, error codes, and examples. With an OpenAPI spec you get: auto-generated docs, client SDKs in any language, contract testing, and mock servers.

## Table of Contents

1. [OpenAPI 3.0 Spec Structure](#1-openapi-30-spec-structure)
2. [Code-First with ogen](#2-code-first-with-ogen)
3. [Spec-First with oapi-codegen](#3-spec-first-with-oapi-codegen)
4. [Serving Swagger UI](#4-serving-swagger-ui)
5. [Validation Middleware](#5-validation-middleware)
6. [Versioning Strategies](#6-versioning-strategies)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. OpenAPI 3.0 Spec Structure

```yaml
# openapi.yaml
openapi: "3.0.3"
info:
  title: User API
  version: "1.0.0"
  description: Manage users in the system
  contact:
    email: api@example.com

servers:
  - url: https://api.example.com/v1
    description: Production
  - url: http://localhost:8080/v1
    description: Local development

# Reusable components:
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    User:
      type: object
      required: [id, name, email, createdAt]
      properties:
        id:
          type: integer
          format: int64
          example: 42
        name:
          type: string
          example: Alice
        email:
          type: string
          format: email
          example: alice@example.com
        age:
          type: integer
          minimum: 0
          maximum: 150
        createdAt:
          type: string
          format: date-time

    CreateUserRequest:
      type: object
      required: [name, email]
      properties:
        name:
          type: string
          minLength: 1
          maxLength: 100
        email:
          type: string
          format: email
        age:
          type: integer
          minimum: 0

    Error:
      type: object
      required: [error]
      properties:
        error:
          type: string

    ValidationError:
      type: object
      required: [error, fields]
      properties:
        error:
          type: string
          example: validation failed
        fields:
          type: object
          additionalProperties:
            type: string
          example:
            email: invalid format
            name: required

  parameters:
    UserID:
      name: userID
      in: path
      required: true
      schema:
        type: integer
        format: int64

    Page:
      name: page
      in: query
      schema:
        type: integer
        default: 1
        minimum: 1

    PageSize:
      name: pageSize
      in: query
      schema:
        type: integer
        default: 20
        minimum: 1
        maximum: 100

  responses:
    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    Unauthorized:
      description: Authentication required
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'

# Security applied globally:
security:
  - bearerAuth: []

# API endpoints:
paths:
  /users:
    get:
      operationId: listUsers
      summary: List users
      tags: [Users]
      parameters:
        - $ref: '#/components/parameters/Page'
        - $ref: '#/components/parameters/PageSize'
        - name: search
          in: query
          schema:
            type: string
      responses:
        "200":
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/User'
                  total:
                    type: integer
                  page:
                    type: integer
                  pageSize:
                    type: integer
        "401":
          $ref: '#/components/responses/Unauthorized'

    post:
      operationId: createUser
      summary: Create a user
      tags: [Users]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        "201":
          description: User created
          headers:
            Location:
              schema:
                type: string
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/User'
        "409":
          description: Email already taken
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        "422":
          description: Validation error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ValidationError'

  /users/{userID}:
    parameters:
      - $ref: '#/components/parameters/UserID'
    get:
      operationId: getUser
      summary: Get a user by ID
      tags: [Users]
      responses:
        "200":
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/User'
        "404":
          $ref: '#/components/responses/NotFound'
    delete:
      operationId: deleteUser
      summary: Delete a user
      tags: [Users]
      responses:
        "204":
          description: User deleted
        "404":
          $ref: '#/components/responses/NotFound'
```

---

## 2. Code-First with ogen

`ogen` generates a complete, type-safe server from an OpenAPI spec.

```bash
go get github.com/ogen-go/ogen/cmd/ogen@latest
ogen --target gen/api --clean openapi.yaml
```

```go
// Implement the generated interface:
type userAPIHandler struct {
    store *store.UserStore
}

// ListUsers implements api.Handler.
func (h *userAPIHandler) ListUsers(ctx context.Context, params api.ListUsersParams) (*api.ListUsersOK, error) {
    page := params.Page.Or(1)
    pageSize := params.PageSize.Or(20)

    users, total := h.store.List((page-1)*pageSize, pageSize)
    result := make([]api.User, len(users))
    for i, u := range users {
        result[i] = api.User{
            ID:    u.ID,
            Name:  u.Name,
            Email: u.Email,
            Age:   api.NewOptInt32(int32(u.Age)),
        }
    }
    return &api.ListUsersOK{
        Data:     result,
        Total:    int32(total),
        Page:     int32(page),
        PageSize: int32(pageSize),
    }, nil
}

// CreateUser implements api.Handler.
func (h *userAPIHandler) CreateUser(ctx context.Context, req *api.CreateUserRequest) (api.CreateUserRes, error) {
    user, err := h.store.Create(ctx, req.Name, req.Email, int(req.Age.Or(0)))
    if err != nil {
        if errors.Is(err, store.ErrEmailTaken) {
            return &api.CreateUserConflict{Error: "email already taken"}, nil
        }
        return nil, err
    }
    return &api.CreateUserCreated{
        Data: api.User{ID: user.ID, Name: user.Name, Email: user.Email},
    }, nil
}

// Wire it up:
func main() {
    handler := &userAPIHandler{store: store.NewUserStore()}
    srv, err := api.NewServer(handler)
    if err != nil { log.Fatal(err) }

    http.ListenAndServe(":8080", srv)
}
```

---

## 3. Spec-First with oapi-codegen

```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
oapi-codegen -generate types,chi-server,spec -package api openapi.yaml > gen/api/api.gen.go
```

```go
// Implement the generated StrictServerInterface:
type userServer struct {
    store *store.UserStore
}

func (s *userServer) ListUsers(ctx context.Context, req api.ListUsersRequestObject) (api.ListUsersResponseObject, error) {
    page := 1
    if req.Params.Page != nil { page = *req.Params.Page }

    users, total := s.store.List((page-1)*20, 20)
    // ... convert and return
    return api.ListUsers200JSONResponse{
        Data:  convertUsers(users),
        Total: total,
    }, nil
}

// main.go:
r := chi.NewRouter()
strictHandler := api.NewStrictHandler(userServer, nil)
api.HandlerFromMux(strictHandler, r)
http.ListenAndServe(":8080", r)
```

---

## 4. Serving Swagger UI

```go
// Embed the OpenAPI spec in the binary:
//go:embed openapi.yaml
var openapiSpec []byte

func RegisterDocs(mux *http.ServeMux) {
    // Serve raw spec:
    mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/yaml")
        w.Write(openapiSpec)
    })

    // Serve Swagger UI (via CDN in HTML, or embed the dist):
    mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, swaggerUIHTML, "/openapi.yaml")
    })
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <title>API Docs</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script>
SwaggerUIBundle({
  url: "%s",
  dom_id: '#swagger-ui',
  presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
  layout: "BaseLayout"
});
</script>
</body>
</html>`

// Or use go-swagger for self-contained docs:
// go get github.com/swaggo/swag/cmd/swag
// go get github.com/swaggo/http-swagger
// Add annotations to handlers:

// @Summary Create a user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 422 {object} ValidationError
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) { ... }
```

---

## 5. Validation Middleware

```go
// Validate incoming requests against the OpenAPI spec:
import "github.com/getkin/kin-openapi/openapi3filter"

func OpenAPIValidation(spec []byte) func(http.Handler) http.Handler {
    loader := openapi3.NewLoader()
    doc, err := loader.LoadFromData(spec)
    if err != nil { panic(err) }
    if err := doc.Validate(loader.Context); err != nil { panic(err) }

    router, err := gorillamux.NewRouter(doc)
    if err != nil { panic(err) }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            route, pathParams, err := router.FindRoute(r)
            if err != nil {
                // Route not in spec — let the app handle it (404)
                next.ServeHTTP(w, r)
                return
            }

            input := &openapi3filter.RequestValidationInput{
                Request:    r,
                PathParams: pathParams,
                Route:      route,
            }
            if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
                var reqErr *openapi3filter.RequestError
                if errors.As(err, &reqErr) {
                    writeJSON(w, http.StatusBadRequest, map[string]string{
                        "error": reqErr.Reason,
                    })
                    return
                }
                writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 6. Versioning Strategies

```yaml
# URL versioning (most common — visible, cacheable):
servers:
  - url: https://api.example.com/v1
  - url: https://api.example.com/v2

# Header versioning (clean URLs, harder to test):
parameters:
  - name: API-Version
    in: header
    required: false
    schema:
      type: string
      enum: [v1, v2]
      default: v1

# Content negotiation:
# Accept: application/vnd.example.v2+json
```

**Deprecation in OpenAPI:**
```yaml
paths:
  /users/{id}/profile:
    get:
      deprecated: true
      description: "Deprecated: use GET /users/{id} which now includes profile fields."
```

```go
// Middleware to add Deprecation header on deprecated routes:
func DeprecationMiddleware(deprecatedPaths map[string]string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if msg, ok := deprecatedPaths[r.URL.Path]; ok {
                w.Header().Set("Deprecation", "true")
                w.Header().Set("Sunset", "2026-12-31")
                w.Header().Set("Link", `</v2`+r.URL.Path+`>; rel="successor-version"`)
                w.Header().Set("Warning", `299 - "`+msg+`"`)
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Summary

- OpenAPI 3.0 YAML: `info`, `servers`, `paths`, `components` — describe every endpoint, schema, and error
- **Code-first** (`ogen`, `oapi-codegen`): write spec, generate server boilerplate and client — ensures spec stays accurate
- **Annotation-first** (`swag`): annotate Go handler functions, generate spec — easier to add to existing code
- Serve `/openapi.yaml` and Swagger UI at `/docs` for consumers and explorers
- `kin-openapi` validates incoming requests against the spec before they reach handlers
- Deprecate endpoints with `deprecated: true` in spec + `Deprecation`/`Sunset` response headers
- API contract is a promise to consumers — version carefully and give migration time

---

## Exercises

### Easy
1. Write a complete OpenAPI 3.0 spec for the Todo API from Chapter 46 (gRPC exercises). Include: `Todo` schema, `ListTodos`, `CreateTodo`, `GetTodo`, `MarkDone` operations. Use `$ref` for shared schemas.
2. Serve the spec and Swagger UI from your Chapter 41 REST API. Navigate to `/docs` in a browser and verify all endpoints appear with correct request/response schemas.
3. Add the `deprecated: true` flag to an old endpoint and add a `Deprecation` response header middleware. Verify the header appears in responses and Swagger UI shows the endpoint as deprecated.

### Medium
4. Use `oapi-codegen` to generate server stubs from your spec. Implement the generated interface with the in-memory store from Chapter 41. Write a test that sends a request with a missing required field and verifies the validation middleware returns `400 Bad Request` (not `422` — OpenAPI validation happens before the handler).
5. Add **examples** to your OpenAPI spec: one happy-path example for each operation (request + response). Add a `400` and `422` error example for `createUser`. Verify they appear in Swagger UI. Examples make APIs 10× easier to consume.
6. Generate a Go client from your OpenAPI spec using `oapi-codegen --generate client`. Write an integration test that uses the generated client (not raw `http.Client`) to call your API. Compare the ergonomics with raw HTTP calls.

### Hard
7. **OpenAPI contract test**: write a test that loads your OpenAPI spec and verifies every response from your integration test suite conforms to the spec schema (not just status code). Use `kin-openapi`'s response validator. This catches cases where you return an undocumented field or wrong type.
8. **Changelog-driven versioning**: build a tool that diffs two OpenAPI specs (v1 vs v2) and categorizes each change as: breaking (removed field, changed type, removed endpoint), non-breaking (added field, new endpoint), or deprecated (marked deprecated). Output a human-readable changelog. Use the `kin-openapi` library to parse both specs.
