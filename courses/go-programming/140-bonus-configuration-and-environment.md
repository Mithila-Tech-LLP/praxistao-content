# Chapter 49: Configuration and Environment Management

Configuration is how a single binary behaves differently in dev, staging, and production without code changes. A well-designed config system is explicit, validated at startup, and never silently uses defaults in production. This chapter covers the full config lifecycle.

## Table of Contents

1. [Config Principles](#1-config-principles)
2. [Environment Variables](#2-environment-variables)
3. [Config Files with Viper](#3-config-files-with-viper)
4. [Structured Config with Validation](#4-structured-config-with-validation)
5. [Secrets Management](#5-secrets-management)
6. [Feature Flags](#6-feature-flags)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Config Principles

**The 12-Factor App** says:
1. Store config in the **environment** (not in code or config files committed to VCS)
2. Strict separation of config from code — no hard-coded URLs, secrets, or timeouts
3. All config should be **explicit**: fail loudly at startup if required config is missing

```
Config sources (priority, highest to lowest):
  CLI flags → Environment variables → Config file → Defaults
```

**What belongs in config:**
- Port numbers, database URLs, secret keys
- Feature flags, timeouts, rate limits
- External service URLs

**What does NOT belong in config:**
- Business logic
- Data (use a database)
- Code (use feature branches)

---

## 2. Environment Variables

```go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

// Typed getters with defaults:
func Env(key, defaultVal string) string {
    if v := os.Getenv(key); v != "" { return v }
    return defaultVal
}

func EnvRequired(key string) (string, error) {
    v := os.Getenv(key)
    if v == "" { return "", fmt.Errorf("required environment variable %q not set", key) }
    return v, nil
}

func EnvInt(key string, defaultVal int) (int, error) {
    v := os.Getenv(key)
    if v == "" { return defaultVal, nil }
    n, err := strconv.Atoi(v)
    if err != nil { return 0, fmt.Errorf("env %q: expected integer, got %q", key, v) }
    return n, nil
}

func EnvDuration(key string, defaultVal time.Duration) (time.Duration, error) {
    v := os.Getenv(key)
    if v == "" { return defaultVal, nil }
    d, err := time.ParseDuration(v)
    if err != nil { return 0, fmt.Errorf("env %q: expected duration (e.g. 30s), got %q", key, v) }
    return d, nil
}

func EnvBool(key string, defaultVal bool) (bool, error) {
    v := os.Getenv(key)
    if v == "" { return defaultVal, nil }
    b, err := strconv.ParseBool(v)  // Accepts: 1, t, T, TRUE, true, 0, f, F, FALSE, false
    if err != nil { return false, fmt.Errorf("env %q: expected bool, got %q", key, v) }
    return b, nil
}
```

### .env files for local development
```go
// Use godotenv to load .env file in development:
// go get github.com/joho/godotenv

func LoadEnvFile() {
    if env := os.Getenv("APP_ENV"); env == "" || env == "development" {
        if err := godotenv.Load(".env"); err != nil {
            // .env is optional — don't fail
            slog.Debug("no .env file found", "err", err)
        }
    }
}

// .env file (never commit to git):
// PORT=8080
// DB_URL=postgres://user:pass@localhost/mydb
// JWT_SECRET=dev-secret-only
// LOG_LEVEL=debug
```

---

## 3. Config Files with Viper

```bash
go get github.com/spf13/viper
```

```go
package config

import (
    "github.com/spf13/viper"
)

func initViper() {
    viper.SetConfigName("config")        // config.yaml, config.json, config.toml
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")             // Current directory
    viper.AddConfigPath("$HOME/.myapp")  // Home directory
    viper.AddConfigPath("/etc/myapp")    // System-wide

    // Env vars override file values:
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")  // APP_PORT overrides port
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))  // db.url → DB_URL

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            log.Fatalf("config file error: %v", err)
        }
    }
}

// config.yaml:
// server:
//   port: 8080
//   read_timeout: 5s
// database:
//   url: postgres://localhost/mydb
//   max_conns: 25
// log:
//   level: info
```

---

## 4. Structured Config with Validation

This is the best pattern: load everything into a typed struct, validate at startup, fail fast.

```go
package config

import (
    "fmt"
    "os"
    "time"
)

type Config struct {
    Env      string
    Server   ServerConfig
    Database DatabaseConfig
    Auth     AuthConfig
    Log      LogConfig
}

type ServerConfig struct {
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type DatabaseConfig struct {
    URL             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

type AuthConfig struct {
    JWTSecret      string
    JWTIssuer      string
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration
}

type LogConfig struct {
    Level  string  // debug, info, warn, error
    Format string  // json, text
}

// Load builds the config from environment variables with validation.
func Load() (*Config, error) {
    var errs []error
    errorf := func(format string, args ...any) {
        errs = append(errs, fmt.Errorf(format, args...))
    }

    // Required values:
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" { errorf("DATABASE_URL is required") }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" { errorf("JWT_SECRET is required") }
    if len(jwtSecret) < 32 { errorf("JWT_SECRET must be at least 32 characters") }

    // Optional with defaults:
    port, err := envInt("PORT", 8080)
    if err != nil { errorf("%v", err) }

    readTimeout, err := envDuration("READ_TIMEOUT", 5*time.Second)
    if err != nil { errorf("%v", err) }

    writeTimeout, err := envDuration("WRITE_TIMEOUT", 10*time.Second)
    if err != nil { errorf("%v", err) }

    maxConns, err := envInt("DB_MAX_CONNS", 25)
    if err != nil { errorf("%v", err) }

    accessTTL, err := envDuration("ACCESS_TOKEN_TTL", 15*time.Minute)
    if err != nil { errorf("%v", err) }

    refreshTTL, err := envDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour)
    if err != nil { errorf("%v", err) }

    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" { logLevel = "info" }
    if !isValidLogLevel(logLevel) { errorf("LOG_LEVEL must be one of: debug, info, warn, error") }

    appEnv := os.Getenv("APP_ENV")
    if appEnv == "" { appEnv = "development" }

    // Collect all errors and fail at once:
    if len(errs) > 0 {
        msgs := make([]string, len(errs))
        for i, e := range errs { msgs[i] = "  - " + e.Error() }
        return nil, fmt.Errorf("configuration errors:\n%s", strings.Join(msgs, "\n"))
    }

    return &Config{
        Env: appEnv,
        Server: ServerConfig{
            Port:         port,
            ReadTimeout:  readTimeout,
            WriteTimeout: writeTimeout,
            IdleTimeout:  120 * time.Second,
        },
        Database: DatabaseConfig{
            URL:             dbURL,
            MaxOpenConns:    maxConns,
            MaxIdleConns:    maxConns / 2,
            ConnMaxLifetime: 30 * time.Minute,
        },
        Auth: AuthConfig{
            JWTSecret:       jwtSecret,
            JWTIssuer:       "myapp",
            AccessTokenTTL:  accessTTL,
            RefreshTokenTTL: refreshTTL,
        },
        Log: LogConfig{Level: logLevel, Format: "json"},
    }, nil
}

func isValidLogLevel(l string) bool {
    switch l { case "debug", "info", "warn", "error": return true }
    return false
}

func envInt(key string, def int) (int, error) {
    v := os.Getenv(key)
    if v == "" { return def, nil }
    n, err := strconv.Atoi(v)
    if err != nil { return 0, fmt.Errorf("%s must be an integer, got %q", key, v) }
    return n, nil
}

func envDuration(key string, def time.Duration) (time.Duration, error) {
    v := os.Getenv(key)
    if v == "" { return def, nil }
    d, err := time.ParseDuration(v)
    if err != nil { return 0, fmt.Errorf("%s must be a duration (e.g. 30s), got %q", key, v) }
    return d, nil
}

// main.go:
func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)  // Exit early with all config errors at once
    }

    // Wire dependencies:
    logger := buildLogger(cfg.Log)
    db := buildDB(cfg.Database)
    srv := buildServer(cfg, db, logger)
    srv.Run()
}
```

---

## 5. Secrets Management

Never put secrets in config files or environment variables in production if you can avoid it.

```go
// AWS Secrets Manager integration:
func LoadSecretFromAWS(ctx context.Context, secretName string) (map[string]string, error) {
    cfg, err := awsconfig.LoadDefaultConfig(ctx)
    if err != nil { return nil, err }

    client := secretsmanager.NewFromConfig(cfg)
    result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    })
    if err != nil { return nil, err }

    var secrets map[string]string
    json.Unmarshal([]byte(*result.SecretString), &secrets)
    return secrets, nil
}

// HashiCorp Vault:
import vault "github.com/hashicorp/vault/api"

func LoadFromVault(secretPath string) (map[string]any, error) {
    client, err := vault.NewClient(vault.DefaultConfig())
    if err != nil { return nil, err }

    secret, err := client.Logical().Read(secretPath)
    if err != nil { return nil, err }
    return secret.Data, nil
}

// At minimum: never log config values that might be secrets:
func (c *Config) String() string {
    return fmt.Sprintf("Config{Env:%s Port:%d DB:***redacted*** JWT:***redacted***}", c.Env, c.Server.Port)
}
```

---

## 6. Feature Flags

```go
// Simple in-process feature flags (no external service):
type FeatureFlags struct {
    NewPaymentFlow    bool
    BetaUserOnboarding bool
    EnableMetrics     bool
}

func LoadFeatureFlags() FeatureFlags {
    return FeatureFlags{
        NewPaymentFlow:    os.Getenv("FEATURE_NEW_PAYMENT") == "true",
        BetaUserOnboarding: os.Getenv("FEATURE_BETA_ONBOARDING") == "true",
        EnableMetrics:     os.Getenv("ENABLE_METRICS") != "false",  // Default on
    }
}

// Use in handler:
func paymentHandler(flags *FeatureFlags) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if flags.NewPaymentFlow {
            newPaymentHandler(w, r)
        } else {
            legacyPaymentHandler(w, r)
        }
    }
}

// Percentage rollout:
func (f *FeatureFlags) IsEnabledForUser(userID int64, featureName string, percent int) bool {
    h := fnv.New32a()
    fmt.Fprintf(h, "%s:%d", featureName, userID)
    return int(h.Sum32()%100) < percent
}
```

---

## Summary

- Config in environment variables: explicit, overridable per deployment, no secrets in git
- Typed config struct: load all env vars at startup, validate everything, fail with all errors at once
- Load order: CLI flags > env vars > config file > defaults
- Use `.env` files locally with `godotenv`; never commit `.env` to version control
- Secrets: use AWS Secrets Manager / Vault in production; rotate frequently; never log secret values
- Feature flags: env var booleans for simple toggles; percentage rollout for gradual releases
- Fail fast at startup with clear error messages — better than mysterious runtime failures

---

## Exercises

### Easy
1. Write a `config.Load()` function for a simple web API with required fields: `PORT`, `DATABASE_URL`, `JWT_SECRET`. Required fields that are missing should fail with a single error listing ALL missing fields (not one by one). Write a test that sets only some env vars and verifies the error message lists the missing ones.
2. Add a `config.Print()` method that logs the config at startup but redacts secrets (`JWT_SECRET → "***"`, `DATABASE_URL → "postgres://***@host/db"`). Test that the output never contains the actual secret value.
3. Implement `EnvDuration` with validation that the value is within a sensible range: timeouts must be between 1ms and 1 hour; return an error outside this range.

### Medium
4. Build a **config watcher**: use `fsnotify` (`go get github.com/fsnotify/fsnotify`) to watch a config file for changes. When it changes, reload the non-sensitive fields (log level, feature flags, rate limits) without restarting the server. Protect the mutable fields with `sync.RWMutex`. Test by changing the config file and verifying the server behavior changes without a restart.
5. Implement a **config diff logger**: when config is loaded, compare it to the previous load (or to defaults). Log each field that changed with its old and new value (redacting secrets). This makes it easy to see what changed between deployments.
6. Integrate **AWS Secrets Manager** in a test environment by mocking it: create a `SecretsProvider` interface with `GetSecret(ctx, name) (string, error)`. Implement `EnvSecretsProvider` (reads from env vars) and `AWSSecretsProvider` (reads from AWS). In tests, use the env provider; in production, use AWS. Pass the provider to `config.Load()`.

### Hard
7. Build a **GrowthBook-style feature flag system**: flags have a name, percentage rollout (0-100), and optional user ID allowlist. Store flag definitions in a YAML file. Implement `IsEnabled(ctx, flagName string) bool` that checks: if the user ID (from context) is in the allowlist, always enabled; otherwise, use the hash-based percentage check. Reload flags on file change using fsnotify.
8. Build a **config validation test**: write a test that loads every config variable your application uses from a `config_test.env` file (with safe test values), calls `config.Load()`, and verifies each field is correctly typed and within valid ranges. This test runs in CI and catches config schema changes that aren't reflected in deployment scripts.
