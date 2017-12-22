package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("Commit: %s, Build Time: %s, Release: %s",
		Commit, BuildTime, Release)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set.")
	}
	router := Router()

	log.Print("Starting the service...")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
