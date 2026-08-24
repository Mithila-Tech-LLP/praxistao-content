# Chapter 28: Mini Project 2 — Concurrent Log Processor

A concurrent log processor that reads log files in parallel, parses entries, aggregates stats, and writes filtered results to output files. This project exercises goroutines, channels, worker pools, `sync.WaitGroup`, `sync.Mutex`, and context cancellation from Volume 2.

## What You'll Build

```
logprocessor/
├── main.go
├── parser/
│   └── parser.go      # Log line parsing
├── worker/
│   └── pool.go        # Worker pool
├── aggregator/
│   └── stats.go       # Thread-safe stats
├── writer/
│   └── writer.go      # Output writer
└── testdata/
    └── gen.go         # Test log generator
```

**Features:**
- Read multiple log files concurrently (one goroutine per file)
- Parse each line (level, timestamp, message, fields)
- Filter by log level (INFO, WARN, ERROR)
- Aggregate statistics: counts by level, error frequency, top endpoints
- Write filtered output to a new file
- Graceful cancellation via `context`

---

## 1. Log Line Format

We'll parse a structured log format that's common in Go services:

```
2024-01-15T14:23:05Z INFO  request completed method=GET path=/api/users status=200 latency=12ms
2024-01-15T14:23:06Z ERROR database query failed error="connection refused" query=GetUser
2024-01-15T14:23:07Z WARN  rate limit approaching limit=1000 current=950 window=1m
```

```go
// parser/parser.go
package parser

import (
	"fmt"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(s string) (Level, error) {
	switch strings.TrimSpace(strings.ToUpper(s)) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN", "WARNING":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return DEBUG, fmt.Errorf("unknown level: %q", s)
	}
}

type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Fields    map[string]string
	Raw       string
	Source    string // which file this came from
}

// ParseLine parses one log line.
// Format: <timestamp> <level> <message> [key=value ...]
func ParseLine(line, source string) (*Entry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	// Split into at most 3 parts: timestamp, level, rest
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("too few fields: %q", line)
	}

	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, fmt.Errorf("bad timestamp %q: %w", parts[0], err)
	}

	level, err := ParseLevel(parts[1])
	if err != nil {
		return nil, fmt.Errorf("bad level: %w", err)
	}

	rest := parts[2]
	fields := map[string]string{}

	// Extract key=value pairs from the end of the message.
	// A token is a field if it contains '=' and no spaces in the value,
	// or the value is quoted.
	tokens := tokenize(rest)
	msgTokens := []string{}
	for _, tok := range tokens {
		if k, v, ok := parseField(tok); ok {
			fields[k] = v
		} else {
			msgTokens = append(msgTokens, tok)
		}
	}

	return &Entry{
		Timestamp: ts,
		Level:     level,
		Message:   strings.Join(msgTokens, " "),
		Fields:    fields,
		Raw:       line,
		Source:    source,
	}, nil
}

// tokenize splits a string respecting quoted values.
func tokenize(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			inQuote = !inQuote
			cur.WriteByte(c)
		case c == ' ' && !inQuote:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

// parseField parses "key=value" or "key=\"quoted value\"".
func parseField(tok string) (key, value string, ok bool) {
	idx := strings.IndexByte(tok, '=')
	if idx <= 0 {
		return "", "", false
	}
	k := tok[:idx]
	// Key must be alphanumeric + underscore
	for _, c := range k {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return "", "", false
		}
	}
	v := tok[idx+1:]
	if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		v = v[1 : len(v)-1]
	}
	return k, v, true
}
```

---

## 2. Thread-Safe Statistics

```go
// aggregator/stats.go
package aggregator

import (
	"fmt"
	"sort"
	"sync"

	"logprocessor/parser"
)

// Stats aggregates log metrics across all files, safe for concurrent use.
type Stats struct {
	mu sync.Mutex

	LevelCounts map[parser.Level]int
	TopErrors   map[string]int   // error message → count
	TopPaths    map[string]int   // path field → count
	TotalLines  int
	ParseErrors int
}

func New() *Stats {
	return &Stats{
		LevelCounts: make(map[parser.Level]int),
		TopErrors:   make(map[string]int),
		TopPaths:    make(map[string]int),
	}
}

func (s *Stats) Record(entry *parser.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalLines++
	s.LevelCounts[entry.Level]++

	if entry.Level >= parser.ERROR {
		key := entry.Message
		if len(key) > 80 {
			key = key[:80]
		}
		s.TopErrors[key]++
	}

	if path, ok := entry.Fields["path"]; ok {
		s.TopPaths[path]++
	}
}

func (s *Stats) RecordParseError() {
	s.mu.Lock()
	s.ParseErrors++
	s.mu.Unlock()
}

type Pair struct {
	Key   string
	Count int
}

func top(m map[string]int, n int) []Pair {
	pairs := make([]Pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})
	if len(pairs) > n {
		pairs = pairs[:n]
	}
	return pairs
}

func (s *Stats) Print() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("\n=== Log Processing Results ===\n")
	fmt.Printf("Total lines:  %d\n", s.TotalLines)
	fmt.Printf("Parse errors: %d\n\n", s.ParseErrors)

	fmt.Println("Counts by level:")
	for _, lvl := range []parser.Level{parser.DEBUG, parser.INFO, parser.WARN, parser.ERROR, parser.FATAL} {
		if c := s.LevelCounts[lvl]; c > 0 {
			fmt.Printf("  %-6s %d\n", lvl, c)
		}
	}

	if len(s.TopErrors) > 0 {
		fmt.Println("\nTop errors:")
		for _, p := range top(s.TopErrors, 5) {
			fmt.Printf("  [%3d] %s\n", p.Count, p.Key)
		}
	}

	if len(s.TopPaths) > 0 {
		fmt.Println("\nTop paths:")
		for _, p := range top(s.TopPaths, 5) {
			fmt.Printf("  [%3d] %s\n", p.Count, p.Key)
		}
	}
}
```

