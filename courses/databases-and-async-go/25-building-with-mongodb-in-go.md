# Chapter 25: Building with MongoDB in Go

This chapter is pure code: complete Go patterns for MongoDB, from connection management to aggregation pipelines, with a full blog engine project at the end.

## Table of Contents

1. Connection Management
2. CRUD Patterns in Go
3. Error Handling
4. Aggregation in Go
5. Transactions in Go
6. Mini Project: Blog Engine with MongoDB
7. Exercises

---

## 1. Connection Management

```go
package mongodb

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
    client *mongo.Client
    db     *mongo.Database
}

func New(uri, dbName string) (*Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(20).
        SetMinPoolSize(5).
        SetMaxConnIdleTime(30 * time.Minute)

    client, err := mongo.Connect(ctx, opts)
    if err != nil {
        return nil, fmt.Errorf("connect: %w", err)
    }

    // Verify connection
    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        return nil, fmt.Errorf("ping: %w", err)
    }

    return &Client{
        client: client,
        db:     client.Database(dbName),
    }, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
    return c.db.Collection(name)
}

func (c *Client) Close(ctx context.Context) error {
    return c.client.Disconnect(ctx)
}
```

---

## 2. CRUD Patterns in Go

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Title     string             `bson:"title"`
    Content   string             `bson:"content"`
    Author    string             `bson:"author"`
    Tags      []string           `bson:"tags"`
    Views     int64              `bson:"views"`
    CreatedAt time.Time          `bson:"createdAt"`
    UpdatedAt time.Time          `bson:"updatedAt"`
}

type PostRepository struct {
    coll *mongo.Collection
}

func NewPostRepository(coll *mongo.Collection) *PostRepository {
    return &PostRepository{coll: coll}
}

// Create inserts a new post and returns its ID.
func (r *PostRepository) Create(ctx context.Context, p *Post) (primitive.ObjectID, error) {
    p.CreatedAt = time.Now()
    p.UpdatedAt = time.Now()

    result, err := r.coll.InsertOne(ctx, p)
    if err != nil {
        return primitive.NilObjectID, err
    }
    return result.InsertedID.(primitive.ObjectID), nil
}

// GetByID retrieves a post by its ObjectID.
func (r *PostRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Post, error) {
    var post Post
    err := r.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&post)
    if errors.Is(err, mongo.ErrNoDocuments) {
        return nil, nil
    }
    return &post, err
}

// ListByAuthor returns posts by a given author, paginated.
func (r *PostRepository) ListByAuthor(ctx context.Context, author string, limit, skip int64) ([]Post, error) {
    opts := options.Find().
        SetSort(bson.D{{Key: "createdAt", Value: -1}}).
        SetLimit(limit).
        SetSkip(skip)

    cursor, err := r.coll.Find(ctx,
        bson.D{{Key: "author", Value: author}},
        opts,
    )
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var posts []Post
    if err := cursor.All(ctx, &posts); err != nil {
        return nil, err
    }
    return posts, nil
}

// FindByTags returns posts that have ALL the given tags.
func (r *PostRepository) FindByTags(ctx context.Context, tags []string) ([]Post, error) {
    cursor, err := r.coll.Find(ctx, bson.D{{Key: "tags", Value: bson.D{{Key: "$all", Value: tags}}}})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var posts []Post
    cursor.All(ctx, &posts)
    return posts, nil
}

// IncrementViews atomically increments the view count.
func (r *PostRepository) IncrementViews(ctx context.Context, id primitive.ObjectID) error {
    _, err := r.coll.UpdateOne(ctx,
        bson.D{{Key: "_id", Value: id}},
        bson.D{
            {Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}},
            {Key: "$set", Value: bson.D{{Key: "updatedAt", Value: time.Now()}}},
        },
    )
    return err
}

// Update replaces specific fields.
func (r *PostRepository) Update(ctx context.Context, id primitive.ObjectID, title, content string) error {
    result, err := r.coll.UpdateOne(ctx,
        bson.D{{Key: "_id", Value: id}},
        bson.D{{Key: "$set", Value: bson.D{
            {Key: "title", Value: title},
            {Key: "content", Value: content},
            {Key: "updatedAt", Value: time.Now()},
        }}},
    )
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return fmt.Errorf("post %s not found", id.Hex())
    }
    return nil
}

// Delete removes a post by ID.
func (r *PostRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
    result, err := r.coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return fmt.Errorf("post %s not found", id.Hex())
    }
    return nil
}
```

---

## 3. Error Handling

```go
import (
    "errors"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func IsDuplicate(err error) bool {
    var writeException mongo.WriteException
    if errors.As(err, &writeException) {
        for _, we := range writeException.WriteErrors {
            if we.Code == 11000 { // duplicate key error
                return true
            }
        }
    }
    return false
}

func IsNotFound(err error) bool {
    return errors.Is(err, mongo.ErrNoDocuments)
}

// Usage:
func (r *PostRepository) CreateUser(ctx context.Context, email string) error {
    _, err := r.coll.InsertOne(ctx, bson.D{{Key: "email", Value: email}})
    if IsDuplicate(err) {
        return fmt.Errorf("email %q already registered", email)
    }
    return err
}
```

---

## 4. Aggregation in Go

```go
// Top authors by total views
func (r *PostRepository) TopAuthors(ctx context.Context, limit int) ([]bson.M, error) {
    pipeline := mongo.Pipeline{
        // Group by author, sum views
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$author"},
            {Key: "totalViews", Value: bson.D{{Key: "$sum", Value: "$views"}}},
            {Key: "postCount", Value: bson.D{{Key: "$sum", Value: 1}}},
        }}},
        // Sort by views descending
        {{Key: "$sort", Value: bson.D{{Key: "totalViews", Value: -1}}}},
        // Limit results
        {{Key: "$limit", Value: limit}},
        // Rename _id to author
        {{Key: "$project", Value: bson.D{
            {Key: "author", Value: "$_id"},
            {Key: "totalViews", Value: 1},
            {Key: "postCount", Value: 1},
            {Key: "_id", Value: 0},
        }}},
    }

    cursor, err := r.coll.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []bson.M
    return results, cursor.All(ctx, &results)
}

