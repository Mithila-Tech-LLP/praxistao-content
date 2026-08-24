# Chapter 67: GoShield Dashboard — REST API and Web Interface

*A security tool without a UI is half a tool. This chapter builds the GoShield web dashboard that security analysts use to view alerts, hunt through events, and respond to incidents.*

---

## Dashboard Goals

The GoShield dashboard provides:
1. **Alert feed** — real-time stream of security alerts with severity
2. **Event hunt** — search across process/file/network events
3. **Agent status** — which endpoints are online
4. **Statistics** — events per hour, top processes, top connections
5. **Alert management** — acknowledge, resolve, add notes

---

## Alerting Module First

```go
// pkg/alerting/alerter.go
package alerting

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/yourname/goshield/pkg/events"
)

// WebhookConfig configures a webhook alert destination
type WebhookConfig struct {
    Name         string
    URL          string
    Severities   []string  // Only send alerts at these severity levels
}

// Alerter dispatches alerts to configured destinations
type Alerter struct {
    webhooks []WebhookConfig
    client   *http.Client
}

// NewAlerter creates an alerter
func NewAlerter() *Alerter {
    return &Alerter{
        client: &http.Client{Timeout: 10 * time.Second},
    }
}

// AddWebhook registers a webhook
func (a *Alerter) AddWebhook(url string, severities []string) {
    a.webhooks = append(a.webhooks, WebhookConfig{
        URL:       url,
        Severities: severities,
    })
}

// Send dispatches an alert to all configured destinations
func (a *Alerter) Send(alert *events.Alert) {
    for _, wh := range a.webhooks {
        if !a.shouldSend(alert, wh.Severities) {
            continue
        }
        go a.sendWebhook(wh, alert)
    }
}

func (a *Alerter) shouldSend(alert *events.Alert, severities []string) bool {
    if len(severities) == 0 {
        return true
    }
    for _, s := range severities {
        if string(alert.Severity) == s {
            return true
        }
    }
    return false
}

// SlackPayload formats an alert as a Slack message
type SlackPayload struct {
    Text        string            `json:"text"`
    Username    string            `json:"username"`
    IconEmoji   string            `json:"icon_emoji"`
    Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
    Color  string `json:"color"`
    Title  string `json:"title"`
    Text   string `json:"text"`
    Footer string `json:"footer"`
    Ts     int64  `json:"ts"`
}

func (a *Alerter) sendWebhook(wh WebhookConfig, alert *events.Alert) {
    severityColors := map[events.Severity]string{
        events.SeverityInfo:     "#36a64f",
        events.SeverityLow:      "#2196F3",
        events.SeverityMedium:   "#FF9800",
        events.SeverityHigh:     "#FF5722",
        events.SeverityCritical: "#F44336",
    }
    
    color, ok := severityColors[alert.Severity]
    if !ok {
        color = "#9E9E9E"
    }
    
    payload := SlackPayload{
        Username:  "GoShield EDR",
        IconEmoji: ":shield:",
        Attachments: []SlackAttachment{
            {
                Color: color,
                Title: fmt.Sprintf("[%s] %s", strings.ToUpper(string(alert.Severity)), alert.RuleName),
                Text: fmt.Sprintf("*Host:* %s\n*Rule:* %s\n*Description:* %s\n*MITRE:* %s",
                    alert.Hostname, alert.RuleID, alert.Description, alert.MITRE),
                Footer: "GoShield EDR",
                Ts:     alert.Timestamp.Unix(),
            },
        },
    }
    
    data, err := json.Marshal(payload)
    if err != nil {
        return
    }
    
    resp, err := a.client.Post(wh.URL, "application/json", bytes.NewReader(data))
    if err != nil {
        fmt.Printf("Webhook error (%s): %v\n", wh.URL, err)
        return
    }
    resp.Body.Close()
}
```

---

## Web Dashboard — HTML + JavaScript

The dashboard is a single-page app served as static files:

