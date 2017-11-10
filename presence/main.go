package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	// 3-legged OAuth2 flow
	conf = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:8000/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	confState = randomState(20)
)

// randomState
func randomState(length int) string {
	raw := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, raw)
	if err != nil {
		log.Fatalf("randomState failed: %v\n", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw)[:length]
}

// index
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte("<html><body>Log in with <a href=\"/login\">GitHub</a></body></html>"))
	if err != nil {
		log.Printf("index: write response: %v\n", err)
	}
}

// login
func login(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(confState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// callback
func callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	fmt.Printf("state: %q\n", state)

	if state != confState {
		log.Printf("callback: invalid state: %q\n", state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	fmt.Printf("code: %q\n", code)

	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	fmt.Printf("token: %+v\n", token)
	if err != nil {
		log.Printf("callback: exchange failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Printf("callback: new request failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	ctx = context.Background()
	client := conf.Client(ctx, token)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("callback: get authenticated user failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()

	user := struct {
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&user)
	if err != nil {
		log.Printf("callback: decode user failed: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("logged in as %+v\n", user)
}

// main
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)
	if err := http.ListenAndServe("127.0.0.1:8000", nil); err != nil {
		log.Fatal(err)
	}
}
