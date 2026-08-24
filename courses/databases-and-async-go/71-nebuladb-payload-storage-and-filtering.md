# Chapter 71: NebulaDB — Payload Storage and Filtering Engine

Every vector needs metadata. In Qdrant, the metadata is called the "payload" — a JSON object attached to each point. This chapter builds NebulaDB's payload layer: storage, field indexing, and the filter evaluation engine that powers pre-filtered HNSW search.

## Table of Contents

1. Payload Storage with BoltDB
2. The Filter Engine — Evaluating Conditions
3. Payload Indexes — B-Trees for Fast Filtering
4. Pre-filter vs Post-filter Strategy
5. The Collection Layer — Wiring Everything Together
6. Exercises

---

## 1. Payload Storage with BoltDB

Each point's payload is a JSON object stored in BoltDB (an embedded key-value store). BoltDB is a single file, ACID-compliant, and great for embedded use cases.

```go
// storage/payload_store.go
package storage

import (
    "encoding/json"
    "fmt"

    bolt "go.etcd.io/bbolt"
    "nebuladb/types"
)

var payloadBucket = []byte("payloads")

// PayloadStore persists JSON payloads keyed by point ID
type PayloadStore struct {
    db *bolt.DB
}

func NewPayloadStore(path string) (*PayloadStore, error) {
    db, err := bolt.Open(path, 0600, nil)
    if err != nil {
        return nil, fmt.Errorf("open payload db: %w", err)
    }

    err = db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(payloadBucket)
        return err
    })
    if err != nil {
        db.Close()
        return nil, err
    }

    return &PayloadStore{db: db}, nil
}

// Set stores or overwrites a payload for a point
func (s *PayloadStore) Set(id uint64, payload types.Payload) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    key := encodeID(id)
    return s.db.Update(func(tx *bolt.Tx) error {
        return tx.Bucket(payloadBucket).Put(key, data)
    })
}

// Get retrieves a payload by point ID
func (s *PayloadStore) Get(id uint64) (types.Payload, error) {
    var payload types.Payload
    err := s.db.View(func(tx *bolt.Tx) error {
        data := tx.Bucket(payloadBucket).Get(encodeID(id))
        if data == nil {
            return fmt.Errorf("point %d not found", id)
        }
        return json.Unmarshal(data, &payload)
    })
    return payload, err
}

// Delete removes a payload
func (s *PayloadStore) Delete(id uint64) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        return tx.Bucket(payloadBucket).Delete(encodeID(id))
    })
}

// GetMany retrieves payloads for multiple IDs in one transaction
func (s *PayloadStore) GetMany(ids []uint64) (map[uint64]types.Payload, error) {
    result := make(map[uint64]types.Payload, len(ids))
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(payloadBucket)
        for _, id := range ids {
            data := b.Get(encodeID(id))
            if data == nil {
                continue
            }
            var p types.Payload
            if err := json.Unmarshal(data, &p); err != nil {
                return err
            }
            result[id] = p
        }
        return nil
    })
    return result, err
}

func (s *PayloadStore) Close() error {
    return s.db.Close()
}

// encodeID encodes a uint64 as an 8-byte big-endian key (BoltDB sorts lexicographically)
func encodeID(id uint64) []byte {
    b := make([]byte, 8)
    b[0] = byte(id >> 56)
    b[1] = byte(id >> 48)
    b[2] = byte(id >> 40)
    b[3] = byte(id >> 32)
    b[4] = byte(id >> 24)
    b[5] = byte(id >> 16)
    b[6] = byte(id >> 8)
    b[7] = byte(id)
    return b
}
```

---

## 2. The Filter Engine — Evaluating Conditions

The filter engine takes a `types.Filter` and evaluates it against a payload. This is what the HNSW search calls via `filterFn`.

