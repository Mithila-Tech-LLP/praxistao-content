# Chapter 34: Building with Vector Databases in Go

A complete semantic search system in Go using pgvector and Qdrant. We'll build a document search engine that understands meaning, not just keywords.

## Table of Contents

1. pgvector in Go — Setup and CRUD
2. Semantic Search with pgvector
3. Qdrant in Go — Collections and Points
4. Hybrid Search (Vector + Filters)
5. Mini Project: Document Search Engine
6. Exercises

---

## 1. pgvector in Go — Setup and CRUD

```bash
# Start PostgreSQL with pgvector
docker run -d \
  -e POSTGRES_PASSWORD=secret \
  -p 5432:5432 \
  pgvector/pgvector:pg16

go get github.com/pgvector/pgvector-go
go get github.com/jackc/pgx/v5
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5"
    "github.com/pgvector/pgvector-go"
    pgxvector "github.com/pgvector/pgvector-go/pgx"
)

func main() {
    ctx := context.Background()

    conn, err := pgx.Connect(ctx, "postgres://postgres:secret@localhost/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close(ctx)

    // Register the vector type with pgx
    if err := pgxvector.RegisterTypes(ctx, conn); err != nil {
        log.Fatal(err)
    }

    // Create table
    conn.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector`)
    conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS documents (
        id        SERIAL PRIMARY KEY,
        title     TEXT NOT NULL,
        content   TEXT NOT NULL,
        embedding vector(3)  -- tiny for demo; use 1536 in prod
    )`)
    conn.Exec(ctx, `CREATE INDEX IF NOT EXISTS docs_embedding_idx
        ON documents USING hnsw (embedding vector_cosine_ops)`)

    // Insert with embedding
    v := pgvector.NewVector([]float32{0.1, 0.5, 0.9})
    _, err = conn.Exec(ctx,
        `INSERT INTO documents (title, content, embedding) VALUES ($1, $2, $3)`,
        "Go Programming", "Go is a statically typed language", v,
    )
    if err != nil {
        log.Fatal("insert:", err)
    }

    // Similarity search
    query := pgvector.NewVector([]float32{0.15, 0.45, 0.88})
    rows, _ := conn.Query(ctx, `
        SELECT id, title, 1 - (embedding <=> $1) AS similarity
        FROM documents
        ORDER BY embedding <=> $1
        LIMIT 5
    `, query)

    for rows.Next() {
        var id int
        var title string
        var similarity float64
        rows.Scan(&id, &title, &similarity)
        fmt.Printf("ID: %d | %s | similarity: %.4f\n", id, title, similarity)
    }
}
```

---

## 2. Semantic Search with pgvector

A complete repository pattern for document storage and retrieval:

