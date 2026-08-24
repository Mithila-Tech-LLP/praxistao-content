# Chapter 17: Mini Project 1 — CLI Task Manager

It's time to apply everything from Chapters 01–16 in a real project. You'll build a command-line task manager — think a simplified `todo` app — that persists data to a JSON file. This project exercises: structs, methods, interfaces, error handling, file I/O, slices, maps, and testing.

## Project Overview

**What you'll build:** A CLI tool called `tasks` that manages a list of to-do items.

```bash
tasks add "Buy groceries"
tasks add "Write Go code" --priority high
tasks list
tasks list --filter pending
tasks done 3
tasks delete 2
tasks clear
```

**What you'll learn by building it:**
- Real file I/O with JSON encoding/decoding
- Structuring a small Go project cleanly
- The flag package for CLI arguments
- Writing tests for file-based persistence
- Error handling in a real workflow

---

## Project Structure

```
tasks/
├── main.go           ← Entry point, command dispatch
├── task.go           ← Task type and business logic
├── store.go          ← File-based persistence
├── display.go        ← Formatting and output
├── task_test.go      ← Tests for task logic
├── store_test.go     ← Tests for persistence
└── go.mod
```

---

## Step 1: Define the Data Model (`task.go`)

Start with what a task IS:

```go
// task.go
package main

import (
    "fmt"
    "time"
)

type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
)

type Status string

const (
    StatusPending  Status = "pending"
    StatusDone     Status = "done"
)

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Priority    Priority  `json:"priority"`
    Status      Status    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NewTask creates a task with defaults applied.
func NewTask(title string, priority Priority) Task {
    if priority == "" {
        priority = PriorityMedium
    }
    return Task{
        Title:     title,
        Priority:  priority,
        Status:    StatusPending,
        CreatedAt: time.Now(),
    }
}

func (t Task) IsDone() bool {
    return t.Status == StatusDone
}

func (t *Task) Complete() {
    now := time.Now()
    t.Status = StatusDone
    t.CompletedAt = &now
}

func (t Task) String() string {
    status := "[ ]"
    if t.IsDone() {
        status = "[✓]"
    }
    return fmt.Sprintf("%s #%d [%s] %s", status, t.ID, t.Priority, t.Title)
}
```

**Key design decisions:**
- `Priority` and `Status` are `string` type aliases — cleaner than raw strings, still JSON-compatible
- `CompletedAt` is `*time.Time` — nil means not completed (optional field pattern from Ch 13)
- Methods have pointer receivers when they mutate (`Complete`), value receivers when they don't (`IsDone`, `String`)

---

## Step 2: File Persistence (`store.go`)

