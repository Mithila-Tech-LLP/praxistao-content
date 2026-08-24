# Chapter 35: SurrealDB — The Multi-Model Database

SurrealDB is a new kind of database that combines relational tables, document storage, graph relationships, and real-time queries — all in one system with one query language. It's like getting PostgreSQL, MongoDB, Neo4j, and Firebase in a single database.

## Table of Contents

1. What Makes SurrealDB Different
2. SurrealQL — The Query Language
3. Tables, Records, and Documents
4. Graph Relationships
5. Real-Time Subscriptions
6. Authentication Built In
7. Exercises

---

## 1. What Makes SurrealDB Different

Traditional approach: you need separate databases for different data patterns.

```
User data → PostgreSQL (relational)
Product catalog → MongoDB (documents)
Social connections → Neo4j (graphs)
Real-time updates → Firebase (subscriptions)
Auth → Separate auth service
```

SurrealDB approach: one database handles all of these.

```
Everything → SurrealDB
```

**Key features:**

| Feature | What It Means |
|---------|--------------|
| Multi-model | Tables, documents, and graphs in one query |
| SurrealQL | SQL-like query language with graph traversal |
| Schema-optional | Schemafull, schemaless, or mixed |
| Real-time | Live queries push updates to clients |
| Built-in auth | Users, scopes, and permissions without extra code |
| Embedded or server | Runs in-process (like SQLite) or as a server |

**When to use SurrealDB:**
- Building a new product where you're not sure if you need relational, document, or graph data
- Real-time applications (chat, collaborative editing, live dashboards)
- Prototyping: zero configuration, zero schema overhead
- Apps where user auth is complex and you want the DB to handle it

**When NOT to use it:**
- You need the ecosystem of PostgreSQL (pgvector, PostGIS, etc.)
- Extreme write performance (ClickHouse, Cassandra are better)
- Your team already knows PostgreSQL deeply

---

## 2. SurrealQL — The Query Language

SurrealQL looks like SQL but with new capabilities:

```sql
-- CREATE (like INSERT but more flexible)
CREATE person:alice SET
    name = "Alice Johnson",
    age = 28,
    email = "alice@example.com",
    tags = ["engineer", "gopher"];

-- SELECT (standard SQL)
SELECT * FROM person WHERE age > 25;

-- SELECT with nested object access
SELECT name, tags[0] AS primary_tag FROM person;

-- UPDATE
UPDATE person:alice SET age = 29;

-- DELETE
DELETE person:alice;

-- UPSERT
UPSERT person:alice SET
    name = "Alice Johnson",
    age = 29;

-- Record IDs: tablename:identifier
-- "person:alice" uniquely identifies Alice's record
-- No need for surrogate keys!
```

SurrealQL also supports math operations, array functions, and custom functions:

```sql
-- Array operations
SELECT * FROM person WHERE "engineer" IN tags;
SELECT array::len(tags) AS tag_count FROM person;

-- Math
SELECT name, age * 2 AS doubled_age FROM person;

-- String functions
SELECT string::uppercase(name) FROM person;

-- Time
SELECT time::now() AS current_time;
SELECT * FROM orders WHERE created_at > time::now() - 7d;
```

---

## 3. Tables, Records, and Documents

SurrealDB has a flexible schema system:

```sql
-- Schemafull table (like PostgreSQL — only defined fields allowed)
DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD username ON user TYPE string;
DEFINE FIELD email    ON user TYPE string ASSERT string::is::email($value);
DEFINE FIELD age      ON user TYPE int    ASSERT $value >= 18;
DEFINE FIELD created  ON user TYPE datetime DEFAULT time::now();
DEFINE INDEX email_unique ON user FIELDS email UNIQUE;

-- Schemaless table (like MongoDB — any fields allowed)
DEFINE TABLE logs SCHEMALESS;

-- Mixed: some enforced fields, anything else allowed
DEFINE TABLE product SCHEMALESS;
DEFINE FIELD price ON product TYPE float ASSERT $value > 0;
-- name, description, images etc. are free-form

-- Create with enforced schema
CREATE user SET username = "alice", email = "alice@example.com", age = 25;
-- This would fail: age = 15 (fails ASSERT $value >= 18)
```

**Document-style storage:**

```sql
CREATE article SET
    title = "Learn SurrealDB",
    content = "SurrealDB is a multi-model database...",
    author = person:alice,  -- reference to another record
    metadata = {
        word_count: 1500,
        read_time: "5 minutes",
        tags: ["database", "tutorial"]
    },
    published = true;
```

