# Chapter 50: Mini Project 4 — Full API Service

This mini-project brings together everything from the Web Development volume (Chapters 40–49): HTTP server, REST API, middleware, authentication, WebSockets, SSE, gRPC, OpenAPI, and configuration. You'll build a **collaborative note-taking API** — the kind of backend that powers a simple Notion-like app.

## What You'll Build

A full-featured API service with:
1. **REST API** for CRUD on notes and tags
2. **JWT Authentication** (register, login, refresh, logout)
3. **WebSocket** for real-time collaborative editing
4. **SSE** for presence notifications (who's online, who's editing what)
5. **gRPC** for an internal search service
6. **OpenAPI** spec served at `/docs`
7. **Structured config** loaded from environment

---

## Project Structure

```
notes-api/
├── cmd/
│   ├── api/main.go          — HTTP/WebSocket/SSE server
│   └── search/main.go       — gRPC search service
├── internal/
│   ├── config/config.go
│   ├── domain/              — pure business types
│   │   ├── note.go
│   │   ├── user.go
│   │   └── tag.go
│   ├── store/               — in-memory storage
│   │   ├── notes.go
│   │   ├── users.go
│   │   └── tags.go
│   ├── auth/                — JWT service
│   │   └── jwt.go
│   ├── handlers/            — HTTP handlers
│   │   ├── auth.go
│   │   ├── notes.go
│   │   └── tags.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logging.go
│   │   └── ratelimit.go
│   ├── realtime/            — WebSocket + SSE
│   │   ├── hub.go
│   │   ├── presence.go
│   │   └── ws_handler.go
│   └── search/              — gRPC search client/server
│       ├── proto/search.proto
│       └── service.go
├── openapi.yaml
├── go.mod
└── .env.example
```

---

## domain/note.go

```go
package domain

import "time"

type NoteStatus string

const (
    NoteStatusDraft     NoteStatus = "draft"
    NoteStatusPublished NoteStatus = "published"
    NoteStatusArchived  NoteStatus = "archived"
)

type Note struct {
    ID        int64      `json:"id"`
    Title     string     `json:"title"`
    Content   string     `json:"content"`
    Status    NoteStatus `json:"status"`
    OwnerID   int64      `json:"ownerId"`
    Tags      []Tag      `json:"tags,omitempty"`
    Version   int        `json:"version"`  // Optimistic locking
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
}

type Tag struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
}

// NoteEvent is broadcast over WebSocket/SSE when a note changes.
type NoteEvent struct {
    Type   string `json:"type"`   // "created", "updated", "deleted", "locked", "unlocked"
    NoteID int64  `json:"noteId"`
    UserID int64  `json:"userId"`
    Title  string `json:"title,omitempty"`
}
```

---

## store/notes.go

