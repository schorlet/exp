package http

import (
	"fmt"
	"net/http"

	"github.com/schorlet/exp/gtimer"
)

func NewServer(todos gtimer.TodoService) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", handleAPI(todos)))
	mux.HandleFunc("/", handleIndex("Hello %s\n"))
	return mux
}

func handleAPI(todos gtimer.TodoService) http.Handler {
	mux := http.NewServeMux()
	handler := TodoHandler{Todos: todos}
	mux.Handle("/todos/", &handler)
	return mux
}

func handleIndex(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, format, "World")
	}
}
