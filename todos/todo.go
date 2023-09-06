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
