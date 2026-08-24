package main

import (
	"net/http"
)

// NewRouter returns an http.Handler with three routes registered:
//
//	GET /hello   → {"message":"hello"}   (200, application/json)
//	GET /ping    → {"status":"ok"}       (200, application/json)
//	GET /version → {"version":"1.0"}     (200, application/json)
func NewRouter() http.Handler {
	// TODO: implement
	return http.NewServeMux()
}

func main() {
	http.ListenAndServe(":8080", NewRouter())
}
