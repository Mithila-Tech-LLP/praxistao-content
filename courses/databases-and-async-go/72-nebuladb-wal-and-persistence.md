# Chapter 72: NebulaDB — WAL, Vector Storage, and Crash Recovery

NebulaDB must survive crashes. If the process dies mid-write, data should not be corrupted or silently lost. This chapter builds the Write-Ahead Log and the vector storage layer, then implements crash recovery.

## Table of Contents

1. The Write-Ahead Log
2. Vector Storage — Memory-Mapped Files
3. Snapshots
4. Crash Recovery — Replaying the WAL
5. The Complete Collection Lifecycle
6. Exercises

---

## 1. The Write-Ahead Log

The WAL is the most important durability primitive. Every write is appended here before being applied — just like Qdrant, PostgreSQL, and VaultDB.

**WAL record format:**

```
[4 bytes: magic] [4 bytes: operation] [8 bytes: data length] [N bytes: data] [4 bytes: CRC32]
```

```go
// storage/wal.go
package storage

import (
    "bufio"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "hash/crc32"
    "io"
    "os"
    "sync"

    "nebuladb/types"
)

const (
    walMagic     = 0xNEBU1A00 // magic number for integrity checking
    opUpsert     = 1
    opDelete     = 2
    recordHeaderSize = 4 + 4 + 8 // magic + optype + datalen
)

// WALRecord is a decoded WAL entry
type WALRecord struct {
    Op      int
    Payload []byte // raw JSON for the operation
}

// WAL is an append-only write-ahead log
type WAL struct {
    f   *os.File
    w   *bufio.Writer
    mu  sync.Mutex
    pos int64 // current write position (for checkpointing)
}

func NewWAL(path string) (*WAL, error) {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    info, err := f.Stat()
    if err != nil {
        f.Close()
        return nil, err
    }

    return &WAL{
        f:   f,
        w:   bufio.NewWriterSize(f, 64*1024),
        pos: info.Size(),
    }, nil
}

// WriteUpsert appends an upsert record to the WAL
func (w *WAL) WriteUpsert(point types.Point) error {
    data, err := json.Marshal(point)
    if err != nil {
        return err
    }
    return w.write(opUpsert, data)
}

// WriteDelete appends a delete record to the WAL
func (w *WAL) WriteDelete(id uint64) error {
    data, err := json.Marshal(map[string]uint64{"id": id})
    if err != nil {
        return err
    }
    return w.write(opDelete, data)
}

func (w *WAL) write(op int, data []byte) error {
    w.mu.Lock()
    defer w.mu.Unlock()

    // Calculate CRC32 over op + data
    checksum := crc32.ChecksumIEEE(append([]byte{byte(op)}, data...))

    // Write: magic(4) + op(4) + datalen(8) + data(N) + checksum(4)
    header := make([]byte, recordHeaderSize)
    binary.LittleEndian.PutUint32(header[0:4], walMagic)
    binary.LittleEndian.PutUint32(header[4:8], uint32(op))
    binary.LittleEndian.PutUint64(header[8:16], uint64(len(data)))

    if _, err := w.w.Write(header); err != nil {
        return err
    }
    if _, err := w.w.Write(data); err != nil {
        return err
    }

    checkBuf := make([]byte, 4)
    binary.LittleEndian.PutUint32(checkBuf, checksum)
    if _, err := w.w.Write(checkBuf); err != nil {
        return err
    }

    // fsync to guarantee durability before returning to caller
    if err := w.w.Flush(); err != nil {
        return err
    }
    return w.f.Sync()
}

// Replay reads all WAL records and calls the provided handler for each
func (w *WAL) Replay(handler func(WALRecord) error) error {
    // Seek to start for reading
    if _, err := w.f.Seek(0, io.SeekStart); err != nil {
        return err
    }

    r := bufio.NewReader(w.f)
    for {
        header := make([]byte, recordHeaderSize)
        if _, err := io.ReadFull(r, header); err != nil {
            if err == io.EOF || err == io.ErrUnexpectedEOF {
                break // normal end of log
            }
            return fmt.Errorf("read header: %w", err)
        }

        magic := binary.LittleEndian.Uint32(header[0:4])
        if magic != walMagic {
            return fmt.Errorf("WAL corruption: bad magic %x at position", magic)
        }

        op := int(binary.LittleEndian.Uint32(header[4:8]))
        dataLen := binary.LittleEndian.Uint64(header[8:16])

        data := make([]byte, dataLen)
        if _, err := io.ReadFull(r, data); err != nil {
            return fmt.Errorf("read data: %w", err)
        }

        checksumBuf := make([]byte, 4)
        if _, err := io.ReadFull(r, checksumBuf); err != nil {
            return fmt.Errorf("read checksum: %w", err)
        }

        // Verify CRC32
        expectedCRC := binary.LittleEndian.Uint32(checksumBuf)
        actualCRC := crc32.ChecksumIEEE(append([]byte{byte(op)}, data...))
        if expectedCRC != actualCRC {
            return fmt.Errorf("WAL corruption: checksum mismatch — data likely truncated during crash")
        }

        if err := handler(WALRecord{Op: op, Payload: data}); err != nil {
            return err
        }
    }
    return nil
}

// Truncate removes processed records after a checkpoint
func (w *WAL) Truncate() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    if err := w.f.Truncate(0); err != nil {
        return err
    }
    _, err := w.f.Seek(0, io.SeekStart)
    return err
}

func (w *WAL) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.w.Flush()
    return w.f.Close()
}
```

