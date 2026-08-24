# Chapter 45: 12-Week Study Plan & Mock Interview Questions

This chapter gives you a concrete week-by-week study plan to prepare for senior software engineering interviews at Google, Meta, Apple, Amazon, Netflix, Stripe, Uber, and Airbnb. It also includes a comprehensive bank of mock interview questions with model answers.

## Table of Contents

1. [How to Use This Plan](#1-how-to-use-this-plan)
2. [12-Week Study Plan](#2-12-week-study-plan)
3. [Algorithms & Data Structures Mock Questions](#3-algorithms--data-structures-mock-questions)
4. [System Design Mock Questions](#4-system-design-mock-questions)
5. [Go Deep Dive Mock Questions](#5-go-deep-dive-mock-questions)
6. [Database Mock Questions](#6-database-mock-questions)
7. [Final Week Checklist](#7-final-week-checklist)

---

## 1. How to Use This Plan

```
Target: 330 hours over 12 weeks (~27 hours/week)
  Weekdays: 3-4 hours/day
  Weekends: 5-6 hours/day

Track your progress:
  Each chapter has exercises — do them
  After each coding problem: time yourself
  Every 2 weeks: do a mock interview (swap with a friend or use Pramp/interviewing.io)

The most common mistake: reading without practicing
  Reading a solution = 20% of the learning
  Writing it from scratch = 80%
  Being able to write it under time pressure = 100%
```

---

## 2. 12-Week Study Plan

### Week 1-2: DSA Foundations
```
Topics:
  Ch 01: How senior interviews work
  Ch 02: Complexity analysis — be able to derive Big O instantly
  Ch 03: Arrays, strings, hashing — sliding window, two pointers, prefix sum

Daily practice:
  1 easy LeetCode (warm up)
  1 medium LeetCode (core practice)
  Target: 14 problems this week
  
Problems to solve:
  - Two Sum, Group Anagrams, Longest Consecutive Sequence
  - Longest Substring Without Repeating Characters
  - Minimum Window Substring (hard — important)
  - Container With Most Water, 3Sum
  - Product of Array Except Self, Subarray Sum Equals K
```

### Week 3-4: Trees, Graphs, Advanced DSA
```
Topics:
  Ch 04: Linked lists — slow/fast pointers, reversal, merge
  Ch 05: Stacks & queues — monotonic stack patterns
  Ch 06: Trees — DFS, BFS, BST, LCA, serialization
  Ch 07: Graphs — BFS, DFS, topological sort, Union-Find
  Ch 08: Shortest paths — Dijkstra, Bellman-Ford

Problems to solve:
  - LRU Cache (medium — very common)
  - Max Sliding Window (hard)
  - Largest Rectangle in Histogram (hard)
  - Binary Tree Right Side View, Level Order
  - Course Schedule I & II (topological sort)
  - Clone Graph, Pacific Atlantic Water Flow
  - Network Delay Time (Dijkstra)
```

### Week 5: Dynamic Programming
```
Topics:
  Ch 09: DP — 5 patterns: linear, grid, interval, knapsack, state machine

Focus problems (hardest for most people):
  - Climbing Stairs, House Robber, Jump Game
  - Unique Paths, Minimum Path Sum
  - Longest Increasing Subsequence (both O(n²) and O(n log n))
  - Coin Change, Word Break
  - Edit Distance (most important DP problem)
  - Regular Expression Matching (hard)
  
Study approach: implement top-down first, then convert to bottom-up
```

### Week 6: Go Deep Dive
```
Topics:
  Ch 10-11: Go memory model, GMP scheduler
  Ch 12-13: Channels, sync package
  Ch 14-15: Context package, GC internals
  Ch 16-17: Concurrency patterns, goroutine leaks & races
  Ch 18-19: Performance profiling, testing

Practice projects:
  1. Implement a thread-safe bounded buffer (producer-consumer)
  2. Implement a concurrent key-value store with eviction
  3. Write benchmarks and find the bottleneck
  4. Detect a race condition with go test -race
```

### Week 7: Node.js, TypeScript, Databases
```
Topics:
  Ch 20-22: Node.js event loop, async JS, TypeScript advanced
  Ch 23-24: SQL mastery, index deep dive
  Ch 25-26: Transactions & MVCC, PostgreSQL internals
  Ch 27-28: NoSQL decision framework, database scaling

Practice:
  Write 10 complex SQL queries (window functions, CTEs, recursive CTEs)
  Set up PostgreSQL locally, run EXPLAIN ANALYZE on slow queries
  Design the data model for a social network from scratch
```

### Week 8: Distributed Systems
```
Topics:
  Ch 29-30: CAP theorem, Raft/etcd/distributed locks
  Ch 31-32: Reliability patterns (retry, circuit breaker, saga, outbox), Kafka
  Ch 33-34: Service communication (REST/gRPC/WebSockets), rate limiting

Practice discussions:
  Explain Raft to yourself out loud (the "rubber duck" test)
  Design a rate limiter that works across 100 servers
  Explain the outbox pattern and when you'd use it
```

### Week 9: System Design
```
Topics:
  Ch 35: System design framework — the 45-minute approach
  Ch 36: URL shortener & Pastebin
  Ch 37: Chat system (WhatsApp) & video platform (YouTube)
  Ch 38: Uber & Stripe

Practice:
  Do 3 timed system design problems (45 minutes each):
    - Design TinyURL (warmup)
    - Design a notification service (intermediate)
    - Design Twitter feed (advanced)
  
  Record yourself — watch it back. Identify where you got stuck.
```

### Week 10: Low-Level Design & Security
```
Topics:
  Ch 39: SOLID principles & design patterns
  Ch 40: LLD — Parking Lot, Elevator, Library
  Ch 41: Networking — TCP/HTTP/TLS
  Ch 42: Security — OWASP, JWT, OAuth2

Practice:
  Code the parking lot from scratch in Go — time yourself
  Design and implement a rate limiter class
  Trace through what happens on a TLS handshake
```

### Week 11: Infrastructure & Behavioral
```
Topics:
  Ch 43: Docker, Kubernetes, Observability
  Ch 44: Behavioral interviews — STAR method

Practice:
  Deploy a Go service to Kubernetes locally (use minikube or kind)
  Write 10 STAR stories from your career
  Do 2 mock behavioral interviews with a friend
  Record yourself answering behavioral questions
```

### Week 12: Mock Interviews & Final Preparation
```
Focus: execution under pressure, not new material

Day 1-2: Full mock coding interviews (45 min each, 2-3 problems)
Day 3:   Full mock system design interview (45 min)
Day 4:   Full mock behavioral interview (45 min)
Day 5:   Review weak areas identified in mocks
Day 6:   Light review of key concepts — no new material
Day 7:   Rest. Your brain consolidates learning while you sleep.

The day before the interview:
  - Review your 10 STAR stories
  - Review the company's engineering blog and recent news
  - Get 8 hours of sleep
  - Prepare your coding environment (IDE, language, templates)
```

---

## 3. Algorithms & Data Structures Mock Questions

**Q: Find all triplets in an array that sum to zero.**
```go
// Approach: sort + two pointers
// Time: O(n²), Space: O(1)
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int
    
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue } // skip duplicates
        
        left, right := i+1, len(nums)-1
        for left < right {
            sum := nums[i] + nums[left] + nums[right]
            if sum == 0 {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                for left < right && nums[left] == nums[left+1] { left++ }
                for left < right && nums[right] == nums[right-1] { right-- }
                left++
                right--
            } else if sum < 0 {
                left++
            } else {
                right--
            }
        }
    }
    return result
}
```

**Q: Implement an LRU Cache.**
```go
type LRUCache struct {
    cap   int
    cache map[int]*Node
    head  *Node // most recently used (dummy)
    tail  *Node // least recently used (dummy)
}

type Node struct {
    key, val   int
    prev, next *Node
}

func Constructor(capacity int) LRUCache {
    head, tail := &Node{}, &Node{}
    head.next, tail.prev = tail, head
    return LRUCache{cap: capacity, cache: make(map[int]*Node), head: head, tail: tail}
}

func (c *LRUCache) Get(key int) int {
    if node, ok := c.cache[key]; ok {
        c.moveToFront(node)
        return node.val
    }
    return -1
}

func (c *LRUCache) Put(key, val int) {
    if node, ok := c.cache[key]; ok {
        node.val = val
        c.moveToFront(node)
        return
    }
    node := &Node{key: key, val: val}
    c.cache[key] = node
    c.addToFront(node)
    if len(c.cache) > c.cap {
        lru := c.tail.prev
        c.remove(lru)
        delete(c.cache, lru.key)
    }
}

func (c *LRUCache) remove(n *Node)      { n.prev.next = n.next; n.next.prev = n.prev }
func (c *LRUCache) addToFront(n *Node)  { n.next = c.head.next; n.prev = c.head; c.head.next.prev = n; c.head.next = n }
func (c *LRUCache) moveToFront(n *Node) { c.remove(n); c.addToFront(n) }
```

---

## 4. System Design Mock Questions

**Q: Design a notification service that sends email, SMS, and push notifications.**

Key points to cover:
```
Requirements:
  - Multiple channels (email, SMS, push)
  - Priority levels (transactional = high, marketing = low)
  - Rate limiting per user
  - Delivery tracking (sent, delivered, failed)
  - Template management

Architecture:
  API Service → Kafka (topic per priority) → Channel Workers
  
  High priority queue: processed immediately (OTP, payment alerts)
  Low priority queue: rate limited, batched (newsletters, promotions)
  
  Channel Workers:
    Email worker → SendGrid/SES
    SMS worker → Twilio
    Push worker → FCM (Android), APNs (iOS)
  
  Each worker:
    - Retries on failure (exponential backoff)
    - Records delivery status in database
    - Webhooks from providers confirm delivery/failure

Data model:
  notification_logs(id, user_id, channel, template_id, status, sent_at, delivered_at)
  user_preferences(user_id, email_enabled, sms_enabled, push_enabled, quiet_hours)
```

**Q: Design a feature flag service.**

```
Core requirements:
  - Toggle features for specific users, groups, or percentages
  - Sub-100ms evaluation time
  - Updates propagate within 30 seconds
  - Audit log of all changes

Design:
  Admin UI → Flag Service → PostgreSQL (source of truth)
                         → Kafka "flag-changes"
                         → SDK agents in each service (local cache)
  
  SDK evaluation (on each request):
    1. Check local in-memory cache (< 1ms, shared across goroutines)
    2. Cache TTL = 30 seconds (acceptable staleness)
    3. Kafka consumer updates cache on flag change events
  
  Targeting rules (evaluated in order):
    1. User-specific override (whitelist for early access)
    2. Group membership (admins, beta users)
    3. Percentage rollout (random but stable per user — hash(user_id) % 100 < percentage)
    4. Default (on/off)
```

---

## 5. Go Deep Dive Mock Questions

**Q: What happens when a goroutine sends on a nil channel?**

"Sending to a nil channel blocks forever — the goroutine is suspended and never woken up. This is a goroutine leak. The same happens for receiving from a nil channel. A nil channel select case is never selected. This behavior is useful when you want to dynamically disable a case in a select — set the channel to nil to remove it from consideration. But in normal use, sending to a nil channel is almost always a bug."

**Q: Explain the difference between `sync.Mutex` and `sync.RWMutex` and when to use each.**

"sync.Mutex has two states: locked and unlocked. Only one goroutine can hold it at a time, regardless of whether it's reading or writing. sync.RWMutex distinguishes reads and writes: multiple goroutines can hold the read lock simultaneously (since reads don't interfere), but a write lock is exclusive. Use RWMutex when: (1) your workload is read-heavy (> ~80% reads), and (2) the critical section is non-trivial (longer than a few nanoseconds). For very short critical sections or write-heavy workloads, the overhead of RWMutex's complexity can make it slower than a plain Mutex. Always measure with benchmarks — don't assume RWMutex is faster."

---

## 6. Database Mock Questions

**Q: You have a query that takes 30 seconds. How do you debug it?**

```
Step 1: EXPLAIN ANALYZE — understand what the planner is doing
  Look for: Seq Scan on large tables, high "Rows Removed by Filter", 
            actual >> estimated rows (stale stats), Nested Loop with high loops

Step 2: Check indexes
  Is there an index on the columns in WHERE and JOIN clauses?
  Is the planner using the index? (if not: might be low cardinality or stale stats)
  Is the index selective enough?

Step 3: Check statistics
  ANALYZE table; — refresh statistics, then re-run EXPLAIN
  Check pg_stats for n_distinct, correlation

Step 4: Look for N+1
  Is this query running in a loop? Is an ORM doing N+1 fetches?

Step 5: Check the query itself
  Is there a function call on an indexed column? (prevents index use)
  Is there a leading wildcard LIKE? (prevents B-tree index)
  Can you rewrite with EXISTS instead of IN?
  
Step 6: Consider schema changes
  Add a covering index (INCLUDE extra columns)
  Add a partial index (WHERE clause for frequently queried subset)
  Denormalize if this query runs millions of times/day
```

---

## 7. Final Week Checklist

### Coding
```
□ Can solve medium LeetCode in 20-25 minutes
□ Can solve hard LeetCode in 35-40 minutes with hints
□ Know Big O for every common data structure operation
□ Can implement: BST, heap, trie from scratch
□ Patterns memorized: sliding window, two pointers, BFS, DFS, DP
```

### System Design
```
□ 45-minute framework is automatic
□ Capacity estimation is fast (minutes, not the whole session)
□ Know tradeoffs: SQL vs NoSQL, push vs pull, sync vs async
□ Can design: URL shortener, chat system, feed system, payment system
□ Deep dives: caching strategies, database sharding, consensus
```

### Go
```
□ GMP scheduler — can explain in 2 minutes
□ Channel axioms — nil, open, closed behaviors
□ Common concurrency patterns: worker pool, pipeline, fan-out/in
□ Goroutine leak detection and prevention
□ Garbage collector phases and escape analysis
□ go test -race, pprof, benchmarks
```

### Behavioral
```
□ 10 STAR stories prepared
□ Can discuss a genuine failure with ownership
□ Can give an example of influencing without authority
□ Questions ready for the interviewer
```

### The Interview Day
```
□ 8 hours of sleep
□ Know the interviewer format (how many rounds, what types)
□ Have water nearby (dry mouth is common)
□ "I'd like to think through this for a moment" is always okay
□ If stuck: explain your thinking out loud — interviewers guide you
□ For system design: draw, then talk. Don't describe without drawing.
```

---

## Summary

- 12 weeks, ~330 hours. Weeks 1-5: DSA. Week 6: Go. Week 7: Databases. Weeks 8-9: Distributed systems and system design. Weeks 10-11: LLD and behavioral. Week 12: mock interviews and consolidation.
- The biggest mistake: reading without writing. Solve every problem from scratch with a timer.
- System design: the 45-minute framework is not optional — practice it until it's automatic.
- Behavioral: 10 concrete STAR stories are enough for 80% of questions. Know the leadership principles of your target company.
- Final week: mock interviews, sleep, review your own stories. No new material.
- In the interview: thinking out loud is your biggest asset. Interviewers want to see your mind work.
