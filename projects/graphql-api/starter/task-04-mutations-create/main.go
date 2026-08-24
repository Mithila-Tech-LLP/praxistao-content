package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{
	{"id": 1, "name": "Alice"},
	{"id": 2, "name": "Bob"},
}
var nextID = 3

// BuildSchema adds a mutation: createUser(name: String!) → User
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
