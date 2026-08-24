---
title: Input Types
number: 7
difficulty: medium
duration: 20-25 minutes
concept: GraphQL Input Types, Structured Mutation Arguments
---

## What You Need to Build

Replace the scalar mutation argument with a structured `UserInput` input type.

## Schema Change

```graphql
input UserInput {
  name: String!
  email: String!
}

type Mutation {
  createUser(input: UserInput!): User
}
```

## User Type Update

Add `email` field to the User type:
```go
type User struct{ ID int; Name, Email string }
```

## Requirements

- `createUser(input: {name: "Alice", email: "alice@example.com"})` creates a user with both fields
- Return the created user with all three fields (id, name, email)

## Hints

<details>
<summary>Hint 1 — Defining the input type</summary>

```go
userInputType := graphql.NewInputObject(graphql.InputObjectConfig{
    Name: "UserInput",
    Fields: graphql.InputObjectConfigFieldMap{
        "name":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
        "email": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
    },
})
```

</details>

<details>
<summary>Hint 2 — Reading the input in the resolver</summary>

```go
input := p.Args["input"].(map[string]interface{})
name := input["name"].(string)
email := input["email"].(string)
```

</details>

## How to Verify

```bash
lncli run
```
