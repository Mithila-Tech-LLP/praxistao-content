# Chapter 27: Web Fuzzer and Directory Bruter — Finding Hidden Attack Surface

*Fuzzing is the art of bombarding inputs with unexpected values to discover hidden endpoints, files, parameters, and vulnerabilities. A thorough fuzz often reveals more attack surface than any manual testing.*

---

## What is Web Fuzzing?

Fuzzing sends large volumes of requests with varying inputs to discover:
- Hidden directories and files (`/admin`, `/backup`, `/.git`)
- Parameters a page accepts (undocumented APIs)
- Values that cause unexpected behavior (SQLi, XSS, crashes)
- Virtual hosts (`Host: admin.example.com`)

---

## Directory and File Discovery

```bash
# gobuster — Go-based, fast
gobuster dir -u http://target.com \
    -w /usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt \
    -x php,html,txt,bak,zip \
    -t 50 \
    -o results.txt

# ffuf — highly flexible
ffuf -w /usr/share/wordlists/dirb/common.txt \
    -u http://target.com/FUZZ \
    -mc 200,301,302,403 \
    -o results.json

# feroxbuster — recursive
feroxbuster -u http://target.com -w wordlist.txt --depth 3

# dirb — simpler
dirb http://target.com /usr/share/wordlists/dirb/common.txt
```

### Key File Targets

```
/.git/          → Source code (git repo exposed!)
/.env           → Environment variables (database passwords!)
/backup.zip     → Full site backup
/wp-admin/      → WordPress admin
/phpinfo.php    → PHP configuration (sensitive info)
/admin/         → Admin panel
/api/           → API endpoints
/swagger.json   → API documentation
/.htaccess      → Apache config
/robots.txt     → Reveals hidden paths (always check!)
/sitemap.xml    → Site structure
```

---

## Parameter Fuzzing

```bash
# ffuf for parameter discovery
ffuf -w params.txt -u "http://target.com/page.php?FUZZ=test" \
    -mc 200 -fs 1234  # -fs: filter by size (skip same-size responses)

# Arjun — smart parameter discoverer
arjun -u http://target.com/page.php

# x8 — parameter discovery
x8 -u "http://target.com/page.php" -w params.txt
```

---

## Virtual Host Fuzzing

Different hostnames on the same IP can reveal separate applications:

```bash
# ffuf for vhost discovery
ffuf -w subdomains.txt \
    -u http://target.com \
    -H "Host: FUZZ.target.com" \
    -mc 200,301,302 \
    -fs 1234  # filter default response size
```

---

## Go: Full-Featured Web Fuzzer

