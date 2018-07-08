package http

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/schorlet/exp/gtimer"
)

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

// TodoHandler handles Todos manipulation via an HTTP server.
type TodoHandler struct {
	Todos gtimer.TodoService
}

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// shiftPath returned values:
	// "todos", "/" := shiftPath(/todos/)
	// "todos", "/:id" := shiftPath(/todos/:id)
	_, tail := shiftPath(r.URL.Path)

	var next http.Handler
	var id string

	switch tail {
	case "/":
		switch r.Method {
		case "GET":
			next = h.GetMany()
		case "POST":
			next = h.Post()
		default:
			next = notAllowed("GET", "POST")
		}
	default:
		// ":id", "/" := shiftPath(/:id)
		id, _ = shiftPath(tail)
		switch r.Method {
		case "GET", "HEAD":
			next = h.Get(id)
		case "PUT":
			next = h.Put(id)
		case "DELETE":
			next = h.Delete(id)
		default:
			next = notAllowed("HEAD", "GET", "PUT", "DELETE")
		}
	}

	next.ServeHTTP(w, r)
}

// Post accepts a Todo encoded in JSON in the request body, saves it
// and returns it in the response body encoded in JSON.
func (h *TodoHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var create gtimer.Todo

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&create); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := h.Todos.Create(create)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

// GetMany selects the Todos according to the TodoFilter
// and returns them in the response body encoded in JSON.
func (h *TodoHandler) GetMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := gtimer.TodoFilter{Status: r.FormValue("status")}

		todos, err := h.Todos.Read(filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.Encode(todos)
	}
}

// Get selects the Todo by its ID and returns it in the response body encoded in JSON.
// A 404 error is returned if the Todo does not exist.
func (h *TodoHandler) Get(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := gtimer.TodoFilter{ID: id}

		todos, err := h.Todos.Read(filter)
		if err != nil {
			switch err {
			case gtimer.ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var buf bytes.Buffer
		if r.Method == "GET" {
			enc := json.NewEncoder(&buf)
			if err = enc.Encode(todos[0]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			etag := fmt.Sprintf(`"%x"`, md5.Sum(buf.Bytes()))
			w.Header().Set("Etag", etag)
			w.Header().Set("Cache-Control", "private, max-age=60")
		}

		content := bytes.NewReader(buf.Bytes())
		modtime := todos[0].Updated
		http.ServeContent(w, r, "todo.json", modtime, content)
	}
}

// Put handles the update of the Todo designated by the ID.
// A 404 error is returned if the Todo does not exist.
func (h *TodoHandler) Put(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update gtimer.Todo

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		update.ID = id
		todo, err := h.Todos.Update(update)
		if err != nil {
			switch err {
			case gtimer.ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

// Delete handles the deletion of the Todo designated by the ID.
// A 404 error is returned if the Todo does not exist.
func (h *TodoHandler) Delete(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.Todos.Delete(id)
		if err != nil {
			switch err {
			case gtimer.ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
