# Chapter 109: Service Communication — REST vs gRPC vs Events

The previous chapter gave you the microservices toolbox and a quick tour of the communication options. You already know how to *build* each one: REST APIs with chi (Chapter 60), gRPC services with protobuf (Chapter 66), and event-driven flows with Kafka, RabbitMQ, and NATS (Volume 8). This chapter answers the harder question: **which one do you pick for each edge in your service graph, and why?** Choosing wrong here is expensive — a synchronous call chain that should have been an event stream will wake you up at 3 AM, and an event pipeline that should have been a simple RPC will bury a one-line answer under three brokers and a consumer group.

## Table of Contents

1. [The Three Styles at a Glance](#1-the-three-styles-at-a-glance)
2. [Synchronous vs Asynchronous — What You Are Really Choosing](#2-synchronous-vs-asynchronous--what-you-are-really-choosing)
3. [Coupling: Temporal, Behavioral, and Schema](#3-coupling-temporal-behavioral-and-schema)
4. [Latency Math — Why Chains Kill You](#4-latency-math--why-chains-kill-you)
5. [Contracts and Schema Evolution](#5-contracts-and-schema-evolution)
6. [The Decision Framework](#6-the-decision-framework)
7. [Hybrid Architectures — How Real Systems Mix All Three](#7-hybrid-architectures--how-real-systems-mix-all-three)
8. [Sync Facades over Async Backends](#8-sync-facades-over-async-backends)
9. [Anti-Patterns](#9-anti-patterns)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Three Styles at a Glance

```
REST (HTTP/1.1 + JSON)                gRPC (HTTP/2 + protobuf)
[Caller] --request--> [Callee]        [Caller] --request--> [Callee]
[Caller] <-response-- [Callee]        [Caller] <-response-- [Callee]
   human-readable, universal             binary, typed, fast, streaming

Events (broker in the middle)
[Producer] --publish--> [Broker] --deliver--> [Consumer A]
                            └-----deliver--> [Consumer B]
   producer does not know or care who consumes, or when
```

| Dimension | REST | gRPC | Events |
|---|---|---|---|
| **Style** | Sync request-response | Sync request-response (+ streaming) | Async publish-subscribe |
| **Payload** | JSON (text) | Protobuf (binary) | Your choice (JSON, protobuf, Avro) |
| **Contract** | OpenAPI (optional, often stale) | `.proto` file (enforced by codegen) | Event schema (enforced only if you make it so) |
| **Typical latency** | 1–10 ms | 0.5–5 ms | 5 ms–seconds (broker hop + consumer lag) |
| **Caller gets an answer?** | Yes, immediately | Yes, immediately | No — fire and forget |
| **Callee must be up?** | Yes | Yes | No — broker buffers |
| **Browser/external clients** | Perfect | Needs gRPC-Web or a gateway | Not directly |
| **Debugging** | curl and read | grpcurl, needs proto files | Trace through broker; harder |
| **Backpressure** | Caller blocks, timeouts fire | Same, plus HTTP/2 flow control | Consumer lag absorbs spikes |

None of these is "best." They occupy different points on a trade-off curve, and a production system of any size uses all three.

---

## 2. Synchronous vs Asynchronous — What You Are Really Choosing

The REST-vs-gRPC choice is a detail. The **sync-vs-async** choice is architectural, and it comes down to one question:

> **Does the caller need the answer to finish its own work?**

- Checkout needs to know **now** whether the card was charged → synchronous.
- Checkout does not need to know whether the confirmation email sent → asynchronous.

### What sync buys you (and costs you)

```go
// Synchronous: the answer is right there. Simple to write, simple to reason about.
func (s *CheckoutService) Checkout(ctx context.Context, cart Cart) (Receipt, error) {
	price, err := s.pricing.Quote(ctx, cart.Items) // blocks until pricing answers
	if err != nil {
		return Receipt{}, fmt.Errorf("quote: %w", err) // and fails if pricing is down
	}
	charge, err := s.payments.Charge(ctx, cart.CustomerID, price.Total)
	if err != nil {
		return Receipt{}, fmt.Errorf("charge: %w", err)
	}
	return Receipt{ChargeID: charge.ID, Total: price.Total}, nil
}
```

You get: an immediate answer, straight-line code, errors at the call site.
You pay: **availability coupling**. If pricing is down, checkout is down. Your availability is the *product* of every service in the chain: three services at 99.9% each give the chain 99.7% — you lose a nine every few hops.

### What async buys you (and costs you)

```go
// Asynchronous: publish and move on. The email service can be down for an hour;
// the message waits in the broker and nobody's checkout fails.
func (s *CheckoutService) afterCharge(ctx context.Context, r Receipt) error {
	return s.publisher.Publish(ctx, "order.completed", OrderCompleted{
		OrderID:    r.OrderID,
		CustomerID: r.CustomerID,
		Total:      r.Total,
		OccurredAt: time.Now().UTC(),
	})
}
```

You get: fault isolation, spike absorption, easy fan-out (analytics and email and loyalty points all consume the same event without the producer changing).
You pay: **no answer**, eventual consistency, duplicate delivery (Chapter 102), harder debugging, and a broker to operate.

Chapter 95 covered this trade-off for background work inside one system. Here the same logic decides the *shape of your service graph*: every synchronous edge is an availability and latency dependency; every asynchronous edge is a consistency and observability cost.

---

## 3. Coupling: Temporal, Behavioral, and Schema

"Loose coupling" is thrown around loosely. It splits into three specific kinds, and each communication style scores differently:

| Coupling type | Question it asks | REST/gRPC | Events |
|---|---|---|---|
| **Temporal** | Must both services be up at the same moment? | Yes — tight | No — broker decouples |
| **Behavioral** | Does the caller know *who* does the work? | Yes — caller names the callee | No — producer names the *fact*, not the audience |
| **Schema** | Do both sides share a data contract? | Yes | Yes — **events do not remove schema coupling** |

Two of these deserve emphasis:

**Behavioral coupling is about direction of knowledge.** `orders` calling `POST /emails/send` means orders knows email exists and commands it. `orders` publishing `order.completed` means orders states a fact; email, analytics, and fraud each decide for themselves to react. Adding a fourth consumer requires **zero changes** to orders. That is the real superpower of events — not speed, not scale: *the producer's code stops changing when the audience grows*.

**Schema coupling never goes away.** A consumer that reads `event.Total` breaks just as hard as an RPC client when you rename the field. Worse: with RPC the break happens at call time and the caller sees the error; with events the break happens in some consumer, minutes later, in another team's logs. Section 5 deals with this.

---

## 4. Latency Math — Why Chains Kill You

Synchronous calls compose *serially* by default. Chapter 108 warned against chains longer than 2–3 hops; here is the arithmetic behind the warning.

```
A --> B --> C --> D        (each hop: p50 = 5ms, p99 = 50ms)

Best case (all p50):   5 + 5 + 5  = 15 ms
Worst case (all p99): 50 + 50 + 50 = 150 ms
```

And it is worse than that: with three independent hops, the chance that *at least one* of them hits its slow tail is roughly `1 - 0.99³ ≈ 3%`. Your chain's p97 is already the sum of tail latencies. **Tail latency compounds across serial hops.** This is why deep synchronous graphs feel fine in dev and melt in production.

Two mitigations, both of which you can do in Go today:

### Parallelize independent calls

```go
import "golang.org/x/sync/errgroup"

// pricing and inventory don't depend on each other — call them concurrently.
func (s *CheckoutService) prepare(ctx context.Context, cart Cart) (Quote, Stock, error) {
	var (
		quote Quote
		stock Stock
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		quote, err = s.pricing.Quote(ctx, cart.Items)
		return err
	})
	g.Go(func() error {
		var err error
		stock, err = s.inventory.Reserve(ctx, cart.Items)
		return err
	})
	if err := g.Wait(); err != nil {
		return Quote{}, Stock{}, err
	}
	return quote, stock, nil
}
```

Latency is now `max(hop latencies)` instead of `sum(hop latencies)`.

### Replace read chains with local data, fed by events

If `orders` calls `users` on every request just to display a name, stop. Have `orders` consume `user.updated` events and keep a small local read model (a table with `user_id, display_name`). The read becomes a local DB hit; the users service can be down all afternoon and order pages still render. You met this idea as CQRS in Chapter 88 — event-carried state transfer is its inter-service form. (When that local copy becomes a cache with real memory pressure, Chapter 113 shows how to distribute it.)

---

## 5. Contracts and Schema Evolution

Every inter-service message — request, response, or event — is a contract that will outlive the code that first wrote it. The question is not "will the schema change?" but "who breaks when it does?"

### The one rule that governs everything

> **You can never deploy both sides atomically.** Producer and consumer deploy at different times, so every schema change must be compatible with the *other* side's old version, in both directions, for at least one release cycle.

That gives you the standard playbook, the **expand–migrate–contract** pattern:

```
1. EXPAND:   add the new field (optional, with a default). Deploy producers.
2. MIGRATE:  update all consumers to read the new field. Deploy them.
3. CONTRACT: remove the old field. Only now. Often "never".
```

### Protobuf: evolution is a first-class feature

```protobuf
message Order {
    string id          = 1;
    string customer_id = 2;
    // string coupon   = 3;  // removed — reserve the number so it's never reused!
    reserved 3;
    reserved "coupon";
    double total       = 4;
    OrderStatus status = 5;   // enums: always keep 0 as UNSPECIFIED
    string currency    = 6;   // added later — old readers simply ignore it
}
```

Safe in protobuf:
- **Adding fields** — old readers skip unknown tags, new readers see zero values from old writers.
- **Removing fields** — if you `reserved` the tag number and name.
- **Renaming fields** — the wire format only uses tag numbers (but codegen'd names change, so it still breaks *code*, just not the wire).

Breaking in protobuf:
- **Reusing a tag number** with a different type — silent data corruption, the worst failure mode in this chapter.
- **Changing a field's type** (with a few narrow exceptions).
- Making semantic changes ("`total` is now in cents, not dollars") — no serialization format saves you from these. Add a new field (`total_cents`) instead.

### JSON/REST: discipline replaces tooling

JSON has no tag numbers and no compiler checking. The same rules apply, enforced by convention:

```go
// Tolerant reader: ignore unknown fields (encoding/json does this by default —
// do NOT enable DisallowUnknownFields on inter-service payloads),
// and give removed/missing fields safe zero values or pointers.
type OrderResponse struct {
	ID         string  `json:"id"`
	Total      float64 `json:"total"`
	// New optional field: pointer distinguishes "absent" from "zero".
	Currency   *string `json:"currency,omitempty"`
}
```

- **Version the API surface, not every field**: `/v1/orders` → `/v2/orders` when you must break (Chapter 60 covered URL versioning).
- Publish an OpenAPI spec and generate types (Chapter 63) so "the contract" is a file in a repo, not tribal knowledge.

### Events: the hardest case

An RPC response has exactly one reader: the caller, right now. An event may have **ten consumers you have never met, replaying messages from last month** (Kafka retention, Chapter 98). So event schemas need the most care:

1. **Envelope with a version field**, so consumers can branch:

```go
type Envelope struct {
	Type       string          `json:"type"`        // "order.completed"
	Version    int             `json:"version"`     // 2
	OccurredAt time.Time       `json:"occurred_at"`
	Payload    json.RawMessage `json:"payload"`     // decode after checking Version
}
```

2. **Only additive changes within a version.** A breaking change means a new event version published *alongside* the old one (dual-publish) until every consumer has migrated.
3. **A schema registry** (Confluent Schema Registry, or a `.proto`/JSON-Schema directory in a shared repo with CI checks) so incompatible producers fail at build time, not in a consumer at runtime.

---

## 6. The Decision Framework

Work through these questions **in order** for each edge in your service graph:

```
Q1. Does the caller need the result to complete its own request?
    NO  ──────────────────────────────► EVENTS
    YES ▼

Q2. Is the caller a browser, mobile app, or third party?
    YES ──────────────────────────────► REST (universal tooling, no codegen burden on clients)
    NO  ▼   (internal service-to-service)

Q3. Is the call on a hot path (high QPS, latency-sensitive),
    or does it need streaming, or strong compile-time contracts?
    YES ──────────────────────────────► gRPC
    NO  ──────────────────────────────► REST or gRPC — pick ONE convention
                                        for internal calls and stay consistent
```

Then sanity-check the answer against the trade-off table:

| If you need… | Choose | Because |
|---|---|---|
| Immediate answer, external client | REST | curl-able, cacheable, no toolchain demands |
| Immediate answer, internal, hot path | gRPC | binary payload, HTTP/2 multiplexing, typed contract |
| Server push / long-lived streams internally | gRPC streaming | built-in; REST needs SSE/WebSocket bolt-ons |
| Fire-and-forget side effects | Events | temporal decoupling, retry via broker |
| Fan-out to N consumers | Events | producers don't change when consumers appear |
| Spike absorption / load leveling | Events | broker buffers; consumers drain at their pace |
| Replayable history / audit | Events (Kafka) | retained log, new consumers read from offset 0 |
| Strict request ordering per key | Events with partitioning (Ch 103) | partition key preserves per-key order |
| Cross-service workflow with rollback | Events + Saga (Ch 91) | no distributed transactions over RPC |

**Default heuristics that survive contact with production:**

- **Queries (reads) are usually sync; facts (writes that others react to) are usually async.** "What is the price?" is an RPC. "An order was placed" is an event.
- If you catch yourself publishing an event and then *waiting* for a reply event — you wanted an RPC.
- If you catch yourself calling another service in a loop, or calling it just to copy its data into your response — you wanted an event-fed local read model.

---

## 7. Hybrid Architectures — How Real Systems Mix All Three

Here is the checkout flow of a realistic e-commerce platform (this is exactly the shape you will build in Major Project 3, Chapter 117):

```
            REST (external)
[Mobile/Web] ────────► [API Gateway] ────────► [Order Service]
                                                  │        │
                                     gRPC (sync)  │        │  Kafka (async)
                              ┌───────────────────┤        └─────────────────┐
                              ▼                   ▼                          ▼
                        [Pricing Svc]       [Payment Svc]            topic: order.placed
                                                                  ┌───────┼──────────┐
                                                                  ▼       ▼          ▼
                                                             [Email]  [Analytics] [Inventory]
```

- **Edge: REST.** Browsers and mobile apps speak HTTP/JSON. The gateway (next two chapters) terminates TLS and authenticates.
- **Hot internal path: gRPC.** Order → Pricing and Order → Payment need answers *now*, thousands of times per second. Typed contracts and binary encoding earn their keep here.
- **Everything downstream of the decision: events.** Once the order is placed, email, analytics, and inventory are reactions, not prerequisites. Order service publishes one event; three teams consume it independently.

The subtle part is the seam between sync and async. The order handler must **atomically** save the order *and* guarantee the event gets published — a crash between the two leaves an order nobody ships. You already know the fix: the **Outbox pattern** from Chapter 90 (insert the event into an `outbox` table in the same DB transaction as the order; a relay publishes it to Kafka). In a microservices world the outbox is not optional — it is the standard bridge from a synchronous write to an asynchronous fact.

```go
// The seam: one transaction, two writes — order + outbox. The relay does the rest.
func (s *OrderService) PlaceOrder(ctx context.Context, cmd PlaceOrder) (Order, error) {
	// 1. Sync validations first: these NEED answers.
	quote, err := s.pricing.Quote(ctx, cmd.Items) // gRPC
	if err != nil {
		return Order{}, fmt.Errorf("pricing unavailable: %w", err)
	}
	if _, err := s.payments.Authorize(ctx, cmd.CustomerID, quote.Total); err != nil { // gRPC
		return Order{}, fmt.Errorf("payment declined: %w", err)
	}

	// 2. One transaction: persist the order AND enqueue the event.
	order := NewOrder(cmd, quote)
	err = s.db.WithTx(ctx, func(tx Tx) error {
		if err := tx.SaveOrder(ctx, order); err != nil {
			return err
		}
		return tx.SaveOutbox(ctx, NewEnvelope("order.placed", 1, order))
	})
	if err != nil {
		return Order{}, err
	}
	return order, nil // email/analytics/inventory happen later, via the relay → Kafka
}
```

---

## 8. Sync Facades over Async Backends

Sometimes the *user* needs an answer but the *work* is asynchronous — video encoding, report generation, payment settlement. The standard shapes, from simplest to fanciest:

**1. Accept + poll.** Return `202 Accepted` with a status URL:

```go
func (h *ReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewString()
	if err := h.publisher.Publish(r.Context(), "report.requested", ReportRequested{ID: id}); err != nil {
		http.Error(w, "try again later", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Location", "/v1/reports/"+id)
	w.WriteHeader(http.StatusAccepted) // 202: "working on it, check Location"
}

// GET /v1/reports/{id} → {"status":"pending"} … {"status":"ready","url":"..."}
```

**2. Push the result.** Same as above, but instead of polling, notify via WebSocket/SSE (Chapters 65, 139) or a webhook when the consumer finishes.

**3. Request–reply over the broker.** The caller publishes with a `reply_to` topic and a correlation ID, then blocks (with a deadline!) on the reply. NATS makes this a one-liner (`nc.Request`, Chapter 100). Use sparingly — you have rebuilt RPC with extra steps, and it is only worth it when everything else already flows through the broker.

The rule: pick the *weakest* mechanism that satisfies the user experience. Polling is unglamorous and bulletproof; start there.

---

## 9. Anti-Patterns

- **The distributed monolith.** Every service calls every other synchronously; nothing works unless everything works. You paid the microservices tax and got the monolith's coupling back as interest. Symptom: a deploy of any one service requires coordinated deploys of three others.
- **Event-as-command.** Publishing `send.welcome.email` to a topic consumed by exactly one service is not event-driven design — it is an RPC hiding from its error handling. Events state facts (`user.registered`); commands name a recipient and expect execution. If you need a command, make the call explicit (RPC or a task queue like asynq, Chapter 104).
- **The chatty read.** Calling `users`, `preferences`, and `subscriptions` on every request to render one profile. Fix with a composite endpoint, a BFF (Chapter 111), or an event-fed read model.
- **Dual-write without an outbox.** `db.Save(order)` then `kafka.Publish(event)` as two separate operations. The crash window between them *will* be hit. Chapter 90 exists for a reason.
- **Sharing a database instead of an API.** Two services reading each other's tables have the tightest schema coupling possible and none of it visible in any contract. Every migration becomes a cross-team incident.
- **gRPC to the browser without a plan.** Browsers cannot speak native gRPC. If your external clients are web apps, front gRPC services with a REST/JSON translation layer (gRPC-Gateway, Chapter 66 §7) or the gateway from Chapter 111.

---

## Summary

- The real decision is **sync vs async**: sync when the caller needs the answer to proceed, async when it does not. REST-vs-gRPC is a second-order choice within "sync."
- **Availability multiplies across sync hops** (three 99.9% services chain to 99.7%) and **tail latency compounds** — keep synchronous chains to 2–3 hops, parallelize independent calls with `errgroup`, and replace chatty reads with event-fed local data.
- Coupling splits into **temporal** (must both be up?), **behavioral** (who knows about whom?), and **schema** (shared contract). Events remove the first two; **nothing removes schema coupling**.
- Schema evolution follows **expand → migrate → contract**, because you can never deploy both sides atomically. Protobuf: add fields freely, `reserved` removed tags, never reuse tag numbers. JSON: tolerant readers, pointers for optional fields. Events: version envelopes, additive-only within a version, dual-publish across breaking changes.
- Default choices: **REST at the edge**, **gRPC for hot internal request-response**, **events for facts, fan-out, and everything the user doesn't wait for** — bridged by the **outbox pattern** at the sync/async seam.
- When users need answers about async work: `202 Accepted` + polling first, push notifications second, broker request-reply only when everything already lives on the broker.

Next up: whatever style you chose, the caller still has to *find* the callee. Chapter 110 covers service discovery and load balancing — how `order-service:50051` actually resolves to a healthy instance.

---

## Exercises

### Easy
1. For each edge, choose REST, gRPC, or events, and justify in one sentence: (a) mobile app → backend "get my order history"; (b) order service → fraud service "score this transaction before accepting it" at 5,000 QPS; (c) order service → warehouse, email, and analytics after an order ships; (d) admin dashboard → any service, used by 5 internal staff.
2. Take the protobuf `Order` message from Section 5 and make three changes: add a `notes` field, remove `status`, and rename `customer_id` to `buyer_id`. Write the resulting `.proto` so that all three changes are wire-compatible, and note which one still breaks compiled code.
3. A chain `A→B→C→D` has per-hop p50 = 4 ms and p99 = 60 ms. Compute the chain's best-case latency, worst-case latency, and the probability that at least one hop exceeds its p99. Then recompute the worst case if B, C, D are independent and called in parallel.

### Medium
4. Implement the **accept + poll** facade from Section 8 end to end: an HTTP handler that returns `202` with a `Location` header, a worker goroutine consuming from a channel (standing in for the broker) that sleeps 2 seconds then marks the job done in a shared in-memory store, and a `GET /reports/{id}` status endpoint. Test the full lifecycle with `httptest`.
5. Build a **versioned event envelope** consumer: define `Envelope` from Section 5, publish (to a channel or NATS) a mix of `order.completed` v1 (`{"total": 12.5}`) and v2 (`{"total_cents": 1250, "currency": "USD"}`) payloads, and write one consumer that decodes both versions into a single internal struct. Add a test that an unknown v3 is logged and skipped, not crashed on.
6. Demonstrate **schema breakage** empirically: write a JSON producer and consumer as separate programs sharing a struct, then rename a field only on the producer side. Show the consumer silently reading zero values. Fix it with the expand–migrate–contract sequence across three "deploys" (three commits), keeping the consumer working at every step.

### Hard
7. Build a mini **event-fed read model**: service A (users) exposes `PUT /users/{id}` and publishes `user.updated` events to NATS or a Redis stream; service B (orders) consumes them into an in-memory `map[string]string` of display names and serves `GET /orders/{id}` embedding the name without ever calling A. Kill service A and verify B still serves complete responses. Then restart A, update a user, and verify B converges.
8. Implement **request–reply over a broker** with correlation IDs: a client publishes `pricing.request` with a `ReplyTo` subject and a UUID, then waits (with a 2-second `context` deadline) for the matching reply; a pricing worker consumes requests and publishes replies. Handle three failure cases in tests: the worker is down (deadline fires), a stale reply arrives with a mismatched correlation ID (ignored), and two concurrent clients get their own answers (no cross-talk).
