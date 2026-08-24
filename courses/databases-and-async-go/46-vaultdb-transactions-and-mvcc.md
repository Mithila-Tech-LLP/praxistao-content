# Chapter 46: VaultDB — Transactions and MVCC

Two users modify the same row at the same time. Without transactions, both changes fight over the row and one gets silently lost. With MVCC, both can proceed simultaneously and see consistent data. This chapter adds transaction safety to VaultDB.

## Table of Contents

1. The Problem With Concurrent Access
2. MVCC — One Row, Many Versions
3. Transaction IDs and Visibility
4. Implementing MVCC in VaultDB
5. Isolation Levels
6. Deadlock Detection
7. Exercises

---

## 1. The Problem With Concurrent Access

Imagine two Go routines running concurrently:

```
Goroutine A: Read balance (=$100), compute $100+$50, write $150
Goroutine B: Read balance (=$100), compute $100+$30, write $130

Result: $130 (B wins) — A's $50 deposit is silently lost!
```

This is called a **lost update**. To prevent it, database systems provide **transactions** with isolation guarantees.

---

## 2. MVCC — One Row, Many Versions

Instead of locking rows so only one transaction can access them, MVCC (Multi-Version Concurrency Control) keeps multiple versions of each row. Each transaction sees the version that was current when the transaction started.

```
Original row (inserted by txn 1):   {age=25, xmin=1, xmax=∞}
Updated by txn 5:                   {age=26, xmin=5, xmax=∞}
  Old version:                      {age=25, xmin=1, xmax=5}

Txn 3 (started before txn 5): sees {age=25} — the version current when txn 3 started
Txn 7 (started after txn 5):  sees {age=26} — the latest committed version
```

**MVCC rules:**
- `xmin`: the transaction ID that created this row version
- `xmax`: the transaction ID that deleted/updated this row version (∞ = still active)
- A row version is **visible** to transaction T if:
  - `xmin` is committed AND `xmin` ≤ T's snapshot
  - `xmax` is not committed OR `xmax` > T's snapshot

---

## 3. Transaction IDs and Visibility

```go
// txn/mvcc.go
package txn

import (
    "sync"
    "sync/atomic"
)

// TxnID is a monotonically increasing transaction identifier
type TxnID uint64

const (
    InvalidTxnID  TxnID = 0
    FrozenTxnID   TxnID = 1  // for pages that predate MVCC
    FirstUserTxnID TxnID = 2
)

// TxnStatus tracks whether a transaction is active, committed, or aborted
type TxnStatus uint8

const (
    StatusActive    TxnStatus = 0
    StatusCommitted TxnStatus = 1
    StatusAborted   TxnStatus = 2
)

// Manager tracks all active and recently completed transactions
type Manager struct {
    mu        sync.Mutex
    nextID    uint64
    status    map[TxnID]TxnStatus
    active    map[TxnID]*Transaction
}

func NewManager() *Manager {
    return &Manager{
        nextID: uint64(FirstUserTxnID),
        status: map[TxnID]TxnStatus{FrozenTxnID: StatusCommitted},
        active: make(map[TxnID]*Transaction),
    }
}

// Transaction represents one open transaction
type Transaction struct {
    ID       TxnID
    snapshot Snapshot  // what this txn can see
    mgr      *Manager
}

// Snapshot captures which transactions were active when this txn started.
// A txn can only see changes from transactions that committed BEFORE this snapshot.
type Snapshot struct {
    minActive TxnID   // smallest active txnID at snapshot time
    maxSeen   TxnID   // largest txnID seen (even if not committed)
    active    map[TxnID]bool // txns active at snapshot time (not yet committed)
}

func (m *Manager) Begin() *Transaction {
    m.mu.Lock()
    defer m.mu.Unlock()

    id := TxnID(atomic.AddUint64(&m.nextID, 1))

    // Build snapshot: which transactions are currently active?
    snap := Snapshot{
        minActive: id,
        maxSeen:   id,
        active:    make(map[TxnID]bool),
    }
    for activeID := range m.active {
        snap.active[activeID] = true
        if activeID < snap.minActive {
            snap.minActive = activeID
        }
    }

    txn := &Transaction{ID: id, snapshot: snap, mgr: m}
    m.status[id] = StatusActive
    m.active[id] = txn
    return txn
}

func (m *Manager) Commit(txn *Transaction) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.status[txn.ID] = StatusCommitted
    delete(m.active, txn.ID)
}

func (m *Manager) Abort(txn *Transaction) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.status[txn.ID] = StatusAborted
    delete(m.active, txn.ID)
}

// IsVisible returns true if a row version with (xmin, xmax) is visible to txn
func (m *Manager) IsVisible(txn *Transaction, xmin, xmax TxnID) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Row must have been created by a committed transaction before this snapshot
    if !m.isCommittedBefore(txn, xmin) {
        return false
    }

    // Row must not have been deleted by a committed transaction before this snapshot
    if xmax != InvalidTxnID && m.isCommittedBefore(txn, xmax) {
        return false
    }

    return true
}

func (m *Manager) isCommittedBefore(txn *Transaction, id TxnID) bool {
    // Our own writes are always visible
    if id == txn.ID {
        return true
    }
    // Transactions that started after our snapshot are not visible
    if id > txn.snapshot.maxSeen {
        return false
    }
    // Transactions that were active when we started are not visible
    if txn.snapshot.active[id] {
        return false
    }
    // The transaction committed before our snapshot started
    return m.status[id] == StatusCommitted
}
```

---

## 4. Implementing MVCC in VaultDB

Add xmin/xmax to each row stored on disk:

