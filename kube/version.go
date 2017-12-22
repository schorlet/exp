package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// BuildTime is a time label of the moment when the binary was built.
var BuildTime string

// Commit is a last commit hash at the moment when the binary was built.
var Commit string

// Release is a semantic version of current build.
var Release string

// Version holds build time flags.
type Version struct {
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	Release   string `json:"release"`
}

// version returns the version of the application.
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
