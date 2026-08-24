# Chapter 36: Building with SurrealDB in Go

Connect Go to SurrealDB, run queries, traverse graphs, and listen to live updates. We'll build a real-time social feed API.

## Table of Contents

1. Connecting to SurrealDB from Go
2. CRUD Operations
3. Graph Queries in Go
4. Live Queries via WebSocket
5. Mini Project: Real-Time Social Feed
6. Exercises

---

## 1. Connecting to SurrealDB from Go

```bash
# Start SurrealDB
docker run -d -p 8000:8000 surrealdb/surrealdb:latest \
    start --user root --pass root memory

go get github.com/surrealdb/surrealdb.go
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    surrealdb "github.com/surrealdb/surrealdb.go"
)

func main() {
    db, err := surrealdb.New("ws://localhost:8000/rpc")
    if err != nil {
        log.Fatal("connect:", err)
    }
    defer db.Close()

    // Sign in as root
    if _, err = db.Signin(map[string]interface{}{
        "user": "root",
        "pass": "root",
    }); err != nil {
        log.Fatal("signin:", err)
    }

    // Select namespace and database
    if _, err = db.Use("myapp", "production"); err != nil {
        log.Fatal("use:", err)
    }

    fmt.Println("Connected to SurrealDB!")
}
```

---

## 2. CRUD Operations

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// Create a record
func createUser(db *surrealdb.DB, id, name, email string, age int) (*User, error) {
    data, err := surrealdb.SmartUnmarshal[User](db.Create("user:"+id, map[string]interface{}{
        "name":  name,
        "email": email,
        "age":   age,
    }))
    if err != nil {
        return nil, err
    }
    return &data, nil
}

// Read a record
func getUser(db *surrealdb.DB, id string) (*User, error) {
    data, err := surrealdb.SmartUnmarshal[User](db.Select("user:" + id))
    if err != nil {
        return nil, err
    }
    return &data, nil
}

// Update a record
func updateUserAge(db *surrealdb.DB, id string, age int) error {
    _, err := db.Merge("user:"+id, map[string]interface{}{
        "age": age,
    })
    return err
}

// Delete a record
func deleteUser(db *surrealdb.DB, id string) error {
    _, err := db.Delete("user:" + id)
    return err
}

// Query with SurrealQL
func getUsersOver(db *surrealdb.DB, minAge int) ([]User, error) {
    result, err := db.Query(
        "SELECT * FROM user WHERE age > $age",
        map[string]interface{}{"age": minAge},
    )
    if err != nil {
        return nil, err
    }

    users, err := surrealdb.SmartUnmarshal[[]User](result, err)
    return users, err
}
```

---

## 3. Graph Queries in Go

```go
type Post struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    Author  string `json:"author"`
}

// Create a "follows" relationship
func followUser(db *surrealdb.DB, fromID, toID string) error {
    _, err := db.Query(
        "RELATE $from->follows->$to",
        map[string]interface{}{
            "from": "user:" + fromID,
            "to":   "user:" + toID,
        },
    )
    return err
}

// Get users that a given user follows
func getFollowing(db *surrealdb.DB, userID string) ([]User, error) {
    result, err := db.Query(
        "SELECT ->follows->user.* AS following FROM $user",
        map[string]interface{}{"user": "user:" + userID},
    )
    if err != nil {
        return nil, err
    }

    type Wrapper struct {
        Following []User `json:"following"`
    }
    wrapped, err := surrealdb.SmartUnmarshal[[]Wrapper](result, err)
    if err != nil || len(wrapped) == 0 {
        return nil, err
    }
    return wrapped[0].Following, nil
}

// Get feed: posts from users that userID follows
func getFeed(db *surrealdb.DB, userID string) ([]Post, error) {
    result, err := db.Query(`
        SELECT ->follows->user->posts.* AS posts
        FROM $user
        ORDER BY posts.created DESC
        LIMIT 20
    `, map[string]interface{}{"user": "user:" + userID})
    if err != nil {
        return nil, err
    }

    type Wrapper struct {
        Posts []Post `json:"posts"`
    }
    wrapped, err := surrealdb.SmartUnmarshal[[]Wrapper](result, err)
    if err != nil || len(wrapped) == 0 {
        return nil, err
    }
    return wrapped[0].Posts, nil
}

// Create a post and link to author
func createPost(db *surrealdb.DB, authorID, title, content string) (*Post, error) {
    result, err := db.Query(`
        LET $post = CREATE posts SET
            title = $title,
            content = $content,
            author = $author,
            created = time::now();
        RELATE $author->posts->$post;
        RETURN $post
    `, map[string]interface{}{
        "title":   title,
        "content": content,
        "author":  "user:" + authorID,
    })
    if err != nil {
        return nil, err
    }

    posts, err := surrealdb.SmartUnmarshal[[]Post](result, err)
    if err != nil || len(posts) == 0 {
        return nil, err
    }
    return &posts[0], nil
}
```

---

## 4. Live Queries via WebSocket

SurrealDB's live queries deliver real-time updates over the same WebSocket connection:

```go
type LiveEvent struct {
    Action string          `json:"action"` // CREATE, UPDATE, DELETE
    ID     string          `json:"id"`
    Result json.RawMessage `json:"result"`
}

