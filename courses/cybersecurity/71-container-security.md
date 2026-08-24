# Chapter 71: Container Security — Docker and Kubernetes

*Containers changed how software is deployed. They also introduced new attack surfaces: container escapes, privileged containers, exposed APIs, and misconfigured Kubernetes RBAC.*

---

## Container Security Model

```
Host Machine
├── Linux Kernel         ← shared by all containers
├── Docker Daemon        ← attack surface (if exposed)
└── Containers
    ├── Namespaces       ← isolation (PID, NET, MNT, UTS, IPC)
    ├── cgroups          ← resource limits
    └── seccomp/AppArmor ← syscall/capability restrictions
```

**Key insight:** Containers share the host kernel. A kernel vulnerability exploitable from inside a container compromises the host.

---

## Docker Security

### Container Escape Vulnerabilities

```bash
# Dangerous: running privileged container
docker run --privileged -it ubuntu bash
# Inside: can access ALL host devices, mount host filesystem
# Escape: mount host root filesystem
mount /dev/sda1 /mnt/host
chroot /mnt/host

# Dangerous: mounting Docker socket
docker run -v /var/run/docker.sock:/var/run/docker.sock -it ubuntu bash
# Inside: can create new privileged containers on HOST
docker run --privileged -v /:/host -it ubuntu bash -c "chroot /host"

# Dangerous: running as root (default)
# If attacker compromises app running as root in container,
# they get root inside container and can attempt escapes
```

### Secure Dockerfile

```dockerfile
# BAD
FROM ubuntu:latest
RUN apt-get install -y python3-pip
COPY . /app
CMD ["python3", "/app/app.py"]

# GOOD
FROM python:3.11-slim AS builder
WORKDIR /build
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

FROM python:3.11-slim
# Don't run as root
RUN useradd -r -s /bin/false appuser
WORKDIR /app
COPY --from=builder /usr/local/lib/python3.11/site-packages/ /usr/local/lib/python3.11/site-packages/
COPY --chown=appuser:appuser . .
USER appuser
# Read-only root filesystem
# Set in docker run: --read-only

EXPOSE 8080
CMD ["python3", "app.py"]
```

### Container Hardening

```bash
# Run with security options
docker run \
    --user 1000:1000 \              # non-root
    --read-only \                   # read-only filesystem
    --tmpfs /tmp \                  # writable temp
    --no-new-privileges \           # can't escalate privileges
    --cap-drop=ALL \                # drop all capabilities
    --cap-add=NET_BIND_SERVICE \    # add only what's needed
    --security-opt seccomp=profile.json \  # custom seccomp
    -e DB_PASSWORD="$(cat /run/secrets/db_pass)" \  # secrets from file
    my-app

# Check Docker Security with Trivy
trivy image my-app:latest        # scan for CVEs
trivy config Dockerfile          # check dockerfile misconfigs
trivy config k8s-deployment.yaml # check k8s configs
```

---

## Kubernetes Security

### RBAC (Role-Based Access Control)

```yaml
# BAD: wildcard permissions
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: bad-role
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
# GOOD: specific read-only
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: production
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
```

### Network Policies (Default: all traffic allowed!)

```yaml
# Default deny all ingress
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: production
spec:
  podSelector: {}  # applies to all pods
  policyTypes:
  - Ingress
  - Egress
---
# Allow only frontend → backend
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - port: 8080
```

### Pod Security (Pod Security Standards)

```yaml
# namespace-level security enforcement
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    # Enforce restricted security profile on all pods in this namespace
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## Kubernetes Attack Scenarios

### Compromised Pod → Cluster Admin

```bash
# Inside compromised pod, if service account has too much RBAC:

# Check service account token
cat /var/run/secrets/kubernetes.io/serviceaccount/token

# Check what permissions we have
curl -k https://kubernetes.default.svc/api/v1/namespaces \
    -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"

# If we can create pods: spawn privileged pod
kubectl run escape --image=ubuntu --privileged=true --restart=Never \
    -- bash -c "nsenter -t 1 -m -u -i -n -p -- bash"
