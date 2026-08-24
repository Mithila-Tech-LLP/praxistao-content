# Chapter 02: How Storage Works — Bits, Bytes, Files, and Disks

Every database ultimately stores data on a disk. Before you can understand databases, you need to understand what a disk is, how data is written to it, and why certain operations are slow. This chapter builds the physical intuition that every database concept rests on.

## Table of Contents

1. Bits and Bytes — The Language of All Storage
2. Files and the File System
3. Hard Drives vs SSDs — Physical Reality
4. Sequential vs Random Access
5. I/O in Go — Reading and Writing Files
6. Why Databases Can't Just Use Files
7. Mini Project: A Key-Value Store in a Single File
8. Exercises

---

## 1. Bits and Bytes — The Language of All Storage

Think of a light switch. It has two states: on or off. Computers store everything using this same principle. A **bit** is a single on/off value — 0 or 1.

One bit alone isn't very useful. But group 8 bits together and you get a **byte**. A byte can represent 256 different values (2⁸ = 256), which is enough to represent any letter, digit, or punctuation mark.

```
Bit:  0 or 1
Byte: 8 bits  → 256 possible values (0–255)
KB:   1024 bytes  (kilobyte)
MB:   1024 KB     (megabyte)
GB:   1024 MB     (gigabyte)
TB:   1024 GB     (terabyte)
```

The letter 'A' is stored as the byte `01000001` (65 in decimal). The number 42 is `00101010`. Every piece of data — a name, a price, a photo — is ultimately a sequence of bytes.

---

## 2. Files and the File System

When you save a file, the operating system (OS) figures out where on the disk to put it. This job belongs to the **file system** (like NTFS on Windows or ext4 on Linux). The file system keeps a directory — a map from file names to physical locations on disk.

Disks are divided into **blocks** (also called **sectors** or **pages**), typically 512 bytes or 4096 bytes in size. The file system allocates blocks to files. When a file grows, more blocks are allocated. When deleted, blocks are freed for reuse.

This has an important implication: **data in a file is not always stored in consecutive blocks on disk**. Large files get fragmented across multiple locations. Reading them means jumping around — which matters for performance.

### File Metadata

Every file has metadata stored separately from its content:
- Name
- Size
- Created/modified timestamps
- Permissions
- Location of data blocks (the inode on Linux)

Reading a file means: look up the inode → find the block locations → read each block from disk.

---

## 3. Hard Drives vs SSDs — Physical Reality

### Hard Disk Drives (HDD)

A hard drive has spinning magnetic platters and a read/write head that physically moves across the platter surface.

```
           Platter (spinning at 7200 RPM)
          ┌─────────────────────────────┐
          │  track 0                    │
          │    track 1                  │
          │      track 2  ← head moves  │
          │        ...                  │
          └─────────────────────────────┘
```

To read data at a given position:
1. **Seek time**: the head physically moves to the right track (~5–15 ms)
2. **Rotational latency**: wait for the platter to spin to the right sector (~2–5 ms)
3. **Transfer time**: read the data (fast once under the head)

Random reads on a spinning disk are painful — each one requires a physical seek. Sequential reads (reading data laid out in order) are much faster because the head doesn't need to move.

### Solid State Drives (SSD)

SSDs have no moving parts. Data is stored in NAND flash cells. There's no seek or rotational latency — any location can be accessed in microseconds.

| Operation        | HDD         | SSD         |
|-----------------|-------------|-------------|
| Random read     | ~10 ms      | ~0.1 ms     |
| Sequential read | ~100 MB/s   | ~500 MB/s   |
| Price per GB    | cheap       | more expensive |

Even SSDs prefer sequential access over random access (due to internal architecture), but the difference is far smaller than HDDs.

**Key insight for databases**: random I/O is expensive. Every database design tries to minimize random disk seeks and maximize sequential reads and writes.

---

## 4. Sequential vs Random Access

Imagine a library. **Sequential access** is reading a book cover-to-cover — you go through pages in order. **Random access** is jumping to page 347, then page 12, then page 891.

For a spinning disk:
- Sequential: the head moves once, then reads continuously. Fast.
- Random: the head moves for every read. Each move costs ~10ms. Reading 1000 random records = ~10 seconds.

For SSDs this gap is smaller, but sequential writes are still faster because SSDs write in large blocks internally.

**What this means for databases:**
- Database pages should be read in order when possible
- Indexes help avoid scanning every row sequentially
- Writes should be batched and sequential where possible (this is why the WAL exists — more on that in Chapter 06)

---

## 5. I/O in Go — Reading and Writing Files

Go's standard library makes file I/O straightforward. Let's see the basics.

### Writing a file

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Create (or overwrite) a file
    file, err := os.Create("notes.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close() // always close the file when done

    // Write bytes to the file
    _, err = file.Write([]byte("Hello, storage!\n"))
    if err != nil {
        fmt.Println("Error writing:", err)
        return
    }

    fmt.Println("File written successfully")
}
```

`defer file.Close()` ensures the file is closed even if the function exits early due to an error. This is a Go idiom you will see everywhere.

### Reading a file

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Read the entire file into memory
    data, err := os.ReadFile("notes.txt")
    if err != nil {
        fmt.Println("Error reading:", err)
        return
    }

    fmt.Println("File contents:", string(data))
}
```

