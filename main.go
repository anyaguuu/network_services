package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World")
}

func (t toDo) ServeHttp(w http.ResponseWriter, r http.Request) {

}
