package main

import (
	"log"
	"net/http"
)

func main() {
	log.Print("Starting the service...")
	router := Router()
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
