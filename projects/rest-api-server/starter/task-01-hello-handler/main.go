package main

import (
	"net/http"
)

// HelloHandler responds to any request with {"message":"hello"}.
// Requirements:
//   - Status code: 200
//   - Header:      Content-Type: application/json
//   - Body:        {"message":"hello"}
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	http.HandleFunc("/hello", HelloHandler)
	http.ListenAndServe(":8080", nil)
}
