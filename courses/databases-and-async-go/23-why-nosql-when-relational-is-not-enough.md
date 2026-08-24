# Chapter 23: Why NoSQL? When Relational Databases Are Not Enough

Imagine you're building a social network where every user has a different set of fields. Some have a bio, some have a website, some have 3 phone numbers, some have none. In a SQL table, you'd need columns for every possible field — mostly NULL. NoSQL was invented to solve this kind of problem.

## Table of Contents

1. The Limitations That Created NoSQL
2. The Five NoSQL Families
3. Document Databases
4. Key-Value Databases
5. Column-Family Databases
6. Graph Databases
7. The CAP Theorem in Practice
8. The Decision Guide
9. Exercises

---

## 1. The Limitations That Created NoSQL

SQL databases are brilliant for structured data with known relationships. They struggle with:

**Schema rigidity:** Every row must fit the same columns. Adding a new field means `ALTER TABLE` on potentially billions of rows.

**Horizontal scaling:** SQL databases (PostgreSQL, MySQL) are primarily designed for one server. Splitting a relational database across 100 servers is complex and expensive.

**Impedance mismatch:** Application objects are rich nested structures (a user with multiple addresses, orders, comments). SQL forces you to split this into many tables and JOIN them back together constantly.

**Specific access patterns:** Some data is fundamentally not relational — social graphs, time-series sensor data, search indexes, leaderboards.

In 2004-2008, companies like Google, Amazon, LinkedIn, and Facebook hit these walls. Their solutions became open-source and changed the database landscape.

---

## 2. The Five NoSQL Families

| Family | Examples | Data Model | Best For |
|--------|----------|------------|----------|
| Document | MongoDB, CouchDB | JSON/BSON documents | Flexible schemas, rich objects |
| Key-Value | Redis, DynamoDB | Key → Value | Caching, sessions, simple lookups |
| Column-Family | Cassandra, HBase | Wide rows of columns | Time-series, IoT, high-write |
| Graph | Neo4j, Amazon Neptune | Nodes and edges | Social networks, fraud detection |
| Search Engine | Elasticsearch, Meilisearch | Inverted index | Full-text search |

---

## 3. Document Databases

Document databases store data as self-contained documents (usually JSON). Each document can have a different structure.

```json
// User document in MongoDB
{
  "_id": "user_123",
  "name": "Alice",
  "email": "alice@example.com",
  "addresses": [
    {"type": "home", "city": "New York", "zip": "10001"},
    {"type": "work", "city": "Brooklyn", "zip": "11201"}
  ],
  "preferences": {
    "notifications": true,
    "theme": "dark"
  }
}
```

No schema means: add any field to any document without migrations. Perfect for:
- User profiles with varying data
- Product catalogs with different attributes per category
- Content management systems
- Rapidly changing schemas during development

**When NOT to use:** Complex multi-document transactions, highly relational data, complex reporting.

---

## 4. Key-Value Databases

A key-value database is the simplest database: a giant hash map. Give it a key, get back a value.

```
SET user:123:session "jwt_token_here"  TTL=3600
GET user:123:session → "jwt_token_here"
DEL user:123:session
```

Redis is the most popular key-value database. Its speed comes from storing everything in RAM. Use cases:
- Caching database query results
- Session storage
- Rate limiting counters
- Pub/Sub messaging
- Leaderboards (sorted sets)

---

## 5. Column-Family Databases

Imagine a spreadsheet where each row can have different columns, and there can be billions of rows efficiently distributed across servers. That's a column-family database.

```
Row key: "alice@example.com"
Columns:
  profile:name       = "Alice"
  profile:age        = 30
  events:2024-01-01  = "login"
  events:2024-01-02  = "purchase"
  events:2024-01-05  = "login"
```

Cassandra is the most popular. It's designed for:
- Write everywhere, no single point of failure
- Billions of rows, petabytes of data
- IoT sensor readings, application logs, activity streams
- Time-series data (measurements indexed by timestamp)

