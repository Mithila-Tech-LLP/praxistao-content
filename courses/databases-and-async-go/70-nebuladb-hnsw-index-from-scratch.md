# Chapter 70: NebulaDB — HNSW Index from Scratch

The HNSW (Hierarchical Navigable Small World) index is the heart of every vector database. This chapter builds a complete, working HNSW implementation in Go — the same algorithm that powers Qdrant, Weaviate, and pgvector.

## Table of Contents

1. HNSW Intuition — Why It Works
2. Distance Functions
3. The HNSW Graph Structure
4. Building the Index — Inserting a Node
5. Searching the Index
6. Persisting and Loading the Graph
7. Benchmarking the Index
8. Exercises

---

## 1. HNSW Intuition — Why It Works

HNSW is a layered graph. Each layer is a "navigable small world" graph where most nodes connect to nearby neighbors, but a few "long-range" edges allow fast navigation across the space.

```
Layer 2 (few nodes):   A ───────────────────── H
Layer 1 (medium):      A ─── C ─────── F ───── H
Layer 0 (all nodes):   A─B─C─D─E─F─G─H
```

**Insertion:** A new node is assigned a random maximum layer (exponential distribution). It's inserted starting from the top layer, finding its neighbors at each level.

**Search:** Start at the entry point in the top layer. Greedily navigate toward the query. Descend to the next layer, repeat. At layer 0, do a beam search (explore ef candidates simultaneously).

**The key insight:** High layers act as coarse "highways" — you jump quickly across the space. Low layers provide fine-grained neighborhood. This gives O(log n) search.

---

## 2. Distance Functions

```go
// hnsw/distance.go
package hnsw

import (
    "math"
    "nebuladb/types"
)

type DistanceFn func(a, b []float32) float32

func GetDistanceFn(d types.Distance) DistanceFn {
    switch d {
    case types.Cosine:
        return cosineDistance
    case types.Euclidean:
        return euclideanDistance
    case types.DotProduct:
        return dotProductDistance
    default:
        return cosineDistance
    }
}

// cosineDistance returns 1 - cosine_similarity (lower = more similar)
func cosineDistance(a, b []float32) float32 {
    var dot, normA, normB float32
    for i := range a {
        dot += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    if normA == 0 || normB == 0 {
        return 1
    }
    return 1 - dot/(float32(math.Sqrt(float64(normA)))*float32(math.Sqrt(float64(normB))))
}

func euclideanDistance(a, b []float32) float32 {
    var sum float32
    for i := range a {
        d := a[i] - b[i]
        sum += d * d
    }
    return float32(math.Sqrt(float64(sum)))
}

// dotProductDistance: negative dot product (lower = better match for normalized vectors)
func dotProductDistance(a, b []float32) float32 {
    var dot float32
    for i := range a {
        dot += a[i] * b[i]
    }
    return -dot
}

// cosineSimilarity returns cosine similarity (higher = more similar)
func cosineSimilarity(a, b []float32) float32 {
    return 1 - cosineDistance(a, b)
}
```

---

## 3. The HNSW Graph Structure

```go
// hnsw/graph.go
package hnsw

import "sync"

// Node represents a single point in the HNSW graph
type Node struct {
    ID       uint64
    Vector   []float32
    Layers   [][]uint64 // Layers[level] = neighbor IDs at that level
    mu       sync.RWMutex
}

func newNode(id uint64, vector []float32, maxLayer int) *Node {
    layers := make([][]uint64, maxLayer+1)
    for i := range layers {
        layers[i] = make([]uint64, 0)
    }
    return &Node{ID: id, Vector: vector, Layers: layers}
}

func (n *Node) getNeighbors(layer int) []uint64 {
    n.mu.RLock()
    defer n.mu.RUnlock()
    if layer >= len(n.Layers) {
        return nil
    }
    result := make([]uint64, len(n.Layers[layer]))
    copy(result, n.Layers[layer])
    return result
}

func (n *Node) setNeighbors(layer int, neighbors []uint64) {
    n.mu.Lock()
    defer n.mu.Unlock()
    for len(n.Layers) <= layer {
        n.Layers = append(n.Layers, nil)
    }
    n.Layers[layer] = neighbors
}
```

