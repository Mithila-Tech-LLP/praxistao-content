# Chapter 33: Vector Databases — Semantic Search and AI Applications

Normal databases find exact matches: "find user with email = 'alice@example.com'". Vector databases find similar things: "find the 10 products most similar to this one" or "find documents that mean the same thing as this question". This chapter explains how that works and why every AI application needs it.

## Table of Contents

1. What is a Vector? (From Coordinates to Meaning)
2. Why Traditional Databases Can't Do This
3. How Vector Search Works — Approximate Nearest Neighbor
4. Popular Vector Databases
5. Embeddings — Turning Text/Images into Vectors
6. Use Cases for Vector Databases
7. Exercises

---

## 1. What is a Vector? (From Coordinates to Meaning)

A vector is just a list of numbers. You already know vectors from school: the point (3, 4) is a vector with 2 dimensions. Databases use vectors with hundreds or thousands of dimensions to represent meaning.

**How does meaning become numbers?**

Imagine placing words on a map where similar words are near each other:

```
         happy ●
      excited ●  ● joyful
                        ● pleased
                               ● content
    sad ●
      ● depressed
           ● miserable
```

In this 2D example:
- "happy" and "joyful" are close → similar meaning
- "happy" and "sad" are far apart → different meaning

Real embedding models do this in 1,536 dimensions (OpenAI) or 768 dimensions (many open-source models). The result is a vector like:

```
"happy" → [0.023, -0.156, 0.891, 0.042, -0.234, ...]  (768 numbers)
"joyful" → [0.019, -0.147, 0.883, 0.039, -0.228, ...]  (very similar!)
"sad"   → [-0.456, 0.234, -0.123, 0.567, 0.089, ...]   (very different)
```

The closeness between two vectors is measured with **cosine similarity** — a number between -1 and 1 where 1 means identical, 0 means unrelated, -1 means opposite.

---

## 2. Why Traditional Databases Can't Do This

**The problem:** Find the 10 vectors most similar to this query vector from a collection of 10 million vectors.

**Naive approach:** Compare the query vector to all 10 million vectors.

```
10 million vectors × 768 dimensions × cosine similarity = very slow
```

At 1 million vectors/second, that's 10 seconds per query. Not acceptable for a search box.

**Why indexes don't help:**

A B-Tree index works by ordering values: 1 < 2 < 3. You can quickly find "all numbers between 5 and 10".

Vectors don't have a simple ordering. You can't sort 768-dimensional points in a line. Closeness in high-dimensional space is complex — no B-Tree can answer "what are the 10 nearest points?" efficiently.

This is why you need specialized vector indexes.

---

## 3. How Vector Search Works — Approximate Nearest Neighbor

**Exact Nearest Neighbor (brute force):** Compare to all vectors. Perfect accuracy, but O(n) time.

**Approximate Nearest Neighbor (ANN):** Find vectors that are very likely to be the nearest, not guaranteed. 99% as accurate, 1000x faster.

### HNSW (Hierarchical Navigable Small World)

The most popular ANN algorithm. Think of it as a skip list but for high-dimensional space:

```
Layer 2 (few nodes, long jumps):    A ────────────────────── E
                                    |                         |
Layer 1 (medium density):          A ──── B ──── D ───────── E
                                   |      |      |            |
Layer 0 (all nodes, small steps):  A ─ B ─ C ─ D ─ E ─ F ─ G
```

Search starts at the top layer (few nodes, long distances) and drills down to the bottom layer for fine-grained neighbors. This gives O(log n) approximate search.

**HNSW parameters:**
- `M` (connections per node): Higher = better accuracy, more memory
- `ef_construction` (build-time exploration): Higher = better index quality, slower to build
- `ef` (query-time exploration): Higher = better accuracy per query, slower queries

### IVF (Inverted File Index)

Clusters vectors into groups. To search: find the nearest cluster centers first, then search within those clusters.

```
All vectors split into 1000 clusters
Query → Find nearest 5 cluster centers → Search only those clusters
Result: search 0.5% of all vectors instead of 100%
```

---

## 4. Popular Vector Databases

| Database | Architecture | Best For |
|----------|-------------|---------|
| **pgvector** | PostgreSQL extension | Already have PostgreSQL; moderate scale |
| **Qdrant** | Purpose-built, Rust | Production AI apps; excellent Go client |
| **Weaviate** | Purpose-built, Go | GraphQL API; knowledge graphs |
| **Pinecone** | Managed cloud | Zero ops; high scale |
| **Chroma** | Embedded Python | Prototyping; dev environments |
| **Milvus** | Distributed | Massive scale (billions of vectors) |