```go
// store.go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
)

var ErrTaskNotFound = errors.New("task not found")

type Store struct {
    path  string
    tasks []Task
    nextID int
}

// NewStore creates or loads a store from the given file path.
func NewStore(path string) (*Store, error) {
    s := &Store{path: path}
    if err := s.load(); err != nil {
        return nil, fmt.Errorf("loading store: %w", err)
    }
    return s, nil
}

// load reads tasks from the JSON file. If the file doesn't exist, starts fresh.
func (s *Store) load() error {
    data, err := os.ReadFile(s.path)
    if errors.Is(err, os.ErrNotExist) {
        s.tasks = []Task{}
        s.nextID = 1
        return nil
    }
    if err != nil {
        return fmt.Errorf("reading file: %w", err)
    }

    type storedData struct {
        Tasks  []Task `json:"tasks"`
        NextID int    `json:"next_id"`
    }
    var d storedData
    if err := json.Unmarshal(data, &d); err != nil {
        return fmt.Errorf("parsing JSON: %w", err)
    }

    s.tasks = d.Tasks
    s.nextID = d.NextID
    if s.nextID == 0 {
        s.nextID = 1
    }
    return nil
}

// save writes tasks to the JSON file atomically.
func (s *Store) save() error {
    type storedData struct {
        Tasks  []Task `json:"tasks"`
        NextID int    `json:"next_id"`
    }
    d := storedData{Tasks: s.tasks, NextID: s.nextID}

    data, err := json.MarshalIndent(d, "", "  ")
    if err != nil {
        return fmt.Errorf("marshaling JSON: %w", err)
    }

    // Write to temp file first, then rename — atomic operation
    tmp := s.path + ".tmp"
    if err := os.WriteFile(tmp, data, 0644); err != nil {
        return fmt.Errorf("writing temp file: %w", err)
    }
    if err := os.Rename(tmp, s.path); err != nil {
        return fmt.Errorf("renaming file: %w", err)
    }
    return nil
}

// Add creates a new task and persists it.
func (s *Store) Add(title string, priority Priority) (Task, error) {
    t := NewTask(title, priority)
    t.ID = s.nextID
    s.nextID++
    s.tasks = append(s.tasks, t)
    if err := s.save(); err != nil {
        return Task{}, fmt.Errorf("saving after add: %w", err)
    }
    return t, nil
}

// List returns tasks, optionally filtered by status.
func (s *Store) List(filter Status) []Task {
    if filter == "" {
        result := make([]Task, len(s.tasks))
        copy(result, s.tasks)
        return result
    }
    var result []Task
    for _, t := range s.tasks {
        if t.Status == filter {
            result = append(result, t)
        }
    }
    return result
}

// Get returns a task by ID.
func (s *Store) Get(id int) (Task, error) {
    for _, t := range s.tasks {
        if t.ID == id {
            return t, nil
        }
    }
    return Task{}, fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
}

// Complete marks a task as done.
func (s *Store) Complete(id int) error {
    for i := range s.tasks {
        if s.tasks[i].ID == id {
            if s.tasks[i].IsDone() {
                return fmt.Errorf("task %d is already done", id)
            }
            s.tasks[i].Complete()
            return s.save()
        }
    }
    return fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
}

// Delete removes a task by ID.
func (s *Store) Delete(id int) error {
    for i, t := range s.tasks {
        if t.ID == id {
            // Remove element i without preserving order:
            s.tasks[i] = s.tasks[len(s.tasks)-1]
            s.tasks = s.tasks[:len(s.tasks)-1]
            return s.save()
        }
    }
    return fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
}

// Clear removes all tasks.
func (s *Store) Clear() error {
    s.tasks = []Task{}
    s.nextID = 1
    return s.save()
}

// Count returns total and done counts.
func (s *Store) Count() (total, done int) {
    total = len(s.tasks)
    for _, t := range s.tasks {
        if t.IsDone() {
            done++
        }
    }
    return
}
```

**Key techniques used:**
- `errors.Is(err, os.ErrNotExist)` — correctly detect missing file (Ch 14)
- Atomic write via temp file + rename — prevents corruption if the process crashes mid-write
- `copy(result, s.tasks)` — return a copy to prevent callers from modifying internal state

---

## Step 3: Display (`display.go`)

```go
// display.go
package main

import (
    "fmt"
    "strings"
)

func printTasks(tasks []Task) {
    if len(tasks) == 0 {
        fmt.Println("No tasks found.")
        return
    }
    
    // Column header:
    fmt.Printf("%-4s %-6s %-8s %-40s %s\n", "ID", "Status", "Priority", "Title", "Created")
    fmt.Println(strings.Repeat("-", 70))
    
    for _, t := range tasks {
        status := "pending"
        if t.IsDone() {
            status = "done"
        }
        created := t.CreatedAt.Format("2006-01-02")
        fmt.Printf("%-4d %-6s %-8s %-40s %s\n",
            t.ID, status, string(t.Priority), truncate(t.Title, 40), created)
    }
}

func printSummary(total, done int) {
    pending := total - done
    fmt.Printf("\n%d tasks total: %d pending, %d done\n", total, pending, done)
}

func truncate(s string, max int) string {
    if len(s) <= max {
        return s
    }
    return s[:max-3] + "..."
}
```

---

## Step 4: Main Entry Point (`main.go`)

