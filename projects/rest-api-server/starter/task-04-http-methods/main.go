package main

import "net/http"

// ItemsHandler handles /items and dispatches based on HTTP method:
//
//	GET    → 200, body: []
//	POST   → 201, body: echo of request body as JSON
//	DELETE → 204, empty body
//	other  → 405 Method Not Allowed, body: {"error":"method not allowed"}
func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	http.HandleFunc("/items", ItemsHandler)
	http.ListenAndServe(":8080", nil)
}
