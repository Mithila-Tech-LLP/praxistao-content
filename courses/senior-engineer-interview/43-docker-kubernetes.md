# Chapter 43: Docker, Kubernetes & Observability

Infrastructure knowledge is expected of senior engineers at top companies. You don't need to be a DevOps engineer, but you need to understand containers, orchestration, and how to observe your services in production.

## Table of Contents

1. [Docker Internals](#1-docker-internals)
2. [Container Best Practices for Go Services](#2-container-best-practices-for-go-services)
3. [Kubernetes Architecture](#3-kubernetes-architecture)
4. [Kubernetes Workloads & Patterns](#4-kubernetes-workloads--patterns)
5. [Observability — The Three Pillars](#5-observability--the-three-pillars)
6. [Structured Logging in Go](#6-structured-logging-in-go)
7. [Metrics with Prometheus](#7-metrics-with-prometheus)
8. [Distributed Tracing with OpenTelemetry](#8-distributed-tracing-with-opentelemetry)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. Docker Internals

### What Docker Actually Does

Docker containers are not VMs. They're isolated processes using Linux kernel features:

```
Linux namespaces: isolation
  pid namespace:    container has its own process tree (PID 1 is your app)
  net namespace:    container has its own network stack and IP
  mnt namespace:    container has its own filesystem view
  uts namespace:    container has its own hostname
  user namespace:   container can have its own UID/GID mapping

Linux cgroups: resource limits
  CPU: limit to 1 CPU out of 8
  Memory: limit to 512MB
  I/O: limit disk read/write rates

Union filesystem (overlayfs):
  Container filesystem = base image layers + writable top layer
  Layers are shared across containers (efficient storage)
  Changes are only written to the top layer
```

### Container Image Layers

```dockerfile
# Each instruction = new layer
FROM golang:1.22-alpine        # base layer (shared across all Go images)
WORKDIR /app
COPY go.mod go.sum ./         # layer: module files (cached if unchanged)
RUN go mod download            # layer: dependencies (cached if go.mod unchanged)
COPY . .                       # layer: source code
RUN go build -o server .       # layer: compiled binary
CMD ["./server"]               # layer: CMD instruction (not a filesystem change)

# Layer caching: Docker reuses unchanged layers
# Order matters: put frequently-changing things LAST
# COPY source code AFTER go mod download (source changes more than deps)
```

---

## 2. Container Best Practices for Go Services

```dockerfile
# Multi-stage build: compile in full image, copy binary to minimal image
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Final image: scratch (empty) or distroless — no shell, no OS tools
# Makes attack surface tiny, image is ~10MB instead of 600MB
FROM scratch
COPY --from=builder /app/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
USER 1000:1000  # never run as root
CMD ["/server"]
```

```yaml
# docker-compose.yml for local development
version: "3.8"
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/mydb
    depends_on:
      db:
        condition: service_healthy
    
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: mydb
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d mydb"]
      interval: 5s
      timeout: 5s
      retries: 5
```

---

## 3. Kubernetes Architecture

```
Control Plane:
  API Server:       REST API for all cluster operations (kubectl talks to this)
  etcd:             Distributed KV store — cluster state (all objects stored here)
  Scheduler:        Assigns pods to nodes based on resource requirements
  Controller Manager: Runs controllers (ReplicaSet, Deployment, Service controllers)
                    Each controller watches state in etcd and reconciles to desired state

Worker Nodes:
  kubelet:          Agent on each node. Receives pod specs, starts/stops containers
  kube-proxy:       Handles networking. Manages iptables rules for Service routing
  Container runtime: containerd or CRI-O. Actually runs containers.

The reconciliation loop:
  "Desired state" stored in etcd
  Controllers constantly watch actual state
  If actual ≠ desired: take action to reconcile
  This is why Kubernetes is self-healing
```

---

## 4. Kubernetes Workloads & Patterns

### Pod, Deployment, Service

```yaml
# Deployment: manages ReplicaSet which manages Pods
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-service
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1    # at most 1 pod down during update
      maxSurge: 1          # at most 1 extra pod during update
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
      - name: my-service
        image: myrepo/my-service:v1.2.3
        ports:
        - containerPort: 8080
        resources:
          requests:
            cpu: "100m"      # 0.1 CPU cores — used for scheduling
            memory: "128Mi"
          limits:
            cpu: "500m"      # 0.5 CPU cores — hard cap
            memory: "256Mi"
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 20
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
---
# Service: stable network endpoint for pods (load balancing)
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP  # internal only; use LoadBalancer for external
```

### Health Probes

```go
// In your Go service — implement health endpoints Kubernetes calls:
func healthLive(w http.ResponseWriter, r *http.Request) {
    // Liveness: am I running? If this fails: restart the pod
    // Only return failure for unrecoverable states (deadlock, corrupted state)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func healthReady(w http.ResponseWriter, r *http.Request) {
    // Readiness: am I ready to serve traffic? If this fails: remove from load balancer
    // Return failure when: database is unreachable, dependency is down, initializing
    if err := db.PingContext(r.Context()); err != nil {
        http.Error(w, "database unreachable", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
```

### Horizontal Pod Autoscaler (HPA)

```yaml
# Auto-scale based on CPU usage
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: my-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-service
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70  # scale up when avg CPU > 70%
```

---

## 5. Observability — The Three Pillars

```
Logs:   What happened? (events, errors, debug info)
Metrics: How is the system behaving? (latency, error rate, throughput, saturation)
Traces: Why did this request take so long? (end-to-end request path)

The three work together:
  Alert fires on metric: "p99 latency > 5 seconds"
  Engineer investigates: looks at logs for errors in that timeframe
  For complex root cause: uses traces to find which service/database is slow
```

---

## 6. Structured Logging in Go

```go
import "log/slog"

// slog: Go 1.21+ structured logging (built-in)
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// Structured log — machine-parseable JSON:
logger.Info("user logged in",
    "user_id", userID,
    "ip", r.RemoteAddr,
    "user_agent", r.UserAgent(),
)
// Output: {"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"user logged in","user_id":"123","ip":"192.168.1.1"}

// Add context to log (request ID, trace ID):
func withRequestLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        reqID := uuid.New().String()
        ctx := context.WithValue(r.Context(), loggerKey, 
            logger.With("request_id", reqID))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func logFromContext(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return l
    }
    return slog.Default()
}

// Log errors with full context:
log := logFromContext(ctx)
if err := db.GetUser(ctx, userID); err != nil {
    log.Error("failed to get user",
        "user_id", userID,
        "error", err,
        "stack", debug.Stack())
    return
}
```

---

## 7. Metrics with Prometheus

```go
import "github.com/prometheus/client_golang/prometheus"
import "github.com/prometheus/client_golang/prometheus/promauto"

// Metric types:
// Counter: monotonically increasing (requests, errors, bytes)
var httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "http_requests_total",
    Help: "Total HTTP requests by method and status",
}, []string{"method", "path", "status"})

// Gauge: current value (connections, goroutines, queue size)
var activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
    Name: "active_connections",
    Help: "Number of active WebSocket connections",
})

// Histogram: distribution of values (request durations)
var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
    Name:    "http_request_duration_seconds",
    Help:    "HTTP request duration in seconds",
    Buckets: prometheus.DefBuckets, // [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10]
}, []string{"method", "path"})

// Instrument your handler:
func instrumentedHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(wrapped.statusCode)
        
        httpRequests.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}

// Expose metrics endpoint for Prometheus to scrape:
http.Handle("/metrics", promhttp.Handler())
```

### Key Metrics to Track (RED Method)

```
Rate:    Requests per second
Errors:  Error rate (4xx/5xx percentage)
Duration: Request latency (p50, p95, p99)

For resources (USE Method):
Utilization: CPU%, memory%, disk%
Saturation: queue length, wait time
Errors: failed operations
```

---

## 8. Distributed Tracing with OpenTelemetry

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

// Initialize tracing (send to Jaeger/Tempo):
func initTracing(serviceName string) (func(), error) {
    exporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    
    tp := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exporter),
        tracesdk.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(serviceName),
        )),
    )
    otel.SetTracerProvider(tp)
    
    return func() { tp.Shutdown(context.Background()) }, nil
}

// Instrument your code:
tracer := otel.Tracer("my-service")

func getUser(ctx context.Context, userID string) (*User, error) {
    ctx, span := tracer.Start(ctx, "getUser",
        trace.WithAttributes(attribute.String("user.id", userID)))
    defer span.End()
    
    user, err := db.GetUser(ctx, userID)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        span.RecordError(err)
        return nil, err
    }
    
    span.SetAttributes(attribute.String("user.email", user.Email))
    return user, nil
}

// HTTP middleware to propagate trace context:
// (otelhttp package handles this automatically)
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

http.Handle("/api/users", otelhttp.NewHandler(usersHandler, "list-users"))
```

---

## 9. Interview Questions & Model Answers

**Q: What is the difference between liveness and readiness probes in Kubernetes?**

"A liveness probe checks if the container is running correctly — if it fails, Kubernetes restarts the pod. Use it to detect deadlocks or corrupted state where the process is running but not making progress. A readiness probe checks if the container is ready to serve traffic — if it fails, Kubernetes removes the pod from the Service's endpoints but doesn't restart it. Use readiness to signal when the service is warming up (loading caches, establishing DB connections) or when a dependency is temporarily unavailable. The critical distinction: liveness failures cause restarts (which can cause cascading failures under load), so only fail liveness for truly unrecoverable states."

**Q: What are the three pillars of observability and why do you need all three?**

"Logs, metrics, and traces work together. Metrics give you the 30,000-foot view: alert when p99 latency spikes. But metrics don't tell you WHY — they just say something is wrong. Logs give you events in context, but searching through GB of logs is slow. Traces connect the dots: they show the end-to-end path of a specific request, including which service and which database query was slow. The workflow is: metric alert fires → look at logs for errors in that window → use trace ID from logs to find the full trace → see exactly which downstream call caused the slowdown. You need all three because no single pillar answers all questions."

---

## Summary

- **Docker:** containers = isolated processes via Linux namespaces + cgroups. Layers are cached.
- **Multi-stage builds:** compile in Go image, copy binary to scratch/distroless. Tiny, secure images.
- **Kubernetes:** control plane (API Server, etcd, Scheduler) + worker nodes (kubelet, kube-proxy).
- **Reconciliation loop:** controllers constantly compare desired vs actual state.
- **Readiness:** remove from load balancer when not ready. **Liveness:** restart when unrecoverable.
- **HPA:** auto-scale pods based on CPU/memory/custom metrics.
- **Logs:** structured JSON (slog in Go). Include request_id and trace_id for correlation.
- **Metrics:** RED (Rate, Errors, Duration) for services. USE (Utilization, Saturation, Errors) for resources.
- **Traces:** end-to-end request path. Propagate trace context via HTTP headers. OpenTelemetry is the standard.