```go
// main.go
package main

import (
    "fmt"
    "os"
    "strconv"
)

const dataFile = ".tasks.json"

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    store, err := NewStore(dataFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    command := os.Args[1]
    args := os.Args[2:]

    switch command {
    case "add":
        runAdd(store, args)
    case "list", "ls":
        runList(store, args)
    case "done":
        runDone(store, args)
    case "delete", "rm":
        runDelete(store, args)
    case "clear":
        runClear(store)
    case "help", "--help", "-h":
        printUsage()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
        printUsage()
        os.Exit(1)
    }
}

func runAdd(store *Store, args []string) {
    if len(args) == 0 {
        fmt.Fprintln(os.Stderr, "Error: add requires a title")
        fmt.Fprintln(os.Stderr, "Usage: tasks add <title> [--priority low|medium|high]")
        os.Exit(1)
    }

    title := args[0]
    priority := PriorityMedium

    // Parse optional --priority flag:
    for i := 1; i < len(args)-1; i++ {
        if args[i] == "--priority" || args[i] == "-p" {
            switch args[i+1] {
            case "low":
                priority = PriorityLow
            case "high":
                priority = PriorityHigh
            case "medium":
                priority = PriorityMedium
            default:
                fmt.Fprintf(os.Stderr, "Unknown priority: %s\n", args[i+1])
                os.Exit(1)
            }
        }
    }

    task, err := store.Add(title, priority)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Added task #%d: %s\n", task.ID, task.Title)
}

func runList(store *Store, args []string) {
    var filter Status
    for i := 0; i < len(args)-1; i++ {
        if args[i] == "--filter" || args[i] == "-f" {
            filter = Status(args[i+1])
        }
    }

    tasks := store.List(filter)
    printTasks(tasks)
    
    total, done := store.Count()
    printSummary(total, done)
}

func runDone(store *Store, args []string) {
    if len(args) == 0 {
        fmt.Fprintln(os.Stderr, "Error: done requires a task ID")
        os.Exit(1)
    }

    id, err := strconv.Atoi(args[0])
    if err != nil {
        fmt.Fprintf(os.Stderr, "Invalid ID %q: must be a number\n", args[0])
        os.Exit(1)
    }

    if err := store.Complete(id); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Task #%d marked as done.\n", id)
}

func runDelete(store *Store, args []string) {
    if len(args) == 0 {
        fmt.Fprintln(os.Stderr, "Error: delete requires a task ID")
        os.Exit(1)
    }

    id, err := strconv.Atoi(args[0])
    if err != nil {
        fmt.Fprintf(os.Stderr, "Invalid ID %q: must be a number\n", args[0])
        os.Exit(1)
    }

    if err := store.Delete(id); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Task #%d deleted.\n", id)
}

func runClear(store *Store) {
    if err := store.Clear(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("All tasks cleared.")
}

func printUsage() {
    fmt.Println(`tasks — a simple CLI task manager

Usage:
  tasks add <title> [--priority low|medium|high]
  tasks list [--filter pending|done]
  tasks done <id>
  tasks delete <id>
  tasks clear
  tasks help

Examples:
  tasks add "Buy groceries"
  tasks add "Write tests" --priority high
  tasks list
  tasks list --filter pending
  tasks done 3
  tasks delete 2`)
}
```

---

## Step 5: Tests (`task_test.go`)

