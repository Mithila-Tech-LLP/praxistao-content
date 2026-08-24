# Chapter 67: Qdrant Deep Dive — Architecture and Internals

Qdrant is one of the fastest-growing vector databases. Written in Rust, it consistently outperforms competitors in benchmarks and has become the go-to choice for production AI applications. This chapter explains *why* Qdrant is so fast and how it works under the hood.

## Table of Contents

1. Why Qdrant Became Popular
2. Core Data Model — Points, Payloads, Collections
3. Segments — The Storage Unit
4. HNSW with Payload Filtering
5. Quantization — Compressing Vectors
6. Named Vectors and Sparse Search
7. Write-Ahead Log and Durability
8. Distributed Mode — Sharding and Replication
9. Snapshots and Backups
10. Exercises

---

## 1. Why Qdrant Became Popular

In 2022-2024, every startup building RAG pipelines, semantic search, or recommendation systems needed a vector database. Qdrant stood out for three reasons:

**Speed:** Qdrant is written in Rust with zero-copy memory access patterns. Its HNSW implementation with compiled payload filters outperforms Python-wrapped alternatives by 3-5x.

**Payload Filtering:** Qdrant was the first to implement *pre-filtering* — applying metadata filters during the HNSW graph traversal rather than after. This is a fundamentally different and faster algorithm.

**Operational Simplicity:** Single binary, Docker-native, gRPC + REST dual API, sensible defaults. You can be in production in 30 minutes.

Qdrant is fully open source (Apache 2.0) with a managed cloud option. The GitHub repo hit 20k+ stars faster than any other vector DB.

---

## 2. Core Data Model — Points, Payloads, Collections

**Collection:** A named group of vectors with a fixed dimension and distance metric. Think of it as a table in SQL.

```
Collection: "products"
  - dimension: 1536
  - distance: Cosine
  - points: 50 million
```

**Point:** The fundamental record. Every point has:
- `id`: unsigned integer or UUID
- `vector`: float32 array (or multiple named vectors)
- `payload`: arbitrary JSON object (the metadata)

```json
{
  "id": 12345,
  "vector": [0.023, -0.156, 0.891, ...],
  "payload": {
    "title": "Running Shoes Model X",
    "category": "footwear",
    "price": 129.99,
    "in_stock": true,
    "tags": ["running", "trail", "waterproof"]
  }
}
```

**Distance Metrics:**

| Metric | Formula | When to Use |
|--------|---------|------------|
| **Cosine** | `1 - (A·B)/(‖A‖‖B‖)` | Text embeddings (direction matters, not magnitude) |
| **Dot Product** | `-(A·B)` | When embeddings are normalized; OpenAI embeddings |
| **Euclidean** | `‖A - B‖²` | Image features, spatial data |
| **Manhattan** | `Σ|Aᵢ - Bᵢ|` | Sparse, high-dimensional data |

---

## 3. Segments — The Storage Unit

Qdrant organizes data into **segments** — the most important architectural concept to understand.

```
Collection
├── Segment 0 (sealed)   ← read-only, fully indexed, fast
│   ├── vector_storage/  ← memmap'd float32 arrays
│   ├── hnsw_index/      ← pre-built graph
│   └── payload_index/   ← B-Tree indexes per payload field
├── Segment 1 (sealed)
└── Segment 2 (growing)  ← write-optimized, sparse index
    ├── vector_storage/
    ├── hnsw_index/
    └── payload_index/
```

**Growing Segment:** Accepts all new writes. Has a partial HNSW index that gets updated incrementally. Fast for writes.

**Sealed Segment:** Read-only. Built when the growing segment reaches a size threshold (default: 20k points). The HNSW index is rebuilt from scratch in the optimal layout. Fast for reads.

**Why segments matter:**
1. Searches run in *parallel* across all segments. 8 sealed segments = 8-thread parallelism.
2. Sealed segments use memory-mapped files (`mmap`) — the OS handles caching. Qdrant never runs out of RAM; it just pages.
3. Compaction merges small sealed segments into larger ones to reduce overhead.

```
Search request
    ↓
Parallel search across all segments
    ↓
Merge top-k results from each segment
    ↓
Return global top-k
```

---

## 4. HNSW with Payload Filtering

This is Qdrant's biggest technical differentiator. Let's understand the problem.

**The Filtering Problem:**

```
Query: "find 10 shoes most similar to this image"
Filter: category = "running" AND price < 150 AND in_stock = true
```

**Approach 1 — Post-filtering (naive):**
1. Run HNSW search, get top 10,000 results
2. Apply filters, keep matching ones
3. Hope you got at least 10 matches

Problem: If only 1% of vectors match the filter, you need to retrieve 100,000 candidates to get 10 results. Extremely slow and wastes resources.