```go
// index/payload_index.go
package index

import (
    "encoding/json"
    "fmt"
    "strings"

    "nebuladb/types"
)

// Evaluator evaluates filter conditions against payload values
type Evaluator struct{}

// Match returns true if the payload satisfies the filter
func (e *Evaluator) Match(payload types.Payload, filter *types.Filter) bool {
    if filter == nil {
        return true
    }

    // All MUST conditions must be true
    for _, cond := range filter.Must {
        if !e.evaluateCondition(payload, cond) {
            return false
        }
    }

    // All MUST_NOT conditions must be false
    for _, cond := range filter.MustNot {
        if e.evaluateCondition(payload, cond) {
            return false
        }
    }

    // At least one SHOULD condition must be true (if any exist)
    if len(filter.Should) > 0 {
        anyMatch := false
        for _, cond := range filter.Should {
            if e.evaluateCondition(payload, cond) {
                anyMatch = true
                break
            }
        }
        if !anyMatch {
            return false
        }
    }

    return true
}

func (e *Evaluator) evaluateCondition(payload types.Payload, cond types.Condition) bool {
    raw, ok := payload[cond.Field]
    if !ok {
        return false
    }

    var fieldValue any
    if err := json.Unmarshal(raw, &fieldValue); err != nil {
        return false
    }

    if cond.Match != nil {
        return e.evaluateMatch(fieldValue, cond.Match)
    }
    if cond.Range != nil {
        return e.evaluateRange(fieldValue, cond.Range)
    }
    return false
}

func (e *Evaluator) evaluateMatch(value any, match *types.MatchCondition) bool {
    // Handle array values (e.g. tags: ["running", "trail"])
    if arr, ok := value.([]any); ok {
        for _, elem := range arr {
            if e.matchSingle(elem, match) {
                return true
            }
        }
        return false
    }
    return e.matchSingle(value, match)
}

func (e *Evaluator) matchSingle(value any, match *types.MatchCondition) bool {
    if match.Any != nil {
        // Match any of the provided values
        for _, target := range match.Any {
            if e.equal(value, target) {
                return true
            }
        }
        return false
    }
    // Exact match
    return e.equal(value, match.Value)
}

func (e *Evaluator) equal(a, b any) bool {
    // Normalize to string for comparison (handles JSON number/string differences)
    return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func (e *Evaluator) evaluateRange(value any, r *types.RangeCondition) bool {
    var num float64
    switch v := value.(type) {
    case float64:
        num = v
    case int:
        num = float64(v)
    case json.Number:
        n, err := v.Float64()
        if err != nil {
            return false
        }
        num = n
    default:
        return false
    }

    if r.Gt != nil && !(num > *r.Gt) {
        return false
    }
    if r.Gte != nil && !(num >= *r.Gte) {
        return false
    }
    if r.Lt != nil && !(num < *r.Lt) {
        return false
    }
    if r.Lte != nil && !(num <= *r.Lte) {
        return false
    }
    return true
}
```

---

## 3. Payload Indexes — B-Trees for Fast Filtering

Scanning all payloads for each filter is O(n). A payload index pre-groups IDs by field value, making filter lookups O(1) for exact matches.

