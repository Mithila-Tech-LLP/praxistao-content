# Chapter 129: Kubernetes for Backend Engineers

Chapter 128 introduced the Kubernetes building blocks: Pods, Deployments, Services, ConfigMaps, and probes. That gets your Go service running in a cluster. This chapter is about running it *well*. As a backend engineer you rarely administer the cluster itself — but you own everything about how your service behaves inside it: how the Go runtime interacts with container limits, how configuration reaches your process, how autoscaling decisions get made, and how to debug a pod at 2 AM with nothing but `kubectl`.

## Table of Contents

1. [Where Chapter 128 Left Off](#1-where-chapter-128-left-off)
2. [kubectl — Your Daily Driver](#2-kubectl--your-daily-driver)
3. [The Go Runtime Inside a Container — GOMAXPROCS and GOMEMLIMIT](#3-the-go-runtime-inside-a-container--gomaxprocs-and-gomemlimit)
4. [ConfigMaps and Secrets in Practice](#4-configmaps-and-secrets-in-practice)
5. [Probes Tuned for Go Services](#5-probes-tuned-for-go-services)
6. [Ingress with TLS](#6-ingress-with-tls)
7. [HPA That Actually Works](#7-hpa-that-actually-works)
8. [Rolling Updates and Rollbacks](#8-rolling-updates-and-rollbacks)
9. [Debugging Pods in Production](#9-debugging-pods-in-production)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Where Chapter 128 Left Off

Quick mental map of the objects you already know, and what this chapter adds:

```
                         Internet
                            |
                        [Ingress]  ← TLS termination, host/path routing (§6)
                            |
                        [Service]  ← stable virtual IP, load balancing
                            |
              +-------------+-------------+
              |             |             |
           [Pod v3]      [Pod v3]      [Pod v3]   ← managed by a Deployment
              ^             ^             ^
              |             |             |
        readiness/liveness probes (§5)   HPA adds/removes pods (§7)
              |
        ConfigMap + Secret injected as env/files (§4)
              |
        Go runtime constrained by requests/limits (§3)
```

Everything in Chapter 128 was "make it run." Everything here is "make it behave predictably under load, config changes, deploys, and failures."

---

## 2. kubectl — Your Daily Driver

You will type `kubectl` hundreds of times a week. Learn the workflow, not just individual commands.

### Contexts and namespaces

A *context* is a (cluster, user, namespace) triple stored in `~/.kube/config`. Most engineers have several: local (kind/minikube), staging, production.

```bash
kubectl config get-contexts            # list all contexts
kubectl config use-context staging     # switch cluster
kubectl config set-context --current --namespace=myapp  # default namespace

# Or per-command:
kubectl get pods -n myapp
kubectl get pods -A                    # all namespaces
```

Tip: install `kubectx` and `kubens` — they make switching instant, and a shell prompt that shows the current context has saved many engineers from running a delete against production.

### The read commands you will use constantly

```bash
# What is running?
kubectl get pods                         # names, READY, STATUS, RESTARTS, AGE
kubectl get pods -o wide                 # + node, pod IP
kubectl get pods -l app=api              # filter by label
kubectl get deploy,svc,ingress           # multiple resource types at once

# Why is it in that state? (the single most useful command)
kubectl describe pod api-7d9f8b6c5-x2j4k
# Look at the Events section at the bottom:
#   Warning  Unhealthy  readiness probe failed: HTTP 503
#   Warning  BackOff    restarting failed container

# Logs
kubectl logs api-7d9f8b6c5-x2j4k                   # current container
kubectl logs api-7d9f8b6c5-x2j4k --previous        # the crashed previous run!
kubectl logs -f deploy/api                          # follow, any pod of the deploy
kubectl logs -l app=api --tail=100 --prefix         # all pods, last 100 lines each

# Full YAML as the API server sees it (includes defaults K8s filled in)
kubectl get pod api-7d9f8b6c5-x2j4k -o yaml
```

### The interact commands

```bash
# Port-forward: talk to a pod from your laptop (great for pprof, admin endpoints)
kubectl port-forward deploy/api 8080:8080
curl localhost:8080/healthz

# Exec into a container (only works if the image has a shell — scratch does not, see §9)
kubectl exec -it api-7d9f8b6c5-x2j4k -- /bin/sh

# Run a one-off pod for network testing inside the cluster
kubectl run tmp --rm -it --image=nicolaka/netshoot -- /bin/bash
# From inside: curl http://api.myapp.svc.cluster.local/healthz

# Apply / delete manifests
kubectl apply -f k8s/                    # apply a whole directory
kubectl diff -f k8s/deployment.yaml      # preview what apply would change
kubectl delete -f k8s/deployment.yaml
```

`kubectl diff` before `kubectl apply` is the Kubernetes equivalent of reviewing a `git diff` before committing. Make it a habit.

---

## 3. The Go Runtime Inside a Container — GOMAXPROCS and GOMEMLIMIT

This is the section most Go engineers learn the hard way. The Go runtime was designed for whole machines; containers give it a *slice* of a machine, and by default the runtime cannot always tell the difference.

### The CPU problem

Recall from Chapter 25: `GOMAXPROCS` controls how many OS threads execute Go code simultaneously. By default it equals the number of CPU cores the *node* has — not the CPU limit of your container.

```
Node: 32 cores
Your container limit: cpu: 500m  (half a core)

Default GOMAXPROCS = 32
  → Go scheduler runs goroutines on up to 32 threads
  → cgroup CPU quota allows only 0.5 cores per 100ms period
  → threads burn the quota in ~1.5ms, then EVERYTHING stops until the
    next period — including goroutines in the middle of serving a request
  → symptom: mysterious 50–100ms latency spikes ("CPU throttling")
```

You can see throttling in container metrics as `container_cpu_cfs_throttled_periods_total` (Prometheus, Chapter 120).

**The fix** depends on your Go version:

- **Go 1.25+**: the runtime reads the cgroup CPU limit itself and sets `GOMAXPROCS` accordingly. Nothing to do.
- **Older versions**: import Uber's automaxprocs library:

```go
import _ "go.uber.org/automaxprocs" // sets GOMAXPROCS from the cgroup limit at startup
```

Or set it explicitly in the manifest using the Downward API:

```yaml
env:
- name: GOMAXPROCS
  valueFrom:
    resourceFieldRef:
      resource: limits.cpu   # rounds up: 500m → 1
```

### The memory problem

The Go garbage collector (Chapter 25) grows the heap based on `GOGC` — by default it lets the live heap double before collecting. It has no idea a cgroup memory limit exists. So:

```
Container limit: memory: 512Mi
Live heap: 300Mi → GC waits until heap reaches ~600Mi → kernel OOM-kills
the container first. Status: OOMKilled, exit code 137.
```

Since Go 1.19, `GOMEMLIMIT` sets a *soft* memory limit: as the heap approaches it, the GC runs more aggressively instead of letting the kernel kill you.

```yaml
resources:
  requests:
    memory: 256Mi
  limits:
    memory: 512Mi
env:
- name: GOMEMLIMIT
  value: "460MiB"    # ~90% of the limit — leave headroom for stacks,
                     # runtime overhead, and non-heap memory
```

Rules of thumb:

| Setting | Value |
|---|---|
| `GOMEMLIMIT` | 85–90% of the container memory limit |
| `GOMAXPROCS` | container CPU limit, rounded up (automatic on Go 1.25+) |
| CPU limit | some teams omit it entirely to avoid throttling; if you set one, set `GOMAXPROCS` to match |
| Memory limit | always set it — an unbounded pod can take down the node |

### Requests vs limits, one more time

- **request** = what the scheduler *reserves* for you. Used for bin-packing pods onto nodes and for HPA math (§7).
- **limit** = the hard ceiling. CPU over the limit → throttled. Memory over the limit → OOMKilled.

A sane starting point for a small Go API: `requests: cpu 100m / memory 128Mi`, `limits: memory 256Mi`, measure under load (Chapter 133), then adjust.

---

## 4. ConfigMaps and Secrets in Practice

Chapter 128 showed the objects; here is how configuration actually reaches your Go process, and the trade-offs.

### Option 1: environment variables

```yaml
containers:
- name: api
  envFrom:
  - configMapRef:
      name: api-config       # every key becomes an env var
  env:
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: api-secrets
        key: database-url
```

```go
port := cmp.Or(os.Getenv("PORT"), "8080") // cmp.Or: first non-empty value (Go 1.22+)
```

Simple and 12-factor (Chapter 93 showed reading these with Viper). The catch: **env vars are frozen at container start**. Change the ConfigMap and nothing happens until pods restart.

### Option 2: mounted files

```yaml
containers:
- name: api
  volumeMounts:
  - name: config
    mountPath: /etc/api
    readOnly: true
volumes:
- name: config
  configMap:
    name: api-config
```

Every key in the ConfigMap becomes a file under `/etc/api/`. When you update the ConfigMap, the kubelet rewrites the files within about a minute — no restart. Combine with fsnotify (you saw this exact pattern for secrets in Chapter 131) to hot-reload things like log level:

```go
watcher, _ := fsnotify.NewWatcher()
watcher.Add("/etc/api")
go func() {
    for ev := range watcher.Events {
        if ev.Has(fsnotify.Write) || ev.Has(fsnotify.Create) {
            level := readLogLevel("/etc/api/LOG_LEVEL")
            slog.SetLogLoggerLevel(level)
            slog.Info("config reloaded", "log_level", level)
        }
    }
}()
```

| | Env vars | Mounted files |
|---|---|---|
| Simplicity | Highest | Medium |
| Updates without restart | No | Yes (~1 min sync) |
| Good for | URLs, ports, feature flags read once | log level, rate limits, TLS certs |

### A note on Secrets

Kubernetes Secrets are base64-*encoded*, not encrypted — anyone with `get secrets` RBAC permission reads them in plaintext. For real secret management (Vault, AWS Secrets Manager, External Secrets Operator), Chapter 131 covers the full picture. What you should do at the Kubernetes level:

- Never commit Secret manifests to git; generate them in CI or sync them with External Secrets.
- Restrict RBAC: application ServiceAccounts should not be able to list Secrets.
- Prefer mounting secrets as files over env vars — env vars leak into `kubectl describe pod` error output, crash dumps, and child processes.

---

## 5. Probes Tuned for Go Services

You wrote `/healthz` and `/readyz` handlers in Chapter 128, and Chapter 134 dissects the probe/shutdown timing in depth. Here, the three tuning mistakes that cause most probe-related incidents:

**Mistake 1: liveness probe checks dependencies.** If `/healthz` pings the database and the database has a 5-minute blip, Kubernetes restarts *every pod in the fleet* — turning a partial outage into a full one, plus cold caches. Liveness answers exactly one question: "is this process stuck?" Return 200 unconditionally.

**Mistake 2: probe timeout shorter than the check.** The default `timeoutSeconds` is 1. A `/readyz` that pings Postgres and Redis can legitimately take longer under load. Set `timeoutSeconds: 2–3` and make your handler enforce its own shorter deadline:

```go
ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
defer cancel()
```

**Mistake 3: no startup allowance.** A Go binary starts in milliseconds, but if it runs migrations or warms a cache first, an eager liveness probe kills it mid-startup, forever (CrashLoopBackOff). Give slow starters a `startupProbe` with a generous `failureThreshold` — see Chapter 134 §9 for the full pattern.

A well-behaved Go service also flips its own readiness during shutdown — return 503 from `/readyz` as soon as SIGTERM arrives, so Kubernetes drains traffic while in-flight requests finish:

```go
var shuttingDown atomic.Bool

mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
    if shuttingDown.Load() {
        http.Error(w, "shutting down", http.StatusServiceUnavailable)
        return
    }
    // ... normal dependency checks
})

// in the signal handler, before srv.Shutdown:
shuttingDown.Store(true)
```

---

## 6. Ingress with TLS

Chapter 128 showed a bare HTTP Ingress. Production traffic needs TLS, and the standard pattern is: terminate TLS at the ingress controller, speak plain HTTP inside the cluster.

With **cert-manager** installed (the de facto standard), a single annotation gets you auto-renewed Let's Encrypt certificates:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/proxy-body-size: "8m"        # default is 1m!
    nginx.ingress.kubernetes.io/proxy-read-timeout: "60"
spec:
  ingressClassName: nginx
  tls:
  - hosts: [api.example.com]
    secretName: api-tls        # cert-manager creates and renews this Secret
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

Two ingress-controller defaults that bite Go APIs:

- `proxy-body-size` defaults to 1 MB — file uploads (Chapter 67) get 413 responses at the ingress before your handler ever runs.
- WebSockets (Chapter 65) need `proxy-read-timeout`/`proxy-send-timeout` raised, or the ingress silently closes idle connections after 60 seconds.

Because TLS terminates at the ingress, your Go server sees the client's real IP only in headers: trust `X-Forwarded-For` *only* when the request comes from your ingress — this matters for the rate limiting you will build in Chapter 132.

---

## 7. HPA That Actually Works

Chapter 128 showed a minimal HorizontalPodAutoscaler. Three things make the difference between an HPA that helps and one that flaps.

### 1. The math runs on *requests*, not limits

```
desiredReplicas = ceil(currentReplicas × currentUtilization / targetUtilization)

utilization = actual CPU usage / CPU *request*
```

If your request is `100m` and target utilization is 70%, scaling triggers at 70 millicores of actual usage. Set requests based on measured usage (Chapter 133 shows how to measure), or the HPA scales far too early or far too late.

### 2. Scale on the right signal

CPU works well for request-crunching Go APIs. It works badly for I/O-bound services that sit at 5% CPU while drowning in slow database calls. For those, scale on requests-per-second or queue depth via custom metrics (Prometheus Adapter or KEDA):

```yaml
metrics:
- type: Pods
  pods:
    metric:
      name: http_requests_per_second   # exposed via Prometheus, Chapter 120
    target:
      type: AverageValue
      averageValue: "200"              # scale to keep ≤200 RPS per pod
```

### 3. Tame the flapping with behavior

By default the HPA can double replicas in seconds and scale down 5 minutes later, over and over. The `behavior` block adds damping:

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
  minReplicas: 3
  maxReplicas: 30
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 0     # scale up immediately
      policies:
      - type: Percent
        value: 100                      # at most double per 30s
        periodSeconds: 30
    scaleDown:
      stabilizationWindowSeconds: 300   # wait 5 min of low usage before shrinking
      policies:
      - type: Pods
        value: 1                        # remove at most 1 pod per minute
        periodSeconds: 60
```

Asymmetric on purpose: scale up fast (users are waiting), scale down slowly (avoid thrashing when traffic oscillates).

Watch it work:

```bash
kubectl get hpa api-hpa --watch
# NAME      REFERENCE        TARGETS    MINPODS  MAXPODS  REPLICAS
# api-hpa   Deployment/api   82%/70%    3        30       6
```

---

## 8. Rolling Updates and Rollbacks

Deploying a new version is just changing the image tag; Kubernetes handles the choreography (Chapter 134 covers the zero-downtime tuning of that choreography).

```bash
# Deploy — either apply an updated manifest, or:
kubectl set image deployment/api api=ghcr.io/myorg/myapp:1.3.0

# Watch it roll
kubectl rollout status deployment/api --timeout=5m
# Waiting for deployment "api" rollout to finish: 1 of 3 updated...
# deployment "api" successfully rolled out

# Something is wrong — check history and roll back
kubectl rollout history deployment/api
kubectl rollout undo deployment/api                    # back to previous
kubectl rollout undo deployment/api --to-revision=4    # or a specific one

# Pause a rollout mid-flight to investigate, then resume
kubectl rollout pause deployment/api
kubectl rollout resume deployment/api

# Restart all pods without changing the image (config change pickup)
kubectl rollout restart deployment/api
```

Two habits that make rollbacks trustworthy:

1. **Immutable tags.** Deploy `:1.3.0` or `:git-sha`, never `:latest`. `rollout undo` restores the old *tag* — if the tag's contents changed underneath, the rollback is a lie.
2. **Change-cause annotations.** `kubectl annotate deployment/api kubernetes.io/change-cause="deploy 1.3.0 (fix order race)"` makes `rollout history` readable.

In practice you rarely run `kubectl set image` by hand — your CI/CD pipeline does (that is exactly what Chapter 130 builds). But when the pipeline is broken at 2 AM, these commands are the manual override.

---

## 9. Debugging Pods in Production

The diagnostic loop, in order:

```
kubectl get pods              → what state is it in?
kubectl describe pod <name>   → Events: why is it in that state?
kubectl logs <name> --previous → what did it say before it died?
kubectl debug / port-forward  → interactive investigation
```

### Decoding pod states

| STATUS | Meaning | First move |
|---|---|---|
| `CrashLoopBackOff` | Container starts, exits, K8s retries with backoff | `kubectl logs <pod> --previous` |
| `OOMKilled` (exit 137) | Exceeded memory limit | Raise limit or set `GOMEMLIMIT` (§3); check for leaks with pprof |
| `ImagePullBackOff` | Registry/tag/credentials problem | `describe pod` — the event names the exact error |
| `Pending` | Unschedulable | `describe pod` — insufficient CPU/memory, or affinity rules can't be satisfied |
| `Running` but `READY 0/1` | Readiness probe failing | `describe pod` events + hit `/readyz` via port-forward |
| Exit code 2 / panic in logs | Your Go code panicked | Read the stack trace; it points at the file and line |

### Debugging a scratch/distroless container

Your production image (Chapter 128) has no shell, so `kubectl exec -it ... -- sh` fails. Use an **ephemeral debug container** — it attaches a tool-filled container into the *running pod's* namespaces:

```bash
kubectl debug -it api-7d9f8b6c5-x2j4k \
  --image=nicolaka/netshoot \
  --target=api

# You now share the pod's network namespace:
curl localhost:8080/healthz          # hit the app from "inside"
netstat -tlpn                        # what ports is it listening on?
nslookup postgres.myapp.svc.cluster.local
tcpdump -i eth0 port 5432            # watch the database traffic
```

You can also copy a crashed pod with a new image for offline poking:

```bash
kubectl debug api-7d9f8b6c5-x2j4k --copy-to=api-debug --image=golang:1.23 -- sleep 1d
```

### Live profiling with pprof

Because your service exposes pprof handlers on an internal port (Chapter 26), `port-forward` turns any production mystery into a local profiling session:

```bash
kubectl port-forward deploy/api 6060:6060
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap     # memory
go tool pprof -http=:8081 "http://localhost:6060/debug/pprof/profile?seconds=30"  # CPU
```

This combination — `kubectl port-forward` + pprof — is how you diagnose the "one pod is using 3× the memory of its siblings" class of problem. You will use it again in Chapter 133 when load tests surface bottlenecks.

### Cluster-level signals

```bash
kubectl get events --sort-by=.lastTimestamp -n myapp   # recent events, all objects
kubectl top pods -n myapp                              # live CPU/memory per pod
kubectl top nodes                                      # is a node saturated?
kubectl get endpoints api                              # which pods is the Service actually routing to?
```

`kubectl get endpoints` deserves special mention: if it shows no addresses, your Service selector does not match your pod labels — the single most common "my service is unreachable" cause.

---

## Summary

- Kubernetes fundamentals (Ch 128) get your service running; this chapter is about making it predictable: runtime tuning, config delivery, autoscaling, rollouts, debugging.
- **`kubectl describe` + Events** answers "why is my pod in this state" more often than any other command; `kubectl logs --previous` shows why the last container died.
- The Go runtime does not automatically respect container limits on older versions: set **`GOMAXPROCS`** to the CPU limit (automatic on Go 1.25+, or `automaxprocs`) and **`GOMEMLIMIT`** to ~90% of the memory limit to avoid CPU throttling and OOMKills.
- ConfigMaps as **env vars** are simple but frozen at startup; as **mounted files** they update live and pair with fsnotify for hot reload.
- **Liveness** = "is the process stuck" (no dependency checks); **readiness** = "can it serve traffic" (check dependencies, flip to 503 during shutdown).
- Terminate **TLS at the ingress** with cert-manager; raise `proxy-body-size` and timeouts for uploads and WebSockets.
- HPA math runs on **resource requests**; add a `behavior` block (fast up, slow down) to stop flapping; scale I/O-bound services on RPS, not CPU.
- `kubectl rollout status / history / undo` — with immutable image tags — make deploys reversible.
- **`kubectl debug --target`** gets you a shell "inside" a scratch container; **port-forward + pprof** profiles production pods live.

Next, Chapter 130 automates everything you just did by hand — building the image, running the tests, and rolling out to the cluster — with a CI/CD pipeline.

## Exercises

### Easy
1. Deploy a Go service to a local cluster (kind or minikube) with `requests: cpu 100m, memory 64Mi` and `limits: memory 128Mi`. Use `kubectl top pods` to observe its actual usage, then use `kubectl describe pod` to find the QoS class Kubernetes assigned.
2. Create a ConfigMap with a `LOG_LEVEL` key and mount it as a file at `/etc/api/LOG_LEVEL`. Change the value with `kubectl edit configmap` and verify (with `kubectl exec` or a log line) that the file inside the pod updates within a minute — without a pod restart.
3. Break your Deployment on purpose three ways — a nonexistent image tag, a memory limit of `16Mi`, and a readiness path of `/nope` — and for each, identify the failure using only `kubectl get pods` and `kubectl describe pod`. Write down which Event line gave it away.

### Medium
4. Demonstrate CPU throttling: run a CPU-heavy Go endpoint (e.g., bcrypt with cost 12) in a pod with `limits: cpu: 500m` on a multi-core node, *without* setting GOMAXPROCS. Measure p99 latency under load. Then set `GOMAXPROCS=1` via the Downward API and measure again. Explain the difference in a comment in your manifest.
5. Add hot-reloadable log levels to a Go service: mount a ConfigMap as a file, watch it with fsnotify, and switch `slog` between `Info` and `Debug` on change. Verify by updating the ConfigMap while sending traffic and watching the log output change.
6. Configure an HPA with a `behavior` block (immediate scale-up, 5-minute scale-down stabilization). Drive traffic with k6 (a preview of Chapter 133), watch `kubectl get hpa --watch`, and record how long after traffic stops the replica count actually drops.

### Hard
7. Reproduce and fix an OOMKill: write a Go handler that appends to a package-level slice on every request (a deliberate leak). Deploy with `limits: memory 128Mi` and no `GOMEMLIMIT`; drive traffic until the pod shows `OOMKilled` and RESTARTS increments. Then (a) set `GOMEMLIMIT=110MiB` and observe how behavior changes (GC death spiral instead of instant kill — check CPU), and (b) find the leak with `kubectl port-forward` + `go tool pprof .../heap` and fix it. Document the timeline of both failure modes.
8. Build a "debug playbook" script: given a deployment name, it prints — in order — pod statuses, the last 10 events, logs from any pod with restarts (`--previous`), current endpoints of the matching Service, and `kubectl top` output. Test it against the three broken deployments from exercise 3. Bonus: launch an ephemeral debug container automatically when a pod is Running but not Ready.