---

## 2. Vector Storage — Memory-Mapped Files

Vectors are the biggest data in the system. We store them in a flat binary file and access them via memory mapping — the OS handles caching automatically, just like Qdrant.

```go
// storage/vector_store.go
package storage

import (
    "encoding/binary"
    "fmt"
    "os"
    "sync"
)

// VectorStore stores float32 vectors in a flat file.
// Layout: [8 bytes: count] [8 bytes: dimension] [vectors...]
// Each vector: dimension * 4 bytes (float32, little-endian)
type VectorStore struct {
    f         *os.File
    dimension int
    index     map[uint64]int64 // id → byte offset in file
    mu        sync.RWMutex
}

const vectorFileHeader = 16 // count(8) + dimension(8)

func NewVectorStore(path string, dimension int) (*VectorStore, error) {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        return nil, err
    }

    vs := &VectorStore{
        f:         f,
        dimension: dimension,
        index:     make(map[uint64]int64),
    }

    info, _ := f.Stat()
    if info.Size() == 0 {
        // New file: write header
        header := make([]byte, vectorFileHeader)
        binary.LittleEndian.PutUint64(header[0:8], 0)
        binary.LittleEndian.PutUint64(header[8:16], uint64(dimension))
        f.Write(header)
    } else {
        // Existing file: build in-memory index by scanning
        vs.buildIndex()
    }

    return vs, nil
}

// Set stores a vector for a given ID (append-only; updates are re-appended)
func (vs *VectorStore) Set(id uint64, vector []float32) error {
    if len(vector) != vs.dimension {
        return fmt.Errorf("dimension mismatch: got %d, want %d", len(vector), vs.dimension)
    }

    vs.mu.Lock()
    defer vs.mu.Unlock()

    offset, err := vs.f.Seek(0, 2) // seek to end
    if err != nil {
        return err
    }

    // Write: [8 bytes: id] [dim * 4 bytes: float32 values]
    buf := make([]byte, 8+vs.dimension*4)
    binary.LittleEndian.PutUint64(buf[0:8], id)
    for i, v := range vector {
        bits := math_float32bits(v)
        binary.LittleEndian.PutUint32(buf[8+i*4:], bits)
    }

    if _, err := vs.f.Write(buf); err != nil {
        return err
    }

    vs.index[id] = offset
    return nil
}

// Get reads a vector by ID
func (vs *VectorStore) Get(id uint64) ([]float32, error) {
    vs.mu.RLock()
    offset, ok := vs.index[id]
    vs.mu.RUnlock()

    if !ok {
        return nil, fmt.Errorf("vector %d not found", id)
    }

    buf := make([]byte, 8+vs.dimension*4)
    if _, err := vs.f.ReadAt(buf, offset); err != nil {
        return nil, err
    }

    vector := make([]float32, vs.dimension)
    for i := range vector {
        bits := binary.LittleEndian.Uint32(buf[8+i*4:])
        vector[i] = math_float32frombits(bits)
    }
    return vector, nil
}

func (vs *VectorStore) buildIndex() {
    offset := int64(vectorFileHeader)
    recordSize := int64(8 + vs.dimension*4)
    buf := make([]byte, 8)

    for {
        if _, err := vs.f.ReadAt(buf, offset); err != nil {
            break
        }
        id := binary.LittleEndian.Uint64(buf)
        vs.index[id] = offset
        offset += recordSize
    }
}

func (vs *VectorStore) Close() error {
    return vs.f.Close()
}

// Pure Go float32 bit conversion (avoids importing math/unsafe)
func math_float32bits(f float32) uint32 {
    buf := [4]byte{}
    buf[0] = *(*byte)(unsafe.Pointer(&f))
    buf[1] = *(*byte)(unsafe.Add(unsafe.Pointer(&f), 1))
    buf[2] = *(*byte)(unsafe.Add(unsafe.Pointer(&f), 2))
    buf[3] = *(*byte)(unsafe.Add(unsafe.Pointer(&f), 3))
    return binary.LittleEndian.Uint32(buf[:])
}

func math_float32frombits(b uint32) float32 {
    buf := [4]byte{}
    binary.LittleEndian.PutUint32(buf[:], b)
    return *(*float32)(unsafe.Pointer(&buf[0]))
}
```

