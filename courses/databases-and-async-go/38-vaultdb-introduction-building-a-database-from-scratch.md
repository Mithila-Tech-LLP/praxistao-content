# Chapter 38: VaultDB — Building a Database From Scratch

We've used databases. Now we build one. VaultDB is a simple but real relational database written entirely in Go — it handles storage, indexing, transactions, and a query language. Building it teaches you more about databases than reading any documentation.

## Table of Contents

1. What We're Building
2. Architecture Overview
3. Project Setup
4. The Fundamental Unit: The Page
5. VaultDB Data Types
6. What Comes Next

---

## 1. What We're Building

VaultDB will support:
- Creating and dropping tables
- INSERT, SELECT, UPDATE, DELETE
- WHERE clauses with conditions
- B-Tree indexes for fast lookups
- ACID transactions with a WAL (Write-Ahead Log)
- A simple binary protocol for client connections

It won't support (to keep it learnable):
- JOINs (too complex for this scope)
- Subqueries
- Stored procedures

By the end, you'll have a working database that can store millions of rows and serve multiple concurrent clients. More importantly, you'll understand exactly why every decision in PostgreSQL or SQLite was made.

---

## 2. Architecture Overview

```
                        ┌─────────────────────────────────────────────┐
                        │                 VaultDB                       │
                        │                                               │
                        │  ┌─────────┐      ┌──────────────────────┐  │
   Client ─────────────►│  │  Wire   │      │    Query Engine      │  │
   (SQL text)           │  │Protocol │─────►│  Lexer → Parser →    │  │
                        │  └─────────┘      │  Planner → Executor  │  │
                        │                   └────────────┬─────────┘  │
                        │                                │             │
                        │                   ┌────────────▼─────────┐  │
                        │                   │  Transaction Manager  │  │
                        │                   │  (MVCC + WAL)         │  │
                        │                   └────────────┬─────────┘  │
                        │                                │             │
                        │                   ┌────────────▼─────────┐  │
                        │                   │    Buffer Pool        │  │
                        │                   │  (in-memory pages)    │  │
                        │                   └────────────┬─────────┘  │
                        │                                │             │
                        │                   ┌────────────▼─────────┐  │
                        │                   │    Storage Engine     │  │
                        │                   │   (disk pages, WAL)   │  │
                        │                   └──────────────────────┘  │
                        └─────────────────────────────────────────────┘
```

**Each layer:**

- **Wire Protocol:** Accepts TCP connections, parses binary messages, sends responses.
- **Query Engine:** Parses SQL text into an AST, creates a query plan, executes it.
- **Transaction Manager:** Ensures atomicity and durability using MVCC and WAL.
- **Buffer Pool:** Caches disk pages in memory. Avoids reading from disk on every access.
- **Storage Engine:** The actual file on disk. Organizes data into fixed-size pages.

---

## 3. Project Setup

```bash
mkdir vaultdb && cd vaultdb
go mod init github.com/yourname/vaultdb
mkdir -p storage btree wal txn query wire
touch main.go
```

Directory structure:
```
vaultdb/
├── main.go           # Server entry point
├── storage/
│   ├── page.go       # Page format
│   ├── disk.go       # File I/O
│   └── buffer.go     # Buffer pool
├── btree/
│   ├── node.go       # B-Tree node layout on a page
│   └── tree.go       # Insert, search, delete
├── wal/
│   └── wal.go        # Write-ahead log
├── txn/
│   └── mvcc.go       # Transaction management
├── query/
│   ├── lexer.go      # SQL tokenizer
│   ├── parser.go     # SQL AST builder
│   ├── planner.go    # Query plan
│   └── executor.go   # Execute the plan
└── wire/
    └── server.go     # TCP server and protocol
```

---

## 4. The Fundamental Unit: The Page

Every database stores data in fixed-size **pages** (SQLite uses 4 KB, PostgreSQL uses 8 KB). We'll use 4 KB.

Why fixed-size pages?
- Predictable I/O: reading one page = one disk seek + one read of known size
- Easy buffer management: buffer pool holds N pages, all same size
- Enables B-Tree indexing: nodes map 1:1 with pages

