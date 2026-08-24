# Chapter 64: GoShield Server — Event Collection, Storage, and API

*The server is the brain of GoShield. It receives events from thousands of agents, stores them efficiently, runs detection, and serves the analyst dashboard. This chapter builds the complete server.*

---

## Server Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    GoShield Server                          │
│                                                             │
│   Agent requests                  Analyst requests         │
│        │                                │                  │
│        ▼                                ▼                  │
│   ┌────────────┐               ┌────────────────┐          │
│   │ /api/events│               │ /api/v1/alerts │          │
│   │ /api/health│               │ /api/v1/hunt   │          │
│   │ /api/status│               │ /api/v1/agents │          │
│   └─────┬──────┘               └───────┬────────┘          │
│         │                              │                   │
│         ▼                              │                   │
│   ┌──────────────┐                     │                   │
│   │  Ingestion   │                     │                   │
│   │  Pipeline    │                     │                   │
│   └──────┬───────┘                     │                   │
│          │                             │                   │
│          ▼                             │                   │
│   ┌──────────────┐              ┌──────▼──────────┐        │
│   │  Detection   │──── alerts ─▶│    Storage      │        │
│   │  Engine      │              │  (SQLite / PG)  │        │
│   └──────────────┘              └─────────────────┘        │
│          │                                                  │
│          ▼                                                  │
│   ┌──────────────┐                                         │
│   │  Alert       │──── webhooks ──▶ Slack / PagerDuty     │
│   │  Manager     │──── email ─────▶ security@company.com  │
│   └──────────────┘                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## Storage Layer — SQLite Schema

```go
// pkg/server/storage.go
package server

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
    "github.com/yourname/goshield/pkg/events"
)

const schema = `
-- Agents registry
CREATE TABLE IF NOT EXISTS agents (
    id          TEXT PRIMARY KEY,
    hostname    TEXT NOT NULL,
    os          TEXT,
    version     TEXT,
    last_seen   DATETIME,
    online      BOOLEAN DEFAULT false,
    events_total INTEGER DEFAULT 0
);

-- All events (partitioned by type for query performance)
CREATE TABLE IF NOT EXISTS process_events (
    id          TEXT PRIMARY KEY,
    agent_id    TEXT,
    hostname    TEXT,
    timestamp   DATETIME,
    action      TEXT,
    pid         INTEGER,
    ppid        INTEGER,
    name        TEXT,
    command_line TEXT,
    username    TEXT,
    exe_path    TEXT,
    sha256      TEXT,
    is_elevated BOOLEAN,
    raw         TEXT  -- full JSON for flexible querying
);

CREATE TABLE IF NOT EXISTS file_events (
    id          TEXT PRIMARY KEY,
    agent_id    TEXT,
    hostname    TEXT,
    timestamp   DATETIME,
    action      TEXT,
    path        TEXT,
    new_path    TEXT,
    sha256      TEXT,
    size        INTEGER,
    pid         INTEGER,
    process     TEXT,
    extension   TEXT,
    raw         TEXT
);

CREATE TABLE IF NOT EXISTS network_events (
    id          TEXT PRIMARY KEY,
    agent_id    TEXT,
    hostname    TEXT,
    timestamp   DATETIME,
    action      TEXT,
    protocol    TEXT,
    src_ip      TEXT,
    src_port    INTEGER,
    dst_ip      TEXT,
    dst_port    INTEGER,
    pid         INTEGER,
    process     TEXT,
    domain      TEXT,
    raw         TEXT
);

CREATE TABLE IF NOT EXISTS alerts (
    id          TEXT PRIMARY KEY,
    agent_id    TEXT,
    hostname    TEXT,
    timestamp   DATETIME,
    rule_id     TEXT,
    rule_name   TEXT,
    severity    TEXT,
    description TEXT,
    event_id    TEXT,
    event_type  TEXT,
    mitre       TEXT,
    resolved    BOOLEAN DEFAULT false,
    resolved_at DATETIME,
    notes       TEXT
);

