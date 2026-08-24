# Chapter 57: Event Sourcing and CQRS in Go

Event sourcing is the pattern that makes systems auditable, debuggable, and time-travelable. Combined with CQRS, it separates "how you write data" from "how you read data" — each optimized for its purpose.

## Table of Contents

1. Event Sourcing Deep Dive
2. The Event Store
3. Aggregates and Domain Events
4. Projections (Read Models)
5. CQRS: Separate Write and Read
6. Mini Project: Bank Account System
7. Exercises

---

## 1. Event Sourcing Deep Dive

**Traditional CRUD:**
```
bank_accounts: {id: "acc-1", owner: "alice", balance: 150}
```
One row. Current state. No history.

**Event Sourcing:**
```
events:
  {id: 1, account_id: "acc-1", type: "AccountOpened",   data: {owner: "alice", initial: 0}}
  {id: 2, account_id: "acc-1", type: "MoneyDeposited",  data: {amount: 200}}
  {id: 3, account_id: "acc-1", type: "MoneyWithdrawn",  data: {amount: 50}}
  
Current balance = 0 + 200 - 50 = 150 (by replaying events)
```

**Benefits:**
- **Audit trail:** Every change is recorded with who did it and when.
- **Time travel:** Reconstruct state at any point in time.
- **Replay:** Rebuild any projection from scratch by replaying events.
- **Debugging:** Events tell you exactly what happened.
- **Integration:** New services can subscribe to the event stream.

**Trade-offs:**
- More complex than CRUD.
- Current state requires replaying events (mitigated by snapshots).
- Schema evolution of events is tricky.

---

## 2. The Event Store

The event store is a specialized database for events: append-only, versioned per aggregate.

```go
package eventsource

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
)

type Event struct {
    ID          int64
    AggregateID string
    Version     int    // monotonically increasing per aggregate
    Type        string
    Data        json.RawMessage
    OccurredAt  time.Time
}

type EventStore struct {
    db *sql.DB
}

func NewEventStore(db *sql.DB) *EventStore {
    db.Exec(`CREATE TABLE IF NOT EXISTS events (
        id           BIGSERIAL PRIMARY KEY,
        aggregate_id VARCHAR(255) NOT NULL,
        version      INT NOT NULL,
        type         VARCHAR(255) NOT NULL,
        data         JSONB NOT NULL,
        occurred_at  TIMESTAMPTZ DEFAULT NOW(),
        UNIQUE (aggregate_id, version)  -- optimistic concurrency control
    )`)
    db.Exec(`CREATE INDEX IF NOT EXISTS events_aggregate_id ON events(aggregate_id)`)
    return &EventStore{db: db}
}

// Append adds events for an aggregate.
// expectedVersion: the version we expect the aggregate to be at.
// If the aggregate has changed (concurrent modification), this fails.
func (es *EventStore) Append(ctx context.Context, aggregateID string, expectedVersion int, events []Event) error {
    tx, err := es.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Optimistic concurrency: check current version
    var currentVersion sql.NullInt64
    tx.QueryRowContext(ctx,
        "SELECT MAX(version) FROM events WHERE aggregate_id = $1",
        aggregateID).Scan(&currentVersion)

    cv := 0
    if currentVersion.Valid {
        cv = int(currentVersion.Int64)
    }
    if cv != expectedVersion {
        return fmt.Errorf("optimistic concurrency conflict: expected version %d, got %d",
            expectedVersion, cv)
    }

    // Append new events
    for i, event := range events {
        _, err := tx.ExecContext(ctx,
            "INSERT INTO events (aggregate_id, version, type, data) VALUES ($1, $2, $3, $4)",
            aggregateID, expectedVersion+i+1, event.Type, event.Data)
        if err != nil {
            return fmt.Errorf("append event: %w", err)
        }
    }

    return tx.Commit()
}

// Load reads all events for an aggregate, optionally from a specific version
func (es *EventStore) Load(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error) {
    rows, err := es.db.QueryContext(ctx,
        "SELECT id, aggregate_id, version, type, data, occurred_at FROM events WHERE aggregate_id = $1 AND version > $2 ORDER BY version",
        aggregateID, fromVersion)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []Event
    for rows.Next() {
        var e Event
        err := rows.Scan(&e.ID, &e.AggregateID, &e.Version, &e.Type, &e.Data, &e.OccurredAt)
        if err != nil {
            return nil, err
        }
        events = append(events, e)
    }
    return events, rows.Err()
}

// LoadFrom returns all events after a given global event ID (for projections)
func (es *EventStore) LoadFrom(ctx context.Context, fromID int64, limit int) ([]Event, error) {
    rows, err := es.db.QueryContext(ctx,
        "SELECT id, aggregate_id, version, type, data, occurred_at FROM events WHERE id > $1 ORDER BY id LIMIT $2",
        fromID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []Event
    for rows.Next() {
        var e Event
        rows.Scan(&e.ID, &e.AggregateID, &e.Version, &e.Type, &e.Data, &e.OccurredAt)
        events = append(events, e)
    }
    return events, rows.Err()
}
```

