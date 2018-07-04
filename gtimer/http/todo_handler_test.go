package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/schorlet/exp/gtimer"
	"github.com/schorlet/exp/gtimer/server"
	"github.com/schorlet/exp/gtimer/storage/mem"
)

func withHandler(fn func(h http.Handler)) {
	store := make(mem.TodoStore)
	service := server.TodoService{Store: store}

	service.Create(gtimer.Todo{ID: "st101", Title: "st101"})
	service.Create(gtimer.Todo{ID: "st102", Title: "st102"})

	handler := TodoHandler{Todos: &service}
	fn(&handler)
}

func hasJSON(header http.Header) error {
	format := header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(format)
	if err != nil {
		return fmt.Errorf("parsing media type: %v", err)
	}
	if mediatype != "application/json" {
		return fmt.Errorf("invalid media type: %s", mediatype)
	}
	return nil
}

func TestTodoGetMany(t *testing.T) {
	withHandler(func(h http.Handler) {
		r, _ := http.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if err := hasJSON(w.HeaderMap); err != nil {
			t.Fatalf("Unexpected content-type: %v", err)
		}

		var todos gtimer.Todos
		dec := json.NewDecoder(w.Body)
		if err := dec.Decode(&todos); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestTodoGet(t *testing.T) {
	withHandler(func(h http.Handler) {
		r, _ := http.NewRequest("GET", "/todos/st101", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if err := hasJSON(w.HeaderMap); err != nil {
			t.Fatalf("Unexpected content-type: %v", err)
		}

		var todo gtimer.Todo
		dec := json.NewDecoder(w.Body)
		if err := dec.Decode(&todo); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if todo.ID != "st101" {
			t.Fatalf("Unexpected todo: %s", todo)
		}
	})
}

func TestTodoGetNotFound(t *testing.T) {
	withHandler(func(h http.Handler) {
		r, _ := http.NewRequest("GET", "/todos/foo", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}
