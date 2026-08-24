# Chapter 81: MongoDB with the Official Go Driver

MongoDB stores data as BSON documents (Binary JSON). Unlike PostgreSQL, there is no fixed schema — each document in a collection can have different fields. This is powerful for content with variable shapes (e-commerce products, user preferences, event logs) and awkward for anything that needs strong relational integrity.

## Table of Contents

1. [Connecting and BSON](#1-connecting-and-bson)
2. [CRUD Operations](#2-crud-operations)
3. [Querying and Filtering](#3-querying-and-filtering)
4. [Indexes](#4-indexes)
5. [Aggregation Pipeline](#5-aggregation-pipeline)
6. [Transactions](#6-transactions)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Connecting and BSON

```go
import (
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func newMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
    opts := options.Client().
        ApplyURI(uri).                         // "mongodb://localhost:27017"
        SetMaxPoolSize(20).
        SetMinPoolSize(5).
        SetConnectTimeout(10 * time.Second).
        SetServerSelectionTimeout(5 * time.Second)
    
    client, err := mongo.Connect(ctx, opts)
    if err != nil { return nil, err }
    
    if err := client.Ping(ctx, nil); err != nil {
        client.Disconnect(ctx)
        return nil, err
    }
    return client, nil
}

// Get a collection
db := client.Database("myapp")
coll := db.Collection("users")
```

### BSON and struct tags

```go
type Address struct {
    Street string `bson:"street"`
    City   string `bson:"city"`
    Zip    string `bson:"zip"`
}

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string             `bson:"name"`
    Email     string             `bson:"email"`
    Age       int                `bson:"age"`
    Address   *Address           `bson:"address,omitempty"`
    Tags      []string           `bson:"tags,omitempty"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// bson.D: ordered document (for queries and updates)
// bson.M: unordered map (simpler for quick filters)
// bson.A: array

filter := bson.M{"email": "alice@example.com"}
filter2 := bson.D{{"email", "alice@example.com"}, {"age", bson.M{"$gte": 18}}}
```

---

## 2. CRUD Operations

```go
// INSERT ONE
func (r *UserRepo) Create(ctx context.Context, user *User) error {
    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now().UTC()
    user.UpdatedAt = time.Now().UTC()
    
    result, err := r.coll.InsertOne(ctx, user)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) { return ErrEmailTaken }
        return fmt.Errorf("insert user: %w", err)
    }
    user.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

// INSERT MANY
func (r *UserRepo) CreateBulk(ctx context.Context, users []*User) error {
    docs := make([]interface{}, len(users))
    for i, u := range users {
        u.ID = primitive.NewObjectID()
        u.CreatedAt = time.Now().UTC()
        docs[i] = u
    }
    _, err := r.coll.InsertMany(ctx, docs)
    return err
}

// FIND ONE
func (r *UserRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
    var user User
    err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if errors.Is(err, mongo.ErrNoDocuments) { return nil, ErrNotFound }
    return &user, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if errors.Is(err, mongo.ErrNoDocuments) { return nil, ErrNotFound }
    return &user, err
}

// FIND MANY with cursor
func (r *UserRepo) ListByTag(ctx context.Context, tag string, limit int64) ([]*User, error) {
    opts := options.Find().
        SetSort(bson.D{{"created_at", -1}}).
        SetLimit(limit)
    
    cursor, err := r.coll.Find(ctx, bson.M{"tags": tag}, opts)
    if err != nil { return nil, err }
    defer cursor.Close(ctx)
    
    var users []*User
    if err := cursor.All(ctx, &users); err != nil { return nil, err }
    return users, nil
}

// UPDATE ONE
func (r *UserRepo) Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
    updates["updated_at"] = time.Now().UTC()
    result, err := r.coll.UpdateOne(ctx,
        bson.M{"_id": id},
        bson.M{"$set": updates},
    )
    if err != nil { return err }
    if result.MatchedCount == 0 { return ErrNotFound }
    return nil
}

// UPSERT
func (r *UserRepo) Upsert(ctx context.Context, email string, data bson.M) error {
    opts := options.Update().SetUpsert(true)
    _, err := r.coll.UpdateOne(ctx,
        bson.M{"email": email},
        bson.M{"$set": data, "$setOnInsert": bson.M{"created_at": time.Now()}},
        opts,
    )
    return err
}

// DELETE
func (r *UserRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
    result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil { return err }
    if result.DeletedCount == 0 { return ErrNotFound }
    return nil
}
```

---

## 3. Querying and Filtering

```go
// Comparison operators
bson.M{"age": bson.M{"$gte": 18, "$lte": 65}}        // 18 <= age <= 65
bson.M{"status": bson.M{"$in": []string{"active", "trial"}}} // status IN
bson.M{"status": bson.M{"$nin": []string{"banned"}}}  // status NOT IN
bson.M{"deleted_at": bson.M{"$exists": false}}         // field does not exist

// Logical operators
bson.M{"$and": []bson.M{
    {"age": bson.M{"$gte": 18}},
    {"email": bson.M{"$regex": "@example\\.com$"}},
}}
bson.M{"$or": []bson.M{
    {"plan": "pro"},
    {"trial_ends": bson.M{"$gt": time.Now()}},
}}

// Array operators
bson.M{"tags": bson.M{"$all": []string{"go", "backend"}}}   // all tags present
bson.M{"tags": bson.M{"$size": 3}}                          // exactly 3 tags
bson.M{"scores": bson.M{"$elemMatch": bson.M{"$gte": 90}}}  // any score >= 90

// Nested field query
bson.M{"address.city": "New York"}

// Text search (requires text index)
bson.M{"$text": bson.M{"$search": "distributed systems golang"}}

// Full-featured list with pagination
func (r *UserRepo) List(ctx context.Context, f UserFilter) ([]*User, int64, error) {
    filter := bson.D{}
    if f.Search != "" {
        filter = append(filter, bson.E{"$text", bson.M{"$search": f.Search}})
    }
    if f.Tag != "" {
        filter = append(filter, bson.E{"tags", f.Tag})
    }
    if f.MinAge > 0 {
        filter = append(filter, bson.E{"age", bson.M{"$gte": f.MinAge}})
    }
    
    total, err := r.coll.CountDocuments(ctx, filter)
    if err != nil { return nil, 0, err }
    
    opts := options.Find().
        SetSort(bson.D{{"created_at", -1}}).
        SetSkip(int64((f.Page - 1) * f.PageSize)).
        SetLimit(int64(f.PageSize))
    
    cursor, err := r.coll.Find(ctx, filter, opts)
    if err != nil { return nil, 0, err }
    defer cursor.Close(ctx)
    
    var users []*User
    return users, total, cursor.All(ctx, &users)
}
```

---

## 4. Indexes

```go
func (r *UserRepo) EnsureIndexes(ctx context.Context) error {
    indexes := []mongo.IndexModel{
        // Unique index on email
        {
            Keys:    bson.D{{"email", 1}},
            Options: options.Index().SetUnique(true),
        },
        // Compound index for common query pattern
        {
            Keys: bson.D{{"tags", 1}, {"created_at", -1}},
        },
        // Text index for full-text search
        {
            Keys:    bson.D{{"name", "text"}, {"bio", "text"}},
            Options: options.Index().SetName("text_search"),
        },
        // TTL index: auto-delete documents after 24 hours
        {
            Keys:    bson.D{{"expires_at", 1}},
            Options: options.Index().SetExpireAfterSeconds(0),
        },
        // Sparse index: only index documents where field exists
        {
            Keys:    bson.D{{"github_id", 1}},
            Options: options.Index().SetSparse(true).SetUnique(true),
        },
    }
    
    _, err := r.coll.Indexes().CreateMany(ctx, indexes)
    return err
}
```

---

## 5. Aggregation Pipeline

The aggregation pipeline is MongoDB's way of doing GROUP BY, JOIN, and transformation in a single query.

```go
// Group users by city, count and average age
func (r *UserRepo) StatsByCity(ctx context.Context) ([]CityStats, error) {
    pipeline := mongo.Pipeline{
        // Stage 1: filter (optional)
        bson.D{{"$match", bson.M{"status": "active"}}},
        // Stage 2: group
        bson.D{{"$group", bson.M{
            "_id":       "$address.city",
            "count":     bson.M{"$sum": 1},
            "avg_age":   bson.M{"$avg": "$age"},
            "min_age":   bson.M{"$min": "$age"},
        }}},
        // Stage 3: reshape the output
        bson.D{{"$project", bson.M{
            "city":    "$_id",
            "count":   1,
            "avg_age": bson.M{"$round": []any{"$avg_age", 1}},
            "_id":     0,
        }}},
        // Stage 4: sort
        bson.D{{"$sort", bson.D{{"count", -1}}}},
        // Stage 5: limit
        bson.D{{"$limit", 20}},
    }
    
    cursor, err := r.coll.Aggregate(ctx, pipeline)
    if err != nil { return nil, err }
    defer cursor.Close(ctx)
    
    var results []CityStats
    return results, cursor.All(ctx, &results)
}

// $lookup: join two collections
func (r *OrderRepo) OrdersWithUsers(ctx context.Context) ([]OrderWithUser, error) {
    pipeline := mongo.Pipeline{
        bson.D{{"$lookup", bson.M{
            "from":         "users",
            "localField":   "user_id",
            "foreignField": "_id",
            "as":           "user",
        }}},
        bson.D{{"$unwind", "$user"}}, // flatten the array to a single doc
        bson.D{{"$project", bson.M{
            "order_id":   "$_id",
            "amount":     1,
            "user.name":  1,
            "user.email": 1,
        }}},
    }
    
    cursor, err := r.coll.Aggregate(ctx, pipeline)
    if err != nil { return nil, err }
    defer cursor.Close(ctx)
    
    var results []OrderWithUser
    return results, cursor.All(ctx, &results)
}
```

---

## 6. Transactions

MongoDB supports multi-document ACID transactions (requires a replica set or sharded cluster):

```go
func (s *OrderService) PlaceOrder(ctx context.Context, order *Order) error {
    session, err := s.client.StartSession()
    if err != nil { return err }
    defer session.EndSession(ctx)
    
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (any, error) {
        // 1. Deduct inventory
        result, err := s.inventory.UpdateOne(sessCtx,
            bson.M{
                "product_id": order.ProductID,
                "quantity":   bson.M{"$gte": order.Quantity},
            },
            bson.M{"$inc": bson.M{"quantity": -order.Quantity}},
        )
        if err != nil { return nil, err }
        if result.MatchedCount == 0 { return nil, ErrInsufficientStock }
        
        // 2. Create order
        if _, err := s.orders.InsertOne(sessCtx, order); err != nil {
            return nil, err
        }
        return nil, nil
    })
    return err
}
```

---

## Summary

- **`bson.M`**: unordered map — great for simple filters
- **`bson.D`**: ordered slice of key-value pairs — use when order matters (aggregation stages)
- Always set `omitempty` on `ObjectID` fields so insertions auto-generate IDs
- **Indexes**: create them at startup with `EnsureIndexes(ctx)`; text indexes for full-text, TTL indexes for auto-expiry
- **Aggregation**: `$match` (filter) → `$group` / `$lookup` / `$project` → `$sort` → `$limit`
- **Transactions**: require replica set; use `session.WithTransaction` for automatic retry on transient errors
- Use `mongo.IsDuplicateKeyError(err)` to detect unique constraint violations

## Exercises

### Easy
1. Build a `ProductRepository` with `Create`, `GetByID`, `Update`, `Delete`, and `List` methods. Products have a `category` (string), `price` (float), and `attributes` (bson.M).
2. Write a query that finds all products in the "electronics" category with price between $100 and $500, sorted by price ascending.
3. Create a TTL index on a `sessions` collection where documents expire 24 hours after `created_at`.

### Medium
4. Implement a **tagging system**: products can have multiple tags. Write queries for: "find all products with tag X", "find products with all of [X, Y, Z]", "find the top 10 most-used tags using aggregation".
5. Build an **audit log** collection: every update to the `users` collection writes an audit event `{user_id, field, old_value, new_value, changed_at}` in a separate collection atomically using a transaction.
6. Write an aggregation that answers "what is the revenue per month for the last 12 months?". Use `$group` with `$dateToString` for month bucketing and `$sum` for totals.

### Hard
7. Implement a **change stream listener**: use `coll.Watch(ctx, pipeline)` to subscribe to all changes in the `orders` collection. When an order status changes to "shipped", trigger a notification. Handle resume tokens so the listener can restart without missing events.
8. Design a **polymorphic document model** where a single `events` collection stores different event types (PageView, Purchase, Signup). Each type has a `type` discriminator field and type-specific fields. Write a Go decoder that unmarshals the BSON into the correct concrete Go struct based on the `type` field.
