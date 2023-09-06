package todos

import (
	"encoding/json"
	"json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

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

// converts all existing ToDos into a JSON object
// sends the JSON obj back to client
func (t toDoList) retrieveAllToDos(w http.ResponseWriter, r http.Request) { // also on slide 22
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	toDos := make([]toDo, 0) // slice with 0 toDos
	// add all todos in todo list to it
	for _, todo := range t.list {
		toDos = append(toDos, todo)
	}

	jsonToDo, err := json.Marshal(toDos)
	if err != nil {
		// should never happen
		slog.Error("retrieveAllToDos: error marshaling ToDos", "ToDos", toDos, "error", err)
		http.Error(w, `"internal server error"`, http.StatusInternalServerError)
		return
	}

	slog.Info("retrieveAllToDos: success")
	w.Write(jsonToDo)
}

// retrieves a single ToDo if it exists and sends back to client in the form
// of json obj
// if not found, return 404 Not Found
func (t toDoList) retrieveToDo(w http.ResponseWriter, r http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// parse the request
	path := strings.Split(r.URL.Path, "/") // should be HOST, Method, URL path i think. so last element (2) is # id
	if len(path) != 3 || path[2] == "" {
		slog.Error("retrieve ToDo: invalid path", "path", path)
		http.Error(w, `"invalid path"`, http.StatusBadRequest)
		return
	}

	// idolate the single id to retrieve
	id, err := strconv.Atoi(path[2])
	if err != nil {
		slog.Error("retrieveToDo: error reading id", "id", path[2], "error", err)
		http.Error(w, `"invalid ToDo id"`, http.StatusBadRequest)
		return
	}

	todo, exists := t.list[id] // get the item
	if exists {
		// send ToDo to the client
		jsonToDo, err := json.Marshal(todo)
		if err != nil {
			// should never happen
			slog.Error("retrieveToDo: error marshaling ToDo", "ToDo", todo)
			http.Error(w, `"internal server error"`, http.StatusInternalServerError)
			return
		}
		w.Write(jsonToDo)

		slog.Info("retrieveToDo: found", "id", id)

		return
	}

	// Id not found
	slog.Info("retrieveToDo: not found", "id", id)
	w.WriteHeader(http.StatusNotFound)
}
