package todos

import (
	"encoding/json"
	"fmt"
	"io"
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
func (t toDoList) handleAllToDos(w http.ResponseWriter, r *http.Request) {
	switch r.Method { // look for method
	case http.MethodGet:
		t.retrieveAllToDos(w, r)
	case http.MethodPost:
		t.createToDo(w, r)
	}
}

// handles requests to access /todos/# resources
// GET, PUT, DELETE
func (t toDoList) HandleSpecificTodo(w http.ResponseWriter, r *http.Request) {
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
func (t toDoList) retrieveAllToDos(w http.ResponseWriter, r *http.Request) { // also on slide 22
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
func (t toDoList) retrieveToDo(w http.ResponseWriter, r *http.Request) {
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

// creates new ToDo with ID that doesn't conflict
// server selects new unique ID and associates ToDo with that ID
func (t toDoList) createToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Orign", "*")

	desc, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		slog.Error("createToDo: error reading ToDo description", "error", err)
		http.Error(w, `"invalid ToDo format"`, http.StatusBadRequest)
		return
	}

	// Convert request to a ToDo!
	var todo toDo
	err = json.Unmarshal(desc, &todo) // converts desc back from json and assigns to todo
	if err != nil {
		slog.Error("createToDo: error unmarshaling ToDo description", "error", err)
		http.Error(w, `"invalid ToDo format"`, http.StatusBadRequest)
		return
	}

	// select next avail ID for new todo
	availID := 0
	_, exists := t.list[availID]
	for exists {
		availID++
		_, exists = t.list[availID]
	}

	slog.Info("createToDo", "id", availID, "ToDo", desc)

	// associates task with uniquely selected ID
	todo.id = availID
	t.list[todo.id] = todo // assigns

	w.WriteHeader(http.StatusCreated)
}

func (t toDoList) createReplaceToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse request
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 3 {
		slog.Error("createReplaceToDo: invalid path", "path", path)
		http.Error(w, `"invalid ToDo ID"`, http.StatusBadRequest)
		return
	}

	desc, err := io.ReadAll(r.Body) // this is how you read the request body
	defer r.Body.Close()            // don't forget this!
	if err != nil {
		slog.Error("createReplaceToDo: error reading ToDo request", "error", err)
		http.Error(w, `"invalid ToDo format"`, http.StatusBadRequest)
		return
	}

	// convert request to a todo
	var givenToDo toDo
	err = json.Unmarshal(desc, &givenToDo)
	if err != nil {
		slog.Error("createReplaceToDo: error unmarshaling ToDo request", "error", err)
		http.Error(w, `"invalid ToDo format"`, http.StatusBadRequest)
		return
	}

	// Checks if the ID given in body matches ID of the resource
	if fmt.Sprintf("%d", givenToDo.id) != path[2] {
		slog.Error("createReplaceToDo: IDs do not match", "pathID", path[2], "toDoID", givenToDo.Id)
		http.Error(w, `"ID in the ToDo does not match the ID in the URL"`, http.StatusBadRequest)
		return
	}

	// Associaes ID to the ToDo and provides corresponding client response
	_, exists := t.list[givenToDo.id]
	t.list[givenToDo.id] = givenToDo
	if exists {
		slog.Info("createReplaceToDo: replacing ToDo", "id", givenToDo.id)
		w.WriteHeader(http.StatusNoContent)
	} else { // didn't already exist, so had to create new one
		slog.Info("createReplaceToDo: creating ToDo", "id", givenToDo.id)
		w.WriteHeader(http.StatusCreated)
	}
}
