package main

import (
	"log/slog"
	"net/http"

	"github.com/anyaguuu/network_services/todos"
)

func main() {
	server := todos.New()

	slog.Info("ToDo server listening on port 5318")
	err := http.ListenAndServe(":5318", server) // params: addr String, http http.handler
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server closed", "error", err)
	} else {
		slog.Info("Server closed", "error", err)
	}
}
