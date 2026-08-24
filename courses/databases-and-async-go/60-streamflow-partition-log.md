# Chapter 60: StreamFlow — The Partition Log

A partition log manages multiple segment files for one partition. When a segment grows too large, a new one is created. Old segments past the retention period are deleted. This is exactly how Kafka's log manager works.

## Table of Contents

1. Why Multiple Segments?
2. The Partition Log
3. Segment Index for Fast Lookup
4. Log Retention
5. Exercises

---

## 1. Why Multiple Segments?

If all messages for a partition were in one file, it would grow forever. Deleting old messages would require reading and rewriting the whole file — expensive.

Instead, Kafka and StreamFlow split the partition into **segments**:

```
Partition "orders/0":
  00000000000000000000.log  (offset 0-999,  1 MB)
  00000000000000001000.log  (offset 1000-1999, 1 MB)
  00000000000000002000.log  (offset 2000+, active)
                               ↑ current segment (still writing)
```

Each segment file is named by its base offset. The active segment grows until it hits the max size (default: 1 GB in Kafka, we'll use 10 MB). Then a new segment is created.

**To delete old messages:** just delete old segment files. Fast, simple, no rewriting needed.

---

## 2. The Partition Log

```go
// log/log.go
package log

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    DefaultMaxSegmentBytes = 10 * 1024 * 1024 // 10 MB
    DefaultRetention       = 7 * 24 * time.Hour
)

type Config struct {
    MaxSegmentBytes int64
    Retention       time.Duration
}

// Log manages all segments for one partition
type Log struct {
    mu       sync.RWMutex
    dir      string
    cfg      Config
    segments []*Segment
    active   *Segment // current write segment
}

func NewLog(dir string, cfg Config) (*Log, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    if cfg.MaxSegmentBytes == 0 {
        cfg.MaxSegmentBytes = DefaultMaxSegmentBytes
    }
    if cfg.Retention == 0 {
        cfg.Retention = DefaultRetention
    }

    log := &Log{dir: dir, cfg: cfg}

    // Open existing segments
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    var baseOffsets []int64
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
            continue
        }
        name := strings.TrimSuffix(entry.Name(), ".log")
        offset, err := strconv.ParseInt(name, 10, 64)
        if err != nil {
            continue
        }
        baseOffsets = append(baseOffsets, offset)
    }
    sort.Slice(baseOffsets, func(i, j int) bool { return baseOffsets[i] < baseOffsets[j] })

    for _, base := range baseOffsets {
        seg, err := OpenSegment(log.segmentPath(base), base)
        if err != nil {
            return nil, fmt.Errorf("open segment %d: %w", base, err)
        }
        log.segments = append(log.segments, seg)
    }

    // Create initial segment if none exist
    if len(log.segments) == 0 {
        seg, err := OpenSegment(log.segmentPath(0), 0)
        if err != nil {
            return nil, err
        }
        log.segments = []*Segment{seg}
    }

    log.active = log.segments[len(log.segments)-1]
    return log, nil
}

func (l *Log) segmentPath(base int64) string {
    return filepath.Join(l.dir, fmt.Sprintf("%020d.log", base))
}

// Append adds a message to the active segment
func (l *Log) Append(key, value []byte) (int64, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    // Roll to new segment if active is too large
    if l.active.Size() >= l.cfg.MaxSegmentBytes {
        if err := l.active.Sync(); err != nil {
            return 0, err
        }
        newBase := l.active.NextOffset()
        seg, err := OpenSegment(l.segmentPath(newBase), newBase)
        if err != nil {
            return 0, fmt.Errorf("create segment: %w", err)
        }
        l.segments = append(l.segments, seg)
        l.active = seg
    }

    return l.active.Append(key, value, time.Now().UnixNano())
}

// Read returns messages starting from offset, up to maxBytes
func (l *Log) Read(offset int64, maxBytes int32) ([]FetchMessage, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()

    // Find the segment containing this offset
    seg := l.findSegment(offset)
    if seg == nil {
        return nil, fmt.Errorf("offset %d not found", offset)
    }

    return seg.Read(offset, maxBytes)
}

// findSegment finds the segment that contains the given offset
// The segment containing offset X is the last one whose baseOffset <= X
func (l *Log) findSegment(offset int64) *Segment {
    // Binary search for the right segment
    idx := sort.Search(len(l.segments), func(i int) bool {
        return l.segments[i].BaseOffset() > offset
    }) - 1

    if idx < 0 {
        return nil
    }
    return l.segments[idx]
}

// LowestOffset returns the offset of the first available message
func (l *Log) LowestOffset() int64 {
    l.mu.RLock()
    defer l.mu.RUnlock()
    if len(l.segments) == 0 {
        return 0
    }
    return l.segments[0].BaseOffset()
}

// HighestOffset returns the offset that will be assigned to the next message
func (l *Log) HighestOffset() int64 {
    l.mu.RLock()
    defer l.mu.RUnlock()
    return l.active.NextOffset()
}

// Truncate removes all segments and creates a fresh log from startOffset
func (l *Log) Truncate(lowestOffset int64) error {
    l.mu.Lock()
    defer l.mu.Unlock()

    var toRemove []*Segment
    var toKeep []*Segment

    for _, seg := range l.segments {
        if seg.NextOffset() <= lowestOffset {
            toRemove = append(toRemove, seg)
        } else {
            toKeep = append(toKeep, seg)
        }
    }

    for _, seg := range toRemove {
        seg.Close()
        os.Remove(seg.path)
    }

    l.segments = toKeep
    return nil
}

// ApplyRetention deletes segments older than the retention period
func (l *Log) ApplyRetention() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    cutoff := time.Now().Add(-l.cfg.Retention).UnixNano()

    var toRemove []*Segment
    var toKeep []*Segment

    for _, seg := range l.segments {
        if seg == l.active {
            toKeep = append(toKeep, seg)
            continue
        }

        // Get last message time from the segment
        lastMsg := seg.lastTimestamp()
        if lastMsg > 0 && lastMsg < cutoff {
            toRemove = append(toRemove, seg)
        } else {
            toKeep = append(toKeep, seg)
        }
    }

    for _, seg := range toRemove {
        seg.Close()
        os.Remove(seg.path)
    }

    l.segments = toKeep
    return nil
}

// Close flushes and closes all segments
func (l *Log) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    for _, seg := range l.segments {
        if err := seg.Close(); err != nil {
            return err
        }
    }
    return nil
}

// Stats returns current log statistics
type LogStats struct {
    Segments    int
    TotalBytes  int64
    LowestOffset int64
    HighestOffset int64
}

func (l *Log) Stats() LogStats {
    l.mu.RLock()
    defer l.mu.RUnlock()
    var total int64
    for _, seg := range l.segments {
        total += seg.Size()
    }
    lo := int64(0)
    if len(l.segments) > 0 {
        lo = l.segments[0].BaseOffset()
    }
    return LogStats{
        Segments:     len(l.segments),
        TotalBytes:   total,
        LowestOffset: lo,
        HighestOffset: l.active.NextOffset(),
    }
}
```

Add `lastTimestamp` to Segment (reads the last message's timestamp):

```go
// In log/segment.go — add this method
func (s *Segment) lastTimestamp() int64 {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if s.size < int64(msgHeaderSize) {
        return 0
    }

    // Seek to near the end and scan backwards isn't efficient
    // In production: maintain a separate index file
    // For simplicity: re-scan from the beginning and return the last ts
    f, err := os.Open(s.path)
    if err != nil {
        return 0
    }
    defer f.Close()

    var lastTS int64
    for {
        var header [msgHeaderSize]byte
        if _, err := io.ReadFull(f, header[:]); err != nil {
            break
        }
        lastTS = int64(binary.BigEndian.Uint64(header[8:16]))
        keyLen := int(binary.BigEndian.Uint16(header[16:18]))
        valLen := int(binary.BigEndian.Uint32(header[18:22]))
        f.Seek(int64(keyLen+valLen), io.SeekCurrent)
    }
    return lastTS
}
```

---

## 3. Segment Index for Fast Lookup

Without an index, finding messages at offset X requires scanning from the beginning. With an index, we jump directly to the right byte position.

```go
// Sparse index: one entry every N messages
// Index file format:
//   8 bytes: message offset
//   8 bytes: byte position in segment

type Index struct {
    mu      sync.RWMutex
    file    *os.File
    entries []indexEntry
}

type indexEntry struct {
    offset   int64
    position int64
}

const indexInterval = 100 // one entry every 100 messages

func (s *Segment) writeIndexEntry(offset, position int64) error {
    if s.index == nil {
        return nil
    }
    buf := make([]byte, 16)
    binary.BigEndian.PutUint64(buf[0:8], uint64(offset))
    binary.BigEndian.PutUint64(buf[8:16], uint64(position))
    _, err := s.index.file.Write(buf)
    return err
}

// FindBytePosition returns the byte position in the segment for a given offset
// Uses the index for O(log n) lookup
func (s *Segment) FindBytePosition(targetOffset int64) (int64, error) {
    if s.index == nil || len(s.index.entries) == 0 {
        return 0, nil // scan from beginning
    }

    // Binary search for the largest indexed offset <= targetOffset
    idx := sort.Search(len(s.index.entries), func(i int) bool {
        return s.index.entries[i].offset > targetOffset
    }) - 1

    if idx < 0 {
        return 0, nil
    }
    return s.index.entries[idx].position, nil
}
```

---

## Summary

- A partition log manages multiple segment files, each named by its base offset.
- When a segment grows past `MaxSegmentBytes`, a new segment is created (segment roll).
- Finding the right segment for an offset: binary search for the largest base offset ≤ target.
- Log retention: delete segments whose last message timestamp is past the retention window.
- A sparse index maps offsets to byte positions, enabling O(log n) seeks instead of O(n) scans.

### Exercises

**Easy:** Write a test that creates a Log, appends 10,000 messages, reads messages from offset 5000, and verifies the first returned message has offset exactly 5000.

**Medium:** Implement the retention policy. Set `MaxSegmentBytes = 1024` (tiny, to trigger segment rolls). Append 10,000 messages, set retention to 1 millisecond, call `ApplyRetention()`, verify that only the active segment remains.

**Hard:** Implement the sparse index fully: every 100 messages, write an `(offset, bytePosition)` entry to a `.index` file alongside the `.log` file. On `Read`, use the index to seek to the right byte position instead of scanning from the beginning. Benchmark the difference.
