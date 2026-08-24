# Chapter 94: asynq — Task Queues and Background Jobs

asynq is a Redis-backed distributed task queue for Go. It handles background job processing with priorities, retries, scheduling, and monitoring. It's simpler to operate than Kafka and perfect for tasks that need reliable background execution: sending emails, resizing images, generating reports, processing webhooks.

## Table of Contents

1. [asynq vs Kafka](#1-asynq-vs-kafka)
2. [Defining Tasks](#2-defining-tasks)
3. [Producer: Enqueuing Jobs](#3-producer-enqueuing-jobs)
4. [Consumer: Processing Jobs](#4-consumer-processing-jobs)
5. [Scheduled Jobs (Cron)](#5-scheduled-jobs-cron)
6. [Monitoring and Retry](#6-monitoring-and-retry)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. asynq vs Kafka

| | asynq | Kafka |
|---|---|---|
| Backend | Redis | Kafka cluster |
| Retention | Until processed + TTL | Time-based (default 7 days) |
| Multiple consumers | Per queue | Per consumer group |
| Ordering | Per queue, FIFO | Per partition |
| Replay | No | Yes |
| Cron | Built-in | External |
| Setup | 1 Redis server | Kafka + ZooKeeper/KRaft |
| Best for | Background tasks, emails, webhooks | Event streams, audit logs, fan-out |

Use asynq when you need **reliable background job execution**. Use Kafka when you need **event streaming with multiple consumers and replay**.

---

## 2. Defining Tasks

```go
import "github.com/hibiken/asynq"

// Task type constants
const (
    TypeEmailWelcome     = "email:welcome"
    TypeEmailOrderConf   = "email:order_confirmation"
    TypeImageResize      = "image:resize"
    TypeReportGenerate   = "report:generate"
    TypeWebhookDeliver   = "webhook:deliver"
)

// Task payloads are typed structs serialized to JSON
type WelcomeEmailPayload struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
}

type OrderConfirmationPayload struct {
    OrderID    string  `json:"order_id"`
    CustomerID string  `json:"customer_id"`
    Total      float64 `json:"total"`
}

type ImageResizePayload struct {
    SourcePath string `json:"source_path"`
    Width      int    `json:"width"`
    Height     int    `json:"height"`
    OutputPath string `json:"output_path"`
}

// Factory functions build asynq.Task from typed payloads
func NewWelcomeEmailTask(p WelcomeEmailPayload) (*asynq.Task, error) {
    payload, err := json.Marshal(p)
    if err != nil { return nil, err }
    return asynq.NewTask(TypeEmailWelcome, payload), nil
}

func NewOrderConfirmationTask(p OrderConfirmationPayload, opts ...asynq.Option) (*asynq.Task, error) {
    payload, err := json.Marshal(p)
    if err != nil { return nil, err }
    return asynq.NewTask(TypeEmailOrderConf, payload, opts...), nil
}
```

---

## 3. Producer: Enqueuing Jobs

```go
// Create a client (producer)
func newAsynqClient(redisAddr string) *asynq.Client {
    return asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}

type TaskProducer struct{ client *asynq.Client }

func (p *TaskProducer) EnqueueWelcomeEmail(ctx context.Context, userID int64, email, name string) error {
    task, err := NewWelcomeEmailTask(WelcomeEmailPayload{
        UserID: userID,
        Email:  email,
        Name:   name,
    })
    if err != nil { return err }
    
    info, err := p.client.EnqueueContext(ctx, task,
        asynq.Queue("emails"),      // route to specific queue
        asynq.MaxRetry(3),          // retry 3 times on failure
        asynq.Timeout(30*time.Second), // task must finish within 30s
    )
    if err != nil { return fmt.Errorf("enqueue: %w", err) }
    
    slog.Info("enqueued welcome email", "task_id", info.ID, "queue", info.Queue)
    return nil
}

func (p *TaskProducer) EnqueueOrderConfirmation(ctx context.Context, orderID string, total float64) error {
    task, err := NewOrderConfirmationTask(
        OrderConfirmationPayload{OrderID: orderID, Total: total},
        asynq.ProcessIn(5*time.Second), // delay 5 seconds before processing
    )
    if err != nil { return err }
    
    _, err = p.client.EnqueueContext(ctx, task,
        asynq.Queue("emails"),
        asynq.Unique(24*time.Hour), // deduplicate: if same task already queued, skip
    )
    return err
}

func (p *TaskProducer) EnqueueImageResize(ctx context.Context, src, dst string, w, h int) error {
    task, _ := json.Marshal(ImageResizePayload{
        SourcePath: src, OutputPath: dst,
        Width: w, Height: h,
    })
    t := asynq.NewTask(TypeImageResize, task,
        asynq.Queue("media"),
        asynq.MaxRetry(5),
        asynq.Retention(24*time.Hour), // keep completed task info for 24h
    )
    _, err := p.client.EnqueueContext(ctx, t)
    return err
}

func (p *TaskProducer) Close() error { return p.client.Close() }
```

---

## 4. Consumer: Processing Jobs

```go
// Handler for each task type
type EmailHandler struct{ mailer Mailer }

func (h *EmailHandler) HandleWelcomeEmail(ctx context.Context, t *asynq.Task) error {
    var p WelcomeEmailPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }
    
    return h.mailer.Send(ctx, p.Email,
        "Welcome to our service!",
        fmt.Sprintf("Hi %s, thanks for joining!", p.Name),
    )
}

func (h *EmailHandler) HandleOrderConfirmation(ctx context.Context, t *asynq.Task) error {
    var p OrderConfirmationPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }
    return h.mailer.SendOrderConfirmation(ctx, p.CustomerID, p.OrderID, p.Total)
}

type ImageHandler struct{ storage Storage }

func (h *ImageHandler) HandleResize(ctx context.Context, t *asynq.Task) error {
    var p ImageResizePayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil { return err }
    
    return h.storage.ResizeImage(ctx, p.SourcePath, p.OutputPath, p.Width, p.Height)
}

// Wire the server
func newAsynqServer(redisAddr string) *asynq.Server {
    return asynq.NewServer(
        asynq.RedisClientOpt{Addr: redisAddr},
        asynq.Config{
            // Priority queues: critical > emails > media > default
            Queues: map[string]int{
                "critical": 6,
                "emails":   3,
                "media":    2,
                "default":  1,
            },
            Concurrency: 10, // process up to 10 tasks concurrently
            
            // Custom retry delay: exponential backoff
            RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
                return time.Duration(math.Pow(2, float64(n))) * time.Second
            },
            
            // Error handler
            ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, t *asynq.Task, err error) {
                slog.Error("task failed",
                    "type", t.Type(),
                    "payload", string(t.Payload()),
                    "err", err,
                )
            }),
        },
    )
}

func startWorker(redisAddr string, emailHandler *EmailHandler, imageHandler *ImageHandler) {
    server := newAsynqServer(redisAddr)
    
    mux := asynq.NewServeMux()
    
    // Register handlers
    mux.HandleFunc(TypeEmailWelcome,   emailHandler.HandleWelcomeEmail)
    mux.HandleFunc(TypeEmailOrderConf, emailHandler.HandleOrderConfirmation)
    mux.HandleFunc(TypeImageResize,    imageHandler.HandleResize)
    
    // Middleware
    mux.Use(loggingMiddleware)
    mux.Use(metricsMiddleware)
    
    if err := server.Run(mux); err != nil {
        log.Fatal(err)
    }
}

func loggingMiddleware(h asynq.Handler) asynq.Handler {
    return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
        start := time.Now()
        err := h.ProcessTask(ctx, t)
        slog.Info("task processed",
            "type", t.Type(),
            "duration", time.Since(start),
            "err", err,
        )
        return err
    })
}
```

---

## 5. Scheduled Jobs (Cron)

```go
// Periodic tasks using asynq's scheduler
func startScheduler(redisAddr string) *asynq.Scheduler {
    scheduler := asynq.NewScheduler(
        asynq.RedisClientOpt{Addr: redisAddr},
        &asynq.SchedulerOpts{
            LogLevel: asynq.WarnLevel,
        },
    )
    
    // Generate daily report at 8 AM UTC
    reportTask := asynq.NewTask("report:daily", nil)
    scheduler.Register("0 8 * * *", reportTask,
        asynq.Queue("default"),
        asynq.Unique(24*time.Hour), // prevent duplicate if scheduler runs multiple instances
    )
    
    // Clean up expired sessions every hour
    cleanupTask := asynq.NewTask("session:cleanup", nil)
    scheduler.Register("0 * * * *", cleanupTask)
    
    // Send weekly newsletter every Monday at 9 AM
    newsletterTask := asynq.NewTask("email:newsletter", nil)
    scheduler.Register("0 9 * * MON", newsletterTask,
        asynq.Queue("emails"),
    )
    
    return scheduler
}

func main() {
    scheduler := startScheduler("localhost:6379")
    
    if err := scheduler.Start(); err != nil { log.Fatal(err) }
    
    // Block until shutdown signal
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
    <-sigs
    scheduler.Shutdown()
}
```

---

## 6. Monitoring and Retry

```go
// Inspect queue state via the Inspector API
func monitorQueues(redisAddr string) {
    inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
    
    queues, _ := inspector.Queues()
    for _, q := range queues {
        info, _ := inspector.GetQueueInfo(q)
        slog.Info("queue stats",
            "queue", q,
            "size",    info.Size,
            "active",  info.Active,
            "pending", info.Pending,
            "failed",  info.Failed,
            "retry",   info.Retry,
        )
    }
    
    // Retry a specific failed task
    // inspector.RunTask("default", taskID)
    
    // Delete all failed tasks in a queue
    // inspector.DeleteAllFailedTasks("emails")
    
    // Pause a queue (stop processing new tasks)
    // inspector.PauseQueue("emails")
}

// Asynqmon: web UI for monitoring (run alongside your worker)
// go run github.com/hibiken/asynqmon@latest --redis-addr=localhost:6379 --port=8888
```

---

## Summary

- asynq = Redis-backed job queue with priorities, retries, scheduling, and deduplication
- **Enqueue**: `client.EnqueueContext(ctx, task, options...)` — returns task ID immediately
- **Workers**: register handler per task type via `mux.HandleFunc`; server distributes work
- **Priority queues**: configure `Queues: map[string]int{"critical": 6, "default": 1}` — higher weight = more workers assigned
- **Scheduler**: cron syntax for recurring tasks; `asynq.Unique` prevents duplicate scheduled tasks
- **Idempotency**: use `asynq.Unique(ttl)` to deduplicate; implement idempotent handlers
- **Monitoring**: `asynq.Inspector` API + `asynqmon` web UI

## Exercises

### Easy
1. Create an asynq worker that handles a `pdf:generate` task. The payload is `{report_id, user_id}`. The handler writes a fake PDF to `/tmp/report_<id>.pdf`. Verify that it runs and the file exists.
2. Enqueue 50 `email:welcome` tasks and observe them being processed by 5 concurrent workers. Print a progress counter.
3. Schedule a `cache:warmup` task to run every 5 minutes using the asynq Scheduler. Verify it fires by logging the current time in the handler.

### Medium
4. Implement **graceful shutdown**: the asynq server should stop accepting new tasks and wait for in-flight tasks to complete when it receives SIGTERM. Test by sending SIGTERM while a long-running task is active.
5. Add a **dead letter queue handler**: after 5 retries, tasks in the `emails` queue are moved to a `emails:failed` queue. Write a separate process that reads from `emails:failed`, logs the failed tasks, and sends an alert.
6. Implement **task chaining**: when a `user:signup` task completes, it enqueues a `email:welcome` task and a `analytics:track_signup` task. The chain only proceeds if the previous step succeeds.

### Hard
7. Build a **workflow engine** using asynq: a workflow has steps, each with a handler. Steps run sequentially. If a step fails, the workflow retries that step. Store workflow state in Redis. Implement `workflow.Run(ctx, steps)` and `workflow.Resume(ctx, workflowID)`.
8. Implement **rate-limited task processing**: the `email:send` handler should process at most 100 tasks per second across all workers. Use a Redis rate limiter (token bucket) inside the handler middleware. Workers that exceed the rate should back off and retry.
