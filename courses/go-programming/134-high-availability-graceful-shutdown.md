# Chapter 134: High Availability, Graceful Shutdown, and Zero-Downtime Deploys

Running your Go service in Kubernetes gets you container orchestration. Running it with high availability means your service survives node failures, rolling updates, and traffic spikes without dropping requests. This chapter covers the full picture: HA configuration, graceful shutdown in Go, the preStop timing race, and zero-downtime rolling updates.

## Table of Contents

1. [What High Availability Means](#1-what-high-availability-means)
2. [HA Checklist](#2-ha-checklist)
3. [PodDisruptionBudget](#3-poddisruptionbudget)
4. [Pod Anti-Affinity](#4-pod-anti-affinity)
5. [Graceful Shutdown in Go](#5-graceful-shutdown-in-go)
6. [The preStop + SIGTERM Race](#6-the-prestop--sigterm-race)
7. [Connection Draining](#7-connection-draining)
8. [Zero-Downtime Rolling Updates](#8-zero-downtime-rolling-updates)
9. [Liveness, Readiness, and Startup Probes](#9-liveness-readiness-and-startup-probes)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What High Availability Means

High availability (HA) means your service keeps running when individual components fail. The core idea: eliminate single points of failure. If one pod dies, traffic shifts to the others with no user impact.

HA is measured in nines:

| Availability | Downtime per year | Downtime per month |
|---|---|---|
| 99.9% (three nines) | 8.7 hours | 43 minutes |
| 99.99% (four nines) | 52 minutes | 4.3 minutes |
| 99.999% (five nines) | 5.25 minutes | 26 seconds |

Moving from 99.9% to 99.99% is not just about writing better code. Most of the gap comes from deployment configuration: how many replicas you run, how Kubernetes handles restarts and evictions, how your rolling update is tuned, and whether your app shuts down cleanly.

A bug in your graceful shutdown code can cost you four nines even if the rest of your infrastructure is perfect.

---

## 2. HA Checklist

Work through this checklist for every service you run in production.

**2+ replicas — never run a single pod**

A single pod means the next `kubectl rollout restart` or node maintenance window causes downtime. Two replicas minimum; three for critical services.

**Readiness probe**

Kubernetes uses the readiness probe to decide whether to send traffic to a pod. A pod that fails its readiness probe is removed from the Service endpoint list. Traffic stops flowing to it until the probe passes again. Without this, Kubernetes routes requests to pods that are still starting up or temporarily unavailable.

**Liveness probe**

If a pod is alive (the process is running) but stuck in a deadlock or infinite loop, it will never recover on its own. The liveness probe detects this and triggers a restart. Keep liveness probes simple — a deadlocked process may not be able to do complex health checks.

**PodDisruptionBudget**

When Kubernetes drains a node for maintenance or an upgrade, it evicts the pods on that node. Without a PodDisruptionBudget (PDB), Kubernetes may evict all your pods simultaneously if they happen to be on the same node. A PDB sets a floor: at least N pods must remain available during voluntary disruptions.

**Pod anti-affinity**

Anti-affinity rules tell the scheduler to spread your pods across nodes or availability zones. If all replicas land on the same node and that node fails, you have zero replicas. Anti-affinity prevents this.

**Resource requests and limits**

Without resource requests, the Kubernetes scheduler cannot make placement decisions. Without limits, a noisy neighbor pod can consume all CPU/memory on a node and cause your pod to be throttled or OOM-killed. Set both.

```yaml
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "256Mi"
```

---

## 3. PodDisruptionBudget

A PDB is a separate Kubernetes resource. It applies to a set of pods selected by label.

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: api-pdb
  namespace: production
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: api
```

`minAvailable: 2` means at least 2 pods must be running at all times during voluntary disruptions. If you have 3 replicas, only 1 can be evicted at a time.

The alternative is `maxUnavailable`:

```yaml
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: api
```

`maxUnavailable: 1` says at most 1 pod can be unavailable at once. This is equivalent if you have 3 replicas (2 must remain = 1 can be down).

PDBs apply only to voluntary disruptions: node drain, cluster upgrades, manual eviction. They do not protect against a node hardware failure (involuntary disruption).

Without a PDB, `kubectl drain node-1` can evict all three of your pods if they all happen to run on `node-1`. The drain waits for nothing.

---

## 4. Pod Anti-Affinity

Anti-affinity tells the scheduler: "prefer not to place me on a node that already has a pod with my label."

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - api
                topologyKey: kubernetes.io/hostname
      containers:
        - name: api
          image: your-registry/api:latest
```

`preferredDuringSchedulingIgnoredDuringExecution` is a soft rule: the scheduler tries to honor it but will place pods on the same node if there is no alternative (e.g., a small cluster). The `weight: 100` makes this the highest-priority preference.

`topologyKey: kubernetes.io/hostname` means "spread across hostnames (nodes)." For multi-zone clusters, use `topology.kubernetes.io/zone` to spread across availability zones.

For stricter guarantees, use `requiredDuringSchedulingIgnoredDuringExecution`:

```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
            - key: app
              operator: In
              values:
                - api
        topologyKey: topology.kubernetes.io/zone
```

This hard rule refuses to schedule the pod if no zone is available that does not already have one. Use this when zone-level isolation is a hard requirement (regulated workloads, strict SLA). The trade-off: if you have fewer zones than replicas, pods will stay pending.

---

## 5. Graceful Shutdown in Go

When Kubernetes terminates a pod it sends `SIGTERM`. Your process has a grace period (default 30 seconds) to finish in-flight requests before Kubernetes sends `SIGKILL` and force-stops everything.

The standard Go pattern:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    router := http.NewServeMux()
    router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    router.HandleFunc("/api/order", handleOrder)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Start server in a goroutine so it does not block.
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("listen error:", err)
        }
    }()

    log.Println("server started on :8080")

    // Block until we receive a signal.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    log.Println("shutting down...")

    // Give in-flight requests 30 seconds to complete.
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("forced shutdown:", err)
    }

    log.Println("server stopped cleanly")
}
```

`srv.Shutdown(ctx)` does three things:

1. Closes the listener — no new connections are accepted.
2. Closes idle connections immediately.
3. Waits for active connections to finish, up to the context deadline.

If the 30-second timeout expires before all requests finish, `Shutdown` returns `context.DeadlineExceeded` and Kubernetes sends `SIGKILL` shortly after. Set your timeout to something realistic: p99 request latency + buffer.

`http.ErrServerClosed` is the normal return value from `ListenAndServe` after `Shutdown` is called. Treat it as success, not an error.

---

## 6. The preStop + SIGTERM Race

Kubernetes does two things simultaneously when it terminates a pod:

1. Sends `SIGTERM` to the container.
2. Removes the pod from the Service endpoints (stops routing traffic to it).

The problem: load balancers (kube-proxy, Envoy, cloud load balancers) take a few seconds to observe the endpoint removal and stop sending requests to the pod. During that window, traffic still arrives at a pod that has already started shutting down.

If your Go server closes its listener the moment it receives `SIGTERM`, those in-flight requests get `connection refused`.

The fix is a `preStop` hook that sleeps before the process receives `SIGTERM`:

```yaml
spec:
  containers:
    - name: api
      image: your-registry/api:latest
      lifecycle:
        preStop:
          exec:
            command: ["/bin/sh", "-c", "sleep 5"]
```

The preStop hook runs before `SIGTERM` is delivered. The 5-second sleep gives kube-proxy and upstream load balancers time to drain the endpoint from their routing tables. After the hook completes, Kubernetes sends `SIGTERM` to the process.

Full shutdown sequence:

```
+---------------------------+----+----+-------------------------------+-------+
|  Pod running, serving     | preStop hook |   Graceful shutdown       | KILL  |
|  traffic normally         | (sleep 5s)   |   (SIGTERM + 30s window)  |       |
+---------------------------+----+----+-------------------------------+-------+
0s                          t    t+5s t+5s                            t+35s  t+35s+
                            |         |                                |
                       K8s sends  SIGTERM                        context
                       preStop    delivered                      deadline /
                       hook       to process                     SIGKILL
                            |         |
                       Load balancer  Go server stops accepting,
                       drains         waits for in-flight requests
                       endpoint
```

The total termination budget is `terminationGracePeriodSeconds` (default 30s). The preStop sleep counts against this budget. If your preStop hook takes 5 seconds and your graceful shutdown needs up to 30 seconds, set `terminationGracePeriodSeconds: 40`.

```yaml
spec:
  terminationGracePeriodSeconds: 40
  containers:
    - name: api
      lifecycle:
        preStop:
          exec:
            command: ["/bin/sh", "-c", "sleep 5"]
```

---

## 7. Connection Draining

`srv.Shutdown` handles HTTP connections. You also need to close database pools and cache clients after the HTTP server stops accepting traffic.

```go
func main() {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }

    rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})

    srv := &http.Server{Addr: ":8080", Handler: buildRouter(db, rdb)}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    log.Println("shutting down HTTP server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("HTTP shutdown error:", err)
    }

    // HTTP server is done — no more requests will touch db or rdb.
    // Now safe to close downstream clients.
    log.Println("closing database pool...")
    if err := db.Close(); err != nil {
        log.Println("db close error:", err)
    }

    log.Println("closing redis client...")
    if err := rdb.Close(); err != nil {
        log.Println("redis close error:", err)
    }

    log.Println("server stopped cleanly")
}
```

Order matters. Close the HTTP server first to stop accepting new requests, then close downstream clients. If you close the database pool before the HTTP server finishes draining, in-flight request handlers may try to use a closed pool and panic or return 500.

---

## 8. Zero-Downtime Rolling Updates

A rolling update replaces old pods with new ones incrementally. The default Kubernetes strategy allows some unavailability during the update. For zero downtime, set `maxUnavailable: 0`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      terminationGracePeriodSeconds: 40
      containers:
        - name: api
          image: your-registry/api:v2
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "sleep 5"]
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
            failureThreshold: 3
```

`maxUnavailable: 0` means no existing pod is terminated until a replacement is ready. `maxSurge: 1` allows one extra pod above the desired count during the update.

The update sequence for 3 replicas:

```
Start:   [v1] [v1] [v1]      (3 running, 0 extra)
Step 1:  [v1] [v1] [v1] [v2] (surge: 4 pods, v2 starting)
Step 2:  [v1] [v1] [v2]      (v2 passed readiness, one v1 terminated)
Step 3:  [v1] [v2] [v2]      (second v2 ready, second v1 terminated)
Step 4:  [v2] [v2] [v2]      (done)
```

The readiness probe is the critical gate. Kubernetes will not terminate an old pod until the new pod's readiness probe passes. If the new pod has a bug and never passes readiness, the rollout stalls — old pods keep serving traffic and the deployment never completes. You retain availability while the problem is diagnosed.

Without `maxUnavailable: 0`, Kubernetes might terminate an old pod and start a new one simultaneously. If the new pod takes 10 seconds to start, you have 10 seconds with one fewer replica.

---

## 9. Liveness, Readiness, and Startup Probes

Three distinct probe types serve different purposes.

**Liveness probe**

Answers: is the process alive and not stuck? If the liveness probe fails, Kubernetes restarts the container.

Use a minimal handler that just returns 200. A deadlocked process may not be able to query a database, so do not add complex checks here. The goal is to detect "process is running but completely broken."

```go
router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
})
```

**Readiness probe**

Answers: is the process ready to serve production traffic? If the readiness probe fails, the pod is removed from the Service endpoints — traffic stops routing to it, but the container is not restarted.

Use a handler that checks actual dependencies: database connectivity, cache warmup, feature flag initialization.

```go
router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        http.Error(w, "database unavailable", http.StatusServiceUnavailable)
        return
    }

    w.WriteHeader(http.StatusOK)
})
```

**Startup probe**

For containers that take a long time to initialize (running database migrations, JVM warmup), the startup probe gives extra time before the liveness probe kicks in. While the startup probe is active, liveness and readiness probes are disabled. Once the startup probe passes, normal liveness/readiness probing begins.

Full probe configuration for a realistic Go service:

```yaml
containers:
  - name: api
    image: your-registry/api:latest
    startupProbe:
      httpGet:
        path: /healthz
        port: 8080
      failureThreshold: 6      # 6 * 5s = 30s startup grace period
      periodSeconds: 5
    livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
      initialDelaySeconds: 0   # startup probe handles the initial delay
      periodSeconds: 10
      failureThreshold: 3      # restart after 3 consecutive failures (30s)
      timeoutSeconds: 2
    readinessProbe:
      httpGet:
        path: /readyz
        port: 8080
      initialDelaySeconds: 0
      periodSeconds: 5
      failureThreshold: 3      # stop routing after 3 failures (15s)
      successThreshold: 1
      timeoutSeconds: 2
```

Key distinctions:

| Probe | On failure | Checks | Endpoint |
|---|---|---|---|
| Liveness | Restart container | Process not stuck | /healthz |
| Readiness | Stop routing traffic | Dependencies reachable | /readyz |
| Startup | Delay liveness | Container finished init | /healthz |

A common mistake: using the same endpoint for liveness and readiness. If your database goes down, the readiness probe should fail (stop routing), but the liveness probe should still pass (do not restart — the restart will not fix a down database). Keep them separate.

---

## Summary

High availability in Kubernetes requires both code changes and deployment configuration:

- Run at least 2 replicas. One replica means every restart is downtime.
- Add a PodDisruptionBudget to prevent simultaneous eviction during node drain or cluster upgrades.
- Use pod anti-affinity to spread replicas across nodes and availability zones.
- Implement graceful shutdown in Go: catch `SIGTERM`/`SIGINT`, call `srv.Shutdown` with a timeout, then close downstream clients in order.
- Add a `preStop` sleep (5 seconds is typical) to give load balancers time to drain traffic before `SIGTERM` arrives.
- Set `terminationGracePeriodSeconds` to `preStop duration + graceful shutdown timeout + buffer`.
- Configure rolling updates with `maxUnavailable: 0, maxSurge: 1` for zero-downtime deploys.
- Use separate `/healthz` (liveness, always 200) and `/readyz` (readiness, checks dependencies) endpoints.
- Use a startup probe for slow-starting containers to avoid premature liveness restarts.

The preStop + graceful shutdown combination is the most commonly missed piece. Without it, you will see a small but consistent rate of 502/503 errors during every deployment.

---

## Exercises

### Easy

Add graceful shutdown to an existing HTTP server.

Start with a basic server:

```go
srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
srv.ListenAndServe()
```

Modify it to:
- Start the server in a goroutine
- Block on a channel that receives `syscall.SIGTERM` and `syscall.SIGINT`
- Call `srv.Shutdown` with a 30-second context after the signal arrives
- Log `"server stopped cleanly"` after `Shutdown` returns without error

---

### Medium

Write a Go HTTP server with separate health endpoints and configure Kubernetes probes for it.

Requirements for the Go server:
- `GET /healthz` always returns `200 OK` with body `"ok"`
- `GET /readyz` pings the database using `db.PingContext`. Returns `200` if the ping succeeds, `503` with body `"db unavailable"` if it fails or times out (2s timeout)
- Graceful shutdown on `SIGTERM`/`SIGINT`

Requirements for the Kubernetes deployment YAML:
- Liveness probe on `/healthz`, `periodSeconds: 10`, `failureThreshold: 3`
- Readiness probe on `/readyz`, `periodSeconds: 5`, `failureThreshold: 3`
- `preStop` sleep of 5 seconds
- `terminationGracePeriodSeconds: 40`

Test it locally by running the server, killing the database, and verifying that `/readyz` returns 503 while `/healthz` returns 200.

---

### Hard

Simulate the preStop + SIGTERM race in a Go test and verify that all in-flight requests complete.

Write a test that:
1. Starts an HTTP server in a goroutine. The handler sleeps 100ms (simulating a slow request) and returns `200 OK`.
2. Sends 10 concurrent requests to the server using a `sync.WaitGroup`. Each goroutine records whether it received `200` or an error.
3. Sends `SIGTERM` to the current process (use `syscall.Kill(os.Getpid(), syscall.SIGTERM)`) 50ms after the requests start (while they are still in-flight).
4. Waits for all 10 requests to complete.
5. Asserts that all 10 requests returned `200 OK` with no errors.

The test should pass only if the graceful shutdown waits for in-flight requests to finish. If you call `srv.Close()` instead of `srv.Shutdown()`, the test should fail because `Close` terminates active connections immediately.

Hint: the server needs the graceful shutdown logic from the chapter. Wire the signal handler in a goroutine before starting the requests. Use `httptest.NewServer` or bind to a random port with `:0` to avoid conflicts.