---

## 4. Graph Relationships

SurrealDB has first-class graph support. You can traverse relationships without JOINs.

```sql
-- Create people
CREATE person:alice SET name = "Alice";
CREATE person:bob   SET name = "Bob";
CREATE person:carol SET name = "Carol";

-- Create relationships (edges in the graph)
RELATE person:alice -> knows -> person:bob   SET since = "2020-01-01";
RELATE person:bob   -> knows -> person:carol SET since = "2021-06-15";
RELATE person:alice -> likes -> article:1;

-- Graph traversal: friends of Alice
SELECT ->knows->person.name AS friends FROM person:alice;
-- Returns: ["Bob"]

-- Friends of friends (2 hops)
SELECT ->knows->person->knows->person.name AS fof FROM person:alice;
-- Returns: ["Carol"]

-- Who does Alice like?
SELECT ->likes->article.title AS liked_articles FROM person:alice;

-- Who likes a specific article?
SELECT <-likes<-person.name AS liked_by FROM article:1;

-- Path: can Alice reach Carol?
SELECT path FROM person:alice SHORTEST person:carol;
```

This is what makes graph databases powerful: traversing relationships in SQL requires complex multi-table JOINs. In SurrealDB it's `->relationship->entity`.

**Practical example — social network:**

```sql
-- User connections
RELATE user:alice -> follows -> user:bob;
RELATE user:bob   -> follows -> user:carol;

-- Get Alice's feed: posts from people Alice follows
SELECT ->follows->user->posts.* AS feed
FROM user:alice;

-- Mutual friends of alice and carol
SELECT ->follows->user<-follows<-user:carol.name AS mutual
FROM user:alice;
```

---

## 5. Real-Time Subscriptions

SurrealDB can push changes to connected clients — no polling needed.

```sql
-- Live query: get notified when any person is created/updated/deleted
LIVE SELECT * FROM person;
-- Returns a UUID (the live query ID)

-- Filtered live query
LIVE SELECT * FROM person WHERE age > 25;

-- Kill a live query
KILL "e63b4278-8f49-4f5e-b5b4-f87f0a5f64bb";
```

In a Go WebSocket handler, you'd listen to these notifications and forward them to connected clients — like a real-time collaborative app or live dashboard.

---

## 6. Authentication Built In

SurrealDB has an authentication layer built into the database itself — no separate auth service needed:

```sql
-- Define a "user" scope (authentication namespace)
DEFINE SCOPE user SESSION 24h
  SIGNUP (
    CREATE user SET
      email = $email,
      password = crypto::argon2::generate($password),  -- bcrypt equivalent
      created = time::now()
  )
  SIGNIN (
    SELECT * FROM user
    WHERE email = $email
    AND crypto::argon2::compare(password, $password)
  );

-- Client signs up
SIGNUP {
    NS: "myapp",
    DB: "production",
    SC: "user",
    email: "alice@example.com",
    password: "supersecret123"
}
-- Returns a JWT token

-- Client signs in
SIGNIN {
    NS: "myapp",
    DB: "production",
    SC: "user",
    email: "alice@example.com",
    password: "supersecret123"
}
-- Returns a JWT token

-- Row-level security: users can only see their own data
DEFINE TABLE messages SCHEMAFULL PERMISSIONS
    FOR select WHERE author = $auth.id
    FOR create WHERE author = $auth.id
    FOR update WHERE author = $auth.id
    FOR delete WHERE author = $auth.id;
```

`$auth` is the current authenticated user — automatically injected by SurrealDB.

---

## Summary

- SurrealDB combines relational, document, graph, and real-time in one database with one query language.
- Record IDs (`person:alice`) are first-class — no auto-increment surrogate keys needed.
- Graph traversal (`->knows->person`) replaces complex multi-table JOINs.
- Live queries push changes to clients in real-time — no polling.
- Built-in auth with scopes, row-level permissions, and `$auth` variable.
- Best for new projects needing flexible data models; not a replacement for PostgreSQL's ecosystem.

### Exercises

**Easy:** Start SurrealDB with Docker and connect with `surreal sql`. Create 3 people records and 2 "knows" relationships. Run a graph traversal query to find friends.

**Medium:** Design a SurrealDB schema for a blog: posts, comments, likes, and user follows. Use graph relationships for likes and follows. Query: "find all posts liked by users that Alice follows."

**Hard:** Implement a real-time notification system using SurrealDB's LIVE queries. When user A follows user B, B gets a notification. Use WebSockets in Go to push these notifications to connected clients.
