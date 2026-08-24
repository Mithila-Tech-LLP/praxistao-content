package main

import "net/http"

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var store []Todo

// CreateTodoHandler reads a JSON body {"title":"..."}, creates a Todo with
// auto-incremented ID (len(store)+1), appends to store, returns 201 + created Todo.
func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	http.HandleFunc("/todos", CreateTodoHandler)
	http.ListenAndServe(":8080", nil)
}
