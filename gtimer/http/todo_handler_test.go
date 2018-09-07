package http

import (
	"bytes"
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

func withHandler(fn func(string, http.Handler)) {
	store := make(mem.TodoStore)
	service := server.TodoService{Store: store}

	service.Create(gtimer.Todo{ID: "st101", Title: "st101"})
	service.Create(gtimer.Todo{ID: "st102", Title: "st102"})

	handler := TodoHandler(&service)
	mux := http.NewServeMux()

	prefix := "/api/todos/"
	mux.Handle(prefix, http.StripPrefix(prefix, handler))

	fn(prefix[:len(prefix)-1], mux)
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

func TestTodoPost(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		create := gtimer.Todo{Title: "st103"}
		buf, _ := json.Marshal(create)

		r, _ := http.NewRequest("POST", prefix+"/", bytes.NewReader(buf))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if err := hasJSON(w.HeaderMap); err != nil {
			t.Fatalf("Unexpected content type: %v", err)
		}

		dec := json.NewDecoder(w.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&create); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
	})
}

func TestTodoPostBadRequest(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		create := struct{ Foo string }{"foo"}
		buf, _ := json.Marshal(create)

		r, _ := http.NewRequest("POST", prefix+"/", bytes.NewReader(buf))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}

func TestTodoGetMany(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("GET", prefix+"/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if err := hasJSON(w.HeaderMap); err != nil {
			t.Fatalf("Unexpected content type: %v", err)
		}

		var todos gtimer.Todos
		dec := json.NewDecoder(w.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&todos); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("Unexpected count of Todos: %d", len(todos))
		}
	})
}

func TestTodoHead(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("HEAD", prefix+"/st101", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if n, _ := io.Copy(ioutil.Discard, w.Body); n != 0 {
			t.Fatalf("Unexpected body size: %d", n)
		}
	})
}

func TestTodoHeadNotFound(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("HEAD", prefix+"/foo", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}

func TestTodoGet(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("GET", prefix+"/st101", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if err := hasJSON(w.HeaderMap); err != nil {
			t.Fatalf("Unexpected content type: %v", err)
		}

		var todo gtimer.Todo
		dec := json.NewDecoder(w.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&todo); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if todo.ID != "st101" {
			t.Fatalf("Unexpected todo: %s", todo)
		}
	})
}

func TestTodoGetNotFound(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("GET", prefix+"/foo", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}

func TestTodoPut(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		update := gtimer.Todo{Title: "st101-1", Status: "active"}
		buf, _ := json.Marshal(update)

		r, _ := http.NewRequest("PUT", prefix+"/st101", bytes.NewReader(buf))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
		if err := hasJSON(w.HeaderMap); err != nil {
			t.Fatalf("Unexpected content type: %v", err)
		}

		dec := json.NewDecoder(w.Body)
		if err := dec.Decode(&update); err != nil {
			t.Fatalf("Unable to decode body: %v", err)
		}
		if update.Title != "st101-1" {
			t.Fatalf("Unexpected title: %s", update.Title)
		}
		if update.Status != "active" {
			t.Fatalf("Unexpected status: %s", update.Status)
		}
	})
}

func TestTodoPutNotFound(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		update := gtimer.Todo{Title: "foo", Status: "foo"}
		buf, _ := json.Marshal(update)

		r, _ := http.NewRequest("PUT", prefix+"/foo", bytes.NewReader(buf))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}

func TestTodoPutBadRequest(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		update := struct{ Foo string }{"foo"}
		buf, _ := json.Marshal(update)

		r, _ := http.NewRequest("PUT", prefix+"/st101", bytes.NewReader(buf))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}

func TestTodoDelete(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("DELETE", prefix+"/st101", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}

func TestTodoDeleteNotFound(t *testing.T) {
	withHandler(func(prefix string, h http.Handler) {
		r, _ := http.NewRequest("DELETE", prefix+"/foo", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Unexpected status code: %d", w.Code)
		}
	})
}
