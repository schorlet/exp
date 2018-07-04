package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

func TestTodoGetMany(t *testing.T) {
	withHandler(func(h http.Handler) {
		r, _ := http.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected status: %d %s", resp.StatusCode, resp.Status)
		}

		var todos gtimer.Todos
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&todos); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestTodoGetOne(t *testing.T) {
	withHandler(func(h http.Handler) {
		r, _ := http.NewRequest("GET", "/todos/st101", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		var todo gtimer.Todo
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&todo); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if todo.ID != "st101" {
			t.Fatalf("Unexpected todo: %s", todo)
		}

		var b bytes.Buffer
		n, err := io.Copy(&b, resp.Body)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if n != 0 {
			t.Fatalf("Unexpected read: %s", b.String())
		}
	})
}

func TestTodoGetNotFound(t *testing.T) {
	withHandler(func(h http.Handler) {
		r, _ := http.NewRequest("GET", "/todos/foo", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}
