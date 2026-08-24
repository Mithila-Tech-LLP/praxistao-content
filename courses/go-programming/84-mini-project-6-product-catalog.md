# Chapter 84: Mini Project 6 — Product Catalog with Full-Text Search

This project synthesizes Vol 6 (Databases & Persistence). You'll build a product catalog API backed by PostgreSQL for storage and OpenSearch for search. Products have JSONB attributes for flexible per-category fields. The API supports full-text search, faceted filtering, and pagination.

**What you'll build**: A REST API with:
- PostgreSQL (primary store, JSONB attributes)
- Full-text search via PostgreSQL `tsvector` (simple case) or OpenSearch (advanced case)
- Paginated listing with facets (categories, price ranges)
- Repository pattern with caching layer (Redis optional)

---

## Project Structure

```
catalog/
├── cmd/
│   └── server/
│       └── main.go
├── domain/
│   ├── product.go
│   └── repository.go
├── postgres/
│   ├── db.go
│   └── product_repository.go
├── handler/
│   └── product_handler.go
├── migrations/
│   └── 001_products.sql
└── go.mod
```

---

## Domain Layer

```go
// domain/product.go
package domain

import (
    "errors"
    "strings"
    "time"
)

var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
)

type Product struct {
    ID          int64              `json:"id"`
    Name        string             `json:"name"`
    Slug        string             `json:"slug"`
    Description string             `json:"description,omitempty"`
    Category    string             `json:"category"`
    Brand       string             `json:"brand,omitempty"`
    Price       float64            `json:"price"`
    InStock     bool               `json:"in_stock"`
    Tags        []string           `json:"tags,omitempty"`
    Attributes  map[string]any     `json:"attributes,omitempty"`
    CreatedAt   time.Time          `json:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at"`
}

type CreateProductRequest struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    Category    string         `json:"category"`
    Brand       string         `json:"brand"`
    Price       float64        `json:"price"`
    InStock     bool           `json:"in_stock"`
    Tags        []string       `json:"tags"`
    Attributes  map[string]any `json:"attributes"`
}

func (r *CreateProductRequest) Validate() error {
    var errs []string
    if strings.TrimSpace(r.Name) == ""     { errs = append(errs, "name is required") }
    if strings.TrimSpace(r.Category) == "" { errs = append(errs, "category is required") }
    if r.Price < 0                          { errs = append(errs, "price must be >= 0") }
    if len(errs) > 0 {
        return fmt.Errorf("%w: %s", ErrInvalidInput, strings.Join(errs, "; "))
    }
    return nil
}

type SearchFilter struct {
    Query     string
    Category  string
    Brand     string
    MinPrice  float64
    MaxPrice  float64
    InStock   *bool
    Tags      []string
    SortBy    string  // "price_asc", "price_desc", "newest", "relevance"
    Page      int
    PageSize  int
}

type SearchResult struct {
    Products   []*Product       `json:"products"`
    Total      int64            `json:"total"`
    Page       int              `json:"page"`
    PageSize   int              `json:"page_size"`
    HasMore    bool             `json:"has_more"`
    Facets     *SearchFacets    `json:"facets,omitempty"`
}

type SearchFacets struct {
    Categories  []FacetCount `json:"categories"`
    Brands      []FacetCount `json:"brands"`
    PriceRanges []FacetCount `json:"price_ranges"`
}

type FacetCount struct {
    Value string `json:"value"`
    Count int    `json:"count"`
}

// domain/repository.go
type ProductRepository interface {
    Create(ctx context.Context, p *Product) error
    GetByID(ctx context.Context, id int64) (*Product, error)
    GetBySlug(ctx context.Context, slug string) (*Product, error)
    Update(ctx context.Context, p *Product) error
    Delete(ctx context.Context, id int64) error
    Search(ctx context.Context, f SearchFilter) (*SearchResult, error)
}
```

---

## Database Migration

```sql
-- migrations/001_products.sql

CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    description TEXT,
    category    TEXT NOT NULL,
    brand       TEXT,
    price       NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    in_stock    BOOLEAN NOT NULL DEFAULT TRUE,
    tags        TEXT[],
    attributes  JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english',
            coalesce(name, '')          || ' ' ||
            coalesce(description, '')   || ' ' ||
            coalesce(brand, '')         || ' ' ||
            coalesce(category, '')      || ' ' ||
            coalesce(array_to_string(tags, ' '), '')
        )
    ) STORED
);

