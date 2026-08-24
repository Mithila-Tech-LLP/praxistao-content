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

func TestSearch_MatchesMultiple(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ search(name: "ali") { id name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list, ok := data["search"].([]interface{})
	if !ok {
		t.Fatal("search field is not a slice")
	}
	if len(list) != 2 {
		t.Errorf("expected 2 results for 'ali', got %d", len(list))
	}
}

func TestSearch_MatchesSingle(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ search(name: "bob") { name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list := data["search"].([]interface{})
	if len(list) != 1 {
		t.Errorf("expected 1 result for 'bob', got %d", len(list))
	}
	user := list[0].(map[string]interface{})
	if user["name"] != "Bob" {
		t.Errorf("expected Bob, got %v", user["name"])
	}
}

func TestSearch_NoMatch(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ search(name: "xyz") { name } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list := data["search"].([]interface{})
	if len(list) != 0 {
		t.Errorf("expected 0 results for 'xyz', got %d", len(list))
	}
}