// Tag frequency: which tags are most used
func (r *PostRepository) TagFrequency(ctx context.Context) ([]bson.M, error) {
    pipeline := mongo.Pipeline{
        // Unwind the tags array into separate documents
        {{Key: "$unwind", Value: "$tags"}},
        // Group by tag, count occurrences
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$tags"},
            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
        }}},
        // Sort by count
        {{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
        {{Key: "$limit", Value: 20}},
    }

    cursor, err := r.coll.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []bson.M
    return results, cursor.All(ctx, &results)
}
```

---

## 5. Transactions in Go

```go
func TransferCredits(ctx context.Context, client *mongo.Client, fromID, toID primitive.ObjectID, amount int) error {
    session, err := client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    _, err = session.WithTransaction(ctx, func(sCtx mongo.SessionContext) (interface{}, error) {
        users := client.Database("app").Collection("users")

        // Debit
        result, err := users.UpdateOne(sCtx,
            bson.D{
                {Key: "_id", Value: fromID},
                {Key: "credits", Value: bson.D{{Key: "$gte", Value: amount}}},
            },
            bson.D{{Key: "$inc", Value: bson.D{{Key: "credits", Value: -amount}}}},
        )
        if err != nil {
            return nil, err
        }
        if result.ModifiedCount == 0 {
            return nil, errors.New("insufficient credits")
        }

        // Credit
        _, err = users.UpdateOne(sCtx,
            bson.D{{Key: "_id", Value: toID}},
            bson.D{{Key: "$inc", Value: bson.D{{Key: "credits", Value: amount}}}},
        )
        return nil, err
    })
    return err
}
```

`session.WithTransaction` automatically retries on transient errors (network issues, write conflicts) and handles commit/abort.

---

## 6. Mini Project: Blog Engine with MongoDB

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var postsColl *mongo.Collection

func main() {
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    postsColl = client.Database("blog").Collection("posts")

    // Create indexes
    postsColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{Key: "author", Value: 1}, {Key: "createdAt", Value: -1}}},
        {Keys: bson.D{{Key: "tags", Value: 1}}},
        {Keys: bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}}},
    })

    mux := http.NewServeMux()
    mux.HandleFunc("GET /posts", listPosts)
    mux.HandleFunc("POST /posts", createPost)
    mux.HandleFunc("GET /posts/", getPost)
    mux.HandleFunc("GET /search", searchPosts)

    log.Println("Blog API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func listPosts(w http.ResponseWriter, r *http.Request) {
    author := r.URL.Query().Get("author")
    filter := bson.D{}
    if author != "" {
        filter = bson.D{{Key: "author", Value: author}}
    }

    cursor, err := postsColl.Find(r.Context(), filter,
        options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(20))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    var posts []Post
    if err := cursor.All(r.Context(), &posts); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if posts == nil {
        posts = []Post{}
    }
    json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title   string   `json:"title"`
        Content string   `json:"content"`
        Author  string   `json:"author"`
        Tags    []string `json:"tags"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }

    post := Post{
        Title:     req.Title,
        Content:   req.Content,
        Author:    req.Author,
        Tags:      req.Tags,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    result, err := postsColl.InsertOne(r.Context(), post)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    post.ID = result.InsertedID.(primitive.ObjectID)
    w.WriteHeader(201)
    json.NewEncoder(w).Encode(post)
}

func getPost(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        http.Error(w, "invalid id", 400)
        return
    }

    var post Post
    err = postsColl.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}).Decode(&post)
    if err == mongo.ErrNoDocuments {
        http.Error(w, "not found", 404)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Increment views in background
    go postsColl.UpdateOne(context.Background(),
        bson.D{{Key: "_id", Value: id}},
        bson.D{{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}}},
    )

    json.NewEncoder(w).Encode(post)
}

func searchPosts(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    if q == "" {
        http.Error(w, "q required", 400)
        return
    }

    cursor, err := postsColl.Find(r.Context(),
        bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: q}}}},
        options.Find().
            SetProjection(bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}}}).
            SetSort(bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}}}).
            SetLimit(10),
    )
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    var results []bson.M
    cursor.All(r.Context(), &results)
    json.NewEncoder(w).Encode(results)
}
```

---

## Summary

- Use `cursor.All(ctx, &slice)` to decode all results at once.
- `bson.D` for ordered documents (queries, updates, pipelines). `bson.M` for unordered maps (when field order doesn't matter).
- Check `result.MatchedCount == 0` after UpdateOne to detect "not found".
- `session.WithTransaction` handles retry logic automatically.
- Error code 11000 = duplicate key. `mongo.ErrNoDocuments` = FindOne with no match.

### Exercises

**Easy:** Add a `PUT /posts/{id}` endpoint that updates a post's title, content, and tags. Return 404 if not found.

**Medium:** Build an aggregation that returns the monthly post count for the last 6 months (use `$dateToString` to group by month).

**Hard:** Add a comments system: each post can have comments. Decide whether to embed comments in the post document or store them in a separate collection. Justify your choice. Implement `POST /posts/{id}/comments` and make the post endpoint include the latest 5 comments.
