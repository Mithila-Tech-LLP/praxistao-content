# Chapter 42: VaultDB — Buffer Pool: Memory Management

Reading a page from disk takes ~8 milliseconds. Reading it from RAM takes ~100 nanoseconds — that's 80,000x faster. The buffer pool is a cache of disk pages in memory. Get this right and your database is fast. Get it wrong and every query hits disk.

## Table of Contents

1. What the Buffer Pool Does
2. Page Table and Frame Layout
3. LRU Replacement Policy
4. Clock-Sweep (PostgreSQL's Algorithm)
5. Pinning and Latching
6. Complete Buffer Pool Implementation
7. Exercises

---

## 1. What the Buffer Pool Does

The buffer pool holds N fixed-size "frames" in memory. Each frame can hold one page:

```
Buffer Pool (holds 8 pages):
Frame 0: page_42   [pinned=0, dirty=false, usage=3]
Frame 1: page_7    [pinned=1, dirty=true,  usage=2]
Frame 2: page_15   [pinned=0, dirty=false, usage=1]
Frame 3: page_99   [pinned=0, dirty=false, usage=0]  ← evict this
Frame 4: [empty]
...
```

**Pin count:** A page with `pinned > 0` cannot be evicted (someone is currently using it).
**Dirty flag:** Set to true when a page is modified. Dirty pages must be flushed to disk before eviction.
**Usage count:** How recently the page was accessed. Used by the replacement policy.

---

## 2. Page Table and Frame Layout

```go
// storage/buffer.go
package storage

import (
    "fmt"
    "sync"
)

const DefaultBufferSize = 1024 // 1024 frames = 4 MB of buffer pool

type FrameID int

type frame struct {
    page     *Page
    pageID   PageID
    pinCount int
    dirty    bool
    usage    uint8  // for clock-sweep: 0 or 1
}

type BufferPool struct {
    mu      sync.Mutex
    dm      *DiskManager
    frames  []frame
    pageMap map[PageID]FrameID  // pageID → frameID
    clock   int                  // clock hand for replacement
}

func NewBufferPool(dm *DiskManager, size int) *BufferPool {
    return &BufferPool{
        dm:      dm,
        frames:  make([]frame, size),
        pageMap: make(map[PageID]FrameID),
    }
}
```

---

## 3. LRU Replacement Policy

LRU (Least Recently Used) evicts the page that hasn't been used for the longest time.

**Problem with pure LRU:** Expensive to implement efficiently (requires a doubly-linked list and a hash map, updated on every access). More importantly, a sequential scan reads every page once — LRU evicts hot pages to make room for cold ones!

**Example of LRU pollution:**
```
Buffer: [A, B, C, D] (A is hottest)
Sequential scan reads: E, F, G, H  (each accessed once)
After scan: [E, F, G, H] (A, B, C, D evicted!)
Next query accesses A → cache miss — must reload from disk
```

PostgreSQL's fix: **clock-sweep with buffer rings for sequential scans**.

---

## 4. Clock-Sweep (PostgreSQL's Algorithm)

Clock-sweep is O(1) and handles sequential scan pollution:

```
Frames arranged in a circle:
[A: usage=1] → [B: usage=1] → [C: usage=0] → [D: usage=1]
                                    ↑ clock hand

Clock hand at C:
- C has usage=0 → evict C, replace with new page
- Move clock hand to D

If clock hand hits a page with usage=1:
- Set usage to 0 (second chance)
- Move to next frame
- If it comes around again with usage=0 → evict
```

Each page access sets `usage = 1`. The clock sweeps around, giving pages a "second chance" before eviction.

```go
// findVictim returns a frame to evict (caller must hold the lock)
func (bp *BufferPool) findVictim() (FrameID, error) {
    n := len(bp.frames)
    // We'll try at most 2 full rotations
    for attempt := 0; attempt < 2*n; attempt++ {
        f := &bp.frames[bp.clock]
        bp.clock = (bp.clock + 1) % n

        if f.pinCount > 0 {
            continue // can't evict pinned frames
        }

        if f.usage > 0 {
            f.usage-- // second chance
            continue
        }

        // usage == 0 → this is our victim
        return FrameID(bp.clock - 1 + n) % FrameID(n), nil
    }
    return -1, fmt.Errorf("buffer pool exhausted: all frames are pinned")
}
```

---

## 5. Pinning and Latching

**Pin:** Prevent a page from being evicted while you're using it. Increment pin count when you start using a page, decrement when done.

**Latch:** A mutex protecting the frame's content from concurrent modification. Not the same as a lock (which is for transaction isolation — covered in the next chapter).

```go
// FetchPage returns a page, loading from disk if necessary.
// The caller must call UnpinPage when done.
func (bp *BufferPool) FetchPage(pageID PageID) (*Page, error) {
    bp.mu.Lock()
    defer bp.mu.Unlock()

    // Check if already in buffer pool
    if fid, ok := bp.pageMap[pageID]; ok {
        f := &bp.frames[fid]
        f.pinCount++
        f.usage = 1
        return f.page, nil
    }

    // Need to load from disk — find a victim frame
    fid, err := bp.findVictim()
    if err != nil {
        return nil, err
    }

    f := &bp.frames[fid]

    // Evict the current occupant
    if f.page != nil {
        if f.dirty {
            if err := bp.dm.WritePage(f.pageID, f.page); err != nil {
                return nil, fmt.Errorf("evict dirty page %d: %w", f.pageID, err)
            }
        }
        delete(bp.pageMap, f.pageID)
    }

    // Load new page
    page, err := bp.dm.ReadPage(pageID)
    if err != nil {
        return nil, fmt.Errorf("fetch page %d: %w", pageID, err)
    }

    f.page = page
    f.pageID = pageID
    f.pinCount = 1
    f.dirty = false
    f.usage = 1

    bp.pageMap[pageID] = fid
    return page, nil
}

// UnpinPage decrements the pin count. Set dirty=true if you modified the page.
func (bp *BufferPool) UnpinPage(pageID PageID, dirty bool) {
    bp.mu.Lock()
    defer bp.mu.Unlock()

    fid, ok := bp.pageMap[pageID]
    if !ok {
        return
    }

    f := &bp.frames[fid]
    if f.pinCount > 0 {
        f.pinCount--
    }
    if dirty {
        f.dirty = true
    }
}

// NewPage allocates a new page on disk and loads it into the buffer pool.
func (bp *BufferPool) NewPage() (PageID, *Page, error) {
    pageID, err := bp.dm.AllocatePage()
    if err != nil {
        return InvalidPageID, nil, err
    }
    page, err := bp.FetchPage(pageID)
    if err != nil {
        return InvalidPageID, nil, err
    }
    return pageID, page, nil
}

// FlushAll writes all dirty pages to disk.
func (bp *BufferPool) FlushAll() error {
    bp.mu.Lock()
    defer bp.mu.Unlock()

    for i := range bp.frames {
        f := &bp.frames[i]
        if f.page != nil && f.dirty {
            if err := bp.dm.WritePage(f.pageID, f.page); err != nil {
                return err
            }
            f.dirty = false
        }
    }
    return bp.dm.Sync()
}

// FlushPage writes a specific dirty page to disk immediately.
func (bp *BufferPool) FlushPage(pageID PageID) error {
    bp.mu.Lock()
    defer bp.mu.Unlock()

    fid, ok := bp.pageMap[pageID]
    if !ok {
        return nil // not in pool, assume clean
    }
    f := &bp.frames[fid]
    if !f.dirty {
        return nil
    }
    if err := bp.dm.WritePage(f.pageID, f.page); err != nil {
        return err
    }
    f.dirty = false
    return nil
}
```

---

## 6. Complete Buffer Pool Implementation

Here's how the layers connect:

```go
// main.go excerpt — how VaultDB uses the layers together
package main

import (
    "fmt"
    "log"
    "github.com/yourname/vaultdb/storage"
    "github.com/yourname/vaultdb/wal"
)

type Database struct {
    dm  *storage.DiskManager
    bp  *storage.BufferPool
    wal *wal.WAL
}

func Open(dbPath string) (*Database, error) {
    dm, err := storage.NewDiskManager(dbPath)
    if err != nil {
        return nil, err
    }

    bp := storage.NewBufferPool(dm, storage.DefaultBufferSize)

    w, err := wal.Open(dbPath + ".wal")
    if err != nil {
        return nil, err
    }

    db := &Database{dm: dm, bp: bp, wal: w}

    // Replay WAL on startup
    if err := db.recover(); err != nil {
        return nil, fmt.Errorf("recovery failed: %w", err)
    }

    return db, nil
}

func (db *Database) recover() error {
    toRedo, _, err := db.wal.Recover()
    if err != nil {
        return err
    }

    for _, rec := range toRedo {
        pageID := storage.PageID(rec.PageID)
        page, err := db.bp.FetchPage(pageID)
        if err != nil {
            return err
        }

        // Apply the change
        switch rec.Type {
        case wal.RecordInsert:
            // Write NewData to slot rec.SlotID
            writeSlot(page, int(rec.SlotID), rec.NewData)
        case wal.RecordDelete:
            // Zero out slot rec.SlotID
            clearSlot(page, int(rec.SlotID))
        case wal.RecordUpdate:
            writeSlot(page, int(rec.SlotID), rec.NewData)
        }

        db.bp.UnpinPage(pageID, true)
    }

    return db.bp.FlushAll()
}

func (db *Database) Close() error {
    db.bp.FlushAll()
    db.wal.Flush()
    return db.dm.Close()
}

func writeSlot(p *storage.Page, slotID int, data []byte) {
    // Implementation from storage package
    _ = p; _ = slotID; _ = data
}

func clearSlot(p *storage.Page, slotID int) {
    _ = p; _ = slotID
}
```

**Buffer pool statistics (for monitoring):**

```go
type BufferStats struct {
    Hits   int64 // page found in buffer pool
    Misses int64 // page not found, loaded from disk
    Dirty  int   // number of dirty pages
    Pinned int   // number of pinned pages
}

func (bp *BufferPool) Stats() BufferStats {
    bp.mu.Lock()
    defer bp.mu.Unlock()
    var stats BufferStats
    for _, f := range bp.frames {
        if f.dirty {
            stats.Dirty++
        }
        if f.pinCount > 0 {
            stats.Pinned++
        }
    }
    return stats
}
```

---

## Summary

- The buffer pool holds pages in memory to avoid repeated disk reads.
- Clock-sweep is PostgreSQL's eviction algorithm: pages get a "second chance" before eviction. Runs in O(1).
- Always `UnpinPage(dirty=true)` after modifying a page — this marks it for flushing.
- Never hold two pinned pages simultaneously when both might need eviction (deadlock risk).
- `FlushAll()` before checkpointing. `FlushAll()` + `WAL.Flush()` before shutdown.

### Exercises

**Easy:** Write a test with a buffer pool of size 3. FetchPage for pages 1, 2, 3. Then FetchPage for page 4 — page 1 (or the least recently used) should be evicted. Verify by checking that FetchPage for page 1 causes a disk read.

**Medium:** Implement `PrefetchPages(ids []PageID)` that launches goroutines to load multiple pages from disk concurrently. Pages are placed in available frames. If no frames are available, block until one is freed.

**Hard:** Implement a "buffer ring" for sequential scans: a dedicated pool of 32 frames that sequential scans use instead of the main pool. This prevents scan-induced cache pollution (the LRU problem described in section 3).