func watchTable(db *surrealdb.DB, table string, handler func(LiveEvent)) error {
    // Start live query
    result, err := db.Query("LIVE SELECT * FROM "+table, nil)
    if err != nil {
        return err
    }

    var liveID string
    if ids, err := surrealdb.SmartUnmarshal[[]string](result, err); err == nil && len(ids) > 0 {
        liveID = ids[0]
    }

    fmt.Printf("Live query started: %s\n", liveID)

    // SurrealDB Go client sends live events via a channel
    notifications := make(chan surrealdb.Notification)
    db.RegisterForLive(liveID, notifications)

    go func() {
        for notif := range notifications {
            var event LiveEvent
            if err := json.Unmarshal(notif.Result, &event); err == nil {
                handler(event)
            }
        }
    }()

    return nil
}
```

---

## 5. Mini Project: Real-Time Social Feed

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    surrealdb "github.com/surrealdb/surrealdb.go"
)

var db *surrealdb.DB

// Server-Sent Events hub for live updates
type Hub struct {
    mu      sync.Mutex
    clients map[string]chan string // userID -> channel
}

var hub = &Hub{clients: make(map[string]chan string)}

func main() {
    var err error
    db, err = surrealdb.New("ws://localhost:8000/rpc")
    if err != nil {
        log.Fatal(err)
    }
    db.Signin(map[string]interface{}{"user": "root", "pass": "root"})
    db.Use("myapp", "prod")

    // Watch for new posts globally and push to followers
    watchTable(db, "posts", func(event LiveEvent) {
        if event.Action == "CREATE" {
            broadcastPost(event)
        }
    })

    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", handleCreateUser)
    mux.HandleFunc("POST /follow", handleFollow)
    mux.HandleFunc("POST /posts", handleCreatePost)
    mux.HandleFunc("GET /feed/{userID}", handleGetFeed)
    mux.HandleFunc("GET /live/{userID}", handleLiveFeed) // SSE endpoint

    log.Println("Social Feed API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var body struct {
        ID    string `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    json.NewDecoder(r.Body).Decode(&body)
    user, err := createUser(db, body.ID, body.Name, body.Email, 25)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(user)
}

func handleFollow(w http.ResponseWriter, r *http.Request) {
    var body struct {
        From string `json:"from"`
        To   string `json:"to"`
    }
    json.NewDecoder(r.Body).Decode(&body)
    if err := followUser(db, body.From, body.To); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(204)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
    var body struct {
        AuthorID string `json:"author_id"`
        Title    string `json:"title"`
        Content  string `json:"content"`
    }
    json.NewDecoder(r.Body).Decode(&body)
    post, err := createPost(db, body.AuthorID, body.Title, body.Content)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(post)
}

func handleGetFeed(w http.ResponseWriter, r *http.Request) {
    userID := r.PathValue("userID")
    posts, err := getFeed(db, userID)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(posts)
}

func handleLiveFeed(w http.ResponseWriter, r *http.Request) {
    userID := r.PathValue("userID")

    // Server-Sent Events
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    ch := make(chan string, 10)
    hub.mu.Lock()
    hub.clients[userID] = ch
    hub.mu.Unlock()

    defer func() {
        hub.mu.Lock()
        delete(hub.clients, userID)
        hub.mu.Unlock()
    }()

    flusher := w.(http.Flusher)
    ctx := r.Context()

    for {
        select {
        case msg := <-ch:
            w.Write([]byte("data: " + msg + "\n\n"))
            flusher.Flush()
        case <-ctx.Done():
            return
        }
    }
}

func broadcastPost(event LiveEvent) {
    // In production: query which users follow the post author, then push only to them
    hub.mu.Lock()
    defer hub.mu.Unlock()
    for _, ch := range hub.clients {
        select {
        case ch <- string(event.Result):
        default: // don't block if client is slow
        }
    }
}
```

Test:
```bash
# Create users
curl -X POST localhost:8080/users -d '{"id":"alice","name":"Alice","email":"alice@example.com"}'
curl -X POST localhost:8080/users -d '{"id":"bob","name":"Bob","email":"bob@example.com"}'

# Follow
curl -X POST localhost:8080/follow -d '{"from":"alice","to":"bob"}'

# Create post as Bob
curl -X POST localhost:8080/posts -d '{"author_id":"bob","title":"Hello World","content":"My first post!"}'

# Alice's feed
curl localhost:8080/feed/alice
# Returns Bob's post!

# Live feed (open in separate terminal)
curl -N localhost:8080/live/alice
# Events arrive as Bob publishes new posts
```

---

## Summary

- Connect to SurrealDB with a WebSocket URL: `ws://localhost:8000/rpc`.
- Use `surrealdb.SmartUnmarshal[T]` to decode query results into Go structs.
- Graph traversal (`->follows->user->posts`) returns nested data in one query — no multiple round trips.
- Live queries + SSE = real-time feed with zero polling.

### Exercises

**Easy:** Create 3 user records in SurrealDB from Go. Write a function that queries all users and prints their names.

**Medium:** Implement a "like" system: users can like posts. Create a `RELATE user->likes->post` relationship. Write a function that returns all posts liked by a given user, and another that returns the like count for a post.

**Hard:** Add a "mutual followers" feature: given two user IDs, find users that both follow. Use SurrealDB graph queries. Cache the result in Redis for 5 minutes to avoid re-querying the graph on every request.
