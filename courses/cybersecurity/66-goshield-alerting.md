# Chapter 66: GoShield — Alerting Engine

*Detections are worthless if nobody sees them. GoShield's alerting engine routes findings to the right people via the right channels — fast enough to matter.*

---

## Alerting Design

```
Detection Engine → Severity Classification → Routing → Delivery

CRITICAL → Slack + PagerDuty + Email (immediate, 24/7 on-call)
HIGH     → Slack + Email (within 30 minutes)
MEDIUM   → Slack channel (review during business hours)
LOW      → Dashboard only (batch review)
INFO     → Log only (no alert)
```

---

## Alert Structure

```go
package alerting

import (
    "time"
    "encoding/json"
)

type Severity string

const (
    Critical Severity = "CRITICAL"
    High     Severity = "HIGH"
    Medium   Severity = "MEDIUM"
    Low      Severity = "LOW"
    Info     Severity = "INFO"
)

type Alert struct {
    ID          string          `json:"id"`
    Time        time.Time       `json:"time"`
    Severity    Severity        `json:"severity"`
    RuleName    string          `json:"rule_name"`
    Description string          `json:"description"`
    Host        string          `json:"host"`
    ProcessName string          `json:"process_name,omitempty"`
    PID         int32           `json:"pid,omitempty"`
    FilePath    string          `json:"file_path,omitempty"`
    NetworkDst  string          `json:"network_dst,omitempty"`
    Evidence    []string        `json:"evidence"`
    ATTACKTechnique string      `json:"attack_technique,omitempty"`  // e.g., T1055 Process Injection
    Suppressed  bool            `json:"suppressed"`
}

func (a *Alert) JSON() string {
    b, _ := json.MarshalIndent(a, "", "  ")
    return string(b)
}
```

---

## Slack Alerter

```go
package alerting

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type SlackAlerter struct {
    webhookURL string
    channel    string
    client     *http.Client
}

func NewSlackAlerter(webhookURL, channel string) *SlackAlerter {
    return &SlackAlerter{
        webhookURL: webhookURL,
        channel:    channel,
        client:     &http.Client{Timeout: 10 * time.Second},
    }
}

func (s *SlackAlerter) Send(alert *Alert) error {
    color := severityColor(alert.Severity)
    
    // Slack Block Kit message
    payload := map[string]interface{}{
        "channel": s.channel,
        "attachments": []map[string]interface{}{
            {
                "color": color,
                "blocks": []map[string]interface{}{
                    {
                        "type": "header",
                        "text": map[string]interface{}{
                            "type": "plain_text",
                            "text": fmt.Sprintf("[%s] %s", alert.Severity, alert.RuleName),
                        },
                    },
                    {
                        "type": "section",
                        "fields": []map[string]interface{}{
                            {"type": "mrkdwn", "text": fmt.Sprintf("*Host:*\n%s", alert.Host)},
                            {"type": "mrkdwn", "text": fmt.Sprintf("*Time:*\n%s", alert.Time.Format("2006-01-02 15:04:05 UTC"))},
                            {"type": "mrkdwn", "text": fmt.Sprintf("*Description:*\n%s", alert.Description)},
                            {"type": "mrkdwn", "text": fmt.Sprintf("*ATT&CK:*\n%s", alert.ATTACKTechnique)},
                        },
                    },
                },
            },
        },
    }
    
    if len(alert.Evidence) > 0 {
        evidenceText := ""
        for _, e := range alert.Evidence {
            evidenceText += "• " + e + "\n"
        }
        payload["attachments"].([]map[string]interface{})[0]["blocks"] = 
            append(payload["attachments"].([]map[string]interface{})[0]["blocks"].([]map[string]interface{}),
            map[string]interface{}{
                "type": "section",
                "text": map[string]interface{}{
                    "type": "mrkdwn",
                    "text": fmt.Sprintf("*Evidence:*\n%s", evidenceText),
                },
            })
    }
    
    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("slack webhook failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("slack returned status %d", resp.StatusCode)
    }
    
    return nil
}

func severityColor(s Severity) string {
    switch s {
    case Critical:
        return "#FF0000"
    case High:
        return "#FF6600"
    case Medium:
        return "#FFAA00"
    case Low:
        return "#00AA00"
    default:
        return "#888888"
    }
}
```

---

## Email Alerter

```go
package alerting

import (
    "fmt"
    "net/smtp"
    "strings"
)

type EmailAlerter struct {
    host     string
    port     int
    username string
    password string
    from     string
    to       []string
}

func NewEmailAlerter(host string, port int, user, pass, from string, to []string) *EmailAlerter {
    return &EmailAlerter{host, port, user, pass, from, to}
}

func (e *EmailAlerter) Send(alert *Alert) error {
    subject := fmt.Sprintf("[GoShield %s] %s on %s", alert.Severity, alert.RuleName, alert.Host)
    
    body := fmt.Sprintf(`GoShield Security Alert

Severity: %s
Rule: %s
Host: %s
Time: %s

Description: %s

ATT&CK Technique: %s

Evidence:
%s

