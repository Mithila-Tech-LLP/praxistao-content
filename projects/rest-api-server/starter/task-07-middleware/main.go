package main

import "net/http"

// Logs records entries in the format "METHOD /path STATUS"
var Logs []string

// LoggingMiddleware wraps next, records "METHOD /path STATUSCODE" in Logs after the call.
// Hint: use a custom ResponseWriter wrapper to capture the status code.
func LoggingMiddleware(next http.Handler) http.Handler {
	// TODO: implement
	return next
}

// AuthMiddleware checks for header "X-API-Key: secret".
// Returns 401 JSON error if missing or wrong. Calls next if correct.
func AuthMiddleware(next http.Handler) http.Handler {
	// TODO: implement
	return next
}

func main() {}
