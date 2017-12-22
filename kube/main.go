package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Printf("Commit: %s, Build Time: %s, Release: %s",
		Commit, BuildTime, Release)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set.")
	}
	router := Router()

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Print("Starting the service...")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	log.Print("The service is shutting down...")
	srv.Shutdown(context.Background())
	log.Print("Done")
}
