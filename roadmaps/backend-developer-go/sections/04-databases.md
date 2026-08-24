---
title: Databases & Storage
---
Almost every backend service exists to read and write data reliably. This section covers how databases actually work under the hood, how to talk to them from Go, and when to reach for something other than a relational database.

### How Databases Work
Before learning any specific database, understand the shared ideas underneath all of them: pages, indexes, and trees on disk.

**Resources:**
- [How Databases Work Inside: Pages, Trees, and Indexes](course:databases-and-async-go#05-how-databases-work-inside-pages-trees-and-indexes)
- [Building Blocks of Every Database](course:databases-and-async-go#06-building-blocks-of-every-database)

### SQL Fundamentals
Tables, joins, and aggregations — the language nearly every backend eventually needs to speak fluently, whatever ORM sits on top of it.

**Resources:**
- [Introduction to SQL](course:databases-and-async-go#07-introduction-to-sql-talking-to-databases-in-plain-english)
- [Creating Tables, Inserting, Querying](course:databases-and-async-go#08-creating-tables-inserting-querying)
- [Joins, Aggregations, Complex Queries](course:databases-and-async-go#09-joins-aggregations-complex-queries)

### PostgreSQL with Go
PostgreSQL is the default choice for most new backend services. Learn to connect to it, write queries, and use it properly from Go code.

**Resources:**
- [Building with PostgreSQL in Go](course:databases-and-async-go#15-building-with-postgresql-in-go)
- [PostgreSQL and sqlx](course:go-programming#73-postgresql-and-sqlx)

### Transactions, Indexes & Schema Design
A database that loses data under concurrent writes, or falls over once a table hits a million rows, is worse than no database at all. This is how you avoid both.

**Resources:**
- [ACID Transactions: Never Lose Data](course:databases-and-async-go#11-acid-transactions-never-lose-data)
- [Indexes: Making Queries Fast](course:databases-and-async-go#10-indexes-making-queries-fast)
- [Schema Design and Normalization](course:databases-and-async-go#12-schema-design-and-normalization)

### NoSQL: Redis & MongoDB
> optional

Not every problem fits a relational table. Redis for caching and fast lookups, MongoDB for flexible document storage — know when each one actually earns its place in an architecture.

**Resources:**
- [Redis: The Speed Demon](course:databases-and-async-go#26-redis-the-speed-demon)
- [MongoDB: Documents All the Way Down](course:databases-and-async-go#24-mongodb-documents-all-the-way-down)

### Practice: Build a Key-Value Store
> branches-from: How Databases Work

Build a small key-value store yourself — the simplest possible database — to see indexing and storage tradeoffs from the inside instead of just reading about them.

**Resources:**
- [Key-Value Store project](project:key-value-store)
