---
title: Error Handling
number: 8
difficulty: medium
duration: 20-25 minutes
concept: GraphQL Errors, Resolver Errors, Validation
---

## What You Need to Build

Return proper GraphQL errors when operations fail.

## Requirements

1. `user(id: 999)` — user not found → return `null` for the field, add to `errors` array
2. `createUser(input: {name: "", email: "a@b.com"})` — empty name → error, no user created
3. `createUser(input: {name: "Alice", email: ""})` — empty email → error, no user created

## How GraphQL Errors Work

In `graphql-go`, returning an error from a resolver adds it to the `errors` array in the response while still returning the rest of the data. The field value becomes `null`.

```go
Resolve: func(p graphql.ResolveParams) (interface{}, error) {
    if notFound {
        return nil, fmt.Errorf("user not found: id=%d", id)
    }
    return user, nil
},
```

## Hints

<details>
<summary>Hint 1 — Checking the errors array in tests</summary>

```go
result := graphql.Do(graphql.Params{Schema: schema, RequestString: query})
if len(result.Errors) == 0 {
    t.Error("expected errors, got none")
}
```

</details>

<details>
<summary>Hint 2 — Validation in createUser resolver</summary>

Check `name == ""` and `email == ""` at the start of the resolver. Return an error without mutating the users slice.

</details>

## How to Verify

```bash
lncli run
```