```html
<!-- web/index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoShield EDR Dashboard</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            background: #0d1117;
            color: #c9d1d9;
            min-height: 100vh;
        }
        
        .navbar {
            background: #161b22;
            border-bottom: 1px solid #30363d;
            padding: 12px 24px;
            display: flex;
            align-items: center;
            gap: 16px;
        }
        
        .navbar h1 {
            font-size: 18px;
            color: #58a6ff;
            font-weight: 600;
        }
        
        .nav-link {
            color: #8b949e;
            text-decoration: none;
            padding: 6px 12px;
            border-radius: 6px;
            cursor: pointer;
        }
        
        .nav-link.active, .nav-link:hover {
            color: #c9d1d9;
            background: #21262d;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 24px;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 16px;
            margin-bottom: 24px;
        }
        
        .stat-card {
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 8px;
            padding: 20px;
        }
        
        .stat-card .label {
            font-size: 12px;
            color: #8b949e;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .stat-card .value {
            font-size: 28px;
            font-weight: 600;
            color: #c9d1d9;
            margin-top: 8px;
        }
        
        .stat-card.critical .value { color: #f85149; }
        .stat-card.high .value { color: #ff7b72; }
        .stat-card.online .value { color: #56d364; }
        
        .panel {
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 8px;
            margin-bottom: 24px;
        }
        
        .panel-header {
            padding: 16px 20px;
            border-bottom: 1px solid #30363d;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .panel-header h2 {
            font-size: 14px;
            font-weight: 600;
            color: #c9d1d9;
        }
        
        .badge {
            padding: 2px 8px;
            border-radius: 10px;
            font-size: 11px;
            font-weight: 600;
        }
        
        .badge-critical { background: #3d1c1c; color: #f85149; border: 1px solid #f85149; }
        .badge-high     { background: #3d2c1c; color: #ff7b72; border: 1px solid #ff7b72; }
        .badge-medium   { background: #3d341c; color: #e3b341; border: 1px solid #e3b341; }
        .badge-low      { background: #1c2a3d; color: #58a6ff; border: 1px solid #58a6ff; }
        .badge-info     { background: #1d2d1d; color: #56d364; border: 1px solid #56d364; }
        
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 13px;
        }
        
        th {
            padding: 10px 16px;
            text-align: left;
            font-size: 11px;
            font-weight: 600;
            color: #8b949e;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            border-bottom: 1px solid #30363d;
        }
        
        td {
            padding: 12px 16px;
            border-bottom: 1px solid #21262d;
            vertical-align: middle;
        }
        
        tr:last-child td { border-bottom: none; }
        
        tr:hover td { background: #1c2128; }
        
        .code {
            font-family: 'SFMono-Regular', Consolas, monospace;
            font-size: 12px;
            color: #79c0ff;
            background: #0d1117;
            padding: 2px 6px;
            border-radius: 4px;
            max-width: 400px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            display: inline-block;
        }
        
        .search-bar {
            display: flex;
            gap: 12px;
            padding: 16px 20px;
            border-bottom: 1px solid #30363d;
        }
        
        .search-bar input, .search-bar select {
            background: #0d1117;
            border: 1px solid #30363d;
            color: #c9d1d9;
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 13px;
            outline: none;
        }
        
        .search-bar input:focus, .search-bar select:focus {
            border-color: #58a6ff;
        }
        
        .search-bar input { flex: 1; }
        
        button {
            background: #21262d;
            border: 1px solid #30363d;
            color: #c9d1d9;
            padding: 8px 16px;
            border-radius: 6px;
            font-size: 13px;
            cursor: pointer;
        }
        
        button:hover { background: #30363d; }
        button.primary { background: #1f6feb; border-color: #1f6feb; color: white; }
        button.primary:hover { background: #388bfd; }
        
        .online-dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 6px; }
        .online-dot.online { background: #56d364; }
        .online-dot.offline { background: #6e7681; }
        
        #page-alerts, #page-hunt, #page-agents { display: none; }
        #page-alerts.active, #page-hunt.active, #page-agents.active { display: block; }
        
        .loading { text-align: center; padding: 40px; color: #8b949e; }
        
        .refresh-indicator { font-size: 11px; color: #8b949e; }
        
        .resolve-btn {
            background: none;
            border: 1px solid #56d364;
            color: #56d364;
            padding: 4px 10px;
            font-size: 11px;
            cursor: pointer;
            border-radius: 4px;
        }
        .resolve-btn:hover { background: #1d2d1d; }
    </style>
</head>
<body>

<nav class="navbar">
    <h1>🛡 GoShield EDR</h1>
    <a class="nav-link active" onclick="showPage('alerts')">Alerts</a>
    <a class="nav-link" onclick="showPage('hunt')">Hunt</a>
    <a class="nav-link" onclick="showPage('agents')">Agents</a>
    <span style="margin-left: auto" class="refresh-indicator" id="refresh-time"></span>
</nav>

<div class="container">

    <!-- Stats Row -->
    <div class="stats-grid" id="stats-grid">
        <div class="stat-card critical">
            <div class="label">Critical Alerts</div>
            <div class="value" id="stat-critical">—</div>
        </div>
        <div class="stat-card high">
            <div class="label">High Alerts</div>
            <div class="value" id="stat-high">—</div>
        </div>
        <div class="stat-card online">
            <div class="label">Online Agents</div>
            <div class="value" id="stat-agents">—</div>
        </div>
        <div class="stat-card">
            <div class="label">Events (24h)</div>
            <div class="value" id="stat-events">—</div>
        </div>
    </div>

    <!-- Alerts Page -->
    <div id="page-alerts" class="active">
        <div class="panel">
            <div class="panel-header">
                <h2>Security Alerts</h2>
                <div style="display:flex;gap:8px">
                    <select id="severity-filter" onchange="loadAlerts()">
                        <option value="">All Severities</option>
                        <option value="critical">Critical</option>
                        <option value="high">High</option>
                        <option value="medium">Medium</option>
                        <option value="low">Low</option>
                    </select>
                    <button onclick="loadAlerts()">Refresh</button>
                </div>
            </div>
            <div id="alerts-content"><div class="loading">Loading...</div></div>
        </div>
    </div>

    <!-- Hunt Page -->
    <div id="page-hunt">
        <div class="panel">
            <div class="panel-header"><h2>Event Hunt — Process Events</h2></div>
            <div class="search-bar">
                <input id="hunt-name" type="text" placeholder="Process name (e.g. powershell)">
                <input id="hunt-cmd" type="text" placeholder="Command line contains...">
                <input id="hunt-agent" type="text" placeholder="Agent/hostname">
                <button class="primary" onclick="hunt()">Search</button>
            </div>
            <div id="hunt-results"><div class="loading">Enter search criteria above</div></div>
        </div>
    </div>

    <!-- Agents Page -->
    <div id="page-agents">
        <div class="panel">
            <div class="panel-header"><h2>Connected Agents</h2></div>
            <div id="agents-content"><div class="loading">Loading...</div></div>
        </div>
    </div>

</div>

<script>
const API_KEY = 'changeme-insecure-default'; // In production, use proper auth
const BASE = '';

function apiHeaders() {
    return { 'X-API-Key': API_KEY, 'Content-Type': 'application/json' };
}

async function loadStats() {
    const resp = await fetch(`${BASE}/api/v1/stats`, { headers: apiHeaders() });
    if (!resp.ok) return;
    const data = await resp.json();
    
    const bySev = data.alerts_by_severity || {};
    document.getElementById('stat-critical').textContent = bySev.critical || 0;
    document.getElementById('stat-high').textContent = bySev.high || 0;
    document.getElementById('stat-agents').textContent = data.online_agents || 0;
    
    const total = (data.process_events_24h || 0) + (data.file_events_24h || 0) + (data.network_events_24h || 0);
    document.getElementById('stat-events').textContent = total.toLocaleString();
    
    document.getElementById('refresh-time').textContent = 'Last updated: ' + new Date().toLocaleTimeString();
}

function severityBadge(sev) {
    return `<span class="badge badge-${sev}">${sev.toUpperCase()}</span>`;
}

async function loadAlerts() {
    const sev = document.getElementById('severity-filter').value;
    const url = `${BASE}/api/v1/alerts${sev ? '?severity=' + sev : ''}`;
    const resp = await fetch(url, { headers: apiHeaders() });
    const data = await resp.json();
    
    const alerts = data.alerts || [];
    
    if (alerts.length === 0) {
        document.getElementById('alerts-content').innerHTML = '<div class="loading">No active alerts</div>';
        return;
    }
    
    let html = `<table>
        <thead><tr>
            <th>Time</th><th>Severity</th><th>Rule</th><th>Host</th><th>Description</th><th></th>
        </tr></thead>
        <tbody>`;
    
    for (const a of alerts) {
        const time = new Date(a.timestamp).toLocaleString();
        html += `<tr>
            <td>${time}</td>
            <td>${severityBadge(a.severity)}</td>
            <td><span class="code">${a.rule_id}</span> ${a.rule_name}</td>
            <td>${a.hostname}</td>
            <td title="${a.description}">${a.description.substring(0, 80)}${a.description.length > 80 ? '...' : ''}</td>
            <td><button class="resolve-btn" onclick="resolveAlert('${a.id}')">Resolve</button></td>
        </tr>`;
    }
    
    html += '</tbody></table>';
    document.getElementById('alerts-content').innerHTML = html;
}

async function resolveAlert(id) {
    await fetch(`${BASE}/api/v1/alerts/${id}/resolve`, {
        method: 'POST',
        headers: apiHeaders(),
        body: JSON.stringify({ notes: 'Resolved from dashboard' })
    });
    loadAlerts();
}

async function hunt() {
    const name = document.getElementById('hunt-name').value;
    const cmd = document.getElementById('hunt-cmd').value;
    const agent = document.getElementById('hunt-agent').value;
    
    let url = `${BASE}/api/v1/events/process?limit=200`;
    if (name) url += '&name=' + encodeURIComponent(name);
    if (cmd) url += '&cmd=' + encodeURIComponent(cmd);
    if (agent) url += '&agent=' + encodeURIComponent(agent);
    
    document.getElementById('hunt-results').innerHTML = '<div class="loading">Searching...</div>';
    
    const resp = await fetch(url, { headers: apiHeaders() });
    const data = await resp.json();
    const evts = data.events || [];
    
    if (evts.length === 0) {
        document.getElementById('hunt-results').innerHTML = '<div class="loading">No results found</div>';
        return;
    }
    
    let html = `<div style="padding:12px 20px;font-size:12px;color:#8b949e">${evts.length} result(s)</div>
        <table>
        <thead><tr>
            <th>Time</th><th>Host</th><th>PID</th><th>User</th><th>Process</th><th>Command Line</th>
        </tr></thead><tbody>`;
    
    for (const e of evts) {
        const time = new Date(e.timestamp).toLocaleString();
        const cmd_truncated = (e.command_line || '').substring(0, 80);
        html += `<tr>
            <td style="font-size:11px;white-space:nowrap">${time}</td>
            <td>${e.hostname}</td>
            <td>${e.pid}</td>
            <td>${e.username || '—'}</td>
            <td><strong>${e.name}</strong></td>
            <td><span class="code" title="${e.command_line}">${cmd_truncated}</span></td>
        </tr>`;
    }
    
    html += '</tbody></table>';
    document.getElementById('hunt-results').innerHTML = html;
}

async function loadAgents() {
    const resp = await fetch(`${BASE}/api/v1/agents`, { headers: apiHeaders() });
    const data = await resp.json();
    const agents = data.agents || [];
    
    if (agents.length === 0) {
        document.getElementById('agents-content').innerHTML = '<div class="loading">No agents registered</div>';
        return;
    }
    
    let html = `<table>
        <thead><tr>
            <th>Status</th><th>Hostname</th><th>OS</th><th>Agent ID</th><th>Last Seen</th><th>Events Total</th>
        </tr></thead><tbody>`;
    
    for (const a of agents) {
        const online = a.online;
        const status = `<span class="online-dot ${online ? 'online' : 'offline'}"></span>${online ? 'Online' : 'Offline'}`;
        const lastSeen = new Date(a.last_seen).toLocaleString();
        html += `<tr>
            <td>${status}</td>
            <td><strong>${a.hostname}</strong></td>
            <td>${a.os || '—'}</td>
            <td><span class="code">${a.id.substring(0, 16)}...</span></td>
            <td style="font-size:12px">${lastSeen}</td>
            <td>${(a.events_total || 0).toLocaleString()}</td>
        </tr>`;
    }
    
    html += '</tbody></table>';
    document.getElementById('agents-content').innerHTML = html;
}

function showPage(page) {
    document.querySelectorAll('.nav-link').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('[id^="page-"]').forEach(el => el.classList.remove('active'));
    
    event.target.classList.add('active');
    document.getElementById('page-' + page).classList.add('active');
    
    if (page === 'alerts') loadAlerts();
    if (page === 'agents') loadAgents();
}

// Initial load
loadStats();
loadAlerts();

// Auto-refresh every 30 seconds
setInterval(() => {
    loadStats();
    loadAlerts();
}, 30000);
</script>

</body>
</html>
```