---

## 3. Aggregates and Domain Events

An **aggregate** is a group of objects treated as a unit. It applies commands (intents to change state) and produces events (records of what happened).

```go
// Domain events
type AccountOpenedEvent struct {
    Owner          string  `json:"owner"`
    InitialBalance float64 `json:"initial_balance"`
}

type MoneyDepositedEvent struct {
    Amount float64 `json:"amount"`
}

type MoneyWithdrawnEvent struct {
    Amount float64 `json:"amount"`
}

type TransferredEvent struct {
    ToAccount string  `json:"to_account"`
    Amount    float64 `json:"amount"`
}

// Aggregate: BankAccount
type BankAccount struct {
    ID      string
    Owner   string
    Balance float64
    Version int  // current event version
}

// Apply applies an event to update the aggregate state (no side effects)
func (a *BankAccount) Apply(event Event) error {
    switch event.Type {
    case "AccountOpened":
        var e AccountOpenedEvent
        if err := json.Unmarshal(event.Data, &e); err != nil {
            return err
        }
        a.Owner = e.Owner
        a.Balance = e.InitialBalance

    case "MoneyDeposited":
        var e MoneyDepositedEvent
        json.Unmarshal(event.Data, &e)
        a.Balance += e.Amount

    case "MoneyWithdrawn":
        var e MoneyWithdrawnEvent
        json.Unmarshal(event.Data, &e)
        a.Balance -= e.Amount
    }
    a.Version = event.Version
    return nil
}

// Commands: validate business rules and produce events

func (a *BankAccount) Deposit(amount float64) ([]Event, error) {
    if amount <= 0 {
        return nil, fmt.Errorf("deposit amount must be positive")
    }
    data, _ := json.Marshal(MoneyDepositedEvent{Amount: amount})
    return []Event{{Type: "MoneyDeposited", Data: data}}, nil
}

func (a *BankAccount) Withdraw(amount float64) ([]Event, error) {
    if amount <= 0 {
        return nil, fmt.Errorf("withdrawal amount must be positive")
    }
    if a.Balance < amount {
        return nil, fmt.Errorf("insufficient funds: have %.2f, need %.2f", a.Balance, amount)
    }
    data, _ := json.Marshal(MoneyWithdrawnEvent{Amount: amount})
    return []Event{{Type: "MoneyWithdrawn", Data: data}}, nil
}

// Repository: load and save aggregates
type BankAccountRepository struct {
    store *EventStore
}

func (r *BankAccountRepository) Load(ctx context.Context, accountID string) (*BankAccount, error) {
    events, err := r.store.Load(ctx, accountID, 0)
    if err != nil {
        return nil, err
    }
    if len(events) == 0 {
        return nil, fmt.Errorf("account %s not found", accountID)
    }

    account := &BankAccount{ID: accountID}
    for _, e := range events {
        if err := account.Apply(e); err != nil {
            return nil, err
        }
    }
    return account, nil
}

func (r *BankAccountRepository) Save(ctx context.Context, account *BankAccount, newEvents []Event) error {
    return r.store.Append(ctx, account.ID, account.Version, newEvents)
}
```

