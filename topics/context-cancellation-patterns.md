---
title: Context Cancellation Patterns
category: Software & Programming
tags: [Go, Concurrency]
duration: 8 min read
relatedCourses: [go-programming, senior-engineer-interview]
relatedProjects: [rest-api-server, grpc-service]
relatedTopics: [goroutine-leaks, worker-pool-patterns]
---

## TL;DR

- `context.Context` threads a single cancellation/deadline signal through an entire call chain — every function that might block should accept one as its first parameter.
- Cancellation is cooperative: calling `cancel()` doesn't forcibly stop anything. It closes `ctx.Done()`, and it's up to your code to check that channel and actually return.
- Four ways to create one: `context.Background()` (root), `context.WithCancel`, `context.WithTimeout`/`WithDeadline`, and `context.WithValue` (for request-scoped data, not control flow).
- Always `defer cancel()` immediately after creating a cancellable context — even if the operation finishes normally, `cancel()` releases the context's internal resources.

## Why Context Exists

Say an HTTP handler calls a database query, which internally calls another service over gRPC, which internally does a cache lookup. If the original HTTP client disconnects, you want all three of those nested operations to stop — not keep running to completion for a response nobody will read.

Before `context.Context` existed (pre-Go 1.7), every library invented its own ad-hoc way to pass a "please stop" signal down a call stack, if it bothered at all. `context` standardized this into one interface everyone agrees on, and it now threads through the standard library itself — `net/http`, `database/sql`, and virtually every serious Go library accept one.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // cancelled automatically if the client disconnects
    result, err := fetchUser(ctx, userID)
    ...
}

func fetchUser(ctx context.Context, id string) (*User, error) {
    row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)
    ...
}
```

## The Four Ways to Get a Context

```go
ctx := context.Background()        // the root — used at program/request entry points
ctx := context.TODO()              // "I haven't decided yet" — same as Background, signals intent to revisit

ctx, cancel := context.WithCancel(parent)               // cancel it manually, whenever you want
ctx, cancel := context.WithTimeout(parent, 5*time.Second) // cancelled automatically after a duration
ctx, cancel := context.WithDeadline(parent, someTime)     // cancelled automatically at a specific time

ctx = context.WithValue(parent, key, value)  // attaches request-scoped data — NOT for control flow
```

`WithCancel`/`WithTimeout`/`WithDeadline` all return a `cancel func()` that you must call — usually via `defer cancel()` right after creation — to release resources tied to that context (an internal timer, for `WithTimeout`/`WithDeadline`) even when everything finished normally without needing cancellation.

## The Core Pattern: select on ctx.Done()

Any function that blocks — on a channel, a network call, a long computation — needs to race that blocking operation against `ctx.Done()`:

```go
func worker(ctx context.Context, jobs <-chan int) {
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                return
            }
            process(job)
        case <-ctx.Done():
            fmt.Println("cancelled:", ctx.Err())
            return
        }
    }
}
```

`ctx.Done()` returns a channel that's closed exactly once, when the context is cancelled or its deadline passes — closing a channel is Go's built-in way to broadcast to any number of listening goroutines simultaneously, which is exactly the semantics cancellation needs. `ctx.Err()` then tells you *why*: `context.Canceled` (someone called `cancel()`) or `context.DeadlineExceeded` (a timeout/deadline fired).

## Propagation: Contexts Form a Tree

Calling `context.WithCancel(parent)` creates a *child* of `parent`. Cancelling the parent automatically cancels every child (and grandchild) derived from it — but cancelling a child has no effect on its parent or siblings.

```
context.Background()
  └── ctx1 = WithTimeout(bg, 10s)         // request-level deadline
        ├── ctx2 = WithCancel(ctx1)       // one sub-operation
        └── ctx3 = WithTimeout(ctx1, 2s)  // another sub-operation, shorter deadline
```

If `ctx1`'s 10-second timeout fires, both `ctx2` and `ctx3` are cancelled too, even though nothing directly cancelled them. This is what makes "cancel the whole request tree from one place" work without manually tracking every derived context.

Note that a child's deadline can only be *tighter* than its parent's, never looser — `WithTimeout(ctx1, 30s)` where `ctx1` already has a 10-second deadline still gets cancelled at 10 seconds, because the parent's deadline still applies.

## context.Value — Handle With Care

`context.WithValue` is for request-scoped metadata that has to cross API boundaries you don't control — a request ID for logging, an authenticated user extracted by middleware. It is explicitly **not** meant for passing optional parameters into functions you do control; those should just be normal function arguments.

```go
type ctxKey string
const userKey ctxKey = "user"

ctx = context.WithValue(ctx, userKey, currentUser)
// later, deep in the call stack:
user, ok := ctx.Value(userKey).(*User)
```

Using an unexported custom type for the key (`ctxKey` above, not a bare `string`) prevents collisions with keys some other package might use — this is a real, commonly-hit bug if you skip it.

## Common Pitfalls

- **Storing a context in a struct field** — contexts are meant to flow through function call chains as an explicit parameter (conventionally named `ctx` and passed first), not be stashed for later use. A stored context can outlive the request it belonged to and silently never get cancelled when you expect.
- **Passing `nil` instead of `context.Background()`** — a nil context will panic the moment something calls `.Done()` or `.Value()` on it. If you truly don't have a context yet, use `context.TODO()`.
- **Forgetting `defer cancel()`** — this doesn't leak your business goroutine, but it does leak the small internal timer goroutine that `WithTimeout`/`WithDeadline` create, until that timeout eventually fires on its own.
- **Ignoring `ctx.Done()` in a long-running loop** — a function that accepts a context but never actually checks `ctx.Done()` anywhere is decorative. The context being cancelled changes nothing if nothing is listening for it.
