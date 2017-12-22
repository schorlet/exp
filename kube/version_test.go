package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersion(t *testing.T) {
	BuildTime = "BuildTime"
	Commit = "Commit"
	Release = "Release"

	w := httptest.NewRecorder()
	version(w, nil)

	resp := w.Result()
	if have, want := resp.StatusCode, http.StatusOK; have != want {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", have, want)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	have := Version{}
	err = json.Unmarshal(body, &have)
	if err != nil {
		t.Errorf("Unable to decode body: %v", err)
	}

	want := Version{BuildTime, Commit, Release}
	if have != want {
		t.Errorf("The version is wrong. Have: %s, want: %s.", have, want)
	}
}
