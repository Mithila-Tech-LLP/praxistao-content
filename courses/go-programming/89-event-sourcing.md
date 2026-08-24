# Chapter 89: Event Sourcing — State as a Stream of Events

Chapter 88 introduced Event Sourcing as an idea: instead of storing current state, store every change as an immutable event and rebuild state by replaying them. This chapter is the practical deep dive — the parts you need before running it in production: a proper event store schema, the decide/apply split inside aggregates, snapshots so replay stays fast, versioning and upcasting so old events keep working with new code, and projections with checkpoints so read models can always be rebuilt. It ends with a complete runnable example you can paste into one file and run.

## Table of Contents

1. [From Concept to Production](#1-from-concept-to-production)
2. [Event Store Design](#2-event-store-design)
3. [The Aggregate — Decide and Apply](#3-the-aggregate--decide-and-apply)
4. [A Generic Event-Sourced Repository](#4-a-generic-event-sourced-repository)
5. [Snapshots](#5-snapshots)
6. [Event Versioning and Upcasting](#6-event-versioning-and-upcasting)
7. [Projections and Checkpoints](#7-projections-and-checkpoints)
8. [A Complete Runnable Example](#8-a-complete-runnable-example)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. From Concept to Production

Think of a bank statement. Your bank does not store just a single number called "balance" — it stores every deposit and withdrawal ever made. The balance is *derived* from that history. If the bank and you disagree, you can replay the statement line by line and settle the argument. That is event sourcing.

```
Traditional (state-oriented):        Event-sourced:

  accounts                             events (append-only)
  ┌────────┬─────────┐                 ┌──────────────────┬─────────┐
  │ id     │ balance │                 │ AccountOpened    │ acc-1   │
  ├────────┼─────────┤                 │ MoneyDeposited   │ +500    │
  │ acc-1  │ 350     │  ◄── UPDATE     │ MoneyWithdrawn   │ -200    │
  └────────┴─────────┘   overwrites    │ MoneyDeposited   │ +50     │
                          history      └──────────────────┴─────────┘
                                        balance = replay = 350
```

The chapter 88 sketch works, but four questions appear the moment you build something real:

1. **How do I store events so concurrent writers can't corrupt a stream?** → event store design (§2)
2. **Where does business logic live if state is just replayed events?** → decide/apply (§3)
3. **What happens when a stream has 100,000 events?** → snapshots (§5)
4. **What happens when the shape of an event changes next year?** → upcasting (§6)

Answer those four and you have a production event-sourced system.

---

## 2. Event Store Design

An event store needs surprisingly little: an append-only table with two orderings — per-stream order (`version`) and global order (`position`).

```sql
CREATE TABLE events (
    position    BIGSERIAL PRIMARY KEY,        -- global order (for projections)
    stream_id   TEXT        NOT NULL,          -- "account-42"
    version     INT         NOT NULL,          -- 1, 2, 3, ... within the stream
    event_type  TEXT        NOT NULL,          -- "MoneyDeposited"
    schema_ver  INT         NOT NULL DEFAULT 1,-- for upcasting (§6)
    data        JSONB       NOT NULL,
    metadata    JSONB       NOT NULL DEFAULT '{}', -- correlation ID, user ID, ...
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (stream_id, version)               -- THE concurrency guard
);

CREATE INDEX idx_events_stream ON events (stream_id, version);
```

The `UNIQUE (stream_id, version)` constraint is the heart of the design. Two writers that both loaded the account at version 7 will both try to append version 8. One insert succeeds; the other violates the constraint and gets rejected. That is **optimistic concurrency** enforced by the database itself — no locks held during business logic.

```go
var ErrConcurrency = errors.New("event store: stream was modified by another writer")

// Append writes events to a stream, expecting the stream to currently be
// at expectedVersion. Events get versions expectedVersion+1, +2, ...
func (s *PostgresEventStore) Append(ctx context.Context, streamID string, expectedVersion int, events []Event) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    for i, e := range events {
        data, err := json.Marshal(e)
        if err != nil { return fmt.Errorf("marshal %T: %w", e, err) }

        _, err = tx.ExecContext(ctx, `
            INSERT INTO events (stream_id, version, event_type, schema_ver, data)
            VALUES ($1, $2, $3, $4, $5)`,
            streamID, expectedVersion+i+1, e.EventType(), e.SchemaVersion(), data,
        )
        if err != nil {
            var pqErr *pq.Error
            if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
                return ErrConcurrency
            }
            return err
        }
    }
    return tx.Commit()
}
```

The two orderings serve two audiences:

| Column | Ordering | Who uses it |
|--------|----------|-------------|
| `version` | within one stream | aggregates (replay one account) |
| `position` | across all streams | projections (replay everything, in order) |

> **Design rule**: events are immutable. No `UPDATE`, no `DELETE` (except GDPR-style crypto-shredding, where you delete the encryption key instead of the event). If an event was wrong, you append a *correcting* event — exactly like an accountant writes a correcting entry instead of using an eraser.

---

## 3. The Aggregate — Decide and Apply

Chapter 88 showed `Apply`. The production-grade aggregate splits its behavior into two strictly separated halves:

- **Decide** (command methods): validate the command against current state and *emit* events. Never mutates state directly.
- **Apply** (event handler): mutate state from an event. Never validates, never fails on valid history.

```
       command                 event                  state
  Deposit(100) ──decide──► MoneyDeposited(100) ──apply──► balance += 100
                 (may                                (must
                 reject)                             never fail)
```

Why so strict? Because `Apply` runs in two situations: right after `decide` emits an event, and *years later* when replaying history. If `Apply` contained validation, changing a business rule would make old, perfectly valid history fail to load.

```go
type BankAccount struct {
    id      string
    owner   string
    balance int64 // cents — never float64 for money
    open    bool
    version int     // last applied version
    pending []Event // emitted but not yet saved
}

// ---- Decide: business rules live here ----

func (a *BankAccount) Deposit(amount int64) error {
    if !a.open      { return errors.New("account is closed") }
    if amount <= 0  { return errors.New("deposit must be positive") }

    a.raise(MoneyDeposited{AccountID: a.id, Amount: amount, At: time.Now().UTC()})
    return nil
}

func (a *BankAccount) Withdraw(amount int64) error {
    if !a.open             { return errors.New("account is closed") }
    if amount <= 0         { return errors.New("withdrawal must be positive") }
    if amount > a.balance  { return fmt.Errorf("insufficient funds: have %d, want %d", a.balance, amount) }

    a.raise(MoneyWithdrawn{AccountID: a.id, Amount: amount, At: time.Now().UTC()})
    return nil
}

// ---- Apply: pure state transition, no validation ----

func (a *BankAccount) Apply(e Event) {
    switch ev := e.(type) {
    case AccountOpened:
        a.id, a.owner, a.open = ev.AccountID, ev.Owner, true
    case MoneyDeposited:
        a.balance += ev.Amount
    case MoneyWithdrawn:
        a.balance -= ev.Amount
    case AccountClosed:
        a.open = false
    }
    a.version++
}

// raise applies the event immediately (so subsequent decisions see the
// new state) AND records it for saving.
func (a *BankAccount) raise(e Event) {
    a.Apply(e)
    a.pending = append(a.pending, e)
}

func (a *BankAccount) PendingEvents() []Event { return a.pending }
func (a *BankAccount) ClearPending()          { a.pending = nil }
```

Note that `raise` applies immediately. If a command emits two events, the second decision must see the state produced by the first — otherwise `Withdraw(100)` twice in one command could both pass the balance check.

---

## 4. A Generic Event-Sourced Repository

With generics (Chapter 23) you can write the load–decide–save cycle once and reuse it for every aggregate:

```go
// Aggregate is anything that can replay events and expose pending ones.
type Aggregate interface {
    Apply(Event)
    Version() int
    PendingEvents() []Event
    ClearPending()
}

type EventStore interface {
    Append(ctx context.Context, streamID string, expectedVersion int, events []Event) error
    Load(ctx context.Context, streamID string, fromVersion int) ([]Event, error)
}

type ESRepository[T Aggregate] struct {
    store  EventStore
    newFn  func() T // constructor for an empty aggregate
    stream func(id string) string
}

func (r *ESRepository[T]) Get(ctx context.Context, id string) (T, error) {
    agg := r.newFn()
    events, err := r.store.Load(ctx, r.stream(id), 0)
    if err != nil            { return agg, err }
    if len(events) == 0      { return agg, ErrNotFound }

    for _, e := range events { agg.Apply(e) }
    return agg, nil
}

func (r *ESRepository[T]) Save(ctx context.Context, id string, agg T) error {
    pending := agg.PendingEvents()
    if len(pending) == 0 { return nil }

    expected := agg.Version() - len(pending) // version before this command ran
    if err := r.store.Append(ctx, r.stream(id), expected, pending); err != nil {
        return err
    }
    agg.ClearPending()
    return nil
}
```

The command handler becomes a three-liner, and retrying on `ErrConcurrency` is trivial because the whole cycle is side-effect free until `Append` commits:

```go
func (h *AccountHandler) HandleWithdraw(ctx context.Context, cmd WithdrawCommand) error {
    for attempt := 0; attempt < 3; attempt++ {
        acc, err := h.repo.Get(ctx, cmd.AccountID)
        if err != nil { return err }

        if err := acc.Withdraw(cmd.Amount); err != nil { return err } // business rejection: no retry

        err = h.repo.Save(ctx, cmd.AccountID, acc)
        if err == nil                            { return nil }
        if !errors.Is(err, ErrConcurrency)       { return err }
        // someone else won the race — reload and try again
    }
    return ErrConcurrency
}
```

---

## 5. Snapshots

Replaying 20 events is instant. Replaying 2 million events for a hot aggregate (a popular product, a busy account) is not. A **snapshot** is a cached copy of aggregate state at some version, stored separately so it can always be thrown away and rebuilt:

```sql
CREATE TABLE snapshots (
    stream_id  TEXT PRIMARY KEY,
    version    INT  NOT NULL,     -- state as of this version
    data       JSONB NOT NULL,
    taken_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```
Load without snapshot:   replay events 1..50,000                (slow)
Load with snapshot:      load snapshot@49,950 + replay 51..50,000 events... 
                         no — replay only 49,951..50,000        (fast)
```

```go
const snapshotEvery = 100

func (r *SnapshottingRepo) Get(ctx context.Context, id string) (*BankAccount, error) {
    acc := &BankAccount{}

    // 1. Try the snapshot
    snap, err := r.loadSnapshot(ctx, id) // returns nil, nil when absent
    if err != nil { return nil, err }
    if snap != nil {
        if err := acc.RestoreSnapshot(snap.Data); err != nil { return nil, err }
        acc.version = snap.Version
    }

    // 2. Replay only events AFTER the snapshot
    events, err := r.store.Load(ctx, "account-"+id, acc.version)
    if err != nil { return nil, err }
    for _, e := range events { acc.Apply(e) }

    return acc, nil
}

func (r *SnapshottingRepo) Save(ctx context.Context, id string, acc *BankAccount) error {
    before := acc.version - len(acc.pending)
    if err := r.store.Append(ctx, "account-"+id, before, acc.pending); err != nil {
        return err
    }
    acc.ClearPending()

    // 3. Snapshot occasionally — best effort, failures are fine
    if acc.version/snapshotEvery > before/snapshotEvery {
        if data, err := acc.Snapshot(); err == nil {
            _ = r.saveSnapshot(ctx, id, acc.version, data)
        }
    }
    return nil
}
```

Three rules keep snapshots honest:

1. **Snapshots are an optimization, never the source of truth.** Delete the table and the system still works, just slower.
2. **Version your snapshot format too.** When the aggregate struct changes, bump the snapshot version and ignore stale snapshots (fall back to full replay).
3. **Snapshot asynchronously or best-effort.** A failed snapshot write must never fail the command.

---

## 6. Event Versioning and Upcasting

Events live forever, but your code evolves. Suppose v1 of `AccountOpened` stored one `name` field, and v2 splits it into `first_name` and `last_name`. You cannot migrate the stored events (immutable!), and you do not want `switch` statements over versions scattered through every `Apply`.

The answer is an **upcaster**: a function that converts old event JSON to the newest shape at read time, before deserialization. Your domain code only ever sees the latest version.

```
  stored bytes (v1) ──► upcast v1→v2 ──► upcast v2→v3 ──► unmarshal ──► Apply
                        (chain runs at load time, one hop per version)
```

```go
// Upcaster transforms raw event JSON from one schema version to the next.
type Upcaster func(data []byte) ([]byte, error)

type EventCodec struct {
    // decoders for the LATEST version of each type
    decoders map[string]func([]byte) (Event, error)
    // upcasters[eventType][fromVersion] converts fromVersion → fromVersion+1
    upcasters map[string]map[int]Upcaster
    latest    map[string]int
}

func (c *EventCodec) Decode(eventType string, schemaVer int, data []byte) (Event, error) {
    // Walk the upcaster chain until we reach the latest version
    for v := schemaVer; v < c.latest[eventType]; v++ {
        up, ok := c.upcasters[eventType][v]
        if !ok { return nil, fmt.Errorf("no upcaster for %s v%d", eventType, v) }
        var err error
        if data, err = up(data); err != nil { return nil, err }
    }
    decode, ok := c.decoders[eventType]
    if !ok { return nil, fmt.Errorf("unknown event type %q", eventType) }
    return decode(data)
}
```

A concrete upcaster — split `name` into first/last:

```go
codec.RegisterUpcaster("AccountOpened", 1, func(data []byte) ([]byte, error) {
    var v1 struct {
        AccountID string `json:"account_id"`
        Name      string `json:"name"`
    }
    if err := json.Unmarshal(data, &v1); err != nil { return nil, err }

    first, last, _ := strings.Cut(v1.Name, " ")
    return json.Marshal(map[string]any{
        "account_id": v1.AccountID,
        "first_name": first,
        "last_name":  last,
    })
})
```

Guidelines for evolving events:

| Change | Strategy |
|--------|----------|
| Add optional field | Just add it — old JSON unmarshals fine, zero value applies |
| Rename / split / merge fields | Upcaster (as above) |
| Change meaning of a field | New event type — do not reuse the old one |
| Event was simply wrong | Append a correcting event; never edit history |

---

## 7. Projections and Checkpoints

Chapter 88 showed projections updating a read table. The missing production piece is the **checkpoint**: each projection remembers the last global `position` it processed, so it can resume after a crash — or start from 0 to rebuild from scratch.

```sql
CREATE TABLE projection_checkpoints (
    projection TEXT PRIMARY KEY,   -- "account_balances"
    position   BIGINT NOT NULL
);
```

```go
type Projection interface {
    Name() string
    Handle(ctx context.Context, e RecordedEvent) error // RecordedEvent = Event + position
}

// Runner polls the global event log and feeds each projection from its
// own checkpoint. Each projection progresses independently.
func (r *ProjectionRunner) Run(ctx context.Context, p Projection) error {
    ticker := time.NewTicker(200 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            pos, err := r.checkpoint(ctx, p.Name())
            if err != nil { return err }

            events, err := r.store.LoadAllFrom(ctx, pos, 500) // position > pos, LIMIT 500
            if err != nil { return err }

            for _, e := range events {
                if err := p.Handle(ctx, e); err != nil {
                    return fmt.Errorf("projection %s at position %d: %w", p.Name(), e.Position, err)
                }
                if err := r.saveCheckpoint(ctx, p.Name(), e.Position); err != nil {
                    return err
                }
            }
        }
    }
}
```

Because the checkpoint advances only after a successful `Handle`, a crash re-delivers the last event — projections must be **idempotent** (`ON CONFLICT DO NOTHING` / `DO UPDATE`, as in Chapter 88).

**Rebuilding** a read model is now a superpower instead of a migration nightmare:

```go
// Deploy new projection code, then:
TRUNCATE account_balances;
UPDATE projection_checkpoints SET position = 0 WHERE projection = 'account_balances';
// the runner replays all of history into the new shape
```

This is why event-sourced teams say "the read model is disposable" — any bug in a projection is fixed by fixing the code and replaying.

---

## 8. A Complete Runnable Example

Everything in one file with an in-memory store — no database needed. Save as `main.go` and `go run` it (Go 1.22+):

```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ---------- Events ----------

type Event interface{ EventType() string }

type AccountOpened struct {
	AccountID string
	Owner     string
	At        time.Time
}
type MoneyDeposited struct {
	AccountID string
	Amount    int64
	At        time.Time
}
type MoneyWithdrawn struct {
	AccountID string
	Amount    int64
	At        time.Time
}

func (AccountOpened) EventType() string  { return "AccountOpened" }
func (MoneyDeposited) EventType() string { return "MoneyDeposited" }
func (MoneyWithdrawn) EventType() string { return "MoneyWithdrawn" }

// ---------- Aggregate ----------

type BankAccount struct {
	id      string
	owner   string
	balance int64
	open    bool
	version int
	pending []Event
}

func OpenAccount(id, owner string) *BankAccount {
	a := &BankAccount{}
	a.raise(AccountOpened{AccountID: id, Owner: owner, At: time.Now().UTC()})
	return a
}

func (a *BankAccount) Deposit(amount int64) error {
	if !a.open     { return errors.New("account is closed") }
	if amount <= 0 { return errors.New("deposit must be positive") }
	a.raise(MoneyDeposited{AccountID: a.id, Amount: amount, At: time.Now().UTC()})
	return nil
}

func (a *BankAccount) Withdraw(amount int64) error {
	if !a.open            { return errors.New("account is closed") }
	if amount <= 0        { return errors.New("withdrawal must be positive") }
	if amount > a.balance { return fmt.Errorf("insufficient funds: have %d, want %d", a.balance, amount) }
	a.raise(MoneyWithdrawn{AccountID: a.id, Amount: amount, At: time.Now().UTC()})
	return nil
}

func (a *BankAccount) Apply(e Event) {
	switch ev := e.(type) {
	case AccountOpened:
		a.id, a.owner, a.open = ev.AccountID, ev.Owner, true
	case MoneyDeposited:
		a.balance += ev.Amount
	case MoneyWithdrawn:
		a.balance -= ev.Amount
	}
	a.version++
}

func (a *BankAccount) raise(e Event) {
	a.Apply(e)
	a.pending = append(a.pending, e)
}

// ---------- In-memory event store ----------

var ErrConcurrency = errors.New("concurrency conflict")

type MemoryEventStore struct {
	mu      sync.Mutex
	streams map[string][]Event
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{streams: make(map[string][]Event)}
}

func (s *MemoryEventStore) Append(streamID string, expectedVersion int, events []Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.streams[streamID]) != expectedVersion {
		return fmt.Errorf("%w: stream %s at %d, expected %d",
			ErrConcurrency, streamID, len(s.streams[streamID]), expectedVersion)
	}
	s.streams[streamID] = append(s.streams[streamID], events...)
	return nil
}

func (s *MemoryEventStore) Load(streamID string) []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Event, len(s.streams[streamID]))
	copy(out, s.streams[streamID])
	return out
}

// ---------- Repository ----------

type AccountRepo struct{ store *MemoryEventStore }

func (r *AccountRepo) Get(id string) (*BankAccount, error) {
	events := r.store.Load("account-" + id)
	if len(events) == 0 { return nil, errors.New("not found") }
	acc := &BankAccount{}
	for _, e := range events { acc.Apply(e) }
	return acc, nil
}

func (r *AccountRepo) Save(id string, acc *BankAccount) error {
	expected := acc.version - len(acc.pending)
	if err := r.store.Append("account-"+id, expected, acc.pending); err != nil {
		return err
	}
	acc.pending = nil
	return nil
}

// ---------- Demo ----------

func main() {
	store := NewMemoryEventStore()
	repo := &AccountRepo{store: store}

	// Open an account and run some commands
	acc := OpenAccount("42", "Alice")
	must(acc.Deposit(500))
	must(acc.Withdraw(200))
	must(acc.Deposit(50))
	must(repo.Save("42", acc))

	// Rehydrate from scratch — state comes purely from replay
	loaded, err := repo.Get("42")
	must(err)
	fmt.Printf("balance=%d version=%d owner=%s\n", loaded.balance, loaded.version, loaded.owner)
	// balance=350 version=4 owner=Alice

	// A business rule rejection emits nothing
	if err := loaded.Withdraw(10_000); err != nil {
		fmt.Println("rejected:", err) // rejected: insufficient funds: have 350, want 10000
	}

	// Optimistic concurrency: two loads, two conflicting saves
	a1, _ := repo.Get("42")
	a2, _ := repo.Get("42")
	must(a1.Deposit(1))
	must(a2.Deposit(1))
	must(repo.Save("42", a1))                     // wins
	fmt.Println("conflict:", repo.Save("42", a2)) // loses with ErrConcurrency

	// Full audit trail, for free
	for i, e := range store.Load("account-42") {
		fmt.Printf("  %d: %s %+v\n", i+1, e.EventType(), e)
	}
}

func must(err error) {
	if err != nil { panic(err) }
}
```

Run it and read the output carefully — the balance was never stored anywhere, yet `repo.Get` reconstructs it perfectly, the audit trail is complete, and the losing concurrent writer is rejected exactly as a PostgreSQL unique constraint would reject it.

---

## Summary

- Event store = append-only table with **two orderings**: `version` (per stream, for aggregates) and `position` (global, for projections)
- `UNIQUE (stream_id, version)` gives optimistic concurrency for free — losers get `ErrConcurrency` and retry the whole load–decide–save cycle
- Aggregates split into **decide** (validates, may reject, emits events) and **apply** (pure state transition, must never fail on valid history)
- **Snapshots** cap replay cost: store state every N events, replay only the tail; snapshots are disposable, never the source of truth
- **Upcasting** converts old event JSON to the latest schema at read time — domain code only ever sees the newest version; never edit stored events
- **Projections** track a checkpoint (last processed global position); crash-safe because handlers are idempotent, and rebuildable by resetting the checkpoint to 0
- Events are immutable forever — fix mistakes with correcting events, not edits

## Exercises

### Easy
1. Add an `AccountClosed` event and a `Close()` command to the runnable example. Closing requires balance = 0. Verify that `Deposit` on a closed account is rejected but replaying a stream that ends in `AccountClosed` still works.
2. Add a `TransactionCount()` method that is computed during `Apply` (increment on every deposit/withdrawal). Confirm it survives rehydration.
3. Write a test that appends 3 events, loads the aggregate, and asserts `version == 3` and the correct balance. Then simulate a concurrency conflict and assert `errors.Is(err, ErrConcurrency)`.

### Medium
4. Add snapshots to the in-memory example: a `map[string]snapshot` holding `{state, version}`, taken every 5 events. Instrument `Get` to print how many events it replayed, and verify it drops from N to <5 after the first snapshot.
5. Implement the `EventCodec` from §6 with a real upcaster: `MoneyDeposited` v1 stores `amount` in rupees (int); v2 stores `amount_cents`. Store some v1 JSON, decode through the codec, and verify the aggregate sees cents.
6. Build an `AccountBalancesProjection` with a checkpoint over the in-memory store (give `MemoryEventStore` a global `[]RecordedEvent` log with positions). Kill and restart the projection mid-stream and prove no event is skipped or double-counted (make the handler idempotent).

### Hard
7. Port the runnable example to PostgreSQL using the schema from §2: `Append` with the unique constraint, `Load`, snapshots table, and a projection with `projection_checkpoints`. Run two writer goroutines hammering the same account and verify the final balance equals the sum of successful commands.
8. Implement **temporal queries**: `BalanceAt(ctx, accountID string, t time.Time)` that replays only events with `occurred_at <= t`. Then build a "monthly statement" generator that emits opening balance, all transactions, and closing balance for any month — entirely from the event stream, with no extra tables.
