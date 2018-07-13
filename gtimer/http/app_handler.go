package http

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/schorlet/exp/gtimer"
)

// NewAppHandler exposes services through a HTTP handler.
func NewAppHandler(todos gtimer.TodoService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex("Hello %s\n"))

	handler := TodoHandler{Todos: todos}
	mux.Handle("/api/todos/", http.StripPrefix("/api/todos/", &handler))

	return mux
}

func handleIndex(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, format, "World")
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func notAllowed(allowed ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
