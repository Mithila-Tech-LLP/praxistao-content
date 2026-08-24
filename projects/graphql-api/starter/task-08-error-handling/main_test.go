package main

import (
	"testing"

	"github.com/graphql-go/graphql"
)

func runQuery(t *testing.T, schema graphql.Schema, query string) *graphql.Result {
	t.Helper()
	result := graphql.Do(graphql.Params{Schema: schema, RequestString: query})
	return result
}

func resetState() {
	users = []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
	}
	nextID = 2
}

func TestUserNotFound_ReturnsError(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ user(id: 999) { id } }`)
	if len(result.Errors) == 0 {
		t.Error("expected errors for non-existent user, got none")
	}
}

func TestCreateUser_EmptyName_ReturnsError(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { createUser(input: {name: "", email: "a@b.com"}) { id } }`)
	if len(result.Errors) == 0 {
		t.Error("expected error for empty name, got none")
	}

	// Ensure no user was created
	listResult := runQuery(t, schema, `{ users { id } }`)
	data := listResult.Data.(map[string]interface{})
	list := data["users"].([]interface{})
	if len(list) != 1 {
		t.Errorf("expected users list unchanged (1 user), got %d", len(list))
	}
}

func TestCreateUser_EmptyEmail_ReturnsError(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { createUser(input: {name: "Alice", email: ""}) { id } }`)
	if len(result.Errors) == 0 {
		t.Error("expected error for empty email, got none")
	}
}

func TestValidUser_Works(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ user(id: 1) { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors for valid user: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	user, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatal("user field is not a map")
	}
	if user["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", user["name"])
	}
}
