# Chapter 05: How Databases Work Inside — Pages, Trees, and Indexes

Every time you search for a product on Amazon or scroll through your Instagram feed, a database finds your data in milliseconds — out of billions of rows. That is not magic. It is clever engineering built on a handful of elegant ideas. In this chapter, we will crack open the database and see exactly how it stores data, finds it fast, and survives a power cut without losing your information.

## Table of Contents

1. [How Databases Store Data on Disk](#1-how-databases-store-data-on-disk)
2. [The Heap File — The Simplest Storage](#2-the-heap-file--the-simplest-storage)
3. [Why Sequential Scanning Is Slow](#3-why-sequential-scanning-is-slow)
4. [What an Index Is](#4-what-an-index-is)
5. [B-Trees — The Data Structure That Makes Databases Fast](#5-b-trees--the-data-structure-that-makes-databases-fast)
6. [How a Query Actually Executes](#6-how-a-query-actually-executes)
7. [The Buffer Pool — Keeping Hot Pages in Memory](#7-the-buffer-pool--keeping-hot-pages-in-memory)
8. [Write-Ahead Logging — Why Databases Are Crash-Safe](#8-write-ahead-logging--why-databases-are-crash-safe)
9. [The Anatomy of a Complete Database System](#9-the-anatomy-of-a-complete-database-system)
10. [Mini Project — A Tiny Page Store in Go](#10-mini-project--a-tiny-page-store-in-go)
11. [Exercises](#exercises)
12. [Summary](#summary)

---

## 1. How Databases Store Data on Disk

### The Filing Cabinet Analogy

Imagine your school's filing cabinet. It does not hold one piece of paper per drawer. That would waste enormous space and take forever to organize. Instead, each drawer holds a hanging folder, and each folder holds many sheets of paper grouped together.

A database does exactly the same thing. It does not store one row of data per spot on your hard drive. It groups rows into fixed-size chunks called **pages** (sometimes called **blocks**).

### What Is a Page?

A **page** is a fixed-size chunk of data — almost always 4,096 bytes (4 KB) or 8,192 bytes (8 KB). Why those sizes? Because your operating system's kernel and your hard drive's hardware both think in chunks of roughly that size. If the database uses the same chunk size, it lines up perfectly with the hardware and avoids wasted work.

Think of a page like a physical page in a notebook. A notebook page has a fixed size. You write as much as fits, then flip to the next page.

```
+--------------------------------------------------+
|  PAGE 0  (4096 bytes)                            |
|  ------------------------------------------------|
|  Row 1: id=1, name="Alice", age=25               |
|  Row 2: id=2, name="Bob", age=30                 |
|  Row 3: id=3, name="Carol", age=22               |
|  ... (more rows until the page is full)          |
+--------------------------------------------------+

+--------------------------------------------------+
|  PAGE 1  (4096 bytes)                            |
|  ------------------------------------------------|
|  Row 4: id=4, name="Dave", age=27                |
|  Row 5: id=5, name="Eve", age=31                 |
|  ...                                             |
+--------------------------------------------------+
```

### Why Fixed Sizes Matter

When every page is exactly 4 KB, the database can do arithmetic to jump to any page instantly:

```
position of page N = N × 4096
```

Page 0 starts at byte 0. Page 1 starts at byte 4096. Page 1000 starts at byte 4,096,000. The database can skip straight there without reading everything in between. This is called **random access** and it is extraordinarily valuable.

### Quick Check

> 1. What is a page in a database?
> 2. If each page is 4096 bytes, where does page 7 start in the file?
> 3. Why does the database use a fixed page size instead of variable sizes?

---

## 2. The Heap File — The Simplest Storage

### The Pile on Your Desk Analogy

Imagine you are running a lemonade stand and you write each sale on a sticky note. At the end of the day, you have a pile of sticky notes. They are not sorted. Sale number 47 might be at the top, sale number 3 might be at the bottom. When your mom asks "how much did you charge Alice?" you have to search through every single note.

That pile is exactly what a **heap file** is.

### What Is a Heap File?

A **heap file** is the most basic way to store rows in a database. New rows are appended to the end of the file. There is no sorting, no ordering, no structure. Rows land wherever there is free space.

```
Heap File: users.db
+------------------+
| id=3, Alice, 25  |  <- inserted first
| id=7, Bob, 30    |  <- inserted second
| id=1, Carol, 22  |  <- inserted third
| id=5, Dave, 27   |  <- inserted fourth
+------------------+
```

Notice: the rows are not in order by `id`. They are in *insertion order*. The database just appended them as they arrived.

### The Good Parts of Heap Files

- **Inserts are fast.** Writing to the end of a file is simple and quick.
- **Simple to implement.** There is very little logic needed.

### The Problem with Heap Files

Searching is slow. If you ask for "the row where id = 5," the database has no idea where that row lives. It must read every single page from the beginning until it finds it. This is called a **full table scan** or **sequential scan**.

---

## 3. Why Sequential Scanning Is Slow

### The Phone Book Analogy

Imagine a phone book where numbers are listed in the order people called in to register — not alphabetically. To find "Bob Smith," you start at page 1 and read every name until you find Bob. If there are a million names, you might read 999,999 names before finding his.

That is a sequential scan. And it is exactly what a database does when it has no index.

### The Numbers

Suppose your `users` table has 1,000,000 rows. Each page holds 100 rows. That means 10,000 pages.

- Reading one page from an SSD takes roughly 0.1 milliseconds.
- Reading all 10,000 pages takes roughly 1 second.
- Reading from a spinning hard disk is 10x slower — about 10 seconds.

One second might sound fine. But what if you have 100 users all making requests at the same time? What if your table has 100 million rows? Sequential scanning does not scale.

### Seeing It in Go

Here is a simple simulation that shows the difference between scanning a slice (like a heap file) versus using a map (like an index). Read this carefully — every line is explained.

```go
package main

import (
	"fmt"
	"time"
)

// Row represents a single database row.
type Row struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// Build a "heap" of 1,000,000 rows in random order.
	// We use a slice because slices, like heap files, have no ordering.
	const totalRows = 1_000_000
	heap := make([]Row, totalRows)
	for i := 0; i < totalRows; i++ {
		heap[i] = Row{ID: totalRows - i, Name: "user", Age: 20 + (i % 50)}
		// Notice: ID goes from 1,000,000 down to 1.
		// The row with ID=5 is near the very end of the slice.
	}

	target := 5

	// --- Sequential Scan (like a heap file with no index) ---
	start := time.Now()
	var found *Row
	for i := range heap {
		if heap[i].ID == target {
			found = &heap[i]
			break
		}
	}
	seqDuration := time.Since(start)

	if found != nil {
		fmt.Printf("Sequential scan found ID=%d in %v\n", found.ID, seqDuration)
	}

	// --- Build a map index (like a B-Tree index) ---
	// A map lets us jump directly to the row by ID.
	index := make(map[int]*Row, totalRows)
	for i := range heap {
		index[heap[i].ID] = &heap[i]
	}

	// --- Index Lookup ---
	start = time.Now()
	found = index[target]
	idxDuration := time.Since(start)

	if found != nil {
		fmt.Printf("Index lookup   found ID=%d in %v\n", found.ID, idxDuration)
	}

	fmt.Printf("\nIndex was roughly %dx faster\n", seqDuration/idxDuration)
}
```

**Line by line:**
- `heap := make([]Row, totalRows)` — creates a slice of one million rows. This is our heap file simulation.
- The loop fills it in reverse order, so the row with `ID=5` is near the very end. The scan has to work hard.
- `for i := range heap` — this is the sequential scan. It visits every element until it finds `ID=5`.
- `index := make(map[int]*Row, totalRows)` — builds a Go map where each key is a row ID and each value is a pointer to the row.
- `found = index[target]` — this is a single hash lookup. It jumps directly to the row.

Run this program and you will see the index lookup is thousands of times faster.

### Quick Check

> 1. What is a sequential scan?
> 2. If a table has 10,000 pages and reading one page takes 0.1ms, how long does a full scan take?
> 3. Why does it matter that rows in a heap file have no ordering?

---

## 4. What an Index Is

### The Library Catalog Analogy

Imagine a library with 500,000 books. The books are shelved in the order they were donated — no sorting at all. A new donation goes to the next open shelf. If you want "Harry Potter," you would need to walk every aisle.

Now imagine a small card catalog near the entrance. Each card lists a book title and tells you the exact shelf number. You look up "Harry Potter" in the catalog (which is sorted alphabetically, so it takes seconds), then walk straight to shelf 42,817.

**The card catalog is an index.**

An **index** is a separate, smaller data structure that maps a specific column's values to the physical location (page number + slot) of the rows. You look something up in the index, it hands you an address, and you jump straight there.

### What Does an Index Store?

For a table with columns `id`, `name`, `email`, and `age`, an index on the `id` column stores entries like:

```
id=1  -> page 0, slot 3
id=2  -> page 0, slot 7
id=3  -> page 0, slot 1
id=4  -> page 1, slot 0
id=5  -> page 1, slot 2
...
```

The index is much smaller than the table (it only stores the indexed column value and an address, not all the row data), and it is kept sorted so lookups are fast.

### The Trade-Off

Indexes make reads faster but writes slightly slower. Every time you insert a row, the database must update the index too. This is like adding a new book to the library and also adding a new card to the catalog. More work per insert, but faster searches.

---

## 5. B-Trees — The Data Structure That Makes Databases Fast

### The Sorted Bookshelf That Fits in a Filing Cabinet

A sorted array works great for searching (you can use binary search — always guess the middle first). But databases cannot use a plain sorted array because inserting a new row in the middle shifts everything, which is very slow on disk.

Databases use a data structure called a **B-Tree** (Balanced Tree). It gives you the search speed of a sorted structure while keeping insertions fast.

### Understanding Trees First

A **tree** is a way of organizing data where each piece of data (called a **node**) can have **children** — pointers to other nodes below it. The top node is called the **root**. Nodes at the very bottom (with no children) are called **leaves**.

```
          [Root: 50]
         /           \
    [25]               [75]
   /    \             /    \
 [10]  [40]        [60]   [90]
```

In this tree, to find the number 60:
1. Start at root: is 60 less than 50? No, go right.
2. At 75: is 60 less than 75? Yes, go left.
3. At 60: found it!

We checked 3 nodes out of 7. That is the power of a tree.

### What Makes a B-Tree Special?

A regular binary tree has one value per node and two children. A **B-Tree** node can hold many values and many children. This is critical for databases because reading one page from disk is expensive — you want to pack as much useful data into each page-read as possible.

A B-Tree node that fits in one 4KB page might hold 100 key-value pairs and 101 child pointers. That is a very wide, shallow tree.

```
B-Tree with 3 levels can index:
  Level 0 (root):    1 node   x 100 keys = 100 keys
  Level 1:         101 nodes  x 100 keys = 10,100 keys
  Level 2 (leaves): ~10,000 nodes x 100 rows each = 1,000,000 rows

So 3 page reads can find any of 1,000,000 rows.
```

Three disk reads to find one row in a million. That is why databases are fast.

### The Structure of a B-Tree Node

Each node in a B-Tree (specifically a B+ Tree, which PostgreSQL and MySQL use) looks like this:

```
+-----------------------------------------------+
|  [key1] [ptr1] [key2] [ptr2] [key3] [ptr3]   |
|                                               |
|  key1=10, key2=30, key3=50                    |
|  ptr1 -> subtree with keys < 10               |
|  ptr2 -> subtree with keys 10-30              |
|  ptr3 -> subtree with keys 30-50              |
|  ptr4 -> subtree with keys > 50               |
+-----------------------------------------------+
```

**Leaf nodes** (at the bottom) hold the actual data or a pointer to where the row lives in the heap file. They also have pointers to the *next* leaf node, forming a linked list. This makes range queries like `WHERE age BETWEEN 20 AND 30` extremely fast — find the first matching leaf, then walk the linked list.

### Visualizing a B-Tree Search for id=5

Suppose our `users` table has millions of rows and an index on `id`:

```
Step 1: Read root page from disk.
  Root contains keys: [100, 200, 300, 400, 500]
  5 < 100, so follow the leftmost child pointer.

Step 2: Read child page from disk.
  Child contains keys: [10, 20, 30, 40, 50]
  5 < 10, so follow the leftmost child pointer.

Step 3: Read leaf page from disk.
  Leaf contains: [1, 2, 3, 4, 5, 6, 7, 8, 9]
  Found id=5! The leaf stores: "page 47, slot 3"

Step 4: Read page 47 from the heap file, slot 3.
  Return the full row: {id:5, name:"Eve", age:31}
```

Four disk reads. Done.

### Quick Check

> 1. What is the difference between a tree node and a leaf node?
> 2. Why does a B-Tree node hold many keys instead of just one?
> 3. Why are leaf nodes linked together in a B+ Tree?

---

## 6. How a Query Actually Executes

### The Restaurant Kitchen Analogy

When you order food at a restaurant, your order passes through several hands: the waiter writes it down, the kitchen manager reads it and decides which station handles it, the cook prepares it, and finally the waiter delivers it. A database query goes through a similar pipeline.

### The Life of `SELECT * FROM users WHERE id = 5`

```
Your Go program
    |
    v
[1. SQL Parser]
    Reads the text "SELECT * FROM users WHERE id = 5"
    Checks the grammar is valid
    Builds a tree structure representing the query
    |
    v
[2. Query Planner / Optimizer]
    Asks: "Does the users table have an index on id?"
    If YES: plan = "Index scan on users.id_index"
    If NO:  plan = "Sequential scan on users"
    Picks the fastest plan
    |
    v
[3. Query Executor]
    Executes the plan step by step
    If using index: asks Storage Engine for B-Tree lookup
    |
    v
[4. Storage Engine]
    Walks the B-Tree to find "id=5" -> page 47, slot 3
    Asks Buffer Pool for page 47
    |
    v
[5. Buffer Pool]
    Is page 47 in memory already? (Cache hit!)
    If yes: return it immediately
    If no: read page 47 from disk, store in memory, return it
    |
    v
[6. Result]
    Row {id:5, name:"Eve", age:31} returned to your Go program
```

Each step is a separate component. Real databases like PostgreSQL have entire source files dedicated to each one. Let us look at a tiny simulation in Go.

```go
package main

import (
	"errors"
	"fmt"
)

// Row represents one row in our users table.
type Row struct {
	ID   int
	Name string
	Age  int
}

// SimpleDB is a toy database with a heap and a map-based index.
type SimpleDB struct {
	heap  []Row          // The heap file (all rows, unsorted)
	index map[int]int    // index[id] = position in heap slice
}

// NewSimpleDB creates an empty database.
func NewSimpleDB() *SimpleDB {
	return &SimpleDB{
		heap:  make([]Row, 0),
		index: make(map[int]int),
	}
}

// Insert adds a row to the database and updates the index.
func (db *SimpleDB) Insert(row Row) {
	// Append row to the heap (just like a heap file).
	position := len(db.heap)
	db.heap = append(db.heap, row)

	// Update the index: map this ID to the heap position.
	// In a real database this updates the B-Tree instead.
	db.index[row.ID] = position

	fmt.Printf("Inserted row id=%d at heap position %d\n", row.ID, position)
}

// FindByID simulates "SELECT * FROM users WHERE id = ?"
// It uses the index to jump directly to the row.
func (db *SimpleDB) FindByID(id int) (Row, error) {
	// Step 1: Look up the ID in the index.
	// In a real database, this walks the B-Tree.
	position, ok := db.index[id]
	if !ok {
		return Row{}, errors.New("row not found")
	}

	// Step 2: Read the row from the heap using the position.
	// In a real database, this reads a specific page and slot.
	row := db.heap[position]
	return row, nil
}

// FullScan simulates a sequential scan (no index used).
func (db *SimpleDB) FullScan(targetID int) (Row, error) {
	// Must read every row, one by one.
	pagesRead := 0
	for _, row := range db.heap {
		pagesRead++ // each row read simulates reading part of a page
		if row.ID == targetID {
			fmt.Printf("Full scan: read %d rows to find id=%d\n", pagesRead, targetID)
			return row, nil
		}
	}
	return Row{}, errors.New("row not found")
}

func main() {
	db := NewSimpleDB()

	// Insert rows in non-sequential order to simulate real insertions.
	db.Insert(Row{ID: 3, Name: "Carol", Age: 22})
	db.Insert(Row{ID: 7, Name: "Bob", Age: 30})
	db.Insert(Row{ID: 1, Name: "Alice", Age: 25})
	db.Insert(Row{ID: 5, Name: "Eve", Age: 31})
	db.Insert(Row{ID: 9, Name: "Dave", Age: 27})

	fmt.Println()

	// Query using the index.
	fmt.Println("--- Index lookup for id=5 ---")
	row, err := db.FindByID(5)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Found: %+v\n", row)
	}

	fmt.Println()

	// Query using a full scan (imagine we had no index).
	fmt.Println("--- Full scan for id=5 ---")
	row, err = db.FullScan(5)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Found: %+v\n", row)
	}
}
```

**Line by line highlights:**
- `SimpleDB` has a `heap` (slice of rows, no ordering) and an `index` (map from id to position). This mirrors a real database's heap file and B-Tree index.
- `Insert` appends to the heap and then updates the index. Both happen on every insert.
- `FindByID` does a two-step lookup: check the index, then jump to that position in the heap. In a real database those two steps involve reading pages.
- `FullScan` ignores the index and reads every row. Notice the counter — it counts how many rows it had to read before finding the target.

### Quick Check

> 1. What does the query planner do?
> 2. What happens when the query executor asks for a page that is already in memory?
> 3. Name the four major steps in the life of a SELECT query.

---

## 7. The Buffer Pool — Keeping Hot Pages in Memory

### The Desk Analogy

Your desk holds a limited number of books open at once. When you need a book that is already on your desk, you grab it instantly. When you need a book that is not on your desk, you must walk to the shelf, get it, and bring it back — but now your desk might be full, so you have to put one of the existing books back first.

The **buffer pool** is the database's desk. It is a fixed amount of RAM (memory) that the database uses to keep recently accessed pages available without going to disk.

### Why This Matters

RAM access takes roughly 0.0001 milliseconds. Disk access takes roughly 0.1 milliseconds. That is a 1000x difference.

If the same page is read again and again (for example, the root page of a B-Tree is accessed by *every* query), keeping it in the buffer pool means those reads are 1000x faster.

### The LRU Policy — Least Recently Used

When the buffer pool is full and needs to bring in a new page, it must **evict** (remove) an old page to make room. The most common strategy is **LRU (Least Recently Used)**: evict the page that was accessed least recently. The intuition is that pages you have not needed in a while are less likely to be needed soon.

```
Buffer pool (capacity: 3 pages)

Initial state: [page2, page5, page8]

Request page10:
  - page10 not in pool (cache miss)
  - evict least recently used: page2
  - load page10 from disk
  - pool is now: [page10, page5, page8]

Request page5:
  - page5 IS in pool (cache hit!)
  - mark page5 as recently used
  - pool order: [page5, page10, page8]

Request page3:
  - page3 not in pool (cache miss)
  - evict least recently used: page8
  - load page3 from disk
  - pool is now: [page3, page5, page10]
```

### A Simple Buffer Pool in Go

```go
package main

import (
	"container/list"
	"fmt"
)

const PageSize = 4096 // 4 KB per page

// Page represents one page of data from disk.
type Page struct {
	PageNumber int
	Data       [PageSize]byte
}

// BufferPool holds recently used pages in memory.
type BufferPool struct {
	capacity int
	pages    map[int]*list.Element // pageNumber -> list element
	lruList  *list.List            // front = most recent, back = least recent
}

// NewBufferPool creates a pool that can hold 'capacity' pages.
func NewBufferPool(capacity int) *BufferPool {
	return &BufferPool{
		capacity: capacity,
		pages:    make(map[int]*list.Element),
		lruList:  list.New(),
	}
}

// fetchFromDisk simulates reading a page from disk.
// In a real database, this would do actual file I/O.
func fetchFromDisk(pageNumber int) *Page {
	fmt.Printf("  [DISK READ] Loading page %d from disk\n", pageNumber)
	page := &Page{PageNumber: pageNumber}
	// Fill with fake data to simulate a real page.
	page.Data[0] = byte(pageNumber % 256)
	return page
}

// GetPage returns the requested page, loading from disk if needed.
func (bp *BufferPool) GetPage(pageNumber int) *Page {
	// Check if the page is already in the pool (cache hit).
	if elem, ok := bp.pages[pageNumber]; ok {
		fmt.Printf("  [CACHE HIT] Page %d is already in memory\n", pageNumber)
		// Move to front (mark as most recently used).
		bp.lruList.MoveToFront(elem)
		return elem.Value.(*Page)
	}

	// Cache miss: need to load from disk.
	page := fetchFromDisk(pageNumber)

	// If pool is full, evict the least recently used page.
	if bp.lruList.Len() >= bp.capacity {
		// The back of the list is the least recently used.
		oldest := bp.lruList.Back()
		if oldest != nil {
			evicted := oldest.Value.(*Page)
			fmt.Printf("  [EVICT]    Removing page %d from pool\n", evicted.PageNumber)
			bp.lruList.Remove(oldest)
			delete(bp.pages, evicted.PageNumber)
		}
	}

	// Add the new page to the front (most recently used position).
	elem := bp.lruList.PushFront(page)
	bp.pages[pageNumber] = elem

	return page
}

func main() {
	// Create a pool that holds 3 pages at once.
	pool := NewBufferPool(3)

	fmt.Println("Requesting page 2:")
	pool.GetPage(2)

	fmt.Println("\nRequesting page 5:")
	pool.GetPage(5)

	fmt.Println("\nRequesting page 8:")
	pool.GetPage(8)

	fmt.Println("\nRequesting page 5 again (should be a cache hit):")
	pool.GetPage(5)

	fmt.Println("\nRequesting page 10 (pool is full, must evict):")
	pool.GetPage(10)

	fmt.Println("\nRequesting page 2 (was evicted, must reload from disk):")
	pool.GetPage(2)
}
```

**Key ideas in this code:**
- `container/list` is Go's built-in doubly-linked list. It lets us efficiently move items to the front and remove from the back.
- `pages map[int]*list.Element` maps a page number to its position in the linked list. This lets us check "is this page in the pool?" in O(1) time.
- `MoveToFront` on a cache hit keeps the LRU order correct without rebuilding anything.
- When the pool is full, `bp.lruList.Back()` gives us the victim to evict instantly.

---

## 8. Write-Ahead Logging — Why Databases Are Crash-Safe

### The Surgeon's Checklist Analogy

Before a surgeon makes a single cut, they write down every step they plan to take. If something goes wrong mid-operation, the next surgeon can read the checklist and know exactly what was done and what was not. Nothing is forgotten.

Databases face the same problem. What if the power goes out in the middle of writing a transaction? Some pages might be updated, others might not. The database would be in an inconsistent, corrupt state.

The solution is **Write-Ahead Logging**, almost always called **WAL**.

### What Is WAL?

The rule is simple and strict:

> **Before you change any page on disk, first write what you are about to do into a log file.**

The log file (called the WAL) is an append-only file. Appending to a file is very fast because you only ever write to the end.

### The WAL Workflow

Here is what happens when you run `UPDATE users SET name='Eve2' WHERE id=5`:

```
Step 1: Write a log entry BEFORE touching any data page.
  Log: "Transaction 42: about to update page 47, slot 3,
        old value: {id:5, name:'Eve', age:31}
        new value: {id:5, name:'Eve2', age:31}"

Step 2: Update the page in the buffer pool (in memory only, not yet on disk).

Step 3: Later, the buffer pool "flushes" the updated page to disk.
  The page is now permanently changed.

Step 4: Write a "commit" record to the log.
  Log: "Transaction 42: COMMITTED"
```

### What Happens If the Power Goes Out?

**Scenario A: Power out after Step 1 but before Step 3.**
- The data page was never written to disk (it was only in RAM, which was lost).
- On restart, the database reads the log: it sees Transaction 42 started but has no COMMIT record.
- Decision: **roll back** — pretend the transaction never happened.
- The data is as if the UPDATE never ran.

**Scenario B: Power out after Step 3 but before Step 4.**
- The data page was written but the COMMIT record was not.
- On restart: no COMMIT record found, so roll back.
- The page on disk gets restored using the "old value" in the log entry.

**Scenario C: Power out after Step 4.**
- Everything is fine. The log confirms the transaction committed, the data page is on disk.
- If needed, the database can **redo** the change using the log.

### WAL Means Two Writes, But That Is Okay

Writing to the WAL adds an extra write per transaction. But appending to a log file is much faster than random writes to data pages. Databases batch up many log entries and write them in one big sequential write — sequential writes are the fastest possible disk operation.

### A Simplified WAL in Go

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// LogEntry represents one record in the WAL.
type LogEntry struct {
	TxID      int    `json:"tx_id"`
	Operation string `json:"operation"` // "UPDATE", "INSERT", "COMMIT", "ROLLBACK"
	TableName string `json:"table"`
	RowID     int    `json:"row_id"`
	OldValue  string `json:"old_value,omitempty"`
	NewValue  string `json:"new_value,omitempty"`
}

// WAL manages writing to the write-ahead log.
type WAL struct {
	file *os.File
}

// OpenWAL opens (or creates) the WAL file.
func OpenWAL(path string) (*WAL, error) {
	// os.O_APPEND means every write goes to the end of the file.
	// os.O_CREATE creates the file if it does not exist.
	// os.O_WRONLY means write-only.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: f}, nil
}

// Write appends a log entry to the WAL file.
// This MUST happen before any data page is modified.
func (w *WAL) Write(entry LogEntry) error {
	// Marshal the entry to JSON for simplicity.
	// Real databases use binary formats for speed.
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	// Append the JSON + newline to the file.
	_, err = fmt.Fprintf(w.file, "%s\n", data)
	return err
}

// Close closes the WAL file.
func (w *WAL) Close() {
	w.file.Close()
}

// ReadWAL reads all log entries back (used during crash recovery).
func ReadWAL(path string) ([]LogEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []LogEntry
	// Each line is one JSON entry.
	start := 0
	for i, b := range data {
		if b == '\n' {
			line := data[start:i]
			start = i + 1
			if len(line) == 0 {
				continue
			}
			var entry LogEntry
			if err := json.Unmarshal(line, &entry); err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func main() {
	walPath := "/tmp/demo.wal"
	os.Remove(walPath) // Start fresh for this demo.

	wal, err := OpenWAL(walPath)
	if err != nil {
		fmt.Println("Error opening WAL:", err)
		return
	}
	defer wal.Close()

	// Simulate an UPDATE transaction.
	txID := 42

	// Step 1: Write the intent to the log BEFORE touching data.
	wal.Write(LogEntry{
		TxID:      txID,
		Operation: "UPDATE",
		TableName: "users",
		RowID:     5,
		OldValue:  `{"id":5,"name":"Eve","age":31}`,
		NewValue:  `{"id":5,"name":"Eve2","age":31}`,
	})
	fmt.Println("WAL: wrote UPDATE intent for tx", txID)

	// Step 2: (In a real database, we would now update the buffer pool page.)
	fmt.Println("DATA: updated page in memory (buffer pool)")

	// Step 3: Write the COMMIT record.
	wal.Write(LogEntry{
		TxID:      txID,
		Operation: "COMMIT",
		TableName: "users",
	})
	fmt.Println("WAL: wrote COMMIT for tx", txID)

	// --- Simulate crash recovery ---
	fmt.Println("\n--- Reading WAL for crash recovery ---")
	entries, err := ReadWAL(walPath)
	if err != nil {
		fmt.Println("Error reading WAL:", err)
		return
	}

	committed := make(map[int]bool)
	for _, e := range entries {
		if e.Operation == "COMMIT" {
			committed[e.TxID] = true
		}
	}

	for _, e := range entries {
		if e.Operation == "UPDATE" || e.Operation == "INSERT" {
			if committed[e.TxID] {
				fmt.Printf("TX %d: REDO %s on %s row %d\n", e.TxID, e.Operation, e.TableName, e.RowID)
			} else {
				fmt.Printf("TX %d: ROLLBACK %s on %s row %d (no commit found)\n", e.TxID, e.Operation, e.TableName, e.RowID)
			}
		}
	}
}
```

This toy WAL shows the core pattern: write intent first, then do the work, then write commit. On recovery, check for commits and redo or rollback accordingly.

### Quick Check

> 1. What does "write-ahead" mean in Write-Ahead Logging?
> 2. If the power goes out after writing the log entry but before updating the data page, what does the database do on restart?
> 3. Why is appending to a log file faster than random writes to data pages?

---

## 9. The Anatomy of a Complete Database System

### The Factory Analogy

A modern factory has separate departments: raw materials storage, the assembly line, quality control, and shipping. Each department is specialized and hands its output to the next. A database is similar — it has distinct components, each with a specific job.

Here is the full picture:

```
+------------------------------------------------------------------+
|                        DATABASE SYSTEM                           |
|                                                                  |
|  +------------------+     +---------------------------+         |
|  |  Query Engine    |     |  Transaction Manager      |         |
|  |                  |     |                           |         |
|  |  SQL Parser      |     |  Concurrency Control      |         |
|  |  Query Planner   |     |  (who can read/write now) |         |
|  |  Query Executor  |     |  Lock Manager             |         |
|  +--------+---------+     +-------------+-------------+         |
|           |                             |                        |
|           v                             v                        |
|  +----------------------------------------------------------+   |
|  |                    Storage Engine                         |   |
|  |                                                          |   |
|  |   Buffer Pool   <-->  B-Tree Index  <-->  Heap Files     |   |
|  |        |                                                  |   |
|  |        v                                                  |   |
|  |   Write-Ahead Log (WAL)                                  |   |
|  |        |                                                  |   |
|  |        v                                                  |   |
|  |   [ Disk / SSD ]                                         |   |
|  +----------------------------------------------------------+   |
+------------------------------------------------------------------+
```

### The Three Main Subsystems

**1. The Query Engine**

This is the brain. It receives SQL text, parses it into a structured format, decides the best execution plan (should we use the index or do a full scan?), and executes that plan step by step. The query planner is one of the most complex parts of a database — PostgreSQL's planner is tens of thousands of lines of code.

**2. The Transaction Manager**

This is the traffic controller. When multiple users query the database at the same time, the transaction manager ensures they do not interfere with each other. It uses **locks** (preventing two transactions from modifying the same row simultaneously) and **MVCC** (Multi-Version Concurrency Control, which lets readers see a consistent snapshot of the data even while writers are making changes). This is what makes database transactions **ACID** compliant.

**3. The Storage Engine**

This is the muscle. It handles everything to do with data on disk: reading and writing pages, managing the buffer pool, maintaining B-Tree indexes, and writing to the WAL. MySQL lets you swap storage engines — the famous InnoDB and MyISAM are both storage engines. PostgreSQL has one tightly integrated storage engine.

### The Flow of a Write Operation

```
Your Go app calls: db.Exec("UPDATE users SET name='Eve2' WHERE id=5")

1. Query Engine parses SQL -> builds an execution plan
2. Transaction Manager starts a transaction, acquires a lock on row id=5
3. Storage Engine: WAL writes the "before" and "after" images
4. Storage Engine: Buffer Pool finds the page, updates it in memory
5. Transaction Manager: releases the lock, marks transaction committed
6. Storage Engine: WAL writes COMMIT record
7. Eventually: Buffer Pool flushes the dirty page to disk
```

---

## 10. Mini Project — A Tiny Page Store in Go

Let us build a minimal page-based storage file. This is the foundation of a real storage engine: the ability to read and write fixed-size pages to and from a file.

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	PageSizeBytes = 4096 // Each page is 4 KB
)

// PageStore manages a file of fixed-size pages.
type PageStore struct {
	file      *os.File
	pageCount int
}

// OpenPageStore opens an existing file or creates a new one.
func OpenPageStore(path string) (*PageStore, error) {
	// os.O_RDWR = read+write, os.O_CREATE = create if missing
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// Figure out how many pages already exist in the file.
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	pageCount := int(info.Size()) / PageSizeBytes

	return &PageStore{file: f, pageCount: pageCount}, nil
}

// AllocatePage adds a new empty page at the end of the file.
// Returns the page number of the newly allocated page.
func (ps *PageStore) AllocatePage() (int, error) {
	newPageNum := ps.pageCount
	// Create an empty (zeroed) page.
	emptyPage := make([]byte, PageSizeBytes)
	// Seek to where the new page should start.
	offset := int64(newPageNum) * PageSizeBytes
	_, err := ps.file.WriteAt(emptyPage, offset)
	if err != nil {
		return 0, err
	}
	ps.pageCount++
	return newPageNum, nil
}

// WritePage writes data into a specific page number.
// Data must be exactly PageSizeBytes bytes.
func (ps *PageStore) WritePage(pageNum int, data []byte) error {
	if len(data) != PageSizeBytes {
		return errors.New("data must be exactly PageSizeBytes bytes")
	}
	if pageNum >= ps.pageCount {
		return errors.New("page number out of range")
	}
	offset := int64(pageNum) * PageSizeBytes
	_, err := ps.file.WriteAt(data, offset)
	return err
}

// ReadPage reads a full page from disk by page number.
func (ps *PageStore) ReadPage(pageNum int) ([]byte, error) {
	if pageNum >= ps.pageCount {
		return nil, errors.New("page number out of range")
	}
	data := make([]byte, PageSizeBytes)
	offset := int64(pageNum) * PageSizeBytes
	_, err := ps.file.ReadAt(data, offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Close closes the underlying file.
func (ps *PageStore) Close() {
	ps.file.Close()
}

func main() {
	storePath := "/tmp/mystore.db"
	os.Remove(storePath) // Start fresh.

	store, err := OpenPageStore(storePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer store.Close()

	// Allocate two pages.
	page0, _ := store.AllocatePage()
	page1, _ := store.AllocatePage()
	fmt.Printf("Allocated page %d and page %d\n", page0, page1)

	// Write a message into page 0.
	// We create a 4096-byte slice and copy our message into the start.
	data0 := make([]byte, PageSizeBytes)
	copy(data0, []byte("Hello from page 0! This is row data."))
	store.WritePage(page0, data0)
	fmt.Println("Wrote data to page 0")

	// Write a different message into page 1.
	data1 := make([]byte, PageSizeBytes)
	copy(data1, []byte("Page 1: more rows stored here."))
	store.WritePage(page1, data1)
	fmt.Println("Wrote data to page 1")

	// Read page 0 back.
	result, err := store.ReadPage(page0)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	// The message starts at the beginning; trim the zero padding.
	msg := string(result[:36])
	fmt.Printf("Read from page 0: %q\n", msg)

	// Demonstrate that we can jump directly to page 1 without reading page 0.
	result1, _ := store.ReadPage(page1)
	msg1 := string(result1[:30])
	fmt.Printf("Read from page 1: %q\n", msg1)

	fmt.Printf("\nTotal pages in store: %d\n", store.pageCount)
	fmt.Printf("Total file size: %d bytes (%d KB)\n",
		store.pageCount*PageSizeBytes,
		store.pageCount*PageSizeBytes/1024)
}
```

**What you built:** A real file on disk organized into fixed 4KB pages. You can allocate new pages, write data to any page, and read any page back by number — jumping directly to it using arithmetic. This is the beating heart of every storage engine.

**Extension ideas:**
- Add a page header (first 8 bytes of each page) that stores how many rows are in the page.
- Add a function to write and read `Row` structs by serializing them to bytes.
- Add a simple in-memory map that acts as a page cache (buffer pool).

---

## Exercises

### Easy

1. **Page arithmetic.** If each page is 4096 bytes, at what byte offset does page 15 start? Write a short Go program that prints the offset for pages 0 through 20.

2. **Heap insert simulation.** Create a Go slice of `Row` structs. Write a function `Insert(rows []Row, r Row) []Row` that appends the row to the slice. Write a second function `ScanForID(rows []Row, id int) *Row` that does a full sequential scan. Insert 10 rows with random IDs (not in order) and search for a specific one.

3. **WAL reading.** Modify the WAL example from this chapter so that it simulates a crash: write two transactions to the WAL, but only write the COMMIT record for the first one (simulate the power going out mid-second-transaction). Then run the recovery logic and confirm it identifies the first as "REDO" and the second as "ROLLBACK."

### Medium

4. **Index vs. no index benchmark.** Write a Go benchmark (using the `testing` package's `Benchmark` functions) that compares sequential scan versus map-based index lookup for a slice of 100,000 rows. Use `go test -bench=.` to run it and report the results in a comment at the top of your file.

5. **Buffer pool with dirty tracking.** Extend the `BufferPool` from this chapter. Add a `dirty` flag to each cached page (a dirty page has been modified but not yet written to disk). Add a `FlushDirtyPages(store *PageStore)` method that writes all dirty pages back to disk and clears their dirty flag.

6. **Simple B-Tree search.** Implement a binary search function over a sorted slice of `int` keys that simulates one "level" of a B-Tree lookup. The function should take a slice of keys (sorted), a target, and return the index of the target or -1 if not found. Then simulate a three-level B-Tree by having each level call the search function on a smaller slice.

### Hard

7. **Persistent key-value store.** Using the `PageStore` from the Mini Project, build a simple persistent key-value store. Each page stores up to 10 key-value pairs (both strings, max 100 bytes each). Implement `Put(key, value string)` and `Get(key string) (string, bool)`. The data must survive the program restarting — store it on disk and reload it on `OpenPageStore`.

8. **WAL-based recovery.** Extend your key-value store from exercise 7 to use a WAL. Before any `Put` operation, write a log entry. On startup, read the WAL and replay committed transactions to bring the store to the correct state. Simulate a crash by writing log entries but terminating the process before the COMMIT entry, then confirm the store ignores those changes on restart.

---

## Summary

- Databases store data in fixed-size **pages** (typically 4KB or 8KB). Fixed sizes allow the database to calculate the exact disk offset of any page instantly using simple multiplication.

- A **heap file** stores rows in insertion order with no sorting. It makes inserts fast but forces a **sequential scan** (reading every page) when searching without an index — this becomes unacceptably slow for large tables.

- An **index** is a separate, sorted data structure (almost always a **B-Tree**) that maps a column's values to the physical location of rows. A B-Tree search on a million-row table requires only 3-4 disk reads.

- The **buffer pool** is a cache of recently used pages kept in RAM. Since RAM is 1000x faster than disk, keeping hot pages in memory dramatically reduces query latency. An LRU eviction policy keeps the most recently used pages resident.

- **Write-Ahead Logging (WAL)** makes databases crash-safe. Every change is recorded in a log file *before* the data page is modified. On crash recovery, the database replays committed transactions and rolls back incomplete ones, guaranteeing no data loss for committed writes.
