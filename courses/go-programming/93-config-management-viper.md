# Chapter 90: Configuration Management with Viper

Every application has configuration: database URLs, API keys, feature flags. You need to read from environment variables in production, config files in development, and command-line flags for tools. Viper handles all three with a clean precedence order and type-safe access.

## Table of Contents

1. [The Configuration Challenge](#1-the-configuration-challenge)
2. [Viper Basics](#2-viper-basics)
3. [Struct Binding](#3-struct-binding)
4. [Multiple Sources and Precedence](#4-multiple-sources-and-precedence)
5. [Secrets and Sensitive Values](#5-secrets-and-sensitive-values)
6. [Hot Reload](#6-hot-reload)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. The Configuration Challenge

Configuration sources (in priority order, highest first):
1. Command-line flags
2. Environment variables
3. Config file
4. Default values

```go
// Bad: scattered os.Getenv calls everywhere
func main() {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL")) // "" if not set
    port := os.Getenv("PORT") // might be empty
    maxConns, _ := strconv.Atoi(os.Getenv("MAX_CONNS")) // might fail
}

// Good: all configuration validated at startup, typed access everywhere
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
}

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatal("invalid config:", err) }
    // All config values are typed and validated
}
```

---

## 2. Viper Basics

```go
import (
    "github.com/spf13/viper"
    "github.com/spf13/pflag"
)

// Set defaults (lowest priority)
viper.SetDefault("server.port", 8080)
viper.SetDefault("server.read_timeout", "5s")
viper.SetDefault("database.max_connections", 25)

// Read from a config file
viper.SetConfigName("config")  // config.yaml, config.json, etc.
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.AddConfigPath("$HOME/.myapp")

if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
        log.Fatal("config file error:", err)
    }
    // Config file not found — use defaults and env vars only
}

// Read from environment variables
viper.SetEnvPrefix("APP")          // APP_SERVER_PORT
viper.AutomaticEnv()               // bind all env vars automatically
viper.SetEnvKeyReplacer(           // APP_DATABASE_MAX_CONNECTIONS → database.max_connections
    strings.NewReplacer(".", "_"))

// Access values
port := viper.GetInt("server.port")           // 8080
timeout := viper.GetDuration("server.read_timeout") // 5s
connStr := viper.GetString("database.url")
```

---

## 3. Struct Binding

Binding to a struct is cleaner than scattered `viper.GetX` calls:

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Auth     AuthConfig     `mapstructure:"auth"`
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
    Debug        bool          `mapstructure:"debug"`
}

type DatabaseConfig struct {
    URL            string        `mapstructure:"url"`
    MaxConnections int           `mapstructure:"max_connections"`
    MinConnections int           `mapstructure:"min_connections"`
    ConnMaxLife    time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
    Address  string `mapstructure:"address"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize int    `mapstructure:"pool_size"`
}

type AuthConfig struct {
    JWTSecret     string        `mapstructure:"jwt_secret"`
    TokenTTL      time.Duration `mapstructure:"token_ttl"`
    RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
}

func Load() (*Config, error) {
    v := viper.New()
    
    // Defaults
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.read_timeout", "5s")
    v.SetDefault("server.write_timeout", "10s")
    v.SetDefault("server.idle_timeout", "120s")
    v.SetDefault("database.max_connections", 25)
    v.SetDefault("database.min_connections", 5)
    v.SetDefault("database.conn_max_lifetime", "5m")
    v.SetDefault("redis.pool_size", 10)
    v.SetDefault("auth.token_ttl", "15m")
    v.SetDefault("auth.refresh_ttl", "7d")
    
    // Config file (optional)
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath(".")
    v.AddConfigPath("/etc/myapp/")
    v.ReadInConfig() // ignore not found
    
    // Environment variables (highest priority after flags)
    v.SetEnvPrefix("APP")
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshal config: %w", err)
    }
    
    return &cfg, validate(&cfg)
}

func validate(cfg *Config) error {
    var errs []string
    if cfg.Database.URL == "" {
        errs = append(errs, "database.url is required (APP_DATABASE_URL)")
    }
    if cfg.Auth.JWTSecret == "" {
        errs = append(errs, "auth.jwt_secret is required (APP_AUTH_JWT_SECRET)")
    }
    if len(cfg.Auth.JWTSecret) < 32 {
        errs = append(errs, "auth.jwt_secret must be at least 32 characters")
    }
    if len(errs) > 0 {
        return fmt.Errorf("invalid configuration:\n  %s", strings.Join(errs, "\n  "))
    }
    return nil
}
```

### config.yaml example

```yaml
server:
  port: 8080
  read_timeout: 5s
  write_timeout: 10s
  debug: false

database:
  url: "postgres://user:pass@localhost:5432/myapp?sslmode=disable"
  max_connections: 25

redis:
  address: "localhost:6379"
  pool_size: 10

auth:
  token_ttl: 15m
  refresh_ttl: 168h  # 7 days
```

---

## 4. Multiple Sources and Precedence

```
pflags (highest)
  ↓
environment variables: APP_SERVER_PORT=9090
  ↓
config.yaml: server.port: 8080
  ↓
defaults: 8080 (lowest)
```

### Binding command-line flags with pflag

```go
import "github.com/spf13/pflag"

func main() {
    pflag.Int("port", 0, "HTTP server port (overrides APP_SERVER_PORT and config file)")
    pflag.String("config", "", "Path to config file")
    pflag.Parse()
    
    viper.BindPFlags(pflag.CommandLine)
    
    cfg, _ := config.Load()
    // pflag --port 9090 wins over APP_SERVER_PORT=8080 wins over config.yaml
}
```

---

## 5. Secrets and Sensitive Values

Secrets should come from a secrets manager, not from config files committed to git.

```go
// Never put secrets in config.yaml — read them from the environment or a vault

// Option 1: Environment variable (simple, works everywhere)
// APP_AUTH_JWT_SECRET=... APP_DATABASE_URL=...

// Option 2: Doppler / Vault (production)
type SecretLoader interface {
    GetSecret(ctx context.Context, name string) (string, error)
}

type VaultLoader struct{ client *vault.Client }

func (v *VaultLoader) GetSecret(ctx context.Context, name string) (string, error) {
    secret, err := v.client.KVv2("secret").Get(ctx, name)
    if err != nil { return "", err }
    return secret.Data["value"].(string), nil
}

func LoadWithSecrets(ctx context.Context, secrets SecretLoader) (*Config, error) {
    cfg, err := Load() // load non-secret values from viper
    if err != nil { return nil, err }
    
    // Override secrets from vault
    cfg.Database.URL, err = secrets.GetSecret(ctx, "myapp/database_url")
    if err != nil { return nil, fmt.Errorf("load db secret: %w", err) }
    
    cfg.Auth.JWTSecret, err = secrets.GetSecret(ctx, "myapp/jwt_secret")
    if err != nil { return nil, fmt.Errorf("load jwt secret: %w", err) }
    
    return cfg, nil
}

// Never log sensitive config values
func (c *Config) SafeLog() map[string]any {
    return map[string]any{
        "server.port":     c.Server.Port,
        "database.url":    maskURL(c.Database.URL),
        "redis.address":   c.Redis.Address,
        "auth.token_ttl":  c.Auth.TokenTTL,
    }
}

func maskURL(dsn string) string {
    u, err := url.Parse(dsn)
    if err != nil { return "***" }
    u.User = url.UserPassword("***", "***")
    return u.String()
}
```

---

## 6. Hot Reload

Viper can watch config files and reload when they change:

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    log.Println("config file changed:", e.Name)
    
    var newCfg Config
    if err := viper.Unmarshal(&newCfg); err != nil {
        log.Println("invalid config, keeping old:", err)
        return
    }
    if err := validate(&newCfg); err != nil {
        log.Println("invalid config, keeping old:", err)
        return
    }
    
    // Swap config atomically
    configMu.Lock()
    currentConfig = &newCfg
    configMu.Unlock()
    
    log.Println("config reloaded successfully")
})
```

---

## Summary

- Viper merges multiple sources: flags > env vars > config file > defaults
- Use `mapstructure` tags to bind Viper values to a typed struct
- Validate at startup — fail fast with a clear error message
- Never store secrets in config files or log them
- `AutomaticEnv()` + `SetEnvKeyReplacer` maps `APP_DATABASE_MAX_CONNECTIONS` → `database.max_connections`
- `WatchConfig()` for runtime config changes without restart

## Exercises

### Easy
1. Create a `Config` struct for a simple web service: port, database URL, log level, and debug mode. Load it with Viper, setting defaults and mapping environment variables. Test that `APP_SERVER_PORT=9090` overrides the default.
2. Write a `config.yaml` with three environments: development, staging, production. Load the right file based on an `APP_ENV` environment variable. Use Viper's `AddConfigPath` to look in `./config/$APP_ENV.yaml`.
3. Add input validation: require that `database.url` starts with `postgres://`, `server.port` is between 1024-65535, and `log.level` is one of `debug|info|warn|error`. Return a collected error listing all violations.

### Medium
4. Implement a **feature flags** config section: `features.new_checkout: true`, `features.beta_search: false`. Write a `Features` struct and a `IsEnabled(feature string) bool` helper. Hot-reload the features when `config.yaml` changes without restarting the server.
5. Build a config hierarchy: `base.yaml` with shared settings, `production.yaml` that overrides only production-specific values. Use `viper.MergeInConfig` to layer them. Test that production values take precedence over base values.
6. Integrate **AWS Secrets Manager** or **HashiCorp Vault** (or a mock): write a `SecretLoader` interface, an implementation that returns hardcoded values for tests, and an implementation that calls the real secrets API. Wire the right one based on `APP_ENV`.

### Hard
7. Build a **dynamic config service**: config values are stored in a PostgreSQL table `config(key text, value text, updated_at timestamptz)`. The service polls this table every 30 seconds and updates its in-memory config. Write a `Get(key string) string` function that reads the latest value under a `sync.RWMutex`.
8. Implement a **config diff logger**: at startup, log the diff between the loaded config and the previous run's config (stored in a file). Redact sensitive fields. This helps answer "what changed since the last deploy?"