```go
package store

import (
    "context"
    "errors"
    "sort"
    "sync"
    "time"

    "notes-api/internal/domain"
)

var (
    ErrNotFound    = errors.New("not found")
    ErrVersionConflict = errors.New("version conflict")
    ErrNotOwner    = errors.New("not the owner")
)

type NoteStore struct {
    mu     sync.RWMutex
    notes  map[int64]*domain.Note
    nextID int64
    locks  map[int64]int64  // noteID → userID holding the edit lock
}

func NewNoteStore() *NoteStore {
    return &NoteStore{
        notes:  make(map[int64]*domain.Note),
        locks:  make(map[int64]int64),
        nextID: 1,
    }
}

func (s *NoteStore) Create(ctx context.Context, ownerID int64, title, content string) (*domain.Note, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UTC()
    n := &domain.Note{
        ID:        s.nextID,
        Title:     title,
        Content:   content,
        Status:    domain.NoteStatusDraft,
        OwnerID:   ownerID,
        Version:   1,
        CreatedAt: now,
        UpdatedAt: now,
    }
    s.notes[s.nextID] = n
    s.nextID++
    return copyNote(n), nil
}

func (s *NoteStore) GetByID(ctx context.Context, id, requestingUserID int64) (*domain.Note, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    n, ok := s.notes[id]
    if !ok { return nil, ErrNotFound }
    if n.OwnerID != requestingUserID && n.Status != domain.NoteStatusPublished {
        return nil, ErrNotFound  // Private notes invisible to non-owners
    }
    return copyNote(n), nil
}

func (s *NoteStore) Update(ctx context.Context, id, userID int64, title, content string, version int) (*domain.Note, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    n, ok := s.notes[id]
    if !ok { return nil, ErrNotFound }
    if n.OwnerID != userID { return nil, ErrNotOwner }
    if n.Version != version { return nil, ErrVersionConflict }

    if title != "" { n.Title = title }
    if content != "" { n.Content = content }
    n.Version++
    n.UpdatedAt = time.Now().UTC()
    return copyNote(n), nil
}

func (s *NoteStore) Delete(ctx context.Context, id, userID int64) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    n, ok := s.notes[id]
    if !ok { return ErrNotFound }
    if n.OwnerID != userID { return ErrNotOwner }
    delete(s.notes, id)
    delete(s.locks, id)
    return nil
}

func (s *NoteStore) ListByOwner(ctx context.Context, ownerID int64, offset, limit int) ([]*domain.Note, int) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var owned []*domain.Note
    for _, n := range s.notes {
        if n.OwnerID == ownerID { owned = append(owned, copyNote(n)) }
    }
    sort.Slice(owned, func(i, j int) bool { return owned[i].UpdatedAt.After(owned[j].UpdatedAt) })

    total := len(owned)
    if offset >= total { return []*domain.Note{}, total }
    end := offset + limit
    if end > total { end = total }
    return owned[offset:end], total
}

// AcquireLock: take an exclusive edit lock on a note.
func (s *NoteStore) AcquireLock(noteID, userID int64) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if holder, locked := s.locks[noteID]; locked && holder != userID {
        return fmt.Errorf("note locked by user %d", holder)
    }
    s.locks[noteID] = userID
    return nil
}

// ReleaseLock: release the edit lock.
func (s *NoteStore) ReleaseLock(noteID, userID int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.locks[noteID] == userID { delete(s.locks, noteID) }
}

func copyNote(n *domain.Note) *domain.Note {
    cp := *n
    return &cp
}
```

---

## handlers/notes.go

