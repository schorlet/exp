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

	mux.Handle("/", handleIndex())
	mux.Handle("/about", statsHandler("about", handleAbout("Hello %s\n")))

	handler := TodoHandler(todos)
	handler = statsHandler("api/todos", handler)
	mux.Handle("/api/todos/", http.StripPrefix("/api/todos/", handler))

	mux.Handle("/debug/vars", basicAuth(expvar.Handler()))

	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/favicon.png", http.NotFound)
	mux.HandleFunc("/opensearch.xml", http.NotFound)

	return mux
}

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "./vuejs" + r.URL.Path
		http.ServeFile(w, r, name)
	}
}

func handleAbout(format string) http.HandlerFunc {
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
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || username != "basic" || password != "basic" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