```go
package vectordb

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    pgvector "github.com/pgvector/pgvector-go"
    pgxvector "github.com/pgvector/pgvector-go/pgx"
)

type Document struct {
    ID        int
    Title     string
    Content   string
    Source    string
    Embedding []float32
}

type SearchResult struct {
    Document
    Similarity float64
}

type DocRepository struct {
    pool *pgxpool.Pool
}

func NewDocRepository(pool *pgxpool.Pool) (*DocRepository, error) {
    ctx := context.Background()
    conn, err := pool.Acquire(ctx)
    if err != nil {
        return nil, err
    }
    defer conn.Release()
    if err := pgxvector.RegisterTypes(ctx, conn.Conn()); err != nil {
        return nil, err
    }
    return &DocRepository{pool: pool}, nil
}

func (r *DocRepository) Upsert(ctx context.Context, doc Document) (int, error) {
    v := pgvector.NewVector(doc.Embedding)
    var id int
    err := r.pool.QueryRow(ctx, `
        INSERT INTO documents (title, content, source, embedding)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (source) DO UPDATE
          SET title = EXCLUDED.title,
              content = EXCLUDED.content,
              embedding = EXCLUDED.embedding
        RETURNING id
    `, doc.Title, doc.Content, doc.Source, v).Scan(&id)
    return id, err
}

func (r *DocRepository) Search(ctx context.Context, queryVec []float32, limit int) ([]SearchResult, error) {
    q := pgvector.NewVector(queryVec)
    rows, err := r.pool.Query(ctx, `
        SELECT id, title, content, source,
               1 - (embedding <=> $1) AS similarity
        FROM documents
        ORDER BY embedding <=> $1
        LIMIT $2
    `, q, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []SearchResult
    for rows.Next() {
        var sr SearchResult
        err := rows.Scan(&sr.ID, &sr.Title, &sr.Content, &sr.Source, &sr.Similarity)
        if err != nil {
            return nil, err
        }
        results = append(results, sr)
    }
    return results, rows.Err()
}

func (r *DocRepository) FilteredSearch(
    ctx context.Context,
    queryVec []float32,
    source string,
    minSimilarity float64,
    limit int,
) ([]SearchResult, error) {
    q := pgvector.NewVector(queryVec)
    rows, err := r.pool.Query(ctx, `
        SELECT id, title, content, source,
               1 - (embedding <=> $1) AS similarity
        FROM documents
        WHERE source = $2
          AND (1 - (embedding <=> $1)) >= $3
        ORDER BY embedding <=> $1
        LIMIT $4
    `, q, source, minSimilarity, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []SearchResult
    for rows.Next() {
        var sr SearchResult
        rows.Scan(&sr.ID, &sr.Title, &sr.Content, &sr.Source, &sr.Similarity)
        results = append(results, sr)
    }
    return results, rows.Err()
}
```

---

## 3. Qdrant in Go — Collections and Points

Qdrant is a purpose-built vector database with a richer feature set than pgvector.

```bash
docker run -d -p 6333:6333 qdrant/qdrant
go get github.com/qdrant/go-client
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    pb "github.com/qdrant/go-client/qdrant"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    conn, err := grpc.NewClient("localhost:6334",
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewCollectionsClient(conn)
    pointsClient := pb.NewPointsClient(conn)
    ctx := context.Background()

    // Create collection
    _, err = client.Create(ctx, &pb.CreateCollection{
        CollectionName: "docs",
        VectorsConfig: &pb.VectorsConfig{
            Config: &pb.VectorsConfig_Params{
                Params: &pb.VectorParams{
                    Size:     3, // 1536 in production
                    Distance: pb.Distance_Cosine,
                },
            },
        },
    })
    if err != nil {
        log.Println("collection may already exist:", err)
    }

    // Insert points
    _, err = pointsClient.Upsert(ctx, &pb.UpsertPoints{
        CollectionName: "docs",
        Points: []*pb.PointStruct{
            {
                Id: &pb.PointId{PointIdOptions: &pb.PointId_Num{Num: 1}},
                Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{
                    Vector: &pb.Vector{Data: []float32{0.1, 0.5, 0.9}},
                }},
                Payload: map[string]*pb.Value{
                    "title":   {Kind: &pb.Value_StringValue{StringValue: "Go Programming"}},
                    "content": {Kind: &pb.Value_StringValue{StringValue: "Go is a statically typed language"}},
                    "source":  {Kind: &pb.Value_StringValue{StringValue: "golang.org"}},
                },
            },
            {
                Id: &pb.PointId{PointIdOptions: &pb.PointId_Num{Num: 2}},
                Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{
                    Vector: &pb.Vector{Data: []float32{0.8, 0.2, 0.1}},
                }},
                Payload: map[string]*pb.Value{
                    "title":   {Kind: &pb.Value_StringValue{StringValue: "Python Basics"}},
                    "content": {Kind: &pb.Value_StringValue{StringValue: "Python is a dynamically typed language"}},
                    "source":  {Kind: &pb.Value_StringValue{StringValue: "python.org"}},
                },
            },
        },
    })
    if err != nil {
        log.Fatal("upsert:", err)
    }

    // Search
    limit := uint64(5)
    results, err := pointsClient.Search(ctx, &pb.SearchPoints{
        CollectionName: "docs",
        Vector:         []float32{0.15, 0.45, 0.88},
        Limit:          limit,
        WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
    })
    if err != nil {
        log.Fatal("search:", err)
    }

    for _, r := range results.GetResult() {
        title := r.GetPayload()["title"].GetStringValue()
        fmt.Printf("Score: %.4f | %s\n", r.GetScore(), title)
    }
}
```

