# Chapter 22: Building with MySQL in Go

We've covered MySQL's internals. Now let's build with it. This chapter is code-first: complete working patterns for everything you'll need in production MySQL applications.

## Table of Contents

1. Production Connection Setup
2. CRUD with Transactions
3. Bulk Inserts
4. Schema Migrations with golang-migrate
5. Mini Project: Product Catalog API
6. Exercises

---

## 1. Production Connection Setup

```go
package db

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
}

func New(cfg Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=UTC&multiStatements=true",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
    )

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("open mysql: %w", err)
    }

    // Connection pool tuning
    db.SetMaxOpenConns(25)          // max connections in pool
    db.SetMaxIdleConns(5)           // idle connections kept alive
    db.SetConnMaxLifetime(5 * time.Minute) // recycle connections
    db.SetConnMaxIdleTime(1 * time.Minute) // close idle connections

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("ping mysql: %w", err)
    }

    return db, nil
}
```

---

## 2. CRUD with Transactions

```go
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/go-sql-driver/mysql"
)

type Product struct {
    ID          int64
    Name        string
    Price       float64
    Stock       int
    Description string
}

// CreateProduct inserts a product and returns its ID.
func CreateProduct(ctx context.Context, db *sql.DB, p Product) (int64, error) {
    result, err := db.ExecContext(ctx,
        "INSERT INTO products (name, price, stock, description) VALUES (?, ?, ?, ?)",
        p.Name, p.Price, p.Stock, p.Description,
    )
    if err != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            return 0, fmt.Errorf("product %q already exists", p.Name)
        }
        return 0, err
    }
    return result.LastInsertId()
}

// GetProduct retrieves a product by ID.
func GetProduct(ctx context.Context, db *sql.DB, id int64) (*Product, error) {
    var p Product
    err := db.QueryRowContext(ctx,
        "SELECT id, name, price, stock, description FROM products WHERE id = ?", id,
    ).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Description)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    return &p, err
}

// UpdateStock atomically adjusts stock by delta (negative to decrease).
func UpdateStock(ctx context.Context, db *sql.DB, productID int64, delta int) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    var current int
    err = tx.QueryRowContext(ctx,
        "SELECT stock FROM products WHERE id = ? FOR UPDATE",
        productID,
    ).Scan(&current)
    if errors.Is(err, sql.ErrNoRows) {
        return fmt.Errorf("product %d not found", productID)
    }
    if err != nil {
        return err
    }

    newStock := current + delta
    if newStock < 0 {
        return fmt.Errorf("insufficient stock: have %d, need %d", current, -delta)
    }

    _, err = tx.ExecContext(ctx,
        "UPDATE products SET stock = ? WHERE id = ?",
        newStock, productID,
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

// ListProducts returns paginated products.
func ListProducts(ctx context.Context, db *sql.DB, limit, offset int) ([]Product, error) {
    rows, err := db.QueryContext(ctx,
        "SELECT id, name, price, stock, description FROM products ORDER BY id LIMIT ? OFFSET ?",
        limit, offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Description); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, rows.Err()
}
```

---

## 3. Bulk Inserts

### Method 1: Multi-row INSERT (Fastest for < 10,000 rows)

```go
func BulkInsertProducts(ctx context.Context, db *sql.DB, products []Product) error {
    if len(products) == 0 {
        return nil
    }

    // Build: INSERT INTO products (name, price) VALUES (?,?),(?,?),...
    query := "INSERT INTO products (name, price, stock) VALUES "
    args := make([]interface{}, 0, len(products)*3)
    
    for i, p := range products {
        if i > 0 {
            query += ","
        }
        query += "(?,?,?)"
        args = append(args, p.Name, p.Price, p.Stock)
    }

    _, err := db.ExecContext(ctx, query, args...)
    return err
}
```

### Method 2: Prepared statement + transaction (Best for large batches)

