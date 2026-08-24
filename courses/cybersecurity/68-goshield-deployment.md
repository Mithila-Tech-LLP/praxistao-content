# Chapter 68: GoShield — Deployment and Packaging

*Great software that can't be deployed easily won't be deployed at all. This chapter covers packaging GoShield for real enterprise deployment: Docker, systemd services, cross-compilation, and automated deployment.*

---

## Deployment Architecture

```
Production Deployment:

┌─────────────────────┐      ┌─────────────────────────┐
│   Agent Hosts       │      │   GoShield Server        │
│                     │      │                          │
│  goshield-agent     │─────→│  Nginx (reverse proxy)  │
│  (systemd service)  │      │    ↓                     │
│                     │      │  goshield-server         │
│  Runs as:           │      │  (Docker container)      │
│    root (needs      │      │    ↓                     │
│    syscall access)  │      │  SQLite database         │
│                     │      │    (persistent volume)   │
└─────────────────────┘      └─────────────────────────┘
         ↑
  1000s of agents
```

---

## Building Release Binaries

Go's cross-compilation makes it trivial to build binaries for any target:

```bash
# Makefile for GoShield
# goshield/Makefile

VERSION := 1.0.0
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

# Build agent for all platforms
build-agents:
	# Linux x86_64 (most servers)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) \
		-o dist/goshield-agent-linux-amd64 ./cmd/agent
	
	# Linux ARM64 (AWS Graviton, Raspberry Pi 4)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) \
		-o dist/goshield-agent-linux-arm64 ./cmd/agent
	
	# Windows x64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) \
		-o dist/goshield-agent-windows-amd64.exe ./cmd/agent
	
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) \
		-o dist/goshield-agent-darwin-amd64 ./cmd/agent
	
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) \
		-o dist/goshield-agent-darwin-arm64 ./cmd/agent

# Build server
build-server:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) \
		-o dist/goshield-server-linux-amd64 ./cmd/server

# Static binary (no libc dependency — runs anywhere)
build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) \
		-a -installsuffix cgo \
		-o dist/goshield-agent-static-linux-amd64 ./cmd/agent

# SHA256 checksums for integrity verification
checksums:
	cd dist && sha256sum * > SHA256SUMS
```

**Note:** GoShield agent uses `gopsutil` which needs CGO for some features. For fully static builds, use the pure-Go alternatives.

---

## Docker — Containerizing the Server

```dockerfile
# server/Dockerfile

# ─── Build stage ───────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-s -w" \
    -o goshield-server ./cmd/server

# ─── Runtime stage ─────────────────────────────────────
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    tzdata

# Non-root user for security
RUN addgroup -S goshield && adduser -S goshield -G goshield

WORKDIR /app

# Copy binary and default configs
COPY --from=builder /app/goshield-server .
COPY configs/server-config.yaml .

# Data directory for SQLite
RUN mkdir -p /app/data && chown goshield:goshield /app/data

USER goshield

EXPOSE 8080

ENV GOSHIELD_DB_PATH=/app/data/goshield.db
ENV GOSHIELD_CONFIG=/app/server-config.yaml

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./goshield-server"]
```

### Docker Compose — Full Stack

```yaml
# docker-compose.yml

version: '3.8'

services:
  goshield-server:
    build:
      context: .
      dockerfile: server/Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - goshield-data:/app/data
      - ./configs/server-config.yaml:/app/server-config.yaml:ro
      - ./rules:/app/rules:ro
    environment:
      - GOSHIELD_DB_PATH=/app/data/goshield.db
      - GOSHIELD_API_KEY=${GOSHIELD_API_KEY}
      - GOSHIELD_SLACK_WEBHOOK=${SLACK_WEBHOOK_URL}
    restart: unless-stopped
    networks:
      - goshield-net
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3

  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - goshield-server
    networks:
      - goshield-net

volumes:
  goshield-data:
    driver: local

networks:
  goshield-net:
    driver: bridge
```

```bash
# Deploy
export GOSHIELD_API_KEY=$(openssl rand -hex 32)
docker-compose up -d

# View logs
docker-compose logs -f goshield-server

# Update
docker-compose pull && docker-compose up -d
```

### Nginx Reverse Proxy Config

```nginx
# nginx/nginx.conf
server {
    listen 443 ssl http2;
    server_name goshield.yourcompany.com;

    ssl_certificate     /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # Rate limiting — prevent agent flooding
    limit_req_zone $binary_remote_addr zone=agents:10m rate=10r/s;
    
    location /api/agent/ {
        limit_req zone=agents burst=20 nodelay;
        proxy_pass http://goshield-server:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    location / {
        proxy_pass http://goshield-server:8080;
    }
}
```