---

## 4. Hybrid Search (Vector + Filters)

Real-world search combines vector similarity with metadata filters. "Find documents about 'database performance' but only from the 'engineering' category":

```go
// Qdrant filtered search
results, err := pointsClient.Search(ctx, &pb.SearchPoints{
    CollectionName: "docs",
    Vector:         queryVector,
    Limit:          10,
    Filter: &pb.Filter{
        Must: []*pb.Condition{
            {
                ConditionOneOf: &pb.Condition_Field{
                    Field: &pb.FieldCondition{
                        Key: "category",
                        Match: &pb.Match{
                            MatchValue: &pb.Match_Keyword{Keyword: "engineering"},
                        },
                    },
                },
            },
        },
    },
    WithPayload: &pb.WithPayloadSelector{
        SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true},
    },
})
```

---

## 5. Mini Project: Document Search Engine

A complete semantic search API that ingests documents and answers natural-language questions.

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
)

type SearchHandler struct {
    repo *DocRepository
}

type IngestRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
    Source  string `json:"source"`
}

type SearchRequest struct {
    Query string `json:"query"`
    Limit int    `json:"limit"`
}

func main() {
    pool, err := pgxpool.New(context.Background(),
        "postgres://postgres:secret@localhost/mydb")
    if err != nil {
        log.Fatal(err)
    }

    repo, err := NewDocRepository(pool)
    if err != nil {
        log.Fatal(err)
    }

    h := &SearchHandler{repo: repo}

    mux := http.NewServeMux()
    mux.HandleFunc("POST /ingest", h.Ingest)
    mux.HandleFunc("POST /search", h.Search)

    log.Println("Document Search API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func (h *SearchHandler) Ingest(w http.ResponseWriter, r *http.Request) {
    var req IngestRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }

    embedding, err := GetEmbedding(req.Title + " " + req.Content)
    if err != nil {
        http.Error(w, "embedding error: "+err.Error(), 500)
        return
    }

    id, err := h.repo.Upsert(r.Context(), Document{
        Title:     req.Title,
        Content:   req.Content,
        Source:    req.Source,
        Embedding: embedding,
    })
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
    var req SearchRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }
    if req.Limit == 0 {
        req.Limit = 5
    }

    embedding, err := GetEmbedding(req.Query)
    if err != nil {
        http.Error(w, "embedding error: "+err.Error(), 500)
        return
    }

    results, err := h.repo.Search(r.Context(), embedding, req.Limit)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(results)
}
```

Test:
```bash
# Ingest documents
curl -X POST localhost:8080/ingest \
  -H "Content-Type: application/json" \
  -d '{"title":"Go Goroutines","content":"Goroutines are lightweight threads managed by the Go runtime","source":"golang.org"}'

# Semantic search
curl -X POST localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{"query":"concurrent programming in Go","limit":3}'
# Returns the goroutines article even though "goroutines" wasn't in the query!
```

---

## Summary

- `pgvector` adds vector search to PostgreSQL with the `vector` type and `<=>` cosine distance operator.
- Use `USING hnsw` index for fast approximate nearest-neighbor search.
- Qdrant is a purpose-built vector DB with payload filtering, richer query API, and gRPC interface.
- Hybrid search = vector similarity + metadata filters — the standard production pattern.
- Always embed the concatenation of multiple fields (title + content) for richer semantic representation.

### Exercises

**Easy:** Set up pgvector. Create a `quotes` table with `text` and `embedding vector(384)`. Manually insert 3 quotes with made-up vectors. Run a similarity search.

**Medium:** Build a "FAQ bot" backend: ingest 10 FAQ entries into pgvector. Given a user question, find the 3 most semantically similar FAQs and return their answers.

**Hard:** Build a RAG (Retrieval-Augmented Generation) system: ingest a local text file (e.g., a Go book chapter) split into 300-word chunks. For a user question, retrieve the top 3 chunks from pgvector, then send them as context to the OpenAI chat API to get an answer grounded in your document.