```go
func BulkInsertLarge(ctx context.Context, db *sql.DB, products []Product) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        "INSERT INTO products (name, price, stock) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, p := range products {
        if _, err := stmt.ExecContext(ctx, p.Name, p.Price, p.Stock); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### Method 3: LOAD DATA INFILE (Fastest for millions of rows)

```go
func BulkLoadFromCSV(ctx context.Context, db *sql.DB, csvPath string) error {
    // Register the file for reading (security feature of MySQL driver)
    mysql.RegisterLocalFile(csvPath)
    
    _, err := db.ExecContext(ctx, fmt.Sprintf(`
        LOAD DATA LOCAL INFILE '%s'
        INTO TABLE products
        FIELDS TERMINATED BY ','
        LINES TERMINATED BY '\n'
        IGNORE 1 ROWS
        (name, price, stock)
    `, csvPath))
    return err
}
```

---

## 4. Schema Migrations with golang-migrate

```bash
go get -tags mysql github.com/golang-migrate/migrate/v4
```

Create migration files:

```bash
mkdir -p migrations
cat > migrations/000001_create_products.up.sql << 'EOF'
CREATE TABLE IF NOT EXISTS products (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    price       DECIMAL(10,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    description TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_products_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
EOF

cat > migrations/000001_create_products.down.sql << 'EOF'
DROP TABLE IF EXISTS products;
EOF
```

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) error {
    // migrate expects: mysql://user:pass@host:port/dbname
    m, err := migrate.New("file://migrations", "mysql://"+dsn)
    if err != nil {
        return fmt.Errorf("create migrator: %w", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migrate up: %w", err)
    }
    return nil
}
```

---

## 5. Mini Project: Product Catalog API

A complete REST API for a product catalog with MySQL:

```go
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/go-sql-driver/mysql"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Product struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

func main() {
    var err error
    db, err = sql.Open("mysql",
        "dev:secret@tcp(localhost:3306)/catalog?parseTime=true&charset=utf8mb4")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)

    if err := db.Ping(); err != nil {
        log.Fatal("MySQL not available:", err)
    }

    // Ensure table exists
    db.Exec(`CREATE TABLE IF NOT EXISTS products (
        id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        name        VARCHAR(255) UNIQUE NOT NULL,
        price       DECIMAL(10,2) NOT NULL,
        stock       INT NOT NULL DEFAULT 0,
        description TEXT,
        created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /products", handleList)
    mux.HandleFunc("POST /products", handleCreate)
    mux.HandleFunc("GET /products/", handleGet)
    mux.HandleFunc("PATCH /products/", handleUpdateStock)

    log.Println("Product catalog API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleList(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    limit := 20
    offset := (page - 1) * limit

    rows, err := db.QueryContext(r.Context(),
        "SELECT id, name, price, stock, IFNULL(description,''), created_at FROM products ORDER BY id LIMIT ? OFFSET ?",
        limit, offset)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Description, &p.CreatedAt); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        products = append(products, p)
    }
    if products == nil {
        products = []Product{}
    }
    json.NewEncoder(w).Encode(products)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name        string  `json:"name"`
        Price       float64 `json:"price"`
        Stock       int     `json:"stock"`
        Description string  `json:"description"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }
    if req.Name == "" || req.Price <= 0 {
        http.Error(w, "name and price required", 400)
        return
    }

    result, err := db.ExecContext(r.Context(),
        "INSERT INTO products (name, price, stock, description) VALUES (?, ?, ?, ?)",
        req.Name, req.Price, req.Stock, req.Description)
    if err != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            http.Error(w, "product name already exists", 409)
            return
        }
        http.Error(w, err.Error(), 500)
        return
    }

    id, _ := result.LastInsertId()
    var p Product
    db.QueryRowContext(r.Context(),
        "SELECT id, name, price, stock, IFNULL(description,''), created_at FROM products WHERE id = ?", id,
    ).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Description, &p.CreatedAt)

    w.WriteHeader(201)
    json.NewEncoder(w).Encode(p)
}

func parseID(r *http.Request, prefix string) (int64, error) {
    idStr := strings.TrimPrefix(r.URL.Path, prefix)
    idStr = strings.Split(idStr, "/")[0]
    return strconv.ParseInt(idStr, 10, 64)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r, "/products/")
    if err != nil {
        http.Error(w, "invalid id", 400)
        return
    }

    var p Product
    err = db.QueryRowContext(r.Context(),
        "SELECT id, name, price, stock, IFNULL(description,''), created_at FROM products WHERE id = ?", id,
    ).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Description, &p.CreatedAt)

    if errors.Is(err, sql.ErrNoRows) {
        http.Error(w, "not found", 404)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(p)
}

func handleUpdateStock(w http.ResponseWriter, r *http.Request) {
    // PATCH /products/{id}/stock
    parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/products/"), "/")
    if len(parts) < 2 || parts[1] != "stock" {
        http.Error(w, "not found", 404)
        return
    }
    id, err := strconv.ParseInt(parts[0], 10, 64)
    if err != nil {
        http.Error(w, "invalid id", 400)
        return
    }

    var req struct{ Delta int `json:"delta"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }

    if err := UpdateStock(r.Context(), db, id, req.Delta); err != nil {
        if strings.Contains(err.Error(), "insufficient") {
            http.Error(w, err.Error(), 409)
        } else if strings.Contains(err.Error(), "not found") {
            http.Error(w, err.Error(), 404)
        } else {
            http.Error(w, err.Error(), 500)
        }
        return
    }
    w.WriteHeader(204)
}
```

Test:
```bash
# Add a product
curl -X POST localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99,"stock":100}'

# List products
curl localhost:8080/products

# Update stock
curl -X PATCH localhost:8080/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{"delta":-5}'
```

---

## Summary

- Always set `parseTime=true`, `charset=utf8mb4`, and `loc=UTC` in the MySQL DSN.
- Use `db.ExecContext` and `db.QueryContext` (with context) for production — they respect timeouts.
- Bulk inserts: multi-row INSERT for small batches, prepared statement + transaction for large batches.
- `result.LastInsertId()` returns the AUTO_INCREMENT value after INSERT.
- Handle MySQL error code 1062 (duplicate entry) separately from generic errors.

### Exercises

**Easy:** Add `DELETE /products/{id}` to the catalog API. Return 404 if not found, 204 on success.

**Medium:** Add product search: `GET /products?q=widget` searches by name using `LIKE %query%`. Add a full-text index on `name` and compare `LIKE` vs `MATCH AGAINST` performance.

**Hard:** Implement optimistic locking: add a `version` column to products. The `PATCH /products/{id}/stock` endpoint should require a `version` field and return 409 Conflict if the version doesn't match (someone else updated first).
