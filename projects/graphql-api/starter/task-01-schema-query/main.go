package main

import "github.com/graphql-go/graphql"

// BuildSchema builds and returns the GraphQL schema.
// It should have a root query with a "user" field that accepts an "id" (Int!)
// argument and returns a User type with "id" (Int) and "name" (String) fields.
// The resolver returns {"id": 1, "name": "Alice"} for any id.
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
