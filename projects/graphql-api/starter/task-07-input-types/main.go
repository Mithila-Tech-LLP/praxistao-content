package main

import "github.com/graphql-go/graphql"

var users = []map[string]interface{}{}
var nextID = 1

// BuildSchema uses a UserInput input type for createUser mutation.
// createUser(input: UserInput!) where UserInput has name and email fields.
func BuildSchema() (graphql.Schema, error) {
	// TODO: implement
	return graphql.Schema{}, nil
}

func main() {}
