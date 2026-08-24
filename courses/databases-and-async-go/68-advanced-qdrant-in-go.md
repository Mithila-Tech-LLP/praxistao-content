# Chapter 68: Advanced Qdrant in Go

Chapter 34 covered the Qdrant basics. This chapter covers everything that makes Qdrant production-grade: named vectors, sparse hybrid search, quantization, collection management, and building a real-world semantic search system.

## Table of Contents

1. Qdrant Go Client — Modern Setup
2. Named Vectors — Multimodal Search
3. Payload Indexes and Filtered Search
4. Sparse + Dense Hybrid Search
5. Quantization in Practice
6. Batch Operations and Performance
7. Collection Management
8. Project: AI-Powered Product Search
9. Exercises

---

## 1. Qdrant Go Client — Modern Setup

```bash
docker run -d -p 6333:6333 -p 6334:6334 \
  -v $(pwd)/qdrant_storage:/qdrant/storage \
  qdrant/qdrant

go get github.com/qdrant/go-client@latest
```

The Qdrant Go client uses gRPC under the hood (port 6334) with a high-level wrapper:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/qdrant/go-client/qdrant"
)

func main() {
    // New high-level client (wraps gRPC)
    client, err := qdrant.NewClient(&qdrant.Config{
        Host: "localhost",
        Port: 6334,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Health check
    health, err := client.HealthCheck(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Qdrant version: %s\n", health.GetVersion())
}
```

---

## 2. Named Vectors — Multimodal Search

A single point can have multiple embeddings from different models. This is powerful for products that have both an image and a text description.

```go
// Create collection with named vectors
err = client.CreateCollection(ctx, &qdrant.CreateCollection{
    CollectionName: "products",
    VectorsConfig: qdrant.NewVectorsConfigMap(map[string]*qdrant.VectorParams{
        "text": {
            Size:     1536,
            Distance: qdrant.Distance_Cosine,
        },
        "image": {
            Size:     512,
            Distance: qdrant.Distance_Cosine,
        },
    }),
})
```

Insert a point with multiple named vectors:

```go
_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
    CollectionName: "products",
    Points: []*qdrant.PointStruct{
        {
            Id: qdrant.NewIDNum(1),
            Vectors: qdrant.NewVectorsMap(map[string][]float32{
                "text":  {0.1, 0.5, 0.9 /* ... 1536 dims */},
                "image": {0.3, 0.2, 0.8 /* ... 512 dims */},
            }),
            Payload: qdrant.NewValueMap(map[string]any{
                "name":     "Trail Running Shoe X",
                "category": "footwear",
                "price":    129.99,
                "in_stock": true,
            }),
        },
    },
})
```

Search by a specific named vector:

```go
// Search by image embedding
results, err := client.Query(ctx, &qdrant.QueryPoints{
    CollectionName: "products",
    Query:          qdrant.NewQuery(0.3, 0.25, 0.75 /* image query vector */),
    Using:          qdrant.PtrOf("image"), // which named vector to use
    Limit:          qdrant.PtrOf(uint64(10)),
    WithPayload:    qdrant.NewWithPayload(true),
})
if err != nil {
    log.Fatal(err)
}

for _, result := range results {
    name := result.GetPayload()["name"].GetStringValue()
    fmt.Printf("Score: %.4f | %s\n", result.GetScore(), name)
}
```

---

## 3. Payload Indexes and Filtered Search

Without an index, payload filtering requires scanning all points. Always create indexes on fields you filter by.

```go
// Create payload indexes
indexes := []struct {
    field     string
    fieldType qdrant.FieldType
}{
    {"category", qdrant.FieldType_FieldTypeKeyword},
    {"price", qdrant.FieldType_FieldTypeFloat},
    {"in_stock", qdrant.FieldType_FieldTypeBool},
    {"tags", qdrant.FieldType_FieldTypeKeyword},
}

for _, idx := range indexes {
    _, err = client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
        CollectionName: "products",
        FieldName:      idx.field,
        FieldType:      qdrant.PtrOf(idx.fieldType),
    })
    if err != nil {
        log.Printf("index %s: %v", idx.field, err)
    }
}
```

Filtered vector search:

```go
results, err := client.Query(ctx, &qdrant.QueryPoints{
    CollectionName: "products",
    Query:          qdrant.NewQuery(queryVector...),
    Limit:          qdrant.PtrOf(uint64(10)),

    // Complex filter: (category = "footwear" OR category = "apparel")
    //                 AND price < 200
    //                 AND in_stock = true
    Filter: &qdrant.Filter{
        Must: []*qdrant.Condition{
            qdrant.NewMatchAny("category", "footwear", "apparel"),
            qdrant.NewRange("price", &qdrant.Range{
                Lt: qdrant.PtrOf(200.0),
            }),
            qdrant.NewMatchBool("in_stock", true),
        },
    },
    WithPayload: qdrant.NewWithPayload(true),
})
```

**Filter conditions reference:**

```go
// Exact match
qdrant.NewMatch("category", "footwear")

