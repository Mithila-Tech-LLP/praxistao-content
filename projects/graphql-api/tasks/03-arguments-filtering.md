---
title: Task 03 — Arguments and Filtering
task: "03"
slug: arguments-filtering
starter: task-03-arguments-filtering
---

## Goal

Add a `search` query that accepts a string argument and returns only users whose name contains the search term (case-insensitive).

## Background

You already used arguments in Task 01 for `user(id: Int!)`. Arguments can be any scalar type — here you will use a non-null `String`:

```go
Args: graphql.FieldConfigArgument{
    "name": &graphql.ArgumentConfig{
        Type: graphql.NewNonNull(graphql.String),
    },
},
```

Inside the resolver, read it back:

```go
query := p.Args["name"].(string)
```

## Your Task

Open `starter/task-03-arguments-filtering/main.go`.

1. Copy the `userType` and `users` slice from Task 02.

2. Add a root query field `"search"` that:
   - Returns `graphql.NewList(userType)`
   - Takes one required argument: `name` of type `String!`
   - Resolves by iterating `users` and keeping only those where the user's name **contains** the search term (case-insensitive)

## Key APIs

```go
import "strings"

strings.Contains(
    strings.ToLower(candidate),
    strings.ToLower(query),
)
```

## Running the Test

```
cd starter/task-03-arguments-filtering
go test -v
```

Tests:
- `TestSearchExact` — `{ search(name: "Alice") { name } }` returns exactly ["Alice"]
- `TestSearchPartial` — `{ search(name: "ali") { name } }` returns ["Alice"] (case-insensitive)
- `TestSearchMultiple` — `{ search(name: "o") { name } }` returns ["Bob", "Carol"]
- `TestSearchNoMatch` — `{ search(name: "xyz") { name } }` returns an empty list