```go
// storage/page.go
package storage

const PageSize = 4096 // 4 KB

// PageID is the page's position in the file (0-indexed)
type PageID uint64

const InvalidPageID PageID = ^PageID(0)

// Page is a fixed-size block of bytes
type Page [PageSize]byte

// Page header layout (first 16 bytes of every page):
// Bytes  0- 1: page type (leaf, internal, overflow, free)
// Bytes  2- 3: number of slots used
// Bytes  4- 7: free space offset (where free space starts)
// Bytes  8-15: parent page ID (for B-Tree nodes)

type PageType uint16

const (
    PageTypeFree     PageType = 0
    PageTypeLeaf     PageType = 1
    PageTypeInternal PageType = 2
    PageTypeOverflow PageType = 3
)

const (
    headerSize    = 16
    slotSize      = 4 // each slot = 2 bytes offset + 2 bytes length
)

type PageHeader struct {
    Type       PageType
    NumSlots   uint16
    FreeOffset uint16
    ParentID   PageID
}

func (p *Page) GetHeader() PageHeader {
    return PageHeader{
        Type:       PageType(uint16(p[0])<<8 | uint16(p[1])),
        NumSlots:   uint16(p[2])<<8 | uint16(p[3]),
        FreeOffset: uint16(p[4])<<8 | uint16(p[5]),
        ParentID:   PageID(uint64(p[8])<<56 | uint64(p[9])<<48 | uint64(p[10])<<40 |
                           uint64(p[11])<<32 | uint64(p[12])<<24 | uint64(p[13])<<16 |
                           uint64(p[14])<<8 | uint64(p[15])),
    }
}

func (p *Page) SetHeader(h PageHeader) {
    p[0] = byte(h.Type >> 8)
    p[1] = byte(h.Type)
    p[2] = byte(h.NumSlots >> 8)
    p[3] = byte(h.NumSlots)
    p[4] = byte(h.FreeOffset >> 8)
    p[5] = byte(h.FreeOffset)
    p[8]  = byte(h.ParentID >> 56)
    p[9]  = byte(h.ParentID >> 48)
    p[10] = byte(h.ParentID >> 40)
    p[11] = byte(h.ParentID >> 32)
    p[12] = byte(h.ParentID >> 24)
    p[13] = byte(h.ParentID >> 16)
    p[14] = byte(h.ParentID >> 8)
    p[15] = byte(h.ParentID)
}

// FreeSpace returns how many bytes are available for new data
func (p *Page) FreeSpace() int {
    h := p.GetHeader()
    slotTableEnd := headerSize + int(h.NumSlots)*slotSize
    return int(h.FreeOffset) - slotTableEnd
}

// Initialize sets up a fresh page with a header
func (p *Page) Initialize(pageType PageType) {
    // Zero out the page
    for i := range p {
        p[i] = 0
    }
    p.SetHeader(PageHeader{
        Type:       pageType,
        NumSlots:   0,
        FreeOffset: PageSize, // free space starts at the end (grows downward)
        ParentID:   InvalidPageID,
    })
}
```

**The page layout: slot array grows down from the top, data grows up from the bottom:**

```
┌─────────────────────────────────────────────────────┐
│ Header (16 bytes)                                   │  ← offset 0
├─────────────────────────────────────────────────────┤
│ Slot 0: (offset=4080, len=16)                       │  ← offset 16
│ Slot 1: (offset=4064, len=15)                       │  ← offset 20
│ Slot 2: (offset=4048, len=16)                       │  ← offset 24
│ ...                                                 │
│              [ FREE SPACE ]                         │
│                                                     │
│ Row 2 data (16 bytes)                               │  ← offset 4048
│ Row 1 data (15 bytes)                               │  ← offset 4064 (padded to 4079)
│ Row 0 data (16 bytes)                               │  ← offset 4080
└─────────────────────────────────────────────────────┘  ← offset 4096
```

This "slot page" design lets rows be variable length and allows rows to be moved during compaction without changing their slot numbers.

---

## 5. VaultDB Data Types

