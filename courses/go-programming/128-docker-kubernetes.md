# Chapter 128: Docker and Kubernetes

Docker packages your Go app into a portable container. Kubernetes (K8s) runs, scales, and heals containers in production. Together they replace the old "SSH into a server and run the binary" deployment model.

## Table of Contents

1. [Dockerfile for Go](#1-dockerfile-for-go)
2. [Docker Compose for Local Dev](#2-docker-compose-for-local-dev)
3. [Kubernetes Fundamentals](#3-kubernetes-fundamentals)
4. [Deployments and Services](#4-deployments-and-services)
5. [ConfigMaps, Secrets, Probes](#5-configmaps-secrets-and-probes)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Dockerfile for Go

```dockerfile
# Multi-stage build: build stage + final minimal image
# Stage 1: build
FROM golang:1.23-alpine AS builder

# CA certificates and timezone data to copy into the scratch image later
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy go.mod and go.sum first — Docker caches this layer if unchanged
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (static, no CGO)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/api

# Stage 2: final minimal image
FROM scratch

# Copy certificates for HTTPS outbound calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/server /server

# Run as non-root user (security best practice)
USER 65534:65534

EXPOSE 8080

# Use exec form — no shell, signals delivered directly to the process
ENTRYPOINT ["/server"]
```

Build and run:
```bash
docker build -t myapp:latest .
docker run -p 8080:8080 -e DATABASE_URL=postgres://... myapp:latest

# Check image size — multi-stage builds are typically 5-20 MB
docker image ls myapp
```

### .dockerignore

```
.git
.github
*.md
*_test.go
.env
*.local
vendor/
dist/
```

---

## 2. Docker Compose for Local Dev

```yaml
# docker-compose.yml
version: '3.9'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://myapp:secret@db:5432/myapp?sslmode=disable
      REDIS_URL: redis://redis:6379
      ENV: development
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./config:/app/config:ro

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: secret
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myapp"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru

  migrate:
    image: migrate/migrate
    command: ["-path=/migrations", "-database=postgres://myapp:secret@db:5432/myapp?sslmode=disable", "up"]
    volumes:
      - ./migrations:/migrations
    depends_on:
      db:
        condition: service_healthy

volumes:
  pgdata:
```

```bash
docker compose up -d      # start everything
docker compose logs -f api # follow API logs
docker compose down -v     # stop and remove volumes
```

---

## 3. Kubernetes Fundamentals

```
Cluster: a set of nodes running K8s

Node: a machine (VM or physical) running Pods

Pod: the smallest deployable unit — one or more containers that share:
  - Network namespace (same IP, port space)
  - Storage volumes
  - Lifecycle

Deployment: manages a set of identical Pods (desired state + rolling updates)

Service: stable network endpoint for a set of Pods (load balanced)

ConfigMap: non-secret configuration (env vars, config files)

Secret: sensitive configuration (passwords, tokens) — base64 encoded

Namespace: virtual cluster for isolation (dev, staging, prod)
```

### K8s object flow

```
You → kubectl apply -f deployment.yaml
        → API Server stores desired state in etcd
          → Scheduler assigns Pods to Nodes
            → Kubelet on each Node starts the containers
              → Service Controller sets up load balancing
```

---

## 4. Deployments and Services

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: production
  labels:
    app: api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1         # allow 1 extra pod during update
      maxUnavailable: 0   # never go below desired replicas
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
      - name: api
        image: myapp:1.2.3  # always use a specific tag, never :latest
        ports:
        - containerPort: 8080
        
        env:
        - name: ENV
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: database-url
        - name: PORT
          value: "8080"
        
        resources:
          requests:
            cpu: 100m      # 0.1 CPU cores
            memory: 128Mi
          limits:
            cpu: 500m      # 0.5 CPU cores
            memory: 512Mi
        
        # Health checks
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          failureThreshold: 3
        
        # Graceful shutdown
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sh", "-c", "sleep 5"]
      
      terminationGracePeriodSeconds: 30
```

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: api
  namespace: production
spec:
  selector:
    app: api           # routes to pods with this label
  ports:
  - protocol: TCP
    port: 80           # external port
    targetPort: 8080   # pod port
  type: ClusterIP      # internal only; use LoadBalancer for external
```

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api
  annotations:
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api
            port:
              number: 80
```

---

## 5. ConfigMaps, Secrets, and Probes

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
data:
  LOG_LEVEL: "info"
  RATE_LIMIT_RPS: "100"
  CACHE_TTL: "5m"

# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: api-secrets
type: Opaque
# Values must be base64-encoded: echo -n "value" | base64
data:
  database-url: cG9zdGdyZXM6Ly8...
  jwt-secret: c2VjcmV0a2V5MTIz...
# In practice: use ExternalSecrets + AWS Secrets Manager or Vault
```

### Health check endpoints in Go

```go
// K8s calls /healthz to check if the app is alive (restart if it fails)
// K8s calls /readyz to check if the app is ready to serve traffic

type HealthChecker struct {
    db    *sqlx.DB
    redis *redis.Client
}

func (h *HealthChecker) Healthz(w http.ResponseWriter, r *http.Request) {
    // Liveness: is the process healthy? (not deadlocked, no fatal errors)
    // Keep this fast and simple — failing it causes a restart
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}

func (h *HealthChecker) Readyz(w http.ResponseWriter, r *http.Request) {
    // Readiness: can this pod serve traffic? (deps connected, warmed up)
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()
    
    type check struct{ name string; err error }
    checks := []check{
        {"database", h.db.PingContext(ctx)},
        {"redis",    h.redis.Ping(ctx).Err()},
    }
    
    var failures []string
    for _, c := range checks {
        if c.err != nil {
            failures = append(failures, fmt.Sprintf("%s: %v", c.name, c.err))
        }
    }
    
    if len(failures) > 0 {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]any{"status": "unhealthy", "failures": failures})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
```

### Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70   # scale up if CPU > 70%
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

## Summary

- **Multi-stage Docker build**: compile in `golang:alpine`, ship in `scratch` — 5-20 MB final image
- Always use specific image tags (`:1.2.3` not `:latest`) for reproducibility
- **Deployment**: manages Pod replicas, rolling updates, rollbacks
- **Resources**: always set `requests` and `limits` — prevents noisy-neighbor issues
- **Liveness probe**: is the process alive? Fast check, failure = restart pod
- **Readiness probe**: ready to serve traffic? Check DB/Redis, failure = stop routing traffic
- **Graceful shutdown**: `preStop` sleep + `terminationGracePeriodSeconds` gives in-flight requests time to complete

## Exercises

### Easy
1. Write a Dockerfile for a Go HTTP service. Use multi-stage builds. Verify the final image size is under 20 MB.
2. Create a `docker-compose.yml` that runs your service with PostgreSQL and Redis. Verify all services start and the API can reach the database.
3. Deploy a simple Go service to a local Minikube cluster. Write the Deployment and Service manifests. Verify the pod is Running.

### Medium
4. Add liveness and readiness probes to your Deployment. Simulate a database unavailability (stop the PostgreSQL container). Verify that the readiness probe fails and K8s stops routing traffic to the pod, but does NOT restart it.
5. Configure a `HorizontalPodAutoscaler`. Run a load test with `k6`. Observe the HPA scaling pods up and down. Verify no request drops during scale-up.
6. Add a `PodDisruptionBudget` that ensures at least 2 pods are always available during voluntary disruptions (node drains). Test by draining a node and verifying the service stays available.

### Hard
7. Implement a **blue-green deployment** in Kubernetes: two Deployments (`api-blue`, `api-green`), one Service selector that can switch between them. Write a script that: deploys the green version, runs smoke tests, and if they pass, switches the Service selector from blue to green.
8. Set up **Kubernetes RBAC**: create a ServiceAccount for your application pod with only the permissions it needs (read ConfigMaps, get Secrets by name). Verify that the pod cannot list all Secrets or access other namespaces.