Actually, let's use `math` package for cleanliness:

```go
import "math"

func math_float32bits(f float32) uint32 {
    return math.Float32bits(f)
}

func math_float32frombits(b uint32) float32 {
    return math.Float32frombits(b)
}
```

---

## 3. Snapshots

A snapshot is a consistent backup of a collection — all vectors, payloads, HNSW index, and config.

```go
// collection/snapshot.go
package collection

import (
    "archive/tar"
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

// Snapshot creates a .tar.gz backup of the collection directory
func (c *Collection) Snapshot(destDir string) (string, error) {
    // First, flush everything to disk
    if err := c.hnswIndex.Save(filepath.Join(c.dir, "hnsw.bin")); err != nil {
        return "", fmt.Errorf("save hnsw: %w", err)
    }

    snapName := fmt.Sprintf("%s-%s.tar.gz", c.name, time.Now().Format("2006-01-02-150405"))
    snapPath := filepath.Join(destDir, snapName)

    out, err := os.Create(snapPath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    gw := gzip.NewWriter(out)
    defer gw.Close()
    tw := tar.NewWriter(gw)
    defer tw.Close()

    // Walk all files in the collection directory
    err = filepath.Walk(c.dir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }

        rel, _ := filepath.Rel(c.dir, path)
        header := &tar.Header{
            Name:    rel,
            Size:    info.Size(),
            Mode:    int64(info.Mode()),
            ModTime: info.ModTime(),
        }
        if err := tw.WriteHeader(header); err != nil {
            return err
        }

        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()
        _, err = io.Copy(tw, f)
        return err
    })

    if err != nil {
        os.Remove(snapPath)
        return "", err
    }

    return snapPath, nil
}

// RestoreSnapshot extracts a .tar.gz snapshot into a collection directory
func RestoreSnapshot(snapPath, targetDir string) error {
    f, err := os.Open(snapPath)
    if err != nil {
        return err
    }
    defer f.Close()

    gr, err := gzip.NewReader(f)
    if err != nil {
        return err
    }
    defer gr.Close()
    tr := tar.NewReader(gr)

    if err := os.MkdirAll(targetDir, 0755); err != nil {
        return err
    }

    for {
        header, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        destPath := filepath.Join(targetDir, header.Name)
        out, err := os.Create(destPath)
        if err != nil {
            return err
        }
        if _, err := io.Copy(out, tr); err != nil {
            out.Close()
            return err
        }
        out.Close()
    }
    return nil
}
```

---

## 4. Crash Recovery — Replaying the WAL

When NebulaDB starts, it replays any unprocessed WAL records:

```go
// collection/collection.go — loadCollection
func loadCollection(name, dir string) (*Collection, error) {
    cfg, err := loadConfig(dir)
    if err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }

    payloadStore, err := storage.NewPayloadStore(filepath.Join(dir, "payload.db"))
    if err != nil {
        return nil, err
    }

    vectorStore, err := storage.NewVectorStore(filepath.Join(dir, "vectors.bin"), cfg.Dimension)
    if err != nil {
        return nil, err
    }

    wal, err := storage.NewWAL(filepath.Join(dir, "wal.log"))
    if err != nil {
        return nil, err
    }

    // Load HNSW index from snapshot if it exists
    hnswPath := filepath.Join(dir, "hnsw.bin")
    var hnswIdx *hnsw.Index
    if _, err := os.Stat(hnswPath); err == nil {
        hnswIdx, err = hnsw.LoadIndex(hnswPath, cfg.Distance)
        if err != nil {
            return nil, fmt.Errorf("load hnsw: %w", err)
        }
    } else {
        hnswIdx = hnsw.NewIndex(cfg.Dimension, cfg.Distance, hnsw.HNSWConfig{
            M:              cfg.HNSW.M,
            EfConstruction: cfg.HNSW.EfConstruction,
            EfSearch:       cfg.HNSW.EfSearch,
        })
    }

    c := &Collection{
        name:         name,
        dir:          dir,
        config:       cfg,
        hnswIndex:    hnswIdx,
        vectorStore:  vectorStore,
        payloadStore: payloadStore,
        indexMgr:     index.NewManager(),
        wal:          wal,
    }

    // Replay WAL to recover any writes that happened after the last HNSW snapshot
    recovered := 0
    err = wal.Replay(func(record storage.WALRecord) error {
        switch record.Op {
        case storage.OpUpsert:
            var point types.Point
            if err := json.Unmarshal(record.Payload, &point); err != nil {
                return err
            }
            // Re-apply the write: vector store + payload store + HNSW
            // (idempotent: upsert with same id overwrites)
            vectorStore.Set(point.ID, point.Vector)
            payloadStore.Set(point.ID, point.Payload)
            hnswIdx.Insert(point.ID, point.Vector)
            c.indexMgr.Index(point.ID, point.Payload)
            c.count.Add(1)
            recovered++
        case storage.OpDelete:
            var rec struct{ ID uint64 }
            if err := json.Unmarshal(record.Payload, &rec); err != nil {
                return err
            }
            vectorStore.Delete(rec.ID)
            payloadStore.Delete(rec.ID)
            // Note: HNSW deletion is lazy — mark as deleted
            recovered++
        }
        return nil
    })
    if err != nil {
        return nil, fmt.Errorf("wal replay: %w", err)
    }

    if recovered > 0 {
        log.Printf("[%s] recovered %d operations from WAL", name, recovered)
    }

    return c, nil
}
```

**Recovery timeline:**

```
Normal operation:
  Write → WAL (fsynced) → VectorStore → PayloadStore → HNSW
  Periodic checkpoint → Save HNSW snapshot → Truncate WAL

Crash recovery:
  Load HNSW snapshot (last checkpoint)
  Replay WAL from beginning (re-apply writes since last checkpoint)
  Done — state is consistent
```

---

## 5. The Complete Collection Lifecycle

```
CollectionManager.Create("products", Config{Dimension: 1536, Distance: Cosine})
    → creates data/products/ directory
    → initializes PayloadStore (BoltDB)
    → initializes VectorStore (binary file)
    → initializes WAL (append-only log)
    → initializes HNSW index (empty)
    → saves config.json

Collection.Upsert(Point{ID: 1, Vector: [...], Payload: {...}})
    → WAL.WriteUpsert (fsynced)
    → VectorStore.Set(1, [...])
    → PayloadStore.Set(1, {...})
    → IndexMgr.Index(1, {...})
    → HNSWIndex.Insert(1, [...])

Collection.Search(SearchRequest{Vector: [...], Filter: {...}, Limit: 10})
    → IndexMgr.FilterIDs(filter)       // pre-filter using payload indexes
    → choose strategy (brute/filtered/post)
    → HNSWIndex.Search(vector, k, filterFn)
    → PayloadStore.GetMany(result IDs)
    → return ScoredPoints

Process crash → restart
    → CollectionManager.loadExisting()
    → loadCollection("products")
    → load HNSWIndex from hnsw.bin (last checkpoint)
    → WAL.Replay() → re-apply missed writes
    → ready to serve
```

---

## Summary

- The **WAL** appends a CRC32-protected record before every write. On crash, replay restores missed writes.
- **VectorStore** stores float32 arrays in a flat binary file. An in-memory `id → offset` map makes reads O(1).
- **Snapshots** are gzip'd tar archives of the collection directory — all files included in one atomic bundle.
- **Recovery** = load last HNSW snapshot + replay WAL. The HNSW snapshot is the checkpoint; the WAL is the tail.
- Write order matters: WAL first, then actual stores. If WAL write fails, nothing is applied. If stores fail after WAL write, recovery replays the WAL record.

### Exercises

**Easy:** Add a `Checkpoint()` method to `Collection` that saves the HNSW index to `hnsw.bin` and truncates the WAL. When should this be called? (Hint: periodic background goroutine.)

**Medium:** The VectorStore uses an append-only layout — updates append a new record with the same ID, and only the last occurrence matters. Over time this wastes space. Implement a `Compact()` method that rewrites the file keeping only the latest record per ID.

**Hard:** The current WAL writes one record per upsert — slow for bulk ingestion. Implement batch WAL writes: accumulate up to 1000 records in memory, then write them all in one `pwrite` + single `fsync`. Measure the throughput improvement on 100k upserts. What happens to durability guarantees?