```go
// index/field_index.go
package index

import (
    "encoding/json"
    "fmt"
    "sort"
    "sync"

    "nebuladb/types"
)

// FieldIndex maintains an in-memory index for a single payload field
type FieldIndex struct {
    field string
    // For keyword fields: value → set of point IDs
    keywords map[string]map[uint64]struct{}
    // For numeric fields: sorted list of (value, id) pairs for range queries
    numerics []numericEntry
    mu       sync.RWMutex
}

type numericEntry struct {
    value float64
    id    uint64
}

func NewFieldIndex(field string) *FieldIndex {
    return &FieldIndex{
        field:    field,
        keywords: make(map[string]map[uint64]struct{}),
    }
}

// Index adds a point's field value to the index
func (fi *FieldIndex) Index(id uint64, payload types.Payload) error {
    raw, ok := payload[fi.field]
    if !ok {
        return nil // field not present in this point
    }

    var value any
    if err := json.Unmarshal(raw, &value); err != nil {
        return err
    }

    fi.mu.Lock()
    defer fi.mu.Unlock()

    switch v := value.(type) {
    case string:
        fi.addKeyword(id, v)
    case float64:
        fi.numerics = append(fi.numerics, numericEntry{value: v, id: id})
    case bool:
        fi.addKeyword(id, fmt.Sprintf("%v", v))
    case []any:
        for _, elem := range v {
            if s, ok := elem.(string); ok {
                fi.addKeyword(id, s)
            }
        }
    }
    return nil
}

func (fi *FieldIndex) addKeyword(id uint64, value string) {
    if fi.keywords[value] == nil {
        fi.keywords[value] = make(map[uint64]struct{})
    }
    fi.keywords[value][id] = struct{}{}
}

// Remove removes a point from the index
func (fi *FieldIndex) Remove(id uint64, payload types.Payload) {
    fi.mu.Lock()
    defer fi.mu.Unlock()

    for _, ids := range fi.keywords {
        delete(ids, id)
    }
    // Remove from numerics
    newNumerics := fi.numerics[:0]
    for _, e := range fi.numerics {
        if e.id != id {
            newNumerics = append(newNumerics, e)
        }
    }
    fi.numerics = newNumerics
}

// Lookup returns IDs matching a condition (nil means "use full scan")
func (fi *FieldIndex) Lookup(cond types.Condition) (map[uint64]struct{}, bool) {
    fi.mu.RLock()
    defer fi.mu.RUnlock()

    if cond.Match != nil {
        return fi.lookupMatch(cond.Match)
    }
    if cond.Range != nil {
        return fi.lookupRange(cond.Range)
    }
    return nil, false
}

func (fi *FieldIndex) lookupMatch(match *types.MatchCondition) (map[uint64]struct{}, bool) {
    if match.Any != nil {
        result := make(map[uint64]struct{})
        for _, val := range match.Any {
            key := fmt.Sprintf("%v", val)
            for id := range fi.keywords[key] {
                result[id] = struct{}{}
            }
        }
        return result, true
    }
    key := fmt.Sprintf("%v", match.Value)
    ids, ok := fi.keywords[key]
    if !ok {
        return make(map[uint64]struct{}), true
    }
    return ids, true
}

func (fi *FieldIndex) lookupRange(r *types.RangeCondition) (map[uint64]struct{}, bool) {
    if len(fi.numerics) == 0 {
        return nil, false
    }

    // Ensure sorted
    sort.Slice(fi.numerics, func(i, j int) bool {
        return fi.numerics[i].value < fi.numerics[j].value
    })

    result := make(map[uint64]struct{})
    for _, e := range fi.numerics {
        if r.Gt != nil && !(e.value > *r.Gt) {
            continue
        }
        if r.Gte != nil && !(e.value >= *r.Gte) {
            continue
        }
        if r.Lt != nil && !(e.value < *r.Lt) {
            continue
        }
        if r.Lte != nil && !(e.value <= *r.Lte) {
            continue
        }
        result[e.id] = struct{}{}
    }
    return result, true
}
```

**PayloadIndexManager** manages multiple field indexes:

```go
// index/manager.go
package index

import (
    "sync"

    "nebuladb/types"
)

type Manager struct {
    indexes map[string]*FieldIndex
    mu      sync.RWMutex
}

func NewManager() *Manager {
    return &Manager{indexes: make(map[string]*FieldIndex)}
}

// CreateIndex registers an index for a field
func (m *Manager) CreateIndex(field string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, exists := m.indexes[field]; !exists {
        m.indexes[field] = NewFieldIndex(field)
    }
}

// Index adds a point to all registered field indexes
func (m *Manager) Index(id uint64, payload types.Payload) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    for _, fi := range m.indexes {
        fi.Index(id, payload)
    }
}

// Remove removes a point from all indexes
func (m *Manager) Remove(id uint64, payload types.Payload) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    for _, fi := range m.indexes {
        fi.Remove(id, payload)
    }
}

// FilterIDs returns the set of IDs matching the filter using available indexes.
// Returns (nil, false) if no index can narrow the candidate set.
func (m *Manager) FilterIDs(filter *types.Filter) (map[uint64]struct{}, bool) {
    if filter == nil || len(filter.Must) == 0 {
        return nil, false
    }

    m.mu.RLock()
    defer m.mu.RUnlock()

    var result map[uint64]struct{}
    for _, cond := range filter.Must {
        fi, ok := m.indexes[cond.Field]
        if !ok {
            continue // no index for this field, skip
        }

        ids, ok := fi.Lookup(cond)
        if !ok {
            continue
        }

        if result == nil {
            // First indexed condition: start with its result set
            result = make(map[uint64]struct{}, len(ids))
            for id := range ids {
                result[id] = struct{}{}
            }
        } else {
            // Intersect with subsequent conditions (AND logic)
            for id := range result {
                if _, ok := ids[id]; !ok {
                    delete(result, id)
                }
            }
        }
    }
    return result, result != nil
}
```

---

## 4. Pre-filter vs Post-filter Strategy

Like Qdrant, NebulaDB auto-selects the search strategy based on filter selectivity:

