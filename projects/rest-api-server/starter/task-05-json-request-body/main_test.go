package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetStore() {
	store = nil
}

func TestCreateTodoHandler_FirstTodo(t *testing.T) {
	resetStore()

	body := `{"title":"Buy groceries"}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTodoHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var got Todo
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}
	if got.Title != "Buy groceries" {
		t.Errorf("expected title 'Buy groceries', got %q", got.Title)
	}
	if got.Done != false {
		t.Errorf("expected done false, got %v", got.Done)
	}
}

func TestCreateTodoHandler_SecondTodoAutoIncrements(t *testing.T) {
	resetStore()

	for _, title := range []string{"First task", "Second task"} {
		body := `{"title":"` + title + `"}`
		req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		CreateTodoHandler(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected 201 for %q, got %d", title, rr.Code)
		}
	}

	if len(store) != 2 {
		t.Fatalf("expected 2 todos in store, got %d", len(store))
	}
	if store[0].ID != 1 {
		t.Errorf("expected first todo ID 1, got %d", store[0].ID)
	}
	if store[1].ID != 2 {
		t.Errorf("expected second todo ID 2, got %d", store[1].ID)
	}
}

func TestCreateTodoHandler_ResponseBodyMatchesTodo(t *testing.T) {
	resetStore()

	body := `{"title":"Write tests"}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateTodoHandler(rr, req)

	var got Todo
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if got.Title != "Write tests" {
		t.Errorf("expected title 'Write tests', got %q", got.Title)
	}
	if got.ID == 0 {
		t.Error("expected non-zero ID in response")
	}
}
