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

func TestBuildSchema_NoError(t *testing.T) {
	_, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() returned error: %v", err)
	}
}

func TestUserQuery(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ user(id: 1) { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("result.Data is not a map")
	}

	user, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatal("user field is not a map")
	}

	if user["id"] != 1 {
		t.Errorf("expected id=1, got %v", user["id"])
	}
	if user["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", user["name"])
	}
}
