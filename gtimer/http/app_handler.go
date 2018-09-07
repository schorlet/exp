package http

import (
	"expvar"
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/schorlet/exp/gtimer"
)

// NewAppHandler exposes services through a HTTP handler.
func NewAppHandler(todos gtimer.TodoService) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", statsHandler("app", handleIndex("Hello %s\n")))

	handler := TodoHandler(todos)
	handler = statsHandler("api/todos", handler)
	mux.Handle("/api/todos/", http.StripPrefix("/api/todos/", handler))

	mux.Handle("/debug/vars", expvar.Handler())

	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/favicon.png", http.NotFound)
	mux.HandleFunc("/opensearch.xml", http.NotFound)

	return mux
}

func handleIndex(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		if rand.Int()%2 == 0 {
			panic("random panic")
		}
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
