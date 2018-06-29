package gtimer

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

func NewServer(db *DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", handleAPI(db)))
	mux.HandleFunc("/", handleIndex("Hello %s\n"))
	return mux
}

func handleAPI(db *DB) http.Handler {
	todos := TodoServiceSQL{DB: db, Todos: TodoSqlite{}}
	handler := TodoHandler{Todos: &todos}

	mux := http.NewServeMux()
	mux.Handle("/todos/", &handler)
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