---

## Systemd Service — Agent on Linux Hosts

```ini
# /etc/systemd/system/goshield-agent.service

[Unit]
Description=GoShield Security Agent
Documentation=https://github.com/yourorg/goshield
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/local/bin/goshield-agent
Restart=always
RestartSec=5s

# Environment
Environment=GOSHIELD_SERVER=https://goshield.yourcompany.com
Environment=GOSHIELD_API_KEY=YOUR_API_KEY_HERE
Environment=GOSHIELD_LOG_LEVEL=info

# Security hardening (keep what the agent needs)
NoNewPrivileges=yes
ProtectHome=no         # agent needs to monitor /home
ProtectSystem=no       # agent needs to monitor /etc

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=goshield-agent

[Install]
WantedBy=multi-user.target
```

```bash
# Install and start
sudo cp goshield-agent-linux-amd64 /usr/local/bin/goshield-agent
sudo chmod 755 /usr/local/bin/goshield-agent
sudo cp goshield-agent.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable goshield-agent
sudo systemctl start goshield-agent

# Check status
sudo systemctl status goshield-agent
sudo journalctl -u goshield-agent -f
```

---

## Mass Deployment with Ansible

For deploying to hundreds of servers:

```yaml
# deploy-agents.yml

- name: Deploy GoShield Agent
  hosts: all
  become: yes
  
  vars:
    goshield_version: "1.0.0"
    goshield_server: "https://goshield.yourcompany.com"
    goshield_api_key: "{{ vault_goshield_api_key }}"  # from Ansible Vault
  
  tasks:
    - name: Copy agent binary
      copy:
        src: "dist/goshield-agent-linux-{{ ansible_architecture | replace('x86_64', 'amd64') }}"
        dest: /usr/local/bin/goshield-agent
        mode: '0755'
        owner: root
        group: root
    
    - name: Create systemd service
      template:
        src: goshield-agent.service.j2
        dest: /etc/systemd/system/goshield-agent.service
      notify: Restart goshield-agent
    
    - name: Enable and start service
      systemd:
        name: goshield-agent
        enabled: yes
        state: started
        daemon_reload: yes
  
  handlers:
    - name: Restart goshield-agent
      systemd:
        name: goshield-agent
        state: restarted
```

```bash
# Deploy to all servers in inventory
ansible-playbook -i inventory/production deploy-agents.yml \
    --ask-vault-pass

# Check deployment
ansible all -m command -a "systemctl is-active goshield-agent"
```

---

## Agent Auto-Registration

When a new agent comes online, it should auto-register with the server:

```go
// cmd/agent/main.go — agent registration
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "runtime"
    "time"
    
    "github.com/shirou/gopsutil/v3/host"
)

type RegistrationRequest struct {
    AgentID  string `json:"agent_id"`
    Hostname string `json:"hostname"`
    OS       string `json:"os"`
    Arch     string `json:"arch"`
    Version  string `json:"version"`
}

var Version = "1.0.0"

func registerAgent(serverURL, apiKey string) error {
    hostInfo, _ := host.Info()
    
    hostname, _ := os.Hostname()
    
    req := RegistrationRequest{
        AgentID:  generateAgentID(hostname),  // deterministic from hostname
        Hostname: hostname,
        OS:       hostInfo.Platform + " " + hostInfo.PlatformVersion,
        Arch:     runtime.GOARCH,
        Version:  Version,
    }
    
    body, _ := json.Marshal(req)
    
    httpReq, err := http.NewRequest("POST", 
        serverURL+"/api/agent/register",
        bytes.NewReader(body))
    if err != nil {
        return err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("X-API-Key", apiKey)
    
    resp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        return fmt.Errorf("registration failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("registration rejected: %d", resp.StatusCode)
    }
    
    return nil
}
```

---

## Agent Health Check and Heartbeat

```go
// Periodic heartbeat so server knows the agent is alive
func sendHeartbeat(serverURL, apiKey, agentID string) error {
    payload := map[string]interface{}{
        "agent_id":  agentID,
        "timestamp": time.Now().Unix(),
        "status":    "running",
    }
    
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", serverURL+"/api/agent/heartbeat", bytes.NewReader(body))
    req.Header.Set("X-API-Key", apiKey)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}

// In agent main loop:
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            if err := sendHeartbeat(serverURL, apiKey, agentID); err != nil {
                log.Printf("Heartbeat failed: %v", err)
            }
        case <-ctx.Done():
            return
        }
    }
}()
```

