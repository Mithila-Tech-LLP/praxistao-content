package main

import "net/http"

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var items = []Item{
	{ID: 1, Name: "apple"},
	{ID: 2, Name: "banana"},
	{ID: 3, Name: "apricot"},
}

// SearchHandler filters items by ?q= substring (case-insensitive).
// No q param → return all items.
// Returns JSON array.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	http.HandleFunc("/search", SearchHandler)
	http.ListenAndServe(":8080", nil)
}
