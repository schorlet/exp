package gtimer

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"

	"github.com/jmoiron/sqlx"
)

func randomString(length int) (string, error) {
	// https://www.commandlinefu.com/commands/view/24071/generate-random-text-based-on-length
	raw := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, raw)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(raw)[:length], nil
}

// CreateTodo creates a Todo with the given Title.
func CreateTodo(e sqlx.Ext, title string) (Todo, error) {
	query := `
			insert into TODO (ID, TITLE)
			values (?, ?)`

	id, err := randomString(12)
	if err != nil {
		return Todo{}, err
	}

	_, err = e.Exec(query, id, title)
	if err != nil {
		return Todo{}, err
	}

	return GetTodo(e, id)
}

// GetTodo returns the Todo with the given ID.
func GetTodo(q sqlx.Queryer, id string) (Todo, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where id = ?`

	var todo Todo
	err := sqlx.Get(q, &todo, query, id)

	return todo, err
}

// GetTodos returns all Todos.
func GetTodos(q sqlx.Queryer) (Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			order by CREATED desc, TITLE asc`

	var todos Todos
	err := sqlx.Select(q, &todos, query)

	return todos, err
}

// GetTodosByStatus returns all Todos with the specified Status.
func GetTodosByStatus(q sqlx.Queryer, status string) (Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where STATUS = ?
			order by CREATED desc, TITLE asc`

	var todos Todos
	err := sqlx.Select(q, &todos, query, status)

	return todos, err
}

// UpdateTodo updates the Title and Status of the given Todo.
func UpdateTodo(e sqlx.Execer, todo Todo) error {
	query := `
			update TODO set TITLE = ?, STATUS = ?
			where ID = ?`

	r, err := e.Exec(query, todo.Title, todo.Status, todo.ID)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		err = sql.ErrNoRows
	}
	return err
}

// DeleteTodo deletes the Todo with the given ID.
func DeleteTodo(e sqlx.Execer, id string) error {
	query := `delete from TODO where ID = ?`

	r, err := e.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err == nil && count == 0 {
		err = sql.ErrNoRows
	}
	return err
}
