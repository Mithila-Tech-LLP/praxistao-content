---
title: Idempotency in APIs
category: Software & Programming
tags: [System Design, APIs]
duration: 7 min read
relatedCourses: [go-programming, senior-engineer-interview]
relatedProjects: [rest-api-server]
relatedTopics: [rate-limiting-algorithms, circuit-breakers-and-retries]
---

## TL;DR

- An operation is idempotent if running it once has the same effect as running it many times — critical for any operation that might get retried after a network failure.
- `GET`, `PUT`, and `DELETE` are supposed to be idempotent by HTTP's own convention; `POST` is not, by default — and "creating a charge" or "creating an order" is exactly the kind of `POST` that most needs to be made idempotent in practice.
- The standard fix is a client-generated **idempotency key**: the server remembers which keys it's already processed and returns the original result for a repeat, instead of doing the operation again.
- This matters specifically because "did my request succeed?" is genuinely unknowable to a client after a timeout — the request might have succeeded and the *response* got lost, not the request itself.

## Why Retries Make This Necessary

A client sends `POST /charge` to charge a customer $50. The request reaches the server, the charge succeeds, the server sends back a 200 — but the response is lost somewhere on the way back (a dropped connection, a load balancer restart, whatever). From the client's point of view, it just sees a timeout. It has no way to distinguish "my request never arrived" from "my request succeeded but the response didn't come back."

A naive client retries the request. If `POST /charge` isn't idempotent, the customer just got charged twice — for a failure that, from the server's perspective, wasn't a failure at all.

This is not a rare edge case: it is the single most common reason idempotency matters in real systems, because **any client that retries on timeout** (which is standard, sensible client behavior) creates exactly this scenario for any non-idempotent operation.

## HTTP's Own Idempotency Contract

The HTTP spec assigns idempotency expectations to methods:

| Method | Idempotent? | Meaning |
|---|---|---|
| `GET` | Yes | Reading data twice has the same effect as reading it once (no effect at all) |
| `PUT` | Yes | "Set this resource to this exact state" — doing it twice leaves it in the same state as doing it once |
| `DELETE` | Yes | Deleting an already-deleted resource is still "not present" — same end state |
| `POST` | **No** | Conventionally means "create a new thing" or "perform an action" — doing it twice, by default, does the thing twice |

This is *convention*, not something HTTP itself enforces — nothing stops you from writing a `PUT` handler with a side effect that isn't actually idempotent. But clients, proxies, and browsers (which will sometimes automatically retry idempotent methods on connection failure, but never automatically retry a `POST`) rely on this contract being honestly implemented.

## The Idempotency Key Pattern

For operations that are inherently "do a thing" (charge a card, place an order, send a notification) rather than "set a resource to a state," the standard fix is a client-supplied idempotency key:

```
POST /charges
Idempotency-Key: 7c9e1f2a-...  (client-generated UUID, unique per logical attempt)

{ "amount": 5000, "currency": "usd" }
```

Server-side logic:

```go
func handleCharge(key string, req ChargeRequest) (ChargeResult, error) {
    if existing, ok := store.Get(key); ok {
        return existing, nil // already processed — return the ORIGINAL result, don't charge again
    }

    result, err := processCharge(req)
    if err != nil {
        return result, err // don't store a key for a failed attempt — allow retry to actually retry
    }

    store.Set(key, result, ttl: 24*time.Hour)
    return result, nil
}
```

The critical detail: the key is generated **once per logical operation** on the client side, before the first attempt — and the *same* key is reused on every retry of that same logical operation. If the client generated a new key on each retry, this scheme provides zero protection, because the server would never recognize the retry as a repeat of anything.

This is exactly the pattern Stripe's API popularized and most payment/order-processing APIs now use directly.

## What Gets Stored, and for How Long

The server needs to remember, per idempotency key: what the result was (so it can return the *same* result on a repeat, not just refuse to redo the work) and for how long (an idempotency key doesn't need to be remembered forever — a TTL of 24 hours is common, since retries happen within seconds-to-minutes of the original attempt in practice, not days later).

A subtlety worth being deliberate about: what happens if a *second, concurrent* request with the same key arrives while the *first* one is still being processed (not yet stored)? A correct implementation needs to handle this — typically by storing an "in-progress" marker atomically before starting the actual work, and having the second request either wait for the first to finish or return a 409/425 indicating "this operation is already in flight."

## Idempotency vs Exactly-Once

It's worth being precise about what idempotency actually buys you: it does not mean the operation only ever *executes* once at the storage layer — it means that no matter how many times the *client-visible* operation is invoked with the same key, the *effect* and the *result returned* are the same as if it happened once. Underneath, the server might do real work on the first call and nothing but a lookup on every repeat — that's the mechanism, not a violation of the guarantee.

This distinction matters because "exactly-once delivery" in messaging systems is a related but different, and genuinely harder, problem — most message brokers only guarantee at-least-once delivery, and it's the *consumer's* idempotency (processing the same message twice safely) that turns at-least-once delivery into effectively-exactly-once processing.

## Common Pitfalls

- **Making the idempotency key optional** — if clients can skip sending one, they'll skip it exactly in the failure-prone conditions (flaky networks, retry logic bolted on later) where it matters most. Require it for any endpoint where duplication is dangerous.
- **Generating a new key on every retry attempt** — this defeats the entire mechanism; the same logical operation must reuse the same key across all its retry attempts.
- **Only checking the key, not storing the actual result** — if a repeat request just gets "200 OK, no body" instead of the original response, clients that need the original result (e.g., the created order's ID) can't get it from a retry.
- **Treating idempotency keys as a substitute for authentication/authorization** — a key only protects against *duplicate processing*; it says nothing about whether the request is allowed at all. Both checks are needed, independently.
