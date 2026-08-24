# Chapter 131: Secrets Management — Vault, AWS Secrets Manager, and Doppler

Secrets are credentials: database passwords, API keys, TLS private keys, OAuth client secrets. Every application needs them. The question is where they live and who can read them. This chapter covers the full spectrum — from the dangerous patterns most codebases start with, through the 12-factor stepping stone, to production-grade solutions with HashiCorp Vault, AWS Secrets Manager, and Doppler.

## Table of Contents

1. [The Problem — What Can Go Wrong](#1-the-problem--what-can-go-wrong)
2. [12-Factor Approach — Environment Variables](#2-12-factor-approach--environment-variables)
3. [HashiCorp Vault](#3-hashicorp-vault)
4. [AWS Secrets Manager](#4-aws-secrets-manager)
5. [Kubernetes External Secrets Operator](#5-kubernetes-external-secrets-operator)
6. [Doppler — Simpler Alternative](#6-doppler--simpler-alternative)
7. [Secrets Rotation](#7-secrets-rotation)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Problem — What Can Go Wrong

### Hardcoded secrets

The most common mistake:

```go
// DO NOT DO THIS
const dbPassword = "super_secret_password_123"
const stripeAPIKey = "sk_live_abc123xyz789"

func connectDB() *sql.DB {
    dsn := fmt.Sprintf("postgres://app:%s@prod-db:5432/myapp", dbPassword)
    db, _ := sql.Open("pgx", dsn)
    return db
}
```

The problem is not just that the file contains a password. It is that git history is forever. Running `git log -p --all -S "super_secret_password_123"` finds it in every commit since it was added, even after deletion. When this repository is cloned — by a new developer, a CI runner, or an attacker who finds a public fork — the credential goes with it.

### .env files committed to git

The second most common mistake:

```bash
# .env -- accidentally committed
DATABASE_URL=postgres://app:real_password@prod-db:5432/myapp
STRIPE_SECRET_KEY=sk_live_abc123xyz789
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

Developers add `.env` to `.gitignore` but forget to remove a file that was already committed. One `git add .` and it is in the history.

### CI environment variables

CI systems like GitHub Actions store secrets in the repository settings — but many teams store non-secret config alongside actual secrets in the same place, visible to all collaborators. A disgruntled contractor with read access to the repository settings can see every secret the pipeline uses. Secrets printed in CI logs via a debug `fmt.Println` are stored in log archives for months.

### The actual risks

- Leaked database password: attacker dumps the entire users table, exfiltrates PII, sells it.
- Leaked API key: attacker calls your payment provider on your behalf, charges customers, or sends millions of spam emails through your email API.
- Leaked AWS secret key: attacker spins up EC2 instances to mine cryptocurrency on your account. AWS bill arrives for $50,000.
- No rotation: once leaked, a static credential stays valid until someone manually rotates it — which may take days if the leak is not discovered immediately.
- No audit trail: you cannot answer "who accessed this secret and when," which is a compliance failure under SOC 2 and GDPR.

---

## 2. 12-Factor Approach — Environment Variables

The [12-factor app](https://12factor.net/config) methodology says secrets belong in environment variables, not code.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }

    db, err := sql.Open("pgx", dbURL)
    if err != nil {
        log.Fatalf("open db: %v", err)
    }
    defer db.Close()

    fmt.Println("connected")
}
```

This is better than hardcoding: the value is not in source code, different environments can have different values, and the app does not need to know where the secret came from.

But environment variables alone are not enough:

| Problem | Environment Variables |
|---------|----------------------|
| Secret not in code | Yes, this is solved |
| Plaintext at rest | Still plaintext — visible in `ps aux`, `/proc/<pid>/environ`, container inspect |
| Rotation | Manual — requires redeploying all services that use the secret |
| Audit trail | None — no log of who read what or when |
| Per-service access control | None — any process on the host with the right permissions sees all env vars |
| Environment isolation | Easy to accidentally paste prod secrets into staging |

Environment variables are a stepping stone, not a destination. They solve the "secrets in git" problem but not rotation, audit, or access control.

---

## 3. HashiCorp Vault

Vault is the industry standard for secrets management. It handles storage, access control, auditing, and — its killer feature — dynamic secrets.

### Core concepts

```
Secrets Engines       Auth Methods          Policies
-----------------     -----------------     --------------------------
KV v2 (key-value)     AppRole               path "secret/data/myapp/*" {
Database (dynamic)    Kubernetes              capabilities = ["read"]
PKI (TLS certs)       AWS IAM               }
Transit (encryption)  Token
SSH
```

- **Secrets engines** are backends that store or generate secrets. KV v2 stores arbitrary key-value pairs. The database engine generates short-lived database credentials on demand.
- **Auth methods** verify the identity of clients. AppRole is common for services. Kubernetes auth lets pods authenticate using their service account JWT.
- **Policies** define what authenticated identities can do. Least-privilege: a service gets read access to its own secrets path and nothing else.
- **Leases** attach an expiry to secrets. Dynamic secrets are automatically revoked when the lease expires.

### KV v2 — storing static secrets

```bash
# Start a dev server (stores in memory, not for production)
vault server -dev

# In another terminal
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='root'  # dev mode root token printed on startup

# Store a secret
vault kv put secret/myapp/database \
  url="postgres://app:password@db:5432/myapp" \
  max_connections=25

# Read it back
vault kv get secret/myapp/database
# === Secret Path ===
# secret/data/myapp/database
#
# ======= Data =======
# Key              Value
# ---              -----
# max_connections  25
# url              postgres://app:password@db:5432/myapp

# Read just one field
vault kv get -field=url secret/myapp/database
```

### Dynamic secrets — Vault's killer feature

Static secrets have a fixed value. If leaked, they remain valid until rotated manually. Dynamic secrets are generated on demand with a short TTL. When the lease expires, Vault revokes the credential automatically. No static password exists anywhere.

Configure the database secrets engine:

```bash
# Enable the database secrets engine
vault secrets enable database

# Configure a PostgreSQL connection
vault write database/config/myapp-db \
    plugin_name=postgresql-database-plugin \
    allowed_roles="myapp-role" \
    connection_url="postgresql://vault:vaultpassword@db:5432/myapp?sslmode=disable" \
    username="vault" \
    password="vaultpassword"

# Create a role -- Vault runs this SQL to create short-lived credentials
vault write database/roles/myapp-role \
    db_name=myapp-db \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' \
        VALID UNTIL '{{expiration}}'; \
        GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
    default_ttl="1h" \
    max_ttl="24h"

# Any service with the right policy gets fresh credentials
vault read database/creds/myapp-role
# Key                Value
# lease_id           database/creds/myapp-role/abc123
# lease_duration     1h
# username           v-approle-xyz789
# password           A1b2C3d4-randomly-generated
```

When the lease expires, Vault drops the PostgreSQL user. A compromised credential is automatically invalidated within an hour.

### Go client

```bash
go get github.com/hashicorp/vault-client-go
```

```go
package vault

import (
    "context"
    "fmt"
    "net/http"
    "time"

    vault "github.com/hashicorp/vault-client-go"
    "github.com/hashicorp/vault-client-go/schema"
)

type Client struct {
    v *vault.Client
}

func NewClient(address, roleID, secretID string) (*Client, error) {
    client, err := vault.New(
        vault.WithAddress(address),
        vault.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
    )
    if err != nil {
        return nil, fmt.Errorf("vault.New: %w", err)
    }

    // Authenticate with AppRole
    resp, err := client.Auth.AppRoleLogin(
        context.Background(),
        schema.AppRoleLoginRequest{
            RoleId:   roleID,
            SecretId: secretID,
        },
    )
    if err != nil {
        return nil, fmt.Errorf("AppRoleLogin: %w", err)
    }
    if err := client.SetToken(resp.Auth.ClientToken); err != nil {
        return nil, fmt.Errorf("SetToken: %w", err)
    }

    return &Client{v: client}, nil
}

// GetSecret reads a KV v2 secret at the given path under the "secret" mount.
func (c *Client) GetSecret(ctx context.Context, path string) (map[string]interface{}, error) {
    resp, err := c.v.Secrets.KvV2Read(
        ctx,
        path,
        vault.WithMountPath("secret"),
    )
    if err != nil {
        return nil, fmt.Errorf("KvV2Read %s: %w", path, err)
    }
    return resp.Data.Data, nil
}
```

Usage:

```go
client, err := vault.NewClient(
    "http://vault:8200",
    os.Getenv("VAULT_ROLE_ID"),
    os.Getenv("VAULT_SECRET_ID"),
)
if err != nil {
    log.Fatal(err)
}

data, err := client.GetSecret(ctx, "myapp/database")
if err != nil {
    log.Fatal(err)
}
dbURL := data["url"].(string)
```

### AppRole authentication for services

AppRole uses two credentials: a RoleID (like a username, not secret — can be in config) and a SecretID (like a password — short-lived, injected by CI/Kubernetes).

```bash
# Create a policy
vault policy write myapp-policy - <<EOF
path "secret/data/myapp/*" {
  capabilities = ["read"]
}
path "database/creds/myapp-role" {
  capabilities = ["read"]
}
EOF

# Enable AppRole and create a role
vault auth enable approle
vault write auth/approle/role/myapp-role \
    token_policies="myapp-policy" \
    token_ttl=1h \
    token_max_ttl=4h \
    secret_id_ttl=24h

# RoleID: stable, can be in non-secret config
vault read auth/approle/role/myapp-role/role-id

# SecretID: short-lived, injected via CI at deploy time
vault write -f auth/approle/role/myapp-role/secret-id
```

### Vault Agent sidecar in Kubernetes

For Kubernetes workloads, the Vault Agent Injector sidecar handles authentication and secret delivery. The app reads secrets from files — no Vault SDK required.

```
+---------------------------Kubernetes Pod---------------------------+
|                                                                    |
|  +-------------------+       +----------------------------+       |
|  |   App Container   |       |  Vault Agent Sidecar       |       |
|  |                   | <---- |                            |       |
|  | reads:            |       | - authenticates via K8s SA |       |
|  | /vault/secrets/db |       | - writes rendered template |       |
|  |                   |       | - renews leases            |       |
|  +-------------------+       +----------+-----------------+       |
|                                         |                         |
+---------------------------------------- | -------------------------+
                                          | authenticate + fetch
                                          v
                                 +--------+---------+
                                 |   Vault Server   |
                                 |                  |
                                 | KV / Database    |
                                 | secrets engine   |
                                 +------------------+
```

Annotation on the pod spec enables injection:

```yaml
annotations:
  vault.hashicorp.com/agent-inject: "true"
  vault.hashicorp.com/role: "myapp-role"
  vault.hashicorp.com/agent-inject-secret-db: "secret/data/myapp/database"
  vault.hashicorp.com/agent-inject-template-db: |
    {{- with secret "secret/data/myapp/database" -}}
    DATABASE_URL={{ .Data.data.url }}
    {{- end }}
```

The sidecar writes the rendered template to `/vault/secrets/db`. The app reads it as a plain file.

---

## 4. AWS Secrets Manager

When you are already on AWS, Secrets Manager integrates with IAM and removes the need to run a Vault cluster.

### Go client

```bash
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
```

```go
package secrets

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type DBCredentials struct {
    URL      string `json:"url"`
    Username string `json:"username"`
    Password string `json:"password"`
    Host     string `json:"host"`
    Port     int    `json:"port"`
    DBName   string `json:"dbname"`
}

func NewSecretsManagerClient(ctx context.Context) (*secretsmanager.Client, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("LoadDefaultConfig: %w", err)
    }
    return secretsmanager.NewFromConfig(cfg), nil
}

// GetDatabaseCredentials fetches a JSON secret from AWS Secrets Manager
// and parses it into a DBCredentials struct.
func GetDatabaseCredentials(ctx context.Context, client *secretsmanager.Client, secretName string) (*DBCredentials, error) {
    result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    })
    if err != nil {
        return nil, fmt.Errorf("GetSecretValue %s: %w", secretName, err)
    }

    var creds DBCredentials
    if err := json.Unmarshal([]byte(*result.SecretString), &creds); err != nil {
        return nil, fmt.Errorf("unmarshal secret: %w", err)
    }
    return &creds, nil
}
```

Usage:

```go
smClient, err := secrets.NewSecretsManagerClient(ctx)
if err != nil {
    log.Fatal(err)
}

creds, err := secrets.GetDatabaseCredentials(ctx, smClient, "prod/myapp/database")
if err != nil {
    log.Fatal(err)
}

dsn := creds.URL  // or build DSN from creds.Host, creds.Port, etc.
```

### Automatic rotation

AWS Secrets Manager rotates secrets on a schedule using a Lambda function. For RDS databases, AWS provides a managed rotation Lambda.

```bash
aws secretsmanager rotate-secret \
  --secret-id prod/myapp/database \
  --rotation-lambda-arn arn:aws:lambda:us-east-1:123456789:function:SecretsManagerRDSRotation \
  --rotation-rules AutomaticallyAfterDays=30
```

When rotation runs: the Lambda creates new credentials in the database, updates the secret value in Secrets Manager, tests the new credentials, then removes the old ones. The app must call `GetSecretValue` on every startup (or on a short cache TTL — 5 to 10 minutes) to get current credentials. Do not cache credentials in a global variable for the lifetime of the process.

### IAM policy — least privilege

Each service gets access to exactly the secrets it needs. Nothing more.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowReadMyAppDatabaseSecret",
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret"
      ],
      "Resource": "arn:aws:secretsmanager:us-east-1:123456789012:secret:prod/myapp/database-*"
    }
  ]
}
```

The trailing `-*` in the ARN accounts for the suffix AWS appends to secret names. This policy grants read access to one secret, in one region, in one account. The service cannot list secrets, access other secrets, or delete anything. Attach this policy to the IAM role assumed by the service: EC2 instance profile, ECS task role, or EKS pod identity.

---

## 5. Kubernetes External Secrets Operator

The External Secrets Operator bridges external secret stores (Vault, AWS Secrets Manager, GCP Secret Manager) into native Kubernetes Secrets. Apps read a normal K8s Secret — no Vault SDK, no AWS SDK needed.

```bash
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace
```

Configure a SecretStore that points to Vault:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: vault-backend
  namespace: myapp
spec:
  provider:
    vault:
      server: "http://vault.vault:8200"
      path: "secret"
      version: "v2"
      auth:
        kubernetes:
          mountPath: "kubernetes"
          role: "myapp-role"
          serviceAccountRef:
            name: "myapp-sa"
```

Create an ExternalSecret that maps Vault paths to Kubernetes Secret keys:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: myapp-db-secret
  namespace: myapp
spec:
  refreshInterval: "5m"
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: myapp-db-secret          # name of the K8s Secret to create
    creationPolicy: Owner
  data:
    - secretKey: DATABASE_URL       # key in the resulting K8s Secret
      remoteRef:
        key: myapp/database         # path in Vault
        property: url               # field within the Vault secret
```

The operator creates and keeps in sync a K8s Secret named `myapp-db-secret`. It re-syncs every 5 minutes, so when a secret is rotated in Vault, the K8s Secret is updated automatically. Pods that mount the secret as an env var pick up the new value on the next restart or live-reload.

---

## 6. Doppler — Simpler Alternative

Doppler is a hosted secrets manager with a simpler operational model. You manage secrets through a dashboard, CLI, or API. It handles environments (dev/staging/prod) per project, per-team access controls, and audit logs without running any infrastructure.

```bash
# Install CLI
brew install dopplerhq/cli/doppler

# Log in and link to a project
doppler login
doppler setup

# Run the app with secrets injected as environment variables
doppler run -- go run ./cmd/api

# Run tests with secrets
doppler run -- go test ./...
```

The CLI fetches secrets from Doppler and injects them before the process starts. The app reads `os.Getenv` as normal. No SDK import, no code change.

For Kubernetes, the Doppler operator syncs secrets into native K8s Secrets:

```yaml
apiVersion: secrets.doppler.com/v1alpha1
kind: DopplerSecret
metadata:
  name: myapp-secrets
  namespace: myapp
spec:
  tokenSecret:
    name: doppler-token-secret    # K8s Secret containing the Doppler service token
  managedSecret:
    name: myapp-env               # K8s Secret to create and keep in sync
    namespace: myapp
```

When to use Doppler versus Vault:

| | Doppler | Vault |
|---|---------|-------|
| Setup time | Minutes | Hours to days |
| Dynamic DB credentials | No | Yes |
| Self-hosted | No (SaaS) | Yes |
| Cost | Free tier + paid | Free (self-hosted) |
| Best for | Small and medium teams | Enterprise, dynamic secrets |

Doppler is a good fit for teams that do not need dynamic credentials and do not have the operational capacity to run Vault.

---

## 7. Secrets Rotation

Storing secrets in Vault or AWS Secrets Manager solves the storage problem. Rotation limits the damage when a secret leaks.

### Update the secret

In Vault:

```bash
vault kv put secret/myapp/database url="postgres://app:new_password@db:5432/myapp"
```

In AWS Secrets Manager:

```bash
aws secretsmanager put-secret-value \
  --secret-id prod/myapp/database \
  --secret-string '{"url":"postgres://app:new_password@db:5432/myapp"}'
```

### Strategy 1 — Rolling restart

The simplest approach. The app reads secrets at startup. After rotating the secret, trigger a rolling deployment. Kubernetes replaces pods one at a time with `maxUnavailable: 0` so there is no downtime.

```bash
kubectl rollout restart deployment/myapp
```

This works for most cases. The window between secret rotation and rollout completion is when both the old and new credentials must be valid simultaneously — rotate the credential in the database before removing the old one.

### Strategy 2 — Live reload with fsnotify

When the Vault Agent Sidecar writes a new secret to a file, the app can detect the change and reload the connection pool without restarting.

```bash
go get github.com/fsnotify/fsnotify
```

```go
package reload

import (
    "log/slog"
    "os"

    "github.com/fsnotify/fsnotify"
)

// WatchSecretFile calls reload whenever the file at path is written or recreated.
// The Vault Agent Injector rewrites the file when it renews a lease.
func WatchSecretFile(path string, onReload func(newContent string)) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }

    go func() {
        defer watcher.Close()
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
                    data, err := os.ReadFile(path)
                    if err != nil {
                        slog.Error("read secret file", "path", path, "err", err)
                        continue
                    }
                    slog.Info("secret file updated, reloading", "path", path)
                    onReload(string(data))
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                slog.Error("secret watcher error", "err", err)
            }
        }
    }()

    return watcher.Add(path)
}
```

Usage — reload the database connection pool when credentials change:

```go
err := reload.WatchSecretFile("/vault/secrets/db", func(content string) {
    newURL := parseDBURL(content)

    newDB, err := sql.Open("pgx", newURL)
    if err != nil {
        slog.Error("failed to open new DB connection", "err", err)
        return
    }

    oldDB := atomicDB.Swap(newDB)

    // Give in-flight queries 30 seconds to finish on the old pool
    time.AfterFunc(30*time.Second, func() { oldDB.Close() })
})
```

### Database credential rotation — graceful handoff

When rotating database passwords, the sequence matters. Both the old and new credentials must be valid at the same time during the transition window:

```
1. Create new credential in the database (or let Vault/AWS do it)
2. Update the secret store with the new value
3. App detects change (file watch or cache TTL expiry)
4. App opens new connection pool with new credential
5. Drain old pool: stop issuing connections from it, wait for in-flight queries
6. Close old pool
7. Revoke old credential
```

If step 7 happens before step 5, in-flight queries fail with authentication errors. Most `database/sql` pools do not support graceful drain out of the box — implement it with an atomic pointer swap and a time-delayed close as shown above.

---

## Summary

| Tool | Best For | Dynamic Secrets | Self-Hosted |
|------|----------|-----------------|-------------|
| Environment variables | Local dev, simple apps | No | N/A |
| HashiCorp Vault | Any environment, full control | Yes | Yes |
| AWS Secrets Manager | AWS-native workloads | Via Lambda rotation | No (managed) |
| External Secrets Operator | Kubernetes, bridging any backend | Depends on backend | Yes (operator) |
| Doppler | Small and medium teams, fast setup | No | No (managed) |

The hardcoded-secret problem is easy to solve. The audit trail, rotation, and least-privilege problems require a dedicated secrets manager. Vault is the most capable but requires operational overhead. AWS Secrets Manager removes that overhead if you are already on AWS. Doppler is the fastest path for smaller teams.

Key practices regardless of tool:
- Never commit secrets to git. Use pre-commit hooks to scan for them (`gitleaks`, `trufflehog`).
- Rotate secrets regularly, not only after a suspected breach.
- Use short TTLs for dynamic credentials whenever possible.
- Apply least-privilege IAM or Vault policies: each service reads only what it needs.
- Always re-fetch secrets on connection errors. Do not cache credentials for the lifetime of the process.

---

## Exercises

### Easy

Write a Go function `GetDatabaseURL(ctx context.Context, client *secretsmanager.Client, secretName string) (string, error)` that:
1. Calls `GetSecretValue` with the given secret name.
2. Parses the JSON secret string as `{"url": "postgres://..."}` into a struct.
3. Returns the URL string.
4. Returns a descriptive error if the `url` key is missing or the value is empty.

Test it by storing `{"url": "postgres://localhost/test"}` in AWS Secrets Manager using LocalStack (`docker run -p 4566:4566 localstack/localstack`), then pointing the Go SDK at `http://localhost:4566`.

