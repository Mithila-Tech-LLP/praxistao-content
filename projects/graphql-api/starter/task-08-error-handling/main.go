package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{
	{"id": 1, "name": "Alice", "email": "alice@example.com"},
}
var nextID = 2

// BuildSchema handles errors: user not found → nil + error; createUser with empty name → error
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
