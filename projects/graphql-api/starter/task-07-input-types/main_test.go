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
	users = []map[string]interface{}{}
	nextID = 1
}

func TestCreateUser_WithInputType(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { createUser(input: {name: "Alice", email: "alice@example.com"}) { id name email } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	user, ok := data["createUser"].(map[string]interface{})
	if !ok {
		t.Fatal("createUser did not return a user map")
	}

	if user["id"] != 1 {
		t.Errorf("expected id=1, got %v", user["id"])
	}
	if user["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", user["name"])
	}
	if user["email"] != "alice@example.com" {
		t.Errorf("expected email=alice@example.com, got %v", user["email"])
	}
}

func TestCreateUser_AppearsInList(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	runQuery(t, schema, `mutation { createUser(input: {name: "Alice", email: "alice@example.com"}) { id } }`)

	result := runQuery(t, schema, `{ users { id name email } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list := data["users"].([]interface{})
	if len(list) != 1 {
		t.Errorf("expected 1 user in list, got %d", len(list))
	}
}
