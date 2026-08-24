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
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
}

func TestUpdateUser_ChangesName(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { updateUser(id: 1, name: "Alicia") { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	user, ok := data["updateUser"].(map[string]interface{})
	if !ok {
		t.Fatal("updateUser did not return a user")
	}
	if user["name"] != "Alicia" {
		t.Errorf("expected name=Alicia, got %v", user["name"])
	}
}

func TestUpdateUser_NonExistentReturnsNull(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { updateUser(id: 999, name: "Ghost") { id name } }`)
	data := result.Data.(map[string]interface{})
	if data["updateUser"] != nil {
		t.Errorf("expected null for non-existent user, got %v", data["updateUser"])
	}
}

func TestDeleteUser_ReturnsTrue(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { deleteUser(id: 1) }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	if data["deleteUser"] != true {
		t.Errorf("expected deleteUser=true, got %v", data["deleteUser"])
	}
}

func TestDeleteUser_NonExistentReturnsFalse(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { deleteUser(id: 999) }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	if data["deleteUser"] != false {
		t.Errorf("expected deleteUser=false, got %v", data["deleteUser"])
	}
}

func TestDeleteUser_ReducesListSize(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	runQuery(t, schema, `mutation { deleteUser(id: 1) }`)

	result := runQuery(t, schema, `{ users { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list := data["users"].([]interface{})
	if len(list) != 1 {
		t.Errorf("expected 1 user after delete, got %d", len(list))
	}
}
