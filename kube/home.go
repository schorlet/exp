package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Router register necessary routes and returns an instance of a router.
func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/home", home).Methods("GET")
	r.HandleFunc("/version", version).Methods("GET")
	return r
}

// home is a HTTP handler which says Hello!
func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello!")
}
