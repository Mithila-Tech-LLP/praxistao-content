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

func TestUsersQuery_ReturnsAll(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ users { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list, ok := data["users"].([]interface{})
	if !ok {
		t.Fatal("users field is not a slice")
	}
	if len(list) != 3 {
		t.Errorf("expected 3 users, got %d", len(list))
	}
}

func TestUserQuery_ReturnsCorrectUser(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ user(id: 2) { name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	user, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatal("user field is not a map")
	}
	if user["name"] != "Bob" {
		t.Errorf("expected name=Bob, got %v", user["name"])
	}
}
