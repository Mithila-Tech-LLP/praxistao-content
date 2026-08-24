# Chapter 34: Vector Databases

> **"A vector database is not just storage — it's a spatial index that turns semantic search from O(N) brute force into something that scales to billions of vectors."**

---

## Table of Contents
1. [Why Vector Databases Exist](#1-why-vector-databases-exist)
2. [Approximate Nearest Neighbor Search](#2-approximate-nearest-neighbor-search)
3. [IVF Indexing](#3-ivf-indexing)
4. [HNSW Indexing](#4-hnsw-indexing)
5. [Similarity Metrics](#5-similarity-metrics)
6. [FAISS — The Building Block](#6-faiss--the-building-block)
7. [ChromaDB — Local First](#7-chromadb--local-first)
8. [Other Options: Pinecone, pgvector](#8-other-options-pinecone-pgvector)
9. [Filtering and Hybrid Search](#9-filtering-and-hybrid-search)
10. [Complete Search System](#10-complete-search-system)
11. [Mini Projects](#11-mini-projects)
12. [Summary and Exercises](#12-summary-and-exercises)

---

## Before You Start

**Prerequisites:** Chapter 33 (Embeddings and Semantic Search)

```bash
pip install faiss-cpu chromadb sentence-transformers numpy pandas
```

---

## 1. Why Vector Databases Exist

### The Problem with Regular Databases

SQL databases index on exact values (B-trees, hash indexes). They cannot efficiently answer "find me 10 rows that are *similar* to this vector."

```
SQL APPROACH (doesn't work):
  WHERE embedding ≈ [0.2, -0.1, 0.5, ...]  ← SQL can't do this

NAIVE APPROACH (works but too slow):
  For each of 1 million documents:
    similarity = cosine(query_vec, doc_vec)  ← 1M computations
  Sort and return top 10

COST:
  1M documents × 384 dimensions × float32 = 1.5 GB to scan
  At 1 GB/s memory bandwidth: 1.5 seconds per query
  
  With 1 BILLION documents: 1,500 seconds. Unusable!
```

Vector databases solve this with **Approximate Nearest Neighbor (ANN)** algorithms that can search billions of vectors in milliseconds.

---

## 2. Approximate Nearest Neighbor Search

```
EXACT vs APPROXIMATE:

Exact NN: guaranteed to find the true nearest neighbor
  → Must scan all vectors
  → O(N) time, O(N×d) space
  → Unusable at scale

Approximate NN: finds vectors that are "close enough"
  → Trade a small accuracy loss for enormous speed gain
  → Common recall target: 99% (miss only 1% of true answers)
  → 100-1000× faster than exact search

RECALL:
  recall@10 = (true answers in our top-10) / 10
  
  If true top-10 is: [A, B, C, D, E, F, G, H, I, J]
  ANN returns:       [A, B, C, D, E, F, G, H, K, L]  ← missed I, J
  recall@10 = 8/10 = 0.8 (80% recall)

For most RAG/search applications, 95-99% recall is fine.
Users rarely notice if the 10th result is slightly suboptimal.
```

---

## 3. IVF Indexing

IVF (Inverted File Index) partitions vectors into clusters, then only searches nearby clusters.

```
IVF ALGORITHM:

BUILD PHASE:
  1. Run k-means on all N vectors → N_clusters clusters
  2. Assign each vector to nearest cluster centroid
  3. Build an inverted file: cluster_id → [vector_ids in that cluster]

  Cluster 1: [vec_42, vec_891, vec_3421, ...]
  Cluster 2: [vec_17, vec_288, vec_9001, ...]
  ...
  Cluster 256: [vec_71, vec_432, ...]

QUERY PHASE:
  1. Find query's nearest nprobe cluster centroids (e.g., nprobe=8)
  2. Only search vectors in those 8 clusters
  3. Return top-k from those candidates
  
  WHY IT WORKS:
  If query is near cluster 5, most similar vectors are likely
  IN cluster 5 or neighboring clusters. Skip the other 248 clusters!

TRADEOFFS:
  • nprobe=1: fastest, lowest recall
  • nprobe=N_clusters: exact search (defeats the purpose)
  • nprobe=8-32: good balance (recommended)
  
  N_clusters:
  • Rule of thumb: sqrt(N) for small N, N/1000 for large N
  • More clusters = faster search but more clusters to scan
```

---

## 4. HNSW Indexing

HNSW (Hierarchical Navigable Small World) builds a layered graph where each vector connects to its nearest neighbors.

```
HNSW STRUCTURE:

Layer 2 (sparse):  A ─── G ─── K
                   │
Layer 1:         A ─── C ─── G ─── I ─── K
                 │     │
Layer 0 (dense): A─B─C─D─E─F─G─H─I─J─K─L─M (all vectors)

QUERY ALGORITHM:
  1. Start at a fixed entry point in the top layer
  2. Greedy search: move to neighbor closest to query
  3. When stuck (no neighbor is closer), drop to lower layer
  4. Repeat until Layer 0
  5. Explore neighborhood in Layer 0 for top-k candidates

WHY IT'S FAST:
  Upper layers have few nodes → fast navigation
  Lower layers have more nodes → fine-grained search
  Like navigating a city: highway → main road → street

PARAMETERS:
  M: connections per node (16-64). Higher = better recall, more memory
  ef_construction: search quality during index build. Higher = slower build, better index
  ef_search: search quality during query. Higher = slower search, better recall
```

---

## 5. Similarity Metrics

```
COSINE SIMILARITY:
  Range: -1 to +1 (or 0 to 2 as distance)
  Best for: text embeddings, normalized vectors
  Formula: cos(θ) = A·B / (||A|| × ||B||)

DOT PRODUCT:
  Range: unbounded
  Best for: when you want magnitude to matter
  (larger vectors = more "confident" representations)
  Note: cosine sim = dot product when vectors are L2-normalized!

L2 (EUCLIDEAN) DISTANCE:
  Range: 0 to ∞ (0 = identical)
  Best for: image embeddings, when absolute positions matter
  Formula: sqrt(Σ(aᵢ - bᵢ)²)

INNER PRODUCT (same as dot product, different name in some libs)

CHOOSING:
  • You're using sentence-transformers? → Cosine (or dot with L2-normalized)
  • You're using OpenAI embeddings? → Cosine  
  • You're building image search? → L2 or dot
  • Default choice when unsure? → Cosine
```

---

## 6. FAISS — The Building Block

FAISS (Facebook AI Similarity Search) is the foundational library for vector search.

```python
import faiss
import numpy as np

# Generate sample data: 10,000 vectors of dimension 128
d = 128         # dimension
N = 10_000      # number of vectors

np.random.seed(42)
vectors = np.random.randn(N, d).astype('float32')

# L2-normalize for cosine similarity
faiss.normalize_L2(vectors)

# ── Option 1: Exact Search (small datasets < 100k) ──────────────
index_flat = faiss.IndexFlatIP(d)  # Inner Product = cosine for normalized vecs
index_flat.add(vectors)

query = np.random.randn(1, d).astype('float32')
faiss.normalize_L2(query)

k = 5
scores, indices = index_flat.search(query, k)
print("Exact search results:")
print(f"  Indices: {indices[0]}")
print(f"  Scores:  {scores[0]}")

# ── Option 2: IVF (medium datasets 100k - 10M) ──────────────────
n_clusters = int(np.sqrt(N))  # Rule of thumb: sqrt(N) clusters

quantizer = faiss.IndexFlatIP(d)  # Used to assign vectors to clusters
index_ivf = faiss.IndexIVFFlat(quantizer, d, n_clusters, faiss.METRIC_INNER_PRODUCT)

index_ivf.train(vectors)  # Learn cluster centroids
index_ivf.add(vectors)
index_ivf.nprobe = 8  # Search 8 nearest clusters per query

scores_ivf, indices_ivf = index_ivf.search(query, k)
print("\nIVF search results:")
print(f"  Indices: {indices_ivf[0]}")
print(f"  Scores:  {scores_ivf[0]}")

# ── Option 3: HNSW (fast, good quality, no training needed) ─────
index_hnsw = faiss.IndexHNSWFlat(d, 32, faiss.METRIC_INNER_PRODUCT)
index_hnsw.hnsw.efConstruction = 64  # Build quality (higher = better but slower)
index_hnsw.hnsw.efSearch = 64        # Search quality

index_hnsw.add(vectors)

scores_hnsw, indices_hnsw = index_hnsw.search(query, k)
print("\nHNSW search results:")
print(f"  Indices: {indices_hnsw[0]}")

# ── Save and load index ──────────────────────────────────────────
faiss.write_index(index_hnsw, "my_index.faiss")
loaded_index = faiss.read_index("my_index.faiss")

# ── GPU acceleration (if you have a GPU) ────────────────────────
# res = faiss.StandardGpuResources()
# gpu_index = faiss.index_cpu_to_gpu(res, 0, index_flat)
```

**FAISS Index Cheat Sheet:**

| Index Type | Size | Speed | Recall | Notes |
|---|---|---|---|---|
| `IndexFlatIP` | < 100k | Slow | 100% | Exact, no training |
| `IndexIVFFlat` | 100k–10M | Fast | 95-99% | Need to train |
| `IndexHNSWFlat` | Any | Fast | 98-99.9% | No training, more memory |
| `IndexIVFPQ` | > 10M | Fastest | 90-95% | Compressed, less memory |

---

## 7. ChromaDB — Local First

ChromaDB is the easiest way to get started. It handles embeddings, storage, and search in one package.

```python
import chromadb
from chromadb.utils import embedding_functions

# ── Setup ────────────────────────────────────────────────────────
# In-memory (no persistence):
client = chromadb.Client()

# Persistent (saves to disk):
client = chromadb.PersistentClient(path="./chroma_db")

# ── Embedding function ───────────────────────────────────────────
# ChromaDB can use sentence-transformers automatically
embed_fn = embedding_functions.SentenceTransformerEmbeddingFunction(
    model_name="all-MiniLM-L6-v2"
)

# ── Create a collection (like a table in SQL) ────────────────────
collection = client.get_or_create_collection(
    name="documents",
    embedding_function=embed_fn,
    metadata={"hnsw:space": "cosine"}  # Use cosine distance
)

# ── Add documents ────────────────────────────────────────────────
documents = [
    "Machine learning is a subset of artificial intelligence.",
    "Neural networks are inspired by the human brain.",
    "Deep learning enables computers to learn from raw data.",
    "Natural language processing deals with text understanding.",
    "Computer vision allows machines to interpret images.",
    "Reinforcement learning trains agents through rewards.",
]

# ChromaDB auto-generates IDs if you don't provide them
collection.add(
    documents=documents,
    ids=[f"doc_{i}" for i in range(len(documents))],
    metadatas=[
        {"topic": "ml_intro", "difficulty": "beginner"},
        {"topic": "deep_learning", "difficulty": "intermediate"},
        {"topic": "deep_learning", "difficulty": "beginner"},
        {"topic": "nlp", "difficulty": "intermediate"},
        {"topic": "cv", "difficulty": "intermediate"},
        {"topic": "rl", "difficulty": "advanced"},
    ]
)

print(f"Collection has {collection.count()} documents")

# ── Query by semantic similarity ─────────────────────────────────
results = collection.query(
    query_texts=["How do computers learn from data?"],
    n_results=3,
    include=["documents", "distances", "metadatas"]
)

print("\nSearch results:")
for doc, dist, meta in zip(
    results['documents'][0],
    results['distances'][0],
    results['metadatas'][0]
):
    print(f"  [{1-dist:.3f}] ({meta['topic']}) {doc}")

# ── Filter by metadata ────────────────────────────────────────────
results_filtered = collection.query(
    query_texts=["How do computers learn?"],
    n_results=3,
    where={"difficulty": "beginner"},  # Only beginner docs
)

print("\nFiltered results (beginner only):")
for doc in results_filtered['documents'][0]:
    print(f"  {doc}")

# ── Update and delete ─────────────────────────────────────────────
collection.update(
    ids=["doc_0"],
    documents=["Machine learning is a type of AI that learns from data."],
    metadatas=[{"topic": "ml_intro", "difficulty": "beginner", "updated": True}]
)

collection.delete(ids=["doc_5"])
print(f"Collection now has {collection.count()} documents")
```

---

## 8. Other Options: Pinecone, pgvector

### Pinecone (Managed Cloud)

```python
from pinecone import Pinecone, ServerlessSpec

pc = Pinecone(api_key="your-api-key")

# Create index
pc.create_index(
    name="my-index",
    dimension=384,
    metric="cosine",
    spec=ServerlessSpec(cloud="aws", region="us-east-1")
)

index = pc.Index("my-index")

# Upsert vectors
index.upsert(vectors=[
    ("id1", [0.1, 0.2, ...], {"category": "tech"}),
    ("id2", [0.3, 0.1, ...], {"category": "science"}),
])

# Query
results = index.query(vector=[0.1, 0.2, ...], top_k=5, filter={"category": "tech"})
```

**Use Pinecone when:** You need managed infrastructure, multi-region replication, or you don't want to maintain infrastructure.

### pgvector (PostgreSQL Extension)

```sql
-- Install the extension
CREATE EXTENSION vector;

-- Create a table with a vector column
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    content TEXT,
    embedding vector(384),  -- 384-dimensional vector
    category TEXT
);

-- Add an HNSW index
CREATE INDEX ON documents USING hnsw (embedding vector_cosine_ops);

-- Insert
INSERT INTO documents (content, embedding) VALUES 
    ('Machine learning intro', '[0.1, 0.2, ...]');

-- Semantic search query
SELECT content, 1 - (embedding <=> '[0.1, 0.2, ...]') as similarity
FROM documents
ORDER BY embedding <=> '[0.1, 0.2, ...]'
LIMIT 5;

-- Filter + search
SELECT content 
FROM documents
WHERE category = 'tech'
ORDER BY embedding <=> '[0.1, ...]'
LIMIT 5;
```

**Use pgvector when:** You already use PostgreSQL and want to add vector search without another service.

---

## 9. Filtering and Hybrid Search

```python
import chromadb
from chromadb.utils import embedding_functions

# ── Metadata Filtering ────────────────────────────────────────────
# Find similar documents, but only from a specific category/date/author
results = collection.query(
    query_texts=["deep learning techniques"],
    n_results=5,
    where={
        "$and": [
            {"topic": {"$in": ["deep_learning", "ml_intro"]}},
            {"difficulty": {"$ne": "advanced"}}
        ]
    }
)

# ── Hybrid Search (Dense + Sparse) ────────────────────────────────
# Combines semantic similarity (dense) with keyword matching (sparse/BM25)

from sklearn.feature_extraction.text import TfidfVectorizer
from scipy.sparse import csr_matrix
import numpy as np

def hybrid_search(query: str, documents: list, embeddings: np.ndarray,
                  model, alpha: float = 0.7, top_k: int = 5):
    """
    Combine semantic (dense) and keyword (sparse) search.
    
    alpha=1.0: pure semantic search
    alpha=0.0: pure keyword search
    alpha=0.7: 70% semantic, 30% keyword (usually best)
    """
    # Dense score (semantic)
    q_emb = model.encode([query], normalize_embeddings=True)
    dense_scores = (embeddings @ q_emb.T).squeeze()
    dense_scores = (dense_scores - dense_scores.min()) / (dense_scores.max() - dense_scores.min() + 1e-8)
    
    # Sparse score (BM25/TF-IDF)
    vectorizer = TfidfVectorizer()
    tfidf_matrix = vectorizer.fit_transform(documents)
    query_vec = vectorizer.transform([query])
    sparse_scores = (tfidf_matrix @ query_vec.T).toarray().squeeze()
    if sparse_scores.max() > 0:
        sparse_scores = sparse_scores / sparse_scores.max()
    
    # Combine
    combined = alpha * dense_scores + (1 - alpha) * sparse_scores
    top_indices = np.argsort(combined)[::-1][:top_k]
    
    return [(combined[i], documents[i]) for i in top_indices]
```

---

## 10. Complete Search System

```python
# complete_search_system.py
"""
End-to-end vector search system:
  • Embed documents with sentence-transformers
  • Store in ChromaDB
  • Search semantically
  • Filter by metadata
"""

from sentence_transformers import SentenceTransformer
import chromadb
from chromadb.utils import embedding_functions
from pathlib import Path
import json

class DocumentSearchSystem:
    def __init__(self, db_path: str = "./search_db"):
        embed_fn = embedding_functions.SentenceTransformerEmbeddingFunction(
            model_name="all-MiniLM-L6-v2"
        )
        self.client = chromadb.PersistentClient(path=db_path)
        self.collection = self.client.get_or_create_collection(
            name="documents",
            embedding_function=embed_fn,
            metadata={"hnsw:space": "cosine"}
        )
        print(f"Loaded collection with {self.collection.count()} documents")
    
    def index_documents(self, documents: list[dict]):
        """
        Index a list of documents.
        Each doc: {"id": str, "text": str, "metadata": dict}
        """
        if not documents:
            return
        
        # Avoid duplicates
        existing_ids = set()
        try:
            existing = self.collection.get()
            existing_ids = set(existing['ids'])
        except:
            pass
        
        new_docs = [d for d in documents if d['id'] not in existing_ids]
        
        if new_docs:
            self.collection.add(
                documents=[d['text'] for d in new_docs],
                ids=[d['id'] for d in new_docs],
                metadatas=[d.get('metadata', {}) for d in new_docs]
            )
            print(f"Indexed {len(new_docs)} new documents")
    
    def search(self, query: str, top_k: int = 5, filters: dict = None) -> list:
        """Search for relevant documents."""
        kwargs = {
            "query_texts": [query],
            "n_results": min(top_k, self.collection.count()),
            "include": ["documents", "distances", "metadatas", "ids"]
        }
        if filters:
            kwargs["where"] = filters
        
        results = self.collection.query(**kwargs)
        
        return [
            {
                "id": id_,
                "document": doc,
                "score": 1 - dist,  # Convert distance to similarity
                "metadata": meta
            }
            for id_, doc, dist, meta in zip(
                results['ids'][0],
                results['documents'][0],
                results['distances'][0],
                results['metadatas'][0]
            )
        ]
    
    def print_results(self, query: str, results: list):
        print(f"\nQuery: '{query}'")
        print("─"*60)
        for r in results:
            print(f"[{r['score']:.3f}] {r['document'][:80]}...")
            if r['metadata']:
                print(f"  metadata: {r['metadata']}")


# Demo usage
if __name__ == "__main__":
    system = DocumentSearchSystem()
    
    # Index some sample documents
    sample_docs = [
        {"id": "1", "text": "PyTorch is an open-source machine learning framework.", 
         "metadata": {"topic": "tools"}},
        {"id": "2", "text": "Transformers revolutionized natural language processing with attention mechanisms.",
         "metadata": {"topic": "architecture"}},
        {"id": "3", "text": "Training large language models requires significant compute resources.",
         "metadata": {"topic": "training"}},
        {"id": "4", "text": "Fine-tuning pre-trained models is more efficient than training from scratch.",
         "metadata": {"topic": "training"}},
        {"id": "5", "text": "Vector databases enable efficient semantic search at scale.",
         "metadata": {"topic": "infrastructure"}},
    ]
    
    system.index_documents(sample_docs)
    
    # Search
    for query in ["how to train AI models", "language understanding", "efficient search"]:
        results = system.search(query, top_k=2)
        system.print_results(query, results)
```

---

## 11. Mini Projects

### Mini Project 1: Movie Recommendation Engine

**What You'll Build:** Recommend movies based on plot similarity.

**Time Estimate:** 1-2 hours

```python
# movie_recommender.py
import pandas as pd
import chromadb
from chromadb.utils import embedding_functions

# Use TMDB 5000 movie dataset from Kaggle, or this small sample:
movies = [
    {"title": "Inception", "plot": "A thief who steals corporate secrets through dream-sharing technology is given the task of planting an idea into the mind of a CEO."},
    {"title": "The Matrix", "plot": "A computer hacker learns about the true nature of reality and joins a rebellion against its controllers."},
    {"title": "Interstellar", "plot": "A team of explorers travel through a wormhole in space to ensure humanity's survival."},
    {"title": "Toy Story", "plot": "A cowboy doll is threatened by a new spaceman figure, and the toys must work together."},
    {"title": "Finding Nemo", "plot": "After his son is captured by a diver, a clownfish sets out on a journey across the ocean."},
    {"title": "The Dark Knight", "plot": "Batman battles the Joker, a criminal mastermind who wants to plunge Gotham City into anarchy."},
    # Add more movies...
]

client = chromadb.Client()
embed_fn = embedding_functions.SentenceTransformerEmbeddingFunction("all-MiniLM-L6-v2")
col = client.create_collection("movies", embedding_function=embed_fn)

col.add(
    documents=[m['plot'] for m in movies],
    ids=[str(i) for i in range(len(movies))],
    metadatas=[{"title": m['title']} for m in movies]
)

def recommend(query_or_title: str, n: int = 3):
    results = col.query(query_texts=[query_or_title], n_results=n+1)
    print(f"\nMovies similar to '{query_or_title}':")
    for doc, meta, dist in zip(results['documents'][0], results['metadatas'][0], results['distances'][0]):
        if meta['title'] != query_or_title:
            print(f"  [{1-dist:.3f}] {meta['title']}: {doc[:80]}...")

recommend("Inception")
recommend("A story about friendship and adventure")
```

---

### Mini Project 2: Code Search Tool

**What You'll Build:** Search your codebase by describing what a function does.

```python
# code_search.py
import ast
import os
from pathlib import Path
import chromadb
from chromadb.utils import embedding_functions

def extract_functions(filepath: str) -> list[dict]:
    """Extract functions with their docstrings from a Python file."""
    try:
        tree = ast.parse(Path(filepath).read_text())
    except:
        return []
    
    functions = []
    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef):
            docstring = ast.get_docstring(node) or ""
            source_lines = Path(filepath).read_text().splitlines()
            # Get function signature
            sig = f"def {node.name}({', '.join(a.arg for a in node.args.args)})"
            
            if docstring:  # Only index functions with docstrings
                functions.append({
                    "name": node.name,
                    "signature": sig,
                    "docstring": docstring,
                    "file": filepath,
                    "line": node.lineno,
                    "text": f"{sig}\n{docstring}"
                })
    return functions

# Index your codebase
client = chromadb.Client()
embed_fn = embedding_functions.SentenceTransformerEmbeddingFunction("all-MiniLM-L6-v2")
col = client.create_collection("code", embedding_function=embed_fn)

# Scan Python files in current directory
all_functions = []
for py_file in Path(".").rglob("*.py"):
    all_functions.extend(extract_functions(str(py_file)))

if all_functions:
    col.add(
        documents=[f['text'] for f in all_functions],
        ids=[f"{f['file']}:{f['name']}" for f in all_functions],
        metadatas=[{"file": f['file'], "line": str(f['line'])} for f in all_functions]
    )
    print(f"Indexed {len(all_functions)} functions")

    # Search
    results = col.query(query_texts=["calculate cosine similarity between vectors"], n_results=3)
    for doc, meta in zip(results['documents'][0], results['metadatas'][0]):
        print(f"\nFile: {meta['file']} (line {meta['line']})")
        print(doc[:200])
```

---

## 12. Summary and Exercises

```
VECTOR DATABASE SELECTION GUIDE:
════════════════════════════════════════════════════════
Use Case                    Recommendation
──────────────────────────────────────────────────────
Learning / prototyping      ChromaDB (in-memory)
Local application           ChromaDB (persistent)
Already using PostgreSQL    pgvector
Production scale            Pinecone or Weaviate
Custom / max control        FAISS

Index Selection:
  < 100k vectors:           IndexFlatIP (exact, simple)
  100k–10M vectors:         HNSW (fast, no training)
  > 10M vectors:            IVF + PQ (compressed, scalable)
  
Similarity metric:
  Most text embeddings:     Cosine (or normalized dot product)
  Images / raw features:    L2 (Euclidean)
════════════════════════════════════════════════════════
```

**Exercise 1:** Benchmark FAISS index types. Create 100k random 384-dim vectors. Measure query time (average of 1000 queries) for: FlatIP, IVFFlat (nprobe=8), IVFFlat (nprobe=32), HNSWFlat. Plot the recall-speed tradeoff.

**Exercise 2:** Implement a "freshness-aware" search: documents have a `timestamp` metadata field. Blend semantic similarity score with a freshness score (newer docs score higher). Formula: `final_score = 0.7 * semantic + 0.3 * freshness_factor`.

**Exercise 3:** Build a "duplicate detector": given a collection of 1000 documents, find all pairs with cosine similarity > 0.95. Compare brute-force vs HNSW index approach in terms of speed.

**Exercise 4:** ChromaDB vs FAISS benchmark: index 50k document embeddings in both. Compare: build time, query time (p50, p95), memory usage. Which would you choose for production?

**Exercise 5:** Implement **chunking strategies** for a long document (>10 pages): (a) fixed-size chunks of 200 chars with 50-char overlap, (b) sentence-level chunks, (c) paragraph-level chunks. Index all three in separate collections. Compare retrieval quality on 10 test questions.

---

← [Chapter 33: Embeddings and Semantic Search](./33-embeddings-and-semantic-search.md) | [Chapter 35: RAG](./35-rag-retrieval-augmented-generation.md) →