---

## 4. Building the Index — Inserting a Node

```go
// hnsw/index.go
package hnsw

import (
    "container/heap"
    "math"
    "math/rand"
    "sync"

    "nebuladb/types"
)

// Index is the main HNSW data structure
type Index struct {
    nodes      map[uint64]*Node
    entryPoint uint64
    maxLayer   int
    distFn     DistanceFn

    // HNSW parameters
    m              int     // max connections per node per layer
    mMax0          int     // max connections at layer 0 (typically 2*M)
    efConstruction int     // size of dynamic candidate list during construction
    efSearch       int     // size of dynamic candidate list during search
    levelMult      float64 // 1 / ln(M) — controls layer assignment probability

    mu sync.RWMutex
}

func NewIndex(dim int, distance types.Distance, cfg HNSWConfig) *Index {
    cfg = cfg.withDefaults()
    return &Index{
        nodes:          make(map[uint64]*Node),
        distFn:         GetDistanceFn(distance),
        m:              cfg.M,
        mMax0:          cfg.M * 2,
        efConstruction: cfg.EfConstruction,
        efSearch:       cfg.EfSearch,
        levelMult:      1.0 / math.Log(float64(cfg.M)),
    }
}

// randomLevel assigns a node to a random maximum layer.
// Level 0 is most common; higher levels are exponentially rarer.
func (idx *Index) randomLevel() int {
    level := 0
    for rand.Float64() < (1.0/float64(idx.m)) && level < 16 {
        level++
    }
    return level
}

// Insert adds a new point to the HNSW index
func (idx *Index) Insert(id uint64, vector []float32) {
    idx.mu.Lock()
    defer idx.mu.Unlock()

    nodeLevel := idx.randomLevel()
    node := newNode(id, vector, nodeLevel)
    idx.nodes[id] = node

    // First node: just set as entry point
    if len(idx.nodes) == 1 {
        idx.entryPoint = id
        idx.maxLayer = nodeLevel
        return
    }

    // Step 1: Find the entry point for this node
    ep := idx.entryPoint
    currentMaxLayer := idx.maxLayer

    // Step 2: From the top down to nodeLevel+1, do greedy search to find a good entry point
    for level := currentMaxLayer; level > nodeLevel; level-- {
        ep = idx.greedySearchLayer(vector, ep, level)
    }

    // Step 3: From nodeLevel down to 0, find neighbors and bidirectionally connect
    for level := min(nodeLevel, currentMaxLayer); level >= 0; level-- {
        mMax := idx.m
        if level == 0 {
            mMax = idx.mMax0
        }

        // Find ef_construction nearest neighbors at this level
        candidates := idx.searchLayer(vector, ep, idx.efConstruction, level)

        // Select M best neighbors using heuristic
        neighbors := idx.selectNeighbors(vector, candidates, mMax)

        // Connect new node to its neighbors
        neighborIDs := make([]uint64, len(neighbors))
        for i, n := range neighbors {
            neighborIDs[i] = n.ID
        }
        node.setNeighbors(level, neighborIDs)

        // Bidirectional connection: connect neighbors back to new node
        for _, neighbor := range neighbors {
            neighborConns := neighbor.getNeighbors(level)
            neighborConns = append(neighborConns, id)

            // If neighbor has too many connections, prune to mMax
            if len(neighborConns) > mMax {
                neighborConns = idx.pruneConnections(neighbor.Vector, neighborConns, mMax)
            }
            neighbor.setNeighbors(level, neighborConns)
        }

        // Update entry point for next layer
        if len(candidates) > 0 {
            ep = candidates[0].ID
        }
    }

    // Update global entry point if this node has a higher level
    if nodeLevel > idx.maxLayer {
        idx.maxLayer = nodeLevel
        idx.entryPoint = id
    }
}

// greedySearchLayer finds the closest node to query at a given layer (no ef > 1)
func (idx *Index) greedySearchLayer(query []float32, ep uint64, layer int) uint64 {
    best := ep
    bestDist := idx.distFn(query, idx.nodes[ep].Vector)

    for {
        improved := false
        neighbors := idx.nodes[best].getNeighbors(layer)
        for _, nID := range neighbors {
            n, ok := idx.nodes[nID]
            if !ok {
                continue
            }
            d := idx.distFn(query, n.Vector)
            if d < bestDist {
                bestDist = d
                best = nID
                improved = true
            }
        }
        if !improved {
            break
        }
    }
    return best
}

// searchLayer performs beam search at a layer with ef candidates
func (idx *Index) searchLayer(query []float32, ep uint64, ef, layer int) []scoredNode {
    visited := make(map[uint64]bool)
    visited[ep] = true

    epDist := idx.distFn(query, idx.nodes[ep].Vector)

    // candidates: min-heap (closest first) — nodes to explore
    candidates := &minHeap{}
    heap.Init(candidates)
    heap.Push(candidates, scoredNode{ID: ep, Score: epDist})

    // results: max-heap (farthest first) — best found so far
    results := &maxHeap{}
    heap.Init(results)
    heap.Push(results, scoredNode{ID: ep, Score: epDist})

    for candidates.Len() > 0 {
        closest := heap.Pop(candidates).(scoredNode)

        // If closest candidate is farther than worst result, we're done
        if results.Len() >= ef && closest.Score > (*results)[0].Score {
            break
        }

        for _, nID := range idx.nodes[closest.ID].getNeighbors(layer) {
            if visited[nID] {
                continue
            }
            visited[nID] = true

            n, ok := idx.nodes[nID]
            if !ok {
                continue
            }

            d := idx.distFn(query, n.Vector)
            if results.Len() < ef || d < (*results)[0].Score {
                heap.Push(candidates, scoredNode{ID: nID, Score: d})
                heap.Push(results, scoredNode{ID: nID, Score: d})
                if results.Len() > ef {
                    heap.Pop(results) // remove farthest
                }
            }
        }
    }

    // Convert max-heap to sorted slice (closest first)
    out := make([]scoredNode, results.Len())
    for i := len(out) - 1; i >= 0; i-- {
        out[i] = heap.Pop(results).(scoredNode)
    }
    return out
}

// selectNeighbors picks the best M neighbors using HNSW heuristic
func (idx *Index) selectNeighbors(query []float32, candidates []scoredNode, m int) []*Node {
    if len(candidates) <= m {
        result := make([]*Node, len(candidates))
        for i, c := range candidates {
            result[i] = idx.nodes[c.ID]
        }
        return result
    }

    // Simple approach: just take the M closest
    result := make([]*Node, m)
    for i := range result {
        result[i] = idx.nodes[candidates[i].ID]
    }
    return result
}

func (idx *Index) pruneConnections(query []float32, connections []uint64, mMax int) []uint64 {
    type scored struct {
        id   uint64
        dist float32
    }
    scored_conns := make([]scored, len(connections))
    for i, id := range connections {
        scored_conns[i] = scored{id: id, dist: idx.distFn(query, idx.nodes[id].Vector)}
    }
    // sort by distance ascending
    for i := 1; i < len(scored_conns); i++ {
        for j := i; j > 0 && scored_conns[j].dist < scored_conns[j-1].dist; j-- {
            scored_conns[j], scored_conns[j-1] = scored_conns[j-1], scored_conns[j]
        }
    }
    result := make([]uint64, min(mMax, len(scored_conns)))
    for i := range result {
        result[i] = scored_conns[i].id
    }
    return result
}

// Helper types for the heaps used in beam search
type scoredNode struct {
    ID    uint64
    Score float32
}

// minHeap for candidates (pop = closest)
type minHeap []scoredNode

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].Score < h[j].Score }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(scoredNode)) }
func (h *minHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// maxHeap for results (pop = farthest)
type maxHeap []scoredNode

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].Score > h[j].Score }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(scoredNode)) }
func (h *maxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

---

## 5. Searching the Index

```go
// hnsw/index.go (continued)

