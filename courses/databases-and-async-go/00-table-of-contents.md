# Databases & Asynchronous Systems in Go
### From Zero to Building Your Own Database, Async System, and Production Service in Go

---

## What You Will Build

By the end of this course, you will have built three real, working systems from scratch using the Go programming language. Not toy demos. Not hello-world scripts. Real systems that mirror what engineers at companies like Google, Stripe, and Shopify build every day.

Here is what you are going to build:

---

### VaultDB — Your Own Database Engine

Imagine you could build your own version of PostgreSQL. Not all of it — but enough that you deeply understand how every database you ever use actually works under the hood.

VaultDB is a production-ready database engine written in Go. It supports SQL queries, stores data to disk, indexes records for fast lookup, handles multiple users at the same time, and recovers safely after a crash. By the time you finish Volume 8, you will have built all of it.

Here is what using VaultDB looks like:

```go
package main

import (
    "fmt"
    "log"

    "github.com/you/vaultdb"
)

func main() {
    // Open (or create) a database file on disk.
    // Think of this like opening a notebook where all your data lives.
    db, err := vaultdb.Open("myapp.vdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close() // Close the notebook when we are done.

    // Create a table. A table is like a spreadsheet with named columns.
    // "id" is a number that uniquely identifies each row.
    // "name" is text (a string). "age" is a whole number.
    err = db.Exec(`
        CREATE TABLE users (
            id   INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            age  INTEGER
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Insert some rows. Each row is one person.
    // We are storing three people in our users table.
    db.Exec(`INSERT INTO users VALUES (1, 'Aisha', 17)`)
    db.Exec(`INSERT INTO users VALUES (2, 'Carlos', 19)`)
    db.Exec(`INSERT INTO users VALUES (3, 'Priya', 16)`)

    // Query rows. This says: "give me the name and age of everyone under 18."
    // The database will find matching rows without scanning every row by hand.
    rows, err := db.Query(`SELECT name, age FROM users WHERE age < 18`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Loop through the results and print each matching row.
    // rows.Next() advances to the next result row and returns true if there is one.
    for rows.Next() {
        var name string
        var age int
        // Scan reads the current row's columns into our Go variables.
        rows.Scan(&name, &age)
        fmt.Printf("%s is %d years old\n", name, age)
    }
    // Output:
    // Aisha is 17 years old
    // Priya is 16 years old
}
```

Every single line of `vaultdb` — the file format, the query parser, the index, the storage engine — you will have written yourself.

---

### StreamFlow — Your Own Async Message Broker

Now imagine you could build your own version of Apache Kafka — the system that companies like LinkedIn, Uber, and Netflix use to move billions of messages per day between their services.

StreamFlow is a production-ready async message broker written in Go. It lets programs publish messages to named topics, and other programs subscribe to those topics and receive messages in order. It stores messages durably on disk, replays old messages on demand, and handles thousands of connections at once.

Here is what using StreamFlow looks like:

```go
// --- PUBLISHER: the program that sends messages ---

package main

import (
    "log"

    "github.com/you/streamflow/client"
)

func main() {
    // Connect to a running StreamFlow broker.
    // A broker is the central server that receives, stores, and delivers messages.
    // "localhost:7700" is the address where our broker is listening.
    sf, err := client.Connect("localhost:7700")
    if err != nil {
        log.Fatal(err)
    }
    defer sf.Close()

    // Publish a message to the "orders" topic.
    // A topic is like a named mailbox channel. Anyone can drop a message in.
    // Any program subscribed to "orders" will receive this message.
    // The message body is just bytes — here we use JSON-formatted text.
    err = sf.Publish("orders", []byte(`{"item": "laptop", "qty": 1}`))
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Order published.")
}
```

```go
// --- SUBSCRIBER: the program that receives messages ---

package main

import (
    "fmt"
    "log"

    "github.com/you/streamflow/client"
)

func main() {
    // Connect to the same broker.
    sf, err := client.Connect("localhost:7700")
    if err != nil {
        log.Fatal(err)
    }
    defer sf.Close()

    // Subscribe to the "orders" topic starting from message offset 0.
    // An offset is a position in the message log — like a line number in a book.
    // Offset 0 means "start from the very first message ever published."
    // This is powerful: we can replay the entire history of a topic any time.
    sub, err := sf.Subscribe("orders", 0)
    if err != nil {
        log.Fatal(err)
    }

    // This loop runs forever, waiting for and processing each incoming message.
    // sub.Messages() returns a Go channel that delivers one message at a time.
    for msg := range sub.Messages() {
        fmt.Printf("Received order at offset %d: %s\n", msg.Offset, msg.Body)
    }
}
```

Again — every byte of StreamFlow is code you will have written.

---

### The Final Project — A Real Analytics Service

In Part 3, you combine both systems into a working analytics service. It looks like this:

- An HTTP server receives events (page views, clicks, purchases) from a web app.
- Each event is published to StreamFlow.
- A worker subscribes to StreamFlow, reads each event, and writes it to VaultDB.
- A dashboard queries VaultDB and returns aggregated stats.

This is the exact same architecture used by real analytics platforms. You will have built the entire stack yourself, from the database to the message broker to the HTTP service.

---

## Who Is This Course For

This course is for anyone who is curious about how software really works. You do not need a computer science degree. You do not need years of experience. You do not need to have ever heard of a database before.

If you can read, think logically, and are willing to type code and run it, you can follow this course from the very first chapter to the very last. Beginners are assumed throughout — every term that is introduced is defined in plain language the moment it appears. You will never be left wondering what a word means.

Specifically, this course is for you if:

- You are new to programming and want to go deep, not just skim the surface.
- You know some programming in any language but have never worked with databases or async systems.
- You know Go but have only used it for small scripts or web handlers.
- You use databases every day but want to truly understand what happens inside them.
- You want to build something you can show to the world — a real database and a real message broker, written by you.

The Go code in this course is explained line by line. If you have never seen Go before, Chapter 1 will catch you up on everything you need. No prior Go knowledge is assumed.

---

## How Long Will This Take

The honest answer: it depends on how much time you put in. Here is a realistic breakdown:

```
PART 1 — DATABASES (Volumes 1-8)
├── Volume 1:  Why Databases?               ~  8 hours
├── Volume 2:  SQL Language                 ~ 10 hours
├── Volume 3:  PostgreSQL                   ~ 12 hours
├── Volume 4:  SQLite & MySQL               ~  8 hours
├── Volume 5:  NoSQL                        ~ 10 hours
├── Volume 6:  Distributed Databases        ~  8 hours
├── Volume 7:  Modern Databases             ~  8 hours
└── Volume 8:  Building VaultDB             ~ 24 hours
                                            ----------
                                  Total:    ~ 88 hours

PART 2 — ASYNC SYSTEMS (Volumes 9-11)
├── Volume 9:  Why Async?                   ~  6 hours
├── Volume 10: Real Async Systems           ~ 12 hours
└── Volume 11: Building StreamFlow          ~ 20 hours
                                            ----------
                                  Total:    ~ 38 hours

PART 3 — INTEGRATION (Volume 12)
└── Volume 12: Putting It All Together      ~ 14 hours
                                            ----------
                                  Total:    ~ 14 hours

GRAND TOTAL:                                ~140 hours
```

At one hour per day, that is about five months. At two hours per day, about two and a half months. Many people find this kind of deep-building work absorbing — they end up spending more time than planned because they are genuinely having fun.

---

## The Learning Progression

```
Week  1-2   |  Part 1 starts: What is data? What is a database?
             |  You store data in files. You hit problems. You understand why databases exist.

Week  3-5   |  SQL: you learn the language used to talk to databases.
             |  By the end, you can write real queries against real data.

Week  6-9   |  PostgreSQL, SQLite, MySQL: you use the big real-world databases.
             |  You connect Go programs to databases. You build a small app.

Week 10-12  |  NoSQL, distributed systems, modern databases.
             |  You see the full landscape of how data is stored at scale.

Week 13-18  |  You build VaultDB. Every chapter adds one piece.
             |  By the end you have a working database you built from scratch.

Week 19-20  |  Part 2 starts: What is async? Why does it matter?
             |  You feel the pain of synchronous systems. You understand the solution.

Week 21-24  |  Kafka, RabbitMQ, Redis Streams — real async systems.
             |  You use them with Go. You understand their internals.

Week 25-29  |  You build StreamFlow. Chapter by chapter the broker comes to life.

Week 30-32  |  Part 3: You wire VaultDB and StreamFlow together.
             |  You build and deploy the analytics service. You are done.
```

---

## Full Table of Contents

---

# PART 1 — DATABASES

---

## Volume 1: Why Databases?

> The question before the answer. Before learning what a database is, you will feel the problem it solves. You will store data in plain files, run into every painful limitation, and arrive at the idea of a database yourself.

### Chapter 01: Data Is Everywhere

What is data? How do programs remember things between runs? The story of a cookie shop that needs to track orders. Writing data to a plain text file in Go. Reading it back. Why this almost works — and where it starts to break down.

**Key topics:** what data is, files vs. memory, reading and writing files in Go, the persistence problem.

---

### Chapter 02: The Problem with Files

Your cookie shop is growing. You have 10,000 orders. Searching a text file takes forever. Two people try to write at the same time and corrupt the file. What happens when the program crashes mid-write? In this chapter you experience every real-world problem that motivates databases, by actually hitting them.

**Key topics:** linear search, file corruption, concurrency problems, crash recovery, the case for structure.

---

### Chapter 03: What Is a Database, Really?

A database is not magic. It is a program that solves exactly the problems from Chapter 02. We define the four things every database must do: store data durably, find it fast, let multiple users work at once, and survive crashes. The ACID properties — Atomicity, Consistency, Isolation, Durability — explained with a bank transfer analogy. "Atomic" means the transfer either fully happens or does not happen at all. There is no halfway.

**Key topics:** database definition, ACID (Atomicity, Consistency, Isolation, Durability), why ACID matters.

---

### Chapter 04: How Databases Store Data

Your hard drive is like a warehouse with numbered shelves. A database figures out which shelf to put data on and how to find it again quickly. Pages and blocks are the fixed-size units a database uses to read and write — like standardized boxes in that warehouse. A simple Go program that mimics this page-based storage.

**Key topics:** disk vs. memory, pages and blocks, heap files, row-oriented storage.

---

### Chapter 05: Finding Data Fast — The Index

Looking for one person in a phone book by scanning every page would take hours. The index at the back of the phone book solves this — it points you directly to the right page. Database indexes work the same way. A B-tree (balanced tree) is the data structure most databases use for their indexes — think of a library card catalog organized so every lookup takes roughly the same small number of steps regardless of how large the catalog grows.

**Key topics:** full table scan vs. index lookup, what an index is, introduction to B-trees, trade-offs of indexing.

---

### Chapter 06: A Map of the Database World

There are many kinds of databases: relational, document, key-value, graph, time-series, and more. A relational database stores data in tables with rows and columns. A document database stores flexible JSON-like objects. A key-value database is like a dictionary — you look things up by a single key. A map of the whole landscape. Which one to pick and when. A preview of everything covered in Part 1.

**Key topics:** relational vs. NoSQL, overview of database categories, when to use each.

---

## Volume 2: SQL Language

> SQL (Structured Query Language) is the language you use to talk to a relational database. It reads almost like English sentences. It is one of the most useful languages ever invented — and one of the simplest to start with. By the end of this volume you will write queries that answer real business questions.

### Chapter 07: Your First SQL Queries

SQL reads like English on purpose. `SELECT name FROM users WHERE age > 18` means exactly what it sounds like: give me the names of all users older than 18. Setting up a SQLite database (the simplest database — it is just one file on your computer, no server needed). Writing your first SELECT, INSERT, UPDATE, and DELETE statements.

**Key topics:** SELECT, FROM, WHERE, INSERT, UPDATE, DELETE, SQLite setup, running SQL from Go.

---

### Chapter 08: Designing Tables

A table is a spreadsheet with rules. Columns define the structure — each column has a name and a data type (TEXT for words, INTEGER for whole numbers, REAL for decimals). Constraints enforce rules — PRIMARY KEY means every row must have a unique identifier, NOT NULL means a column cannot be left empty. Designing a database for a school: students, classes, grades.

**Key topics:** CREATE TABLE, data types, PRIMARY KEY, NOT NULL, DEFAULT, NULL semantics.

---

### Chapter 09: Joining Tables

Real data lives across multiple tables. Joining connects them. The school database has a `students` table and a `grades` table. To find each student's average grade you need to combine data from both tables. An INNER JOIN returns only rows that match in both tables. A LEFT JOIN returns all rows from the left table plus any matches from the right — unmatched right rows show up as empty (NULL).

**Key topics:** normalization (brief intro), INNER JOIN, LEFT JOIN, JOIN ON, multi-table queries.

---

### Chapter 10: Aggregations and Groups

Counting, summing, averaging. `COUNT` counts rows. `SUM` adds up a column. `AVG` computes the average. `MIN` and `MAX` find the smallest and largest values. The `GROUP BY` clause: instead of one result per row, produce one result per group of rows. Finding the average grade per class. The `HAVING` clause filters groups — like WHERE but applied after grouping.

**Key topics:** aggregate functions, GROUP BY, HAVING, the difference between WHERE and HAVING.

---

### Chapter 11: Advanced SQL

Subqueries: a query nested inside another query, like a calculation inside a calculation. Window functions: compute a running total or rank students within their class without collapsing rows into groups — the `OVER` clause defines the window of rows each calculation sees. Common Table Expressions (CTEs): give a subquery a name using the `WITH` keyword so the main query stays readable. These are the tools that make SQL genuinely powerful for complex analysis.

**Key topics:** subqueries, EXISTS, window functions (OVER, PARTITION BY, ROW_NUMBER), CTEs (WITH).

---

### Chapter 12: SQL in Go

Connecting a Go program to a SQLite database using Go's standard `database/sql` package. The standard pattern: open a connection with `sql.Open`, execute queries with `db.Query` or `db.Exec`, read each result row with `rows.Scan`, wrap multiple operations in a transaction so they succeed or fail together. Writing a small Go CLI that queries the school database interactively.

**Key topics:** `database/sql`, `sql.Open`, `db.Query`, `rows.Scan`, `db.Exec`, transactions in Go, prepared statements.

---

## Volume 3: PostgreSQL

> PostgreSQL is the most capable open-source relational database in the world. It has been in active development since 1986. Companies like Instagram, Spotify, and GitHub store billions of rows in it. This volume teaches you to use it like a professional.

### Chapter 13: PostgreSQL vs. SQLite

SQLite is a file your program opens directly. PostgreSQL is a separate server your program connects to over a network. This server architecture means multiple programs can connect simultaneously, and the database can run on a different machine entirely. Installing and running PostgreSQL locally. Connecting from Go using the `pgx` driver. Creating users, databases, and schemas — a schema is like a namespace that groups related tables together.

**Key topics:** client-server architecture, PostgreSQL setup, `pgx` driver, schemas, roles.

---

### Chapter 14: PostgreSQL Data Types

PostgreSQL has far more data types than SQLite. Arrays let one column hold a list of values. JSONB stores arbitrary JSON documents that can be indexed and queried — giving you a relational database that can also behave like a document database. UUID is a universally unique identifier, useful for distributed systems where you cannot rely on auto-incrementing IDs. TIMESTAMP WITH TIME ZONE stores moments in time correctly across time zones.

**Key topics:** TEXT, JSONB, UUID, ARRAY, TIMESTAMP WITH TIME ZONE, ENUM, domain types.

---

### Chapter 15: Indexes in PostgreSQL

PostgreSQL offers multiple index types for different situations. A B-tree index works for comparisons (greater than, less than, equal). A Hash index works only for equality but is faster for that case. A GIN index (Generalized Inverted Index) works for full-text search and arrays. How to use `EXPLAIN ANALYZE` — PostgreSQL's tool for showing exactly what it does when it runs your query, including which indexes it uses and how long each step takes.

**Key topics:** index types, EXPLAIN ANALYZE, index-only scans, partial indexes, expression indexes.

---

### Chapter 16: Transactions and Concurrency

Two users book the last seat on a plane at the same time. Only one should succeed. PostgreSQL prevents conflicts using MVCC — Multi-Version Concurrency Control — which keeps multiple versions of each row so readers and writers do not block each other. Transaction isolation levels control how much one transaction can see of another's work. Deadlocks happen when two transactions each wait for the other to release a lock — PostgreSQL detects and breaks them automatically.

**Key topics:** transaction isolation levels, MVCC (Multi-Version Concurrency Control), SELECT FOR UPDATE, deadlocks.

---

### Chapter 17: PostgreSQL Performance

Connection pooling: opening a new database connection is expensive, so you keep a pool of connections open and reuse them. PgBouncer is the standard connection pooler for PostgreSQL. Vacuuming: PostgreSQL keeps old row versions for MVCC, and VACUUM reclaims the space they occupy. Table partitioning splits a very large table into smaller physical pieces while the database presents them as one. Materialized views precompute expensive queries and store the result.

**Key topics:** connection pooling, VACUUM, table partitioning, materialized views, `pg_stat_statements`.

---

### Chapter 18: Building a Real App with PostgreSQL

A complete Go REST API backed by PostgreSQL: users, posts, and comments. Migrations are versioned SQL scripts that evolve the database schema over time — every schema change is tracked in code, just like application code. Handling errors from the database gracefully. Testing with a real PostgreSQL instance using testcontainers (a Go library that starts Docker containers during tests). Deploying to a server.

**Key topics:** migrations, REST API with database, testing with `testcontainers-go`, deployment.

---

## Volume 4: SQLite & MySQL

> Two databases you will encounter constantly in the real world. SQLite runs inside your app with no server required — every iOS app, Android app, and browser uses it. MySQL powers WordPress, GitHub, and much of the web.

### Chapter 19: SQLite Deep Dive

SQLite is the most deployed database in the world — by a very large margin. It runs inside your application process and stores everything in a single file. There is no server, no installation, no configuration. WAL mode (Write-Ahead Logging) is a SQLite option that makes concurrent reads and writes faster. Building a Go CLI tool that uses SQLite as its storage backend.

**Key topics:** SQLite file format, WAL mode, appropriate use cases, `modernc.org/sqlite` in Go.

---

### Chapter 20: SQLite Internals (Preview of VaultDB)

Opening the SQLite file in a hex editor and reading the raw bytes directly. A hex editor shows each byte as two hexadecimal digits — this is what data looks like at its lowest level. How pages are laid out. How the B-tree occupies those pages. How a row is encoded as a sequence of bytes. This chapter is a preview of what you will build in Volume 8 — you are seeing the destination before you start the journey.

**Key topics:** SQLite file format spec, page header, cell format, freelist, overflow pages.

---

### Chapter 21: MySQL Setup and Core Differences

MySQL prioritizes speed and broad compatibility. Setting up MySQL locally and connecting from Go with the `go-sql-driver/mysql` package. Key differences from PostgreSQL: AUTO_INCREMENT instead of SERIAL for auto-incrementing IDs, different string type defaults, a less strict SQL mode by default (which can silently truncate data unless you configure it). Knowing the differences prevents bugs when moving between databases.

**Key topics:** MySQL setup, `go-sql-driver/mysql`, SQL mode, ENUM vs CHECK, AUTO_INCREMENT.

---

### Chapter 22: MySQL in Production

Replication: one primary server accepts all writes and sends a stream of changes to replica servers, which stay in sync. Applications read from replicas (spreading load) and write to the primary. ProxySQL sits in front of MySQL and routes queries automatically — writes go to the primary, reads go to replicas. InnoDB is MySQL's main storage engine — it implements transactions, row-level locking, and the B-tree storage.

**Key topics:** MySQL replication, ProxySQL, InnoDB, binary log, choosing between MySQL and PostgreSQL.

---

## Volume 5: NoSQL

> Relational databases are not the only way to store data. NoSQL (Not Only SQL) databases trade some of the guarantees of relational databases for speed, flexibility, or scale. Understanding when to reach for NoSQL — and when not to — is one of the most important skills in backend engineering.

### Chapter 23: Why NoSQL?

The relational model forces data into rows and columns with a fixed schema. Sometimes your data does not fit neatly: a social network where every user has a different set of profile fields, or a product catalog where a laptop has different attributes than a shirt. The CAP theorem states that a distributed system can guarantee at most two of three properties: Consistency (every read sees the latest write), Availability (every request gets a response), and Partition tolerance (the system keeps working when network links fail). BASE (Basically Available, Soft state, Eventually consistent) is the alternative to ACID embraced by many NoSQL systems.

**Key topics:** schema-less data, the CAP theorem, BASE vs. ACID, overview of NoSQL categories.

---

### Chapter 24: Redis — The In-Memory Database

Redis stores everything in memory — RAM — making it blindingly fast compared to disk-backed databases. Redis data structures go beyond simple key-value: lists (ordered sequences), sets (unique items), sorted sets (items with a score for ranking), and hashes (like a dictionary inside a key). TTL (Time To Live) lets you set an expiry on any key so Redis automatically deletes it after a given duration. Connecting from Go with the `go-redis` library.

**Key topics:** Redis data structures, caching patterns, TTL, atomic operations, Pub/Sub, `go-redis`.

---

### Chapter 25: MongoDB — Document Databases

MongoDB stores documents — JSON-like objects — instead of rows. Each document can have different fields. There is no schema you must define upfront. BSON is the binary format MongoDB uses internally to store documents — it is like JSON but with additional types like dates and binary data. The aggregation pipeline is MongoDB's way of transforming and analyzing documents in a series of stages: filter, group, sort, project, and more.

**Key topics:** BSON documents, CRUD in MongoDB, aggregation pipeline, indexes, `mongo-driver`.

---

### Chapter 26: Cassandra — Wide-Column Stores

Cassandra was designed to never go down, even during a data center outage. It replicates data across many nodes organized in a ring — there is no single primary server, so there is no single point of failure. Eventual consistency means that a write to one node will eventually propagate to all others, but a read immediately after the write might not yet see it. The partition key determines which node stores a row. The clustering key determines the order of rows within a partition.

**Key topics:** ring topology, eventual consistency, partition key, clustering key, CQL, `gocql`.

---

### Chapter 27: Neo4j — Graph Databases

Some data is naturally a graph: social networks (people connected to people), maps (cities connected by roads), fraud detection (accounts connected by shared phone numbers). A graph database stores nodes (things) and edges (connections between things). Cypher is Neo4j's query language — you write patterns like `(alice)-[:FOLLOWS]->(bob)` to describe the relationships you are looking for. Finding friends-of-friends or the shortest path between two cities is trivial in Cypher and painful in SQL.

**Key topics:** graph model, Cypher query language, use cases for graph databases, `neo4j-go-driver`.

---

### Chapter 28: Choosing the Right Database

A decision framework for picking the right database for a given problem. For each situation, which database fits best and why. Polyglot persistence: most large systems use multiple databases — PostgreSQL for the main relational data, Redis for caching, Kafka for event streaming. The most common mistakes: picking a NoSQL database because it sounds modern, using Redis as a primary database without understanding its data loss risks, using MongoDB when your data is relational and you just need to learn SQL joins.

**Key topics:** decision framework, polyglot persistence, common mistakes, designing for the future.

---

## Volume 6: Distributed Databases

> What happens when one machine is not enough? When you have so much data, or so many users, that a single server cannot keep up? Distributed databases spread data across many machines. They introduce new problems — and fascinating solutions.

### Chapter 29: The Problem of Distribution

Sending data across a network introduces latency (delay) and the possibility of failure — packets can be lost, reordered, or duplicated. Clocks on different machines drift and are never exactly in sync — two machines might disagree on the current time by dozens of milliseconds. The Two Generals Problem is a classic thought experiment that proves some coordination problems are fundamentally unsolvable over an unreliable network. This chapter is about understanding why distributed systems are hard before seeing how real systems cope.

**Key topics:** network partitions, clock skew, the Two Generals Problem, why distributed systems are hard.

---

### Chapter 30: Replication

Replication means keeping identical copies of your data on multiple machines. If one machine dies, the others continue serving requests. Leader-follower replication: one machine (the leader) accepts all writes and streams changes to followers, which apply the same changes in order. Synchronous replication waits for the follower to confirm before the leader acknowledges the write — safe but slower. Asynchronous replication acknowledges the write immediately — faster but the follower might lag behind.

**Key topics:** replication, leader election, synchronous vs. asynchronous, replication lag, failover.

---

### Chapter 31: Sharding

Sharding splits your data across machines so each machine holds only a fraction of it. Instead of one server with 10 terabytes, you have 10 servers each with 1 terabyte. The shard key determines which machine a row lives on. Consistent hashing is an algorithm that makes adding and removing machines painless — when you add a new shard, only a fraction of the data needs to move, rather than reshuffling everything. The problems sharding creates: queries that need data from multiple shards require coordination, and distributed transactions (changing data on two shards atomically) are much harder.

**Key topics:** horizontal sharding, shard key, consistent hashing, cross-shard queries, resharding.

---

### Chapter 32: Consensus Algorithms

How do distributed machines agree on anything when messages can be lost and machines can crash? Raft is a consensus algorithm designed to be understandable — unlike its predecessor Paxos, which is notoriously difficult to explain. Raft works like a voting committee: one node is the leader, it proposes changes, and the change is committed when a majority (quorum) of nodes vote to accept it. If the leader dies, the remaining nodes elect a new one. Implementing a tiny Raft-like system in Go to make the algorithm concrete.

**Key topics:** Raft algorithm, leader election, log replication, quorum, etcd, CockroachDB.

---

## Volume 7: Modern Databases

> The database world is evolving fast. New systems are being built for cloud infrastructure, analytical workloads, vector search, and serverless architectures. This volume tours the frontier.

### Chapter 33: NewSQL — The Best of Both Worlds

NewSQL databases are distributed like NoSQL but fully SQL-compatible and ACID-compliant like traditional relational databases. CockroachDB and TiDB achieve this by using the Raft consensus algorithm to keep multiple replicas in sync while exposing a standard SQL interface. You can run SQL queries, have full transactions, and survive a node failure — all at once. Geo-partitioning lets you pin data to specific geographic regions to comply with data residency laws.

**Key topics:** NewSQL definition, CockroachDB architecture, TiDB, distributed transactions, geo-partitioning.

---

### Chapter 34: Analytical Databases — OLAP vs. OLTP

Online Transaction Processing (OLTP) is optimized for fast individual operations: insert one order, look up one user. Online Analytical Processing (OLAP) is optimized for slow, complex queries over huge datasets: total revenue by country by month for the past three years. OLAP databases use columnar storage — instead of storing all columns of a row together, they store all values of each column together, which makes aggregating one column across millions of rows very fast. DuckDB brings this power to an embeddable, serverless database you can run inside a Go program.

**Key topics:** OLTP vs. OLAP, columnar storage, DuckDB, ClickHouse, data warehouses.

---

### Chapter 35: Time-Series Databases

A time-series is a sequence of measurements taken over time: temperature every second, stock price every millisecond, server CPU usage every 10 seconds. Ordinary databases can store time-series data but are not optimized for the common operations: inserting millions of timestamped rows per second, computing a moving average over a sliding window of time, or automatically deleting data older than 90 days (retention policy). TimescaleDB is a PostgreSQL extension that adds these capabilities. InfluxDB is a standalone time-series database with its own query language.

**Key topics:** time-series data, TimescaleDB, InfluxDB, continuous aggregates, retention policies.

---

### Chapter 36: Vector Databases

A vector is a list of numbers that represents meaning — a sentence converted into a list of 1,536 numbers by an AI model, for example. Two sentences with similar meaning produce vectors that are numerically close to each other. Vector databases are optimized for finding the closest vectors to a query vector — the nearest-neighbor search. This powers semantic search (find documents that mean the same thing as this query, not just documents with the same words), image similarity, and AI recommendation systems. Pgvector is a PostgreSQL extension that adds vector search to standard PostgreSQL.

**Key topics:** embeddings, vector similarity, pgvector, approximate nearest neighbor search, use cases in AI.

---

### Chapter 37: Serverless and Edge Databases

Serverless databases charge per query and scale to zero when idle — you pay nothing when your app has no users. Edge databases replicate your data to servers physically close to your users around the world — a user in Tokyo reads from a Tokyo server rather than one in Virginia, reducing latency dramatically. Turso (based on libSQL, a fork of SQLite) and Neon (serverless PostgreSQL) are two of the most interesting new entries in this space. Trade-offs: no ops overhead, but different performance characteristics and sometimes different SQL compatibility.

**Key topics:** serverless databases, edge computing, Turso, Neon, connection model differences.

---

## Volume 8: Building VaultDB — Your Own Database

> This is the heart of Part 1. Everything you have learned so far was preparation for this volume. You are going to build VaultDB: a real, working, disk-backed, SQL-capable database engine in Go. One chapter at a time, one piece at a time, until it all works together.

### Chapter 38: VaultDB Architecture Overview

Before writing a single line, you need a map. The storage engine reads and writes pages to disk. The B-tree provides the index. The page cache keeps recently used pages in memory so you do not hit disk for every read. The SQL parser turns a string into a structured object. The query planner decides the best execution strategy. The execution engine carries out the plan. The transaction manager coordinates concurrent access. A diagram of how a SQL query flows through all these layers from text to bytes on disk and back.

**Key topics:** database architecture, component overview, data flow, what you will build and in what order.

---

### Chapter 39: The Page — VaultDB's Unit of Storage

Every database reads and writes data in fixed-size chunks called pages — typically 4,096 bytes, matching the size of a block on most hard drives and SSDs. Aligning to this size means every page read is one disk operation, never splitting across two blocks. Define the VaultDB page format: a header (page number, page type, free space pointer) followed by the data area. Write Go code to read and write pages to a file using `encoding/binary` to control the exact byte layout.

**Key topics:** page format, page header, Go `os.File` operations, binary encoding with `encoding/binary`.

**Mini Project:** A Go program that writes 100 pages to a file and reads them back in random order, verifying the content of each page.

---

### Chapter 40: The B-Tree — VaultDB's Index Heart

The B-tree is the data structure at the core of almost every relational database. It keeps keys sorted so that finding a key takes O(log n) time — a tree with a million keys requires at most 20 comparisons. Implement a B-tree in Go from scratch. Internal nodes store keys and pointers to child nodes. Leaf nodes store the actual data. When a node fills up, it splits in two and pushes the middle key up to the parent. When a node drops below half full, it borrows from a sibling or merges.

**Key topics:** B-tree node types (internal vs. leaf), insertion algorithm, search algorithm, deletion, node splitting.

**Mini Project:** A standalone B-tree in Go that can store and retrieve 10,000 integer keys, with a test that verifies every key is retrievable after all insertions.

---

### Chapter 41: The Heap File — Storing Rows

A heap file is an unordered collection of rows stored in pages — "heap" here means disorganized, not the programming heap. Rows are stored wherever there is free space. Implement VaultDB's heap file: insert a row (find a page with space, add the row, update the free space pointer), delete a row (mark it with a tombstone — a flag that says "this slot is empty"), read a row by its row ID (the page number and slot number combined). The tuple format encodes one row as a header (column count, null bitmap) followed by the column values.

**Key topics:** heap file, tuple format, slot array, free space management, row ID (page, slot).

---

### Chapter 42: The Catalog — VaultDB's Schema Store

The catalog is the database's memory of itself: which tables exist, what columns each table has, what data types those columns are, which indexes are defined. This metadata must itself be stored durably — if the catalog is lost, the database cannot interpret any of the data files. Implement the VaultDB catalog stored in a special system table that is bootstrapped on first startup and loaded from disk on every subsequent startup.

**Key topics:** system catalog, metadata storage, schema persistence, bootstrapping.

---

### Chapter 43: The SQL Parser

Parsing turns a string like `SELECT name FROM users WHERE age > 18` into a structured Go object your code can work with. A tokenizer (also called a lexer) splits the string into tokens: keywords like SELECT, identifiers like "users", operators like >, and literals like 18. A recursive-descent parser reads tokens in order and builds an Abstract Syntax Tree (AST) — a tree-shaped Go struct that represents the structure of the query. Build a parser that handles SELECT, INSERT, UPDATE, DELETE, CREATE TABLE, and DROP TABLE.

**Key topics:** tokenizer (lexer), recursive descent parsing, abstract syntax tree (AST), parser errors.

**Mini Project:** A parser that can parse 20 SQL statements and print the resulting AST for each one.

---

### Chapter 44: The Query Planner

Given a parsed SQL statement (an AST), the planner decides how to execute it. Should it use an index or scan the entire table? If the query filters by an indexed column, an index scan might read only 10 rows. Without the index, a sequential scan reads every row in the table. Should it apply the WHERE clause filter before or after performing a JOIN? Rule-based planning applies a fixed set of rules (always push filters down, always use an index when available). Implement a simple rule-based planner for VaultDB.

**Key topics:** logical plan, physical plan, predicate pushdown, index selection, plan nodes.

---

### Chapter 45: The Execution Engine

The execution engine takes the physical plan from the planner and actually runs it. VaultDB uses the iterator model (also called the volcano model): each plan node is an object with a `Next()` method that returns one row at a time. A SeqScan node reads rows from a heap file one by one. A Filter node wraps a SeqScan and skips rows that do not match the predicate. A Projection node wraps a Filter and returns only the requested columns. A NestedLoopJoin node takes two child iterators and returns the cross product of matching rows.

**Key topics:** volcano/iterator model, SeqScan, IndexScan, Filter, Projection, NestedLoopJoin.

---

### Chapter 46: The Write-Ahead Log (WAL)

If VaultDB crashes in the middle of writing a row, the page on disk might be half-written — corrupted. The Write-Ahead Log solves this: before changing any data page, write a log record describing the change to a separate log file. If the program crashes, the log file survives. On restart, VaultDB replays completed transactions from the log (redo) and removes incomplete ones (undo). The Log Sequence Number (LSN) is a monotonically increasing number that orders all log records. A checkpoint periodically flushes dirty pages to disk and records the LSN up to which recovery can skip.

**Key topics:** WAL format, log sequence number (LSN), redo logging, crash recovery, checkpoint.

---

### Chapter 47: Transactions in VaultDB

Implement `BEGIN`, `COMMIT`, and `ROLLBACK`. A transaction has an ID. When it begins, it records its start LSN. Every change it makes is logged with its transaction ID. When it commits, it writes a COMMIT record to the WAL. When it rolls back, it undoes its changes using the undo log. A lock manager grants shared locks (for reads) and exclusive locks (for writes) and blocks conflicting requests until the lock is released. Two-phase locking: acquire all locks before releasing any, which guarantees serializability.

**Key topics:** transaction ID, lock manager, two-phase locking, undo log, COMMIT and ROLLBACK.

---

### Chapter 48: VaultDB Complete — Wire It All Together

Every piece is built. Now connect them. Build the VaultDB TCP server: it listens for connections, reads SQL queries, passes them through the parser, planner, and execution engine, and writes the results back. Build the Go client library with a clean API: `Open`, `Exec`, `Query`, `Close`. Write an integration test suite that tests the full system end-to-end. The system you built from scratch now passes real SQL queries.

**Key topics:** TCP server for SQL, wire protocol, client library, integration testing, what to build next.

**Mini Project:** A command-line REPL (Read-Eval-Print Loop) for VaultDB — a program that accepts SQL typed at the terminal, executes it, and prints the results, just like the `psql` command-line client for PostgreSQL.

---

# PART 2 — ASYNC SYSTEMS

---

## Volume 9: Why Async?

> Async means "not happening at the same time." When a program sends a message and immediately moves on — without waiting for a reply — that is async. This volume builds the intuition for why async systems exist and what problems they solve.

### Chapter 49: The Problem with Waiting

When your program calls a function, it waits for the result before doing anything else. This is synchronous execution. Imagine a restaurant where every waiter stands frozen next to the customer's table until the food is ready. The kitchen would be fine but the dining room would be chaos — all the waiters are blocked and new customers cannot be seated. Programs hit the same wall at scale: a web server that calls a slow email service synchronously will block every request until that call finishes. Goroutines are Go's lightweight threads that help — but they are not always enough.

**Key topics:** synchronous vs. asynchronous, blocking I/O, the cost of waiting, goroutines as a partial solution.

---

### Chapter 50: Queues — The Core Idea

A queue is a line. First in, first out — the first item added is the first item removed. A post office drop box is a queue: you drop a letter in (producer), a postal worker picks it up later (consumer). The drop box decouples the sender from the receiver — the sender does not need to wait for the postal worker to be present, and the postal worker does not need to wait for customers to arrive. This simple idea — a buffer between producer and consumer — is the foundation of every async system. Backpressure is what happens when the queue fills up faster than it is emptied, and how well-designed systems handle it.

**Key topics:** queue abstraction, producer-consumer pattern, decoupling, backpressure, queue depth.

---

### Chapter 51: When to Go Async

Not everything should be async. Async adds complexity: you must think about what happens when the consumer is slow, when messages are delivered more than once, or when the consumer crashes mid-processing. Async is worth it when: the producer generates work faster than the consumer can handle it, you need to retry failed operations without the caller waiting, or one event must be delivered to many different systems independently. At-least-once delivery means a message might be delivered more than once — your consumer must be idempotent (safe to run multiple times with the same input) to handle this.

**Key topics:** async trade-offs, idempotency, at-least-once vs. exactly-once delivery, use cases.

---

## Volume 10: Real Async Systems

> Now you use the real tools. Kafka, RabbitMQ, Redis Streams, NATS — each one is a different answer to the async question, with different trade-offs. You will use all of them from Go.

### Chapter 52: Apache Kafka — The Foundation

Kafka is a distributed append-only log. Messages are written to the end of the log and never modified or deleted during their retention period. A consumer tracks its position in the log using an offset — a number that says "I have processed everything up to message 500,000." Crucially, the consumer owns its offset: it can rewind to any past offset and replay messages. This makes Kafka unlike a traditional queue where a consumed message is gone forever. Topics are named logs. Partitions split a topic's log across multiple servers for parallelism.

**Key topics:** Kafka topics, partitions, offsets, consumer groups, producers, brokers, using `confluent-kafka-go`.

---

### Chapter 53: Kafka Patterns in Go

Event sourcing stores the history of state changes (events) rather than the current state — to reconstruct current state, replay all events. CQRS (Command Query Responsibility Segregation) separates write operations (commands) from read operations (queries) — the write side publishes events to Kafka, multiple read sides consume those events and maintain their own optimized read models. The saga pattern coordinates multi-step distributed transactions through a sequence of events and compensating actions when a step fails. Writing Go producers and consumer groups. Handling errors, retries, and dead-letter queues (a special topic for messages that failed too many times).

**Key topics:** event sourcing, CQRS, saga, consumer groups in Go, error handling, DLQ.

---

### Chapter 54: RabbitMQ — Message Queues

RabbitMQ uses a different model from Kafka. Messages go to a queue and are deleted once a consumer acknowledges receiving them — there is no offset, no replay. This simpler model is ideal for task queues: one job goes in, one worker picks it up, processes it, and acknowledges it. Exchanges are the routing layer: a publisher sends to an exchange, not directly to a queue. The exchange routes messages to queues based on configurable rules. A direct exchange routes by exact routing key. A fanout exchange broadcasts to all bound queues. Using the `amqp091-go` library from Go.

**Key topics:** RabbitMQ model, exchanges, bindings, routing keys, acknowledgments, `amqp091-go`.

---

### Chapter 55: Redis Streams

Redis Streams are a Kafka-like append-only log built into Redis. They are lighter weight than Kafka — lower throughput but simpler to operate when you are already running Redis. `XADD` appends a message to a stream. `XREAD` reads messages forward from a given ID. `XACK` acknowledges that a message has been processed. Consumer groups in Redis Streams work like Kafka consumer groups: multiple consumers in a group split the messages, and each message goes to exactly one consumer in the group.

**Key topics:** Redis Streams commands, consumer groups, pending messages, ACK, `go-redis` streams API.

---

### Chapter 56: NATS — Fast Lightweight Messaging

NATS is a tiny, extremely fast messaging system written in Go. Its core is a simple publish-subscribe model with no persistence — if no subscriber is listening when a message is published, the message is gone. JetStream adds persistence and replay capabilities, making it comparable to Kafka for durable messaging. NATS is popular in cloud-native and IoT systems because its server binary is small and its protocol is simple. The `nats.go` client library from Go is clean and idiomatic.

**Key topics:** NATS pub/sub, request-reply, JetStream, `nats.go`, when to choose NATS over Kafka.

---

### Chapter 57: Async Patterns and Pitfalls

The outbox pattern solves a fundamental problem: how do you atomically update your database and publish an event to a message broker? If you do the database write and then crash before publishing the event, the event is lost. The outbox pattern writes the event to an "outbox" table in the same database transaction as the main write, then a separate process reads from the outbox and publishes to the broker. Poison pills are messages that always cause the consumer to crash — a bad message trapped in a retry loop that crashes your service indefinitely. Distributed tracing tracks a request across many async hops.

**Key topics:** idempotency, the outbox pattern, transactional outbox, poison pills, DLQ, distributed tracing.

---

### Chapter 58: Testing Async Systems

Async code is difficult to test because the producer and consumer run independently — you cannot simply call a function and check its return value. Polling with timeout: after publishing a message, repeatedly check whether the expected side effect has appeared, up to a maximum wait time. Testcontainers starts a real Kafka or RabbitMQ server inside a Docker container as part of your test setup, so tests run against the real thing. In-memory fakes replace the real broker with a simple in-process implementation that delivers messages synchronously during tests.

**Key topics:** testing async code, testcontainers for Kafka, event-driven test patterns, flaky test prevention.

---

## Volume 11: Building StreamFlow — Your Own Async System

> Now you build it. StreamFlow is a persistent, log-structured message broker in Go, inspired by Kafka. By the end of this volume you will have a working message broker that your programs can publish to and subscribe from.

### Chapter 59: StreamFlow Architecture Overview

The same planning approach as VaultDB's Chapter 38. The commit log is the storage core — an append-only file of message records. The offset index provides fast lookup of any message by offset. The topic manager maintains the set of topics and their associated log files. The TCP server accepts connections from clients. The client library wraps the protocol in a clean Go API. A diagram. A plan. The order in which each component will be built.

**Key topics:** StreamFlow architecture, component overview, design decisions, what you will build.

---

### Chapter 60: The Commit Log — StreamFlow's Storage Heart

StreamFlow stores messages in an append-only log file. Sequential writes are the fastest possible disk operation — there is no seeking, no fragmentation. Each record in the log has a header (the offset number, the timestamp, the body length in bytes) followed by the body bytes. Implement the log: append a record and return its offset, read a record by its exact byte position, scan forward from a given offset reading records one by one. Log rotation creates a new segment file when the current one reaches a maximum size.

**Key topics:** append-only log, record format, sequential writes, log segment files, log rotation.

**Mini Project:** A log file that can store and retrieve 1,000 messages by offset, with a test that verifies each message body is retrieved correctly.

---

### Chapter 61: The Offset Index

Reading message number 500,000 directly requires knowing its byte position in the log file — you cannot calculate it because messages have variable lengths. Build a sparse index: every 1,000 messages, record a pair of (offset, byte position). To find message 500,000, binary search the index for the largest offset less than or equal to 500,000, then scan forward from that byte position in the log file. Sparse means the index does not have an entry for every message — it is much smaller than a dense index and fast enough in practice.

**Key topics:** sparse index, binary search, index file format, trade-offs vs. dense index.

---

### Chapter 62: Topics and Partitions

A topic is a named stream of messages — like a named channel. A partition is a subdivision of a topic: one topic can have three partitions, and each partition is an independent log file. Partitions are the unit of parallelism — a consumer group with three workers can assign one partition to each worker, processing three messages simultaneously. Implement the StreamFlow topic manager: create a topic with a specified number of partitions, route a publish request to the correct partition log based on the message key, list all topics and their partition counts.

**Key topics:** topic abstraction, partitioning strategy, partition assignment, topic metadata.

---

### Chapter 63: The StreamFlow Server

Build the TCP server that clients connect to. Design a binary wire protocol: a request starts with a one-byte command code, followed by a length-prefixed payload. The PUBLISH command includes the topic name and message body. The FETCH command includes the topic name, partition number, and starting offset. SUBSCRIBE registers the connection for push delivery of new messages. Handle each client connection in its own goroutine. Read requests in a loop, dispatch to the appropriate handler, write the response.

**Key topics:** TCP server in Go, binary wire protocol, goroutine-per-connection, connection lifecycle.

---

### Chapter 64: The StreamFlow Client Library

Build the Go client library that application code uses. `Connect` dials the broker's TCP address and returns a client. `Publish` sends a message to a topic. `Subscribe` registers interest in a topic from a given offset and returns a channel — application code ranges over this channel to receive messages as they arrive. `Close` gracefully shuts down the connection. Reconnection logic: if the connection drops, the client automatically reconnects and resumes from its last acknowledged offset.

**Key topics:** client library design, TCP client in Go, channel-based API, reconnection logic.

---

### Chapter 65: Consumer Groups

A consumer group is a named set of consumers that cooperate to process a topic together. Each partition is assigned to exactly one consumer in the group at a time — so messages from that partition are processed by exactly one consumer. When a new consumer joins the group, the broker rebalances: it reassigns partitions so the work is evenly distributed. When a consumer leaves (or crashes), its partitions are reassigned to the remaining consumers. Implement group state tracking, the rebalance protocol, and offset commits (consumers periodically report their progress to the broker so they can resume after a restart).

**Key topics:** consumer group protocol, partition ownership, rebalancing, offset commits.

---

### Chapter 66: Durability and Crash Recovery

The commit log is inherently durable — it is an append-only file and completed writes survive crashes. The in-memory state — which topics exist, which consumer groups are registered, which consumer owns which partition — must also survive crashes. Implement a metadata log: every time the topic manager makes a change, it writes a metadata log record describing the change. On startup, StreamFlow replays the metadata log to reconstruct the in-memory state. `fsync` forces the operating system to flush buffered writes to the physical disk, guaranteeing durability.

**Key topics:** metadata log, recovery protocol, fsync, durability guarantees.

---

### Chapter 67: StreamFlow Complete — Wire It All Together and Test

All pieces are built. Integrate them. Run a full end-to-end test: a producer publishes 100,000 messages to a topic with three partitions, a consumer group of three workers consumes them all, the StreamFlow server is killed and restarted mid-run, the consumers detect the disconnect and reconnect, and consumption resumes from where it left off. Verify that no messages are lost or processed twice. Run throughput benchmarks.

**Key topics:** integration testing, crash testing, throughput benchmarks, what to build next.

**Mini Project:** A StreamFlow-powered chat room — multiple Go programs each connect to StreamFlow, publish messages to a "messages" topic, and subscribe to receive all messages, effectively simulating a multi-user group chat.

---

# PART 3 — INTEGRATION

---

## Volume 12: Putting It All Together

> You have a database. You have a message broker. Now you build the real thing: a production analytics service that uses both. This volume is about architecture, deployment, and the engineering judgment that comes from having built the foundations yourself.

### Chapter 68: System Design — The Analytics Service

The system you are building tracks events from a web application — a user viewed a product, added it to their cart, completed a purchase — and answers questions like "what were the top 10 products by revenue this week?" Design each component: the ingestion server (receives events over HTTP), the processor (reads events from StreamFlow and writes to VaultDB), the query API (reads from VaultDB and returns aggregated stats). Draw the data flow. Identify where failures can occur and how the system handles them.

**Key topics:** system design process, event-driven architecture, component boundaries, data flow.

---

### Chapter 69: The Ingestion Layer

Build the HTTP event ingestion server in Go using the standard `net/http` package. Accept JSON event payloads, validate the required fields, and publish each event to StreamFlow. Handle high throughput: instead of publishing one event per HTTP request, collect events into a batch and publish the batch in a single StreamFlow write. Backpressure handling: if StreamFlow is slow, the HTTP server should signal this to callers with a 429 (Too Many Requests) response rather than queuing indefinitely in memory. Load test with `k6` to verify the server handles 1,000 requests per second.

**Key topics:** HTTP server in Go, JSON validation, batching, backpressure handling, load testing.

---

### Chapter 70: The Processing Layer

Build the StreamFlow consumer that reads events and writes them to VaultDB. Schema evolution: the event format may change over time — new fields are added, old fields are renamed. The processor must handle old and new event formats gracefully. Idempotent writes: if the same event is delivered twice (at-least-once delivery), writing it twice must not corrupt the analytics data — use the event ID as a deduplication key. Monitor consumer lag: the difference between the latest published offset and the consumer's current offset indicates how far behind the processor is.

**Key topics:** stream processing, idempotent writes, schema evolution, consumer lag monitoring, DLQ.

---

### Chapter 71: The Query Layer

Build the HTTP query API that reads from VaultDB and returns analytics results. A SQL query builder constructs SQL strings programmatically rather than through string concatenation — this prevents SQL injection and makes complex queries readable. Redis caching: wrap expensive VaultDB queries in a Redis cache with a short TTL — the first request computes the result and caches it, subsequent requests within the TTL window read from Redis instead. Pagination returns large result sets in fixed-size pages rather than all at once, protecting the server from memory exhaustion.

**Key topics:** query API design, SQL query builder, Redis caching, pagination, API testing.

---

### Chapter 72: Deployment and Operations

Package each component as a Docker image. Write a `docker-compose.yml` that starts all five pieces — VaultDB, StreamFlow, the ingestion server, the processor, and the query API — with a single command. Health checks allow Docker to detect and restart unhealthy containers automatically. Structured logging with Go's `slog` package emits JSON log lines that can be searched and filtered by log management systems. Prometheus metrics expose internal counters (events ingested per second, consumer lag, query latency) that can be visualized in Grafana. Run the full system and watch it handle 10,000 events per second.

**Key topics:** Docker, docker-compose, health checks, structured logging with `slog`, Prometheus metrics, running at scale.

---

---

# PART 4 — VECTOR DATABASES: DEEP DIVE

---

> The world changed when large language models arrived. Every AI application — RAG pipelines, semantic search, recommendation systems, multimodal search — needs a vector database. Part 4 goes deep on Qdrant (one of the fastest-growing open-source databases in history), then builds **NebulaDB**: a complete, working Qdrant-inspired vector database from scratch in Go. By the end, you will understand exactly what Qdrant does internally.

---

## Volume 13: Qdrant — Understanding a Production Vector Database

### Chapter 67: Qdrant Deep Dive — Architecture and Internals

Why Qdrant became the default choice for AI applications. Core data model: Points (id + vector + payload), Collections, and distance metrics. The Segment architecture — growing vs. sealed segments, parallel search, memory-mapped storage. Qdrant's killer feature: HNSW with payload filtering. How pre-filtering, filtered HNSW traversal, and post-filtering work. Auto-strategy selection. Quantization: scalar (int8), product quantization, and binary quantization. Named vectors and sparse + dense hybrid search. WAL for durability. Distributed mode: sharding with consistent hashing, replication with configurable consistency. Snapshots and backups.

**Key topics:** Segments, HNSW with payload filtering, quantization (SQ8/PQ/binary), named vectors, sparse vectors, WAL, sharding, replication, snapshots.

---

### Chapter 68: Advanced Qdrant in Go

Everything beyond the basics. The modern high-level Go client. Named vectors for multimodal search (text + image embeddings per point). Creating and using payload indexes (keyword, numeric, bool, geo, text). Complex filter conditions: `must`, `must_not`, `should`, `match_any`, `range`, geo radius. Sparse + dense hybrid search with Reciprocal Rank Fusion. Scalar quantization with oversampling for rescoring. Bulk ingestion: raising indexing thresholds, batch upserts, restoring quality settings. Collection management: info, snapshots, scroll. **Project:** AI-powered product search — REST API with semantic vector search, category filtering, and price range filtering.

**Key topics:** named vectors, payload indexes, hybrid search, RRF fusion, quantization oversampling, bulk ingestion, collection management.

---

## Volume 14: NebulaDB — Building a Vector Database from Scratch

> NebulaDB is a Qdrant-inspired vector database written in Go. It has Collections, Points (id + vector + payload), a hand-built HNSW index, payload filtering, payload indexes, WAL for crash safety, snapshots, and a REST HTTP API. Every component is built from first principles.

### Chapter 69: NebulaDB Introduction — Architecture and Project Setup

What NebulaDB includes vs. what Qdrant adds. The full architecture diagram: HTTP API → Collection Manager → (HNSW Index, VectorStore, PayloadStore, PayloadIndexes, WAL). Data flow for upsert and search. Core types: `Point`, `Filter`, `Condition`, `SearchRequest`. The `CollectionManager`: create, get, delete, load existing. `Config` with HNSW parameters. `main.go` with graceful shutdown.

**Key topics:** NebulaDB architecture, CollectionManager, core types, project structure.

---

### Chapter 70: NebulaDB — HNSW Index from Scratch

HNSW intuition: layered graph, highways vs. neighborhood. Distance functions: cosine, Euclidean, dot product. The Node structure with per-layer adjacency lists. Insertion: random level assignment, greedy search for entry point, bidirectional connection, pruning excess neighbors. Search: greedy upper-layer navigation, beam search at layer 0 with `ef` candidates using a max-heap + min-heap pair. Filtered search: `filterFn` skips nodes from results while still using them for graph navigation. Persistence with `gob.Encoder`.

**Key topics:** HNSW layer structure, insertion algorithm, beam search, filtered traversal, gob persistence, benchmarking.

---

### Chapter 71: NebulaDB — Payload Storage and Filtering Engine

PayloadStore with BoltDB: `Set`, `Get`, `Delete`, `GetMany`, big-endian ID encoding for lexicographic sort. The Filter Engine: evaluating `Must`, `MustNot`, `Should` conditions. Match (exact, `match_any`), range conditions, array field matching. FieldIndex: in-memory keyword HashMap and sorted numeric slice. PayloadIndexManager: register indexes, `Index` all fields on insert, `FilterIDs` for intersection pre-filtering. Auto-strategy selection: brute force (<1% selectivity), filtered HNSW (1-50%), post-filter HNSW (>50%). Wiring everything together in `Collection`.

**Key topics:** BoltDB payload storage, filter evaluation, field indexes, selectivity-based strategy selection.

---

### Chapter 72: NebulaDB — WAL, Vector Storage, and Crash Recovery

WAL format: magic + optype + length + CRC32. `WriteUpsert`, `WriteDelete`, `Replay` with corruption detection. VectorStore: flat binary file with in-memory `id → offset` map, fixed-record layout, `Set` and `Get`. Snapshots: `tar.gz` archives of the collection directory with `RestoreSnapshot`. Crash recovery: load HNSW snapshot (last checkpoint), replay WAL tail to recover missed writes. The complete write order: WAL (fsynced) → VectorStore → PayloadStore → HNSW. Checkpoint: save HNSW snapshot + truncate WAL.

**Key topics:** WAL with CRC32, flat binary VectorStore, tar.gz snapshots, crash recovery lifecycle.

---

### Chapter 73: NebulaDB — HTTP API

REST API design: collections, points, search, and index endpoints. Go 1.22 path parameters (`{name}`, `{id}`). Consistent JSON response envelope. `withLogging` middleware. Collection endpoints: create, list, get, delete, snapshot. Points endpoints: upsert (batch), get, delete, scroll. Search endpoint with filter. `createIndex` endpoint. Integration tests with `httptest.NewRecorder()` — no network needed, fully isolated.

**Key topics:** Go 1.22 routing, REST endpoints, JSON envelope, middleware, httptest integration tests.

---

### Chapter 74: NebulaDB — Performance and Benchmarking

Benchmarking methodology: QPS, p99 latency, recall@10, indexing throughput. Writing a repeatable Go benchmark. Profiling with pprof: finding hotspots (distance function, map lookups, heap operations). SIMD-friendly distance computations: loop unrolling, pre-normalize + dot product. Concurrent segment search pattern. HNSW parameter tuning: M, ef_construction, ef_search — recall vs. latency trade-offs. Bulk ingestion optimization: skip incremental HNSW updates, rebuild at end. NebulaDB vs. Qdrant benchmark results and why Qdrant is 6-8x faster (Rust + SIMD + compact data structures).

**Key topics:** benchmarking, pprof profiling, distance optimization, HNSW tuning, bulk ingestion, NebulaDB vs. Qdrant comparison.

---

### Chapter 75: NebulaDB — Major Project: Complete Vector Database

Assembling all components into a fully working NebulaDB. **Project: Semantic Code Search Engine** — parse a Go repository with `go/ast`, extract function signatures + docstrings, embed with OpenAI, store in NebulaDB, search by natural language intent ("find functions that parse HTTP requests"). Deploying with Docker and docker-compose. What you'd add to make NebulaDB production-grade: quantization, distributed mode, named vectors, gRPC, mmap'd HNSW, sparse vectors. When to use pgvector vs. Qdrant vs. NebulaDB.

**Key topics:** complete system assembly, semantic code search, Docker deployment, production gap analysis, vector database selection guide.

---

## Summary — What You Have Built

By the time you reach the end of Chapter 75, here is what you have done:

- **Understood databases** from first principles — not just how to use them, but how they work inside. You can look at PostgreSQL, MySQL, SQLite, Redis, MongoDB, Cassandra, or Kafka and explain exactly what they are doing and why.

- **Built VaultDB** — a disk-backed, SQL-capable, ACID-compliant database engine in Go. It has a B-tree index, a heap file for rows, a SQL parser, a query planner, a volcano-model execution engine, a write-ahead log, and a lock-based transaction manager. It is roughly 5,000-8,000 lines of Go that you wrote and understand completely.

- **Built StreamFlow** — a persistent, log-structured message broker in Go. It has an append-only commit log, a sparse offset index, topics, partitions, consumer groups, a binary TCP wire protocol, and crash recovery via a metadata log. It is roughly 3,000-5,000 lines of Go that you wrote and understand completely.

- **Built a production analytics service** that wires both together: an HTTP ingestion API, a stream processor, a query API with Redis caching, Docker packaging with health checks, and Prometheus metrics. It handles real throughput under load.

- **Understood Qdrant deeply** — not just how to use its Go client, but its entire internals: segments, filtered HNSW traversal, payload indexing strategy selection, quantization, WAL, distributed sharding, and snapshots.

- **Built NebulaDB** — a Qdrant-inspired vector database from scratch in Go. It has a hand-built HNSW index with filtered traversal, BoltDB payload storage, field indexes for pre-filtering, a WAL with CRC32 integrity, vector storage, tar.gz snapshots, crash recovery, and a REST HTTP API. Everything Qdrant does architecturally, NebulaDB does in readable Go.

- **Developed engineering judgment** — the ability to look at any new system and ask the right questions: How does it store data? How does it handle crashes? How does it scale? How does it ensure consistency? These questions will serve you for the rest of your career.

---

## A Note on the Journey

This is not a course you finish and put on a shelf. It is a course that changes how you see software. Every database you use after this, you will understand differently. You will know that PostgreSQL's legendary reliability comes from its WAL, its MVCC, and its decades of careful engineering. You will know that Kafka's power comes from the append-only log — the simplest data structure in computing, applied in the most sophisticated way. You will know this because you built versions of both things yourself, from nothing, line by line.

That knowledge does not go away. It is yours.

Now turn to Chapter 1 and let us begin.

---

*Databases & Asynchronous Systems in Go — From Zero to Building Your Own Database, Async System, and Production Service in Go*