**Approach 2 — Pre-filtering + brute force (pgvector default):**
1. Apply SQL WHERE clause to get all matching IDs (e.g., 50,000 out of 1M)
2. Brute-force compare query vector to all 50,000

Works for small filtered sets, but O(n) for the filtered subset.

**Approach 3 — Qdrant's filtered HNSW:**

Qdrant builds a payload index (B-Tree per field) alongside the HNSW graph. During graph traversal, each candidate node is checked against the filter *inline*:

```
HNSW traversal step:
  candidates = get_neighbors(current_node)
  for candidate in candidates:
    if matches_filter(candidate.payload, filter):
      add_to_result_set(candidate)
    else:
      skip_but_still_use_for_navigation
```

This preserves graph connectivity while skipping non-matching nodes in the result set. The key insight: filtered-out nodes still *guide the graph traversal* — they're just not included in results.

**When Qdrant switches strategies automatically:**

```
selectivity = matched_points / total_points

if selectivity > 0.5:
    # Majority matches → post-filter is fine
    use post_filter_hnsw()
elif selectivity > 0.01:
    # Moderate selectivity → filtered HNSW
    use filtered_hnsw_traversal()
else:
    # Very selective → brute force on payload index
    use brute_force_on_filtered_ids()
```

Qdrant estimates selectivity at query time using payload index statistics, then picks the optimal strategy. No configuration needed.

**Payload Indexes:**

```
POST /collections/{name}/index

{
  "field_name": "category",
  "field_schema": "keyword"
}
```

| Schema Type | Index Type | Operators |
|------------|-----------|----------|
| `keyword` | HashMap | `==`, `!=`, `in []` |
| `integer` | B-Tree | `<`, `>`, `range`, `==` |
| `float` | B-Tree | `<`, `>`, `range` |
| `bool` | Bitmap | `==` |
| `text` | Inverted index | full-text `match` |
| `geo` | R-Tree (spatial) | `radius`, `polygon` |
| `datetime` | B-Tree | `<`, `>`, `range` |

---

## 5. Quantization — Compressing Vectors

1M vectors × 1536 dimensions × 4 bytes = **6 GB RAM**. For 10M vectors that's 60 GB. Quantization compresses vectors to reduce memory usage while preserving search quality.

**Scalar Quantization (SQ8):**

Compress each float32 (4 bytes) to a uint8 (1 byte) — 4x compression.

```
Original:  0.023, -0.156, 0.891, ...   (float32, 4 bytes each)
Quantized: 132,   52,     240, ...     (uint8, 1 byte each)

Mapping: min_val → 0, max_val → 255, linear scale
```

Accuracy loss: typically < 1% for cosine similarity. Qdrant rescores the top-k results using original vectors.

**Product Quantization (PQ):**

Split the vector into M subvectors. For each subvector, maintain a codebook of 256 centroid vectors. Encode each subvector as a single byte (index into codebook).

```
1536-dim vector, M=96 subspaces of 16 dims each
Codebook size: 96 × 256 × 16 × 4 bytes = 24 MB (shared)
Per-vector cost: 96 bytes instead of 6144 bytes → 64x compression
```

Tradeoff: Lower accuracy than SQ8. Used when RAM is severely constrained.

**Binary Quantization:**

Each float32 → 1 bit (sign only). 32x compression from float32.

```
0.023  → 1  (positive)
-0.156 → 0  (negative)
0.891  → 1  (positive)
```

Distance computed with XOR + popcount (hardware-accelerated). Works surprisingly well for OpenAI and Cohere embeddings (which are trained for it). 40x+ speedup on query time.

**Configuration example:**

```json
PUT /collections/products
{
  "vectors": {
    "size": 1536,
    "distance": "Cosine",
    "quantization_config": {
      "scalar": {
        "type": "int8",
        "quantile": 0.99,
        "always_ram": true
      }
    }
  }
}
```

---

## 6. Named Vectors and Sparse Search

**Named Vectors:** A single point can have multiple embeddings from different models.

```json
{
  "id": 1,
  "vectors": {
    "image":    [0.1, 0.5, 0.9, ...],   // CLIP image embedding (512 dims)
    "text":     [0.023, -0.156, ...],    // text embedding (1536 dims)
    "colbert":  [[0.1, 0.2], [0.3, 0.4], ...] // multi-vector (late interaction)
  },
  "payload": { "sku": "SHOE-001" }
}
```

You can search by any named vector independently. This lets you build multimodal search ("find by image" or "find by text description") without duplicating storage.

**Sparse Vectors (Hybrid Search):**

Dense vectors capture semantic similarity. Sparse vectors (like BM25/SPLADE) capture keyword/exact-match signals. Hybrid search combines both.

```json
{
  "id": 1,
  "vectors": {
    "dense": [0.023, -0.156, 0.891, ...],
    "sparse": {
      "indices": [142, 890, 2341],
      "values":  [0.8, 0.3, 0.5]
    }
  }
}
```