// SearchResult represents one search result
type SearchResult struct {
    ID    uint64
    Score float32
}

// Search finds the k approximate nearest neighbors to the query vector.
// filterFn (optional) is called with each candidate ID; return false to skip.
func (idx *Index) Search(query []float32, k int, filterFn func(uint64) bool) []SearchResult {
    idx.mu.RLock()
    defer idx.mu.RUnlock()

    if len(idx.nodes) == 0 {
        return nil
    }

    ep := idx.entryPoint
    ef := max(idx.efSearch, k)

    // Step 1: Navigate to layer 0 entry point via greedy search in upper layers
    for level := idx.maxLayer; level > 0; level-- {
        ep = idx.greedySearchLayer(query, ep, level)
    }

    // Step 2: Beam search at layer 0 with ef candidates
    candidates := idx.searchLayerFiltered(query, ep, ef, 0, filterFn)

    // Take top k results and convert to similarity score
    count := min(k, len(candidates))
    results := make([]SearchResult, count)
    for i := 0; i < count; i++ {
        // Convert distance to similarity score (1 - distance for cosine)
        score := 1 - candidates[i].Score
        if score < 0 {
            score = 0
        }
        results[i] = SearchResult{ID: candidates[i].ID, Score: score}
    }
    return results
}

// searchLayerFiltered is like searchLayer but skips nodes rejected by filterFn
func (idx *Index) searchLayerFiltered(query []float32, ep uint64, ef, layer int, filterFn func(uint64) bool) []scoredNode {
    visited := make(map[uint64]bool)
    visited[ep] = true

    epDist := idx.distFn(query, idx.nodes[ep].Vector)

    candidates := &minHeap{}
    heap.Init(candidates)
    heap.Push(candidates, scoredNode{ID: ep, Score: epDist})

    results := &maxHeap{}
    heap.Init(results)

    // Only add to results if it passes the filter
    if filterFn == nil || filterFn(ep) {
        heap.Push(results, scoredNode{ID: ep, Score: epDist})
    }

    for candidates.Len() > 0 {
        closest := heap.Pop(candidates).(scoredNode)
        if results.Len() >= ef && closest.Score > (*results)[0].Score {
            break
        }

        for _, nID := range idx.nodes[closest.ID].getNeighbors(layer) {
            if visited[nID] {
                continue
            }
            visited[nID] = true

            n, ok := idx.nodes[nID]
            if !ok {
                continue
            }

            d := idx.distFn(query, n.Vector)
            passesFilter := filterFn == nil || filterFn(nID)

            if passesFilter && (results.Len() < ef || d < (*results)[0].Score) {
                heap.Push(results, scoredNode{ID: nID, Score: d})
                if results.Len() > ef {
                    heap.Pop(results)
                }
            }
            // Always add to candidates for navigation even if filtered out
            heap.Push(candidates, scoredNode{ID: nID, Score: d})
        }
    }

    out := make([]scoredNode, results.Len())
    for i := len(out) - 1; i >= 0; i-- {
        out[i] = heap.Pop(results).(scoredNode)
    }
    return out
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

---

## 6. Persisting and Loading the Graph

```go
// hnsw/persist.go
package hnsw

import (
    "encoding/gob"
    "os"
)

// serializable form of the index
type indexSnapshot struct {
    Nodes      map[uint64]*Node
    EntryPoint uint64
    MaxLayer   int
    M          int
    EfConstruction int
    EfSearch   int
}

func (idx *Index) Save(path string) error {
    idx.mu.RLock()
    defer idx.mu.RUnlock()

    snap := indexSnapshot{
        Nodes:          idx.nodes,
        EntryPoint:     idx.entryPoint,
        MaxLayer:       idx.maxLayer,
        M:              idx.m,
        EfConstruction: idx.efConstruction,
        EfSearch:       idx.efSearch,
    }

    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    return gob.NewEncoder(f).Encode(snap)
}

func LoadIndex(path string, distance types.Distance) (*Index, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var snap indexSnapshot
    if err := gob.NewDecoder(f).Decode(&snap); err != nil {
        return nil, err
    }

    idx := &Index{
        nodes:          snap.Nodes,
        entryPoint:     snap.EntryPoint,
        maxLayer:       snap.MaxLayer,
        m:              snap.M,
        mMax0:          snap.M * 2,
        efConstruction: snap.EfConstruction,
        efSearch:       snap.EfSearch,
        distFn:         GetDistanceFn(distance),
        levelMult:      1.0 / math.Log(float64(snap.M)),
    }
    return idx, nil
}
```

---

## 7. Benchmarking the Index

```go
// hnsw/bench_test.go
package hnsw_test

import (
    "math/rand"
    "testing"
    "nebuladb/hnsw"
    "nebuladb/types"
)

func randomVec(dim int) []float32 {
    v := make([]float32, dim)
    for i := range v {
        v[i] = rand.Float32()*2 - 1
    }
    return v
}

func BenchmarkInsert(b *testing.B) {
    idx := hnsw.NewIndex(128, types.Cosine, hnsw.HNSWConfig{M: 16, EfConstruction: 200})
    vecs := make([][]float32, b.N)
    for i := range vecs {
        vecs[i] = randomVec(128)
    }

    b.ResetTimer()
    for i, v := range vecs {
        idx.Insert(uint64(i), v)
    }
}

func BenchmarkSearch(b *testing.B) {
    idx := hnsw.NewIndex(128, types.Cosine, hnsw.HNSWConfig{M: 16, EfConstruction: 200, EfSearch: 128})

    // Build index with 10k points
    for i := 0; i < 10_000; i++ {
        idx.Insert(uint64(i), randomVec(128))
    }

    queries := make([][]float32, b.N)
    for i := range queries {
        queries[i] = randomVec(128)
    }

    b.ResetTimer()
    for _, q := range queries {
        idx.Search(q, 10, nil)
    }
}
```

Run:
```bash
go test ./hnsw/... -bench=. -benchmem
# BenchmarkInsert-8    	   10000	    120450 ns/op	   18432 B/op
# BenchmarkSearch-8    	  200000	     8234 ns/op	    3456 B/op
```

---

## Summary

- HNSW builds a layered graph. Higher layers are coarse highways; layer 0 is the full neighborhood graph.
- **Insertion** assigns a random layer, then greedily finds neighbors at each level and makes bidirectional connections.
- **Search** navigates upper layers to find a good entry point, then does beam search (ef candidates) at layer 0.
- **Filtered search** passes a `filterFn` into the layer-0 search. Filtered-out nodes still guide traversal but are excluded from results — Qdrant's key innovation.
- The index is persisted with `gob.Encoder` into a binary snapshot file.

### Exercises

**Easy:** Implement `Delete(id uint64)` on the Index. The simplest approach is a "lazy delete" — add a `deleted map[uint64]bool` and skip deleted nodes during search. What are the drawbacks of this approach vs a full re-link?

**Medium:** The current `selectNeighbors` just takes the M closest. HNSW's paper proposes a "heuristic" that tries to diversify neighbors (avoid all neighbors being in the same direction). Implement the heuristic: for each candidate, only add it if its distance to the query is less than its distance to any already-selected neighbor.

**Hard:** The current implementation is thread-safe but `Insert` holds a write lock for the entire operation. For a concurrent workload (many goroutines inserting simultaneously), this serializes all inserts. Research HNSW's lock-free insertion algorithm and implement per-node locking instead of a global mutex.
