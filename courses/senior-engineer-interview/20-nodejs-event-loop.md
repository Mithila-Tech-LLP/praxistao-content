# Chapter 20: JavaScript Event Loop & Node.js Runtime

This chapter covers the 10% of this course dedicated to Node.js. As a Go-primary engineer who also works with Node.js, you need to understand Node's concurrency model, V8, and the event loop deeply enough to answer senior-level interview questions and debug production issues.

## Table of Contents

1. [V8 Engine & Node.js Architecture](#1-v8-engine--nodejs-architecture)
2. [The Event Loop](#2-the-event-loop)
3. [Call Stack, Callback Queue & Microtask Queue](#3-call-stack-and-queues)
4. [Worker Threads & Cluster Mode](#4-worker-threads--cluster-mode)
5. [Common Node.js Interview Questions](#5-common-nodejs-interview-questions)
6. [Summary](#summary)

---

## 1. V8 Engine & Node.js Architecture

Node.js is a JavaScript runtime built on:
- **V8** — Google's JavaScript engine (same one in Chrome). Compiles JS to machine code using JIT compilation.
- **libuv** — C library providing the event loop, thread pool, and async I/O on all platforms.
- **Node.js APIs** — bindings that expose OS-level APIs (file system, network, etc.) to JavaScript.

```
User Code (JavaScript)
     ↓
Node.js Core APIs (fs, http, crypto, etc.)
     ↓
V8 (JavaScript execution, JIT compilation)
libuv (event loop, async I/O, thread pool)
     ↓
Operating System
```

**Key insight:** JavaScript in Node.js is single-threaded. V8 runs one thread. But libuv handles I/O asynchronously using the OS's async I/O mechanisms (epoll, kqueue) so that the single JS thread is never blocked waiting for I/O.

---

## 2. The Event Loop

The event loop is the mechanism that allows Node.js to perform non-blocking I/O despite being single-threaded. It runs in phases:

```
   ┌───────────────────────────┐
┌─>│           timers          │  ← setTimeout, setInterval callbacks
│  └─────────────┬─────────────┘
│  ┌─────────────┴─────────────┐
│  │     pending callbacks     │  ← I/O callbacks deferred to next loop
│  └─────────────┬─────────────┘
│  ┌─────────────┴─────────────┐
│  │       idle, prepare       │  ← internal use only
│  └─────────────┬─────────────┘
│  ┌─────────────┴─────────────┐
│  │           poll            │  ← retrieve new I/O events, execute I/O callbacks
│  └─────────────┬─────────────┘
│  ┌─────────────┴─────────────┐
│  │           check           │  ← setImmediate callbacks
│  └─────────────┬─────────────┘
│  ┌─────────────┴─────────────┐
└──┤      close callbacks      │  ← socket.on('close', ...) callbacks
   └───────────────────────────┘
```

### Phase Descriptions

**Timers phase:** Executes callbacks from `setTimeout` and `setInterval` whose delay has expired.

**Poll phase:** The heart of the event loop. Retrieves new I/O events (file reads, network connections, etc.) and executes their callbacks. If the poll queue is empty and setImmediate callbacks exist, it moves to the check phase. Otherwise it waits for I/O events.

**Check phase:** Executes `setImmediate` callbacks.

**Between each phase:** Node.js drains the microtask queue (Promises, queueMicrotask). Microtasks always run before the next event loop phase.

---

## 3. Call Stack and Queues

Understanding the execution order is essential for interview questions.

```
EXECUTION ORDER:
1. Synchronous code runs on the call stack
2. When call stack empties: drain the microtask queue completely
3. Move to next event loop phase (timers/poll/check)
4. Before each event loop task: drain microtask queue again
```

```javascript
console.log('1 - sync'); // call stack

setTimeout(() => console.log('2 - setTimeout'), 0); // macrotask queue

Promise.resolve().then(() => console.log('3 - promise')); // microtask queue

queueMicrotask(() => console.log('4 - queueMicrotask')); // microtask queue

setImmediate(() => console.log('5 - setImmediate')); // check phase

console.log('6 - sync');

// OUTPUT:
// 1 - sync         (synchronous)
// 6 - sync         (synchronous)
// 3 - promise      (microtask — runs before event loop phases)
// 4 - queueMicrotask (microtask — same queue)
// 2 - setTimeout   (event loop: timers phase)
// 5 - setImmediate (event loop: check phase)
```

### Why Promises Execute Before setTimeout

```javascript
// Common interview question:
setTimeout(() => console.log('timeout'), 0);
Promise.resolve().then(() => console.log('promise'));

// OUTPUT: promise, then timeout
// WHY: Promise callbacks (microtasks) drain before the event loop
// moves to its next phase (timers). setTimeout is a macrotask
// that waits for its phase.
```

### async/await is Syntactic Sugar for Promises

```javascript
// These are equivalent:

// async/await version:
async function fetchUser(id) {
    const user = await fetch(`/users/${id}`);
    const data = await user.json();
    return data;
}

// Promise chain version:
function fetchUser(id) {
    return fetch(`/users/${id}`)
        .then(user => user.json())
        .then(data => data);
}

// Under the hood, await suspends the async function and schedules
// the continuation as a Promise .then() callback — a microtask.
// The function resumes when the awaited Promise resolves.
```

---

## 4. Worker Threads & Cluster Mode

Node.js has two ways to use multiple CPU cores.

### Worker Threads (CPU-bound work)

```javascript
const { Worker, isMainThread, parentPort } = require('worker_threads');

if (isMainThread) {
    // Main thread: spawn workers
    const worker = new Worker(__filename);
    worker.on('message', (result) => console.log('Result:', result));
    worker.postMessage({ input: 1000000 });
} else {
    // Worker thread: do CPU-bound work
    parentPort.on('message', ({ input }) => {
        let sum = 0;
        for (let i = 0; i < input; i++) sum += i; // CPU-intensive
        parentPort.postMessage(sum);
    });
}
```

**When to use:** CPU-bound work that would block the event loop — image processing, heavy cryptography, large computation. Workers share the V8 heap (unlike cluster processes) and communicate via message passing or SharedArrayBuffer.

### Cluster Mode (Network load balancing)

```javascript
const cluster = require('cluster');
const os = require('os');
const http = require('http');

if (cluster.isMaster) {
    // Master process: fork workers
    const numCPUs = os.cpus().length;
    for (let i = 0; i < numCPUs; i++) {
        cluster.fork();
    }
    cluster.on('exit', (worker) => {
        console.log(`Worker ${worker.process.pid} died. Forking a new one...`);
        cluster.fork(); // auto-restart
    });
} else {
    // Worker process: run the HTTP server
    http.createServer((req, res) => {
        res.end(`Hello from worker ${process.pid}`);
    }).listen(8000);
}
```

**When to use:** Scale an HTTP server across multiple CPU cores. Each worker is a separate Node.js process (separate V8 instances). The OS/cluster module handles distributing incoming connections.

---

## 5. Common Node.js Interview Questions

**Q: Why is Node.js fast for I/O-bound workloads despite being single-threaded?**

"Because it never blocks on I/O. When Node does a file read or network request, it hands the request to libuv, which uses the OS's async I/O mechanisms (epoll on Linux). The single JS thread continues processing other events. When the I/O completes, its callback is placed in the event loop queue and runs when the thread is free. This is very similar to how Go's goroutines work — but implemented differently (callback queue vs goroutine scheduler). The key insight: single-threaded doesn't mean slow when the thread is never waiting."

**Q: What is the difference between setTimeout(fn, 0) and setImmediate(fn)?**

"Both execute 'as soon as possible', but at different phases of the event loop. setTimeout(fn, 0) fires in the timers phase on the next loop iteration (after the minimum delay, even with 0). setImmediate fires in the check phase. When called from within the main module (not inside an I/O callback), the order can vary due to OS timer precision. But when called from within an I/O callback, setImmediate always fires before setTimeout — the check phase comes before the next timers phase."

**Q: What causes memory leaks in Node.js?**

"Five common causes: (1) Event listeners not removed when objects are destroyed — `EventEmitter` holds references, preventing GC. (2) Global variables accidentally created without `const/let/var` — they live forever. (3) Closures capturing large objects unexpectedly — even if the closure is only called once, the referenced objects can't be GC'd. (4) Caches without eviction policies — the cache grows unbounded. (5) Promises/async operations that reference large objects and are never resolved or rejected — they hang in memory indefinitely."

**Q: When would you use Worker Threads vs Cluster?**

"Worker threads for CPU-bound tasks that run within the context of a single request — like image resizing or JSON parsing of huge payloads. Workers share memory and have lower overhead than forking. Cluster for scaling an HTTP server across CPU cores — each cluster worker is an independent process that handles requests independently. Cluster is about horizontal scaling within one machine; Worker threads are about using multiple cores for a single task."

---

## Summary

- Node.js is single-threaded JS + async I/O via libuv. Never blocking I/O = high throughput for I/O-bound work.
- **Event loop phases:** timers → pending callbacks → poll → check → close. Poll is where I/O callbacks execute.
- **Microtasks** (Promises, queueMicrotask) drain completely before each event loop phase.
- **Order:** sync code → microtasks → timers/setImmediate/I/O callbacks.
- `async/await` is syntactic sugar for Promise chains — continuations are microtasks.
- **Worker threads:** CPU-bound work, shared heap, message passing.
- **Cluster mode:** scale HTTP servers across CPU cores, separate processes.
