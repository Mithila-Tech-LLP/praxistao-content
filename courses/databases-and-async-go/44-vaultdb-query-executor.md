# Chapter 44: VaultDB — Query Executor

The executor takes an AST from the parser and actually runs the query — reading rows from storage, evaluating conditions, inserting into pages, writing WAL records. This is where all the layers come together.

## Table of Contents

1. The Executor Interface
2. Evaluating Expressions
3. Executing SELECT
4. Executing INSERT
5. Executing UPDATE and DELETE
6. Executing CREATE TABLE
7. Putting It All Together
8. Exercises

---

## 1. The Executor Interface

```go
// query/executor.go
package query

import (
    "fmt"
    "github.com/yourname/vaultdb/storage"
    "github.com/yourname/vaultdb/wal"
)

// Result holds the output of a query
type Result struct {
    Columns []string
    Rows    []storage.Row
    Affected int64 // for INSERT/UPDATE/DELETE
}

type Executor struct {
    dm      *storage.DiskManager
    bp      *storage.BufferPool
    wal     *wal.WAL
    catalog *storage.Catalog
    txnID   uint64
}

func NewExecutor(dm *storage.DiskManager, bp *storage.BufferPool, w *wal.WAL, catalog *storage.Catalog) *Executor {
    return &Executor{dm: dm, bp: bp, wal: w, catalog: catalog}
}

// Execute runs a parsed statement and returns its result
func (e *Executor) Execute(stmt Statement) (*Result, error) {
    switch s := stmt.(type) {
    case *SelectStmt:
        return e.execSelect(s)
    case *InsertStmt:
        return e.execInsert(s)
    case *UpdateStmt:
        return e.execUpdate(s)
    case *DeleteStmt:
        return e.execDelete(s)
    case *CreateTableStmt:
        return e.execCreateTable(s)
    }
    return nil, fmt.Errorf("unknown statement type: %T", stmt)
}

// getTable finds a table definition in the catalog
func (e *Executor) getTable(name string) (*storage.TableDef, error) {
    for i, t := range e.catalog.Tables {
        if t.Name == name {
            return &e.catalog.Tables[i], nil
        }
    }
    return nil, fmt.Errorf("table %q not found", name)
}

func (e *Executor) colIndex(tbl *storage.TableDef, colName string) (int, error) {
    for i, c := range tbl.Columns {
        if c.Name == colName {
            return i, nil
        }
    }
    return -1, fmt.Errorf("column %q not found in table %q", colName, tbl.Name)
}
```

---

## 2. Evaluating Expressions

The executor evaluates an AST expression against a row:

```go
// evalExpr evaluates an expression against a row, returning a Value
func (e *Executor) evalExpr(expr Expr, row storage.Row, tbl *storage.TableDef) (storage.Value, error) {
    switch ex := expr.(type) {
    case *IntLiteral:
        return storage.IntVal(ex.Value), nil
    case *FloatLiteral:
        return storage.FloatVal(ex.Value), nil
    case *StrLiteral:
        return storage.StringVal(ex.Value), nil
    case *BoolLiteral:
        return storage.BoolVal(ex.Value), nil
    case *NullLiteral:
        return storage.NullVal(), nil

    case *ColumnRef:
        idx, err := e.colIndex(tbl, ex.Name)
        if err != nil {
            return storage.NullVal(), err
        }
        if idx >= len(row) {
            return storage.NullVal(), nil
        }
        return row[idx], nil

    case *BinaryOp:
        left, err := e.evalExpr(ex.Left, row, tbl)
        if err != nil {
            return storage.NullVal(), err
        }
        right, err := e.evalExpr(ex.Right, row, tbl)
        if err != nil {
            return storage.NullVal(), err
        }

        switch ex.Op {
        case "AND":
            return storage.BoolVal(left.AsBool() && right.AsBool()), nil
        case "OR":
            return storage.BoolVal(left.AsBool() || right.AsBool()), nil
        }

        cmp := left.Compare(right)
        switch ex.Op {
        case "=":
            return storage.BoolVal(cmp == 0), nil
        case "!=", "<>":
            return storage.BoolVal(cmp != 0), nil
        case "<":
            return storage.BoolVal(cmp < 0), nil
        case "<=":
            return storage.BoolVal(cmp <= 0), nil
        case ">":
            return storage.BoolVal(cmp > 0), nil
        case ">=":
            return storage.BoolVal(cmp >= 0), nil
        }
        return storage.NullVal(), fmt.Errorf("unknown operator: %s", ex.Op)
    }

    return storage.NullVal(), fmt.Errorf("unknown expr type: %T", expr)
}

// matches returns true if the row satisfies the WHERE condition
func (e *Executor) matches(where Expr, row storage.Row, tbl *storage.TableDef) (bool, error) {
    if where == nil {
        return true, nil
    }
    val, err := e.evalExpr(where, row, tbl)
    if err != nil {
        return false, err
    }
    return val.AsBool(), nil
}
```

