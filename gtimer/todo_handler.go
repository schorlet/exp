package gtimer

import (
	"encoding/json"
	"net/http"
	"strings"
)

type TodoHandler struct {
	Todos TodoService
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
		case "GET":
			next = h.Get(id)
		case "PUT":
			next = h.Put(id)
		case "DELETE":
			next = h.Delete(id)
		default:
			next = notAllowed("GET", "PUT", "DELETE")
		}
	}

	next.ServeHTTP(w, r)
}

func (h *TodoHandler) GetMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := TodoFilter{Status: r.FormValue("status")}

		todos, err := h.Todos.Read(filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todos)
	}
}

func (h *TodoHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var create Todo

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

		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

func (h *TodoHandler) Get(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := TodoFilter{ID: id}

		todos, err := h.Todos.Read(filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todos[0])
	}
}

func (h *TodoHandler) Put(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update Todo

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		update.ID = id
		todo, err := h.Todos.Update(update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

func (h *TodoHandler) Delete(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.Todos.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