```go
// collection/search.go
package collection

import (
    "nebuladb/hnsw"
    "nebuladb/index"
    "nebuladb/types"
)

const (
    bruteForceThreshold = 0.01 // < 1% matches → brute force on filtered IDs
    postFilterThreshold = 0.5  // > 50% matches → post-filter HNSW
)

func (c *Collection) search(query types.Vector, filter *types.Filter, limit int) ([]hnsw.SearchResult, error) {
    totalPoints := c.pointCount()
    if totalPoints == 0 {
        return nil, nil
    }

    // Try to get candidate IDs from payload index
    filteredIDs, hasIndex := c.indexMgr.FilterIDs(filter)

    var selectivity float64
    if hasIndex {
        selectivity = float64(len(filteredIDs)) / float64(totalPoints)
    } else {
        selectivity = 1.0 // no index: assume all points match
    }

    evaluator := &index.Evaluator{}

    switch {
    case hasIndex && selectivity < bruteForceThreshold:
        // Very few matches: brute force over the filtered set
        return c.bruteForceSearch(query, filteredIDs, filter, evaluator, limit)

    case hasIndex && selectivity < postFilterThreshold:
        // Moderate selectivity: use filtered HNSW (check filter inline)
        filterFn := func(id uint64) bool {
            _, inSet := filteredIDs[id]
            if !inSet {
                return false
            }
            payload, err := c.payloadStore.Get(id)
            if err != nil {
                return false
            }
            return evaluator.Match(payload, filter)
        }
        return c.hnswIndex.Search(query, limit, filterFn), nil

    default:
        // High selectivity or no index: HNSW with post-filter
        results := c.hnswIndex.Search(query, limit*10, nil)
        return c.postFilter(results, filter, evaluator, limit)
    }
}

// bruteForceSearch compares the query to all filtered IDs directly
func (c *Collection) bruteForceSearch(
    query types.Vector,
    ids map[uint64]struct{},
    filter *types.Filter,
    evaluator *index.Evaluator,
    limit int,
) ([]hnsw.SearchResult, error) {
    type scored struct {
        id    uint64
        score float32
    }
    var results []scored

    distFn := hnsw.GetDistanceFn(c.config.Distance)

    for id := range ids {
        vec, err := c.vectorStore.Get(id)
        if err != nil {
            continue
        }
        if filter != nil {
            payload, err := c.payloadStore.Get(id)
            if err != nil || !evaluator.Match(payload, filter) {
                continue
            }
        }
        dist := distFn(query, vec)
        results = append(results, scored{id: id, score: 1 - dist})
    }

    // Sort by score descending
    for i := 1; i < len(results); i++ {
        for j := i; j > 0 && results[j].score > results[j-1].score; j-- {
            results[j], results[j-1] = results[j-1], results[j]
        }
    }

    count := limit
    if count > len(results) {
        count = len(results)
    }

    out := make([]hnsw.SearchResult, count)
    for i := range out {
        out[i] = hnsw.SearchResult{ID: results[i].id, Score: results[i].score}
    }
    return out, nil
}

func (c *Collection) postFilter(
    results []hnsw.SearchResult,
    filter *types.Filter,
    evaluator *index.Evaluator,
    limit int,
) ([]hnsw.SearchResult, error) {
    if filter == nil {
        if len(results) > limit {
            return results[:limit], nil
        }
        return results, nil
    }

    var filtered []hnsw.SearchResult
    for _, r := range results {
        payload, err := c.payloadStore.Get(r.ID)
        if err != nil {
            continue
        }
        if evaluator.Match(payload, filter) {
            filtered = append(filtered, r)
            if len(filtered) >= limit {
                break
            }
        }
    }
    return filtered, nil
}
```

---

## 5. The Collection Layer — Wiring Everything Together

