# Chapter 24: MongoDB — Documents All the Way Down

MongoDB is the most popular NoSQL database. It stores data as flexible JSON-like documents, requires no fixed schema, and scales horizontally. This chapter teaches you MongoDB from first principles through the query language to internal architecture.

## Table of Contents

1. MongoDB Concepts — Documents, Collections, Databases
2. The MongoDB Query Language
3. Aggregation Pipeline — The MongoDB Power Tool
4. Indexes in MongoDB
5. Schema Design — Embedding vs Referencing
6. Transactions in MongoDB
7. MongoDB Architecture
8. Building with MongoDB in Go
9. Exercises

---

## 1. MongoDB Concepts

**Document:** A JSON-like record. The basic unit of data. Can have any fields.
```json
{
  "_id": ObjectId("6502a7b8c9d1e2f3a4b5c6d7"),
  "title": "Introduction to Go",
  "author": "Alice",
  "tags": ["go", "programming"],
  "views": 1523,
  "published": true,
  "publishedAt": ISODate("2024-01-15T10:00:00Z")
}
```

**Collection:** A group of documents (like a SQL table, but with no fixed schema). Documents in the same collection can have different fields.

**Database:** A group of collections.

**ObjectId:** MongoDB's default `_id` type — a 12-byte unique ID containing a timestamp, machine ID, and counter. You can use any type as `_id`, including strings or integers.

---

## 2. The MongoDB Query Language

MongoDB uses JSON-like queries instead of SQL.

### Basic CRUD

```javascript
// In the mongo shell:

// Insert
db.posts.insertOne({ title: "Hello", author: "Alice", views: 0 })
db.posts.insertMany([
  { title: "Post 1", author: "Bob", views: 100 },
  { title: "Post 2", author: "Alice", views: 200 }
])

// Find (SELECT)
db.posts.find({})                          // all documents
db.posts.find({ author: "Alice" })         // where author = Alice
db.posts.find({ views: { $gt: 100 } })     // where views > 100
db.posts.findOne({ title: "Hello" })       // single document

// Comparison operators
$eq, $ne     // equal, not equal
$gt, $gte    // greater than, greater than or equal
$lt, $lte    // less than, less than or equal
$in          // in a list: { status: { $in: ["active", "pending"] } }
$nin         // not in a list

// Logical operators
$and, $or, $not, $nor

// Query examples
db.posts.find({ author: "Alice", views: { $gte: 100 } })      // AND
db.posts.find({ $or: [{ author: "Alice" }, { views: { $gt: 500 } }] })  // OR

// Update
db.posts.updateOne(
  { _id: ObjectId("...") },          // filter
  { $set: { views: 1600 } }          // update (only changes views, keeps other fields)
)
db.posts.updateMany(
  { author: "Alice" },               // update all Alice's posts
  { $inc: { views: 10 } }            // increment views by 10
)

// Delete
db.posts.deleteOne({ _id: ObjectId("...") })
db.posts.deleteMany({ published: false })  // delete all drafts
```

### Array Operators

```javascript
// Find posts that have the "go" tag
db.posts.find({ tags: "go" })

// Find posts that have ALL of these tags
db.posts.find({ tags: { $all: ["go", "programming"] } })

// Find posts that have ANY of these tags
db.posts.find({ tags: { $in: ["go", "python", "rust"] } })

// Add a tag
db.posts.updateOne({ _id: id }, { $push: { tags: "tutorial" } })

// Remove a tag
db.posts.updateOne({ _id: id }, { $pull: { tags: "tutorial" } })

// Add only if not present
db.posts.updateOne({ _id: id }, { $addToSet: { tags: "tutorial" } })
```

---

## 3. Aggregation Pipeline — The MongoDB Power Tool

The aggregation pipeline processes documents through a series of stages, like Unix pipes.

```javascript
// Find top 5 authors by total views
db.posts.aggregate([
  // Stage 1: filter published posts only
  { $match: { published: true } },
  
  // Stage 2: group by author, sum views
  { $group: {
    _id: "$author",
    totalViews: { $sum: "$views" },
    postCount: { $count: {} }
  }},
  
  // Stage 3: sort by total views descending
  { $sort: { totalViews: -1 } },
  
  // Stage 4: take top 5
  { $limit: 5 },
  
  // Stage 5: reshape output
  { $project: {
    author: "$_id",
    totalViews: 1,
    postCount: 1,
    _id: 0
  }}
])
```

Common aggregation stages:
- `$match` — filter (like WHERE)
- `$group` — group and aggregate (like GROUP BY)
- `$sort` — sort
- `$limit` / `$skip` — pagination
- `$project` — reshape documents (select/rename fields)
- `$lookup` — join with another collection
- `$unwind` — expand arrays into separate documents
- `$addFields` — add computed fields

```javascript
// $lookup: join posts with users collection
db.posts.aggregate([
  {
    $lookup: {
      from: "users",           // collection to join
      localField: "author_id", // field in posts
      foreignField: "_id",     // field in users
      as: "author"             // output field (array)
    }
  },
  { $unwind: "$author" },      // flatten the author array to an object
  { $project: {
    title: 1,
    "author.name": 1,
    views: 1
  }}
])
```

---

## 4. Indexes in MongoDB

```javascript
// Single field index
db.posts.createIndex({ author: 1 })   // 1 = ascending, -1 = descending

// Compound index
db.posts.createIndex({ author: 1, publishedAt: -1 })

// Unique index
db.users.createIndex({ email: 1 }, { unique: true })

// TTL index: automatically delete documents after N seconds
db.sessions.createIndex({ createdAt: 1 }, { expireAfterSeconds: 3600 })

// Text index for full-text search
db.posts.createIndex({ title: "text", content: "text" })
db.posts.find({ $text: { $search: "golang database" } })

// Explain a query
db.posts.find({ author: "Alice" }).explain("executionStats")
```

