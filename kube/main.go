package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("Starting the service...")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set.")
	}
	router := Router()
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
