package gtimer

import (
	"fmt"
	"time"
)

// Todo struct.
type Todo struct {
	ID      string    `json:"id"      db:"ID"`
	Title   string    `json:"title"   db:"TITLE"`
	Status  string    `json:"status"  db:"STATUS"`
	Created time.Time `json:"created" db:"CREATED"`
	Updated time.Time `json:"updated" db:"UPDATED"`
}

func (t Todo) String() string {
	return fmt.Sprintf("Todo{ID:%s, Title:%s, Status:%s, Created:%s, Updated:%s}",
		t.ID, t.Title, t.Status, t.Created, t.Updated)
}

// Todos slice.
type Todos []Todo

// TodoFilter struct.
type TodoFilter struct {
	ID     string
	Status string
}
