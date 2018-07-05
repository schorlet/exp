package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/schorlet/exp/gtimer"
	"github.com/schorlet/exp/gtimer/http"
	"github.com/schorlet/exp/gtimer/server"
	"github.com/schorlet/exp/gtimer/storage/sqlite"
)

func main() {
	log.Println("gtimer starting ...")

	// database
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = ":memory:"
	}
	db := sqlx.MustConnect("sqlite3", url)
	defer db.Close()

	// storage
	var store sqlite.TodoStore
	store.MustDefine(db)

	// service
	service := server.TodoService{DB: db, Store: store}
	if err := initTodos(&service); err != nil {
		log.Panic(err)
	}

	//server
	server := http.NewServer(&service)
	if err := http.ListenAndServe("localhost:8000", server); err != nil {
		log.Panic(err)
	}
}

func initTodos(service gtimer.TodoService) error {
	_, err := service.Create(gtimer.Todo{ID: "st101", Title: "st101"})
	if err != nil {
		return fmt.Errorf("Unable to create st101: %v", err)
	}

	_, err = service.Create(gtimer.Todo{ID: "st102", Title: "st102"})
	if err != nil {
		return fmt.Errorf("Unable to create st102: %v", err)
	}

	return nil
}