```go
package main

import (
    "bufio"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type FuzzConfig struct {
    URL       string
    Wordlist  string
    Workers   int
    Timeout   time.Duration
    Method    string
    Headers   map[string]string
    FilterCodes []int
    FilterSize  int
    Extensions []string
}

type FuzzResult struct {
    URL        string
    StatusCode int
    Length     int64
    Duration   time.Duration
    Words      int
}

func newHTTPClient(timeout time.Duration) *http.Client {
    return &http.Client{
        Timeout: timeout,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }
}

func fuzzWorker(
    client *http.Client,
    config *FuzzConfig,
    jobs <-chan string,
    results chan<- FuzzResult,
    wg *sync.WaitGroup,
    counter *atomic.Int64,
) {
    defer wg.Done()
    
    for word := range jobs {
        targets := []string{word}
        for _, ext := range config.Extensions {
            targets = append(targets, word+"."+ext)
        }
        
        for _, target := range targets {
            targetURL := strings.Replace(config.URL, "FUZZ", target, 1)
            
            start := time.Now()
            req, err := http.NewRequest(config.Method, targetURL, nil)
            if err != nil {
                continue
            }
            
            for k, v := range config.Headers {
                req.Header.Set(k, v)
            }
            
            resp, err := client.Do(req)
            if err != nil {
                counter.Add(1)
                continue
            }
            resp.Body.Close()
            
            dur := time.Since(start)
            counter.Add(1)
            
            // Apply filters
            if config.FilterSize > 0 && int(resp.ContentLength) == config.FilterSize {
                continue
            }
            if len(config.FilterCodes) > 0 {
                filtered := false
                for _, code := range config.FilterCodes {
                    if resp.StatusCode == code {
                        filtered = true
                        break
                    }
                }
                if filtered {
                    continue
                }
            }
            
            // Only show interesting status codes
            if resp.StatusCode != 404 {
                results <- FuzzResult{
                    URL:        targetURL,
                    StatusCode: resp.StatusCode,
                    Length:     resp.ContentLength,
                    Duration:   dur,
                }
            }
        }
    }
}

func run(config *FuzzConfig) {
    client := newHTTPClient(config.Timeout)
    
    f, err := os.Open(config.Wordlist)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Cannot open wordlist: %v\n", err)
        return
    }
    defer f.Close()
    
    jobs    := make(chan string, config.Workers)
    results := make(chan FuzzResult, config.Workers)
    var wg sync.WaitGroup
    var counter atomic.Int64
    
    // Progress reporter
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            fmt.Fprintf(os.Stderr, "\rRequests: %d", counter.Load())
        }
    }()
    
    for i := 0; i < config.Workers; i++ {
        wg.Add(1)
        go fuzzWorker(client, config, jobs, results, &wg, &counter)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Feed wordlist
    go func() {
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            word := strings.TrimSpace(scanner.Text())
            if word != "" && !strings.HasPrefix(word, "#") {
                jobs <- word
            }
        }
        close(jobs)
    }()
    
    // Print header
    fmt.Printf("\n%-60s %-6s %-10s %s\n", "URL", "Status", "Length", "Duration")
    fmt.Println(strings.Repeat("─", 90))
    
    for r := range results {
        statusColor := ""
        if r.StatusCode >= 200 && r.StatusCode < 300 {
            statusColor = "✓"
        } else if r.StatusCode >= 300 && r.StatusCode < 400 {
            statusColor = "→"
        } else if r.StatusCode == 403 {
            statusColor = "✗"
        }
        fmt.Printf("%-60s %s%-6d %-10d %s\n",
            r.URL, statusColor, r.StatusCode, r.Length, r.Duration.Round(time.Millisecond))
    }
    
    fmt.Printf("\nTotal requests: %d\n", counter.Load())
}

func main() {
    var config FuzzConfig
    var filterCodes string
    var extensions string
    var headers string
    
    flag.StringVar(&config.URL, "u", "", "URL with FUZZ placeholder")
    flag.StringVar(&config.Wordlist, "w", "", "Wordlist path")
    flag.IntVar(&config.Workers, "t", 50, "Workers")
    flag.DurationVar(&config.Timeout, "timeout", 5*time.Second, "Request timeout")
    flag.StringVar(&config.Method, "X", "GET", "HTTP method")
    flag.StringVar(&filterCodes, "fc", "", "Filter status codes (e.g. 404,400)")
    flag.IntVar(&config.FilterSize, "fs", 0, "Filter by response size")
    flag.StringVar(&extensions, "e", "", "Extensions (e.g. php,html,txt)")
    flag.StringVar(&headers, "H", "", "Custom header (e.g. 'Cookie: session=abc')")
    flag.Parse()
    
    if config.URL == "" || config.Wordlist == "" {
        fmt.Println("Usage: gofuzz -u http://target.com/FUZZ -w wordlist.txt")
        flag.PrintDefaults()
        return
    }
    
    // Parse headers
    config.Headers = make(map[string]string)
    if headers != "" {
        parts := strings.SplitN(headers, ":", 2)
        if len(parts) == 2 {
            config.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
        }
    }
    
    // Parse filter codes
    if filterCodes != "" {
        for _, code := range strings.Split(filterCodes, ",") {
            var c int
            fmt.Sscanf(strings.TrimSpace(code), "%d", &c)
            config.FilterCodes = append(config.FilterCodes, c)
        }
    }
    
    // Parse extensions
    if extensions != "" {
        config.Extensions = strings.Split(extensions, ",")
    }
    
    fmt.Printf("Fuzzing: %s\n", config.URL)
    run(&config)
}
```

```bash
# Usage:
go run fuzzer.go -u http://192.168.1.100/FUZZ -w /usr/share/wordlists/dirb/common.txt -t 100 -e php,html,txt

# With custom header (authenticated fuzzing)
go run fuzzer.go -u http://target.com/api/FUZZ -w api-endpoints.txt \
    -H "Authorization: Bearer TOKEN"

# Filter out 403 responses
go run fuzzer.go -u http://target.com/FUZZ -w wordlist.txt -fc 403,404
```

---

## Summary

| Tool | Best for | Speed |
|------|---------|-------|
| `gobuster` | Directory/file discovery | Fast |
| `ffuf` | Flexible fuzzing (params, vhosts) | Very fast |
| `feroxbuster` | Recursive discovery | Fast |
| `arjun` | Parameter discovery | Medium |
| Custom Go fuzzer | Full control, custom logic | Fast |

---

## Exercises

1. Fuzz DVWA with gobuster and find all hidden directories
2. Check for exposed `.git` directories on test servers — if found, use `git-dumper` to extract the source
3. Extend the Go fuzzer to support POST body fuzzing (send different values as JSON body)
4. Build a recursive fuzzer: when a directory is found, automatically fuzz inside it too