CREATE INDEX idx_products_search    ON products USING GIN (search_vector);
CREATE INDEX idx_products_category  ON products (category);
CREATE INDEX idx_products_brand     ON products (brand) WHERE brand IS NOT NULL;
CREATE INDEX idx_products_price     ON products (price);
CREATE INDEX idx_products_in_stock  ON products (in_stock) WHERE in_stock = TRUE;
CREATE INDEX idx_products_tags      ON products USING GIN (tags);
CREATE INDEX idx_products_attrs     ON products USING GIN (attributes);
```

---

## PostgreSQL Repository

```go
// postgres/product_repository.go
package postgres

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "strings"

    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"

    "catalog/domain"
)

type ProductRepository struct {
    db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

func slugify(name string) string {
    s := strings.ToLower(strings.TrimSpace(name))
    var out strings.Builder
    prevDash := false
    for _, r := range s {
        switch {
        case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
            out.WriteRune(r)
            prevDash = false
        default:
            if !prevDash {
                out.WriteByte('-')
                prevDash = true
            }
        }
    }
    return strings.Trim(out.String(), "-")
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
    p.Slug = slugify(p.Name)
    
    attrs, err := json.Marshal(p.Attributes)
    if err != nil { return fmt.Errorf("marshal attributes: %w", err) }
    
    query := `
        INSERT INTO products (name, slug, description, category, brand, price, in_stock, tags, attributes)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at`
    
    err = r.db.QueryRowContext(ctx, query,
        p.Name, p.Slug, p.Description, p.Category, p.Brand,
        p.Price, p.InStock, pq.Array(p.Tags), attrs,
    ).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
    
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            p.Slug = fmt.Sprintf("%s-%d", p.Slug, p.ID)
            return r.Create(ctx, p)
        }
        return fmt.Errorf("create product: %w", err)
    }
    return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
    return r.scanOne(ctx, "SELECT * FROM products WHERE id = $1", id)
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
    return r.scanOne(ctx, "SELECT * FROM products WHERE slug = $1", slug)
}

func (r *ProductRepository) scanOne(ctx context.Context, query string, args ...any) (*domain.Product, error) {
    var row struct {
        ID          int64           `db:"id"`
        Name        string          `db:"name"`
        Slug        string          `db:"slug"`
        Description sql.NullString  `db:"description"`
        Category    string          `db:"category"`
        Brand       sql.NullString  `db:"brand"`
        Price       float64         `db:"price"`
        InStock     bool            `db:"in_stock"`
        Tags        pq.StringArray  `db:"tags"`
        Attributes  []byte          `db:"attributes"`
        CreatedAt   time.Time       `db:"created_at"`
        UpdatedAt   time.Time       `db:"updated_at"`
    }
    
    if err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&row); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
        return nil, err
    }
    
    var attrs map[string]any
    json.Unmarshal(row.Attributes, &attrs)
    
    return &domain.Product{
        ID:          row.ID,
        Name:        row.Name,
        Slug:        row.Slug,
        Description: row.Description.String,
        Category:    row.Category,
        Brand:       row.Brand.String,
        Price:       row.Price,
        InStock:     row.InStock,
        Tags:        []string(row.Tags),
        Attributes:  attrs,
        CreatedAt:   row.CreatedAt,
        UpdatedAt:   row.UpdatedAt,
    }, nil
}