---

## Server-side: Dead Agent Detection

```go
// Server: mark agents offline if no heartbeat in 2 minutes
func (s *Server) detectDeadAgents() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        cutoff := time.Now().Add(-2 * time.Minute)
        
        _, err := s.db.Exec(`
            UPDATE agents SET status = 'offline' 
            WHERE last_heartbeat < ? AND status = 'online'
        `, cutoff.Unix())
        
        if err != nil {
            log.Printf("Dead agent cleanup error: %v", err)
            continue
        }
        
        // Alert if critical agents go offline
        var offline []string
        rows, _ := s.db.Query(`
            SELECT hostname FROM agents 
            WHERE status = 'offline' AND critical = 1
        `)
        for rows.Next() {
            var h string
            rows.Scan(&h)
            offline = append(offline, h)
        }
        rows.Close()
        
        if len(offline) > 0 {
            s.alerter.SendAlert(fmt.Sprintf(
                "⚠️ Critical agents offline: %s", 
                strings.Join(offline, ", ")))
        }
    }
}
```

---

## Production Checklist

```
Agent:
  ✓ TLS verification enabled (don't use InsecureSkipVerify!)
  ✓ API key stored in environment variable, not config file
  ✓ Runs as root but with systemd security hardening
  ✓ Automatic restart on failure (Restart=always)
  ✓ Logs to systemd journal
  ✓ Binary is signed and checksum verified on deploy

Server:
  ✓ Running behind Nginx (TLS termination)
  ✓ TLS 1.2+ only, strong cipher suites
  ✓ Rate limiting on agent endpoints
  ✓ SQLite with WAL mode + regular backups
  ✓ API keys rotated regularly
  ✓ Alerting configured (Slack/email)
  ✓ Health endpoint for monitoring

Network:
  ✓ Agents communicate outbound only to GOSHIELD_SERVER
  ✓ Server accessible only from agent network + SOC
  ✓ Dashboard accessible only from SOC network
```

---

## Quick Start Script

```bash
#!/bin/bash
# install-goshield.sh — one-liner agent installer

set -euo pipefail

GOSHIELD_SERVER="${GOSHIELD_SERVER:-https://goshield.yourcompany.com}"
GOSHIELD_API_KEY="${GOSHIELD_API_KEY:-}"
ARCH=$(uname -m | sed 's/x86_64/amd64/')
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ -z "$GOSHIELD_API_KEY" ]; then
    echo "Error: GOSHIELD_API_KEY not set"
    exit 1
fi

echo "[*] Downloading GoShield agent..."
curl -sSL "$GOSHIELD_SERVER/download/agent/$OS/$ARCH" \
    -o /usr/local/bin/goshield-agent
chmod 755 /usr/local/bin/goshield-agent

echo "[*] Installing systemd service..."
cat > /etc/systemd/system/goshield-agent.service << EOF
[Unit]
Description=GoShield Security Agent
After=network.target

[Service]
ExecStart=/usr/local/bin/goshield-agent
Restart=always
Environment=GOSHIELD_SERVER=$GOSHIELD_SERVER
Environment=GOSHIELD_API_KEY=$GOSHIELD_API_KEY

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable goshield-agent
systemctl start goshield-agent
systemctl is-active goshield-agent && echo "[+] GoShield agent running"
```

```bash
# Deploy to any Linux server in one line:
curl -sSL https://goshield.yourcompany.com/install.sh | \
    GOSHIELD_API_KEY=your-key sudo bash
```

---

## Summary

| Component | Deployment method |
|-----------|-----------------|
| Server | Docker + Docker Compose |
| TLS | Nginx reverse proxy |
| Agent (single) | Binary + systemd service |
| Agent (mass) | Ansible playbook |
| Monitoring | Health endpoint + heartbeat |
| Updates | `docker-compose pull` + `ansible-playbook` |

GoShield is now production-ready: agents auto-register, send heartbeats, and recover from crashes; the server is containerized and sits behind a TLS-terminating proxy; Ansible deploys to hundreds of hosts in minutes.

---

## Exercises

1. Write a Dockerfile for the GoShield agent and build it. Verify the resulting image size.
2. Deploy the GoShield server using Docker Compose. Verify the health endpoint.
3. Write the Ansible playbook for deploying agents. Test against a local VM.
4. Add an endpoint `/download/agent/{os}/{arch}` to the GoShield server that serves the appropriate agent binary
5. Implement agent config updates: server can push new rules to agents without restarting them