---

## 4. Projections (Read Models)

Projections rebuild read models from the event stream:

```go
// AccountSummary: denormalized read model for the "list accounts" query
type AccountSummary struct {
    AccountID string
    Owner     string
    Balance   float64
    TxCount   int
    UpdatedAt time.Time
}

type AccountProjection struct {
    db           *sql.DB
    lastEventID  int64
}

func NewAccountProjection(db *sql.DB) *AccountProjection {
    db.Exec(`CREATE TABLE IF NOT EXISTS account_summaries (
        account_id VARCHAR(255) PRIMARY KEY,
        owner      VARCHAR(255),
        balance    DECIMAL(15,2),
        tx_count   INT DEFAULT 0,
        updated_at TIMESTAMPTZ
    )`)
    return &AccountProjection{db: db}
}

// Handle processes one event and updates the read model
func (p *AccountProjection) Handle(event Event) error {
    switch event.Type {
    case "AccountOpened":
        var e AccountOpenedEvent
        json.Unmarshal(event.Data, &e)
        _, err := p.db.Exec(`
            INSERT INTO account_summaries (account_id, owner, balance, tx_count, updated_at)
            VALUES ($1, $2, $3, 0, $4)
            ON CONFLICT (account_id) DO NOTHING`,
            event.AggregateID, e.Owner, e.InitialBalance, event.OccurredAt)
        return err

    case "MoneyDeposited":
        var e MoneyDepositedEvent
        json.Unmarshal(event.Data, &e)
        _, err := p.db.Exec(`
            UPDATE account_summaries
            SET balance = balance + $2, tx_count = tx_count + 1, updated_at = $3
            WHERE account_id = $1`,
            event.AggregateID, e.Amount, event.OccurredAt)
        return err

    case "MoneyWithdrawn":
        var e MoneyWithdrawnEvent
        json.Unmarshal(event.Data, &e)
        _, err := p.db.Exec(`
            UPDATE account_summaries
            SET balance = balance - $2, tx_count = tx_count + 1, updated_at = $3
            WHERE account_id = $1`,
            event.AggregateID, e.Amount, event.OccurredAt)
        return err
    }
    return nil
}

// Run continuously polls the event store for new events
func (p *AccountProjection) Run(ctx context.Context, store *EventStore) {
    for {
        events, err := store.LoadFrom(ctx, p.lastEventID, 100)
        if err != nil {
            log.Printf("projection error: %v", err)
            time.Sleep(time.Second)
            continue
        }
        for _, e := range events {
            if err := p.Handle(e); err != nil {
                log.Printf("handle event %d: %v", e.ID, err)
                continue
            }
            p.lastEventID = e.ID
        }
        if len(events) < 100 {
            time.Sleep(100 * time.Millisecond)
        }
    }
}
```

---

## 5. CQRS: Separate Write and Read