**When NOT to use:** Ad-hoc queries, anything requiring JOINs, frequent schema changes.

---

## 6. Graph Databases

Graph databases store data as **nodes** (things) and **edges** (relationships between things). Relationships are first-class citizens — not foreign keys.

```
Node: Alice (type: Person)
Node: Bob (type: Person)
Node: Acme Corp (type: Company)

Edge: Alice -[FRIENDS_WITH]-> Bob
Edge: Alice -[WORKS_AT]-> Acme Corp
Edge: Bob -[WORKS_AT]-> Acme Corp
```

Neo4j's query language (Cypher):
```cypher
// Find friends-of-friends
MATCH (alice:Person {name: "Alice"})-[:FRIENDS_WITH]->(friend)-[:FRIENDS_WITH]->(fof)
WHERE NOT (alice)-[:FRIENDS_WITH]->(fof)
RETURN fof.name
```

This query would require 3 JOINs in SQL and is slow at scale. In a graph database, it traverses edges directly — fast regardless of dataset size.

Use cases: social graphs, fraud detection, recommendation engines, knowledge graphs.

---

## 7. The CAP Theorem in Practice

We covered CAP in Chapter 06. Here's how real databases choose:

| Database | CAP Choice | What this means |
|----------|------------|-----------------|
| PostgreSQL | CA (single node) → CP (distributed) | Strongly consistent; may be unavailable during partition |
| MongoDB | CP (default) or AP (tunable) | Configurable consistency |
| Redis | AP (single node), CP (cluster) | May serve stale data from slave |
| Cassandra | AP | Always available; may serve stale data |
| Elasticsearch | AP | Eventually consistent |
| Neo4j | CA (single) / CP (cluster) | Consistent writes; complex partitioning |

**The key insight:** There's no "best" choice. It depends on what your application needs. Banking? You need CP — stale data means financial errors. Social feed? AP is fine — showing a slightly old post count doesn't matter.

---

## 8. The Decision Guide

Answer these questions to choose your database:

**Q1: Is your data highly relational?** (Many foreign key relationships, complex JOINs needed)
→ Yes: PostgreSQL or MySQL
→ No: Continue

**Q2: Do you need flexible/varying schemas?**
→ Yes: MongoDB
→ No: Continue

**Q3: Is this primarily for caching, sessions, or pub/sub?**
→ Yes: Redis
→ No: Continue

**Q4: Is this very high-write workload (millions of writes/sec) across many servers?**
→ Yes: Cassandra
→ No: Continue

**Q5: Is this primarily for analytics on billions of rows?**
→ Yes: ClickHouse
→ No: Continue

**Q6: Is the data primarily a graph (many interconnected nodes)?**
→ Yes: Neo4j or Amazon Neptune
→ No: Start with PostgreSQL — it handles most cases

**The default answer is PostgreSQL.** It can handle most workloads including JSON documents (JSONB), full-text search, and time-series data. Only reach for specialized databases when you've hit a specific limit.

---

## Summary

- NoSQL emerged to solve three problems: schema rigidity, horizontal scaling, and impedance mismatch.
- Five families: document (MongoDB), key-value (Redis), column-family (Cassandra), graph (Neo4j), search (Elasticsearch).
- Document databases shine for flexible schemas and hierarchical data.
- Column-family databases shine for massive write-heavy distributed workloads.
- The CAP theorem forces every distributed database to choose between consistency and availability.
- **Default to PostgreSQL** and switch to a specialized database when you hit a specific scaling or modeling wall.

### Exercises

**Easy:** Draw a schema for a blog using SQL (tables, relationships). Now draw the same data as MongoDB documents. Which approach requires more JOINs to display a post with author info and comments?

**Medium:** Look up the DynamoDB "single-table design" pattern. Why does DynamoDB require denormalization that would be considered bad practice in SQL?

**Hard:** Research the Google Bigtable paper (2006) and Amazon Dynamo paper (2007). These two papers inspired the entire NoSQL movement. Summarize the key design decisions of each and which NoSQL databases they inspired.
