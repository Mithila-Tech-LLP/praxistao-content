# Chapter 21: Async JavaScript — Promises, Async/Await & Error Handling

This chapter covers the async patterns that come up in Node.js interviews: how Promises work internally, the common combinators, and error handling patterns that production JavaScript requires.

## Table of Contents

1. [How Promises Work Internally](#1-how-promises-work-internally)
2. [Promise Combinators](#2-promise-combinators)
3. [Async/Await Patterns](#3-asyncawait-patterns)
4. [Error Handling](#4-error-handling)
5. [Common Patterns in Node.js Services](#5-common-patterns-in-nodejs-services)
6. [Summary](#summary)

---

## 1. How Promises Work Internally

A Promise is an object that represents the eventual result of an async operation. It is in one of three states: pending, fulfilled, or rejected.

```javascript
// Promise internal state machine:
const p = new Promise((resolve, reject) => {
    // executor runs synchronously
    if (success) {
        resolve(value);  // pending → fulfilled
    } else {
        reject(error);   // pending → rejected
    }
});

p.then(value => {}) // callback runs as microtask after fulfillment
 .catch(err => {})  // callback runs as microtask after rejection
```

### Promise Chaining

```javascript
fetch('/api/users')
    .then(response => {
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return response.json();  // returning a Promise here chains correctly
    })
    .then(users => {
        return users.filter(u => u.active);
    })
    .then(activeUsers => {
        console.log(activeUsers);
    })
    .catch(err => {
        // catches ANY rejection from the entire chain
        console.error('Failed:', err);
    });
```

---

## 2. Promise Combinators

```javascript
// Promise.all: wait for ALL to resolve; fails fast if ANY rejects
const [users, posts, comments] = await Promise.all([
    fetchUsers(),
    fetchPosts(),
    fetchComments(),
]);
// If any one fails, the whole thing rejects immediately

// Promise.allSettled: wait for ALL, get outcomes regardless of success/failure
const results = await Promise.allSettled([
    fetchUsers(),
    fetchPosts(),
    fetchMightFail(),
]);
results.forEach(r => {
    if (r.status === 'fulfilled') console.log(r.value);
    else console.error(r.reason);
});
// Useful when you want ALL results even if some fail

// Promise.race: resolves/rejects with the FIRST to complete
const result = await Promise.race([
    fetch('/api/primary'),
    fetch('/api/fallback'),
]);
// Useful for timeouts:
const data = await Promise.race([
    fetchData(),
    new Promise((_, reject) => setTimeout(() => reject(new Error('timeout')), 5000))
]);

// Promise.any: resolves with the FIRST to fulfill; rejects if ALL reject
const fastest = await Promise.any([
    fetch('/api/server1'),
    fetch('/api/server2'),
    fetch('/api/server3'),
]);
// Unlike race, ignores rejections unless all fail
```

---

## 3. Async/Await Patterns

```javascript
// Sequential: each awaits the previous
async function sequential() {
    const a = await fetchA(); // waits for A
    const b = await fetchB(); // THEN waits for B
    // Total time = time(A) + time(B)
}

// Parallel: start both, then await both
async function parallel() {
    const aPromise = fetchA(); // starts A
    const bPromise = fetchB(); // starts B (doesn't wait for A)
    const [a, b] = await Promise.all([aPromise, bPromise]);
    // Total time = max(time(A), time(B))
}

// Sequential vs parallel is a common interview mistake:
async function mistake() {
    // This is SEQUENTIAL even though it looks "parallel"
    const a = await fetchA(); // waits
    const b = await fetchB(); // then waits
}

// Retry with async/await
async function withRetry(fn, maxRetries = 3, delay = 1000) {
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
        try {
            return await fn();
        } catch (err) {
            if (attempt === maxRetries) throw err;
            await new Promise(resolve => setTimeout(resolve, delay * attempt)); // exponential backoff
        }
    }
}
```

---

## 4. Error Handling

```javascript
// Basic try/catch with async/await
async function fetchUser(id) {
    try {
        const response = await fetch(`/api/users/${id}`);
        if (!response.ok) {
            throw new Error(`User not found: ${response.status}`);
        }
        return await response.json();
    } catch (err) {
        if (err.name === 'TypeError') {
            throw new Error('Network error: ' + err.message);
        }
        throw err; // re-throw unexpected errors
    }
}

// Unhandled promise rejections (production danger!)
// Node.js terminates the process on unhandled rejections in newer versions
process.on('unhandledRejection', (reason, promise) => {
    console.error('Unhandled rejection:', reason);
    process.exit(1); // exit cleanly rather than undefined behavior
});
```

---

## 5. Common Patterns in Node.js Services

```javascript
// Pattern: async middleware in Express
app.get('/users/:id', async (req, res, next) => {
    try {
        const user = await userService.getById(req.params.id);
        res.json(user);
    } catch (err) {
        next(err); // pass to Express error handler
    }
});

// Wrapper to avoid try/catch boilerplate:
const asyncHandler = fn => (req, res, next) =>
    Promise.resolve(fn(req, res, next)).catch(next);

app.get('/users/:id', asyncHandler(async (req, res) => {
    const user = await userService.getById(req.params.id);
    res.json(user);
}));
```

---

## Summary

- Promise states: pending → fulfilled or rejected (one-way transitions).
- `Promise.all`: all must succeed. `Promise.allSettled`: get all results. `Promise.race`: first wins. `Promise.any`: first to succeed.
- `async/await` is syntactic sugar for Promises. `await` inside an async function pauses that function but doesn't block the thread.
- Start independent Promises before awaiting them for parallelism.
- Always handle errors: try/catch with async/await, or .catch() on Promises. Handle `unhandledRejection` in production.
