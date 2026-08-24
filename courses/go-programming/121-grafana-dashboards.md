# Chapter 121: Grafana — Dashboards, Panels, and Alerts

Metrics without a good dashboard are just numbers in a database. Grafana turns Prometheus metrics into dashboards that make it immediately obvious when something is wrong. This chapter covers PromQL from scratch, the anatomy of a Grafana dashboard, panel types, reusable variables, and how to set up alerts that page you at the right time — not for every blip, but for sustained problems that need human intervention.

## Table of Contents

1. [Grafana Overview](#1-grafana-overview)
2. [PromQL Fundamentals](#2-promql-fundamentals)
3. [The RED and USE Methods](#3-the-red-and-use-methods)
4. [Dashboard Structure](#4-dashboard-structure)
5. [Panel Types](#5-panel-types)
6. [Dashboard JSON](#6-dashboard-json)
7. [Dashboard Variables](#7-dashboard-variables)
8. [Grafana Alerts](#8-grafana-alerts)
9. [Alert State Machine](#9-alert-state-machine)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Grafana Overview

Grafana is a visualization layer. It does not store data — it connects to data sources and renders what they return. The most common data sources for Go services are Prometheus (metrics), Loki (logs), and Tempo or Jaeger (traces).

```
  Go App  ──metrics──►  Prometheus  ──────►  Grafana (dashboards + alerts)
  Go App  ──logs────►   Loki        ──────►  Grafana (logs explorer)
  Go App  ──traces──►   Tempo       ──────►  Grafana (trace viewer)
                                             │
                                             └──alerts──► Slack / PagerDuty
```

All three signals flow into a single Grafana instance. That means when an alert fires on a metric spike, you can pivot directly to the log explorer filtered to that time window, then drill into a trace for the slowest request — without leaving the tool.

Key concepts:

| Concept | What it is |
|---------|------------|
| **Data source** | Connection to Prometheus, Loki, a database, etc. |
| **Dashboard** | A named collection of panels, shareable by URL |
| **Panel** | A single visualization (graph, stat, table) |
| **Query** | PromQL (or LogQL, TraceQL) expression driving a panel |
| **Variable** | A template placeholder like `$environment` or `$service` |
| **Alert rule** | A PromQL condition that fires a notification when violated |

---

## 2. PromQL Fundamentals

Prometheus stores time series. Every metric is identified by its name plus a set of labels. PromQL lets you filter, aggregate, and compute rates over those series.

```promql
# --- Counters ---
# Rate of HTTP requests per second over last 5 minutes
rate(http_requests_total[5m])

# Total error rate (5xx) as a fraction of all requests
rate(http_requests_total{status=~"5.."}[5m])
  /
rate(http_requests_total[5m])

# Requests per second broken down by path
sum by (path) (rate(http_requests_total[5m]))

# --- Histograms ---
# P99 request latency
histogram_quantile(0.99,
  rate(http_request_duration_seconds_bucket[5m])
)

# P50, P95, P99 in one query (use in table panel)
histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# --- Gauges ---
# Current active connections (gauge: read directly, no rate())
active_connections

# Memory usage as percentage of limit
container_memory_usage_bytes / container_memory_limit_bytes * 100

# --- Vector matching ---
# Error ratio per service (when you have a 'service' label)
sum by (service) (rate(http_requests_total{status=~"5.."}[5m]))
  /
sum by (service) (rate(http_requests_total[5m]))
```

Core functions to memorize:

| Function | When to use | Notes |
|----------|-------------|-------|
| `rate(m[5m])` | Counters — per-second rate of increase | Never call `rate()` on a gauge |
| `histogram_quantile(0.99, ...)` | Latency percentiles from histograms | Always wrap the inner metric with `rate()` |
| `sum by (label)` | Aggregate series, keep one label dimension | Use `without (label)` to drop a single label instead |
| `{label=~"regex"}` | Filter by label using regex | `=~` matches, `!~` excludes |
| `increase(m[1h])` | Total count over a window (not per-second) | Useful for "requests in the last hour" stat panels |

---

## 3. The RED and USE Methods

Two complementary frameworks tell you what to graph. Use both — RED for every service your users touch, USE for every infrastructure resource that service depends on.

**RED** — for every user-facing service:

| Signal | What to measure | PromQL |
|--------|-----------------|--------|
| Rate | Requests per second | `sum(rate(http_requests_total[5m]))` |
| Errors | Error rate (%) | `100 * sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))` |
| Duration | P99 latency | `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` |

**USE** — for every resource (CPU, memory, database connection pool, disk):

| Signal | What to measure |
|--------|-----------------|
| Utilization | Percentage of time the resource is busy |
| Saturation | Work queued or waiting (run queue depth, pool wait time) |
| Errors | Error events from the resource (disk errors, OOM kills, pool timeouts) |

RED answers "are my users having a bad experience?" USE answers "why is the service struggling?" A practical dashboard combines both: RED panels in the top row, USE panels below them.

---

## 4. Dashboard Structure

A dashboard is divided into rows. Each row groups panels by concern.

```
╔══════════════════════════════════════════════════════════════════╗
║  Dashboard: Order Service                    [env: prod] [5m]    ║
╠══════════════════════════════════════════════════════════════════╣
║  ─── Row: Traffic ─────────────────────────────────────────────  ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           ║
║  │ Request Rate │  │  Error Rate  │  │ P99 Latency  │           ║
║  │  (time ser.) │  │  (time ser.) │  │  (time ser.) │           ║
║  └──────────────┘  └──────────────┘  └──────────────┘           ║
║  ─── Row: Saturation ───────────────────────────────────────────  ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           ║
║  │  Active Conn │  │  DB Pool     │  │  CPU Usage   │           ║
║  │  (stat)      │  │  (gauge)     │  │  (time ser.) │           ║
║  └──────────────┘  └──────────────┘  └──────────────┘           ║
╚══════════════════════════════════════════════════════════════════╝
```

- **Rows** collapse and expand. Collapse rows you don't need during an incident to reduce noise.
- **Panels** are positioned on a 24-column grid. Each panel has `x`, `y`, `w`, `h` coordinates.
- **Time range selector** (top right) affects every panel simultaneously. During an incident, narrow the range to zoom in on the event window.
- **Refresh interval** (e.g. `30s`) controls how often all panels re-query their data sources.

Suggested row structure for a typical Go service:

```
Row 1: Traffic       → Rate, Error Rate, P99 Latency
Row 2: Saturation    → Active connections, DB pool, goroutine count
Row 3: Errors        → Error count by type, 4xx vs 5xx breakdown
Row 4: Business      → Orders placed, payments processed, queue depth
```

---

## 5. Panel Types

Choosing the wrong panel type makes data harder to read. Match the panel to the shape of the question.

| Panel Type | Best For | Example |
|------------|----------|---------|
| **Time series** | Trends and patterns over time | Request rate, P99 latency over 24h |
| **Stat** | A single current value with a threshold color | Current error rate: 0.3% |
| **Gauge** | Percentage or utilization with a visual fill | DB pool usage: 72% |
| **Bar gauge** | Ranked list comparison | Top 5 slowest endpoints by P99 |
| **Table** | Per-dimension breakdown, multiple columns | P50 / P95 / P99 per endpoint |
| **Heatmap** | Distribution changing over time | Latency bucket spread — reveals bimodal distributions |
| **Logs** | Raw log lines from Loki | Jump from a metric spike into correlated log lines |

For most RED dashboards you will use time series for Rate / Duration, stat for current Error Rate, and a table for per-endpoint breakdown. Heatmaps are the best way to spot latency distribution shifts — a jump from unimodal to bimodal shows up immediately even when the P99 looks stable.

---

## 6. Dashboard JSON

Every Grafana dashboard is JSON under the hood. You can export, version-control, and import dashboards as code. Here is a minimal but complete example for an order service:

```json
{
  "uid": "order-service-overview",
  "title": "Order Service",
  "refresh": "30s",
  "time": { "from": "now-1h", "to": "now" },
  "panels": [
    {
      "id": 1,
      "type": "timeseries",
      "title": "Request Rate",
      "gridPos": { "x": 0, "y": 0, "w": 6, "h": 8 },
      "targets": [
        {
          "expr": "sum(rate(http_requests_total{service=\"order-service\"}[5m]))",
          "legendFormat": "req/s"
        }
      ]
    },
    {
      "id": 2,
      "type": "timeseries",
      "title": "Error Rate",
      "gridPos": { "x": 6, "y": 0, "w": 6, "h": 8 },
      "targets": [
        {
          "expr": "sum(rate(http_requests_total{service=\"order-service\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{service=\"order-service\"}[5m]))",
          "legendFormat": "error rate"
        }
      ]
    },
    {
      "id": 3,
      "type": "timeseries",
      "title": "P99 Latency",
      "gridPos": { "x": 12, "y": 0, "w": 6, "h": 8 },
      "targets": [
        {
          "expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket{service=\"order-service\"}[5m]))",
          "legendFormat": "p99"
        }
      ]
    },
    {
      "id": 4,
      "type": "stat",
      "title": "Active Connections",
      "gridPos": { "x": 18, "y": 0, "w": 6, "h": 8 },
      "targets": [
        {
          "expr": "active_connections{service=\"order-service\"}",
          "legendFormat": "connections"
        }
      ]
    }
  ]
}
```

Key fields:

- `uid` — stable identifier used in URLs and API calls. Set this explicitly so the dashboard URL never changes.
- `gridPos` — position on the 24-column grid. `w: 6` means the panel occupies one quarter of the row width.
- `targets[].expr` — the PromQL query. One panel can have multiple targets (multiple lines on the same graph).
- `targets[].legendFormat` — the label shown in the legend. Use `{{label_name}}` to interpolate a Prometheus label value.

Store dashboards in your repository under `grafana/dashboards/` and load them via Grafana's provisioning system or the `grafana/grafana` Helm chart's `dashboards` values.

---

## 7. Dashboard Variables

Variables make a single dashboard work across all environments and services. Instead of hardcoding `service="order-service"` in every query, you use `$service` and let the user pick from a dropdown.

Define a `$environment` variable with fixed values:

```json
{
  "name": "environment",
  "type": "custom",
  "options": [
    { "text": "prod",    "value": "prod" },
    { "text": "staging", "value": "staging" },
    { "text": "dev",     "value": "dev" }
  ]
}
```

Define a `$service` variable that queries Prometheus for all label values:

```json
{
  "name": "service",
  "type": "query",
  "datasource": "prometheus",
  "query": "label_values(http_requests_total, service)"
}
```

Then use the variables in every PromQL expression:

```promql
rate(http_requests_total{environment="$environment", service="$service"}[5m])
```

```promql
histogram_quantile(0.99,
  rate(http_request_duration_seconds_bucket{
    environment="$environment",
    service="$service"
  }[5m])
)
```

When the user changes the `$environment` dropdown to `staging`, every panel in the dashboard re-queries with `environment="staging"` simultaneously. This means one dashboard serves all environments, all services — no copy-paste sprawl.

Variables also support multi-value selection (useful for comparing two services side-by-side) and `All` as a special value that expands to a regex matching all label values.

---

## 8. Grafana Alerts

Grafana unified alerting lets you define alert rules as YAML and store them alongside your dashboards in version control. Rules are evaluated against Prometheus on a configurable interval.

```yaml
# grafana/alerts/order-service.yaml
apiVersion: 1
groups:
  - orgId: 1
    name: order-service
    folder: Services
    interval: 1m
    rules:
      - uid: order-p99-high
        title: "Order Service P99 Latency > 1s"
        condition: C
        data:
          - refId: A
            queryType: ""
            relativeTimeRange:
              from: 600
              to: 0
            datasourceUid: prometheus
            model:
              expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket{service="order-service"}[5m]))
              intervalMs: 1000
              maxDataPoints: 43200
              refId: A
          - refId: C
            datasourceUid: __expr__
            model:
              conditions:
                - evaluator:
                    params: [1.0]
                    type: gt
                  operator:
                    type: and
                  query:
                    params: [A]
              refId: C
              type: classic_conditions
        for: 5m
        labels:
          severity: warning
          service: order-service
        annotations:
          summary: "P99 latency is above 1s for 5 minutes"
          runbook: "https://wiki.internal/runbooks/order-service-latency"
```

The rule has two parts: query `A` fetches the metric from Prometheus, and expression `C` is the threshold evaluator (a special `__expr__` data source built into Grafana). The `for: 5m` field prevents a single transient spike from firing a page.

**Contact points and routing:**

| Destination | Use for | Config needed |
|-------------|---------|---------------|
| Slack | `severity: warning` | Incoming webhook URL |
| PagerDuty | `severity: critical` | Routing integration key |
| Email | Team digests | SMTP settings in Grafana |

Notification policies route alerts to contact points based on label matchers. A typical policy:

```
root policy → catch-all → Slack #ops-alerts
  └── severity=critical → PagerDuty (team on-call)
  └── service=payment-service, severity=warning → Slack #payments-team
```

Keep `severity: warning` for "this needs attention in the next hour" and `severity: critical` for "wake someone up now." Routing them to different destinations enforces that contract.

---

## 9. Alert State Machine

Every Grafana alert rule moves through a defined set of states. Understanding the states prevents confusion about why an alert hasn't fired yet.

```
                    threshold crossed
      OK  ──────────────────────────────►  Pending
                                              │
                          still violated      │  threshold OK again
                          for 'for' duration  │  before 'for' duration
                                              ▼
      OK  ◄──────────────────────────────  Alerting
           resolved (threshold OK again)
```

| State | Meaning | Notification sent? |
|-------|---------|-------------------|
| **OK** | Metric is within acceptable range | Yes — "resolved" when returning from Alerting |
| **Pending** | Threshold crossed; waiting for `for` duration | No |
| **Alerting** | Threshold violated for the full `for` duration | Yes — fires the alert |
| **No Data** | Query returned no data points | Configurable — treat as Alerting for uptime checks |

The `Pending` state is what separates a well-tuned alert from a noisy one. A spike that lasts 90 seconds with `for: 5m` set never pages anyone. The same spike sustained for six minutes does. Set `for` to something meaningful: for latency alerts, 2–5 minutes is typical; for disk-full warnings, 10–15 minutes is reasonable since there's no instant fix anyway.

**No Data** deserves special attention. If your query returns nothing because the service is completely down, you want that treated as Alerting, not silently ignored. Configure `execErrState: Alerting` and `noDataState: Alerting` for uptime-critical rules.

---

## Summary

- Grafana connects to Prometheus, Loki, and Tempo — one tool for all three observability signals.
- `rate()` is for counters, direct read is for gauges, and `histogram_quantile()` computes percentiles from histograms.
- The RED method (Rate, Errors, Duration) gives you the core panels for every user-facing service.
- The USE method (Utilization, Saturation, Errors) covers infrastructure resources like CPU and connection pools.
- Dashboards are 24-column grids of panels grouped into rows by concern: Traffic, Saturation, Errors, Business.
- Panel type should match the question: time series for trends, stat for current value, table for per-dimension breakdown, heatmap for distribution shifts.
- Dashboards are JSON — store them in version control, provision them via Helm or Grafana's provisioning API.
- Variables (`$environment`, `$service`) make one dashboard reusable across all environments without copy-paste.
- Alert rules consist of a PromQL query plus a threshold condition. The `for` duration filters out transient spikes.
- The alert state machine has four states: OK → Pending → Alerting → OK. Notifications fire when a rule enters Alerting and again when it resolves back to OK.
- Route `severity: warning` to Slack and `severity: critical` to PagerDuty to enforce the "wake someone up" threshold.

---

## Exercises

### Easy

1. Write the four core PromQL queries for a new `payment-service`: request rate, error rate as a fraction, P99 latency from a histogram, and current active goroutines (a gauge).
2. For each of the following metrics, identify which Grafana panel type fits best: (a) CPU utilization over the last 24 hours, (b) current DB connection pool usage as a percentage, (c) P50 / P95 / P99 latency broken down by endpoint, (d) total orders processed in the last hour (single number).
3. Create a basic Grafana dashboard JSON with three panels for a service of your choice: request rate (time series), error rate (time series), and active connections (stat). Include the `uid`, `title`, `refresh`, and `time` fields at the top level.

### Medium

4. Build a complete dashboard JSON for an `inventory-service` with eight panels organized into two rows: a Traffic row (request rate, error rate, P99 latency, P50 latency) and a Saturation row (active connections, DB pool utilization as a gauge, goroutine count, memory usage percentage). Use `gridPos` coordinates so all panels fit in a 24-column layout.
5. Add `$environment` (custom variable with `prod`, `staging`, `dev`) and `$service` (query variable from `label_values(http_requests_total, service)`) to a dashboard JSON. Update all panel queries to use both variables. Show the full `templating.list` array in the JSON.
6. Write a Grafana unified alerting YAML file with two rules for a `payment-service`: (a) error rate above 1% for 5 minutes with `severity: critical`, and (b) P99 latency above 500ms for 10 minutes with `severity: warning`. Include `labels`, `annotations`, and a `runbook` URL in each rule.

### Hard

7. Implement a Go function `GenerateDashboard(svc ServiceMeta) []byte` that takes a `ServiceMeta` struct (service name, environment, list of endpoint paths) and generates a valid Grafana dashboard JSON. The output should include a Traffic row with RED panels for each endpoint path (using `sum by (path)` queries) and a Saturation row with standard USE panels. Write the struct definitions and the JSON marshaling logic using `encoding/json` — no string templates.
8. Design and implement an SLO dashboard for a service with a 99.9% monthly availability target. The dashboard should include: (a) current error budget remaining as a stat panel (compute from `1 - error_rate` over the current month window), (b) error budget burn rate over 1-hour and 6-hour windows as time series panels, and (c) a Grafana alert that fires when the 1-hour burn rate exceeds 14.4x (meaning the entire 30-day budget would be exhausted in about 2 days at the current rate). Write out the PromQL expressions, the alert YAML, and explain how the 14.4x multiplier is derived from the SLO target and the alerting window.
