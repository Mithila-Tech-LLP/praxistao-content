---
title: Task 04 — Mutations: Create
task: "04"
slug: mutations-create
starter: task-04-mutations-create
---

## Goal

Add a root mutation type and implement `createUser(name: String!) → User`, which appends a new user to the in-memory store.

## Background

GraphQL separates reads (queries) from writes (mutations). A mutation is just a field on a special root type:

```go
rootMutation := graphql.NewObject(graphql.ObjectConfig{
    Name: "Mutation",
    Fields: graphql.Fields{
        "createUser": &graphql.Field{ ... },
    },
})
```

Pass it to the schema builder alongside the query:

```go
graphql.NewSchema(graphql.SchemaConfig{
    Query:    rootQuery,
    Mutation: rootMutation,
})
```

Client syntax for mutations:

```graphql
mutation {
    createUser(name: "Dave") {
        id
        name
    }
}
```

## Your Task

Open `starter/task-04-mutations-create/main.go`.

1. Copy `userType` and `users` from Task 02. Add:
   ```go
   var nextID = 4
   ```

2. Build a root mutation object with one field `"createUser"` that:
   - Takes a required `name: String!` argument
   - Creates `map[string]interface{}{"id": nextID, "name": name}`
   - Increments `nextID`
   - Appends the new user to `users`
   - Returns the new user

3. Pass `Mutation: rootMutation` to `graphql.NewSchema`.

## Running the Test

```
cd starter/task-04-mutations-create
go test -v
```

Tests:
- `TestCreateUser` — mutation returns the new user with correct name and an id >= 4
- `TestCreateUserIncreasesCount` — after the mutation, `{ users { id name } }` returns 4 users