---

## 3. Worker Pool

```go
// worker/pool.go
package worker

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"logprocessor/aggregator"
	"logprocessor/parser"
)

type Job struct {
	FilePath    string
	MinLevel    parser.Level
	OutputPath  string
}

type Result struct {
	FilePath string
	Lines    int
	Errors   int
	Err      error
}

// Pool processes log files concurrently.
type Pool struct {
	workers int
	jobs    chan Job
	results chan Result
	stats   *aggregator.Stats
}

func New(workers int, stats *aggregator.Stats) *Pool {
	return &Pool{
		workers: workers,
		jobs:    make(chan Job, workers*2),
		results: make(chan Result, workers*2),
		stats:   stats,
	}
}

func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.work(ctx)
		}()
	}
	// Close results when all workers finish
	go func() {
		wg.Wait()
		close(p.results)
	}()
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Close() {
	close(p.jobs)
}

func (p *Pool) Results() <-chan Result {
	return p.results
}

func (p *Pool) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			result := p.processFile(ctx, job)
			select {
			case p.results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *Pool) processFile(ctx context.Context, job Job) Result {
	result := Result{FilePath: job.FilePath}

	f, err := os.Open(job.FilePath)
	if err != nil {
		result.Err = fmt.Errorf("open %s: %w", job.FilePath, err)
		return result
	}
	defer f.Close()

	var out *os.File
	if job.OutputPath != "" {
		out, err = os.Create(job.OutputPath)
		if err != nil {
			result.Err = fmt.Errorf("create output %s: %w", job.OutputPath, err)
			return result
		}
		defer out.Close()
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return result
		default:
		}

		line := scanner.Text()
		entry, err := parser.ParseLine(line, job.FilePath)
		if err != nil {
			result.Errors++
			p.stats.RecordParseError()
			continue
		}
		if entry == nil {
			continue
		}

		result.Lines++
		p.stats.Record(entry)

		if out != nil && entry.Level >= job.MinLevel {
			fmt.Fprintln(out, entry.Raw)
		}
	}

	if err := scanner.Err(); err != nil {
		result.Err = fmt.Errorf("scan %s: %w", job.FilePath, err)
	}
	return result
}
```

---

## 4. Test Data Generator

```go
// testdata/gen.go
//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	levels   = []string{"DEBUG", "INFO", "INFO", "INFO", "WARN", "ERROR"}
	paths    = []string{"/api/users", "/api/orders", "/api/products", "/health", "/metrics"}
	methods  = []string{"GET", "POST", "PUT", "DELETE"}
	errors_  = []string{
		"connection refused",
		"context deadline exceeded",
		"database query failed",
		"rate limit exceeded",
		"invalid token",
	}
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: gen.go <output-file> <line-count>")
		os.Exit(1)
	}

	f, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n := 10_000
	fmt.Sscan(os.Args[2], &n)

	ts := time.Now().Add(-time.Duration(n) * time.Second)
	for i := 0; i < n; i++ {
		ts = ts.Add(time.Second)
		lvl := levels[rand.Intn(len(levels))]

		switch lvl {
		case "INFO":
			path := paths[rand.Intn(len(paths))]
			method := methods[rand.Intn(len(methods))]
			latency := rand.Intn(500)
			status := 200
			if rand.Float32() < 0.05 {
				status = 500
			}
			fmt.Fprintf(f, "%s %s request completed method=%s path=%s status=%d latency=%dms\n",
				ts.Format(time.RFC3339), lvl, method, path, status, latency)
		case "WARN":
			fmt.Fprintf(f, "%s %s rate limit approaching limit=1000 current=%d window=1m\n",
				ts.Format(time.RFC3339), lvl, 900+rand.Intn(99))
		case "ERROR":
			msg := errors_[rand.Intn(len(errors_))]
			fmt.Fprintf(f, "%s %s %s error=%q query=GetUser\n",
				ts.Format(time.RFC3339), lvl, msg, msg)
		default:
			fmt.Fprintf(f, "%s %s processing request id=%d\n",
				ts.Format(time.RFC3339), lvl, i)
		}
	}

	fmt.Printf("wrote %d lines to %s\n", n, os.Args[1])
}
```