# This enters host namespaces = full host access
```

### Container Escape via CVE

```
Common container escape CVEs:
- CVE-2022-0847 (DirtyPipe) — Linux kernel, writable /proc/self/mem
- CVE-2019-5736 (runc) — overwrite runc binary from container
- CVE-2020-15257 (containerd-shim) — abstract socket exposure
```

---

## Go: Container Security Scanner

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type ContainerFinding struct {
    Severity string
    Issue    string
    Detail   string
}

func checkDockerRunning() bool {
    cmd := exec.Command("docker", "info")
    return cmd.Run() == nil
}

func checkPrivilegedContainers() []ContainerFinding {
    var findings []ContainerFinding
    
    cmd := exec.Command("docker", "inspect", "--format",
        "{{.Name}} {{.HostConfig.Privileged}}", 
        "$(docker ps -q)")
    
    out, err := cmd.Output()
    if err != nil {
        return findings
    }
    
    for _, line := range strings.Split(string(out), "\n") {
        if strings.Contains(line, "true") {
            parts := strings.Fields(line)
            if len(parts) >= 1 {
                findings = append(findings, ContainerFinding{
                    Severity: "CRITICAL",
                    Issue:    "Privileged container",
                    Detail:   parts[0],
                })
            }
        }
    }
    return findings
}

func checkDockerSocket() []ContainerFinding {
    var findings []ContainerFinding
    
    out, _ := exec.Command("docker", "ps", "--format", "{{.ID}}").Output()
    containers := strings.Split(strings.TrimSpace(string(out)), "\n")
    
    for _, containerID := range containers {
        if containerID == "" {
            continue
        }
        
        out, _ := exec.Command("docker", "inspect", "--format",
            "{{range .Mounts}}{{.Source}} {{end}}", containerID).Output()
        
        if strings.Contains(string(out), "/var/run/docker.sock") {
            findings = append(findings, ContainerFinding{
                Severity: "CRITICAL",
                Issue:    "Docker socket mounted inside container",
                Detail:   containerID,
            })
        }
    }
    return findings
}

func checkRunningAsRoot() []ContainerFinding {
    var findings []ContainerFinding
    
    out, _ := exec.Command("docker", "ps", "--format", "{{.ID}}").Output()
    containers := strings.Split(strings.TrimSpace(string(out)), "\n")
    
    for _, containerID := range containers {
        if containerID == "" {
            continue
        }
        
        out, _ := exec.Command("docker", "exec", containerID, "id", "-u").Output()
        uid := strings.TrimSpace(string(out))
        
        if uid == "0" {
            nameOut, _ := exec.Command("docker", "inspect", "--format",
                "{{.Name}}", containerID).Output()
            findings = append(findings, ContainerFinding{
                Severity: "HIGH",
                Issue:    "Container running as root (UID 0)",
                Detail:   strings.TrimSpace(string(nameOut)),
            })
        }
    }
    return findings
}

func checkSensitiveEnvVars() []ContainerFinding {
    var findings []ContainerFinding
    
    out, _ := exec.Command("docker", "ps", "--format", "{{.ID}}").Output()
    containers := strings.Split(strings.TrimSpace(string(out)), "\n")
    
    sensitiveVars := []string{"PASSWORD", "SECRET", "KEY", "TOKEN", "CREDENTIAL"}
    
    for _, containerID := range containers {
        if containerID == "" {
            continue
        }
        
        envOut, _ := exec.Command("docker", "inspect", "--format",
            "{{range .Config.Env}}{{.}}\n{{end}}", containerID).Output()
        
        for _, line := range strings.Split(string(envOut), "\n") {
            for _, sensitive := range sensitiveVars {
                if strings.Contains(strings.ToUpper(line), sensitive) {
                    findings = append(findings, ContainerFinding{
                        Severity: "HIGH",
                        Issue:    "Sensitive data in environment variable",
                        Detail:   containerID + ": " + strings.Split(line, "=")[0],
                    })
                    break
                }
            }
        }
    }
    return findings
}

func main() {
    if !checkDockerRunning() {
        fmt.Println("Docker not running or not accessible")
        os.Exit(1)
    }
    
    var allFindings []ContainerFinding
    allFindings = append(allFindings, checkPrivilegedContainers()...)
    allFindings = append(allFindings, checkDockerSocket()...)
    allFindings = append(allFindings, checkRunningAsRoot()...)
    allFindings = append(allFindings, checkSensitiveEnvVars()...)
    
    if len(allFindings) == 0 {
        fmt.Println("No obvious container security issues found")
        return
    }
    
    for _, f := range allFindings {
        fmt.Printf("[%s] %s: %s\n", f.Severity, f.Issue, f.Detail)
    }
}
```

---

## Summary

| Risk | Attack | Defense |
|------|--------|---------|
| Privileged container | Mount host FS | Never use --privileged |
| Docker socket mount | Create privileged containers | Don't mount socket |
| Root in container | Escape to host via kernel vuln | Non-root user |
| RBAC wildcard | Cluster admin from compromised pod | Least privilege RBAC |
| No network policy | Lateral movement | Default-deny NetworkPolicy |

---

## Exercises

1. Run `trivy image nginx:latest` — how many CVEs does a fresh nginx image have?
2. Build a hardened Docker image for a simple Go web server — pass trivy with zero HIGH/CRITICAL CVEs
3. Set up a Kubernetes cluster (minikube) and implement default-deny NetworkPolicy
4. Attempt a container escape in your lab using the Docker socket mount technique
