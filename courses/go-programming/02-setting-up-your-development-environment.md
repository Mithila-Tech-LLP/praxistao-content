# Chapter 02: Setting Up Your Development Environment

Before writing a single line of Go, you need a working environment. This chapter walks you through every step: installing Go, configuring your editor, understanding the workspace, and running your first program. By the end, your machine will be ready for everything in this course.

## Table of Contents

1. [Installing Go](#1-installing-go)
2. [Understanding the Go Workspace](#2-understanding-the-go-workspace)
3. [Go Modules — Modern Dependency Management](#3-go-modules--modern-dependency-management)
4. [Setting Up Your Editor](#4-setting-up-your-editor)
5. [Essential Go Tools](#5-essential-go-tools)
6. [Your First Go Program](#6-your-first-go-program)
7. [Docker Setup for This Course](#7-docker-setup-for-this-course)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Installing Go

**Step 1: Download Go**

Go to the official website: **https://go.dev/dl/**

Download the installer for your operating system:
- **macOS**: `.pkg` file (e.g., `go1.23.0.darwin-amd64.pkg`)
- **Linux**: `.tar.gz` file (e.g., `go1.23.0.linux-amd64.tar.gz`)
- **Windows**: `.msi` file (e.g., `go1.23.0.windows-amd64.msi`)

**Step 2: Install**

*macOS*: Double-click the `.pkg` file and follow the installer. Go is installed to `/usr/local/go`.

*Linux*:
```bash
# Remove any previous Go installation
rm -rf /usr/local/go

# Extract the archive
tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Add to PATH in your ~/.bashrc or ~/.zshrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

*Windows*: Run the `.msi` installer. It adds Go to your PATH automatically.

**Step 3: Verify the installation**

Open a terminal and run:
```bash
go version
```

You should see something like:
```
go version go1.23.0 linux/amd64
```

**Check environment variables:**
```bash
go env GOPATH    # Where Go puts downloaded packages
go env GOROOT    # Where Go itself is installed
go env GOBIN     # Where `go install` puts executables
```

### Quick Check
> 1. What command verifies that Go installed correctly?
> 2. What is the difference between GOROOT and GOPATH?
> 3. Where do downloaded packages get stored?

---

## 2. Understanding the Go Workspace

Before Go modules (pre-2019), all Go code lived in a single `GOPATH` directory. This was confusing and limiting. Today, you can put your Go projects **anywhere** on your computer thanks to Go modules.

```
Your computer:
  
  ~/Desktop/
  ├── my-projects/
  │   ├── blog-api/          ← A Go project (has go.mod)
  │   │   ├── go.mod
  │   │   ├── go.sum
  │   │   ├── main.go
  │   │   └── ...
  │   └── ticket-system/     ← Another Go project
  │       ├── go.mod
  │       └── ...
  └── other-stuff/
  
  ~/go/                      ← GOPATH (auto-created)
  ├── pkg/
  │   └── mod/               ← Downloaded modules cache
  │       ├── github.com/
  │       ├── google.golang.org/
  │       └── ...
  └── bin/                   ← Installed CLI tools (golangci-lint, etc.)
```

**Key rule**: Every Go project starts with a `go.mod` file in its root directory. This file tells Go the project's name and its dependencies. We cover this in depth in the next section.

### Quick Check
> 1. Where is the Go module cache located on your computer?
> 2. What file makes a directory a Go project?
> 3. Is it required to put your Go project inside GOPATH?

---

## 3. Go Modules — Modern Dependency Management

A **module** is a collection of related Go packages with versioning. Every module has a `go.mod` file.

**Creating a new module:**
```bash
# Create your project directory
mkdir blog-api
cd blog-api

# Initialize a module
go mod init github.com/yourname/blog-api
```

This creates `go.mod`:
```
module github.com/yourname/blog-api

go 1.23
```

The module path (`github.com/yourname/blog-api`) is the unique name for your project. For personal projects it doesn't have to be an actual GitHub URL, but by convention it should match where you'd publish it.

**Adding dependencies:**
```bash
# Add a dependency (downloads it + updates go.mod + go.sum)
go get github.com/go-chi/chi/v5

# Tidy: remove unused deps, add missing deps
go mod tidy
```

After running `go get`, your `go.mod` looks like:
```
module github.com/yourname/blog-api

go 1.23

require (
    github.com/go-chi/chi/v5 v5.0.11
)
```

And `go.sum` is created — it contains cryptographic checksums of every dependency to ensure integrity:
```
github.com/go-chi/chi/v5 v5.0.11 h1:BnpYbFZ3T3S1WMpD79r7R5ThWX40TaFB7L31Y8xqSwA=
github.com/go-chi/chi/v5 v5.0.11/go.mod h1:DslCQbL2OYiznFReuXYUmovuv+y9CNkIUcp28tkdF9g=
```

**Key commands:**
```bash
go mod init <module-name>   # Create a new module
go get <package>@<version>  # Add or update a dependency
go mod tidy                 # Clean up go.mod and go.sum
go mod download             # Download all dependencies
go mod verify               # Verify checksums
go list -m all              # List all dependencies
```

### Quick Check
> 1. What is a Go module?
> 2. What command adds a new dependency to your project?
> 3. What does `go mod tidy` do?

---

## 4. Setting Up Your Editor

**Recommended: VS Code** (free, excellent Go support)

1. Download VS Code from https://code.visualstudio.com/
2. Open VS Code, go to Extensions (Ctrl+Shift+X)
3. Search for and install: **"Go"** by the Go Team at Google

When you open a `.go` file for the first time, VS Code will prompt you to install Go tools. Click **Install All**. This installs:
- `gopls` — the official Go language server (autocomplete, error detection, go to definition)
- `dlv` — Delve debugger
- `staticcheck` — static analysis
- `gotest` — test runner integration

**Your VS Code settings for Go** (add to `settings.json`):
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "editor.formatOnSave": true,
    "[go]": {
        "editor.defaultFormatter": "golang.go",
        "editor.codeActionsOnSave": {
            "source.organizeImports": "explicit"
        }
    }
}
```

**Alternative editors:**
- **GoLand** (JetBrains): Best Go IDE, paid (~$70/year), excellent refactoring
- **Vim/Neovim**: With `vim-go` or nvim-lspconfig + gopls
- **Zed**: New fast editor, good Go support

### Quick Check
> 1. What is `gopls` and why is it important?
> 2. How do you get VS Code to format Go code automatically on save?
> 3. What is Delve used for?

---

## 5. Essential Go Tools

These command-line tools are essential for Go development:

**Built into Go:**
```bash
go build       # Compile your program
go run         # Compile and run (for development)
go test        # Run tests
go fmt         # Format code (run this before committing!)
go vet         # Find common mistakes
go doc         # Read documentation in terminal
go generate    # Run code generators
go install     # Install a binary into GOBIN
```

**Install these tools:**
```bash
# golangci-lint: Runs many linters at once
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# goose: Database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# air: Live reload for development (like nodemon)
go install github.com/air-verse/air@latest

# dlv: Debugger (usually installed by VS Code automatically)
go install github.com/go-delve/delve/cmd/dlv@latest

# mockery: Generate mocks for interfaces
go install github.com/vektra/mockery/v2@latest

# sqlc: Generate type-safe Go code from SQL (Chapter 75)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**Using `air` for live reload:**
During development, you want your server to restart automatically when you change code. `air` does this:
```bash
# In your project root:
air
# Now when you edit any .go file, air recompiles and restarts
```

**Using `go fmt`:**
Go has one official formatter. There is no argument about code style. `go fmt` is always right:
```bash
go fmt ./...    # Format all Go files in current directory and subdirectories
```

**Using `go vet`:**
Catches common mistakes that compile but are probably wrong:
```bash
go vet ./...
```

### Quick Check
> 1. What is the difference between `go run` and `go build`?
> 2. What does `air` do and why is it useful during development?
> 3. What does `go vet` catch that `go build` doesn't?

---

## 6. Your First Go Program

Let's set up your first Go project from scratch:

```bash
# Create project directory
mkdir hello-go
cd hello-go

# Initialize module
go mod init hello-go

# Create main Go file
touch main.go
```

Open `main.go` in your editor and write:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

Run it:
```bash
go run main.go
```

Output:
```
Hello, Go!
```

**Breaking it down line by line:**

`package main` — Every Go file belongs to a package. The `main` package is special: it's the entry point for an executable program. All other package names are for libraries.

`import "fmt"` — Import the `fmt` package from the standard library. `fmt` stands for "format" and provides functions for printing.

`func main()` — The `main` function is where your program starts executing. Every executable Go program must have exactly one `main` function in the `main` package.

`fmt.Println("Hello, Go!")` — Call the `Println` function from the `fmt` package. `Println` prints a line to standard output with a newline at the end.

**Build and run as a binary:**
```bash
go build -o hello .    # Compile to binary named "hello"
./hello                # Run the binary
```

Output:
```
Hello, Go!
```

The resulting `hello` binary is a standalone executable. No Go installation needed on the target machine.

**Project structure for this course:**
For all projects in this course, we'll use this structure:
```
project-name/
├── cmd/
│   └── server/
│       └── main.go     ← Entry point
├── internal/
│   ├── domain/         ← Business logic
│   ├── handler/        ← HTTP handlers
│   ├── repository/     ← Database access
│   └── service/        ← Application services
├── pkg/                ← Reusable packages (if any)
├── migrations/         ← Database migrations
├── docker-compose.yml  ← Local services (DB, Redis, etc.)
├── go.mod
├── go.sum
└── README.md
```

### Quick Check
> 1. What is the `main` package and why is it special?
> 2. What is the difference between `go run` and `go build -o`?
> 3. What does the `internal/` directory convention mean in Go projects?

---

## 7. Docker Setup for This Course

Many chapters in this course require databases (PostgreSQL, Redis, MongoDB), message brokers (Kafka, RabbitMQ), and other services. We'll use Docker to run all of these locally without installing anything permanently.

**Install Docker Desktop:**
- macOS/Windows: https://www.docker.com/products/docker-desktop/
- Linux: Follow the instructions at https://docs.docker.com/engine/install/

**Verify Docker works:**
```bash
docker --version
docker compose version
```

**The `docker-compose.yml` you'll use throughout the course:**

Save this as `docker-compose.yml` in your project root whenever you need databases and services:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: goapp
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: appdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  mongodb:
    image: mongo:7
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: secret

  opensearch:
    image: opensearchproject/opensearch:2
    environment:
      - discovery.type=single-node
      - DISABLE_SECURITY_PLUGIN=true
    ports:
      - "9200:9200"

volumes:
  postgres_data:
```

**Useful Docker commands:**
```bash
# Start all services in the background
docker compose up -d

# Stop all services
docker compose down

# Stop and delete all data volumes
docker compose down -v

# See logs from a specific service
docker compose logs postgres

# Connect to PostgreSQL inside Docker
docker compose exec postgres psql -U goapp -d appdb
```

You don't need to start all services right now. We'll start only what's needed for each chapter.

### Quick Check
> 1. Why do we use Docker for databases instead of installing them directly?
> 2. What does `docker compose up -d` do?
> 3. What does `docker compose down -v` do differently from `docker compose down`?

---

## Summary

- **Go installation**: Download from go.dev/dl, verify with `go version`
- **Go workspace**: Your projects can live anywhere; modules handle dependency isolation
- **Go modules**: Every project has `go.mod` (dependencies) and `go.sum` (checksums). `go mod init`, `go get`, `go mod tidy` are the key commands
- **Editor setup**: VS Code + Go extension + all tools. `gopls` provides autocomplete and error detection
- **Essential tools**: `go fmt` (format), `go vet` (lint), `air` (live reload), `golangci-lint` (comprehensive linting)
- **First program**: `package main` + `func main()` + `go run` = running Go code
- **Project structure**: `cmd/`, `internal/`, `pkg/` convention for larger projects
- **Docker**: Runs databases and services locally without permanent installation

You are now ready to write Go. Next chapter: Hello World in depth, packages, and modules.

---

## Exercises

### Easy
1. Install Go and verify the version. Run `go env` and list what GOPATH, GOROOT, and GOBIN are set to on your machine.
2. Create a new Go module called `my-first-go` and write a program that prints your name.
3. Run `go fmt` on your file after deliberately misindenting it. What changes?

### Medium
4. Tool exploration: Install `air`. Create a simple Go program that prints "Server running...". Start it with `air`. Modify the print message while `air` is running. Observe it automatically restart.
5. Module dependencies: Create a new module and add `github.com/fatih/color` as a dependency (`go get github.com/fatih/color`). Write a program that prints "Hello, Go!" in green. Look at your `go.mod` and `go.sum` — what changed?
6. Build for different platforms: Go supports cross-compilation easily. Build your hello-world program for Linux (from macOS or Windows) with `GOOS=linux GOARCH=amd64 go build -o hello-linux .` What is the resulting file size? Compare to the binary for your native OS.

### Hard
7. Standard library exploration: Open the official Go standard library documentation at pkg.go.dev/std. Pick three packages you think will be useful for backend development and summarize: (a) what each package does, (b) its most important 2-3 functions/types, (c) give a small code example using each. This is practice for reading Go documentation, which you'll do throughout this course.
8. Project scaffolding: Using only your terminal and a text editor (no Go code generators), create a complete project scaffold for a "user authentication service" following the `cmd/internal/pkg` structure. It should have: `cmd/server/main.go`, `internal/domain/user.go` (user struct), `internal/handler/user.go` (placeholder HTTP handler), `internal/repository/user.go` (placeholder DB access), `internal/service/user.go` (placeholder business logic), `docker-compose.yml` with PostgreSQL. The files don't need to contain real code — just `package` declarations and comments explaining what each file will do.