Query:
```json
POST /collections/docs/points/query
{
  "prefetch": [
    { "query": { "sparse": {"indices": [142], "values": [1.0]} }, "using": "sparse", "limit": 100 },
    { "query": [0.023, -0.156, ...], "using": "dense", "limit": 100 }
  ],
  "query": { "fusion": "rrf" }
}
```

RRF (Reciprocal Rank Fusion) merges the two result lists by ranking position, not score. This avoids the problem of incomparable score scales between dense and sparse.

---

## 7. Write-Ahead Log and Durability

Every mutating operation is written to the WAL before being applied:

```
Client → POST /points/upsert
            ↓
    1. Write to WAL (fsync)
            ↓
    2. Apply to growing segment
            ↓
    3. Return 200 OK
```

On crash, Qdrant replays the WAL from the last checkpoint to recover. The growing segment is rebuilt; sealed segments are untouched (immutable).

**WAL format (simplified):**

```
[timestamp | operation_type | payload_bytes | checksum]

Operation types:
  UPSERT_POINTS
  DELETE_POINTS
  CREATE_COLLECTION
  UPDATE_COLLECTION
  CREATE_INDEX
```

WAL files are rotated when they exceed a size threshold (default: 32 MB). Old WAL files are deleted once the data is durably in sealed segments.

---

## 8. Distributed Mode — Sharding and Replication

For collections too large for one node, Qdrant distributes data.

**Sharding:**

```
Collection (10M points) split across 4 shards:
  Node 1: Shard 0 (0-2.5M points)
  Node 2: Shard 1 (2.5-5M points)
  Node 3: Shard 2 (5-7.5M points)
  Node 4: Shard 3 (7.5-10M points)
```

Shard assignment: `hash(point_id) % num_shards`. Each node is responsible for a subset of points.

**Search in distributed mode:**

```
Query arrives at any node (coordinator role rotates)
    ↓
Coordinator fans out search to all shard leaders
    ↓
Each shard returns its local top-k
    ↓
Coordinator merges results → global top-k
    ↓
Returns to client
```

**Replication:**

Each shard can have multiple replicas for fault tolerance.

```
Shard 0: Leader (Node 1) + Replicas (Node 2, Node 3)
```

Write consistency options:
- `weak` (default): write acknowledged after leader applies it
- `majority`: write acknowledged after majority of replicas apply
- `quorum`: write acknowledged after all replicas apply

Read consistency options:
- `local` (default): read from whatever replica handles the connection
- `majority`: read from majority, compare, return consistent result

---

## 9. Snapshots and Backups

Snapshots are consistent point-in-time backups of a collection.

```bash
# Create snapshot
POST /collections/products/snapshots
# Returns: {"name": "products-2024-01-15-03-00-00.snapshot"}

# List snapshots
GET /collections/products/snapshots

# Download snapshot
GET /collections/products/snapshots/products-2024-01-15-03-00-00.snapshot

# Restore from snapshot
PUT /collections/products/snapshots/recover
{"location": "s3://my-bucket/products-2024-01-15-03-00-00.snapshot"}
```

Internally, a snapshot is a tarball containing:
- All sealed segment files (vector storage + HNSW index + payload index)
- WAL tail (uncommitted writes in the growing segment)
- Collection metadata (config, dimension, distance metric)

Snapshots are created with copy-on-write semantics — no need to pause writes during backup.

---

## Summary

- Qdrant's core unit is the **Point** (id + vector + payload). Points are grouped into **Collections**.
- Data is stored in **Segments** (growing for writes, sealed for reads). Searches run in parallel across all segments.
- **Payload filtering** during HNSW traversal is Qdrant's key innovation. It auto-selects the optimal strategy based on filter selectivity.
- **Quantization** reduces RAM by 4-64x with minimal accuracy loss. Binary quantization is 40x faster.
- **Named vectors** enable multimodal search. **Sparse + dense** hybrid search gives the best of semantic and keyword matching.
- The **WAL** ensures durability. **Segments** are immutable once sealed — crash recovery just replays the WAL.
- **Sharding** scales horizontally. **Replication** provides fault tolerance with configurable consistency.

### Exercises

**Easy:** Explain why Qdrant uses segments instead of one big HNSW index. What are the benefits for reads and writes? What happens during segment sealing?

**Medium:** Design the payload index strategy for an e-commerce search: 10M products, users typically filter by `category` (50 values), `price` (range), and `in_stock` (bool). Which selectivity pattern does each filter create? Which Qdrant filter strategy (pre-filter, filtered HNSW, brute force) does each trigger?

**Hard:** Qdrant's distributed mode uses consistent hashing for shard assignment. Research how consistent hashing works and explain: (1) why `hash(point_id) % num_shards` is problematic when you add a new shard, and (2) how Qdrant's "resharding" operation moves points to balance the new distribution.
