---
title: Concurrency
---
Concurrency is Go's signature feature and the reason it's everywhere in backend infrastructure. A backend server routinely has thousands of requests in flight at once — Go gives you goroutines and channels to handle that without the complexity of manual thread management.

### Goroutines
A goroutine is a function that runs concurrently with the rest of your program, and Go makes starting one as cheap as a single keyword. Understand what that actually buys you and where the sharp edges are.

**Resources:**
- [Goroutines](course:go-programming#18-goroutines)

### Channels & select
Channels are how goroutines talk to each other safely. `select` lets a goroutine wait on multiple channels at once — the backbone of timeouts, cancellation, and fan-in/fan-out patterns.

**Resources:**
- [Channels](course:go-programming#19-channels)
- [Select, Timeouts, and Non-Blocking Operations](course:go-programming#20-select-timeouts-and-non-blocking-operations)

### sync Package & Context
Not everything needs a channel — mutexes and wait groups handle simpler cases. `context.Context` is how Go threads cancellation and deadlines through a whole call chain, which every real HTTP handler needs.

**Resources:**
- [sync Package](course:go-programming#21-sync-package)
- [Context and Cancellation](course:go-programming#22-context-and-cancellation)

### Concurrency Patterns & Pitfalls
> optional

Worker pools, pipelines, fan-out/fan-in — the recurring shapes concurrent Go code takes — plus the failure modes (goroutine leaks, races, deadlocks) that catch people who skip this.

**Resources:**
- [Concurrency Patterns](course:senior-engineer-interview#16-concurrency-patterns)
- [Goroutine Leaks, Races, and Deadlocks](course:senior-engineer-interview#17-goroutine-leaks-races-deadlocks)

### Practice: Build a Concurrent Log Processor
> branches-from: sync Package & Context

Apply goroutines, channels, and `sync` to a program that actually needs them — processing many log entries concurrently instead of one at a time.

**Resources:**
- [Mini Project: Concurrent Log Processor](course:go-programming#28-mini-project-2-concurrent-log-processor)