// Match any of these values
qdrant.NewMatchAny("tags", "running", "trail", "outdoor")

// Numeric range
qdrant.NewRange("price", &qdrant.Range{Gte: qdrant.PtrOf(50.0), Lt: qdrant.PtrOf(200.0)})

// Bool
qdrant.NewMatchBool("in_stock", true)

// Nested: must NOT match
filter := &qdrant.Filter{
    MustNot: []*qdrant.Condition{
        qdrant.NewMatch("category", "clearance"),
    },
}

// Should (OR logic)
filter := &qdrant.Filter{
    Should: []*qdrant.Condition{
        qdrant.NewMatch("category", "footwear"),
        qdrant.NewMatch("category", "apparel"),
    },
    MinShould: &qdrant.MinShould{Condition: 1},
}

// Geo filter (points within 10km of a location)
qdrant.NewGeoRadius("location", &qdrant.GeoRadius{
    Center: &qdrant.GeoPoint{Lon: 77.5946, Lat: 12.9716},
    Radius: 10_000, // meters
})
```

---

## 4. Sparse + Dense Hybrid Search

The most powerful search pattern: dense vectors find semantic matches, sparse vectors ensure keyword/exact-match results aren't missed.

```go
// Create collection with both dense and sparse vectors
err = client.CreateCollection(ctx, &qdrant.CreateCollection{
    CollectionName: "articles",
    VectorsConfig: qdrant.NewVectorsConfigMap(map[string]*qdrant.VectorParams{
        "dense": {
            Size:     1536,
            Distance: qdrant.Distance_Cosine,
        },
    }),
    SparseVectorsConfig: map[string]*qdrant.SparseVectorParams{
        "sparse": {
            Index: &qdrant.SparseIndexConfig{
                OnDisk: qdrant.PtrOf(false),
            },
        },
    },
})
```

Insert with both vector types:

```go
client.Upsert(ctx, &qdrant.UpsertPoints{
    CollectionName: "articles",
    Points: []*qdrant.PointStruct{
        {
            Id: qdrant.NewIDNum(1),
            Vectors: qdrant.NewVectorsMap(map[string][]float32{
                "dense": getDenseEmbedding("Go is great for concurrent programming"),
            }),
            Payload: qdrant.NewValueMap(map[string]any{
                "title": "Go Concurrency Patterns",
                "body":  "Go is great for concurrent programming...",
            }),
        },
    },
})