---

## API Endpoint for Agents List

Add to `api.go`:
```go
func (s *Server) handleGetAgents(w http.ResponseWriter, r *http.Request) {
    rows, err := s.db.conn.Query(
        `SELECT id, hostname, os, version, last_seen, online, events_total 
         FROM agents ORDER BY last_seen DESC`)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    var agents []events.AgentStatus
    for rows.Next() {
        a := events.AgentStatus{}
        rows.Scan(&a.AgentID, &a.Hostname, &a.OS, &a.Version, 
                  &a.LastSeen, &a.Online, &a.EventsTotal)
        agents = append(agents, a)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "agents": agents,
        "count": len(agents),
    })
}
```

---

## Running GoShield

```bash
# Start server
GOSHIELD_API_KEY=mysecret go run cmd/server/main.go

# Open dashboard
open http://localhost:8080

# Deploy agent on a server
GOSHIELD_SERVER=http://your-server:8080 GOSHIELD_API_KEY=mysecret ./goshield-agent
```

---

## Summary

GoShield now has:
- File integrity monitoring (Chapter 60)
- Process monitoring (Chapter 61)
- Network monitoring (Chapter 62)
- Central server with SQLite storage (Chapter 64)
- Rule-based detection engine (Chapter 65)
- Slack/webhook alerting (Chapter 66)
- Web dashboard with hunt capabilities (this chapter)

This is a functional enterprise EDR system.

---

## Exercises

1. Run the complete GoShield stack and trigger an alert by running `powershell.exe -EncodedCommand` or creating a `.php` file in `/var/www`
2. Add a "timeline" view to the dashboard that shows events for a specific agent over the last hour
3. Add real-time updates to the alerts page using Server-Sent Events (SSE)
4. Add a "MITRE ATT&CK" view that maps alerts to ATT&CK techniques and shows coverage
5. Build a Grafana dashboard instead of the custom HTML — connect it to the GoShield API
