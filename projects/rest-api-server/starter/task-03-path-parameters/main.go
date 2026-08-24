package main

import (
	"net/http"
)

// UserHandler is registered at /users/ and extracts the user ID from the path.
// For a request to /users/42 it responds with:
//
//	{"id":"42","name":"User 42"}
//
// Requirements:
//   - Extract the ID by stripping the /users/ prefix from r.URL.Path
//   - Return 200 with Content-Type: application/json
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/", UserHandler)
	http.ListenAndServe(":8080", mux)
}
