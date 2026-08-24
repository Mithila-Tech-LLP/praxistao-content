# Chapter 126: Alerting — Rules, Thresholds, PagerDuty, and On-Call

Metrics and dashboards tell you what happened after you look. Alerting tells you when something needs attention right now. A well-designed alerting pipeline wakes the right person, with the right context, at the right time — and stays quiet when nothing needs doing. A poorly designed one trains engineers to ignore pages.

## Table of Contents

1. [The Alerting Pipeline](#1-the-alerting-pipeline)
2. [Writing Prometheus Alert Rules](#2-writing-prometheus-alert-rules)
3. [Alert Anatomy](#3-alert-anatomy)
4. [Alertmanager Configuration](#4-alertmanager-configuration)
5. [SLO-Based Alerting](#5-slo-based-alerting)
6. [PagerDuty Integration](#6-pagerduty-integration)
7. [On-Call Runbooks](#7-on-call-runbooks)
8. [Alert Fatigue](#8-alert-fatigue)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Alerting Pipeline

Prometheus does not send notifications itself. It evaluates rules and hands off to Alertmanager, which handles routing, deduplication, silencing, and delivery.

```
+----------------+       scrape        +------------------+
|  Your Go App   | ----------------->  |    Prometheus    |
| /metrics       |   (every 15s)       |                  |
+----------------+                     |  evaluates alert |
                                        |  rules every 15s |
                                        +--------+---------+
                                                 |
                                          alert fires
                                          (pending -> firing)
                                                 |
                                                 v
                                        +------------------+
                                        |  Alertmanager    |
                                        |                  |
                                        |  - groups alerts |
                                        |  - deduplicates  |
                                        |  - inhibits      |
                                        |  - silences      |
                                        +--------+---------+
                                                 |
                              +------------------+------------------+
                              |                                     |
                              v                                     v
                     +----------------+                   +------------------+
                     |   PagerDuty    |                   |      Slack       |
                     |  (critical)    |                   |    (warning)     |
                     +----------------+                   +------------------+
```

The key stages:

1. **Scrape** — Prometheus pulls `/metrics` from each target on a fixed interval (default 15s).
2. **Evaluate** — Prometheus runs each alert rule against the scraped data on the same interval.
3. **Pending** — When a rule expression is true, the alert enters `pending` state. It stays there for the `for` duration.
4. **Firing** — After `for` elapses, the alert becomes `firing` and Prometheus sends it to Alertmanager.
5. **Route** — Alertmanager matches the alert against routes using label selectors and sends to the configured receiver.
6. **Deliver** — The receiver (PagerDuty, Slack, email) gets the notification.

---

## 2. Writing Prometheus Alert Rules

Alert rules live in YAML files and are loaded by Prometheus via the `rule_files` config key.

```yaml
# alerts.yml
groups:
  - name: api_alerts
    interval: 15s   # how often to evaluate this group (defaults to Prometheus global)
    rules:

      # Alert 1: Error rate above 1% for 5 minutes
      - alert: HighErrorRate
        expr: |
          rate(http_requests_total{status=~"5.."}[5m])
          /
          rate(http_requests_total[5m])
          > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate on {{ $labels.instance }}"
          description: >
            Error rate is {{ $value | humanizePercentage }} on
            {{ $labels.instance }} (job={{ $labels.job }}).
          runbook_url: "https://wiki.example.com/runbooks/high-error-rate"

      # Alert 2: 99th percentile latency above 500ms for 10 minutes
      - alert: HighLatency
        expr: |
          histogram_quantile(
            0.99,
            rate(http_request_duration_seconds_bucket[5m])
          ) > 0.5
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High p99 latency on {{ $labels.instance }}"
          description: >
            p99 latency is {{ $value | humanizeDuration }} on
            {{ $labels.instance }}. SLO threshold is 500ms.
          runbook_url: "https://wiki.example.com/runbooks/high-latency"

      # Alert 3: Service unreachable for 1 minute
      - alert: ServiceDown
        expr: up{job="api"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.instance }} is down"
          description: >
            Prometheus cannot scrape {{ $labels.instance }}.
            The instance has been unreachable for at least 1 minute.
          runbook_url: "https://wiki.example.com/runbooks/service-down"

      # Alert 4: Pod restarted more than 5 times in the past hour
      - alert: PodCrashLooping
        expr: |
          kube_pod_container_status_restarts_total{namespace="production"}
          - kube_pod_container_status_restarts_total{namespace="production"} offset 1h
          > 5
        for: 0m   # fire immediately — crash loops are always urgent
        labels:
          severity: critical
        annotations:
          summary: "Pod {{ $labels.pod }} is crash looping"
          description: >
            Container {{ $labels.container }} in pod {{ $labels.pod }}
            has restarted {{ $value }} times in the last hour.
          runbook_url: "https://wiki.example.com/runbooks/pod-crash-loop"
```

Wire this file into your `prometheus.yml`:

```yaml
# prometheus.yml (excerpt)
rule_files:
  - "alerts.yml"
  - "slo_alerts.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets: ["alertmanager:9093"]
```

---

## 3. Alert Anatomy

Every Prometheus alert has five parts. Understanding each one prevents the most common alerting mistakes.

```yaml
- alert: HighErrorRate          # (1) alert name
  expr: |                       # (2) PromQL expression
    rate(http_requests_total{status=~"5.."}[5m])
    / rate(http_requests_total[5m]) > 0.01
  for: 5m                       # (3) pending duration
  labels:                       # (4) labels (used for routing)
    severity: critical
    env: production
    team: backend
  annotations:                  # (5) human-readable context
    summary: "High error rate on {{ $labels.instance }}"
    description: >
      {{ $labels.job }} is returning {{ $value | humanizePercentage }}
      errors on {{ $labels.instance }}.
    runbook_url: "https://wiki.example.com/runbooks/high-error-rate"
```

**Alert name** — used by Alertmanager for grouping and deduplication. Alerts with the same name and same labels are treated as one incident. Name it after what is broken, not what threshold was crossed: `HighErrorRate` not `ErrorRateOver1Percent`.

**expr** — any valid PromQL expression. When the expression returns a non-empty result set, the alert is active. The series labels from the expression become available as `$labels` in annotations.

**for** — the pending duration. The expression must be true continuously for this long before the alert fires. This is the most important safeguard against flapping. A 1-second spike in error rate is not an incident; a 5-minute sustained elevation is. Set `for: 0m` only when immediate firing is genuinely warranted (crash loops, complete service outage).

**labels** — key-value pairs attached to the alert. Alertmanager uses these for routing (`severity: critical` → PagerDuty) and for grouping related alerts into a single notification. Labels added here merge with labels from the PromQL expression.

**annotations** — free-form metadata for the responder. Unlike labels, annotations are not used for routing. Use `{{ $labels.instance }}` to reference the metric labels and `{{ $value }}` for the threshold value that triggered the alert. Useful template functions:
- `{{ $value | humanize }}` — 1234567 becomes "1.23M"
- `{{ $value | humanizePercentage }}` — 0.015 becomes "1.5%"
- `{{ $value | humanizeDuration }}` — 0.75 becomes "750ms"

A well-annotated alert gives the on-call engineer enough context to start investigating without opening a dashboard first.

---

## 4. Alertmanager Configuration

Alertmanager receives alerts from Prometheus, groups them, applies routing rules, and sends notifications. It handles the operational complexity so Prometheus does not have to.

```yaml
# alertmanager.yml
global:
  # Default SMTP settings (optional, only if using email)
  smtp_smarthost: "smtp.example.com:587"
  smtp_from: "alertmanager@example.com"

  # If an alert stops being sent by Prometheus and has no explicit
  # end time, Alertmanager marks it resolved after this long
  resolve_timeout: 5m

route:
  # Group alerts that share these labels into one notification batch.
  # This prevents 50 separate pages when a whole cluster goes down.
  group_by: [alertname, env]

  # Wait this long after the first alert in a group before sending.
  # Gives time for related alerts to arrive and be grouped together.
  group_wait: 30s

  # After the first notification, wait this long before sending
  # updates for the same group (e.g., more alerts joined the group).
  group_interval: 5m

  # After a notification, wait this long before re-notifying
  # if the alert is still firing and unresolved.
  repeat_interval: 4h

  # Default receiver if no child route matches
  receiver: slack-receiver

  # Child routes are evaluated in order; first match wins
  routes:
    - matchers:
        - severity = "critical"
      receiver: pagerduty-receiver
      # Critical alerts repeat more frequently
      repeat_interval: 1h
      # Stop here once this route matches (the default).
      # Set continue: true to also evaluate later sibling routes.
      continue: false

    - matchers:
        - severity = "warning"
      receiver: slack-receiver

receivers:
  - name: slack-receiver
    slack_configs:
      - api_url: "https://hooks.slack.com/services/T00/B00/XXXX"
        channel: "#alerts-warning"
        send_resolved: true
        title: >-
          [{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}]
          {{ .CommonLabels.alertname }}
        text: >-
          {{ range .Alerts }}
          *Alert:* {{ .Annotations.summary }}
          *Description:* {{ .Annotations.description }}
          *Severity:* {{ .Labels.severity }}
          *Runbook:* {{ .Annotations.runbook_url }}
          {{ end }}

  - name: pagerduty-receiver
    pagerduty_configs:
      - routing_key: "your-pagerduty-integration-key-here"
        send_resolved: true
        severity: '{{ if eq .CommonLabels.severity "critical" }}critical{{ else }}warning{{ end }}'
        description: '{{ .CommonAnnotations.summary }}'
        details:
          runbook: '{{ .CommonAnnotations.runbook_url }}'
          env: '{{ .CommonLabels.env }}'

# Inhibition rules: suppress lower-priority alerts when a
# higher-priority alert is already firing for the same service.
inhibit_rules:
  # If a critical alert is firing for a given alertname,
  # suppress any warning alert with the same alertname.
  - source_matchers:
      - severity = "critical"
    target_matchers:
      - severity = "warning"
    # Both source and target must share these labels
    # for the inhibition to apply.
    equal: [alertname, env]

  # If ServiceDown is firing, suppress HighErrorRate and HighLatency
  # for the same instance — the root cause is already paging.
  - source_matchers:
      - alertname = "ServiceDown"
    target_matchers:
      - alertname =~ "HighErrorRate|HighLatency"
    equal: [instance]
```

The inhibition rule is critical for noise reduction. When a whole service is down, you do not want separate pages for high error rate and high latency — those are symptoms of the same root cause. Inhibit the symptoms when the root cause alert is already firing.

---

## 5. SLO-Based Alerting

Threshold alerts have a fundamental problem: they treat all violations as equal. A 2% error rate for 30 seconds is not the same problem as a 0.1% error rate sustained for 6 hours. The first might be a deploy; the second is eating through your error budget and will breach your SLO by end of day.

Error budget burn rate gives you a single number that captures both magnitude and duration. A burn rate of 1 means you are consuming your error budget at exactly the rate that matches your SLO period. A burn rate of 14.4 means you will exhaust your entire monthly error budget in about 2 days (30 days / 14.4 ≈ 50 hours).

For a 99.9% availability SLO over 30 days, the error budget is 0.1% of all requests, or about 43.2 minutes of downtime.

```yaml
# slo_recording_rules.yml
# Recording rules compute expensive expressions once and store the result.
groups:
  - name: slo_recording_rules
    interval: 30s
    rules:
      # 5-minute error rate (used by fast-burn alert)
      - record: job:http_requests_errors:rate5m
        expr: rate(http_requests_total{status=~"5.."}[5m])

      - record: job:http_requests_total:rate5m
        expr: rate(http_requests_total[5m])

      # 1-hour error rate (used by slow-burn alert)
      - record: job:http_requests_errors:rate1h
        expr: rate(http_requests_total{status=~"5.."}[1h])

      - record: job:http_requests_total:rate1h
        expr: rate(http_requests_total[1h])

      # 6-hour error rate (used by slow-burn alert second window)
      - record: job:http_requests_errors:rate6h
        expr: rate(http_requests_total{status=~"5.."}[6h])

      - record: job:http_requests_total:rate6h
        expr: rate(http_requests_total[6h])
```

```yaml
# slo_alerts.yml
groups:
  - name: slo_burn_alerts
    rules:

      # Fast burn: burn_rate > 14.4 means you consume 2% of the
      # 30-day error budget in a single hour (14.4 x 1h / 720h = 2%).
      # Uses two windows to reduce false positives from brief spikes.
      - alert: SLOFastBurn
        expr: |
          (
            job:http_requests_errors:rate5m{job="api"}
            / job:http_requests_total:rate5m{job="api"}
          ) / 0.001 > 14.4
          and
          (
            job:http_requests_errors:rate1h{job="api"}
            / job:http_requests_total:rate1h{job="api"}
          ) / 0.001 > 14.4
        for: 2m
        labels:
          severity: critical
          slo: availability
        annotations:
          summary: "SLO fast burn: error budget exhausting rapidly"
          description: >
            Error budget burn rate is {{ $value | humanize }}x the SLO threshold.
            At this rate the 30-day error budget is exhausted in roughly
            30/burn-rate days (e.g. burn rate 14.4 = ~2 days).
          runbook_url: "https://wiki.example.com/runbooks/slo-burn"

      # Slow burn: burn_rate > 1 over 6h means you are on track
      # to exhaust the full monthly budget within 30 days.
      # This does not page — it notifies Slack for daytime awareness.
      - alert: SLOSlowBurn
        expr: |
          (
            job:http_requests_errors:rate1h{job="api"}
            / job:http_requests_total:rate1h{job="api"}
          ) / 0.001 > 3
          and
          (
            job:http_requests_errors:rate6h{job="api"}
            / job:http_requests_total:rate6h{job="api"}
          ) / 0.001 > 3
        for: 15m
        labels:
          severity: warning
          slo: availability
        annotations:
          summary: "SLO slow burn: error budget being consumed"
          description: >
            Error budget burn rate is {{ $value | humanize }}x over the last 6h.
            Investigate before this becomes a paging incident.
          runbook_url: "https://wiki.example.com/runbooks/slo-burn"
```

The `0.001` divisor is the SLO error rate (1 - 0.999). Dividing the current error rate by this gives the burn rate multiplier. The number `14.4` comes from deciding how much budget you are willing to lose before paging: a common choice is "page if 2% of the 30-day budget burns in 1 hour." Budget consumed = burn rate × window / period, so burn rate = 2% × 720h / 1h = 14.4.

Fast burn uses `severity: critical` and routes to PagerDuty. Slow burn uses `severity: warning` and routes to Slack. This means engineers sleep through slow burns but get woken for fast ones — which is exactly right.

---

## 6. PagerDuty Integration

PagerDuty connects Alertmanager alerts to on-call schedules and escalation policies.

**Getting the routing key:**

In PagerDuty, create a new service or use an existing one. Under "Integrations", add the "Prometheus" integration. Copy the "Integration Key" — this is your `routing_key` in Alertmanager.

```yaml
# alertmanager.yml (pagerduty receiver)
receivers:
  - name: pagerduty-critical
    pagerduty_configs:
      - routing_key: "a1b2c3d4e5f6..."   # from PagerDuty integration
        send_resolved: true               # closes the PagerDuty incident automatically
        severity: critical
        description: '{{ .CommonAnnotations.summary }}'
        client: "Alertmanager"
        client_url: "http://alertmanager.example.com"
        details:
          alertname: '{{ .CommonLabels.alertname }}'
          env: '{{ .CommonLabels.env }}'
          runbook: '{{ .CommonAnnotations.runbook_url }}'
```

**Severity mapping to PagerDuty priority:**

| Alertmanager severity | PagerDuty severity | PagerDuty priority |
|-----------------------|--------------------|--------------------|
| critical              | critical           | P1                 |
| warning               | warning            | P2                 |
| info                  | info               | P3                 |

`send_resolved: true` is important. Without it, Alertmanager sends a notification when an alert fires but sends nothing when it resolves. PagerDuty incidents stay open until manually acknowledged. With `send_resolved: true`, PagerDuty automatically resolves the incident when the alert clears.

**Escalation policies in PagerDuty:**

Configure in PagerDuty UI (not in Alertmanager):

```
Escalation Policy: Backend On-Call
  Level 1: Primary on-call (notify immediately via push + SMS)
            Escalate after: 15 minutes if unacknowledged
  Level 2: Secondary on-call (notify via push + SMS)
            Escalate after: 15 minutes if unacknowledged
  Level 3: Engineering Manager (notify via phone call)
            Escalate after: 30 minutes if unacknowledged
```

Alertmanager's `repeat_interval` and PagerDuty's escalation policy serve different purposes. `repeat_interval` controls how often Alertmanager re-sends the same alert if it is still firing. PagerDuty's escalation policy controls who gets notified when the incident is not acknowledged.

---

## 7. On-Call Runbooks

Every alert that pages a human must have a runbook. A runbook is a step-by-step guide that tells the on-call engineer what to do when they receive the alert at 3am without context.

The `runbook_url` annotation in the alert points to the runbook document (Confluence, Notion, GitHub wiki, internal docs).

**Runbook template:**

```markdown
# Runbook: HighErrorRate

## What is this alert?
The HTTP error rate (5xx responses) has exceeded 1% for more than 5 minutes
on the API service.

## What does it mean?
Clients are receiving server errors. This directly impacts user experience
and is consuming the error budget for the availability SLO.

## Immediate steps
1. Check Grafana dashboard: https://grafana.example.com/d/api-overview
2. Check recent deployments: `kubectl rollout history deployment/api -n production`
3. Check pod logs: `kubectl logs -l app=api -n production --tail=100`
4. Check database connectivity:
   `kubectl exec -it <pod> -n production -- pg_isready -h db.example.com`
5. If a bad deploy is identified, rollback:
   `kubectl rollout undo deployment/api -n production`

## Escalation steps
If you cannot identify the root cause within 15 minutes:
- Escalate to the database team if logs show connection errors
- Escalate to the platform team if pods are crash looping
- Page the engineering manager if revenue impact is visible in dashboards

## Related alerts
- ServiceDown — if this fires, HighErrorRate is a symptom; focus on ServiceDown
- HighLatency — may indicate a slow dependency causing timeouts

## Post-incident
File an incident report within 24 hours:
https://wiki.example.com/incidents/new
```

A runbook that says "investigate and fix" is useless. A runbook should let someone who has never seen this service before take the right first steps in the first five minutes.

---

## 8. Alert Fatigue

Alert fatigue happens when engineers receive so many alerts that they stop treating them as signals. Symptoms:
- Pages are acknowledged and silenced without investigation
- "We get this alert all the time, it usually resolves itself"
- The Alertmanager silence list grows without corresponding fixes
- On-call handoff includes "just ignore the X alert"

Every alert that pages an engineer implies: "this requires a human to take action right now." If that is not true, the alert should not page.

**Audit your silences:**

Alertmanager's API exposes active and expired silences. A long list of recurring silences for the same alert is a sign the alert should be deleted or downgraded, not silenced.

```bash
# List all active silences
curl -s http://alertmanager:9093/api/v2/silences | jq '.[] | select(.status.state == "active") | {comment: .comment, createdBy: .createdBy, matchers: .matchers}'

# Count silences per alertname to find the most-silenced alerts
curl -s http://alertmanager:9093/api/v2/silences | \
  jq -r '.[].matchers[] | select(.name == "alertname") | .value' | \
  sort | uniq -c | sort -rn
```

If an alert shows up in the silence list more than twice in a month, that alert is not actionable. Either:
1. Fix the underlying condition so the alert stops firing.
2. Raise the threshold or extend the `for` duration so it only fires on real problems.
3. Downgrade severity from critical to warning.
4. Delete the alert entirely.

**Toil reduction:**

"Toil" is operational work that is manual, repetitive, and does not permanently fix the underlying problem. Silencing the same alert every Monday morning because batch jobs produce elevated error rates is toil. The fix is a time-based inhibition rule that suppresses the alert during the batch window, not a recurring manual silence.

**Monthly alert review checklist:**

```
For each alert that fired in the past month:
  [ ] Was action taken when it fired? (check incident tickets)
  [ ] Was it silenced without a ticket? (sign of noise)
  [ ] Did it fire more than 10 times? (might need threshold adjustment)
  [ ] Does the runbook cover the actual steps taken? (update if not)
  [ ] Could it be inhibited by a more fundamental alert? (reduce noise)
```

The goal is not zero alerts — it is that every alert that fires results in a ticket or a fixed runbook. When that is true, alert fatigue goes away.

---

## Summary

The alerting pipeline runs Prometheus → Alertmanager → notification channels. Prometheus evaluates rules and fires alerts after the `for` duration; Alertmanager routes, deduplicates, inhibits, and delivers.

Good alert rules have a clear `for` duration to prevent flapping, severity labels for routing, and annotations with templated context so the responder knows what broke and where to start.

Alertmanager routes alerts by label matchers. Critical alerts go to PagerDuty and wake people up. Warning alerts go to Slack and are handled during business hours. Inhibit rules suppress symptom alerts when root-cause alerts are already firing.

SLO-based alerting with error budget burn rates is more reliable than raw thresholds. Fast burn alerts (burn_rate > 14.4 over a 5m + 1h window) catch acute incidents. Slow burn alerts (burn_rate > 3 over 1h + 6h windows) catch gradual degradation before it becomes an incident.

Every paging alert must have a runbook. Every page must require human action. Alerts that do not meet both criteria cause alert fatigue, which causes real incidents to be missed.

---

## Exercises

### Easy

Write a complete Prometheus alert rule for `ServiceDown` that fires when `up{job="api"} == 0` for 1 minute. Include correct severity labels and an annotation for `summary` that includes the instance name using template syntax. Include a `runbook_url` annotation.

### Medium

Configure a complete `alertmanager.yml` that:
1. Routes alerts with `severity: critical` to a PagerDuty receiver using a routing key of your choice.
2. Routes alerts with `severity: warning` to a Slack receiver with a webhook URL of your choice.
3. Adds an inhibit rule so that if a critical alert is firing for a given `alertname` and `env`, any warning alert with the same `alertname` and `env` is suppressed.
4. Sets `group_by: [alertname, env]`, `group_wait: 30s`, `group_interval: 5m`, and `repeat_interval: 4h` on the root route.

### Hard

Design a full SLO alerting setup for an API with a 99.9% availability SLO (30-day window):

1. Write recording rules that compute the 5-minute, 1-hour, and 6-hour error rates and request rates for `job="api"`.
2. Using those recording rules, write a fast-burn alert that fires when the error budget burn rate exceeds 14.4 over both a 5-minute and a 1-hour window, with `severity: critical`.
3. Write a slow-burn alert that fires when the burn rate exceeds 3 over both a 1-hour and a 6-hour window, with `severity: warning`.
4. Add the appropriate Alertmanager routing so the fast-burn alert pages PagerDuty immediately and the slow-burn alert sends a Slack message to `#slo-watch`.
5. Explain in comments why you use two windows for each burn-rate alert instead of one.