```go
// Updated EncodeRow for MVCC
// Layout: [xmin: 8 bytes] [xmax: 8 bytes] [type-len-value for each column...]

func EncodeMVCCRow(row storage.Row, xmin, xmax txn.TxnID) []byte {
    var buf []byte

    // Prepend xmin and xmax
    buf = append(buf, make([]byte, 8)...)
    binary.BigEndian.PutUint64(buf[0:], uint64(xmin))

    xmaxBytes := make([]byte, 8)
    binary.BigEndian.PutUint64(xmaxBytes, uint64(xmax))
    buf = append(buf, xmaxBytes...)

    // Append encoded column data
    buf = append(buf, storage.EncodeRow(row)...)
    return buf
}

func DecodeMVCCRow(data []byte, numCols int) (row storage.Row, xmin, xmax txn.TxnID, err error) {
    if len(data) < 16 {
        return nil, 0, 0, fmt.Errorf("MVCC row too short")
    }
    xmin = txn.TxnID(binary.BigEndian.Uint64(data[0:8]))
    xmax = txn.TxnID(binary.BigEndian.Uint64(data[8:16]))
    row, err = storage.DecodeRow(data[16:], numCols)
    return
}
```

**Modified SELECT to filter by MVCC visibility:**

```go
func (e *Executor) execSelectMVCC(stmt *SelectStmt, t *txn.Transaction) (*Result, error) {
    tbl, err := e.getTable(stmt.Table)
    if err != nil {
        return nil, err
    }

    heap := storage.NewHeap(e.dm, tbl.RootPageID)
    rawRows, err := heap.ScanAllRaw()  // returns raw bytes
    if err != nil {
        return nil, err
    }

    var visible []storage.Row
    for _, raw := range rawRows {
        row, xmin, xmax, err := DecodeMVCCRow(raw, len(tbl.Columns))
        if err != nil {
            continue
        }
        if e.txnMgr.IsVisible(t, xmin, xmax) {
            visible = append(visible, row)
        }
    }

    // Filter by WHERE, project columns... (same as before)
    return &Result{Rows: visible}, nil
}
```

---

## 5. Isolation Levels

Different applications need different levels of isolation:

```go
type IsolationLevel int

const (
    ReadCommitted  IsolationLevel = 1 // See any committed change, even if after our start
    RepeatableRead IsolationLevel = 2 // See only data committed before we started
    Serializable   IsolationLevel = 3 // As if all transactions ran one at a time
)
```

**Read Committed (default in PostgreSQL):**
Each statement in a transaction gets a fresh snapshot. You might see different data if you run the same SELECT twice in one transaction.

```
T1: SELECT balance → $100
(T2 commits, changes balance to $50)
T1: SELECT balance → $50 (different result! phantom read)
```

**Repeatable Read:**
Snapshot is taken once at the start of the transaction. Running the same query twice gives the same result.

```
T1 starts, snapshot captures current state
T1: SELECT balance → $100
(T2 commits, changes balance to $50)
T1: SELECT balance → $100 (still! snapshot protects us)
```

Our MVCC implementation defaults to Repeatable Read — the snapshot is fixed when `Begin()` is called.

---

## 6. Deadlock Detection

With transactions, deadlocks can occur: T1 waits for T2, T2 waits for T1. Both wait forever.

Simple detection: track the wait-for graph and find cycles.

```go
type LockManager struct {
    mu      sync.Mutex
    locks   map[RowKey]TxnID        // row → who holds the lock
    waiting map[TxnID]TxnID         // txn → txn it's waiting for
}

type RowKey struct {
    TableName string
    RowID     storage.RowID
}

func (lm *LockManager) AcquireLock(txn *txn.Transaction, row RowKey) error {
    lm.mu.Lock()

    holder, locked := lm.locks[row]
    if !locked || holder == txn.ID {
        lm.locks[row] = txn.ID
        lm.mu.Unlock()
        return nil
    }

    // Row is locked by another transaction — detect deadlock
    if lm.wouldDeadlock(txn.ID, holder) {
        lm.mu.Unlock()
        return fmt.Errorf("deadlock detected: txn %d waits for txn %d", txn.ID, holder)
    }

    lm.waiting[txn.ID] = holder
    lm.mu.Unlock()

    // Wait for holder to release (simplified — use condition variable in production)
    return nil
}

func (lm *LockManager) wouldDeadlock(waiter, holder txn.TxnID) bool {
    visited := make(map[txn.TxnID]bool)
    current := holder
    for {
        if current == waiter {
            return true // cycle found
        }
        if visited[current] {
            return false
        }
        visited[current] = true
        next, waiting := lm.waiting[current]
        if !waiting {
            return false
        }
        current = next
    }
}

func (lm *LockManager) ReleaseLocks(txnID txn.TxnID) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    for row, holder := range lm.locks {
        if holder == txnID {
            delete(lm.locks, row)
        }
    }
    delete(lm.waiting, txnID)
}
```

---

## Summary

- MVCC keeps multiple versions of each row. Readers see old versions; writers create new versions. No readers block writers, no writers block readers.
- Each row has `xmin` (creating transaction) and `xmax` (deleting/updating transaction).
- A transaction's snapshot determines which row versions it can see.
- Isolation levels control how "current" the snapshot is: Read Committed = fresh per statement, Repeatable Read = fixed at transaction start.
- Deadlock detection: track the "waits-for" graph, abort one transaction if a cycle is found.

### Exercises

**Easy:** Write a test demonstrating that two concurrent reads don't block each other in MVCC. Use goroutines and verify both reads complete without waiting for the other.

**Medium:** Implement Repeatable Read by ensuring that re-reading the same row in the same transaction always returns the same value, even if another transaction commits an update in between.

**Hard:** Implement a full Serializable isolation check using SSI (Serializable Snapshot Isolation): track read and write sets per transaction; at commit time, check for dangerous read-write-write cycles across transactions and abort the younger transaction if a cycle is detected.
