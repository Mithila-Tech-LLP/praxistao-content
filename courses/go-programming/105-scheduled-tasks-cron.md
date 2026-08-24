# Chapter 105: Scheduled Tasks, Cron Jobs, and Distributed Schedulers

Every production system eventually needs work that runs on a schedule: nightly reports, hourly cache warmups, weekly email digests. In a single-process app, a cron job is trivial. In a multi-instance deployment, it becomes a distributed coordination problem. This chapter covers the full spectrum — from a simple in-process cron to a Kubernetes `CronJob` with leader election.

## Table of Contents

1. [Types of Scheduled Work](#1-types-of-scheduled-work)
2. [robfig/cron Library](#2-robfigcron-library)
3. [asynq Scheduler](#3-asynq-scheduler)
4. [The Distributed Scheduler Problem](#4-the-distributed-scheduler-problem)
5. [Solution 1: Distributed Lock](#5-solution-1-distributed-lock)
6. [Solution 2: Dedicated Scheduler Instance](#6-solution-2-dedicated-scheduler-instance)
7. [Kubernetes CronJob](#7-kubernetes-cronjob)
8. [Leader Election Pattern](#8-leader-election-pattern)
9. [Monitoring Scheduled Jobs](#9-monitoring-scheduled-jobs)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Types of Scheduled Work

```
One-time delayed task:
  "Send this email in 10 minutes"
  "Expire this token after 24 hours"
  → Use: asynq.ProcessIn(), time.AfterFunc, or a delayed queue

Recurring cron job:
  "Generate daily report at 8 AM"
  "Clean up expired sessions every hour"
  → Use: robfig/cron, asynq Scheduler, Kubernetes CronJob

Rate-limited background processing:
  "Process at most 100 invoices per minute"
  "Send 50 push notifications per second"
  → Use: worker pool + ticker + token bucket (see Chapter 92)

Event-triggered delayed work:
  "If no activity for 30 minutes, close the session"
  "Retry failed webhooks with exponential backoff"
  → Use: asynq with ProcessIn() and retry policies
```

Choosing the right tool:

```
Single instance app, simple schedule  → robfig/cron
Multi-instance app, Redis available   → asynq Scheduler + distributed lock
Kubernetes deployment                 → Kubernetes CronJob (cleanest solution)
Complex orchestration / workflows     → asynq + custom state machine
```

---

## 2. robfig/cron Library

`robfig/cron` is the standard Go library for in-process cron scheduling. It parses cron expressions and runs functions in goroutines.

```go
import "github.com/robfig/cron/v3"

func main() {
    // WithSeconds() enables 6-field expressions: sec min hour dom month dow
    // Without it: 5-field (min hour dom month dow) — standard cron syntax
    c := cron.New(cron.WithSeconds())

    // Cron expression format (with seconds):
    //   ┌──────────── second (0-59)
    //   │ ┌────────── minute (0-59)
    //   │ │ ┌──────── hour (0-23)
    //   │ │ │ ┌────── day of month (1-31)
    //   │ │ │ │ ┌──── month (1-12 or JAN-DEC)
    //   │ │ │ │ │ ┌── day of week (0-7 or SUN-SAT, 0 and 7 = Sunday)
    //   │ │ │ │ │ │
    //   * * * * * *

    // Every day at 8:00 AM UTC
    c.AddFunc("0 0 8 * * *", func() {
        log.Println("running daily report:", time.Now())
        if err := generateDailyReport(context.Background()); err != nil {
            log.Printf("daily report failed: %v", err)
        }
    })

    // Every hour at :00
    c.AddFunc("0 0 * * * *", func() {
        cleanupExpiredSessions(context.Background())
    })

    // Every 5 minutes
    c.AddFunc("0 */5 * * * *", func() {
        warmupCache(context.Background())
    })

    // Every Monday at 9 AM
    c.AddFunc("0 0 9 * * MON", func() {
        sendWeeklyNewsletter(context.Background())
    })

    c.Start()

    // Block until shutdown
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
    <-sigs

    // Stop() waits for running jobs to complete before returning
    ctx := c.Stop()
    <-ctx.Done()
    log.Println("scheduler stopped gracefully")
}
```

### Cron expression cheat sheet

```
Expression           Meaning
───────────────────────────────────────────────────────
0 0 8 * * *          Daily at 8:00:00 AM
0 30 8 * * *         Daily at 8:30:00 AM
0 0 * * * *          Every hour on the hour
0 */15 * * * *       Every 15 minutes
0 0 9 * * MON        Every Monday at 9:00 AM
0 0 0 1 * *          First day of every month at midnight
0 0 0 * * SAT,SUN    Every Saturday and Sunday at midnight
@every 30s           Every 30 seconds (special syntax — requires no WithSeconds())
@every 1h30m         Every 1.5 hours
@daily               Every day at midnight (equivalent to 0 0 0 * * *)
@weekly              Every Sunday at midnight
```

### Named job entries and removal

```go
type Scheduler struct {
    cron    *cron.Cron
    entries map[string]cron.EntryID
    mu      sync.Mutex
}

func NewScheduler() *Scheduler {
    return &Scheduler{
        cron:    cron.New(cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(log.Default()))),
        entries: make(map[string]cron.EntryID),
    }
}

func (s *Scheduler) Add(name, spec string, fn func()) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    id, err := s.cron.AddFunc(spec, fn)
    if err != nil {
        return fmt.Errorf("add cron %q: %w", name, err)
    }
    s.entries[name] = id
    return nil
}

func (s *Scheduler) Remove(name string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if id, ok := s.entries[name]; ok {
        s.cron.Remove(id)
        delete(s.entries, name)
    }
}
```

---

## 3. asynq Scheduler

The asynq Scheduler (covered briefly in Chapter 94) integrates scheduling with asynq's task queue — scheduled tasks are enqueued and processed by your existing worker fleet. This separates schedule triggering from task execution.

```go
import "github.com/hibiken/asynq"

func startScheduler(redisAddr string) {
    scheduler := asynq.NewScheduler(
        asynq.RedisClientOpt{Addr: redisAddr},
        &asynq.SchedulerOpts{
            LogLevel: asynq.WarnLevel,
            // PostEnqueueFunc is called after each task is enqueued
            PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
                if err != nil {
                    slog.Error("scheduler failed to enqueue task",
                        "type", info.Type,
                        "err", err,
                    )
                    return
                }
                slog.Info("scheduled task enqueued",
                    "type", info.Type,
                    "id", info.ID,
                    "next_run", info.NextProcessAt,
                )
            },
        },
    )

    // Register recurring tasks

    // Daily report at 8 AM UTC
    reportTask := asynq.NewTask("report:daily", nil)
    if _, err := scheduler.Register("0 8 * * *", reportTask,
        asynq.Queue("default"),
        asynq.Unique(24*time.Hour), // prevent duplicate if scheduler restarts
    ); err != nil {
        log.Fatalf("register report task: %v", err)
    }

    // Session cleanup every hour
    cleanupTask := asynq.NewTask("session:cleanup", nil)
    scheduler.Register("0 * * * *", cleanupTask)

    // Weekly newsletter every Monday at 9 AM
    type NewsletterPayload struct{ BatchSize int }
    payload, _ := json.Marshal(NewsletterPayload{BatchSize: 500})
    newsletterTask := asynq.NewTask("email:newsletter", payload)
    scheduler.Register("0 9 * * MON", newsletterTask,
        asynq.Queue("emails"),
        asynq.MaxRetry(2),
        asynq.Timeout(30*time.Minute),
    )

    // Use @every for interval-based schedules
    cacheTask := asynq.NewTask("cache:warmup", nil)
    scheduler.Register("@every 1h", cacheTask)

    if err := scheduler.Run(); err != nil {
        log.Fatalf("scheduler: %v", err)
    }
    // scheduler.Run() blocks — call Shutdown() from a signal handler
}
```

The asynq scheduler and worker are separate processes (or goroutines):

```
Scheduler process           Worker process (N instances)
─────────────────           ──────────────────────────────
At 8:00 AM:                 
  Enqueue "report:daily" →  Redis Queue → Worker picks up task
                                          → Runs generateDailyReport()
```

This split means the scheduler doesn't execute the work — it just enqueues it. Workers handle execution, retries, and concurrency.

---

## 4. The Distributed Scheduler Problem

When you deploy multiple instances of your app, every instance runs the same cron setup. Without coordination:

```
Deployment: 3 instances of order-service

Instance 1: cron fires "send daily digest" at 8:00:00
Instance 2: cron fires "send daily digest" at 8:00:00
Instance 3: cron fires "send daily digest" at 8:00:00

Result:
  - 3× the emails sent to every user
  - 3× the load on the email provider
  - 3× the database queries
  - Angry customers
```

This is the **distributed scheduler problem**: you want exactly one instance to run the scheduled job, but multiple instances are running. Solutions:

```
Option A: Distributed lock before running the job
  → Simple, works with any scheduler
  → One instance grabs the lock, others skip

Option B: Dedicated scheduler process / pod
  → Only one instance is the scheduler
  → Others are pure workers

Option C: Leader election
  → Instances elect one leader dynamically
  → Leader runs the scheduler; if it dies, a new leader is elected
```

---

## 5. Solution 1: Distributed Lock

Before running any scheduled job, acquire a Redis lock. Only the instance that holds the lock executes the work. Others see the lock exists and skip.

```go
import "github.com/redis/go-redis/v9"

type DistributedScheduler struct {
    cron   *cron.Cron
    rdb    *redis.Client
    nodeID string // unique per instance: hostname or UUID
    logger *slog.Logger
}

func NewDistributedScheduler(rdb *redis.Client) *DistributedScheduler {
    hostname, _ := os.Hostname()
    return &DistributedScheduler{
        cron:   cron.New(cron.WithSeconds()),
        rdb:    rdb,
        nodeID: hostname,
        logger: slog.Default(),
    }
}

// withLock wraps fn in a distributed lock.
// The lock key includes the job name; TTL should exceed the job's expected duration.
func (s *DistributedScheduler) withLock(jobName string, ttl time.Duration, fn func()) func() {
    return func() {
        ctx := context.Background()
        lockKey := fmt.Sprintf("cron:lock:%s", jobName)

        // SET key nodeID EX ttl NX
        // Only sets if the key does not exist → atomic "grab the lock"
        acquired, err := s.rdb.SetNX(ctx, lockKey, s.nodeID, ttl).Result()
        if err != nil {
            s.logger.Error("lock acquisition failed", "job", jobName, "err", err)
            return
        }
        if !acquired {
            // Another instance holds the lock — this fire is a no-op
            s.logger.Debug("skipping job — lock held by another instance", "job", jobName)
            return
        }

        s.logger.Info("running scheduled job", "job", jobName, "node", s.nodeID)
        defer func() {
            // Release the lock only if we still own it (prevents releasing another node's lock)
            luaScript := `
                if redis.call("GET", KEYS[1]) == ARGV[1] then
                    return redis.call("DEL", KEYS[1])
                else
                    return 0
                end`
            s.rdb.Eval(ctx, luaScript, []string{lockKey}, s.nodeID)
        }()

        fn()
    }
}

func (s *DistributedScheduler) Register(spec, jobName string, fn func()) error {
    // TTL of 55 seconds for a per-minute job — expires before next fire
    lockTTL := 55 * time.Second

    _, err := s.cron.AddFunc(spec, s.withLock(jobName, lockTTL, fn))
    return err
}

// Usage
func startDistributedCron(rdb *redis.Client) {
    sched := NewDistributedScheduler(rdb)

    sched.Register("0 0 8 * * *", "daily-report", func() {
        if err := generateDailyReport(context.Background()); err != nil {
            slog.Error("daily report failed", "err", err)
        }
    })

    sched.Register("0 0 * * * *", "session-cleanup", func() {
        cleanupExpiredSessions(context.Background())
    })

    sched.cron.Start()
}
```

The lock TTL is critical:
- Too short: lock expires during a slow job run — another instance picks it up → double execution
- Too long: if the job crashes, the lock stays held until TTL expires → next fires are skipped

Set TTL slightly less than the cron interval. For a job that runs every hour and takes at most 10 minutes, a 55-minute TTL works.

---

## 6. Solution 2: Dedicated Scheduler Instance

The simplest production solution: separate your app binary into roles.

```go
// main.go
func main() {
    role := os.Getenv("APP_ROLE") // "worker" | "scheduler" | "api"

    switch role {
    case "scheduler":
        runScheduler()
    case "worker":
        runWorker()
    case "api":
        runAPIServer()
    default:
        // Development: run everything in one process
        go runScheduler()
        go runWorker()
        runAPIServer()
    }
}

func runScheduler() {
    // Only ONE instance of this process runs in production
    slog.Info("starting scheduler")
    startScheduler(os.Getenv("REDIS_ADDR"))
}

func runWorker() {
    // Multiple instances run in parallel
    slog.Info("starting worker")
    startWorker(os.Getenv("REDIS_ADDR"))
}
```

In Docker Compose:

```yaml
# docker-compose.yml
services:
  scheduler:
    image: myapp:latest
    environment:
      APP_ROLE: scheduler
      REDIS_ADDR: redis:6379
    deploy:
      replicas: 1    # exactly one scheduler

  worker:
    image: myapp:latest
    environment:
      APP_ROLE: worker
      REDIS_ADDR: redis:6379
    deploy:
      replicas: 5    # scale workers freely
```

This is the pattern the asynq Scheduler is designed for: the scheduler process enqueues tasks, workers process them. Scaling workers doesn't cause duplicate scheduling.

---

## 7. Kubernetes CronJob

For recurring batch work in Kubernetes, the `CronJob` resource is the idiomatic solution. Kubernetes manages execution, history, and concurrency policy.

```yaml
# k8s/daily-report-cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: daily-report
  namespace: production
spec:
  schedule: "0 8 * * *"         # 8:00 AM UTC daily (standard 5-field cron)
  timeZone: "UTC"                # explicit timezone (Kubernetes 1.27+)

  concurrencyPolicy: Forbid      # if previous run is still active, skip the new one
                                 # Replace: kill previous, start new
                                 # Allow: run in parallel (dangerous for non-idempotent jobs)

  successfulJobsHistoryLimit: 3  # keep last 3 successful job records
  failedJobsHistoryLimit: 5      # keep last 5 failed job records (for debugging)

  startingDeadlineSeconds: 300   # if missed by more than 5 minutes, skip this fire

  jobTemplate:
    spec:
      backoffLimit: 2            # retry up to 2 times on failure
      activeDeadlineSeconds: 3600 # kill job if it runs longer than 1 hour

      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: report-generator
              image: myapp:latest
              command: ["./myapp", "--run-once", "daily-report"]
              env:
                - name: DATABASE_URL
                  valueFrom:
                    secretKeyRef:
                      name: app-secrets
                      key: database-url
                - name: REDIS_ADDR
                  value: "redis:6379"
              resources:
                requests:
                  memory: "256Mi"
                  cpu: "100m"
                limits:
                  memory: "512Mi"
                  cpu: "500m"
```

```go
// The app binary supports a --run-once flag for CronJob invocation
func main() {
    runOnce := flag.String("run-once", "", "run a named job once and exit")
    flag.Parse()

    if *runOnce != "" {
        if err := runNamedJob(context.Background(), *runOnce); err != nil {
            slog.Error("job failed", "job", *runOnce, "err", err)
            os.Exit(1)
        }
        os.Exit(0)
    }

    // Normal long-running server mode
    runServer()
}

func runNamedJob(ctx context.Context, name string) error {
    switch name {
    case "daily-report":
        return generateDailyReport(ctx)
    case "session-cleanup":
        return cleanupExpiredSessions(ctx)
    case "weekly-newsletter":
        return sendWeeklyNewsletter(ctx)
    default:
        return fmt.Errorf("unknown job: %s", name)
    }
}
```

`concurrencyPolicy: Forbid` is almost always the right choice. It prevents overlapping runs of the same job, which is dangerous for non-idempotent work like "process all unpaid invoices."

---

## 8. Leader Election Pattern

When you can't use a dedicated scheduler pod (e.g., you have a single deployment type that must scale), leader election designates one instance as the scheduler dynamically.

### Redis-based leader election

```go
type LeaderElector struct {
    rdb      *redis.Client
    nodeID   string
    leaseKey string
    leaseTTL time.Duration
    isLeader atomic.Bool
    logger   *slog.Logger
}

func NewLeaderElector(rdb *redis.Client, leaseKey string) *LeaderElector {
    hostname, _ := os.Hostname()
    return &LeaderElector{
        rdb:      rdb,
        nodeID:   fmt.Sprintf("%s-%d", hostname, os.Getpid()),
        leaseKey: leaseKey,
        leaseTTL: 15 * time.Second,
        logger:   slog.Default(),
    }
}

// Run continuously attempts to acquire or renew the leader lease.
// When leadership is gained, onBecomeLeader() is called.
// When leadership is lost, onLoseLeadership() is called.
func (e *LeaderElector) Run(ctx context.Context, onBecomeLeader, onLoseLeadership func()) {
    ticker := time.NewTicker(5 * time.Second) // try to elect every 5 seconds
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            e.resign(context.Background())
            return
        case <-ticker.C:
            wasLeader := e.isLeader.Load()
            nowLeader := e.tryAcquireOrRenew(ctx)

            if !wasLeader && nowLeader {
                e.logger.Info("became leader", "node", e.nodeID)
                onBecomeLeader()
            } else if wasLeader && !nowLeader {
                e.logger.Warn("lost leadership", "node", e.nodeID)
                onLoseLeadership()
            }
        }
    }
}

func (e *LeaderElector) tryAcquireOrRenew(ctx context.Context) bool {
    // Try to set as new leader (NX = only if not exists)
    acquired, err := e.rdb.SetNX(ctx, e.leaseKey, e.nodeID, e.leaseTTL).Result()
    if err != nil {
        e.isLeader.Store(false)
        return false
    }
    if acquired {
        e.isLeader.Store(true)
        return true
    }

    // Not a new acquisition — check if we're the current holder and renew
    current, err := e.rdb.Get(ctx, e.leaseKey).Result()
    if err != nil {
        e.isLeader.Store(false)
        return false
    }
    if current == e.nodeID {
        // Renew our lease
        e.rdb.Expire(ctx, e.leaseKey, e.leaseTTL)
        e.isLeader.Store(true)
        return true
    }

    e.isLeader.Store(false)
    return false
}

func (e *LeaderElector) resign(ctx context.Context) {
    luaScript := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        end
        return 0`
    e.rdb.Eval(ctx, luaScript, []string{e.leaseKey}, e.nodeID)
    e.isLeader.Store(false)
}

// Usage
func main() {
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    elector := NewLeaderElector(rdb, "scheduler:leader-lease")

    var cronInstance *cron.Cron
    var cronMu sync.Mutex

    go elector.Run(context.Background(),
        // Become leader: start the scheduler
        func() {
            cronMu.Lock()
            defer cronMu.Unlock()
            cronInstance = cron.New(cron.WithSeconds())
            cronInstance.AddFunc("0 0 8 * * *", func() { generateDailyReport(context.Background()) })
            cronInstance.Start()
        },
        // Lose leadership: stop the scheduler
        func() {
            cronMu.Lock()
            defer cronMu.Unlock()
            if cronInstance != nil {
                cronInstance.Stop()
                cronInstance = nil
            }
        },
    )

    // ... rest of the app
}
```

---

## 9. Monitoring Scheduled Jobs

A job that silently fails to run is worse than a job that fails loudly. Monitor both execution and absence.

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    jobLastRunTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "scheduled_job_last_run_timestamp_seconds",
        Help: "Unix timestamp of the last successful run of each scheduled job",
    }, []string{"job"})

    jobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "scheduled_job_duration_seconds",
        Help:    "Duration of scheduled job runs",
        Buckets: prometheus.DefBuckets,
    }, []string{"job", "status"})

    jobRunsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "scheduled_job_runs_total",
        Help: "Total number of scheduled job runs",
    }, []string{"job", "status"})
)

// instrumentedJob wraps a job function with metrics
func instrumentedJob(name string, fn func() error) func() {
    return func() {
        start := time.Now()
        err := fn()
        duration := time.Since(start)
        status := "success"

        if err != nil {
            status = "error"
            slog.Error("scheduled job failed", "job", name, "err", err, "duration", duration)
        } else {
            jobLastRunTime.WithLabelValues(name).Set(float64(time.Now().Unix()))
            slog.Info("scheduled job completed", "job", name, "duration", duration)
        }

        jobDuration.WithLabelValues(name, status).Observe(duration.Seconds())
        jobRunsTotal.WithLabelValues(name, status).Inc()
    }
}

// Usage with robfig/cron
func setupCronWithMetrics(c *cron.Cron) {
    c.AddFunc("0 0 8 * * *", instrumentedJob("daily-report", func() error {
        return generateDailyReport(context.Background())
    }))

    c.AddFunc("0 0 * * * *", instrumentedJob("session-cleanup", func() error {
        return cleanupExpiredSessions(context.Background())
    }))
}
```

### Alerting on missed runs

In Prometheus/Alertmanager, alert when `scheduled_job_last_run_timestamp_seconds` falls outside the expected window:

```yaml
# prometheus/alerts.yaml
groups:
  - name: scheduled_jobs
    rules:
      # Alert if daily-report hasn't run in 26 hours (allows for 2h drift)
      - alert: ScheduledJobMissed
        expr: |
          (time() - scheduled_job_last_run_timestamp_seconds{job="daily-report"}) > 93600
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Scheduled job {{ $labels.job }} hasn't run in over 26 hours"
          description: "Last run was {{ $value | humanizeDuration }} ago"

      # Alert if a job is taking too long (potential hung run)
      - alert: ScheduledJobSlowRun
        expr: |
          histogram_quantile(0.99,
            rate(scheduled_job_duration_seconds_bucket[1h])
          ) > 300
        for: 10m
        annotations:
          summary: "Scheduled job {{ $labels.job }} p99 exceeds 5 minutes"
```

### Dead man's switch with healthcheck.io / Better Uptime

A simple pattern: ping a URL at the end of a successful run. If the ping is missing, you get an alert without Prometheus:

```go
func dailyReportWithPing(ctx context.Context) func() error {
    return func() error {
        if err := generateDailyReport(ctx); err != nil {
            return err
        }

        // Ping dead man's switch endpoint — if this doesn't come, alert fires
        pingURL := os.Getenv("HEALTHCHECK_PING_URL")
        if pingURL != "" {
            resp, err := http.Get(pingURL)
            if err != nil {
                slog.Warn("healthcheck ping failed", "err", err)
            } else {
                resp.Body.Close()
            }
        }
        return nil
    }
}
```

---

## Summary

- **Three types of scheduled work**: one-time delayed tasks, recurring cron jobs, rate-limited background processing
- **`robfig/cron`**: `cron.New(cron.WithSeconds())`, `c.AddFunc(spec, fn)`, `c.Start()` / `c.Stop()` — clean graceful shutdown via `ctx := c.Stop(); <-ctx.Done()`
- **asynq Scheduler**: separates scheduling (enqueue) from execution (worker) — scales workers independently, natural integration with existing task queue
- **The distributed problem**: multiple instances all fire the same cron → duplicate side effects
- **Solution 1 — Redis lock**: `SetNX(lockKey, nodeID, ttl)` — only the instance that grabs the lock runs the job; release with Lua script to avoid releasing another node's lock
- **Solution 2 — dedicated instance**: `APP_ROLE=scheduler` runs one instance; workers are separate — simplest and most reliable
- **Kubernetes CronJob**: `concurrencyPolicy: Forbid` prevents overlapping runs; `startingDeadlineSeconds` handles missed fires; `backoffLimit` handles transient failures
- **Leader election**: Redis `SetNX` + periodic renewal; `onBecomeLeader` starts cron, `onLoseLeadership` stops it
- **Monitoring**: gauge for `last_run_timestamp`, histogram for duration, alert when gap exceeds expected interval

---

## Exercises

### Easy
1. Set up `robfig/cron` with a job that runs every 10 seconds and prints the current time. Verify graceful shutdown: send SIGTERM and confirm the scheduler waits for the in-flight run to complete before exiting.
2. Write a cron expression cheat sheet from memory for: every 15 minutes, daily at midnight, weekdays at 9 AM, first day of month at noon. Use `cron.New()` without `WithSeconds()` and add all four. Verify the next run times with `c.Entries()`.
3. Implement `instrumentedJob` from section 9 and wrap a sample job. Run it 5 times, with 2 failures (return a non-nil error). Verify the counter shows `{status="success"}: 3` and `{status="error"}: 2`.

### Medium
4. Implement `DistributedScheduler` from section 5. Start 3 goroutines, each simulating an app instance with its own scheduler. Run a 1-second recurring job. Use an `atomic.Int64` counter to count actual executions. After 10 seconds, assert that the counter is between 9 and 11 (not 27–33 from triple execution).
5. Build the `LeaderElector` from section 8. Start 3 goroutines simulating instances. Kill the current leader by cancelling its context. Verify that a new leader is elected within 15 seconds (the lease TTL) and the scheduler continues running.
6. Write a `runNamedJob` CLI entrypoint that accepts `--run-once <job-name>` and runs the named job exactly once, exiting 0 on success and 1 on failure. Write a test that invokes it via `exec.Command` and checks the exit code and output for a job that succeeds and one that fails.

### Hard
7. Implement a **fault-tolerant distributed scheduler**: the scheduler stores the last run time of each job in Redis (`HSET scheduler:last_runs <job> <timestamp>`). On startup, check if any job's last run is overdue (should have run since last startup). If so, run it immediately. This handles the case where the scheduler pod was down during a scheduled fire.
8. Build a **job execution history store**: every job run writes a record to PostgreSQL (`job_name`, `started_at`, `finished_at`, `status`, `error`). Expose a `GET /admin/jobs` endpoint that returns the last 10 runs per job. Add an alert query that finds jobs with 3 consecutive failures and marks them as `suspended` so they stop firing until manually re-enabled.
