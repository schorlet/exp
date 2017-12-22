package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// Router register necessary routes and returns an instance of a router.
func Router() *mux.Router {
	isReady := &atomic.Value{}
	isReady.Store(false)
	// the application is ready after 1 second.
	go func() {
		time.Sleep(time.Second)
		isReady.Store(true)
		log.Printf("Ready to serve traffic")
	}()

	r := mux.NewRouter()
	r.HandleFunc("/home", home).Methods("GET")
	r.HandleFunc("/version", version).Methods("GET")
	r.HandleFunc("/liveness", liveness).Methods("GET")
	r.HandleFunc("/readiness", readiness(isReady)).Methods("GET")
	return r
}

// home is a HTTP handler which says Hello!
func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello!")
}
