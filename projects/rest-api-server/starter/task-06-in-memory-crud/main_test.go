package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetDB() {
	db = map[int]Todo{}
	nextID = 1
}

func createTodo(t *testing.T, title string) Todo {
	t.Helper()
	body := `{"title":"` + title + `"}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("createTodo: expected 201, got %d", rr.Code)
	}
	var todo Todo
	if err := json.NewDecoder(rr.Body).Decode(&todo); err != nil {
		t.Fatalf("createTodo: invalid JSON: %v", err)
	}
	return todo
}

func TestListTodos_Empty(t *testing.T) {
	resetDB()

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var got []Todo
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty list, got %d items", len(got))
	}
}

func TestCreateTodo(t *testing.T) {
	resetDB()

	todo := createTodo(t, "Learn Go")
	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}
	if todo.Title != "Learn Go" {
		t.Errorf("expected title 'Learn Go', got %q", todo.Title)
	}
}

func TestListTodos_Populated(t *testing.T) {
	resetDB()

	createTodo(t, "Task A")
	createTodo(t, "Task B")

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	var got []Todo
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 todos, got %d", len(got))
	}
}

func TestGetTodoByID(t *testing.T) {
	resetDB()
	created := createTodo(t, "Specific task")

	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var got Todo
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
}

func TestGetTodoByID_NotFound(t *testing.T) {
	resetDB()

	req := httptest.NewRequest(http.MethodGet, "/todos/99", nil)
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateTodo(t *testing.T) {
	resetDB()
	createTodo(t, "Old title")

	body := `{"title":"New title","done":true}`
	req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var got Todo
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Title != "New title" {
		t.Errorf("expected title 'New title', got %q", got.Title)
	}
	if !got.Done {
		t.Error("expected done to be true")
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	resetDB()

	body := `{"title":"Ghost","done":false}`
	req := httptest.NewRequest(http.MethodPut, "/todos/99", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestDeleteTodo(t *testing.T) {
	resetDB()
	createTodo(t, "To be deleted")

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr.Code)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", rr.Body.String())
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	resetDB()

	req := httptest.NewRequest(http.MethodDelete, "/todos/99", nil)
	rr := httptest.NewRecorder()
	TodosHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