--
GoShield EDR | Do not reply
`, alert.Severity, alert.RuleName, alert.Host,
        alert.Time.Format("2006-01-02 15:04:05 UTC"),
        alert.Description,
        alert.ATTACKTechnique,
        strings.Join(alert.Evidence, "\n"))
    
    msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain\r\n\r\n%s",
        e.from, strings.Join(e.to, ","), subject, body)
    
    auth := smtp.PlainAuth("", e.username, e.password, e.host)
    addr := fmt.Sprintf("%s:%d", e.host, e.port)
    
    return smtp.SendMail(addr, auth, e.from, e.to, []byte(msg))
}
```

---

## Alert Router — Routing by Severity

```go
package alerting

import (
    "fmt"
    "log"
    "sync"
    "time"
)

type Alerter interface {
    Send(alert *Alert) error
}

type AlertRouter struct {
    alerters       map[Severity][]Alerter
    dedupWindow    time.Duration
    recentAlerts   map[string]time.Time
    mu             sync.Mutex
    suppressionRules []SuppressionRule
}

type SuppressionRule struct {
    RuleName string
    Duration time.Duration
    Reason   string
}

func NewAlertRouter(dedupWindow time.Duration) *AlertRouter {
    return &AlertRouter{
        alerters:     make(map[Severity][]Alerter),
        dedupWindow:  dedupWindow,
        recentAlerts: make(map[string]time.Time),
    }
}

func (r *AlertRouter) Register(severity Severity, alerter Alerter) {
    r.alerters[severity] = append(r.alerters[severity], alerter)
}

// Suppress a specific rule for duration (for maintenance windows)
func (r *AlertRouter) Suppress(ruleName string, d time.Duration, reason string) {
    r.suppressionRules = append(r.suppressionRules, SuppressionRule{
        RuleName: ruleName,
        Duration: d,
        Reason:   reason,
    })
    log.Printf("Alert suppressed: %s for %s (%s)", ruleName, d, reason)
}

func (r *AlertRouter) Route(alert *Alert) {
    // Check suppression
    for _, rule := range r.suppressionRules {
        if rule.RuleName == alert.RuleName {
            alert.Suppressed = true
            log.Printf("Alert suppressed: %s", alert.RuleName)
            return
        }
    }
    
    // Deduplicate: same host+rule within window → skip
    dedupKey := alert.Host + ":" + alert.RuleName
    r.mu.Lock()
    if lastSeen, ok := r.recentAlerts[dedupKey]; ok {
        if time.Since(lastSeen) < r.dedupWindow {
            r.mu.Unlock()
            return
        }
    }
    r.recentAlerts[dedupKey] = time.Now()
    r.mu.Unlock()
    
    // Route to registered alerters for this severity
    alerters := r.alerters[alert.Severity]
    
    for _, alerter := range alerters {
        go func(a Alerter) {
            if err := a.Send(alert); err != nil {
                log.Printf("Alert delivery failed: %v", err)
            }
        }(alerter)
    }
    
    fmt.Printf("[%s] %s: %s\n", alert.Severity, alert.RuleName, alert.Description)
}
```

---

## Webhook Alerter (Generic)

```go
package alerting

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
)

type WebhookAlerter struct {
    url     string
    headers map[string]string
    client  *http.Client
}

func NewWebhookAlerter(url string, headers map[string]string) *WebhookAlerter {
    return &WebhookAlerter{
        url:     url,
        headers: headers,
        client:  &http.Client{Timeout: 10 * time.Second},
    }
}

func (w *WebhookAlerter) Send(alert *Alert) error {
    body, err := json.Marshal(alert)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("POST", w.url, bytes.NewReader(body))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    for k, v := range w.headers {
        req.Header.Set(k, v)
    }
    
    resp, err := w.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}
```

---

## Wiring it All Together

```go
package main

import (
    "time"
    "your/goshield/alerting"
)

func setupAlerting() *alerting.AlertRouter {
    router := alerting.NewAlertRouter(5 * time.Minute)
    
    slack := alerting.NewSlackAlerter(
        "https://hooks.slack.com/services/YOUR/WEBHOOK",
        "#security-alerts",
    )
    
    email := alerting.NewEmailAlerter(
        "smtp.company.com", 587,
        "goshield@company.com", "password",
        "goshield@company.com",
        []string{"security-team@company.com"},
    )
    
    // CRITICAL and HIGH go to Slack + Email
    router.Register(alerting.Critical, slack)
    router.Register(alerting.Critical, email)
    router.Register(alerting.High, slack)
    router.Register(alerting.High, email)
    
    // MEDIUM to Slack only
    router.Register(alerting.Medium, slack)
    
    return router
}
```

---

## Summary

| Component | Purpose |
|-----------|---------|
| Alert struct | Standardized format for all findings |
| Slack alerter | Real-time team notification |
| Email alerter | Formal record, off-hours notification |
| Alert router | Severity-based routing, deduplication |
| Suppression | Mute alerts during maintenance |

---

## Exercises

1. Set up a Slack incoming webhook (free) and test the Go Slack alerter
2. Add deduplication that tracks alert counts instead of just suppressing repeats
3. Implement a PagerDuty alerter (REST API) for CRITICAL alerts
4. Build an alert suppression API endpoint: `POST /suppress?rule=ProcessInjection&duration=1h`