Look for `IXSCAN` (index scan) vs `COLLSCAN` (collection scan) in explain output.

---

## 5. Schema Design — Embedding vs Referencing

MongoDB's biggest design decision: should related data be embedded in the same document or stored in a separate collection?

### Embedding (Denormalized)

```json
// Post with embedded comments
{
  "_id": ObjectId("..."),
  "title": "Introduction to MongoDB",
  "content": "...",
  "comments": [
    { "author": "Bob", "text": "Great post!", "createdAt": "2024-01-15" },
    { "author": "Carol", "text": "Very helpful", "createdAt": "2024-01-16" }
  ]
}
```

**Embed when:**
- Data is always accessed together (post with its comments)
- Child data is bounded in size (< 100 items)
- Child data is "owned" by the parent

### Referencing (Normalized)

```json
// Post references user by ID
{
  "_id": ObjectId("post_123"),
  "title": "Introduction to MongoDB",
  "author_id": ObjectId("user_456")  // reference to users collection
}
```

**Reference when:**
- Child data is large or unbounded
- Child data is accessed independently
- Many-to-many relationships

---

## 6. Transactions in MongoDB

MongoDB supports multi-document ACID transactions since version 4.0:

```javascript
const session = client.startSession()
session.startTransaction()

try {
  db.accounts.updateOne(
    { _id: fromId },
    { $inc: { balance: -amount } },
    { session }
  )
  db.accounts.updateOne(
    { _id: toId },
    { $inc: { balance: amount } },
    { session }
  )
  session.commitTransaction()
} catch (err) {
  session.abortTransaction()
  throw err
} finally {
  session.endSession()
}
```

Transactions in MongoDB have overhead. Design your schema to avoid needing them (embedding usually eliminates the need for multi-document transactions).

---

## 7. MongoDB Architecture

- **mongod:** The main database process
- **Replica set:** 3+ mongod nodes (primary + secondaries). Writes go to primary, replicated to secondaries via oplog. Automatic failover if primary dies.
- **Sharding:** Data partitioned across multiple replica sets (shards) by a shard key. Allows horizontal scaling to petabytes.
- **mongos:** A router process that routes queries to the correct shards.

```
Client → mongos → Shard 1 (replica set)
                → Shard 2 (replica set)
                → Shard 3 (replica set)
```

---

## 8. Building with MongoDB in Go

```bash
go get go.mongodb.org/mongo-driver/mongo
```

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string             `bson:"title" json:"title"`
    Author      string             `bson:"author" json:"author"`
    Tags        []string           `bson:"tags" json:"tags"`
    Views       int                `bson:"views" json:"views"`
    PublishedAt time.Time          `bson:"publishedAt" json:"published_at"`
}

var posts *mongo.Collection

func main() {
    ctx := context.Background()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("cannot connect to MongoDB:", err)
    }
    fmt.Println("Connected to MongoDB!")

    posts = client.Database("blog").Collection("posts")

    // Create indexes
    posts.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{Key: "author", Value: 1}}},
        {Keys: bson.D{{Key: "tags", Value: 1}}},
        {
            Keys:    bson.D{{Key: "publishedAt", Value: -1}},
            Options: options.Index().SetName("publishedAt_desc"),
        },
    })

    // Insert a document
    post := Post{
        Title:       "Introduction to MongoDB",
        Author:      "Alice",
        Tags:        []string{"mongodb", "nosql", "go"},
        Views:       0,
        PublishedAt: time.Now(),
    }
    result, err := posts.InsertOne(ctx, post)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Inserted ID:", result.InsertedID)

    // Find by author
    cursor, err := posts.Find(ctx, bson.D{{Key: "author", Value: "Alice"}})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(ctx)

    var results []Post
    if err := cursor.All(ctx, &results); err != nil {
        log.Fatal(err)
    }
    for _, p := range results {
        fmt.Printf("Post: %s (%d views)\n", p.Title, p.Views)
    }

    // Update views
    posts.UpdateOne(ctx,
        bson.D{{Key: "_id", Value: result.InsertedID}},
        bson.D{{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}}},
    )

    // Aggregation: top authors by views
    pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$author"},
            {Key: "totalViews", Value: bson.D{{Key: "$sum", Value: "$views"}}},
        }}},
        {{Key: "$sort", Value: bson.D{{Key: "totalViews", Value: -1}}}},
        {{Key: "$limit", Value: 5}},
    }
    aggCursor, err := posts.Aggregate(ctx, pipeline)
    if err != nil {
        log.Fatal(err)
    }
    var aggResults []bson.M
    aggCursor.All(ctx, &aggResults)
    fmt.Println("Top authors:", aggResults)
}
```

---

## Summary

- MongoDB stores **documents** (JSON-like) in **collections** with no fixed schema.
- The **aggregation pipeline** is MongoDB's equivalent of SQL GROUP BY + JOINs — powerful and expressive.
- **Embedding** vs **referencing**: embed when data is always accessed together and bounded in size.
- Always create indexes on fields you filter/sort by; use `.explain()` to verify they're used.
- Transactions are available but have overhead — good schema design usually avoids needing them.

### Exercises

**Easy:** Use the MongoDB Go driver to insert 10 blog posts and find all posts by a specific author. Use `cursor.All()` to decode them into a slice.

**Medium:** Build an aggregation pipeline that finds the top 3 most-used tags across all posts. Each post has a `tags` array — you'll need `$unwind` to expand them first.

**Hard:** Design a MongoDB schema for a todo application where tasks can be nested (a task can have sub-tasks). Decide whether to embed or reference. Build a Go function that retrieves a task and all its sub-tasks efficiently.
