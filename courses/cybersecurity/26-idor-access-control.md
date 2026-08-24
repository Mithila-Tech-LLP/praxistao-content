# Chapter 26: IDOR and Broken Access Control — The Most Common Critical Bug

*IDOR (Insecure Direct Object Reference) is the #1 finding in bug bounty programs. It's simple: the server gives you access to an object if you just guess its ID — even if you shouldn't have access.*

---

## What is IDOR?

An IDOR exists when:
1. An object is referenced by a predictable identifier (number, UUID, etc.)
2. The server doesn't verify that the requesting user should have access to that object

```
User Alice fetches her profile:
GET /api/users/1001/profile
→ Returns Alice's data

Alice changes 1001 to 1002:
GET /api/users/1002/profile
→ Returns Bob's data  ← IDOR!
```

---

## Types of IDOR

### Simple Numeric IDOR

```
GET /invoices/5021    → your invoice
GET /invoices/5022    → someone else's invoice
GET /invoices/1       → invoice #1 (could be admin?)
```

### UUID-based IDOR

UUIDs are harder to guess, but:
- They may appear in other responses (friend lists, activity feeds)
- APIs might return them in bulk

```
GET /documents/550e8400-e29b-41d4-a716-446655440000
GET /documents/550e8400-e29b-41d4-a716-446655440001  ← sequential UUID variant
```

### Encoded IDOR

```
GET /profile?id=MTAwMQ==   → base64 decode → 1001
GET /profile?id=MTAwMg==   → base64 decode → 1002
```

---

## Testing for IDOR

### Manual Steps

1. Create two accounts (account A and account B)
2. Perform an action as A (create document, order, profile)
3. Capture the request in Burp Suite
4. Note the ID referencing your object
5. Swap to account B's session
6. Try accessing A's object ID from B's session

### Burp Suite Approach

```
1. Log in as user A
2. Navigate to your data: GET /api/orders/12345
3. Burp → Repeater → Change 12345 to 12344, 12346, etc.
4. Or: Use Intruder to iterate IDs automatically

Position: GET /api/orders/§12345§ HTTP/1.1
Payload: Numbers 1 to 10000
```

### IDOR in HTTP Methods

Don't just test GET — try POST, PUT, DELETE:

```
DELETE /api/posts/9999    → can you delete someone else's post?
PUT /api/users/1002/email → can you change someone else's email?
POST /api/payments/confirm?id=5678 → can you confirm someone else's payment?
```

### IDOR in Indirect References

```json
// API response includes other users' IDs
GET /api/messages
{
  "messages": [
    {"id": 1, "from_user": 1001, "to_user": 1002, "text": "..."},
    {"id": 2, "from_user": 1003, "to_user": 1001, "text": "..."}
  ]
}
// Now you know user IDs 1002 and 1003 exist
// Try: GET /api/users/1002/profile
```

---

## Go: Vulnerable vs Secure Implementation

```go
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
)

type Order struct {
    ID     int    `json:"id"`
    UserID int    `json:"user_id"`
    Amount int    `json:"amount"`
    Items  []string `json:"items"`
}

// VULNERABLE: No authorization check
func getOrderVulnerable(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderID, _ := strconv.Atoi(vars["id"])
    
    order := db.GetOrder(orderID)  // No check: is this the requester's order?
    if order == nil {
        http.Error(w, "Not found", 404)
        return
    }
    
    json.NewEncoder(w).Encode(order)  // Returns ANY user's order
}

// SECURE: Check ownership before returning
func getOrderSecure(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user from session
    session := getSession(r)
    if session == nil {
        http.Error(w, "Unauthorized", 401)
        return
    }
    
    vars := mux.Vars(r)
    orderID, _ := strconv.Atoi(vars["id"])
    
    order := db.GetOrder(orderID)
    if order == nil {
        http.Error(w, "Not found", 404)
        return
    }
    
    // CRITICAL: verify ownership
    if order.UserID != session.UserID {
        // Don't reveal the order exists
        http.Error(w, "Not found", 404)  // 404 not 403 — don't leak existence
        return
    }
    
    json.NewEncoder(w).Encode(order)
}

// Pattern: Always scope queries to the authenticated user
func listOrdersSecure(w http.ResponseWriter, r *http.Request) {
    session := getSession(r)
    if session == nil {
        http.Error(w, "Unauthorized", 401)
        return
    }
    
    // CRITICAL: WHERE user_id = authenticated_user_id
    orders := db.GetOrdersForUser(session.UserID)  // user can only see their orders
    json.NewEncoder(w).Encode(orders)
}
```

---

## Horizontal vs Vertical Privilege Escalation

**Horizontal:** Access another user's data at the same privilege level
- Alice accessing Bob's order (both regular users)
- Most IDOR bugs

**Vertical:** Access higher-privilege functionality
- Regular user accessing admin endpoints
- GET /admin/users/list without being admin

```bash
# Testing vertical escalation
# Log in as regular user, capture admin-looking endpoints:
GET /admin/dashboard      → 403 as user
GET /api/admin/users      → 403 as user
POST /api/users/1/promote → 403 as user

# Try with different HTTP methods, content types, headers
GET /api/admin/users → 403
POST /api/admin/users → 200?  (method matters!)

# Try with extra headers sometimes seen in internal access
X-Internal-Request: true
X-Forwarded-For: 127.0.0.1
```

---

## Mass Assignment

When an API accepts extra fields it shouldn't:

```javascript
// User updates their profile
PUT /api/users/1001
{"name": "Alice", "email": "alice@example.com"}

// What if you add an admin field?
PUT /api/users/1001
{"name": "Alice", "email": "alice@example.com", "is_admin": true, "balance": 999999}

// Vulnerable server just binds all fields to the model
```

```go
// VULNERABLE: bind all input to model
func updateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)  // attacker controls everything!
    db.Save(&user)
}

// SECURE: explicit field binding
func updateUserSecure(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        // is_admin NOT included — cannot be set this way
    }
    json.NewDecoder(r.Body).Decode(&input)
    
    session := getSession(r)
    user := db.GetUser(session.UserID)
    user.Name = input.Name
    user.Email = input.Email
    // user.IsAdmin stays unchanged
    db.Save(user)
}
```

---

## Summary

| Vulnerability | Test by | Fix by |
|---------------|---------|--------|
| IDOR (horizontal) | Swap IDs between accounts | Check ownership before returning |
| Vertical privesc | Access admin URLs as user | Role-based access control (RBAC) |
| Mass assignment | Add extra fields to POST/PUT | Explicitly list allowed fields |
| IDOR in methods | Try PUT/DELETE on others' IDs | Authorize all methods, not just GET |

---

## Exercises

1. On DVWA or a test app: find an IDOR in the order/profile system. Access another user's data.
2. Build a Go REST API with proper IDOR protection: orders belong to users, and no user can read another's.
3. Test the mass assignment vulnerability: build an API that's vulnerable, then fix it.
4. Write a Burp Intruder configuration that iterates numeric IDs in a response and flags when non-404 responses appear.