-- Indices for common queries
CREATE INDEX IF NOT EXISTS idx_proc_events_time ON process_events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_proc_events_agent ON process_events(agent_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_proc_events_name ON process_events(name, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_file_events_time ON file_events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_file_events_path ON file_events(path, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_net_events_time ON network_events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_net_events_dst ON network_events(dst_ip, dst_port, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_alerts_time ON alerts(timestamp DESC, resolved);
CREATE INDEX IF NOT EXISTS idx_alerts_severity ON alerts(severity, resolved, timestamp DESC);
`

// DB wraps the database connection
type DB struct {
    conn *sql.DB
}

// NewDB initializes the SQLite database
func NewDB(path string) (*DB, error) {
    conn, err := sql.Open("sqlite3", path+"?_journal=WAL&_timeout=5000")
    if err != nil {
        return nil, fmt.Errorf("open db: %w", err)
    }
    
    // WAL mode for concurrent reads + writes
    conn.SetMaxOpenConns(1)  // SQLite only supports one writer
    
    if _, err := conn.Exec(schema); err != nil {
        return nil, fmt.Errorf("create schema: %w", err)
    }
    
    return &DB{conn: conn}, nil
}

// StoreEvent stores any event type
func (db *DB) StoreEvent(event interface{}) error {
    switch e := event.(type) {
    case *events.ProcessEvent:
        return db.storeProcessEvent(e)
    case *events.FileEvent:
        return db.storeFileEvent(e)
    case *events.NetworkEvent:
        return db.storeNetworkEvent(e)
    case *events.Alert:
        return db.storeAlert(e)
    default:
        return fmt.Errorf("unknown event type: %T", event)
    }
}

func (db *DB) storeProcessEvent(e *events.ProcessEvent) error {
    raw, _ := json.Marshal(e)
    _, err := db.conn.Exec(`
        INSERT OR IGNORE INTO process_events 
        (id, agent_id, hostname, timestamp, action, pid, ppid, name, command_line, username, exe_path, sha256, is_elevated, raw)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        e.ID, e.AgentID, e.Hostname, e.Timestamp,
        e.Action, e.PID, e.PPID, e.Name, e.CommandLine,
        e.Username, e.ExePath, e.SHA256, e.IsElevated, string(raw),
    )
    return err
}

func (db *DB) storeFileEvent(e *events.FileEvent) error {
    raw, _ := json.Marshal(e)
    _, err := db.conn.Exec(`
        INSERT OR IGNORE INTO file_events
        (id, agent_id, hostname, timestamp, action, path, new_path, sha256, size, pid, process, extension, raw)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        e.ID, e.AgentID, e.Hostname, e.Timestamp,
        e.Action, e.Path, e.NewPath, e.SHA256,
        e.Size, e.PID, e.Process, e.Extension, string(raw),
    )
    return err
}

func (db *DB) storeNetworkEvent(e *events.NetworkEvent) error {
    raw, _ := json.Marshal(e)
    _, err := db.conn.Exec(`
        INSERT OR IGNORE INTO network_events
        (id, agent_id, hostname, timestamp, action, protocol, src_ip, src_port, dst_ip, dst_port, pid, process, domain, raw)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        e.ID, e.AgentID, e.Hostname, e.Timestamp,
        e.Action, e.Protocol, e.SrcIP, e.SrcPort,
        e.DstIP, e.DstPort, e.PID, e.Process, e.Domain, string(raw),
    )
    return err
}

func (db *DB) storeAlert(a *events.Alert) error {
    raw, _ := json.Marshal(a)
    _, err := db.conn.Exec(`
        INSERT OR IGNORE INTO alerts
        (id, agent_id, hostname, timestamp, rule_id, rule_name, severity, description, event_id, event_type, mitre, resolved)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        a.ID, a.AgentID, a.Hostname, a.Timestamp,
        a.RuleID, a.RuleName, a.Severity, a.Description,
        a.EventID, a.EventType, a.MITRE, false,
    )
    _ = raw
    return err
}

// QueryResult for search/hunt
type QueryResult struct {
    Total  int           `json:"total"`
    Events []interface{} `json:"events"`
}

// SearchProcessEvents queries process events
func (db *DB) SearchProcessEvents(agentID, nameFilter, cmdFilter string, since time.Time, limit int) ([]*events.ProcessEvent, error) {
    query := `SELECT raw FROM process_events WHERE timestamp > ?`
    args := []interface{}{since}
    
    if agentID != "" {
        query += " AND agent_id = ?"
        args = append(args, agentID)
    }
    if nameFilter != "" {
        query += " AND name LIKE ?"
        args = append(args, "%"+nameFilter+"%")
    }
    if cmdFilter != "" {
        query += " AND command_line LIKE ?"
        args = append(args, "%"+cmdFilter+"%")
    }
    query += " ORDER BY timestamp DESC LIMIT ?"
    args = append(args, limit)
    
    rows, err := db.conn.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []*events.ProcessEvent
    for rows.Next() {
        var raw string
        if err := rows.Scan(&raw); err != nil {
            continue
        }
        var e events.ProcessEvent
        if err := json.Unmarshal([]byte(raw), &e); err == nil {
            results = append(results, &e)
        }
    }
    return results, nil
}

// GetAlerts returns active alerts
func (db *DB) GetAlerts(resolved bool, severity string, limit int) ([]*events.Alert, error) {
    query := `SELECT id, agent_id, hostname, timestamp, rule_id, rule_name, severity, 
              description, event_id, event_type, mitre, resolved
              FROM alerts WHERE resolved = ?`
    args := []interface{}{resolved}
    
    if severity != "" {
        query += " AND severity = ?"
        args = append(args, severity)
    }
    query += " ORDER BY timestamp DESC LIMIT ?"
    args = append(args, limit)
    
    rows, err := db.conn.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var alerts []*events.Alert
    for rows.Next() {
        a := &events.Alert{}
        err := rows.Scan(&a.ID, &a.AgentID, &a.Hostname, &a.Timestamp,
            &a.RuleID, &a.RuleName, &a.Severity, &a.Description,
            &a.EventID, &a.EventType, &a.MITRE, &a.Resolved)
        if err == nil {
            alerts = append(alerts, a)
        }
    }
    return alerts, nil
}

// PruneOldEvents deletes events older than retentionDays
func (db *DB) PruneOldEvents(retentionDays int) error {
    cutoff := time.Now().AddDate(0, 0, -retentionDays)
    tables := []string{"process_events", "file_events", "network_events"}
    for _, table := range tables {
        _, err := db.conn.Exec(fmt.Sprintf("DELETE FROM %s WHERE timestamp < ?", table), cutoff)
        if err != nil {
            return err
        }
    }
    return nil
}
```

---

## HTTP API Server

```go
// pkg/server/api.go
package server

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/yourname/goshield/pkg/events"
)

// Server is the GoShield API server
type Server struct {
    db       *DB
    detector *Detector  // We'll build this in Ch 65
    alerter  *Alerter   // We'll build this in Ch 66
    apiKey   string
    router   *mux.Router
}

// NewServer creates an API server
func NewServer(db *DB, detector *Detector, alerter *Alerter, apiKey string) *Server {
    s := &Server{
        db:       db,
        detector: detector,
        alerter:  alerter,
        apiKey:   apiKey,
        router:   mux.NewRouter(),
    }
    s.registerRoutes()
    return s
}

func (s *Server) registerRoutes() {
    // Agent endpoints (authenticated by API key)
    agent := s.router.PathPrefix("/api/agent").Subrouter()
    agent.Use(s.agentAuthMiddleware)
    agent.HandleFunc("/events", s.handleIngestEvents).Methods("POST")
    agent.HandleFunc("/health", s.handleAgentHealth).Methods("POST")
    
    // Analyst endpoints (also authenticated)
    api := s.router.PathPrefix("/api/v1").Subrouter()
    api.Use(s.analystAuthMiddleware)
    api.HandleFunc("/alerts", s.handleGetAlerts).Methods("GET")
    api.HandleFunc("/alerts/{id}/resolve", s.handleResolveAlert).Methods("POST")
    api.HandleFunc("/events/process", s.handleSearchProcessEvents).Methods("GET")
    api.HandleFunc("/events/file", s.handleSearchFileEvents).Methods("GET")
    api.HandleFunc("/events/network", s.handleSearchNetworkEvents).Methods("GET")
    api.HandleFunc("/agents", s.handleGetAgents).Methods("GET")
    api.HandleFunc("/stats", s.handleGetStats).Methods("GET")
    
    // Dashboard (static files)
    s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
}

// agentAuthMiddleware validates agent API key
func (s *Server) agentAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := r.Header.Get("X-API-Key")
        if key != s.apiKey {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// analystAuthMiddleware — same for now, could add JWT later
func (s *Server) analystAuthMiddleware(next http.Handler) http.Handler {
    return s.agentAuthMiddleware(next)
}

// handleIngestEvents receives events from agents
func (s *Server) handleIngestEvents(w http.ResponseWriter, r *http.Request) {
    // Events are sent as a JSON array of raw event objects
    var rawEvents []json.RawMessage
    if err := json.NewDecoder(r.Body).Decode(&rawEvents); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    stored := 0
    for _, raw := range rawEvents {
        // Peek at the "type" field to determine event type
        var base events.Base
        if err := json.Unmarshal(raw, &base); err != nil {
            continue
        }
        
        var event interface{}
        switch base.Type {
        case events.EventTypeProcess:
            var e events.ProcessEvent
            if err := json.Unmarshal(raw, &e); err == nil {
                event = &e
            }
        case events.EventTypeFile:
            var e events.FileEvent
            if err := json.Unmarshal(raw, &e); err == nil {
                event = &e
            }
        case events.EventTypeNetwork:
            var e events.NetworkEvent
            if err := json.Unmarshal(raw, &e); err == nil {
                event = &e
            }
        }
        
        if event == nil {
            continue
        }
        
        // Store the event
        if err := s.db.StoreEvent(event); err != nil {
            continue
        }
        stored++
        
        // Run detection on each event
        if alerts := s.detector.Evaluate(event); len(alerts) > 0 {
            for _, alert := range alerts {
                s.db.StoreEvent(alert)
                s.alerter.Send(alert)
            }
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"stored": stored})
}

// handleGetAlerts returns active alerts
func (s *Server) handleGetAlerts(w http.ResponseWriter, r *http.Request) {
    resolved := r.URL.Query().Get("resolved") == "true"
    severity := r.URL.Query().Get("severity")
    limit := 100
    if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
        limit = l
    }
    
    alerts, err := s.db.GetAlerts(resolved, severity, limit)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "alerts": alerts,
        "count":  len(alerts),
    })
}

// handleSearchProcessEvents searches process events
func (s *Server) handleSearchProcessEvents(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    
    since := time.Now().Add(-24 * time.Hour)
    if s := q.Get("since"); s != "" {
        if t, err := time.Parse(time.RFC3339, s); err == nil {
            since = t
        }
    }
    
    limit := 500
    if l, err := strconv.Atoi(q.Get("limit")); err == nil && l > 0 && l <= 10000 {
        limit = l
    }
    
    results, err := s.db.SearchProcessEvents(
        q.Get("agent"),
        q.Get("name"),
        q.Get("cmd"),
        since,
        limit,
    )
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "events": results,
        "count":  len(results),
    })
}

// handleResolveAlert marks an alert as resolved
func (s *Server) handleResolveAlert(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    
    var body struct {
        Notes string `json:"notes"`
    }
    json.NewDecoder(r.Body).Decode(&body)
    
    _, err := s.db.conn.Exec(
        "UPDATE alerts SET resolved=true, resolved_at=?, notes=? WHERE id=?",
        time.Now(), body.Notes, id,
    )
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}

// handleGetStats returns summary statistics
func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
    stats := map[string]interface{}{}
    
    // Count events in last 24h
    since := time.Now().Add(-24 * time.Hour)
    for _, table := range []string{"process_events", "file_events", "network_events"} {
        var count int
        s.db.conn.QueryRow(
            "SELECT COUNT(*) FROM "+table+" WHERE timestamp > ?", since,
        ).Scan(&count)
        stats[table+"_24h"] = count
    }
    
    // Count unresolved alerts by severity
    rows, _ := s.db.conn.Query(
        "SELECT severity, COUNT(*) FROM alerts WHERE resolved=false GROUP BY severity",
    )
    if rows != nil {
        defer rows.Close()
        alertsBySeverity := map[string]int{}
        for rows.Next() {
            var sev string
            var count int
            rows.Scan(&sev, &count)
            alertsBySeverity[sev] = count
        }
        stats["alerts_by_severity"] = alertsBySeverity
    }
    
    // Agent count
    var agentCount int
    s.db.conn.QueryRow("SELECT COUNT(*) FROM agents WHERE last_seen > ?", 
        time.Now().Add(-5*time.Minute)).Scan(&agentCount)
    stats["online_agents"] = agentCount
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

// Handler returns the HTTP handler for the server
func (s *Server) Handler() http.Handler {
    return s.router
}
```

---

## Main Server Entry Point

```go
// cmd/server/main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/yourname/goshield/pkg/server"
)

func main() {
    dbPath := getenv("GOSHIELD_DB", "/var/lib/goshield/events.db")
    apiKey := getenv("GOSHIELD_API_KEY", "changeme-insecure-default")
    listenAddr := getenv("GOSHIELD_LISTEN", ":8080")
    rulesDir := getenv("GOSHIELD_RULES", "./rules")
    
    // Initialize database
    db, err := server.NewDB(dbPath)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    
    // Initialize detection engine
    detector, err := server.NewDetector(rulesDir)
    if err != nil {
        log.Fatalf("Failed to initialize detector: %v", err)
    }
    log.Printf("Loaded %d detection rules", detector.RuleCount())
    
    // Initialize alerter (webhooks, email)
    alerter := server.NewAlerter()
    if webhookURL := os.Getenv("GOSHIELD_WEBHOOK_URL"); webhookURL != "" {
        alerter.AddWebhook(webhookURL, []string{"high", "critical"})
    }
    
    // Create API server
    srv := server.NewServer(db, detector, alerter, apiKey)
    
    httpServer := &http.Server{
        Addr:         listenAddr,
        Handler:      srv.Handler(),
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    
    // Start background tasks
    go func() {
        // Prune old events daily
        ticker := time.NewTicker(24 * time.Hour)
        for range ticker.C {
            if err := db.PruneOldEvents(90); err != nil {
                log.Printf("Prune error: %v", err)
            }
        }
    }()
    
    // Start HTTP server
    go func() {
        fmt.Printf("GoShield Server listening on %s\n", listenAddr)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()
    
    // Wait for shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
    <-sigCh
    
    fmt.Println("Shutting down GoShield server...")
}

func getenv(key, defaultValue string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultValue
}
```

---

## Testing the Server

```bash
# Start the server
GOSHIELD_API_KEY=mysecret go run cmd/server/main.go

# Send a test event
curl -X POST http://localhost:8080/api/agent/events \
  -H "X-API-Key: mysecret" \
  -H "Content-Type: application/json" \
  -d '[{
    "id": "test-001",
    "agent_id": "agent-001",
    "hostname": "webserver-01",
    "timestamp": "2024-01-15T10:30:00Z",
    "type": "process",
    "action": "create",
    "pid": 1234,
    "ppid": 1000,
    "name": "powershell.exe",
    "command_line": "powershell.exe -EncodedCommand aGVsbG8=",
    "username": "www-data",
    "exe_path": "/usr/bin/powershell"
  }]'

# Get alerts
curl -H "X-API-Key: mysecret" http://localhost:8080/api/v1/alerts

# Search process events
curl -H "X-API-Key: mysecret" "http://localhost:8080/api/v1/events/process?name=powershell"
```

---

## Summary

The GoShield server is now complete with:
- SQLite storage with proper indices for fast queries
- REST API for agents (ingest) and analysts (search/alerts)
- Background data retention management
- Detection engine integration (next chapter)
- Alerting integration (next chapter)

---

## Exercises

1. Build and run the server. Send it test events from the command line.
2. Add a `/api/v1/hunt` endpoint that takes a free-form SQL query (be careful: SQL injection!)
3. Add pagination to the alerts endpoint (cursor-based, not offset)
4. Implement agent registration: when an agent first connects, store it in the agents table
5. Add TLS: generate a self-signed cert and make the server accept HTTPS
