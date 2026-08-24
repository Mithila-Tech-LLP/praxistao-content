# Chapter 11: Go Concurrency — Goroutines, Channels, and Security Tools

*Concurrency is what makes security tools fast. A port scanner that checks one port at a time would take hours. With goroutines, you check thousands simultaneously. This chapter shows how Go's concurrency model powers real security tools.*

---

## Why Concurrency Matters for Security Tools

- **Port scanners:** Check 65535 ports — 1 at a time = hours; 1000 at a time = seconds
- **Password crackers:** Try thousands of passwords per second
- **Web fuzzers:** Send thousands of HTTP requests simultaneously
- **Network sniffers:** Capture packets while simultaneously processing them
- **Monitoring agents:** Watch files, processes, and network simultaneously

---

## Goroutines — Lightweight Threads

A goroutine is a function that runs concurrently with other functions. The `go` keyword starts one.

```go
package main

import (
    "fmt"
    "time"
)

func scan(ip string) {
    // simulates a slow network operation
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Scanned %s\n", ip)
}

func main() {
    ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
    
    // Sequential: takes 300ms
    for _, ip := range ips {
        scan(ip)
    }
    
    // Concurrent: takes ~100ms (all 3 run in parallel)
    for _, ip := range ips {
        go scan(ip)   // ← goroutine starts, doesn't wait
    }
    
    time.Sleep(200 * time.Millisecond) // wait for goroutines to finish
    fmt.Println("Done")
}
```

**Goroutines vs threads:**
- OS threads: ~1MB stack, kernel overhead
- Goroutines: ~2KB stack, multiplexed onto OS threads by Go runtime
- You can run 100,000+ goroutines comfortably; 100,000 OS threads would crash the system

---

## The Problem: Race Conditions

When goroutines share data, they can interfere with each other.

```go
// BROKEN — data race
var openPorts []int

for port := 1; port <= 100; port++ {
    go func(p int) {
        if isOpen(p) {
            openPorts = append(openPorts, p)  // multiple goroutines write simultaneously!
        }
    }(port)
}
```

Two goroutines calling `append` at the same time can corrupt the slice.

**Detect races:**
```bash
go run -race main.go
# Output: WARNING: DATA RACE
```

---

## Channels — Goroutine Communication

Channels are typed pipes. One goroutine sends, another receives.

```go
ch := make(chan int)        // unbuffered
ch := make(chan int, 100)   // buffered (capacity 100)

// Send
ch <- 42

// Receive
value := <-ch

// Close
close(ch)

// Range over channel (until closed)
for value := range ch {
    fmt.Println(value)
}
```

### Worker Pool Pattern — The Foundation of Fast Scanners

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

func worker(host string, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for port := range jobs {
        addr := fmt.Sprintf("%s:%d", host, port)
        conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
        if err == nil {
            conn.Close()
            results <- port   // send open port to results channel
        }
    }
}

func main() {
    host := "scanme.nmap.org"
    numWorkers := 200
    
    jobs := make(chan int, numWorkers)
    results := make(chan int, numWorkers)
    var wg sync.WaitGroup
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(host, jobs, results, &wg)
    }
    
    // Feed jobs (in goroutine so we can collect results concurrently)
    go func() {
        for port := 1; port <= 10000; port++ {
            jobs <- port
        }
        close(jobs)
    }()
    
    // Close results when all workers done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect
    var openPorts []int
    for port := range results {
        openPorts = append(openPorts, port)
    }
    
    fmt.Printf("Open ports: %v\n", openPorts)
}
```

**Why this pattern works:**
1. Fixed number of workers (avoids spawning 65535 goroutines)
2. Jobs channel controls work distribution
3. Results channel collects output safely (no races)
4. `sync.WaitGroup` waits for all workers to finish

---

## Mutex — Protecting Shared State

When you can't avoid sharing state, use a mutex (mutual exclusion lock).

```go
package main

import (
    "fmt"
    "sync"
)

type ScanStats struct {
    mu      sync.Mutex
    Open    int
    Closed  int
    Errors  int
}

func (s *ScanStats) RecordOpen() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Open++
}

func (s *ScanStats) RecordClosed() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Closed++
}

func main() {
    stats := &ScanStats{}
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            if n%2 == 0 {
                stats.RecordOpen()
            } else {
                stats.RecordClosed()
            }
        }(i)
    }
    
    wg.Wait()
    fmt.Printf("Open: %d, Closed: %d\n", stats.Open, stats.Closed)
    // Always: Open: 500, Closed: 500
}
```

**sync.RWMutex:** For data read often, written rarely:
```go
var mu sync.RWMutex

// Multiple readers at once
mu.RLock()
value := sharedData
mu.RUnlock()

// Only one writer at a time
mu.Lock()
sharedData = newValue
mu.Unlock()
```

---

## Context — Cancellation and Timeouts

Context lets you cancel goroutines and set deadlines.

```go
package main

import (
    "context"
    "fmt"
    "net"
    "time"
)

