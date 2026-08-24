package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{
	{"id": 1, "name": "Alice"},
	{"id": 2, "name": "Bob"},
	{"id": 3, "name": "Alicia"},
}

// BuildSchema adds a "search(name: String!)" query that filters users by name substring (case-insensitive).
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