---

## 5. Main — Wiring It Together

```go
// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"logprocessor/aggregator"
	"logprocessor/parser"
	"logprocessor/worker"
)

func main() {
	var (
		workers  = flag.Int("workers", 4, "number of parallel workers")
		minLevel = flag.String("level", "INFO", "minimum log level to include in output")
		outDir   = flag.String("out", "filtered", "output directory for filtered logs")
	)
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "usage: logprocessor [flags] <file> [file...]")
		fmt.Fprintln(os.Stderr, "  glob patterns are supported: logs/*.log")
		os.Exit(1)
	}

	// Expand globs
	var expanded []string
	for _, pattern := range files {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad pattern %q: %v\n", pattern, err)
			os.Exit(1)
		}
		expanded = append(expanded, matches...)
	}
	if len(expanded) == 0 {
		fmt.Fprintln(os.Stderr, "no files matched")
		os.Exit(1)
	}

	minLvl, err := parser.ParseLevel(*minLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid level %q: %v\n", *minLevel, err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir: %v\n", err)
		os.Exit(1)
	}

	// Graceful shutdown on Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	stats := aggregator.New()
	pool := worker.New(*workers, stats)
	pool.Start(ctx)

	start := time.Now()
	fmt.Printf("Processing %d file(s) with %d workers (min level: %s)...\n",
		len(expanded), *workers, minLvl)

	// Submit jobs
	go func() {
		for _, path := range expanded {
			outPath := filepath.Join(*outDir, filepath.Base(path))
			pool.Submit(worker.Job{
				FilePath:   path,
				MinLevel:   minLvl,
				OutputPath: outPath,
			})
		}
		pool.Close()
	}()

	// Collect results
	var totalLines, totalErrors, failedFiles int
	for result := range pool.Results() {
		if result.Err != nil {
			fmt.Fprintf(os.Stderr, "  ERROR %s: %v\n", result.FilePath, result.Err)
			failedFiles++
			continue
		}
		fmt.Printf("  ✓ %-40s  %6d lines  %d parse errors\n",
			filepath.Base(result.FilePath), result.Lines, result.Errors)
		totalLines += result.Lines
		totalErrors += result.Errors
	}

	if ctx.Err() != nil {
		fmt.Println("\n[interrupted]")
	}

	elapsed := time.Since(start)
	fmt.Printf("\nProcessed %d lines in %s (%.0f lines/sec)\n",
		totalLines, elapsed.Round(time.Millisecond), float64(totalLines)/elapsed.Seconds())
	if failedFiles > 0 {
		fmt.Printf("Failed files: %d\n", failedFiles)
	}

	stats.Print()
}
```

---

## 6. Putting It Together

```bash
# Initialize module
go mod init logprocessor

# Generate test data: 3 files × 10,000 lines each
go run testdata/gen.go logs/app1.log 10000
go run testdata/gen.go logs/app2.log 10000
go run testdata/gen.go logs/app3.log 10000

# Process with 4 workers, output only WARN and above
go run . -workers 4 -level WARN -out filtered/ logs/*.log

# Sample output:
# Processing 3 file(s) with 4 workers (min level: WARN)...
#   ✓ app1.log                                   10000 lines  0 parse errors
#   ✓ app2.log                                   10000 lines  0 parse errors
#   ✓ app3.log                                   10000 lines  0 parse errors
#
# Processed 30000 lines in 98ms (306122 lines/sec)
#
# === Log Processing Results ===
# Total lines:  30000
# Parse errors: 0
#
# Counts by level:
#   DEBUG  2499
#   INFO   21245
#   WARN   3015
#   ERROR  3241
#
# Top errors:
#   [654] connection refused
#   [651] context deadline exceeded
#   [652] database query failed
#   ...
```

---

## 7. Tests