---

## 3. Executing SELECT

```go
func (e *Executor) execSelect(stmt *SelectStmt) (*Result, error) {
    tbl, err := e.getTable(stmt.Table)
    if err != nil {
        return nil, err
    }

    heap := storage.NewHeap(e.dm, tbl.RootPageID)
    allRows, err := heap.ScanAll(len(tbl.Columns))
    if err != nil {
        return nil, err
    }

    // Filter rows
    var filtered []storage.Row
    for _, row := range allRows {
        ok, err := e.matches(stmt.Where, row, tbl)
        if err != nil {
            return nil, err
        }
        if ok {
            filtered = append(filtered, row)
        }
    }

    // Project columns
    selectAll := len(stmt.Columns) == 1 && stmt.Columns[0] == "*"
    var colNames []string
    var projected []storage.Row

    if selectAll {
        colNames = make([]string, len(tbl.Columns))
        for i, c := range tbl.Columns {
            colNames[i] = c.Name
        }
        projected = filtered
    } else {
        colNames = stmt.Columns
        for _, row := range filtered {
            var projRow storage.Row
            for _, colName := range stmt.Columns {
                idx, err := e.colIndex(tbl, colName)
                if err != nil {
                    return nil, err
                }
                if idx < len(row) {
                    projRow = append(projRow, row[idx])
                } else {
                    projRow = append(projRow, storage.NullVal())
                }
            }
            projected = append(projected, projRow)
        }
    }

    // LIMIT
    if stmt.Limit > 0 && len(projected) > stmt.Limit {
        projected = projected[:stmt.Limit]
    }

    return &Result{Columns: colNames, Rows: projected}, nil
}
```

---

## 4. Executing INSERT

```go
func (e *Executor) execInsert(stmt *InsertStmt) (*Result, error) {
    tbl, err := e.getTable(stmt.Table)
    if err != nil {
        return nil, err
    }

    // Build row with correct column order
    row := make(storage.Row, len(tbl.Columns))
    for i := range row {
        row[i] = storage.NullVal()
    }

    var columns []string
    if len(stmt.Columns) == 0 {
        // No column list: values correspond to columns in definition order
        for _, c := range tbl.Columns {
            columns = append(columns, c.Name)
        }
    } else {
        columns = stmt.Columns
    }

    if len(columns) != len(stmt.Values) {
        return nil, fmt.Errorf("column count (%d) does not match value count (%d)", len(columns), len(stmt.Values))
    }

    for i, colName := range columns {
        val, err := e.evalExpr(stmt.Values[i], nil, tbl)
        if err != nil {
            return nil, err
        }
        idx, err := e.colIndex(tbl, colName)
        if err != nil {
            return nil, err
        }
        row[idx] = val
    }

    heap := storage.NewHeap(e.dm, tbl.RootPageID)

    // Write WAL record before modifying the page
    e.txnID++
    data := storage.EncodeRow(row)
    _, walErr := e.wal.Append(wal.Record{
        TxnID:   e.txnID,
        Type:    wal.RecordInsert,
        NewData: data,
    })
    if walErr != nil {
        return nil, walErr
    }

    // Insert into heap
    _, err = heap.InsertRow(row)
    if err != nil {
        return nil, err
    }

    // Write COMMIT record and flush
    e.wal.Append(wal.Record{TxnID: e.txnID, Type: wal.RecordCommit})
    e.wal.Flush()

    return &Result{Affected: 1}, nil
}
```

---

## 5. Executing UPDATE and DELETE

```go
func (e *Executor) execUpdate(stmt *UpdateStmt) (*Result, error) {
    tbl, err := e.getTable(stmt.Table)
    if err != nil {
        return nil, err
    }

    heap := storage.NewHeap(e.dm, tbl.RootPageID)
    allRows, err := heap.ScanAll(len(tbl.Columns))
    if err != nil {
        return nil, err
    }

    var affected int64
    e.txnID++

    for i, row := range allRows {
        ok, err := e.matches(stmt.Where, row, tbl)
        if err != nil {
            return nil, err
        }
        if !ok {
            continue
        }

        // Apply updates to a copy
        updated := make(storage.Row, len(row))
        copy(updated, row)

        for _, assignment := range stmt.Sets {
            val, err := e.evalExpr(assignment.Value, row, tbl)
            if err != nil {
                return nil, err
            }
            idx, err := e.colIndex(tbl, assignment.Column)
            if err != nil {
                return nil, err
            }
            updated[idx] = val
        }

        oldData := storage.EncodeRow(row)
        newData := storage.EncodeRow(updated)

        // WAL before modifying
        e.wal.Append(wal.Record{
            TxnID:   e.txnID,
            Type:    wal.RecordUpdate,
            OldData: oldData,
            NewData: newData,
        })

        // In production: use RowID from scan. For simplicity, re-insert.
        _ = i
        heap.InsertRow(updated)
        affected++
    }

    e.wal.Append(wal.Record{TxnID: e.txnID, Type: wal.RecordCommit})
    e.wal.Flush()

    return &Result{Affected: affected}, nil
}

func (e *Executor) execDelete(stmt *DeleteStmt) (*Result, error) {
    tbl, err := e.getTable(stmt.Table)
    if err != nil {
        return nil, err
    }

    heap := storage.NewHeap(e.dm, tbl.RootPageID)
    allRows, err := heap.ScanAll(len(tbl.Columns))
    if err != nil {
        return nil, err
    }

    var affected int64
    e.txnID++

    for _, row := range allRows {
        ok, err := e.matches(stmt.Where, row, tbl)
        if err != nil {
            return nil, err
        }
        if !ok {
            continue
        }
        e.wal.Append(wal.Record{
            TxnID:   e.txnID,
            Type:    wal.RecordDelete,
            OldData: storage.EncodeRow(row),
        })
        // In production: heap.DeleteRow(rid)
        affected++
    }

    e.wal.Append(wal.Record{TxnID: e.txnID, Type: wal.RecordCommit})
    e.wal.Flush()

    return &Result{Affected: affected}, nil
}
```

