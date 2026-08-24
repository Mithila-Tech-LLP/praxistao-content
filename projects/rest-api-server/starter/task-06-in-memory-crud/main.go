package main

import "net/http"

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db = map[int]Todo{}
var nextID = 1

// TodosHandler routes:
//
//	GET  /todos       → list all todos as JSON array
//	POST /todos       → create todo, return 201 + todo
//	GET  /todos/{id}  → return todo or 404
//	PUT  /todos/{id}  → update todo or 404
//	DELETE /todos/{id}→ delete todo (204) or 404
func TodosHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func main() {
	http.HandleFunc("/todos", TodosHandler)
	http.HandleFunc("/todos/", TodosHandler)
	http.ListenAndServe(":8080", nil)
}