func (r *ProductRepository) Search(ctx context.Context, f domain.SearchFilter) (*domain.SearchResult, error) {
    if f.Page < 1   { f.Page = 1 }
    if f.PageSize < 1 || f.PageSize > 100 { f.PageSize = 20 }
    
    where := []string{"TRUE"}
    args := []any{}
    argIdx := 1
    addArg := func(v any) string {
        args = append(args, v)
        s := fmt.Sprintf("$%d", argIdx)
        argIdx++
        return s
    }
    
    if f.Query != "" {
        where = append(where,
            fmt.Sprintf("search_vector @@ websearch_to_tsquery('english', %s)", addArg(f.Query)))
    }
    if f.Category != "" {
        where = append(where, fmt.Sprintf("category = %s", addArg(f.Category)))
    }
    if f.Brand != "" {
        where = append(where, fmt.Sprintf("brand = %s", addArg(f.Brand)))
    }
    if f.MinPrice > 0 {
        where = append(where, fmt.Sprintf("price >= %s", addArg(f.MinPrice)))
    }
    if f.MaxPrice > 0 {
        where = append(where, fmt.Sprintf("price <= %s", addArg(f.MaxPrice)))
    }
    if f.InStock != nil {
        where = append(where, fmt.Sprintf("in_stock = %s", addArg(*f.InStock)))
    }
    if len(f.Tags) > 0 {
        where = append(where, fmt.Sprintf("tags @> %s", addArg(pq.Array(f.Tags))))
    }
    
    orderBy := "created_at DESC"
    switch f.SortBy {
    case "price_asc":   orderBy = "price ASC"
    case "price_desc":  orderBy = "price DESC"
    case "newest":      orderBy = "created_at DESC"
    case "relevance":
        if f.Query != "" {
            orderBy = fmt.Sprintf("ts_rank(search_vector, websearch_to_tsquery('english', '%s')) DESC", f.Query)
        }
    }
    
    whereClause := strings.Join(where, " AND ")
    
    // Count query
    var total int64
    countSQL := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
    if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
        return nil, fmt.Errorf("count: %w", err)
    }
    
    // Facets: categories and brands
    facets, err := r.queryFacets(ctx, whereClause, args)
    if err != nil { return nil, err }
    
    // Main query
    offset := (f.Page - 1) * f.PageSize
    mainSQL := fmt.Sprintf(`
        SELECT id, name, slug, description, category, brand, price, in_stock, tags, attributes, created_at, updated_at
        FROM products
        WHERE %s
        ORDER BY %s
        LIMIT %d OFFSET %d`,
        whereClause, orderBy, f.PageSize, offset)
    
    rows, err := r.db.QueryxContext(ctx, mainSQL, args...)
    if err != nil { return nil, fmt.Errorf("search query: %w", err) }
    defer rows.Close()
    
    var products []*domain.Product
    for rows.Next() {
        var row struct {
            ID          int64          `db:"id"`
            Name        string         `db:"name"`
            Slug        string         `db:"slug"`
            Description sql.NullString `db:"description"`
            Category    string         `db:"category"`
            Brand       sql.NullString `db:"brand"`
            Price       float64        `db:"price"`
            InStock     bool           `db:"in_stock"`
            Tags        pq.StringArray `db:"tags"`
            Attributes  []byte         `db:"attributes"`
            CreatedAt   time.Time      `db:"created_at"`
            UpdatedAt   time.Time      `db:"updated_at"`
        }
        if err := rows.StructScan(&row); err != nil { return nil, err }
        
        var attrs map[string]any
        json.Unmarshal(row.Attributes, &attrs)
        
        products = append(products, &domain.Product{
            ID: row.ID, Name: row.Name, Slug: row.Slug,
            Description: row.Description.String,
            Category: row.Category, Brand: row.Brand.String,
            Price: row.Price, InStock: row.InStock,
            Tags: []string(row.Tags), Attributes: attrs,
            CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
        })
    }
    
    return &domain.SearchResult{
        Products: products,
        Total:    total,
        Page:     f.Page,
        PageSize: f.PageSize,
        HasMore:  int64(offset+f.PageSize) < total,
        Facets:   facets,
    }, nil
}

func (r *ProductRepository) queryFacets(ctx context.Context, whereClause string, args []any) (*domain.SearchFacets, error) {
    facetSQL := fmt.Sprintf(`
        SELECT 'category' AS type, category AS value, COUNT(*) AS cnt
        FROM products WHERE %s GROUP BY category
        UNION ALL
        SELECT 'brand', brand, COUNT(*)
        FROM products WHERE %s AND brand IS NOT NULL GROUP BY brand
        ORDER BY type, cnt DESC`, whereClause, whereClause)
    
    rows, err := r.db.QueryxContext(ctx, facetSQL, append(args, args...)...)
    if err != nil { return nil, err }
    defer rows.Close()
    
    facets := &domain.SearchFacets{}
    for rows.Next() {
        var row struct {
            Type  string `db:"type"`
            Value string `db:"value"`
            Count int    `db:"cnt"`
        }
        rows.StructScan(&row)
        fc := domain.FacetCount{Value: row.Value, Count: row.Count}
        switch row.Type {
        case "category": facets.Categories = append(facets.Categories, fc)
        case "brand":    facets.Brands = append(facets.Brands, fc)
        }
    }
    return facets, rows.Err()
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
    attrs, _ := json.Marshal(p.Attributes)
    _, err := r.db.ExecContext(ctx, `
        UPDATE products SET name=$1, description=$2, category=$3, brand=$4,
            price=$5, in_stock=$6, tags=$7, attributes=$8, updated_at=NOW()
        WHERE id=$9`,
        p.Name, p.Description, p.Category, p.Brand,
        p.Price, p.InStock, pq.Array(p.Tags), attrs, p.ID)
    return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
    return err
}
```

---

## HTTP Handler

```go
// handler/product_handler.go
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "catalog/domain"
)

type ProductHandler struct {
    repo domain.ProductRepository
}

func NewProductHandler(repo domain.ProductRepository) *ProductHandler {
    return &ProductHandler{repo: repo}
}

