package gtimer

import (
	"encoding/json"
	"net/http"
	"strings"
)

type TodoHandler struct {
	*DB
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
	_, tail := shiftPath(r.URL.Path)
	var next http.Handler
	var id string

	switch tail {
	case "/":
		switch r.Method {
		case "GET":
			next = h.GetTodos()
		case "POST":
			next = h.CreateTodo()
		default:
			next = notAllowed("GET", "POST")
		}
	default:
		id, _ = shiftPath(tail)
		switch r.Method {
		case "GET":
			next = h.GetTodo(id)
		case "PUT":
			next = h.UpdateTodo(id)
		case "DELETE":
			next = h.DeleteTodo(id)
		default:
			next = notAllowed("GET", "PUT", "DELETE")
		}
	}

	next.ServeHTTP(w, r)
}

func (h *TodoHandler) GetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// status := r.FormValue("status")
		todos, err := GetTodos(h.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todos)
	}
}

func (h *TodoHandler) CreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body Todo

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := CreateTodo(h.DB, body.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

func (h *TodoHandler) GetTodo(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todo, err := GetTodo(h.DB, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

func (h *TodoHandler) UpdateTodo(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body Todo
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body.ID = id
		todo, err := UpdateTodo(h.DB, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(todo)
	}
}

func (h *TodoHandler) DeleteTodo(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := DeleteTodo(h.DB, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
