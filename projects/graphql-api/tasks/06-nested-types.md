---
title: Nested Types
number: 6
difficulty: medium
duration: 20-25 minutes
concept: GraphQL Object Types, Nested Resolvers
---

## What You Need to Build

Add a `Post` type and resolve its nested `author` field.

## Schema Addition

```graphql
type Post {
  id: Int!
  title: String!
  authorId: Int!
  author: User
}

type Query {
  # existing: user, users, search
  posts: [Post!]!
}
```

## Package-Level Data

```go
var posts = []map[string]interface{}{
    {"id": 1, "title": "Hello World",   "authorId": 1},
    {"id": 2, "title": "Go is great",   "authorId": 2},
    {"id": 3, "title": "GraphQL rocks", "authorId": 1},
}
```

## What to Implement

- `postType` GraphQL object with `id`, `title`, `authorId`, and `author` fields
- `author` resolver: look up the user whose `id` matches `authorId`
- Add `posts` to the root query

## Requirements

Query `{ posts { title author { name } } }` must return posts with nested author names.

## Hints

<details>
<summary>Hint 1 — Resolving the author field</summary>

In the `author` field resolver, the source (`p.Source`) is the post map. Read `authorId` from it, then find the matching user in the users slice.

</details>

<details>
<summary>Hint 2 — Field resolver signature</summary>

```go
Resolve: func(p graphql.ResolveParams) (interface{}, error) {
    post := p.Source.(map[string]interface{})
    authorId := post["authorId"].(int)
    // find user with that id
},
```

</details>

## How to Verify

```bash
lncli run
```