---

## 6. Executing CREATE TABLE

```go
func (e *Executor) execCreateTable(stmt *CreateTableStmt) (*Result, error) {
    // Check if table already exists
    for _, t := range e.catalog.Tables {
        if t.Name == stmt.Table {
            return nil, fmt.Errorf("table %q already exists", stmt.Table)
        }
    }

    // Allocate a root page for the new table
    rootPageID, rootPage, err := e.bp.NewPage()
    if err != nil {
        return nil, err
    }
    rootPage.Initialize(storage.PageTypeLeaf)
    e.bp.UnpinPage(rootPageID, true)

    // Build table definition
    tbl := storage.TableDef{
        Name:       stmt.Table,
        RootPageID: rootPageID,
    }
    for _, col := range stmt.Columns {
        tbl.Columns = append(tbl.Columns, storage.ColumnDef{
            Name:   col.Name,
            TypeID: col.TypeID,
        })
    }

    // Add to catalog and persist
    e.catalog.Tables = append(e.catalog.Tables, tbl)
    if err := e.dm.WriteCatalog(e.catalog); err != nil {
        return nil, err
    }

    return &Result{}, nil
}
```

---

## 7. Putting It All Together

```go
// Demo: run queries against VaultDB
func main() {
    db, _ := Open("mydb.vault")
    defer db.Close()

    exec := query.NewExecutor(db.dm, db.bp, db.wal, db.catalog)

    // Create table
    runQuery(exec, "CREATE TABLE users (id INT, name VARCHAR, age INT)")

    // Insert rows
    runQuery(exec, "INSERT INTO users (id, name, age) VALUES (1, 'Alice', 25)")
    runQuery(exec, "INSERT INTO users (id, name, age) VALUES (2, 'Bob', 30)")
    runQuery(exec, "INSERT INTO users (id, name, age) VALUES (3, 'Carol', 22)")

    // Select all
    result, _ := runQuery(exec, "SELECT * FROM users")
    printResult(result)

    // Select with WHERE
    result, _ = runQuery(exec, "SELECT name FROM users WHERE age > 24")
    printResult(result)
}

func runQuery(exec *query.Executor, sql string) (*query.Result, error) {
    stmt, err := query.Parse(sql)
    if err != nil {
        fmt.Println("parse error:", err)
        return nil, err
    }
    result, err := exec.Execute(stmt)
    if err != nil {
        fmt.Println("execute error:", err)
        return nil, err
    }
    return result, nil
}

func printResult(r *query.Result) {
    if r == nil {
        return
    }
    fmt.Println(r.Columns)
    for _, row := range r.Rows {
        for i, v := range row {
            if i > 0 {
                fmt.Print(" | ")
            }
            fmt.Print(v.String())
        }
        fmt.Println()
    }
    fmt.Printf("(%d rows)\n\n", len(r.Rows))
}
```

---

## Summary

- The executor walks the AST and calls the appropriate storage operations.
- Expression evaluation turns AST nodes into runtime `Value`s by reading from the current row.
- Before every write (INSERT/UPDATE/DELETE), we append to the WAL. After all writes, we write COMMIT + flush.
- SELECT is pure read: scan all rows, filter by WHERE, project requested columns, apply LIMIT.
- CREATE TABLE allocates a root page and adds the table to the catalog.

### Exercises

**Easy:** Add a `SHOW TABLES` statement (no AST changes needed). The executor looks at `e.catalog.Tables` and returns them as a result set with columns `[name, num_columns, root_page_id]`.

**Medium:** Add support for the `COUNT(*)` aggregate: if the SELECT list is `[count(*)]`, instead of returning rows, return a single row with the count of matching rows.

**Hard:** Add ORDER BY support. After filtering rows, sort by the specified column. Support both ascending (default) and descending (`ORDER BY col DESC`). Handle mixed types gracefully.
