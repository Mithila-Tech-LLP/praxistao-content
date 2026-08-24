---
title: Task 05 — Mutations: Update and Delete
task: "05"
slug: mutations-update-delete
starter: task-05-mutations-update-delete
---

## Goal

Extend the mutation root with `updateUser` (returns the updated user, or null if not found) and `deleteUser` (returns a Boolean).

## Background

Mutations can return any output type — including existing object types, scalars, or null. Returning `nil, nil` from a resolver maps to `null` in the JSON response, which is the standard way to signal "not found" for a nullable field.

`graphql.Boolean` is the scalar type for Go's `bool`.

## Your Task

Open `starter/task-05-mutations-update-delete/main.go`. Start from your Task 04 solution (with `createUser` already in place) and add two more mutation fields.

### `updateUser(id: Int!, name: String!) → User`

1. Read `id` and `name` from `p.Args`.
2. Loop over `users`; when `u["id"].(int) == id`, set `users[i]["name"] = name` and return `users[i], nil`.
3. If no match, return `nil, nil` (null).

### `deleteUser(id: Int!) → Boolean`

1. Read `id` from `p.Args`.
2. Loop over `users`; when a match is found, remove that element:
   ```go
   users = append(users[:i], users[i+1:]...)
   return true, nil
   ```
3. If no match, return `false, nil`.

## Running the Test

```
cd starter/task-05-mutations-update-delete
go test -v
```

Tests:
- `TestUpdateUser` — updates Alice's name to "Alicia", query confirms the change
- `TestUpdateNonExistent` — updating id=999 returns null (no errors in `result.Errors`)
- `TestDeleteUser` — deletes Bob, list shrinks to 2 users
- `TestDeleteNonExistent` — deleting id=999 returns false
