---
title: Task 01 — Schema and Simple Query
task: "01"
slug: schema-query
starter: task-01-schema-query
---

## Goal

Define your first GraphQL schema in Go, wire up a single `user` query, and return hardcoded data from a resolver.

## Background

A GraphQL schema is a contract: it declares every type your API knows about and every operation a client can perform. At minimum you need three pieces:

1. **An object type** — describes the shape of a resource (`User`, with fields `id` and `name`).
2. **A root query type** — the entry point for all read operations.
3. **A schema** — wraps the root query (and later a mutation root).

In `graphql-go` you build all three programmatically:

```go
userType := graphql.NewObject(graphql.ObjectConfig{
    Name: "User",
    Fields: graphql.Fields{
        "id":   &graphql.Field{Type: graphql.Int},
        "name": &graphql.Field{Type: graphql.String},
    },
})
```

A **resolver** is just a Go function that returns the data for a field:

```go
Resolve: func(p graphql.ResolveParams) (interface{}, error) {
    return map[string]interface{}{"id": 1, "name": "Alice"}, nil
},
```

## Your Task

Open `starter/task-01-schema-query/main.go` and implement `BuildSchema()`.

1. Define `userType` as a package-level variable — a `graphql.NewObject` with:
   - `Name`: `"User"`
   - Fields: `id` (`graphql.Int`) and `name` (`graphql.String`)

2. Inside `BuildSchema`, create a root query object (`graphql.NewObject`) with:
   - `Name`: `"Query"`
   - A field named `"user"` that:
     - Returns `userType`
     - Accepts one argument: `id` of type `Int!` (non-null Int)
     - Has a resolver that returns `map[string]interface{}{"id": 1, "name": "Alice"}` for any id

3. Return `graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})`.

## Key APIs

| What you need | How to get it |
|---|---|
| Non-null wrapper | `graphql.NewNonNull(graphql.Int)` |
| Declare an argument | `graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: ...}}` |
| Build the schema | `graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})` |

## Running the Test

```
cd starter/task-01-schema-query
go test -v
```

Both tests must pass:
- `TestSchemaBuilds` — `BuildSchema()` returns no error
- `TestUserQuery` — `{ user(id: 1) { id name } }` returns id=1 and name="Alice"
