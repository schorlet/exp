package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type Version struct {
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	Release   string `json:"release"`
}

// version is a HTTP handler which returns the version of the application.
func version(w http.ResponseWriter, _ *http.Request) {
	v := Version{BuildTime, Commit, Release}
	body, err := json.Marshal(v)
	if err != nil {
		log.Printf("Could not encode version: %v", err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable),
			http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
