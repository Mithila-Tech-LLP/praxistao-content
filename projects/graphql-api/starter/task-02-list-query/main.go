package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{
	{"id": 1, "name": "Alice"},
	{"id": 2, "name": "Bob"},
	{"id": 3, "name": "Charlie"},
}

// BuildSchema adds a "users" query that returns all users as a list.
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement — include both "user(id)" and "users" queries
	return graphql.Schema{}, nil
}

func main() {}