### Seeking — Jumping to a specific position

For databases, we often need to jump to a specific byte position in a file. This is called **seeking**.

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    file, _ := os.Open("notes.txt")
    defer file.Close()

    // Seek to byte position 7
    _, err := file.Seek(7, 0) // 0 = seek from start
    if err != nil {
        fmt.Println("Error seeking:", err)
        return
    }

    // Read 8 bytes from position 7
    buf := make([]byte, 8)
    n, _ := file.Read(buf)
    fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
}
```

Seeking is what databases do constantly — jump to a B-tree node at byte offset 16384, read 4096 bytes, process it. The efficiency of seeks determines a large part of database performance.

### Buffered I/O

Raw file reads and writes go all the way to the OS kernel. For small, frequent writes, buffered I/O batches them up:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, _ := os.Create("buffered.txt")
    defer file.Close()

    writer := bufio.NewWriter(file)

    // These writes go into an in-memory buffer first
    for i := 0; i < 100; i++ {
        fmt.Fprintf(writer, "line %d\n", i)
    }

    // Flush the buffer to disk in one go
    writer.Flush()

    fmt.Println("Done")
}
```

Buffered writes are faster because they batch small writes into larger disk operations.

---

## 6. Why Databases Can't Just Use Files

If files exist, why do we need databases at all? Let's try to build a simple user database with plain text files.

```go
// Store users as "id,name,email\n" in a text file
// users.txt:
// 1,Alice,alice@example.com
// 2,Bob,bob@example.com
// 3,Charlie,charlie@example.com
```

**Problem 1: Finding a user by ID requires reading the whole file**

To find user with ID 5000, you must read lines 1 through 4999 first. With a million users, that's slow.

**Problem 2: No concurrent access protection**

If two programs write to the file at the same time, the file gets corrupted. Databases handle this with locks and transactions.

**Problem 3: No crash safety**

If the power goes out mid-write, the file might be half-written. Databases use write-ahead logging to recover from crashes.

**Problem 4: No query language**

"Find all users who signed up this month from France" requires custom code. SQL gives you a standard language.

**Problem 5: Scaling**

A text file can't be efficiently split across multiple servers. Databases have replication and sharding built in.

These limitations are exactly what databases solve. Every feature of a database — indexes, transactions, SQL, replication — exists to solve one of these file-based problems.

---

## 7. Mini Project: A Key-Value Store in a Single File

Let's build a tiny key-value store that saves data to a file. This is exactly what early databases like dbm looked like.

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

const dbFile = "store.db"

// Set writes key=value to the file (append-only)
func Set(key, value string) error {
    file, err := os.OpenFile(dbFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = fmt.Fprintf(file, "%s=%s\n", key, value)
    return err
}

// Get reads the file from start and returns the LAST value for key
// (later writes override earlier ones)
func Get(key string) (string, bool) {
    file, err := os.Open(dbFile)
    if err != nil {
        return "", false
    }
    defer file.Close()

    var found string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 && parts[0] == key {
            found = parts[1] // keep updating — last one wins
        }
    }

    if found == "" {
        return "", false
    }
    return found, true
}

func main() {
    Set("name", "Alice")
    Set("age", "25")
    Set("name", "Alicia") // update name

    name, _ := Get("name")
    age, _ := Get("age")
    fmt.Println("name:", name) // Alicia
    fmt.Println("age:", age)   // 25
}
```

This is a real, working key-value store. But it has problems:
- `Get` is O(n) — reads the entire file every time
- The file grows forever (we never compact it)
- No crash protection

These are exactly the problems that real storage engines (like RocksDB's LSM-tree) solve. We will build a proper storage engine in Chapter 39.

---

## Summary

- A **bit** is a 0 or 1. A **byte** is 8 bits. All data is bytes.
- **File systems** map file names to physical disk locations.
- **HDDs** are slow for random access (moving head). **SSDs** are much faster.
- **Sequential access** (reading data in order) is always faster than **random access**.
- Go's `os` package handles file I/O with `Read`, `Write`, and `Seek`.
- Plain files fail at scale, concurrency, crash safety, and querying — databases solve all four.

### Quick Check

1. Why is random disk access slower than sequential access on an HDD?
2. What does `file.Seek(100, 0)` do?
3. Name two problems with storing a database in a plain text file.

### Exercises

**Easy:** Write a Go program that creates a file, writes 5 lines to it, then reads it back and prints each line.

**Medium:** Extend the mini key-value store to support a `Delete(key)` operation. How would you mark a key as deleted in an append-only log?

**Hard:** The current `Get` reads the whole file every time. Add an in-memory index (a `map[string]int64`) that maps each key to its byte offset in the file. How do you keep this index up to date after a restart?