// Sparse vectors stored separately via UpdateVectors
client.UpdateVectors(ctx, &qdrant.UpdatePointVectors{
    CollectionName: "articles",
    Points: []*qdrant.PointVectors{
        {
            Id: qdrant.NewIDNum(1),
            Vectors: qdrant.NewVectorsMap(map[string][]float32{
                // sparse vectors require special handling via named sparse config
                // indices/values come from SPLADE or BM25 model
            }),
        },
    },
})
```

Hybrid search with Reciprocal Rank Fusion:

```go
results, err := client.Query(ctx, &qdrant.QueryPoints{
    CollectionName: "articles",

    // Prefetch: run both searches independently
    Prefetch: []*qdrant.PrefetchQuery{
        {
            // Dense semantic search
            Query: qdrant.NewQuery(denseQueryVector...),
            Using: qdrant.PtrOf("dense"),
            Limit: qdrant.PtrOf(uint64(50)),
        },
        {
            // Sparse keyword search (using sparse vector)
            Query: qdrant.NewQuerySparse(sparseIndices, sparseValues),
            Using: qdrant.PtrOf("sparse"),
            Limit: qdrant.PtrOf(uint64(50)),
        },
    },

    // Fuse the two result lists with RRF
    Query: qdrant.NewQueryFusion(qdrant.Fusion_Rrf),
    Limit: qdrant.PtrOf(uint64(10)),
})
```

RRF formula: `score(d) = Σ 1 / (k + rank_i(d))` where k=60 is a constant and rank_i is position in list i.

---

## 5. Quantization in Practice

```go
// Scalar quantization: 4x memory reduction, <1% accuracy loss
err = client.CreateCollection(ctx, &qdrant.CreateCollection{
    CollectionName: "large_collection",
    VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
        Size:     1536,
        Distance: qdrant.Distance_Cosine,
        QuantizationConfig: &qdrant.QuantizationConfig{
            Quantization: &qdrant.QuantizationConfig_Scalar{
                Scalar: &qdrant.ScalarQuantization{
                    Type:     qdrant.QuantizationType_Int8,
                    Quantile: qdrant.PtrOf(float32(0.99)),
                    AlwaysRam: qdrant.PtrOf(true), // keep quantized index in RAM
                },
            },
        },
    }),
})
```

Query with quantization oversampling (fetch more candidates, rescore with originals):

```go
results, err := client.Query(ctx, &qdrant.QueryPoints{
    CollectionName: "large_collection",
    Query:          qdrant.NewQuery(queryVector...),
    Limit:          qdrant.PtrOf(uint64(10)),
    Params: &qdrant.SearchParams{
        QuantizationSearchParams: &qdrant.QuantizationSearchParams{
            Rescore:    qdrant.PtrOf(true),        // rescore with original vectors
            Oversampling: qdrant.PtrOf(float64(2.0)), // fetch 2x candidates before rescoring
        },
    },
})
```

---

## 6. Batch Operations and Performance

**Batch upsert (critical for ingestion speed):**

```go
package ingestion

import (
    "context"
    "log"

    "github.com/qdrant/go-client/qdrant"
)

const batchSize = 100

func BatchIngest(client *qdrant.Client, collection string, docs []Document) error {
    ctx := context.Background()

    for i := 0; i < len(docs); i += batchSize {
        end := i + batchSize
        if end > len(docs) {
            end = len(docs)
        }
        batch := docs[i:end]

        points := make([]*qdrant.PointStruct, len(batch))
        for j, doc := range batch {
            points[j] = &qdrant.PointStruct{
                Id:      qdrant.NewIDNum(uint64(doc.ID)),
                Vectors: qdrant.NewVectors(doc.Embedding...),
                Payload: qdrant.NewValueMap(map[string]any{
                    "title":    doc.Title,
                    "category": doc.Category,
                    "price":    doc.Price,
                }),
            }
        }

        _, err := client.Upsert(ctx, &qdrant.UpsertPoints{
            CollectionName: collection,
            Points:         points,
            Wait:           qdrant.PtrOf(false), // async — don't wait for index update
        })
        if err != nil {
            return err
        }

        log.Printf("ingested %d/%d", end, len(docs))
    }
    return nil
}
```

**Optimizing HNSW for ingestion:**

```go
// Before bulk ingestion: reduce ef_construction for speed
client.UpdateCollection(ctx, &qdrant.UpdateCollection{
    CollectionName: "products",
    OptimizersConfig: &qdrant.OptimizersConfigDiff{
        IndexingThreshold: qdrant.PtrOf(uint64(100_000)), // don't build HNSW until 100k points
    },
    HnswConfig: &qdrant.HnswConfigDiff{
        EfConstruct: qdrant.PtrOf(uint64(64)), // faster build
    },
})

