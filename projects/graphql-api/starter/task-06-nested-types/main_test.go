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

func TestPosts_WithNestedAuthor(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ posts { title author { name } } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list, ok := data["posts"].([]interface{})
	if !ok {
		t.Fatal("posts field is not a slice")
	}
	if len(list) != 3 {
		t.Errorf("expected 3 posts, got %d", len(list))
	}

	// Check that each post has an author with a name
	for _, item := range list {
		post := item.(map[string]interface{})
		author, ok := post["author"].(map[string]interface{})
		if !ok {
			t.Errorf("post %v has no author map", post["title"])
			continue
		}
		if author["name"] == "" || author["name"] == nil {
			t.Errorf("post %v has empty author name", post["title"])
		}
	}
}

func TestPosts_AliceIsAuthorOfCorrectPosts(t *testing.T) {
	schema, err := BuildSchema()
	if err != nil {
		t.Fatalf("BuildSchema() error: %v", err)
	}

	result := runQuery(t, schema, `{ posts { title author { name } } }`)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}

	data := result.Data.(map[string]interface{})
	list := data["posts"].([]interface{})

	alicePosts := 0
	for _, item := range list {
		post := item.(map[string]interface{})
		author := post["author"].(map[string]interface{})
		if author["name"] == "Alice" {
			alicePosts++
		}
	}

	// posts with authorId=1 (Alice): "Hello World" and "GraphQL rocks"
	if alicePosts != 2 {
		t.Errorf("expected 2 posts by Alice, got %d", alicePosts)
	}
}
