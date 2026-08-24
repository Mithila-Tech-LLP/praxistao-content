# Chapter 91: Monitoring with Prometheus and Grafana

Your GoChain testnet is now a real, deployed thing — a Docker Compose stack of nodes running on a cloud VM, reachable on ports 8080 (API/HTTP) and 9000 (P2P), with a faucet handing out test coins. But right now, the only way to know what it is doing is to `ssh` in and read logs one line at a time. This chapter gives GoChain eyes: real Prometheus metrics wired into the actual mining, mempool, and networking code you have already written, and a Grafana dashboard that turns those numbers into a live picture of the network's health.

## Table of Contents

1. [Observability Is a Habit, Not a Checkbox](#1-observability-is-a-habit-not-a-checkbox)
2. [Metrics, Logs, and Traces](#2-metrics-logs-and-traces)
3. [How Prometheus Works: the Pull Model](#3-how-prometheus-works-the-pull-model)
4. [The Metric Types GoChain Needs](#4-the-metric-types-gochain-needs)
5. [Adding client_golang and a metrics Package](#5-adding-client_golang-and-a-metrics-package)
6. [Instrumenting Block Height and Mempool Size](#6-instrumenting-block-height-and-mempool-size)
7. [Instrumenting Peer Count](#7-instrumenting-peer-count)
8. [Instrumenting Mining Hashrate](#8-instrumenting-mining-hashrate)
9. [Instrumenting API Request Duration](#9-instrumenting-api-request-duration)
10. [Exposing /metrics and Scraping It with Prometheus](#10-exposing-metrics-and-scraping-it-with-prometheus)
11. [Building the Grafana Dashboard](#11-building-the-grafana-dashboard)
12. [A Word on Alerting](#12-a-word-on-alerting)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. Observability Is a Habit, Not a Checkbox

Every earlier volume ended when the code worked: tests passed, a demo ran, you moved on. Running a live testnet is different. GoChain is now a process that is supposed to keep running, unattended, for days or weeks, serving other people's wallets and other people's nodes. When something goes wrong at 3 a.m. — a node stops mining, peers silently drop off, the API starts timing out — nobody is sitting at a terminal watching it happen.

**Observability** is the property of a system that lets you answer "what is it doing right now, and why" from the outside, without attaching a debugger or guessing from scattered log lines. It is not something you add once and forget; it is an ongoing habit, the same way writing tests became a habit back in Volume 1. Every new feature you add to GoChain from this point on should come with the question "how will I know if this breaks in production?" attached to it, the same reflex as "does this have a test?"

Think of it like the dashboard in a car. You do not open the hood and inspect the engine every time you want to know your speed, fuel level, or engine temperature — a small panel in front of you surfaces exactly the numbers that matter, continuously, so a problem is visible the moment it starts, not after the car stalls on the highway. This chapter builds GoChain's dashboard.

---

## 2. Metrics, Logs, and Traces

Observability tooling generally splits into three categories, and it helps to know where "metrics" sits before diving into Prometheus specifically.

- **Logs** are discrete, timestamped text events — "block 4021 mined," "peer 93.184.1.2 disconnected." You already have these from every earlier volume's `log.Printf` calls. Logs are great for *what exactly happened*, but terrible for *trends over time* — nobody wants to eyeball ten thousand log lines to see whether mempool size has been creeping up all week.
- **Metrics** are numeric measurements, sampled at regular intervals and stored as a time series — "mempool size was 12 at 10:00:00, 15 at 10:00:15, 9 at 10:00:30." Metrics are cheap to store, cheap to query, and exactly what you want for dashboards, trends, and alerting. This chapter is entirely about metrics.
- **Traces** follow a single request as it moves through multiple services (a request hits the API, which calls the mempool, which calls storage), timing each hop. GoChain's architecture is simple enough that traces are not part of this course, but the concept is worth knowing — it is the natural next step once a system has many independently deployed services instead of one binary.

For a single testnet of GoChain nodes, metrics give the best signal for the least effort, which is why they come first.

---

## 3. How Prometheus Works: the Pull Model

**Prometheus** is an open-source monitoring system that stores metrics as time series and lets you query them. Unlike a logging pipeline where your application *pushes* data somewhere, Prometheus uses a **pull model**: your application exposes its current metric values on a plain HTTP endpoint (`/metrics`), and Prometheus itself, running as a separate process, periodically visits that endpoint and "scrapes" (reads) the current values, storing each one with a timestamp.

```
   +------------------+          scrape every 15s          +----------------+
   |  gochain node    |  <-------------------------------- |   Prometheus   |
   |  :8080/metrics   |  ------------------------------->   |    server      |
   +------------------+     plain-text metric snapshot     +----------------+
                                                                     |
                                                                     | PromQL queries
                                                                     v
                                                             +----------------+
                                                             |    Grafana     |
                                                             | (dashboards)   |
                                                             +----------------+
```

This diagram is the entire chapter in miniature: GoChain exposes numbers on an HTTP endpoint, Prometheus repeatedly fetches and stores them, and Grafana queries Prometheus's stored history to draw graphs. Your job in this chapter is exactly the left-hand box — making `/metrics` exist and contain the right numbers — plus configuring the two boxes on the right to point at it.

The pull model has a practical advantage for a testnet: Prometheus does not need to know anything about *when* a node changes state, only *where* to find it. As long as a GoChain node's `/metrics` endpoint is reachable, Prometheus will keep sampling it, even through restarts, without any code in GoChain needing to know Prometheus exists beyond exposing that one endpoint.

---

## 4. The Metric Types GoChain Needs

Prometheus's client libraries support a handful of metric types; GoChain needs exactly two of them.

A **gauge** is a value that can go up or down freely, like a speedometer — it always reports "the current value right now." Block height, mempool size, peer count, and mining hashrate are all gauges: each one is a single current number that changes over time in either direction.

A **histogram** buckets observed values into ranges and counts how many observations fell into each range, letting you later ask questions like "what fraction of API requests finished in under 100ms?" instead of only ever seeing an average. API request duration is a histogram, because a single "average response time" hides important detail — a few very slow requests can be invisible in an average but obvious in a histogram's tail buckets.

```
 gauge (block height over time)          histogram (request durations, bucketed)
 -----------------------------           ---------------------------------------
   42 |        __/‾‾\                     count
      |       /      \___                  |
   40 |   ___/            \                 |   ██
      |  /                                  |   ██  ██
   38 |_/                                   |   ██  ██  ██  ▁▁
      +------------------> time             +---------------------> bucket (seconds)
                                                 <0.01 <0.05 <0.1  <0.5
```

Here are the exact five metrics this chapter instruments, matching the names GoChain will use everywhere from here on — dashboards, alerts, and any future volume that reads them:

| Metric name | Type | What it measures |
|---|---|---|
| `gochain_block_height` | gauge | Current height of the local chain |
| `gochain_mempool_size` | gauge | Number of pending transactions waiting to be mined |
| `gochain_peer_count` | gauge | Number of currently connected P2P peers |
| `gochain_mining_hashrate` | gauge | Estimated hashes per second from the local miner |
| `gochain_api_request_duration_seconds` | histogram, labeled by `endpoint` | How long each API request took, bucketed |

---

## 5. Adding client_golang and a metrics Package

Prometheus's official Go client library, `prometheus/client_golang`, does the heavy lifting of storing metric values thread-safely and rendering them in the exact text format Prometheus expects. Add it to GoChain's module:

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

Rather than scattering metric definitions across every package that needs one, we give GoChain a single new package, `gochain/metrics`, that every other package imports. This mirrors the same reasoning as Volume 8's `storage.Store` interface: one shared place, one source of truth, no risk of two packages accidentally defining `gochain_block_height` twice with slightly different help text.

```go
// gochain/metrics/metrics.go
//
// Package metrics defines every Prometheus metric GoChain exposes. Every
// other package (core, network, consensus, api) imports this package and
// calls methods on these shared variables — it never creates its own
// competing metric with the same name.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// BlockHeight reports the height of the highest block the local node
	// has accepted. A gauge because it only ever needs to reflect "right
	// now," not a running total.
	BlockHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gochain_block_height",
		Help: "Current height of the local blockchain",
	})

	// MempoolSize reports how many transactions are currently waiting,
	// unmined, in this node's mempool (Volume 5, Chapter 34).
	MempoolSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gochain_mempool_size",
		Help: "Number of pending transactions in the mempool",
	})

	// PeerCount reports how many P2P connections (Volume 7) this node
	// currently has open, in either direction.
	PeerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gochain_peer_count",
		Help: "Number of connected P2P peers",
	})

	// MiningHashrate reports the miner's estimated hashes-per-second,
	// recomputed periodically rather than on every single hash attempt
	// (updating a metric a million times a second would itself become a
	// performance problem).
	MiningHashrate = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gochain_mining_hashrate",
		Help: "Estimated hashes per second from the local miner",
	})

	// APIRequestDuration is a histogram, labeled by endpoint, so
	// "/balance" and "/sendTransaction" get separate buckets instead of
	// being averaged together into one meaningless number.
	APIRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gochain_api_request_duration_seconds",
		Help:    "Duration of API requests in seconds, labeled by endpoint",
		Buckets: prometheus.DefBuckets, // 5ms .. 10s, Prometheus's sensible default spread
	}, []string{"endpoint"})
)
```

`promauto.NewGauge` and `promauto.NewHistogramVec` both do two things at once: they construct the metric object, and they automatically register it with Prometheus's default global registry, so we never have to remember a separate `prometheus.MustRegister(...)` call. `prometheus.DefBuckets` is a built-in set of histogram bucket boundaries (5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s) that comfortably covers everything from a fast in-memory lookup to a slow disk-backed query — fine defaults for GoChain's API, which does nothing exotic enough to need custom buckets.

---

## 6. Instrumenting Block Height and Mempool Size

Metrics are only useful if they are updated at exactly the moment the underlying value actually changes — not on a timer that might lag reality, but right inside the code path that already knows the new value. `core.Blockchain.AddBlock` (Chapter 18) and `MineBlock` (Chapter 25) are the only two places a new block ever becomes part of the chain, so that is where `BlockHeight` gets set:

```go
// gochain/core/blockchain.go

import "github.com/you/gochain/metrics"

// AddBlock appends a validated block to the chain and persists it. This is
// called both when we mine a block ourselves and when we accept one
// synced from a peer (Chapter 49) — either way, height changes here.
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.ValidateBlock(b); err != nil {
		return err
	}

	bc.blocks = append(bc.blocks, b)
	if err := bc.store.PutBlock(b); err != nil {
		return err
	}

	// The metric is updated in the exact same place the height itself
	// changes, so it can never silently drift out of sync with reality.
	metrics.BlockHeight.Set(float64(b.Height))

	return nil
}
```

The mempool follows the same pattern. `core.Mempool.Add` (Chapter 34) inserts a pending transaction, and `Remove` (called both when a transaction is mined into a block and when it is evicted) takes one out — both are the right place to refresh `MempoolSize`:

```go
// gochain/core/mempool.go

import "github.com/you/gochain/metrics"

func (mp *Mempool) Add(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if err := mp.validate(tx); err != nil {
		return err
	}
	mp.txs[tx.IDHex()] = tx

	// len(mp.txs) is already the source of truth for mempool size —
	// we are just also reporting it to Prometheus.
	metrics.MempoolSize.Set(float64(len(mp.txs)))
	return nil
}

func (mp *Mempool) Remove(txID string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	delete(mp.txs, txID)
	metrics.MempoolSize.Set(float64(len(mp.txs)))
}
```

Notice the shape repeating: find the one function (or small handful of functions) that is the single source of truth for a value, and set the gauge right there, under the same lock that protects the underlying data. Never update a metric from a separate polling goroutine when the real code path is already known and already holding the right lock — that only invites the metric to lag or race against reality.

---

## 7. Instrumenting Peer Count

`network.Node` (Chapter 46) tracks connected peers in a map, adding an entry when a handshake (Chapter 47) completes and removing one when a connection drops. Both of those events are exactly where `PeerCount` belongs:

```go
// gochain/network/node.go

import "github.com/you/gochain/metrics"

func (n *Node) addPeer(p *Peer) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.peers[p.Address] = p
	metrics.PeerCount.Set(float64(len(n.peers)))
}

func (n *Node) removePeer(addr string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.peers, addr)
	metrics.PeerCount.Set(float64(len(n.peers)))
}
```

This single number, watched over time on a live testnet, is one of the most useful early-warning signals you have. A healthy node's peer count should hover around a stable, non-zero value; a node whose peer count drops to zero and stays there is isolated from the network — either an eclipse attack (Chapter 51) in progress, a firewall misconfiguration, or a crashed peer — long before anyone would notice by reading logs.

---

## 8. Instrumenting Mining Hashrate

Hashrate is different from the previous three metrics: there is no single function call that "is" the hashrate the way `len(mp.txs)` "is" the mempool size. Instead, the concurrent miner from Chapter 27 counts how many nonces it has tried, and we periodically convert that running count into a rate.

```go
// gochain/consensus/pow.go

import (
	"time"

	"github.com/you/gochain/metrics"
)

// Run searches for a valid nonce across multiple goroutines (Chapter 27).
// Each worker increments a shared atomic counter every time it tries a
// nonce; a separate reporter goroutine turns that counter into a
// hashes-per-second estimate once a second, which is a fine enough
// resolution for a dashboard without adding measurable overhead to the
// hot hashing loop itself.
func (pow *ProofOfWork) Run(ctx context.Context) (*Block, error) {
	var attempts int64 // updated with atomic.AddInt64 by every worker

	reportCtx, cancelReport := context.WithCancel(ctx)
	defer cancelReport()

	go pow.reportHashrate(reportCtx, &attempts)

	// ... existing goroutine fan-out from Chapter 27, each worker calling
	// atomic.AddInt64(&attempts, 1) once per nonce it tries ...

	return pow.waitForSolution(ctx, &attempts)
}

func (pow *ProofOfWork) reportHashrate(ctx context.Context, attempts *int64) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var last int64
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := atomic.LoadInt64(attempts)
			// Hashes tried in the last second == this second's rate,
			// since the ticker fires exactly once a second.
			metrics.MiningHashrate.Set(float64(current - last))
			last = current
		}
	}
}
```

This is a legitimate exception to "update the metric where the value changes": hashrate is fundamentally a *rate*, not an instantaneous fact, so it has to be computed by sampling a counter over a fixed time window rather than set directly. The `reportHashrate` goroutine is canceled the moment mining stops (via `reportCtx`), so a node that stops mining correctly shows its hashrate dropping to whatever the last reported second's rate was, rather than reporting a stale high number forever.

---

## 9. Instrumenting API Request Duration

The API server from Chapter 70 handles requests through Go's standard `net/http`. Rather than adding timing code inside every single handler, we wrap the whole router in one **middleware** — a function that wraps an `http.Handler` and runs code before and after the wrapped handler, exactly the same pattern you would use to add logging or authentication.

```go
// gochain/api/middleware.go

import (
	"net/http"
	"time"

	"github.com/you/gochain/metrics"
)

// MetricsMiddleware wraps every API route so request duration is measured
// identically everywhere, without every handler needing to remember to
// do it itself. endpoint is a stable label like "/balance" — never the
// raw URL with query parameters, which would create a new, unbounded
// time series for every distinct request and slowly overwhelm Prometheus.
func MetricsMiddleware(endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start).Seconds()

		metrics.APIRequestDuration.
			WithLabelValues(endpoint).
			Observe(duration)
	}
}
```

Wiring it in means wrapping each route registration exactly once, at startup:

```go
// gochain/api/server.go

mux := http.NewServeMux()
mux.HandleFunc("/balance", MetricsMiddleware("/balance", handleBalance))
mux.HandleFunc("/sendTransaction", MetricsMiddleware("/sendTransaction", handleSendTransaction))
mux.HandleFunc("/getBlock", MetricsMiddleware("/getBlock", handleGetBlock))
```

`WithLabelValues(endpoint)` looks up (or creates, the first time) the specific histogram for that one endpoint label — this is exactly why `APIRequestDuration` was declared as a `HistogramVec` rather than a plain `Histogram` back in Section 5: a `Vec` (vector) is a family of metrics, one per unique combination of label values, all sharing one name. `/balance` and `/sendTransaction` end up as two separate time series under the same metric name, distinguishable in Grafana by their `endpoint` label.

---

## 10. Exposing /metrics and Scraping It with Prometheus

With every metric now updated at the right place, the last piece of code is the HTTP endpoint Prometheus actually scrapes. Because GoChain's API server already listens on port 8080 (the port opened in Chapter 88's firewall rules), `/metrics` is simply one more route on that same server — no new port to open, no new firewall rule needed:

```go
// gochain/api/server.go

import "github.com/prometheus/client_golang/prometheus/promhttp"

func StartServer(addr string) error {
	mux := http.NewServeMux()

	// ... existing /balance, /sendTransaction, /getBlock routes ...

	// promhttp.Handler() renders every registered metric (Section 5) in
	// the plain-text exposition format Prometheus expects — one line per
	// metric sample, human-readable enough to curl directly.
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, mux)
}
```

You can see this working immediately, without Prometheus involved at all, with `curl`:

```bash
curl http://localhost:8080/metrics | grep gochain_

# gochain_block_height 142
# gochain_mempool_size 3
# gochain_peer_count 4
# gochain_mining_hashrate 18452
# gochain_api_request_duration_seconds_bucket{endpoint="/balance",le="0.005"} 12
# gochain_api_request_duration_seconds_bucket{endpoint="/balance",le="0.01"} 15
# ...
```

Now point an actual Prometheus server at it. Add a `prometheus` service to the `docker-compose.yml` from Chapter 87, alongside your existing `gochain` node services and `gochain-explorer`:

```yaml
# docker-compose.yml — additions for Chapter 91
services:
  # ... existing gochain-node-1, gochain-node-2, gochain-explorer, faucet ...

  prometheus:
    image: prom/prometheus:v2.53.0
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
    ports:
      - "9090:9090" # Prometheus's own web UI, for local debugging only
    depends_on:
      - gochain-node-1
      - gochain-node-2

  grafana:
    image: grafana/grafana:11.1.0
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
    depends_on:
      - prometheus

volumes:
  grafana-data:
```

`prometheus.yml` is Prometheus's own configuration file, telling it *where* to scrape and *how often*:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s # how often Prometheus polls every target below

scrape_configs:
  - job_name: "gochain-nodes"
    static_configs:
      # Compose's internal DNS resolves service names to container IPs,
      # so we reference nodes by name, not by IP address — the same
      # pattern Chapter 87 already relies on for node-to-node discovery.
      - targets:
          - "gochain-node-1:8080"
          - "gochain-node-2:8080"
```

`scrape_interval: 15s` means every fifteen seconds, Prometheus fetches `/metrics` from each target in the list and stores whatever it finds, timestamped with the moment it fetched it. Ports 9090 and 3000 are only exposed for local debugging over an SSH tunnel — Chapter 93 puts Grafana and the explorer behind a real domain and HTTPS instead of leaving raw ports open to the internet.

---

## 11. Building the Grafana Dashboard

**Grafana** is a dashboard tool that queries a data source — here, Prometheus — using **PromQL**, Prometheus's own query language, and renders the results as graphs, single-number panels, and tables. Rather than reading raw numbers off `/metrics`, Grafana is what turns "gochain_block_height 142" into a live, auto-updating line chart.

After starting the `grafana` service, visit `http://<your-vm-ip>:3000` (default login `admin`/`admin`, which you should change immediately), add Prometheus as a data source pointing at `http://prometheus:9090` (again, Compose's internal service name), and build a dashboard with one panel per metric:

```
 +--------------------+  +--------------------+  +--------------------+
 |   Block Height     |  |   Mempool Size      |  |    Peer Count       |
 |   (time series)     |  |   (time series)      |  |   (time series)     |
 |     ___/‾‾\___      |  |    _/\_    _/\_      |  |   ‾‾‾\___/‾‾‾       |
 +--------------------+  +--------------------+  +--------------------+

 +--------------------+  +----------------------------------------------+
 |  Mining Hashrate    |  |      API p95 Request Duration by Endpoint    |
 |   (time series)     |  |      (time series, one line per endpoint)    |
 |    ‾‾\__/‾‾\__      |  |    /balance  ----                            |
 +--------------------+  |    /sendTransaction  ----                    |
                         +----------------------------------------------+
```

Each panel's PromQL query is short:

- **Block Height** — `gochain_block_height` (a raw gauge needs no aggregation; just plot it).
- **Mempool Size** — `gochain_mempool_size`.
- **Peer Count** — `gochain_peer_count`.
- **Mining Hashrate** — `gochain_mining_hashrate`.
- **API p95 Duration** — `histogram_quantile(0.95, sum(rate(gochain_api_request_duration_seconds_bucket[5m])) by (le, endpoint))`. This computes the 95th-percentile request duration over a rolling 5-minute window: `rate(...[5m])` turns the histogram's raw cumulative bucket counts into a per-second rate, and `histogram_quantile` estimates the duration below which 95% of requests fall, per `endpoint`.

A dashboard is just JSON underneath; Grafana can import one directly. A trimmed example of one panel's definition, saved as `gochain-dashboard.json` and imported via Grafana's "Import Dashboard" screen:

```json
{
  "title": "GoChain Testnet",
  "panels": [
    {
      "title": "Block Height",
      "type": "timeseries",
      "targets": [
        { "expr": "gochain_block_height", "legendFormat": "{{instance}}" }
      ]
    },
    {
      "title": "API p95 Duration by Endpoint",
      "type": "timeseries",
      "targets": [
        {
          "expr": "histogram_quantile(0.95, sum(rate(gochain_api_request_duration_seconds_bucket[5m])) by (le, endpoint))",
          "legendFormat": "{{endpoint}}"
        }
      ]
    }
  ]
}
```

`legendFormat: "{{instance}}"` and `"{{endpoint}}"` tell Grafana to label each line on the graph using that metric's label values — so with two nodes scraped, the Block Height panel automatically draws two separate lines, one per node, without any extra configuration.

---

## 12. A Word on Alerting

A dashboard nobody is looking at is not much better than logs nobody is reading. Both Prometheus and Grafana support **alerting rules** — conditions that, when true for a sustained period, fire a notification (email, Slack, a webhook) instead of waiting for a human to notice a graph looking wrong. A reasonable first alert for GoChain: fire if `gochain_peer_count == 0` for more than five minutes, since an isolated node is almost always a real problem worth waking someone up for. This course does not build a full alerting pipeline, but the dashboard from this chapter is the exact foundation any alerting rule would query against — the same PromQL expressions, just evaluated continuously instead of viewed on demand.

---

## Summary

- Observability is an ongoing engineering habit: every new feature should ship with an answer to "how would I know if this broke in production?"
- Metrics, logs, and traces serve different purposes; GoChain's testnet needs metrics first, for trends and dashboards.
- Prometheus uses a **pull model** — your app exposes `/metrics`, Prometheus scrapes it on a schedule, Grafana queries Prometheus's stored history.
- GoChain exposes five metrics: `gochain_block_height`, `gochain_mempool_size`, `gochain_peer_count`, and `gochain_mining_hashrate` as gauges, plus `gochain_api_request_duration_seconds` as a histogram labeled by endpoint.
- Each metric is updated at the single source-of-truth code path that already knows the new value — `AddBlock`, `Mempool.Add`/`Remove`, peer connect/disconnect, and an HTTP middleware — never from a disconnected polling loop.
- Hashrate is a rate, not an instantaneous fact, so it is computed by sampling an atomic counter once a second rather than set directly.
- `/metrics` rides on the same port 8080 the API already uses, needing no new firewall rule.
- Grafana turns Prometheus's stored numbers into live dashboards using PromQL, including percentile queries like `histogram_quantile` that a plain average would hide.

---

## Exercises

### Easy

1. Run `curl http://localhost:8080/metrics | grep gochain_` against a local GoChain node with this chapter's instrumentation added, and paste the five metric lines it prints. Explain, in your own words, what each one currently reports.

2. Add a Grafana panel for `gochain_mempool_size` using the `timeseries` visualization, and submit a handful of test transactions to your local testnet while watching it update live. Screenshot or describe what you see.

3. Change `scrape_interval` in `prometheus.yml` from `15s` to `5s`, restart the `prometheus` service, and describe what trade-off you are making by scraping more frequently (consider both node load and dashboard responsiveness).

### Medium

4. Add a sixth metric, `gochain_utxo_set_size` (a gauge), updated wherever Volume 8's UTXO index (Chapter 56) adds or removes an entry. Wire it into the Grafana dashboard next to Mempool Size.

5. Write a PromQL query that computes the *rate of new blocks per minute* over the last 15 minutes, using `gochain_block_height` (hint: `rate()` works on counters and gauges that only increase; think about what function computes "how much did this gauge change over this window" instead).

6. The `MetricsMiddleware` in Section 9 takes `endpoint` as a hardcoded string argument at each call site. Refactor it to derive the endpoint label automatically from the registered route pattern instead, and explain why hardcoding a label from the *raw request URL* (rather than the route pattern) would have been a mistake.

### Hard

7. Simulate a node becoming isolated: use `docker network disconnect` to cut a running `gochain-node` container off from the Compose network, and watch `gochain_peer_count` fall to zero in Grafana. Measure and report how long it takes from disconnection to the metric reflecting zero, and explain what determines that delay.

8. Design (on paper, or as an actual Prometheus alerting rule in `prometheus.yml`) an alert that fires if `gochain_mining_hashrate` drops by more than 50% compared to its value 10 minutes ago, without firing on the normal, expected variance of a healthy miner. Explain your threshold choice.

9. Investigate what happens to `gochain_api_request_duration_seconds` if a single malformed or attacker-crafted request causes a handler to hang indefinitely (never calling `next(w, r)`'s inner logic to completion). Does the current middleware design in Section 9 protect the histogram from being skewed by such a request? Propose a fix (for example, a request timeout) and explain where it should be added.