```go
// CQRS Command Service (write side)
type CommandService struct {
    repo *BankAccountRepository
}

func (s *CommandService) OpenAccount(ctx context.Context, owner string, initial float64) (string, error) {
    id := generateID()
    data, _ := json.Marshal(AccountOpenedEvent{Owner: owner, InitialBalance: initial})
    events := []Event{{Type: "AccountOpened", Data: data}}
    account := &BankAccount{ID: id, Version: 0}
    if err := s.repo.Save(ctx, account, events); err != nil {
        return "", err
    }
    return id, nil
}

func (s *CommandService) Deposit(ctx context.Context, accountID string, amount float64) error {
    account, err := s.repo.Load(ctx, accountID)
    if err != nil {
        return err
    }
    events, err := account.Deposit(amount)
    if err != nil {
        return err
    }
    return s.repo.Save(ctx, account, events)
}

func (s *CommandService) Withdraw(ctx context.Context, accountID string, amount float64) error {
    account, err := s.repo.Load(ctx, accountID)
    if err != nil {
        return err
    }
    events, err := account.Withdraw(amount)
    if err != nil {
        return err
    }
    return s.repo.Save(ctx, account, events)
}

// CQRS Query Service (read side)
type QueryService struct {
    db *sql.DB
}

func (s *QueryService) GetBalance(ctx context.Context, accountID string) (float64, error) {
    var balance float64
    err := s.db.QueryRowContext(ctx,
        "SELECT balance FROM account_summaries WHERE account_id = $1",
        accountID).Scan(&balance)
    return balance, err
}

func (s *QueryService) ListAccounts(ctx context.Context) ([]AccountSummary, error) {
    rows, err := s.db.QueryContext(ctx,
        "SELECT account_id, owner, balance, tx_count FROM account_summaries ORDER BY balance DESC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var summaries []AccountSummary
    for rows.Next() {
        var s AccountSummary
        rows.Scan(&s.AccountID, &s.Owner, &s.Balance, &s.TxCount)
        summaries = append(summaries, s)
    }
    return summaries, nil
}
```

---

## 6. Mini Project: Bank Account System

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "postgres://localhost/eventsource?sslmode=disable")
    
    store := eventsource.NewEventStore(db)
    repo := &eventsource.BankAccountRepository{Store: store}
    cmdSvc := &eventsource.CommandService{Repo: repo}
    
    projection := eventsource.NewAccountProjection(db)
    go projection.Run(context.Background(), store)
    
    querySvc := &eventsource.QueryService{DB: db}

    mux := http.NewServeMux()

    mux.HandleFunc("POST /accounts", func(w http.ResponseWriter, r *http.Request) {
        var req struct{ Owner string; Initial float64 }
        json.NewDecoder(r.Body).Decode(&req)
        id, err := cmdSvc.OpenAccount(r.Context(), req.Owner, req.Initial)
        if err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{"account_id": id})
    })

    mux.HandleFunc("POST /accounts/{id}/deposit", func(w http.ResponseWriter, r *http.Request) {
        var req struct{ Amount float64 }
        json.NewDecoder(r.Body).Decode(&req)
        if err := cmdSvc.Deposit(r.Context(), r.PathValue("id"), req.Amount); err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        w.WriteHeader(204)
    })

    mux.HandleFunc("POST /accounts/{id}/withdraw", func(w http.ResponseWriter, r *http.Request) {
        var req struct{ Amount float64 }
        json.NewDecoder(r.Body).Decode(&req)
        if err := cmdSvc.Withdraw(r.Context(), r.PathValue("id"), req.Amount); err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        w.WriteHeader(204)
    })

    mux.HandleFunc("GET /accounts/{id}/balance", func(w http.ResponseWriter, r *http.Request) {
        balance, err := querySvc.GetBalance(r.Context(), r.PathValue("id"))
        if err != nil {
            http.Error(w, err.Error(), 404)
            return
        }
        json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
    })

    log.Println("Bank API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

---

## Summary

- Event sourcing stores all changes as immutable events. Current state = replay of all events.
- Optimistic concurrency: use `(aggregate_id, version)` unique constraint to detect conflicts.
- Aggregates validate business rules and produce events. They never directly mutate persistent state.
- Projections rebuild read models from the event stream — they can be rebuilt from scratch any time.
- CQRS: write side uses event sourcing, read side uses optimized projections (denormalized, fast queries).

### Exercises

**Easy:** Extend the bank account to support transfers: `Transfer(fromID, toID, amount)`. This creates two events atomically: `MoneyWithdrawn` on the source account and `MoneyDeposited` on the destination.

**Medium:** Add snapshots: every 50 events, save a snapshot of the aggregate state. Loading an account should start from the most recent snapshot instead of replaying from the beginning.

**Hard:** Add a "transaction history" projection: a table `transactions(account_id, type, amount, balance_after, occurred_at)` populated from deposit/withdraw events. Expose `GET /accounts/{id}/transactions` that returns paginated history.