```go
// collection/collection.go
package collection

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "sync/atomic"

    "nebuladb/hnsw"
    "nebuladb/index"
    "nebuladb/storage"
    "nebuladb/types"
)

type Collection struct {
    name         string
    dir          string
    config       Config
    hnswIndex    *hnsw.Index
    vectorStore  *storage.VectorStore
    payloadStore *storage.PayloadStore
    indexMgr     *index.Manager
    wal          *storage.WAL
    count        atomic.Int64
    mu           sync.RWMutex
}

func newCollection(name, dir string, cfg Config) (*Collection, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }
    cfg.HNSW = cfg.HNSW.withDefaults()

    payloadStore, err := storage.NewPayloadStore(filepath.Join(dir, "payload.db"))
    if err != nil {
        return nil, fmt.Errorf("payload store: %w", err)
    }

    vectorStore, err := storage.NewVectorStore(filepath.Join(dir, "vectors.bin"), cfg.Dimension)
    if err != nil {
        return nil, fmt.Errorf("vector store: %w", err)
    }

    wal, err := storage.NewWAL(filepath.Join(dir, "wal.log"))
    if err != nil {
        return nil, fmt.Errorf("wal: %w", err)
    }

    hnswIdx := hnsw.NewIndex(cfg.Dimension, cfg.Distance, hnsw.HNSWConfig{
        M:              cfg.HNSW.M,
        EfConstruction: cfg.HNSW.EfConstruction,
        EfSearch:       cfg.HNSW.EfSearch,
    })

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

    if err := saveConfig(dir, cfg); err != nil {
        return nil, err
    }
    return c, nil
}

// Upsert inserts or updates a point
func (c *Collection) Upsert(point types.Point) error {
    // 1. Write to WAL first (durability)
    if err := c.wal.WriteUpsert(point); err != nil {
        return fmt.Errorf("wal: %w", err)
    }

    // 2. Store vector
    if err := c.vectorStore.Set(point.ID, point.Vector); err != nil {
        return fmt.Errorf("vector store: %w", err)
    }

    // 3. Store payload
    if err := c.payloadStore.Set(point.ID, point.Payload); err != nil {
        return fmt.Errorf("payload store: %w", err)
    }

    // 4. Update payload indexes
    c.indexMgr.Index(point.ID, point.Payload)

    // 5. Insert into HNSW index
    c.hnswIndex.Insert(point.ID, point.Vector)

    c.count.Add(1)
    return nil
}

// Search performs a vector similarity search with optional filtering
func (c *Collection) Search(req types.SearchRequest) ([]types.ScoredPoint, error) {
    if len(req.Vector) != c.config.Dimension {
        return nil, fmt.Errorf("query vector dimension %d != collection dimension %d",
            len(req.Vector), c.config.Dimension)
    }

    results, err := c.search(req.Vector, req.Filter, req.Limit)
    if err != nil {
        return nil, err
    }

    scored := make([]types.ScoredPoint, len(results))
    for i, r := range results {
        p := types.ScoredPoint{
            Point: types.Point{ID: r.ID},
            Score: r.Score,
        }

        if req.WithPayload {
            payload, err := c.payloadStore.Get(r.ID)
            if err == nil {
                p.Payload = payload
            }
        }
        scored[i] = p
    }
    return scored, nil
}

func (c *Collection) pointCount() int64 {
    return c.count.Load()
}

func (c *Collection) Close() error {
    // Save HNSW index to disk before closing
    if err := c.hnswIndex.Save(filepath.Join(c.dir, "hnsw.bin")); err != nil {
        return err
    }
    c.vectorStore.Close()
    c.payloadStore.Close()
    c.wal.Close()
    return nil
}
```

---

## Summary

- **PayloadStore** (BoltDB) stores JSON payloads keyed by point ID. ACID-compliant, embedded, single file.
- The **Filter Engine** evaluates `Must`, `MustNot`, and `Should` conditions against a payload. Supports exact match, multi-value match, and numeric range.
- **FieldIndex** maintains an in-memory reverse index per field: keyword → `{id set}`, numeric → sorted `(value, id)` pairs.
- **FilterIDs** intersects results from multiple indexed fields to produce a narrow candidate set before HNSW search.
- **Strategy selection**: brute force (<1% selectivity), filtered HNSW (1-50%), post-filter HNSW (>50%).

### Exercises

**Easy:** Add a `text` match condition that checks if a string field *contains* a substring (case-insensitive). Add it to `evaluateCondition` and test it with the `Evaluator.Match` function.

**Medium:** The numeric field index is unsorted until `lookupRange` sorts it. This means every range query re-sorts. Implement an `append + sort on insert` approach so the list is always sorted, making range queries O(log n) via binary search.

**Hard:** The `FilterIDs` method only uses `Must` conditions for pre-filtering. Extend it to handle `Should` conditions: if all SHOULD conditions have indexes, compute the union of their ID sets. Combine this union with the MUST intersection. Make sure the logic is correct when MUST conditions exist alongside SHOULD conditions.