### Medium

Set up a local Vault dev server and write a Go program that reads a secret from it.

1. Run `vault server -dev`. Note the root token printed to stdout.
2. In another terminal: `vault kv put secret/myapp/db password="dev_password_123"`.
3. Write a Go program that:
   - Reads `VAULT_TOKEN` from the environment.
   - Creates a Vault client pointing to `http://127.0.0.1:8200`.
   - Sets the token on the client using `client.SetToken`.
   - Reads the secret at path `myapp/db` from the `secret` mount using `client.Secrets.KvV2Read`.
   - Prints the `password` field.
4. Verify it works. Then rotate the secret with `vault kv put secret/myapp/db password="new_password_456"` and re-run. Confirm the program reads the new value.

Do not hardcode the token in the source file — read it only from `os.Getenv("VAULT_TOKEN")`.

### Hard

Implement a `SecretProvider` interface with two concrete implementations so that production code uses Vault and tests use environment variables — without build tags or conditional compilation.

Define the interface:

```go
type SecretProvider interface {
    GetSecret(ctx context.Context, key string) (string, error)
}
```

Implement `EnvSecretProvider`:
- `GetSecret` calls `os.Getenv(key)`.
- Returns `fmt.Errorf("secret not found: %s", key)` if the value is empty.

Implement `VaultSecretProvider`:
- Constructor accepts a Vault address, mount path, and AppRole credentials (roleID, secretID).
- `GetSecret(ctx, key)` reads from KV v2 at the given key path.
- Cache each secret in a `sync.Map` with a 5-minute TTL. Re-fetch from Vault only when the TTL has expired. Store the value and the `time.Time` it was fetched together in a struct.

Wire it in `main.go`:
- If `VAULT_ADDR` is set, construct and use `VaultSecretProvider`.
- Otherwise fall back to `EnvSecretProvider`.

Write a unit test for a function `BuildDSN(ctx context.Context, sp SecretProvider) (string, error)` that calls `sp.GetSecret(ctx, "DATABASE_URL")` and returns the result. The test uses `EnvSecretProvider` by calling `t.Setenv("DATABASE_URL", "postgres://localhost/test")` before invoking `BuildDSN`. No Vault server is required to run the test.