```go
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "notes-api/internal/domain"
    "notes-api/internal/middleware"
    "notes-api/internal/store"
)

type NoteHandler struct {
    notes  *store.NoteStore
    events chan<- domain.NoteEvent
}

func NewNoteHandler(notes *store.NoteStore, events chan<- domain.NoteEvent) *NoteHandler {
    return &NoteHandler{notes: notes, events: events}
}

func (h *NoteHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Route("/{noteID}", func(r chi.Router) {
        r.Get("/", h.Get)
        r.Patch("/", h.Update)
        r.Delete("/", h.Delete)
        r.Post("/lock", h.AcquireLock)
        r.Delete("/lock", h.ReleaseLock)
        r.Put("/status", h.UpdateStatus)
    })
    return r
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
    claims := middleware.ClaimsFromCtx(r.Context())

    var req struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := decodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    if req.Title == "" {
        writeError(w, http.StatusUnprocessableEntity, "title is required")
        return
    }

    note, err := h.notes.Create(r.Context(), claims.UserID, req.Title, req.Content)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to create note")
        return
    }

    // Broadcast event asynchronously:
    select {
    case h.events <- domain.NoteEvent{Type: "created", NoteID: note.ID, UserID: claims.UserID, Title: note.Title}:
    default:
    }

    w.Header().Set("Location", "/api/v1/notes/"+strconv.FormatInt(note.ID, 10))
    writeJSON(w, http.StatusCreated, map[string]any{"data": note})
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
    claims := middleware.ClaimsFromCtx(r.Context())
    noteID, err := parseID(r, "noteID")
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid note ID")
        return
    }

    var req struct {
        Title   string `json:"title"`
        Content string `json:"content"`
        Version int    `json:"version"`
    }
    if err := decodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    note, err := h.notes.Update(r.Context(), noteID, claims.UserID, req.Title, req.Content, req.Version)
    if err != nil {
        switch {
        case errors.Is(err, store.ErrNotFound):
            writeError(w, http.StatusNotFound, "note not found")
        case errors.Is(err, store.ErrNotOwner):
            writeError(w, http.StatusForbidden, "not the note owner")
        case errors.Is(err, store.ErrVersionConflict):
            writeJSON(w, http.StatusConflict, map[string]any{
                "error":          "version conflict",
                "currentVersion": note,
            })
        default:
            writeError(w, http.StatusInternalServerError, "failed to update note")
        }
        return
    }

    select {
    case h.events <- domain.NoteEvent{Type: "updated", NoteID: note.ID, UserID: claims.UserID}:
    default:
    }

    writeJSON(w, http.StatusOK, map[string]any{"data": note})
}

// AcquireLock takes an exclusive edit lock, preventing others from editing simultaneously.
func (h *NoteHandler) AcquireLock(w http.ResponseWriter, r *http.Request) {
    claims := middleware.ClaimsFromCtx(r.Context())
    noteID, _ := parseID(r, "noteID")

    if err := h.notes.AcquireLock(noteID, claims.UserID); err != nil {
        writeError(w, http.StatusConflict, err.Error())
        return
    }

    select {
    case h.events <- domain.NoteEvent{Type: "locked", NoteID: noteID, UserID: claims.UserID}:
    default:
    }

    w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request, param string) (int64, error) {
    return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}
```

---

## realtime/hub.go

```go
package realtime

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"

    "notes-api/internal/domain"
    "notes-api/internal/middleware"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
    clients    map[int64]*wsClient  // userID → client
    noteEvents <-chan domain.NoteEvent
    mu         sync.RWMutex
}

type wsClient struct {
    userID int64
    send   chan []byte
    conn   *websocket.Conn
}

func NewHub(events <-chan domain.NoteEvent) *Hub {
    return &Hub{
        clients:    make(map[int64]*wsClient),
        noteEvents: events,
    }
}

func (h *Hub) Run() {
    for event := range h.noteEvents {
        data, _ := json.Marshal(event)
        h.mu.RLock()
        for _, client := range h.clients {
            select {
            case client.send <- data:
            default:
            }
        }
        h.mu.RUnlock()
    }
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
    claims := middleware.ClaimsFromCtx(r.Context())

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil { return }

    client := &wsClient{
        userID: claims.UserID,
        send:   make(chan []byte, 64),
        conn:   conn,
    }

    h.mu.Lock()
    h.clients[claims.UserID] = client
    h.mu.Unlock()

    defer func() {
        h.mu.Lock()
        delete(h.clients, claims.UserID)
        h.mu.Unlock()
        close(client.send)
        conn.Close()
    }()

    // Write pump:
    go func() {
        ticker := time.NewTicker(54 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case msg := <-client.send:
                conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
                conn.WriteMessage(websocket.TextMessage, msg)
            case <-ticker.C:
                conn.WriteMessage(websocket.PingMessage, nil)
            }
        }
    }()

    // Read pump (discard — clients only receive from this endpoint):
    conn.SetReadLimit(512)
    conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    conn.SetPongHandler(func(string) error {
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    for {
        if _, _, err := conn.ReadMessage(); err != nil { return }
    }
}
```

---

## cmd/api/main.go

