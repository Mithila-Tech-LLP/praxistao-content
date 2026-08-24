package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{
	{"id": 1, "name": "Alice"},
	{"id": 2, "name": "Bob"},
}

// BuildSchema adds mutations: updateUser(id, name) and deleteUser(id)
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
