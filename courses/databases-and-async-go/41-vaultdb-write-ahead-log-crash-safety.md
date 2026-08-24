# Chapter 41: VaultDB — Write-Ahead Log (WAL) for Crash Safety

Your database is doing thousands of writes. The power goes out. When you restart, is your data intact? Without a WAL, probably not. With a WAL, always yes.

## Table of Contents

1. The Problem: Crashes at the Wrong Moment
2. The WAL Principle
3. WAL Record Format
4. Writing to the WAL
5. Recovery on Startup
6. Checkpointing
7. Exercises

---

## 1. The Problem: Crashes at the Wrong Moment

Imagine writing a row that spans two pages:
1. Write page 1 to disk ✓
2. **POWER OUTAGE**
3. Page 2 never written ✗

Your data is now partially written. The row is split: half on disk, half lost. This is called a **torn write** or **partial write**. It leaves the database in an inconsistent state.

**Even worse:** An UPDATE might:
1. Modify page A (new value)
2. **CRASH**
3. Page A has new value but the old value is gone

If you wrote page A but crashed before "committing" the change, you have data the user never confirmed.

---

## 2. The WAL Principle

**"Write the log before writing the data."**

Before changing any page on disk, first append a log record describing the change. If we crash mid-way:
- The log record exists → we can redo the change (or undo it if the transaction wasn't committed)
- The page change may or may not be on disk → doesn't matter, log tells us what to do

On startup, VaultDB:
1. Reads the WAL from start to end
2. **Redo** all committed transactions (apply their changes if not already on disk)
3. **Undo** all uncommitted transactions (roll back partial changes)

This guarantees that after restart, the database is in exactly the state it was in at the last committed transaction.

```
Timeline:
T1 BEGIN
T1: update row → WAL: UPDATE page=5 slot=3 old=... new=...
T1: insert row → WAL: INSERT page=6 slot=0 new=...
T1 COMMIT       → WAL: COMMIT txnID=1
                         ↑ ONLY NOW do we guarantee durability
Pages may still be dirty (in memory) but log is on disk (fsync)
```

---

## 3. WAL Record Format

```go
// wal/wal.go
package wal

import (
    "encoding/binary"
    "fmt"
    "io"
    "os"
    "sync"
)

type RecordType uint8

const (
    RecordInsert   RecordType = 1
    RecordUpdate   RecordType = 2
    RecordDelete   RecordType = 3
    RecordCommit   RecordType = 4
    RecordAbort    RecordType = 5
    RecordCheckpoint RecordType = 6
)

// LSN (Log Sequence Number) is a unique, monotonically increasing ID for each WAL record.
// It's the byte offset in the WAL file.
type LSN uint64

const InvalidLSN LSN = ^LSN(0)

// Record is one WAL entry.
// Binary format:
//   4 bytes: total record length (including header)
//   8 bytes: transaction ID
//   1 byte:  record type
//   8 bytes: page ID
//   2 bytes: slot ID
//   N bytes: payload (old data for update, new data for insert/update)
//   4 bytes: CRC32 checksum (to detect corrupt records)

type Record struct {
    TxnID   uint64
    Type    RecordType
    PageID  uint64
    SlotID  uint16
    OldData []byte // for UPDATE: the data before change
    NewData []byte // for INSERT/UPDATE: the new data
}

const recordHeaderSize = 4 + 8 + 1 + 8 + 2 // len + txnid + type + pageid + slotid
```

---

## 4. Writing to the WAL

```go
type WAL struct {
    mu      sync.Mutex
    file    *os.File
    offset  int64 // current write position
}

func Open(filename string) (*WAL, error) {
    f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }
    info, _ := f.Stat()
    return &WAL{file: f, offset: info.Size()}, nil
}

// Append writes a record to the WAL and returns its LSN (byte offset)
func (w *WAL) Append(rec Record) (LSN, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    data := encodeRecord(rec)

    lsn := LSN(w.offset)
    n, err := w.file.Write(data)
    if err != nil {
        return InvalidLSN, fmt.Errorf("wal write: %w", err)
    }
    w.offset += int64(n)
    return lsn, nil
}

// Flush ensures all written records are on disk (fsync).
// Must be called before returning "committed" to the client.
func (w *WAL) Flush() error {
    return w.file.Sync()
}

func encodeRecord(rec Record) []byte {
    oldLen := len(rec.OldData)
    newLen := len(rec.NewData)
    // total = header + 4(oldLen) + oldData + 4(newLen) + newData + 4(crc)
    total := recordHeaderSize + 4 + oldLen + 4 + newLen + 4

    buf := make([]byte, total)
    off := 0

    binary.BigEndian.PutUint32(buf[off:], uint32(total))
    off += 4

    binary.BigEndian.PutUint64(buf[off:], rec.TxnID)
    off += 8

    buf[off] = byte(rec.Type)
    off++

    binary.BigEndian.PutUint64(buf[off:], rec.PageID)
    off += 8

    binary.BigEndian.PutUint16(buf[off:], rec.SlotID)
    off += 2

    binary.BigEndian.PutUint32(buf[off:], uint32(oldLen))
    off += 4
    copy(buf[off:], rec.OldData)
    off += oldLen

    binary.BigEndian.PutUint32(buf[off:], uint32(newLen))
    off += 4
    copy(buf[off:], rec.NewData)
    off += newLen

    // CRC32 of everything before the checksum
    crc := crc32Sum(buf[:off])
    binary.BigEndian.PutUint32(buf[off:], crc)

    return buf
}

func crc32Sum(data []byte) uint32 {
    // Simple FNV-1a as placeholder (use hash/crc32 in production)
    h := uint32(2166136261)
    for _, b := range data {
        h ^= uint32(b)
        h *= 16777619
    }
    return h
}
```

---

## 5. Recovery on Startup

```go
// Recover reads the WAL from the beginning and returns:
// - committed: all records for committed transactions (to REDO)
// - aborted: txn IDs that should be undone
func (w *WAL) Recover() ([]Record, []uint64, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    // Seek to start
    if _, err := w.file.Seek(0, io.SeekStart); err != nil {
        return nil, nil, err
    }

    var all []Record
    committed := make(map[uint64]bool)
    aborted := make(map[uint64]bool)

    for {
        rec, err := readNextRecord(w.file)
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("WAL: stopping recovery at corrupt record: %v\n", err)
            break // stop at first corrupt record
        }

        switch rec.Type {
        case RecordCommit:
            committed[rec.TxnID] = true
        case RecordAbort:
            aborted[rec.TxnID] = true
        default:
            all = append(all, rec)
        }
    }

    // Filter: only keep records from committed transactions
    var toRedo []Record
    var toAbort []uint64

    txnSeen := make(map[uint64]bool)
    for _, rec := range all {
        if committed[rec.TxnID] {
            toRedo = append(toRedo, rec)
        } else if !aborted[rec.TxnID] {
            // Transaction was in progress when crash happened
            if !txnSeen[rec.TxnID] {
                txnSeen[rec.TxnID] = true
                toAbort = append(toAbort, rec.TxnID)
            }
        }
    }

    return toRedo, toAbort, nil
}

func readNextRecord(r io.Reader) (Record, error) {
    // Read 4-byte length
    var lenBuf [4]byte
    if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
        if err == io.EOF {
            return Record{}, io.EOF
        }
        return Record{}, fmt.Errorf("read record length: %w", err)
    }

    total := int(binary.BigEndian.Uint32(lenBuf[:]))
    if total < recordHeaderSize+8 {
        return Record{}, fmt.Errorf("record too short: %d bytes", total)
    }

    // Read the rest
    rest := make([]byte, total-4)
    if _, err := io.ReadFull(r, rest); err != nil {
        return Record{}, fmt.Errorf("read record body: %w", err)
    }

    off := 0
    txnID := binary.BigEndian.Uint64(rest[off:])
    off += 8
    recType := RecordType(rest[off])
    off++
    pageID := binary.BigEndian.Uint64(rest[off:])
    off += 8
    slotID := binary.BigEndian.Uint16(rest[off:])
    off += 2

    oldLen := int(binary.BigEndian.Uint32(rest[off:]))
    off += 4
    oldData := make([]byte, oldLen)
    copy(oldData, rest[off:off+oldLen])
    off += oldLen

    newLen := int(binary.BigEndian.Uint32(rest[off:]))
    off += 4
    newData := make([]byte, newLen)
    copy(newData, rest[off:off+newLen])

    return Record{
        TxnID:   txnID,
        Type:    recType,
        PageID:  pageID,
        SlotID:  slotID,
        OldData: oldData,
        NewData: newData,
    }, nil
}
```

---

## 6. Checkpointing

The WAL grows forever if we never truncate it. **Checkpointing** flushes all dirty pages to disk, then truncates the WAL.

After a checkpoint, recovery only needs to replay WAL records written after the checkpoint — not from the beginning of time.

```go
// Checkpoint flushes dirty pages and records the checkpoint LSN
func (w *WAL) Checkpoint(flushDirtyPages func() error) error {
    // 1. Wait for all active transactions to complete
    // 2. Flush all dirty pages from buffer pool to disk
    if err := flushDirtyPages(); err != nil {
        return fmt.Errorf("flush pages: %w", err)
    }

    // 3. Write a CHECKPOINT record to the WAL
    _, err := w.Append(Record{Type: RecordCheckpoint})
    if err != nil {
        return err
    }

    // 4. Sync the WAL
    if err := w.Flush(); err != nil {
        return err
    }

    // 5. Truncate WAL up to this point
    // (In production: keep WAL, just move the "recovery start" pointer)
    return nil
}
```

**The durability guarantee:**
- `WAL.Append(CommitRecord)` + `WAL.Flush()` = transaction is durable, even if we crash immediately after.
- On restart: read WAL, redo committed transactions, undo uncommitted ones.
- **No data loss for committed transactions, ever.**

---

## Summary

- The WAL is an append-only log of all changes to the database.
- Write to WAL first, then update the page in the buffer pool.
- Before returning "committed" to the client, call `WAL.Flush()` (fsync) to guarantee the record is on disk.
- On crash recovery: redo all committed transactions, undo all uncommitted ones.
- Checkpoints flush dirty pages so recovery can start from a recent point instead of the beginning of the WAL.

### Exercises

**Easy:** Write a Go test that: (1) opens a WAL, (2) appends 5 INSERT records, (3) appends a COMMIT for txn 1, (4) appends 3 more records for txn 2 (no commit), (5) calls `Recover()`, (6) verifies that only txn 1's records are returned for redo, and txn 2's ID is in the to-undo list.

**Medium:** Implement `Apply(rec Record, dm *storage.DiskManager) error` that reads the target page, applies the change (for INSERT: write NewData to the slot, for DELETE: zero out the slot), and writes the page back. This is the "redo" step in recovery.

**Hard:** Implement a full recovery integration test: (1) create a database, (2) insert 10 rows in a transaction, (3) commit, (4) "simulate a crash" by killing the process without flushing pages, (5) re-open the database (WAL replay happens in constructor), (6) verify all 10 rows are visible via a scan.
