package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{
	{"id": 1, "name": "Alice"},
	{"id": 2, "name": "Bob"},
}

var posts = []map[string]interface{}{
	{"id": 1, "title": "Hello World", "authorId": 1},
	{"id": 2, "title": "Go is great", "authorId": 2},
	{"id": 3, "title": "GraphQL rocks", "authorId": 1},
}

// BuildSchema adds Post type with nested author resolver.
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
