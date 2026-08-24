# Chapter 123: Kibana and OpenSearch — Log Management

Kubernetes pods come and go. When a pod crashes, its logs vanish with it — unless you shipped them somewhere first. Centralized log management solves this: every log line from every container is collected, indexed, and made searchable. This chapter covers the full log pipeline from a Go service writing structured JSON logs, through a log shipper (Fluent Bit), into OpenSearch, and out through Kibana/OpenSearch Dashboards for search and visualization. Structured logs are not optional here — they are what makes the whole thing queryable.

## Table of Contents

1. [Why Centralized Logs?](#1-why-centralized-logs)
2. [The Log Pipeline](#2-the-log-pipeline)
3. [Structured Logging in Go](#3-structured-logging-in-go)
4. [Fluent Bit Configuration](#4-fluent-bit-configuration)
5. [OpenSearch Index Strategy](#5-opensearch-index-strategy)
6. [KQL and Lucene Queries](#6-kql-and-lucene-queries)
7. [Index Templates and Mappings](#7-index-templates-and-mappings)
8. [Index Lifecycle Management (ILM)](#8-index-lifecycle-management-ilm)
9. [OpenSearch Dashboards](#9-opensearch-dashboards)
10. [Grok Parsing for Unstructured Logs](#10-grok-parsing-for-unstructured-logs)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Centralized Logs?

In a local monolith, `kubectl logs pod-name` is enough. In production Kubernetes with dozens of services and many replicas, it breaks down immediately:

- **Pod crashes**: the container is gone, and so are its logs. You have no idea what it printed before it died.
- **Multiple replicas**: a request failed — but which of the five `order-service` pods handled it? You'd have to check each one.
- **Ephemeral containers**: debug containers launched with `kubectl debug` disappear the moment the session ends, taking any output with them.
- **No cross-service correlation**: a `request_id` flows through `order-service`, `payment-service`, and `inventory-service`. Without centralized logs, you have to query three separate pod log streams and mentally join them.

The solution is to ship every log line to a central store the moment it is written. Pods become stateless with respect to logs.

```
  pod-1 (order-service)  ──┐
  pod-2 (order-service)  ──┤
  pod-3 (user-service)   ──┼──►  Fluent Bit  ──►  OpenSearch  ──►  Kibana
  pod-4 (payment-svc)    ──┤
  pod-5 (payment-svc)    ──┘
```

Every pod writes to `stdout`. Fluent Bit, running as a DaemonSet on each node, tails the container log files and forwards them. OpenSearch indexes and stores them. Kibana lets you search and visualize.

---

## 2. The Log Pipeline

Here is the full path a log line travels from your Go application to your browser:

```
Go App (slog JSON)
       │
       │ stdout/stderr
       ▼
Fluent Bit (DaemonSet per node)
  - INPUT:  tail /var/log/containers/*.log
  - FILTER: add Kubernetes metadata (pod name, namespace, labels)
  - FILTER: parse JSON log fields
  - OUTPUT: OpenSearch HTTP API
       │
       ▼
OpenSearch Cluster
  - Index: app-logs-2025.01.15
  - Shards: 1 primary, 1 replica
       │
       ▼
Kibana / OpenSearch Dashboards
  - Discover: search and filter logs
  - Dashboard: visualize log volumes, error rates
  - Alerts: notify on error spikes
```

The components and their roles:

| Component | Role | Alternatives |
|-----------|------|-------------|
| Fluent Bit | Lightweight log shipper (DaemonSet) | Fluentd, Filebeat, Logstash |
| OpenSearch | Search and analytics engine | Elasticsearch |
| Kibana/OS Dashboards | Visualization UI | Grafana (with Loki plugin) |

Fluent Bit is preferred over Fluentd in Kubernetes because it uses far less memory (~450KB vs ~40MB). It is written in C and designed for embedded and container environments.

---

## 3. Structured Logging in Go

Unstructured logs are strings. Structured logs are documents. OpenSearch indexes documents — so the format of your log output directly determines what you can query.

Go 1.21 ships `log/slog` with built-in JSON output:

```go
import "log/slog"

// Production setup: JSON output to stdout
func newLogger() *slog.Logger {
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelInfo,
        AddSource: false,
    }))
}

// Every log line becomes a JSON document in OpenSearch
slog.Info("order placed",
    "request_id", reqID,
    "order_id",   order.ID,
    "user_id",    order.UserID,
    "total",      order.Total,
    "duration_ms", time.Since(start).Milliseconds(),
)
```

Each call to `slog.Info` produces one JSON line on stdout:

```json
{
  "time": "2025-01-15T14:32:01.123Z",
  "level": "INFO",
  "msg": "order placed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "order_id": "ord_789",
  "user_id": 42,
  "total": 149.99,
  "duration_ms": 47
}
```

This JSON maps directly to OpenSearch fields with no parsing step. Each key becomes an independently searchable, filterable, and aggregatable field. You can write a KQL query like `user_id: 42 AND duration_ms > 500` because the values are typed fields, not substrings inside a blob of text.

Compare with unstructured output:

```go
// Bad: unstructured — requires grok parsing to extract fields
log.Printf("order placed: user=%d total=%.2f duration=%dms", user.ID, total, ms)
// Output: "order placed: user=42 total=149.99 duration=47ms"
// OpenSearch sees one big string; you can't filter by user_id
```

With unstructured logs, `user=42` is a character sequence buried in a message string. OpenSearch has no way to know that `42` is a number representing a user — it is just text.

A useful pattern is to pass a logger through context so every handler automatically carries request-scoped fields:

```go
// Attach request fields once at the handler boundary
log := slog.With(
    "request_id", reqID,
    "service",    "order-service",
)

// All subsequent calls carry those fields automatically
log.Info("validating order", "items", len(req.Items))
log.Info("payment charged", "amount", req.Total)
log.Error("inventory check failed", "err", err)
```

---

## 4. Fluent Bit Configuration

Fluent Bit runs as a DaemonSet — one pod per node — and reads every container's log file from `/var/log/containers/`.

```ini
# fluent-bit.conf
[SERVICE]
    Flush         5
    Log_Level     info
    Parsers_File  parsers.conf

[INPUT]
    Name              tail
    Tag               kube.*
    Path              /var/log/containers/*.log
    Parser            docker
    DB                /var/log/flb_kube.db
    Mem_Buf_Limit     5MB
    Skip_Long_Lines   On
    Refresh_Interval  10

[FILTER]
    Name                kubernetes
    Match               kube.*
    Kube_URL            https://kubernetes.default.svc:443
    Kube_CA_File        /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    Kube_Token_File     /var/run/secrets/kubernetes.io/serviceaccount/token
    Kube_Tag_Prefix     kube.var.log.containers.
    Merge_Log           On    # parse JSON log body into top-level fields
    Keep_Log            Off
    Annotations         Off
    Labels              On

[FILTER]
    Name    grep
    Match   kube.*
    Exclude log ^$   # drop empty log lines

[OUTPUT]
    Name            opensearch
    Match           kube.*
    Host            opensearch.logging.svc.cluster.local
    Port            9200
    Index           app-logs
    Generate_ID     On
    Logstash_Format On          # creates daily indices: app-logs-2025.01.15
    Logstash_Prefix app-logs
    Time_Key        @timestamp
    Retry_Limit     5
```

Key configuration points:

- **`kubernetes` filter**: calls the Kubernetes API to enrich every log line with `kubernetes.pod_name`, `kubernetes.namespace_name`, `kubernetes.container_name`, and `kubernetes.labels.*`. You never have to add these fields in your Go code.
- **`Merge_Log On`**: if the raw log body is valid JSON (which it will be with `slog`), the filter merges its keys into the top-level document. Without this, your `user_id` would be nested inside a `log` field instead of being a top-level searchable field.
- **`Logstash_Format On`**: instructs Fluent Bit to create daily indices named `app-logs-YYYY.MM.DD` rather than writing everything into a single `app-logs` index. This is essential for ILM (section 8).
- **`DB`**: Fluent Bit tracks the file offset in a SQLite database so it resumes from the right position if it restarts. Without this, it would re-read files from the beginning.
- **`Parser docker`**: works for Docker-based nodes, which write JSON log files. Most modern Kubernetes clusters use containerd instead — for those, set `Parser cri` (or `multiline.parser cri` in newer Fluent Bit versions), because containerd writes logs in the CRI text format, not JSON.

---

## 5. OpenSearch Index Strategy

OpenSearch stores documents in indices. For log data the standard approach is daily indices.

**Naming convention**: `app-logs-2025.01.15`

The date suffix lets you:
- Delete old indices directly (`DELETE app-logs-2025.01.01`) without touching recent data
- Apply ILM policies that move or delete indices by age
- Run queries scoped to a date range (`app-logs-2025.01.*`)

Set the initial index configuration once via an index template (section 7). For the physical settings:

```json
{
  "settings": {
    "number_of_shards": 2,
    "number_of_replicas": 1,
    "refresh_interval": "5s"
  }
}
```

**Shard sizing rule**: aim for 10–50 GB per shard. A shard that is too large degrades query performance; too many tiny shards wastes overhead. For a low-traffic service producing 1 GB/day, a single shard is fine. A high-traffic service producing 100 GB/day needs at least 2–3 shards per daily index.

`refresh_interval: "5s"` means newly indexed documents become searchable within 5 seconds. The default is 1s, which adds unnecessary I/O for log workloads where near-real-time is good enough.

---

## 6. KQL and Lucene Queries

The Kibana/OpenSearch Dashboards search bar accepts KQL (Kibana Query Language) by default. KQL is designed to be readable; Lucene is the underlying query syntax for more advanced use.

**Basic KQL syntax**:

```
# Find errors for a specific service
service.name: "order-service" AND level: "ERROR"

# Find a specific request by ID
request_id: "550e8400-e29b-41d4-a716-446655440000"

# Range queries
duration_ms > 500

# Wildcard (use sparingly — slow on large indices)
msg: "order*"

# Boolean combinations
level: "ERROR" AND (service.name: "order-service" OR service.name: "payment-service")

# Exists check — field is present in the document
order_id: *
```

KQL vs Lucene comparison:

| Feature | KQL | Lucene |
|---------|-----|--------|
| Syntax | Simple, field: value | Classic query string |
| AND/OR | `AND` / `OR` (capitalized) | `AND` / `OR` or `+`/`-` |
| Range | `duration_ms > 500` | `duration_ms:[500 TO *]` |
| Wildcard | `service.name: order*` | `service.name:order*` |
| Nested | `kubernetes.namespace: prod` | `kubernetes.namespace:prod` |
| Default | KQL in newer Kibana/OSD | Lucene in older versions |

**Practical debugging queries**:

```
# 1. All errors in the last hour — set time range in the picker, then:
level: "ERROR"

# 2. Slow requests over 500ms
duration_ms > 500 AND level: "INFO"

# 3. All activity for a specific user
user_id: 42

# 4. Trace a request across services
request_id: "abc-123"
# finds it in order-service, payment-service, inventory-service — all in one view

# 5. 5xx errors by service
level: "ERROR" AND http.status_code >= 500
```

The `request_id` query in item 4 is the main payoff of structured logging. One search surfaces every log line from every service that touched that request.

---

## 7. Index Templates and Mappings

Without explicit mappings, OpenSearch uses dynamic mapping — it guesses field types from the first document it sees. If the first document has `user_id: "42"` (a string), OpenSearch maps the field as `text`. Later documents with `user_id: 42` (an integer) will still be indexed as text, breaking range queries and aggregations.

Index templates apply settings and mappings automatically to every new index that matches a pattern:

```json
{
  "index_patterns": ["app-logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 2,
      "number_of_replicas": 1
    },
    "mappings": {
      "properties": {
        "@timestamp":  { "type": "date" },
        "level":       { "type": "keyword" },
        "msg":         { "type": "text", "analyzer": "standard" },
        "service":     { "type": "keyword" },
        "request_id":  { "type": "keyword" },
        "user_id":     { "type": "long" },
        "order_id":    { "type": "keyword" },
        "duration_ms": { "type": "integer" },
        "http": {
          "properties": {
            "method":      { "type": "keyword" },
            "path":        { "type": "keyword" },
            "status_code": { "type": "integer" }
          }
        },
        "kubernetes": {
          "properties": {
            "pod_name":       { "type": "keyword" },
            "namespace_name": { "type": "keyword" },
            "container_name": { "type": "keyword" }
          }
        }
      }
    }
  }
}
```

The four mapping types you will use most:

| Type | Use for | Supports |
|------|---------|---------|
| `keyword` | IDs, service names, log levels, enum values | Exact match, aggregation, sorting |
| `text` | Log messages, descriptions | Full-text search with tokenization |
| `date` | Timestamps | Time-range queries, date histograms |
| `long` / `integer` | Counts, durations, status codes | Range queries, aggregations |

The rule of thumb: if you will filter by exact value (e.g., `level: "ERROR"`), use `keyword`. If you will search for words within the value (e.g., `msg: "connection refused"`), use `text`.

Apply the template via the OpenSearch REST API before the first log arrives:

```
PUT _index_template/app-logs-template
{ ...the JSON above... }
```

---

## 8. Index Lifecycle Management (ILM)

Logs are high-volume and time-bounded. You care a lot about yesterday's errors; you rarely need logs from 90 days ago. ILM automates moving indices through phases as they age, reducing storage cost automatically.

```json
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_age": "1d",
            "max_size": "50gb"
          },
          "set_priority": { "priority": 100 }
        }
      },
      "warm": {
        "min_age": "1d",
        "actions": {
          "readonly": {},
          "shrink":     { "number_of_shards": 1 },
          "forcemerge": { "max_num_segments": 1 },
          "set_priority": { "priority": 50 }
        }
      },
      "cold": {
        "min_age": "7d",
        "actions": {
          "readonly": {},
          "set_priority": { "priority": 0 }
        }
      },
      "delete": {
        "min_age": "90d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

What each phase does:

```
day 0 ──► HOT  (SSD, full replicas, active writes, 1 day or 50GB)
              │
day 1 ──► WARM (slower disk, read-only, merged to 1 shard, forcemerged)
              │
day 7 ──► COLD (object storage / cold disk, read-only, lowest priority)
              │
day 90 ─► DELETE
```

- **Hot**: active writes, stored on fast SSD, full replica count. Rollover triggers when the index is 1 day old or 50 GB, whichever comes first.
- **Warm**: index becomes read-only. Shrink reduces it to 1 shard (no point keeping multiple shards for a finished day). Forcemerge compacts the Lucene segments, reducing file handles and speeding up read queries.
- **Cold**: moved to cheaper storage. Still queryable but slow. Replicas may be dropped to 0 to save space.
- **Delete**: index is permanently removed. Adjust `min_age` based on compliance requirements — some industries require 1 year or more.

Apply the policy with `PUT _ilm/policy/app-logs-ilm` and attach it to the index template's settings.

> **Naming note**: "ILM" (Index Lifecycle Management) and the JSON format above are Elasticsearch terminology. OpenSearch's built-in equivalent is called **ISM** (Index State Management) — the concept is identical (states with transitions and actions), but the API is `PUT _plugins/_ism/policies/app-logs-policy` and the policy JSON uses `states` and `transitions` instead of `phases`. If you run OpenSearch, use ISM; if you run Elasticsearch with Kibana, use ILM.

---

## 9. OpenSearch Dashboards

OpenSearch Dashboards (the open-source fork of Kibana) has four areas useful for log management:

**Discover**

The primary day-to-day debugging tool. Set a time range (top right), type a KQL query in the search bar, and see matching log documents in a table. Click any document to expand it and see every field — including the Kubernetes metadata Fluent Bit injected. Export results to CSV for offline analysis.

**Typical debugging workflow**:

1. Set time range to the last 1 hour
2. Add filter: `level: "ERROR"`
3. Spot the time when errors spiked from the histogram
4. Narrow time range to that spike
5. Add filter: `kubernetes.namespace_name: "production"`
6. Identify the failing service from `service` field
7. Copy a `request_id` from one of the error documents
8. Clear level filter, add `request_id: "<that-id>"` — now see the complete trace of that request across all services, in chronological order

**Dashboard**

Build reusable panels from saved searches and aggregations. Useful panels for log monitoring:

- **Bar chart**: error count per service over time (X axis: time, Y axis: count, split by `service` keyword)
- **Pie chart**: log level distribution — what fraction of logs are ERROR vs WARN vs INFO
- **Data table**: top 10 slowest requests (terms aggregation on `request_id`, max of `duration_ms`)
- **TSVB (Time Series Visual Builder)**: log ingestion rate over time — useful to spot when a service suddenly goes silent (it may have crashed)

**Alerting**

Define monitors that query OpenSearch on a schedule and fire when a condition is met. Example: if `COUNT(level: "ERROR") > 100` in a 1-minute window, send a webhook to Slack. Alerting is available natively in OpenSearch Dashboards under the Alerting plugin.

---

## 10. Grok Parsing for Unstructured Logs

Not every service you inherit will use structured logging. Legacy applications often write freeform strings. Fluent Bit can extract fields from these using regex parsers (Grok-style).

```ini
# parsers.conf

[PARSER]
    Name         nginx_access
    Format       regex
    Regex        ^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*?)(?: +\S*)?)?" (?<code>[^ ]*) (?<size>[^ ]*)
    Time_Key     time
    Time_Format  %d/%b/%Y:%H:%M:%S %z

[PARSER]
    Name         go_log
    Format       regex
    Regex        ^(?<time>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (?<level>\w+) (?<message>.*)$
    Time_Key     time
    Time_Format  %Y/%m/%d %H:%M:%S
```

Apply the parser in Fluent Bit with a filter:

```ini
[FILTER]
    Name     parser
    Match    kube.*
    Key_Name log
    Parser   nginx_access
```

`Key_Name log` tells Fluent Bit to apply the parser to the `log` field (the raw log line). If the regex matches, the named capture groups become top-level fields in the document.

**Limitations of grok parsing**:

- Brittle: any change to the log format breaks the regex silently — logs still ship, but fields are missing
- Slow: regex parsing at high volume consumes measurable CPU
- Maintenance burden: every new log format needs a new parser and testing

Use grok as a migration aid while you add structured logging to the legacy service. Once the service emits JSON, remove the parser and rely on `Merge_Log On`.

---

## Summary

- Pods are ephemeral — centralized log shipping is required in any serious Kubernetes deployment
- The standard pipeline is: Go (`slog` JSON) → stdout → Fluent Bit (DaemonSet) → OpenSearch → Kibana/OS Dashboards
- `log/slog` with `NewJSONHandler` produces one JSON document per log call; each field is independently queryable in OpenSearch
- Fluent Bit's `kubernetes` filter enriches every log with pod name, namespace, and labels automatically
- `Merge_Log On` in Fluent Bit merges JSON log bodies into top-level document fields — essential for structured logs to be queryable
- Daily index naming (`app-logs-YYYY.MM.DD`) via `Logstash_Format On` enables easy deletion and ILM management
- Set explicit index mappings via index templates before logs arrive; dynamic mapping guesses types and gets it wrong
- `keyword` for exact-match fields, `text` for full-text search, `date` for timestamps, `integer`/`long` for numbers
- ILM automates the hot → warm → cold → delete lifecycle, reducing storage cost without manual cleanup
- KQL queries like `request_id: "abc-123"` across all services in Discover are the practical payoff of structured logging
- Grok parsing works for legacy unstructured logs but is brittle; migrate to structured JSON as soon as possible

---

## Exercises

### Easy

1. Modify an existing Go HTTP handler to use `slog.NewJSONHandler` and log a request with at least five fields (`request_id`, `method`, `path`, `status_code`, `duration_ms`). Print the output to stdout and verify the JSON is valid with `jq`.
2. Write five KQL queries for an `app-logs-*` index: one for errors, one for a specific user, one range query on `duration_ms`, one to find a specific `request_id`, and one to find all 5xx responses.
3. Write a minimal Fluent Bit `[INPUT]` and `[OUTPUT]` config that tails a local log file (`/tmp/app.log`) and writes matching lines to stdout. Test it by appending a JSON line to the file and watching Fluent Bit forward it.

### Medium

4. Write an OpenSearch index template for `app-logs-*` with correct mappings for `@timestamp` (date), `level` (keyword), `msg` (text), `user_id` (long), `duration_ms` (integer), and a nested `http` object with `method`, `path`, and `status_code`. Apply it to a local OpenSearch instance and verify via `GET _index_template/app-logs-template`.
5. Create an ILM policy named `app-logs-ilm` that moves indices to warm after 1 day (shrink to 1 shard, forcemerge), cold after 7 days, and deletes them after 30 days. Attach the policy to your index template and verify the policy is applied to a new index.
6. In OpenSearch Dashboards, build a dashboard with two panels: a bar chart showing error count per service over the last 24 hours (split series by `service` keyword), and a data table showing the top 10 requests by `duration_ms`. Save the dashboard as "Service Health".

### Hard

7. Build a complete local log pipeline using Docker Compose: a Go service that writes structured JSON logs on a 1-second interval, a Fluent Bit container that reads from the Go service's log file and ships to OpenSearch, and an OpenSearch + Dashboards container pair. Verify that logs appear in Discover within 10 seconds of being written, and that `user_id` and `duration_ms` are indexed as the correct numeric types.
8. Take an nginx access log file with the combined log format and write a Fluent Bit `parsers.conf` that extracts `remote`, `method`, `path`, `status_code` (as integer), and `response_size` (as integer). Ship the parsed logs to OpenSearch and write a KQL query that returns all requests where `status_code >= 400`. Verify that `status_code > 400` (a range query, which only works on numeric types) returns the same results as `status_code: 400 OR status_code: 404 OR ...` — confirming the field was mapped as integer and not keyword.
9. Implement an OpenSearch alerting monitor that queries the `app-logs-*` index every minute, counts documents where `level: "ERROR"` in the last 60 seconds, and fires a webhook notification if the count exceeds 100. Configure the webhook to post a message to a Slack incoming webhook URL. Write a Go script that generates 150 error log lines in under 10 seconds to trigger the alert, and confirm the Slack message arrives.