```go
// task_test.go
package main

import (
    "testing"
)

func TestNewTask(t *testing.T) {
    tests := []struct {
        name     string
        title    string
        priority Priority
        wantPri  Priority
    }{
        {"with priority", "Buy milk", PriorityHigh, PriorityHigh},
        {"default priority", "Buy milk", "", PriorityMedium},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            task := NewTask(tt.title, tt.priority)
            if task.Title != tt.title {
                t.Errorf("Title = %q; want %q", task.Title, tt.title)
            }
            if task.Priority != tt.wantPri {
                t.Errorf("Priority = %q; want %q", task.Priority, tt.wantPri)
            }
            if task.Status != StatusPending {
                t.Errorf("Status = %q; want pending", task.Status)
            }
            if task.CreatedAt.IsZero() {
                t.Error("CreatedAt should not be zero")
            }
            if task.CompletedAt != nil {
                t.Error("CompletedAt should be nil for new tasks")
            }
        })
    }
}

func TestTask_Complete(t *testing.T) {
    task := NewTask("Do something", PriorityLow)
    if task.IsDone() {
        t.Error("new task should not be done")
    }

    task.Complete()
    if !task.IsDone() {
        t.Error("task should be done after Complete()")
    }
    if task.CompletedAt == nil {
        t.Error("CompletedAt should be set after Complete()")
    }
}

func TestTask_String(t *testing.T) {
    task := Task{ID: 5, Title: "Test task", Priority: PriorityHigh, Status: StatusPending}
    s := task.String()
    if s == "" {
        t.Error("String() should not be empty")
    }
    // Should contain ID and title:
    if !contains(s, "5") || !contains(s, "Test task") {
        t.Errorf("String() = %q; should contain ID and title", s)
    }
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr ||
        len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, sub string) bool {
    for i := 0; i <= len(s)-len(sub); i++ {
        if s[i:i+len(sub)] == sub {
            return true
        }
    }
    return false
}
```

```go
// store_test.go
package main

import (
    "errors"
    "path/filepath"
    "testing"
)

func newTestStore(t *testing.T) *Store {
    t.Helper()
    dir := t.TempDir()
    store, err := NewStore(filepath.Join(dir, "test_tasks.json"))
    if err != nil {
        t.Fatalf("creating test store: %v", err)
    }
    return store
}

func TestStore_Add(t *testing.T) {
    store := newTestStore(t)
    
    task, err := store.Add("Test task", PriorityHigh)
    if err != nil {
        t.Fatalf("Add: %v", err)
    }
    if task.ID != 1 {
        t.Errorf("first task ID = %d; want 1", task.ID)
    }
    if task.Title != "Test task" {
        t.Errorf("Title = %q; want Test task", task.Title)
    }
    if task.Priority != PriorityHigh {
        t.Errorf("Priority = %q; want high", task.Priority)
    }
    
    // Second task gets ID 2:
    task2, err := store.Add("Another task", PriorityLow)
    if err != nil {
        t.Fatalf("Add second: %v", err)
    }
    if task2.ID != 2 {
        t.Errorf("second task ID = %d; want 2", task2.ID)
    }
}

func TestStore_Persistence(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "tasks.json")
    
    // Create store and add tasks:
    store1, _ := NewStore(path)
    store1.Add("Task 1", PriorityHigh)
    store1.Add("Task 2", PriorityLow)
    
    // Load from same file — should see the tasks:
    store2, err := NewStore(path)
    if err != nil {
        t.Fatalf("loading persisted store: %v", err)
    }
    
    tasks := store2.List("")
    if len(tasks) != 2 {
        t.Errorf("loaded %d tasks; want 2", len(tasks))
    }
    if tasks[0].Title != "Task 1" {
        t.Errorf("tasks[0].Title = %q; want Task 1", tasks[0].Title)
    }
}

func TestStore_Complete(t *testing.T) {
    store := newTestStore(t)
    task, _ := store.Add("Task", PriorityMedium)
    
    if err := store.Complete(task.ID); err != nil {
        t.Fatalf("Complete: %v", err)
    }
    
    got, _ := store.Get(task.ID)
    if !got.IsDone() {
        t.Error("task should be done after Complete")
    }
    
    // Completing already-done task should error:
    err := store.Complete(task.ID)
    if err == nil {
        t.Error("expected error completing already-done task")
    }
}

func TestStore_Delete(t *testing.T) {
    store := newTestStore(t)
    task, _ := store.Add("Task to delete", PriorityMedium)
    
    if err := store.Delete(task.ID); err != nil {
        t.Fatalf("Delete: %v", err)
    }
    
    _, err := store.Get(task.ID)
    if !errors.Is(err, ErrTaskNotFound) {
        t.Errorf("expected ErrTaskNotFound after delete, got %v", err)
    }
    
    // Deleting non-existent task should error:
    err = store.Delete(999)
    if !errors.Is(err, ErrTaskNotFound) {
        t.Errorf("expected ErrTaskNotFound, got %v", err)
    }
}

func TestStore_List_Filter(t *testing.T) {
    store := newTestStore(t)
    store.Add("Task 1", PriorityHigh)
    task2, _ := store.Add("Task 2", PriorityLow)
    store.Add("Task 3", PriorityMedium)
    store.Complete(task2.ID)
    
    pending := store.List(StatusPending)
    if len(pending) != 2 {
        t.Errorf("pending count = %d; want 2", len(pending))
    }
    
    done := store.List(StatusDone)
    if len(done) != 1 {
        t.Errorf("done count = %d; want 1", len(done))
    }
    
    all := store.List("")
    if len(all) != 3 {
        t.Errorf("total count = %d; want 3", len(all))
    }
}

func TestStore_Clear(t *testing.T) {
    store := newTestStore(t)
    store.Add("Task 1", PriorityHigh)
    store.Add("Task 2", PriorityLow)
    
    if err := store.Clear(); err != nil {
        t.Fatalf("Clear: %v", err)
    }
    
    tasks := store.List("")
    if len(tasks) != 0 {
        t.Errorf("expected 0 tasks after clear, got %d", len(tasks))
    }
    
    // After clear, IDs should restart from 1:
    newTask, _ := store.Add("New task", PriorityMedium)
    if newTask.ID != 1 {
        t.Errorf("ID after clear = %d; want 1", newTask.ID)
    }
}
```

