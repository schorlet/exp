package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	r := Router()
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/home")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for /home is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusOK)
	}

	res, err = http.Post(ts.URL+"/home", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status code for /home is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusMethodNotAllowed)
	}

	res, err = http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code for /home is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusNotFound)
	}
}

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	home(w, nil)

	resp := w.Result()
	if have, want := resp.StatusCode, http.StatusOK; have != want {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", have, want)
	}

	greeting, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if have, want := string(greeting), "Hello! Your request was processed.\n"; have != want {
		t.Errorf("The greeting is wrong. Have: %s, want: %s.", have, want)
	}
}
