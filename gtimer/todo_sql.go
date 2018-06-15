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
func CreateTodo(db *sqlx.DB, title string) (Todo, error) {
	query := `
			insert into TODO (ID, TITLE)
			values (?, ?)`

	id, err := randomString(12)
	if err != nil {
		return Todo{}, err
	}

	err = withTx(db, func(ext sqlx.Ext) error {
		_, ere := ext.Exec(query, id, title)
		return ere
	})
	if err != nil {
		return Todo{}, err
	}

	return GetTodo(db, id)
}

// GetTodo returns the Todo with the given ID.
func GetTodo(db *sqlx.DB, id string) (Todo, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where id = ?`

	var todo Todo
	err := db.Get(&todo, query, id)

	return todo, err
}

// GetTodos returns all Todos.
func GetTodos(db *sqlx.DB) (Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			order by CREATED desc, TITLE asc`

	var todos Todos
	err := db.Select(&todos, query)

	return todos, err
}

// GetTodosByStatus returns all Todos with the specified Status.
func GetTodosByStatus(db *sqlx.DB, status string) (Todos, error) {
	query := `
			select ID, TITLE, STATUS, CREATED, UPDATED
			from TODO
			where STATUS = ?
			order by CREATED desc, TITLE asc`

	var todos Todos
	err := db.Select(&todos, query, status)

	return todos, err
}

// UpdateTodo updates the Title and Status of the given Todo.
func UpdateTodo(db *sqlx.DB, todo Todo) error {
	query := `
			update TODO set TITLE = ?, STATUS = ?
			where ID = ?`

	return withTx(db, func(ext sqlx.Ext) error {
		r, err := ext.Exec(query, todo.Title, todo.Status, todo.ID)
		if err != nil {
			return err
		}

		count, err := r.RowsAffected()
		if err == nil && count == 0 {
			err = sql.ErrNoRows
		}

		return err
	})
}

// DeleteTodo deletes the Todo with the given ID.
func DeleteTodo(db *sqlx.DB, id string) error {
	query := `delete from TODO where ID = ?`

	return withTx(db, func(ext sqlx.Ext) error {
		r, err := ext.Exec(query, id)
		if err != nil {
			return err
		}

		count, err := r.RowsAffected()
		if err == nil && count == 0 {
			err = sql.ErrNoRows
		}

		return err
	})
}
