# Chapter 72: Mini Project 5 — Blog Platform REST API

A full-featured blog platform REST API with posts, tags, comments, user authentication, and pagination. This consolidates the entire Web Development volume: chi router, JWT auth, middleware, OpenAPI, PostgreSQL with SQLC, and WebSockets for live comment notifications.

## What You'll Build

```
blog/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── domain/
│   │   ├── post.go
│   │   ├── comment.go
│   │   └── user.go
│   ├── handler/
│   │   ├── post.go
│   │   ├── comment.go
│   │   └── auth.go
│   ├── store/
│   │   └── store.go
│   ├── middleware/
│   │   └── auth.go
│   └── config/
│       └── config.go
└── migrations/
    └── 001_initial.sql
```

---

## 1. API Design

```
Authentication
  POST   /auth/register      Create account
  POST   /auth/login         Get access + refresh tokens
  POST   /auth/refresh       Refresh access token

Posts
  GET    /posts              List posts (paginated, filterable by tag)
  POST   /posts              Create post (auth required)
  GET    /posts/:slug        Get single post with comments
  PUT    /posts/:id          Update post (owner only)
  DELETE /posts/:id          Delete post (owner only)

Comments
  POST   /posts/:id/comments Add comment (auth required)
  DELETE /comments/:id       Delete comment (owner only)

Tags
  GET    /tags               List all tags with post counts

Users
  GET    /users/:id          Public profile
  PUT    /users/me           Update own profile (auth required)

Real-time
  WS     /ws/posts/:id       Live comments via WebSocket
```

---

## 2. Domain Types

```go
// internal/domain/post.go
package domain

import "time"

type Post struct {
    ID          int64     `json:"id"`
    Slug        string    `json:"slug"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Summary     string    `json:"summary"`
    AuthorID    int64     `json:"authorId"`
    Author      *User     `json:"author,omitempty"`
    Tags        []Tag     `json:"tags,omitempty"`
    Published   bool      `json:"published"`
    PublishedAt *time.Time `json:"publishedAt,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type Tag struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Slug  string `json:"slug"`
    Count int    `json:"count,omitempty"`
}

type CreatePostRequest struct {
    Title   string   `json:"title"`
    Content string   `json:"content"`
    Summary string   `json:"summary"`
    Tags    []string `json:"tags"`    // tag names or slugs
    Publish bool     `json:"publish"` // publish immediately?
}

func (r *CreatePostRequest) Validate() ValidationErrors {
    var errs ValidationErrors
    if r.Title = strings.TrimSpace(r.Title); r.Title == "" {
        errs = append(errs, ValidationError{Field: "title", Message: "required"})
    } else if len(r.Title) > 200 {
        errs = append(errs, ValidationError{Field: "title", Message: "max 200 characters"})
    }
    if r.Content = strings.TrimSpace(r.Content); r.Content == "" {
        errs = append(errs, ValidationError{Field: "content", Message: "required"})
    }
    return errs
}

type PostFilter struct {
    Tag       string
    AuthorID  int64
    Published *bool
    Page      int
    PageSize  int
}

type PostPage struct {
    Posts      []*Post `json:"posts"`
    Total      int     `json:"total"`
    Page       int     `json:"page"`
    PageSize   int     `json:"pageSize"`
    HasMore    bool    `json:"hasMore"`
}
```

```go
// internal/domain/user.go
package domain

import "time"

type User struct {
    ID           int64     `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email,omitempty"` // hidden in public views
    Bio          string    `json:"bio"`
    AvatarURL    string    `json:"avatarUrl"`
    PostCount    int       `json:"postCount,omitempty"`
    CreatedAt    time.Time `json:"createdAt"`
}

type Comment struct {
    ID        int64     `json:"id"`
    PostID    int64     `json:"postId"`
    AuthorID  int64     `json:"authorId"`
    Author    *User     `json:"author,omitempty"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"createdAt"`
}
```

---

## 3. Store Interface

```go
// internal/store/store.go
package store

import (
    "context"
    "blog/internal/domain"
)

type PostStore interface {
    Create(ctx context.Context, post *domain.Post) error
    GetByID(ctx context.Context, id int64) (*domain.Post, error)
    GetBySlug(ctx context.Context, slug string) (*domain.Post, error)
    Update(ctx context.Context, post *domain.Post) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, filter domain.PostFilter) (*domain.PostPage, error)
}

type CommentStore interface {
    Create(ctx context.Context, comment *domain.Comment) error
    ListByPost(ctx context.Context, postID int64) ([]*domain.Comment, error)
    Delete(ctx context.Context, id, authorID int64) error
}

type UserStore interface {
    Create(ctx context.Context, user *domain.User, passwordHash string) error
    GetByID(ctx context.Context, id int64) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, string, error) // returns user + password hash
    GetByUsername(ctx context.Context, username string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
}