// ... bulk ingest all data ...

// After ingestion: restore quality settings
client.UpdateCollection(ctx, &qdrant.UpdateCollection{
    CollectionName: "products",
    OptimizersConfig: &qdrant.OptimizersConfigDiff{
        IndexingThreshold: qdrant.PtrOf(uint64(20_000)), // resume auto-indexing
    },
    HnswConfig: &qdrant.HnswConfigDiff{
        EfConstruct: qdrant.PtrOf(uint64(200)), // higher quality index
    },
})
```

---

## 7. Collection Management

```go
// List all collections
collections, _ := client.ListCollections(ctx)
for _, c := range collections {
    fmt.Println(c)
}

// Get collection info (size, segments, etc.)
info, _ := client.GetCollection(ctx, "products")
fmt.Printf("Points: %d\n", info.GetPointsCount())
fmt.Printf("Segments: %d\n", info.GetSegmentsCount())
fmt.Printf("Status: %s\n", info.GetStatus())

// Create a snapshot for backup
snapshot, _ := client.CreateSnapshot(ctx, "products")
fmt.Printf("Snapshot: %s\n", snapshot.GetName())

// Delete collection
client.DeleteCollection(ctx, "old_products")

// Scroll through all points (for export/migration)
var offset *qdrant.PointId
for {
    results, nextOffset, err := client.Scroll(ctx, &qdrant.ScrollPoints{
        CollectionName: "products",
        Limit:          qdrant.PtrOf(uint32(100)),
        Offset:         offset,
        WithPayload:    qdrant.NewWithPayload(true),
        WithVectors:    qdrant.NewWithVectorsEnable(true),
    })
    if err != nil || len(results) == 0 {
        break
    }
    // process results...
    offset = nextOffset
}
```

---

## 8. Project: AI-Powered Product Search

A complete Go HTTP server providing semantic product search with category filtering and hybrid ranking.

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/qdrant/go-client/qdrant"
)

type Product struct {
    ID       uint64  `json:"id"`
    Name     string  `json:"name"`
    Category string  `json:"category"`
    Price    float64 `json:"price"`
    InStock  bool    `json:"in_stock"`
}

type SearchRequest struct {
    Query    string  `json:"query"`
    Category string  `json:"category,omitempty"`
    MaxPrice float64 `json:"max_price,omitempty"`
    Limit    uint64  `json:"limit,omitempty"`
}

type SearchResponse struct {
    Results []SearchHit `json:"results"`
    Total   int         `json:"total"`
}

type SearchHit struct {
    Product
    Score float32 `json:"score"`
}

type Server struct {
    qdrant *qdrant.Client
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
    var req SearchRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }
    if req.Limit == 0 {
        req.Limit = 10
    }

    // Get embedding for query
    embedding, err := getEmbedding(os.Getenv("OPENAI_KEY"), req.Query)
    if err != nil {
        http.Error(w, "embedding error", 500)
        return
    }

    // Build filter
    var mustConditions []*qdrant.Condition
    mustConditions = append(mustConditions, qdrant.NewMatchBool("in_stock", true))

    if req.Category != "" {
        mustConditions = append(mustConditions, qdrant.NewMatch("category", req.Category))
    }
    if req.MaxPrice > 0 {
        mustConditions = append(mustConditions, qdrant.NewRange("price", &qdrant.Range{
            Lt: qdrant.PtrOf(req.MaxPrice),
        }))
    }

    qr := &qdrant.QueryPoints{
        CollectionName: "products",
        Query:          qdrant.NewQuery(embedding...),
        Limit:          qdrant.PtrOf(req.Limit),
        WithPayload:    qdrant.NewWithPayload(true),
    }
    if len(mustConditions) > 0 {
        qr.Filter = &qdrant.Filter{Must: mustConditions}
    }

    results, err := s.qdrant.Query(r.Context(), qr)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    var hits []SearchHit
    for _, res := range results {
        p := res.GetPayload()
        hits = append(hits, SearchHit{
            Product: Product{
                ID:       res.GetId().GetNum(),
                Name:     p["name"].GetStringValue(),
                Category: p["category"].GetStringValue(),
                Price:    p["price"].GetDoubleValue(),
                InStock:  p["in_stock"].GetBoolValue(),
            },
            Score: res.GetScore(),
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SearchResponse{Results: hits, Total: len(hits)})
}

func main() {
    client, err := qdrant.NewClient(&qdrant.Config{Host: "localhost", Port: 6334})
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    s := &Server{qdrant: client}

    // Setup collection on startup
    setupCollection(context.Background(), client)

    mux := http.NewServeMux()
    mux.HandleFunc("POST /search", s.Search)

    log.Println("Product search API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func setupCollection(ctx context.Context, client *qdrant.Client) {
    err := client.CreateCollection(ctx, &qdrant.CreateCollection{
        CollectionName: "products",
        VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
            Size:     1536,
            Distance: qdrant.Distance_Cosine,
            QuantizationConfig: &qdrant.QuantizationConfig{
                Quantization: &qdrant.QuantizationConfig_Scalar{
                    Scalar: &qdrant.ScalarQuantization{
                        Type:      qdrant.QuantizationType_Int8,
                        AlwaysRam: qdrant.PtrOf(true),
                    },
                },
            },
        }),
    })
    if err != nil {
        log.Println("collection may already exist:", err)
        return
    }

    // Create payload indexes
    for _, field := range []string{"category", "price", "in_stock"} {
        ft := qdrant.FieldType_FieldTypeKeyword
        if field == "price" {
            ft = qdrant.FieldType_FieldTypeFloat
        } else if field == "in_stock" {
            ft = qdrant.FieldType_FieldTypeBool
        }
        client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
            CollectionName: "products",
            FieldName:      field,
            FieldType:      qdrant.PtrOf(ft),
        })
    }
}
```

