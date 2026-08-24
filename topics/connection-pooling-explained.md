---
title: Connection Pooling, Explained
category: Software & Programming
tags: [Databases, Performance]
duration: 6 min read
relatedCourses: [databases-and-async-go, go-programming]
relatedProjects: [rest-api-server, key-value-store]
relatedTopics: [b-tree-indexes-explained, circuit-breakers-and-retries]
---

## TL;DR

- Opening a new database connection is expensive (a TCP handshake, authentication, sometimes TLS negotiation) — expensive enough that doing it on every single query would dominate your actual query latency.
- A connection pool opens a fixed set of connections up front and hands them out to whoever needs one, returning them to the pool when done, instead of opening and closing a connection per request.
- Pool size is a real tuning parameter with a real ceiling — it's bounded by how many concurrent connections the *database* can handle, not just how much concurrency your application wants.
- Go's `database/sql` package has a built-in pool; the four settings that actually matter are max open, max idle, and the two connection-lifetime limits.

## Why Not Just Open a Connection Per Request?

```go
// the naive, expensive way
func handleRequest() {
    conn, _ := sql.Open("postgres", dsn) // TCP handshake + auth, every single time
    defer conn.Close()
    conn.Query("SELECT ...")
}
```

Establishing a new database connection typically costs single-digit to tens of milliseconds — a TCP handshake, then the database's own authentication handshake, sometimes TLS negotiation on top. If your actual query only takes 2ms, spending 15ms setting up the connection to run it means the connection setup is the dominant cost of the entire request, not the query. Under real concurrent load, this also means constantly opening and closing connections, which is itself real work for the database server to do repeatedly.

## The Pool

A connection pool solves this by keeping a set of already-established, already-authenticated connections open and ready:

```
Application requests a connection
        |
        v
   [ Pool ] -- has an idle connection? --> hand it out immediately
        |
        no idle connection, but under max limit?
        |
        v
   open a new one, hand it out
        |
        (if at max limit: caller waits, or gets an error, depending on config)

... query runs using that connection ...

Connection returned to the pool (not closed) -> becomes idle, ready for reuse
```

The expensive setup (handshake, auth) happens once per connection, not once per query. Reusing connections turns "open connection + run query + close connection" into just "run query" for the common case where an idle connection is already available.

## Go's database/sql Pool

Go's standard library `database/sql` package includes a connection pool built in — you don't manage it manually, but you do need to configure it:

```go
db, err := sql.Open("postgres", dsn)

db.SetMaxOpenConns(25)                  // hard ceiling on total connections (idle + in-use)
db.SetMaxIdleConns(25)                  // how many idle connections to keep ready, not close
db.SetConnMaxLifetime(30 * time.Minute) // force-recycle a connection after this long, even if healthy
db.SetConnMaxIdleTime(5 * time.Minute)  // close an idle connection if unused this long
```

- **`MaxOpenConns`**: the actual ceiling. Once this many connections are in use, the next request for a connection blocks (up to the query's context deadline) until one is returned to the pool. This is the number that has to be chosen with the *database's* own connection limit in mind — not just "how much concurrency does my app want."
- **`MaxIdleConns`**: how many unused connections stay open, ready to be handed out instantly, versus being closed after use. Setting this lower than `MaxOpenConns` means connections get closed and reopened more often once load drops — usually you want this close to (or equal to) `MaxOpenConns`.
- **`ConnMaxLifetime`**: forces connections to be recycled periodically even if nothing's wrong with them — useful for working with infrastructure that expects periodic reconnection (a load balancer in front of the database, a database that's periodically failed-over) and for avoiding subtle long-lived-connection issues that accumulate over hours or days.
- **`ConnMaxIdleTime`**: closes connections that have been idle too long, so the pool doesn't hold onto more open connections than current load actually needs.

## Why Pool Size Has a Real Ceiling

It's tempting to set `MaxOpenConns` very high "so nothing ever has to wait." This is a mistake, because the actual constraint isn't your application's appetite for concurrency — it's the **database's** connection limit. PostgreSQL's default `max_connections` is 100 (shared across *all* clients, not per application instance); every additional connection also costs the database server memory and scheduling overhead, whether or not it's actively running a query.

Concretely: if you run 10 instances of your application, each configured for `MaxOpenConns = 100`, you're allowing up to 1,000 total connections against a database that might only be configured to accept 100 — the pool doesn't know or care about other application instances, so this coordination has to be a deliberate capacity-planning decision, not a per-instance default.

```
10 app instances × MaxOpenConns=100 = up to 1,000 potential connections
                                        vs.
                database's max_connections = 100 (shared across everyone)
```

This is exactly why many real production setups put a dedicated connection pooler (PgBouncer is the standard one for PostgreSQL) *in front of* the database, multiplexing many application-level connections onto a smaller number of actual database connections — a pool of pools, effectively, because the database's own connection ceiling is a genuinely scarce shared resource.

## Common Pitfalls

- **Setting `MaxOpenConns` without checking the database's actual connection limit** — across multiple application instances, this is the single most common way to accidentally exhaust the database's connection capacity under load.
- **Not returning connections to the pool** — in Go specifically, forgetting to close a `*sql.Rows` (or an error path that returns early without doing so) leaks that connection out of the pool until it's garbage collected, effectively shrinking your available pool size over time under load.
- **`MaxIdleConns` set much lower than `MaxOpenConns`** — under bursty traffic, this causes connections to be closed right after use and reopened again shortly after, defeating much of the point of pooling during exactly the bursts where it matters most.
- **No `ConnMaxLifetime`** — long-lived connections behind infrastructure that expects periodic rotation (DNS-based failover, a load balancer with its own idle timeout) can silently go stale, causing intermittent errors that are hard to reproduce because they depend on how long a given connection happened to live.
