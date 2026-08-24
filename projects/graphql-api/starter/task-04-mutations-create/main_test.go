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
	nextID = 3
}

func TestCreateUser_FirstUser(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `mutation { createUser(name: "Dave") { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	user, ok := data["createUser"].(map[string]interface{})
	if !ok {
		t.Fatal("createUser did not return a user")
	}
	if user["id"] != 3 {
		t.Errorf("expected id=3, got %v", user["id"])
	}
	if user["name"] != "Dave" {
		t.Errorf("expected name=Dave, got %v", user["name"])
	}
}

func TestCreateUser_SecondUserIncrementsID(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	runQuery(t, schema, `mutation { createUser(name: "Dave") { id } }`)
	result := runQuery(t, schema, `mutation { createUser(name: "Eve") { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	user := data["createUser"].(map[string]interface{})
	if user["id"] != 4 {
		t.Errorf("expected id=4 for second create, got %v", user["id"])
	}
}

func TestCreateUser_ListShowsAllUsers(t *testing.T) {
	resetState()
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	runQuery(t, schema, `mutation { createUser(name: "Dave") { id } }`)
	runQuery(t, schema, `mutation { createUser(name: "Eve") { id } }`)

	result := runQuery(t, schema, `{ users { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list := data["users"].([]interface{})
	if len(list) != 4 {
		t.Errorf("expected 4 users after 2 creates, got %d", len(list))
	}
}