Test it:
```bash
# Search running shoes under $150
curl -X POST localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{"query":"lightweight running shoes for trails","category":"footwear","max_price":150,"limit":5}'

# Semantic search finds "Trail Running Shoe X" even if the description says
# "mountain jogging footwear" — because they're semantically similar.
```

---

## Summary

- Use the high-level `qdrant.NewClient()` for a clean API; it uses gRPC under the hood.
- **Named vectors** enable one point to have multiple embeddings (text + image) — search by either independently.
- Always create **payload indexes** on fields you filter by. Without indexes, filtering scans all points.
- **Hybrid search** (dense + sparse + RRF) beats either approach alone — use it for production search.
- **Scalar quantization** is the best default: 4x memory reduction, <1% accuracy loss, transparent to queries.
- For bulk ingestion: raise `indexing_threshold`, bulk insert, then restore thresholds to trigger index build.

### Exercises

**Easy:** Set up Qdrant with Docker. Create a `movies` collection with `title`, `genre`, and `year` payload fields. Insert 5 movies with fake 3-dim vectors. Run a filtered search for movies with `genre = "action"`.

**Medium:** Build a multi-modal image+text search: create a collection with `"text"` (384-dim) and `"image"` (128-dim) named vectors. Insert 10 items. Build a Go HTTP server with two endpoints: `/search/text` and `/search/image` that use the appropriate named vector.

**Hard:** Build a complete RAG pipeline in Go: (1) Chunk a 50-page text file into 300-word segments, (2) embed each chunk with a local model (Ollama), (3) store in Qdrant with `source`, `page`, `chunk_index` payload, (4) build an HTTP endpoint that accepts a question, finds the top 3 relevant chunks via Qdrant, feeds them to an LLM, and returns a grounded answer.
