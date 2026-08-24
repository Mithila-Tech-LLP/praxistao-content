package main

import "net/http"

// WriteJSON writes status code and encodes v as JSON response.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	// TODO: implement
}

// WriteError writes status code and {"error": message} JSON body.
func WriteError(w http.ResponseWriter, status int, message string) {
	// TODO: implement
}

// SafeDivideHandler reads ?a=&b= query params.
// Returns {"result": N} on success (200).
// Returns 400 {"error":"missing parameter"} if a or b missing.
// Returns 400 {"error":"invalid number"} if not parseable as int.
// Returns 422 {"error":"division by zero"} if b == 0.
func SafeDivideHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	http.HandleFunc("/divide", SafeDivideHandler)
	http.ListenAndServe(":8080", nil)
}
