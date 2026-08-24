# Chapter 91: Saga Pattern — Distributed Transactions

Chapter 90 introduced the saga in one page: a sequence of local transactions, each with a compensating action that undoes it. That version had a hidden weakness — the saga lived only in memory, so a crash halfway through left the system in limbo. This chapter builds sagas properly: why two-phase commit is not the answer, how to classify saga steps (compensatable, pivot, retriable), sagas as persisted state machines that survive crashes, failure handling in depth, and a complete runnable order/payment/inventory example.

## Table of Contents

1. [Why Not Two-Phase Commit](#1-why-not-two-phase-commit)
2. [Anatomy of a Saga — Compensatable, Pivot, Retriable](#2-anatomy-of-a-saga--compensatable-pivot-retriable)
3. [Choreography in Depth](#3-choreography-in-depth)
4. [Orchestration in Depth](#4-orchestration-in-depth)
5. [The Saga State Machine — Surviving Crashes](#5-the-saga-state-machine--surviving-crashes)
6. [Failure Handling](#6-failure-handling)
7. [Worked Example — Order, Payment, Inventory](#7-worked-example--order-payment-inventory)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Not Two-Phase Commit

If the order, payment, and inventory data lived in one PostgreSQL database, you would wrap everything in one ACID transaction (Chapter 78) and be done. In a microservice world each service owns its own database, and the classic answer — **two-phase commit (2PC)** — asks a coordinator to first get every participant to *prepare* (lock resources, promise to commit), then tell all of them to *commit*.

```
2PC:                                      Why it fails in practice:
  Coordinator ──prepare──► Payments        - every participant holds LOCKS
  Coordinator ──prepare──► Inventory         while waiting for the slowest one
  Coordinator ──prepare──► Orders          - coordinator crash = everyone
  (all voted yes)                            stuck holding locks (blocking)
  Coordinator ──commit───► everyone        - Kafka, Redis, most cloud DBs
                                             simply don't support it
```

2PC trades availability for consistency and does it badly: throughput collapses under lock contention, and one slow or dead node freezes everyone. Almost no modern message broker or managed database even implements the protocol.

The saga takes the opposite trade: **each step commits immediately** (locally, with full ACID inside that one service), and if a later step fails, previously committed steps are *semantically undone* by compensating transactions. You give up atomic isolation across services and gain availability and loose coupling.

A saga guarantees **ACD** — Atomicity (eventually: all steps or all compensations), Consistency, Durability — but **not Isolation**: other transactions can observe the saga's intermediate states. §6 shows how to live with that.

---

## 2. Anatomy of a Saga — Compensatable, Pivot, Retriable

Not every step in a saga is equal. Classifying them is the single most useful design exercise:

| Kind | Definition | Example |
|------|-----------|---------|
| **Compensatable** | Can be undone by a compensating transaction | Reserve inventory → release it |
| **Pivot** | The go/no-go point. Once it succeeds, the saga must complete — it is never compensated | Charging the customer's card |
| **Retriable** | Comes after the pivot; cannot fail permanently, only be retried until it succeeds | Create shipment, send email |

```
  ReserveInventory ──► ChargePayment ──► CreateShipment ──► NotifyCustomer
  [compensatable]        [PIVOT]          [retriable]        [retriable]
        │
        ▼ on failure of a later compensatable/pivot step:
  ReleaseInventory  (compensations run in REVERSE order)
```

Rules that fall out of this classification:

1. **Order the steps so risky, compensatable steps come first** and hard-to-undo steps come last. Reserving stock is easy to undo; refunding a charged card is slow, costs fees, and annoys customers — so charge as late as possible.
2. **Everything after the pivot must be retriable.** If `CreateShipment` could fail permanently after the customer paid, you have designed a saga that can strand money. Fix the design, not the code.
3. **Compensations must be idempotent and must not fail on business grounds.** `ReleaseInventory` can hit a network error (retry it) but must never respond "no, I refuse to release."

Compensation is *semantic* undo, not a database rollback. Cancelling a sent email is impossible — the compensation for "order confirmation sent" is "send a cancellation email." The forward action and its compensation together must leave the system in an acceptable state, not necessarily the exact prior state.

---

## 3. Choreography in Depth

In choreography there is no coordinator: each service listens for events, does its local transaction, and publishes the next event. Chapter 90 sketched the happy path; the full picture includes the compensation flow:

```
                 OrderPlaced
                      │
                      ▼
              InventoryService ── stock ok ──► InventoryReserved
                      │                              │
                stock missing                        ▼
                      │                       PaymentService ── ok ──► PaymentCompleted
                      ▼                              │                       │
             InventoryRejected                 card declined                 ▼
                      │                              │                 OrderService
                      ▼                              ▼                 (confirm order)
                OrderService                  PaymentFailed
              (cancel order)                         │
                                                     ▼
                                            InventoryService
                                            (release reservation)
                                                     │
                                                     ▼
                                               OrderService
                                              (cancel order)
```

Each arrow is an event on a broker (Kafka, RabbitMQ, or Redis Streams from Chapter 97), and each service publishes via its outbox (Chapter 90) so the local write and the event are atomic:

```go
// InventoryService: forward step + compensation, both event handlers
func (s *InventoryService) OnOrderPlaced(ctx context.Context, e OrderPlacedEvent) error {
    return s.uow.Run(ctx, func(tx Tx) error {
        ok, err := s.reserveTx(ctx, tx, e.OrderID, e.Items)
        if err != nil { return err } // infra error → retry via redelivery

        if !ok { // business failure → compensating event, NOT an error
            return s.outbox.WriteTx(ctx, tx, InventoryRejectedEvent{
                OrderID: e.OrderID, Reason: "insufficient stock",
            })
        }
        return s.outbox.WriteTx(ctx, tx, InventoryReservedEvent{OrderID: e.OrderID})
    })
}

// Compensation: triggered by a failure LATER in the saga
func (s *InventoryService) OnPaymentFailed(ctx context.Context, e PaymentFailedEvent) error {
    return s.uow.Run(ctx, func(tx Tx) error {
        // Idempotent: releasing an already-released reservation is a no-op
        if err := s.releaseTx(ctx, tx, e.OrderID); err != nil { return err }
        return s.outbox.WriteTx(ctx, tx, InventoryReleasedEvent{OrderID: e.OrderID})
    })
}
```

Note the crucial distinction in error handling: an **infrastructure error** (DB down) returns an error so the broker redelivers; a **business failure** (out of stock) publishes a compensating *event* and returns nil — redelivering "out of stock" forever would help nobody.

Choreography shines for short sagas (2–3 steps). Beyond that it degrades: nobody can answer "where is order 123 in the saga?" without mentally simulating five services, and adding a step means touching several codebases. That is when you reach for orchestration.

---

## 4. Orchestration in Depth

An orchestrator owns the saga definition — the steps, their order, and their compensations — and drives participants through commands. Participants stay dumb: they execute a command and reply.

Instead of chapter 90's hand-rolled if/else chain, define the saga declaratively as data:

```go
// SagaStep pairs a forward action with its compensation.
type SagaStep struct {
    Name       string
    Action     func(ctx context.Context, s *SagaData) error
    Compensate func(ctx context.Context, s *SagaData) error // nil = nothing to undo
}

type SagaData struct {
    OrderID    string
    CustomerID string
    Items      []OrderItem
    Total      int64
    PaymentRef string // filled in by the payment step, needed by its compensation
}

func NewOrderSagaSteps(inv Inventory, pay Payments, ship Shipping) []SagaStep {
    return []SagaStep{
        {
            Name:       "reserve-inventory",
            Action:     func(ctx context.Context, s *SagaData) error { return inv.Reserve(ctx, s.OrderID, s.Items) },
            Compensate: func(ctx context.Context, s *SagaData) error { return inv.Release(ctx, s.OrderID) },
        },
        {
            Name: "charge-payment", // PIVOT
            Action: func(ctx context.Context, s *SagaData) error {
                ref, err := pay.Charge(ctx, s.CustomerID, s.Total, s.OrderID /* idempotency key */)
                if err != nil { return err }
                s.PaymentRef = ref
                return nil
            },
            Compensate: func(ctx context.Context, s *SagaData) error { return pay.Refund(ctx, s.PaymentRef) },
        },
        {
            Name:       "create-shipment", // retriable — after the pivot
            Action:     func(ctx context.Context, s *SagaData) error { return ship.Create(ctx, s.OrderID) },
            Compensate: nil,
        },
    }
}
```

The generic executor walks forward; on failure it walks the *completed* steps backward:

```go
func RunSaga(ctx context.Context, steps []SagaStep, data *SagaData) error {
    for i, step := range steps {
        if err := step.Action(ctx, data); err != nil {
            compErr := compensate(ctx, steps[:i], data) // only steps that completed
            if compErr != nil {
                return fmt.Errorf("step %q failed (%w); compensation ALSO failed: %v — manual intervention required",
                    step.Name, err, compErr)
            }
            return fmt.Errorf("saga aborted at step %q: %w", step.Name, err)
        }
    }
    return nil
}

func compensate(ctx context.Context, completed []SagaStep, data *SagaData) error {
    var errs []error
    for i := len(completed) - 1; i >= 0; i-- { // reverse order
        step := completed[i]
        if step.Compensate == nil { continue }
        if err := step.Compensate(ctx, data); err != nil {
            errs = append(errs, fmt.Errorf("compensate %q: %w", step.Name, err))
        }
    }
    return errors.Join(errs...)
}
```

This is clean and testable — but it still evaporates if the process dies between steps. Fixing that requires persistence.

---

## 5. The Saga State Machine — Surviving Crashes

A production saga is a **state machine whose current state is stored in a database**. Before executing a step, the orchestrator records where it is; after the step, it records the outcome. A crashed orchestrator (or any replica) picks up unfinished sagas and resumes.

```sql
CREATE TABLE sagas (
    saga_id     UUID PRIMARY KEY,
    saga_type   TEXT NOT NULL,              -- "order-fulfillment"
    data        JSONB NOT NULL,             -- SagaData
    step        INT  NOT NULL DEFAULT 0,    -- index of the step being executed
    direction   TEXT NOT NULL DEFAULT 'forward',  -- forward | compensating
    status      TEXT NOT NULL DEFAULT 'running',  -- running | completed | compensated | failed
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```
                        ┌────────────────────────────────────────┐
                        ▼                                        │
  running/forward ──step ok──► step++ ──last step?──► completed  │
        │                                     no ────────────────┘
     step fails
        │
        ▼                       ┌───────────────────────────────────┐
  running/compensating ──comp ok──► step-- ──first step?──► compensated
        │                                          no ──────────────┘
   comp fails (after retries)
        ▼
      failed  (alert a human)
```

The resumable executor persists the cursor around every transition:

```go
func (o *Orchestrator) Resume(ctx context.Context, sagaID uuid.UUID) error {
    saga, err := o.repo.Lock(ctx, sagaID) // SELECT ... FOR UPDATE SKIP LOCKED
    if err != nil { return err }

    steps := o.definitions[saga.Type]

    for saga.Status == "running" {
        if saga.Direction == "forward" {
            err := steps[saga.Step].Action(ctx, &saga.Data)
            switch {
            case err == nil && saga.Step == len(steps)-1:
                saga.Status = "completed"
            case err == nil:
                saga.Step++
            default:
                saga.Direction = "compensating" // stay on this step; it did not complete
                saga.Step--                     // compensate from the previous completed step
            }
        } else { // compensating
            if saga.Step < 0 {
                saga.Status = "compensated"
            } else {
                if c := steps[saga.Step].Compensate; c != nil {
                    if err := c(ctx, &saga.Data); err != nil {
                        return o.repo.Save(ctx, saga) // leave running; retry later with backoff
                    }
                }
                saga.Step--
            }
        }
        if err := o.repo.Save(ctx, saga); err != nil { return err } // checkpoint EVERY transition
    }
    return nil
}
```

A background sweeper makes crash recovery automatic:

```go
// Any replica can adopt sagas whose owner died mid-flight.
func (o *Orchestrator) Sweep(ctx context.Context) {
    for ids := range o.repo.PollStuck(ctx, 30*time.Second) { // running & not updated recently
        for _, id := range ids {
            go o.Resume(ctx, id) // SKIP LOCKED prevents two replicas grabbing one saga
        }
    }
}
```

One subtlety: the crash might have happened *after* a step's side effect but *before* the checkpoint. On resume, that step runs again — so **every step and every compensation must be idempotent**. The standard trick is an idempotency key derived from `(sagaID, stepName)`, exactly the pattern the payment step used above (`s.OrderID` as the charge's idempotency key). This is the same discipline Chapter 102 applies to message consumers.

---

## 6. Failure Handling

Sagas have three distinct failure classes, each with its own response:

| Failure | Example | Response |
|---------|---------|----------|
| **Business rejection** | Card declined, out of stock | Compensate — this is the saga working as designed |
| **Transient infra error** | Timeout, 503, deadlock | Retry the same step with exponential backoff; do not compensate |
| **Compensation failure** | Refund API down for an hour | Retry forever with backoff; after a threshold, park the saga as `failed` and page a human |

```go
func withRetry(fn func(ctx context.Context) error, attempts int) func(context.Context) error {
    return func(ctx context.Context) error {
        var err error
        for i := range attempts {
            if err = fn(ctx); err == nil { return nil }
            var b BusinessError
            if errors.As(err, &b) { return err } // business rejection: compensate, don't retry
            select {
            case <-time.After(time.Duration(1<<i) * 200 * time.Millisecond):
            case <-ctx.Done():
                return ctx.Err()
            }
        }
        return err
    }
}
```

Distinguishing the first two classes is worth real design effort: retrying a declined card is useless (and may get you blocked by the payment provider); compensating because of a 2-second network blip cancels orders that would have succeeded.

### Living without isolation

Because sagas commit each step immediately, other requests see intermediate states: an order that is "placed" but whose payment later fails. Standard countermeasures:

- **Semantic lock**: mark the record with a pending state (`status = 'PAYMENT_PENDING'`). Other transactions treat pending rows as untouchable. The final step (or compensation) clears the flag.
- **Commutative updates**: design updates so order does not matter — `balance = balance + 100` commutes; `balance = 500` does not.
- **Reread the value**: before the pivot step, re-verify critical data (price, stock) rather than trusting values read at saga start.

### Timeouts

A saga stuck waiting on a reply forever is a resource leak and a semantic lock held indefinitely. Give each step a deadline; on expiry treat it as a transient failure (retry) and eventually as a business failure (compensate):

```go
stepCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
err := step.Action(stepCtx, data)
```

Careful: a timed-out charge may still have succeeded on the provider's side. The compensation (refund by idempotency key) must handle "charge exists" and "charge never happened" equally gracefully — which is exactly why compensations key off the same idempotency key as the forward action.

---

## 7. Worked Example — Order, Payment, Inventory

A complete, runnable orchestration saga. Services are in-memory with injectable failures so you can watch compensation happen. Save as `main.go` and `go run` it (Go 1.22+):

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ---------- Domain errors ----------

type BusinessError struct{ Reason string }

func (e BusinessError) Error() string { return e.Reason }

// ---------- Participants ----------

type Inventory struct {
	mu       sync.Mutex
	stock    map[string]int
	reserved map[string]map[string]int // orderID -> sku -> qty
}

func NewInventory(stock map[string]int) *Inventory {
	return &Inventory{stock: stock, reserved: map[string]map[string]int{}}
}

func (inv *Inventory) Reserve(_ context.Context, orderID string, items map[string]int) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	if _, done := inv.reserved[orderID]; done { return nil } // idempotent
	for sku, qty := range items {
		if inv.stock[sku] < qty {
			return BusinessError{Reason: fmt.Sprintf("insufficient stock for %s", sku)}
		}
	}
	for sku, qty := range items { inv.stock[sku] -= qty }
	inv.reserved[orderID] = items
	return nil
}

func (inv *Inventory) Release(_ context.Context, orderID string) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	items, ok := inv.reserved[orderID]
	if !ok { return nil } // idempotent: never reserved, or already released
	for sku, qty := range items { inv.stock[sku] += qty }
	delete(inv.reserved, orderID)
	return nil
}

type Payments struct {
	mu      sync.Mutex
	charges map[string]int64 // idempotency key -> amount
	Decline bool             // test hook: simulate card declined
}

func NewPayments() *Payments { return &Payments{charges: map[string]int64{}} }

func (p *Payments) Charge(_ context.Context, key string, amount int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, done := p.charges[key]; done { return nil } // idempotent
	if p.Decline {
		return BusinessError{Reason: "card declined"}
	}
	p.charges[key] = amount
	return nil
}

func (p *Payments) Refund(_ context.Context, key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.charges, key) // refunding a non-existent charge is a no-op
	return nil
}

type Shipping struct{ created sync.Map }

func (s *Shipping) Create(_ context.Context, orderID string) error {
	s.created.Store(orderID, true)
	return nil
}

// ---------- Saga machinery ----------

type SagaData struct {
	OrderID string
	Items   map[string]int
	Total   int64
}

type SagaStep struct {
	Name       string
	Action     func(context.Context, *SagaData) error
	Compensate func(context.Context, *SagaData) error
}

func RunSaga(ctx context.Context, steps []SagaStep, data *SagaData) error {
	for i, step := range steps {
		fmt.Printf("  → %s\n", step.Name)
		if err := step.Action(ctx, data); err != nil {
			fmt.Printf("  ✗ %s failed: %v — compensating\n", step.Name, err)
			for j := i - 1; j >= 0; j-- {
				if steps[j].Compensate == nil { continue }
				fmt.Printf("  ↩ undo %s\n", steps[j].Name)
				if cerr := steps[j].Compensate(ctx, data); cerr != nil {
					return fmt.Errorf("compensation %q failed: %v (original: %w)", steps[j].Name, cerr, err)
				}
			}
			return fmt.Errorf("saga aborted at %q: %w", step.Name, err)
		}
	}
	return nil
}

// ---------- Wiring ----------

func orderSaga(inv *Inventory, pay *Payments, ship *Shipping) []SagaStep {
	return []SagaStep{
		{
			Name:       "reserve-inventory",
			Action:     func(ctx context.Context, s *SagaData) error { return inv.Reserve(ctx, s.OrderID, s.Items) },
			Compensate: func(ctx context.Context, s *SagaData) error { return inv.Release(ctx, s.OrderID) },
		},
		{
			Name:       "charge-payment", // pivot
			Action:     func(ctx context.Context, s *SagaData) error { return pay.Charge(ctx, "order:"+s.OrderID, s.Total) },
			Compensate: func(ctx context.Context, s *SagaData) error { return pay.Refund(ctx, "order:"+s.OrderID) },
		},
		{
			Name:   "create-shipment", // retriable
			Action: func(ctx context.Context, s *SagaData) error { return ship.Create(ctx, s.OrderID) },
		},
	}
}

func main() {
	ctx := context.Background()
	inv := NewInventory(map[string]int{"widget": 5})
	pay := NewPayments()
	ship := &Shipping{}
	steps := orderSaga(inv, pay, ship)

	// Happy path
	fmt.Println("order-1 (happy path):")
	err := RunSaga(ctx, steps, &SagaData{OrderID: "order-1", Items: map[string]int{"widget": 2}, Total: 2000})
	fmt.Println("  result:", errOrOK(err), "| stock:", inv.stock["widget"]) // ok | stock: 3

	// Payment fails → inventory must be released
	fmt.Println("order-2 (card declined):")
	pay.Decline = true
	err = RunSaga(ctx, steps, &SagaData{OrderID: "order-2", Items: map[string]int{"widget": 2}, Total: 2000})
	fmt.Println("  result:", errOrOK(err), "| stock:", inv.stock["widget"]) // aborted | stock: 3 (released!)
	pay.Decline = false

	// Business rejection at step 1 → nothing to compensate
	fmt.Println("order-3 (not enough stock):")
	err = RunSaga(ctx, steps, &SagaData{OrderID: "order-3", Items: map[string]int{"widget": 99}, Total: 99000})
	fmt.Println("  result:", errOrOK(err), "| stock:", inv.stock["widget"]) // aborted | stock: 3

	var b BusinessError
	fmt.Println("  was a business rejection:", errors.As(err, &b))
}

func errOrOK(err error) string {
	if err != nil { return err.Error() }
	return "ok"
}
```

Run it and trace the output: order-2 reserves stock (5→3... then 3→1), the charge is declined, and compensation returns the stock to 3. The invariant to internalize — **stock lost = stock in successfully completed orders, always** — holds no matter where the saga fails, because every completed step has an idempotent undo.

To turn this into the real thing: persist `SagaData` + step cursor with the `sagas` table from §5, replace direct method calls with commands/replies over a broker (Chapters 98–101), and publish every state change through the outbox from Chapter 90. Chapter 106's ticket booking project does exactly that.

---

## Summary

- 2PC holds locks across services and blocks on coordinator failure — sagas commit each step locally and undo with **compensating transactions** instead
- Classify steps: **compensatable** (undoable, put them first), **pivot** (go/no-go, never compensated), **retriable** (after the pivot, must eventually succeed)
- Compensation is *semantic* undo, runs in reverse order, and must be idempotent and unable to fail on business grounds
- **Choreography**: services react to events, compensations are events too; great for 2–3 steps, opaque beyond that
- **Orchestration**: saga defined as data (`[]SagaStep`), driven by an executor; persist the step cursor in a `sagas` table so any replica can resume after a crash
- Steps re-run after crashes, so every action and compensation needs an **idempotency key** (typically `sagaID + stepName`)
- Distinguish business rejections (compensate immediately) from transient errors (retry with backoff) from compensation failures (retry, then page a human)
- No isolation between steps — use semantic locks, commutative updates, and pre-pivot rereads

## Exercises

### Easy
1. Add a `send-confirmation` retriable step to the worked example that just prints a message. Verify it never runs when payment is declined.
2. Make the `Shipping` service fail once via a test hook, and wrap its `Action` with the `withRetry` helper from §6. Verify the saga completes on the second attempt without compensating anything.
3. Call `RunSaga` twice with the same `OrderID` on the happy path. Using the idempotency guards already in the services, verify stock is only decremented once and the customer is only charged once.

### Medium
4. Add a `sagas` table (or an in-memory `map[string]SagaRecord`) and refactor `RunSaga` into the resumable state machine from §5: persist `{step, direction, status}` after every transition. Kill the process (return early) after `charge-payment`, "restart," call `Resume`, and verify the shipment is still created.
5. Implement a **semantic lock**: give orders a status column (`PENDING → CONFIRMED / CANCELLED`). A `GetOrder` query must report pending orders as "processing." Verify that a concurrently running query never sees a half-finished order as confirmed.
6. Rebuild the worked example as a **choreography saga**: three goroutines (order, inventory, payment) communicating via channels carrying events (`OrderPlaced`, `InventoryReserved`, `InventoryRejected`, `PaymentCompleted`, `PaymentFailed`, `InventoryReleased`). Verify the compensation flow triggers on payment failure with no central coordinator.

### Hard
7. Implement the full crash-safe orchestrator against PostgreSQL: `sagas` table, `Resume` with `SELECT ... FOR UPDATE SKIP LOCKED`, and a `Sweep` goroutine that adopts sagas idle for >10 seconds. Run two orchestrator replicas, kill one mid-saga (randomly), and prove every saga ends `completed` or `compensated` — never stuck — under 100 concurrent orders.
8. Build a **saga observability layer**: record every step transition to a `saga_events` table with timestamps, expose `GET /sagas/{id}` returning the full timeline, and add Prometheus metrics (`saga_completed_total`, `saga_compensated_total`, `saga_step_duration_seconds`). Then introduce a 10% random failure rate in payments and use your metrics to compute the compensation rate.