```go
// parser/parser_test.go
package parser_test

import (
	"testing"
	"time"

	"logprocessor/parser"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantLvl parser.Level
		wantMsg string
		wantFields map[string]string
		wantErr bool
	}{
		{
			name:    "info with fields",
			line:    "2024-01-15T14:23:05Z INFO request completed method=GET path=/api/users status=200",
			wantLvl: parser.INFO,
			wantMsg: "request completed",
			wantFields: map[string]string{
				"method": "GET",
				"path":   "/api/users",
				"status": "200",
			},
		},
		{
			name:    "error with quoted field",
			line:    `2024-01-15T14:23:06Z ERROR database failed error="connection refused"`,
			wantLvl: parser.ERROR,
			wantMsg: "database failed",
			wantFields: map[string]string{"error": "connection refused"},
		},
		{
			name:    "bad timestamp",
			line:    "not-a-time INFO something",
			wantErr: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: false, // returns nil, nil
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := parser.ParseLine(tc.line, "test.log")
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.line == "" {
				if entry != nil {
					t.Error("expected nil entry for empty line")
				}
				return
			}
			if entry.Level != tc.wantLvl {
				t.Errorf("level: got %v, want %v", entry.Level, tc.wantLvl)
			}
			if entry.Message != tc.wantMsg {
				t.Errorf("message: got %q, want %q", entry.Message, tc.wantMsg)
			}
			for k, want := range tc.wantFields {
				if got := entry.Fields[k]; got != want {
					t.Errorf("field %q: got %q, want %q", k, got, want)
				}
			}
			_ = time.Time{} // ensure time package used
		})
	}
}
```

```go
// worker/pool_test.go
package worker_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"logprocessor/aggregator"
	"logprocessor/parser"
	"logprocessor/worker"
)

func TestPool(t *testing.T) {
	// Write a temp log file
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")
	content := `2024-01-15T14:23:05Z INFO request ok path=/health
2024-01-15T14:23:06Z ERROR db failed error="timeout"
2024-01-15T14:23:07Z WARN high latency latency=900ms
`
	if err := os.WriteFile(logFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	stats := aggregator.New()
	pool := worker.New(2, stats)
	pool.Start(context.Background())

	outFile := filepath.Join(dir, "out.log")
	pool.Submit(worker.Job{
		FilePath:   logFile,
		MinLevel:   parser.WARN,
		OutputPath: outFile,
	})
	pool.Close()

	var results []worker.Result
	for r := range pool.Results() {
		results = append(results, r)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Lines != 3 {
		t.Errorf("lines: got %d, want 3", results[0].Lines)
	}

	// Output file should have WARN and ERROR only (2 lines)
	out, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	lines := 0
	for _, c := range out {
		if c == '\n' {
			lines++
		}
	}
	if lines != 2 {
		t.Errorf("output lines: got %d, want 2\n%s", lines, out)
	}
}
```

---

## 8. Key Concepts Applied

| Concept (Vol 2) | Where Used |
|-----------------|------------|
| Goroutines | One goroutine per worker; job submitter runs in its own goroutine |
| Channels | `jobs chan Job` (work queue), `results chan Result` (output collection) |
| `sync.WaitGroup` | Wait for all workers before closing `results` |
| `sync.Mutex` | `Stats.mu` protects counters from concurrent writes |
| `select` with context | `work()` cancels cleanly when `ctx.Done()` fires |
| Context cancellation | `signal.NotifyContext` → `Cancel()` → worker `select` sees `ctx.Done()` |
| Buffered channels | `jobs` and `results` are buffered to decouple producer/consumer speed |

**Channel patterns used:**
```
main goroutine            workers (×N)              collector
─────────────────         ─────────────             ─────────────
Submit(job) → jobs ──→   processFile()  → results ──→  range results
Close()    → jobs      (closed when wg done → results closed)
```

---

## Summary

- Worker pool pattern: fixed number of goroutines reading from a shared `jobs` channel
- Thread-safe stats: `sync.Mutex` wraps all write operations on shared counters
- Context propagation: every blocking operation checks `ctx.Done()` for graceful cancellation
- Buffered channels as queues: decouple the submitter speed from worker speed
- `signal.NotifyContext` is the idiomatic way to handle OS signals in Go 1.20+

## Exercises

### Easy
1. Add a `--dry-run` flag that processes files and prints stats but writes no output files.
2. Add a progress bar: print a line like `[2/5 files done]` after each file completes.
3. Add a `--since` flag that skips log entries older than a given duration (e.g., `--since 24h`).

### Medium
4. Add a `--format json` flag that outputs entries as JSON instead of raw log lines.
5. Implement a **fan-in merger**: instead of one output file per input file, merge all filtered entries into a single chronologically sorted output file. Use a heap to merge N sorted streams.
6. Add **metrics per worker**: track how many lines and parse errors each worker processed. Print a per-worker breakdown at the end.

### Hard
7. Replace the file-per-goroutine model with a **line-level pipeline**: a reader goroutine sends lines to a `parse` stage, which sends parsed entries to a `filter` stage, which sends matching entries to a `write` stage. Measure whether the pipeline or the file-per-worker approach is faster on large files. Explain why.
8. Add **tail mode** (`--tail`): after processing existing files, watch them for new lines using `fsnotify` and process new entries in real-time. Print a summary every 5 seconds showing the rate of entries per level.