```go
package main

import (
    "context"
    "log"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"

    "notes-api/internal/auth"
    "notes-api/internal/config"
    "notes-api/internal/domain"
    "notes-api/internal/handlers"
    "notes-api/internal/middleware"
    "notes-api/internal/realtime"
    "notes-api/internal/store"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: parseLevel(cfg.Log.Level),
    }))
    slog.SetDefault(logger)

    // Stores:
    userStore := store.NewUserStore()
    noteStore := store.NewNoteStore()
    tagStore := store.NewTagStore()

    // Auth:
    jwtSvc := auth.NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer)
    authSvc := auth.NewAuthService(jwtSvc, userStore)

    // Real-time:
    events := make(chan domain.NoteEvent, 256)
    hub := realtime.NewHub(events)
    sseBroker := realtime.NewSSEBroker(events)
    go hub.Run()
    go sseBroker.Run()

    // Handlers:
    authHandler := handlers.NewAuthHandler(authSvc)
    noteHandler := handlers.NewNoteHandler(noteStore, events)
    tagHandler := handlers.NewTagHandler(tagStore)

    // Router:
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger(logger))
    r.Use(middleware.Recover(logger))
    r.Use(middleware.CORS(middleware.DefaultCORS))

    // Public:
    r.Post("/auth/register", authHandler.Register)
    r.Post("/auth/login", authHandler.Login)
    r.Post("/auth/refresh", authHandler.Refresh)
    r.Get("/health", healthHandler)
    r.Get("/docs", docsHandler)
    r.Get("/openapi.yaml", openapiHandler)

    // Protected:
    r.Group(func(r chi.Router) {
        r.Use(middleware.Authenticate(jwtSvc.ValidateToken))
        r.Delete("/auth/logout", authHandler.Logout)
        r.Mount("/api/v1/notes", noteHandler.Routes())
        r.Mount("/api/v1/tags", tagHandler.Routes())
        r.Get("/ws", hub.ServeWS)
        r.Get("/events", sseBroker.Handler)
    })

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:      r,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
    }

    go func() {
        logger.Info("server starting", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("graceful shutdown: %v", err)
    }
    logger.Info("server stopped")
}
```

---

## .env.example

```bash
# Application
APP_ENV=development
PORT=8080
LOG_LEVEL=debug

# Auth
JWT_SECRET=dev-secret-key-minimum-32-chars-!!
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h

# Database (used in later chapters when we add PostgreSQL)
DATABASE_URL=postgres://user:pass@localhost:5432/notesdb?sslmode=disable

# Feature flags
FEATURE_COLLABORATIVE_EDITING=true
FEATURE_PUBLIC_NOTES=false
```

---

## Running and Testing

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api

# Register:
curl -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@test.com","password":"secret123"}'

# Login:
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@test.com","password":"secret123"}' | jq -r '.accessToken')

# Create note:
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"My First Note","content":"Hello world!"}'

# Stream events (in another terminal):
curl -N -H "Authorization: Bearer $TOKEN" http://localhost:8080/events

# Connect via WebSocket:
# wscat -c "ws://localhost:8080/ws" -H "Authorization: Bearer $TOKEN"

# Run all tests:
go test ./... -race -cover
```

---

## Extension Challenges

### Intermediate
1. Add **full-text search** to notes: build a simple in-memory inverted index (word → set of noteIDs). Update the index on Create/Update/Delete. Add `GET /api/v1/notes/search?q=word` that returns matching notes for the authenticated user. Cap results at 20.

2. Add **note sharing**: `POST /api/v1/notes/{id}/share` with `{"userID": N, "permission": "read"|"edit"}`. Store share grants. Modify `GetByID` to check both ownership and share grants. Broadcast a `"shared"` event over SSE.

### Advanced
3. Build the **gRPC search service** (`cmd/search/main.go`) with a `Search(query)` RPC that the API server calls instead of the local inverted index. The search service maintains its own in-memory index, updated via a gRPC streaming call from the API server whenever notes change. This demonstrates the service boundary pattern.

4. Add **operational readiness**: `GET /ready` returns 200 only if all stores are accessible and the event bus is running (check `cap(events) - len(events) > 10` to verify the bus isn't backed up). `GET /live` returns 200 always (just means the process is running). Use these as Kubernetes readiness and liveness probes.