func (h *ProductHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/",           h.List)
    r.Post("/",          h.Create)
    r.Get("/{id}",       h.GetByID)
    r.Get("/slug/{slug}", h.GetBySlug)
    r.Put("/{id}",       h.Update)
    r.Delete("/{id}",    h.Delete)
    return r
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    
    page, _ := strconv.Atoi(q.Get("page"))
    if page < 1 { page = 1 }
    pageSize, _ := strconv.Atoi(q.Get("page_size"))
    if pageSize < 1 { pageSize = 20 }
    
    filter := domain.SearchFilter{
        Query:    q.Get("q"),
        Category: q.Get("category"),
        Brand:    q.Get("brand"),
        SortBy:   q.Get("sort"),
        Page:     page,
        PageSize: pageSize,
    }
    
    if minP, err := strconv.ParseFloat(q.Get("min_price"), 64); err == nil { filter.MinPrice = minP }
    if maxP, err := strconv.ParseFloat(q.Get("max_price"), 64); err == nil { filter.MaxPrice = maxP }
    if inStock := q.Get("in_stock"); inStock != "" {
        b, _ := strconv.ParseBool(inStock)
        filter.InStock = &b
    }
    if tags := q["tag"]; len(tags) > 0 { filter.Tags = tags }
    
    result, err := h.repo.Search(r.Context(), filter)
    if err != nil {
        http.Error(w, "search failed", http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req domain.CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    if err := req.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }
    
    p := &domain.Product{
        Name:        req.Name,
        Description: req.Description,
        Category:    req.Category,
        Brand:       req.Brand,
        Price:       req.Price,
        InStock:     req.InStock,
        Tags:        req.Tags,
        Attributes:  req.Attributes,
    }
    
    if err := h.repo.Create(r.Context(), p); err != nil {
        http.Error(w, "create failed", http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    p, err := h.repo.GetByID(r.Context(), id)
    if errors.Is(err, domain.ErrNotFound) {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, "fetch failed", http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
    p, err := h.repo.GetBySlug(r.Context(), chi.URLParam(r, "slug"))
    if errors.Is(err, domain.ErrNotFound) {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, "fetch failed", http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    
    p, err := h.repo.GetByID(r.Context(), id)
    if errors.Is(err, domain.ErrNotFound) {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    
    var req domain.CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    
    p.Name = req.Name
    p.Description = req.Description
    p.Category = req.Category
    p.Brand = req.Brand
    p.Price = req.Price
    p.InStock = req.InStock
    p.Tags = req.Tags
    p.Attributes = req.Attributes
    
    if err := h.repo.Update(r.Context(), p); err != nil {
        http.Error(w, "update failed", http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err := h.repo.Delete(r.Context(), id); err != nil {
        http.Error(w, "delete failed", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(v)
}
```

---

## Main

```go
// cmd/server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "catalog/handler"
    "catalog/postgres"
)

func main() {
    db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil { log.Fatal(err) }
    defer db.Close()
    
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    repo := postgres.NewProductRepository(db)
    h := handler.NewProductHandler(repo)
    
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Mount("/products", h.Routes())
    
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    
    go func() {
        fmt.Println("listening on :8080")
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    <-ctx.Done()
    
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    srv.Shutdown(shutdownCtx)
}
```

---

## Testing the API

```bash
# Create a product
curl -X POST http://localhost:8080/products \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Sony WH-1000XM5",
    "category": "headphones",
    "brand": "Sony",
    "price": 349.99,
    "in_stock": true,
    "tags": ["wireless", "noise-cancelling"],
    "attributes": {
      "battery_hours": 30,
      "connectivity": ["bluetooth", "aux"]
    }
  }'

# Full-text search
curl "http://localhost:8080/products?q=noise+cancelling&category=headphones"

# Filter with facets
curl "http://localhost:8080/products?min_price=100&max_price=500&sort=price_asc&tag=wireless"

# Get by slug
curl http://localhost:8080/products/slug/sony-wh-1000xm5
```

---

## Extension Challenges

1. **Add Redis caching**: wrap `ProductRepository` with a `CachingProductRepository`. Cache individual products by ID and slug with a 5-minute TTL. Invalidate on update/delete.

2. **Add OpenSearch**: sync products to OpenSearch using the pattern from Ch 82. Use OpenSearch for the `/products?q=...` endpoint and PostgreSQL for everything else.

3. **Bulk import**: add `POST /products/import` that accepts a JSON array of products and inserts them in a single database transaction. Return a summary of created/failed counts.

4. **Category hierarchy**: add a `parent_category` field. Build a query that returns all products in a category and all its descendants using a recursive CTE (`WITH RECURSIVE`).
