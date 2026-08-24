---
title: Task 02 — List Query
task: "02"
slug: list-query
starter: task-02-list-query
---

## Goal

Add an in-memory store of users and expose a `users` query that returns the whole list.

## Background

GraphQL lists are written as `[User]` in the schema language. In `graphql-go` you wrap any output type with `graphql.NewList(...)`:

```go
"users": &graphql.Field{
    Type: graphql.NewList(userType),
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        // return a slice
    },
},
```

The resolver can return `[]interface{}` or any Go slice — `graphql-go` uses reflection to iterate it, calling each field's resolver for every element.

## Your Task

Open `starter/task-02-list-query/main.go`.

1. Declare a package-level variable:
   ```go
   var users = []map[string]interface{}{
       {"id": 1, "name": "Alice"},
       {"id": 2, "name": "Bob"},
       {"id": 3, "name": "Carol"},
   }
   ```

2. Keep the `user(id: Int!) → User` field from Task 01, but update its resolver to look up the user by id in the `users` slice. Return `nil, nil` if not found.

3. Add a new root query field `"users"` that:
   - Returns `graphql.NewList(userType)`
   - Takes no arguments
   - Resolves by converting `users` to `[]interface{}` and returning it

## Hint

To convert `[]map[string]interface{}` to `[]interface{}`:

```go
result := make([]interface{}, len(users))
for i, u := range users {
    result[i] = u
}
return result, nil
```

## Running the Test

```
cd starter/task-02-list-query
go test -v
```

Tests:
- `TestUserByID` — `{ user(id: 2) { name } }` returns "Bob"
- `TestListUsers` — `{ users { id name } }` returns all 3 users
