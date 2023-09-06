package todos

type toDo struct {
	id          int
	Description string
}

type toDoList struct {
	list map[int]toDo // so this is a list of maps, keys being id and mapped to a toDo
}