func scanWithContext(ctx context.Context, host string, port int) (bool, error) {
    dialer := &net.Dialer{}
    conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        return false, err
    }
    conn.Close()
    return true, nil
}

func main() {
    // Cancel after 5 seconds total
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    results := make(chan struct{port int; open bool}, 100)
    
    for port := 1; port <= 1000; port++ {
        go func(p int) {
            open, _ := scanWithContext(ctx, "192.168.1.1", p)
            select {
            case results <- struct{port int; open bool}{p, open}:
            case <-ctx.Done():  // stop if context cancelled
                return
            }
        }(port)
    }
    
    timeout := time.After(6 * time.Second)
    for {
        select {
        case r := <-results:
            if r.open {
                fmt.Printf("Port %d is open\n", r.port)
            }
        case <-ctx.Done():
            fmt.Println("Scan timed out or cancelled")
            return
        case <-timeout:
            return
        }
    }
}
```

---

## Select — Multiplexing Channels

`select` waits on multiple channel operations.

```go
select {
case msg := <-ch1:
    fmt.Println("from ch1:", msg)
case msg := <-ch2:
    fmt.Println("from ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("timeout — no message")
case <-ctx.Done():
    fmt.Println("cancelled")
}
```

Common pattern in security tools:
```go
for {
    select {
    case event := <-fileEvents:
        processFileEvent(event)
    case event := <-networkEvents:
        processNetworkEvent(event)
    case <-ticker.C:
        heartbeat()
    case <-ctx.Done():
        return  // clean shutdown
    }
}
```

---

## Real Example: Concurrent HTTP Fuzzer

```go
package main

import (
    "bufio"
    "flag"
    "fmt"
    "net/http"
    "os"
    "sync"
    "time"
)

type FuzzResult struct {
    Path   string
    Status int
    Length int64
}

func fuzzWorker(
    client *http.Client,
    baseURL string,
    paths <-chan string,
    results chan<- FuzzResult,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    for path := range paths {
        url := baseURL + "/" + path
        resp, err := client.Get(url)
        if err != nil {
            continue
        }
        resp.Body.Close()
        
        // Filter: only interesting status codes
        if resp.StatusCode != 404 {
            results <- FuzzResult{
                Path:   path,
                Status: resp.StatusCode,
                Length: resp.ContentLength,
            }
        }
    }
}

func main() {
    target := flag.String("u", "", "Target URL")
    wordlist := flag.String("w", "/usr/share/wordlists/dirb/common.txt", "Wordlist")
    workers := flag.Int("t", 50, "Threads")
    flag.Parse()

    if *target == "" {
        fmt.Fprintln(os.Stderr, "Usage: fuzzer -u http://target.com -w wordlist.txt")
        os.Exit(1)
    }

    client := &http.Client{
        Timeout: 5 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse  // don't follow redirects
        },
    }

    f, err := os.Open(*wordlist)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    paths := make(chan string, *workers)
    results := make(chan FuzzResult, *workers)
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < *workers; i++ {
        wg.Add(1)
        go fuzzWorker(client, *target, paths, results, &wg)
    }

    // Close results when done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Feed wordlist
    go func() {
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            path := scanner.Text()
            if path != "" && path[0] != '#' {
                paths <- path
            }
        }
        close(paths)
    }()

    // Print results
    fmt.Printf("%-50s %-10s %s\n", "Path", "Status", "Length")
    fmt.Println("---")
    for r := range results {
        fmt.Printf("%-50s %-10d %d\n", r.Path, r.Status, r.Length)
    }
}
```

```bash
go run fuzzer.go -u http://192.168.100.101 -w /usr/share/wordlists/dirb/common.txt -t 100
```

---

## Common Concurrency Patterns Summary

| Pattern | Use case | Key primitives |
|---------|---------|----------------|
| **Worker pool** | Bounded concurrent work (scanning, fuzzing) | `channel`, `WaitGroup` |
| **Fan-out** | Send work to many goroutines | Multiple goroutines reading same channel |
| **Fan-in** | Aggregate from many goroutines | Multiple goroutines writing same channel |
| **Pipeline** | Chain processing stages | Channel between each stage |
| **Timeout** | Don't wait forever | `context.WithTimeout`, `time.After` |
| **Rate limiting** | Don't overwhelm target | `time.Tick` + channel |
| **Graceful shutdown** | Clean stop | `context.WithCancel`, `signal.Notify` |

---

## Exercises

1. Build a ping sweeper using goroutines: given a /24 subnet, find all hosts that are up (ICMP or TCP connect to port 80)
2. Add rate limiting to the HTTP fuzzer: max 100 requests/second regardless of worker count
3. Implement a graceful shutdown in the port scanner: when Ctrl+C is pressed, wait for in-flight scans to complete before exiting
4. Build a DNS resolver: take a list of 1000 hostnames, resolve them all concurrently, output IP addresses
5. Add a `--timeout` flag to the fuzzer that cancels the entire scan after N seconds