---

## Step 6: Build and Run

```bash
# Initialize the module:
cd tasks
go mod init tasks

# Build the binary:
go build -o tasks .

# Try it:
./tasks add "Buy groceries"
./tasks add "Write Go code" --priority high
./tasks add "Read a book" --priority low
./tasks list
./tasks done 1
./tasks list --filter done
./tasks delete 2
./tasks list
./tasks clear

# Run tests:
go test ./... -v

# Run tests with coverage:
go test ./... -cover

# Build for different platforms:
GOOS=linux GOARCH=amd64 go build -o tasks-linux .
GOOS=windows GOARCH=amd64 go build -o tasks.exe .
```

---

## What You Just Practiced

| Concept | Where used |
|---------|-----------|
| Structs + methods | `Task`, `Store` |
| Pointer receivers | `Task.Complete()`, all `Store` methods |
| Value receivers | `Task.IsDone()`, `Task.String()` |
| Interfaces (`error`) | All error returns |
| Error wrapping (`%w`) | `load`, `save`, `Add`, etc. |
| Sentinel errors | `ErrTaskNotFound` |
| Slices | `Store.tasks`, `List` results |
| Pointer to value | `CompletedAt *time.Time` |
| JSON struct tags | All `Task` fields |
| File I/O | `os.ReadFile`, `os.WriteFile`, `os.Rename` |
| `errors.Is` | `ErrNotExist`, `ErrTaskNotFound` checks |
| Table-driven tests | `TestNewTask` |
| `t.TempDir()` | `newTestStore` helper |
| Test helpers with `t.Helper()` | `newTestStore` |

---

## Extension Challenges

Once the base project works, try these:

1. **Due dates**: Add a `DueAt *time.Time` field. Show "(overdue)" next to tasks past their due date. Add `tasks add "Buy milk" --due 2024-12-25`.

2. **Tags**: Add `Tags []string`. Allow filtering by tag: `tasks list --tag work`. Add `tasks tag 3 work urgent`.

3. **Edit**: Implement `tasks edit <id> --title "New title" --priority high` — modify an existing task.

4. **Sort**: Add `--sort priority` and `--sort created` flags to the list command.

5. **Export**: Add `tasks export` that writes tasks to CSV or Markdown format.

6. **Undo**: Keep a small undo stack (last 5 operations). Implement `tasks undo` that reverts the most recent add/delete/complete.

7. **Multiple lists**: Support named lists: `tasks --list work add "Meeting"`, `tasks --list personal list`.
