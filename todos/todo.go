package todos

import "net/http"

type toDo struct {
	id          int
	Description string
}

type toDoList struct {
	list map[int]toDo // so this is a list of maps, keys being id and mapped to a toDo
}

// creates new toDoList server
func New() http.Handler { // so returns http handler
	// create new todolist
	var toDoList = toDoList{make(map[int]toDo)} // just initialize

	// set handlers for appropriate paths
	mux := http.NewServeMux() // a mux takes in a method and figures out which path after a request is set in, then sets up GoRoutine
	mux.HandleFunc("/todos", toDoList.handleAllToDos)
	mux.HandleFunc("/todos/", toDoList.HandleSpecificTodo)

	return mux
}

// method of the toDoList struct
// Handles requests to access /todos resources
// # represents numeric id for a ToDo
// supports GET and POST requests
func (t toDoList) handleAllToDos(w http.ResponseWriter, r http.Request) {
	switch r.Method { // look for method
	case http.MethodGet:
		t.retrieveAllToDos(w, r)
	case http.MethodPost:
		t.createToDo(w, r)
	}
}

// handles requests to access /todos/# resources
// GET, PUT, DELETE
func (t ToDoList) HandleSpecificTodo(w http.ResponseWriter, r http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.retrieveTodo(w, r)
	case http.MethodPut:
		t.createReplaceToDo(w, r)
	case http.MethodDelete:
		t.deleteToDo(w, r)
	}
}
