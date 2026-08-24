# Chapter 38: System Design — Ride-Sharing (Uber) & Payment Processing (Stripe)

These are the most complex system design questions. Uber tests geospatial, real-time matching, and coordination. Stripe tests financial correctness, distributed systems reliability, and security.

## Table of Contents

1. [Design Ride-Sharing (Uber)](#1-design-ride-sharing-uber)
2. [Design Payment Processing (Stripe)](#2-design-payment-processing-stripe)

---

## 1. Design Ride-Sharing (Uber)

### 1.1 Requirements

```
Functional:
  - Rider requests a ride (picks up location, destination)
  - System finds nearby available drivers
  - Match rider to driver
  - Real-time location tracking (rider sees driver moving on map)
  - Fare estimation and payment at trip end
  - Driver acceptance/rejection of requests

Non-functional:
  - 5M trips/day globally
  - Real-time location updates: 5-second intervals from each driver
  - Matching latency: <10 seconds
  - High availability (3 nines = 8.7 hours downtime/year)
```

### 1.2 Capacity Estimation

```
Drivers:
  1M active drivers, each sending location every 5 seconds
  = 200,000 location updates/second (writes)
  
Ride requests:
  5M trips/day ÷ 86,400 = ~60 requests/second

Location storage:
  Driver location: driver_id + (lat, lng) + timestamp ≈ 50 bytes
  200,000/second × 50 bytes = 10MB/second of location data
  For real-time matching, we only need the CURRENT location, not history
  
Matching:
  Within 5 miles radius: typically 100-500 drivers
  Need to query "drivers within radius R of point (lat, lng)"
```

### 1.3 The Core Challenge: Geospatial Proximity Search

```
Problem: Find all available drivers within 5km of the rider
         with thousands of driver location updates per second

Naive approach: store (lat, lng) for each driver, query by lat/lng range
  SELECT * FROM drivers 
  WHERE lat BETWEEN 37.77-0.05 AND 37.77+0.05 
  AND lng BETWEEN -122.41-0.05 AND -122.41+0.05
  AND status = 'available'
  
  Problem: lat/lng range queries are expensive. Index can only be used for one dimension.

Better: Geohash (divides Earth into hierarchical rectangular cells)
  Each location maps to a string like "9q8yy9mf" (precision = length)
  "9q8yy9m" = 153m × 153m cell
  Nearby locations share a prefix!
  
  Query all drivers in adjacent geohash cells:
  geohash = "9q8yy9m"
  neighbors = geohash.Neighbors(geohash) // returns 8 adjacent cells
  query Redis: SRANDMEMBER drivers:9q8yy9m, drivers:9q8yy8n, ...
```

### 1.4 Location Storage Architecture

```
For real-time matching: Redis (fast reads, driver locations update every 5s)
For trip history/analytics: PostgreSQL + S3

Redis data structure:
  Key: "drivers:9q8yy9m"  (geohash prefix for a cell)
  Value: Sorted Set of driver_ids with score = timestamp
  
  ZADD drivers:9q8yy9m <timestamp> driver_123
  
  Driver location update:
  1. Compute new geohash from (lat, lng)
  2. If geohash changed (moved cells): SREM from old cell, ZADD to new cell
  3. If same cell: ZADD with updated timestamp (upsert)

Finding nearby drivers:
  1. Compute rider's geohash
  2. Get this cell + 8 neighbors
  3. SUNION of all 9 driver sets
  4. Calculate exact distance for each driver
  5. Return closest N available drivers
```

### 1.5 High-Level Architecture

```
[Driver App] ─── location updates ──▶ [Location Service] ──▶ [Redis Geohash]
                                              │
                                       [Kafka "driver-locations"]
                                              │
                                      [Location History DB]

[Rider App] ──── ride request ──▶ [Matching Service]
                                         │
                                   Queries Redis
                                   (nearby drivers)
                                         │
                               [Driver Selection Algorithm]
                                         │
                               [Trip Service] ──▶ [PostgreSQL]
                                         │
                               [Notification Service]
                               (push to driver app)
                                         │
                               Driver accepts/rejects
                                         │
                               [Driver tracking stream]
                               (rider sees driver on map)
```

### 1.6 Driver-Rider Matching Algorithm

```
Factors in matching:
  1. Distance (primary factor: minimize pickup time)
  2. Driver rating
  3. Vehicle type (UberX, UberXL, Black)
  4. ETA accuracy (account for traffic)

Algorithm:
  1. Find top 10 nearby available drivers by distance
  2. Calculate actual ETA using routing service (Google Maps API)
  3. Score each driver: score = ETA × rating_factor
  4. Send request to top-scored driver
  5. Driver has 15 seconds to accept/reject
  6. If rejected/timeout: move to next driver on list
  7. If all drivers reject: expand search radius, retry

Surge pricing:
  Monitor demand/supply ratio per geohash cell
  If (requests in last 5 min) / (available drivers) > threshold → surge
  Surge multiplier = f(demand/supply ratio)
```

### 1.7 Real-Time Driver Tracking for Riders

```
After matching:
  Rider needs to see driver location update every 3-5 seconds
  
  Option 1: Rider polls /trip/{id}/driver-location every 5s
    Simple but polling is wasteful
    
  Option 2: WebSocket connection from rider app to Location Service
    Rider opens WS after match
    Location Service pushes updates as driver location changes
    Efficient, real-time
    
  Option 3: Server-Sent Events (SSE)
    Server pushes updates, client can't send messages
    Simpler than WS, sufficient for one-directional tracking
    
Recommended: WebSocket (can also send messages like "driver is waiting")
```

---

## 2. Design Payment Processing (Stripe)

### 2.1 Requirements

```
Functional:
  - Create payment intents (charge a card)
  - Support multiple payment methods (card, ACH, bank transfer)
  - Payment state machine: pending → processing → succeeded/failed
  - Webhooks: notify merchants of payment events
  - Refunds
  - Idempotent API (safe to retry)

Non-functional:
  - Financial accuracy: no double charges, no lost payments
  - Exactly-once processing
  - Compliance: PCI-DSS (cardholder data security)
  - High availability: payments must work even during partial outages
  - Audit trail: every state change must be recorded
```

### 2.2 The Core Challenge: Financial Correctness

```
Two types of failures that are catastrophic:
  1. Double charge: customer charged twice for one purchase → legal liability
  2. Lost payment: merchant expects money that never arrives → revenue loss

Solution: idempotency keys + payment state machine

Payment state machine:
  CREATED → PENDING → PROCESSING → SUCCEEDED
                     ↓
                  FAILED
                     ↓
                  REFUNDED (if refund requested)

Every state transition is recorded with a timestamp:
  payments(id, amount, currency, status, idempotency_key, created_at)
  payment_events(id, payment_id, from_status, to_status, metadata, created_at)
```

### 2.3 Idempotency

```go
// Every payment API call must include an idempotency key
// The same key always produces the same result

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *CreatePaymentReq) (*Payment, error) {
    key := req.IdempotencyKey
    
    // Check if this key was already processed
    existing, err := h.db.GetByIdempotencyKey(ctx, key)
    if err == nil {
        return existing, nil // return cached result
    }
    
    // Not seen before: process the payment
    payment := &Payment{
        ID:             uuid.New(),
        Amount:         req.Amount,
        Currency:       req.Currency,
        Status:         "pending",
        IdempotencyKey: key,
    }
    
    // Write to database with UNIQUE constraint on idempotency_key
    // If concurrent request with same key arrives, one will get a unique violation
    // → return the first request's result
    err = h.db.CreatePayment(ctx, payment)
    if err != nil && isDuplicateKeyError(err) {
        existing, _ = h.db.GetByIdempotencyKey(ctx, key)
        return existing, nil
    }
    
    // Async: send to payment processor (bank)
    h.queue.Publish(ctx, "payments.process", payment.ID)
    
    return payment, nil
}
```

### 2.4 PCI-DSS Compliance

```
PCI-DSS (Payment Card Industry Data Security Standard) requirements:
  - Card numbers (PANs) must NEVER be stored in plaintext
  - Limited access to cardholder data
  - Audit logs for all access to card data
  - Network segmentation: card data servers on isolated network
  - Encryption in transit (TLS 1.2+) and at rest

Tokenization:
  Real card number: 4111111111111111 (dangerous)
  Token:            tok_abc123def456 (safe to store everywhere)
  
  How it works:
    Card number is sent directly to payment processor (Stripe, Adyen)
    Processor returns a token
    Your servers only ever store the token
    Token is used for future charges (never the actual card number)
    
  Stripe.js / Stripe Elements:
    Card fields render in a Stripe-hosted iframe
    Card data goes directly from browser → Stripe servers
    Your servers NEVER see the raw card number (even in transit)
```

### 2.5 Payment Processing Flow

```
1. Merchant's frontend sends card to Stripe.js
   Stripe.js tokenizes → returns payment_method token

2. Merchant's backend calls Stripe:
   POST /v1/payment_intents
   { amount: 2000, currency: "usd", payment_method: "pm_xxx", confirm: true }
   With idempotency key in header

3. Stripe:
   a. Validates request
   b. Sends authorization request to card network (Visa/Mastercard)
   c. Card network routes to issuing bank
   d. Bank authorizes (or declines)
   e. Stripe captures funds (or schedules capture)
   f. Returns result to merchant
   
4. Stripe sends webhook to merchant:
   POST https://merchant.com/stripe-webhook
   { type: "payment_intent.succeeded", data: { ... } }
   
5. Merchant fulfills order upon webhook receipt
   (NOT upon API response — webhooks are more reliable for async events)
```

### 2.6 Webhook Delivery

```
Webhooks must be delivered reliably:
  - At-least-once delivery (retry on failure)
  - Exponential backoff: retry after 1s, 2s, 4s, ... up to 72 hours
  - Idempotent processing: merchant endpoint must handle duplicate webhooks
  
Merchant verifies webhook authenticity:
  Stripe signs each webhook with HMAC-SHA256 using webhook secret
  Merchant verifies signature before processing
  
  // Stripe-Signature header contains: timestamp + hmac
  // Merchant checks: expected_sig == computed_sig AND timestamp is recent
  // Prevents replay attacks (timestamp check) and spoofing (HMAC check)
```

### 2.7 Architecture

```
[Browser] ──── Stripe.js tokenize ──▶ [Stripe Card Vault] (PCI-compliant zone)
                                                │
[Merchant API] ──── payment intent ──▶ [Stripe API Layer]
                                                │
                                    [Payment State Machine]
                                                │
                                   [Bank/Card Network Gateway]
                                                │
                              ┌────────────────┬┴────────────────┐
                              │                │                 │
                          [Webhook]      [Settlement]      [Reconciliation]
                          Delivery       (ACH batch)         Service
                              │
                   [Merchant Webhook Endpoint]
```

---

## Summary

### Uber
- Geohash for geospatial indexing: nearby cells share prefixes → query adjacent cells
- Redis sorted sets for driver locations: fast updates, fast neighborhood queries
- Matching: find top candidates → calculate ETA → dispatch to best driver → fallback chain
- Real-time tracking: WebSockets from driver app through Location Service to rider app

### Stripe
- Idempotency keys are critical: same key = same result, prevents double charges
- Payment state machine: every state transition logged, audit trail is sacred
- PCI-DSS: tokenization, Stripe Elements ensure raw card numbers never touch your servers
- Webhook delivery: at-least-once with exponential backoff; merchant must be idempotent
- HMAC signature verification: prevents spoofed webhooks