### pgvector (PostgreSQL Extension)

Best starting point if you already use PostgreSQL:

```sql
-- Install extension
CREATE EXTENSION vector;

-- Create a table with a vector column (1536 dimensions for OpenAI embeddings)
CREATE TABLE documents (
    id      SERIAL PRIMARY KEY,
    content TEXT,
    embedding vector(1536)
);

-- Create HNSW index for fast similarity search
CREATE INDEX ON documents USING hnsw (embedding vector_cosine_ops);

-- Insert a document with its embedding
INSERT INTO documents (content, embedding)
VALUES ('Go is a statically typed language', '[0.023, -0.156, ...]');

-- Find 5 most similar documents to a query vector
SELECT content, 1 - (embedding <=> '[0.019, -0.147, ...]') AS similarity
FROM documents
ORDER BY embedding <=> '[0.019, -0.147, ...]'  -- <=> is cosine distance
LIMIT 5;
```

### Qdrant

Purpose-built vector database with excellent performance:

```bash
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

Features:
- HNSW index with payload filtering (filter by metadata while doing vector search)
- Named vectors (multiple embeddings per document)
- Sparse + dense vector hybrid search
- Snapshots for backups

---

## 5. Embeddings — Turning Text/Images into Vectors

An embedding model takes raw data (text, image, audio) and outputs a vector. The model is trained so that semantically similar inputs produce similar vectors.

**Popular embedding models:**

| Model | Dimensions | Notes |
|-------|-----------|-------|
| OpenAI text-embedding-3-small | 1536 | Excellent quality, paid API |
| OpenAI text-embedding-3-large | 3072 | Best quality, more expensive |
| sentence-transformers/all-MiniLM-L6 | 384 | Free, open-source, runs locally |
| nomic-embed-text | 768 | Free, open-source, very good quality |

**How to get embeddings in Go:**

```go
package embedding

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

type EmbeddingResponse struct {
    Data []struct {
        Embedding []float32 `json:"embedding"`
    } `json:"data"`
}

func GetEmbedding(text string) ([]float32, error) {
    body, _ := json.Marshal(map[string]interface{}{
        "model": "text-embedding-3-small",
        "input": text,
    })

    req, _ := http.NewRequest("POST", "https://api.openai.com/v1/embeddings",
        bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result EmbeddingResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    if len(result.Data) == 0 {
        return nil, fmt.Errorf("no embeddings returned")
    }
    return result.Data[0].Embedding, nil
}
```

---

## 6. Use Cases for Vector Databases

**Semantic Search:**
Instead of keyword matching (user searches "running shoes" → find "jogging sneakers" too because they're semantically similar).

**RAG (Retrieval-Augmented Generation):**
Your private documents + an AI chatbot. Embed all your docs, store in vector DB. When user asks a question: find the most relevant doc chunks (vector search), feed them to an LLM, get an answer grounded in your data.

**Recommendations:**
"Users who bought X also liked Y" — embed user purchase history as a vector, find similar users, recommend what they bought.

**Duplicate Detection:**
Embed all images/texts. Two inputs with very similar vectors = likely duplicates.

**Anomaly Detection:**
Embed log lines. Log lines very dissimilar from all known normal logs = potential security incident.

---

## Summary

- Vectors are lists of numbers that represent meaning. Similar things have vectors that are mathematically close.
- Traditional B-Tree indexes can't do nearest-neighbor search in high dimensions.
- HNSW is the standard ANN algorithm: fast, memory-efficient, ~99% accuracy.
- Start with pgvector if you already use PostgreSQL. Use Qdrant for a dedicated vector store.
- Embeddings convert text/images into vectors using a trained model. OpenAI API is easiest; open-source models run locally for free.

### Exercises

**Easy:** Explain in your own words why cosine similarity is better than Euclidean distance for comparing text embeddings. (Hint: think about what happens if one text is twice as long as another.)

**Medium:** Set up pgvector with Docker and PostgreSQL. Create a `notes` table with an `embedding vector(384)` column. Insert 5 fake notes with manually entered vectors. Run a similarity search.

**Hard:** Get the sentence-transformers library running locally (via Python or Ollama). Embed 100 Wikipedia article summaries. Store in pgvector. Build a Go HTTP server that accepts a question and returns the 3 most semantically similar articles.