```go
// storage/types.go
package storage

import (
    "encoding/binary"
    "fmt"
    "math"
)

type TypeID uint8

const (
    TypeInt    TypeID = 1
    TypeFloat  TypeID = 2
    TypeBool   TypeID = 3
    TypeString TypeID = 4
    TypeNull   TypeID = 5
)

// Value is a typed value stored in a row
type Value struct {
    Type TypeID
    Data []byte
}

func IntVal(n int64) Value {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(n))
    return Value{Type: TypeInt, Data: b}
}

func FloatVal(f float64) Value {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, math.Float64bits(f))
    return Value{Type: TypeFloat, Data: b}
}

func BoolVal(b bool) Value {
    if b {
        return Value{Type: TypeBool, Data: []byte{1}}
    }
    return Value{Type: TypeBool, Data: []byte{0}}
}

func StringVal(s string) Value {
    return Value{Type: TypeString, Data: []byte(s)}
}

func NullVal() Value {
    return Value{Type: TypeNull, Data: nil}
}

func (v Value) AsInt() int64 {
    return int64(binary.BigEndian.Uint64(v.Data))
}

func (v Value) AsFloat() float64 {
    return math.Float64frombits(binary.BigEndian.Uint64(v.Data))
}

func (v Value) AsBool() bool {
    return v.Data[0] != 0
}

func (v Value) AsString() string {
    return string(v.Data)
}

func (v Value) String() string {
    switch v.Type {
    case TypeInt:
        return fmt.Sprintf("%d", v.AsInt())
    case TypeFloat:
        return fmt.Sprintf("%g", v.AsFloat())
    case TypeBool:
        if v.AsBool() {
            return "true"
        }
        return "false"
    case TypeString:
        return v.AsString()
    case TypeNull:
        return "NULL"
    }
    return "?"
}

// Compare returns -1, 0, or 1
func (v Value) Compare(other Value) int {
    switch v.Type {
    case TypeInt:
        a, b := v.AsInt(), other.AsInt()
        if a < b { return -1 }
        if a > b { return 1 }
        return 0
    case TypeFloat:
        a, b := v.AsFloat(), other.AsFloat()
        if a < b { return -1 }
        if a > b { return 1 }
        return 0
    case TypeString:
        a, b := v.AsString(), other.AsString()
        if a < b { return -1 }
        if a > b { return 1 }
        return 0
    }
    return 0
}
```

---

## 6. What Comes Next

Over the next 10 chapters, we'll build VaultDB layer by layer:

| Chapter | What We Build |
|---------|--------------|
| 39 | Disk manager: read/write pages to files |
| 40 | B-Tree: the data structure powering every index |
| 41 | Write-Ahead Log: crash recovery |
| 42 | Buffer Pool: caching pages in RAM |
| 43 | SQL Parser: turning text into a query plan |
| 44 | Query Executor: running SELECT, INSERT, UPDATE, DELETE |
| 45 | Transaction Manager: MVCC for concurrent access |
| 46 | Wire Protocol: accepting connections from Go/Python clients |
| 47 | Performance: benchmarks and profiling |
| 48 | Final integration: complete working VaultDB |

---

## Summary

- Databases store data in fixed-size **pages** (4 KB or 8 KB). This enables predictable I/O and efficient B-Tree indexing.
- VaultDB's architecture: Wire Protocol → Query Engine → Transaction Manager → Buffer Pool → Storage Engine.
- Each layer has a single responsibility. Understanding each layer individually makes the whole system comprehensible.
- The slot page format lets rows be variable-length while maintaining stable row identifiers.

### Exercises

**Easy:** Write a Go program that creates a `Page`, initializes it as a leaf page, and prints its header values. Verify that `FreeSpace()` returns exactly `PageSize - headerSize` for a fresh page.

**Medium:** Implement `InsertSlot(data []byte) (slotID int, err error)` on `Page`: allocate space from the bottom, add a slot entry at the top, and return the slot number. Handle the case where there's not enough free space.

**Hard:** Implement `GetSlotData(slotID int) []byte` and `DeleteSlot(slotID int)` on `Page`. Deletion should mark the slot as empty (set length to 0) without moving other data. Implement `Compact()` to reclaim space from deleted slots.