type TagStore interface {
    GetOrCreate(ctx context.Context, name string) (*domain.Tag, error)
    ListWithCounts(ctx context.Context) ([]*domain.Tag, error)
    SetPostTags(ctx context.Context, postID int64, tagIDs []int64) error
    GetByPost(ctx context.Context, postID int64) ([]*domain.Tag, error)
}

// Store combines all stores
type Store struct {
    Posts    PostStore
    Comments CommentStore
    Users    UserStore
    Tags     TagStore
}
```

---

## 4. Handlers

```go
// internal/handler/post.go
package handler

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/go-chi/chi/v5"
    "blog/internal/domain"
    "blog/internal/store"
)

type PostHandler struct {
    store *store.Store
}

func NewPostHandler(s *store.Store) *PostHandler {
    return &PostHandler{store: s}
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
    filter := domain.PostFilter{
        Tag:      r.URL.Query().Get("tag"),
        Page:     pageParam(r, 1),
        PageSize: pageSizeParam(r, 20, 100),
    }
    published := true
    filter.Published = &published // only show published posts in public listing

    page, err := h.store.Posts.List(r.Context(), filter)
    if err != nil {
        writeServerError(w, err)
        return
    }
    writeJSON(w, http.StatusOK, page)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
    user := userFromContext(r.Context())
    
    var req domain.CreatePostRequest
    if err := decodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    if errs := req.Validate(); len(errs) > 0 {
        writeValidationError(w, errs)
        return
    }

    // Generate slug from title
    slug := slugify(req.Title)
    now := time.Now()
    
    post := &domain.Post{
        Slug:      slug,
        Title:     req.Title,
        Content:   req.Content,
        Summary:   req.Summary,
        AuthorID:  user.ID,
        Published: req.Publish,
    }
    if req.Publish {
        post.PublishedAt = &now
    }

    if err := h.store.Posts.Create(r.Context(), post); err != nil {
        writeServerError(w, err)
        return
    }

    // Attach tags
    for _, tagName := range req.Tags {
        tag, err := h.store.Tags.GetOrCreate(r.Context(), tagName)
        if err != nil { continue }
        _ = tag // collected IDs would be passed to SetPostTags
    }

    writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "slug")
    
    post, err := h.store.Posts.GetBySlug(r.Context(), slug)
    if err != nil {
        writeNotFound(w, "post not found")
        return
    }
    
    // Load comments and tags
    comments, _ := h.store.Comments.ListByPost(r.Context(), post.ID)
    tags, _ := h.store.Tags.GetByPost(r.Context(), post.ID)
    post.Tags = tags
    
    writeJSON(w, http.StatusOK, map[string]any{
        "post":     post,
        "comments": comments,
    })
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
    user := userFromContext(r.Context())
    id := int64Param(r, "id")
    
    post, err := h.store.Posts.GetByID(r.Context(), id)
    if err != nil {
        writeNotFound(w, "post not found")
        return
    }
    if post.AuthorID != user.ID {
        writeForbidden(w)
        return
    }

    var req domain.CreatePostRequest
    if err := decodeJSON(r, &req); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    post.Title = req.Title
    post.Content = req.Content
    post.Summary = req.Summary
    
    if err := h.store.Posts.Update(r.Context(), post); err != nil {
        writeServerError(w, err)
        return
    }
    writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
    user := userFromContext(r.Context())
    id := int64Param(r, "id")
    
    post, err := h.store.Posts.GetByID(r.Context(), id)
    if err != nil {
        writeNotFound(w, "post not found")
        return
    }
    if post.AuthorID != user.ID {
        writeForbidden(w)
        return
    }
    if err := h.store.Posts.Delete(r.Context(), id); err != nil {
        writeServerError(w, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func slugify(title string) string {
    slug := strings.ToLower(title)
    slug = strings.Map(func(r rune) rune {
        if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' { return r }
        if r == ' ' || r == '-' { return '-' }
        return -1
    }, slug)
    // Remove consecutive dashes
    for strings.Contains(slug, "--") {
        slug = strings.ReplaceAll(slug, "--", "-")
    }
    return strings.Trim(slug, "-")
}

func pageParam(r *http.Request, def int) int {
    v, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || v < 1 { return def }
    return v
}

func pageSizeParam(r *http.Request, def, max int) int {
    v, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
    if err != nil || v < 1 { return def }
    if v > max { return max }
    return v
}

func int64Param(r *http.Request, name string) int64 {
    v, _ := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
    return v
}
```

---

## 5. Router

```go
// cmd/server/main.go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    chimw "github.com/go-chi/chi/v5/middleware"
    "blog/internal/handler"
    "blog/internal/middleware"
    "blog/internal/config"
    "blog/internal/store/postgres"
)

func main() {
    cfg, err := config.Load()
    if err != nil { slog.Error("config", "err", err); os.Exit(1) }

    db, err := postgres.Connect(cfg.DatabaseURL)
    if err != nil { slog.Error("db", "err", err); os.Exit(1) }
    defer db.Close()

    s := postgres.NewStore(db)
    
    authMW := middleware.NewAuth(cfg.JWTSecret)
    postH  := handler.NewPostHandler(s)
    commentH := handler.NewCommentHandler(s)
    authH  := handler.NewAuthHandler(s, cfg.JWTSecret)
    tagH   := handler.NewTagHandler(s)

    r := chi.NewRouter()
    r.Use(chimw.RequestID)
    r.Use(chimw.RealIP)
    r.Use(chimw.Logger)
    r.Use(chimw.Recoverer)
    r.Use(chimw.Timeout(30 * time.Second))

    // Public routes
    r.Post("/auth/register", authH.Register)
    r.Post("/auth/login",    authH.Login)
    r.Post("/auth/refresh",  authH.Refresh)

    r.Get("/posts",         postH.List)
    r.Get("/posts/{slug}",  postH.Get)
    r.Get("/tags",          tagH.List)
    r.Get("/users/{id}",    authH.GetUser)

    // Authenticated routes
    r.Group(func(r chi.Router) {
        r.Use(authMW.Authenticate)

        r.Post("/posts",              postH.Create)
        r.Put("/posts/{id}",          postH.Update)
        r.Delete("/posts/{id}",       postH.Delete)

        r.Post("/posts/{id}/comments", commentH.Create)
        r.Delete("/comments/{id}",     commentH.Delete)

        r.Put("/users/me", authH.UpdateMe)
    })

    // WebSocket for live comments
    hub := handler.NewHub()
    go hub.Run()
    r.Get("/ws/posts/{id}", handler.ServeWS(hub))

    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      r,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    go func() {
        slog.Info("blog API starting", "port", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server", "err", err)
            os.Exit(1)
        }
    }()

    <-ctx.Done()
    slog.Info("shutting down")
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer shutdownCancel()
    srv.Shutdown(shutdownCtx)
}
```

---

## 6. Database Migration

```sql
-- migrations/001_initial.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    username     VARCHAR(50)  UNIQUE NOT NULL,
    email        VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    bio          TEXT         NOT NULL DEFAULT '',
    avatar_url   TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE posts (
    id           BIGSERIAL PRIMARY KEY,
    slug         VARCHAR(200) UNIQUE NOT NULL,
    title        VARCHAR(200) NOT NULL,
    content      TEXT         NOT NULL,
    summary      TEXT         NOT NULL DEFAULT '',
    author_id    BIGINT       NOT NULL REFERENCES users(id),
    published    BOOLEAN      NOT NULL DEFAULT false,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX posts_author_id_idx ON posts(author_id);
CREATE INDEX posts_published_idx ON posts(published, published_at DESC);

CREATE TABLE tags (
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE post_tags (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  BIGINT NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE comments (
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT      NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id  BIGINT      NOT NULL REFERENCES users(id),
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX comments_post_id_idx ON comments(post_id, created_at);

-- Full text search on posts
ALTER TABLE posts ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('english', title || ' ' || content)
    ) STORED;

CREATE INDEX posts_search_idx ON posts USING GIN(search_vector);
```

---

## 7. Key Concepts Applied

| Concept (Vol 5-6) | Where Used |
|-------------------|------------|
| chi Router | Route registration with URL params |
| JWT authentication | Auth middleware, access/refresh tokens |
| Middleware chains | RequestID, Logger, Recover, Auth |
| JSON decode + validate | All request handlers |
| Pagination | PostFilter, PostPage |
| SQLC + PostgreSQL | Store layer |
| Full-text search | GIN index, `search_vector` |
| WebSocket | Live comment hub |
| Graceful shutdown | signal.NotifyContext |

---

## Exercises

### Easy
1. Add a `GET /posts/search?q=<query>` endpoint that uses the `search_vector` full-text index.
2. Add `GET /posts/{slug}/related` that returns 5 posts sharing the most tags with the given post.
3. Add pagination metadata (`X-Total-Count`, `Link` headers) to the `GET /posts` endpoint.

### Medium
4. Add **draft previews**: authenticated authors can GET their own unpublished posts. Non-authenticated users and other users see 404 for unpublished posts.
5. Add **comment reactions** (like/dislike): `POST /comments/{id}/reactions` with a `type` field. Store in a `comment_reactions` table; deduplicate per user per comment.
6. Implement **post scheduling**: posts can have a future `publish_at` time. A background goroutine runs every minute and publishes due posts. Handle time zones correctly.

### Hard
7. Add **RSS and Atom feeds**: `GET /feed.xml` and `GET /feed.atom` that return valid XML feeds of the 20 most recent published posts. Use Go's `encoding/xml`.
8. Implement **full OpenAPI 3.0 spec** for this API. Generate the Go handler interfaces with `oapi-codegen`. Ensure every endpoint's request and response shapes are specified and validated. Add Swagger UI at `/docs`.
