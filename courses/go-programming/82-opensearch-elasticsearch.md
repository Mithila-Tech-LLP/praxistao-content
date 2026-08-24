# Chapter 82: OpenSearch and Elasticsearch Integration

OpenSearch (the AWS fork of Elasticsearch) is a distributed search and analytics engine. It excels at full-text search, log analytics, and faceted navigation. While PostgreSQL's built-in full-text search works well up to millions of rows, OpenSearch handles billions with sub-second latency and richer query features.

## Table of Contents

1. [When to Use OpenSearch vs PostgreSQL FTS](#1-when-to-use-opensearch-vs-postgresql-fts)
2. [Indexes and Mappings](#2-indexes-and-mappings)
3. [Indexing Documents in Go](#3-indexing-documents-in-go)
4. [Search Queries](#4-search-queries)
5. [Bulk Operations and Index Sync](#5-bulk-operations-and-index-sync)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. When to Use OpenSearch vs PostgreSQL FTS

| | PostgreSQL FTS | OpenSearch |
|---|---|---|
| Setup | None — built in | Separate service |
| Scale | Up to ~50M rows easily | Billions of docs |
| Query language | tsquery | Full JSON DSL |
| Facets/aggregations | GROUP BY | Native, fast |
| Fuzzy/autocomplete | Limited | Excellent |
| Consistency | ACID | Eventually consistent |
| Source of truth | Yes | No — sync from primary DB |

**Rule**: Use PostgreSQL FTS for internal search on moderate data. Add OpenSearch when you need fuzzy matching, autocomplete, faceted navigation, or log analytics at scale.

---

## 2. Indexes and Mappings

An OpenSearch **index** is like a table. **Mappings** define the field types.

```json
PUT /products
{
  "settings": {
    "number_of_shards": 2,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "product_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding", "en_stemmer"]
        }
      },
      "filter": {
        "en_stemmer": {
          "type": "stemmer",
          "language": "english"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id":          { "type": "keyword" },
      "name":        { "type": "text", "analyzer": "product_analyzer",
                       "fields": { "keyword": { "type": "keyword" } } },
      "description": { "type": "text", "analyzer": "product_analyzer" },
      "category":    { "type": "keyword" },
      "brand":       { "type": "keyword" },
      "price":       { "type": "float" },
      "tags":        { "type": "keyword" },
      "in_stock":    { "type": "boolean" },
      "created_at":  { "type": "date" }
    }
  }
}
```

### Field type cheat sheet

| Type | Use for |
|------|---------|
| `text` | Full-text search (tokenized and analyzed) |
| `keyword` | Exact match, filtering, sorting, aggregations |
| `float` / `integer` | Numeric comparisons |
| `date` | Date range queries |
| `boolean` | Filter-only boolean |

Use `"fields"` to index a field as both `text` (for search) and `keyword` (for aggregation/sort).

---

## 3. Indexing Documents in Go

```go
import opensearch "github.com/opensearch-project/opensearch-go/v2"

func newOpenSearchClient(addrs []string) (*opensearch.Client, error) {
    cfg := opensearch.Config{
        Addresses: addrs, // []string{"http://localhost:9200"}
    }
    return opensearch.NewClient(cfg)
}

type ProductDoc struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Category    string    `json:"category"`
    Brand       string    `json:"brand,omitempty"`
    Price       float64   `json:"price"`
    Tags        []string  `json:"tags,omitempty"`
    InStock     bool      `json:"in_stock"`
    CreatedAt   time.Time `json:"created_at"`
}

type ProductIndex struct {
    client *opensearch.Client
    index  string
}

// Index a single document
func (idx *ProductIndex) Index(ctx context.Context, doc ProductDoc) error {
    data, err := json.Marshal(doc)
    if err != nil { return err }
    
    resp, err := idx.client.Index(
        idx.index,
        bytes.NewReader(data),
        idx.client.Index.WithDocumentID(doc.ID),
        idx.client.Index.WithContext(ctx),
        idx.client.Index.WithRefresh("false"), // don't wait for segment refresh
    )
    if err != nil { return fmt.Errorf("index request: %w", err) }
    defer resp.Body.Close()
    
    if resp.IsError() {
        return fmt.Errorf("index error: %s", resp.String())
    }
    return nil
}

// Delete a document
func (idx *ProductIndex) Delete(ctx context.Context, id string) error {
    resp, err := idx.client.Delete(idx.index, id,
        idx.client.Delete.WithContext(ctx))
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.IsError() && resp.StatusCode != 404 {
        return fmt.Errorf("delete error: %s", resp.String())
    }
    return nil
}
```

---

## 4. Search Queries

```go
type SearchRequest struct {
    Query     string
    Category  string
    MinPrice  float64
    MaxPrice  float64
    InStock   *bool
    Tags      []string
    Page      int
    PageSize  int
}

type SearchResult struct {
    Total    int64
    Products []ProductDoc
}

func (idx *ProductIndex) Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
    // Build query
    must := []map[string]any{}
    filter := []map[string]any{}
    
    // Full-text match
    if req.Query != "" {
        must = append(must, map[string]any{
            "multi_match": map[string]any{
                "query":     req.Query,
                "fields":    []string{"name^3", "description", "brand^2"},
                "fuzziness": "AUTO",
                "type":      "best_fields",
            },
        })
    }
    
    // Category filter
    if req.Category != "" {
        filter = append(filter, map[string]any{
            "term": map[string]any{"category": req.Category},
        })
    }
    
    // Price range
    if req.MinPrice > 0 || req.MaxPrice > 0 {
        priceRange := map[string]any{}
        if req.MinPrice > 0 { priceRange["gte"] = req.MinPrice }
        if req.MaxPrice > 0 { priceRange["lte"] = req.MaxPrice }
        filter = append(filter, map[string]any{
            "range": map[string]any{"price": priceRange},
        })
    }
    
    // In-stock filter
    if req.InStock != nil {
        filter = append(filter, map[string]any{
            "term": map[string]any{"in_stock": *req.InStock},
        })
    }
    
    // Tags filter
    if len(req.Tags) > 0 {
        filter = append(filter, map[string]any{
            "terms": map[string]any{"tags": req.Tags},
        })
    }
    
    if len(must) == 0 {
        must = append(must, map[string]any{"match_all": map[string]any{}})
    }
    
    from := (req.Page - 1) * req.PageSize
    
    body := map[string]any{
        "from": from,
        "size": req.PageSize,
        "query": map[string]any{
            "bool": map[string]any{
                "must":   must,
                "filter": filter,
            },
        },
        "aggs": map[string]any{
            "categories": map[string]any{
                "terms": map[string]any{"field": "category", "size": 20},
            },
            "price_stats": map[string]any{
                "stats": map[string]any{"field": "price"},
            },
        },
        "highlight": map[string]any{
            "fields": map[string]any{
                "name":        map[string]any{},
                "description": map[string]any{"number_of_fragments": 2},
            },
        },
    }
    
    data, _ := json.Marshal(body)
    resp, err := idx.client.Search(
        idx.client.Search.WithIndex(idx.index),
        idx.client.Search.WithBody(bytes.NewReader(data)),
        idx.client.Search.WithContext(ctx),
    )
    if err != nil { return nil, fmt.Errorf("search request: %w", err) }
    defer resp.Body.Close()
    
    if resp.IsError() { return nil, fmt.Errorf("search error: %s", resp.String()) }
    
    var result struct {
        Hits struct {
            Total struct{ Value int64 } `json:"total"`
            Hits  []struct {
                Source    ProductDoc            `json:"_source"`
                Highlight map[string][]string   `json:"highlight"`
            } `json:"hits"`
        } `json:"hits"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }
    
    products := make([]ProductDoc, len(result.Hits.Hits))
    for i, hit := range result.Hits.Hits {
        products[i] = hit.Source
    }
    
    return &SearchResult{
        Total:    result.Hits.Total.Value,
        Products: products,
    }, nil
}
```

---

## 5. Bulk Operations and Index Sync

### Bulk indexing

```go
// Bulk index up to 1000 documents at once
func (idx *ProductIndex) BulkIndex(ctx context.Context, docs []ProductDoc) error {
    var buf bytes.Buffer
    
    for _, doc := range docs {
        // Action metadata line
        meta := map[string]any{
            "index": map[string]any{
                "_index": idx.index,
                "_id":    doc.ID,
            },
        }
        if err := json.NewEncoder(&buf).Encode(meta); err != nil { return err }
        
        // Document line
        if err := json.NewEncoder(&buf).Encode(doc); err != nil { return err }
    }
    
    resp, err := idx.client.Bulk(
        bytes.NewReader(buf.Bytes()),
        idx.client.Bulk.WithContext(ctx),
        idx.client.Bulk.WithRefresh("false"),
    )
    if err != nil { return fmt.Errorf("bulk request: %w", err) }
    defer resp.Body.Close()
    
    if resp.IsError() { return fmt.Errorf("bulk error: %s", resp.String()) }
    return nil
}
```

### Sync pattern: PostgreSQL → OpenSearch

The standard pattern for keeping OpenSearch in sync with PostgreSQL:

```go
// Option 1: Sync on write (simplest, works well for low write rates)
func (s *ProductService) Create(ctx context.Context, product *Product) error {
    // 1. Write to PostgreSQL (source of truth)
    if err := s.repo.Create(ctx, product); err != nil { return err }
    
    // 2. Index in OpenSearch (best-effort, async)
    go func() {
        doc := toProductDoc(product)
        if err := s.searchIdx.Index(context.Background(), doc); err != nil {
            s.logger.Error("failed to index product", "id", product.ID, "err", err)
            // Could push to a retry queue instead
        }
    }()
    return nil
}

// Option 2: CDC via logical replication (for high write rates)
// Use a tool like Debezium or write-a-tail with pg_logical
// Debezium reads the PostgreSQL WAL and publishes changes to Kafka
// A consumer reads from Kafka and bulk-indexes into OpenSearch
// This ensures eventually consistent sync even if the app crashes mid-write
```

---

## Summary

- OpenSearch = Elasticsearch-compatible search engine; great for full-text, facets, log analytics
- **`keyword`** vs **`text`**: keyword = exact match, text = analyzed/tokenized
- Always use `"fields"` to index a field as both `text` and `keyword` when you need to both search and aggregate on it
- Build queries with **bool query**: `must` (scored), `filter` (no scoring, cached)
- **Bulk API** for high-throughput indexing — never index one document at a time in a loop
- **Sync strategy**: write-through for low write rates; CDC (Debezium + Kafka) for high write rates

## Exercises

### Easy
1. Set up a local OpenSearch instance with Docker. Create an `articles` index with mappings for `title` (text), `body` (text), `author` (keyword), and `published_at` (date).
2. Write a Go function that indexes 100 articles and then searches for articles containing "distributed systems" with pagination.
3. Add a category facet to your search results using a `terms` aggregation on the `category` field. Display category counts alongside search results.

### Medium
4. Implement **autocomplete**: use the `completion` field type and the `suggest` API to build a product name autocomplete that works as the user types. Return up to 5 suggestions.
5. Build a **relevance tuning experiment**: index 50 products. Run the query "noise cancelling headphones" with different `fields` boost configurations (`name^1`, `name^3`, `name^5`). Compare the ranking of results.
6. Implement a **sync worker** that reads from a PostgreSQL `changes` table (populated by a trigger) and bulk-indexes the changed records into OpenSearch every 5 seconds. Handle failed syncs by tracking the last successful sync position.

### Hard
7. Implement **index aliasing with zero-downtime reindex**: when the mapping changes, create a new index `products_v2`, reindex all documents from `products_v1`, then atomically swap the `products` alias from v1 to v2.
8. Build a **search analytics dashboard**: log every search query, response time, and result count to an `analytics` index. Write aggregations to answer: "top 20 search queries", "searches with zero results", "average response time by hour". Display the results in a terminal table using `go-term-markdown`.
