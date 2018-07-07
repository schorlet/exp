package mem

import (
	"testing"

	"github.com/schorlet/exp/gtimer/storage"
)

func todoTester(fn storage.TodoTest) func(*testing.T) {
	return func(t *testing.T) {
		store := make(TodoStore)
		fn(t, nil, store)
	}
}

func TestMem(t *testing.T) {
	storage.TodoTestSuite(t, todoTester)
}
