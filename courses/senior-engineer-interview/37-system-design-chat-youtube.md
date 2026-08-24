# Chapter 37: System Design — Chat System (WhatsApp) & Video Platform (YouTube)

These two designs test very different skills: chat requires real-time delivery with presence detection, while YouTube requires massive media ingestion, transcoding, and CDN delivery. Both are common at top company interviews.

## Table of Contents

1. [Design a Chat System (WhatsApp)](#1-design-a-chat-system-whatsapp)
2. [Design a Video Platform (YouTube)](#2-design-a-video-platform-youtube)

---

## 1. Design a Chat System (WhatsApp)

### 1.1 Requirements

```
Functional:
  - 1-on-1 messaging
  - Group chat (up to 100 members)
  - Message delivery receipts (sent, delivered, read)
  - Online/offline status (presence)
  - Message history (last 30 days on server, full history locally)
  
Non-functional:
  - 500M daily active users
  - Each user sends 40 messages/day
  - Low latency delivery (<500ms)
  - No message loss (at-least-once delivery)
```

### 1.2 Capacity Estimation

```
Messages:
  500M users × 40 messages/day = 20B messages/day
  20B ÷ 86,400 = ~230,000 messages/second
  
Storage:
  Each message ≈ 100 bytes (text only)
  20B × 100 bytes = 2TB/day
  30 days: 60TB of messages
  
Connections:
  500M DAU: at peak, ~50M concurrent WebSocket connections
  Need: connection servers with very high file descriptor limits
```

### 1.3 Key Design Decision: WebSockets for Real-Time

```
Options for real-time messaging:
  Short polling: client asks "any messages?" every N seconds — wasteful, laggy
  Long polling: client holds request open — better, but HTTP overhead per message
  WebSocket: persistent bidirectional TCP connection — best for chat
  SSE: server → client only — not enough for bidirectional chat

Choice: WebSocket for mobile/web clients

Scaling WebSocket connections:
  WebSocket servers are stateful (each server holds client connections)
  Problem: message for user B must reach B's WebSocket server (not just any server)
  
  Solution: message routing
    User B is connected to WS Server 3
    Service registry maps user_id → server_id (stored in Redis)
    When message arrives for B: lookup which WS server holds B's connection
    Forward message to that WS server via internal message bus (Redis pub/sub or Kafka)
```

### 1.4 High-Level Architecture

```
[Mobile App] ←──WebSocket──→ [WS Connection Server 1]
[Mobile App] ←──WebSocket──→ [WS Connection Server 2]
[Mobile App] ←──WebSocket──→ [WS Connection Server 3]
                                     │
                            [Redis Pub/Sub]
                            (route messages between WS servers)
                                     │
                          [Message Service]
                                     │
                    ┌────────────────┼────────────────┐
                    │                │                │
             [Cassandra DB]   [Media Service]  [Push Notification]
             (messages)       (images/videos)   (when offline)
```

### 1.5 Message Flow

```
Alice sends message to Bob:
1. Alice's app sends message via WebSocket to WS Server 1 (Alice's server)
2. WS Server 1 assigns message ID, writes to Cassandra (message store)
3. WS Server 1 sends ack to Alice: "message stored" (sent receipt)
4. WS Server 1 publishes to Redis channel for Bob's user_id
5. WS Server 3 (Bob's server) is subscribed to Bob's channel
6. WS Server 3 delivers message to Bob via WebSocket
7. Bob's app sends "delivered" receipt → "read" receipt back to Alice via same path

Offline user (Bob's phone is off):
4b. Redis pub/sub finds no active subscription for Bob
5b. Message Service sends push notification (FCM/APNs)
6b. When Bob comes online, his app fetches undelivered messages from Cassandra
```

### 1.6 Data Model for Messages

```sql
-- Cassandra is ideal: high write volume, time-series access pattern
-- Read: "give me all messages in conversation X after timestamp T"

CREATE TABLE messages (
    conversation_id UUID,
    message_id TIMEUUID,         -- includes timestamp, monotonic
    sender_id BIGINT,
    content TEXT,
    media_url TEXT,
    status TEXT,                 -- 'sent', 'delivered', 'read'
    PRIMARY KEY (conversation_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);

-- Fetch latest 50 messages in a conversation:
SELECT * FROM messages WHERE conversation_id = ? LIMIT 50;

-- User's conversations list (separate table for efficient lookup):
CREATE TABLE user_conversations (
    user_id BIGINT,
    conversation_id UUID,
    last_message_id TIMEUUID,
    last_message_text TEXT,
    unread_count INT,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id, updated_at)
) WITH CLUSTERING ORDER BY (updated_at DESC);
```

### 1.7 Presence (Online/Offline Status)

```
Presence is eventually consistent — it's OK if it's a few seconds stale

Implementation:
  On connect: SET user:{id}:online true EX 30 (Redis, 30s TTL)
  Heartbeat: client pings every 10s, server refreshes TTL
  On disconnect: key expires naturally after 30s (no active delete)
  
  Other users checking status:
  GET user:{id}:online → exists = online, null = offline
  
  Presence feed:
  When user comes online: publish to Redis channel "presence:{friend_id}"
  Friends' WS servers receive event, push to clients
```

---

## 2. Design a Video Platform (YouTube)

### 2.1 Requirements

```
Functional:
  - Upload video (any format, up to 10GB)
  - Transcode to multiple formats/resolutions (360p, 720p, 1080p, 4K)
  - Stream video to users (adaptive bitrate)
  - Video search
  - View count, likes, comments
  - Subscriptions and notifications

Non-functional:
  - 2B DAU
  - 500 hours of video uploaded every minute
  - Hundreds of millions of videos served daily
  - Videos must be durably stored forever
  - Upload and transcoding must not block streaming
```

### 2.2 Capacity Estimation

```
Uploads:
  500 hours/minute = 500 × 60 = 30,000 seconds of video/minute
  30,000 seconds × 1Mbps avg = 30,000 Mb/minute ≈ 3.75 GB/second of raw video ingested

Storage:
  After transcoding, each video is stored in ~5 resolutions
  Space per second of video: ~2MB (compressed across resolutions)
  30,000 seconds/minute × 2MB = 60GB/minute of stored video
  Per day: 60 × 60 × 24 = 86.4TB/day

Streaming:
  2B DAU, avg 30 min/day → 1B hours of video watched/day
  At avg 720p (1 Mbps): 1B × 3600 seconds × 1Mbps = 3.6 × 10^15 bits/day = petabits/day
  This is why CDN is not optional — it's the core infrastructure
```

### 2.3 Video Upload & Transcoding Pipeline

```
Raw upload:
  Large file → chunk it on client side (5MB chunks, resumable upload)
  Upload chunks to API server → stream to object storage (S3) raw bucket
  
Why chunked:
  Connection drops → resume from last chunk, not start over
  Parallel upload of multiple chunks = faster upload
  
After upload complete:
  1. API server publishes "video.uploaded" event to Kafka
  2. Transcoding workers pick up the event
  3. Each worker transcodes one resolution (360p, 720p, 1080p)
  4. Transcoded videos saved to S3 output bucket
  5. Thumbnail generated
  6. Metadata (title, description, duration) indexed in PostgreSQL
  7. Event published: "video.ready" → notification service sends email to creator

Transcoding farm:
  Stateless workers, horizontally scalable
  CPU-heavy: use specialized instances (c5 on AWS)
  Parallelism: different workers transcode different resolutions simultaneously
  FFmpeg is the standard tool for video transcoding
```

### 2.4 High-Level Architecture

```
UPLOAD PATH:
[Creator App] ──chunks──▶ [Upload API] ──▶ [S3 Raw Bucket]
                                │
                            [Kafka "uploads"]
                                │
                     [Transcoding Workers × N]
                          FFmpeg farm
                                │
                         [S3 Output Bucket]
                                │
                         [CDN Origin Pull]

PLAYBACK PATH:
[Viewer] ──▶ [CDN Edge Node] (serves 99% of traffic)
                │ (cache miss)
          [CDN Origin] ──▶ [S3 Output Bucket]

VIDEO METADATA:
[Viewer] ──▶ [API Gateway] ──▶ [Video Service] ──▶ [PostgreSQL + Redis]
```

### 2.5 Adaptive Bitrate Streaming (HLS/DASH)

```
Videos are not served as single files — they're segmented:
  1080p/: segment0.ts, segment1.ts, segment2.ts, ... (10-second chunks)
  720p/:  segment0.ts, segment1.ts, ...
  360p/:  segment0.ts, ...
  
Manifest file (playlist.m3u8):
  Lists available resolutions and segment URLs
  Player downloads manifest, picks resolution based on bandwidth
  Player downloads segments one by one, prefetching ahead
  
Adaptive bitrate:
  Player monitors download speed
  If segments take too long to download (buffering): switch to lower resolution
  If download is fast: switch to higher resolution
  
CDN delivers segments:
  Segment files are small (1-5MB each) and static
  Perfect for CDN caching: a video watched by 1M people → CDN caches segments
  99% of traffic served from CDN edge, not origin S3
```

### 2.6 View Count — Handling Extreme Scale

```
Problem: a viral video gets 10M views in 1 hour = 2,778 views/second
         Each view triggers an UPDATE view_count = view_count + 1
         Database can't handle this many point updates
         
Solution: approximate counting with eventual consistency
  In-memory counter per API server: increment local counter
  Periodic flush (every 30s): write aggregate to Redis INCR
  Background job: flush Redis counts to PostgreSQL every 5 minutes
  
  User sees: ~99% accurate, up to 5 minutes behind
  That's acceptable for view counts
  
For exact counts (e.g., monetization/reporting):
  Write all view events to Kafka
  Kafka consumer reads events and writes to analytics database
  Report generated from analytics DB (exact, but query is expensive)
```

### 2.7 Search

```
Video metadata stored in Elasticsearch:
  title, description, tags, transcript, creator_name
  
  Indexed fields:
  - Full-text search on title, description (with relevance scoring)
  - Filtered by: category, upload_date, duration, language
  
Ranking signals (beyond text relevance):
  - View count (viral videos rank higher)
  - Like/dislike ratio
  - Watch time (users watch 90% of video vs 10%)
  - Recency (fresh content gets a boost)
  - Creator subscriber count
  
These signals are pre-computed and stored as numeric fields in Elasticsearch
Combined with full-text relevance score using function_score queries
```

---

## Summary

### Chat System
- WebSockets for real-time bidirectional communication
- WS servers are stateful; use Redis pub/sub to route messages between them
- Cassandra for message storage (high write volume, time-series queries)
- Redis for presence tracking (TTL-based, refresh every heartbeat)
- Push notifications for offline users (FCM/APNs)

### Video Platform
- Chunked upload → S3 raw → Kafka → transcoding workers → S3 output → CDN
- Adaptive bitrate streaming (HLS/DASH): segmented files, player picks resolution
- CDN is the core infrastructure — serves 99% of video traffic
- View counts: approximate with in-memory + Redis, flush to DB periodically
- Search: Elasticsearch with text + ranking signals (views, watch time, recency)
